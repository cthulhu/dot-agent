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
	Codex   = "codex"
	Gemini  = "gemini"
	Copilot = "copilot"
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
			Gemini:  DefaultGemini(),
			Copilot: DefaultCopilot(),
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

func DefaultGemini() config.AssistantEntry {
	return config.AssistantEntry{
		Source: "assistants/gemini",
		Target: "~/.gemini",
		Ignore: []string{
			".env",
			"**/.env",
			"oauth_creds.json",
			"mcp-oauth-tokens.json",
			"a2a-oauth-tokens.json",
			"google_accounts.json",
			"installation_id",
			"policy_integrity.json",
			"trustedFolders.json",
			"projects.json",
			"**/tmp/**",
			"**/history/**",
			"**/logs/**",
			"**/chats/**",
			"**/checkpoints/**",
			"**/*.log",
		},
	}
}

func DefaultCopilot() config.AssistantEntry {
	return config.AssistantEntry{
		Source: "assistants/copilot",
		Target: "~/.copilot",
		Ignore: []string{
			".env",
			"**/.env",
			"config.json",
			"session-store.db",
			"**/session-state/**",
			"**/command-history-state/**",
			"**/logs/**",
			"**/ide/**",
			"**/plugin-data/**",
			"**/*.log",
		},
	}
}

func KnownNames() []string {
	return []string{Claude, Cursor, Hermes, Codex, Gemini, Copilot}
}

func KnownNamesString() string {
	return "claude, cursor, hermes, codex, gemini, copilot"
}

func IsKnown(name string) bool {
	switch name {
	case Claude, Cursor, Hermes, Codex, Gemini, Copilot:
		return true
	default:
		return false
	}
}

func DefaultEntry(name string) (config.AssistantEntry, bool) {
	switch name {
	case Claude:
		return DefaultClaude(), true
	case Cursor:
		return DefaultCursor(), true
	case Hermes:
		return DefaultHermes(), true
	case Codex:
		return DefaultCodex(), true
	case Gemini:
		return DefaultGemini(), true
	case Copilot:
		return DefaultCopilot(), true
	default:
		return config.AssistantEntry{}, false
	}
}

// MergeMissingAssistants adds built-in defaults for any known assistant not yet in the manifest.
func MergeMissingAssistants(m *config.Manifest) []string {
	if m.Assistants == nil {
		m.Assistants = make(map[string]config.AssistantEntry)
	}
	var added []string
	for _, name := range KnownNames() {
		if _, ok := m.Assistants[name]; ok {
			continue
		}
		entry, ok := DefaultEntry(name)
		if !ok {
			continue
		}
		m.Assistants[name] = entry
		added = append(added, name)
	}
	return added
}

func WriteDefaultManifest(sourceDir string) error {
	path := filepath.Join(sourceDir, "dot-agent.yaml")
	if _, err := os.Stat(path); err == nil {
		return nil
	}
	return config.WriteManifest(path, DefaultManifest())
}
