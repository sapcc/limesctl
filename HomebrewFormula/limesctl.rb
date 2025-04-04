# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class Limesctl < Formula
  desc "Command-line interface for Limes"
  homepage "https://github.com/sapcc/limesctl"
  version "3.5.0"
  license "Apache-2.0"

  on_macos do
    if Hardware::CPU.intel?
      url "https://github.com/sapcc/limesctl/releases/download/v3.5.0/limesctl-3.5.0-darwin-amd64.tar.gz"
      sha256 "e30fe450d2d7195849f271dc14686585e77a582a8c530cdc7c0942ce2ee2f791"

      def install
        bin.install "limesctl"
        bash_completion.install "completions/limesctl.bash" => "limesctl"
        fish_completion.install "completions/limesctl.fish"
        zsh_completion.install "completions/limesctl.zsh" => "_limesctl"
      end
    end
    if Hardware::CPU.arm?
      url "https://github.com/sapcc/limesctl/releases/download/v3.5.0/limesctl-3.5.0-darwin-arm64.tar.gz"
      sha256 "facc5b65a6f4b74927490a6e200d6841fb43fd5a89a7d87079013a39692749e3"

      def install
        bin.install "limesctl"
        bash_completion.install "completions/limesctl.bash" => "limesctl"
        fish_completion.install "completions/limesctl.fish"
        zsh_completion.install "completions/limesctl.zsh" => "_limesctl"
      end
    end
  end

  on_linux do
    if Hardware::CPU.intel?
      if Hardware::CPU.is_64_bit?
        url "https://github.com/sapcc/limesctl/releases/download/v3.5.0/limesctl-3.5.0-linux-amd64.tar.gz"
        sha256 "e8895f2348a37f280e84177c029caa752772545e36084c191c8895455887ae7d"

        def install
          bin.install "limesctl"
          bash_completion.install "completions/limesctl.bash" => "limesctl"
          fish_completion.install "completions/limesctl.fish"
          zsh_completion.install "completions/limesctl.zsh" => "_limesctl"
        end
      end
    end
    if Hardware::CPU.arm?
      if Hardware::CPU.is_64_bit?
        url "https://github.com/sapcc/limesctl/releases/download/v3.5.0/limesctl-3.5.0-linux-arm64.tar.gz"
        sha256 "924d2be78ed73f241e5c0ec3b9cd6144e575648a6d393747413d1c1b4a12d49b"

        def install
          bin.install "limesctl"
          bash_completion.install "completions/limesctl.bash" => "limesctl"
          fish_completion.install "completions/limesctl.fish"
          zsh_completion.install "completions/limesctl.zsh" => "_limesctl"
        end
      end
    end
  end

  test do
    system "#{bin}/limesctl --version"
  end
end
