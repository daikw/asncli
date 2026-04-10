package asana

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type SectionList struct {
	Data     []Section `json:"data"`
	NextPage *Page     `json:"next_page,omitempty"`
}

type CreateSectionRequest struct {
	Name string `json:"name"`
}

func (c *Client) ListSectionsForProject(ctx context.Context, projectGID string) (*SectionList, error) {
	query := map[string]string{
		"opt_fields": "name",
	}
	var resp responseEnvelope
	if err := c.doRaw(ctx, http.MethodGet, "/projects/"+projectGID+"/sections", query, nil, &resp); err != nil {
		return nil, err
	}

	var list SectionList
	if err := json.Unmarshal(resp.Data, &list.Data); err != nil {
		return nil, err
	}
	list.NextPage = resp.NextPage
	return &list, nil
}

func (c *Client) CreateSectionForProject(ctx context.Context, projectGID string, req CreateSectionRequest) (*Section, error) {
	query := map[string]string{
		"opt_fields": "name",
	}
	payload := map[string]CreateSectionRequest{"data": req}
	var section Section
	if err := c.do(ctx, http.MethodPost, "/projects/"+projectGID+"/sections", query, payload, &section); err != nil {
		return nil, err
	}
	return &section, nil
}

type MoveTaskToSectionRequest struct {
	Task       string `json:"task"`
	InsertBefore string `json:"insert_before,omitempty"`
	InsertAfter  string `json:"insert_after,omitempty"`
}

func (c *Client) AddTaskToSection(ctx context.Context, sectionGID string, taskGID string) error {
	payload := map[string]map[string]string{"data": {"task": taskGID}}
	var resp responseEnvelope
	return c.doRaw(ctx, http.MethodPost, fmt.Sprintf("/sections/%s/addTask", sectionGID), nil, payload, &resp)
}
