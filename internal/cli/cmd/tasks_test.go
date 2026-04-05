package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/michalvavra/asncli/internal/asana"
	"github.com/michalvavra/asncli/internal/cli"
	"github.com/michalvavra/asncli/internal/config"
)

type tasksClient struct {
	list              *asana.TaskList
	task              *asana.Task
	createdTask       *asana.Task
	updatedTask       *asana.Task
	stories           *asana.StoryList
	subtasks          *asana.SubtaskList
	attachments       *asana.AttachmentList
	createdStory      *asana.Story
	updatedStory      *asana.Story
	gotParams         asana.ListTasksParams
	gotCreateReq      asana.CreateTaskRequest
	gotUpdateGID      string
	gotUpdateReq      asana.UpdateTaskRequest
	gotCreateStoryGID string
	gotCreateStoryReq asana.CreateStoryRequest
	gotUpdateStoryGID string
	gotUpdateStoryReq asana.UpdateStoryRequest
	gotDeleteStoryGID string
	attachment        *asana.Attachment
}

func (f *tasksClient) GetAttachment(ctx context.Context, gid string) (*asana.Attachment, error) {
	return f.attachment, nil
}

func (f *tasksClient) ListTasks(ctx context.Context, params asana.ListTasksParams) (*asana.TaskList, error) {
	f.gotParams = params
	return f.list, nil
}

func (f *tasksClient) GetTask(ctx context.Context, gid string) (*asana.Task, error) {
	return f.task, nil
}

func (f *tasksClient) GetTaskStories(ctx context.Context, taskGID string) (*asana.StoryList, error) {
	if f.stories != nil {
		return f.stories, nil
	}
	return &asana.StoryList{}, nil
}

func (f *tasksClient) GetSubtasks(ctx context.Context, taskGID string) (*asana.SubtaskList, error) {
	if f.subtasks != nil {
		return f.subtasks, nil
	}
	return &asana.SubtaskList{}, nil
}

func (f *tasksClient) GetTaskAttachments(ctx context.Context, taskGID string) (*asana.AttachmentList, error) {
	if f.attachments != nil {
		return f.attachments, nil
	}
	return &asana.AttachmentList{}, nil
}

func (f *tasksClient) CreateTask(ctx context.Context, req asana.CreateTaskRequest) (*asana.Task, error) {
	f.gotCreateReq = req
	return f.createdTask, nil
}

func (f *tasksClient) UpdateTask(ctx context.Context, gid string, req asana.UpdateTaskRequest) (*asana.Task, error) {
	f.gotUpdateGID = gid
	f.gotUpdateReq = req
	return f.updatedTask, nil
}

func (f *tasksClient) CreateStory(ctx context.Context, taskGID string, req asana.CreateStoryRequest) (*asana.Story, error) {
	f.gotCreateStoryGID = taskGID
	f.gotCreateStoryReq = req
	return f.createdStory, nil
}

func (f *tasksClient) UpdateStory(ctx context.Context, storyGID string, req asana.UpdateStoryRequest) (*asana.Story, error) {
	f.gotUpdateStoryGID = storyGID
	f.gotUpdateStoryReq = req
	return f.updatedStory, nil
}

func (f *tasksClient) DeleteStory(ctx context.Context, storyGID string) error {
	f.gotDeleteStoryGID = storyGID
	return nil
}

func TestTasksListJSON(t *testing.T) {
	buf := &bytes.Buffer{}
	ctx := &cli.Context{
		Stdout: buf,
		Stderr: &bytes.Buffer{},
		JSON:   true,
		Client: &tasksClient{list: &asana.TaskList{Data: []asana.Task{{GID: "1", Name: "Task"}}}},
	}

	cmd := TasksListCmd{}
	if err := cmd.Run(context.Background(), ctx); err != nil {
		t.Fatalf("TasksListCmd.Run returned unexpected error: %v", err)
	}

	var env cli.Envelope
	if err := json.Unmarshal(buf.Bytes(), &env); err != nil {
		t.Fatalf("failed to unmarshal JSON output: %v", err)
	}
	if env.Data == nil {
		t.Fatal("envelope data is nil, want task list")
	}
}

func TestTasksListHumanReadable(t *testing.T) {
	buf := &bytes.Buffer{}
	ctx := &cli.Context{
		Stdout: buf,
		Stderr: &bytes.Buffer{},
		Client: &tasksClient{list: &asana.TaskList{Data: []asana.Task{{
			GID:       "1",
			Name:      "Task",
			Completed: false,
			Assignee:  &asana.User{Name: "Alex"},
		}}}},
	}

	cmd := TasksListCmd{}
	if err := cmd.Run(context.Background(), ctx); err != nil {
		t.Fatalf("TasksListCmd.Run returned unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "GID") || !strings.Contains(out, "Task") || !strings.Contains(out, "Alex") {
		t.Errorf("output missing expected content\ngot: %q\nwant: contains GID, Task, and Alex", out)
	}
}

func TestTasksListWithAssigneeUsesDefaultWorkspace(t *testing.T) {
	buf := &bytes.Buffer{}
	client := &tasksClient{list: &asana.TaskList{Data: []asana.Task{{GID: "1", Name: "Task"}}}}
	ctx := &cli.Context{
		Stdout: buf,
		Stderr: &bytes.Buffer{},
		Client: client,
		Config: &config.Config{DefaultWorkspace: "default-ws"},
	}

	cmd := TasksListCmd{Assignee: "me"}
	if err := cmd.Run(context.Background(), ctx); err != nil {
		t.Fatalf("TasksListCmd.Run returned unexpected error: %v", err)
	}

	if got, want := client.gotParams.Workspace, "default-ws"; got != want {
		t.Errorf("workspace = %q, want %q", got, want)
	}
}

func TestTasksListWithAssigneeExplicitWorkspace(t *testing.T) {
	buf := &bytes.Buffer{}
	client := &tasksClient{list: &asana.TaskList{Data: []asana.Task{{GID: "1", Name: "Task"}}}}
	ctx := &cli.Context{
		Stdout: buf,
		Stderr: &bytes.Buffer{},
		Client: client,
		Config: &config.Config{DefaultWorkspace: "default-ws"},
	}

	cmd := TasksListCmd{Assignee: "me", Workspace: "explicit-ws"}
	if err := cmd.Run(context.Background(), ctx); err != nil {
		t.Fatalf("TasksListCmd.Run returned unexpected error: %v", err)
	}

	if got, want := client.gotParams.Workspace, "explicit-ws"; got != want {
		t.Errorf("workspace = %q, want %q", got, want)
	}
}

func TestTasksListWithAssigneeNoWorkspace(t *testing.T) {
	ctx := &cli.Context{
		Stdout: &bytes.Buffer{},
		Stderr: &bytes.Buffer{},
		Client: &tasksClient{list: &asana.TaskList{}},
		Config: &config.Config{},
	}

	cmd := TasksListCmd{Assignee: "me"}
	err := cmd.Run(context.Background(), ctx)
	if err == nil {
		t.Fatal("TasksListCmd.Run should return error when assignee is set but no workspace configured")
	}
	if !strings.Contains(err.Error(), "workspace is required") {
		t.Errorf("error = %q, want to contain 'workspace is required'", err.Error())
	}
}

func TestTasksGetHumanReadable(t *testing.T) {
	buf := &bytes.Buffer{}
	cfValue := 3.5
	ctx := &cli.Context{
		Stdout: buf,
		Stderr: &bytes.Buffer{},
		Client: &tasksClient{task: &asana.Task{
			GID:          "9",
			Name:         "Title",
			Completed:    true,
			Assignee:     &asana.User{Name: "Pat"},
			DueOn:        "2024-01-01",
			PermalinkURL: "https://app.asana.com/0/123/9",
			Notes:        "note",
			Projects:     []asana.Project{{Name: "Alpha"}},
			Tags:         []asana.Tag{{Name: "Urgent"}},
			Followers:    []asana.User{{Name: "Lee"}},
			Parent:       &asana.TaskCompact{Name: "Parent"},
			CustomFields: []asana.CustomField{{
				Name:         "Size",
				DisplayValue: "L",
			}, {
				Name:        "Estimate",
				NumberValue: &cfValue,
			}},
		}},
	}

	cmd := TasksGetCmd{GID: "9"}
	if err := cmd.Run(context.Background(), ctx); err != nil {
		t.Fatalf("TasksGetCmd.Run returned unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "GID") || !strings.Contains(out, "9") || !strings.Contains(out, "Title") || !strings.Contains(out, "Pat") {
		t.Errorf("output missing task details\ngot: %q\nwant: contains GID, 9, Title, and Pat", out)
	}
	if !strings.Contains(out, "Permalink") || !strings.Contains(out, "https://app.asana.com/0/123/9") {
		t.Errorf("output missing permalink\ngot: %q\nwant: contains Permalink and URL", out)
	}
	if !strings.Contains(out, "Projects") || !strings.Contains(out, "Alpha") || !strings.Contains(out, "Urgent") {
		t.Errorf("output missing projects/tags\ngot: %q\nwant: contains Projects, Alpha, and Urgent", out)
	}
}

func TestTasksGetJSON(t *testing.T) {
	buf := &bytes.Buffer{}
	ctx := &cli.Context{
		Stdout: buf,
		Stderr: &bytes.Buffer{},
		JSON:   true,
		Client: &tasksClient{task: &asana.Task{GID: "9", Name: "Title"}},
	}

	cmd := TasksGetCmd{GID: "9"}
	if err := cmd.Run(context.Background(), ctx); err != nil {
		t.Fatalf("TasksGetCmd.Run returned unexpected error: %v", err)
	}

	var env cli.Envelope
	if err := json.Unmarshal(buf.Bytes(), &env); err != nil {
		t.Fatalf("failed to unmarshal JSON output: %v", err)
	}
	if env.Data == nil {
		t.Fatal("envelope data is nil, want task")
	}
}

func TestTasksUpdateNoFields(t *testing.T) {
	cmd := TasksUpdateCmd{GID: "1"}
	ctx := &cli.Context{}
	if err := cmd.Run(context.Background(), ctx); err == nil {
		t.Fatal("TasksUpdateCmd.Run should return error when no fields specified")
	}
}

func TestTasksCreateHumanReadable(t *testing.T) {
	buf := &bytes.Buffer{}
	client := &tasksClient{createdTask: &asana.Task{GID: "new-task-123", Name: "New Task"}}
	ctx := &cli.Context{
		Stdout: buf,
		Stderr: &bytes.Buffer{},
		Client: client,
	}

	cmd := TasksCreateCmd{Name: "New Task", Notes: "Some notes", Assignee: "me", DueOn: "2024-12-31"}
	if err := cmd.Run(context.Background(), ctx); err != nil {
		t.Fatalf("TasksCreateCmd.Run returned unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "created") || !strings.Contains(out, "new-task-123") {
		t.Errorf("output = %q, want to contain 'created' and 'new-task-123'", out)
	}

	if client.gotCreateReq.Name != "New Task" {
		t.Errorf("request name = %q, want %q", client.gotCreateReq.Name, "New Task")
	}
	if client.gotCreateReq.Notes != "Some notes" {
		t.Errorf("request notes = %q, want %q", client.gotCreateReq.Notes, "Some notes")
	}
	if client.gotCreateReq.Assignee != "me" {
		t.Errorf("request assignee = %q, want %q", client.gotCreateReq.Assignee, "me")
	}
	if client.gotCreateReq.DueOn != "2024-12-31" {
		t.Errorf("request due_on = %q, want %q", client.gotCreateReq.DueOn, "2024-12-31")
	}
}

func TestTasksCreateJSON(t *testing.T) {
	buf := &bytes.Buffer{}
	client := &tasksClient{createdTask: &asana.Task{GID: "json-task", Name: "JSON Task"}}
	ctx := &cli.Context{
		Stdout: buf,
		Stderr: &bytes.Buffer{},
		JSON:   true,
		Client: client,
	}

	cmd := TasksCreateCmd{Name: "JSON Task"}
	if err := cmd.Run(context.Background(), ctx); err != nil {
		t.Fatalf("TasksCreateCmd.Run returned unexpected error: %v", err)
	}

	var env cli.Envelope
	if err := json.Unmarshal(buf.Bytes(), &env); err != nil {
		t.Fatalf("failed to unmarshal JSON output: %v", err)
	}
	if env.Data == nil {
		t.Fatal("envelope data is nil, want task")
	}
}

func TestTasksCreateWithProject(t *testing.T) {
	buf := &bytes.Buffer{}
	client := &tasksClient{createdTask: &asana.Task{GID: "proj-task"}}
	ctx := &cli.Context{
		Stdout: buf,
		Stderr: &bytes.Buffer{},
		Client: client,
	}

	cmd := TasksCreateCmd{Name: "Task with Project", Project: "proj-123"}
	if err := cmd.Run(context.Background(), ctx); err != nil {
		t.Fatalf("TasksCreateCmd.Run returned unexpected error: %v", err)
	}

	if len(client.gotCreateReq.Projects) != 1 || client.gotCreateReq.Projects[0] != "proj-123" {
		t.Errorf("request projects = %v, want [proj-123]", client.gotCreateReq.Projects)
	}
}

func TestTasksCreateTrimsWhitespace(t *testing.T) {
	buf := &bytes.Buffer{}
	client := &tasksClient{createdTask: &asana.Task{GID: "trimmed-task"}}
	ctx := &cli.Context{
		Stdout: buf,
		Stderr: &bytes.Buffer{},
		Client: client,
	}

	cmd := TasksCreateCmd{Name: "  Trimmed Name  ", Notes: "  Notes  ", Assignee: "  me  ", DueOn: "  2024-01-01  "}
	if err := cmd.Run(context.Background(), ctx); err != nil {
		t.Fatalf("TasksCreateCmd.Run returned unexpected error: %v", err)
	}

	if client.gotCreateReq.Name != "Trimmed Name" {
		t.Errorf("request name = %q, want %q", client.gotCreateReq.Name, "Trimmed Name")
	}
	if client.gotCreateReq.Notes != "Notes" {
		t.Errorf("request notes = %q, want %q", client.gotCreateReq.Notes, "Notes")
	}
	if client.gotCreateReq.Assignee != "me" {
		t.Errorf("request assignee = %q, want %q", client.gotCreateReq.Assignee, "me")
	}
	if client.gotCreateReq.DueOn != "2024-01-01" {
		t.Errorf("request due_on = %q, want %q", client.gotCreateReq.DueOn, "2024-01-01")
	}
}

func TestTasksCreateInvalidClient(t *testing.T) {
	ctx := &cli.Context{
		Stdout: &bytes.Buffer{},
		Stderr: &bytes.Buffer{},
		Client: struct{}{},
	}

	cmd := TasksCreateCmd{Name: "Test"}
	err := cmd.Run(context.Background(), ctx)
	if err == nil {
		t.Fatal("TasksCreateCmd.Run should return error for unsupported client type")
	}
	if !strings.Contains(err.Error(), "client does not support") {
		t.Errorf("error = %q, want to contain 'client does not support'", err.Error())
	}
}

func TestTasksUpdateHumanReadable(t *testing.T) {
	buf := &bytes.Buffer{}
	client := &tasksClient{updatedTask: &asana.Task{GID: "upd-123", Name: "Updated"}}
	ctx := &cli.Context{
		Stdout: buf,
		Stderr: &bytes.Buffer{},
		Client: client,
	}

	cmd := TasksUpdateCmd{GID: "upd-123", Name: "Updated Name"}
	if err := cmd.Run(context.Background(), ctx); err != nil {
		t.Fatalf("TasksUpdateCmd.Run returned unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "updated") || !strings.Contains(out, "upd-123") {
		t.Errorf("output = %q, want to contain 'updated' and 'upd-123'", out)
	}
}

func TestTasksUpdateJSON(t *testing.T) {
	buf := &bytes.Buffer{}
	client := &tasksClient{updatedTask: &asana.Task{GID: "json-upd", Name: "JSON Updated"}}
	ctx := &cli.Context{
		Stdout: buf,
		Stderr: &bytes.Buffer{},
		JSON:   true,
		Client: client,
	}

	cmd := TasksUpdateCmd{GID: "json-upd", Name: "Updated"}
	if err := cmd.Run(context.Background(), ctx); err != nil {
		t.Fatalf("TasksUpdateCmd.Run returned unexpected error: %v", err)
	}

	var env cli.Envelope
	if err := json.Unmarshal(buf.Bytes(), &env); err != nil {
		t.Fatalf("failed to unmarshal JSON output: %v", err)
	}
	if env.Data == nil {
		t.Fatal("envelope data is nil, want task")
	}
}

func TestTasksUpdateAllFields(t *testing.T) {
	buf := &bytes.Buffer{}
	client := &tasksClient{updatedTask: &asana.Task{GID: "all-fields"}}
	ctx := &cli.Context{
		Stdout: buf,
		Stderr: &bytes.Buffer{},
		Client: client,
	}

	completed := true
	cmd := TasksUpdateCmd{
		GID:       "all-fields",
		Name:      "New Name",
		Notes:     "New Notes",
		Assignee:  "user-123",
		DueOn:     "2025-01-01",
		Completed: &completed,
	}
	if err := cmd.Run(context.Background(), ctx); err != nil {
		t.Fatalf("TasksUpdateCmd.Run returned unexpected error: %v", err)
	}

	if client.gotUpdateGID != "all-fields" {
		t.Errorf("got GID = %q, want %q", client.gotUpdateGID, "all-fields")
	}
	if client.gotUpdateReq.Name == nil || *client.gotUpdateReq.Name != "New Name" {
		t.Errorf("request name = %v, want 'New Name'", client.gotUpdateReq.Name)
	}
	if client.gotUpdateReq.Notes == nil || *client.gotUpdateReq.Notes != "New Notes" {
		t.Errorf("request notes = %v, want 'New Notes'", client.gotUpdateReq.Notes)
	}
	if client.gotUpdateReq.Assignee == nil || *client.gotUpdateReq.Assignee != "user-123" {
		t.Errorf("request assignee = %v, want 'user-123'", client.gotUpdateReq.Assignee)
	}
	if client.gotUpdateReq.DueOn == nil || *client.gotUpdateReq.DueOn != "2025-01-01" {
		t.Errorf("request due_on = %v, want '2025-01-01'", client.gotUpdateReq.DueOn)
	}
	if client.gotUpdateReq.Completed == nil || *client.gotUpdateReq.Completed != true {
		t.Errorf("request completed = %v, want true", client.gotUpdateReq.Completed)
	}
}

func TestTasksUpdateInvalidClient(t *testing.T) {
	ctx := &cli.Context{
		Stdout: &bytes.Buffer{},
		Stderr: &bytes.Buffer{},
		Client: struct{}{},
	}

	cmd := TasksUpdateCmd{GID: "1", Name: "Test"}
	err := cmd.Run(context.Background(), ctx)
	if err == nil {
		t.Fatal("TasksUpdateCmd.Run should return error for unsupported client type")
	}
	if !strings.Contains(err.Error(), "client does not support") {
		t.Errorf("error = %q, want to contain 'client does not support'", err.Error())
	}
}

func TestTasksListInvalidClient(t *testing.T) {
	ctx := &cli.Context{
		Stdout: &bytes.Buffer{},
		Stderr: &bytes.Buffer{},
		Client: struct{}{},
	}

	cmd := TasksListCmd{}
	err := cmd.Run(context.Background(), ctx)
	if err == nil {
		t.Fatal("TasksListCmd.Run should return error for unsupported client type")
	}
	if !strings.Contains(err.Error(), "client does not support") {
		t.Errorf("error = %q, want to contain 'client does not support'", err.Error())
	}
}

func TestTasksGetInvalidClient(t *testing.T) {
	ctx := &cli.Context{
		Stdout: &bytes.Buffer{},
		Stderr: &bytes.Buffer{},
		Client: struct{}{},
	}

	cmd := TasksGetCmd{GID: "123"}
	err := cmd.Run(context.Background(), ctx)
	if err == nil {
		t.Fatal("TasksGetCmd.Run should return error for unsupported client type")
	}
	if !strings.Contains(err.Error(), "client does not support") {
		t.Errorf("error = %q, want to contain 'client does not support'", err.Error())
	}
}

func TestTasksGetNoAssignee(t *testing.T) {
	buf := &bytes.Buffer{}
	ctx := &cli.Context{
		Stdout: buf,
		Stderr: &bytes.Buffer{},
		Client: &tasksClient{task: &asana.Task{
			GID:       "no-assignee",
			Name:      "Unassigned Task",
			Completed: false,
		}},
	}

	cmd := TasksGetCmd{GID: "no-assignee"}
	if err := cmd.Run(context.Background(), ctx); err != nil {
		t.Fatalf("TasksGetCmd.Run returned unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "Unassigned Task") {
		t.Errorf("output = %q, want to contain 'Unassigned Task'", out)
	}
}

func TestTasksListNoAssignee(t *testing.T) {
	buf := &bytes.Buffer{}
	ctx := &cli.Context{
		Stdout: buf,
		Stderr: &bytes.Buffer{},
		Client: &tasksClient{list: &asana.TaskList{Data: []asana.Task{{
			GID:       "1",
			Name:      "Task Without Assignee",
			Completed: false,
		}}}},
	}

	cmd := TasksListCmd{}
	if err := cmd.Run(context.Background(), ctx); err != nil {
		t.Fatalf("TasksListCmd.Run returned unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "Task Without Assignee") {
		t.Errorf("output = %q, want to contain 'Task Without Assignee'", out)
	}
}

type tasksSearchClientStub struct {
	list         *asana.TaskList
	gotWorkspace string
	gotParams    asana.SearchTasksParams
}

func (c *tasksSearchClientStub) SearchTasks(ctx context.Context, workspaceGID string, params asana.SearchTasksParams) (*asana.TaskList, error) {
	c.gotWorkspace = workspaceGID
	c.gotParams = params
	return c.list, nil
}

func TestTasksSearchHumanReadable(t *testing.T) {
	buf := &bytes.Buffer{}
	client := &tasksSearchClientStub{list: &asana.TaskList{Data: []asana.Task{{
		GID:       "search-1",
		Name:      "Found Task",
		Completed: false,
		Assignee:  &asana.User{Name: "Finder"},
	}}}}
	ctx := &cli.Context{
		Stdout: buf,
		Stderr: &bytes.Buffer{},
		Client: client,
		Config: &config.Config{DefaultWorkspace: "ws-search"},
	}

	cmd := TasksSearchCmd{Text: "test query"}
	if err := cmd.Run(context.Background(), ctx); err != nil {
		t.Fatalf("TasksSearchCmd.Run returned unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "Found Task") || !strings.Contains(out, "Finder") {
		t.Errorf("output = %q, want to contain 'Found Task' and 'Finder'", out)
	}

	if client.gotWorkspace != "ws-search" {
		t.Errorf("workspace = %q, want %q", client.gotWorkspace, "ws-search")
	}
	if client.gotParams.Text != "test query" {
		t.Errorf("text param = %q, want %q", client.gotParams.Text, "test query")
	}
}

func TestTasksSearchJSON(t *testing.T) {
	buf := &bytes.Buffer{}
	client := &tasksSearchClientStub{list: &asana.TaskList{Data: []asana.Task{{GID: "json-search", Name: "JSON Found"}}}}
	ctx := &cli.Context{
		Stdout: buf,
		Stderr: &bytes.Buffer{},
		JSON:   true,
		Client: client,
		Config: &config.Config{DefaultWorkspace: "ws-json"},
	}

	cmd := TasksSearchCmd{Text: "json search"}
	if err := cmd.Run(context.Background(), ctx); err != nil {
		t.Fatalf("TasksSearchCmd.Run returned unexpected error: %v", err)
	}

	var env cli.Envelope
	if err := json.Unmarshal(buf.Bytes(), &env); err != nil {
		t.Fatalf("failed to unmarshal JSON output: %v", err)
	}
	if env.Data == nil {
		t.Fatal("envelope data is nil, want task list")
	}
}

func TestTasksSearchNoWorkspace(t *testing.T) {
	ctx := &cli.Context{
		Stdout: &bytes.Buffer{},
		Stderr: &bytes.Buffer{},
		Client: &tasksSearchClientStub{list: &asana.TaskList{}},
		Config: &config.Config{},
	}

	cmd := TasksSearchCmd{Text: "query"}
	err := cmd.Run(context.Background(), ctx)
	if err == nil {
		t.Fatal("TasksSearchCmd.Run should return error when no workspace configured")
	}
	if !strings.Contains(err.Error(), "workspace is required") {
		t.Errorf("error = %q, want to contain 'workspace is required'", err.Error())
	}
}

func TestTasksSearchInvalidClient(t *testing.T) {
	ctx := &cli.Context{
		Stdout: &bytes.Buffer{},
		Stderr: &bytes.Buffer{},
		Client: struct{}{},
		Config: &config.Config{DefaultWorkspace: "ws"},
	}

	cmd := TasksSearchCmd{Text: "query"}
	err := cmd.Run(context.Background(), ctx)
	if err == nil {
		t.Fatal("TasksSearchCmd.Run should return error for unsupported client type")
	}
	if !strings.Contains(err.Error(), "client does not support") {
		t.Errorf("error = %q, want to contain 'client does not support'", err.Error())
	}
}

func TestTasksSearchInvalidFilter(t *testing.T) {
	ctx := &cli.Context{
		Stdout: &bytes.Buffer{},
		Stderr: &bytes.Buffer{},
		Client: &tasksSearchClientStub{list: &asana.TaskList{}},
		Config: &config.Config{DefaultWorkspace: "ws"},
	}

	cmd := TasksSearchCmd{Filter: []string{"invalid-no-equals"}}
	err := cmd.Run(context.Background(), ctx)
	if err == nil {
		t.Fatal("TasksSearchCmd.Run should return error for invalid filter format")
	}
	if !strings.Contains(err.Error(), "failed to parse filter") {
		t.Errorf("error = %q, want to contain 'failed to parse filter'", err.Error())
	}
}

func TestTasksSearchWithFilters(t *testing.T) {
	buf := &bytes.Buffer{}
	client := &tasksSearchClientStub{list: &asana.TaskList{Data: []asana.Task{}}}
	ctx := &cli.Context{
		Stdout: buf,
		Stderr: &bytes.Buffer{},
		Client: client,
		Config: &config.Config{DefaultWorkspace: "ws"},
	}

	completed := true
	isBlocking := true
	cmd := TasksSearchCmd{
		Text:        "query",
		AssigneeAny: []string{"user1", "user2"},
		ProjectsAny: []string{"proj1"},
		Completed:   &completed,
		IsBlocking:  &isBlocking,
		SortBy:      "due_date",
		Filter:      []string{"custom=value"},
	}
	if err := cmd.Run(context.Background(), ctx); err != nil {
		t.Fatalf("TasksSearchCmd.Run returned unexpected error: %v", err)
	}

	if client.gotParams.AssigneeAny != "user1,user2" {
		t.Errorf("assignee_any = %q, want %q", client.gotParams.AssigneeAny, "user1,user2")
	}
	if client.gotParams.ProjectsAny != "proj1" {
		t.Errorf("projects_any = %q, want %q", client.gotParams.ProjectsAny, "proj1")
	}
	if client.gotParams.Completed == nil || *client.gotParams.Completed != true {
		t.Errorf("completed = %v, want true", client.gotParams.Completed)
	}
	if client.gotParams.IsBlocking == nil || *client.gotParams.IsBlocking != true {
		t.Errorf("is_blocking = %v, want true", client.gotParams.IsBlocking)
	}
	if client.gotParams.SortBy != "due_date" {
		t.Errorf("sort_by = %q, want %q", client.gotParams.SortBy, "due_date")
	}
	if client.gotParams.Extra["custom"] != "value" {
		t.Errorf("extra[custom] = %q, want %q", client.gotParams.Extra["custom"], "value")
	}
}

func TestFormatCustomFieldsEmpty(t *testing.T) {
	result := formatCustomFields(nil)
	if result != "none" {
		t.Errorf("formatCustomFields(nil) = %q, want %q", result, "none")
	}

	result = formatCustomFields([]asana.CustomField{})
	if result != "none" {
		t.Errorf("formatCustomFields([]) = %q, want %q", result, "none")
	}
}

func TestFormatCustomFieldsVariousTypes(t *testing.T) {
	numVal := 42.5
	fields := []asana.CustomField{
		{Name: "Display", DisplayValue: "displayed"},
		{Name: "Enum", EnumValue: &asana.CustomFieldEnum{Name: "Option A"}},
		{Name: "Text", TextValue: "text value"},
		{Name: "Number", NumberValue: &numVal},
		{Name: "Empty"},
	}

	result := formatCustomFields(fields)

	if !strings.Contains(result, "Display=displayed") {
		t.Errorf("result = %q, want to contain 'Display=displayed'", result)
	}
	if !strings.Contains(result, "Enum=Option A") {
		t.Errorf("result = %q, want to contain 'Enum=Option A'", result)
	}
	if !strings.Contains(result, "Text=text value") {
		t.Errorf("result = %q, want to contain 'Text=text value'", result)
	}
	if !strings.Contains(result, "Number=42.5") {
		t.Errorf("result = %q, want to contain 'Number=42.5'", result)
	}
	if !strings.Contains(result, "Empty=empty") {
		t.Errorf("result = %q, want to contain 'Empty=empty'", result)
	}
}

func TestFormatProjectNamesEmpty(t *testing.T) {
	result := formatProjectNames(nil)
	if result != "none" {
		t.Errorf("formatProjectNames(nil) = %q, want %q", result, "none")
	}

	result = formatProjectNames([]asana.Project{})
	if result != "none" {
		t.Errorf("formatProjectNames([]) = %q, want %q", result, "none")
	}

	result = formatProjectNames([]asana.Project{{GID: "1"}})
	if result != "none" {
		t.Errorf("formatProjectNames([{GID:1}]) = %q, want %q", result, "none")
	}
}

func TestFormatTagNamesEmpty(t *testing.T) {
	result := formatTagNames(nil)
	if result != "none" {
		t.Errorf("formatTagNames(nil) = %q, want %q", result, "none")
	}

	result = formatTagNames([]asana.Tag{})
	if result != "none" {
		t.Errorf("formatTagNames([]) = %q, want %q", result, "none")
	}

	result = formatTagNames([]asana.Tag{{GID: "1"}})
	if result != "none" {
		t.Errorf("formatTagNames([{GID:1}]) = %q, want %q", result, "none")
	}
}

func TestFormatUserNamesEmpty(t *testing.T) {
	result := formatUserNames(nil)
	if result != "none" {
		t.Errorf("formatUserNames(nil) = %q, want %q", result, "none")
	}

	result = formatUserNames([]asana.User{})
	if result != "none" {
		t.Errorf("formatUserNames([]) = %q, want %q", result, "none")
	}

	result = formatUserNames([]asana.User{{GID: "1"}})
	if result != "none" {
		t.Errorf("formatUserNames([{GID:1}]) = %q, want %q", result, "none")
	}
}

func TestFormatMembershipSectionsEmpty(t *testing.T) {
	result := formatMembershipSections(nil)
	if result != "none" {
		t.Errorf("formatMembershipSections(nil) = %q, want %q", result, "none")
	}

	result = formatMembershipSections([]asana.Membership{})
	if result != "none" {
		t.Errorf("formatMembershipSections([]) = %q, want %q", result, "none")
	}

	result = formatMembershipSections([]asana.Membership{{Section: asana.Section{GID: "1"}}})
	if result != "none" {
		t.Errorf("formatMembershipSections([{Section:{GID:1}}]) = %q, want %q", result, "none")
	}
}

func TestTasksCommentAddHumanReadable(t *testing.T) {
	buf := &bytes.Buffer{}
	client := &tasksClient{createdStory: &asana.Story{GID: "story-123", Text: "hello"}}
	ctx := &cli.Context{
		Stdout: buf,
		Stderr: &bytes.Buffer{},
		Client: client,
	}

	cmd := TasksCommentAddCmd{GID: "task-1", Text: "hello"}
	if err := cmd.Run(context.Background(), ctx); err != nil {
		t.Fatalf("TasksCommentAddCmd.Run returned unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "commented") || !strings.Contains(out, "story-123") {
		t.Errorf("output = %q, want to contain 'commented' and 'story-123'", out)
	}
	if client.gotCreateStoryGID != "task-1" {
		t.Errorf("task GID = %q, want %q", client.gotCreateStoryGID, "task-1")
	}
	if client.gotCreateStoryReq.Text != "hello" {
		t.Errorf("text = %q, want %q", client.gotCreateStoryReq.Text, "hello")
	}
}

func TestTasksCommentAddJSON(t *testing.T) {
	buf := &bytes.Buffer{}
	client := &tasksClient{createdStory: &asana.Story{GID: "story-json", Text: "test"}}
	ctx := &cli.Context{
		Stdout: buf,
		Stderr: &bytes.Buffer{},
		JSON:   true,
		Client: client,
	}

	cmd := TasksCommentAddCmd{GID: "task-1", Text: "test"}
	if err := cmd.Run(context.Background(), ctx); err != nil {
		t.Fatalf("TasksCommentAddCmd.Run returned unexpected error: %v", err)
	}

	var env cli.Envelope
	if err := json.Unmarshal(buf.Bytes(), &env); err != nil {
		t.Fatalf("failed to unmarshal JSON output: %v", err)
	}
	if env.Data == nil {
		t.Fatal("envelope data is nil, want story")
	}
}

func TestTasksCommentAddInvalidClient(t *testing.T) {
	ctx := &cli.Context{
		Stdout: &bytes.Buffer{},
		Stderr: &bytes.Buffer{},
		Client: struct{}{},
	}

	cmd := TasksCommentAddCmd{GID: "task-1", Text: "test"}
	err := cmd.Run(context.Background(), ctx)
	if err == nil {
		t.Fatal("TasksCommentAddCmd.Run should return error for unsupported client type")
	}
	if !strings.Contains(err.Error(), "client does not support") {
		t.Errorf("error = %q, want to contain 'client does not support'", err.Error())
	}
}

func TestTasksCommentUpdateHumanReadable(t *testing.T) {
	buf := &bytes.Buffer{}
	client := &tasksClient{updatedStory: &asana.Story{GID: "story-456", Text: "updated"}}
	ctx := &cli.Context{
		Stdout: buf,
		Stderr: &bytes.Buffer{},
		Client: client,
	}

	cmd := TasksCommentUpdateCmd{GID: "story-456", Text: "updated"}
	if err := cmd.Run(context.Background(), ctx); err != nil {
		t.Fatalf("TasksCommentUpdateCmd.Run returned unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "updated") || !strings.Contains(out, "story-456") {
		t.Errorf("output = %q, want to contain 'updated' and 'story-456'", out)
	}
	if client.gotUpdateStoryGID != "story-456" {
		t.Errorf("story GID = %q, want %q", client.gotUpdateStoryGID, "story-456")
	}
	if client.gotUpdateStoryReq.Text != "updated" {
		t.Errorf("text = %q, want %q", client.gotUpdateStoryReq.Text, "updated")
	}
}

func TestTasksCommentUpdateJSON(t *testing.T) {
	buf := &bytes.Buffer{}
	client := &tasksClient{updatedStory: &asana.Story{GID: "story-json-upd", Text: "new text"}}
	ctx := &cli.Context{
		Stdout: buf,
		Stderr: &bytes.Buffer{},
		JSON:   true,
		Client: client,
	}

	cmd := TasksCommentUpdateCmd{GID: "story-json-upd", Text: "new text"}
	if err := cmd.Run(context.Background(), ctx); err != nil {
		t.Fatalf("TasksCommentUpdateCmd.Run returned unexpected error: %v", err)
	}

	var env cli.Envelope
	if err := json.Unmarshal(buf.Bytes(), &env); err != nil {
		t.Fatalf("failed to unmarshal JSON output: %v", err)
	}
	if env.Data == nil {
		t.Fatal("envelope data is nil, want story")
	}
}

func TestTasksCommentUpdateInvalidClient(t *testing.T) {
	ctx := &cli.Context{
		Stdout: &bytes.Buffer{},
		Stderr: &bytes.Buffer{},
		Client: struct{}{},
	}

	cmd := TasksCommentUpdateCmd{GID: "story-1", Text: "test"}
	err := cmd.Run(context.Background(), ctx)
	if err == nil {
		t.Fatal("TasksCommentUpdateCmd.Run should return error for unsupported client type")
	}
	if !strings.Contains(err.Error(), "client does not support") {
		t.Errorf("error = %q, want to contain 'client does not support'", err.Error())
	}
}

func TestTasksCommentDeleteHumanReadable(t *testing.T) {
	buf := &bytes.Buffer{}
	client := &tasksClient{}
	ctx := &cli.Context{
		Stdout: buf,
		Stderr: &bytes.Buffer{},
		Client: client,
	}

	cmd := TasksCommentDeleteCmd{GID: "story-789"}
	if err := cmd.Run(context.Background(), ctx); err != nil {
		t.Fatalf("TasksCommentDeleteCmd.Run returned unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "deleted") || !strings.Contains(out, "story-789") {
		t.Errorf("output = %q, want to contain 'deleted' and 'story-789'", out)
	}
	if client.gotDeleteStoryGID != "story-789" {
		t.Errorf("story GID = %q, want %q", client.gotDeleteStoryGID, "story-789")
	}
}

func TestTasksCommentDeleteJSON(t *testing.T) {
	buf := &bytes.Buffer{}
	client := &tasksClient{}
	ctx := &cli.Context{
		Stdout: buf,
		Stderr: &bytes.Buffer{},
		JSON:   true,
		Client: client,
	}

	cmd := TasksCommentDeleteCmd{GID: "story-del-json"}
	if err := cmd.Run(context.Background(), ctx); err != nil {
		t.Fatalf("TasksCommentDeleteCmd.Run returned unexpected error: %v", err)
	}

	var env cli.Envelope
	if err := json.Unmarshal(buf.Bytes(), &env); err != nil {
		t.Fatalf("failed to unmarshal JSON output: %v", err)
	}
	if env.Data == nil {
		t.Fatal("envelope data is nil, want deletion response")
	}
}

func TestTasksCommentDeleteInvalidClient(t *testing.T) {
	ctx := &cli.Context{
		Stdout: &bytes.Buffer{},
		Stderr: &bytes.Buffer{},
		Client: struct{}{},
	}

	cmd := TasksCommentDeleteCmd{GID: "story-1"}
	err := cmd.Run(context.Background(), ctx)
	if err == nil {
		t.Fatal("TasksCommentDeleteCmd.Run should return error for unsupported client type")
	}
	if !strings.Contains(err.Error(), "client does not support") {
		t.Errorf("error = %q, want to contain 'client does not support'", err.Error())
	}
}

func TestTasksCommentAddTrimsWhitespace(t *testing.T) {
	buf := &bytes.Buffer{}
	client := &tasksClient{createdStory: &asana.Story{GID: "trimmed-story"}}
	ctx := &cli.Context{
		Stdout: buf,
		Stderr: &bytes.Buffer{},
		Client: client,
	}

	cmd := TasksCommentAddCmd{GID: "task-1", Text: "  hello world  "}
	if err := cmd.Run(context.Background(), ctx); err != nil {
		t.Fatalf("TasksCommentAddCmd.Run returned unexpected error: %v", err)
	}

	if client.gotCreateStoryReq.Text != "hello world" {
		t.Errorf("text = %q, want %q", client.gotCreateStoryReq.Text, "hello world")
	}
}

func TestTasksAttachmentGetHumanReadable(t *testing.T) {
	buf := &bytes.Buffer{}
	client := &tasksClient{attachment: &asana.Attachment{
		GID:         "att-1",
		Name:        "file.png",
		Host:        "asana",
		CreatedAt:   "2024-01-01T00:00:00Z",
		DownloadURL: "https://example.com/file.png",
		ViewURL:     "https://example.com/view",
		Parent:      &asana.TaskCompact{GID: "task-1", Name: "Parent Task"},
	}}
	ctx := &cli.Context{
		Stdout: buf,
		Stderr: &bytes.Buffer{},
		Client: client,
	}

	cmd := TasksAttachmentGetCmd{GID: "att-1"}
	if err := cmd.Run(context.Background(), ctx); err != nil {
		t.Fatalf("TasksAttachmentGetCmd.Run returned unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "att-1") || !strings.Contains(out, "file.png") || !strings.Contains(out, "Parent Task") {
		t.Errorf("output missing expected content\ngot: %q", out)
	}
}

func TestTasksAttachmentGetJSON(t *testing.T) {
	buf := &bytes.Buffer{}
	client := &tasksClient{attachment: &asana.Attachment{GID: "att-json", Name: "doc.pdf"}}
	ctx := &cli.Context{
		Stdout: buf,
		Stderr: &bytes.Buffer{},
		JSON:   true,
		Client: client,
	}

	cmd := TasksAttachmentGetCmd{GID: "att-json"}
	if err := cmd.Run(context.Background(), ctx); err != nil {
		t.Fatalf("TasksAttachmentGetCmd.Run returned unexpected error: %v", err)
	}

	var env cli.Envelope
	if err := json.Unmarshal(buf.Bytes(), &env); err != nil {
		t.Fatalf("failed to unmarshal JSON output: %v", err)
	}
	if env.Data == nil {
		t.Fatal("envelope data is nil, want attachment")
	}
}

func TestTasksAttachmentGetInvalidClient(t *testing.T) {
	ctx := &cli.Context{
		Stdout: &bytes.Buffer{},
		Stderr: &bytes.Buffer{},
		Client: struct{}{},
	}

	cmd := TasksAttachmentGetCmd{GID: "att-1"}
	err := cmd.Run(context.Background(), ctx)
	if err == nil {
		t.Fatal("TasksAttachmentGetCmd.Run should return error for unsupported client type")
	}
	if !strings.Contains(err.Error(), "client does not support") {
		t.Errorf("error = %q, want to contain 'client does not support'", err.Error())
	}
}
