# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class Limesctl < Formula
  desc "Command-line interface for Limes"
  homepage "https://github.com/sapcc/limesctl"
  version "3.1.0"
  license "Apache-2.0"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/sapcc/limesctl/releases/download/v3.1.0/limesctl-3.1.0-darwin-arm64.tar.gz"
      sha256 "7e73e7d016b433bf58e86b1b82c8ac17feafc507d4a426e78efe58e042738815"

      def install
        bin.install "limesctl"
        bash_completion.install "completions/limesctl.bash" => "limesctl"
        zsh_completion.install "completions/limesctl.zsh" => "_limesctl"
        fish_completion.install "completions/limesctl.fish"
      end
    end
    if Hardware::CPU.intel?
      url "https://github.com/sapcc/limesctl/releases/download/v3.1.0/limesctl-3.1.0-darwin-amd64.tar.gz"
      sha256 "bdf95c759b67dd02b3014ee6cf65f07a8393de9f5f530dc492dabd4de3f7a8dd"

      def install
        bin.install "limesctl"
        bash_completion.install "completions/limesctl.bash" => "limesctl"
        zsh_completion.install "completions/limesctl.zsh" => "_limesctl"
        fish_completion.install "completions/limesctl.fish"
      end
    end
  end

  on_linux do
    if Hardware::CPU.arm? && Hardware::CPU.is_64_bit?
      url "https://github.com/sapcc/limesctl/releases/download/v3.1.0/limesctl-3.1.0-linux-arm64.tar.gz"
      sha256 "5b54f8ce378eeca2c1cb6c2b7be54a3d6eb9d6e53c3e057abf53ee2ff891f2c8"

      def install
        bin.install "limesctl"
        bash_completion.install "completions/limesctl.bash" => "limesctl"
        zsh_completion.install "completions/limesctl.zsh" => "_limesctl"
        fish_completion.install "completions/limesctl.fish"
      end
    end
    if Hardware::CPU.intel?
      url "https://github.com/sapcc/limesctl/releases/download/v3.1.0/limesctl-3.1.0-linux-amd64.tar.gz"
      sha256 "2420438e7cd9cf9e4939bda4dd1842814813b2c39062a7e8f617499e2a27bedf"

      def install
        bin.install "limesctl"
        bash_completion.install "completions/limesctl.bash" => "limesctl"
        zsh_completion.install "completions/limesctl.zsh" => "_limesctl"
        fish_completion.install "completions/limesctl.fish"
      end
    end
  end

  test do
    system "#{bin}/limesctl --version"
  end
end
