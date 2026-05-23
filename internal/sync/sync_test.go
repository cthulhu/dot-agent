package sync_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cthulhu/dot-agent/internal/assistant"
	"github.com/cthulhu/dot-agent/internal/config"
	"github.com/cthulhu/dot-agent/internal/sync"
)

func TestAddAndApply(t *testing.T) {
	root := t.TempDir()
	sourceRoot := filepath.Join(root, "source")
	localClaude := filepath.Join(root, "home", ".claude")
	if err := os.MkdirAll(filepath.Join(localClaude, "rules"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(localClaude, "settings.json"), []byte(`{"theme":"dark"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(localClaude, "rules", "go.md"), []byte("# Go rules"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(localClaude, ".env"), []byte("SECRET=1"), 0o644); err != nil {
		t.Fatal(err)
	}

	entry := assistant.DefaultClaude()
	entry.Target = localClaude
	if err := os.MkdirAll(filepath.Join(sourceRoot, entry.Source), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := config.WriteManifest(filepath.Join(sourceRoot, "dot-agent.yaml"), assistant.DefaultManifest()); err != nil {
		t.Fatal(err)
	}

	result, err := sync.Add(sourceRoot, entry, sync.Options{})
	if err != nil {
		t.Fatal(err)
	}
	foundBlocked := false
	for _, a := range result.Actions {
		if a.Action == "blocked" && a.RelPath == ".env" {
			foundBlocked = true
		}
	}
	if !foundBlocked {
		t.Fatal("expected .env to be blocked during add")
	}

	settingsPath := filepath.Join(sourceRoot, entry.Source, "settings.json")
	if _, err := os.Stat(settingsPath); err != nil {
		t.Fatalf("expected settings.json in source: %v", err)
	}

	applyHome := filepath.Join(root, "apply-home", ".claude")
	entry.Target = applyHome
	_, err = sync.Apply(sourceRoot, entry, sync.Options{Force: true})
	if err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(filepath.Join(applyHome, "settings.json"))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != `{"theme":"dark"}` {
		t.Fatalf("unexpected settings content: %q", got)
	}
}

func TestCompareDrift(t *testing.T) {
	root := t.TempDir()
	sourceRoot := filepath.Join(root, "source")
	local := filepath.Join(root, "local", ".claude")
	entry := config.AssistantEntry{
		Source: "assistants/claude",
		Target: local,
		Ignore: assistant.DefaultClaude().Ignore,
	}
	srcDir := filepath.Join(sourceRoot, entry.Source)
	if err := os.MkdirAll(srcDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(local, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(srcDir, "a.txt"), []byte("source"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(local, "b.txt"), []byte("local"), 0o644); err != nil {
		t.Fatal(err)
	}

	report, err := sync.Compare(sourceRoot, entry)
	if err != nil {
		t.Fatal(err)
	}
	if len(report.Entries) != 2 {
		t.Fatalf("expected 2 drift entries, got %d", len(report.Entries))
	}
}
