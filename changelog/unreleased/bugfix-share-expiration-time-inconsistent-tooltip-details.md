Bugfix: Share expiration time now matches between tooltip and share details

A share's expiration time shown in the share details panel no longer
disagrees with the time shown in the hover tooltip. Previously the details
panel always showed 23:59 regardless of when the share was actually set to
expire, which could make a share look like it expires later in the day than
it really does.

https://github.com/owncloud/ocis/pull/13298
