Change: Remove username field in OCS

Tags: ocs

We use the incoming userid as both the `id` and the `on_premises_sam_account_name` for new accounts in the accounts
service. The userid in OCS requests is in fact the username, not our internal account id. We need to enforce the userid
as our internal account id though, because the account id is part of various `path` formats.

https://github.com/owncloud/ocis/pull/709
https://github.com/owncloud/ocis/pull/816
