package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ConfluenceURL   string `yaml:"confluence_url"`
	Email           string `yaml:"email"`
	APIToken        string `yaml:"api_token"`
}

func Load() (*Config, error) {
	cfg := &Config{}

	// Try to load from XDG config file first
	configPath := getConfigPath()
	if _, err := os.Stat(configPath); err == nil {
		data, err := os.ReadFile(configPath)
		if err != nil {
			return nil, fmt.Errorf("reading config file: %w", err)
		}

		if err := yaml.Unmarshal(data, cfg); err != nil {
			return nil, fmt.Errorf("parsing config file: %w", err)
		}
	}

	// Environment variables override config file
	if url := os.Getenv("CONFLUENCE_URL"); url != "" {
		cfg.ConfluenceURL = url
	}
	if email := os.Getenv("CONFLUENCE_EMAIL"); email != "" {
		cfg.Email = email
	}
	if token := os.Getenv("CONFLUENCE_API_TOKEN"); token != "" {
		cfg.APIToken = token
	}

	// Validate required fields
	if cfg.ConfluenceURL == "" {
		return nil, fmt.Errorf("confluence_url not set (check config file or CONFLUENCE_URL env var)")
	}
	if cfg.Email == "" {
		return nil, fmt.Errorf("email not set (check config file or CONFLUENCE_EMAIL env var)")
	}
	if cfg.APIToken == "" {
		return nil, fmt.Errorf("api_token not set (check config file or CONFLUENCE_API_TOKEN env var)")
	}

	return cfg, nil
}

func getConfigPath() string {
	configDir := os.Getenv("XDG_CONFIG_HOME")
	if configDir == "" {
		home, _ := os.UserHomeDir()
		configDir = filepath.Join(home, ".config")
	}
	return filepath.Join(configDir, "confluence-md", "config.yaml")
}
