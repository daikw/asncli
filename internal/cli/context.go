package cli

import (
	"fmt"
	"io"

	"github.com/michalvavra/asncli/internal/asana"
	"github.com/michalvavra/asncli/internal/auth"
	"github.com/michalvavra/asncli/internal/config"
)

type Context struct {
	Stdout        io.Writer
	Stderr        io.Writer
	JSON          bool
	TokenSource   auth.TokenSource
	Store         auth.TokenStore
	Client        any
	ClientFactory func(auth.TokenSource) *asana.Client
	Renderer      Renderer
	Config        *config.Config
}

func (c *Context) ClientOrDefault() any {
	if c.Client != nil {
		return c.Client
	}
	if c.ClientFactory != nil {
		return c.ClientFactory(c.TokenSource)
	}
	return asana.NewClient(c.TokenSource, nil)
}

func (c *Context) RendererOrDefault() Renderer {
	if c.Renderer != nil {
		return c.Renderer
	}
	return NewRenderer(c.Stdout, c.Stderr, c.JSON)
}

// ResolveWorkspace returns the workspace to use, checking in order:
// 1. Explicit flag value (if provided)
// 2. Environment variable ASNCLI_DEFAULT_WORKSPACE
// 3. Config file default_workspace
// Returns error if no workspace is configured.
func (c *Context) ResolveWorkspace(flagValue string) (string, error) {
	cfg := c.Config
	if cfg == nil {
		var err error
		cfg, err = config.Load()
		if err != nil {
			return "", fmt.Errorf("failed to load config: %w", err)
		}
	}

	workspace := config.ResolveWorkspace(flagValue, cfg)
	if workspace == "" {
		return "", fmt.Errorf("workspace is required\nSet a default with: asn config set-workspace\nOr use: --workspace <gid>\nOr set: ASNCLI_DEFAULT_WORKSPACE=<gid>")
	}
	return workspace, nil
}
