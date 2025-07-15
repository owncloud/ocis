Bugfix: Add default store to postprocessing

Postprocessing did not have a default store especially `database` and `table` are needed to talk to nats-js

https://github.com/owncloud/ocis/pull/6578
