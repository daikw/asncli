package testutil

import (
	"bytes"
	"testing"

	"github.com/michalvavra/asncli/internal/cli"
)

// NewTestContext creates a CLI context for testing with buffered stdout/stderr.
func NewTestContext(t *testing.T, json bool, client interface{}) (*cli.Context, *bytes.Buffer, *bytes.Buffer) {
	t.Helper()
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	ctx := &cli.Context{
		Stdout: stdout,
		Stderr: stderr,
		JSON:   json,
		Client: client,
	}
	return ctx, stdout, stderr
}

// AssertNoError fails the test if err is not nil.
func AssertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// AssertError fails the test if err is nil.
func AssertError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("expected error but got nil")
	}
}
