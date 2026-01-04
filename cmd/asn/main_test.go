package main

import (
	"io"
	"os"
	"strings"
	"testing"
)

func TestRunInvalidArgs(t *testing.T) {
	stdout := newTempFile(t)
	stderr := newTempFile(t)

	exitCode := run([]string{"asn", "--nope"}, stdout, stderr)
	if got, want := exitCode, 2; got != want {
		t.Errorf("exit code = %d, want %d for invalid argument", got, want)
	}

	if _, err := stderr.Seek(0, 0); err != nil {
		t.Fatalf("failed to seek stderr: %v", err)
	}
	b, err := io.ReadAll(stderr)
	if err != nil {
		t.Fatalf("failed to read stderr: %v", err)
	}
	if got, want := string(b), "Try 'asn --help'"; !strings.Contains(got, want) {
		t.Errorf("stderr = %q, want to contain %q", got, want)
	}
}

func TestRunUnknownCommand(t *testing.T) {
	stdout := newTempFile(t)
	stderr := newTempFile(t)

	exitCode := run([]string{"asn", "unknowncommand"}, stdout, stderr)
	if exitCode == 0 {
		t.Error("exit code = 0, want non-zero for unknown command")
	}
}

func TestRunJSONFlag(t *testing.T) {
	stdout := newTempFile(t)
	stderr := newTempFile(t)

	// --json without a valid command should still parse
	exitCode := run([]string{"asn", "--json", "unknowncommand"}, stdout, stderr)
	if exitCode == 0 {
		t.Error("exit code = 0, want non-zero for unknown command")
	}
}

func newTempFile(t *testing.T) *os.File {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "asncli-test-*")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	return f
}
