package asana

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetSection(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/sections/s1" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if got := r.URL.Query().Get("opt_fields"); got != "name,project.gid,project.name" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"gid":     "s1",
				"name":    "In review",
				"project": map[string]string{"gid": "p1", "name": "Board"},
			},
		})
	}))
	defer server.Close()

	client := NewClient(staticToken{token: "token"}, server.Client()).WithBaseURL(server.URL)
	section, err := client.GetSection(context.Background(), "s1")
	if err != nil {
		t.Fatalf("GetSection returned unexpected error: %v", err)
	}
	if section.GID != "s1" || section.Name != "In review" {
		t.Errorf("section = %+v, want gid s1 / name In review", section)
	}
	if section.Project == nil || section.Project.GID != "p1" {
		t.Errorf("section.Project = %+v, want gid p1", section.Project)
	}
}

func TestAddTaskToSection(t *testing.T) {
	var gotBody map[string]map[string]string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/sections/s1/addTask" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{}})
	}))
	defer server.Close()

	client := NewClient(staticToken{token: "token"}, server.Client()).WithBaseURL(server.URL)
	if err := client.AddTaskToSection(context.Background(), "s1", "t1"); err != nil {
		t.Fatalf("AddTaskToSection returned unexpected error: %v", err)
	}
	if got := gotBody["data"]["task"]; got != "t1" {
		t.Errorf("request data.task = %q, want %q", got, "t1")
	}
}
