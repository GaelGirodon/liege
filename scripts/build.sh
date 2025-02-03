#!/usr/bin/env bash

if [[ ! "$1" =~ ^(windows|linux)$ ]];
then echo -e "\033[1;31mUsage: build.sh <linux|windows>\033[0m"; exit 1; fi

set -e; trap 'echo -e "\033[1;36m$ $BASH_COMMAND\033[0m"' debug

export GOOS="$1"
export GOARCH="amd64"

ext=$([[ "$GOOS" == "windows" ]] && echo ".exe" || echo "")
go build -ldflags="-s -w" -o "./bin/liege$ext" ./cmd/liege.go
