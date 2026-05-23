package cmd

import (
	"fmt"
	"os"
	"path/filepath"
)

func resolveSourceDirAbs() (string, error) {
	sourceDir, err := resolveSourceDir()
	if err != nil {
		return "", err
	}
	sourceDir, err = filepath.Abs(sourceDir)
	if err != nil {
		return "", err
	}
	if _, err := os.Stat(sourceDir); err != nil {
		return "", fmt.Errorf("source directory not found: %s (run dot-agent init first)", sourceDir)
	}
	return sourceDir, nil
}
