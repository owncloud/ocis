Enhancement: SSE for messaging

So far, sse has only been used to exchange messages between the server and the client.
In order to be able to send more content to the client, we have moved the endpoint to a separate service and are now also using it for other notifications like:

* notify postprocessing state changes.
* notify file locking and unlocking.

https://github.com/owncloud/ocis/pull/6992
