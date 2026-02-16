Bugfix: Fix search report

There were multiple issues with REPORT search responses from webdav. Also we want it to be consistent with PROPFIND responses.
*   the `remote.php` prefix was missing from the href (added even though not necessary)
*   the ids were formatted wrong, they should look different for shares and spaces.
*   the name of the resource was missing
*   the shareid was missing (for shares)
*   the prop `shareroot` (containing the name of the share root) was missing
*   the permissions prop was empty

https://github.com/owncloud/web/issues/7557
https://github.com/owncloud/ocis/pull/4485
