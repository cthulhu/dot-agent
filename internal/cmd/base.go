package cmd

import (
	"fmt"

	"github.com/cthulhu/dot-agent/internal/paths"
	"github.com/cthulhu/dot-agent/internal/skills"
	"github.com/spf13/cobra"
)

var baseCmd = &cobra.Command{
	Use:   "base",
	Short: "Manage the local base skills library (~/.dot-agent/base)",
}

var baseInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize ~/.dot-agent/base/skills",
	Run: func(cmd *cobra.Command, args []string) {
		skillsDir, err := skills.InitBase()
		if err != nil {
			fatal(err)
		}
		baseDir, err := paths.BaseDir()
		if err != nil {
			fatal(err)
		}
		fmt.Printf("Initialized base library at %s\n", baseDir)
		fmt.Printf("Skills directory: %s\n", skillsDir)
	},
}

func init() {
	baseCmd.AddCommand(baseInitCmd)
	rootCmd.AddCommand(baseCmd)
}
