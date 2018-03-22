# rig HTTP API v1

rig's HTTP API is the primary way through which clients and agents communicate with rig. It is a JSON-based RESTful API with authorization implemented using JWT.

## Table of contents

- [**Introduction**](#introduction)
    - [Request format](#request-format)
    - [Response format](#response-format)
    - [Authorization](#authorization)
- [**Accounts API**](#accounts-api)
- [**Streams API**](#streams-api)
- [**Agents API**](#agents-api)

# Introduction

## Request format

All requests are required to have `Content-Type: application/json`.

## Response format

All responses, success or error, return JSON responses in their body.

### Success responses

Success responses use HTTP status code `200` and carry this structure in their body:

- `message`: `string`. A message explaining the result of the operation.
- `payload`: `object|null`. The actual data returned by the operation. `payload` can be null if there are no data required to be returned. When present, you can safely parse `payload` as valid JSON. It can be any valid JSON data type.

Responses documented later in this document mention the structure of the `payload` property of the response body.

### Error responses

Error responses carry HTTP status codes `4xx` or `5xx` depending on the nature of error encountered. They return JSON responses in their body with information about the specific. Error response bodies contain this structure:

- `message`: `string`. A message about the error in general.
- `errorCode`: `number`. A numeric error code that denotes the type of error. These are described in [`rest/errors.go`](rest/errors.go). Will be `1000` and up.
- `details`: `object|null`. If the error has anything specific to return, it is provided here. It can be `null` if there is nothing specific to return.

An example error body looks like this:

```json
{
    "message": "Invalid credentials",
    "errorCode": 1003,
    "details": null
}
```

#### Error codes

Here is a list of error codes that are returned:

- `1000`: Not found. Can indicate the requested resource is not known or defined.
- `1001`: Method not allowed. Used when you use an undefined HTTP method on a known resource. e.g., if you use `GET` where you should use `POST`.
- `1002`: Content-Type should be `application/json`. Used when incoming request does not match Content-Type header.
- `1003`: Invalid credentials. Can also mean given access token does not have the scope to access requested resource.
- `1004`: Missing required input. Used when a required input parameter was not specified in the body of the request.
- `1005`: Invalid input. Used when a provided input parameter does not pass validation criteria.
- `1006`: Unable to create resource. Used when a resource creation request failed. This can be due to bad input data (e.g., duplicate username when creating new account) or some server side error.
- `1007`: Unable to update resource. Used when a resource failed to update either because of data conflict or server side error.
- `1008`: Unable to fetch resource. Used when the server is not able to send the requested resource, usually due to server side issues.
- `1009`: Resource is already in use. Used when a requested resource is being used elsewhere and cannot be shared. A typical scenario for this is when a stream with a given stream key is already live and user tries to stream from another location.

## Authorization

Authorization uses "Bearer" token format. For resources that require authorization, you will need to get an access token by authenticating with the respective auth endpoints. Once you have the token, set the `Authorization` HTTP header with the format `Bearer xxxx` where `xxxx` is the token string. e.g.,

```
Authorization: Bearer youractualtokenstringhere
```

# Clients

Clients are user facing applications. This includes web app, desktop apps or any other app through which a user may access xplex. Authorizations for client apps need user to authenticate to access endpoints that require authorization.

# Accounts API

## Authenticate user

Authenticate user using local authentication strategy, with username and password. This will give you the token you need to access rest of the restricted client APIs.

- URL: `/accounts/auth/local`
- Method: `POST`
- Authorization required: No

**Request params:**

- `username`: `string`. Unique. Minimum 3 characters.
- `password`: `string`. Minimum 8 characters.

**Success response payload:**

- `username`: `string`. User's chosen username.
- `token`: `string`. Access token to use for requests that require authorization.

**Example request body:**

```json
{
    "username": "foobar",
    "password": "Hop0nTh4tM4ng0Tree"
}
```

**Example success response:**

```json
{
    "message": "Authentication successful",
    "payload": {
        "username": "foobar",
        "token": "xxxxxxxx.xxxxxxxx.xxxxxxx"
    }
}
```

## Create account

Create a new user account.

- URL: `/accounts/`
- Method: `POST`
- Authorization required: No

**Request params:**

- `username`: `string`. Unique. Minimum 3 characters.
- `password`: `string`. Minimum 8 characters.
- `email`: `string`. Unique

**Success response payload:**

- `userID`: `string`. Newly created user's ID. UUID v4.
- `username`: `string`. User's chosen username.
- `email`: `string`. User's chosen email.

**Example request body:**

```json
{
    "username": "foobar",
    "password": "Hop0nTh4tM4ng0Tree",
    "email": "foobar@fancycompany.net"
}
```

**Example success response:**

```json
{
    "message": "User created",
    "payload": {
        "userID": "7217b0b6-d735-4ba9-bda8-419b1b397ede",
        "username": "foobar",
        "email": "foobar@fancycompany.net"
    }
}
```

## Change password

Used to change a logged in user's password. This is not forgot password endpoint.

- URL: `/accounts/password`
- Method: `POST`
- Authorization required: Yes

**Request params:**

- `oldPassword`: `string`. Current password of the user. Minimum 8 characters.
- `newPassword`: `string`. New password. Minimum 8 characters.

**Success response payload:**

- `username`: `string`. User's chosen username.

**Example request body:**

```json
{
    "oldPassword": "Hop0nTh4tM4ng0Tree",
    "newPassword": "Ar3YouComingToTheTree"
}
```

**Example success response:**

```json
{
    "message": "User password changed",
    "payload": {
        "username": "foobar"
    }
}
```

# Streams API

## Create new stream

Create a new stream configuration for the current user and returns the stream data.

- URL: `/streams/`
- Method: `POST`
- Authorization required: Yes

**Request params:**

None.

**Success response payload:**

- `id`: `string`. Stream's ID. This will be used for all stream related operations.
- `streamKey`: `string`. The stream key that user puts in the URL of the stream.
- `destinations`: `null`. Initially there will be no destinations set for the stream.
- `provisionStatus`: `null`. Initially the stream will not be provisioned.
- `isStreaming`: `boolean`. Initially `false`.
- `isActive`: `boolean`. Initially `true`.

**Example request body:**

None.

**Example success response:**

```json
{
    "message": "New stream created",
    "payload": {
        "id": "30960952-cf23-46f8-bf91-9f738c27be2e",
        "streamKey": "_mYrirpvU-OkH5LoVDci4g",
        "destinations": null,
        "provisionStatus": null,
        "isStreaming": false,
        "isActive": true
    }
}
```

## Add stream destination

Add a new egress destination for a stream. This should be used by clients when they have stream key for a given service already available. This is not an OAuth interface.

- URL: `/streams/{streamID}/destination`
- Method: `POST`
- Authorization required: Yes

**Request params:**

- `service`: `string`. A valid value for a service supported by xplex. See [supported services][#supported-services] for detail.
- `streamKey`: `string`. The streaming key provided by that service.

**Success response payload:**

- `id`: `number`. Destinantion ID
- `service`: `string`. The service value matching the request.
- `streamKey`: `string`. The streaming key for that service matching the request.
- `rtmpUrl`: `string`. A fully formed RTMP URL of the ingestion server of that service.
- `isActive`: `boolean`. Initially `true`.

**Example request body:**

```json
{
    "service": "YouTube",
    "streamKey": "simaimsimim_siimimimiasdsr"
}
```

**Example success response:**

```json
{
    "message": "Stream destination added",
    "payload": {
        "id": 1,
        "service": "YouTube",
        "streamKey": "simaimsimim_siimimimiasdsr",
        "rtmpUrl": "rtmp://a.rtmp.youtube.com/live2/simaimsimim_siimimimiasdsr",
        "isActive": true
    }
}
```

## Change stream key

Change streaming key for a an existing xplex stream. This changes only the stream key used in xplex's streaming URL and does not change keys for any destinations for that stream.

- URL: `/streams/{streamID}/changeKey`
- Method: `POST`
- Authorization required: Yes

**Request params:**

None.

**Success response payload:**

- `key`: `string`. The new stream key.
- `streamID`: `string`. The stream ID for which the key was changed.

**Example request body:**

None.

**Example success response:**

```json
{
    "message": "New stream key",
    "payload": {
        "key": "_I55odu7WyOBlqkciXADiA",
        "streamID": "30960952-cf23-46f8-bf91-9f738c27be2e"
    }
}
```

## Stream detail

Get information for a specific stream of a user.

- URL: `/streams/{streamID}`
- Method: `GET`
- Authorization required: Yes

**Success response payload:**

- `id`: `string`. Stream's ID.
- `streamKey`: `string`. The stream key that user puts in the URL of the stream.
- `destinations`: `array|null`. Destinations conffigured for the stream. Can be `null` if no destinations are found. When destinations are present, this is the structure of each item:
    - `id`: `number`. Destinantion ID
    - `service`: `string`. The service value matching the request.
    - `streamKey`: `string`. The streaming key for that service matching the request.
    - `rtmpUrl`: `string`. A fully formed RTMP URL of the ingestion server of that service.
    - `isActive`: `boolean`. Initially `true`.
- `provisionStatus`: `object|null`. Provision status of the stream
- `isStreaming`: `boolean`. Whether the stream is currently live.
- `isActive`: `boolean`. Whether the stream is active, implying whether this stream URL can be streamed to.

**Example success response:**

```json
{
    "message": "Stream detail",
    "payload": {
        "id": "30960952-cf23-46f8-bf91-9f738c27be2e",
        "streamKey": "_I55odu7WyOBlqkciXADiA",
        "destinations": [
            {
                "id": 1,
                "service": "YouTube",
                "streamKey": "simaimsimim_siimimimiasdsr",
                "rtmpUrl": "rtmp://a.rtmp.youtube.com/live2/simaimsimim_siimimimiasdsr",
                "isActive": true
            }
        ],
        "provisionStatus": null,
        "isStreaming": false,
        "isActive": true
    }
}
```

## List all streams

List all streams configured by current user.

- URL: `/streams/`
- Method: `GET`
- Authorization required: Yes

**Success response payload:**

Payload is an array of these items:

- `id`: `string`. Stream's ID.
- `streamKey`: `string`. The stream key that user puts in the URL of the stream.
- `destinations`: `null`.
- `provisionStatus`: `null`.
- `isStreaming`: `boolean`. Whether the stream is currently live.
- `isActive`: `boolean`. Whether the stream is active, implying whether this stream URL can be streamed to.

**Example success response:**

```json
{
    "message": "Stream list for user",
    "payload": [
        {
            "id": "30960952-cf23-46f8-bf91-9f738c27be2e",
            "streamKey": "_I55odu7WyOBlqkciXADiA",
            "destinations": null,
            "provisionStatus": null,
            "isStreaming": false,
            "isActive": true
        }
    ]
}
```

## Supported services

These are the values of supported streaming services used in rig's API. These values are case-sensitive.

- `YouTube`
- `Twitch`

# Agents API

This part of the API is reserved for agents. User clients cannot access agent endpoints. Authorization tokens for agent are obtained separately (TODO).

## Stream config

Get stream configuration by stream key.

- URL: `/agents/config/{streamKey}`
- Method: `GET`
- Authorization required: Yes (Not implemented)

**Success response payload:**

Payload is an array of destination configurations in this format:

- `id`: `number`. Destination ID
- `service`: `string`. Service name for this destination
- `streamKey`: `string`. The streaming key configured for this service.
- `rtmpUrl`: `string`. A fully formed RTMP URL of the ingestion server of that service. This is what agent will use to stream to.
- `isActive`: `boolean`. Determines whether this destination is active and should be streamed to.

**Example success response:**

```json
{
    "message": "Stream Config",
    "payload": [
        {
            "id": 1,
            "service": "YouTube",
            "streamKey": "simaimsimim_siimimimiasdsr",
            "rtmpUrl": "rtmp://a.rtmp.youtube.com/live2/simaimsimim_siimimimiasdsr",
            "isActive": true
        },
        {
            "id": 2,
            "service": "Twitch",
            "streamKey": "fxxbxx_live_moofooboo",
            "rtmpUrl": "rtmp://live.twitch.tv/app/fxxbxx_live_moofooboo",
            "isActive": true
        }
    ]
}
```
