#!/usr/bin/env bash

set -e; trap '[[ ! $BASH_COMMAND =~ (read|\[\[) ]] && echo -e "\033[1;36m$ $BASH_COMMAND\033[0m"' debug

old_version=$(grep -m 1 -oP '(?<=Version = ")[0-9.]+(?=")' ./internal/console/console.go)
v=(${old_version//./ })
new_version="${v[0]}.$((v[1] + 1)).${v[2]}"
sed -i "s/$old_version/$new_version/g" ./internal/console/console.go
read -rp "Update the changelog..."
git add CHANGELOG.md ./internal/console/console.go
read -rp "Commit and tag? [Yn] " yn && [[ $yn != Y* ]] && exit 0
git commit -m "Release $new_version"
git tag -a "$new_version" -m "$new_version"
