package cmd

import (
	"fmt"

	"github.com/cthulhu/dot-agent/internal/sync"
	"github.com/spf13/cobra"
)

var (
	applyDryRun bool
	applyBackup bool
	applyForce  bool
)

var applyCmd = &cobra.Command{
	Use:   "apply [claude|cursor|hermes|codex|gemini|copilot]",
	Short: "Apply source repo config to local assistant directories",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		m, sourceDir, err := loadManifest()
		if err != nil {
			fatal(err)
		}

		names, err := resolveAssistantArgs(m, args)
		if err != nil {
			fatal(err)
		}

		opts := sync.Options{DryRun: applyDryRun, Backup: applyBackup, Force: applyForce}
		for _, name := range names {
			entry, err := m.ResolveAssistant(name)
			if err != nil {
				fatal(err)
			}

			if !applyForce && !applyDryRun {
				drift, err := sync.Compare(sourceDir, entry)
				if err != nil {
					fatal(err)
				}
				hasModified := false
				for _, e := range drift.Entries {
					if e.Status == "modified" {
						hasModified = true
						fmt.Printf("Warning: local %s differs from source (%s). Use --force to apply.\n", name, e.RelPath)
					}
				}
				if hasModified {
					continue
				}
			}

			fmt.Printf("Applying %s...\n", name)
			result, err := sync.Apply(sourceDir, entry, opts)
			if err != nil {
				fatal(err)
			}
			sync.PrintResult(result)
		}
	},
}

func init() {
	applyCmd.Flags().BoolVar(&applyDryRun, "dry-run", false, "show what would change without writing")
	applyCmd.Flags().BoolVar(&applyBackup, "backup", false, "backup overwritten local files")
	applyCmd.Flags().BoolVar(&applyForce, "force", false, "apply even when local config differs")
	rootCmd.AddCommand(applyCmd)
}
