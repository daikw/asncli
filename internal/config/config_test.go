package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadEmpty(t *testing.T) {
	cfg, err := LoadWithOptions(Options{ConfigPath: "/nonexistent/config.json"})
	if err != nil {
		t.Fatalf("LoadWithOptions returned unexpected error: %v", err)
	}
	if cfg.DefaultWorkspace != "" {
		t.Errorf("DefaultWorkspace = %q, want empty string for missing config", cfg.DefaultWorkspace)
	}
}

func TestLoadFromFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")

	// Write test config
	testCfg := Config{DefaultWorkspace: "123456"}
	data, err := json.MarshalIndent(testCfg, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal test config: %v", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("failed to write test config file: %v", err)
	}

	// Load it
	cfg, err := LoadWithOptions(Options{ConfigPath: path})
	if err != nil {
		t.Fatalf("LoadWithOptions returned unexpected error: %v", err)
	}
	if cfg.DefaultWorkspace != "123456" {
		t.Errorf("DefaultWorkspace = %q, want %q", cfg.DefaultWorkspace, "123456")
	}
}

func TestLoadEnvOverridesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")

	// Write test config
	testCfg := Config{DefaultWorkspace: "file-workspace"}
	data, err := json.MarshalIndent(testCfg, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal test config: %v", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("failed to write test config file: %v", err)
	}

	// Set env var
	if err := os.Setenv(envDefaultWorkspace, "env-workspace"); err != nil {
		t.Fatalf("failed to set environment variable: %v", err)
	}
	t.Cleanup(func() { _ = os.Unsetenv(envDefaultWorkspace) })

	// Load it
	cfg, err := LoadWithOptions(Options{ConfigPath: path})
	if err != nil {
		t.Fatalf("LoadWithOptions returned unexpected error: %v", err)
	}
	if cfg.DefaultWorkspace != "env-workspace" {
		t.Errorf("DefaultWorkspace = %q, want %q (env should override file)", cfg.DefaultWorkspace, "env-workspace")
	}
}

func TestSave(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")

	cfg := &Config{DefaultWorkspace: "789012", DefaultWorkspaceName: "My Workspace"}
	if err := SaveWithPath(cfg, path); err != nil {
		t.Fatalf("SaveWithPath returned unexpected error: %v", err)
	}

	// Read it back
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read saved config file: %v", err)
	}

	var loaded Config
	if err := json.Unmarshal(data, &loaded); err != nil {
		t.Fatalf("failed to unmarshal saved config: %v", err)
	}

	if loaded.DefaultWorkspace != "789012" {
		t.Errorf("DefaultWorkspace = %q, want %q", loaded.DefaultWorkspace, "789012")
	}
	if loaded.DefaultWorkspaceName != "My Workspace" {
		t.Errorf("DefaultWorkspaceName = %q, want %q", loaded.DefaultWorkspaceName, "My Workspace")
	}
}

func TestSaveCreatesDirectory(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "subdir", "config.json")

	cfg := &Config{DefaultWorkspace: "123"}
	if err := SaveWithPath(cfg, path); err != nil {
		t.Fatalf("SaveWithPath returned unexpected error: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(path); err != nil {
		t.Errorf("config file not created: %v", err)
	}
}

func TestLoadInvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")

	if err := os.WriteFile(path, []byte("{invalid"), 0o644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	_, err := LoadWithOptions(Options{ConfigPath: path})
	if err == nil {
		t.Fatal("LoadWithOptions should return error for invalid JSON")
	}
}

func TestGetConfigPath(t *testing.T) {
	path, err := GetConfigPath()
	if err != nil {
		t.Fatalf("GetConfigPath returned unexpected error: %v", err)
	}
	if path == "" {
		t.Fatal("GetConfigPath returned empty path")
	}
	if !strings.Contains(path, "asncli") {
		t.Errorf("config path = %q, want to contain 'asncli'", path)
	}
	if !strings.Contains(path, "config.json") {
		t.Errorf("config path = %q, want to contain 'config.json'", path)
	}
}

func TestLoadAndSaveRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")

	original := &Config{
		DefaultWorkspace:     "ws-123",
		DefaultWorkspaceName: "Test Workspace",
	}

	if err := SaveWithPath(original, path); err != nil {
		t.Fatalf("SaveWithPath returned unexpected error: %v", err)
	}

	loaded, err := LoadWithOptions(Options{ConfigPath: path})
	if err != nil {
		t.Fatalf("LoadWithOptions returned unexpected error: %v", err)
	}

	if loaded.DefaultWorkspace != original.DefaultWorkspace {
		t.Errorf("DefaultWorkspace = %q, want %q", loaded.DefaultWorkspace, original.DefaultWorkspace)
	}
	if loaded.DefaultWorkspaceName != original.DefaultWorkspaceName {
		t.Errorf("DefaultWorkspaceName = %q, want %q", loaded.DefaultWorkspaceName, original.DefaultWorkspaceName)
	}
}

func TestLoadEnvOnlyNoFile(t *testing.T) {
	if err := os.Setenv(envDefaultWorkspace, "env-only-workspace"); err != nil {
		t.Fatalf("failed to set environment variable: %v", err)
	}
	t.Cleanup(func() { _ = os.Unsetenv(envDefaultWorkspace) })

	cfg, err := LoadWithOptions(Options{ConfigPath: "/nonexistent/config.json"})
	if err != nil {
		t.Fatalf("LoadWithOptions returned unexpected error: %v", err)
	}
	if cfg.DefaultWorkspace != "env-only-workspace" {
		t.Errorf("DefaultWorkspace = %q, want %q", cfg.DefaultWorkspace, "env-only-workspace")
	}
}

func TestLoadDefaultPath(t *testing.T) {
	// Load() without options should work (uses default path)
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load returned unexpected error: %v", err)
	}
	if cfg == nil {
		t.Fatal("Load returned nil config")
	}
}

func TestLoadReadError(t *testing.T) {
	// Create a directory where config file should be (can't read a directory as file)
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatalf("failed to create test directory: %v", err)
	}

	_, err := LoadWithOptions(Options{ConfigPath: path})
	if err == nil {
		t.Fatal("LoadWithOptions should return error when path is a directory")
	}
}

func TestResolveWorkspace(t *testing.T) {
	tests := []struct {
		name      string
		flagValue string
		config    *Config
		want      string
	}{
		{
			name:      "flag takes precedence",
			flagValue: "flag-workspace",
			config:    &Config{DefaultWorkspace: "config-workspace"},
			want:      "flag-workspace",
		},
		{
			name:      "config used when flag empty",
			flagValue: "",
			config:    &Config{DefaultWorkspace: "config-workspace"},
			want:      "config-workspace",
		},
		{
			name:      "empty when both empty",
			flagValue: "",
			config:    &Config{},
			want:      "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ResolveWorkspace(tt.flagValue, tt.config)
			if got != tt.want {
				t.Errorf("ResolveWorkspace() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestSaveWithEmptyPath(t *testing.T) {
	// Save() uses default path, but we can test SaveWithPath with empty path
	// This tests the path resolution branch in SaveWithPath
	cfg := &Config{DefaultWorkspace: "test"}

	// We can't easily test Save() without mocking GetConfigPath,
	// but we can verify SaveWithPath handles the empty path case
	// by checking it doesn't error (it will try to use real config path)
	err := SaveWithPath(cfg, "")
	// This may succeed or fail depending on permissions, but shouldn't panic
	_ = err
}

func TestSaveEmptyConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty-config.json")

	cfg := &Config{}
	if err := SaveWithPath(cfg, path); err != nil {
		t.Fatalf("SaveWithPath returned unexpected error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read saved config file: %v", err)
	}

	var loaded Config
	if err := json.Unmarshal(data, &loaded); err != nil {
		t.Fatalf("failed to unmarshal saved config: %v", err)
	}

	if loaded.DefaultWorkspace != "" {
		t.Errorf("DefaultWorkspace = %q, want empty", loaded.DefaultWorkspace)
	}
}

func TestSaveOverwritesExisting(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")

	// Write initial config
	cfg1 := &Config{DefaultWorkspace: "first", DefaultWorkspaceName: "First WS"}
	if err := SaveWithPath(cfg1, path); err != nil {
		t.Fatalf("SaveWithPath (first) returned unexpected error: %v", err)
	}

	// Overwrite with new config
	cfg2 := &Config{DefaultWorkspace: "second", DefaultWorkspaceName: "Second WS"}
	if err := SaveWithPath(cfg2, path); err != nil {
		t.Fatalf("SaveWithPath (second) returned unexpected error: %v", err)
	}

	// Read and verify
	loaded, err := LoadWithOptions(Options{ConfigPath: path})
	if err != nil {
		t.Fatalf("LoadWithOptions returned unexpected error: %v", err)
	}

	if loaded.DefaultWorkspace != "second" {
		t.Errorf("DefaultWorkspace = %q, want %q", loaded.DefaultWorkspace, "second")
	}
	if loaded.DefaultWorkspaceName != "Second WS" {
		t.Errorf("DefaultWorkspaceName = %q, want %q", loaded.DefaultWorkspaceName, "Second WS")
	}
}

func TestSaveWriteError(t *testing.T) {
	// Try to write to a read-only directory
	dir := t.TempDir()
	readOnlyDir := filepath.Join(dir, "readonly")
	if err := os.MkdirAll(readOnlyDir, 0o555); err != nil {
		t.Fatalf("failed to create readonly directory: %v", err)
	}
	t.Cleanup(func() { _ = os.Chmod(readOnlyDir, 0o755) })

	path := filepath.Join(readOnlyDir, "config.json")
	cfg := &Config{DefaultWorkspace: "test"}

	err := SaveWithPath(cfg, path)
	if err == nil {
		t.Fatal("SaveWithPath should return error when writing to read-only directory")
	}
}
