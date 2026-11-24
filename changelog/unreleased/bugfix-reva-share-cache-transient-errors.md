Bugfix: Handle transient storage errors in the received share cache

The share manager now handles transient and internal storage errors in the
retry path of the received share cache. The retry loop uses the new `cas`
package of reva.

The NATS client now reads the credentials file once for each connection
attempt.

https://github.com/owncloud/ocis/pull/13549
https://github.com/owncloud/reva/pull/831
https://github.com/owncloud/reva/pull/829
