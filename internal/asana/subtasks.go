package asana

import (
	"context"
	"encoding/json"
	"net/http"
)

type SubtaskList struct {
	Data     []Task `json:"data"`
	NextPage *Page  `json:"next_page,omitempty"`
}

func (c *Client) GetSubtasks(ctx context.Context, taskGID string) (*SubtaskList, error) {
	query := map[string]string{
		"opt_fields": "name,completed,assignee.name,due_on",
	}

	var resp responseEnvelope
	if err := c.doRaw(ctx, http.MethodGet, "/tasks/"+taskGID+"/subtasks", query, nil, &resp); err != nil {
		return nil, err
	}

	var list SubtaskList
	if err := json.Unmarshal(resp.Data, &list.Data); err != nil {
		return nil, err
	}
	list.NextPage = resp.NextPage
	return &list, nil
}
