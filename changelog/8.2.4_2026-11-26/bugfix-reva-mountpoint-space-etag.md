Bugfix: Expose the etag of shared drives

The `mountpoint` drives of received shares in `/graph/v1.0/me/drives` now carry
a `root.eTag`, so clients can detect changes in shared folders by comparing the
etag of the drive root. The shares jail and mountpoint drives fall back to a
stable default etag when the shared resource has no etag.

https://github.com/owncloud/ocis/pull/13597
https://github.com/owncloud/reva/pull/863
