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

function buildAgentRelease {
    echo "Building xplex-rig agent release"
    CGO_ENABLED=0 GOOS=linux go build -o bin/rig-agent -a -ldflags '-extldflags "-static"' ./agent/
}

function buildServerDev {
    echo "Building xplex-rig server dev"
    go build -o bin/rig-server ./server/
}

function buildAgentDev {
    echo "Building xplex-rig agent dev"
    go build -o bin/rig-agent ./agent/
}

function buildDev {
    cleanRebuild
    buildServerDev
    buildAgentDev
}

function buildRelease {
    cleanRebuild
    buildServerRelease
    buildAgentRelease
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
