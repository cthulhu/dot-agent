class DotAgent < Formula
  desc "Sync AI coding assistant configuration across machines using git"
  homepage "https://github.com/cthulhu/dot-agent"
  url "https://github.com/cthulhu/dot-agent/archive/refs/tags/v0.8.4.tar.gz"
  # sha256 "TODO"
  license "MIT"
  head "https://github.com/cthulhu/dot-agent.git", branch: "main"

  depends_on "go" => :build
  depends_on "git"

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w"), "./cmd/dot-agent"
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/dot-agent --version")
    assert_match "dot-agent", shell_output("#{bin}/dot-agent --help")
  end
end
