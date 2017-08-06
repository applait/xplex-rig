# rig-agent HTTP API

rig-agents expose a small HTTP 1.1 API through which rig-server and load balancers could interact with them. This document states the REST endpoints exposed by rig-agents:

## Heartbeat

Respond to heartbeat requests

- **URL:** `/heartbeat`
- **Method:** `GET`
- **Success Response:**

    - Status Code: `200`
    - Content:
        - `msg` {string} - `"Heartbeat"`
        - `status` {number} - `200`
        - `payload` {object}
            - `timestamp` {number} - Current Unix timestamp

    - Example:
        ```json
        "msg": "Heartbeat",
        "status": 200,
        "payload": {
            "timestamp": 1502033100533
        }
        ```
