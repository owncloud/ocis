# storage-publiclink

## Brute Force Protection

The brute force protection will prevent access to public links if wrong passwords are entered. The implementation is very similar to a rate limiter, but taking into account only wrong password attempts.

By default, you're allowed a maximum of 5 failed attempts in 1 hour:

* `STORAGE_PUBLICLINK_BRUTEFORCE_MAXATTEMPTS=5` 
* `STORAGE_PUBLICLINK_BRUTEFORCE_TIMEGAP=1h`

You can adjust those values to your liking in order to define the failure rate threshold (5 failures per hour, by default).

If the failure rate threshold is exceeded, the public link will be blocked until such rate goes below the threshold. This means that it will remain blocked for an undefined time: a couple of seconds in the best case, or up to `STORAGE_PUBLICLINK_BRUTEFORCE_TIME` in the worst case.

If the public link is blocked by the brute force protection, it will be blocked for all the users.

In case of multiple service replicas, the brute force protection won't share any data among the replicas and the failure rate will apply per replica. This means that a replica might be blocked due to high failure rate while the rest work fine.

As said, this feature is enabled by default, with a 5 failures per hour rate. If you want to disable this feature, set the related configuration values to 0.
