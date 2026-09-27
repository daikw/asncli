package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/michalvavra/asncli/internal/asana"
	"github.com/michalvavra/asncli/internal/cli"
)

type moveClient struct {
	section       *asana.Section
	task          *asana.Task
	addErr        error
	addCalls      int
	gotSectionGID string
	gotTaskGID    string
}

func (f *moveClient) GetSection(ctx context.Context, gid string) (*asana.Section, error) {
	return f.section, nil
}

func (f *moveClient) GetTask(ctx context.Context, gid string) (*asana.Task, error) {
	return f.task, nil
}

func (f *moveClient) AddTaskToSection(ctx context.Context, sectionGID, taskGID string) error {
	f.addCalls++
	f.gotSectionGID = sectionGID
	f.gotTaskGID = taskGID
	return f.addErr
}

func boardTask(sectionGID string) *asana.Task {
	return &asana.Task{
		GID: "t1",
		Memberships: []asana.Membership{
			{Project: asana.Project{GID: "other"}, Section: asana.Section{GID: "x", Name: "Elsewhere"}},
			{Project: asana.Project{GID: "p1"}, Section: asana.Section{GID: sectionGID, Name: "ToDo"}},
		},
	}
}

func reviewSection() *asana.Section {
	return &asana.Section{GID: "s-review", Name: "レビュー中", Project: &asana.Project{GID: "p1"}}
}

func runMove(t *testing.T, client *moveClient, jsonOut bool) (string, error) {
	t.Helper()
	buf := &bytes.Buffer{}
	ctx := &cli.Context{Stdout: buf, Stderr: &bytes.Buffer{}, JSON: jsonOut, Client: client}
	cmd := TasksMoveCmd{GID: "t1", Section: "s-review"}
	err := cmd.Run(context.Background(), ctx)
	return buf.String(), err
}

func TestTasksMoveJSON(t *testing.T) {
	client := &moveClient{section: reviewSection(), task: boardTask("s-todo")}
	out, err := runMove(t, client, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client.addCalls != 1 || client.gotSectionGID != "s-review" || client.gotTaskGID != "t1" {
		t.Errorf("addTask calls=%d section=%q task=%q, want 1 / s-review / t1", client.addCalls, client.gotSectionGID, client.gotTaskGID)
	}
	var env struct {
		Data map[string]string `json:"data"`
	}
	if err := json.Unmarshal([]byte(out), &env); err != nil {
		t.Fatalf("failed to unmarshal JSON output %q: %v", out, err)
	}
	got := env.Data
	want := map[string]string{
		"task_gid":         "t1",
		"project_gid":      "p1",
		"from_section_gid": "s-todo",
		"section_gid":      "s-review",
		"status":           "moved",
	}
	for k, v := range want {
		if got[k] != v {
			t.Errorf("%s = %q, want %q", k, got[k], v)
		}
	}
}

func TestTasksMoveHumanReadable(t *testing.T) {
	client := &moveClient{section: reviewSection(), task: boardTask("s-todo")}
	out, err := runMove(t, client, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "moved t1") || !strings.Contains(out, "ToDo") || !strings.Contains(out, "レビュー中") {
		t.Errorf("output = %q, want to mention task, source and destination sections", out)
	}
}

func TestTasksMoveAlreadyInSectionIsNoop(t *testing.T) {
	client := &moveClient{section: reviewSection(), task: boardTask("s-review")}
	out, err := runMove(t, client, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client.addCalls != 0 {
		t.Errorf("addTask calls = %d, want 0 (re-adding would reorder the task to the top)", client.addCalls)
	}
	if !strings.Contains(out, `"already_in_section"`) {
		t.Errorf("output = %q, want status already_in_section", out)
	}
}

func TestTasksMoveRejectsTaskOutsideSectionProject(t *testing.T) {
	task := &asana.Task{GID: "t1", Memberships: []asana.Membership{
		{Project: asana.Project{GID: "other"}, Section: asana.Section{GID: "x"}},
	}}
	client := &moveClient{section: reviewSection(), task: task}
	_, err := runMove(t, client, true)
	if err == nil || !strings.Contains(err.Error(), "p1") {
		t.Fatalf("error = %v, want error mentioning project p1", err)
	}
	if client.addCalls != 0 {
		t.Errorf("addTask calls = %d, want 0", client.addCalls)
	}
}

func TestTasksMoveAddTaskError(t *testing.T) {
	client := &moveClient{section: reviewSection(), task: boardTask("s-todo"), addErr: errors.New("boom")}
	if _, err := runMove(t, client, true); err == nil || !strings.Contains(err.Error(), "boom") {
		t.Fatalf("error = %v, want wrapped addTask error", err)
	}
}

func TestTasksMoveInvalidClient(t *testing.T) {
	ctx := &cli.Context{Stdout: &bytes.Buffer{}, Stderr: &bytes.Buffer{}, Client: struct{}{}}
	cmd := TasksMoveCmd{GID: "t1", Section: "s1"}
	if err := cmd.Run(context.Background(), ctx); err == nil {
		t.Fatal("want error for client without move support")
	}
}
