Bugfix: Fix concurrent map access when listing shares

Listing shares could abort the whole oCIS process with the Go runtime error
"fatal error: concurrent map read and map write". This is an unrecoverable
runtime fault, so it could not be caught and recovered from, and it dropped all
connections that were in flight at that moment.

When converting CS3 shares into Graph DriveItems, the worker goroutines read
from the shared driveItems map while the collecting loop was already writing
results into the very same map. The results channel is buffered, so the workers
never blocked when handing over an item and the collecting loop started writing
while workers were still running. The results are now collected after all
workers have finished, which removes the overlapping access.

Note that lowering the maximum concurrency did not avoid this, as a single
worker could still run concurrently with the collecting loop.

https://github.com/owncloud/ocis/pull/12670
