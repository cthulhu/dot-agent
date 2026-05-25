package assistant_test

import (
	"github.com/cthulhu/dot-agent/internal/assistant"
	"github.com/cthulhu/dot-agent/internal/config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Assistant", func() {
	Describe("Default Configurations", func() {
		It("Hermes should ignore secrets", func() {
			entry := assistant.DefaultHermes()
			Expect(entry.Ignore).To(ContainElements(".env", "auth.json"))
		})

		It("Codex should ignore secrets and have correct target", func() {
			entry := assistant.DefaultCodex()
			Expect(entry.Ignore).To(ContainElements("auth.json", "history.jsonl"))
			Expect(entry.Target).To(Equal("~/.codex"))
		})

		It("Gemini should ignore secrets and tmp, and have correct target", func() {
			entry := assistant.DefaultGemini()
			Expect(entry.Ignore).To(ContainElements(".env", "oauth_creds.json", "**/tmp/**"))
			Expect(entry.Target).To(Equal("~/.gemini"))
		})

		It("Copilot should ignore secrets and sessions, and have correct target", func() {
			entry := assistant.DefaultCopilot()
			Expect(entry.Ignore).To(ContainElements("config.json", "**/session-state/**"))
			Expect(entry.Target).To(Equal("~/.copilot"))
		})

		It("Antigravity should ignore secrets and local logs/history, and have correct target", func() {
			entry := assistant.DefaultAntigravity()
			Expect(entry.Ignore).To(ContainElements(".env", "installation_id", "history.jsonl", "**/brain/**", "**/knowledge/**", "**/*.lock"))
			Expect(entry.Target).To(Equal("~/.gemini/antigravity-cli"))
		})
	})

	Describe("Assistant Registry", func() {
		It("KnownNames should include all assistants", func() {
			names := assistant.KnownNames()
			Expect(names).To(ContainElements(
				assistant.Claude,
				assistant.Cursor,
				assistant.Hermes,
				assistant.Codex,
				assistant.Gemini,
				assistant.Copilot,
				assistant.Antigravity,
			))
		})

		It("MergeMissingAssistants should add missing assistants to manifest", func() {
			m := &config.Manifest{
				Version: 1,
				Assistants: map[string]config.AssistantEntry{
					assistant.Claude: assistant.DefaultClaude(),
				},
			}
			added := assistant.MergeMissingAssistants(m)
			Expect(added).To(HaveLen(6))
			Expect(m.Assistants).To(HaveKey(assistant.Gemini))
			Expect(m.Assistants).To(HaveKey(assistant.Copilot))
			Expect(m.Assistants).To(HaveKey(assistant.Antigravity))

			again := assistant.MergeMissingAssistants(m)
			Expect(again).To(BeEmpty())
		})
	})
})
