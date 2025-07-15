Enhancement: Initiator-IDs

Allows sending a header `Initiator-ID` on http requests. This id will be added to sse events so clients can figure out if their particular instance was triggering the event. Additionally this adds the etag of the file/folder to all sse events.

https://github.com/owncloud/ocis/pull/8936
https://github.com/owncloud/ocis/pull/8701
