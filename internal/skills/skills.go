package skills

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/cthulhu/dot-agent/internal/assistant"
	"github.com/cthulhu/dot-agent/internal/config"
	"github.com/cthulhu/dot-agent/internal/paths"
)

const SkillFileName = "SKILL.md"

type Entry struct {
	Name string
	Path string
}

type Options struct {
	DryRun   bool
	Force    bool
	FromRepo bool
}

func baseSkillsDir() (string, error) {
	return paths.BaseSkillsDir()
}

func Pull(m *config.Manifest, sourceRoot, assistantName, skillName string, opts Options) error {
	srcRoot, err := ResolveAssistantSkillsDir(m, sourceRoot, assistantName, opts.FromRepo)
	if err != nil {
		return err
	}
	baseDir, err := baseSkillsDir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(baseDir, 0o755); err != nil && !opts.DryRun {
		return err
	}
	return copyToBase(srcRoot, baseDir, skillName, opts)
}

func Push(m *config.Manifest, sourceRoot, assistantName, skillName string, opts Options) error {
	dstRoot, err := ResolveAssistantSkillsDir(m, sourceRoot, assistantName, opts.FromRepo)
	if err != nil {
		return err
	}
	baseDir, err := baseSkillsDir()
	if err != nil {
		return err
	}
	return copyFromBase(baseDir, dstRoot, skillName, opts)
}

func List(dir string) ([]Entry, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var skills []Entry
	for _, e := range entries {
		if !e.IsDir() || strings.HasPrefix(e.Name(), ".") {
			continue
		}
		skillPath := filepath.Join(dir, e.Name())
		if !isSkillDir(skillPath) {
			continue
		}
		skills = append(skills, Entry{Name: e.Name(), Path: skillPath})
	}
	sort.Slice(skills, func(i, j int) bool {
		return skills[i].Name < skills[j].Name
	})
	return skills, nil
}

func isSkillDir(dir string) bool {
	info, err := os.Stat(filepath.Join(dir, SkillFileName))
	return err == nil && !info.IsDir()
}

func ResolveAssistantSkillsDir(m *config.Manifest, sourceRoot, name string, fromRepo bool) (string, error) {
	if !assistant.SupportsSkills(name) {
		return "", fmt.Errorf("assistant %q does not support skills (use %s)", name, skillsAssistantsString())
	}
	entry, err := m.ResolveAssistant(name)
	if err != nil {
		return "", err
	}
	rel, _ := assistant.SkillsRelPath(name)
	if fromRepo {
		return filepath.Join(sourceRoot, filepath.FromSlash(entry.Source), rel), nil
	}
	target, err := paths.ExpandPath(entry.Target)
	if err != nil {
		return "", err
	}
	return filepath.Join(target, rel), nil
}

func ResolveSkillPath(skillsRoot, name string) (string, error) {
	if name == "" {
		return "", fmt.Errorf("skill name is required")
	}
	if strings.Contains(name, string(os.PathSeparator)) || strings.Contains(name, "/") {
		return "", fmt.Errorf("invalid skill name %q", name)
	}
	return filepath.Join(skillsRoot, name), nil
}

func InitBase() (string, error) {
	skillsDir, err := baseSkillsDir()
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(skillsDir, 0o755); err != nil {
		return "", err
	}
	return skillsDir, nil
}

// CopyToBase copies one or all skills from an assistant skills directory into the base library.
func CopyToBase(srcRoot, baseDir, skillName string, opts Options) error {
	return copyToBase(srcRoot, baseDir, skillName, opts)
}

// CopyFromBase copies one or all skills from the base library into an assistant skills directory.
func CopyFromBase(baseDir, dstRoot, skillName string, opts Options) error {
	return copyFromBase(baseDir, dstRoot, skillName, opts)
}

func copyToBase(srcRoot, baseDir, skillName string, opts Options) error {
	if skillName == "" || skillName == "--all" {
		return pullAll(srcRoot, baseDir, opts)
	}
	src, err := ResolveSkillPath(srcRoot, skillName)
	if err != nil {
		return err
	}
	dst, err := ResolveSkillPath(baseDir, skillName)
	if err != nil {
		return err
	}
	return copySkill(src, dst, opts)
}

func copyFromBase(baseDir, dstRoot, skillName string, opts Options) error {
	if skillName == "" || skillName == "--all" {
		return pushAll(baseDir, dstRoot, opts)
	}
	src, err := ResolveSkillPath(baseDir, skillName)
	if err != nil {
		return err
	}
	dst, err := ResolveSkillPath(dstRoot, skillName)
	if err != nil {
		return err
	}
	return copySkill(src, dst, opts)
}

func CopyBetween(m *config.Manifest, sourceRoot, fromAssistant, toAssistant, skillName string, opts Options) error {
	srcRoot, err := ResolveAssistantSkillsDir(m, sourceRoot, fromAssistant, opts.FromRepo)
	if err != nil {
		return err
	}
	dstRoot, err := ResolveAssistantSkillsDir(m, sourceRoot, toAssistant, false)
	if err != nil {
		return err
	}

	if skillName == "" || skillName == "--all" {
		return copyAllSkills(srcRoot, dstRoot, opts)
	}

	src, err := ResolveSkillPath(srcRoot, skillName)
	if err != nil {
		return err
	}
	dst, err := ResolveSkillPath(dstRoot, skillName)
	if err != nil {
		return err
	}
	return copySkill(src, dst, opts)
}

func Remove(skillName string, opts Options) error {
	baseDir, err := baseSkillsDir()
	if err != nil {
		return err
	}
	skillPath, err := ResolveSkillPath(baseDir, skillName)
	if err != nil {
		return err
	}
	if _, err := os.Stat(skillPath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("skill %q not found in base library", skillName)
		}
		return err
	}
	if !isSkillDir(skillPath) {
		return fmt.Errorf("%q is not a valid skill directory (missing %s)", skillName, SkillFileName)
	}
	if opts.DryRun {
		fmt.Printf("would remove: %s\n", skillPath)
		return nil
	}
	return os.RemoveAll(skillPath)
}

func pullAll(srcRoot, baseDir string, opts Options) error {
	skills, err := List(srcRoot)
	if err != nil {
		return err
	}
	if len(skills) == 0 {
		return fmt.Errorf("no skills found in %s", srcRoot)
	}
	for _, skill := range skills {
		dst := filepath.Join(baseDir, skill.Name)
		if err := copySkill(skill.Path, dst, opts); err != nil {
			return err
		}
	}
	return nil
}

func pushAll(baseDir, dstRoot string, opts Options) error {
	skills, err := List(baseDir)
	if err != nil {
		return err
	}
	if len(skills) == 0 {
		return fmt.Errorf("no skills found in base library")
	}
	for _, skill := range skills {
		dst := filepath.Join(dstRoot, skill.Name)
		if err := copySkill(skill.Path, dst, opts); err != nil {
			return err
		}
	}
	return nil
}

func copyAllSkills(srcRoot, dstRoot string, opts Options) error {
	skills, err := List(srcRoot)
	if err != nil {
		return err
	}
	if len(skills) == 0 {
		return fmt.Errorf("no skills found in %s", srcRoot)
	}
	for _, skill := range skills {
		dst := filepath.Join(dstRoot, skill.Name)
		if err := copySkill(skill.Path, dst, opts); err != nil {
			return err
		}
	}
	return nil
}

func copySkill(src, dst string, opts Options) error {
	if _, err := os.Stat(src); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("skill not found: %s", src)
		}
		return err
	}
	if !isSkillDir(src) {
		return fmt.Errorf("not a valid skill directory (missing %s): %s", SkillFileName, src)
	}

	if _, err := os.Stat(dst); err == nil {
		if !opts.Force {
			return fmt.Errorf("destination already exists: %s (use --force to overwrite)", dst)
		}
		if opts.DryRun {
			fmt.Printf("would overwrite: %s -> %s\n", src, dst)
			return nil
		}
		if err := os.RemoveAll(dst); err != nil {
			return err
		}
	} else if !os.IsNotExist(err) {
		return err
	}

	if opts.DryRun {
		fmt.Printf("would copy: %s -> %s\n", src, dst)
		return nil
	}

	return copyDir(src, dst)
}

func copyDir(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		return copyFile(path, target)
	})
}

func copyFile(src, dst string) error {
	info, err := os.Stat(src)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, info.Mode().Perm())
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

func skillsAssistantsString() string {
	var names []string
	for _, name := range assistant.KnownNames() {
		if assistant.SupportsSkills(name) {
			names = append(names, name)
		}
	}
	sort.Strings(names)
	return strings.Join(names, ", ")
}
