# Eventhistory service

The `eventhistory` consumes all events from the configured event systems, stores them and allows to retrieve them via an eventid

## Consuming

The `eventhistory` services consumes all events from the configured event sytem. Running it without an event sytem is not possible.

## Storing

The `eventhistory` stores each consumed event in the configured store. Possible stores are `inmemory` and ? but not ?.

## Retrieving

Other services can call the `eventhistory` service via a grpc call to retrieve events. The request must contain the eventid that should be retrieved
