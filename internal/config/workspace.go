package config

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// WorkspaceLister can list workspaces for the current user.
type WorkspaceLister interface {
	ListWorkspaceMembershipsForUser(ctx context.Context, userGID string) (WorkspaceMembershipList, error)
}

// WorkspaceMembership represents a workspace membership.
type WorkspaceMembership struct {
	GID        string
	Name       string
	IsActive   bool
	IsAdmin    bool
	IsGuest    bool
	IsViewOnly bool
}

// WorkspaceMembershipList is a list of workspace memberships.
type WorkspaceMembershipList struct {
	Data []WorkspaceMembership
}

// ResolveWorkspace determines which workspace to use based on:
// 1. Explicit flag value (if provided)
// 2. Environment variable ASNCLI_DEFAULT_WORKSPACE
// 3. Config file default_workspace
// Returns empty string if no workspace is configured.
func ResolveWorkspace(flagValue string, cfg *Config) string {
	if flagValue != "" {
		return flagValue
	}
	// Environment variable is already merged into cfg by Load()
	return cfg.DefaultWorkspace
}

// PromptWorkspace shows a list of workspaces and prompts the user to select one.
func PromptWorkspace(ctx context.Context, lister WorkspaceLister, stdout, stderr io.Writer) (string, error) {
	// Get workspace memberships
	memberships, err := lister.ListWorkspaceMembershipsForUser(ctx, "me")
	if err != nil {
		return "", fmt.Errorf("failed to list workspaces: %w", err)
	}

	if len(memberships.Data) == 0 {
		return "", fmt.Errorf("no workspaces found")
	}

	// If only one workspace, auto-select it
	if len(memberships.Data) == 1 {
		workspace := memberships.Data[0]
		fmt.Fprintf(stdout, "Using your only workspace: %s (%s)\n", workspace.Name, workspace.GID)
		return workspace.GID, nil
	}

	// Show workspace list
	fmt.Fprintln(stdout, "Select default workspace:")
	for i, membership := range memberships.Data {
		status := ""
		if !membership.IsActive {
			status = " (inactive)"
		}
		fmt.Fprintf(stdout, "  %d. %s (%s)%s\n", i+1, membership.Name, membership.GID, status)
	}

	// Prompt for selection
	fmt.Fprint(stderr, "Choice: ")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read input: %w", err)
	}

	// Parse selection
	input = strings.TrimSpace(input)
	choice, err := strconv.Atoi(input)
	if err != nil {
		return "", fmt.Errorf("invalid choice: %s", input)
	}

	if choice < 1 || choice > len(memberships.Data) {
		return "", fmt.Errorf("choice out of range: %d", choice)
	}

	selected := memberships.Data[choice-1]
	return selected.GID, nil
}
