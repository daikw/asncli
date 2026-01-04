package cmd

import (
	"bytes"
	"os"
	"testing"
)

func TestPromptTokenNonTerminal(t *testing.T) {
	withStdin(t, "token\n", func() {
		var prompt bytes.Buffer
		got, err := promptToken(&prompt)
		if err != nil {
			t.Fatalf("promptToken returned unexpected error: %v", err)
		}
		if want := "token\n"; got != want {
			t.Errorf("promptToken() = %q, want %q", got, want)
		}
	})
}

func withStdin(t *testing.T, input string, fn func()) {
	t.Helper()
	orig := os.Stdin
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe for stdin: %v", err)
	}
	if _, err := w.Write([]byte(input)); err != nil {
		_ = r.Close()
		_ = w.Close()
		t.Fatalf("failed to write to stdin pipe: %v", err)
	}
	if err := w.Close(); err != nil {
		_ = r.Close()
		t.Fatalf("failed to close stdin pipe writer: %v", err)
	}
	os.Stdin = r
	t.Cleanup(func() {
		os.Stdin = orig
		_ = r.Close()
	})
	fn()
}
