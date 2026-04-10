package cmd

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/michalvavra/asncli/internal/asana"
	"github.com/michalvavra/asncli/internal/cli"
)

type TasksCmd struct {
	List            TasksListCmd            `cmd:"" help:"List tasks."`
	Search          TasksSearchCmd          `cmd:"" help:"Search tasks."`
	Get             TasksGetCmd             `cmd:"" help:"Get a task."`
	Create          TasksCreateCmd          `cmd:"" help:"Create a task."`
	Update          TasksUpdateCmd          `cmd:"" help:"Update a task."`
	Comments        TasksCommentsCmd        `cmd:"" help:"List comments on a task."`
	CommentAdd      TasksCommentAddCmd      `cmd:"comment-add" help:"Add a comment to a task."`
	CommentUpdate   TasksCommentUpdateCmd   `cmd:"comment-update" help:"Update a comment."`
	CommentDelete   TasksCommentDeleteCmd   `cmd:"comment-delete" help:"Delete a comment."`
	Subtasks        TasksSubtasksCmd        `cmd:"" help:"List subtasks of a task."`
	Attachments     TasksAttachmentsCmd     `cmd:"" help:"List attachments on a task."`
	AttachmentGet   TasksAttachmentGetCmd   `cmd:"attachment-get" help:"Get an attachment."`
	Stories         TasksStoriesCmd         `cmd:"" help:"List all stories (activity log) on a task."`
	AddFollower     TasksAddFollowerCmd     `cmd:"add-follower" help:"Add followers to a task."`
	RemoveFollower  TasksRemoveFollowerCmd  `cmd:"remove-follower" help:"Remove followers from a task."`
}

type TasksListCmd struct {
	Project        string `help:"Project GID."`
	Assignee       string `help:"Assignee GID or 'me'."`
	Workspace      string `help:"Workspace GID (uses default if not specified when assignee is set)."`
	Limit          int    `help:"Maximum number of tasks."`
	CompletedSince string `help:"Only tasks completed since this time (RFC3339)."`
}

type TasksSearchCmd struct {
	Workspace         string   `help:"Workspace GID (uses default if not specified)."`
	Text              string   `help:"Full-text search on name and description."`
	ResourceSubtype   string   `help:"Filter by resource subtype."`
	AssigneeAny       []string `sep:"," help:"Assignee GIDs."`
	AssigneeNot       []string `sep:"," help:"Exclude assignee GIDs."`
	PortfoliosAny     []string `sep:"," help:"Portfolio GIDs."`
	ProjectsAny       []string `sep:"," help:"Project GIDs."`
	ProjectsNot       []string `sep:"," help:"Exclude project GIDs."`
	ProjectsAll       []string `sep:"," help:"Project GIDs required."`
	SectionsAny       []string `sep:"," help:"Section GIDs."`
	SectionsNot       []string `sep:"," help:"Exclude section GIDs."`
	SectionsAll       []string `sep:"," help:"Section GIDs required."`
	TagsAny           []string `sep:"," help:"Tag GIDs."`
	TagsNot           []string `sep:"," help:"Exclude tag GIDs."`
	TagsAll           []string `sep:"," help:"Tag GIDs required."`
	TeamsAny          []string `sep:"," help:"Team GIDs."`
	FollowersAny      []string `sep:"," help:"Follower GIDs."`
	FollowersNot      []string `sep:"," help:"Exclude follower GIDs."`
	CreatedByAny      []string `sep:"," help:"Creator GIDs."`
	CreatedByNot      []string `sep:"," help:"Exclude creator GIDs."`
	AssignedByAny     []string `sep:"," help:"Assigned by GIDs."`
	AssignedByNot     []string `sep:"," help:"Exclude assigned by GIDs."`
	LikedByNot        []string `sep:"," help:"Exclude liked by GIDs."`
	CommentedOnByNot  []string `sep:"," help:"Exclude commented by GIDs."`
	DueOn             string   `help:"Due on date (YYYY-MM-DD or 'null')."`
	DueOnBefore       string   `help:"Due on before date (YYYY-MM-DD)."`
	DueOnAfter        string   `help:"Due on after date (YYYY-MM-DD)."`
	DueAtBefore       string   `help:"Due at before (RFC3339)."`
	DueAtAfter        string   `help:"Due at after (RFC3339)."`
	StartOn           string   `help:"Start on date (YYYY-MM-DD or 'null')."`
	StartOnBefore     string   `help:"Start on before date (YYYY-MM-DD)."`
	StartOnAfter      string   `help:"Start on after date (YYYY-MM-DD)."`
	CreatedOn         string   `help:"Created on date (YYYY-MM-DD or 'null')."`
	CreatedOnBefore   string   `help:"Created on before date (YYYY-MM-DD)."`
	CreatedOnAfter    string   `help:"Created on after date (YYYY-MM-DD)."`
	CreatedAtBefore   string   `help:"Created at before (RFC3339)."`
	CreatedAtAfter    string   `help:"Created at after (RFC3339)."`
	CompletedOn       string   `help:"Completed on date (YYYY-MM-DD or 'null')."`
	CompletedOnBefore string   `help:"Completed on before date (YYYY-MM-DD)."`
	CompletedOnAfter  string   `help:"Completed on after date (YYYY-MM-DD)."`
	CompletedAtBefore string   `help:"Completed at before (RFC3339)."`
	CompletedAtAfter  string   `help:"Completed at after (RFC3339)."`
	ModifiedOn        string   `help:"Modified on date (YYYY-MM-DD or 'null')."`
	ModifiedOnBefore  string   `help:"Modified on before date (YYYY-MM-DD)."`
	ModifiedOnAfter   string   `help:"Modified on after date (YYYY-MM-DD)."`
	ModifiedAtBefore  string   `help:"Modified at before (RFC3339)."`
	ModifiedAtAfter   string   `help:"Modified at after (RFC3339)."`
	IsBlocking        *bool    `help:"Only tasks that block other tasks."`
	IsBlocked         *bool    `help:"Only tasks blocked by others."`
	HasAttachment     *bool    `help:"Only tasks with attachments."`
	Completed         *bool    `help:"Filter by completed status."`
	IsSubtask         *bool    `help:"Only subtasks."`
	SortBy            string   `help:"Sort by due_date, created_at, completed_at, likes, modified_at."`
	SortAscending     *bool    `help:"Sort ascending."`
	Limit             int      `help:"Maximum number of tasks."`
	Filter            []string `help:"Extra filters as key=value."`
}

type TasksGetCmd struct {
	GID string `arg:"" help:"Task GID."`
}

type TasksCommentsCmd struct {
	GID string `arg:"" help:"Task GID."`
}

type TasksSubtasksCmd struct {
	GID string `arg:"" help:"Task GID."`
}

type TasksAttachmentsCmd struct {
	GID string `arg:"" help:"Task GID."`
}

type TasksAttachmentGetCmd struct {
	GID string `arg:"" help:"Attachment GID."`
}

type TasksCommentAddCmd struct {
	GID  string `arg:"" help:"Task GID."`
	Text string `help:"Comment text." required:""`
}

type TasksCommentUpdateCmd struct {
	GID  string `arg:"" help:"Comment (story) GID."`
	Text string `help:"New comment text." required:""`
}

type TasksCommentDeleteCmd struct {
	GID string `arg:"" help:"Comment (story) GID."`
}

type TasksStoriesCmd struct {
	GID string `arg:"" help:"Task GID."`
}

type TasksAddFollowerCmd struct {
	GID       string   `arg:"" help:"Task GID."`
	Followers []string `arg:"" help:"Follower GIDs to add."`
}

type TasksRemoveFollowerCmd struct {
	GID       string   `arg:"" help:"Task GID."`
	Followers []string `arg:"" help:"Follower GIDs to remove."`
}

type TasksCreateCmd struct {
	Name         string   `help:"Task name." required:""`
	Notes        string   `help:"Task notes."`
	Assignee     string   `help:"Assignee GID or 'me'."`
	Project      string   `help:"Project GID."`
	Parent       string   `help:"Parent task GID (creates a subtask)."`
	Section      string   `help:"Section GID (adds task to section after creation)."`
	DueOn        string   `help:"Due date (YYYY-MM-DD)."`
	CustomFields []string `sep:"," help:"Custom field values as GID=value pairs."`
}

type TasksUpdateCmd struct {
	GID       string `arg:"" help:"Task GID."`
	Name      string `help:"Task name."`
	Notes     string `help:"Task notes."`
	Assignee  string `help:"Assignee GID or 'me'."`
	DueOn     string `help:"Due date (YYYY-MM-DD)."`
	Completed *bool  `help:"Set completion status."`
}

type tasksListClient interface {
	ListTasks(ctx context.Context, params asana.ListTasksParams) (*asana.TaskList, error)
}

type tasksSearchClient interface {
	SearchTasks(ctx context.Context, workspaceGID string, params asana.SearchTasksParams) (*asana.TaskList, error)
}

type tasksGetClient interface {
	GetTask(ctx context.Context, gid string) (*asana.Task, error)
	GetTaskStories(ctx context.Context, taskGID string) (*asana.StoryList, error)
	GetSubtasks(ctx context.Context, taskGID string) (*asana.SubtaskList, error)
	GetTaskAttachments(ctx context.Context, taskGID string) (*asana.AttachmentList, error)
}

type tasksCommentsClient interface {
	GetTaskStories(ctx context.Context, taskGID string) (*asana.StoryList, error)
}

type tasksSubtasksClient interface {
	GetSubtasks(ctx context.Context, taskGID string) (*asana.SubtaskList, error)
}

type tasksAttachmentsClient interface {
	GetTaskAttachments(ctx context.Context, taskGID string) (*asana.AttachmentList, error)
}

type tasksAttachmentGetClient interface {
	GetAttachment(ctx context.Context, gid string) (*asana.Attachment, error)
}

type tasksCommentAddClient interface {
	CreateStory(ctx context.Context, taskGID string, req asana.CreateStoryRequest) (*asana.Story, error)
}

type tasksCommentUpdateClient interface {
	UpdateStory(ctx context.Context, storyGID string, req asana.UpdateStoryRequest) (*asana.Story, error)
}

type tasksCommentDeleteClient interface {
	DeleteStory(ctx context.Context, storyGID string) error
}

type tasksCreateClient interface {
	CreateTask(ctx context.Context, req asana.CreateTaskRequest) (*asana.Task, error)
	AddTaskToSection(ctx context.Context, sectionGID string, taskGID string) error
}

type tasksUpdateClient interface {
	UpdateTask(ctx context.Context, gid string, req asana.UpdateTaskRequest) (*asana.Task, error)
}

type tasksAddFollowerClient interface {
	AddFollowers(ctx context.Context, taskGID string, followers []string) (*asana.Task, error)
}

type tasksRemoveFollowerClient interface {
	RemoveFollowers(ctx context.Context, taskGID string, followers []string) (*asana.Task, error)
}

type tasksStoriesClient interface {
	GetTaskStories(ctx context.Context, taskGID string) (*asana.StoryList, error)
}

func (cmd *TasksListCmd) Run(ctx context.Context, c *cli.Context) error {
	// Resolve workspace if assignee is provided (Asana requires workspace with assignee)
	workspace := cmd.Workspace
	if cmd.Assignee != "" {
		var err error
		workspace, err = c.ResolveWorkspace(cmd.Workspace)
		if err != nil {
			return err
		}
	}

	clientAny := c.ClientOrDefault()
	client, ok := clientAny.(tasksListClient)
	if !ok {
		return fmt.Errorf("failed to list tasks: client does not support listing tasks")
	}
	params := asana.ListTasksParams{
		Project:        cmd.Project,
		Assignee:       cmd.Assignee,
		Workspace:      workspace,
		Limit:          cmd.Limit,
		CompletedSince: cmd.CompletedSince,
		OptFields:      "name,completed,assignee.name,assignee.gid",
	}
	list, err := client.ListTasks(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to list tasks: %w", err)
	}

	renderer := c.RendererOrDefault()
	if c.JSON {
		return renderer.Envelope(list.Data, list.NextPage, nil)
	}

	rows := make([][]string, 0, len(list.Data))
	for _, task := range list.Data {
		assignee := ""
		if task.Assignee != nil {
			assignee = task.Assignee.Name
		}
		rows = append(rows, []string{task.GID, task.Name, assignee, fmt.Sprintf("%t", task.Completed)})
	}
	return renderer.Table([]string{"GID", "NAME", "ASSIGNEE", "COMPLETED"}, rows)
}

func (cmd *TasksSearchCmd) Run(ctx context.Context, c *cli.Context) error {
	// Resolve workspace
	workspace, err := c.ResolveWorkspace(cmd.Workspace)
	if err != nil {
		return err
	}

	clientAny := c.ClientOrDefault()
	client, ok := clientAny.(tasksSearchClient)
	if !ok {
		return fmt.Errorf("failed to search tasks: client does not support searching tasks")
	}
	extra := map[string]string{}
	for _, filter := range cmd.Filter {
		parts := strings.SplitN(filter, "=", 2)
		if len(parts) != 2 || parts[0] == "" {
			return fmt.Errorf("failed to parse filter %q, want key=value", filter)
		}
		extra[parts[0]] = parts[1]
	}
	params := asana.SearchTasksParams{
		Text:              cmd.Text,
		ResourceSubtype:   cmd.ResourceSubtype,
		AssigneeAny:       strings.Join(cmd.AssigneeAny, ","),
		AssigneeNot:       strings.Join(cmd.AssigneeNot, ","),
		PortfoliosAny:     strings.Join(cmd.PortfoliosAny, ","),
		ProjectsAny:       strings.Join(cmd.ProjectsAny, ","),
		ProjectsNot:       strings.Join(cmd.ProjectsNot, ","),
		ProjectsAll:       strings.Join(cmd.ProjectsAll, ","),
		SectionsAny:       strings.Join(cmd.SectionsAny, ","),
		SectionsNot:       strings.Join(cmd.SectionsNot, ","),
		SectionsAll:       strings.Join(cmd.SectionsAll, ","),
		TagsAny:           strings.Join(cmd.TagsAny, ","),
		TagsNot:           strings.Join(cmd.TagsNot, ","),
		TagsAll:           strings.Join(cmd.TagsAll, ","),
		TeamsAny:          strings.Join(cmd.TeamsAny, ","),
		FollowersAny:      strings.Join(cmd.FollowersAny, ","),
		FollowersNot:      strings.Join(cmd.FollowersNot, ","),
		CreatedByAny:      strings.Join(cmd.CreatedByAny, ","),
		CreatedByNot:      strings.Join(cmd.CreatedByNot, ","),
		AssignedByAny:     strings.Join(cmd.AssignedByAny, ","),
		AssignedByNot:     strings.Join(cmd.AssignedByNot, ","),
		LikedByNot:        strings.Join(cmd.LikedByNot, ","),
		CommentedOnByNot:  strings.Join(cmd.CommentedOnByNot, ","),
		DueOn:             cmd.DueOn,
		DueOnBefore:       cmd.DueOnBefore,
		DueOnAfter:        cmd.DueOnAfter,
		DueAtBefore:       cmd.DueAtBefore,
		DueAtAfter:        cmd.DueAtAfter,
		StartOn:           cmd.StartOn,
		StartOnBefore:     cmd.StartOnBefore,
		StartOnAfter:      cmd.StartOnAfter,
		CreatedOn:         cmd.CreatedOn,
		CreatedOnBefore:   cmd.CreatedOnBefore,
		CreatedOnAfter:    cmd.CreatedOnAfter,
		CreatedAtBefore:   cmd.CreatedAtBefore,
		CreatedAtAfter:    cmd.CreatedAtAfter,
		CompletedOn:       cmd.CompletedOn,
		CompletedOnBefore: cmd.CompletedOnBefore,
		CompletedOnAfter:  cmd.CompletedOnAfter,
		CompletedAtBefore: cmd.CompletedAtBefore,
		CompletedAtAfter:  cmd.CompletedAtAfter,
		ModifiedOn:        cmd.ModifiedOn,
		ModifiedOnBefore:  cmd.ModifiedOnBefore,
		ModifiedOnAfter:   cmd.ModifiedOnAfter,
		ModifiedAtBefore:  cmd.ModifiedAtBefore,
		ModifiedAtAfter:   cmd.ModifiedAtAfter,
		IsBlocking:        cmd.IsBlocking,
		IsBlocked:         cmd.IsBlocked,
		HasAttachment:     cmd.HasAttachment,
		Completed:         cmd.Completed,
		IsSubtask:         cmd.IsSubtask,
		SortBy:            cmd.SortBy,
		SortAscending:     cmd.SortAscending,
		Limit:             cmd.Limit,
		OptFields:         asana.TaskDetailOptFields,
		Extra:             extra,
	}
	list, err := client.SearchTasks(ctx, workspace, params)
	if err != nil {
		return fmt.Errorf("failed to search tasks: %w", err)
	}

	renderer := c.RendererOrDefault()
	if c.JSON {
		return renderer.Envelope(list.Data, list.NextPage, nil)
	}

	rows := make([][]string, 0, len(list.Data))
	for _, task := range list.Data {
		assignee := ""
		if task.Assignee != nil {
			assignee = task.Assignee.Name
		}
		rows = append(rows, []string{task.GID, task.Name, assignee, fmt.Sprintf("%t", task.Completed)})
	}
	return renderer.Table([]string{"GID", "NAME", "ASSIGNEE", "COMPLETED"}, rows)
}

func (cmd *TasksGetCmd) Run(ctx context.Context, c *cli.Context) error {
	clientAny := c.ClientOrDefault()
	client, ok := clientAny.(tasksGetClient)
	if !ok {
		return fmt.Errorf("failed to get task: client does not support getting tasks")
	}
	task, err := client.GetTask(ctx, cmd.GID)
	if err != nil {
		return fmt.Errorf("failed to get task: %w", err)
	}

	renderer := c.RendererOrDefault()
	if c.JSON {
		return renderer.JSON(task)
	}

	assignee := ""
	if task.Assignee != nil {
		assignee = task.Assignee.Name
	}
	customFields := formatCustomFields(task.CustomFields)
	projects := formatProjectNames(task.Projects)
	sections := formatMembershipSections(task.Memberships)
	tags := formatTagNames(task.Tags)
	followers := formatUserNames(task.Followers)
	parent := ""
	if task.Parent != nil {
		parent = task.Parent.Name
	}
	// Fetch comments, subtasks, and attachments
	stories, _ := client.GetTaskStories(ctx, cmd.GID)
	subtaskList, _ := client.GetSubtasks(ctx, cmd.GID)
	attachmentList, _ := client.GetTaskAttachments(ctx, cmd.GID)

	comments := formatComments(stories)
	subtasks := formatSubtasks(subtaskList)
	attachments := formatAttachments(attachmentList)

	rows := [][]string{
		{"GID", task.GID},
		{"Name", task.Name},
		{"Assignee", assignee},
		{"Completed", fmt.Sprintf("%t", task.Completed)},
		{"Due", task.DueOn},
		{"Due At", task.DueAt},
		{"Start On", task.StartOn},
		{"Start At", task.StartAt},
		{"Created", task.CreatedAt},
		{"Modified", task.ModifiedAt},
		{"Completed At", task.CompletedAt},
		{"Permalink", task.PermalinkURL},
		{"Projects", projects},
		{"Sections", sections},
		{"Tags", tags},
		{"Followers", followers},
		{"Parent", parent},
		{"Custom Fields", customFields},
		{"Subtasks", subtasks},
		{"Attachments", attachments},
		{"Notes", task.Notes},
		{"Comments", comments},
	}
	return renderer.Table([]string{"FIELD", "VALUE"}, rows)
}

func (cmd *TasksCreateCmd) Run(ctx context.Context, c *cli.Context) error {
	clientAny := c.ClientOrDefault()
	client, ok := clientAny.(tasksCreateClient)
	if !ok {
		return fmt.Errorf("failed to create task: client does not support creating tasks")
	}
	req := asana.CreateTaskRequest{
		Name:     strings.TrimSpace(cmd.Name),
		Notes:    strings.TrimSpace(cmd.Notes),
		Assignee: strings.TrimSpace(cmd.Assignee),
		DueOn:    strings.TrimSpace(cmd.DueOn),
		Parent:   strings.TrimSpace(cmd.Parent),
	}
	if cmd.Project != "" {
		req.Projects = []string{cmd.Project}
	}
	if len(cmd.CustomFields) > 0 {
		cf := make(map[string]string, len(cmd.CustomFields))
		for _, kv := range cmd.CustomFields {
			parts := strings.SplitN(kv, "=", 2)
			if len(parts) != 2 || parts[0] == "" {
				return fmt.Errorf("failed to parse custom field %q, want GID=value", kv)
			}
			cf[parts[0]] = parts[1]
		}
		req.CustomFields = cf
	}

	task, err := client.CreateTask(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}

	if cmd.Section != "" {
		if err := client.AddTaskToSection(ctx, cmd.Section, task.GID); err != nil {
			return fmt.Errorf("task created (%s) but failed to add to section: %w", task.GID, err)
		}
	}

	renderer := c.RendererOrDefault()
	if c.JSON {
		return renderer.JSON(task)
	}
	return renderer.Message("created %s\n", task.GID)
}

func (cmd *TasksUpdateCmd) Run(ctx context.Context, c *cli.Context) error {
	clientAny := c.ClientOrDefault()
	client, ok := clientAny.(tasksUpdateClient)
	if !ok {
		return fmt.Errorf("failed to update task: client does not support updating tasks")
	}
	req := asana.UpdateTaskRequest{}
	if strings.TrimSpace(cmd.Name) != "" {
		name := strings.TrimSpace(cmd.Name)
		req.Name = &name
	}
	if strings.TrimSpace(cmd.Notes) != "" {
		notes := strings.TrimSpace(cmd.Notes)
		req.Notes = &notes
	}
	if strings.TrimSpace(cmd.Assignee) != "" {
		assignee := strings.TrimSpace(cmd.Assignee)
		req.Assignee = &assignee
	}
	if strings.TrimSpace(cmd.DueOn) != "" {
		dueOn := strings.TrimSpace(cmd.DueOn)
		req.DueOn = &dueOn
	}
	if cmd.Completed != nil {
		req.Completed = cmd.Completed
	}

	if req.Name == nil && req.Notes == nil && req.Assignee == nil && req.DueOn == nil && req.Completed == nil {
		return errors.New("failed to update task: no fields to update")
	}

	task, err := client.UpdateTask(ctx, cmd.GID, req)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	renderer := c.RendererOrDefault()
	if c.JSON {
		return renderer.JSON(task)
	}
	return renderer.Message("updated %s\n", task.GID)
}

func (cmd *TasksCommentsCmd) Run(ctx context.Context, c *cli.Context) error {
	clientAny := c.ClientOrDefault()
	client, ok := clientAny.(tasksCommentsClient)
	if !ok {
		return fmt.Errorf("failed to list comments: client does not support listing stories")
	}
	stories, err := client.GetTaskStories(ctx, cmd.GID)
	if err != nil {
		return fmt.Errorf("failed to list comments: %w", err)
	}

	// Filter to only comments
	comments := make([]asana.Story, 0)
	for _, s := range stories.Data {
		if s.Type == "comment" {
			comments = append(comments, s)
		}
	}

	renderer := c.RendererOrDefault()
	if c.JSON {
		return renderer.Envelope(comments, stories.NextPage, nil)
	}

	rows := make([][]string, 0, len(comments))
	for _, s := range comments {
		author := ""
		if s.CreatedBy != nil {
			author = s.CreatedBy.Name
		}
		rows = append(rows, []string{s.GID, s.CreatedAt, author, s.Text})
	}
	return renderer.Table([]string{"GID", "DATE", "AUTHOR", "TEXT"}, rows)
}

func (cmd *TasksCommentAddCmd) Run(ctx context.Context, c *cli.Context) error {
	clientAny := c.ClientOrDefault()
	client, ok := clientAny.(tasksCommentAddClient)
	if !ok {
		return fmt.Errorf("failed to add comment: client does not support creating stories")
	}
	req := asana.CreateStoryRequest{
		Text: strings.TrimSpace(cmd.Text),
	}
	story, err := client.CreateStory(ctx, cmd.GID, req)
	if err != nil {
		return fmt.Errorf("failed to add comment: %w", err)
	}

	renderer := c.RendererOrDefault()
	if c.JSON {
		return renderer.JSON(story)
	}
	return renderer.Message("commented %s\n", story.GID)
}

func (cmd *TasksCommentUpdateCmd) Run(ctx context.Context, c *cli.Context) error {
	clientAny := c.ClientOrDefault()
	client, ok := clientAny.(tasksCommentUpdateClient)
	if !ok {
		return fmt.Errorf("failed to update comment: client does not support updating stories")
	}
	req := asana.UpdateStoryRequest{
		Text: strings.TrimSpace(cmd.Text),
	}
	story, err := client.UpdateStory(ctx, cmd.GID, req)
	if err != nil {
		return fmt.Errorf("failed to update comment: %w", err)
	}

	renderer := c.RendererOrDefault()
	if c.JSON {
		return renderer.JSON(story)
	}
	return renderer.Message("updated %s\n", story.GID)
}

func (cmd *TasksCommentDeleteCmd) Run(ctx context.Context, c *cli.Context) error {
	clientAny := c.ClientOrDefault()
	client, ok := clientAny.(tasksCommentDeleteClient)
	if !ok {
		return fmt.Errorf("failed to delete comment: client does not support deleting stories")
	}
	if err := client.DeleteStory(ctx, cmd.GID); err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}

	renderer := c.RendererOrDefault()
	if c.JSON {
		return renderer.JSON(map[string]string{"gid": cmd.GID})
	}
	return renderer.Message("deleted %s\n", cmd.GID)
}

func (cmd *TasksSubtasksCmd) Run(ctx context.Context, c *cli.Context) error {
	clientAny := c.ClientOrDefault()
	client, ok := clientAny.(tasksSubtasksClient)
	if !ok {
		return fmt.Errorf("failed to list subtasks: client does not support listing subtasks")
	}
	list, err := client.GetSubtasks(ctx, cmd.GID)
	if err != nil {
		return fmt.Errorf("failed to list subtasks: %w", err)
	}

	renderer := c.RendererOrDefault()
	if c.JSON {
		return renderer.Envelope(list.Data, list.NextPage, nil)
	}

	rows := make([][]string, 0, len(list.Data))
	for _, task := range list.Data {
		assignee := ""
		if task.Assignee != nil {
			assignee = task.Assignee.Name
		}
		rows = append(rows, []string{task.GID, task.Name, assignee, fmt.Sprintf("%t", task.Completed)})
	}
	return renderer.Table([]string{"GID", "NAME", "ASSIGNEE", "COMPLETED"}, rows)
}

func (cmd *TasksAttachmentsCmd) Run(ctx context.Context, c *cli.Context) error {
	clientAny := c.ClientOrDefault()
	client, ok := clientAny.(tasksAttachmentsClient)
	if !ok {
		return fmt.Errorf("failed to list attachments: client does not support listing attachments")
	}
	list, err := client.GetTaskAttachments(ctx, cmd.GID)
	if err != nil {
		return fmt.Errorf("failed to list attachments: %w", err)
	}

	renderer := c.RendererOrDefault()
	if c.JSON {
		return renderer.Envelope(list.Data, list.NextPage, nil)
	}

	rows := make([][]string, 0, len(list.Data))
	for _, a := range list.Data {
		rows = append(rows, []string{a.GID, a.Name, a.Host, a.CreatedAt})
	}
	return renderer.Table([]string{"GID", "NAME", "HOST", "CREATED"}, rows)
}

func (cmd *TasksAttachmentGetCmd) Run(ctx context.Context, c *cli.Context) error {
	clientAny := c.ClientOrDefault()
	client, ok := clientAny.(tasksAttachmentGetClient)
	if !ok {
		return fmt.Errorf("failed to get attachment: client does not support getting attachments")
	}
	attachment, err := client.GetAttachment(ctx, cmd.GID)
	if err != nil {
		return fmt.Errorf("failed to get attachment: %w", err)
	}

	renderer := c.RendererOrDefault()
	if c.JSON {
		return renderer.JSON(attachment)
	}

	parent := ""
	if attachment.Parent != nil {
		parent = fmt.Sprintf("%s (%s)", attachment.Parent.Name, attachment.Parent.GID)
	}
	rows := [][]string{
		{"GID", attachment.GID},
		{"Name", attachment.Name},
		{"Host", attachment.Host},
		{"Parent", parent},
		{"Created", attachment.CreatedAt},
		{"Download URL", attachment.DownloadURL},
		{"View URL", attachment.ViewURL},
		{"Permanent URL", attachment.PermanentURL},
		{"Size", fmt.Sprintf("%d", attachment.Size)},
	}
	return renderer.Table([]string{"FIELD", "VALUE"}, rows)
}

func (cmd *TasksStoriesCmd) Run(ctx context.Context, c *cli.Context) error {
	clientAny := c.ClientOrDefault()
	client, ok := clientAny.(tasksStoriesClient)
	if !ok {
		return fmt.Errorf("failed to list stories: client does not support listing stories")
	}
	stories, err := client.GetTaskStories(ctx, cmd.GID)
	if err != nil {
		return fmt.Errorf("failed to list stories: %w", err)
	}

	renderer := c.RendererOrDefault()
	if c.JSON {
		return renderer.Envelope(stories.Data, stories.NextPage, nil)
	}

	rows := make([][]string, 0, len(stories.Data))
	for _, s := range stories.Data {
		author := ""
		if s.CreatedBy != nil {
			author = s.CreatedBy.Name
		}
		rows = append(rows, []string{s.GID, s.Type, s.CreatedAt, author, s.Text})
	}
	return renderer.Table([]string{"GID", "TYPE", "DATE", "AUTHOR", "TEXT"}, rows)
}

func (cmd *TasksAddFollowerCmd) Run(ctx context.Context, c *cli.Context) error {
	clientAny := c.ClientOrDefault()
	client, ok := clientAny.(tasksAddFollowerClient)
	if !ok {
		return fmt.Errorf("failed to add followers: client does not support adding followers")
	}
	_, err := client.AddFollowers(ctx, cmd.GID, cmd.Followers)
	if err != nil {
		return fmt.Errorf("failed to add followers: %w", err)
	}

	renderer := c.RendererOrDefault()
	if c.JSON {
		return renderer.JSON(map[string]any{"task": cmd.GID, "added": cmd.Followers})
	}
	return renderer.Message("added %d follower(s) to %s\n", len(cmd.Followers), cmd.GID)
}

func (cmd *TasksRemoveFollowerCmd) Run(ctx context.Context, c *cli.Context) error {
	clientAny := c.ClientOrDefault()
	client, ok := clientAny.(tasksRemoveFollowerClient)
	if !ok {
		return fmt.Errorf("failed to remove followers: client does not support removing followers")
	}
	_, err := client.RemoveFollowers(ctx, cmd.GID, cmd.Followers)
	if err != nil {
		return fmt.Errorf("failed to remove followers: %w", err)
	}

	renderer := c.RendererOrDefault()
	if c.JSON {
		return renderer.JSON(map[string]any{"task": cmd.GID, "removed": cmd.Followers})
	}
	return renderer.Message("removed %d follower(s) from %s\n", len(cmd.Followers), cmd.GID)
}

func formatComments(stories *asana.StoryList) string {
	if stories == nil || len(stories.Data) == 0 {
		return "none"
	}
	comments := make([]string, 0)
	for _, s := range stories.Data {
		if s.Type != "comment" {
			continue
		}
		author := "unknown"
		if s.CreatedBy != nil {
			author = s.CreatedBy.Name
		}
		comments = append(comments, fmt.Sprintf("[%s] %s: %s", s.CreatedAt, author, s.Text))
	}
	if len(comments) == 0 {
		return "none"
	}
	return strings.Join(comments, "\n")
}

func formatSubtasks(list *asana.SubtaskList) string {
	if list == nil || len(list.Data) == 0 {
		return "none"
	}
	items := make([]string, 0, len(list.Data))
	for _, t := range list.Data {
		status := "[ ]"
		if t.Completed {
			status = "[x]"
		}
		items = append(items, fmt.Sprintf("%s %s (%s)", status, t.Name, t.GID))
	}
	return strings.Join(items, "\n")
}

func formatAttachments(list *asana.AttachmentList) string {
	if list == nil || len(list.Data) == 0 {
		return "none"
	}
	items := make([]string, 0, len(list.Data))
	for _, a := range list.Data {
		items = append(items, fmt.Sprintf("%s (%s)", a.Name, a.GID))
	}
	return strings.Join(items, ", ")
}

func formatCustomFields(fields []asana.CustomField) string {
	if len(fields) == 0 {
		return "none"
	}
	values := make([]string, 0, len(fields))
	for _, field := range fields {
		value := field.DisplayValue
		if value == "" && field.EnumValue != nil {
			value = field.EnumValue.Name
		}
		if value == "" && field.TextValue != "" {
			value = field.TextValue
		}
		if value == "" && field.NumberValue != nil {
			value = fmt.Sprintf("%g", *field.NumberValue)
		}
		if value == "" {
			value = "empty"
		}
		values = append(values, fmt.Sprintf("%s=%s", field.Name, value))
	}
	return strings.Join(values, ", ")
}

func formatProjectNames(projects []asana.Project) string {
	if len(projects) == 0 {
		return "none"
	}
	names := make([]string, 0, len(projects))
	for _, project := range projects {
		if project.Name != "" {
			names = append(names, project.Name)
		}
	}
	if len(names) == 0 {
		return "none"
	}
	return strings.Join(names, ", ")
}

func formatMembershipSections(memberships []asana.Membership) string {
	if len(memberships) == 0 {
		return "none"
	}
	names := make([]string, 0, len(memberships))
	for _, membership := range memberships {
		if membership.Section.Name != "" {
			names = append(names, membership.Section.Name)
		}
	}
	if len(names) == 0 {
		return "none"
	}
	return strings.Join(names, ", ")
}

func formatTagNames(tags []asana.Tag) string {
	if len(tags) == 0 {
		return "none"
	}
	names := make([]string, 0, len(tags))
	for _, tag := range tags {
		if tag.Name != "" {
			names = append(names, tag.Name)
		}
	}
	if len(names) == 0 {
		return "none"
	}
	return strings.Join(names, ", ")
}

func formatUserNames(users []asana.User) string {
	if len(users) == 0 {
		return "none"
	}
	names := make([]string, 0, len(users))
	for _, user := range users {
		if user.Name != "" {
			names = append(names, user.Name)
		}
	}
	if len(names) == 0 {
		return "none"
	}
	return strings.Join(names, ", ")
}
