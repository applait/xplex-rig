#!/bin/sh

set -e

cleanRebuild () {
    rm -f bin/xplex-rig
    mkdir -p bin/
}

buildServerRelease () {
    echo "Building xplex-rig server release"
    CGO_ENABLED=0 GOOS=linux go build -o bin/xplex-rig -a -ldflags '-extldflags "-static"' .
    echo "Compiled to './bin/xplex-rig'"
}

buildServerDev () {
    echo "Building xplex-rig server dev"
    go build -o bin/xplex-rig .
    echo "Compiled to './bin/xplex-rig'"
}

buildDev () {
    cleanRebuild
    buildServerDev
}

buildRelease () {
    cleanRebuild
    buildServerRelease
}

migrate () {
  if [ -z ${DATABASE_URL+x} ] || [ -z ${DATABASE_URL} ]; then
    echo "Set Postgresql connection URI in the 'DATABASE_URL' environment variable"
    exit 1
  fi
    if [ ! -f ./bin/migrate ]; then
        echo "Downloading migrate cli"
        curl -Ls https://github.com/golang-migrate/migrate/releases/download/v3.2.0/migrate.linux-amd64.tar.gz | tar xz
        mv ./migrate.linux-amd64 ./bin/migrate
        echo "Installed migrate cli at ./bin/migrate"
    fi
    opts="--path ./migrations/ --database ${DATABASE_URL}"
    ./bin/migrate $opts $@
}

case "$1" in
    "release")
    buildRelease
    ;;
    "dev")
    buildDev
    ;;
    "migrate")
    migrate $2 $3 $4 $5
    ;;
    *)
    echo "Usage: ./build.sh [dev|release|migrate [up|down] [num]|help]"
    echo ""
    echo "For using 'migrate', export database connection string as DATABASE_URL variable, e.g.:"
    echo "  $ export DATABASE_URL=postgres://user:pass@host/db"
    echo "Arguments given to './build.sh migrate' are passed on to the migrate CLI tool"
    echo ""
    echo "Binaries are put in './bin/'"
    ;;
esac
