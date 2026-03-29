package asana

import (
	"context"
	"encoding/json"
	"net/http"
)

type Attachment struct {
	GID          string       `json:"gid"`
	Name         string       `json:"name"`
	DownloadURL  string       `json:"download_url,omitempty"`
	Host         string       `json:"host,omitempty"`
	CreatedAt    string       `json:"created_at,omitempty"`
	ViewURL      string       `json:"view_url,omitempty"`
	PermanentURL string       `json:"permanent_url,omitempty"`
	Size         int64        `json:"size,omitempty"`
	ResourceType string       `json:"resource_type,omitempty"`
	Parent       *TaskCompact `json:"parent,omitempty"`
}

type AttachmentList struct {
	Data     []Attachment `json:"data"`
	NextPage *Page        `json:"next_page,omitempty"`
}

const AttachmentDetailOptFields = "name,download_url,host,created_at,view_url,permanent_url,size,resource_type,parent.gid,parent.name"

func (c *Client) GetAttachment(ctx context.Context, gid string) (*Attachment, error) {
	query := map[string]string{
		"opt_fields": AttachmentDetailOptFields,
	}
	var attachment Attachment
	if err := c.do(ctx, http.MethodGet, "/attachments/"+gid, query, nil, &attachment); err != nil {
		return nil, err
	}
	return &attachment, nil
}

func (c *Client) GetTaskAttachments(ctx context.Context, taskGID string) (*AttachmentList, error) {
	query := map[string]string{
		"opt_fields": "name,download_url,host,created_at",
	}

	var resp responseEnvelope
	if err := c.doRaw(ctx, http.MethodGet, "/tasks/"+taskGID+"/attachments", query, nil, &resp); err != nil {
		return nil, err
	}

	var list AttachmentList
	if err := json.Unmarshal(resp.Data, &list.Data); err != nil {
		return nil, err
	}
	list.NextPage = resp.NextPage
	return &list, nil
}
