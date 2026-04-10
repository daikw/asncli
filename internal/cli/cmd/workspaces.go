package cmd

import (
	"context"
	"fmt"

	"github.com/michalvavra/asncli/internal/asana"
	"github.com/michalvavra/asncli/internal/cli"
)

type WorkspacesCmd struct {
	List WorkspacesListCmd `cmd:"" help:"List workspaces for the current user."`
}

type WorkspacesListCmd struct{}

type workspacesListClient interface {
	ListWorkspaceMembershipsForUser(ctx context.Context, userGID string) (*asana.WorkspaceMembershipList, error)
}

func (cmd *WorkspacesListCmd) Run(ctx context.Context, c *cli.Context) error {
	clientAny := c.ClientOrDefault()
	client, ok := clientAny.(workspacesListClient)
	if !ok {
		return fmt.Errorf("failed to list workspaces: client does not support listing workspaces")
	}
	list, err := client.ListWorkspaceMembershipsForUser(ctx, "me")
	if err != nil {
		return fmt.Errorf("failed to list workspaces: %w", err)
	}

	renderer := c.RendererOrDefault()
	if c.JSON {
		return renderer.Envelope(list.Data, list.NextPage, nil)
	}

	rows := make([][]string, 0, len(list.Data))
	for _, m := range list.Data {
		rows = append(rows, []string{m.Workspace.GID, m.Workspace.Name, fmt.Sprintf("%t", m.IsActive)})
	}
	return renderer.Table([]string{"GID", "NAME", "ACTIVE"}, rows)
}
