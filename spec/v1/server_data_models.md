# rig-server data models

- Version: v1
- Spec status: Experimental

These are the data models used by rig-server.

## StreamConfiguration

Stream routing configuration per stream key.

- `id` {number} [primary] - Programmatically generated StreamConfiguration id
- `key` {string} [unique] - Unique stream key
- `routeTo` {array} - Array of destinations to route stream to
    - `[]provider` {string} - Provider name, e.g., `YouTube Live`, `Twitch`
    - `[]streamKey` {string} - User's stream key for given service
- `allowPublishing` {boolean} - Whether stream key is allowed to publish. Set to `false` when quota has expired.
- `isLive` {boolean} - Whether stream is currently live
- `userId` {number} - `belongsTo` reference to user to whom this stream configuration belongs.

## User

User details

- `id` {number} [primary]- Programmatically generated user id
- `username` {string] [unique] - User's chosen username
- `password` string} - Hashed user password
- `salt` {string} - Generated salt for hash of the user

## Agent

Details of each agent.

- `hostname` {string} [primary] - Unique hostname for each agent
- `cluster` {string} - Hostname of cluster the agent belongs to.
- `region` {string} - Region where the agent is located
- `isActive` {boolean} - Whether agent is active or not
- `capabilities` {array<`string`>} - List of capabilities for each agent, useful for future purposes. Currently, everyone has `RTMP` capability.
