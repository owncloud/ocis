Enhancement: Disable public link sharing for vault resources

The `graph` service now rejects creating, updating, and setting passwords on
public links when the target resource lives in the vault storage provider.
Requests targeting a vault resource return `400 Bad Request` with the message
`public links are not allowed for vault resources`.

https://github.com/owncloud/ocis/pull/12321
