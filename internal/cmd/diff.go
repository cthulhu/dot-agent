package cmd

import (
	"fmt"

	"github.com/cthulhu/dot-agent/internal/sync"
	"github.com/spf13/cobra"
)

var diffCmd = &cobra.Command{
	Use:   "diff [claude|cursor|hermes|codex]",
	Short: "Show differences between source repo and local config",
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

		for _, name := range names {
			entry, err := m.ResolveAssistant(name)
			if err != nil {
				fatal(err)
			}
			fmt.Printf("=== %s ===\n", name)
			out, err := sync.Diff(sourceDir, entry)
			if err != nil {
				fatal(err)
			}
			if out == "" {
				fmt.Println("No differences.")
			} else {
				fmt.Print(out)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(diffCmd)
}
