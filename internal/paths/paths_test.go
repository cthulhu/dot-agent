package paths_test

import (
	"path/filepath"
	"runtime"
	"strings"

	"github.com/cthulhu/dot-agent/internal/paths"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Paths", func() {
	It("should expand home directory paths", func() {
		home, err := paths.HomeDir()
		Expect(err).NotTo(HaveOccurred())

		got, err := paths.ExpandPath("~/.claude")
		Expect(err).NotTo(HaveOccurred())

		want := filepath.Join(home, ".claude")
		Expect(got).To(Equal(want))
	})

	It("should provide a default source directory", func() {
		dir, err := paths.DefaultSourceDir()
		Expect(err).NotTo(HaveOccurred())
		Expect(dir).To(ContainSubstring("dot-agent"))

		if runtime.GOOS == "windows" {
			Expect(strings.ToLower(dir)).To(ContainSubstring("local"))
		}
	})

	It("should provide a configuration directory", func() {
		dir, err := paths.ConfigDir()
		Expect(err).NotTo(HaveOccurred())
		Expect(dir).To(ContainSubstring("dot-agent"))
	})
})
