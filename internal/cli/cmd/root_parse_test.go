package cmd

import (
	"testing"

	"github.com/alecthomas/kong"
)

func parseRoot(t *testing.T, args ...string) *Root {
	t.Helper()
	root := &Root{}
	parser, err := kong.New(root)
	if err != nil {
		t.Fatalf("kong.New: %v", err)
	}
	if _, err := parser.Parse(args); err != nil {
		t.Fatalf("parse %v: %v", args, err)
	}
	return root
}

func TestParseUpdateCompleted(t *testing.T) {
	tests := []struct {
		args []string
		want *bool
	}{
		{[]string{"tasks", "update", "1", "--name=x"}, nil},
		{[]string{"tasks", "update", "1", "--completed"}, ptr(true)},
		{[]string{"tasks", "update", "1", "--completed=true"}, ptr(true)},
		{[]string{"tasks", "update", "1", "--completed=false"}, ptr(false)},
	}
	for _, tt := range tests {
		got := parseRoot(t, tt.args...).Tasks.Update.Completed
		switch {
		case tt.want == nil && got != nil:
			t.Errorf("%v: Completed = %v, want nil", tt.args, *got)
		case tt.want != nil && (got == nil || *got != *tt.want):
			t.Errorf("%v: Completed = %v, want %v", tt.args, got, *tt.want)
		}
	}
}

func TestParseSectionsListProjectFlag(t *testing.T) {
	root := parseRoot(t, "sections", "list", "--project=p1")
	if root.Sections.List.Project != "p1" {
		t.Errorf("Project = %q, want p1", root.Sections.List.Project)
	}
}

func TestParseTasksMove(t *testing.T) {
	root := parseRoot(t, "tasks", "move", "t1", "--section=s1")
	if root.Tasks.Move.GID != "t1" || root.Tasks.Move.Section != "s1" {
		t.Errorf("move = %+v, want t1 / s1", root.Tasks.Move)
	}
}

func ptr[T any](v T) *T { return &v }
