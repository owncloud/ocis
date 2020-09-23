# Changelog for [unreleased] \(UNRELEASED)

The following sections list the changes in ocis-accounts unreleased.

[unreleased]: https://github.com/owncloud/ocis/accounts/compare/v0.4.1...master

## Summary

-   Bugfix - Initialize roleService client in GRPC server: [#114](https://github.com/owncloud/ocis/accounts/pull/114)
-   Bugfix - Cleanup separated indices in memory: [#224](https://github.com/owncloud/product/issues/224)
-   Change - Set user role on builtin users: [#102](https://github.com/owncloud/ocis/accounts/pull/102)
-   Change - Add new builtin admin user: [#102](https://github.com/owncloud/ocis/accounts/pull/102)
-   Change - We make use of the roles cache to enforce permission checks: [#100](https://github.com/owncloud/ocis/accounts/pull/100)
-   Change - We make use of the roles manager to enforce permission checks: [#108](https://github.com/owncloud/ocis/accounts/pull/108)
-   Enhancement - Add create account form: [#148](https://github.com/owncloud/product/issues/148)
-   Enhancement - Add delete accounts action: [#148](https://github.com/owncloud/product/issues/148)
-   Enhancement - Add enable/disable capabilities to the WebUI: [#118](https://github.com/owncloud/product/issues/118)
-   Enhancement - Improve visual appearance of accounts UI: [#222](https://github.com/owncloud/product/issues/222)

## Details

-   Bugfix - Initialize roleService client in GRPC server: [#114](https://github.com/owncloud/ocis/accounts/pull/114)

     We fixed the initialization of the GRPC server by also providing a roleService client and a
     roleManager instance.

     <https://github.com/owncloud/ocis/accounts/pull/114>


-   Bugfix - Cleanup separated indices in memory: [#224](https://github.com/owncloud/product/issues/224)

     The accounts service was creating a bleve index instance in the service handler, thus creating
     separate in memory indices for the http and grpc servers. We moved the service handler creation
     out of the server creation so that the service handler, thus also the bleve index, is a shared
     instance of the servers.

     This fixes a bug that accounts created through the web ui were not able to sign in until a service
     restart.

     <https://github.com/owncloud/product/issues/224>
     <https://github.com/owncloud/ocis/accounts/pull/117>
     <https://github.com/owncloud/ocis/accounts/pull/118>


-   Change - Set user role on builtin users: [#102](https://github.com/owncloud/ocis/accounts/pull/102)

     We now set the default `user` role on our builtin users.

     <https://github.com/owncloud/ocis/accounts/pull/102>


-   Change - Add new builtin admin user: [#102](https://github.com/owncloud/ocis/accounts/pull/102)

     We added a new builtin user `moss` and assigned the admin role.

     <https://github.com/owncloud/ocis/accounts/pull/102>


-   Change - We make use of the roles cache to enforce permission checks: [#100](https://github.com/owncloud/ocis/accounts/pull/100)

     The roles cache and its cache update middleware are used to make permission checks possible.
     The permission checks take place in the accounts handler.

     <https://github.com/owncloud/ocis/accounts/pull/100>


-   Change - We make use of the roles manager to enforce permission checks: [#108](https://github.com/owncloud/ocis/accounts/pull/108)

     The roles cache and its cache update middleware have been replaced with a roles manager in
     ocis-pkg/v2. We've switched over to the new roles manager implementation, to prepare for
     permission checks on grpc requests as well.

     <https://github.com/owncloud/ocis/accounts/pull/108>
     <https://github.com/owncloud/ocis-pkg/pull/60>


-   Enhancement - Add create account form: [#148](https://github.com/owncloud/product/issues/148)

     We've added a form to create new users above the accounts list.

     <https://github.com/owncloud/product/issues/148>
     <https://github.com/owncloud/ocis/accounts/pull/115>


-   Enhancement - Add delete accounts action: [#148](https://github.com/owncloud/product/issues/148)

     We've added an action into the actions dropdown to enable admins to delete users.

     <https://github.com/owncloud/product/issues/148>
     <https://github.com/owncloud/ocis/accounts/pull/115>


-   Enhancement - Add enable/disable capabilities to the WebUI: [#118](https://github.com/owncloud/product/issues/118)

     We've added batch actions into the accounts listing to provide options to enable and disable
     accounts.

     <https://github.com/owncloud/product/issues/118>
     <https://github.com/owncloud/ocis/accounts/pull/109>


-   Enhancement - Improve visual appearance of accounts UI: [#222](https://github.com/owncloud/product/issues/222)

     We aligned the visual appearance of the accounts UI with default ocis-web apps (full width,
     style of batch actions), added icons to buttons, extracted the buttons from the batch actions
     dropdown into individual buttons, improved the wording added a confirmation widget for the
     user deletion and removed the uid and gid columns.

     <https://github.com/owncloud/product/issues/222>
     <https://github.com/owncloud/ocis/accounts/pull/116>

# Changelog for [0.4.1] \(2020-08-27)

The following sections list the changes in ocis-accounts 0.4.1.

[0.4.1]: https://github.com/owncloud/ocis/accounts/compare/v0.3.0...v0.4.1

## Summary

-   Bugfix - Adapting to new settings API for fetching roles: [#96](https://github.com/owncloud/ocis/accounts/pull/96)
-   Change - Create account api-call implicitly adds "default-user" role: [#173](https://github.com/owncloud/product/issues/173)

## Details

-   Bugfix - Adapting to new settings API for fetching roles: [#96](https://github.com/owncloud/ocis/accounts/pull/96)

     We fixed the usage of the ocis-settings endpoint for fetching roles.

     <https://github.com/owncloud/ocis/accounts/pull/96>


-   Change - Create account api-call implicitly adds "default-user" role: [#173](https://github.com/owncloud/product/issues/173)

     When calling CreateAccount default-user-role is implicitly added.

     <https://github.com/owncloud/product/issues/173>

# Changelog for [0.3.0] \(2020-08-20)

The following sections list the changes in ocis-accounts 0.3.0.

[0.3.0]: https://github.com/owncloud/ocis/accounts/compare/v0.4.0...v0.3.0

## Summary

-   Bugfix - Atomic Requests: [#82](https://github.com/owncloud/ocis/accounts/pull/82)
-   Bugfix - Unescape value for prefix query: [#76](https://github.com/owncloud/ocis/accounts/pull/76)
-   Change - Adapt to new ocis-settings data model: [#87](https://github.com/owncloud/ocis/accounts/pull/87)
-   Change - Add permissions for language to default roles: [#88](https://github.com/owncloud/ocis/accounts/pull/88)

## Details

-   Bugfix - Atomic Requests: [#82](https://github.com/owncloud/ocis/accounts/pull/82)

     Operations on the file system level are now atomic. This happens only on the provisioning API.

     <https://github.com/owncloud/ocis/accounts/pull/82>


-   Bugfix - Unescape value for prefix query: [#76](https://github.com/owncloud/ocis/accounts/pull/76)

     Prefix queries also need to unescape token values like `'some ''ol string'` to `some 'ol
     string` before using it in a prefix query

     <https://github.com/owncloud/ocis/accounts/pull/76>


-   Change - Adapt to new ocis-settings data model: [#87](https://github.com/owncloud/ocis/accounts/pull/87)

     Ocis-settings introduced UUIDs and less verbose endpoint and message type names. This PR
     adjusts ocis-accounts accordingly.

     <https://github.com/owncloud/ocis/accounts/pull/87>
     <https://github.com/owncloud/ocis/settings/pull/46>


-   Change - Add permissions for language to default roles: [#88](https://github.com/owncloud/ocis/accounts/pull/88)

     Ocis-settings has default roles and exposes the respective bundle uuids. We now added
     permissions for reading/writing the preferred language to the default roles.

     <https://github.com/owncloud/ocis/accounts/pull/88>

# Changelog for [0.4.0] \(2020-08-20)

The following sections list the changes in ocis-accounts 0.4.0.

[0.4.0]: https://github.com/owncloud/ocis/accounts/compare/v0.2.0...v0.4.0

## Summary

-   Change - Add role selection to accounts UI: [#103](https://github.com/owncloud/product/issues/103)

## Details

-   Change - Add role selection to accounts UI: [#103](https://github.com/owncloud/product/issues/103)

     We added a role selection dropdown for each account in the accounts UI. As a first iteration,
     this doesn't require account management permissions.

     <https://github.com/owncloud/product/issues/103>
     <https://github.com/owncloud/ocis/accounts/pull/89>

# Changelog for [0.2.0] \(2020-08-19)

The following sections list the changes in ocis-accounts 0.2.0.

[0.2.0]: https://github.com/owncloud/ocis/accounts/compare/v0.1.1...v0.2.0

## Summary

-   Bugfix - Add write mutexes: [#71](https://github.com/owncloud/ocis/accounts/pull/71)
-   Bugfix - Fix the accountId and groupId mismatch in DeleteGroup Method: [#60](https://github.com/owncloud/ocis/accounts/pull/60)
-   Bugfix - Fix index mapping: [#73](https://github.com/owncloud/ocis/accounts/issues/73)
-   Bugfix - Use NewNumericRangeInclusiveQuery for numeric literals: [#28](https://github.com/owncloud/ocis-glauth/issues/28)
-   Bugfix - Prevent segfault when no password is set: [#65](https://github.com/owncloud/ocis/accounts/pull/65)
-   Bugfix - Update account return value not used: [#70](https://github.com/owncloud/ocis/accounts/pull/70)
-   Bugfix - Build docker images with alpine:latest instead of alpine:edge: [#64](https://github.com/owncloud/ocis/accounts/pull/64)
-   Change - Align structure of this extension with other extensions: [#51](https://github.com/owncloud/ocis/accounts/pull/51)
-   Change - Change api errors: [#11](https://github.com/owncloud/ocis/accounts/issues/11)
-   Change - Enable accounts on creation: [#43](https://github.com/owncloud/ocis/accounts/issues/43)
-   Change - Fix index update on create/update: [#57](https://github.com/owncloud/ocis/accounts/issues/57)
-   Change - Pass around the correct logger throughout the code: [#41](https://github.com/owncloud/ocis/accounts/issues/41)
-   Change - Remove timezone setting: [#33](https://github.com/owncloud/ocis/accounts/pull/33)
-   Change - Tighten screws on usernames and email addresses: [#65](https://github.com/owncloud/ocis/accounts/pull/65)
-   Enhancement - Add early version of cli tools for user-management: [#69](https://github.com/owncloud/ocis/accounts/pull/69)
-   Enhancement - Update accounts API: [#30](https://github.com/owncloud/ocis/accounts/pull/30)
-   Enhancement - Add simple user listing UI: [#51](https://github.com/owncloud/ocis/accounts/pull/51)

## Details

-   Bugfix - Add write mutexes: [#71](https://github.com/owncloud/ocis/accounts/pull/71)

     Concurrent account or groups writes would corrupt the json file on disk, because the different
     goroutines would be treated as a single thread from the os. We introduce a mutex for account and
     group file writes each. This locks the update frequency for all accounts/groups and could be
     further improved by using a concurrent map of mutexes with a mutex per account / group. PR
     welcome.

     <https://github.com/owncloud/ocis/accounts/pull/71>


-   Bugfix - Fix the accountId and groupId mismatch in DeleteGroup Method: [#60](https://github.com/owncloud/ocis/accounts/pull/60)

     We've fixed a bug in deleting the groups.

     The accountId and GroupId were swapped when removing the member from a group after deleting the
     group.

     <https://github.com/owncloud/ocis/accounts/pull/60>


-   Bugfix - Fix index mapping: [#73](https://github.com/owncloud/ocis/accounts/issues/73)

     The index mapping was not being used because we were not using the right blevesearch TypeField,
     leading to username like properties like `preferred_name` and
     `on_premises_sam_account_name` to be case sensitive.

     <https://github.com/owncloud/ocis/accounts/issues/73>


-   Bugfix - Use NewNumericRangeInclusiveQuery for numeric literals: [#28](https://github.com/owncloud/ocis-glauth/issues/28)

     Some LDAP properties like `uidnumber` and `gidnumber` are numeric. When an OS tries to look up a
     user it will not only try to lookup the user by username, but also by the `uidnumber`:
     `(&(objectclass=posixAccount)(uidnumber=20000))`. The accounts backend for glauth was
     sending that as a string query `uid_number eq '20000'` and has been changed to send it as
     `uid_number eq 20000`. The removed quotes allow the parser in ocis-accounts to identify the
     numeric literal and use the NewNumericRangeInclusiveQuery instead of a TermQuery.

     <https://github.com/owncloud/ocis-glauth/issues/28>
     <https://github.com/owncloud/ocis/accounts/pull/68>
     <https://github.com/owncloud/ocis-glauth/pull/29>


-   Bugfix - Prevent segfault when no password is set: [#65](https://github.com/owncloud/ocis/accounts/pull/65)

     Passwords are stored in a dedicated child struct of an account. We fixed several segfault
     conditions where the methods would try to unset a password when that child struct was not
     existing.

     <https://github.com/owncloud/ocis/accounts/pull/65>


-   Bugfix - Update account return value not used: [#70](https://github.com/owncloud/ocis/accounts/pull/70)

     In order to return a value using the micro go code we need to override the `out` value.

     <https://github.com/owncloud/ocis/accounts/pull/70>


-   Bugfix - Build docker images with alpine:latest instead of alpine:edge: [#64](https://github.com/owncloud/ocis/accounts/pull/64)

     ARM builds were failing when built on alpine:edge, so we switched to alpine:latest instead.

     <https://github.com/owncloud/ocis/accounts/pull/64>


-   Change - Align structure of this extension with other extensions: [#51](https://github.com/owncloud/ocis/accounts/pull/51)

     We aim to have a similar project structure for all our ocis extensions. This extension was
     different with regard to the structure of the server command and naming of some flag names.

     <https://github.com/owncloud/ocis/accounts/pull/51>


-   Change - Change api errors: [#11](https://github.com/owncloud/ocis/accounts/issues/11)

     Replaced the plain golang errors with the error model from the micro framework.

     <https://github.com/owncloud/ocis/accounts/issues/11>


-   Change - Enable accounts on creation: [#43](https://github.com/owncloud/ocis/accounts/issues/43)

     Accounts have been created with the account_enabled flag set to false. Now when they are
     created accounts will be enabled per default.

     <https://github.com/owncloud/ocis/accounts/issues/43>


-   Change - Fix index update on create/update: [#57](https://github.com/owncloud/ocis/accounts/issues/57)

     We fixed a bug in creating/updating accounts and groups, that caused new entities not to show up
     in list queries.

     <https://github.com/owncloud/ocis/accounts/issues/57>
     <https://github.com/owncloud/ocis/accounts/pull/59>


-   Change - Pass around the correct logger throughout the code: [#41](https://github.com/owncloud/ocis/accounts/issues/41)

     Pass around the logger to have consistent log formatting, log level, etc.

     <https://github.com/owncloud/ocis/accounts/issues/41>
     <https://github.com/owncloud/ocis/accounts/pull/48>


-   Change - Remove timezone setting: [#33](https://github.com/owncloud/ocis/accounts/pull/33)

     We had a timezone setting in our profile settings bundle. As we're not dealing with a timezone
     yet it would be confusing for the user to have a timezone setting available. We removed it, until
     we have a timezone implementation available in ocis-web.

     <https://github.com/owncloud/ocis/accounts/pull/33>


-   Change - Tighten screws on usernames and email addresses: [#65](https://github.com/owncloud/ocis/accounts/pull/65)

     In order to match accounts to the OIDC claims we currently rely on the email address or username
     to be present. We force both to match the [W3C recommended
     regex](https://www.w3.org/TR/2016/REC-html51-20161101/sec-forms.html#valid-e-mail-address)
     with usernames having to start with a character or `_`. This allows the username to be presented
     and used in ACLs when integrating the os with the glauth LDAP service of ocis.

     <https://github.com/owncloud/ocis/accounts/pull/65>


-   Enhancement - Add early version of cli tools for user-management: [#69](https://github.com/owncloud/ocis/accounts/pull/69)

     Following commands are available:

     List, ls List existing accounts add, create, Create a new account update Make changes to an
     existing account remove, rm Removes an existing account inspect Show detailed data on an
     existing account

     See --help for details.

     Note that not all account-attributes have an effect yet. This is due to ocis being in an early
     development stage.

     <https://github.com/owncloud/product/issues/115>
     <https://github.com/owncloud/ocis/accounts/pull/69>


-   Enhancement - Update accounts API: [#30](https://github.com/owncloud/ocis/accounts/pull/30)

     We updated the api to allow fetching users not onyl by UUID, but also by identity (OpenID issuer
     and subject) email, username and optionally a password.

     <https://github.com/owncloud/ocis/accounts/pull/30>


-   Enhancement - Add simple user listing UI: [#51](https://github.com/owncloud/ocis/accounts/pull/51)

     We added an extension for ocis-web that shows a simple list of all existing users.

     <https://github.com/owncloud/ocis/accounts/pull/51>

# Changelog for [0.1.1] \(2020-04-29)

The following sections list the changes in ocis-accounts 0.1.1.

[0.1.1]: https://github.com/owncloud/ocis/accounts/compare/v0.1.0...v0.1.1

## Summary

-   Enhancement - Logging is configurable: [#24](https://github.com/owncloud/ocis/accounts/pull/24)

## Details

-   Enhancement - Logging is configurable: [#24](https://github.com/owncloud/ocis/accounts/pull/24)

     ACCOUNTS_LOG_\* env-vars or cli-flags can be used for logging configuration. See --help for
     more details.

     <https://github.com/owncloud/ocis/accounts/pull/24>

# Changelog for [0.1.0] \(2020-03-18)

The following sections list the changes in ocis-accounts 0.1.0.

[0.1.0]: https://github.com/owncloud/ocis/accounts/compare/500e303cb544ed93d84153f01219d77eeee44929...v0.1.0

## Summary

-   Change - Initial release of basic version: [#1](https://github.com/owncloud/ocis/accounts/issues/1)
-   Enhancement - Configuration: [#15](https://github.com/owncloud/ocis/accounts/pull/15)

## Details

-   Change - Initial release of basic version: [#1](https://github.com/owncloud/ocis/accounts/issues/1)

     Just prepared an initial basic version.

     <https://github.com/owncloud/ocis/accounts/issues/1>


-   Enhancement - Configuration: [#15](https://github.com/owncloud/ocis/accounts/pull/15)

     Extensions should be responsible of configuring themselves. We use Viper for config loading
     from default paths. Environment variables **WILL** take precedence over config files.

     <https://github.com/owncloud/ocis/accounts/pull/15>
