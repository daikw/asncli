package asana

import (
	"context"
	"encoding/json"
	"net/http"
)

type Story struct {
	GID       string `json:"gid"`
	CreatedAt string `json:"created_at"`
	CreatedBy *User  `json:"created_by,omitempty"`
	Text      string `json:"text"`
	Type      string `json:"type"`
}

type StoryList struct {
	Data     []Story `json:"data"`
	NextPage *Page   `json:"next_page,omitempty"`
}

func (c *Client) GetTaskStories(ctx context.Context, taskGID string) (*StoryList, error) {
	query := map[string]string{
		"opt_fields": "created_at,created_by.name,text,type",
	}

	var resp responseEnvelope
	if err := c.doRaw(ctx, http.MethodGet, "/tasks/"+taskGID+"/stories", query, nil, &resp); err != nil {
		return nil, err
	}

	var list StoryList
	if err := json.Unmarshal(resp.Data, &list.Data); err != nil {
		return nil, err
	}
	list.NextPage = resp.NextPage
	return &list, nil
}
