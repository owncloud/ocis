Enhancement: Optional events in graph service

We've changed the graph service so that you also can start it without any
event bus.
Therefore you need to set `GRAPH_EVENTS_ENDPOINT` to an empty string.
The graph API will not emit any events in this case.

https://github.com/owncloud/ocis/pull/55555
