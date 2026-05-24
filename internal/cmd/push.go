package cmd

import (
	"fmt"

	"github.com/cthulhu/dot-agent/internal/git"
	"github.com/spf13/cobra"
)

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Commit and push source repo to remote",
	Run: func(cmd *cobra.Command, args []string) {
		if err := git.RequireGit(); err != nil {
			fatal(err)
		}

		_, sourceDir, err := loadManifest()
		if err != nil {
			fatal(err)
		}

		if err := ensureGitRemote(sourceDir); err != nil {
			fatal(err)
		}

		if err := git.Commit(sourceDir, "dot-agent: sync assistant config"); err != nil {
			fatal(err)
		}
		if err := git.Push(sourceDir); err != nil {
			fatal(err)
		}
		fmt.Println("Pushed to remote.")
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)
}
