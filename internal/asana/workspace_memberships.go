package asana

import (
	"context"
	"encoding/json"
	"net/http"
)

type WorkspaceMembership struct {
	GID        string    `json:"gid"`
	Workspace  Workspace `json:"workspace"`
	IsActive   bool      `json:"is_active,omitempty"`
	IsAdmin    bool      `json:"is_admin,omitempty"`
	IsGuest    bool      `json:"is_guest,omitempty"`
	IsViewOnly bool      `json:"is_view_only,omitempty"`
}

type WorkspaceMembershipList struct {
	Data     []WorkspaceMembership `json:"data"`
	NextPage *Page                 `json:"next_page,omitempty"`
}

func (c *Client) ListWorkspaceMembershipsForUser(ctx context.Context, userGID string) (*WorkspaceMembershipList, error) {
	query := map[string]string{
		"opt_fields": "is_active,is_admin,is_guest,is_view_only,workspace.name,workspace.gid",
	}
	var resp responseEnvelope
	if err := c.doRaw(ctx, http.MethodGet, "/users/"+userGID+"/workspace_memberships", query, nil, &resp); err != nil {
		return nil, err
	}

	var list WorkspaceMembershipList
	if err := json.Unmarshal(resp.Data, &list.Data); err != nil {
		return nil, err
	}
	list.NextPage = resp.NextPage
	return &list, nil
}
