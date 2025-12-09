Enhancement: Implement brute force protection for public links

Public links will be protected by default, allowing up to 5 wrong password
attempts per hour. If such rate is exceeded, the link will be blocked for
all the users until the failure rate goes below the configured threshold
(5 failures per hour by default, as said).

The failure rate is configurable, so it can be 10 failures each 2 hours
or 3 failures per minute.

Note that the protection will apply per service replica, so one replica
might be blocked while another replica is fully functional.

https://github.com/owncloud/ocis/pull/11864
https://github.com/owncloud/reva/pull/460
