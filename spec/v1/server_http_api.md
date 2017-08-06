# rig-server HTTP API

Version: v1

rig-server HTTP API is the interface for frontends and rig-agents to interact with rig. This document states the endpoints exposed by the rig-server HTTP API along with their desired purpose.

## Users

### Create new User

Create a new user for xplex.

- **URL:** `/users`
- **Method:** `POST`
- **Data params:**

    **Required**:

    `username=<string>`: User's chosen username. Must be unique.

    `email=<string>` : User's email address

    `password=<string>` : User's chosen password

    **Optional**:

    `send_notification=[true|false]` Default: `true` : Whether to send email notifications to user on signup

    `auto_confirm=[true|false]` Default: `false` : Whether to auto confirm user. Default `false`.
