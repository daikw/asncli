package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/michalvavra/asncli/internal/asana"
	"github.com/michalvavra/asncli/internal/cli"
)

type workspacesClient struct {
	list       *asana.WorkspaceMembershipList
	gotUserGID string
}

func (f *workspacesClient) ListWorkspaceMembershipsForUser(ctx context.Context, userGID string) (*asana.WorkspaceMembershipList, error) {
	f.gotUserGID = userGID
	return f.list, nil
}

func TestWorkspacesListHumanReadable(t *testing.T) {
	buf := &bytes.Buffer{}
	client := &workspacesClient{list: &asana.WorkspaceMembershipList{Data: []asana.WorkspaceMembership{
		{GID: "wm-1", Workspace: asana.Workspace{GID: "ws-1", Name: "My Company"}, IsActive: true},
		{GID: "wm-2", Workspace: asana.Workspace{GID: "ws-2", Name: "Side Project"}, IsActive: false},
	}}}
	ctx := &cli.Context{
		Stdout: buf,
		Stderr: &bytes.Buffer{},
		Client: client,
	}

	cmd := WorkspacesListCmd{}
	if err := cmd.Run(context.Background(), ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "My Company") || !strings.Contains(out, "Side Project") {
		t.Errorf("output = %q, want to contain workspace names", out)
	}
	if !strings.Contains(out, "ws-1") || !strings.Contains(out, "ws-2") {
		t.Errorf("output = %q, want to contain workspace GIDs", out)
	}
	if client.gotUserGID != "me" {
		t.Errorf("user GID = %q, want %q", client.gotUserGID, "me")
	}
}

func TestWorkspacesListJSON(t *testing.T) {
	buf := &bytes.Buffer{}
	client := &workspacesClient{list: &asana.WorkspaceMembershipList{Data: []asana.WorkspaceMembership{
		{GID: "wm-1", Workspace: asana.Workspace{GID: "ws-1", Name: "Test"}},
	}}}
	ctx := &cli.Context{
		Stdout: buf,
		Stderr: &bytes.Buffer{},
		JSON:   true,
		Client: client,
	}

	cmd := WorkspacesListCmd{}
	if err := cmd.Run(context.Background(), ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var env cli.Envelope
	if err := json.Unmarshal(buf.Bytes(), &env); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}
	if env.Data == nil {
		t.Fatal("envelope data is nil")
	}
}

func TestWorkspacesListInvalidClient(t *testing.T) {
	ctx := &cli.Context{
		Stdout: &bytes.Buffer{},
		Stderr: &bytes.Buffer{},
		Client: struct{}{},
	}

	cmd := WorkspacesListCmd{}
	err := cmd.Run(context.Background(), ctx)
	if err == nil {
		t.Fatal("should return error for unsupported client type")
	}
	if !strings.Contains(err.Error(), "client does not support") {
		t.Errorf("error = %q, want to contain 'client does not support'", err.Error())
	}
}
