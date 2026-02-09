class LinearCli < Formula
  desc "Comprehensive command-line interface for Linear's API"
  homepage "https://github.com/roboalchemist/linear-cli"
  url "https://github.com/roboalchemist/linear-cli/archive/refs/tags/v0.1.1.tar.gz"
  sha256 "0cedf76b65a52223c35301419ae1ac45e86f1c02c8ccbc7702f0c1f6c61c1f76"
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
