package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cthulhu/dot-agent/internal/assistant"
	"github.com/cthulhu/dot-agent/internal/config"
	"github.com/cthulhu/dot-agent/internal/git"
	"github.com/cthulhu/dot-agent/internal/paths"
	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Validate git, paths, and manifest",
	Run: func(cmd *cobra.Command, args []string) {
		ok := true
		report := func(msg string) { fmt.Println("ok:", msg) }
		fail := func(msg string) {
			fmt.Println("fail:", msg)
			ok = false
		}

		if err := git.RequireGit(); err != nil {
			fail(err.Error())
		} else {
			report("git found in PATH")
		}

		home, err := paths.HomeDir()
		if err != nil {
			fail(fmt.Sprintf("home dir: %v", err))
		} else {
			report(fmt.Sprintf("home directory: %s", home))
		}

		cfgDir, err := paths.ConfigDir()
		if err != nil {
			fail(fmt.Sprintf("config dir: %v", err))
		} else {
			report(fmt.Sprintf("config directory: %s", cfgDir))
		}

		sourceDir, err := resolveSourceDir()
		if err != nil {
			fail(fmt.Sprintf("source dir: %v", err))
		} else {
			report(fmt.Sprintf("source directory: %s", sourceDir))
		}

		if _, err := os.Stat(filepath.Join(sourceDir, ".git")); err != nil {
			fail(fmt.Sprintf("source is not a git repo: %s", sourceDir))
		} else {
			report("source git repo exists")
		}

		if _, err := git.RemoteURL(sourceDir); err != nil {
			userCfg, cfgErr := paths.LoadUserConfig()
			if cfgErr != nil {
				fail(fmt.Sprintf("remote: %v", err))
			} else if userCfg != nil && userCfg.RemoteURL != "" {
				fmt.Printf("warn: git remote not set; push/pull will configure origin from config: %s\n", userCfg.RemoteURL)
			} else {
				fail("no git remote configured; run dot-agent init --repo <url> or git remote add origin <url>")
			}
		} else {
			report("git remote origin configured")
		}

		manifestPath := paths.ManifestPath(sourceDir)
		m, err := config.LoadManifest(manifestPath)
		if err != nil {
			fail(fmt.Sprintf("manifest: %v", err))
		} else {
			report(fmt.Sprintf("manifest valid (%d assistants)", len(m.Assistants)))
			for _, name := range assistant.KnownNames() {
				entry, err := m.ResolveAssistant(name)
				if err != nil {
					fail(fmt.Sprintf("assistant %s: %v", name, err))
					continue
				}
				target, err := paths.ExpandPath(entry.Target)
				if err != nil {
					fail(fmt.Sprintf("%s target path: %v", name, err))
					continue
				}
				if _, err := os.Stat(target); os.IsNotExist(err) {
					fmt.Printf("warn: %s target not found (will be created on apply): %s\n", name, target)
				} else if err != nil {
					fail(fmt.Sprintf("%s target: %v", name, err))
				} else {
					report(fmt.Sprintf("%s target exists: %s", name, target))
				}
			}
		}

		if !ok {
			os.Exit(1)
		}
		fmt.Println("\nAll checks passed.")
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
