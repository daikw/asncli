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
	list           *asana.CustomFieldList
	field          *asana.CustomFieldDefinition
	updatedField   *asana.CustomFieldDefinition
	gotWorkspace   string
	gotParams      asana.ListCustomFieldsParams
	gotUpdateGID   string
	gotUpdateReq   asana.UpdateCustomFieldRequest
}

func (c *customFieldsClientStub) ListCustomFieldsForWorkspace(ctx context.Context, workspaceGID string, params asana.ListCustomFieldsParams) (*asana.CustomFieldList, error) {
	c.gotWorkspace = workspaceGID
	c.gotParams = params
	return c.list, nil
}

func (c *customFieldsClientStub) GetCustomField(ctx context.Context, gid string) (*asana.CustomFieldDefinition, error) {
	return c.field, nil
}

func (c *customFieldsClientStub) UpdateCustomField(ctx context.Context, gid string, req asana.UpdateCustomFieldRequest) (*asana.CustomFieldDefinition, error) {
	c.gotUpdateGID = gid
	c.gotUpdateReq = req
	return c.updatedField, nil
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

func TestCustomFieldsGetHumanReadable(t *testing.T) {
	buf := &bytes.Buffer{}
	client := &customFieldsClientStub{field: &asana.CustomFieldDefinition{
		GID:             "cf-1",
		Name:            "Priority",
		ResourceSubtype: "enum",
		Description:     "Task priority level",
		EnumOptions: []asana.CustomFieldEnumOption{
			{Name: "High"},
			{Name: "Medium"},
			{Name: "Low"},
		},
	}}
	ctx := &cli.Context{
		Stdout: buf,
		Stderr: &bytes.Buffer{},
		Client: client,
	}

	cmd := CustomFieldsGetCmd{GID: "cf-1"}
	if err := cmd.Run(context.Background(), ctx); err != nil {
		t.Fatalf("CustomFieldsGetCmd.Run returned unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "cf-1") || !strings.Contains(out, "Priority") || !strings.Contains(out, "enum") {
		t.Errorf("output missing expected content\ngot: %q", out)
	}
	if !strings.Contains(out, "High") || !strings.Contains(out, "Medium") || !strings.Contains(out, "Low") {
		t.Errorf("output missing enum options\ngot: %q", out)
	}
}

func TestCustomFieldsGetJSON(t *testing.T) {
	buf := &bytes.Buffer{}
	client := &customFieldsClientStub{field: &asana.CustomFieldDefinition{GID: "cf-json", Name: "Size"}}
	ctx := &cli.Context{
		Stdout: buf,
		Stderr: &bytes.Buffer{},
		JSON:   true,
		Client: client,
	}

	cmd := CustomFieldsGetCmd{GID: "cf-json"}
	if err := cmd.Run(context.Background(), ctx); err != nil {
		t.Fatalf("CustomFieldsGetCmd.Run returned unexpected error: %v", err)
	}

	var env cli.Envelope
	if err := json.Unmarshal(buf.Bytes(), &env); err != nil {
		t.Fatalf("failed to unmarshal JSON output: %v", err)
	}
	if env.Data == nil {
		t.Fatal("envelope data is nil, want custom field")
	}
}

func TestCustomFieldsGetInvalidClient(t *testing.T) {
	ctx := &cli.Context{
		Stdout: &bytes.Buffer{},
		Stderr: &bytes.Buffer{},
		Client: struct{}{},
	}

	cmd := CustomFieldsGetCmd{GID: "cf-1"}
	err := cmd.Run(context.Background(), ctx)
	if err == nil {
		t.Fatal("CustomFieldsGetCmd.Run should return error for unsupported client type")
	}
}

func TestCustomFieldsUpdateHumanReadable(t *testing.T) {
	buf := &bytes.Buffer{}
	client := &customFieldsClientStub{updatedField: &asana.CustomFieldDefinition{GID: "cf-upd", Name: "Updated"}}
	ctx := &cli.Context{
		Stdout: buf,
		Stderr: &bytes.Buffer{},
		Client: client,
	}

	cmd := CustomFieldsUpdateCmd{GID: "cf-upd", Name: "Updated"}
	if err := cmd.Run(context.Background(), ctx); err != nil {
		t.Fatalf("CustomFieldsUpdateCmd.Run returned unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "updated") || !strings.Contains(out, "cf-upd") {
		t.Errorf("output = %q, want to contain 'updated' and 'cf-upd'", out)
	}
	if client.gotUpdateGID != "cf-upd" {
		t.Errorf("GID = %q, want %q", client.gotUpdateGID, "cf-upd")
	}
	if client.gotUpdateReq.Name == nil || *client.gotUpdateReq.Name != "Updated" {
		t.Errorf("name = %v, want 'Updated'", client.gotUpdateReq.Name)
	}
}

func TestCustomFieldsUpdateJSON(t *testing.T) {
	buf := &bytes.Buffer{}
	client := &customFieldsClientStub{updatedField: &asana.CustomFieldDefinition{GID: "cf-json-upd", Name: "J"}}
	ctx := &cli.Context{
		Stdout: buf,
		Stderr: &bytes.Buffer{},
		JSON:   true,
		Client: client,
	}

	cmd := CustomFieldsUpdateCmd{GID: "cf-json-upd", Name: "J"}
	if err := cmd.Run(context.Background(), ctx); err != nil {
		t.Fatalf("CustomFieldsUpdateCmd.Run returned unexpected error: %v", err)
	}

	var env cli.Envelope
	if err := json.Unmarshal(buf.Bytes(), &env); err != nil {
		t.Fatalf("failed to unmarshal JSON output: %v", err)
	}
	if env.Data == nil {
		t.Fatal("envelope data is nil, want custom field")
	}
}

func TestCustomFieldsUpdateNoFields(t *testing.T) {
	cmd := CustomFieldsUpdateCmd{GID: "cf-1"}
	ctx := &cli.Context{}
	if err := cmd.Run(context.Background(), ctx); err == nil {
		t.Fatal("CustomFieldsUpdateCmd.Run should return error when no fields specified")
	}
}

func TestCustomFieldsUpdateInvalidClient(t *testing.T) {
	ctx := &cli.Context{
		Stdout: &bytes.Buffer{},
		Stderr: &bytes.Buffer{},
		Client: struct{}{},
	}

	cmd := CustomFieldsUpdateCmd{GID: "cf-1", Name: "test"}
	err := cmd.Run(context.Background(), ctx)
	if err == nil {
		t.Fatal("CustomFieldsUpdateCmd.Run should return error for unsupported client type")
	}
}
