Bugfix: Rework default role provisioning

We fixed a race condition in the default role assignment code that could lead to
users loosing privileges. When authenticating before the settings service was fully
running.

https://github.com/owncloud/ocis/issues/3900
