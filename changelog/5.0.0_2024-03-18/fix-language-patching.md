Bugfix: Fix patching of language

User would not be able to patch their preferred language when the ldap backend is set to `read-only`.
This makes no sense as language is stored elsewhere.

https://github.com/owncloud/ocis/pull/8182
