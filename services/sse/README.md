# SSE

The `sse` service is responsible for sending sse (Server-Sent Events) to a user. See [What is Server-Sent Events](https://medium.com/yemeksepeti-teknoloji/what-is-server-sent-events-sse-and-how-to-implement-it-904938bffd73) for a simple introduction and examples of server sent events.

## The Log Service Ecosystem

Log services like the `userlog`, `clientlog` and `sse` are responsible for composing notifications for a certain audience.
  -   The `userlog` service translates and adjusts messages to be human readable.
  -   The `clientlog` service composes machine readable messages, so clients can act without the need to query the server.
  -   The `sse` service is only responsible for sending these messages. It does not care about their form or language.

## Subscribing

Clients can subscribe to the `/sse` endpoint to be informed by the server when an event happens. The `sse` endpoint will respect language changes of the user without needing to reconnect. Note that SSE has a limitation of six open connections per browser which can be reached if one has opened various tabs of the Web UI pointing to the same Infinite Scale instance.

## Keep SSE Connections Alive

Some intermediate proxies drop connections after an idle time with no activity. If this is the case, configure the `SSE_KEEPALIVE_INTERVAL` envvar. This will send periodic SSE comments to keep connections open.

