package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/michalvavra/asncli/internal/asana"
	"github.com/michalvavra/asncli/internal/auth"
	"github.com/michalvavra/asncli/internal/cli"
)

type storeRecorder struct {
	setService string
	setUser    string
	setToken   string
	setCalled  bool

	deleteService string
	deleteUser    string
	deleteCalled  bool
}

func (s *storeRecorder) Get(service, user string) (string, error) { return "", auth.ErrNoToken }

func (s *storeRecorder) Set(service, user, token string) error {
	s.setCalled = true
	s.setService = service
	s.setUser = user
	s.setToken = token
	return nil
}

func (s *storeRecorder) Delete(service, user string) error {
	s.deleteCalled = true
	s.deleteService = service
	s.deleteUser = user
	return nil
}

type authClientStub struct {
	user        *asana.User
	memberships *asana.WorkspaceMembershipList
}

func (a authClientStub) GetMe(ctx context.Context) (*asana.User, error) { return a.user, nil }
func (a authClientStub) ListWorkspaceMembershipsForUser(ctx context.Context, userGID string) (*asana.WorkspaceMembershipList, error) {
	return a.memberships, nil
}

func TestAuthLoginStoresToken(t *testing.T) {
	store := &storeRecorder{}
	buf := &bytes.Buffer{}
	ctx := &cli.Context{Stdout: buf, Stderr: &bytes.Buffer{}, Store: store}

	cmd := AuthLoginCmd{Token: "abc"}
	if err := cmd.Run(context.Background(), ctx); err != nil {
		t.Fatalf("AuthLoginCmd.Run returned unexpected error: %v", err)
	}
	if !store.setCalled {
		t.Error("store Set was not called")
	}
	if got, want := store.setToken, "abc"; got != want {
		t.Errorf("stored token = %q, want %q", got, want)
	}
}

func TestAuthLogoutDeletesToken(t *testing.T) {
	store := &storeRecorder{}
	ctx := &cli.Context{Stdout: &bytes.Buffer{}, Stderr: &bytes.Buffer{}, Store: store}

	cmd := AuthLogoutCmd{}
	if err := cmd.Run(context.Background(), ctx); err != nil {
		t.Fatalf("AuthLogoutCmd.Run returned unexpected error: %v", err)
	}
	if !store.deleteCalled {
		t.Error("store Delete was not called")
	}
	if got, want := store.deleteService, auth.DefaultService; got != want {
		t.Errorf("delete service = %q, want %q", got, want)
	}
	if got, want := store.deleteUser, auth.DefaultUser; got != want {
		t.Errorf("delete user = %q, want %q", got, want)
	}
}

func TestAuthStatusJSON(t *testing.T) {
	buf := &bytes.Buffer{}
	ctx := &cli.Context{
		Stdout: buf,
		Stderr: &bytes.Buffer{},
		JSON:   true,
		Client: authClientStub{
			user: &asana.User{GID: "1", Name: "Jane"},
			memberships: &asana.WorkspaceMembershipList{
				Data: []asana.WorkspaceMembership{
					{GID: "wm1", Workspace: asana.Workspace{GID: "w1", Name: "Main"}},
				},
			},
		},
	}

	cmd := AuthStatusCmd{}
	if err := cmd.Run(context.Background(), ctx); err != nil {
		t.Fatalf("AuthStatusCmd.Run returned unexpected error: %v", err)
	}

	var env cli.Envelope
	if err := json.Unmarshal(buf.Bytes(), &env); err != nil {
		t.Fatalf("failed to unmarshal JSON output: %v", err)
	}
	if env.Data == nil {
		t.Fatal("envelope data is nil, want auth status")
	}
}

func TestAuthLoginJSON(t *testing.T) {
	store := &storeRecorder{}
	buf := &bytes.Buffer{}
	ctx := &cli.Context{Stdout: buf, Stderr: &bytes.Buffer{}, Store: store, JSON: true}

	cmd := AuthLoginCmd{Token: "abc"}
	if err := cmd.Run(context.Background(), ctx); err != nil {
		t.Fatalf("AuthLoginCmd.Run returned unexpected error: %v", err)
	}

	var env cli.Envelope
	if err := json.Unmarshal(buf.Bytes(), &env); err != nil {
		t.Fatalf("failed to unmarshal JSON output: %v", err)
	}
	data, ok := env.Data.(map[string]any)
	if !ok {
		t.Fatalf("envelope data type = %T, want map[string]any", env.Data)
	}
	if got, want := data["status"], "stored"; got != want {
		t.Errorf("status = %v, want %v", got, want)
	}
}

func TestAuthLogoutJSON(t *testing.T) {
	store := &storeRecorder{}
	buf := &bytes.Buffer{}
	ctx := &cli.Context{Stdout: buf, Stderr: &bytes.Buffer{}, Store: store, JSON: true}

	cmd := AuthLogoutCmd{}
	if err := cmd.Run(context.Background(), ctx); err != nil {
		t.Fatalf("AuthLogoutCmd.Run returned unexpected error: %v", err)
	}

	var env cli.Envelope
	if err := json.Unmarshal(buf.Bytes(), &env); err != nil {
		t.Fatalf("failed to unmarshal JSON output: %v", err)
	}
	data, ok := env.Data.(map[string]any)
	if !ok {
		t.Fatalf("envelope data type = %T, want map[string]any", env.Data)
	}
	if got, want := data["status"], "removed"; got != want {
		t.Errorf("status = %v, want %v", got, want)
	}
}

func TestAuthStatusHumanReadable(t *testing.T) {
	buf := &bytes.Buffer{}
	ctx := &cli.Context{
		Stdout: buf,
		Stderr: &bytes.Buffer{},
		Client: authClientStub{
			user: &asana.User{GID: "1", Name: "Jane"},
			memberships: &asana.WorkspaceMembershipList{
				Data: []asana.WorkspaceMembership{
					{GID: "wm1", Workspace: asana.Workspace{GID: "w1", Name: "Main"}, IsActive: true},
				},
			},
		},
	}

	cmd := AuthStatusCmd{}
	if err := cmd.Run(context.Background(), ctx); err != nil {
		t.Fatalf("AuthStatusCmd.Run returned unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "authenticated as Jane") || !strings.Contains(out, "WORKSPACE_GID") || !strings.Contains(out, "w1") {
		t.Errorf("output missing auth status details\ngot: %q\nwant: contains 'authenticated as Jane', 'WORKSPACE_GID', and 'w1'", out)
	}
}

type storeRecorderWithToken struct {
	token string
}

func (s *storeRecorderWithToken) Get(service, user string) (string, error) {
	if s.token != "" {
		return s.token, nil
	}
	return "", auth.ErrNoToken
}

func (s *storeRecorderWithToken) Set(service, user, token string) error {
	s.token = token
	return nil
}

func (s *storeRecorderWithToken) Delete(service, user string) error {
	s.token = ""
	return nil
}

func TestAuthLoginStoresTokenHuman(t *testing.T) {
	store := &storeRecorder{}
	buf := &bytes.Buffer{}
	ctx := &cli.Context{Stdout: buf, Stderr: &bytes.Buffer{}, Store: store}

	cmd := AuthLoginCmd{Token: "test-token"}
	if err := cmd.Run(context.Background(), ctx); err != nil {
		t.Fatalf("AuthLoginCmd.Run returned unexpected error: %v", err)
	}

	if !store.setCalled {
		t.Error("store Set was not called")
	}
	if store.setToken != "test-token" {
		t.Errorf("stored token = %q, want %q", store.setToken, "test-token")
	}

	out := buf.String()
	if !strings.Contains(out, "stored") {
		t.Errorf("output = %q, want to contain 'stored'", out)
	}
}

func TestAuthLogoutHuman(t *testing.T) {
	store := &storeRecorder{}
	buf := &bytes.Buffer{}
	ctx := &cli.Context{Stdout: buf, Stderr: &bytes.Buffer{}, Store: store}

	cmd := AuthLogoutCmd{}
	if err := cmd.Run(context.Background(), ctx); err != nil {
		t.Fatalf("AuthLogoutCmd.Run returned unexpected error: %v", err)
	}

	if !store.deleteCalled {
		t.Error("store Delete was not called")
	}

	out := buf.String()
	if !strings.Contains(out, "removed") {
		t.Errorf("output = %q, want to contain 'removed'", out)
	}
}

type authClientStubError struct {
	userErr        error
	membershipsErr error
}

func (a authClientStubError) GetMe(ctx context.Context) (*asana.User, error) {
	if a.userErr != nil {
		return nil, a.userErr
	}
	return &asana.User{GID: "1", Name: "Test"}, nil
}

func (a authClientStubError) ListWorkspaceMembershipsForUser(ctx context.Context, userGID string) (*asana.WorkspaceMembershipList, error) {
	if a.membershipsErr != nil {
		return nil, a.membershipsErr
	}
	return &asana.WorkspaceMembershipList{Data: []asana.WorkspaceMembership{}}, nil
}

func TestAuthStatusGetMeError(t *testing.T) {
	ctx := &cli.Context{
		Stdout: &bytes.Buffer{},
		Stderr: &bytes.Buffer{},
		Client: authClientStubError{userErr: context.DeadlineExceeded},
	}

	cmd := AuthStatusCmd{}
	err := cmd.Run(context.Background(), ctx)
	if err == nil {
		t.Fatal("AuthStatusCmd.Run should return error when GetMe fails")
	}
	if !strings.Contains(err.Error(), "failed to get user") {
		t.Errorf("error = %q, want to contain 'failed to get user'", err.Error())
	}
}

func TestAuthStatusListMembershipsError(t *testing.T) {
	ctx := &cli.Context{
		Stdout: &bytes.Buffer{},
		Stderr: &bytes.Buffer{},
		Client: authClientStubError{membershipsErr: context.DeadlineExceeded},
	}

	cmd := AuthStatusCmd{}
	err := cmd.Run(context.Background(), ctx)
	if err == nil {
		t.Fatal("AuthStatusCmd.Run should return error when ListWorkspaceMembershipsForUser fails")
	}
	if !strings.Contains(err.Error(), "failed to list workspace memberships") {
		t.Errorf("error = %q, want to contain 'failed to list workspace memberships'", err.Error())
	}
}

func TestAuthStatusInvalidClient(t *testing.T) {
	ctx := &cli.Context{
		Stdout: &bytes.Buffer{},
		Stderr: &bytes.Buffer{},
		Client: struct{}{},
	}

	cmd := AuthStatusCmd{}
	err := cmd.Run(context.Background(), ctx)
	if err == nil {
		t.Fatal("AuthStatusCmd.Run should return error for unsupported client type")
	}
	if !strings.Contains(err.Error(), "client does not support") {
		t.Errorf("error = %q, want to contain 'client does not support'", err.Error())
	}
}

// TestAuthLoginNoStore removed - the command creates a default keyring store when none provided

// TestAuthLogoutNoStore removed - the command creates a default keyring store when none provided

type storeSetError struct{}

func (s *storeSetError) Get(service, user string) (string, error) { return "", auth.ErrNoToken }
func (s *storeSetError) Set(service, user, token string) error    { return context.DeadlineExceeded }
func (s *storeSetError) Delete(service, user string) error        { return nil }

type storeDeleteError struct{}

func (s *storeDeleteError) Get(service, user string) (string, error) { return "", auth.ErrNoToken }
func (s *storeDeleteError) Set(service, user, token string) error    { return nil }
func (s *storeDeleteError) Delete(service, user string) error        { return context.DeadlineExceeded }

func TestAuthLoginStoreError(t *testing.T) {
	ctx := &cli.Context{
		Stdout: &bytes.Buffer{},
		Stderr: &bytes.Buffer{},
		Store:  &storeSetError{},
	}

	cmd := AuthLoginCmd{Token: "test"}
	err := cmd.Run(context.Background(), ctx)
	if err == nil {
		t.Fatal("AuthLoginCmd.Run should return error when store Set fails")
	}
	// Error is returned directly from store.Set without wrapping
	if !strings.Contains(err.Error(), "deadline exceeded") {
		t.Errorf("error = %q, want to contain 'deadline exceeded'", err.Error())
	}
}

func TestAuthLogoutStoreError(t *testing.T) {
	ctx := &cli.Context{
		Stdout: &bytes.Buffer{},
		Stderr: &bytes.Buffer{},
		Store:  &storeDeleteError{},
	}

	cmd := AuthLogoutCmd{}
	err := cmd.Run(context.Background(), ctx)
	if err == nil {
		t.Fatal("AuthLogoutCmd.Run should return error when store Delete fails")
	}
	// Error is returned directly from store.Delete without wrapping
	if !strings.Contains(err.Error(), "deadline exceeded") {
		t.Errorf("error = %q, want to contain 'deadline exceeded'", err.Error())
	}
}

func TestAuthStatusMultipleWorkspaces(t *testing.T) {
	buf := &bytes.Buffer{}
	ctx := &cli.Context{
		Stdout: buf,
		Stderr: &bytes.Buffer{},
		Client: authClientStub{
			user: &asana.User{GID: "1", Name: "Multi User"},
			memberships: &asana.WorkspaceMembershipList{
				Data: []asana.WorkspaceMembership{
					{GID: "wm1", Workspace: asana.Workspace{GID: "w1", Name: "Workspace 1"}, IsActive: true},
					{GID: "wm2", Workspace: asana.Workspace{GID: "w2", Name: "Workspace 2"}, IsActive: true, IsAdmin: true},
					{GID: "wm3", Workspace: asana.Workspace{GID: "w3", Name: "Workspace 3"}, IsGuest: true},
				},
			},
		},
	}

	cmd := AuthStatusCmd{}
	if err := cmd.Run(context.Background(), ctx); err != nil {
		t.Fatalf("AuthStatusCmd.Run returned unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "Workspace 1") || !strings.Contains(out, "Workspace 2") || !strings.Contains(out, "Workspace 3") {
		t.Errorf("output missing workspaces\ngot: %q\nwant: contains all workspace names", out)
	}
}
