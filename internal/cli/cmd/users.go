package cmd

import (
	"context"
	"fmt"

	"github.com/michalvavra/asncli/internal/asana"
	"github.com/michalvavra/asncli/internal/cli"
)

type UsersCmd struct {
	List UsersListCmd `cmd:"" help:"List users in a workspace."`
	Get  UsersGetCmd  `cmd:"" help:"Get a user."`
}

type UsersListCmd struct {
	Workspace string `help:"Workspace GID (uses default if not specified)."`
	Limit     int    `help:"Maximum number of users."`
}

type UsersGetCmd struct {
	GID string `arg:"" help:"User GID or 'me'."`
}

type usersListClient interface {
	ListUsersInWorkspace(ctx context.Context, workspaceGID string, limit int) (*asana.UserList, error)
}

type usersGetClient interface {
	GetUser(ctx context.Context, userGID string) (*asana.UserDetail, error)
}

func (cmd *UsersListCmd) Run(ctx context.Context, c *cli.Context) error {
	workspace, err := c.ResolveWorkspace(cmd.Workspace)
	if err != nil {
		return err
	}

	clientAny := c.ClientOrDefault()
	client, ok := clientAny.(usersListClient)
	if !ok {
		return fmt.Errorf("failed to list users: client does not support listing users")
	}
	list, err := client.ListUsersInWorkspace(ctx, workspace, cmd.Limit)
	if err != nil {
		return fmt.Errorf("failed to list users: %w", err)
	}

	renderer := c.RendererOrDefault()
	if c.JSON {
		return renderer.Envelope(list.Data, list.NextPage, nil)
	}

	rows := make([][]string, 0, len(list.Data))
	for _, u := range list.Data {
		rows = append(rows, []string{u.GID, u.Name, u.Email})
	}
	return renderer.Table([]string{"GID", "NAME", "EMAIL"}, rows)
}

func (cmd *UsersGetCmd) Run(ctx context.Context, c *cli.Context) error {
	clientAny := c.ClientOrDefault()
	client, ok := clientAny.(usersGetClient)
	if !ok {
		return fmt.Errorf("failed to get user: client does not support getting users")
	}
	user, err := client.GetUser(ctx, cmd.GID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	renderer := c.RendererOrDefault()
	if c.JSON {
		return renderer.JSON(user)
	}

	rows := [][]string{
		{"GID", user.GID},
		{"Name", user.Name},
		{"Email", user.Email},
	}
	return renderer.Table([]string{"FIELD", "VALUE"}, rows)
}
