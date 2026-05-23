package paths_test

import (
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/cthulhu/dot-agent/internal/paths"
)

func TestExpandPathHome(t *testing.T) {
	home, err := paths.HomeDir()
	if err != nil {
		t.Fatal(err)
	}

	got, err := paths.ExpandPath("~/.claude")
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(home, ".claude")
	if got != want {
		t.Fatalf("ExpandPath(~/.claude) = %q, want %q", got, want)
	}
}

func TestDefaultSourceDir(t *testing.T) {
	dir, err := paths.DefaultSourceDir()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(dir, "dot-agent") {
		t.Fatalf("expected dot-agent in path, got %q", dir)
	}
	if runtime.GOOS == "windows" {
		if !strings.Contains(strings.ToLower(dir), "local") {
			t.Fatalf("expected LOCALAPPDATA-style path on windows, got %q", dir)
		}
	}
}

func TestConfigDir(t *testing.T) {
	dir, err := paths.ConfigDir()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(dir, "dot-agent") {
		t.Fatalf("expected dot-agent in config path, got %q", dir)
	}
}
