#!/usr/bin/env bash

set -eu
[ "${BASH_VERSINFO[0]}" -ge 3 ] && set -o pipefail

DIR=$(dirname "$0")
ROOTDIR=$(cd "$DIR/../" && pwd )
GORELEASER_VERSION=0.164.0
GOLANGCI_LINT_VERSION=1.40.1

os=$(uname -s)
case "$os" in
  Darwin*)
        os=darwin
        ;;
  Linux*)
        os=linux
        ;;
  *)
        echo "unsupported: $os"
        exit 1
        ;;
esac

__install_golangci_lint() {
    local asset="golangci-lint-${GOLANGCI_LINT_VERSION}-${os}-amd64.tar.gz"
    local url="https://github.com/golangci/golangci-lint/releases/download/v${GOLANGCI_LINT_VERSION}/golangci-lint-${GOLANGCI_LINT_VERSION}-${os}-amd64.tar.gz"
    echo "Download $url"

    curl -fsLJO "$url"
    tar -C "${ROOTDIR}"/.tools/bin --strip-components=1 -zxf ${asset} "golangci-lint-${GOLANGCI_LINT_VERSION}-${os}-amd64/golangci-lint"
    rm -rf ${asset}
}

__install_goreleaser() {
    local os
    os="$(uname -s)"
    local asset="goreleaser_${os}_x86_64.tar.gz"
    local url="https://github.com/goreleaser/goreleaser/releases/download/v${GORELEASER_VERSION}/${asset}"
    echo "Download $url"

    curl -fsLJO "$url"
    tar -C "${ROOTDIR}"/.tools/bin -zxf "${asset}"
    rm -rf "${asset}"
}

__install_gotools() {
    go install golang.org/x/tools/cmd/goimports
}

rm -rf "$ROOTDIR/.tools"
mkdir -p "$ROOTDIR/.tools/bin"

__install_golangci_lint

__install_goreleaser

__install_gotools

