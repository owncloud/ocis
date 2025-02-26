Enhancement: update reva to v2.2.0

Updated reva to version 2.2.0. This update includes:

* Bugfix [cs3org/reva#3373](https://github.com/cs3org/reva/pull/3373):  Fix the permissions attribute in propfind responses
* Bugfix [cs3org/reva#2721](https://github.com/cs3org/reva/pull/2721):  Fix locking and public link scope checker to make the WOPI server work
* Bugfix [cs3org/reva#2668](https://github.com/cs3org/reva/pull/2668):  Minor cleanup
* Bugfix [cs3org/reva#2692](https://github.com/cs3org/reva/pull/2692):  Ensure that the host in the ocs config endpoint has no protocol
* Bugfix [cs3org/reva#2709](https://github.com/cs3org/reva/pull/2709):  Decomposed FS: return precondition failed if already locked
* Change [cs3org/reva#2687](https://github.com/cs3org/reva/pull/2687):  Allow link with no or edit permission
* Change [cs3org/reva#2658](https://github.com/cs3org/reva/pull/2658):  Small clean up of the ocdav code
* Change [cs3org/reva#2691](https://github.com/cs3org/reva/pull/2691):  Decomposed FS: return a reference to the parent
* Enhancement [cs3org/reva#2708](https://github.com/cs3org/reva/pull/2708):  Rework LDAP configuration of user and group providers
* Enhancement [cs3org/reva#2665](https://github.com/cs3org/reva/pull/2665):  Add embeddable ocdav go micro service
* Enhancement [cs3org/reva#2715](https://github.com/cs3org/reva/pull/2715):  Introduced quicklinks
* Enhancement [cs3org/reva#3370](https://github.com/cs3org/reva/pull/3370):  Enable all spaces members to list public shares
* Enhancement [cs3org/reva#3370](https://github.com/cs3org/reva/pull/3370):  Enable space members to list shares inside the space
* Enhancement [cs3org/reva#2717](https://github.com/cs3org/reva/pull/2717):  Add definitions for user and group events

https://github.com/owncloud/ocis/pull/3397
https://github.com/owncloud/ocis/pull/3430
https://github.com/owncloud/ocis/pull/3476
https://github.com/owncloud/ocis/pull/3482
https://github.com/owncloud/ocis/pull/3497
https://github.com/owncloud/ocis/pull/3513
https://github.com/owncloud/ocis/pull/3514
