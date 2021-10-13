Change: Configure users and metadata storage separately

We've fixed the configuration behaviour of the user and metadata service writing in the same
directory when using oCIS storage.

Therefore we needed to separate the configuration of the users and metadata storage so that they
now can be configured totally separate.

https://github.com/owncloud/ocis/pull/2598
