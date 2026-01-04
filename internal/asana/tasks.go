package asana

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Task struct {
	GID          string        `json:"gid"`
	Name         string        `json:"name"`
	Completed    bool          `json:"completed"`
	Assignee     *User         `json:"assignee"`
	Notes        string        `json:"notes,omitempty"`
	HTMLNotes    string        `json:"html_notes,omitempty"`
	DueOn        string        `json:"due_on,omitempty"`
	DueAt        string        `json:"due_at,omitempty"`
	StartOn      string        `json:"start_on,omitempty"`
	StartAt      string        `json:"start_at,omitempty"`
	CreatedAt    string        `json:"created_at,omitempty"`
	ModifiedAt   string        `json:"modified_at,omitempty"`
	CompletedAt  string        `json:"completed_at,omitempty"`
	PermalinkURL string        `json:"permalink_url,omitempty"`
	Projects     []Project     `json:"projects,omitempty"`
	Memberships  []Membership  `json:"memberships,omitempty"`
	Tags         []Tag         `json:"tags,omitempty"`
	Followers    []User        `json:"followers,omitempty"`
	Parent       *TaskCompact  `json:"parent,omitempty"`
	CustomFields []CustomField `json:"custom_fields,omitempty"`
}

type TaskCompact struct {
	GID  string `json:"gid"`
	Name string `json:"name"`
}

type Tag struct {
	GID  string `json:"gid"`
	Name string `json:"name"`
}

type Membership struct {
	Project Project `json:"project"`
	Section Section `json:"section"`
}

type Section struct {
	GID  string `json:"gid"`
	Name string `json:"name"`
}

const TaskDetailOptFields = "name,completed,assignee.name,assignee.gid,assignee.email,notes,html_notes,due_on,due_at,start_on,start_at,created_at,modified_at,completed_at,permalink_url,projects.name,projects.gid,memberships.project.name,memberships.section.name,tags.name,tags.gid,followers.name,followers.gid,parent.name,parent.gid,custom_fields.name,custom_fields.type,custom_fields.display_value,custom_fields.text_value,custom_fields.number_value,custom_fields.enum_value.name,custom_fields.enum_value.color,custom_fields.multi_enum_values.name,custom_fields.date_value.date,custom_fields.date_value.date_time"

type TaskList struct {
	Data     []Task `json:"data"`
	NextPage *Page  `json:"next_page,omitempty"`
}

type ListTasksParams struct {
	Project        string
	Assignee       string
	Workspace      string
	Limit          int
	CompletedSince string
	OptFields      string
}

type SearchTasksParams struct {
	Text              string            `query:"text"`
	ResourceSubtype   string            `query:"resource_subtype"`
	AssigneeAny       string            `query:"assignee.any"`
	AssigneeNot       string            `query:"assignee.not"`
	PortfoliosAny     string            `query:"portfolios.any"`
	ProjectsAny       string            `query:"projects.any"`
	ProjectsNot       string            `query:"projects.not"`
	ProjectsAll       string            `query:"projects.all"`
	SectionsAny       string            `query:"sections.any"`
	SectionsNot       string            `query:"sections.not"`
	SectionsAll       string            `query:"sections.all"`
	TagsAny           string            `query:"tags.any"`
	TagsNot           string            `query:"tags.not"`
	TagsAll           string            `query:"tags.all"`
	TeamsAny          string            `query:"teams.any"`
	FollowersAny      string            `query:"followers.any"`
	FollowersNot      string            `query:"followers.not"`
	CreatedByAny      string            `query:"created_by.any"`
	CreatedByNot      string            `query:"created_by.not"`
	AssignedByAny     string            `query:"assigned_by.any"`
	AssignedByNot     string            `query:"assigned_by.not"`
	LikedByNot        string            `query:"liked_by.not"`
	CommentedOnByNot  string            `query:"commented_on_by.not"`
	DueOn             string            `query:"due_on"`
	DueOnBefore       string            `query:"due_on.before"`
	DueOnAfter        string            `query:"due_on.after"`
	DueAtBefore       string            `query:"due_at.before"`
	DueAtAfter        string            `query:"due_at.after"`
	StartOn           string            `query:"start_on"`
	StartOnBefore     string            `query:"start_on.before"`
	StartOnAfter      string            `query:"start_on.after"`
	CreatedOn         string            `query:"created_on"`
	CreatedOnBefore   string            `query:"created_on.before"`
	CreatedOnAfter    string            `query:"created_on.after"`
	CreatedAtBefore   string            `query:"created_at.before"`
	CreatedAtAfter    string            `query:"created_at.after"`
	CompletedOn       string            `query:"completed_on"`
	CompletedOnBefore string            `query:"completed_on.before"`
	CompletedOnAfter  string            `query:"completed_on.after"`
	CompletedAtBefore string            `query:"completed_at.before"`
	CompletedAtAfter  string            `query:"completed_at.after"`
	ModifiedOn        string            `query:"modified_on"`
	ModifiedOnBefore  string            `query:"modified_on.before"`
	ModifiedOnAfter   string            `query:"modified_on.after"`
	ModifiedAtBefore  string            `query:"modified_at.before"`
	ModifiedAtAfter   string            `query:"modified_at.after"`
	IsBlocking        *bool             `query:"is_blocking"`
	IsBlocked         *bool             `query:"is_blocked"`
	HasAttachment     *bool             `query:"has_attachment"`
	Completed         *bool             `query:"completed"`
	IsSubtask         *bool             `query:"is_subtask"`
	SortBy            string            `query:"sort_by"`
	SortAscending     *bool             `query:"sort_ascending"`
	Limit             int               `query:"limit"`
	OptFields         string            `query:"opt_fields"`
	Extra             map[string]string `query:"-"` // Merged separately
}

type CreateTaskRequest struct {
	Name      string   `json:"name"`
	Notes     string   `json:"notes,omitempty"`
	Assignee  string   `json:"assignee,omitempty"`
	Projects  []string `json:"projects,omitempty"`
	DueOn     string   `json:"due_on,omitempty"`
	Completed bool     `json:"completed,omitempty"`
}

type UpdateTaskRequest struct {
	Name      *string `json:"name,omitempty"`
	Notes     *string `json:"notes,omitempty"`
	Assignee  *string `json:"assignee,omitempty"`
	DueOn     *string `json:"due_on,omitempty"`
	Completed *bool   `json:"completed,omitempty"`
}

func (c *Client) ListTasks(ctx context.Context, params ListTasksParams) (*TaskList, error) {
	query := map[string]string{}
	if params.Project != "" {
		query["project"] = params.Project
	}
	if params.Assignee != "" {
		query["assignee"] = params.Assignee
	}
	if params.Workspace != "" {
		query["workspace"] = params.Workspace
	}
	if params.Limit > 0 {
		query["limit"] = fmt.Sprintf("%d", params.Limit)
	}
	if params.CompletedSince != "" {
		query["completed_since"] = params.CompletedSince
	}
	if params.OptFields != "" {
		query["opt_fields"] = params.OptFields
	}

	var resp responseEnvelope
	if err := c.doRaw(ctx, http.MethodGet, "/tasks", query, nil, &resp); err != nil {
		return nil, err
	}

	var list TaskList
	if err := json.Unmarshal(resp.Data, &list.Data); err != nil {
		return nil, err
	}
	list.NextPage = resp.NextPage
	return &list, nil
}

func (c *Client) SearchTasks(ctx context.Context, workspaceGID string, params SearchTasksParams) (*TaskList, error) {
	if params.OptFields == "" {
		params.OptFields = TaskDetailOptFields
	}

	query := buildQuery(params)

	var resp responseEnvelope
	if err := c.doRaw(ctx, http.MethodGet, "/workspaces/"+workspaceGID+"/tasks/search", query, nil, &resp); err != nil {
		return nil, err
	}

	var list TaskList
	if err := json.Unmarshal(resp.Data, &list.Data); err != nil {
		return nil, err
	}
	list.NextPage = resp.NextPage
	return &list, nil
}

func (c *Client) GetTask(ctx context.Context, gid string) (*Task, error) {
	query := map[string]string{
		"opt_fields": TaskDetailOptFields,
	}
	var task Task
	if err := c.do(ctx, http.MethodGet, "/tasks/"+gid, query, nil, &task); err != nil {
		return nil, err
	}
	return &task, nil
}

func (c *Client) CreateTask(ctx context.Context, req CreateTaskRequest) (*Task, error) {
	query := map[string]string{
		"opt_fields": "name,completed,assignee.name,assignee.gid,notes,due_on",
	}
	payload := map[string]CreateTaskRequest{"data": req}
	var task Task
	if err := c.do(ctx, http.MethodPost, "/tasks", query, payload, &task); err != nil {
		return nil, err
	}
	return &task, nil
}

func (c *Client) UpdateTask(ctx context.Context, gid string, req UpdateTaskRequest) (*Task, error) {
	query := map[string]string{
		"opt_fields": "name,completed,assignee.name,assignee.gid,notes,due_on",
	}
	payload := map[string]UpdateTaskRequest{"data": req}
	var task Task
	if err := c.do(ctx, http.MethodPut, "/tasks/"+gid, query, payload, &task); err != nil {
		return nil, err
	}
	return &task, nil
}
