xplex rig
============

Internal API server for xplex.

## Pre-requisites

- NodeJS 8.2+
- PostgreSQL 9.6+ (get it running)
- Etcd store

## Install

- Copy `config.sample.js` to `config.js` and edit values, with DB credentials.
- Install application dependencies: `npm install`
- Start server: `npm start`
- Start agent: `npm run agent`

## Tests

- Run: `npm test`

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
