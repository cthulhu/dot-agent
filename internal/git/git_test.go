package git_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/cthulhu/dot-agent/internal/git"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Git", func() {
	var tmpDir string

	BeforeEach(func() {
		var err error
		tmpDir, err = os.MkdirTemp("", "dot-agent-git-test-*")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		Expect(os.RemoveAll(tmpDir)).To(Succeed())
	})

	It("should verify git is in PATH", func() {
		Expect(git.RequireGit()).To(Succeed())
	})

	Describe("Repository Operations", func() {
		var repoDir string

		BeforeEach(func() {
			repoDir = filepath.Join(tmpDir, "repo")
			Expect(git.Init(repoDir)).To(Succeed())

			// Configure local git user for commits to work
			_, _ = gitRun(repoDir, "config", "user.email", "test@example.com")
			_, _ = gitRun(repoDir, "config", "user.name", "Test User")
		})

		It("should initialize a git repository", func() {
			Expect(filepath.Join(repoDir, ".git")).To(BeADirectory())
		})

		It("should report dirty status and handle commits", func() {
			By("checking initial status")
			Expect(git.IsDirty(repoDir)).To(BeFalse())

			By("adding a new file")
			Expect(os.WriteFile(filepath.Join(repoDir, "test.txt"), []byte("hello"), 0o644)).To(Succeed())
			Expect(git.IsDirty(repoDir)).To(BeTrue())

			By("adding all changes")
			Expect(git.AddAll(repoDir)).To(Succeed())
			Expect(git.StatusPorcelain(repoDir)).To(ContainSubstring("A  test.txt"))

			By("committing changes")
			err := git.Commit(repoDir, "initial commit")
			if err != nil && strings.Contains(err.Error(), "identity unknown") {
				Skip("Git user not configured, skipping commit test")
			}
			Expect(err).NotTo(HaveOccurred())
			Expect(git.IsDirty(repoDir)).To(BeFalse())
		})

		It("should report the current branch", func() {
			By("creating an initial commit")
			Expect(os.WriteFile(filepath.Join(repoDir, "init.txt"), []byte("init"), 0o644)).To(Succeed())
			Expect(git.AddAll(repoDir)).To(Succeed())
			err := git.Commit(repoDir, "init")
			if err != nil && strings.Contains(err.Error(), "identity unknown") {
				Skip("Git user not configured, skipping branch test")
			}
			Expect(err).NotTo(HaveOccurred())

			branch, err := git.CurrentBranch(repoDir)
			Expect(err).NotTo(HaveOccurred())
			Expect(branch).To(Or(Equal("main"), Equal("master")))
		})

		It("should manage remotes", func() {
			remoteURL := "https://github.com/cthulhu/dot-agent.git"

			By("setting a remote")
			Expect(git.SetRemote(repoDir, remoteURL)).To(Succeed())

			By("getting the remote URL")
			url, err := git.RemoteURL(repoDir)
			Expect(err).NotTo(HaveOccurred())
			Expect(url).To(Equal(remoteURL))

			By("updating the remote URL")
			newURL := "https://github.com/cthulhu/new-repo.git"
			Expect(git.SetRemote(repoDir, newURL)).To(Succeed())
			url, err = git.RemoteURL(repoDir)
			Expect(err).NotTo(HaveOccurred())
			Expect(url).To(Equal(newURL))
		})
	})

	Describe("Clone", func() {
		It("should clone a local repository", func() {
			srcDir := filepath.Join(tmpDir, "src")
			destDir := filepath.Join(tmpDir, "dest")

			By("creating a source repo")
			Expect(git.Init(srcDir)).To(Succeed())
			_, _ = gitRun(srcDir, "config", "user.email", "test@example.com")
			_, _ = gitRun(srcDir, "config", "user.name", "Test User")

			Expect(os.WriteFile(filepath.Join(srcDir, "README.md"), []byte("# Test"), 0o644)).To(Succeed())
			Expect(git.AddAll(srcDir)).To(Succeed())

			err := git.Commit(srcDir, "initial commit")
			if err != nil && strings.Contains(err.Error(), "identity unknown") {
				Skip("Git user not configured, skipping clone test")
			}
			Expect(err).NotTo(HaveOccurred())

			By("cloning the source repo")
			Expect(git.Clone(srcDir, destDir)).To(Succeed())
			Expect(filepath.Join(destDir, ".git")).To(BeADirectory())
			Expect(filepath.Join(destDir, "README.md")).To(BeAnExistingFile())
		})
	})
})

func gitRun(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	return string(out), err
}
