package assistant

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/cthulhu/dot-agent/internal/config"
)

// Known assistant names supported in v1.
const (
	Claude      = "claude"
	Cursor      = "cursor"
	Hermes      = "hermes"
	Codex       = "codex"
	Gemini      = "gemini"
	Copilot     = "copilot"
	Antigravity = "antigravity"
)

type assistantInfo struct {
	displayName string
	defaultFn   func() config.AssistantEntry
}

var registry = map[string]assistantInfo{
	Claude:      {displayName: "Claude Code", defaultFn: DefaultClaude},
	Cursor:      {displayName: "Cursor", defaultFn: DefaultCursor},
	Hermes:      {displayName: "Hermes Agent", defaultFn: DefaultHermes},
	Codex:       {displayName: "OpenAI Codex", defaultFn: DefaultCodex},
	Gemini:      {displayName: "Gemini CLI", defaultFn: DefaultGemini},
	Copilot:     {displayName: "GitHub Copilot CLI", defaultFn: DefaultCopilot},
	Antigravity: {displayName: "Antigravity", defaultFn: DefaultAntigravity},
}

// DefaultManifest returns the built-in dot-agent.yaml content for a fresh repo.
func DefaultManifest() *config.Manifest {
	m := &config.Manifest{
		Version:    1,
		Assistants: make(map[string]config.AssistantEntry),
	}
	for name, info := range registry {
		m.Assistants[name] = info.defaultFn()
	}
	return m
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

func DefaultAntigravity() config.AssistantEntry {
	return config.AssistantEntry{
		Source: "assistants/antigravity",
		Target: "~/.gemini/antigravity-cli",
		Ignore: []string{
			".env",
			"**/.env",
			"installation_id",
			"last_check.timestamp",
			"cli.log",
			"history.jsonl",
			"**/bin/**",
			"**/brain/**",
			"**/cache/**",
			"**/conversations/**",
			"**/implicit/**",
			"**/knowledge/**",
			"**/log/**",
			"**/updater/**",
			"**/*.log",
			"**/*.lock",
		},
	}
}

func KnownNames() []string {
	names := make([]string, 0, len(registry))
	for name := range registry {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func KnownNamesString() string {
	return strings.Join(KnownNames(), ", ")
}

func DisplayNamesString() string {
	names := make([]string, 0, len(registry))
	for _, info := range registry {
		names = append(names, info.displayName)
	}
	sort.Strings(names)
	return strings.Join(names, ", ")
}

func IsKnown(name string) bool {
	_, ok := registry[name]
	return ok
}

func DefaultEntry(name string) (config.AssistantEntry, bool) {
	info, ok := registry[name]
	if !ok {
		return config.AssistantEntry{}, false
	}
	return info.defaultFn(), true
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
