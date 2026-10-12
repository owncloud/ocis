Bugfix: Password protection is enforced in the "Shared via link" list

Creating a password-protected folder also creates a hidden folder which holds the
actual content and carries the password-protected link. That hidden folder is
listed under Shares / Shared via link, but its entry pointed straight at the
content in the owner's personal space instead of at the link. Anyone looking at
that list could therefore open the folder without being asked for the password.

The entry now opens the password-protected link, so the password is requested
just like it is when opening the folder from the regular file list. If the link
cannot be determined, the entry is no longer clickable at all.

https://github.com/owncloud/ocis/pull/12835
