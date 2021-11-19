Enhancement: Update REVA to v1.16.0

Updated REVA to v1.16.0
This update includes:
 * Fix [cs3org/reva#2245](https://github.com/cs3org/reva/pull/2245): Don't announce search-files capability
 * Fix [cs3org/reva#2247](https://github.com/cs3org/reva/pull/2247): Merge user ACLs from EOS to sys ACLs
 * Fix [cs3org/reva#2279](https://github.com/cs3org/reva/pull/2279): Return the inode of the version folder for files when listing in EOS
 * Fix [cs3org/reva#2294](https://github.com/cs3org/reva/pull/2294): Fix HTTP return code when path is invalid
 * Fix [cs3org/reva#2231](https://github.com/cs3org/reva/pull/2231): Fix share permission on a single file in sql share driver (cbox pkg)
 * Fix [cs3org/reva#2230](https://github.com/cs3org/reva/pull/2230): Fix open by default app and expose default app
 * Fix [cs3org/reva#2265](https://github.com/cs3org/reva/pull/2265): Fix nil pointer exception when resolving members of a group (rest driver)
 * Fix [cs3org/reva#1214](https://github.com/cs3org/reva/pull/1214): Fix restoring versions
 * Fix [cs3org/reva#2254](https://github.com/cs3org/reva/pull/2254): Fix spaces propfind
 * Fix [cs3org/reva#2260](https://github.com/cs3org/reva/pull/2260): Fix unset quota xattr on darwin
 * Fix [cs3org/reva#5776](https://github.com/cs3org/reva/pull/5776): Enforce permissions in public share apps
 * Fix [cs3org/reva#2767](https://github.com/cs3org/reva/pull/2767): Fix status code for WebDAV mkcol requests where an ancestor is missing
 * Fix [cs3org/reva#2287](https://github.com/cs3org/reva/pull/2287): Add public link access via mount-ID:token/relative-path to the scope
 * Fix [cs3org/reva#2244](https://github.com/cs3org/reva/pull/2244): Fix the permissions response for shared files in the cbox sql driver
 * Enh [cs3org/reva#2219](https://github.com/cs3org/reva/pull/2219): Add virtual view tests
 * Enh [cs3org/reva#2230](https://github.com/cs3org/reva/pull/2230): Add priority to app providers
 * Enh [cs3org/reva#2258](https://github.com/cs3org/reva/pull/2258): Improved error messages from the AppProviders
 * Enh [cs3org/reva#2119](https://github.com/cs3org/reva/pull/2119): Add authprovider owncloudsql
 * Enh [cs3org/reva#2211](https://github.com/cs3org/reva/pull/2211): Enhance the cbox share sql driver to store accepted group shares
 * Enh [cs3org/reva#2212](https://github.com/cs3org/reva/pull/2212): Filter root path according to the agent that makes the request
 * Enh [cs3org/reva#2237](https://github.com/cs3org/reva/pull/2237): Skip get user call in eosfs in case previous ones also failed
 * Enh [cs3org/reva#2266](https://github.com/cs3org/reva/pull/2266): Callback for the EOS UID cache to retry fetch for failed keys
 * Enh [cs3org/reva#2215](https://github.com/cs3org/reva/pull/2215): Aggregrate resource info properties for virtual views
 * Enh [cs3org/reva#2271](https://github.com/cs3org/reva/pull/2271): Revamp the favorite manager and add the cbox sql driver
 * Enh [cs3org/reva#2248](https://github.com/cs3org/reva/pull/2248): Cache whether a user home was created or not
 * Enh [cs3org/reva#2282](https://github.com/cs3org/reva/pull/2282): Return a proper NOT_FOUND error when a user or group is not found
 * Enh [cs3org/reva#2268](https://github.com/cs3org/reva/pull/2268): Add the reverseproxy http service
 * Enh [cs3org/reva#2207](https://github.com/cs3org/reva/pull/2207): Enable users to list all spaces
 * Enh [cs3org/reva#2286](https://github.com/cs3org/reva/pull/2286): Add trace ID to middleware loggers
 * Enh [cs3org/reva#2251](https://github.com/cs3org/reva/pull/2251): Mentix service inference
 * Enh [cs3org/reva#2218](https://github.com/cs3org/reva/pull/2218): Allow filtering of mime types supported by app providers
 * Enh [cs3org/reva#2213](https://github.com/cs3org/reva/pull/2213): Add public link share type to propfind response
 * Enh [cs3org/reva#2253](https://github.com/cs3org/reva/pull/2253): Support the file editor role for public links
 * Enh [cs3org/reva#2208](https://github.com/cs3org/reva/pull/2208): Reduce redundant stat calls when statting by resource ID
 * Enh [cs3org/reva#2235](https://github.com/cs3org/reva/pull/2235): Specify a list of allowed folders/files to be archived
 * Enh [cs3org/reva#2267](https://github.com/cs3org/reva/pull/2267): Restrict the paths where share creation is allowed
 * Enh [cs3org/reva#2252](https://github.com/cs3org/reva/pull/2252): Add the xattr sys.acl to SysACL (eosgrpc)
 * Enh [cs3org/reva#2239](https://github.com/cs3org/reva/pull/2239): Update toml configs


https://github.com/owncloud/ocis/pull/2737
https://github.com/owncloud/ocis/pull/2726
https://github.com/owncloud/ocis/pull/2790
https://github.com/owncloud/ocis/pull/2797
