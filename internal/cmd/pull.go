package cmd

import (
	"fmt"

	"github.com/cthulhu/dot-agent/internal/git"
	"github.com/cthulhu/dot-agent/internal/sync"
	"github.com/spf13/cobra"
)

var pullApply bool

var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull latest config from remote",
	Run: func(cmd *cobra.Command, args []string) {
		if err := git.RequireGit(); err != nil {
			fatal(err)
		}

		m, sourceDir, err := loadManifest()
		if err != nil {
			fatal(err)
		}

		if err := git.Pull(sourceDir); err != nil {
			fatal(err)
		}
		fmt.Println("Pulled latest changes.")

		if !pullApply {
			return
		}

		names, err := m.AssistantNames(nil)
		if err != nil {
			fatal(err)
		}
		opts := sync.Options{Force: true}
		for _, name := range names {
			entry, err := m.ResolveAssistant(name)
			if err != nil {
				fatal(err)
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
	pullCmd.Flags().BoolVar(&pullApply, "apply", false, "apply config to local after pull")
	rootCmd.AddCommand(pullCmd)
}
