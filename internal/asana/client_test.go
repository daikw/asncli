package asana

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type staticToken struct {
	token string
}

func (s staticToken) Token(ctx context.Context) (string, error) { return s.token, nil }

func TestGetMe(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer test-token" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if r.URL.Path != "/users/me" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]string{"gid": "1", "name": "Jane"},
		})
	}))
	defer server.Close()

	client := NewClient(staticToken{token: "test-token"}, server.Client()).WithBaseURL(server.URL)
	user, err := client.GetMe(context.Background())
	if err != nil {
		t.Fatalf("GetMe returned unexpected error: %v", err)
	}
	if got, want := user.GID, "1"; got != want {
		t.Errorf("GetMe().GID = %q, want %q", got, want)
	}
	if got, want := user.Name, "Jane"; got != want {
		t.Errorf("GetMe().Name = %q, want %q", got, want)
	}
}

func TestListTasksQuery(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/tasks" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if r.URL.Query().Get("project") != "p1" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if r.URL.Query().Get("assignee") != "me" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data":      []map[string]string{{"gid": "10", "name": "Task"}},
			"next_page": map[string]string{"offset": "next", "uri": "uri"},
		})
	}))
	defer server.Close()

	client := NewClient(staticToken{token: "token"}, server.Client()).WithBaseURL(server.URL)
	list, err := client.ListTasks(context.Background(), ListTasksParams{Project: "p1", Assignee: "me"})
	if err != nil {
		t.Fatalf("ListTasks returned unexpected error: %v", err)
	}
	if got, want := len(list.Data), 1; got != want {
		t.Errorf("ListTasks returned %d tasks, want %d", got, want)
	}
	if got, want := list.Data[0].GID, "10"; got != want {
		t.Errorf("ListTasks first task GID = %q, want %q", got, want)
	}
	if list.NextPage == nil {
		t.Fatal("ListTasks NextPage is nil, want non-nil for pagination")
	}
	if got, want := list.NextPage.Offset, "next"; got != want {
		t.Errorf("ListTasks NextPage.Offset = %q, want %q", got, want)
	}
}

func TestListProjectsQuery(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/projects" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if r.URL.Query().Get("workspace") != "w1" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if r.URL.Query().Get("archived") != "true" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]string{{"gid": "p1", "name": "Project"}},
		})
	}))
	defer server.Close()

	client := NewClient(staticToken{token: "token"}, server.Client()).WithBaseURL(server.URL)
	archived := true
	list, err := client.ListProjects(context.Background(), ListProjectsParams{
		Workspace: "w1",
		Archived:  &archived,
	})
	if err != nil {
		t.Fatalf("ListProjects returned unexpected error: %v", err)
	}
	if got, want := len(list.Data), 1; got != want {
		t.Errorf("ListProjects returned %d projects, want %d", got, want)
	}
	if got, want := list.Data[0].GID, "p1"; got != want {
		t.Errorf("ListProjects first project GID = %q, want %q", got, want)
	}
}

func TestGetProject(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/projects/7" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]string{"gid": "7", "name": "Example"},
		})
	}))
	defer server.Close()

	client := NewClient(staticToken{token: "token"}, server.Client()).WithBaseURL(server.URL)
	project, err := client.GetProject(context.Background(), "7")
	if err != nil {
		t.Fatalf("GetProject returned unexpected error: %v", err)
	}
	if got, want := project.GID, "7"; got != want {
		t.Errorf("GetProject GID = %q, want %q", got, want)
	}
}

func TestCreateTask(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if r.URL.Path != "/tasks" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if r.URL.Query().Get("opt_fields") == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var payload map[string]CreateTaskRequest
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if payload["data"].Name != "New task" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]string{"gid": "99", "name": "New task"},
		})
	}))
	defer server.Close()

	client := NewClient(staticToken{token: "token"}, server.Client()).WithBaseURL(server.URL)
	task, err := client.CreateTask(context.Background(), CreateTaskRequest{Name: "New task"})
	if err != nil {
		t.Fatalf("CreateTask returned unexpected error: %v", err)
	}
	if got, want := task.GID, "99"; got != want {
		t.Errorf("CreateTask GID = %q, want %q", got, want)
	}
}

func TestUpdateTask(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if r.URL.Path != "/tasks/42" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		var payload map[string]UpdateTaskRequest
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if payload["data"].Notes == nil || *payload["data"].Notes != "note" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]string{"gid": "42", "name": "Updated"},
		})
	}))
	defer server.Close()

	client := NewClient(staticToken{token: "token"}, server.Client()).WithBaseURL(server.URL)
	notes := "note"
	task, err := client.UpdateTask(context.Background(), "42", UpdateTaskRequest{Notes: &notes})
	if err != nil {
		t.Fatalf("UpdateTask returned unexpected error: %v", err)
	}
	if got, want := task.GID, "42"; got != want {
		t.Errorf("UpdateTask GID = %q, want %q", got, want)
	}
}

func TestErrorResponseMessage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"errors": []map[string]string{{"message": "bad request"}},
		})
	}))
	defer server.Close()

	client := NewClient(staticToken{token: "token"}, server.Client()).WithBaseURL(server.URL)
	_, err := client.GetTask(context.Background(), "1")
	if err == nil {
		t.Fatal("GetTask should return error for bad request")
	}
	if got, want := err.Error(), "response: bad request"; got != want {
		t.Errorf("GetTask error = %q, want %q", got, want)
	}
}

func TestUnauthorized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	client := NewClient(staticToken{token: "bad"}, server.Client()).WithBaseURL(server.URL)
	_, err := client.GetMe(context.Background())
	if err == nil {
		t.Fatal("GetMe should return error for unauthorized request")
	}
	if !errors.Is(err, ErrUnauthorized) {
		t.Errorf("GetMe error = %v, want ErrUnauthorized", err)
	}
}
