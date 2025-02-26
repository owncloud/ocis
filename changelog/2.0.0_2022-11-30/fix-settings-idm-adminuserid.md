Bugfix: Fix the idm and settings extensions' admin user id configuration option

We've fixed the admin user id configuration of the settings and idm extensions.
The have previously only been configurable via the oCIS shared configuration and
therefore have been undocumented for the extensions. This config option is now part
of both extensions' configuration and can now also be used when the extensions are
compiled standalone.

https://github.com/owncloud/ocis/pull/3799
