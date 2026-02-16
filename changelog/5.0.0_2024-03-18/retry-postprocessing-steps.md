Enhancement: Retry antivirus postprocessing step in case of problems

The antivirus postprocessing step will now be retried for a configurable amount of times in case it can't get a result from clamav.

https://github.com/owncloud/ocis/pull/7874
