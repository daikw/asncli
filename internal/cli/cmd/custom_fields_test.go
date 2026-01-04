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

type customFieldsClientStub struct {
	list         *asana.CustomFieldList
	gotWorkspace string
	gotParams    asana.ListCustomFieldsParams
}

func (c *customFieldsClientStub) ListCustomFieldsForWorkspace(ctx context.Context, workspaceGID string, params asana.ListCustomFieldsParams) (*asana.CustomFieldList, error) {
	c.gotWorkspace = workspaceGID
	c.gotParams = params
	return c.list, nil
}

func TestCustomFieldsListJSON(t *testing.T) {
	buf := &bytes.Buffer{}
	client := &customFieldsClientStub{
		list: &asana.CustomFieldList{
			Data: []asana.CustomFieldDefinition{{GID: "cf1", Name: "Size", Type: "text"}},
			NextPage: &asana.Page{
				Offset: "next",
				URI:    "uri",
			},
		},
	}
	ctx := &cli.Context{
		Stdout: buf,
		Stderr: &bytes.Buffer{},
		JSON:   true,
		Client: client,
	}

	cmd := CustomFieldsListCmd{Workspace: "w1", Limit: 10, Offset: "next", OptFields: "name"}
	if err := cmd.Run(context.Background(), ctx); err != nil {
		t.Fatalf("CustomFieldsListCmd.Run returned unexpected error: %v", err)
	}

	if got, want := client.gotWorkspace, "w1"; got != want {
		t.Errorf("workspace = %q, want %q", got, want)
	}
	if got, want := client.gotParams.Limit, 10; got != want {
		t.Errorf("limit = %d, want %d", got, want)
	}
	if got, want := client.gotParams.Offset, "next"; got != want {
		t.Errorf("offset = %q, want %q", got, want)
	}
	if got, want := client.gotParams.OptFields, "name"; got != want {
		t.Errorf("opt_fields = %q, want %q", got, want)
	}

	var env cli.Envelope
	if err := json.Unmarshal(buf.Bytes(), &env); err != nil {
		t.Fatalf("failed to unmarshal JSON output: %v", err)
	}
	if env.Data == nil || env.Next == nil {
		t.Errorf("envelope incomplete: Data = %v, Next = %v, want both non-nil", env.Data, env.Next)
	}
}

func TestCustomFieldsListHumanReadable(t *testing.T) {
	buf := &bytes.Buffer{}
	ctx := &cli.Context{
		Stdout: buf,
		Stderr: &bytes.Buffer{},
		Client: &customFieldsClientStub{
			list: &asana.CustomFieldList{
				Data: []asana.CustomFieldDefinition{
					{GID: "cf1", Name: "Size", Type: "text"},
					{GID: "cf2", Name: "Priority", ResourceSubtype: "enum", Type: "multi_enum"},
				},
			},
		},
	}

	cmd := CustomFieldsListCmd{Workspace: "w1"}
	if err := cmd.Run(context.Background(), ctx); err != nil {
		t.Fatalf("CustomFieldsListCmd.Run returned unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "GID") || !strings.Contains(out, "Size") || !strings.Contains(out, "text") || !strings.Contains(out, "enum") {
		t.Errorf("output missing expected content\ngot: %q\nwant: contains GID, Size, text, and enum", out)
	}
}

func TestCustomFieldsListUnsupportedClient(t *testing.T) {
	ctx := &cli.Context{
		Stdout: &bytes.Buffer{},
		Stderr: &bytes.Buffer{},
		Client: struct{}{},
	}

	cmd := CustomFieldsListCmd{Workspace: "w1"}
	if err := cmd.Run(context.Background(), ctx); err == nil {
		t.Fatal("CustomFieldsListCmd.Run should return error for unsupported client type")
	}
}
