class DotAgent < Formula
  desc "Sync AI coding assistant configuration across machines using git"
  homepage "https://github.com/cthulhu/dot-agent"
  url "https://github.com/cthulhu/dot-agent/archive/refs/tags/v0.4.1.tar.gz"
  sha256 "9ef70f70147576819d95bbfe9f936a7154e81fa16347928039a012aea176e244"
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
