xplex rig
============

Internal API server for xplex.

[![CircleCI](https://circleci.com/gh/applait/xplex-rig/tree/master.svg?style=svg&circle-token=e8251cf207d23cb87a38799f6e70f165bc5a1ad7)](https://circleci.com/gh/applait/xplex-rig/tree/master)

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

## Specifications

Interface specifications are present in the respective versioned directory:
 - [**v1 specification**](spec/v1)

## Components

- HTTP API for communication in both rig-server and rig-agent. TLS for all communications.
- rig-server controller and rig-server workers exchange tasks through etcd v2.
- rig-agent collects stats from local nginx-rtmp server and pushes periodically to etcd v2.
- Edge nodes are built of nginx-rtmp, rig-agent and ffmpeg. Rig-agent controls incoming streams. Spins up rig-agent workers - `ffmpeg` processes that pull from nginx and push to specified routes.

### Rig-server

Main controller of business logic internally.

- HTTP API for frontends and agents. TLS for all communications.
- etcd v2 connection to coordinate between rig-agent and rig-server-workers
- Built in 2 parts:
    - Controller: Handles request/response of the HTTP API. Responds to incoming queries. Delegates non-realtime, asynchronous or long-running tasks to workers via etcd v2.
    - Worker: Asynchronous workers that consume etcd v2 task queue and keep on performing tasks. Use for periodic tasks, notifications and quota checks.

What it does:

- Manage user auth and user records. Is the OAuth2 provider.
- Manage streaming configurations per user. Send when requested.
- Authenticate incoming media streams.
- Send out user notifications - email, social media etc.
- Manage rate limiting and auditing of user streams in edge clusters. Disconnect when beyond quota.
- (Near future) Control edge clusters and nodes. Scale capacity up and down as needed.

### Rig-agent

- HTTP API for communicating with rig-server
- Perform stream authentication by looking up data with rig-server
- Spin up and control media worker processes (`ffmpeg`)
- Report slot availability and usage status via etcd.
