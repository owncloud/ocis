Bugfix: Fix index mapping

The index mapping was not being used because we were not using the right blevesearch TypeField, leading to username like properties like `preferred_name` and `on_premises_sam_account_name` to be case sensitive.

https://github.com/owncloud/ocis-accounts/issues/73
