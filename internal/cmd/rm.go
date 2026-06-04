package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cthulhu/dot-agent/internal/config"
	"github.com/cthulhu/dot-agent/internal/git"
	"github.com/cthulhu/dot-agent/internal/paths"
	"github.com/spf13/cobra"
)

var rmDryRun bool

var rmCmd = &cobra.Command{
	Use:   "rm <assistant>",
	Short: "Remove an assistant from the dot-agent source repo",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		m, sourceDir, err := loadManifest()
		if err != nil {
			fatal(err)
		}

		name := args[0]
		entry, err := m.ResolveAssistant(name)
		if err != nil {
			fatal(err)
		}

		assistantSourceDir := filepath.Join(sourceDir, filepath.FromSlash(entry.Source))

		if rmDryRun {
			fmt.Printf("Would remove %s files from %s\n", name, assistantSourceDir)
			fmt.Printf("Would remove %s from dot-agent.yaml\n", name)
			return
		}

		fmt.Printf("Removing %s...\n", name)

		// 1. Remove files from source repo
		if _, err := os.Stat(assistantSourceDir); err == nil {
			if err := os.RemoveAll(assistantSourceDir); err != nil {
				fatal(err)
			}
		}

		// 2. Remove from manifest
		delete(m.Assistants, name)
		if err := config.WriteManifest(paths.ManifestPath(sourceDir), m); err != nil {
			fatal(err)
		}

		// 3. Stage changes
		if err := git.AddAll(sourceDir); err != nil {
			fatal(err)
		}
		fmt.Printf("Removed %s and staged changes in git.\n", name)
	},
}

func init() {
	rmCmd.Flags().BoolVar(&rmDryRun, "dry-run", false, "show what would be removed without actually removing")
	rootCmd.AddCommand(rmCmd)
}
