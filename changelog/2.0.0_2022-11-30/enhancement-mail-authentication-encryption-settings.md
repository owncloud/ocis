Enhancement: Add configuration options for mail authentication and encryption

We've added configuration options to configure the authentication and encryption
for sending mails in the notifications service.

Furthermore there is now a distinguished configuration option for the username to use
for authentication against the mail server. This allows you to customize the sender address
to your liking. For example sender addresses like `my oCIS instance <ocis@owncloud.test>` are now possible, too.

https://github.com/owncloud/ocis/pull/4443
