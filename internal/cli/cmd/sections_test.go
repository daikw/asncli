package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/michalvavra/asncli/internal/asana"
	"github.com/michalvavra/asncli/internal/cli"
)

type sectionsClient struct {
	list            *asana.SectionList
	createdSection  *asana.Section
	gotProjectGID   string
	gotCreateReq    asana.CreateSectionRequest
}

func (f *sectionsClient) ListSectionsForProject(ctx context.Context, projectGID string) (*asana.SectionList, error) {
	f.gotProjectGID = projectGID
	return f.list, nil
}

func (f *sectionsClient) CreateSectionForProject(ctx context.Context, projectGID string, req asana.CreateSectionRequest) (*asana.Section, error) {
	f.gotProjectGID = projectGID
	f.gotCreateReq = req
	return f.createdSection, nil
}

func TestSectionsListHumanReadable(t *testing.T) {
	buf := &bytes.Buffer{}
	client := &sectionsClient{list: &asana.SectionList{Data: []asana.Section{
		{GID: "sec-1", Name: "To Do"},
		{GID: "sec-2", Name: "In Progress"},
	}}}
	ctx := &cli.Context{
		Stdout: buf,
		Stderr: &bytes.Buffer{},
		Client: client,
	}

	cmd := SectionsListCmd{Project: "proj-1"}
	if err := cmd.Run(context.Background(), ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "To Do") || !strings.Contains(out, "In Progress") {
		t.Errorf("output = %q, want to contain section names", out)
	}
	if client.gotProjectGID != "proj-1" {
		t.Errorf("project GID = %q, want %q", client.gotProjectGID, "proj-1")
	}
}

func TestSectionsListJSON(t *testing.T) {
	buf := &bytes.Buffer{}
	client := &sectionsClient{list: &asana.SectionList{Data: []asana.Section{{GID: "sec-1", Name: "Done"}}}}
	ctx := &cli.Context{
		Stdout: buf,
		Stderr: &bytes.Buffer{},
		JSON:   true,
		Client: client,
	}

	cmd := SectionsListCmd{Project: "proj-1"}
	if err := cmd.Run(context.Background(), ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var env cli.Envelope
	if err := json.Unmarshal(buf.Bytes(), &env); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}
	if env.Data == nil {
		t.Fatal("envelope data is nil")
	}
}

func TestSectionsListInvalidClient(t *testing.T) {
	ctx := &cli.Context{
		Stdout: &bytes.Buffer{},
		Stderr: &bytes.Buffer{},
		Client: struct{}{},
	}

	cmd := SectionsListCmd{Project: "proj-1"}
	err := cmd.Run(context.Background(), ctx)
	if err == nil {
		t.Fatal("should return error for unsupported client type")
	}
	if !strings.Contains(err.Error(), "client does not support") {
		t.Errorf("error = %q, want to contain 'client does not support'", err.Error())
	}
}

func TestSectionsCreateHumanReadable(t *testing.T) {
	buf := &bytes.Buffer{}
	client := &sectionsClient{createdSection: &asana.Section{GID: "new-sec", Name: "Review"}}
	ctx := &cli.Context{
		Stdout: buf,
		Stderr: &bytes.Buffer{},
		Client: client,
	}

	cmd := SectionsCreateCmd{Project: "proj-1", Name: "Review"}
	if err := cmd.Run(context.Background(), ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "created") || !strings.Contains(out, "new-sec") {
		t.Errorf("output = %q, want to contain 'created' and 'new-sec'", out)
	}
	if client.gotCreateReq.Name != "Review" {
		t.Errorf("name = %q, want %q", client.gotCreateReq.Name, "Review")
	}
}

func TestSectionsCreateJSON(t *testing.T) {
	buf := &bytes.Buffer{}
	client := &sectionsClient{createdSection: &asana.Section{GID: "json-sec", Name: "Test"}}
	ctx := &cli.Context{
		Stdout: buf,
		Stderr: &bytes.Buffer{},
		JSON:   true,
		Client: client,
	}

	cmd := SectionsCreateCmd{Project: "proj-1", Name: "Test"}
	if err := cmd.Run(context.Background(), ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var env cli.Envelope
	if err := json.Unmarshal(buf.Bytes(), &env); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}
	if env.Data == nil {
		t.Fatal("envelope data is nil")
	}
}

func TestSectionsCreateInvalidClient(t *testing.T) {
	ctx := &cli.Context{
		Stdout: &bytes.Buffer{},
		Stderr: &bytes.Buffer{},
		Client: struct{}{},
	}

	cmd := SectionsCreateCmd{Project: "proj-1", Name: "Test"}
	err := cmd.Run(context.Background(), ctx)
	if err == nil {
		t.Fatal("should return error for unsupported client type")
	}
	if !strings.Contains(err.Error(), "client does not support") {
		t.Errorf("error = %q, want to contain 'client does not support'", err.Error())
	}
}

func TestSectionsCreateTrimsWhitespace(t *testing.T) {
	buf := &bytes.Buffer{}
	client := &sectionsClient{createdSection: &asana.Section{GID: "trim-sec"}}
	ctx := &cli.Context{
		Stdout: buf,
		Stderr: &bytes.Buffer{},
		Client: client,
	}

	cmd := SectionsCreateCmd{Project: "proj-1", Name: "  Trimmed  "}
	if err := cmd.Run(context.Background(), ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if client.gotCreateReq.Name != "Trimmed" {
		t.Errorf("name = %q, want %q", client.gotCreateReq.Name, "Trimmed")
	}
}
