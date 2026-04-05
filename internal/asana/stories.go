package asana

import (
	"context"
	"encoding/json"
	"fmt"
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

type CreateStoryRequest struct {
	Text string `json:"text"`
}

type UpdateStoryRequest struct {
	Text string `json:"text"`
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

func (c *Client) CreateStory(ctx context.Context, taskGID string, req CreateStoryRequest) (*Story, error) {
	query := map[string]string{
		"opt_fields": "created_at,created_by.name,text,type",
	}
	payload := map[string]CreateStoryRequest{"data": req}
	var story Story
	if err := c.do(ctx, http.MethodPost, "/tasks/"+taskGID+"/stories", query, payload, &story); err != nil {
		return nil, err
	}
	return &story, nil
}

func (c *Client) UpdateStory(ctx context.Context, storyGID string, req UpdateStoryRequest) (*Story, error) {
	query := map[string]string{
		"opt_fields": "created_at,created_by.name,text,type",
	}
	payload := map[string]UpdateStoryRequest{"data": req}
	var story Story
	if err := c.do(ctx, http.MethodPut, "/stories/"+storyGID, query, payload, &story); err != nil {
		return nil, err
	}
	return &story, nil
}

func (c *Client) DeleteStory(ctx context.Context, storyGID string) error {
	var resp responseEnvelope
	if err := c.doRaw(ctx, http.MethodDelete, fmt.Sprintf("/stories/%s", storyGID), nil, nil, &resp); err != nil {
		return err
	}
	return nil
}
