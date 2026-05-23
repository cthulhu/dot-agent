package paths

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	ConfigDirName  = "dot-agent"
	ManifestName   = "dot-agent.yaml"
	UserConfigName = "config.yaml"
)

// UserConfig stores local dot-agent settings (source repo path, remote URL).
type UserConfig struct {
	SourcePath string `yaml:"source_path"`
	RemoteURL  string `yaml:"remote_url,omitempty"`
}

func HomeDir() (string, error) {
	return os.UserHomeDir()
}

func ConfigDir() (string, error) {
	if runtime.GOOS == "windows" {
		dir := os.Getenv("APPDATA")
		if dir == "" {
			home, err := HomeDir()
			if err != nil {
				return "", err
			}
			dir = filepath.Join(home, "AppData", "Roaming")
		}
		return filepath.Join(dir, ConfigDirName), nil
	}
	home, err := HomeDir()
	if err != nil {
		return "", err
	}
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, ConfigDirName), nil
	}
	return filepath.Join(home, ".config", ConfigDirName), nil
}

func DefaultSourceDir() (string, error) {
	if runtime.GOOS == "windows" {
		dir := os.Getenv("LOCALAPPDATA")
		if dir == "" {
			home, err := HomeDir()
			if err != nil {
				return "", err
			}
			dir = filepath.Join(home, "AppData", "Local")
		}
		return filepath.Join(dir, ConfigDirName, "source"), nil
	}
	home, err := HomeDir()
	if err != nil {
		return "", err
	}
	if xdg := os.Getenv("XDG_DATA_HOME"); xdg != "" {
		return filepath.Join(xdg, ConfigDirName, "source"), nil
	}
	return filepath.Join(home, ".local", "share", ConfigDirName, "source"), nil
}

// ExpandPath resolves ~ and environment variables in a path string.
func ExpandPath(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("empty path")
	}
	if strings.HasPrefix(path, "~") {
		home, err := HomeDir()
		if err != nil {
			return "", err
		}
		if path == "~" {
			return home, nil
		}
		if strings.HasPrefix(path, "~/") || strings.HasPrefix(path, "~\\") {
			return filepath.Join(home, path[2:]), nil
		}
		return "", fmt.Errorf("invalid home path: %q", path)
	}
	return filepath.Clean(os.ExpandEnv(path)), nil
}

func UserConfigPath() (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, UserConfigName), nil
}

func LoadUserConfig() (*UserConfig, error) {
	path, err := UserConfigPath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var cfg UserConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	return &cfg, nil
}

func SaveUserConfig(cfg *UserConfig) error {
	path, err := UserConfigPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

func ManifestPath(sourceDir string) string {
	return filepath.Join(sourceDir, ManifestName)
}
