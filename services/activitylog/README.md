# Activitylog

The `activitylog` service is responsible for storing events (activities) per resource.

## The Log Service Ecosystem

Log services like the `activitylog`, `userlog`, `clientlog` and `sse` are responsible for composing notifications for a specific audience.
  -   The `userlog` service translates and adjusts messages to be human readable.
  -   The `clientlog` service composes machine readable messages, so clients can act without the need to query the server.
  -   The `sse` service is only responsible for sending these messages. It does not care about their form or language.
  -   The `activitylog` service stores events per resource. These can be retrieved to show item activities

## Activitylog Store

The `activitylog` stores activities for each resource. It works in conjunction with the `eventhistory` service to keep the data it needs to store to a minimum.
