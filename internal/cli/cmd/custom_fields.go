package cmd

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/michalvavra/asncli/internal/asana"
	"github.com/michalvavra/asncli/internal/cli"
)

type CustomFieldsCmd struct {
	List   CustomFieldsListCmd   `cmd:"" help:"List custom fields."`
	Get    CustomFieldsGetCmd    `cmd:"" help:"Get a custom field."`
	Update CustomFieldsUpdateCmd `cmd:"" help:"Update a custom field."`
}

type CustomFieldsListCmd struct {
	Workspace string `help:"Workspace GID (uses default if not specified)."`
	Limit     int    `help:"Maximum number of custom fields."`
	Offset    string `help:"Pagination offset."`
	OptFields string `help:"Comma-separated list of fields to include."`
}

type CustomFieldsGetCmd struct {
	GID string `arg:"" help:"Custom field GID."`
}

type CustomFieldsUpdateCmd struct {
	GID         string `arg:"" help:"Custom field GID."`
	Name        string `help:"Custom field name."`
	Description string `help:"Custom field description."`
}

type customFieldsListClient interface {
	ListCustomFieldsForWorkspace(ctx context.Context, workspaceGID string, params asana.ListCustomFieldsParams) (*asana.CustomFieldList, error)
}

type customFieldsGetClient interface {
	GetCustomField(ctx context.Context, gid string) (*asana.CustomFieldDefinition, error)
}

type customFieldsUpdateClient interface {
	UpdateCustomField(ctx context.Context, gid string, req asana.UpdateCustomFieldRequest) (*asana.CustomFieldDefinition, error)
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

func (cmd *CustomFieldsGetCmd) Run(ctx context.Context, c *cli.Context) error {
	clientAny := c.ClientOrDefault()
	client, ok := clientAny.(customFieldsGetClient)
	if !ok {
		return fmt.Errorf("failed to get custom field: client does not support getting custom fields")
	}
	field, err := client.GetCustomField(ctx, cmd.GID)
	if err != nil {
		return fmt.Errorf("failed to get custom field: %w", err)
	}

	renderer := c.RendererOrDefault()
	if c.JSON {
		return renderer.JSON(field)
	}

	fieldType := field.ResourceSubtype
	if fieldType == "" {
		fieldType = field.Type
	}
	enumOptions := formatEnumOptions(field.EnumOptions)
	rows := [][]string{
		{"GID", field.GID},
		{"Name", field.Name},
		{"Type", fieldType},
		{"Description", field.Description},
		{"Enum Options", enumOptions},
	}
	return renderer.Table([]string{"FIELD", "VALUE"}, rows)
}

func (cmd *CustomFieldsUpdateCmd) Run(ctx context.Context, c *cli.Context) error {
	clientAny := c.ClientOrDefault()
	client, ok := clientAny.(customFieldsUpdateClient)
	if !ok {
		return fmt.Errorf("failed to update custom field: client does not support updating custom fields")
	}
	req := asana.UpdateCustomFieldRequest{}
	if strings.TrimSpace(cmd.Name) != "" {
		name := strings.TrimSpace(cmd.Name)
		req.Name = &name
	}
	if strings.TrimSpace(cmd.Description) != "" {
		desc := strings.TrimSpace(cmd.Description)
		req.Description = &desc
	}
	if req.Name == nil && req.Description == nil {
		return errors.New("failed to update custom field: no fields to update")
	}

	field, err := client.UpdateCustomField(ctx, cmd.GID, req)
	if err != nil {
		return fmt.Errorf("failed to update custom field: %w", err)
	}

	renderer := c.RendererOrDefault()
	if c.JSON {
		return renderer.JSON(field)
	}
	return renderer.Message("updated %s\n", field.GID)
}

func formatEnumOptions(options []asana.CustomFieldEnumOption) string {
	if len(options) == 0 {
		return "none"
	}
	names := make([]string, 0, len(options))
	for _, opt := range options {
		names = append(names, opt.Name)
	}
	return strings.Join(names, ", ")
}
