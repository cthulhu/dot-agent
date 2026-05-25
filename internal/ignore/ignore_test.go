package ignore_test

import (
	"path/filepath"

	"github.com/cthulhu/dot-agent/internal/ignore"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Ignore", func() {
	Describe("Matcher", func() {
		var matcher *ignore.Matcher

		BeforeEach(func() {
			matcher = ignore.New("**/cache/**", "**/*.log")
		})

		DescribeTable("Ignored paths",
			func(path string, expected bool) {
				Expect(matcher.Ignored(path)).To(Equal(expected))
			},
			Entry("cache/foo.txt", "cache/foo.txt", true),
			Entry("rules/my-rule.md", "rules/my-rule.md", false),
			Entry("debug/out.log", "debug/out.log", true),
		)

		It("should handle OS-native separators", func() {
			path := filepath.Join("cache", "foo.txt")
			Expect(matcher.Ignored(path)).To(BeTrue())
		})
	})

	Describe("Blocked", func() {
		DescribeTable("Blocked paths",
			func(path string, expected bool) {
				Expect(ignore.Blocked(path)).To(Equal(expected))
			},
			Entry(".env", ".env", true),
			Entry("subdir/credentials.json", "subdir/credentials.json", true),
			Entry("auth.json", "auth.json", true),
			Entry("settings.json", "settings.json", false),
		)
	})
})
