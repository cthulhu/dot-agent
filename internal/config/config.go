package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Manifest struct {
	Version    int                       `yaml:"version"`
	Assistants map[string]AssistantEntry `yaml:"assistants"`
}

type AssistantEntry struct {
	Source string   `yaml:"source"`
	Target string   `yaml:"target"`
	Ignore []string `yaml:"ignore,omitempty"`
}

func LoadManifest(path string) (*Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var m Manifest
	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("parse manifest: %w", err)
	}
	return &m, Validate(m)
}

func Validate(m Manifest) error {
	if m.Version != 1 {
		return fmt.Errorf("unsupported manifest version %d (expected 1)", m.Version)
	}
	if len(m.Assistants) == 0 {
		return fmt.Errorf("manifest has no assistants defined")
	}
	for name, entry := range m.Assistants {
		if entry.Source == "" {
			return fmt.Errorf("assistant %q: source path is required", name)
		}
		if entry.Target == "" {
			return fmt.Errorf("assistant %q: target path is required", name)
		}
	}
	return nil
}

func WriteManifest(path string, m *Manifest) error {
	data, err := yaml.Marshal(m)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

func (m *Manifest) ResolveAssistant(name string) (AssistantEntry, error) {
	entry, ok := m.Assistants[name]
	if !ok {
		return AssistantEntry{}, fmt.Errorf("unknown assistant %q (known: claude, cursor, hermes, codex)", name)
	}
	return entry, nil
}

func (m *Manifest) AssistantNames(filter []string) ([]string, error) {
	if len(filter) == 0 {
		names := make([]string, 0, len(m.Assistants))
		for name := range m.Assistants {
			names = append(names, name)
		}
		return names, nil
	}
	for _, name := range filter {
		if _, ok := m.Assistants[name]; !ok {
			return nil, fmt.Errorf("assistant %q not in manifest", name)
		}
	}
	return filter, nil
}
