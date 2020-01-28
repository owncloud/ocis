# Changelog for unreleased

The following sections list the changes for unreleased.

## Summary

 * Chg #1: Initial release of basic version
 * Chg #6: Start multiple services with dedicated commands

## Details

 * Change #1: Initial release of basic version

   Just prepared an initial basic version to start a reva server and start integrating with the
   go-micro base dextension framework of ownCloud Infinite Scale.

   https://github.com/owncloud/ocis-reva/issues/1

 * Change #6: Start multiple services with dedicated commands

   The initial version would only allow us to use a set of reva configurations to start multiple
   services. We use a more opinionated set of commands to start dedicated services that allows us
   to configure them individually. It allowcs us to switch eg. the user backend to LDAP and fully it
   on the cli.

   https://github.com/owncloud/ocis-reva/issues/6


