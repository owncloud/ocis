# Clientlog Service

The `clientlog` service is responsible for composing machine readable notifications for clients. Clients are apps and web interfaces.

## The Log Service Ecosystem

Log services like the `userlog`, `clientlog` and `sse` are responsible for composing notifications for a certain audience.
  -   The `userlog` service translates and adjusts messages to be human readable.
  -   The `clientlog` service composes machine readable messages, so clients can act without the need to query the server.
  -   The `sse` service is only responsible for sending these messages. It does not care about their form or language.

## Clientlog Events

The messages the `clientlog` service sends are intended for the use by clients, not by users. The client might for example be informed that a file has finished post-processing. With that, the client can make the file available to the user without additional server queries.
