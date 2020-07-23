Change: mint new username property in the reva token

An accounts username is now taken from the on_premises_sam_account_name property instead of the preferred_name.
Furthermore the group name (also from on_premises_sam_account_name property) is now minted into the token as well.

https://github.com/owncloud/ocis-proxy/pull/62
