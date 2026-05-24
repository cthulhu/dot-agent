class DotAgent < Formula
  desc "Sync AI coding assistant configuration across machines using git"
  homepage "https://github.com/cthulhu/dot-agent"
  url "https://github.com/cthulhu/dot-agent/archive/refs/tags/v0.1.0.tar.gz"
  sha256 "3d9d69b4796a973a9fb5610d59c3cde719923194670e687981146aa5fb92d13e"
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
