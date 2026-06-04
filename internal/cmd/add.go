package cmd

import (
	"fmt"
	"strings"

	"github.com/cthulhu/dot-agent/internal/assistant"
	"github.com/cthulhu/dot-agent/internal/config"
	"github.com/cthulhu/dot-agent/internal/git"
	"github.com/cthulhu/dot-agent/internal/paths"
	"github.com/cthulhu/dot-agent/internal/sync"
	"github.com/spf13/cobra"
)

var addDryRun bool

var addCmd = &cobra.Command{
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
		manifestUpdated := false
		for _, name := range names {
			entry, err := m.ResolveAssistant(name)
			if err != nil {
				// If missing but known, add it now
				if defaultEntry, ok := assistant.DefaultEntry(name); ok {
					fmt.Printf("Adding %s to manifest...\n", name)
					if m.Assistants == nil {
						m.Assistants = make(map[string]config.AssistantEntry)
					}
					m.Assistants[name] = defaultEntry
					entry = defaultEntry
					manifestUpdated = true
				} else {
					fatal(err)
				}
			}

			fmt.Printf("Adding %s...\n", name)
			result, err := sync.Add(sourceDir, entry, opts)
			if err != nil {
				fatal(err)
			}
			sync.PrintResult(result)
		}

		if manifestUpdated && !addDryRun {
			if err := config.WriteManifest(paths.ManifestPath(sourceDir), m); err != nil {
				fatal(err)
			}
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
		// If not known, it must be in manifest
		if _, err := m.ResolveAssistant(name); err != nil {
			return nil, fmt.Errorf("unknown assistant %q (use %s)", name, assistant.KnownNamesString())
		}
	}
	return []string{name}, nil
}

func init() {
	addCmd.Use = fmt.Sprintf("add [%s]", strings.ReplaceAll(assistant.KnownNamesString(), ", ", "|"))
	addCmd.Flags().BoolVar(&addDryRun, "dry-run", false, "show what would be captured without writing")
	rootCmd.AddCommand(addCmd)
}
