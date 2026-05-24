package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	DefaultHost     string
	DefaultInsecure bool
	WorkflowsDir    string
}

func DefaultConfig() *Config {
	home, _ := os.UserHomeDir()
	return &Config{
		DefaultHost:     "https://localhost:8006",
		DefaultInsecure: false,
		WorkflowsDir:    filepath.Join(home, ".pmgwire", "workflows"),
	}
}

func (c *Config) EnsureDirs() error {
	dirs := []string{
		filepath.Dir(c.WorkflowsDir),
		c.WorkflowsDir,
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	return nil
}