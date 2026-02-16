Enhancement: Update antivirus service

We update the antivirus icap client library and optimize the antivirus scanning service.
ANTIVIRUS_ICAP_TIMEOUT is now deprecated and ANTIVIRUS_ICAP_SCAN_TIMEOUT should be used instead.

ANTIVIRUS_ICAP_SCAN_TIMEOUT supports human durations like `1s`, `1m`, `1h` and `1d`.

https://github.com/owncloud/ocis/pull/8062
https://github.com/owncloud/ocis/issues/6764
