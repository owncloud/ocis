Bugfix: Update the admin user role assignment to enforce the config

The admin user role assigment was not updated after the first assignment. We now read the assigned role during init and update the admin user ID accordingly if the role is not assigned.
This is especially needed when the OCIS_ADMIN_USER_ID is set after the autoprovisioning of the admin user when it originates from an external Identity Provider.

https://github.com/owncloud/ocis/pull/8897
