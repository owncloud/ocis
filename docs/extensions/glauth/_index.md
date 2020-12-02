---
title: "GLAuth"
date: 2018-05-02T00:00:00+00:00
weight: 20
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/extensions/glauth
geekdocFilePath: _index.md
---

This service provides a [glauth](https://github.com/glauth/glauth) based LDAP proxy for ocis which can be used by clients or other extensions. It allows applications relying on LDAP to access the accounts stored in the ocis accounts service. It can be used to make firewalls or identity providers aware of all users, including guest accounts.

We are using it to make eos aware of all accounts so the native ACLs can be used to persist share information in the storage.
