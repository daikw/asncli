package cmd

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/michalvavra/asncli/internal/asana"
	"github.com/michalvavra/asncli/internal/cli"
	"github.com/michalvavra/asncli/internal/config"
)

type configClientStub struct {
	workspaces []asana.WorkspaceMembership
}

func (c *configClientStub) ListWorkspaceMembershipsForUser(ctx context.Context, userGID string) (*asana.WorkspaceMembershipList, error) {
	return &asana.WorkspaceMembershipList{Data: c.workspaces}, nil
}

func TestConfigGetWorkspaceEmpty(t *testing.T) {
	var stdout, stderr bytes.Buffer
	c := &cli.Context{
		Stdout: &stdout,
		Stderr: &stderr,
		Config: &config.Config{},
	}

	cmd := ConfigGetWorkspaceCmd{}
	if err := cmd.Run(context.Background(), c); err != nil {
		t.Fatalf("ConfigGetWorkspaceCmd.Run returned unexpected error: %v", err)
	}

	got := stdout.String()
	if !strings.Contains(got, "no default workspace set") {
		t.Errorf("output = %q, want to contain 'no default workspace set'", got)
	}
}

func TestConfigGetWorkspaceSet(t *testing.T) {
	var stdout, stderr bytes.Buffer
	c := &cli.Context{
		Stdout: &stdout,
		Stderr: &stderr,
		Config: &config.Config{DefaultWorkspace: "123456", DefaultWorkspaceName: "My Workspace"},
	}

	cmd := ConfigGetWorkspaceCmd{}
	if err := cmd.Run(context.Background(), c); err != nil {
		t.Fatalf("ConfigGetWorkspaceCmd.Run returned unexpected error: %v", err)
	}

	got := stdout.String()
	if !strings.Contains(got, "123456") {
		t.Errorf("output = %q, want to contain workspace GID '123456'", got)
	}
	if !strings.Contains(got, "My Workspace") {
		t.Errorf("output = %q, want to contain workspace name 'My Workspace'", got)
	}
}

func TestConfigShowJSON(t *testing.T) {
	var stdout, stderr bytes.Buffer
	c := &cli.Context{
		Stdout: &stdout,
		Stderr: &stderr,
		JSON:   true,
		Config: &config.Config{DefaultWorkspace: "789012", DefaultWorkspaceName: "Test Workspace"},
	}

	cmd := ConfigShowCmd{}
	if err := cmd.Run(context.Background(), c); err != nil {
		t.Fatalf("ConfigShowCmd.Run returned unexpected error: %v", err)
	}

	got := stdout.String()
	if !strings.Contains(got, "789012") {
		t.Errorf("output = %q, want to contain workspace GID '789012'", got)
	}
	if !strings.Contains(got, "Test Workspace") {
		t.Errorf("output = %q, want to contain workspace name 'Test Workspace'", got)
	}
	if !strings.Contains(got, "config_path") {
		t.Errorf("output = %q, want to contain 'config_path'", got)
	}
}

func TestConfigShowTable(t *testing.T) {
	var stdout, stderr bytes.Buffer
	c := &cli.Context{
		Stdout: &stdout,
		Stderr: &stderr,
		Config: &config.Config{DefaultWorkspace: "123", DefaultWorkspaceName: "My WS"},
	}

	cmd := ConfigShowCmd{}
	if err := cmd.Run(context.Background(), c); err != nil {
		t.Fatalf("ConfigShowCmd.Run returned unexpected error: %v", err)
	}

	got := stdout.String()
	if !strings.Contains(got, "SETTING") {
		t.Errorf("output = %q, want to contain table header 'SETTING'", got)
	}
	if !strings.Contains(got, "default_workspace") {
		t.Errorf("output = %q, want to contain 'default_workspace'", got)
	}
}

func TestConfigGetWorkspaceJSON(t *testing.T) {
	var stdout, stderr bytes.Buffer
	c := &cli.Context{
		Stdout: &stdout,
		Stderr: &stderr,
		JSON:   true,
		Config: &config.Config{DefaultWorkspace: "ws-id", DefaultWorkspaceName: "WS Name"},
	}

	cmd := ConfigGetWorkspaceCmd{}
	if err := cmd.Run(context.Background(), c); err != nil {
		t.Fatalf("ConfigGetWorkspaceCmd.Run returned unexpected error: %v", err)
	}

	got := stdout.String()
	if !strings.Contains(got, "ws-id") {
		t.Errorf("output = %q, want to contain 'ws-id'", got)
	}
}

func TestConfigGetWorkspaceEmptyJSON(t *testing.T) {
	var stdout, stderr bytes.Buffer
	c := &cli.Context{
		Stdout: &stdout,
		Stderr: &stderr,
		JSON:   true,
		Config: &config.Config{},
	}

	cmd := ConfigGetWorkspaceCmd{}
	if err := cmd.Run(context.Background(), c); err != nil {
		t.Fatalf("ConfigGetWorkspaceCmd.Run returned unexpected error: %v", err)
	}

	got := stdout.String()
	if !strings.Contains(got, "null") {
		t.Errorf("output = %q, want to contain 'null' for empty workspace", got)
	}
}

func TestConfigGetWorkspaceOnlyGID(t *testing.T) {
	var stdout, stderr bytes.Buffer
	c := &cli.Context{
		Stdout: &stdout,
		Stderr: &stderr,
		Config: &config.Config{DefaultWorkspace: "only-gid"},
	}

	cmd := ConfigGetWorkspaceCmd{}
	if err := cmd.Run(context.Background(), c); err != nil {
		t.Fatalf("ConfigGetWorkspaceCmd.Run returned unexpected error: %v", err)
	}

	got := stdout.String()
	if !strings.Contains(got, "only-gid") {
		t.Errorf("output = %q, want to contain workspace GID 'only-gid'", got)
	}
}

func TestConfigSetWorkspaceInvalidClient(t *testing.T) {
	var stdout, stderr bytes.Buffer
	c := &cli.Context{
		Stdout: &stdout,
		Stderr: &stderr,
		Client: struct{}{},
	}

	cmd := ConfigSetWorkspaceCmd{}
	err := cmd.Run(context.Background(), c)
	if err == nil {
		t.Fatal("ConfigSetWorkspaceCmd.Run should return error for unsupported client type")
	}
	if !strings.Contains(err.Error(), "client does not support") {
		t.Errorf("error = %q, want to contain 'client does not support'", err.Error())
	}
}

func TestWorkspaceListerAdapter(t *testing.T) {
	memberships := config.WorkspaceMembershipList{
		Data: []config.WorkspaceMembership{
			{GID: "ws-1", Name: "Workspace 1"},
		},
	}
	adapter := &workspaceListerAdapter{memberships: memberships}

	got, err := adapter.ListWorkspaceMembershipsForUser(context.Background(), "me")
	if err != nil {
		t.Fatalf("ListWorkspaceMembershipsForUser returned unexpected error: %v", err)
	}
	if len(got.Data) != 1 {
		t.Errorf("workspaces count = %d, want 1", len(got.Data))
	}
	if got.Data[0].GID != "ws-1" {
		t.Errorf("first workspace GID = %q, want %q", got.Data[0].GID, "ws-1")
	}
}

func TestConfigSetWorkspaceListerError(t *testing.T) {
	var stdout, stderr bytes.Buffer
	c := &cli.Context{
		Stdout: &stdout,
		Stderr: &stderr,
		Client: &configClientStubError{err: context.DeadlineExceeded},
	}

	cmd := ConfigSetWorkspaceCmd{}
	err := cmd.Run(context.Background(), c)
	if err == nil {
		t.Fatal("ConfigSetWorkspaceCmd.Run should return error when listing workspaces fails")
	}
	if !strings.Contains(err.Error(), "failed to list workspaces") {
		t.Errorf("error = %q, want to contain 'failed to list workspaces'", err.Error())
	}
}

type configClientStubError struct {
	err error
}

func (c *configClientStubError) ListWorkspaceMembershipsForUser(ctx context.Context, userGID string) (*asana.WorkspaceMembershipList, error) {
	return nil, c.err
}

func TestConfigGetWorkspaceJSONEmpty(t *testing.T) {
	var stdout, stderr bytes.Buffer
	c := &cli.Context{
		Stdout: &stdout,
		Stderr: &stderr,
		JSON:   true,
		Config: &config.Config{},
	}

	cmd := ConfigGetWorkspaceCmd{}
	if err := cmd.Run(context.Background(), c); err != nil {
		t.Fatalf("ConfigGetWorkspaceCmd.Run returned unexpected error: %v", err)
	}

	got := stdout.String()
	// Should have null values for empty workspace
	if !strings.Contains(got, "null") {
		t.Errorf("output = %q, want to contain 'null' for empty config", got)
	}
	if !strings.Contains(got, "workspace_id") {
		t.Errorf("output = %q, want to contain 'workspace_id'", got)
	}
}

func TestConfigGetWorkspaceJSONSet(t *testing.T) {
	var stdout, stderr bytes.Buffer
	c := &cli.Context{
		Stdout: &stdout,
		Stderr: &stderr,
		JSON:   true,
		Config: &config.Config{DefaultWorkspace: "json-ws-id", DefaultWorkspaceName: "JSON WS"},
	}

	cmd := ConfigGetWorkspaceCmd{}
	if err := cmd.Run(context.Background(), c); err != nil {
		t.Fatalf("ConfigGetWorkspaceCmd.Run returned unexpected error: %v", err)
	}

	got := stdout.String()
	if !strings.Contains(got, "json-ws-id") {
		t.Errorf("output = %q, want to contain 'json-ws-id'", got)
	}
	if !strings.Contains(got, "JSON WS") {
		t.Errorf("output = %q, want to contain 'JSON WS'", got)
	}
}

func TestConfigShowTableEmpty(t *testing.T) {
	var stdout, stderr bytes.Buffer
	c := &cli.Context{
		Stdout: &stdout,
		Stderr: &stderr,
		Config: &config.Config{},
	}

	cmd := ConfigShowCmd{}
	if err := cmd.Run(context.Background(), c); err != nil {
		t.Fatalf("ConfigShowCmd.Run returned unexpected error: %v", err)
	}

	got := stdout.String()
	if !strings.Contains(got, "Config file:") {
		t.Errorf("output = %q, want to contain 'Config file:'", got)
	}
	if !strings.Contains(got, "default_workspace") {
		t.Errorf("output = %q, want to contain 'default_workspace'", got)
	}
}

func TestConfigShowJSONEmpty(t *testing.T) {
	var stdout, stderr bytes.Buffer
	c := &cli.Context{
		Stdout: &stdout,
		Stderr: &stderr,
		JSON:   true,
		Config: &config.Config{},
	}

	cmd := ConfigShowCmd{}
	if err := cmd.Run(context.Background(), c); err != nil {
		t.Fatalf("ConfigShowCmd.Run returned unexpected error: %v", err)
	}

	got := stdout.String()
	if !strings.Contains(got, "config_path") {
		t.Errorf("output = %q, want to contain 'config_path'", got)
	}
	if !strings.Contains(got, "default_workspace") {
		t.Errorf("output = %q, want to contain 'default_workspace'", got)
	}
}
