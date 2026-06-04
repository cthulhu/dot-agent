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

		if err := writeDefaultREADME(sourceDir); err != nil {
			fatal(err)
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
	var dir string
	var err error
	if initPath != "" {
		dir, err = paths.ExpandPath(initPath)
	} else if sourceFlag != "" {
		dir, err = paths.ExpandPath(sourceFlag)
	} else {
		dir, err = paths.DefaultSourceDir()
	}
	if err != nil {
		return "", err
	}
	return filepath.Abs(dir)
}

func writeDefaultREADME(sourceDir string) error {
	readme := filepath.Join(sourceDir, "README.md")
	if _, err := os.Stat(readme); err == nil {
		return nil
	}
	// Also check for lowercase readme.md
	if _, err := os.Stat(filepath.Join(sourceDir, "readme.md")); err == nil {
		return nil
	}

	content := `# dot-agent configuration

This repository contains configuration and skills for your AI assistants, managed by [dot-agent](https://github.com/cthulhu/dot-agent).

## Prerequisites

Install **dot-agent**:

### Homebrew (macOS / Linux)

` + "```bash" + `
brew tap cthulhu/dot-agent https://github.com/cthulhu/dot-agent
brew install dot-agent
` + "```" + `

### Chocolatey (Windows)

` + "```powershell" + `
# Register the GitHub Packages source (run as Administrator)
choco source add -n="cthulhu" -s="https://nuget.pkg.github.com/cthulhu/index.json"

# Install dot-agent from the custom source
choco install dot-agent --source="cthulhu"
` + "```" + `

## Usage

### Sync to a new machine

` + "```bash" + `
dot-agent init --repo <this-repo-url>
dot-agent pull --apply
` + "```" + `

### Capture local changes

` + "```bash" + `
# Capture changes for a specific assistant
dot-agent add claude

# Push changes to the remote repository
dot-agent push
` + "```" + `

### Apply changes from repo to local

` + "```bash" + `
# Pull latest changes from the remote repository
dot-agent pull

# Apply changes to local assistant directories
dot-agent apply
` + "```" + `

## Repository Layout

- ` + "`dot-agent.yaml`" + `: Main configuration file
- ` + "`assistants/`" + `: Directory containing assistant-specific configurations and skills
`

	return os.WriteFile(readme, []byte(content), 0o644)
}

func init() {
	initCmd.Flags().StringVar(&initRepoURL, "repo", "", "git remote URL to clone")
	initCmd.Flags().StringVar(&initPath, "path", "", "local path for source repo (default: OS-specific dot-agent source dir)")
	rootCmd.AddCommand(initCmd)
}
