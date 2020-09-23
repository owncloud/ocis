Bugfix: redirect to the provided uri

The phoenix client was not set as trusted therefore when logging out the user was redirected to a default page instead of the provided url.

https://github.com/owncloud/ocis/konnectd/issues/26
