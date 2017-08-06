# xplex-rig v1 API spec

This directory contains specifications of version 1 of xplex-rig's API.

## HTTP APIs

Most of the components communicate over HTTP 1.1. When a HTTP request is involved, this is the generic structure followed for request and responses.

### Request

`POST` requests contain data in request body. `GET` requests contain data as URL parameters.

In general:

- `GET` is used to retrieve resource
- `POST` is used to create and update resources.
- `DELETE` is used to delete data.


### Response

All API responses are returned as `Content-Type: application/json`, with this general data structure:

```js
{
    "msg": String, // A unique message per endpoint action, used to describe current response. e.g., "User created".
    "status": Number, // HTTP status code. `200` for all successful requests.
    "payload": Array || Object || null // The data returned for the request. If response does not have any associated data to be returned, this is set to `null`.
}
```

#### Error response

Any error response has this structure:

- Status Code: `4xx` or `5xx`
- Content:
    - `msg` {string} - Generic error message
    - `status` {number} - Same as HTTP status code
    - `payload` {array<string>} - Array of strings for each error encountered in detail

For example:

```json
"msg": "Cannot create user",
"status": 400,
"payload": [
    "username field is required",
    "email is malformed
]
```
