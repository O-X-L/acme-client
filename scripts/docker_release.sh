#!/usr/bin/env bash

if [ -z "$1" ]
then
  echo "USAGE:"
  echo " 1 > Version"
  exit 1
fi

set -euo pipefail

VERSION="$1"

cd "$(dirname "$0")/../docker"

docker push "oxlorg/acme-client:${VERSION}"
docker push "oxlorg/acme-client:latest"
