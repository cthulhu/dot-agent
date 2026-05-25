package config_test

import (
	"os"
	"path/filepath"

	"github.com/cthulhu/dot-agent/internal/config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {
	var tmpDir string

	BeforeEach(func() {
		var err error
		tmpDir, err = os.MkdirTemp("", "dot-agent-config-test-*")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		Expect(os.RemoveAll(tmpDir)).To(Succeed())
	})

	Describe("Manifest", func() {
		It("should successfully load and validate a valid manifest", func() {
			path := filepath.Join(tmpDir, "dot-agent.yaml")
			content := `
version: 1
assistants:
  claude:
    source: assistants/claude
    target: ~/.claude
    ignore:
      - "*.log"
`
			Expect(os.WriteFile(path, []byte(content), 0o644)).To(Succeed())

			m, err := config.LoadManifest(path)
			Expect(err).NotTo(HaveOccurred())
			Expect(m.Version).To(Equal(1))
			Expect(m.Assistants).To(HaveKey("claude"))
			
			entry := m.Assistants["claude"]
			Expect(entry.Source).To(Equal("assistants/claude"))
			Expect(entry.Target).To(Equal("~/.claude"))
			Expect(entry.Ignore).To(ContainElement("*.log"))
		})

		It("should fail to load a non-existent manifest", func() {
			path := filepath.Join(tmpDir, "non-existent.yaml")
			m, err := config.LoadManifest(path)
			Expect(err).To(HaveOccurred())
			Expect(m).To(BeNil())
		})

		It("should fail on unsupported manifest version", func() {
			path := filepath.Join(tmpDir, "invalid-version.yaml")
			content := "version: 2\nassistants: {}"
			Expect(os.WriteFile(path, []byte(content), 0o644)).To(Succeed())

			_, err := config.LoadManifest(path)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("unsupported manifest version 2"))
		})

		It("should fail if no assistants are defined", func() {
			path := filepath.Join(tmpDir, "no-assistants.yaml")
			content := "version: 1\nassistants: {}"
			Expect(os.WriteFile(path, []byte(content), 0o644)).To(Succeed())

			_, err := config.LoadManifest(path)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("manifest has no assistants defined"))
		})

		DescribeTable("Validation errors for assistant entries",
			func(content string, expectedErr string) {
				path := filepath.Join(tmpDir, "invalid-assistant.yaml")
				Expect(os.WriteFile(path, []byte(content), 0o644)).To(Succeed())

				_, err := config.LoadManifest(path)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(expectedErr))
			},
			Entry("missing source", "version: 1\nassistants:\n  test:\n    target: ~/test", "assistant \"test\": source path is required"),
			Entry("missing target", "version: 1\nassistants:\n  test:\n    source: path/test", "assistant \"test\": target path is required"),
		)
	})

	Describe("Resolution", func() {
		var m *config.Manifest

		BeforeEach(func() {
			m = &config.Manifest{
				Version: 1,
				Assistants: map[string]config.AssistantEntry{
					"claude": {Source: "s1", Target: "t1"},
					"cursor": {Source: "s2", Target: "t2"},
				},
			}
		})

		It("should resolve an existing assistant", func() {
			entry, err := m.ResolveAssistant("claude")
			Expect(err).NotTo(HaveOccurred())
			Expect(entry.Source).To(Equal("s1"))
		})

		It("should return an error for an unknown assistant with available names", func() {
			_, err := m.ResolveAssistant("unknown")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("assistant \"unknown\" not found in manifest"))
			Expect(err.Error()).To(ContainSubstring("available: claude, cursor"))
		})

		It("should return assistant names filtered or all", func() {
			names, err := m.AssistantNames(nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(names).To(ConsistOf("claude", "cursor"))

			filtered, err := m.AssistantNames([]string{"claude"})
			Expect(err).NotTo(HaveOccurred())
			Expect(filtered).To(Equal([]string{"claude"}))
		})

		It("should fail to return names if filter contains unknown assistant", func() {
			_, err := m.AssistantNames([]string{"unknown"})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("assistant \"unknown\" not in manifest"))
		})
	})

	It("should write manifest to file", func() {
		path := filepath.Join(tmpDir, "output.yaml")
		m := &config.Manifest{
			Version: 1,
			Assistants: map[string]config.AssistantEntry{
				"test": {Source: "src", Target: "tgt"},
			},
		}

		Expect(config.WriteManifest(path, m)).To(Succeed())
		Expect(path).To(BeAnExistingFile())

		data, err := os.ReadFile(path)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(data)).To(ContainSubstring("source: src"))
	})
})
