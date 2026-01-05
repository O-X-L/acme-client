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

docker build -f Dockerfile -t "oxlorg/acme-client:${VERSION}" --build-arg "VERSION=${VERSION}" --no-cache --network=host --progress=plain .
docker build -f Dockerfile -t "oxlorg/acme-client:latest" --build-arg "VERSION=${VERSION}" --network=host --progress=plain .
