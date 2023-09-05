# Clientlog service

The `clientlog` service is responsible for composing machine readable notifications for clients

## The `...log` service ecosystem

`...log` services (`userlog`, `clientlog`) are responsible for composing notifications for a certain audience.
  -  `userlog` service translates and adjust messages to be human readable
  -  `clientlog` service composes machine readable messages so clients can act without needing to query the server
  -  `sse` service is only responsible for sending these messages. It does not care about their form or language

## `clientlog` events

The messages the `clientlog` service sends are meant to be used by clients, not by users. The client might for example be informed that a file is finished postprocessing, so it can make the file available to the user without needing to make another call to the server.
