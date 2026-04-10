package asana

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type UserDetail struct {
	GID   string `json:"gid"`
	Name  string `json:"name"`
	Email string `json:"email,omitempty"`
}

type UserList struct {
	Data     []UserDetail `json:"data"`
	NextPage *Page        `json:"next_page,omitempty"`
}

func (c *Client) GetMe(ctx context.Context) (*User, error) {
	var user User
	if err := c.do(ctx, http.MethodGet, "/users/me", nil, nil, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (c *Client) ListUsersInWorkspace(ctx context.Context, workspaceGID string, limit int) (*UserList, error) {
	query := map[string]string{
		"opt_fields": "name,email",
	}
	if limit > 0 {
		query["limit"] = fmt.Sprintf("%d", limit)
	}
	var resp responseEnvelope
	if err := c.doRaw(ctx, http.MethodGet, "/workspaces/"+workspaceGID+"/users", query, nil, &resp); err != nil {
		return nil, err
	}

	var list UserList
	if err := json.Unmarshal(resp.Data, &list.Data); err != nil {
		return nil, err
	}
	list.NextPage = resp.NextPage
	return &list, nil
}

func (c *Client) GetUser(ctx context.Context, userGID string) (*UserDetail, error) {
	query := map[string]string{
		"opt_fields": "name,email",
	}
	var user UserDetail
	if err := c.do(ctx, http.MethodGet, "/users/"+userGID, query, nil, &user); err != nil {
		return nil, err
	}
	return &user, nil
}
