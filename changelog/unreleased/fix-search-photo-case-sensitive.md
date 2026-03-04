Bugfix: Fix case-sensitive photo metadata search

Searching for photo metadata fields like camera make/model was case-sensitive, so searching for "google" would not match a camera make stored as "Google". Changed the photo string field analyzer from `keyword` to `lowercaseKeyword` so both indexed values and search terms are lowercased. Existing Bleve indexes need to be rebuilt after this change.

https://github.com/owncloud/ocis/pull/12078
