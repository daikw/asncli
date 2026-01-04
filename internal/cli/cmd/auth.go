package cmd

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/term"

	"github.com/michalvavra/asncli/internal/asana"
	"github.com/michalvavra/asncli/internal/auth"
	"github.com/michalvavra/asncli/internal/cli"
)

type AuthCmd struct {
	Login  AuthLoginCmd  `cmd:"" help:"Store a personal access token."`
	Logout AuthLogoutCmd `cmd:"" help:"Remove the stored token."`
	Status AuthStatusCmd `cmd:"" help:"Show current authentication status."`
}

type AuthLoginCmd struct {
	Token string `help:"Personal access token."`
}

type AuthLogoutCmd struct{}

type AuthStatusCmd struct{}

type AuthStatus struct {
	User                 *asana.User                 `json:"user"`
	WorkspaceMemberships []asana.WorkspaceMembership `json:"workspace_memberships"`
}

type authClient interface {
	GetMe(ctx context.Context) (*asana.User, error)
	ListWorkspaceMembershipsForUser(ctx context.Context, userGID string) (*asana.WorkspaceMembershipList, error)
}

func (cmd *AuthLoginCmd) Run(ctx context.Context, c *cli.Context) error {
	store := c.Store
	if store == nil {
		store = auth.NewKeyringStore()
	}
	token := strings.TrimSpace(cmd.Token)
	if token == "" {
		readToken, err := promptToken(c.Stderr)
		if err != nil {
			return err
		}
		token = readToken
	}
	if token == "" {
		return errors.New("failed to login: token is required")
	}
	if err := store.Set(auth.DefaultService, auth.DefaultUser, token); err != nil {
		return err
	}

	renderer := c.RendererOrDefault()
	if c.JSON {
		return renderer.JSON(map[string]string{"status": "stored"})
	}
	return renderer.Message("token stored\n")
}

func (cmd *AuthLogoutCmd) Run(ctx context.Context, c *cli.Context) error {
	store := c.Store
	if store == nil {
		store = auth.NewKeyringStore()
	}
	if err := store.Delete(auth.DefaultService, auth.DefaultUser); err != nil {
		return err
	}

	renderer := c.RendererOrDefault()
	if c.JSON {
		return renderer.JSON(map[string]string{"status": "removed"})
	}
	return renderer.Message("token removed\n")
}

func (cmd *AuthStatusCmd) Run(ctx context.Context, c *cli.Context) error {
	clientAny := c.ClientOrDefault()
	client, ok := clientAny.(authClient)
	if !ok {
		return errors.New("failed to get auth status: client does not support auth status")
	}
	me, err := client.GetMe(ctx)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	memberships, err := client.ListWorkspaceMembershipsForUser(ctx, "me")
	if err != nil {
		return fmt.Errorf("failed to list workspace memberships: %w", err)
	}
	if memberships == nil {
		memberships = &asana.WorkspaceMembershipList{}
	}

	renderer := c.RendererOrDefault()
	if c.JSON {
		data := AuthStatus{
			User:                 me,
			WorkspaceMemberships: memberships.Data,
		}
		return renderer.JSON(data)
	}

	if err := renderer.Message("authenticated as %s (GID: %s)\n\n", me.Name, me.GID); err != nil {
		return err
	}

	rows := make([][]string, 0, len(memberships.Data))
	for _, membership := range memberships.Data {
		rows = append(rows, []string{
			membership.Workspace.GID,
			membership.Workspace.Name,
			fmt.Sprintf("%t", membership.IsActive),
			fmt.Sprintf("%t", membership.IsAdmin),
			fmt.Sprintf("%t", membership.IsGuest),
			fmt.Sprintf("%t", membership.IsViewOnly),
		})
	}
	return renderer.Table([]string{"WORKSPACE_GID", "WORKSPACE_NAME", "ACTIVE", "ADMIN", "GUEST", "VIEW_ONLY"}, rows)
}

func promptToken(prompt io.Writer) (string, error) {
	if term.IsTerminal(int(os.Stdin.Fd())) {
		fmt.Fprint(prompt, "Personal access token: ")
		bytes, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Fprintln(prompt)
		if err != nil {
			return "", err
		}
		return string(bytes), nil
	}

	reader := bufio.NewReader(os.Stdin)
	return reader.ReadString('\n')
}
