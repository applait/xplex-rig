#!/bin/sh

set -e

function cleanRebuild {
    rm -rf bin/
    mkdir -p bin/
}

function buildServerRelease {
    echo "Building xplex-rig server release"
    CGO_ENABLED=0 GOOS=linux go build -o bin/rig-server -a -ldflags '-extldflags "-static"' ./server/
}

function buildServerDev {
    echo "Building xplex-rig server dev"
    go build -o bin/rig-server ./server/
}

function buildDev {
    cleanRebuild
    buildServerDev
}

function buildRelease {
    cleanRebuild
    buildServerRelease
}

case "$1" in
    "release")
    buildRelease
    ;;
    "dev")
    buildDev
    ;;
    *)
    echo "Usage: ./build.sh [dev|release|help]"
    echo "Binaries are put in './bin/'"
    ;;
esac
