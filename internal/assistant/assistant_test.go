package assistant_test

import (
	"testing"

	"github.com/cthulhu/dot-agent/internal/assistant"
	"github.com/cthulhu/dot-agent/internal/config"
)

func TestDefaultHermesIgnoresSecrets(t *testing.T) {
	entry := assistant.DefaultHermes()
	foundEnv := false
	foundAuth := false
	for _, p := range entry.Ignore {
		if p == ".env" {
			foundEnv = true
		}
		if p == "auth.json" {
			foundAuth = true
		}
	}
	if !foundEnv || !foundAuth {
		t.Fatalf("expected hermes defaults to ignore .env and auth.json, got %v", entry.Ignore)
	}
}

func TestDefaultCodexIgnoresSecrets(t *testing.T) {
	entry := assistant.DefaultCodex()
	foundAuth := false
	foundHistory := false
	for _, p := range entry.Ignore {
		if p == "auth.json" {
			foundAuth = true
		}
		if p == "history.jsonl" {
			foundHistory = true
		}
	}
	if !foundAuth || !foundHistory {
		t.Fatalf("expected codex defaults to ignore auth.json and history.jsonl, got %v", entry.Ignore)
	}
	if entry.Target != "~/.codex" {
		t.Fatalf("expected codex target ~/.codex, got %q", entry.Target)
	}
}

func TestDefaultGeminiIgnoresSecrets(t *testing.T) {
	entry := assistant.DefaultGemini()
	foundEnv := false
	foundOAuth := false
	foundTmp := false
	for _, p := range entry.Ignore {
		if p == ".env" {
			foundEnv = true
		}
		if p == "oauth_creds.json" {
			foundOAuth = true
		}
		if p == "**/tmp/**" {
			foundTmp = true
		}
	}
	if !foundEnv || !foundOAuth || !foundTmp {
		t.Fatalf("expected gemini defaults to ignore secrets and tmp, got %v", entry.Ignore)
	}
	if entry.Target != "~/.gemini" {
		t.Fatalf("expected gemini target ~/.gemini, got %q", entry.Target)
	}
}

func TestKnownNamesIncludesAllAssistants(t *testing.T) {
	names := assistant.KnownNames()
	want := map[string]bool{
		assistant.Claude: false,
		assistant.Cursor: false,
		assistant.Hermes: false,
		assistant.Codex:  false,
		assistant.Gemini: false,
	}
	for _, n := range names {
		if _, ok := want[n]; ok {
			want[n] = true
		}
	}
	for name, found := range want {
		if !found {
			t.Fatalf("expected %q in KnownNames, got %v", name, names)
		}
	}
}

func TestMergeMissingAssistants(t *testing.T) {
	m := &config.Manifest{
		Version: 1,
		Assistants: map[string]config.AssistantEntry{
			assistant.Claude: assistant.DefaultClaude(),
		},
	}
	added := assistant.MergeMissingAssistants(m)
	if len(added) != 4 {
		t.Fatalf("expected 4 assistants added, got %v", added)
	}
	if _, ok := m.Assistants[assistant.Gemini]; !ok {
		t.Fatal("expected gemini in manifest after merge")
	}
	if again := assistant.MergeMissingAssistants(m); len(again) != 0 {
		t.Fatalf("expected no assistants added on second merge, got %v", again)
	}
}
