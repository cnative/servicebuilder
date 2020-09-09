#!/usr/bin/env bash

set -eu
[ "${BASH_VERSINFO[0]}" -ge 3 ] && set -o pipefail

DIR=$(dirname "$0")
ROOTDIR=$(cd "$DIR/../" && pwd )

if [ -r "$ROOTDIR/.tools/checksum.txt" ]; then
    install_checksum=$(cksum "$ROOTDIR/scripts/install_tools.sh")
    current_checksum=$(cat "$ROOTDIR/.tools/checksum.txt")
    if [ "$install_checksum" == "$current_checksum" ]; then
        exit 0
    fi
fi

"$ROOTDIR/scripts/install_tools.sh" # this will remove the current .tools folder if present and install fresh
cksum "$ROOTDIR/scripts/install_tools.sh" > "$ROOTDIR/.tools/checksum.txt"
