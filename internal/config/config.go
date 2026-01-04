package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	envDefaultWorkspace = "ASNCLI_DEFAULT_WORKSPACE"
	configDirName       = "asncli"
	configFileName      = "config.json"
)

// Config holds application configuration.
type Config struct {
	DefaultWorkspace     string `json:"default_workspace,omitempty"`
	DefaultWorkspaceName string `json:"default_workspace_name,omitempty"`
}

// Options for loading configuration.
type Options struct {
	// ConfigPath overrides the default config file path (useful for testing).
	ConfigPath string
}

// Load loads configuration from file and environment variables.
// Environment variables take precedence over file settings.
// Returns empty config if file doesn't exist (not an error).
func Load() (*Config, error) {
	return LoadWithOptions(Options{})
}

// LoadWithOptions loads configuration with custom options.
func LoadWithOptions(opts Options) (*Config, error) {
	cfg := &Config{}

	// Get config file path
	path := opts.ConfigPath
	if path == "" {
		var err error
		path, err = GetConfigPath()
		if err != nil {
			return nil, fmt.Errorf("failed to get config path: %w", err)
		}
	}

	// Read config file if it exists
	data, err := os.ReadFile(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
		// File doesn't exist, use empty config
	} else {
		if err := json.Unmarshal(data, cfg); err != nil {
			return nil, fmt.Errorf("failed to parse config file: %w", err)
		}
	}

	// Environment variables override file settings
	if workspace := os.Getenv(envDefaultWorkspace); workspace != "" {
		cfg.DefaultWorkspace = workspace
	}

	return cfg, nil
}

// Save saves configuration to the config file.
func Save(cfg *Config) error {
	return SaveWithPath(cfg, "")
}

// SaveWithPath saves configuration to a specific path.
// If path is empty, uses the default config path.
func SaveWithPath(cfg *Config, path string) error {
	if path == "" {
		var err error
		path, err = GetConfigPath()
		if err != nil {
			return fmt.Errorf("failed to get config path: %w", err)
		}
	}

	// Create config directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Marshal config to JSON
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to file
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GetConfigPath returns the path to the config file.
func GetConfigPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, configDirName, configFileName), nil
}
