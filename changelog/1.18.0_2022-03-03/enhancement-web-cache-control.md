Enhancement: Re-Enabling web cache control

We've re-enable browser caching headers (`Expires` and `Last-Modified`) for the web service, this was disabled due to a problem in the fileserver used before.
Since we're now using our own fileserver implementation this works again and is enabled by default.

https://github.com/owncloud/ocis/pull/3109
