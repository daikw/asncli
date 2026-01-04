package cli

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestOutputPrintJSON(t *testing.T) {
	buf := &bytes.Buffer{}
	out := Output{Stdout: buf}

	if err := out.PrintJSON(map[string]string{"status": "ok"}); err != nil {
		t.Fatalf("PrintJSON returned unexpected error: %v", err)
	}

	var env Envelope
	if err := json.Unmarshal(buf.Bytes(), &env); err != nil {
		t.Fatalf("failed to unmarshal JSON output: %v", err)
	}
	data, ok := env.Data.(map[string]any)
	if !ok {
		t.Fatalf("envelope data type = %T, want map[string]any", env.Data)
	}
	if got, want := data["status"], "ok"; got != want {
		t.Errorf("data[status] = %v, want %v", got, want)
	}
}

func TestOutputPrintEnvelope(t *testing.T) {
	buf := &bytes.Buffer{}
	out := Output{Stdout: buf}

	if err := out.PrintEnvelope(Envelope{
		Data:     map[string]string{"message": "hi"},
		Next:     map[string]string{"offset": "next"},
		Warnings: []string{"warn"},
	}); err != nil {
		t.Fatalf("PrintEnvelope returned unexpected error: %v", err)
	}

	var env Envelope
	if err := json.Unmarshal(buf.Bytes(), &env); err != nil {
		t.Fatalf("failed to unmarshal JSON output: %v", err)
	}
	if env.Data == nil || env.Next == nil {
		t.Fatalf("envelope incomplete: Data = %v, Next = %v, want both non-nil", env.Data, env.Next)
	}
	if got, want := len(env.Warnings), 1; got != want {
		t.Errorf("warnings count = %d, want %d", got, want)
	}
	if got, want := env.Warnings[0], "warn"; got != want {
		t.Errorf("first warning = %q, want %q", got, want)
	}
}

func TestOutputTable(t *testing.T) {
	buf := &bytes.Buffer{}
	out := Output{Stdout: buf}

	w := out.Table()
	if w == nil {
		t.Fatal("Table returned nil writer")
	}
}
