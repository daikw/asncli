package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/michalvavra/asncli/internal/asana"
	"github.com/michalvavra/asncli/internal/cli"
)

type SectionsCmd struct {
	List   SectionsListCmd   `cmd:"" help:"List sections in a project."`
	Create SectionsCreateCmd `cmd:"" help:"Create a section in a project."`
}

type SectionsListCmd struct {
	Project string `arg:"" help:"Project GID."`
}

type SectionsCreateCmd struct {
	Project string `arg:"" help:"Project GID."`
	Name    string `help:"Section name." required:""`
}

type sectionsListClient interface {
	ListSectionsForProject(ctx context.Context, projectGID string) (*asana.SectionList, error)
}

type sectionsCreateClient interface {
	CreateSectionForProject(ctx context.Context, projectGID string, req asana.CreateSectionRequest) (*asana.Section, error)
}

func (cmd *SectionsListCmd) Run(ctx context.Context, c *cli.Context) error {
	clientAny := c.ClientOrDefault()
	client, ok := clientAny.(sectionsListClient)
	if !ok {
		return fmt.Errorf("failed to list sections: client does not support listing sections")
	}
	list, err := client.ListSectionsForProject(ctx, cmd.Project)
	if err != nil {
		return fmt.Errorf("failed to list sections: %w", err)
	}

	renderer := c.RendererOrDefault()
	if c.JSON {
		return renderer.Envelope(list.Data, list.NextPage, nil)
	}

	rows := make([][]string, 0, len(list.Data))
	for _, s := range list.Data {
		rows = append(rows, []string{s.GID, s.Name})
	}
	return renderer.Table([]string{"GID", "NAME"}, rows)
}

func (cmd *SectionsCreateCmd) Run(ctx context.Context, c *cli.Context) error {
	clientAny := c.ClientOrDefault()
	client, ok := clientAny.(sectionsCreateClient)
	if !ok {
		return fmt.Errorf("failed to create section: client does not support creating sections")
	}
	req := asana.CreateSectionRequest{
		Name: strings.TrimSpace(cmd.Name),
	}
	section, err := client.CreateSectionForProject(ctx, cmd.Project, req)
	if err != nil {
		return fmt.Errorf("failed to create section: %w", err)
	}

	renderer := c.RendererOrDefault()
	if c.JSON {
		return renderer.JSON(section)
	}
	return renderer.Message("created %s\n", section.GID)
}
