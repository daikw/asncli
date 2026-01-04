package config

import (
	"bytes"
	"context"
	"os"
	"strings"
	"testing"
)

type stubLister struct {
	memberships WorkspaceMembershipList
	err         error
}

func (s *stubLister) ListWorkspaceMembershipsForUser(ctx context.Context, userGID string) (WorkspaceMembershipList, error) {
	return s.memberships, s.err
}

func TestPromptWorkspaceSingleWorkspace(t *testing.T) {
	lister := &stubLister{
		memberships: WorkspaceMembershipList{
			Data: []WorkspaceMembership{
				{GID: "ws-1", Name: "Only Workspace", IsActive: true},
			},
		},
	}

	var stdout, stderr bytes.Buffer
	gid, err := PromptWorkspace(context.Background(), lister, &stdout, &stderr)
	if err != nil {
		t.Fatalf("PromptWorkspace returned unexpected error: %v", err)
	}
	if gid != "ws-1" {
		t.Errorf("PromptWorkspace returned gid = %q, want %q", gid, "ws-1")
	}
	if !strings.Contains(stdout.String(), "Using your only workspace") {
		t.Errorf("output = %q, want to contain 'Using your only workspace'", stdout.String())
	}
}

func TestPromptWorkspaceNoWorkspaces(t *testing.T) {
	lister := &stubLister{
		memberships: WorkspaceMembershipList{Data: nil},
	}

	var stdout, stderr bytes.Buffer
	_, err := PromptWorkspace(context.Background(), lister, &stdout, &stderr)
	if err == nil {
		t.Fatal("PromptWorkspace should return error when no workspaces found")
	}
	if !strings.Contains(err.Error(), "no workspaces found") {
		t.Errorf("error = %q, want to contain 'no workspaces found'", err.Error())
	}
}

func TestPromptWorkspaceListerError(t *testing.T) {
	lister := &stubLister{
		err: context.DeadlineExceeded,
	}

	var stdout, stderr bytes.Buffer
	_, err := PromptWorkspace(context.Background(), lister, &stdout, &stderr)
	if err == nil {
		t.Fatal("PromptWorkspace should return error when lister fails")
	}
	if !strings.Contains(err.Error(), "failed to list workspaces") {
		t.Errorf("error = %q, want to contain 'failed to list workspaces'", err.Error())
	}
}

func TestPromptWorkspaceMultipleWorkspaces(t *testing.T) {
	lister := &stubLister{
		memberships: WorkspaceMembershipList{
			Data: []WorkspaceMembership{
				{GID: "ws-1", Name: "Workspace 1", IsActive: true},
				{GID: "ws-2", Name: "Workspace 2", IsActive: false},
			},
		},
	}

	withStdin(t, "2\n", func() {
		var stdout, stderr bytes.Buffer
		gid, err := PromptWorkspace(context.Background(), lister, &stdout, &stderr)
		if err != nil {
			t.Fatalf("PromptWorkspace returned unexpected error: %v", err)
		}
		if gid != "ws-2" {
			t.Errorf("PromptWorkspace returned gid = %q, want %q", gid, "ws-2")
		}
		if !strings.Contains(stdout.String(), "Select default workspace") {
			t.Errorf("output = %q, want to contain 'Select default workspace'", stdout.String())
		}
		if !strings.Contains(stdout.String(), "(inactive)") {
			t.Errorf("output = %q, want to contain '(inactive)' for inactive workspace", stdout.String())
		}
	})
}

func TestPromptWorkspaceInvalidChoice(t *testing.T) {
	lister := &stubLister{
		memberships: WorkspaceMembershipList{
			Data: []WorkspaceMembership{
				{GID: "ws-1", Name: "Workspace 1", IsActive: true},
				{GID: "ws-2", Name: "Workspace 2", IsActive: true},
			},
		},
	}

	withStdin(t, "abc\n", func() {
		var stdout, stderr bytes.Buffer
		_, err := PromptWorkspace(context.Background(), lister, &stdout, &stderr)
		if err == nil {
			t.Fatal("PromptWorkspace should return error for non-numeric input")
		}
		if !strings.Contains(err.Error(), "invalid choice") {
			t.Errorf("error = %q, want to contain 'invalid choice'", err.Error())
		}
	})
}

func TestPromptWorkspaceOutOfRange(t *testing.T) {
	lister := &stubLister{
		memberships: WorkspaceMembershipList{
			Data: []WorkspaceMembership{
				{GID: "ws-1", Name: "Workspace 1", IsActive: true},
				{GID: "ws-2", Name: "Workspace 2", IsActive: true},
			},
		},
	}

	withStdin(t, "5\n", func() {
		var stdout, stderr bytes.Buffer
		_, err := PromptWorkspace(context.Background(), lister, &stdout, &stderr)
		if err == nil {
			t.Fatal("PromptWorkspace should return error for choice > workspace count")
		}
		if !strings.Contains(err.Error(), "choice out of range") {
			t.Errorf("error = %q, want to contain 'choice out of range'", err.Error())
		}
	})
}

func TestPromptWorkspaceZeroChoice(t *testing.T) {
	lister := &stubLister{
		memberships: WorkspaceMembershipList{
			Data: []WorkspaceMembership{
				{GID: "ws-1", Name: "Workspace 1", IsActive: true},
				{GID: "ws-2", Name: "Workspace 2", IsActive: true},
			},
		},
	}

	withStdin(t, "0\n", func() {
		var stdout, stderr bytes.Buffer
		_, err := PromptWorkspace(context.Background(), lister, &stdout, &stderr)
		if err == nil {
			t.Fatal("PromptWorkspace should return error for zero choice")
		}
		if !strings.Contains(err.Error(), "choice out of range") {
			t.Errorf("error = %q, want to contain 'choice out of range'", err.Error())
		}
	})
}

func withStdin(t *testing.T, input string, fn func()) {
	t.Helper()
	orig := os.Stdin
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe for stdin: %v", err)
	}
	if _, err := w.Write([]byte(input)); err != nil {
		_ = r.Close()
		_ = w.Close()
		t.Fatalf("failed to write to stdin pipe: %v", err)
	}
	if err := w.Close(); err != nil {
		_ = r.Close()
		t.Fatalf("failed to close stdin pipe writer: %v", err)
	}
	os.Stdin = r
	t.Cleanup(func() {
		os.Stdin = orig
		_ = r.Close()
	})
	fn()
}
