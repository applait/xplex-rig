xplex rig
============

Application controller plane for xplex

What it does:

- Manage user auth and user records.
- Manage streaming configurations per user.
- Authenticate incoming media streams and hands out config.
- Send out user notifications - email, social media etc.
- Manage rate limiting and auditing of user streams in edge clusters. Disconnect when beyond quota.

## Install

- [Install Golang 1.9+](https://golang.org/doc/install)
- Install [dep](https://golang.github.io/dep/docs/installation.html)

In project root, run:

```sh
$ dep ensure
```

Compile debug builds using `build.sh`:

```sh
$ ./build.sh dev
```

Compile static binaries for release profile:

```sh
$ ./build.sh release
```

Compiled binaries are put in `./bin/`.

## Run migrations

Migrations are run using [golang-migrate/migrate](https://github.com/golang-migrate/migrate). `./build.sh` script has a `migrate` command which downloads and calls the migrate CLI.

First, set the `DATABASE_URL` environment variable with the Postgres database URL:

```sh
export DATABASE_URL="postgres://user:name@host/db"
```

Then, run the migrations:

```sh
$ ./build.sh migrate up
```

This will install the `migrate` CLI if it doesn't exist in `./bin/migrate` and run the `up` migrations.

Undo all migrations:

```sh
$ ./build.sh migrate down
```

Run specific number of migrations up or down, e.g. for `1` migration:

```sh
$ ./build.sh migrate up 1
$ ./build.sh migrate down 1
```

