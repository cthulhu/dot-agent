package cmd

import (
	"fmt"
	"os"

	"github.com/cthulhu/dot-agent/internal/config"
	"github.com/cthulhu/dot-agent/internal/paths"
	"github.com/spf13/cobra"
)

var (
	sourceFlag string
	version    = "0.1.0"
)

var rootCmd = &cobra.Command{
	Use:   "dot-agent",
	Short: "Sync AI coding assistant configs via git",
	Long:  "dot-agent manages Claude Code, Cursor, and Hermes Agent configuration in a git repo and applies it across machines.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&sourceFlag, "source", "", "path to dot-agent source git repo (default: from config or OS default)")
	rootCmd.SetVersionTemplate("dot-agent {{.Version}}\n")
	rootCmd.Version = version
}

func resolveSourceDir() (string, error) {
	if sourceFlag != "" {
		return paths.ExpandPath(sourceFlag)
	}
	userCfg, err := paths.LoadUserConfig()
	if err != nil {
		return "", err
	}
	if userCfg != nil && userCfg.SourcePath != "" {
		return paths.ExpandPath(userCfg.SourcePath)
	}
	return paths.DefaultSourceDir()
}

func loadManifest() (*config.Manifest, string, error) {
	sourceDir, err := resolveSourceDir()
	if err != nil {
		return nil, "", err
	}
	manifestPath := paths.ManifestPath(sourceDir)
	m, err := config.LoadManifest(manifestPath)
	if err != nil {
		return nil, sourceDir, fmt.Errorf("load manifest at %s: %w (run dot-agent init first)", manifestPath, err)
	}
	return m, sourceDir, nil
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(1)
}
