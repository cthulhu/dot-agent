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

func TestKnownNamesIncludesHermes(t *testing.T) {
	names := assistant.KnownNames()
	found := false
	for _, n := range names {
		if n == assistant.Hermes {
			found = true
		}
	}
	if !found {
		t.Fatal("expected hermes in KnownNames")
	}
}
