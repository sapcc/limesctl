#!/bin/sh
set -euo pipefail

OS_TYPE=$(uname -s)
if ! [ "$OS_TYPE" == "Darwin" ]  && ! [ "$OS_TYPE" == "Linux" ]; then
	printf "\e[1;31m==> This script only works on macOS and Linux, '${OS_TYPE}' is not supported\e[0m\n"
	exit 1
fi

# Borrowed from https://gist.github.com/lukechilds/a83e1d7127b78fef38c2914c4ececc3c
get_latest_release() {
	curl --silent "https://api.github.com/repos/$1/releases/latest" | # Get latest release from GitHub API
	grep '"tag_name":' | # Get tag line
	sed -E 's/.*"v([^"]+)".*/\1/' # Pluck version number
}
VERSION="$(get_latest_release 'sapcc/limesctl')"

TEMP_DIR="$(mktemp -d)"

download() {
	local binary_archive="limesctl-${VERSION}-${OS_TYPE}_amd64.tgz"
	curl -L "https://github.com/sapcc/limesctl/releases/download/v${VERSION}/${binary_archive}" -o ${TEMP_DIR}/limesctl.tgz
	tar -xzf ${TEMP_DIR}/limesctl.tgz -C $TEMP_DIR
}

install() {
	sudo mv -f ${TEMP_DIR}/limesctl /usr/local/bin/limesctl
}

cleanup() {
	rm -rf $TEMP_DIR
}

main() {
	printf "\e[1;34m==> Downloading limesctl for $(uname -s)\e[0m\n"
	download
	printf "\e[1;34m==> Installing limesctl\e[0m\n"
	install
	printf "\e[1;32m==> limesctl v${VERSION} successfully installed as $(which limesctl)\e[0m\n"
	cleanup
}

main
