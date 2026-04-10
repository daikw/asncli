package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/michalvavra/asncli/internal/asana"
	"github.com/michalvavra/asncli/internal/cli"
	"github.com/michalvavra/asncli/internal/config"
)

type usersClient struct {
	list          *asana.UserList
	user          *asana.UserDetail
	gotWorkspace  string
	gotLimit      int
	gotUserGID    string
}

func (f *usersClient) ListUsersInWorkspace(ctx context.Context, workspaceGID string, limit int) (*asana.UserList, error) {
	f.gotWorkspace = workspaceGID
	f.gotLimit = limit
	return f.list, nil
}

func (f *usersClient) GetUser(ctx context.Context, userGID string) (*asana.UserDetail, error) {
	f.gotUserGID = userGID
	return f.user, nil
}

func TestUsersListHumanReadable(t *testing.T) {
	buf := &bytes.Buffer{}
	client := &usersClient{list: &asana.UserList{Data: []asana.UserDetail{
		{GID: "u1", Name: "Alice", Email: "alice@example.com"},
		{GID: "u2", Name: "Bob", Email: "bob@example.com"},
	}}}
	ctx := &cli.Context{
		Stdout: buf,
		Stderr: &bytes.Buffer{},
		Client: client,
		Config: &config.Config{DefaultWorkspace: "ws-1"},
	}

	cmd := UsersListCmd{}
	if err := cmd.Run(context.Background(), ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "Alice") || !strings.Contains(out, "Bob") {
		t.Errorf("output = %q, want to contain user names", out)
	}
	if !strings.Contains(out, "alice@example.com") {
		t.Errorf("output = %q, want to contain email", out)
	}
	if client.gotWorkspace != "ws-1" {
		t.Errorf("workspace = %q, want %q", client.gotWorkspace, "ws-1")
	}
}

func TestUsersListJSON(t *testing.T) {
	buf := &bytes.Buffer{}
	client := &usersClient{list: &asana.UserList{Data: []asana.UserDetail{{GID: "u1", Name: "Alice"}}}}
	ctx := &cli.Context{
		Stdout: buf,
		Stderr: &bytes.Buffer{},
		JSON:   true,
		Client: client,
		Config: &config.Config{DefaultWorkspace: "ws-1"},
	}

	cmd := UsersListCmd{}
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

func TestUsersListNoWorkspace(t *testing.T) {
	ctx := &cli.Context{
		Stdout: &bytes.Buffer{},
		Stderr: &bytes.Buffer{},
		Client: &usersClient{list: &asana.UserList{}},
		Config: &config.Config{},
	}

	cmd := UsersListCmd{}
	err := cmd.Run(context.Background(), ctx)
	if err == nil {
		t.Fatal("should return error when no workspace configured")
	}
	if !strings.Contains(err.Error(), "workspace is required") {
		t.Errorf("error = %q, want to contain 'workspace is required'", err.Error())
	}
}

func TestUsersListInvalidClient(t *testing.T) {
	ctx := &cli.Context{
		Stdout: &bytes.Buffer{},
		Stderr: &bytes.Buffer{},
		Client: struct{}{},
		Config: &config.Config{DefaultWorkspace: "ws-1"},
	}

	cmd := UsersListCmd{}
	err := cmd.Run(context.Background(), ctx)
	if err == nil {
		t.Fatal("should return error for unsupported client type")
	}
}

func TestUsersGetHumanReadable(t *testing.T) {
	buf := &bytes.Buffer{}
	client := &usersClient{user: &asana.UserDetail{GID: "u1", Name: "Alice", Email: "alice@example.com"}}
	ctx := &cli.Context{
		Stdout: buf,
		Stderr: &bytes.Buffer{},
		Client: client,
	}

	cmd := UsersGetCmd{GID: "u1"}
	if err := cmd.Run(context.Background(), ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "Alice") || !strings.Contains(out, "alice@example.com") {
		t.Errorf("output = %q, want to contain user details", out)
	}
}

func TestUsersGetJSON(t *testing.T) {
	buf := &bytes.Buffer{}
	client := &usersClient{user: &asana.UserDetail{GID: "u1", Name: "Alice"}}
	ctx := &cli.Context{
		Stdout: buf,
		Stderr: &bytes.Buffer{},
		JSON:   true,
		Client: client,
	}

	cmd := UsersGetCmd{GID: "u1"}
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

func TestUsersGetInvalidClient(t *testing.T) {
	ctx := &cli.Context{
		Stdout: &bytes.Buffer{},
		Stderr: &bytes.Buffer{},
		Client: struct{}{},
	}

	cmd := UsersGetCmd{GID: "u1"}
	err := cmd.Run(context.Background(), ctx)
	if err == nil {
		t.Fatal("should return error for unsupported client type")
	}
}
