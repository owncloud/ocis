Bugfix: Race condition in config parsing

There was a race condition in the config parsing when configuring the storage services caused by services overwriting a pointer to a config value. We fixed it by setting sane defaults.

https://github.com/owncloud/ocis/pull/2574
