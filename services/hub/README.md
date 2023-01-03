# Hub Service

The hub service provides `server sent events` (`sse`) functionality to clients

## Subscribing

A client can use the `hub/sse` endpoint to subscribe to events. This will open a http(s) connection for the server to send events to subscribed clients. These events can inform clients about various changes on the server like: file uploads, shares, space memberships, etc.

## Available Events
For a complete and up-to-date list of available events see `/services/hub/pkg/service/events.go`.

Note that for the time being, the `hub` service only serves the `UploadReady` event which is emitted when postprocessing a file has finished and the file is available for user access.
