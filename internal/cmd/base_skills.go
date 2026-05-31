package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cthulhu/dot-agent/internal/assistant"
	"github.com/cthulhu/dot-agent/internal/paths"
	"github.com/cthulhu/dot-agent/internal/skills"
	"github.com/spf13/cobra"
)

var (
	skillsDryRun   bool
	skillsForce    bool
	skillsFromRepo bool
)

var baseSkillsCmd = &cobra.Command{
	Use:   "skills",
	Short: "Manage skills in the base library",
}

var baseSkillsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List skills in the base library",
	Run: func(cmd *cobra.Command, args []string) {
		skillsDir, err := paths.BaseSkillsDir()
		if err != nil {
			fatal(err)
		}
		entries, err := skills.List(skillsDir)
		if err != nil {
			fatal(err)
		}
		if len(entries) == 0 {
			fmt.Printf("No skills in base library (%s)\n", skillsDir)
			fmt.Println("Pull a skill from an assistant: dot-agent base skills pull <assistant> <skill>")
			return
		}
		fmt.Printf("Base skills (%s):\n", skillsDir)
		for _, e := range entries {
			fmt.Printf("  %s\n", e.Name)
		}
	},
}

var baseSkillsShowCmd = &cobra.Command{
	Use:   "show <skill>",
	Short: "Show a skill from the base library",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		skillsDir, err := paths.BaseSkillsDir()
		if err != nil {
			fatal(err)
		}
		skillPath, err := skills.ResolveSkillPath(skillsDir, args[0])
		if err != nil {
			fatal(err)
		}
		data, err := os.ReadFile(filepath.Join(skillPath, skills.SkillFileName))
		if err != nil {
			fatal(fmt.Errorf("read skill %q: %w", args[0], err))
		}
		fmt.Printf("%s\n", skillPath)
		fmt.Println(strings.Repeat("-", 40))
		fmt.Print(string(data))
	},
}

var baseSkillsPullCmd = &cobra.Command{
	Use:   "pull <assistant> [skill|--all]",
	Short: "Copy skill(s) from an assistant into the base library",
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		m, sourceDir, err := loadManifest()
		if err != nil {
			fatal(err)
		}
		if err := validateSkillsAssistant(args[0]); err != nil {
			fatal(err)
		}
		skillName := ""
		if len(args) > 1 {
			skillName = args[1]
		}
		opts := skills.Options{DryRun: skillsDryRun, Force: skillsForce, FromRepo: skillsFromRepo}
		if err := skills.Pull(m, sourceDir, args[0], skillName, opts); err != nil {
			fatal(err)
		}
		if !skillsDryRun {
			fmt.Printf("Pulled skill(s) from %s into base library\n", args[0])
		}
	},
}

var baseSkillsPushCmd = &cobra.Command{
	Use:   "push <assistant> [skill|--all]",
	Short: "Copy skill(s) from the base library into an assistant",
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		m, sourceDir, err := loadManifest()
		if err != nil {
			fatal(err)
		}
		if err := validateSkillsAssistant(args[0]); err != nil {
			fatal(err)
		}
		skillName := ""
		if len(args) > 1 {
			skillName = args[1]
		}
		opts := skills.Options{DryRun: skillsDryRun, Force: skillsForce, FromRepo: skillsFromRepo}
		if err := skills.Push(m, sourceDir, args[0], skillName, opts); err != nil {
			fatal(err)
		}
		if !skillsDryRun {
			fmt.Printf("Pushed skill(s) to %s from base library\n", args[0])
		}
	},
}

var baseSkillsRemoveCmd = &cobra.Command{
	Use:   "remove <skill>",
	Short: "Remove a skill from the base library",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		opts := skills.Options{DryRun: skillsDryRun}
		if err := skills.Remove(args[0], opts); err != nil {
			fatal(err)
		}
		if !skillsDryRun {
			fmt.Printf("Removed %s from base library\n", args[0])
		}
	},
}

var baseSkillsCopyCmd = &cobra.Command{
	Use:   "copy",
	Short: "Copy skills between assistants or via the base library",
}

var baseSkillsCopyFromCmd = &cobra.Command{
	Use:   "from <assistant> <skill|--all>",
	Short: "Copy skill(s) from an assistant into the base library",
	Args:  cobra.RangeArgs(2, 2),
	Run: func(cmd *cobra.Command, args []string) {
		baseSkillsPullCmd.Run(cmd, []string{args[0], args[1]})
	},
}

var baseSkillsCopyToCmd = &cobra.Command{
	Use:   "to <assistant> <skill|--all>",
	Short: "Copy skill(s) from the base library into an assistant",
	Args:  cobra.RangeArgs(2, 2),
	Run: func(cmd *cobra.Command, args []string) {
		baseSkillsPushCmd.Run(cmd, []string{args[0], args[1]})
	},
}

var baseSkillsCopyBetweenCmd = &cobra.Command{
	Use:   "between <from-assistant> <to-assistant> <skill|--all>",
	Short: "Copy skill(s) directly from one assistant to another",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		m, sourceDir, err := loadManifest()
		if err != nil {
			fatal(err)
		}
		if err := validateSkillsAssistant(args[0]); err != nil {
			fatal(err)
		}
		if err := validateSkillsAssistant(args[1]); err != nil {
			fatal(err)
		}
		opts := skills.Options{DryRun: skillsDryRun, Force: skillsForce, FromRepo: skillsFromRepo}
		if err := skills.CopyBetween(m, sourceDir, args[0], args[1], args[2], opts); err != nil {
			fatal(err)
		}
		if !skillsDryRun {
			fmt.Printf("Copied skill(s) from %s to %s\n", args[0], args[1])
		}
	},
}

func validateSkillsAssistant(name string) error {
	if !assistant.IsKnown(name) {
		return fmt.Errorf("unknown assistant %q (use %s)", name, assistant.KnownNamesString())
	}
	if !assistant.SupportsSkills(name) {
		return fmt.Errorf("assistant %q does not support skills", name)
	}
	return nil
}

func init() {
	skillsAssistants := skillsAssistantsPipeList()
	baseSkillsPullCmd.Use = fmt.Sprintf("pull <%s> [skill|--all]", skillsAssistants)
	baseSkillsPushCmd.Use = fmt.Sprintf("push <%s> [skill|--all]", skillsAssistants)
	baseSkillsCopyFromCmd.Use = fmt.Sprintf("from <%s> <skill|--all>", skillsAssistants)
	baseSkillsCopyToCmd.Use = fmt.Sprintf("to <%s> <skill|--all>", skillsAssistants)
	baseSkillsCopyBetweenCmd.Use = fmt.Sprintf("between <%s> <%s> <skill|--all>", skillsAssistants, skillsAssistants)

	for _, c := range []*cobra.Command{
		baseSkillsPullCmd,
		baseSkillsPushCmd,
		baseSkillsCopyFromCmd,
		baseSkillsCopyToCmd,
		baseSkillsCopyBetweenCmd,
	} {
		c.Flags().BoolVar(&skillsDryRun, "dry-run", false, "show what would change without writing")
		c.Flags().BoolVar(&skillsForce, "force", false, "overwrite existing skills")
		c.Flags().BoolVar(&skillsFromRepo, "repo", false, "use source repo instead of local assistant directory")
	}

	baseSkillsRemoveCmd.Flags().BoolVar(&skillsDryRun, "dry-run", false, "show what would change without writing")

	baseSkillsCopyCmd.AddCommand(baseSkillsCopyFromCmd)
	baseSkillsCopyCmd.AddCommand(baseSkillsCopyToCmd)
	baseSkillsCopyCmd.AddCommand(baseSkillsCopyBetweenCmd)

	baseSkillsCmd.AddCommand(baseSkillsListCmd)
	baseSkillsCmd.AddCommand(baseSkillsShowCmd)
	baseSkillsCmd.AddCommand(baseSkillsPullCmd)
	baseSkillsCmd.AddCommand(baseSkillsPushCmd)
	baseSkillsCmd.AddCommand(baseSkillsRemoveCmd)
	baseSkillsCmd.AddCommand(baseSkillsCopyCmd)
	baseCmd.AddCommand(baseSkillsCmd)
}

func skillsAssistantsPipeList() string {
	var names []string
	for _, name := range assistant.KnownNames() {
		if assistant.SupportsSkills(name) {
			names = append(names, name)
		}
	}
	return strings.Join(names, "|")
}
