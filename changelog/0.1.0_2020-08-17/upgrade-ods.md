Bugfix: Fix multiple submits on string and number form elements

We had a bug with keyboard event listeners triggering multiple submits on input fields.
This was recently fixed in the ownCloud design system (ODS). We rolled out that bugfix
to the settings ui as well.

https://github.com/owncloud/owncloud-design-system/issues/745
https://github.com/owncloud/owncloud-design-system/pull/768
https://github.com/owncloud/ocis-settings/pulls/31
