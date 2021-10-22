Enhancement: Update reva to v1.14.1-0.20211022092730-84a5191f1c5a

Updated reva to v1.14.1-0.20211022092730-84a5191f1c5a
This update includes:
 * Fix [cs3org/reva#2168](https://github.com/cs3org/reva/pull/2168): Override provider if was previously registered
 * Fix [cs3org/reva#2173](https://github.com/cs3org/reva/pull/2173): Fix archiver max size reached error
 * Fix [cs3org/reva#2167](https://github.com/cs3org/reva/pull/2167): Handle nil quota in decomposedfs
 * Fix [cs3org/reva#2153](https://github.com/cs3org/reva/pull/2153): Restrict EOS project spaces sharing permissions to admins and writers
 * Fix [cs3org/reva#2179](https://github.com/cs3org/reva/pull/2179): Fix the returned permissions for webdav uploads
 * Chg [cs3org/reva#2479](https://github.com/cs3org/reva/pull/2479): Make apps able to work with public shares
 * Enh [cs3org/reva#2174](https://github.com/cs3org/reva/pull/2174): Inherit ACLs for files from parent directories
 * Enh [cs3org/reva#2152](https://github.com/cs3org/reva/pull/2152): Add a reference parameter to the getQuota request
 * Enh [cs3org/reva#2171](https://github.com/cs3org/reva/pull/2171): Add optional claim parameter to machine auth
 * Enh [cs3org/reva#2135](https://github.com/cs3org/reva/pull/2135): Nextcloud test improvements
 * Enh [cs3org/reva#2180](https://github.com/cs3org/reva/pull/2180): Remove OCDAV options namespace parameter
 * Enh [cs3org/reva#2170](https://github.com/cs3org/reva/pull/2170): Handle propfind requests for existing files
 * Enh [cs3org/reva#2165](https://github.com/cs3org/reva/pull/2165): Allow access to recycle bin for arbitrary paths outside homes
 * Enh [cs3org/reva#2189](https://github.com/cs3org/reva/pull/2189): Add user settings capability
 * Enh [cs3org/reva#2162](https://github.com/cs3org/reva/pull/2162): Implement the UpdateStorageSpace method

https://github.com/owncloud/ocis/pull/2658
https://github.com/owncloud/ocis/pull/2536
https://github.com/owncloud/ocis/pull/2650
