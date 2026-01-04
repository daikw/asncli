package cmd

import (
	"context"
	"fmt"

	"github.com/michalvavra/asncli/internal/asana"
	"github.com/michalvavra/asncli/internal/cli"
)

type CustomFieldsCmd struct {
	List CustomFieldsListCmd `cmd:"" help:"List custom fields."`
}

type CustomFieldsListCmd struct {
	Workspace string `help:"Workspace GID (uses default if not specified)."`
	Limit     int    `help:"Maximum number of custom fields."`
	Offset    string `help:"Pagination offset."`
	OptFields string `help:"Comma-separated list of fields to include."`
}

type customFieldsListClient interface {
	ListCustomFieldsForWorkspace(ctx context.Context, workspaceGID string, params asana.ListCustomFieldsParams) (*asana.CustomFieldList, error)
}

func (cmd *CustomFieldsListCmd) Run(ctx context.Context, c *cli.Context) error {
	// Resolve workspace
	workspace, err := c.ResolveWorkspace(cmd.Workspace)
	if err != nil {
		return err
	}

	clientAny := c.ClientOrDefault()
	client, ok := clientAny.(customFieldsListClient)
	if !ok {
		return fmt.Errorf("failed to list custom fields: client does not support listing custom fields")
	}
	params := asana.ListCustomFieldsParams{
		Limit:     cmd.Limit,
		Offset:    cmd.Offset,
		OptFields: cmd.OptFields,
	}
	list, err := client.ListCustomFieldsForWorkspace(ctx, workspace, params)
	if err != nil {
		return fmt.Errorf("failed to list custom fields: %w", err)
	}

	renderer := c.RendererOrDefault()
	if c.JSON {
		return renderer.Envelope(list.Data, list.NextPage, nil)
	}

	rows := make([][]string, 0, len(list.Data))
	for _, field := range list.Data {
		fieldType := field.ResourceSubtype
		if fieldType == "" {
			fieldType = field.Type
		}
		rows = append(rows, []string{field.GID, field.Name, fieldType})
	}
	return renderer.Table([]string{"GID", "NAME", "TYPE"}, rows)
}
