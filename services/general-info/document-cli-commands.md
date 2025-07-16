---
title: Document CLI Commands
date: 2025-01-09T00:00:00+00:00
weight: 20
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/services/general-info
geekdocFilePath: document-cli-commands.md
geekdocCollapseSection: true
---

Any CLI command that is added to Infinite Scale must be documented in the dev docs and the [admin docs](https://doc.owncloud.com/ocis/latest/maintenance/commands/commands.html). Note that the admin docs mainly differentiate between online and offline commands as the docs structure is different. Any command documented in the dev docs is properly integrated into the admin docs. The following description is for dev docs, admin docs derive from it.

Note that ANY CLI command needs documentation, but it can be decided that a CLI command will not be added to the admin docs (the reasons should be really valid for such a case).

## Type of CLI Commands

There are three types of CLI commands that require different documentation locations:

1. Commands that depend on a service dependent like\
`ocis storage-users uploads`
2. Commands that are service independent like\
`ocis trash purge-empty-dirs` or `ocis revisions purge`
3. `curl` commands that can be one of the above.


## Rules

* Add any service dependent command into the repsective `README.md` _of the service_.
* Add any service independent command into `ocis/README.md`

## Tips

For examples, see either `ocis/README.md` or\
one of the respective service readme's like in\
`services/storage-users/README.md` or `services/auth-app/README.md`. 
