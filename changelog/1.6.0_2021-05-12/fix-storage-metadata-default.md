Bugfix: Fix STORAGE_METADATA_ROOT default value override

The way the value was being set ensured that it was NOT being overridden where it should have been. This patch ensures the correct loading order of values.

https://github.com/owncloud/ocis/pull/1956
