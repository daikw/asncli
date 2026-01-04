package cli

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestDefaultRendererJSON(t *testing.T) {
	buf := &bytes.Buffer{}
	r := NewRenderer(buf, &bytes.Buffer{}, true)

	if err := r.JSON(map[string]string{"ok": "yes"}); err != nil {
		t.Fatalf("Renderer.JSON returned unexpected error: %v", err)
	}

	var env Envelope
	if err := json.Unmarshal(buf.Bytes(), &env); err != nil {
		t.Fatalf("failed to unmarshal JSON output: %v", err)
	}
	data, ok := env.Data.(map[string]any)
	if !ok {
		t.Fatalf("envelope data type = %T, want map[string]any", env.Data)
	}
	if got, want := data["ok"], "yes"; got != want {
		t.Errorf("data[ok] = %v, want %v", got, want)
	}
}

func TestDefaultRendererTable(t *testing.T) {
	buf := &bytes.Buffer{}
	r := NewRenderer(buf, &bytes.Buffer{}, false)

	if err := r.Table([]string{"A", "B"}, [][]string{{"1", "2"}}); err != nil {
		t.Fatalf("Renderer.Table returned unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "A") || !strings.Contains(out, "1") {
		t.Errorf("table output = %q, want to contain header 'A' and value '1'", out)
	}
}

func TestDefaultRendererMessage(t *testing.T) {
	buf := &bytes.Buffer{}
	r := NewRenderer(buf, &bytes.Buffer{}, false)

	if err := r.Message("hello %s\n", "world"); err != nil {
		t.Fatalf("Renderer.Message returned unexpected error: %v", err)
	}
	if got, want := buf.String(), "hello world\n"; got != want {
		t.Errorf("message output = %q, want %q", got, want)
	}
}

func TestDefaultRendererEnvelope(t *testing.T) {
	buf := &bytes.Buffer{}
	r := NewRenderer(buf, &bytes.Buffer{}, true)

	if err := r.Envelope(map[string]string{"key": "value"}, "next-cursor", []string{"warning1"}); err != nil {
		t.Fatalf("Renderer.Envelope returned unexpected error: %v", err)
	}

	var env Envelope
	if err := json.Unmarshal(buf.Bytes(), &env); err != nil {
		t.Fatalf("failed to unmarshal JSON output: %v", err)
	}
	if env.Next != "next-cursor" {
		t.Errorf("envelope next cursor = %v, want %q", env.Next, "next-cursor")
	}
	if len(env.Warnings) != 1 || env.Warnings[0] != "warning1" {
		t.Errorf("envelope warnings = %v, want [warning1]", env.Warnings)
	}
}
