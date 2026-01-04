package asana

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type errorTokenSource struct {
	err error
}

func (e errorTokenSource) Token(ctx context.Context) (string, error) { return "", e.err }

func TestSearchTasksDefaultOptFields(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/workspaces/w1/tasks/search" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		q := r.URL.Query()
		if got := q.Get("opt_fields"); got != TaskDetailOptFields {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if got := q.Get("text"); got != "query" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if got := q.Get("completed"); got != "true" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if got := q.Get("limit"); got != "3" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data":      []map[string]string{{"gid": "1", "name": "Task"}},
			"next_page": map[string]string{"offset": "next", "uri": "uri"},
		})
	}))
	defer server.Close()

	completed := true
	client := NewClient(staticToken{token: "token"}, server.Client()).WithBaseURL(server.URL)
	list, err := client.SearchTasks(context.Background(), "w1", SearchTasksParams{
		Text:      "query",
		Completed: &completed,
		Limit:     3,
	})
	if err != nil {
		t.Fatalf("SearchTasks returned unexpected error: %v", err)
	}
	if got, want := len(list.Data), 1; got != want {
		t.Errorf("SearchTasks returned %d tasks, want %d", got, want)
	}
	if list.NextPage == nil {
		t.Error("SearchTasks NextPage is nil, want non-nil for pagination")
	}
}

func TestListCustomFieldsForWorkspaceQuery(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/workspaces/w1/custom_fields" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		q := r.URL.Query()
		if got := q.Get("limit"); got != "2" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if got := q.Get("offset"); got != "next" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if got := q.Get("opt_fields"); got != "name" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data":      []map[string]string{{"gid": "cf1", "name": "Size"}},
			"next_page": map[string]string{"offset": "next", "uri": "uri"},
		})
	}))
	defer server.Close()

	client := NewClient(staticToken{token: "token"}, server.Client()).WithBaseURL(server.URL)
	list, err := client.ListCustomFieldsForWorkspace(context.Background(), "w1", ListCustomFieldsParams{
		Limit:     2,
		Offset:    "next",
		OptFields: "name",
	})
	if err != nil {
		t.Fatalf("ListCustomFieldsForWorkspace returned unexpected error: %v", err)
	}
	if got, want := len(list.Data), 1; got != want {
		t.Errorf("ListCustomFieldsForWorkspace returned %d fields, want %d", got, want)
	}
	if list.NextPage == nil {
		t.Error("ListCustomFieldsForWorkspace NextPage is nil, want non-nil for pagination")
	}
}

func TestListWorkspaceMembershipsForUser(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/me/workspace_memberships" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if got := r.URL.Query().Get("opt_fields"); got == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{{
				"gid":          "wm1",
				"workspace":    map[string]string{"gid": "w1", "name": "Main"},
				"is_active":    true,
				"is_admin":     false,
				"is_guest":     false,
				"is_view_only": false,
			}},
		})
	}))
	defer server.Close()

	client := NewClient(staticToken{token: "token"}, server.Client()).WithBaseURL(server.URL)
	list, err := client.ListWorkspaceMembershipsForUser(context.Background(), "me")
	if err != nil {
		t.Fatalf("ListWorkspaceMembershipsForUser returned unexpected error: %v", err)
	}
	if got, want := len(list.Data), 1; got != want {
		t.Errorf("ListWorkspaceMembershipsForUser returned %d memberships, want %d", got, want)
	}
	if got, want := list.Data[0].Workspace.GID, "w1"; got != want {
		t.Errorf("first workspace GID = %q, want %q", got, want)
	}
}

func TestWithBaseURLTrailingSlash(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/me" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]string{"gid": "1", "name": "Jane"},
		})
	}))
	defer server.Close()

	client := NewClient(staticToken{token: "token"}, server.Client()).WithBaseURL(server.URL + "/")
	_, err := client.GetMe(context.Background())
	if err != nil {
		t.Errorf("GetMe with trailing slash base URL returned error: %v", err)
	}
}

func TestTokenSourceError(t *testing.T) {
	want := errors.New("token error")
	client := NewClient(errorTokenSource{err: want}, nil)
	_, err := client.GetMe(context.Background())
	if !errors.Is(err, want) {
		t.Errorf("GetMe error = %v, want %v", err, want)
	}
}
