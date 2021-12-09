Bugfix: Add the gatewaysvc to all shared configuration in REVA services

We've fixed the configuration for REVA services which didn't have a gatewaysvc in their
shared configuration. This could lead to default gatewaysvc addresses in the auth middleware. Now it is set everywhere.

https://github.com/owncloud/ocis/pull/2597
