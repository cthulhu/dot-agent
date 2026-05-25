package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cthulhu/dot-agent/internal/assistant"
	"github.com/cthulhu/dot-agent/internal/git"
	"github.com/cthulhu/dot-agent/internal/paths"
	"github.com/spf13/cobra"
)

var (
	initRepoURL  string
	initPath     string
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize or clone the dot-agent source git repo",
	Run: func(cmd *cobra.Command, args []string) {
		if err := git.RequireGit(); err != nil {
			fatal(err)
		}

		sourceDir, err := resolveInitSourceDir()
		if err != nil {
			fatal(err)
		}

		if initRepoURL != "" {
			if _, err := os.Stat(filepath.Join(sourceDir, ".git")); err == nil {
				fatal(fmt.Errorf("source directory %s already exists; remove it or use --path", sourceDir))
			}
			if err := git.Clone(initRepoURL, sourceDir); err != nil {
				fatal(err)
			}
			fmt.Printf("Cloned %s into %s\n", initRepoURL, sourceDir)
		} else {
			if err := git.Init(sourceDir); err != nil {
				fatal(err)
			}
			fmt.Printf("Initialized git repo at %s\n", sourceDir)
		}

		for _, name := range assistant.KnownNames() {
			dir := filepath.Join(sourceDir, "assistants", name)
			if err := os.MkdirAll(dir, 0o755); err != nil {
				fatal(err)
			}
		}

		if err := assistant.WriteDefaultManifest(sourceDir); err != nil {
			fatal(err)
		}

		gitignore := filepath.Join(sourceDir, ".gitignore")
		if _, err := os.Stat(gitignore); os.IsNotExist(err) {
			if err := os.WriteFile(gitignore, []byte("# OS files\n.DS_Store\nThumbs.db\n"), 0o644); err != nil {
				fatal(err)
			}
		}

		userCfg := &paths.UserConfig{SourcePath: sourceDir}
		if initRepoURL != "" {
			userCfg.RemoteURL = initRepoURL
			if err := git.SetRemote(sourceDir, initRepoURL); err != nil {
				fatal(err)
			}
		}
		if err := paths.SaveUserConfig(userCfg); err != nil {
			fatal(err)
		}

		cfgPath, _ := paths.UserConfigPath()
		fmt.Printf("Wrote manifest and saved config to %s\n", cfgPath)
		if initRepoURL == "" {
			fmt.Println("Note: configure a remote before push/pull:")
			fmt.Println("  dot-agent init --repo git@github.com:you/dot-agent.git")
			fmt.Println("  or: git -C <source> remote add origin <url>")
		}

		var nextCmds []string
		for _, name := range assistant.KnownNames() {
			nextCmds = append(nextCmds, "dot-agent add "+name)
		}
		fmt.Printf("Next: %s\n", strings.Join(nextCmds, " && "))
		fmt.Println("Tip: dot-agent cd opens a shell in your source repo")
	},
}

func resolveInitSourceDir() (string, error) {
	if initPath != "" {
		return paths.ExpandPath(initPath)
	}
	if sourceFlag != "" {
		return paths.ExpandPath(sourceFlag)
	}
	return paths.DefaultSourceDir()
}

func init() {
	initCmd.Flags().StringVar(&initRepoURL, "repo", "", "git remote URL to clone")
	initCmd.Flags().StringVar(&initPath, "path", "", "local path for source repo (default: OS-specific dot-agent source dir)")
	rootCmd.AddCommand(initCmd)
}
