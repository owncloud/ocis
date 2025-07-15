Bugfix: Fix opening images in media viewer for some usernames

We've fixed the opening of images in the media viewer for user names containing special characters (eg. `@`) which will be URL-escaped. Before this fix users could not see the image in the media viewer. Now the user name is correctly escaped and the user can view the image in the media viewer.

https://github.com/owncloud/ocis/pull/2738
