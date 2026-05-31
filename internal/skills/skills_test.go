package skills_test

import (
	"os"
	"path/filepath"

	"github.com/cthulhu/dot-agent/internal/assistant"
	"github.com/cthulhu/dot-agent/internal/config"
	"github.com/cthulhu/dot-agent/internal/skills"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Skills", func() {
	var root string
	var baseSkills string
	var cursorSkills string

	BeforeEach(func() {
		var err error
		root, err = os.MkdirTemp("", "dot-agent-skills-*")
		Expect(err).NotTo(HaveOccurred())

		baseSkills = filepath.Join(root, "base", "skills")
		cursorSkills = filepath.Join(root, "home", ".cursor", "skills-cursor")
		Expect(os.MkdirAll(baseSkills, 0o755)).To(Succeed())
		Expect(os.MkdirAll(filepath.Join(cursorSkills, "demo"), 0o755)).To(Succeed())
		Expect(os.WriteFile(filepath.Join(cursorSkills, "demo", "SKILL.md"), []byte("# demo"), 0o644)).To(Succeed())
		Expect(os.MkdirAll(filepath.Join(cursorSkills, "demo", "scripts"), 0o755)).To(Succeed())
		Expect(os.WriteFile(filepath.Join(cursorSkills, "demo", "scripts", "run.sh"), []byte("#!/bin/sh"), 0o755)).To(Succeed())
	})

	AfterEach(func() {
		Expect(os.RemoveAll(root)).To(Succeed())
	})

	It("should list skills in a directory", func() {
		entries, err := skills.List(cursorSkills)
		Expect(err).NotTo(HaveOccurred())
		Expect(entries).To(HaveLen(1))
		Expect(entries[0].Name).To(Equal("demo"))
	})

	It("should copy a skill into base and back out", func() {
		err := skills.CopyToBase(cursorSkills, baseSkills, "demo", skills.Options{Force: true})
		Expect(err).NotTo(HaveOccurred())
		Expect(filepath.Join(baseSkills, "demo", "SKILL.md")).To(BeAnExistingFile())

		geminiSkills := filepath.Join(root, "home", ".gemini", "skills")
		err = skills.CopyFromBase(baseSkills, geminiSkills, "demo", skills.Options{Force: true})
		Expect(err).NotTo(HaveOccurred())
		Expect(filepath.Join(geminiSkills, "demo", "SKILL.md")).To(BeAnExistingFile())
	})

	It("should copy a skill between assistants", func() {
		claudeSkills := filepath.Join(root, "home", ".claude", "skills")
		Expect(os.MkdirAll(claudeSkills, 0o755)).To(Succeed())

		m := &config.Manifest{
			Version: 1,
			Assistants: map[string]config.AssistantEntry{
				assistant.Cursor: {
					Source: "assistants/cursor",
					Target: filepath.Join(root, "home", ".cursor"),
				},
				assistant.Claude: {
					Source: "assistants/claude",
					Target: filepath.Join(root, "home", ".claude"),
				},
			},
		}

		err := skills.CopyBetween(m, filepath.Join(root, "source"), assistant.Cursor, assistant.Claude, "demo", skills.Options{Force: true})
		Expect(err).NotTo(HaveOccurred())

		Expect(filepath.Join(claudeSkills, "demo", "SKILL.md")).To(BeAnExistingFile())
		Expect(filepath.Join(claudeSkills, "demo", "scripts", "run.sh")).To(BeAnExistingFile())
	})

	It("should reject invalid skill names", func() {
		_, err := skills.ResolveSkillPath(baseSkills, "../escape")
		Expect(err).To(HaveOccurred())
	})
})
