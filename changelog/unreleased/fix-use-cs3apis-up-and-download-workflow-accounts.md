Bugfix: Use the CS3api up- and download workflow for the accounts service

We've fixed the interaction of the accounts service with the metadata storage
after bypassing the InitiateUpload and InitiateDownload have been removed
from various storage drivers. The accounts service now uses the proper
CS3apis workflow for up- and downloads.

https://github.com/owncloud/ocis/pull/2837
https://github.com/cs3org/reva/pull/2309
