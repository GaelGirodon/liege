#!/usr/bin/env bash

APP_NAME=liege
MAIN_FILE=cmd/liege.go
VERSION_FILE=internal/console/console.go
BIN_PATH=bin
BIN_FILE=bin/liege
export GOARCH="amd64"

build() {
    build_linux && build_windows
}

build_linux() {
    export GOOS=linux
    go build -ldflags="-s -w" -o "$BIN_FILE" "$MAIN_FILE"
}

build_windows() {
    export GOOS=windows
    go build -ldflags="-s -w" -o "${BIN_FILE}.exe" "$MAIN_FILE"
}

test() {
    go test ./... -v
}

version() {
    grep -m 1 -oP '(?<=Version = ")[0-9.]+(?=")' "$VERSION_FILE"
}

release() {
    old_v=$(version)
    v=(${old_v//./ })
    new_v="${v[0]}.$((v[1] + 1)).${v[2]}"
    echo "Bump version number from $old_v to $new_v"
    sed -i "s/$old_v/$new_v/g" "$VERSION_FILE"
    read -rp "Update the changelog..."
    git add "$VERSION_FILE" CHANGELOG.md
    read -rp "Commit and tag? [Yn] " yn && [[ $yn != Y* ]] && return 0
    git commit -m "Release $new_v"
    git tag -a "$new_v" -m "$new_v"
    echo "Push -> 'git push --follow-tags'"
}

docker_build() {
    docker build --pull -t "gaelgirodon/${APP_NAME}" .
    docker tag "gaelgirodon/${APP_NAME}" "gaelgirodon/${APP_NAME}:$(version)"
}

docker_push() {
    docker push "gaelgirodon/${APP_NAME}:$(version)"
    docker push "gaelgirodon/${APP_NAME}:latest"
}

clean() {
    if [[ -d "$BIN_PATH" ]]; then rm -ri "$BIN_PATH"; fi
}
