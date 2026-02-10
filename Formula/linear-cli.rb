class LinearCli < Formula
  desc "Comprehensive command-line interface for Linear's API"
  homepage "https://github.com/roboalchemist/linear-cli"
  url "https://github.com/roboalchemist/linear-cli/archive/refs/tags/v0.3.3.tar.gz"
  sha256 "490453c479b425b9fb2081d99b3baf7fca645985542d8232bf6da9d12d7aa8f8"
  license "MIT"
  head "https://github.com/roboalchemist/linear-cli.git", branch: "master"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w -X github.com/roboalchemist/linear-cli/cmd.version=#{version}")
  end

  test do
    # Test version output
    assert_match "linear-cli version #{version}", shell_output("#{bin}/linear-cli --version")

    # Test help command
    assert_match "A comprehensive CLI tool for Linear", shell_output("#{bin}/linear-cli --help")
  end
end
