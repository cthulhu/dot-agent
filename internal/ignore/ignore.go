package ignore

import (
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
)

var blockedBasenames = map[string]struct{}{
	".env":             {},
	"auth.json":        {},
	"credentials.json": {},
	"credentials.yaml":  {},
	"credentials.yml":   {},
	".netrc":            {},
	"id_rsa":            {},
	"id_ed25519":        {},
	".aws/credentials":  {},
}

// Matcher evaluates glob patterns relative to a root directory.
type Matcher struct {
	patterns []string
}

func New(patterns ...string) *Matcher {
	return &Matcher{patterns: patterns}
}

func (m *Matcher) Ignored(relPath string) (bool, error) {
	rel := filepath.ToSlash(relPath)
	for _, pattern := range m.patterns {
		match, err := doublestar.PathMatch(pattern, rel)
		if err != nil {
			return false, err
		}
		if match {
			return true, nil
		}
	}
	return false, nil
}

// Blocked reports whether a relative path must never be synced (secrets).
func Blocked(relPath string) bool {
	rel := filepath.ToSlash(strings.TrimPrefix(relPath, "./"))
	base := filepath.Base(rel)
	if _, ok := blockedBasenames[base]; ok {
		return true
	}
	lower := strings.ToLower(rel)
	for blocked := range blockedBasenames {
		if strings.HasSuffix(lower, "/"+blocked) || lower == blocked {
			return true
		}
	}
	if strings.Contains(lower, "secret") && (strings.HasSuffix(lower, ".json") || strings.HasSuffix(lower, ".yaml")) {
		return true
	}
	return false
}
