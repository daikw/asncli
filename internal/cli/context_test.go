package cli

import (
	"bytes"
	"context"
	"testing"

	"github.com/michalvavra/asncli/internal/asana"
	"github.com/michalvavra/asncli/internal/auth"
	"github.com/michalvavra/asncli/internal/config"
)

type stubTokenSource struct{}

func (stubTokenSource) Token(ctx context.Context) (string, error) { return "token", nil }

type stubRenderer struct{}

func (s *stubRenderer) JSON(data any) error { return nil }
func (s *stubRenderer) Envelope(data any, next any, warnings []string) error {
	return nil
}
func (s *stubRenderer) Table(headers []string, rows [][]string) error { return nil }
func (s *stubRenderer) Message(format string, args ...any) error      { return nil }

func TestContextClientOrDefault_UsesExplicitClient(t *testing.T) {
	want := struct{}{}
	ctx := Context{Client: want}

	got := ctx.ClientOrDefault()
	if got != want {
		t.Errorf("ClientOrDefault() = %#v, want %#v", got, want)
	}
}

func TestContextClientOrDefault_UsesFactory(t *testing.T) {
	want := asana.NewClient(stubTokenSource{}, nil)
	ctx := Context{
		TokenSource: stubTokenSource{},
		ClientFactory: func(source auth.TokenSource) *asana.Client {
			if source == nil {
				t.Fatal("ClientFactory received nil token source")
			}
			return want
		},
	}

	got := ctx.ClientOrDefault()
	if got != want {
		t.Errorf("ClientOrDefault() = %#v, want %#v", got, want)
	}
}

func TestContextClientOrDefault_Default(t *testing.T) {
	ctx := Context{Stdout: &bytes.Buffer{}, Stderr: &bytes.Buffer{}}

	got := ctx.ClientOrDefault()
	if _, ok := got.(*asana.Client); !ok {
		t.Errorf("ClientOrDefault() type = %T, want *asana.Client", got)
	}
}

func TestContextRendererOrDefault_UsesExplicitRenderer(t *testing.T) {
	want := &stubRenderer{}
	ctx := Context{Renderer: want}

	got := ctx.RendererOrDefault()
	if got != want {
		t.Errorf("RendererOrDefault() = %#v, want %#v", got, want)
	}
}

func TestContextRendererOrDefault_Default(t *testing.T) {
	ctx := Context{Stdout: &bytes.Buffer{}, Stderr: &bytes.Buffer{}}

	got := ctx.RendererOrDefault()
	if _, ok := got.(DefaultRenderer); !ok {
		t.Errorf("RendererOrDefault() type = %T, want DefaultRenderer", got)
	}
}

func TestContextResolveWorkspace_FlagValue(t *testing.T) {
	ctx := Context{
		Stdout: &bytes.Buffer{},
		Stderr: &bytes.Buffer{},
		Config: &config.Config{DefaultWorkspace: "config-ws"},
	}

	got, err := ctx.ResolveWorkspace("flag-ws")
	if err != nil {
		t.Fatalf("ResolveWorkspace returned unexpected error: %v", err)
	}
	if got != "flag-ws" {
		t.Errorf("ResolveWorkspace() = %q, want %q (flag should take precedence)", got, "flag-ws")
	}
}

func TestContextResolveWorkspace_ConfigValue(t *testing.T) {
	ctx := Context{
		Stdout: &bytes.Buffer{},
		Stderr: &bytes.Buffer{},
		Config: &config.Config{DefaultWorkspace: "config-ws"},
	}

	got, err := ctx.ResolveWorkspace("")
	if err != nil {
		t.Fatalf("ResolveWorkspace returned unexpected error: %v", err)
	}
	if got != "config-ws" {
		t.Errorf("ResolveWorkspace() = %q, want %q (should fall back to config)", got, "config-ws")
	}
}

func TestContextResolveWorkspace_NoWorkspace(t *testing.T) {
	ctx := Context{
		Stdout: &bytes.Buffer{},
		Stderr: &bytes.Buffer{},
		Config: &config.Config{},
	}

	_, err := ctx.ResolveWorkspace("")
	if err == nil {
		t.Fatal("ResolveWorkspace should return error when no workspace configured")
	}
}

func TestContextResolveWorkspace_LoadsConfigFromFile(t *testing.T) {
	// When Config is nil and flag is provided, should still work
	ctx := Context{
		Stdout: &bytes.Buffer{},
		Stderr: &bytes.Buffer{},
		Config: nil, // Force loading
	}

	// Providing a flag value should work even without config
	got, err := ctx.ResolveWorkspace("flag-provided")
	if err != nil {
		t.Fatalf("ResolveWorkspace returned unexpected error: %v", err)
	}
	if got != "flag-provided" {
		t.Errorf("ResolveWorkspace() = %q, want %q", got, "flag-provided")
	}
}
