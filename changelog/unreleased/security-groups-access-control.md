Security: Restrict group member data to authorized users

The groups API could return group membership details to callers who should not have had access to it. Access to group member data is now properly restricted based on the caller's permissions.

https://github.com/owncloud/ocis/pull/12573
