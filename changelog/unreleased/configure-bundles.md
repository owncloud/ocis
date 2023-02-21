Enhancement: Make the settings bundles part of the service config

We added the settings bundles to the config. The default roles are still unchanged. You can now override the defaults by replacing the whole bundles list via json config files. The config file is loaded from a specified path which can be configured with `SETTINGS_BUNDLES_PATH`.

https://github.com/owncloud/ocis/pull/5589
https://github.com/owncloud/ocis/pull/5607
