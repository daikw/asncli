package cmd

import (
	"context"
	"fmt"

	"github.com/michalvavra/asncli/internal/asana"
	"github.com/michalvavra/asncli/internal/cli"
)

type ProjectsCmd struct {
	List ProjectsListCmd `cmd:"" help:"List projects."`
	Get  ProjectsGetCmd  `cmd:"" help:"Get a project."`
}

type ProjectsListCmd struct {
	Workspace string `help:"Workspace GID (uses default if not specified)."`
	Archived  *bool  `help:"Filter by archived state."`
	Limit     int    `help:"Maximum number of projects."`
}

type ProjectsGetCmd struct {
	GID string `arg:"" help:"Project GID."`
}

type projectsListClient interface {
	ListProjects(ctx context.Context, params asana.ListProjectsParams) (*asana.ProjectList, error)
}

type projectsGetClient interface {
	GetProject(ctx context.Context, gid string) (*asana.Project, error)
}

func (cmd *ProjectsListCmd) Run(ctx context.Context, c *cli.Context) error {
	// Resolve workspace
	workspace, err := c.ResolveWorkspace(cmd.Workspace)
	if err != nil {
		return err
	}

	clientAny := c.ClientOrDefault()
	client, ok := clientAny.(projectsListClient)
	if !ok {
		return fmt.Errorf("failed to list projects: client does not support listing projects")
	}
	params := asana.ListProjectsParams{
		Workspace: workspace,
		Archived:  cmd.Archived,
		Limit:     cmd.Limit,
		OptFields: "name,archived,color,created_at,workspace.name,workspace.gid",
	}
	list, err := client.ListProjects(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to list projects: %w", err)
	}

	renderer := c.RendererOrDefault()
	if c.JSON {
		return renderer.Envelope(list.Data, list.NextPage, nil)
	}

	rows := make([][]string, 0, len(list.Data))
	for _, project := range list.Data {
		workspace := ""
		if project.Workspace != nil {
			workspace = project.Workspace.Name
		}
		rows = append(rows, []string{
			project.GID,
			project.Name,
			fmt.Sprintf("%t", project.Archived),
			workspace,
		})
	}
	return renderer.Table([]string{"GID", "NAME", "ARCHIVED", "WORKSPACE"}, rows)
}

func (cmd *ProjectsGetCmd) Run(ctx context.Context, c *cli.Context) error {
	clientAny := c.ClientOrDefault()
	client, ok := clientAny.(projectsGetClient)
	if !ok {
		return fmt.Errorf("failed to get project: client does not support getting projects")
	}
	project, err := client.GetProject(ctx, cmd.GID)
	if err != nil {
		return fmt.Errorf("failed to get project: %w", err)
	}

	renderer := c.RendererOrDefault()
	if c.JSON {
		return renderer.JSON(project)
	}

	workspace := ""
	if project.Workspace != nil {
		workspace = project.Workspace.Name
	}
	return renderer.Message("GID: %s\nName: %s\nArchived: %t\nColor: %s\nWorkspace: %s\nCreated: %s\n",
		project.GID, project.Name, project.Archived, project.Color, workspace, project.CreatedAt)
}
