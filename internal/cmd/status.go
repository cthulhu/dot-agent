package cmd

import (
	"fmt"

	"github.com/cthulhu/dot-agent/internal/git"
	"github.com/cthulhu/dot-agent/internal/sync"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show git status and config drift vs local",
	Run: func(cmd *cobra.Command, args []string) {
		m, sourceDir, err := loadManifest()
		if err != nil {
			fatal(err)
		}

		if err := git.RequireGit(); err != nil {
			fatal(err)
		}

		porcelain, err := git.StatusPorcelain(sourceDir)
		if err != nil {
			fatal(err)
		}
		fmt.Println("Git working tree:")
		if porcelain == "" {
			fmt.Println("  clean")
		} else {
			fmt.Println(porcelain)
		}

		names, err := m.AssistantNames(nil)
		if err != nil {
			fatal(err)
		}

		for _, name := range names {
			entry, err := m.ResolveAssistant(name)
			if err != nil {
				fatal(err)
			}
			fmt.Printf("\nDrift (%s):\n", name)
			report, err := sync.Compare(sourceDir, entry)
			if err != nil {
				fatal(err)
			}
			sync.PrintDrift(report)
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
