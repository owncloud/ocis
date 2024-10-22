Bugfix: increase event processing workers

We increased the number of go routines that pull events from the queue to three and made the number off workers configurable. Furthermore, the postprocessing delay no longer introduces a sleep that slows down pulling of events, but asynchronously triggers the next step.

https://github.com/owncloud/ocis/pull/10385
https://github.com/owncloud/ocis/pull/10368
