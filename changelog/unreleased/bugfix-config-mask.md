Bugfix: Mask configs that hold secrets

Envvars and their config structs can hold secrets.
The ServiceAccount config is now masked.

https://github.com/owncloud/ocis/pull/12397

