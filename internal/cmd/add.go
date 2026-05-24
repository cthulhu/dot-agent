package cmd

import (
	"fmt"

	"github.com/cthulhu/dot-agent/internal/assistant"
	"github.com/cthulhu/dot-agent/internal/config"
	"github.com/cthulhu/dot-agent/internal/git"
	"github.com/cthulhu/dot-agent/internal/sync"
	"github.com/spf13/cobra"
)

var addDryRun bool

var addCmd = &cobra.Command{
	Use:   "add [claude|cursor|hermes|codex|gemini]",
	Short: "Capture local assistant config into the source repo",
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

		opts := sync.Options{DryRun: addDryRun}
		for _, name := range names {
			entry, err := m.ResolveAssistant(name)
			if err != nil {
				fatal(err)
			}
			fmt.Printf("Adding %s...\n", name)
			result, err := sync.Add(sourceDir, entry, opts)
			if err != nil {
				fatal(err)
			}
			sync.PrintResult(result)
		}

		if !addDryRun {
			if err := git.AddAll(sourceDir); err != nil {
				fatal(err)
			}
			fmt.Println("Staged changes in git.")
		}
	},
}

func resolveAssistantArgs(m *config.Manifest, args []string) ([]string, error) {
	if len(args) == 0 {
		return m.AssistantNames(nil)
	}
	name := args[0]
	if !assistant.IsKnown(name) {
		return nil, fmt.Errorf("unknown assistant %q (use %s)", name, assistant.KnownNamesString())
	}
	return m.AssistantNames([]string{name})
}

func init() {
	addCmd.Flags().BoolVar(&addDryRun, "dry-run", false, "show what would be captured without writing")
	rootCmd.AddCommand(addCmd)
}
