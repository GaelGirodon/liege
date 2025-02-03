#!/usr/bin/env bash

set -e; trap 'echo -e "\033[1;36m$ $BASH_COMMAND\033[0m"' debug

version=$(grep -m 1 -oP '(?<=Version = ")[0-9.]+(?=")' ./internal/console/console.go)
docker build --pull -t gaelgirodon/liege .
docker tag gaelgirodon/liege "gaelgirodon/liege:$version"
