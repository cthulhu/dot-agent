package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var sourcePathCmd = &cobra.Command{
	Use:   "source-path",
	Short: "Print the resolved dot-agent source directory",
	Run: func(cmd *cobra.Command, args []string) {
		sourceDir, err := resolveSourceDirAbs()
		if err != nil {
			fatal(err)
		}
		fmt.Println(sourceDir)
	},
}

func init() {
	rootCmd.AddCommand(sourcePathCmd)
}
