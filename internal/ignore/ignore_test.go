package ignore_test

import (
	"path/filepath"
	"testing"

	"github.com/cthulhu/dot-agent/internal/ignore"
)

func TestMatcherIgnored(t *testing.T) {
	m := ignore.New("**/cache/**", "**/*.log")

	cases := []struct {
		path   string
		ignore bool
	}{
		{"cache/foo.txt", true},
		{"rules/my-rule.md", false},
		{"debug/out.log", true},
	}
	for _, tc := range cases {
		got, err := m.Ignored(tc.path)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.ignore {
			t.Fatalf("Ignored(%q) = %v, want %v", tc.path, got, tc.ignore)
		}
	}

	// Paths with OS-native separators must match the same patterns.
	got, err := m.Ignored(filepath.Join("cache", "foo.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if !got {
		t.Fatalf("expected Ignored(%q) = true", filepath.Join("cache", "foo.txt"))
	}
}

func TestBlocked(t *testing.T) {
	if !ignore.Blocked(".env") {
		t.Fatal("expected .env to be blocked")
	}
	if !ignore.Blocked("subdir/credentials.json") {
		t.Fatal("expected credentials.json to be blocked")
	}
	if !ignore.Blocked("auth.json") {
		t.Fatal("expected auth.json to be blocked")
	}
	if ignore.Blocked("settings.json") {
		t.Fatal("did not expect settings.json to be blocked")
	}
}
