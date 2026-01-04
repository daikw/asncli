package asana

import (
	"context"
	"os"
	"testing"
)

func TestIntegrationGetMe(t *testing.T) {
	token := os.Getenv("ASANA_TOKEN")
	if token == "" {
		t.Skip("ASANA_TOKEN not set")
	}

	client := NewClient(staticToken{token: token}, nil)
	user, err := client.GetMe(context.Background())
	if err != nil {
		t.Fatalf("GetMe() = _, %v; want nil", err)
	}
	if user.GID == "" || user.Name == "" {
		t.Fatalf("GetMe() = %#v; want GID and Name non-empty", user)
	}
}

func TestIntegrationListTasksProject(t *testing.T) {
	token := os.Getenv("ASANA_TOKEN")
	project := os.Getenv("ASANA_PROJECT")
	if token == "" || project == "" {
		t.Skip("ASANA_TOKEN or ASANA_PROJECT not set")
	}

	client := NewClient(staticToken{token: token}, nil)
	list, err := client.ListTasks(context.Background(), ListTasksParams{
		Project:  project,
		Assignee: "me",
		Limit:    5,
	})
	if err != nil {
		t.Fatalf("ListTasks() = _, %v; want nil", err)
	}
	if list == nil {
		t.Fatalf("ListTasks() = nil, _; want non-nil")
	}
}
