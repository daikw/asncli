package asana

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type CustomField struct {
	GID          string           `json:"gid"`
	Name         string           `json:"name"`
	Type         string           `json:"type,omitempty"`
	DisplayValue string           `json:"display_value,omitempty"`
	TextValue    string           `json:"text_value,omitempty"`
	NumberValue  *float64         `json:"number_value,omitempty"`
	EnumValue    *CustomFieldEnum `json:"enum_value,omitempty"`
}

type CustomFieldEnum struct {
	Name string `json:"name"`
}

type CustomFieldEnumOption struct {
	GID     string `json:"gid,omitempty"`
	Name    string `json:"name"`
	Color   string `json:"color,omitempty"`
	Enabled *bool  `json:"enabled,omitempty"`
}

type CustomFieldDefinition struct {
	GID                 string                  `json:"gid"`
	Name                string                  `json:"name"`
	Type                string                  `json:"type,omitempty"`
	ResourceSubtype     string                  `json:"resource_subtype,omitempty"`
	Description         string                  `json:"description,omitempty"`
	EnumOptions         []CustomFieldEnumOption `json:"enum_options,omitempty"`
	Precision           *int                    `json:"precision,omitempty"`
	Format              string                  `json:"format,omitempty"`
	CurrencyCode        string                  `json:"currency_code,omitempty"`
	CustomLabel         string                  `json:"custom_label,omitempty"`
	CustomLabelPosition string                  `json:"custom_label_position,omitempty"`
	InputRestrictions   []string                `json:"input_restrictions,omitempty"`
}

type CustomFieldList struct {
	Data     []CustomFieldDefinition `json:"data"`
	NextPage *Page                   `json:"next_page,omitempty"`
}

type ListCustomFieldsParams struct {
	Limit     int
	Offset    string
	OptFields string
}

const CustomFieldDetailOptFields = "name,type,resource_subtype,description,enum_options.name,enum_options.color,enum_options.enabled,precision,format,currency_code,custom_label,custom_label_position"

type UpdateCustomFieldRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

func (c *Client) GetCustomField(ctx context.Context, gid string) (*CustomFieldDefinition, error) {
	query := map[string]string{
		"opt_fields": CustomFieldDetailOptFields,
	}
	var field CustomFieldDefinition
	if err := c.do(ctx, http.MethodGet, "/custom_fields/"+gid, query, nil, &field); err != nil {
		return nil, err
	}
	return &field, nil
}

func (c *Client) UpdateCustomField(ctx context.Context, gid string, req UpdateCustomFieldRequest) (*CustomFieldDefinition, error) {
	query := map[string]string{
		"opt_fields": CustomFieldDetailOptFields,
	}
	payload := map[string]UpdateCustomFieldRequest{"data": req}
	var field CustomFieldDefinition
	if err := c.do(ctx, http.MethodPut, "/custom_fields/"+gid, query, payload, &field); err != nil {
		return nil, err
	}
	return &field, nil
}

func (c *Client) ListCustomFieldsForWorkspace(ctx context.Context, workspaceGID string, params ListCustomFieldsParams) (*CustomFieldList, error) {
	query := map[string]string{}
	if params.Limit > 0 {
		query["limit"] = fmt.Sprintf("%d", params.Limit)
	}
	if params.Offset != "" {
		query["offset"] = params.Offset
	}
	if params.OptFields != "" {
		query["opt_fields"] = params.OptFields
	}

	var resp responseEnvelope
	if err := c.doRaw(ctx, http.MethodGet, "/workspaces/"+workspaceGID+"/custom_fields", query, nil, &resp); err != nil {
		return nil, err
	}

	var list CustomFieldList
	if err := json.Unmarshal(resp.Data, &list.Data); err != nil {
		return nil, err
	}
	list.NextPage = resp.NextPage
	return &list, nil
}
