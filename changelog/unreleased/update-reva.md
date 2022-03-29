Enhancement: update reva to v2.1.0

Updated reva to version 2.1.0. This update includes:

  * Fix [cs3org/reva#2636](https://github.com/cs3org/reva/pull/2636): Delay reconnect log for events
  * Fix [cs3org/reva#2645](https://github.com/cs3org/reva/pull/2645): Avoid warning about missing .flock files
  * Fix [cs3org/reva#2625](https://github.com/cs3org/reva/pull/2625): Fix locking on publik links and the decomposed filesystem
  * Fix [cs3org/reva#2643](https://github.com/cs3org/reva/pull/2643): Emit linkaccessfailed event when share is nil
  * Fix [cs3org/reva#2646](https://github.com/cs3org/reva/pull/2646): Replace public mountpoint fileid with grant fileid in ocdav
  * Fix [cs3org/reva#2612](https://github.com/cs3org/reva/pull/2612): Adjust the scope handling to support the spaces architecture
  * Fix [cs3org/reva#2621](https://github.com/cs3org/reva/pull/2621): Send events only if response code is `OK`
  * Chg [cs3org/reva#2574](https://github.com/cs3org/reva/pull/2574): Switch NATS backend
  * Chg [cs3org/reva#2667](https://github.com/cs3org/reva/pull/2667): Allow LDAP groups to have no gidNumber
  * Chg [cs3org/reva#3233](https://github.com/cs3org/reva/pull/3233): Improve quota handling
  * Chg [cs3org/reva#2600](https://github.com/cs3org/reva/pull/2600): Use the cs3 share api to manage spaces
  * Enh [cs3org/reva#2644](https://github.com/cs3org/reva/pull/2644): Add new public share manager
  * Enh [cs3org/reva#2626](https://github.com/cs3org/reva/pull/2626): Add new share manager
  * Enh [cs3org/reva#2624](https://github.com/cs3org/reva/pull/2624): Add etags to virtual spaces
  * Enh [cs3org/reva#2639](https://github.com/cs3org/reva/pull/2639): File Events
  * Enh [cs3org/reva#2627](https://github.com/cs3org/reva/pull/2627): Add events for sharing action
  * Enh [cs3org/reva#2664](https://github.com/cs3org/reva/pull/2664): Add grantID to mountpoint
  * Enh [cs3org/reva#2622](https://github.com/cs3org/reva/pull/2622): Allow listing shares in spaces via the OCS API
  * Enh [cs3org/reva#2623](https://github.com/cs3org/reva/pull/2623): Add space aliases
  * Enh [cs3org/reva#2647](https://github.com/cs3org/reva/pull/2647): Add space specific events
  * Enh [cs3org/reva#3345](https://github.com/cs3org/reva/pull/3345): Add the spaceid to propfind responses
  * Enh [cs3org/reva#2616](https://github.com/cs3org/reva/pull/2616): Add etag to spaces response
  * Enh [cs3org/reva#2628](https://github.com/cs3org/reva/pull/2628): Add spaces aware trash-bin API

https://github.com/owncloud/ocis/pull/3330
https://github.com/owncloud/ocis/pull/3405
https://github.com/owncloud/ocis/pull/3416
