Bugfix: Fix Activitylog issues

Fixes multiple activititylog issues. There was an error about `max payload exceeded` when there were too many activities on one folder. Listing would take very long even with a limit activated. All of these
issues are now fixed.

https://github.com/owncloud/ocis/pull/10376
