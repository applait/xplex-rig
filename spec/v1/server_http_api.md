# rig-server HTTP API

Version: v1

rig-server HTTP API is the interface for frontends and rig-agents to interact with rig. This document states the endpoints exposed by the rig-server HTTP API along with their desired purpose.


## Users

### Create new User

Create a new user for xplex.

- **URL:** `/users`
- **Method:** `POST`

- **Data params:**
    - **Required**:
        - `username=<string>`: User's chosen username. Must be unique.
        - `email=<string>` : User's email address
        - `password=<string>` : User's chosen password
    - **Optional**:
        - `send_notification=[true|false]` Default: `true` : Whether to send email notifications to user on signup
        - `auto_confirm=[true|false]` Default: `false` : Whether to auto confirm user. Default `false`.

- **Success Response:**

    - Status Code: `200`
    - Content:
        - `msg` {string} - `"User created"`
        - `status` {number} - `200`
        - `payload` {object}
            - `id` {number} - Generated user id
            - `username` {string} - Chosen username
            - `notification_sent` {boolean} - Whether signup notification was dispatched or not
            - `confirmation_required` {boolean} - Whether user needs to confirm account to activate
    - Example:
        ```json
        "msg": "User Created",
        "status": 200,
        "payload": {
            "id": 42,
            "username": "chosenusername",
            "notification_sent": true,
            "confirmation_required": true
        }
        ```

### Get user information by user ID

Get information about existing user by user ID

- **URL:** `/users/:id`
- **Method:** `GET`
- **Success Response:**

    - Status Code: `200`
    - Content:
        - `msg` {string} - `"User information"`
        - `status` {number} - `200`
        - `payload` {object}
            - `id` {number} - User id
            - `username` {string} - username
    - Example:
        ```json
        "msg": "User Information",
        "status": 200,
        "payload": {
            "id": 42,
            "username": "chosenusername"
        }
        ```
