package asana

import (
	"reflect"
	"testing"
)

func TestBuildQuery(t *testing.T) {
	tests := []struct {
		name   string
		params any
		want   map[string]string
	}{
		{
			name: "simple strings",
			params: struct {
				Text      string `query:"text"`
				Workspace string `query:"workspace"`
			}{
				Text:      "hello",
				Workspace: "12345",
			},
			want: map[string]string{
				"text":      "hello",
				"workspace": "12345",
			},
		},
		{
			name: "empty strings omitted",
			params: struct {
				Text      string `query:"text"`
				Workspace string `query:"workspace"`
			}{
				Text:      "hello",
				Workspace: "",
			},
			want: map[string]string{
				"text": "hello",
			},
		},
		{
			name: "dotted query keys",
			params: struct {
				AssigneeAny string `query:"assignee.any"`
				AssigneeNot string `query:"assignee.not"`
			}{
				AssigneeAny: "123",
				AssigneeNot: "456",
			},
			want: map[string]string{
				"assignee.any": "123",
				"assignee.not": "456",
			},
		},
		{
			name: "integers",
			params: struct {
				Limit  int `query:"limit"`
				Offset int `query:"offset"`
			}{
				Limit:  50,
				Offset: 0,
			},
			want: map[string]string{
				"limit": "50",
			},
		},
		{
			name: "bool pointers",
			params: struct {
				Completed  *bool `query:"completed"`
				IsBlocking *bool `query:"is_blocking"`
				IsBlocked  *bool `query:"is_blocked"`
			}{
				Completed:  boolPtr(true),
				IsBlocking: boolPtr(false),
				IsBlocked:  nil,
			},
			want: map[string]string{
				"completed":   "true",
				"is_blocking": "false",
			},
		},
		{
			name: "extra map merged",
			params: struct {
				Text  string            `query:"text"`
				Extra map[string]string `query:"-"`
			}{
				Text: "hello",
				Extra: map[string]string{
					"custom_field.123": "value",
					"archived":         "true",
				},
			},
			want: map[string]string{
				"text":             "hello",
				"custom_field.123": "value",
				"archived":         "true",
			},
		},
		{
			name: "fields without query tag ignored",
			params: struct {
				Text      string `query:"text"`
				ProjectID string
				CreatedAt string
			}{
				Text:      "hello",
				ProjectID: "123",
				CreatedAt: "2024-01-01",
			},
			want: map[string]string{
				"text": "hello",
			},
		},
		{
			name: "skip fields with dash tag",
			params: struct {
				Text     string `query:"text"`
				Internal string `query:"-"`
			}{
				Text:     "hello",
				Internal: "skip me",
			},
			want: map[string]string{
				"text": "hello",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildQuery(tt.params)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}

func boolPtr(b bool) *bool {
	return &b
}
