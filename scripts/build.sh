#!/usr/bin/env bash

if [[ -z "$MODE_TEST" ]]
then
  MODE_TEST=0
fi

set -euo pipefail

cd "$(dirname "$0")/.."
PATH_OUT="$(pwd)/build"
mkdir -p "$PATH_OUT"
cd ./src
CGO_ENABLED=0 go build -o "${PATH_OUT}/acme" ./cmd/main.go

echo ''
echo '### DONE ###'
echo ''

ls -l "$PATH_OUT"

echo ''
echo ''
