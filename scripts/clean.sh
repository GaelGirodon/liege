#!/usr/bin/env bash

set -e; trap 'echo -e "\033[1;36m$ $BASH_COMMAND\033[0m"' debug

if [[ -d ./bin ]]; then rm -r ./bin; fi
if [[ -d ./dist ]]; then rm -r ./dist; fi
