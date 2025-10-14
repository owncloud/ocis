Bugfix: Fix search by containing special characters

As the OData query parser interprets characters like '@' or '-' in a special
way. Search request for users or groups needs to be quoted. We fixed the libregraph
users and groups endpoints to handle quoted search terms correctly.

https://github.com/owncloud/ocis/pull/8050
https://github.com/owncloud/ocis/pull/8035
https://github.com/owncloud/ocis/issues/7990
