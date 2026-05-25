class DotAgent < Formula
  desc "Sync AI coding assistant configuration across machines using git"
  homepage "https://github.com/cthulhu/dot-agent"
  url "https://github.com/cthulhu/dot-agent/archive/refs/tags/v0.5.0.tar.gz"
  sha256 "8e6854d046e71b68340f738cffb6d76e6d0bb8ce34b63b17c6adc1c765f6ca26"
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
