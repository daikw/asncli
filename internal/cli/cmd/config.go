package cmd

import (
	"context"
	"fmt"

	"github.com/michalvavra/asncli/internal/asana"
	"github.com/michalvavra/asncli/internal/cli"
	"github.com/michalvavra/asncli/internal/config"
)

type ConfigCmd struct {
	SetWorkspace ConfigSetWorkspaceCmd `cmd:"" name:"set-workspace" help:"Set default workspace interactively."`
	GetWorkspace ConfigGetWorkspaceCmd `cmd:"" name:"get-workspace" help:"Show current default workspace."`
	Show         ConfigShowCmd         `cmd:"" help:"Show all configuration."`
}

type ConfigSetWorkspaceCmd struct{}

type ConfigGetWorkspaceCmd struct{}

type ConfigShowCmd struct{}

type configClient interface {
	ListWorkspaceMembershipsForUser(ctx context.Context, userGID string) (*asana.WorkspaceMembershipList, error)
}

func (cmd *ConfigSetWorkspaceCmd) Run(ctx context.Context, c *cli.Context) error {
	clientAny := c.ClientOrDefault()
	client, ok := clientAny.(configClient)
	if !ok {
		return fmt.Errorf("failed to set workspace: client does not support listing workspaces")
	}

	// Get workspace memberships from Asana
	asanaMemberships, err := client.ListWorkspaceMembershipsForUser(ctx, "me")
	if err != nil {
		return fmt.Errorf("failed to list workspaces: %w", err)
	}

	// Convert to config format
	memberships := config.WorkspaceMembershipList{
		Data: make([]config.WorkspaceMembership, len(asanaMemberships.Data)),
	}
	for i, m := range asanaMemberships.Data {
		memberships.Data[i] = config.WorkspaceMembership{
			GID:        m.Workspace.GID,
			Name:       m.Workspace.Name,
			IsActive:   m.IsActive,
			IsAdmin:    m.IsAdmin,
			IsGuest:    m.IsGuest,
			IsViewOnly: m.IsViewOnly,
		}
	}

	// Create a lister adapter
	lister := &workspaceListerAdapter{memberships: memberships}

	// Prompt user to select workspace
	workspaceGID, err := config.PromptWorkspace(ctx, lister, c.Stdout, c.Stderr)
	if err != nil {
		return err
	}

	// Load current config
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Find workspace name
	workspaceName := ""
	for _, m := range memberships.Data {
		if m.GID == workspaceGID {
			workspaceName = m.Name
			break
		}
	}

	// Update and save
	cfg.DefaultWorkspace = workspaceGID
	cfg.DefaultWorkspaceName = workspaceName
	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	renderer := c.RendererOrDefault()
	if c.JSON {
		return renderer.JSON(map[string]string{
			"default_workspace": workspaceGID,
			"workspace_name":    workspaceName,
		})
	}
	return renderer.Message("default workspace set to: %s (%s)\n", workspaceName, workspaceGID)
}

func (cmd *ConfigGetWorkspaceCmd) Run(ctx context.Context, c *cli.Context) error {
	cfg := c.Config
	if cfg == nil {
		var err error
		cfg, err = config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}
	}

	renderer := c.RendererOrDefault()
	if c.JSON {
		result := map[string]interface{}{
			"workspace_id":   cfg.DefaultWorkspace,
			"workspace_name": cfg.DefaultWorkspaceName,
		}
		if cfg.DefaultWorkspace == "" {
			result["workspace_id"] = nil
			result["workspace_name"] = nil
		}
		return renderer.JSON(result)
	}

	if cfg.DefaultWorkspace == "" {
		return renderer.Message("no default workspace set\nSet one with: asn config set-workspace\n")
	}
	if cfg.DefaultWorkspaceName != "" {
		return renderer.Message("%s (%s)\n", cfg.DefaultWorkspaceName, cfg.DefaultWorkspace)
	}
	return renderer.Message("%s\n", cfg.DefaultWorkspace)
}

func (cmd *ConfigShowCmd) Run(ctx context.Context, c *cli.Context) error {
	cfg := c.Config
	if cfg == nil {
		var err error
		cfg, err = config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}
	}

	configPath, err := config.GetConfigPath()
	if err != nil {
		return fmt.Errorf("failed to get config path: %w", err)
	}

	renderer := c.RendererOrDefault()
	if c.JSON {
		return renderer.JSON(map[string]interface{}{
			"config_path":            configPath,
			"default_workspace":      cfg.DefaultWorkspace,
			"default_workspace_name": cfg.DefaultWorkspaceName,
		})
	}

	if err := renderer.Message("Config file: %s\n\n", configPath); err != nil {
		return err
	}

	rows := [][]string{
		{"default_workspace", cfg.DefaultWorkspace},
		{"default_workspace_name", cfg.DefaultWorkspaceName},
	}
	return renderer.Table([]string{"SETTING", "VALUE"}, rows)
}

// workspaceListerAdapter adapts static membership data to the WorkspaceLister interface.
type workspaceListerAdapter struct {
	memberships config.WorkspaceMembershipList
}

func (a *workspaceListerAdapter) ListWorkspaceMembershipsForUser(ctx context.Context, userGID string) (config.WorkspaceMembershipList, error) {
	return a.memberships, nil
}
