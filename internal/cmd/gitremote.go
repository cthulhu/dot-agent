package cmd

import (
	"fmt"

	"github.com/cthulhu/dot-agent/internal/git"
	"github.com/cthulhu/dot-agent/internal/paths"
)

func ensureGitRemote(sourceDir string) error {
	if _, err := git.RemoteURL(sourceDir); err == nil {
		return nil
	}

	userCfg, err := paths.LoadUserConfig()
	if err != nil {
		return err
	}
	if userCfg != nil && userCfg.RemoteURL != "" {
		if err := git.SetRemote(sourceDir, userCfg.RemoteURL); err != nil {
			return err
		}
		fmt.Printf("Configured git remote origin: %s\n", userCfg.RemoteURL)
		return nil
	}

	return fmt.Errorf(`no git remote configured for %s

Set a remote with either:
  dot-agent init --repo git@github.com:you/dot-agent.git
  git -C %q remote add origin <url>

Or save remote_url in your dot-agent config`, sourceDir, sourceDir)
}
