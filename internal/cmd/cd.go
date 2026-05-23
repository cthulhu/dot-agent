package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"
)

const (
	envSubshell   = "DOT_AGENT_SUBSHELL"
	envSourceDir  = "DOT_AGENT_SOURCE_DIR"
)

var cdCmd = &cobra.Command{
	Use:   "cd",
	Short: "Launch a shell in the source directory",
	Long: `Launch an interactive shell in the dot-agent source repo.

This does not change the current directory of your existing shell. When you
exit the subshell, you return to where you were.

To change directory in your current shell instead, run:

  cd "$(dot-agent source-path)"`,
	Run: func(cmd *cobra.Command, args []string) {
		if os.Getenv(envSubshell) != "" {
			fatal(fmt.Errorf("already in a dot-agent subshell; exit before running dot-agent cd again"))
		}

		sourceDir, err := resolveSourceDirAbs()
		if err != nil {
			fatal(err)
		}

		shell, err := userShell()
		if err != nil {
			fatal(err)
		}

		shellCmd := exec.Command(shell)
		shellCmd.Dir = sourceDir
		shellCmd.Stdin = os.Stdin
		shellCmd.Stdout = os.Stdout
		shellCmd.Stderr = os.Stderr
		shellCmd.Env = append(os.Environ(),
			envSubshell+"=1",
			envSourceDir+"="+sourceDir,
		)

		if err := shellCmd.Run(); err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				os.Exit(exitErr.ExitCode())
			}
			fatal(err)
		}
	},
}

func userShell() (string, error) {
	if shell := os.Getenv("SHELL"); shell != "" {
		return shell, nil
	}
	if runtime.GOOS == "windows" {
		if shell := os.Getenv("ComSpec"); shell != "" {
			return shell, nil
		}
		return "cmd.exe", nil
	}
	return "/bin/sh", nil
}

func init() {
	rootCmd.AddCommand(cdCmd)
}
