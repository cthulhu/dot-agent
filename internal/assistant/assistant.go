package assistant

import (
	"os"
	"path/filepath"

	"github.com/cthulhu/dot-agent/internal/config"
)

// Known assistant names supported in v1.
const (
	Claude = "claude"
	Cursor = "cursor"
	Hermes = "hermes"
	Codex  = "codex"
)

// DefaultManifest returns the built-in dot-agent.yaml content for a fresh repo.
func DefaultManifest() *config.Manifest {
	return &config.Manifest{
		Version: 1,
		Assistants: map[string]config.AssistantEntry{
			Claude: DefaultClaude(),
			Cursor: DefaultCursor(),
			Hermes: DefaultHermes(),
			Codex:  DefaultCodex(),
		},
	}
}

func DefaultClaude() config.AssistantEntry {
	return config.AssistantEntry{
		Source: "assistants/claude",
		Target: "~/.claude",
		Ignore: []string{
			"**/.credentials*",
			"**/cache/**",
			"**/debug/**",
			"**/projects/**",
			"**/statsig/**",
			"**/backups/**",
			"**/plugins/marketplaces/**/.git/**",
			"**/*.log",
		},
	}
}

func DefaultCursor() config.AssistantEntry {
	return config.AssistantEntry{
		Source: "assistants/cursor",
		Target: "~/.cursor",
		Ignore: []string{
			"**/projects/**",
			"**/extensions/**",
			"**/agent-transcripts/**",
			"**/terminals/**",
			"prompt_history.json",
			"statsig-cache.json",
			"ide_state.json",
			"agent-cli-state.json",
			"**/*.log",
		},
	}
}

func DefaultHermes() config.AssistantEntry {
	return config.AssistantEntry{
		Source: "assistants/hermes",
		Target: "~/.hermes",
		Ignore: []string{
			".env",
			"auth.json",
			"**/sessions/**",
			"**/logs/**",
			"hermes-agent/**",
			"**/*.log",
		},
	}
}

func DefaultCodex() config.AssistantEntry {
	return config.AssistantEntry{
		Source: "assistants/codex",
		Target: "~/.codex",
		Ignore: []string{
			"auth.json",
			"history.jsonl",
			"**/sessions/**",
			"**/logs/**",
			"**/cache/**",
			"**/*.log",
		},
	}
}

func KnownNames() []string {
	return []string{Claude, Cursor, Hermes, Codex}
}

func KnownNamesString() string {
	return "claude, cursor, hermes, codex"
}

func IsKnown(name string) bool {
	switch name {
	case Claude, Cursor, Hermes, Codex:
		return true
	default:
		return false
	}
}

func WriteDefaultManifest(sourceDir string) error {
	path := filepath.Join(sourceDir, "dot-agent.yaml")
	if _, err := os.Stat(path); err == nil {
		return nil
	}
	return config.WriteManifest(path, DefaultManifest())
}
