---
title: Accounts
date: 2018-05-02T00:00:00+00:00
weight: 20
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/extensions/accounts
geekdocFilePath: _index.md
geekdocCollapseSection: true
---

## Abstract
oCIS needs to be able to identify users. Without a non reassignable and persistent account ID share metadata cannot be reliably persisted. `accounts` allows exchanging oidc claims for a uuid. Using a uuid allows users to change the login, mail or even openid connect provider without breaking any persisted metadata that might have been attached to it.

- persists accounts
- uses graph api properties
- ldap can be synced using the onpremise* attributes

## Table of Contents

{{< toc-tree >}}


## Adding users from the commad line.

The fastest way to create a user is by executing the following command:


    ./ocis accounts add --username johndoe --mail jondoe@none.com --password p123p123 --description  "John Doe's account" --displayname "John Doe" --enabled 

For detailed information just type:


    ./ocis accounts add --help
    NAME:
       ocis accounts add - Create a new account
    
    USAGE:
       ocis accounts add [command options] [arguments...]
    
    OPTIONS:
       --grpc-namespace value                Set the base namespace for the grpc namespace (default: "com.owncloud.api") [$ACCOUNTS_GRPC_NAMESPACE]
       --name value                          service name (default: "accounts") [$ACCOUNTS_NAME]
       --enabled                             Enable the account (default: false)
       --displayname value                   Set the displayname for the account
       --username value                      Username will be written to preferred-name and on_premises_sam_account_name
       --preferred-name value                Set the preferred-name for the account
       --on-premises-sam-account-name value  Set the on-premises-sam-account-name
       --uidnumber value                     Set the uidnumber for the account (default: 0)
       --gidnumber value                     Set the gidnumber for the account (default: 0)
       --mail value                          Set the mail for the account
       --description value                   Set the description for the account
       --password value                      Set the password for the account
       --password-policies value             Possible policies: DisableStrongPassword, DisablePasswordExpiration
       --force-password-change               Force password change on next sign-in (default: false)
       --force-password-change-mfa           Force password change on next sign-in with mfa (default: false)
       --help                                Show the help (default: false)


