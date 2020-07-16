Bugfix: Fix the accountId and groupId mismatch in DeleteGroup Method

We've fixed a bug in deleting the groups.

The accountId and GroupId were swapped when removing the member from a group after deleting
the group.

https://github.com/owncloud/ocis-accounts/pull/60
