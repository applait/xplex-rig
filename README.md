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
- Install [Glide](https://github.com/Masterminds/glide)

In project root, run:

```sh
$ glide install
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

