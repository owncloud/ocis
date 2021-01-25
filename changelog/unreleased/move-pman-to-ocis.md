Change: Move runtime code on refs/pman over to owncloud/ocis/ocis

Tags: ocis, runtime

Currently, the runtime is under my own private account. For future-proofing we don't want oCIS critical components to depend on external repositories, so we're including refs/pman module as an oCIS package instead.

https://github.com/owncloud/ocis/pull/1483
