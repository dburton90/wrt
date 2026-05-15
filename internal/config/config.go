package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

type file struct {
	TaskRootPath     string `toml:"task_root"`
	UsernameOverride string `toml:"username"`
}

type Config struct {
	f file
}

var loaded *Config

func Load() (*Config, error) {
	if loaded != nil {
		return loaded, nil
	}

	cfgPath := filepath.Join(os.Getenv("HOME"), ".config", "wrt", "config.toml")

	var f file
	if _, err := os.Stat(cfgPath); err == nil {
		if _, err := toml.DecodeFile(cfgPath, &f); err != nil {
			return nil, fmt.Errorf("reading config %s: %w", cfgPath, err)
		}
	}

	loaded = &Config{f: f}
	return loaded, nil
}

func (c *Config) TaskRoot() (string, error) {
	if c.f.TaskRootPath == "" {
		return "", fmt.Errorf(
			"task_root is not configured\n\nCreate ~/.config/wrt/config.toml with:\n  task_root = \"/path/to/your/tasks\"",
		)
	}
	return c.f.TaskRootPath, nil
}

func (c *Config) Username() string {
	u := c.f.UsernameOverride
	if u == "" {
		u = os.Getenv("USER")
	}
	return sanitize(u)
}

func sanitize(u string) string {
	return strings.ToLower(strings.ReplaceAll(u, " ", ""))
}
