package sync_test

import (
	"os"
	"path/filepath"

	"github.com/cthulhu/dot-agent/internal/assistant"
	"github.com/cthulhu/dot-agent/internal/config"
	"github.com/cthulhu/dot-agent/internal/sync"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Sync", func() {
	var root string
	var sourceRoot string
	var localClaude string

	BeforeEach(func() {
		var err error
		root, err = os.MkdirTemp("", "dot-agent-test-*")
		Expect(err).NotTo(HaveOccurred())

		sourceRoot = filepath.Join(root, "source")
		localClaude = filepath.Join(root, "home", ".claude")

		Expect(os.MkdirAll(filepath.Join(localClaude, "rules"), 0o755)).To(Succeed())
		Expect(os.WriteFile(filepath.Join(localClaude, "settings.json"), []byte(`{"theme":"dark"}`), 0o644)).To(Succeed())
		Expect(os.WriteFile(filepath.Join(localClaude, "rules", "go.md"), []byte("# Go rules"), 0o644)).To(Succeed())
		Expect(os.WriteFile(filepath.Join(localClaude, ".env"), []byte("SECRET=1"), 0o644)).To(Succeed())
	})

	AfterEach(func() {
		Expect(os.RemoveAll(root)).To(Succeed())
	})

	It("should successfully Add and Apply configurations", func() {
		entry := assistant.DefaultClaude()
		entry.Target = localClaude
		Expect(os.MkdirAll(filepath.Join(sourceRoot, entry.Source), 0o755)).To(Succeed())
		Expect(config.WriteManifest(filepath.Join(sourceRoot, "dot-agent.yaml"), assistant.DefaultManifest())).To(Succeed())

		// Test Add
		result, err := sync.Add(sourceRoot, entry, sync.Options{})
		Expect(err).NotTo(HaveOccurred())

		foundBlocked := false
		for _, a := range result.Actions {
			if a.Action == "blocked" && a.RelPath == ".env" {
				foundBlocked = true
			}
		}
		Expect(foundBlocked).To(BeTrue(), "expected .env to be blocked during add")

		settingsPath := filepath.Join(sourceRoot, entry.Source, "settings.json")
		Expect(settingsPath).To(BeAnExistingFile())

		// Test Apply
		applyHome := filepath.Join(root, "apply-home", ".claude")
		entry.Target = applyHome
		_, err = sync.Apply(sourceRoot, entry, sync.Options{Force: true})
		Expect(err).NotTo(HaveOccurred())

		got, err := os.ReadFile(filepath.Join(applyHome, "settings.json"))
		Expect(err).NotTo(HaveOccurred())
		Expect(string(got)).To(Equal(`{"theme":"dark"}`))
	})

	It("should accurately compare drift between source and local", func() {
		local := filepath.Join(root, "local", ".claude")
		entry := config.AssistantEntry{
			Source: "assistants/claude",
			Target: local,
			Ignore: assistant.DefaultClaude().Ignore,
		}
		srcDir := filepath.Join(sourceRoot, entry.Source)
		Expect(os.MkdirAll(srcDir, 0o755)).To(Succeed())
		Expect(os.MkdirAll(local, 0o755)).To(Succeed())
		Expect(os.WriteFile(filepath.Join(srcDir, "a.txt"), []byte("source"), 0o644)).To(Succeed())
		Expect(os.WriteFile(filepath.Join(local, "b.txt"), []byte("local"), 0o644)).To(Succeed())

		report, err := sync.Compare(sourceRoot, entry)
		Expect(err).NotTo(HaveOccurred())
		Expect(report.Entries).To(HaveLen(2))
	})
})
