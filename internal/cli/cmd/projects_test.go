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

type projectsClient struct {
	list    *asana.ProjectList
	project *asana.Project
}

func (p projectsClient) ListProjects(ctx context.Context, params asana.ListProjectsParams) (*asana.ProjectList, error) {
	return p.list, nil
}

func (p projectsClient) GetProject(ctx context.Context, gid string) (*asana.Project, error) {
	return p.project, nil
}

func TestProjectsListJSON(t *testing.T) {
	buf := &bytes.Buffer{}
	ctx := &cli.Context{
		Stdout: buf,
		Stderr: &bytes.Buffer{},
		JSON:   true,
		Client: projectsClient{list: &asana.ProjectList{Data: []asana.Project{{GID: "1", Name: "Alpha"}}}},
		Config: &config.Config{DefaultWorkspace: "123456"},
	}

	cmd := ProjectsListCmd{}
	if err := cmd.Run(context.Background(), ctx); err != nil {
		t.Fatalf("ProjectsListCmd.Run returned unexpected error: %v", err)
	}

	var env cli.Envelope
	if err := json.Unmarshal(buf.Bytes(), &env); err != nil {
		t.Fatalf("failed to unmarshal JSON output: %v", err)
	}
	if env.Data == nil {
		t.Fatal("envelope data is nil, want project list")
	}
}

func TestProjectsListHumanReadable(t *testing.T) {
	buf := &bytes.Buffer{}
	ctx := &cli.Context{
		Stdout: buf,
		Stderr: &bytes.Buffer{},
		Client: projectsClient{list: &asana.ProjectList{Data: []asana.Project{{
			GID:       "1",
			Name:      "Alpha",
			Archived:  false,
			Workspace: &asana.Workspace{Name: "Work"},
		}}}},
		Config: &config.Config{DefaultWorkspace: "123456"},
	}

	cmd := ProjectsListCmd{}
	if err := cmd.Run(context.Background(), ctx); err != nil {
		t.Fatalf("ProjectsListCmd.Run returned unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "GID") || !strings.Contains(out, "Alpha") || !strings.Contains(out, "Work") {
		t.Errorf("output missing expected content\ngot: %q\nwant: contains GID, Alpha, and Work", out)
	}
}

func TestProjectsGetHumanReadable(t *testing.T) {
	buf := &bytes.Buffer{}
	ctx := &cli.Context{
		Stdout: buf,
		Stderr: &bytes.Buffer{},
		Client: projectsClient{project: &asana.Project{
			GID:       "2",
			Name:      "Beta",
			Archived:  true,
			Color:     "red",
			CreatedAt: "2024-01-01T00:00:00Z",
			Workspace: &asana.Workspace{Name: "Main"},
		}},
	}

	cmd := ProjectsGetCmd{GID: "2"}
	if err := cmd.Run(context.Background(), ctx); err != nil {
		t.Fatalf("ProjectsGetCmd.Run returned unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "GID: 2") || !strings.Contains(out, "Name: Beta") || !strings.Contains(out, "Workspace: Main") {
		t.Errorf("output missing project details\ngot: %q\nwant: contains 'GID: 2', 'Name: Beta', and 'Workspace: Main'", out)
	}
}

func TestProjectsGetJSON(t *testing.T) {
	buf := &bytes.Buffer{}
	ctx := &cli.Context{
		Stdout: buf,
		Stderr: &bytes.Buffer{},
		JSON:   true,
		Client: projectsClient{project: &asana.Project{GID: "2", Name: "Beta"}},
	}

	cmd := ProjectsGetCmd{GID: "2"}
	if err := cmd.Run(context.Background(), ctx); err != nil {
		t.Fatalf("ProjectsGetCmd.Run returned unexpected error: %v", err)
	}

	var env cli.Envelope
	if err := json.Unmarshal(buf.Bytes(), &env); err != nil {
		t.Fatalf("failed to unmarshal JSON output: %v", err)
	}
	if env.Data == nil {
		t.Fatal("envelope data is nil, want project")
	}
}
