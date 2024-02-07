Bugfix: Remove invalid environment variables

We have removed two spaces related environment variables (whether project spaces and the share jail are enabled) and hardcoded the only allowed options. Misusing those variables would have resulted in invalid config.

https://github.com/owncloud/ocis/pull/8303
