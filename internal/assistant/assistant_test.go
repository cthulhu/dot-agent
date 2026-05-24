package assistant_test

import (
	"testing"

	"github.com/cthulhu/dot-agent/internal/assistant"
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

func TestKnownNamesIncludesAllAssistants(t *testing.T) {
	names := assistant.KnownNames()
	want := map[string]bool{
		assistant.Claude: false,
		assistant.Cursor: false,
		assistant.Hermes: false,
		assistant.Codex:  false,
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
