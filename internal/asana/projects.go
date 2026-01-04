package asana

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Project struct {
	GID       string     `json:"gid"`
	Name      string     `json:"name"`
	Archived  bool       `json:"archived"`
	Color     string     `json:"color,omitempty"`
	CreatedAt string     `json:"created_at,omitempty"`
	Workspace *Workspace `json:"workspace,omitempty"`
}

type ProjectList struct {
	Data     []Project `json:"data"`
	NextPage *Page     `json:"next_page,omitempty"`
}

type ListProjectsParams struct {
	Workspace string
	Archived  *bool
	Limit     int
	OptFields string
}

func (c *Client) ListProjects(ctx context.Context, params ListProjectsParams) (*ProjectList, error) {
	query := map[string]string{}
	if params.Workspace != "" {
		query["workspace"] = params.Workspace
	}
	if params.Archived != nil {
		query["archived"] = fmt.Sprintf("%t", *params.Archived)
	}
	if params.Limit > 0 {
		query["limit"] = fmt.Sprintf("%d", params.Limit)
	}
	if params.OptFields != "" {
		query["opt_fields"] = params.OptFields
	}

	var resp responseEnvelope
	if err := c.doRaw(ctx, http.MethodGet, "/projects", query, nil, &resp); err != nil {
		return nil, err
	}

	var list ProjectList
	if err := json.Unmarshal(resp.Data, &list.Data); err != nil {
		return nil, err
	}
	list.NextPage = resp.NextPage
	return &list, nil
}

func (c *Client) GetProject(ctx context.Context, gid string) (*Project, error) {
	query := map[string]string{
		"opt_fields": "name,archived,color,created_at,workspace.name,workspace.gid",
	}
	var project Project
	if err := c.do(ctx, http.MethodGet, "/projects/"+gid, query, nil, &project); err != nil {
		return nil, err
	}
	return &project, nil
}
