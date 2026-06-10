package cmd

import (
	"fmt"
	"strings"

	"github.com/cthulhu/dot-agent/internal/assistant"
	"github.com/cthulhu/dot-agent/internal/git"
	"github.com/cthulhu/dot-agent/internal/sync"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Short: "Show git status and config drift vs local",
	Args:  cobra.MaximumNArgs(1),
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

		names, err := resolveAssistantArgs(m, args)
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
	statusCmd.Use = fmt.Sprintf("status [%s]", strings.ReplaceAll(assistant.KnownNamesString(), ", ", "|"))
	rootCmd.AddCommand(statusCmd)
}
