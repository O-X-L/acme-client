#!/usr/bin/env bash

set -euo pipefail

cd "$(dirname "$0")/.."

export MODE_TEST=1

BASE_DIR="$(pwd)"

go run gotest.tools/gotestsum@latest --format pkgname ./...
