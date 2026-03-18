package asana

import (
	"context"
	"encoding/json"
	"net/http"
)

type Attachment struct {
	GID         string `json:"gid"`
	Name        string `json:"name"`
	DownloadURL string `json:"download_url,omitempty"`
	Host        string `json:"host,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
}

type AttachmentList struct {
	Data     []Attachment `json:"data"`
	NextPage *Page        `json:"next_page,omitempty"`
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
