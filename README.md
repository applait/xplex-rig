xplex rig
============

Internal API server for xplex.

## Pre-requisites

- `build-essentials`
- NodeJS 6.9.0+
- MongoDB 3.4 (get it running)
- `node-gyp`: `npm install -g node-gyp`

## Install

- Copy `config.sample.js` to `config.js` and edit values, with DB credentials.
- Install application dependencies: `npm install`
- Start server: `npm start`
- Start agent: `npm run agent`

## Tests

- Run: `npm test`
