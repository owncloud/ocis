Enhancement: Add `ocis shares clean-corrupt-public-shares` maintenance command

A single public-share entry with a nil/empty `resource_id` makes the json
public-share manager's `ListPublicShares` panic with a nil-pointer dereference.
Because the manager reads all entries and filters them in memory, that one bad
entry poisons the endpoint for the whole tenant: every Members/permissions panel
load and every password-protected link creation fails.

The new `ocis shares clean-corrupt-public-shares` command detects and removes
such corrupt entries. It reads the raw persistence (so it never triggers the
panic itself) and writes back through the same metadata storage path the manager
uses, recomputing blob size, mtime and etag automatically. It defaults to
`--dry-run` and supports the `jsoncs3` and `json` public-share drivers.

https://github.com/owncloud/ocis/pull/12494
