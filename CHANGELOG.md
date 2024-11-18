# Table of Contents

* [Changelog for unreleased](#changelog-for-unreleased-unreleased)
* [Changelog for 5.0.9](#changelog-for-509-2024-11-14)
* [Changelog for 7.0.0-rc.2](#changelog-for-700-rc2-2024-11-12)
* [Changelog for 7.0.0-rc.1](#changelog-for-700-rc1-2024-11-07)
* [Changelog for 6.6.1](#changelog-for-661-2024-10-24)
* [Changelog for 6.6.0](#changelog-for-660-2024-10-21)
* [Changelog for 6.5.0](#changelog-for-650-2024-10-01)
* [Changelog for 5.0.8](#changelog-for-508-2024-09-30)
* [Changelog for 6.4.0](#changelog-for-640-2024-09-12)
* [Changelog for 5.0.7](#changelog-for-507-2024-09-04)
* [Changelog for 6.3.0](#changelog-for-630-2024-08-20)
* [Changelog for 6.2.0](#changelog-for-620-2024-07-30)
* [Changelog for 5.0.6](#changelog-for-506-2024-07-17)
* [Changelog for 6.1.0](#changelog-for-610-2024-07-08)
* [Changelog for 6.0.0](#changelog-for-600-2024-06-19)
* [Changelog for 5.0.5](#changelog-for-505-2024-05-22)
* [Changelog for 5.0.4](#changelog-for-504-2024-05-13)
* [Changelog for 5.0.3](#changelog-for-503-2024-05-02)
* [Changelog for 5.0.2](#changelog-for-502-2024-04-17)
* [Changelog for 5.0.1](#changelog-for-501-2024-04-10)
* [Changelog for 4.0.7](#changelog-for-407-2024-03-27)
* [Changelog for 5.0.0](#changelog-for-500-2024-03-18)
* [Changelog for 4.0.6](#changelog-for-406-2024-02-07)
* [Changelog for 4.0.5](#changelog-for-405-2023-12-21)
* [Changelog for 4.0.4](#changelog-for-404-2023-12-07)
* [Changelog for 4.0.3](#changelog-for-403-2023-11-24)
* [Changelog for 4.0.2](#changelog-for-402-2023-09-28)
* [Changelog for 4.0.1](#changelog-for-401-2023-09-01)
* [Changelog for 4.0.0](#changelog-for-400-2023-08-21)
* [Changelog for 3.0.0](#changelog-for-300-2023-06-06)
* [Changelog for 2.0.0](#changelog-for-200-2022-11-30)
* [Changelog for 1.20.0](#changelog-for-1200-2022-04-13)
* [Changelog for 1.19.0](#changelog-for-1190-2022-03-29)
* [Changelog for 1.19.1](#changelog-for-1191-2022-03-29)
* [Changelog for 1.18.0](#changelog-for-1180-2022-03-03)
* [Changelog for 1.17.0](#changelog-for-1170-2022-02-16)
* [Changelog for 1.16.0](#changelog-for-1160-2021-12-10)
* [Changelog for 1.15.0](#changelog-for-1150-2021-11-19)
* [Changelog for 1.14.0](#changelog-for-1140-2021-10-27)
* [Changelog for 1.13.0](#changelog-for-1130-2021-10-13)
* [Changelog for 1.12.0](#changelog-for-1120-2021-09-14)
* [Changelog for 1.11.0](#changelog-for-1110-2021-08-24)
* [Changelog for 1.10.0](#changelog-for-1100-2021-08-06)
* [Changelog for 1.9.0](#changelog-for-190-2021-07-13)
* [Changelog for 1.8.0](#changelog-for-180-2021-06-28)
* [Changelog for 1.7.0](#changelog-for-170-2021-06-04)
* [Changelog for 1.6.0](#changelog-for-160-2021-05-12)
* [Changelog for 1.5.0](#changelog-for-150-2021-04-21)
* [Changelog for 1.4.0](#changelog-for-140-2021-03-30)
* [Changelog for 1.3.0](#changelog-for-130-2021-03-09)
* [Changelog for 1.2.0](#changelog-for-120-2021-02-17)
* [Changelog for 1.1.0](#changelog-for-110-2021-01-22)
* [Changelog for 1.0.0](#changelog-for-100-2020-12-17)

# Changelog for [unreleased] (UNRELEASED)

The following sections list the changes for unreleased.

[unreleased]: https://github.com/owncloud/ocis/compare/v5.0.9...master

## Summary

* Bugfix - Log GRPC requests in debug mode: [#10438](https://github.com/owncloud/ocis/pull/10438)
* Bugfix - Fix impersonated request user mismatch: [#10548](https://github.com/owncloud/ocis/pull/10548)
* Bugfix - Fix federated sharing when using an external IDP: [#10567](https://github.com/owncloud/ocis/pull/10567)
* Bugfix - Fix node cache ttl updates: [#10575](https://github.com/owncloud/ocis/pull/10575)
* Bugfix - We now limit the number of workers of the jsoncs3 share manager: [#10578](https://github.com/owncloud/ocis/pull/10578)
* Bugfix - Set MaxConcurrency to 1: [#10580](https://github.com/owncloud/ocis/pull/10580)
* Bugfix - Reuse go-micro service clients: [#10582](https://github.com/owncloud/ocis/pull/10582)
* Bugfix - Make collaboration service use a gateway selector: [#10584](https://github.com/owncloud/ocis/pull/10584)
* Bugfix - Return an error if we can't get the keys and ensure they're cached: [#10590](https://github.com/owncloud/ocis/pull/10590)
* Bugfix - Fix status code for thumbnail requests: [#10592](https://github.com/owncloud/ocis/pull/10592)
* Bugfix - Fix the activity field mapping: [#10593](https://github.com/owncloud/ocis/pull/10593)
* Enhancement - Update web to v11.0.3: [#10569](https://github.com/owncloud/ocis/pull/10569)

## Details

* Bugfix - Log GRPC requests in debug mode: [#10438](https://github.com/owncloud/ocis/pull/10438)

   When log level is set to debug we will now also log grpc requests.

   https://github.com/owncloud/ocis/pull/10438

* Bugfix - Fix impersonated request user mismatch: [#10548](https://github.com/owncloud/ocis/pull/10548)

   We fixed a user id and name mismatch in the impersonated auth-app API request

   https://github.com/owncloud/ocis/issues/10292
   https://github.com/owncloud/ocis/pull/10548

* Bugfix - Fix federated sharing when using an external IDP: [#10567](https://github.com/owncloud/ocis/pull/10567)

   We fixed a bug that caused federated sharing to fail, when the federated oCIS
   instances where sharing the same external IDP.

   https://github.com/owncloud/ocis/pull/10567
   https://github.com/cs3org/reva/pull/4933

* Bugfix - Fix node cache ttl updates: [#10575](https://github.com/owncloud/ocis/pull/10575)

   We now only udpate the TTL of the node that is created or updated.

   https://github.com/owncloud/ocis/pull/10575

* Bugfix - We now limit the number of workers of the jsoncs3 share manager: [#10578](https://github.com/owncloud/ocis/pull/10578)

   We now restrict the number of workers that look up shares to 5. The number can
   be changed with `SHARING_USER_JSONCS3_MAX_CONCURRENCY` or
   `OCIS_MAX_CONCURRENCY`.

   https://github.com/owncloud/ocis/pull/10578
   https://github.com/owncloud/ocis/pull/10552

* Bugfix - Set MaxConcurrency to 1: [#10580](https://github.com/owncloud/ocis/pull/10580)

   Set MaxConcurrency for frontend and userlog and sharing services to 1. Too many
   workers will negatively impact performance on small machines.

   https://github.com/owncloud/ocis/pull/10580
   https://github.com/owncloud/ocis/pull/10557

* Bugfix - Reuse go-micro service clients: [#10582](https://github.com/owncloud/ocis/pull/10582)

   Go micro clients must not be reinitialized. The internal selector will spawn a
   new go routine to watch for registry changes.

   https://github.com/owncloud/ocis/pull/10582

* Bugfix - Make collaboration service use a gateway selector: [#10584](https://github.com/owncloud/ocis/pull/10584)

   https://github.com/owncloud/ocis/pull/10584

* Bugfix - Return an error if we can't get the keys and ensure they're cached: [#10590](https://github.com/owncloud/ocis/pull/10590)

   Previously, there was an issue where we could get an error while getting the
   public keys from the /hosting/discovery endpoint but we're returning a wrong
   success value instead. This is fixed now and we're returning the error.

   In addition, the public keys weren't being cached, so we hit the
   /hosting/discovery endpoint every time we need to use the public keys. The keys
   are now cached so we don't need to hit the endpoint more than what we need.

   https://github.com/owncloud/ocis/pull/10590

* Bugfix - Fix status code for thumbnail requests: [#10592](https://github.com/owncloud/ocis/pull/10592)

   We fixed the status code returned by the thumbnails service when the image
   source for a thumbnail exceeds the configured maximum dimensions or file size.
   The service now returns a 403 Forbidden status code instead of a 500 Internal
   Server Error status code.

   https://github.com/owncloud/ocis/issues/10589
   https://github.com/owncloud/ocis/pull/10592

* Bugfix - Fix the activity field mapping: [#10593](https://github.com/owncloud/ocis/pull/10593)

   https://github.com/owncloud/ocis/issues/10228
   https://github.com/owncloud/ocis/pull/10593
   Fixed
   the
   activity
   field
   mapping

* Enhancement - Update web to v11.0.3: [#10569](https://github.com/owncloud/ocis/pull/10569)

   Tags: web

   We updated ownCloud Web to v11.0.3. Please refer to the changelog (linked) for
   details on the web release.

   - Bugfix [owncloud/web#11870](https://github.com/owncloud/web/issues/11870):
   Preview image retries postprocessing - Bugfix
   [owncloud/web#11883](https://github.com/owncloud/web/issues/11883): Preview app
   Shared with me page - Bugfix
   [owncloud/web#11897](https://github.com/owncloud/web/issues/11897): "Save as" /
   "Open" when embed delegate authentication is enabled - Bugfix
   [owncloud/web#11900](https://github.com/owncloud/web/issues/11900): App top bar
   does not show location when shared file is opened - Bugfix
   [owncloud/web#11900](https://github.com/owncloud/web/issues/11900): Open from
   app and Save As feature broken when opened via shared file - Bugfix
   [owncloud/web#11904](https://github.com/owncloud/web/issues/11904): Public
   folder reload

   https://github.com/owncloud/ocis/pull/10569
   https://github.com/owncloud/web/releases/tag/v11.0.3

# Changelog for [5.0.9] (2024-11-14)

The following sections list the changes for 5.0.9.

[5.0.9]: https://github.com/owncloud/ocis/compare/v7.0.0-rc.2...v5.0.9

## Summary

* Bugfix - Thumbnail request limit: [#10280](https://github.com/owncloud/ocis/pull/10280)
* Bugfix - Restart Postprocessing properly: [#10439](https://github.com/owncloud/ocis/pull/10439)
* Change - Define maximum input image dimensions and size when generating previews: [#10270](https://github.com/owncloud/ocis/pull/10270)

## Details

* Bugfix - Thumbnail request limit: [#10280](https://github.com/owncloud/ocis/pull/10280)

   The `THUMBNAILS_MAX_CONCURRENT_REQUESTS` setting was not working correctly.
   Previously it was just limiting the number of concurrent thumbnail downloads.
   Now the limit is applied to the number thumbnail generations requests.
   Additionally the webdav service is now returning a "Retry-After" header when it
   is hitting the ratelimit of the thumbnail service.

   https://github.com/owncloud/ocis/pull/10280
   https://github.com/owncloud/ocis/pull/10270
   https://github.com/owncloud/ocis/pull/10225

* Bugfix - Restart Postprocessing properly: [#10439](https://github.com/owncloud/ocis/pull/10439)

   Properly differentiate between resume and restart postprocessing.

   https://github.com/owncloud/ocis/pull/10439

* Change - Define maximum input image dimensions and size when generating previews: [#10270](https://github.com/owncloud/ocis/pull/10270)

   This is a general hardening change to limit processing time and resources of the
   thumbnailer.

   https://github.com/owncloud/ocis/pull/10270
   https://github.com/owncloud/ocis/pull/9360
   https://github.com/owncloud/ocis/pull/9035
   https://github.com/owncloud/ocis/pull/9069

# Changelog for [7.0.0-rc.2] (2024-11-12)

The following sections list the changes for 7.0.0-rc.2.

[7.0.0-rc.2]: https://github.com/owncloud/ocis/compare/v7.0.0-rc.1...v7.0.0-rc.2

## Summary

* Bugfix - Fix idp guest role default assignment: [#10511](https://github.com/owncloud/ocis/pull/10511)
* Bugfix - Remove mbreaker: [#10524](https://github.com/owncloud/ocis/pull/10524)
* Bugfix - Bump Reva to v2.26.5: [#10552](https://github.com/owncloud/ocis/pull/10552)

## Details

* Bugfix - Fix idp guest role default assignment: [#10511](https://github.com/owncloud/ocis/pull/10511)

   We fixed an idp guest role default assignment.

   https://github.com/owncloud/ocis/issues/10474
   https://github.com/owncloud/ocis/pull/10511

* Bugfix - Remove mbreaker: [#10524](https://github.com/owncloud/ocis/pull/10524)

   The circuit breaker is not handle correctly and leads therefore to more issues
   than it solves. We removed it.

   https://github.com/owncloud/ocis/pull/10524

* Bugfix - Bump Reva to v2.26.5: [#10552](https://github.com/owncloud/ocis/pull/10552)

  * Fix [cs3org/reva#4926](https://github.com/cs3org/reva/issues/4926): Make etag always match content on downloads
  * Fix [cs3org/reva#4920](https://github.com/cs3org/reva/issues/4920): Return correct status codes for simple uploads
  * Fix [cs3org/reva#4924](https://github.com/cs3org/reva/issues/4924): Fix sync propagation
  * Fix [cs3org/reva#4916](https://github.com/cs3org/reva/issues/4916): Improve posixfs stability and performanc

   https://github.com/owncloud/ocis/pull/10552
   https://github.com/owncloud/ocis/pull/10539

# Changelog for [7.0.0-rc.1] (2024-11-07)

The following sections list the changes for 7.0.0-rc.1.

[7.0.0-rc.1]: https://github.com/owncloud/ocis/compare/v6.6.1...v7.0.0-rc.1

## Summary

* Bugfix - Generate short tokens to be used as access tokens for WOPI: [#10391](https://github.com/owncloud/ocis/pull/10391)
* Bugfix - Fix put relative wopi operation for microsoft: [#10403](https://github.com/owncloud/ocis/pull/10403)
* Bugfix - Make SSE keepalive interval configurable: [#10411](https://github.com/owncloud/ocis/pull/10411)
* Bugfix - Removed 'OCM_OCM_PROVIDER_AUTHORIZER_VERIFY_REQUEST_HOSTNAME' setting: [#10425](https://github.com/owncloud/ocis/pull/10425)
* Bugfix - Micro registry cache fixes: [#10429](https://github.com/owncloud/ocis/pull/10429)
* Bugfix - Fix the memlimit loglevel: [#10433](https://github.com/owncloud/ocis/pull/10433)
* Bugfix - Restart Postprocessing properly: [#10439](https://github.com/owncloud/ocis/pull/10439)
* Bugfix - Allow to configure data server URL for ocm: [#10440](https://github.com/owncloud/ocis/pull/10440)
* Bugfix - Respect proxy url when validating proofkeys: [#10462](https://github.com/owncloud/ocis/pull/10462)
* Bugfix - Return wopi lock header in get lock response: [#10470](https://github.com/owncloud/ocis/pull/10470)
* Bugfix - 'ocis backup consistency' fixed for file revisions: [#10493](https://github.com/owncloud/ocis/pull/10493)
* Bugfix - Wait for services to be ready before registering them: [#10498](https://github.com/owncloud/ocis/pull/10498)
* Bugfix - Fix 0-byte file uploads: [#10500](https://github.com/owncloud/ocis/pull/10500)
* Bugfix - Fixed `sharedWithMe` response for OCM shares: [#10501](https://github.com/owncloud/ocis/pull/10501)
* Bugfix - Fix gateway nats checks: [#10502](https://github.com/owncloud/ocis/pull/10502)
* Enhancement - Create thumbnails for GGP MIME types: [#10304](https://github.com/owncloud/ocis/pull/10304)
* Enhancement - Include a product name in the collaboration service: [#10335](https://github.com/owncloud/ocis/pull/10335)
* Enhancement - Add web extensions to the ocis_full example: [#10399](https://github.com/owncloud/ocis/pull/10399)
* Enhancement - Bump reva to 2.26.4: [#10419](https://github.com/owncloud/ocis/pull/10419)
* Enhancement - Remove deprecated CLI commands: [#10432](https://github.com/owncloud/ocis/pull/10432)
* Enhancement - Bump cs3api: [#10449](https://github.com/owncloud/ocis/pull/10449)
* Enhancement - Update web to v11.0.2: [#10467](https://github.com/owncloud/ocis/pull/10467)
* Enhancement - Bump reva to latest: [#10472](https://github.com/owncloud/ocis/pull/10472)
* Enhancement - Concurrent userlog processing: [#10504](https://github.com/owncloud/ocis/pull/10504)
* Enhancement - Concurrent autoaccept for shares: [#10507](https://github.com/owncloud/ocis/pull/10507)

## Details

* Bugfix - Generate short tokens to be used as access tokens for WOPI: [#10391](https://github.com/owncloud/ocis/pull/10391)

   Currently, the access tokens being used might be too long. In particular,
   Microsoft Office Online complains about the URL (which contains the access
   token) is too long and refuses to work.

   https://github.com/owncloud/ocis/pull/10391

* Bugfix - Fix put relative wopi operation for microsoft: [#10403](https://github.com/owncloud/ocis/pull/10403)

   We fixed a bug in the put relative wopi operation for microsoft. The response
   now contains the correct properties.

   https://github.com/owncloud/ocis/pull/10403

* Bugfix - Make SSE keepalive interval configurable: [#10411](https://github.com/owncloud/ocis/pull/10411)

   To prevent intermediate proxies from closing the SSE connection admins can now
   configure a `SSE_KEEPALIVE_INTERVAL`.

   https://github.com/owncloud/ocis/pull/10411

* Bugfix - Removed 'OCM_OCM_PROVIDER_AUTHORIZER_VERIFY_REQUEST_HOSTNAME' setting: [#10425](https://github.com/owncloud/ocis/pull/10425)

   The config option 'OCM_OCM_PROVIDER_AUTHORIZER_VERIFY_REQUEST_HOSTNAME' was
   removed from the OCM service. The additional security provided by this setting
   is somewhat questionable and only provided in very specific setups.

   We are not going through the normal deprecation process for this setting, as it
   was never really working anyway. If you have this setting in your configuration,
   it will be ignored. You can safely remove it.

   https://github.com/owncloud/ocis/issues/10355
   https://github.com/owncloud/ocis/pull/10425

* Bugfix - Micro registry cache fixes: [#10429](https://github.com/owncloud/ocis/pull/10429)

   We now invalidate cache entries when any of the nodes was not updated.

   https://github.com/owncloud/ocis/pull/10429

* Bugfix - Fix the memlimit loglevel: [#10433](https://github.com/owncloud/ocis/pull/10433)

   We set the memlimit default loglevel to error.

   https://github.com/owncloud/ocis/issues/10427
   https://github.com/owncloud/ocis/pull/10433

* Bugfix - Restart Postprocessing properly: [#10439](https://github.com/owncloud/ocis/pull/10439)

   Properly differentiate between resume and restart postprocessing.

   https://github.com/owncloud/ocis/pull/10439

* Bugfix - Allow to configure data server URL for ocm: [#10440](https://github.com/owncloud/ocis/pull/10440)

   We introduced the `OCM_OCM_STORAGE_DATA_SERVER_URL` setting to fix a bug when
   downloading files from an OCM share. Before the data server URL defaulted to the
   listen address of the OCM server, which did not work when using 0.0.0.0 as the
   listen address.

   https://github.com/owncloud/ocis/issues/10358
   https://github.com/owncloud/ocis/pull/10440

* Bugfix - Respect proxy url when validating proofkeys: [#10462](https://github.com/owncloud/ocis/pull/10462)

   We fixed a bug where the proxied wopi URL was not used when validating
   proofkeys. This caused the validation to fail when the proxy was used.

   https://github.com/owncloud/ocis/pull/10462

* Bugfix - Return wopi lock header in get lock response: [#10470](https://github.com/owncloud/ocis/pull/10470)

   We fixed a bug where the wopi lock header was not returned in the get lock
   response. This is now fixed and the wopi validator tests are passing.

   https://github.com/owncloud/ocis/pull/10470

* Bugfix - 'ocis backup consistency' fixed for file revisions: [#10493](https://github.com/owncloud/ocis/pull/10493)

   A bug was fixed that caused the 'ocis backup consistency' command to incorrectly
   report inconistencies when file revisions with a zero value for the nano-second
   part of the timestamp were present.

   https://github.com/owncloud/ocis/issues/9498
   https://github.com/owncloud/ocis/pull/10493

* Bugfix - Wait for services to be ready before registering them: [#10498](https://github.com/owncloud/ocis/pull/10498)

   https://github.com/owncloud/ocis/pull/10498

* Bugfix - Fix 0-byte file uploads: [#10500](https://github.com/owncloud/ocis/pull/10500)

   We fixed an issue where 0-byte files upload did not return the Location header.

   https://github.com/owncloud/ocis/issues/10469
   https://github.com/owncloud/ocis/pull/10500

* Bugfix - Fixed `sharedWithMe` response for OCM shares: [#10501](https://github.com/owncloud/ocis/pull/10501)

   OCM shares returned in the `sharedWithMe` response did not have the `mimeType`
   property populated correctly.

   https://github.com/owncloud/ocis/issues/10495
   https://github.com/owncloud/ocis/pull/10501

* Bugfix - Fix gateway nats checks: [#10502](https://github.com/owncloud/ocis/pull/10502)

   We now only check if nats is available when the gateway actually uses it.
   Furthermore, we added a backoff for checking the readys endpoint.

   https://github.com/owncloud/ocis/pull/10502

* Enhancement - Create thumbnails for GGP MIME types: [#10304](https://github.com/owncloud/ocis/pull/10304)

   Creates thumbnails for newly added ggp files

   https://github.com/owncloud/ocis/pull/10304

* Enhancement - Include a product name in the collaboration service: [#10335](https://github.com/owncloud/ocis/pull/10335)

   The product name will allow using a different app name. For example, a "CoolBox"
   app name might use a branded Collabora instance by using "Collabora" as product
   name.

   https://github.com/owncloud/ocis/pull/10335
   https://github.com/owncloud/ocis/pull/10490

* Enhancement - Add web extensions to the ocis_full example: [#10399](https://github.com/owncloud/ocis/pull/10399)

   We added some of the web extensions from ownCloud to the ocis_full docker
   compose example.

   - importer - draw-io - external-sites - json-viewer - unzip - progressbars

   These can be enabled in the .env file one by one.

   Read more about ocis extensions in
   https://github.com/owncloud/web-extensions/blob/main/README.md

   https://github.com/owncloud/ocis/pull/10399

* Enhancement - Bump reva to 2.26.4: [#10419](https://github.com/owncloud/ocis/pull/10419)

  *   Bugfix [cs3org/reva#4917](https://github.com/cs3org/reva/pull/4917): Fix 0-byte file uploads
  *   Bugfix [cs3org/reva#4918](https://github.com/cs3org/reva/pull/4918): Fix app templates

   Bump reva to 2.26.3

  *   Bugfix [cs3org/reva#4908](https://github.com/cs3org/reva/pull/4908): Add checksum to OCM storageprovider responses
  *   Enhancement [cs3org/reva#4910](https://github.com/cs3org/reva/pull/4910): Bump cs3api
  *   Enhancement [cs3org/reva#4909](https://github.com/cs3org/reva/pull/4909): Bump cs3api
  *   Enhancement [cs3org/reva#4906](https://github.com/cs3org/reva/pull/4906): Bump cs3api

   Bump reva to 2.26.2

  *   Enhancement [cs3org/reva#4897](https://github.com/cs3org/reva/pull/4897): Fix remaining quota calculation
  *   Bugfix      [cs3org/reva#4902](https://github.com/cs3org/reva/pull/4902): Fix quota calculation

   https://github.com/owncloud/ocis/pull/10419

* Enhancement - Remove deprecated CLI commands: [#10432](https://github.com/owncloud/ocis/pull/10432)

   We removed the following deprecated CLI commands: `ocis storage-users uploads
   list` `ocis storage-users uploads clean`

   https://github.com/owncloud/ocis/issues/10428
   https://github.com/owncloud/ocis/pull/10432

* Enhancement - Bump cs3api: [#10449](https://github.com/owncloud/ocis/pull/10449)

   https://github.com/owncloud/ocis/pull/10449

* Enhancement - Update web to v11.0.2: [#10467](https://github.com/owncloud/ocis/pull/10467)

   Tags: web

   We updated ownCloud Web to v11.0.2. Please refer to the changelog (linked) for
   details on the web release.

   - Bugfix [owncloud/web#11803](https://github.com/owncloud/web/issues/11803):
   Files appearing in file list - Bugfix
   [owncloud/web#11804](https://github.com/owncloud/web/pull/11804): Add missing
   translations - Bugfix
   [owncloud/web#11806](https://github.com/owncloud/web/issues/11806): Folder size
   mismatch - Bugfix
   [owncloud/web#11813](https://github.com/owncloud/web/pull/11813): Preview image
   retries - Bugfix
   [owncloud/web#11817](https://github.com/owncloud/web/pull/11817): Respect post
   logout url - Bugfix
   [owncloud/web#11830](https://github.com/owncloud/web/issues/11830): Right side
   bar cut off - Bugfix
   [owncloud/web#11831](https://github.com/owncloud/web/pull/11831): Sidebar
   showing wrong shares - Bugfix
   [owncloud/web#11853](https://github.com/owncloud/web/issues/11853): Context menu
   "Open app in new tab" broken - Bugfix
   [owncloud/web#11008](https://github.com/owncloud/web/issues/11008): Show lock
   information in file details

   https://github.com/owncloud/ocis/pull/10467
   https://github.com/owncloud/ocis/pull/10503
   https://github.com/owncloud/web/releases/tag/v11.0.1
   https://github.com/owncloud/web/releases/tag/v11.0.2

* Enhancement - Bump reva to latest: [#10472](https://github.com/owncloud/ocis/pull/10472)

   https://github.com/owncloud/ocis/pull/10472

* Enhancement - Concurrent userlog processing: [#10504](https://github.com/owncloud/ocis/pull/10504)

   We now start multiple go routines that process events. The default of 5
   goroutines can be changed with the new `USERLOG_MAX_CONCURRENCY` environment
   variable.

   https://github.com/owncloud/ocis/pull/10504

* Enhancement - Concurrent autoaccept for shares: [#10507](https://github.com/owncloud/ocis/pull/10507)

   Shares for groups are now concurrently accepted. Tha default of 25 goroutinges
   can be changed with the new `FRONTEND_MAX_CONCURRENCY` environment variable.

   https://github.com/owncloud/ocis/pull/10507
   https://github.com/owncloud/ocis/pull/10476

# Changelog for [6.6.1] (2024-10-24)

The following sections list the changes for 6.6.1.

[6.6.1]: https://github.com/owncloud/ocis/compare/v6.6.0...v6.6.1

## Summary

* Bugfix - Fix panic when stopping the nats: [#10363](https://github.com/owncloud/ocis/pull/10363)
* Bugfix - Disable download activity: [#10368](https://github.com/owncloud/ocis/pull/10368)
* Bugfix - Fix Activitylog issues: [#10376](https://github.com/owncloud/ocis/pull/10376)
* Bugfix - Security fixes: [#10376](https://github.com/owncloud/ocis/pull/10376)
* Bugfix - Make antivirus workers configurable: [#10383](https://github.com/owncloud/ocis/pull/10383)
* Bugfix - Increase event processing workers: [#10385](https://github.com/owncloud/ocis/pull/10385)
* Bugfix - Fix envvar deprecations for next production release: [#10386](https://github.com/owncloud/ocis/pull/10386)
* Bugfix - Fix healthchecks: [#10405](https://github.com/owncloud/ocis/pull/10405)

## Details

* Bugfix - Fix panic when stopping the nats: [#10363](https://github.com/owncloud/ocis/pull/10363)

   The nats server itself runs signal handling that the Shutdown() call in the ocis
   code is redundant and led to a panic.

   https://github.com/owncloud/ocis/issues/10360
   https://github.com/owncloud/ocis/pull/10363

* Bugfix - Disable download activity: [#10368](https://github.com/owncloud/ocis/pull/10368)

   We disable the download activity until we have a proper solution for it.

   https://github.com/owncloud/ocis/issues/10293
   https://github.com/owncloud/ocis/pull/10368

* Bugfix - Fix Activitylog issues: [#10376](https://github.com/owncloud/ocis/pull/10376)

   Fixes multiple activititylog issues. There was an error about `max payload
   exceeded` when there were too many activities on one folder. Listing would take
   very long even with a limit activated. All of these issues are now fixed.

   https://github.com/owncloud/ocis/pull/10376

* Bugfix - Security fixes: [#10376](https://github.com/owncloud/ocis/pull/10376)

   We polished some of the sonarcloud issues.

   https://github.com/owncloud/ocis/pull/10376

* Bugfix - Make antivirus workers configurable: [#10383](https://github.com/owncloud/ocis/pull/10383)

   We made the number of go routines that pull events from the queue configurable.

   https://github.com/owncloud/ocis/pull/10383

* Bugfix - Increase event processing workers: [#10385](https://github.com/owncloud/ocis/pull/10385)

   We increased the number of go routines that pull events from the queue to three
   and made the number off workers configurable. Furthermore, the postprocessing
   delay no longer introduces a sleep that slows down pulling of events, but
   asynchronously triggers the next step.

   https://github.com/owncloud/ocis/pull/10385
   https://github.com/owncloud/ocis/pull/10368

* Bugfix - Fix envvar deprecations for next production release: [#10386](https://github.com/owncloud/ocis/pull/10386)

   Some envvar deprecations were incomplete. One was missed to be removed, one had
   missing information.

   https://github.com/owncloud/ocis/pull/10386

* Bugfix - Fix healthchecks: [#10405](https://github.com/owncloud/ocis/pull/10405)

   We needed to replace 0.0.0.0 bind addresses by outbound IP addresses in the
   healthcheck routine.

   https://github.com/owncloud/ocis/pull/10405

# Changelog for [6.6.0] (2024-10-21)

The following sections list the changes for 6.6.0.

[6.6.0]: https://github.com/owncloud/ocis/compare/v6.5.0...v6.6.0

## Summary

* Bugfix - Fix health and ready endpoints: [#10163](https://github.com/owncloud/ocis/pull/10163)
* Bugfix - Always treat LDAP attribute names case-insensitively: [#10204](https://github.com/owncloud/ocis/pull/10204)
* Bugfix - Fix delete share panic: [#10219](https://github.com/owncloud/ocis/pull/10219)
* Bugfix - Continue listing shares on error: [#10243](https://github.com/owncloud/ocis/pull/10243)
* Bugfix - Avoid re-creating thumbnails: [#10251](https://github.com/owncloud/ocis/pull/10251)
* Bugfix - Graph service now supports `OCIS_LDAP_USER_SCHEMA_DISPLAYNAME` env var: [#10257](https://github.com/owncloud/ocis/issues/10257)
* Bugfix - Kept historical resource naming in activity: [#10266](https://github.com/owncloud/ocis/pull/10266)
* Bugfix - Fix panic when sharing with groups: [#10279](https://github.com/owncloud/ocis/pull/10279)
* Bugfix - Thumbnail request limit: [#10280](https://github.com/owncloud/ocis/pull/10280)
* Bugfix - Forbid the ocm space sharing: [#10287](https://github.com/owncloud/ocis/pull/10287)
* Bugfix - Use secure config defaults for OCM: [#10307](https://github.com/owncloud/ocis/pull/10307)
* Enhancement - Add OCM wellknown configuration: [#9815](https://github.com/owncloud/ocis/pull/9815)
* Enhancement - Load IDP logo from theme: [#10274](https://github.com/owncloud/ocis/pull/10274)
* Enhancement - WebOffice Templates: [#10276](https://github.com/owncloud/ocis/pull/10276)
* Enhancement - Remove Deprecations: [#10305](https://github.com/owncloud/ocis/pull/10305)
* Enhancement - Allow to use libvips for generating thumbnails: [#10310](https://github.com/owncloud/ocis/pull/10310)
* Enhancement - Bump various dependencies: [#10352](https://github.com/owncloud/ocis/pull/10352)
* Enhancement - Update web to v11.0.0: [#10357](https://github.com/owncloud/ocis/pull/10357)
* Enhancement - Bump reva to 2.26.0: [#10364](https://github.com/owncloud/ocis/pull/10364)

## Details

* Bugfix - Fix health and ready endpoints: [#10163](https://github.com/owncloud/ocis/pull/10163)

   We added new checks to the `/readyz` and `/healthz` endpoints to ensure that the
   services are ready and healthy. This change ensures that the endpoints return
   the correct status codes, which is needed to stabilize the k8s deployments.

   https://github.com/owncloud/ocis/issues/10316
   https://github.com/owncloud/ocis/issues/10281
   https://github.com/owncloud/ocis/pull/10163
   https://github.com/owncloud/ocis/pull/10301
   https://github.com/owncloud/ocis/pull/10302
   https://github.com/owncloud/ocis/pull/10303
   https://github.com/owncloud/ocis/pull/10308
   https://github.com/owncloud/ocis/pull/10323
   https://github.com/owncloud/ocis/pull/10163
   https://github.com/owncloud/ocis/pull/10333

* Bugfix - Always treat LDAP attribute names case-insensitively: [#10204](https://github.com/owncloud/ocis/pull/10204)

   We fixes a bug where some LDAP attributes (e.g. owncloudUUID) were not treated
   case-insensitively.

   https://github.com/owncloud/ocis/issues/10200
   https://github.com/owncloud/ocis/pull/10204

* Bugfix - Fix delete share panic: [#10219](https://github.com/owncloud/ocis/pull/10219)

   Fixes a panic when deleting an ocm share

   https://github.com/owncloud/ocis/pull/10219

* Bugfix - Continue listing shares on error: [#10243](https://github.com/owncloud/ocis/pull/10243)

   We now continue listing received shares when one of the shares cannot be statted
   or converted to a driveItem.

   https://github.com/owncloud/ocis/pull/10243

* Bugfix - Avoid re-creating thumbnails: [#10251](https://github.com/owncloud/ocis/pull/10251)

   We fixed a bug that caused the system to re-create thumbnails for images, even
   if a thumbnail already existed in the cache.

   https://github.com/owncloud/ocis/pull/10251

* Bugfix - Graph service now supports `OCIS_LDAP_USER_SCHEMA_DISPLAYNAME` env var: [#10257](https://github.com/owncloud/ocis/issues/10257)

   To align with the other services the graph service now supports the
   `OCIS_LDAP_USER_SCHEMA_DISPLAYNAME` environment variable to configure the LDAP
   attribute that is used for display name attribute of users.

   `LDAP_USER_SCHEMA_DISPLAY_NAME` is now deprecated and will be removed in a
   future release.

   https://github.com/owncloud/ocis/issues/10257

* Bugfix - Kept historical resource naming in activity: [#10266](https://github.com/owncloud/ocis/pull/10266)

   Kept historical resource naming after renaming in activity for shares and public
   links.

   https://github.com/owncloud/ocis/issues/10210
   https://github.com/owncloud/ocis/pull/10266

* Bugfix - Fix panic when sharing with groups: [#10279](https://github.com/owncloud/ocis/pull/10279)

   We fixed a bug which caused a panic when sharing with groups, this only happened
   under a heavy load. Besides the bugfix, we also reduced the number of share auto
   accept log messages to avoid flooding the logs.

   https://github.com/owncloud/ocis/issues/10258
   https://github.com/owncloud/ocis/pull/10279

* Bugfix - Thumbnail request limit: [#10280](https://github.com/owncloud/ocis/pull/10280)

   The `THUMBNAILS_MAX_CONCURRENT_REQUESTS` setting was not working correctly.
   Previously it was just limiting the number of concurrent thumbnail downloads.
   Now the limit is applied to the number thumbnail generations requests.
   Additionally the webdav service is now returning a "Retry-After" header when it
   is hitting the ratelimit of the thumbnail service.

   https://github.com/owncloud/ocis/pull/10280
   https://github.com/owncloud/ocis/pull/10225

* Bugfix - Forbid the ocm space sharing: [#10287](https://github.com/owncloud/ocis/pull/10287)

   We forbid adding the federated users as members of the space via items invite.

   https://github.com/owncloud/ocis/issues/10051
   https://github.com/owncloud/ocis/pull/10287

* Bugfix - Use secure config defaults for OCM: [#10307](https://github.com/owncloud/ocis/pull/10307)

   https://github.com/owncloud/ocis/pull/10307

* Enhancement - Add OCM wellknown configuration: [#9815](https://github.com/owncloud/ocis/pull/9815)

   We now configure the `wellknown` service for OCM.

   https://github.com/owncloud/ocis/pull/9815

* Enhancement - Load IDP logo from theme: [#10274](https://github.com/owncloud/ocis/pull/10274)

   We now load the IDP logo from the theme file.

   https://github.com/owncloud/web/issues/11603
   https://github.com/owncloud/ocis/pull/10274

* Enhancement - WebOffice Templates: [#10276](https://github.com/owncloud/ocis/pull/10276)

   We are now able to use templates for WebOffice and use them as a starting point
   for new documents.

   We are supporting the following mime types:

   ## OnlyOffice

   - **MimeType:** `application/vnd.ms-word.template.macroenabled.12`
   **TargetExtension:** `docx`

   - **MimeType:** `application/vnd.oasis.opendocument.text-template`
   **TargetExtension:** `docx`

   - **MimeType:**
   `application/vnd.openxmlformats-officedocument.wordprocessingml.template`
   **TargetExtension:** `docx`

   - **MimeType:** `application/vnd.oasis.opendocument.spreadsheet-template`
   **TargetExtension:** `xlsx`

   - **MimeType:** `application/vnd.ms-excel.template.macroenabled.12`
   **TargetExtension:** `xlsx`

   - **MimeType:**
   `application/vnd.openxmlformats-officedocument.spreadsheetml.template`
   **TargetExtension:** `xlsx`

   - **MimeType:** `application/vnd.oasis.opendocument.presentation-template`
   **TargetExtension:** `pptx`

   - **MimeType:** `application/vnd.ms-powerpoint.template.macroenabled.12`
   **TargetExtension:** `pptx`

   - **MimeType:**
   `application/vnd.openxmlformats-officedocument.presentationml.template`
   **TargetExtension:** `pptx`

   ## Collabora

   - **MimeType:** `application/vnd.oasis.opendocument.spreadsheet-template`
   **TargetExtension:** `ods`

   - **MimeType:** `application/vnd.oasis.opendocument.text-template`
   **TargetExtension:** `odt`

   - **MimeType:** `application/vnd.oasis.opendocument.presentation-template`
   **TargetExtension:** `odp`

   https://github.com/owncloud/ocis/issues/9785
   https://github.com/owncloud/ocis/pull/10276

* Enhancement - Remove Deprecations: [#10305](https://github.com/owncloud/ocis/pull/10305)

   Remove deprecated stores/caches/registries and envvars from the codebase.

   https://github.com/owncloud/ocis/pull/10305

* Enhancement - Allow to use libvips for generating thumbnails: [#10310](https://github.com/owncloud/ocis/pull/10310)

   To improve performance (and to be able to support a wider range of images
   formats in the future) the thumbnails service is now able to utilize libvips
   (https://www.libvips.org/) for generating thumbnails. Enabling the use of
   libvips is implemented as a build-time option which is currently disabled for
   the "bare-metal" build of the ocis binary and enabled for the docker image
   builds.

   https://github.com/owncloud/ocis/pull/10310

* Enhancement - Bump various dependencies: [#10352](https://github.com/owncloud/ocis/pull/10352)

   https://github.com/owncloud/ocis/pull/10352

* Enhancement - Update web to v11.0.0: [#10357](https://github.com/owncloud/ocis/pull/10357)

   Tags: web

   We updated ownCloud Web to v11.0.0. Please refer to the changelog (linked) for
   details on the web release.

   - Change [owncloud/web#11709](https://github.com/owncloud/web/pull/11709):
   Remove importer as default app - Enhancement
   [owncloud/web#11668](https://github.com/owncloud/web/pull/11668): Allow setting
   view mode for apps via query - Enhancement
   [owncloud/web#11731](https://github.com/owncloud/web/pull/11731): File size
   warning in editors - Enhancement
   [owncloud/web#11737](https://github.com/owncloud/web/pull/11737): Add not found
   page - Enhancement
   [owncloud/web#11750](https://github.com/owncloud/web/pull/11750): Create
   documents from templates - Bugfix
   [owncloud/web#11604](https://github.com/owncloud/web/pull/11604): User filters
   after page reload - Bugfix
   [owncloud/web#11645](https://github.com/owncloud/web/pull/11645): Hide copy
   permanent link action on public pages - Bugfix
   [owncloud/web#11677](https://github.com/owncloud/web/pull/11677): Missing tags
   on "Shared with me" page - Bugfix
   [owncloud/web#11678](https://github.com/owncloud/web/pull/11678): Undefined
   request IDs - Bugfix
   [owncloud/web#11688](https://github.com/owncloud/web/pull/11688): Deleting
   federated connections - Bugfix
   [owncloud/web#11706](https://github.com/owncloud/web/pull/11706): Escape HTML
   characters in activities and notification view - Bugfix
   [owncloud/web#11707](https://github.com/owncloud/web/pull/11707): Prevent not
   allowed characters in shortcut name - Bugfix
   [owncloud/web#11712](https://github.com/owncloud/web/pull/11712): Details panel
   wrong WebDAV URL of received shares - Bugfix
   [owncloud/web#11725](https://github.com/owncloud/web/pull/11725): Accessing
   disabled password-protected space does not show error - Bugfix
   [owncloud/web#11726](https://github.com/owncloud/web/pull/11726): Application
   menu not operable in Safari browser - Bugfix
   [owncloud/web#11758](https://github.com/owncloud/web/pull/11758): Navigating
   into folders that have been shared externally - Bugfix
   [owncloud/web#11795](https://github.com/owncloud/web/pull/11795): Sharing label
   for locked files

   https://github.com/owncloud/ocis/pull/10357
   https://github.com/owncloud/web/releases/tag/v11.0.0

* Enhancement - Bump reva to 2.26.0: [#10364](https://github.com/owncloud/ocis/pull/10364)

  *   Bugfix [cs3org/reva#4880](https://github.com/cs3org/reva/pull/4880): Kept historical resource naming in activity
  *   Bugfix [cs3org/reva#4874](https://github.com/cs3org/reva/pull/4874): Fix rename activity
  *   Bugfix [cs3org/reva#4881](https://github.com/cs3org/reva/pull/4881): Log levels
  *   Bugfix [cs3org/reva#4884](https://github.com/cs3org/reva/pull/4884): Fix OCM upload crush
  *   Bugfix [cs3org/reva#4872](https://github.com/cs3org/reva/pull/4872): Return 409 conflict when a file was already created
  *   Bugfix [cs3org/reva#4887](https://github.com/cs3org/reva/pull/4887): Fix ShareCache concurrency panic
  *   Bugfix [cs3org/reva#4876](https://github.com/cs3org/reva/pull/4876): Fix share jail mountpoint parent id
  *   Bugfix [cs3org/reva#4879](https://github.com/cs3org/reva/pull/4879): Fix trash-bin propfind panic
  *   Bugfix [cs3org/reva#4888](https://github.com/cs3org/reva/pull/4888): Fix upload session bugs
  *   Bugfix [cs3org/reva#4560](https://github.com/cs3org/reva/pull/4560): Always select next before making CS3 calls for propfinds
  *   Enhancement [cs3org/reva#4893](https://github.com/cs3org/reva/pull/4893): Bump dependencies and go to 1.22.8
  *   Enhancement [cs3org/reva#4890](https://github.com/cs3org/reva/pull/4890): Bump golangci-lint to 1.61.0
  *   Enhancement [cs3org/reva#4886](https://github.com/cs3org/reva/pull/4886): Add new Mimetype ggp
  *   Enhancement [cs3org/reva#4809](https://github.com/cs3org/reva/pull/4809): Implement OCM well-known endpoint
  *   Enhancement [cs3org/reva#4889](https://github.com/cs3org/reva/pull/4889): Improve posixfs stability and performance
  *   Enhancement [cs3org/reva#4882](https://github.com/cs3org/reva/pull/4882): Indicate template conversion capabilities on apps

   https://github.com/owncloud/ocis/pull/10364
   https://github.com/owncloud/ocis/pull/10347
   https://github.com/owncloud/ocis/pull/10321
   https://github.com/owncloud/ocis/pull/10236
   https://github.com/owncloud/ocis/pull/10216
   https://github.com/owncloud/ocis/pull/10315

# Changelog for [6.5.0] (2024-10-01)

The following sections list the changes for 6.5.0.

[6.5.0]: https://github.com/owncloud/ocis/compare/v5.0.8...v6.5.0

## Summary

* Bugfix - Fixed the ocm email template: [#10030](https://github.com/owncloud/ocis/pull/10030)
* Bugfix - Fixed activity filter depth: [#10031](https://github.com/owncloud/ocis/pull/10031)
* Bugfix - Fixed proxy build info: [#10039](https://github.com/owncloud/ocis/pull/10039)
* Bugfix - Fixed the ocm tocken: [#10050](https://github.com/owncloud/ocis/pull/10050)
* Bugfix - Fix ocm space sharing: [#10060](https://github.com/owncloud/ocis/pull/10060)
* Bugfix - Fix the error code for ocm space sharing: [#10079](https://github.com/owncloud/ocis/pull/10079)
* Bugfix - Added LinkUpdated activity: [#10085](https://github.com/owncloud/ocis/pull/10085)
* Bugfix - Fix Activities leak: [#10092](https://github.com/owncloud/ocis/pull/10092)
* Bugfix - Include additional logs in the collaboration service: [#10101](https://github.com/owncloud/ocis/pull/10101)
* Bugfix - Added ShareUpdate activity: [#10104](https://github.com/owncloud/ocis/pull/10104)
* Bugfix - Fixed the collaboration service registration: [#10107](https://github.com/owncloud/ocis/pull/10107)
* Bugfix - CheckFileInfo will return a 404 error if the target file isn't found: [#10112](https://github.com/owncloud/ocis/pull/10112)
* Bugfix - Forbid Activities for Sharees: [#10136](https://github.com/owncloud/ocis/pull/10136)
* Bugfix - Always select next gateway client: [#10141](https://github.com/owncloud/ocis/pull/10141)
* Bugfix - Remove duplicate CSP header from responses: [#10146](https://github.com/owncloud/ocis/pull/10146)
* Bugfix - Fixed the missing folder variable: [#10150](https://github.com/owncloud/ocis/pull/10150)
* Bugfix - Fix activity limit: [#10165](https://github.com/owncloud/ocis/pull/10165)
* Bugfix - Fix email translations: [#10171](https://github.com/owncloud/ocis/pull/10171)
* Bugfix - Fix Activities translation: [#10175](https://github.com/owncloud/ocis/pull/10175)
* Enhancement - Allow to maintain the last sign-in timestamp of a user: [#9942](https://github.com/owncloud/ocis/pull/9942)
* Enhancement - Add an Activity for FileUpdated: [#10072](https://github.com/owncloud/ocis/pull/10072)
* Enhancement - Remove METADATA_BACKEND: [#10113](https://github.com/owncloud/ocis/pull/10113)
* Enhancement - Load CSP configuration file if it exists: [#10139](https://github.com/owncloud/ocis/pull/10139)
* Enhancement - FileDownloaded Activity: [#10161](https://github.com/owncloud/ocis/pull/10161)
* Enhancement - Add WOPI host URLs to the collaboration service: [#10174](https://github.com/owncloud/ocis/pull/10174)
* Enhancement - Update web to v10.3.0: [#10177](https://github.com/owncloud/ocis/pull/10177)
* Enhancement - Bump reva to 2.25.0: [#10194](https://github.com/owncloud/ocis/pull/10194)

## Details

* Bugfix - Fixed the ocm email template: [#10030](https://github.com/owncloud/ocis/pull/10030)

   The golang conditional construction moved out from the transifex template.

   https://github.com/owncloud/ocis/pull/10030

* Bugfix - Fixed activity filter depth: [#10031](https://github.com/owncloud/ocis/pull/10031)

   Fixed activity filter 'depth:-1'

   https://github.com/owncloud/ocis/issues/9850
   https://github.com/owncloud/ocis/pull/10031

* Bugfix - Fixed proxy build info: [#10039](https://github.com/owncloud/ocis/pull/10039)

   The version string for the proxy service has been changed to 'version'.

   https://github.com/owncloud/ocis/pull/10039

* Bugfix - Fixed the ocm tocken: [#10050](https://github.com/owncloud/ocis/pull/10050)

   We now pass the JWT secret to the reva runtime.

   https://github.com/owncloud/ocis/pull/10050

* Bugfix - Fix ocm space sharing: [#10060](https://github.com/owncloud/ocis/pull/10060)

   We prevent adding the federated users as members of the space.

   https://github.com/owncloud/ocis/issues/10051
   https://github.com/owncloud/ocis/pull/10060

* Bugfix - Fix the error code for ocm space sharing: [#10079](https://github.com/owncloud/ocis/pull/10079)

   We fixed the error code for ocm space sharing

   https://github.com/owncloud/ocis/issues/10051
   https://github.com/owncloud/ocis/pull/10079

* Bugfix - Added LinkUpdated activity: [#10085](https://github.com/owncloud/ocis/pull/10085)

   Added the LinkUpdated activity in the space context

   https://github.com/owncloud/ocis/issues/10012
   https://github.com/owncloud/ocis/pull/10085

* Bugfix - Fix Activities leak: [#10092](https://github.com/owncloud/ocis/pull/10092)

   Fix activities endpoint by preventing unauthorized users to get activities

   https://github.com/owncloud/ocis/pull/10092

* Bugfix - Include additional logs in the collaboration service: [#10101](https://github.com/owncloud/ocis/pull/10101)

   More logs have been added in the middlware of the collaboration service to debug
   401 error codes. Any error that happens in that middleware should have its
   corresponding log entry

   https://github.com/owncloud/ocis/pull/10101

* Bugfix - Added ShareUpdate activity: [#10104](https://github.com/owncloud/ocis/pull/10104)

   Added the ShareUpdate activity in the space context.

   https://github.com/owncloud/ocis/issues/10011
   https://github.com/owncloud/ocis/pull/10104

* Bugfix - Fixed the collaboration service registration: [#10107](https://github.com/owncloud/ocis/pull/10107)

   Fixed an issue when the collaboration service registers apps also for binary and
   unknown mime types.

   https://github.com/owncloud/ocis/issues/10086
   https://github.com/owncloud/ocis/pull/10107

* Bugfix - CheckFileInfo will return a 404 error if the target file isn't found: [#10112](https://github.com/owncloud/ocis/pull/10112)

   Previously, the request failed with a 500 error code, but it it will fail with a
   404 error code

   https://github.com/owncloud/ocis/pull/10112

* Bugfix - Forbid Activities for Sharees: [#10136](https://github.com/owncloud/ocis/pull/10136)

   Sharees may not see item activities. We now bind it to ListGrants permission.

   https://github.com/owncloud/ocis/pull/10136

* Bugfix - Always select next gateway client: [#10141](https://github.com/owncloud/ocis/pull/10141)

   We now use the gateway selector to always select the next gateway client. This
   ensures that we can always connect to the gateway during up- and downscaling.

   https://github.com/owncloud/ocis/pull/10141
   https://github.com/owncloud/ocis/pull/10133

* Bugfix - Remove duplicate CSP header from responses: [#10146](https://github.com/owncloud/ocis/pull/10146)

   The web service was adding a CSP on its own, and that one has been removed. The
   proxy service will take care of the CSP header.

   https://github.com/owncloud/ocis/pull/10146

* Bugfix - Fixed the missing folder variable: [#10150](https://github.com/owncloud/ocis/pull/10150)

   We fixed the missing folder variable when folder renamed.

   https://github.com/owncloud/ocis/issues/10148
   https://github.com/owncloud/ocis/pull/10150

* Bugfix - Fix activity limit: [#10165](https://github.com/owncloud/ocis/pull/10165)

   When requesting a limit on activities, ocis would limit first, then filter and
   sort. Now it filters and sorts first, then limits.

   https://github.com/owncloud/ocis/pull/10165

* Bugfix - Fix email translations: [#10171](https://github.com/owncloud/ocis/pull/10171)

   Email translations would not use custom translation pathes. This is now fixed.

   https://github.com/owncloud/ocis/pull/10171

* Bugfix - Fix Activities translation: [#10175](https://github.com/owncloud/ocis/pull/10175)

   Fix the panic for the translation-sync in the activities service.

   https://github.com/owncloud/ocis/pull/10175

* Enhancement - Allow to maintain the last sign-in timestamp of a user: [#9942](https://github.com/owncloud/ocis/pull/9942)

   When the LDAP identity backend is configured to have write access to the
   database we're now able to maintain the ocLastSignInTimestamp attribute for the
   users.

   This attribute is return in the 'signinActivity/lastSuccessfulSignInDateTime'
   properity of the user objects. It is also possible to $filter on this attribute.

   Use e.g. '$filter=signinActivity/lastSuccessfulSignInDateTime le
   2023-12-31T00:00:00Z' to search for users that have not signed in since
   2023-12-31. Note: To use this type of filter the underlying LDAP server must
   support the '<=' filter. Which is currently not the case of the built-in LDAP
   server (idm).

   https://github.com/owncloud/ocis/pull/9942
   https://github.com/owncloud/ocis/pull/10111

* Enhancement - Add an Activity for FileUpdated: [#10072](https://github.com/owncloud/ocis/pull/10072)

   Previously FileUpdated has also triggered a FileAdded Activity

   https://github.com/owncloud/ocis/pull/10072

* Enhancement - Remove METADATA_BACKEND: [#10113](https://github.com/owncloud/ocis/pull/10113)

   Removes the deprecated XXX_METADATA_BACKEND envvars

   https://github.com/owncloud/ocis/pull/10113

* Enhancement - Load CSP configuration file if it exists: [#10139](https://github.com/owncloud/ocis/pull/10139)

   The Content Security Policy (CSP) configuration file is now loaded by default if
   it exists. The configuration file looked for should be located at
   `$OCIS_BASE_DATA_PATH/proxy/csp.yaml`. If the file does not exist, the default
   CSP configuration is used.

   https://github.com/owncloud/ocis/issues/10021
   https://github.com/owncloud/ocis/pull/10139

* Enhancement - FileDownloaded Activity: [#10161](https://github.com/owncloud/ocis/pull/10161)

   Add an activity when a file gets downloaded via public link

   https://github.com/owncloud/ocis/pull/10161

* Enhancement - Add WOPI host URLs to the collaboration service: [#10174](https://github.com/owncloud/ocis/pull/10174)

   We added the WOPI host urls to create a better integration with WOPI clients.
   This allows the WOPI apps to display links to our sharing and versions panel in
   the UI.

   https://github.com/owncloud/ocis/pull/10174

* Enhancement - Update web to v10.3.0: [#10177](https://github.com/owncloud/ocis/pull/10177)

   Tags: web

   We updated ownCloud Web to v10.3.0. Please refer to the changelog (linked) for
   details on the web release.

  * Bugfix [owncloud/web#11557](https://github.com/owncloud/web/pull/11557): OCM token clipboard copy
  * Bugfix [owncloud/web#11560](https://github.com/owncloud/web/pull/11560): OCM local instance check
  * Bugfix [owncloud/web#11583](https://github.com/owncloud/web/pull/11583): Thumbnails for GeoGebra slides not showing up
  * Bugfix [owncloud/web#11584](https://github.com/owncloud/web/pull/11584): Logout issues on token renewal failure
  * Bugfix [owncloud/web#11633](https://github.com/owncloud/web/pull/11633): App version downloads
  * Bugfix [owncloud/web#11642](https://github.com/owncloud/web/pull/11642): Wrong webdav URL in sidebar
  * Bugfix [owncloud/web#11643](https://github.com/owncloud/web/pull/11643): Renaming space in projects view files table does not work
  * Bugfix [owncloud/web#11653](https://github.com/owncloud/web/pull/11653): Hide share type switch for project spaces
  * Bugfix [owncloud/web#11658](https://github.com/owncloud/web/pull/11658): File name truncation
  * Enhancement [owncloud/web#11553](https://github.com/owncloud/web/pull/11553): Copy quick link action removal
  * Enhancement [owncloud/web#11553](https://github.com/owncloud/web/pull/11553): Internal link removal
  * Enhancement [owncloud/web#11558](https://github.com/owncloud/web/pull/11558): Add split confirm button to create link modal
  * Enhancement [owncloud/web#11561](https://github.com/owncloud/web/pull/11561): Add versions to the left sidebar bottom
  * Enhancement [owncloud/web#11574](https://github.com/owncloud/web/pull/11574): Accessibility improvements
  * Enhancement [owncloud/web#11580](https://github.com/owncloud/web/pull/11580): Show min oCIS version in app details (app store)
  * Enhancement [owncloud/web#11586](https://github.com/owncloud/web/pull/11586): Add a "Save As" function to the app top bar
  * Enhancement [owncloud/web#11606](https://github.com/owncloud/web/pull/11606): Move permanent link indicator
  * Enhancement [owncloud/web#11606](https://github.com/owncloud/web/pull/11606): Redesign sidebar link section in sharing panel
  * Enhancement [owncloud/web#11614](https://github.com/owncloud/web/pull/11614): Soothe right sidebar panel transitions
  * Enhancement [owncloud/web#11631](https://github.com/owncloud/web/pull/11631): Preview loading performance
  * Enhancement [owncloud/web#11644](https://github.com/owncloud/web/pull/11644): Add cancel button to unsaved changes dialog
  * Enhancement [owncloud/web#11646](https://github.com/owncloud/web/pull/11646): File type icon for .ggs files
  * Enhancement [owncloud/web#11661](https://github.com/owncloud/web/pull/11661): Remove link type "Uploader"

   https://github.com/owncloud/ocis/pull/10177
   https://github.com/owncloud/web/releases/tag/v10.3.0

* Enhancement - Bump reva to 2.25.0: [#10194](https://github.com/owncloud/ocis/pull/10194)

  *   Bugfix [cs3org/reva#4854](https://github.com/cs3org/reva/pull/4854): Added ShareUpdate activity
  *   Bugfix [cs3org/reva#4865](https://github.com/cs3org/reva/pull/4865): Better response codes for app new endpoint
  *   Bugfix [cs3org/reva#4858](https://github.com/cs3org/reva/pull/4858): Better response codes for app new endpoint
  *   Bugfix [cs3org/reva#4867](https://github.com/cs3org/reva/pull/4867): Fix remaining space calculation for S3 blobstore
  *   Bugfix [cs3org/reva#4852](https://github.com/cs3org/reva/pull/4852): Populate public link user correctly
  *   Bugfix [cs3org/reva#4859](https://github.com/cs3org/reva/pull/4859): Fixed the collaboration service registration
  *   Bugfix [cs3org/reva#4835](https://github.com/cs3org/reva/pull/4835): Fix sharejail stat id
  *   Bugfix [cs3org/reva#4856](https://github.com/cs3org/reva/pull/4856): Fix time conversion
  *   Bugfix [cs3org/reva#4851](https://github.com/cs3org/reva/pull/4851): Use gateway selector in sciencemesh
  *   Bugfix [cs3org/reva#4850](https://github.com/cs3org/reva/pull/4850): Write upload session info atomically
  *   Enhancement [cs3org/reva#4866](https://github.com/cs3org/reva/pull/4866): Unit test the json ocm invite manager
  *   Enhancement [cs3org/reva#4847](https://github.com/cs3org/reva/pull/4847): Add IsVersion to UploadReadyEvent
  *   Enhancement [cs3org/reva#4868](https://github.com/cs3org/reva/pull/4868): Improve metadata client errors
  *   Enhancement [cs3org/reva#4848](https://github.com/cs3org/reva/pull/4848): Add trashbin support to posixfs alongside other improvements

   https://github.com/owncloud/ocis/pull/10194
   https://github.com/owncloud/ocis/pull/10172
   https://github.com/owncloud/ocis/pull/10157
   https://github.com/owncloud/ocis/pull/9817

# Changelog for [5.0.8] (2024-09-30)

The following sections list the changes for 5.0.8.

[5.0.8]: https://github.com/owncloud/ocis/compare/v6.4.0...v5.0.8

## Summary

* Bugfix - Update reva to v2.19.8: [#10138](https://github.com/owncloud/ocis/pull/10138)

## Details

* Bugfix - Update reva to v2.19.8: [#10138](https://github.com/owncloud/ocis/pull/10138)

   We updated reva to v2.19.8

  *   Fix [cs3org/reva#4761](https://github.com/cs3org/reva/pull/4761): Quotes in dav Content-Disposition header
  *   Fix [cs3org/reva#4853](https://github.com/cs3org/reva/pull/4853): Write upload session info atomically
  *   Enh [cs3org/reva#4701](https://github.com/cs3org/reva/pull/4701): Extend service account permissions

   https://github.com/owncloud/ocis/pull/10138
   https://github.com/owncloud/ocis/pull/10103

# Changelog for [6.4.0] (2024-09-12)

The following sections list the changes for 6.4.0.

[6.4.0]: https://github.com/owncloud/ocis/compare/v5.0.7...v6.4.0

## Summary

* Bugfix - Set capability response `disable_self_password_change` correctly: [#9853](https://github.com/owncloud/ocis/pull/9853)
* Bugfix - Activity Translations: [#9856](https://github.com/owncloud/ocis/pull/9856)
* Bugfix - The user attributes `userType` and `memberOf` are readonly: [#9867](https://github.com/owncloud/ocis/pull/9867)
* Bugfix - Use key to get specific trash item: [#9879](https://github.com/owncloud/ocis/pull/9879)
* Bugfix - Fix response code when upload a file over locked: [#9894](https://github.com/owncloud/ocis/pull/9894)
* Bugfix - List OCM permissions as graph drive item permissions: [#9905](https://github.com/owncloud/ocis/pull/9905)
* Bugfix - Fix listing ocm shares: [#9925](https://github.com/owncloud/ocis/pull/9925)
* Bugfix - Allow update of ocm shares: [#9980](https://github.com/owncloud/ocis/pull/9980)
* Change - Remove store service: [#9890](https://github.com/owncloud/ocis/pull/9890)
* Enhancement - We now set the configured protocol transport for service metadata: [#9490](https://github.com/owncloud/ocis/pull/9490)
* Enhancement - Microsoft Office365 and Office Online support: [#9686](https://github.com/owncloud/ocis/pull/9686)
* Enhancement - Added a new role space editor without versions: [#9880](https://github.com/owncloud/ocis/pull/9880)
* Enhancement - Improve revisions purge: [#9891](https://github.com/owncloud/ocis/pull/9891)
* Enhancement - Allow setting default locale of activitylog: [#9892](https://github.com/owncloud/ocis/pull/9892)
* Enhancement - Graph translation path: [#9902](https://github.com/owncloud/ocis/pull/9902)
* Enhancement - Added a new roles viewer/editor with ListGrants: [#9943](https://github.com/owncloud/ocis/pull/9943)
* Enhancement - Handle OCM invite generated event: [#9966](https://github.com/owncloud/ocis/pull/9966)
* Enhancement - Update web to v10.2.0: [#9988](https://github.com/owncloud/ocis/pull/9988)
* Enhancement - Allow blob as connect-src in default CSP: [#9993](https://github.com/owncloud/ocis/pull/9993)
* Enhancement - Unified Roles Management: [#10013](https://github.com/owncloud/ocis/pull/10013)
* Enhancement - Bump reva to v2.24.1: [#10028](https://github.com/owncloud/ocis/pull/10028)

## Details

* Bugfix - Set capability response `disable_self_password_change` correctly: [#9853](https://github.com/owncloud/ocis/pull/9853)

   The capability value `disable_self_password_change` was not being set correctly
   when `user.passwordProfile` is configured as a read-only attribute.

   https://github.com/owncloud/enterprise/issues/6849
   https://github.com/owncloud/ocis/pull/9853

* Bugfix - Activity Translations: [#9856](https://github.com/owncloud/ocis/pull/9856)

   Translations for activities did not show up in transifex

   https://github.com/owncloud/ocis/pull/9856

* Bugfix - The user attributes `userType` and `memberOf` are readonly: [#9867](https://github.com/owncloud/ocis/pull/9867)

   The graph API now treats the user attributes `userType` and `memberOf` as
   read-only. They are not meant be updated directly by the client.

   https://github.com/owncloud/ocis/issues/9858
   https://github.com/owncloud/ocis/pull/9867

* Bugfix - Use key to get specific trash item: [#9879](https://github.com/owncloud/ocis/pull/9879)

   The activitylog and clientlog services now only fetch the specific trash item
   instead of getting all items in trash and filtering them on their side. This
   reduces the load on the storage users service because it no longer has to
   assemble a full trash listing.

   https://github.com/owncloud/ocis/pull/9879

* Bugfix - Fix response code when upload a file over locked: [#9894](https://github.com/owncloud/ocis/pull/9894)

   We fixed a bug where the response code was incorrect when uploading a file over
   a locked file.

   https://github.com/owncloud/ocis/issues/7638
   https://github.com/owncloud/ocis/pull/9894

* Bugfix - List OCM permissions as graph drive item permissions: [#9905](https://github.com/owncloud/ocis/pull/9905)

   The libre graph API now returns OCM shares when listing driveItem permissions.

   https://github.com/owncloud/ocis/issues/9898
   https://github.com/owncloud/ocis/pull/9905

* Bugfix - Fix listing ocm shares: [#9925](https://github.com/owncloud/ocis/pull/9925)

   The libre graph API now returns an etag, the role and the creation time for ocm
   shares. It also includes ocm shares in the sharedByMe endpoint.

   https://github.com/owncloud/ocis/pull/9925
   https://github.com/owncloud/ocis/pull/9920

* Bugfix - Allow update of ocm shares: [#9980](https://github.com/owncloud/ocis/pull/9980)

   We fixed a bug that prevented ocm shares to be updated or removed.

   https://github.com/owncloud/ocis/issues/9926
   https://github.com/owncloud/ocis/pull/9980

* Change - Remove store service: [#9890](https://github.com/owncloud/ocis/pull/9890)

   We have removed the unused store service.

   https://github.com/owncloud/ocis/issues/1357
   https://github.com/owncloud/ocis/pull/9890

* Enhancement - We now set the configured protocol transport for service metadata: [#9490](https://github.com/owncloud/ocis/pull/9490)

   This allows configuring services to listan on `tcp` or `unix` sockets and
   clients to use the `dns`, `kubernetes` or `unix` protocol URIs instead of
   service names.

   https://github.com/owncloud/ocis/pull/9490
   https://github.com/cs3org/reva/pull/4744

* Enhancement - Microsoft Office365 and Office Online support: [#9686](https://github.com/owncloud/ocis/pull/9686)

   Add support for Microsoft Office365 Cloud and Microsoft Office Online on
   premises. You can use the cloud feature either within a Microsoft
   [CSP](https://learn.microsoft.com/en-us/partner-center/enroll/csp-overview)
   partnership or via the ownCloud office365 proxy subscription. Please contact
   sales@owncloud.com to get more information about the ownCloud office365 proxy
   subscription.

   https://github.com/owncloud/ocis/pull/9686

* Enhancement - Added a new role space editor without versions: [#9880](https://github.com/owncloud/ocis/pull/9880)

   We add a new role space editor without list and restore version permissions.

   https://github.com/owncloud/ocis/issues/9699
   https://github.com/owncloud/ocis/pull/9880

* Enhancement - Improve revisions purge: [#9891](https://github.com/owncloud/ocis/pull/9891)

   The `revisions purge` command would time out on big spaces. We have improved
   performance by parallelizing the process.

   https://github.com/owncloud/ocis/pull/9891

* Enhancement - Allow setting default locale of activitylog: [#9892](https://github.com/owncloud/ocis/pull/9892)

   Allows setting the default locale via `OCIS_DEFAULT_LANGUAGE` envvar

   https://github.com/owncloud/ocis/pull/9892

* Enhancement - Graph translation path: [#9902](https://github.com/owncloud/ocis/pull/9902)

   Add `GRAPH_TRANSLATION_PATH` envvar like in other l10n services

   https://github.com/owncloud/ocis/pull/9902

* Enhancement - Added a new roles viewer/editor with ListGrants: [#9943](https://github.com/owncloud/ocis/pull/9943)

   We add a new roles space viewer/editor with ListGrants permissions.

   https://github.com/owncloud/ocis/issues/9701
   https://github.com/owncloud/ocis/pull/9943

* Enhancement - Handle OCM invite generated event: [#9966](https://github.com/owncloud/ocis/pull/9966)

   Both the notification and audit services now handle the OCM invite generated
   event.

   - The notification service is responsible for sending an email to the invited
   user. - The audit service is responsible for logging the event.

   https://github.com/owncloud/ocis/issues/9583
   https://github.com/owncloud/ocis/pull/9966
   https://github.com/cs3org/reva/pull/4832

* Enhancement - Update web to v10.2.0: [#9988](https://github.com/owncloud/ocis/pull/9988)

   Tags: web

   We updated ownCloud Web to v10.2.0. Please refer to the changelog (linked) for
   details on the web release.

  * Bugfix [owncloud/web#11512](https://github.com/owncloud/web/pull/11512): OCM invite generation body format
  * Bugfix [owncloud/web#11526](https://github.com/owncloud/web/pull/11526): Logout on access token renewal failure
  * Enhancement [owncloud/web#11377](https://github.com/owncloud/web/pull/11377): Replace custom datepicker with native html element
  * Enhancement [owncloud/web#11387](https://github.com/owncloud/web/pull/11387): Display disabled role permissions
  * Enhancement [owncloud/web#11394](https://github.com/owncloud/web/pull/11394): Mark external shares
  * Enhancement [owncloud/web#11484](https://github.com/owncloud/web/pull/11484): Hide versions panel with insufficient permissions
  * Enhancement [owncloud/web#11502](https://github.com/owncloud/web/pull/11502): Support a tags in actions
  * Enhancement [owncloud/web#11508](https://github.com/owncloud/web/pull/11508): Improve tiles view performance
  * Enhancement [owncloud/web#11515](https://github.com/owncloud/web/pull/11515): Add default actions extension point
  * Enhancement [owncloud/web#11518](https://github.com/owncloud/web/pull/11518): Add select all checkbox to tiles view

   https://github.com/owncloud/ocis/pull/9988
   https://github.com/owncloud/web/releases/tag/v10.2.0

* Enhancement - Allow blob as connect-src in default CSP: [#9993](https://github.com/owncloud/ocis/pull/9993)

   We added 'blob:' to the default connect-src items in the default CSP rules.

   https://github.com/owncloud/ocis/pull/9993

* Enhancement - Unified Roles Management: [#10013](https://github.com/owncloud/ocis/pull/10013)

   Improved management of unified roles with the introduction of default
   enabled/disabled states and a new command for listing available roles. It is
   important to note that a disabled role does not lose previously assigned
   permissions; it only means that the role is not available for new assignments.

   The following roles are now enabled by default:

   - UnifiedRoleViewerID - UnifiedRoleSpaceViewer - UnifiedRoleEditor -
   UnifiedRoleSpaceEditor - UnifiedRoleFileEditor - UnifiedRoleEditorLite -
   UnifiedRoleManager

   The following roles are now disabled by default:

   - UnifiedRoleSecureViewer

   To enable the UnifiedRoleSecureViewer role, you must provide a list of all
   available roles through one of the following methods:

   - Using the GRAPH_AVAILABLE_ROLES environment variable. - Setting the
   available_roles configuration value.

   To enable a role, include the UID of the role in the list of available roles.

   A new command has been introduced to simplify the process of finding out which
   UID belongs to which role. The command is:

   ```
   $ ocis graph list-unified-roles
   ```

   The output of this command includes the following information for each role:

   - uid: The unique identifier of the role. - Description: A short description of
   the role. - Enabled: Whether the role is enabled or not.

   https://github.com/owncloud/ocis/issues/9698
   https://github.com/owncloud/ocis/pull/10013
   https://github.com/owncloud/ocis/pull/9727

* Enhancement - Bump reva to v2.24.1: [#10028](https://github.com/owncloud/ocis/pull/10028)

  *   Bugfix [cs3org/reva#4843](https://github.com/cs3org/reva/pull/4843): Allow update of ocm shares
  *   Bugfix [cs3org/reva#4820](https://github.com/cs3org/reva/pull/4820): Fix response code when upload a file over locked
  *   Bugfix [cs3org/reva#4837](https://github.com/cs3org/reva/pull/4837): Fix OCM userid encoding
  *   Bugfix [cs3org/reva#4823](https://github.com/cs3org/reva/pull/4823): Return etag for ocm shares
  *   Bugfix [cs3org/reva#4822](https://github.com/cs3org/reva/pull/4822): Allow listing directory trash items by key
  *   Enhancement [cs3org/reva#4816](https://github.com/cs3org/reva/pull/4816): Ignore resharing requests
  *   Enhancement [cs3org/reva#4817](https://github.com/cs3org/reva/pull/4817): Added a new role space editor without versions
  *   Enhancement [cs3org/reva#4829](https://github.com/cs3org/reva/pull/4829): Added a new roles viewer/editor with ListGrants
  *   Enhancement [cs3org/reva#4828](https://github.com/cs3org/reva/pull/4828): New event: UserSignedIn
  *   Enhancement [cs3org/reva#4836](https://github.com/cs3org/reva/pull/4836): Publish an event when an OCM invite is generated

   https://github.com/owncloud/ocis/pull/10028
   https://github.com/owncloud/ocis/pull/9980
   https://github.com/owncloud/ocis/pull/9981
   https://github.com/owncloud/ocis/pull/9981
   https://github.com/owncloud/ocis/pull/9920
   https://github.com/owncloud/ocis/pull/9879
   https://github.com/owncloud/ocis/pull/9860

# Changelog for [5.0.7] (2024-09-04)

The following sections list the changes for 5.0.7.

[5.0.7]: https://github.com/owncloud/ocis/compare/v6.3.0...v5.0.7

## Summary

* Enhancement - Add virus filter to sessions command: [#9041](https://github.com/owncloud/ocis/pull/9041)
* Enhancement - Assimilate `clean` into `sessions` command: [#9828](https://github.com/owncloud/ocis/pull/9828)
* Enhancement - Update web to v8.0.5: [#9958](https://github.com/owncloud/ocis/pull/9958)

## Details

* Enhancement - Add virus filter to sessions command: [#9041](https://github.com/owncloud/ocis/pull/9041)

   Allow filtering upload session by virus status (has-virus=true/false)

   https://github.com/owncloud/ocis/pull/9041

* Enhancement - Assimilate `clean` into `sessions` command: [#9828](https://github.com/owncloud/ocis/pull/9828)

   We deprecated `ocis storage-user uploads clean` and added the same logic to
   `ocis storage-users uploads session --clean`

   https://github.com/owncloud/ocis/pull/9828

* Enhancement - Update web to v8.0.5: [#9958](https://github.com/owncloud/ocis/pull/9958)

   Tags: web

   We updated ownCloud Web to v8.0.5. Please refer to the changelog (linked) for
   details on the web release.

   - Bugfix [owncloud/web#11395](https://github.com/owncloud/web/pull/11395):
   Missing space members for group memberships - Bugfix
   [owncloud/web#11263](https://github.com/owncloud/web/pull/11263): Show more
   toggle in space members view not reactive - Bugfix
   [owncloud/web#11263](https://github.com/owncloud/web/pull/11263): Space show
   links from other spaces - Bugfix
   [owncloud/web#11303](https://github.com/owncloud/web/pull/11303): Uploading
   nested folders

   https://github.com/owncloud/ocis/pull/9958
   https://github.com/owncloud/web/releases/tag/v8.0.5

# Changelog for [6.3.0] (2024-08-20)

The following sections list the changes for 6.3.0.

[6.3.0]: https://github.com/owncloud/ocis/compare/v6.2.0...v6.3.0

## Summary

* Bugfix - Ignore address for kubernetes registry: [#9490](https://github.com/owncloud/ocis/pull/9490)
* Bugfix - Use bool type for web embed delegatedAuthentication: [#9692](https://github.com/owncloud/ocis/pull/9692)
* Bugfix - Repair nats-js-kv registry: [#9734](https://github.com/owncloud/ocis/pull/9734)
* Bugfix - Use less selectors that watch the registry: [#9741](https://github.com/owncloud/ocis/pull/9741)
* Bugfix - We fixed the client config generation for the built in IDP: [#9770](https://github.com/owncloud/ocis/pull/9770)
* Bugfix - Change ocmproviders config defaultpath: [#9778](https://github.com/owncloud/ocis/pull/9778)
* Bugfix - Web theme color contrasts: [#10726](https://github.com/owncloud/web/issues/10726)
* Enhancement - New WOPI operations added to the collaboration service: [#9505](https://github.com/owncloud/ocis/pull/9505)
* Enhancement - Allow configuring grpc max connection age: [#9657](https://github.com/owncloud/ocis/pull/9657)
* Enhancement - Tracing improvements in the collaboration service: [#9684](https://github.com/owncloud/ocis/pull/9684)
* Enhancement - Local WEB App configuration: [#9691](https://github.com/owncloud/ocis/pull/9691)
* Enhancement - Bump tusd pkg to v2: [#9714](https://github.com/owncloud/ocis/pull/9714)
* Enhancement - Gateways should directly talk to themselves: [#9714](https://github.com/owncloud/ocis/pull/9714)
* Enhancement - Support Skyhigh Security ICAP as an ICAP server: [#9720](https://github.com/owncloud/ocis/issues/9720)
* Enhancement - Added generic way to translate composite entities: [#9722](https://github.com/owncloud/ocis/pull/9722)
* Enhancement - Add an API to auth-app service: [#9755](https://github.com/owncloud/ocis/pull/9755)
* Enhancement - Bump go-micro plugins pkg: [#9756](https://github.com/owncloud/ocis/pull/9756)
* Enhancement - Allow querying federated user roles for sharing: [#9765](https://github.com/owncloud/ocis/pull/9765)
* Enhancement - Refactor the connector in the collaboration service: [#9771](https://github.com/owncloud/ocis/pull/9771)
* Enhancement - Add OCIS_ENABLE_OCM env var: [#9784](https://github.com/owncloud/ocis/pull/9784)
* Enhancement - OCM related adjustments in graph: [#9788](https://github.com/owncloud/ocis/pull/9788)
* Enhancement - Update web to v10.1.0: [#9832](https://github.com/owncloud/ocis/pull/9832)
* Enhancement - Bump reva to 2.23.0: [#9852](https://github.com/owncloud/ocis/pull/9852)

## Details

* Bugfix - Ignore address for kubernetes registry: [#9490](https://github.com/owncloud/ocis/pull/9490)

   We no longer pass an address to the go micro kubernetes registry implementation.
   This causes the implementation to autodetect the namespace and not hardcode it
   to `default`.

   https://github.com/owncloud/ocis/pull/9490

* Bugfix - Use bool type for web embed delegatedAuthentication: [#9692](https://github.com/owncloud/ocis/pull/9692)

   https://github.com/owncloud/ocis/pull/9692

* Bugfix - Repair nats-js-kv registry: [#9734](https://github.com/owncloud/ocis/pull/9734)

   The registry would always send traffic to only one pod. This is now fixed and
   load should be spread evenly. Also implements watcher method so the cache can
   use it. Internally, it can now distinguish services by version and will
   aggregate all nodes of the same version into a single service, as expected by
   the registry cache and watcher.

   https://github.com/owncloud/ocis/pull/9734
   https://github.com/owncloud/ocis/pull/9726
   https://github.com/owncloud/ocis/pull/9656

* Bugfix - Use less selectors that watch the registry: [#9741](https://github.com/owncloud/ocis/pull/9741)

   The proxy now shares the service selector for all host lookups.

   https://github.com/owncloud/ocis/pull/9741

* Bugfix - We fixed the client config generation for the built in IDP: [#9770](https://github.com/owncloud/ocis/pull/9770)

   We now use the OCIS_URL to generate the web client registration configuration.
   It does not make sense use the OCIS_ISSUER_URL if the idp was configured to run
   on a different domain.

   https://github.com/owncloud/ocis/pull/9770

* Bugfix - Change ocmproviders config defaultpath: [#9778](https://github.com/owncloud/ocis/pull/9778)

   We moved the default location of the `ocmproviders.json` config file out of the
   data directory of the ocm service to the ocis config directory.

   https://github.com/owncloud/ocis/pull/9778

* Bugfix - Web theme color contrasts: [#10726](https://github.com/owncloud/web/issues/10726)

   Web theme colors have been enhanced so they have at least a 4.5:1 contrast ratio
   because of a11y reasons.

   https://github.com/owncloud/web/issues/10726
   https://github.com/owncloud/web/pull/11331
   https://github.com/owncloud/ocis/pull/9752

* Enhancement - New WOPI operations added to the collaboration service: [#9505](https://github.com/owncloud/ocis/pull/9505)

   PutRelativeFile, DeleteFile and RenameFile operations have been added to the
   collaboration service. GetFileInfo operation will now report the support of
   these operations to the WOPI app

   https://github.com/owncloud/ocis/pull/9505

* Enhancement - Allow configuring grpc max connection age: [#9657](https://github.com/owncloud/ocis/pull/9657)

   We added a GRPC_MAX_CONNECTION_AGE env var that allows limiting the lifespan of
   connections. A closed connection triggers grpc clients to do a new DNS lookup to
   pick up new IPs.

   https://github.com/owncloud/ocis/pull/9657

* Enhancement - Tracing improvements in the collaboration service: [#9684](https://github.com/owncloud/ocis/pull/9684)

   Uploads and downloads through the collaboration service will be traced. The
   openInApp request will also be linked properly with other requests in the
   tracing. In addition, the collaboration service will include some additional
   information in the traces. Filtering based on that information might be an
   option.

   https://github.com/owncloud/ocis/pull/9684

* Enhancement - Local WEB App configuration: [#9691](https://github.com/owncloud/ocis/pull/9691)

   We've added a new feature which allows configuring applications individually
   instead of using the global apps.yaml file. With that, each application can have
   its own configuration file, which will be loaded by the WEB service.

   The local configuration has the highest priority and will override the global
   configuration. The Following order of precedence is used: local.config >
   global.config > manifest.config.

   Besides the configuration, the application now be disabled by setting the
   `disabled` field to `true` in one of the configuration files.

   https://github.com/owncloud/ocis/issues/9687
   https://github.com/owncloud/ocis/pull/9691

* Enhancement - Bump tusd pkg to v2: [#9714](https://github.com/owncloud/ocis/pull/9714)

   Bumps the tusd pkg to v2.4.0

   https://github.com/owncloud/ocis/pull/9714

* Enhancement - Gateways should directly talk to themselves: [#9714](https://github.com/owncloud/ocis/pull/9714)

   The CS3 gateway can directly to itself when it wants to talk to the registries
   running in the same reva runtime.

   https://github.com/owncloud/ocis/pull/9714

* Enhancement - Support Skyhigh Security ICAP as an ICAP server: [#9720](https://github.com/owncloud/ocis/issues/9720)

   We have upgraded the antivirus ICAP client library, bringing enhanced
   performance and reliability to our antivirus scanning service. With this update,
   the Skyhigh Security ICAP can now be used as an ICAP server, providing robust
   and scalable antivirus solutions.

   https://github.com/owncloud/ocis/issues/9720
   https://github.com/fschade/icap-client/pull/6

* Enhancement - Added generic way to translate composite entities: [#9722](https://github.com/owncloud/ocis/pull/9722)

   Added a generic way to translate the necessary fields in composite entities. The
   function takes the entity, translation function and fields to translate that are
   described by the TranslateField function. The function supports nested structs
   and slices of structs.

   https://github.com/owncloud/ocis/issues/9700
   https://github.com/owncloud/ocis/pull/9722

* Enhancement - Add an API to auth-app service: [#9755](https://github.com/owncloud/ocis/pull/9755)

   Adds an API to create, list and delete app tokens. Includes an impersonification
   feature for migration scenarios.

   https://github.com/owncloud/ocis/pull/9755

* Enhancement - Bump go-micro plugins pkg: [#9756](https://github.com/owncloud/ocis/pull/9756)

   Bump plugins pkg to include fix for cache delete

   https://github.com/owncloud/ocis/pull/9756

* Enhancement - Allow querying federated user roles for sharing: [#9765](https://github.com/owncloud/ocis/pull/9765)

   When listing permissions clients can now fetch the list of available federated
   sharing roles by sending a `GET
   /graph/v1beta1/drives/{driveid}/items/{itemid}/permissions?$filter=@libre.graph.permissions.roles.allowedValues/rolePermissions/any(p:contains(p/condition,
   '@Subject.UserType=="Federated"'))` request. Note that this is the only
   supported filter expression. Federated sharing roles will be omitted from
   requests without this filter.

   https://github.com/owncloud/ocis/pull/9765

* Enhancement - Refactor the connector in the collaboration service: [#9771](https://github.com/owncloud/ocis/pull/9771)

   This will simplify and homogenize the code around the connector

   https://github.com/owncloud/ocis/pull/9771

* Enhancement - Add OCIS_ENABLE_OCM env var: [#9784](https://github.com/owncloud/ocis/pull/9784)

   We added a new `OCIS_ENABLE_OCM` env var that will enable all ocm flags.

   https://github.com/owncloud/ocis/pull/9784

* Enhancement - OCM related adjustments in graph: [#9788](https://github.com/owncloud/ocis/pull/9788)

   The /users enpdoint of the graph service was changed with respect to how it
   handles OCM federeated users: - The 'userType' property is now alway returned.
   As new usertype 'Federated' was introduced. To indicate that the user is a
   federated user. - Supported for filtering users by 'userType' as added. Queries
   like "$filter=userType eq 'Federated'" are now possible. - Federated users are
   only returned when explicitly requested via filter. When no filter is provider
   only 'Member' users are returned.

   https://github.com/owncloud/ocis/issues/9702
   https://github.com/owncloud/ocis/pull/9788
   https://github.com/owncloud/ocis/pull/9757

* Enhancement - Update web to v10.1.0: [#9832](https://github.com/owncloud/ocis/pull/9832)

   Tags: web

   We updated ownCloud Web to v10.1.0. Please refer to the changelog (linked) for
   details on the web release.

   - Bugfix [owncloud/web#11263](https://github.com/owncloud/web/pull/11263) Show
   more toggle in space members view not reactive - Bugfix
   [owncloud/web#11299](https://github.com/owncloud/web/pull/11299) Uploading
   nested folders - Bugfix
   [owncloud/web#11312](https://github.com/owncloud/web/pull/11312) Toggling
   checkboxes via keyboard - Bugfix
   [owncloud/web#11313](https://github.com/owncloud/web/pull/11313) Prevent
   horizontal table scroll - Bugfix
   [owncloud/web#11342](https://github.com/owncloud/web/pull/11342) Keyboard
   actions for disabled resources - Bugfix
   [owncloud/web#11348](https://github.com/owncloud/web/pull/11348) OCM page reload
   - Bugfix [owncloud/web#11353](https://github.com/owncloud/web/pull/11353)
   Closing an app opened via in-app open feature stays open - Enhancement
   [owncloud/web#11287](https://github.com/owncloud/web/pull/11287) Add quota
   information to account page - Enhancement
   [owncloud/web#11302](https://github.com/owncloud/web/pull/11302) App Store app -
   Enhancement [owncloud/web#11310](https://github.com/owncloud/web/pull/11310)
   Redesign share link modal - Enhancement
   [owncloud/web#11315](https://github.com/owncloud/web/pull/11315) Accessibility -
   Enhancement [owncloud/web#11329](https://github.com/owncloud/web/pull/11329)
   Files as links - Enhancement
   [owncloud/web#11344](https://github.com/owncloud/web/pull/11344) Unstick top bar

   https://github.com/owncloud/ocis/pull/9832
   https://github.com/owncloud/web/releases/tag/v10.1.0

* Enhancement - Bump reva to 2.23.0: [#9852](https://github.com/owncloud/ocis/pull/9852)

  *   Fix [cs3org/reva#4741](https://github.com/cs3org/reva/pull/4741): Always find unique providers
  *   Fix [cs3org/reva#4741](https://github.com/cs3org/reva/pull/4762): Blanks in dav Content-Disposition header
  *   Fix [cs3org/reva#4741](https://github.com/cs3org/reva/pull/4775): Fixed the response code when copying the shared from to personal
  *   Fix [cs3org/reva#4741](https://github.com/cs3org/reva/pull/4633): Allow all users to create internal links
  *   Fix [cs3org/reva#4741](https://github.com/cs3org/reva/pull/4771): Deleting resources via their id
  *   Fix [cs3org/reva#4741](https://github.com/cs3org/reva/pull/4768): Fixed the file name validation if nodeid is used
  *   Fix [cs3org/reva#4741](https://github.com/cs3org/reva/pull/4758): Fix moving locked files, enable handling locked files via ocdav
  *   Fix [cs3org/reva#4741](https://github.com/cs3org/reva/pull/4774): Fix micro ocdav service init and registration
  *   Fix [cs3org/reva#4741](https://github.com/cs3org/reva/pull/4776): Fix response code for DEL file that in postprocessing
  *   Fix [cs3org/reva#4741](https://github.com/cs3org/reva/pull/4746): Uploading the same file multiple times leads to orphaned blobs
  *   Fix [cs3org/reva#4741](https://github.com/cs3org/reva/pull/4778): Zero byte uploads
  *   Chg [cs3org/reva#4741](https://github.com/cs3org/reva/pull/4759): Updated to the latest version of the go-cs3apis
  *   Chg [cs3org/reva#4741](https://github.com/cs3org/reva/pull/4773): Ocis bumped
  *   Enh [cs3org/reva#4741](https://github.com/cs3org/reva/pull/4766): Set archiver output format via query parameter
  *   Enh [cs3org/reva#4741](https://github.com/cs3org/reva/pull/4763): Improve posixfs storage driver

   https://github.com/owncloud/ocis/pull/9852
   https://github.com/owncloud/ocis/pull/9763
   https://github.com/owncloud/ocis/pull/9714
   https://github.com/owncloud/ocis/pull/9715

# Changelog for [6.2.0] (2024-07-30)

The following sections list the changes for 6.2.0.

[6.2.0]: https://github.com/owncloud/ocis/compare/v5.0.6...v6.2.0

## Summary

* Bugfix - Fix restarting of postprocessing: [#6945](https://github.com/owncloud/ocis/pull/6945)
* Bugfix - Fix crash on empty tracing provider: [#9622](https://github.com/owncloud/ocis/pull/9622)
* Bugfix - Fixed the file name validation if nodeid is used: [#9634](https://github.com/owncloud/ocis/pull/9634)
* Bugfix - Fix a missing SecureView permission attribute in the REPORT response: [#9638](https://github.com/owncloud/ocis/pull/9638)
* Bugfix - Fixed the channel lock in a workers pool: [#9647](https://github.com/owncloud/ocis/pull/9647)
* Bugfix - Missing invitation in permission responses: [#9652](https://github.com/owncloud/ocis/pull/9652)
* Bugfix - Repair nats-js-kv registry: [#9662](https://github.com/owncloud/ocis/pull/9662)
* Bugfix - Fix panic: [#9673](https://github.com/owncloud/ocis/pull/9673)
* Bugfix - Fixed the response code when copying the shared from to personal: [#9677](https://github.com/owncloud/ocis/pull/9677)
* Bugfix - Fixed response code for DELETE file that is in postprocessing: [#9689](https://github.com/owncloud/ocis/pull/9689)
* Change - Remove unavailable web config options: [#9679](https://github.com/owncloud/ocis/pull/9679)
* Enhancement - Introduce auth-app service: [#9079](https://github.com/owncloud/ocis/pull/9079)
* Enhancement - Add support for proof keys for the collaboration service: [#9366](https://github.com/owncloud/ocis/pull/9366)
* Enhancement - Log user agent and remote addr on auth errors: [#9475](https://github.com/owncloud/ocis/pull/9475)
* Enhancement - Add missing WOPI features: [#9580](https://github.com/owncloud/ocis/pull/9580)
* Enhancement - Bump commitID for web: [#9631](https://github.com/owncloud/ocis/pull/9631)
* Enhancement - Remove oidc-go dependency: [#9641](https://github.com/owncloud/ocis/pull/9641)
* Enhancement - Improve the collaboration service logging: [#9653](https://github.com/owncloud/ocis/pull/9653)
* Enhancement - Fix trash command: [#9665](https://github.com/owncloud/ocis/pull/9665)
* Enhancement - Added the debugging to full ocis docker example: [#9666](https://github.com/owncloud/ocis/pull/9666)
* Enhancement - Add locking support for MS Office Online Server: [#9685](https://github.com/owncloud/ocis/pull/9685)
* Enhancement - Bump reva to v.2.22.0: [#9690](https://github.com/owncloud/ocis/pull/9690)
* Enhancement - Add `--diff` to the `ocis init` command: [#9693](https://github.com/owncloud/ocis/pull/9693)
* Enhancement - Update web to v10.0.0: [#9707](https://github.com/owncloud/ocis/pull/9707)

## Details

* Bugfix - Fix restarting of postprocessing: [#6945](https://github.com/owncloud/ocis/pull/6945)

   We fixed a bug where non-admin requests to admin resources would get 401
   Unauthorized. Now, the server sends 403 Forbidden response.

   https://github.com/owncloud/ocis/issues/5938
   https://github.com/owncloud/ocis/pull/6945

* Bugfix - Fix crash on empty tracing provider: [#9622](https://github.com/owncloud/ocis/pull/9622)

   We have fixed a bug that causes a crash when OCIS_TRACING_ENABLED is set to
   true, but no tracing Endpoints or Collectors have been provided.a

   https://github.com/owncloud/ocis/issues/7012
   https://github.com/owncloud/ocis/pull/9622

* Bugfix - Fixed the file name validation if nodeid is used: [#9634](https://github.com/owncloud/ocis/pull/9634)

   We have fixed the file name validation if nodeid is used

   https://github.com/owncloud/ocis/issues/9568
   https://github.com/owncloud/ocis/pull/9634

* Bugfix - Fix a missing SecureView permission attribute in the REPORT response: [#9638](https://github.com/owncloud/ocis/pull/9638)

   We fixed a missing SecureView permission attribute in the REPORT response.

   https://github.com/owncloud/ocis/issues/9607
   https://github.com/owncloud/ocis/pull/9638

* Bugfix - Fixed the channel lock in a workers pool: [#9647](https://github.com/owncloud/ocis/pull/9647)

   We fixed an error when the users can't see more than 10 shares

   https://github.com/owncloud/ocis/issues/9642
   https://github.com/owncloud/ocis/pull/9647

* Bugfix - Missing invitation in permission responses: [#9652](https://github.com/owncloud/ocis/pull/9652)

   We have fixed a bug where the `invitation` property was missing in the response
   when creating, listing or updating graph permissions.

   https://github.com/owncloud/ocis/issues/9571
   https://github.com/owncloud/ocis/pull/9652

* Bugfix - Repair nats-js-kv registry: [#9662](https://github.com/owncloud/ocis/pull/9662)

   The registry would always send traffic to only one pod. This is now fixed and
   load should be spread evenly. Also implements watcher method so the cache can
   use it.

   https://github.com/owncloud/ocis/pull/9662
   https://github.com/owncloud/ocis/pull/9654
   https://github.com/owncloud/ocis/pull/9620

* Bugfix - Fix panic: [#9673](https://github.com/owncloud/ocis/pull/9673)

   Fixes panic occuring when the nats-js-kv is not properly initialized.

   https://github.com/owncloud/ocis/pull/9673

* Bugfix - Fixed the response code when copying the shared from to personal: [#9677](https://github.com/owncloud/ocis/pull/9677)

   We fixed the response code when copying the file from shares to personal space
   with a secure view role.

   https://github.com/owncloud/ocis/issues/9482
   https://github.com/owncloud/ocis/pull/9677

* Bugfix - Fixed response code for DELETE file that is in postprocessing: [#9689](https://github.com/owncloud/ocis/pull/9689)

   We fixed the response code when DELETE and MOVE requests to the file that is
   still in post-processing.

   https://github.com/owncloud/ocis/issues/9432
   https://github.com/owncloud/ocis/pull/9689

* Change - Remove unavailable web config options: [#9679](https://github.com/owncloud/ocis/pull/9679)

   We've removed config options from the web package, that are no longer available
   in web. Check the web changelog for more details.

   https://github.com/owncloud/ocis/pull/9679
   https://github.com/owncloud/web/pull/11256
   https://github.com/owncloud/web/pull/10122
   https://github.com/owncloud/web/pull/11260

* Enhancement - Introduce auth-app service: [#9079](https://github.com/owncloud/ocis/pull/9079)

   Introduce a new service, auth-app, that provides authentication and
   authorization services for applications.

   https://github.com/owncloud/ocis/pull/9079

* Enhancement - Add support for proof keys for the collaboration service: [#9366](https://github.com/owncloud/ocis/pull/9366)

   Proof keys support will be enabled by default in order to ensure that all the
   requests come from a trusted source. Since proof keys must be set in the WOPI
   app (OnlyOffice, Collabora...), it's possible to disable the verification of the
   proof keys via configuration.

   https://github.com/owncloud/ocis/pull/9366

* Enhancement - Log user agent and remote addr on auth errors: [#9475](https://github.com/owncloud/ocis/pull/9475)

   The proxy will now log `user_agent`, `client.address`, `network.peer.address`
   and `network.peer.port` to help operations debug authentication errors. The
   latter three follow the [Semantic Conventions 1.26.0 / General /
   Attributes](https://opentelemetry.io/docs/specs/semconv/general/attributes/)
   naming to better integrate with log aggregation tools.

   https://github.com/owncloud/ocis/pull/9475

* Enhancement - Add missing WOPI features: [#9580](https://github.com/owncloud/ocis/pull/9580)

   We added the feature to disable the chat for onlyoffice and added the missing
   language parameters to the wopi app url.

   https://github.com/owncloud/ocis/pull/9580

* Enhancement - Bump commitID for web: [#9631](https://github.com/owncloud/ocis/pull/9631)

   Bump the web commitID to current master

   https://github.com/owncloud/ocis/pull/9631

* Enhancement - Remove oidc-go dependency: [#9641](https://github.com/owncloud/ocis/pull/9641)

   Removes the kgol/oidc-go dependency because it was flagged by dependabot.
   Luckily us we only used it for importing the strings "profile" and "email".

   https://github.com/owncloud/ocis/pull/9641

* Enhancement - Improve the collaboration service logging: [#9653](https://github.com/owncloud/ocis/pull/9653)

   We added more debug log information to the collaboration service. This is vital
   for scenarios when we need to debug in remote setups.

   https://github.com/owncloud/ocis/pull/9653

* Enhancement - Fix trash command: [#9665](https://github.com/owncloud/ocis/pull/9665)

   The `ocis trash purge-empty-dirs` command should work on any storage provider,
   not just `storage/users`.

   https://github.com/owncloud/ocis/pull/9665

* Enhancement - Added the debugging to full ocis docker example: [#9666](https://github.com/owncloud/ocis/pull/9666)

   Added the debugging to full ocis docker example

   https://github.com/owncloud/ocis/pull/9666

* Enhancement - Add locking support for MS Office Online Server: [#9685](https://github.com/owncloud/ocis/pull/9685)

   We added support for the special kind of lock tokens that MS Office Online
   Server uses to lock files via the Wopi protocol. It will only be active if you
   set the `COLLABORATION_APP_NAME` environment variable to
   `MicrosoftOfficeOnline`.

   https://github.com/owncloud/ocis/pull/9685

* Enhancement - Bump reva to v.2.22.0: [#9690](https://github.com/owncloud/ocis/pull/9690)

  *   Bugfix [cs3org/reva#4741](https://github.com/cs3org/reva/pull/4741): Always find unique providers
  *   Bugfix [cs3org/reva#4762](https://github.com/cs3org/reva/pull/4762): Blanks in dav Content-Disposition header
  *   Bugfix [cs3org/reva#4775](https://github.com/cs3org/reva/pull/4775): Fixed the response code when copying the shared from to personal
  *   Bugfix [cs3org/reva#4633](https://github.com/cs3org/reva/pull/4633): Allow all users to create internal links
  *   Bugfix [cs3org/reva#4771](https://github.com/cs3org/reva/pull/4771): Deleting resources via their id
  *   Bugfix [cs3org/reva#4768](https://github.com/cs3org/reva/pull/4768): Fixed the file name validation if nodeid is used
  *   Bugfix [cs3org/reva#4758](https://github.com/cs3org/reva/pull/4758): Fix moving locked files, enable handling locked files via ocdav
  *   Bugfix [cs3org/reva#4774](https://github.com/cs3org/reva/pull/4774): Fix micro ocdav service init and registration
  *   Bugfix [cs3org/reva#4776](https://github.com/cs3org/reva/pull/4776): Fix response code for DEL file that in postprocessing
  *   Bugfix [cs3org/reva#4746](https://github.com/cs3org/reva/pull/4746): Uploading the same file multiple times leads to orphaned blobs
  *   Bugfix [cs3org/reva#4778](https://github.com/cs3org/reva/pull/4778): Zero byte uploads
  *   Change [cs3org/reva#4759](https://github.com/cs3org/reva/pull/4759): Updated to the latest version of the go-cs3apis
  *   Change [cs3org/reva#4773](https://github.com/cs3org/reva/pull/4773): Ocis bumped
  *   Enhancement [cs3org/reva#4766](https://github.com/cs3org/reva/pull/4766): Set archiver output format via query parameter
  *   Enhancement [cs3org/reva#4763](https://github.com/cs3org/reva/pull/4763): Improve posixfs storage driver

   https://github.com/owncloud/ocis/pull/9690
   https://github.com/owncloud/ocis/pull/9662
   https://github.com/owncloud/ocis/pull/9621
   https://github.com/owncloud/ocis/pull/9677
   https://github.com/owncloud/ocis/pull/9689

* Enhancement - Add `--diff` to the `ocis init` command: [#9693](https://github.com/owncloud/ocis/pull/9693)

   We have added a new flag `--diff` to the `ocis init` command to show the diff of
   the configuration files. This is useful to see what has changed in the
   configuration files when you run the `ocis init` command. The diff is stored to
   the ocispath in the config folder as ocis.config.patch and can be applied using
   the linux `patch` command.

   https://github.com/owncloud/ocis/issues/3645
   https://github.com/owncloud/ocis/pull/9693

* Enhancement - Update web to v10.0.0: [#9707](https://github.com/owncloud/ocis/pull/9707)

   Tags: web

   We updated ownCloud Web to v10.0.0. Please refer to the changelog (linked) for
   details on the web release.

   - Bugfix [owncloud/web#11174](https://github.com/owncloud/web/pull/11174)
   Downloading files via the app top bar doesn't reflect the current state - Bugfix
   [owncloud/web#11186](https://github.com/owncloud/web/pull/11186) Images
   stretched in preview app in Safari browser - Bugfix
   [owncloud/web#11194](https://github.com/owncloud/web/pull/11194) UI breaks when
   tags are numbers - Bugfix
   [owncloud/web#11253](https://github.com/owncloud/web/pull/11253) Open dropdown
   menu does not deselect other items in admin settings app - Change
   [owncloud/web#11251](https://github.com/owncloud/web/pull/11251) Removal of
   Deprecated Config Options - Change
   [owncloud/web#11252](https://github.com/owncloud/web/pull/11252) Remove draw-io
   as default app - Change
   [owncloud/web#11277](https://github.com/owncloud/web/pull/11277) Remove set as
   description space action - Enhancement
   [owncloud/web#11166](https://github.com/owncloud/web/pull/11166) Add share role
   icon to shared with me table - Enhancement
   [owncloud/web#11258](https://github.com/owncloud/web/pull/11258) Application
   menu extension point - Enhancement
   [owncloud/web#11279](https://github.com/owncloud/web/pull/11279) Move quota info
   to general info in user menu - Enhancement
   [owncloud/web#11280](https://github.com/owncloud/web/pull/11280) Add edit
   description button to space info

   https://github.com/owncloud/ocis/pull/9707
   https://github.com/owncloud/web/releases/tag/v10.0.0

# Changelog for [5.0.6] (2024-07-17)

The following sections list the changes for 5.0.6.

[5.0.6]: https://github.com/owncloud/ocis/compare/v6.1.0...v5.0.6

## Summary

* Bugfix - Allow all uploads to restart: [#9506](https://github.com/owncloud/ocis/pull/9506)
* Bugfix - Fix the email notification service: [#9514](https://github.com/owncloud/ocis/pull/9514)
* Enhancement - Limit concurrent thumbnail requests: [#9199](https://github.com/owncloud/ocis/pull/9199)
* Enhancement - Update web to v8.0.4: [#9429](https://github.com/owncloud/ocis/pull/9429)
* Enhancement - Add cli to purge revisions: [#9497](https://github.com/owncloud/ocis/pull/9497)

## Details

* Bugfix - Allow all uploads to restart: [#9506](https://github.com/owncloud/ocis/pull/9506)

   On postprocessing-restart, allow all uploads to restart even if one fails.

   https://github.com/owncloud/ocis/pull/9506

* Bugfix - Fix the email notification service: [#9514](https://github.com/owncloud/ocis/pull/9514)

   We fixed an error in the notification service that caused the email notification
   to fail when the user's display name contained special characters.

   https://github.com/owncloud/ocis/issues/9402
   https://github.com/owncloud/ocis/pull/9514

* Enhancement - Limit concurrent thumbnail requests: [#9199](https://github.com/owncloud/ocis/pull/9199)

   The number of concurrent requests to the thumbnail service can be limited now to
   have more control over the consumed system resources.

   https://github.com/owncloud/ocis/pull/9199

* Enhancement - Update web to v8.0.4: [#9429](https://github.com/owncloud/ocis/pull/9429)

   Tags: web

   We updated ownCloud Web to v8.0.4. Please refer to the changelog (linked) for
   details on the web release.

   - Bugfix [owncloud/web#10814](https://github.com/owncloud/web/issues/10814):
   Vertical scroll for OcModal on small screens - Bugfix
   [owncloud/web#10918](https://github.com/owncloud/web/issues/10918): Resource
   deselection on right-click - Bugfix
   [owncloud/web#10920](https://github.com/owncloud/web/pull/10920): Resources with
   name consist of number won't show up in trash bin - Bugfix
   [owncloud/web#10941](https://github.com/owncloud/web/issues/10941): Space not
   updating on navigation - Bugfix
   [owncloud/web#11063](https://github.com/owncloud/web/issues/11063): Enforce
   shortcut URL protocol - Bugfix
   [owncloud/web#11092](https://github.com/owncloud/web/issues/11092): Browser
   confirmation dialog after closing editor - Bugfix
   [owncloud/web#11091](https://github.com/owncloud/web/issues/11091): Button focus
   when closing editor - Bugfix
   [owncloud/web#10942](https://github.com/owncloud/web/issues/10942): Keyboard
   navigation breaking - Bugfix
   [owncloud/web#11086](https://github.com/owncloud/web/pull/11086): Opening public
   links with an expired token

   https://github.com/owncloud/ocis/pull/9429
   https://github.com/owncloud/ocis/pull/9510
   https://github.com/owncloud/web/releases/tag/v8.0.3
   https://github.com/owncloud/web/releases/tag/v8.0.4

* Enhancement - Add cli to purge revisions: [#9497](https://github.com/owncloud/ocis/pull/9497)

   Adds a cli that allows removing all revisions for a storage-provider.

   https://github.com/owncloud/ocis/pull/9497

# Changelog for [6.1.0] (2024-07-08)

The following sections list the changes for 6.1.0.

[6.1.0]: https://github.com/owncloud/ocis/compare/v6.0.0...v6.1.0

## Summary

* Bugfix - Fix sharing-ng permission listings for personal and virtual drive items: [#9438](https://github.com/owncloud/ocis/pull/9438)
* Bugfix - Add inotify-tools and bash packages to docker files: [#9440](https://github.com/owncloud/ocis/pull/9440)
* Bugfix - Allow all uploads to restart: [#9465](https://github.com/owncloud/ocis/pull/9465)
* Bugfix - Fix the email notification service: [#9467](https://github.com/owncloud/ocis/pull/9467)
* Bugfix - Fix Password Reset: [#9479](https://github.com/owncloud/ocis/pull/9479)
* Bugfix - Fixed the email template: [#9484](https://github.com/owncloud/ocis/pull/9484)
* Bugfix - Polish secure view: [#9532](https://github.com/owncloud/ocis/pull/9532)
* Enhancement - Rudimentary OCM support in graph: [#8909](https://github.com/owncloud/ocis/pull/8909)
* Enhancement - Activitylog API: [#9361](https://github.com/owncloud/ocis/pull/9361)
* Enhancement - Add the backchannel logout event: [#9447](https://github.com/owncloud/ocis/pull/9447)
* Enhancement - Add fail flag to consistency check: [#9447](https://github.com/owncloud/ocis/pull/9447)
* Enhancement - Configurable OCM timeouts: [#9450](https://github.com/owncloud/ocis/pull/9450)
* Enhancement - Deprecate gateway environment variables: [#9451](https://github.com/owncloud/ocis/pull/9451)
* Enhancement - Allow reindexing all spaces: [#9456](https://github.com/owncloud/ocis/pull/9456)
* Enhancement - Autoprovision group memberships: [#9458](https://github.com/owncloud/ocis/pull/9458)
* Enhancement - Allow disable versioning: [#9473](https://github.com/owncloud/ocis/pull/9473)
* Enhancement - Empty trash directories: [#9483](https://github.com/owncloud/ocis/pull/9483)
* Enhancement - Various fixes for the activitylog service: [#9485](https://github.com/owncloud/ocis/pull/9485)
* Enhancement - Add cli to purge revisions: [#9497](https://github.com/owncloud/ocis/pull/9497)
* Enhancement - Update web to v9.1.0: [#9547](https://github.com/owncloud/ocis/pull/9547)
* Enhancement - Bump reva to v2.21.0: [#9556](https://github.com/owncloud/ocis/pull/9556)

## Details

* Bugfix - Fix sharing-ng permission listings for personal and virtual drive items: [#9438](https://github.com/owncloud/ocis/pull/9438)

   Fixes an issue where the sharing-ng service was not able to list permissions for
   personal and virtual drive items.

   https://github.com/owncloud/ocis/issues/8922
   https://github.com/owncloud/ocis/pull/9438

* Bugfix - Add inotify-tools and bash packages to docker files: [#9440](https://github.com/owncloud/ocis/pull/9440)

   We need both packages to make posixfs work. Later, once the golang package is
   fixed to not depend on bash any more, bash can be removed again.

   https://github.com/owncloud/ocis/pull/9440

* Bugfix - Allow all uploads to restart: [#9465](https://github.com/owncloud/ocis/pull/9465)

   On postprocessing-restart, allow all uploads to restart even if one fails.

   https://github.com/owncloud/ocis/pull/9465

* Bugfix - Fix the email notification service: [#9467](https://github.com/owncloud/ocis/pull/9467)

   We fixed an error in the notification service that caused the email notification
   to fail when the user's display name contained special characters.

   https://github.com/owncloud/ocis/issues/9402
   https://github.com/owncloud/ocis/pull/9467

* Bugfix - Fix Password Reset: [#9479](https://github.com/owncloud/ocis/pull/9479)

   The `ocis idm resetpassword` always used the hardcoded `admin` name for the
   user. Now user name can be specified via the `--user-name` (`-u`) flag.

   https://github.com/owncloud/ocis/pull/9479

* Bugfix - Fixed the email template: [#9484](https://github.com/owncloud/ocis/pull/9484)

   Fixed the email template when the description was marked as a link.

   https://github.com/owncloud/ocis/issues/8424
   https://github.com/owncloud/ocis/pull/9484

* Bugfix - Polish secure view: [#9532](https://github.com/owncloud/ocis/pull/9532)

   We fixed a bug where viewing pdf files in secure view mode was not possible.
   Secure view access on space roots was dropped because of unwanted side effects.

   https://github.com/owncloud/ocis/pull/9532

* Enhancement - Rudimentary OCM support in graph: [#8909](https://github.com/owncloud/ocis/pull/8909)

   We now allow creating and accepting OCM shares.

   https://github.com/owncloud/ocis/pull/8909

* Enhancement - Activitylog API: [#9361](https://github.com/owncloud/ocis/pull/9361)

   Adds an api to the `activitylog` service which allows retrieving data by clients
   to show item activities

   https://github.com/owncloud/ocis/pull/9361

* Enhancement - Add the backchannel logout event: [#9447](https://github.com/owncloud/ocis/pull/9447)

   We've added the backchannel logout event

   https://github.com/owncloud/ocis/issues/9355
   https://github.com/owncloud/ocis/pull/9447

* Enhancement - Add fail flag to consistency check: [#9447](https://github.com/owncloud/ocis/pull/9447)

   We added a `--fail` flag to the `ocis backup consistency` command. If set to
   true, the command will return a non-zero exit code if any inconsistencies are
   found. This allows you to use the command in scripts and CI/CD pipelines to
   ensure that backups are consistent.

   https://github.com/owncloud/ocis/pull/9447

* Enhancement - Configurable OCM timeouts: [#9450](https://github.com/owncloud/ocis/pull/9450)

   We added `OCM_OCM_INVITE_MANAGER_TOKEN_EXPIRATION` and
   `OCM_OCM_INVITE_MANAGER_TIMEOUT` to allow changing the default invite token
   duration as well as the request timeout for requests made to other instances.

   https://github.com/owncloud/ocis/pull/9450

* Enhancement - Deprecate gateway environment variables: [#9451](https://github.com/owncloud/ocis/pull/9451)

   Deprecate service specific `_GATEWAY_NAME` env vars. It makes no sense to point
   one specific service to a different gateway.

   https://github.com/owncloud/ocis/pull/9451

* Enhancement - Allow reindexing all spaces: [#9456](https://github.com/owncloud/ocis/pull/9456)

   Adds a `--all-spaces` flag to the `ocis search index` command to allow
   reindexing all spaces at once.

   https://github.com/owncloud/ocis/pull/9456

* Enhancement - Autoprovision group memberships: [#9458](https://github.com/owncloud/ocis/pull/9458)

   When PROXY_AUTOPROVISION_ACCOUNTS is enabled it is now possible to automatically
   maintain the group memberships of users via a configurable OIDC claim.

   https://github.com/owncloud/ocis/issues/5538
   https://github.com/owncloud/ocis/pull/9458

* Enhancement - Allow disable versioning: [#9473](https://github.com/owncloud/ocis/pull/9473)

   Adds new configuration options to disable versioning for the storage providers

   https://github.com/owncloud/ocis/pull/9473

* Enhancement - Empty trash directories: [#9483](https://github.com/owncloud/ocis/pull/9483)

   We have added a cli-command that allows cleaning up empty directories in the
   trashbins folder structure in decomposedFS.

   https://github.com/owncloud/ocis/issues/9393
   https://github.com/owncloud/ocis/issues/9271
   https://github.com/owncloud/ocis/pull/9483

* Enhancement - Various fixes for the activitylog service: [#9485](https://github.com/owncloud/ocis/pull/9485)

   First round of fixes to make the activitylog service more robust and reliable.

   https://github.com/owncloud/ocis/pull/9485
   https://github.com/owncloud/ocis/pull/9467

* Enhancement - Add cli to purge revisions: [#9497](https://github.com/owncloud/ocis/pull/9497)

   Adds a cli that allows removing all revisions for a storage-provider.

   https://github.com/owncloud/ocis/pull/9497

* Enhancement - Update web to v9.1.0: [#9547](https://github.com/owncloud/ocis/pull/9547)

   Tags: web

   We updated ownCloud Web to v9.1.0. Please refer to the changelog (linked) for
   details on the web release.

  * Bugfix [owncloud/web#11058](https://github.com/owncloud/web/pull/11058): Resetting user after logout
  * Bugfix [owncloud/web#11059](https://github.com/owncloud/web/pull/11059): Admin settings UI update after save
  * Bugfix [owncloud/web#11068](https://github.com/owncloud/web/pull/11068): Editor save after token renewal
  * Bugfix [owncloud/web#11132](https://github.com/owncloud/web/pull/11132): Trash bin breaking on navigation
  * Bugfix [owncloud/web#11135](https://github.com/owncloud/web/issues/11135): Tooltips in trashbin covered
  * Bugfix [owncloud/web#11137](https://github.com/owncloud/web/pull/11137): Duplicated elements on public link page
  * Bugfix [owncloud/web#11139](https://github.com/owncloud/web/pull/11139): Secure view default action
  * Enhancement [owncloud/web#5387](https://github.com/owncloud/web/issues/5387): Accessibility improvements
  * Enhancement [owncloud/web#10996](https://github.com/owncloud/web/pull/10996): Activities sidebar app panel
  * Enhancement [owncloud/web#11054](https://github.com/owncloud/web/pull/11054): Consistent initial loading spinner
  * Enhancement [owncloud/web#11057](https://github.com/owncloud/web/pull/11057): Add action drop down to app top bar
  * Enhancement [owncloud/web#11060](https://github.com/owncloud/web/pull/11060): Decrease text editor loading times
  * Enhancement [owncloud/web#11077](https://github.com/owncloud/web/pull/11077): Reduce network load on token renewal
  * Enhancement [owncloud/web#11085](https://github.com/owncloud/web/pull/11085): Open file directly from app
  * Enhancement [owncloud/web#11093](https://github.com/owncloud/web/pull/11093): Enable default autosave in editors

   https://github.com/owncloud/ocis/pull/9547
   https://github.com/owncloud/web/releases/tag/v9.1.0

* Enhancement - Bump reva to v2.21.0: [#9556](https://github.com/owncloud/ocis/pull/9556)

  *   Bugfix [cs3org/reva#4740](https://github.com/cs3org/reva/pull/4740): Disallow reserved filenames
  *   Bugfix [cs3org/reva#4748](https://github.com/cs3org/reva/pull/4748): Quotes in dav Content-Disposition header
  *   Bugfix [cs3org/reva#4750](https://github.com/cs3org/reva/pull/4750): Validate a space path
  *   Enhancement [cs3org/reva#4737](https://github.com/cs3org/reva/pull/4737): Add the backchannel logout event
  *   Enhancement [cs3org/reva#4749](https://github.com/cs3org/reva/pull/4749): DAV error codes
  *   Enhancement [cs3org/reva#4742](https://github.com/cs3org/reva/pull/4742): Expose disable-versioning configuration option
  *   Enhancement [cs3org/reva#4739](https://github.com/cs3org/reva/pull/4739): Improve posixfs storage driver
  *   Enhancement [cs3org/reva#4738](https://github.com/cs3org/reva/pull/4738): Add GetServiceUserToken() method to utils pkg

   https://github.com/owncloud/ocis/pull/9556
   https://github.com/owncloud/ocis/pull/9473

# Changelog for [6.0.0] (2024-06-19)

The following sections list the changes for 6.0.0.

[6.0.0]: https://github.com/owncloud/ocis/compare/v5.0.5...v6.0.0

## Summary

* Bugfix - Fix an error when lock/unlock a public shared file: [#8472](https://github.com/owncloud/ocis/pull/8472)
* Bugfix - Fix the docker-compose wopi: [#8483](https://github.com/owncloud/ocis/pull/8483)
* Bugfix - Fix remove/update share permissions: [#8529](https://github.com/owncloud/ocis/pull/8529)
* Bugfix - Correct the default mapping of roles: [#8534](https://github.com/owncloud/ocis/pull/8534)
* Bugfix - Fix graph drive invite: [#8538](https://github.com/owncloud/ocis/pull/8538)
* Bugfix - Fix the mount points naming: [#8543](https://github.com/owncloud/ocis/pull/8543)
* Bugfix - We now always select the next clients when autoaccepting shares: [#8570](https://github.com/owncloud/ocis/pull/8570)
* Bugfix - Always select next before making calls: [#8578](https://github.com/owncloud/ocis/pull/8578)
* Bugfix - Fix sharing invite on virtual drive: [#8609](https://github.com/owncloud/ocis/pull/8609)
* Bugfix - Prevent copying a file to a parent folder: [#8649](https://github.com/owncloud/ocis/pull/8649)
* Bugfix - Disable Multipart uploads: [#8666](https://github.com/owncloud/ocis/pull/8666)
* Bugfix - Internal links shouldn't have a password: [#8668](https://github.com/owncloud/ocis/pull/8668)
* Bugfix - Fix uploading via a public link: [#8702](https://github.com/owncloud/ocis/pull/8702)
* Bugfix - Mask user email in output: [#8726](https://github.com/owncloud/ocis/issues/8726)
* Bugfix - Fix restarting of postprocessing: [#8782](https://github.com/owncloud/ocis/pull/8782)
* Bugfix - Fix the create personal space cache: [#8799](https://github.com/owncloud/ocis/pull/8799)
* Bugfix - Fix removing groups from space: [#8803](https://github.com/owncloud/ocis/pull/8803)
* Bugfix - Validate conditions for sharing roles by resource type: [#8815](https://github.com/owncloud/ocis/pull/8815)
* Bugfix - Fix creating the drive item: [#8817](https://github.com/owncloud/ocis/pull/8817)
* Bugfix - Fix unmount item from share: [#8827](https://github.com/owncloud/ocis/pull/8827)
* Bugfix - Fix creating new WOPI documents on public shares: [#8828](https://github.com/owncloud/ocis/pull/8828)
* Bugfix - Nats reconnects: [#8880](https://github.com/owncloud/ocis/pull/8880)
* Bugfix - Update the admin user role assignment to enforce the config: [#8897](https://github.com/owncloud/ocis/pull/8897)
* Bugfix - Fix affected users on sses: [#8928](https://github.com/owncloud/ocis/pull/8928)
* Bugfix - Fix well-known rewrite endpoint: [#8946](https://github.com/owncloud/ocis/pull/8946)
* Bugfix - Crash when processing crafted TIFF files: [#8981](https://github.com/owncloud/ocis/pull/8981)
* Bugfix - Fix collaboration registry setting: [#9105](https://github.com/owncloud/ocis/pull/9105)
* Bugfix - Service startup of WOPI example: [#9127](https://github.com/owncloud/ocis/pull/9127)
* Bugfix - Fix the status code for multiple mount and unmount share: [#9193](https://github.com/owncloud/ocis/pull/9193)
* Bugfix - Don't show thumbnails for secureview shares: [#9299](https://github.com/owncloud/ocis/pull/9299)
* Bugfix - Fix share update: [#9301](https://github.com/owncloud/ocis/pull/9301)
* Bugfix - Fix the error translation from utils: [#9331](https://github.com/owncloud/ocis/pull/9331)
* Bugfix - Fix the settings metedata tests: [#9341](https://github.com/owncloud/ocis/pull/9341)
* Bugfix - The hidden shares have been excluded from a search result: [#9371](https://github.com/owncloud/ocis/pull/9371)
* Bugfix - Encode Registry Keys: [#9385](https://github.com/owncloud/ocis/pull/9385)
* Change - Change the default store for presigned keys to nats-js-kv: [#8419](https://github.com/owncloud/ocis/pull/8419)
* Change - Disable resharing by default for deprecation: [#8653](https://github.com/owncloud/ocis/pull/8653)
* Change - The `filesystem` backend for the settings service has been removed: [#9138](https://github.com/owncloud/ocis/pull/9138)
* Change - Define maximum input image dimensions and size when generating previews: [#9360](https://github.com/owncloud/ocis/pull/9360)
* Enhancement - Introduce staticroutes package & remove well-known OIDC middleware: [#6095](https://github.com/owncloud/ocis/issues/6095)
* Enhancement - Graphs endpoint for mounting and unmounting shares: [#7885](https://github.com/owncloud/ocis/pull/7885)
* Enhancement - Add epub reader to web default apps: [#8410](https://github.com/owncloud/ocis/pull/8410)
* Enhancement - Change Cors default settings: [#8518](https://github.com/owncloud/ocis/pull/8518)
* Enhancement - Custom WEB App Loading: [#8523](https://github.com/owncloud/ocis/pull/8523)
* Enhancement - Update to go 1.22: [#8586](https://github.com/owncloud/ocis/pull/8586)
* Enhancement - Send more sse events: [#8587](https://github.com/owncloud/ocis/pull/8587)
* Enhancement - Send SSE when file is locked/unlocked: [#8602](https://github.com/owncloud/ocis/pull/8602)
* Enhancement - Add the spaceID to sse: [#8614](https://github.com/owncloud/ocis/pull/8614)
* Enhancement - The graph endpoints for listing permission works for spaces now: [#8642](https://github.com/owncloud/ocis/pull/8642)
* Enhancement - Bump keycloak: [#8687](https://github.com/owncloud/ocis/pull/8687)
* Enhancement - Make IDP cookies same site strict: [#8716](https://github.com/owncloud/ocis/pull/8716)
* Enhancement - Make server side space templates production ready: [#8723](https://github.com/owncloud/ocis/pull/8723)
* Enhancement - Sharing NG role names and descriptions: [#8743](https://github.com/owncloud/ocis/pull/8743)
* Enhancement - Ability to Change Share Item Visibility in Graph API: [#8750](https://github.com/owncloud/ocis/pull/8750)
* Enhancement - Enable web extension drawio by default: [#8760](https://github.com/owncloud/ocis/pull/8760)
* Enhancement - Remove resharing: [#8762](https://github.com/owncloud/ocis/pull/8762)
* Enhancement - Add CSP and other security related headers to oCIS: [#8777](https://github.com/owncloud/ocis/pull/8777)
* Enhancement - Add FileTouched SSE Event: [#8778](https://github.com/owncloud/ocis/pull/8778)
* Enhancement - Prepare runners to start the services: [#8802](https://github.com/owncloud/ocis/pull/8802)
* Enhancement - Sharing SSEs: [#8854](https://github.com/owncloud/ocis/pull/8854)
* Enhancement - Secure viewer share role: [#8907](https://github.com/owncloud/ocis/pull/8907)
* Enhancement - Add Link SSEs: [#8908](https://github.com/owncloud/ocis/pull/8908)
* Enhancement - ShareeIDs in SSEs: [#8915](https://github.com/owncloud/ocis/pull/8915)
* Enhancement - Allow to resolve public shares without the ocs tokeninfo endpoint: [#8926](https://github.com/owncloud/ocis/pull/8926)
* Enhancement - Initiator-IDs: [#8936](https://github.com/owncloud/ocis/pull/8936)
* Enhancement - Add endpoint for getting drive items: [#8939](https://github.com/owncloud/ocis/pull/8939)
* Enhancement - Improve infected file handling: [#8947](https://github.com/owncloud/ocis/pull/8947)
* Enhancement - Configurable claims for auto-provisioning user accounts: [#8952](https://github.com/owncloud/ocis/pull/8952)
* Enhancement - Bump nats-js-kv pkg: [#8953](https://github.com/owncloud/ocis/pull/8953)
* Enhancement - Graph permission created date time: [#8954](https://github.com/owncloud/ocis/pull/8954)
* Enhancement - Add virus filter to sessions command: [#9041](https://github.com/owncloud/ocis/pull/9041)
* Enhancement - Assimilate `clean` into `sessions` command: [#9041](https://github.com/owncloud/ocis/pull/9041)
* Enhancement - Add remote item id to WebDAV report responses: [#9094](https://github.com/owncloud/ocis/issues/9094)
* Enhancement - Theme Processing and Logo Customization: [#9133](https://github.com/owncloud/ocis/pull/9133)
* Enhancement - Add watermark text: [#9144](https://github.com/owncloud/ocis/pull/9144)
* Enhancement - Update selected attributes of autoprovisioned users: [#9166](https://github.com/owncloud/ocis/pull/9166)
* Enhancement - Limit concurrent thumbnail requests: [#9199](https://github.com/owncloud/ocis/pull/9199)
* Enhancement - The storage-users doc updated: [#9228](https://github.com/owncloud/ocis/pull/9228)
* Enhancement - Docker compose example for ClamAV: [#9229](https://github.com/owncloud/ocis/pull/9229)
* Enhancement - Add command to check ocis backup consistency: [#9238](https://github.com/owncloud/ocis/pull/9238)
* Enhancement - Web server compression: [#9287](https://github.com/owncloud/ocis/pull/9287)
* Enhancement - Add secureview flag when listing apps via http: [#9289](https://github.com/owncloud/ocis/pull/9289)
* Enhancement - Activitylog Service: [#9327](https://github.com/owncloud/ocis/pull/9327)
* Enhancement - Update web to v9.0.0-alpha.7: [#9395](https://github.com/owncloud/ocis/pull/9395)
* Enhancement - Bump Reva to v2.20.0: [#9415](https://github.com/owncloud/ocis/pull/9415)

## Details

* Bugfix - Fix an error when lock/unlock a public shared file: [#8472](https://github.com/owncloud/ocis/pull/8472)

   We fixed a bug when anonymous user with viewer role in public link of a folder
   can lock/unlock a file inside it

   https://github.com/owncloud/ocis/issues/7785
   https://github.com/owncloud/ocis/pull/8472

* Bugfix - Fix the docker-compose wopi: [#8483](https://github.com/owncloud/ocis/pull/8483)

   We fixed an issue when Collabora is not available time by time after running the
   docker-compose wopi deployment

   https://github.com/owncloud/ocis/issues/8474
   https://github.com/owncloud/ocis/pull/8483

* Bugfix - Fix remove/update share permissions: [#8529](https://github.com/owncloud/ocis/pull/8529)

   This is a workaround that should prevent removing or changing the share
   permissions when the file is locked. These limitations have to be removed after
   the wopi server will be able to unlock the file properly. These limitations are
   not spread on the files inside the shared folder.

   https://github.com/owncloud/ocis/issues/8273
   https://github.com/owncloud/ocis/pull/8529
   https://github.com/cs3org/reva/pull/4534

* Bugfix - Correct the default mapping of roles: [#8534](https://github.com/owncloud/ocis/pull/8534)

   The default config for the OIDC role mapping was incorrect. Lightweight users
   are now assignable.

   https://github.com/owncloud/ocis/pull/8534

* Bugfix - Fix graph drive invite: [#8538](https://github.com/owncloud/ocis/pull/8538)

   We fixed the issue when sharing of personal drive is allowed via graph

   https://github.com/owncloud/ocis/issues/8494
   https://github.com/owncloud/ocis/pull/8538

* Bugfix - Fix the mount points naming: [#8543](https://github.com/owncloud/ocis/pull/8543)

   We fixed a bug that caused inconsistent naming when multiple users share the
   resource with same name to another user.

   https://github.com/owncloud/ocis/issues/8471
   https://github.com/owncloud/ocis/pull/8543

* Bugfix - We now always select the next clients when autoaccepting shares: [#8570](https://github.com/owncloud/ocis/pull/8570)

   https://github.com/owncloud/ocis/pull/8570

* Bugfix - Always select next before making calls: [#8578](https://github.com/owncloud/ocis/pull/8578)

   We now select the next client more often to spread out load

   https://github.com/owncloud/ocis/pull/8578

* Bugfix - Fix sharing invite on virtual drive: [#8609](https://github.com/owncloud/ocis/pull/8609)

   We fixed the issue when sharing of virtual drive with other users was allowed

   https://github.com/owncloud/ocis/issues/8495
   https://github.com/owncloud/ocis/pull/8609

* Bugfix - Prevent copying a file to a parent folder: [#8649](https://github.com/owncloud/ocis/pull/8649)

   When copying a file to a parent folder, the file would be copied to the parent
   folder, but the file would not be removed from the original folder.

   https://github.com/owncloud/ocis/issues/1230
   https://github.com/owncloud/ocis/pull/8649
   https://github.com/cs3org/reva/pull/4571
   %60

* Bugfix - Disable Multipart uploads: [#8666](https://github.com/owncloud/ocis/pull/8666)

   Disables multiparts uploads as they lead to high memory consumption

   https://github.com/owncloud/ocis/pull/8666

* Bugfix - Internal links shouldn't have a password: [#8668](https://github.com/owncloud/ocis/pull/8668)

   Internal links shouldn't have a password when create/update

   https://github.com/owncloud/ocis/issues/8619
   https://github.com/owncloud/ocis/pull/8668

* Bugfix - Fix uploading via a public link: [#8702](https://github.com/owncloud/ocis/pull/8702)

   Fix http error when uploading via a public link

   https://github.com/owncloud/ocis/issues/8699
   https://github.com/owncloud/ocis/pull/8702

* Bugfix - Mask user email in output: [#8726](https://github.com/owncloud/ocis/issues/8726)

   We have fixed a bug where the user email was not masked in the output and the
   user emails could be enumerated through the sharee search. This is the ocis side
   which adds an suiting config option to mask user emails in the output.

   https://github.com/owncloud/ocis/issues/8726
   https://github.com/cs3org/reva/pull/4603
   https://github.com/owncloud/ocis/pull/8764

* Bugfix - Fix restarting of postprocessing: [#8782](https://github.com/owncloud/ocis/pull/8782)

   When an upload is not found, the logic to restart postprocessing was bunked.
   Additionally we extended the upload sessions command to be able to restart the
   uploads without using a second command.

   NOTE: This also includes a breaking fix for the deprecated `ocis storage-users
   uploads list` command

   https://github.com/owncloud/ocis/pull/8782

* Bugfix - Fix the create personal space cache: [#8799](https://github.com/owncloud/ocis/pull/8799)

   We fixed a problem with the config for the create personal space cache which
   resulted in the cache never being used.

   https://github.com/owncloud/ocis/pull/8799

* Bugfix - Fix removing groups from space: [#8803](https://github.com/owncloud/ocis/pull/8803)

   We fixed a bug when unable to remove groups from space via graph

   https://github.com/owncloud/ocis/issues/8768
   https://github.com/owncloud/ocis/pull/8803

* Bugfix - Validate conditions for sharing roles by resource type: [#8815](https://github.com/owncloud/ocis/pull/8815)

   We improved the validation of the allowed sharing roles for specific resource
   type for various sharing related graph API endpoints. This allows e.g. the web
   client to restrict the sharing roles presented to the user based on the type of
   the resource that is being shared.

   https://github.com/owncloud/ocis/issues/8331
   https://github.com/owncloud/ocis/pull/8815

* Bugfix - Fix creating the drive item: [#8817](https://github.com/owncloud/ocis/pull/8817)

   We fixed the issue when creating a drive item with random item id was allowed

   https://github.com/owncloud/ocis/issues/8724
   https://github.com/owncloud/ocis/pull/8817

* Bugfix - Fix unmount item from share: [#8827](https://github.com/owncloud/ocis/pull/8827)

   We fixed the status code returned for the request to delete a driveitem.

   https://github.com/owncloud/ocis/issues/8731
   https://github.com/owncloud/ocis/pull/8827

* Bugfix - Fix creating new WOPI documents on public shares: [#8828](https://github.com/owncloud/ocis/pull/8828)

   Creating a new Office document in a publicly shared folder is now possible.

   https://github.com/owncloud/ocis/issues/8691
   https://github.com/owncloud/ocis/pull/8828

* Bugfix - Nats reconnects: [#8880](https://github.com/owncloud/ocis/pull/8880)

   We fixed the reconnect handling of the natjs kv registry.

   https://github.com/owncloud/ocis/pull/8880

* Bugfix - Update the admin user role assignment to enforce the config: [#8897](https://github.com/owncloud/ocis/pull/8897)

   The admin user role assigment was not updated after the first assignment. We now
   read the assigned role during init and update the admin user ID accordingly if
   the role is not assigned. This is especially needed when the OCIS_ADMIN_USER_ID
   is set after the autoprovisioning of the admin user when it originates from an
   external Identity Provider.

   https://github.com/owncloud/ocis/pull/8897

* Bugfix - Fix affected users on sses: [#8928](https://github.com/owncloud/ocis/pull/8928)

   The AffectedUsers field of sses now only reports affected users.

   https://github.com/owncloud/ocis/pull/8928

* Bugfix - Fix well-known rewrite endpoint: [#8946](https://github.com/owncloud/ocis/pull/8946)

   https://github.com/owncloud/ocis/issues/8703
   https://github.com/owncloud/ocis/pull/8946

* Bugfix - Crash when processing crafted TIFF files: [#8981](https://github.com/owncloud/ocis/pull/8981)

   Fix for a vulnerability with low severity in disintegration/imaging.

   https://github.com/owncloud/ocis/pull/8981
   https://github.com/advisories/GHSA-q7pp-wcgr-pffx

* Bugfix - Fix collaboration registry setting: [#9105](https://github.com/owncloud/ocis/pull/9105)

   Fixed the collaboration service GRPC namespace

   https://github.com/owncloud/ocis/pull/9105

* Bugfix - Service startup of WOPI example: [#9127](https://github.com/owncloud/ocis/pull/9127)

   We fixed a bug in the service startup of the appprovider-onlyoffice in the
   ocis_wopi deployment example.

   https://github.com/owncloud/ocis/pull/9127

* Bugfix - Fix the status code for multiple mount and unmount share: [#9193](https://github.com/owncloud/ocis/pull/9193)

   We fixed the status code for multiple mount and unmount share.

   https://github.com/owncloud/ocis/issues/8876
   https://github.com/owncloud/ocis/pull/9193

* Bugfix - Don't show thumbnails for secureview shares: [#9299](https://github.com/owncloud/ocis/pull/9299)

   We have fixed a bug where thumbnails were shown for secureview shares.

   https://github.com/owncloud/ocis/issues/9249
   https://github.com/owncloud/ocis/pull/9299

* Bugfix - Fix share update: [#9301](https://github.com/owncloud/ocis/pull/9301)

   We fixed the response code when the role/permission is empty on the share update

   https://github.com/owncloud/ocis/issues/8747
   https://github.com/owncloud/ocis/pull/9301

* Bugfix - Fix the error translation from utils: [#9331](https://github.com/owncloud/ocis/pull/9331)

   We've fixed the error translation from the statusCodeError type to CS3 Status
   because the FromCS3Status function converts a CS3 status code into a
   corresponding local Error representation.

   https://github.com/owncloud/ocis/issues/9151
   https://github.com/owncloud/ocis/pull/9331

* Bugfix - Fix the settings metedata tests: [#9341](https://github.com/owncloud/ocis/pull/9341)

   We fix the settings metedata tests that had the data race

   https://github.com/owncloud/ocis/issues/9372
   https://github.com/owncloud/ocis/pull/9341

* Bugfix - The hidden shares have been excluded from a search result: [#9371](https://github.com/owncloud/ocis/pull/9371)

   The hidden shares have been excluded from a search result.

   https://github.com/owncloud/ocis/issues/7383
   https://github.com/owncloud/ocis/pull/9371

* Bugfix - Encode Registry Keys: [#9385](https://github.com/owncloud/ocis/pull/9385)

   Encode the keys of the natsjskv registry as they have always been.

   https://github.com/owncloud/ocis/pull/9385

* Change - Change the default store for presigned keys to nats-js-kv: [#8419](https://github.com/owncloud/ocis/pull/8419)

   We wrapped the store service in a micro store implementation and changed the
   default to the built-in NATS instance.

   https://github.com/owncloud/ocis/pull/8419

* Change - Disable resharing by default for deprecation: [#8653](https://github.com/owncloud/ocis/pull/8653)

   We disabled the resharing feature by default. This feature will be removed from
   the product in the next major release. The resharing feature is not recommended
   for use and should be disabled. Existing reshares will continue to work.

   https://github.com/owncloud/ocis/pull/8653

* Change - The `filesystem` backend for the settings service has been removed: [#9138](https://github.com/owncloud/ocis/pull/9138)

   The only remaining backend for the settings service is `metadata`, which has
   been the default backend since ocis 2.0

   https://github.com/owncloud/ocis/pull/9138

* Change - Define maximum input image dimensions and size when generating previews: [#9360](https://github.com/owncloud/ocis/pull/9360)

   This is a general hardening change to limit processing time and resources of the
   thumbnailer.

   https://github.com/owncloud/ocis/pull/9360
   https://github.com/owncloud/ocis/pull/9035
   https://github.com/owncloud/ocis/pull/9069

* Enhancement - Introduce staticroutes package & remove well-known OIDC middleware: [#6095](https://github.com/owncloud/ocis/issues/6095)

   We have introduced a new static routes package to the proxy. This package is
   responsible for serving static files and oidc well-known endpoint
   `/.well-known/openid-configuration`. We have removed the well-known middleware
   for OIDC and moved it to the newly introduced static routes module in the proxy.

   https://github.com/owncloud/ocis/issues/6095
   https://github.com/owncloud/ocis/pull/8541

* Enhancement - Graphs endpoint for mounting and unmounting shares: [#7885](https://github.com/owncloud/ocis/pull/7885)

   Functionality for mounting (accepting) and unmounting (rejecting) received
   shares has been added to the graph API.

   https://github.com/owncloud/ocis/pull/7885

* Enhancement - Add epub reader to web default apps: [#8410](https://github.com/owncloud/ocis/pull/8410)

   We've added the new epub reader app to the web default apps, so it will be
   enabled and usable by default.

   https://github.com/owncloud/ocis/pull/8410

* Enhancement - Change Cors default settings: [#8518](https://github.com/owncloud/ocis/pull/8518)

   We have changed the default CORS settings to set `Access-Control-Allow-Origin`
   to the `OCIS_URL` if not explicitely set and `Access-Control-Allow-Credentials`
   to `false` if not explicitely set.

   https://github.com/owncloud/ocis/issues/8514
   https://github.com/owncloud/ocis/pull/8518

* Enhancement - Custom WEB App Loading: [#8523](https://github.com/owncloud/ocis/pull/8523)

   We've added a new feature which allows the administrator of the environment to
   provide custom web applications to the users. This feature is useful for
   organizations that have specific web applications that they want to provide to
   their users.

   The users will then be able to access these custom web applications from the web
   ui. For a detailed description of the feature, please read the WEB service
   README.md file.

   https://github.com/owncloud/ocis/issues/8392
   https://github.com/owncloud/ocis/pull/8523

* Enhancement - Update to go 1.22: [#8586](https://github.com/owncloud/ocis/pull/8586)

   We have updated go to version 1.22.

   https://github.com/owncloud/ocis/pull/8586

* Enhancement - Send more sse events: [#8587](https://github.com/owncloud/ocis/pull/8587)

   We added sse events for `ItemTrashed`, `ItemRestored`,`ContainerCreated` and
   `FileRenamed`

   https://github.com/owncloud/ocis/pull/8587

* Enhancement - Send SSE when file is locked/unlocked: [#8602](https://github.com/owncloud/ocis/pull/8602)

   Send sse events when a file is locked or unlocked.

   https://github.com/owncloud/ocis/pull/8602

* Enhancement - Add the spaceID to sse: [#8614](https://github.com/owncloud/ocis/pull/8614)

   Adds the spaceID to all clientlog sse messages

   https://github.com/owncloud/ocis/pull/8614
   https://github.com/owncloud/ocis/pull/8624

* Enhancement - The graph endpoints for listing permission works for spaces now: [#8642](https://github.com/owncloud/ocis/pull/8642)

   We enhanced the 'graph/v1beta1/drives/{{driveid}}/items/{{itemid}}/permissions'
   endpoint to list permission of the space when the 'itemid' refers to a space
   root.

   https://github.com/owncloud/ocis/issues/8352
   https://github.com/owncloud/ocis/pull/8642

* Enhancement - Bump keycloak: [#8687](https://github.com/owncloud/ocis/pull/8687)

   Bumps keycloak version

   https://github.com/owncloud/ocis/issues/8569
   https://github.com/owncloud/ocis/pull/8687

* Enhancement - Make IDP cookies same site strict: [#8716](https://github.com/owncloud/ocis/pull/8716)

   To enhance the security of our application and prevent Cross-Site Request
   Forgery (CSRF) attacks, we have updated the SameSite attribute of the build in
   Identity Provider (IDP) cookies to Strict.

   This change restricts the browser from sending these cookies with any cross-site
   requests, thereby limiting the exposure of the user's session to potential
   threats.

   This update does not impact the existing functionality of the application but
   provides an additional layer of security where needed.

   https://github.com/owncloud/ocis/pull/8716

* Enhancement - Make server side space templates production ready: [#8723](https://github.com/owncloud/ocis/pull/8723)

   Fixes several smaller bugs and adds some improvements to space templates,
   introduced with https://github.com/owncloud/ocis/pull/8558

   https://github.com/owncloud/ocis/pull/8723

* Enhancement - Sharing NG role names and descriptions: [#8743](https://github.com/owncloud/ocis/pull/8743)

   We've adjusted the display names and descriptions of the sharing NG roles to
   align with the previously agreed upon terms.

   https://github.com/owncloud/ocis/pull/8743

* Enhancement - Ability to Change Share Item Visibility in Graph API: [#8750](https://github.com/owncloud/ocis/pull/8750)

   Introduce the `PATCH /graph/v1beta1/drives/{driveID}/items/{itemID}` Graph API
   endpoint which allows updating individual Drive Items.

   At the moment, only the share visibility is considered changeable, but in the
   future, more properties can be added to this endpoint.

   This enhancement is needed for the user interface, allowing specific shares to
   be hidden or unhidden as needed, thereby improving the user experience.

   https://github.com/owncloud/ocis/issues/8654
   https://github.com/owncloud/ocis/pull/8750

* Enhancement - Enable web extension drawio by default: [#8760](https://github.com/owncloud/ocis/pull/8760)

   Enable web extension drawio by default

   https://github.com/owncloud/ocis/pull/8760

* Enhancement - Remove resharing: [#8762](https://github.com/owncloud/ocis/pull/8762)

   Removed resharing feature from codebase

   https://github.com/owncloud/ocis/pull/8762

* Enhancement - Add CSP and other security related headers to oCIS: [#8777](https://github.com/owncloud/ocis/pull/8777)

   General hardening of oCIS

   https://github.com/owncloud/ocis/pull/8777
   https://github.com/owncloud/ocis/pull/9025
   https://github.com/owncloud/ocis/pull/9167
   https://github.com/owncloud/ocis/pull/9313

* Enhancement - Add FileTouched SSE Event: [#8778](https://github.com/owncloud/ocis/pull/8778)

   Send an sse when a file is touched (aka 0 byte upload)

   https://github.com/owncloud/ocis/pull/8778

* Enhancement - Prepare runners to start the services: [#8802](https://github.com/owncloud/ocis/pull/8802)

   The runners will improve and make service startup easier. The runner's behavior
   is more predictable with clear expectations.

   https://github.com/owncloud/ocis/pull/8802

* Enhancement - Sharing SSEs: [#8854](https://github.com/owncloud/ocis/pull/8854)

   Added server side events for item moved, share created/updated/removed, space
   membership created/removed.

   https://github.com/owncloud/ocis/pull/8854
   https://github.com/owncloud/ocis/pull/8875

* Enhancement - Secure viewer share role: [#8907](https://github.com/owncloud/ocis/pull/8907)

   A new share role "Secure viewer" has been added. This role is applicable for
   files, folders and spaces and only allows viewing them (and their content).

   https://github.com/owncloud/ocis/pull/8907

* Enhancement - Add Link SSEs: [#8908](https://github.com/owncloud/ocis/pull/8908)

   Add sses for link created/updated/removed.

   https://github.com/owncloud/ocis/pull/8908

* Enhancement - ShareeIDs in SSEs: [#8915](https://github.com/owncloud/ocis/pull/8915)

   We will now send a list of userIDs (one or in case of a group share multiple) on
   share related SSEs

   https://github.com/owncloud/ocis/pull/8915

* Enhancement - Allow to resolve public shares without the ocs tokeninfo endpoint: [#8926](https://github.com/owncloud/ocis/pull/8926)

   Instead of querying the /v1.php/apps/files_sharing/api/v1/tokeninfo/ endpoint, a
   client can now resolve public and internal links by sending a PROPFIND request
   to /dav/public-files/{sharetoken}

  * authenticated clients accessing an internal link are redirected to the "real" resource (`/dav/spaces/{target-resource-id}
  * authenticated clients are able to resolve public links like before. For password protected links they need to supply the password even if they have access to the underlying resource by other means.
  * unauthenticated clients accessing an internal link get a 401 returned with  WWW-Authenticate set to Bearer (so that the client knows that it need to get a token via the IDP login page.
  * unauthenticated clients accessing a password protected link get a 401 returned with an error message to indicate the requirement for needing the link's password.

   https://github.com/owncloud/ocis/issues/8858
   https://github.com/owncloud/ocis/pull/8926
   https://github.com/cs3org/reva/pull/4653

* Enhancement - Initiator-IDs: [#8936](https://github.com/owncloud/ocis/pull/8936)

   Allows sending a header `Initiator-ID` on http requests. This id will be added
   to sse events so clients can figure out if their particular instance was
   triggering the event. Additionally this adds the etag of the file/folder to all
   sse events.

   https://github.com/owncloud/ocis/pull/8936
   https://github.com/owncloud/ocis/pull/8701

* Enhancement - Add endpoint for getting drive items: [#8939](https://github.com/owncloud/ocis/pull/8939)

   An endpoint for getting drive items via ID has been added.

   https://github.com/owncloud/ocis/issues/8915
   https://github.com/owncloud/ocis/pull/8939

* Enhancement - Improve infected file handling: [#8947](https://github.com/owncloud/ocis/pull/8947)

   Reworks virus handling.Shows scandate and outcome on ocis storage-users uploads
   sessions. Avoids retrying infected files on ocis postprocessing restart.

   https://github.com/owncloud/ocis/pull/8947

* Enhancement - Configurable claims for auto-provisioning user accounts: [#8952](https://github.com/owncloud/ocis/pull/8952)

   We introduce the new environment variables "PROXY_AUTOPROVISION_CLAIM_USERNAME",
   "PROXY_AUTOPROVISION_CLAIM_EMAIL", and "PROXY_AUTOPROVISION_CLAIM_DISPLAYNAME"
   which can be used to configure the OIDC claims that should be used for
   auto-provisioning user accounts.

   The automatic fallback to use the 'email' claim value as the username when the
   'preferred_username' claim is not set, has been removed.

   Also it is now possible to autoprovision users without an email address.

   https://github.com/owncloud/ocis/issues/8635
   https://github.com/owncloud/ocis/issues/6909
   https://github.com/owncloud/ocis/pull/8952

* Enhancement - Bump nats-js-kv pkg: [#8953](https://github.com/owncloud/ocis/pull/8953)

   Uses official nats-js-kv package now. Moves away from custom fork.

   https://github.com/owncloud/ocis/pull/8953

* Enhancement - Graph permission created date time: [#8954](https://github.com/owncloud/ocis/pull/8954)

   We've added the created date time to graph permission objects.

   https://github.com/owncloud/ocis/issues/8749
   https://github.com/owncloud/ocis/pull/8954

* Enhancement - Add virus filter to sessions command: [#9041](https://github.com/owncloud/ocis/pull/9041)

   Allow filtering upload session by virus status (has-virus=true/false)

   https://github.com/owncloud/ocis/pull/9041

* Enhancement - Assimilate `clean` into `sessions` command: [#9041](https://github.com/owncloud/ocis/pull/9041)

   We deprecated `ocis storage-user uploads clean` and added the same logic to
   `ocis storage-users uploads session --clean`

   https://github.com/owncloud/ocis/pull/9041

* Enhancement - Add remote item id to WebDAV report responses: [#9094](https://github.com/owncloud/ocis/issues/9094)

   The remote item id has been added to WebDAV `REPORT` responses.

   https://github.com/owncloud/ocis/issues/9094
   https://github.com/owncloud/ocis/pull/9095

* Enhancement - Theme Processing and Logo Customization: [#9133](https://github.com/owncloud/ocis/pull/9133)

   We have made significant improvements to the theme processing in Infinite Scale.
   The changes include:

   - Enhanced the way themes are composed. Now, the final theme is a combination of
   the built-in theme and the custom theme provided by the administrator via
   `WEB_ASSET_THEMES_PATH` and `WEB_UI_THEME_PATH`. - Introduced a new mechanism to
   load custom assets. This is particularly useful when a single asset, such as a
   logo, needs to be overwritten. - Fixed the logo customization option.
   Previously, small theme changes would copy the entire theme. Now, only the
   changed keys are considered, making the process more efficient. - Default themes
   are now part of ocis. This change simplifies the theme management process for
   web.

   These changes enhance the robustness of the theme handling in Infinite Scale and
   provide a better user experience.

   https://github.com/owncloud/ocis/issues/8966
   https://github.com/owncloud/ocis/pull/9133

* Enhancement - Add watermark text: [#9144](https://github.com/owncloud/ocis/pull/9144)

   We've added the watermark text for the Secure View mode.

   https://github.com/owncloud/ocis/pull/9144

* Enhancement - Update selected attributes of autoprovisioned users: [#9166](https://github.com/owncloud/ocis/pull/9166)

   When autoprovisioning is enabled, we now update autoprovisioned users when their
   display name or email address claims change.

   https://github.com/owncloud/ocis/issues/8955
   https://github.com/owncloud/ocis/pull/9166

* Enhancement - Limit concurrent thumbnail requests: [#9199](https://github.com/owncloud/ocis/pull/9199)

   The number of concurrent requests to the thumbnail service can be limited now to
   have more control over the consumed system resources.

   https://github.com/owncloud/ocis/pull/9199

* Enhancement - The storage-users doc updated: [#9228](https://github.com/owncloud/ocis/pull/9228)

   The storage-users doc was updated, added the details to the 'Restore Trash-Bins
   Items' section.

   https://github.com/owncloud/ocis/pull/9228

* Enhancement - Docker compose example for ClamAV: [#9229](https://github.com/owncloud/ocis/pull/9229)

   This PR adds a docker compose example for running a local oCIS together with
   ClamAV as virus scanner. The example is for demonstration purposes only and
   should not be used in production.

   https://github.com/owncloud/ocis/pull/9229

* Enhancement - Add command to check ocis backup consistency: [#9238](https://github.com/owncloud/ocis/pull/9238)

   Adds a command that checks the consistency of an ocis backup.

   https://github.com/owncloud/ocis/pull/9238

* Enhancement - Web server compression: [#9287](https://github.com/owncloud/ocis/pull/9287)

   We've added a compression middleware to the web server to reduce the request
   size when delivering static files. This speeds up loading times in web clients.

   https://github.com/owncloud/web/issues/7964
   https://github.com/owncloud/ocis/pull/9287

* Enhancement - Add secureview flag when listing apps via http: [#9289](https://github.com/owncloud/ocis/pull/9289)

   To allow clients to see which application supports secure view, we add a flag to
   the http response when the app service name matches a configured secure view app
   provider. The app can be configured by setting
   `FRONTEND_APP_HANDLER_SECURE_VIEW_APP_ADDR` to the address of the registered CS3
   app provider.

   https://github.com/owncloud/ocis/pull/9289
   https://github.com/owncloud/ocis/pull/9280
   https://github.com/owncloud/ocis/pull/9277

* Enhancement - Activitylog Service: [#9327](https://github.com/owncloud/ocis/pull/9327)

   Adds a new service `activitylog` which stores events (activities) per resource.
   This data can be retrieved by clients to show item activities

   https://github.com/owncloud/ocis/pull/9327

* Enhancement - Update web to v9.0.0-alpha.7: [#9395](https://github.com/owncloud/ocis/pull/9395)

   Tags: web

   We updated ownCloud Web to v9.0.0-alpha.7. Please refer to the changelog
   (linked) for details on the web release.

  * Bugfix [owncloud/web#10377](https://github.com/owncloud/web/pull/10377): User data not updated while altering own user
  * Bugfix [owncloud/web#10417](https://github.com/owncloud/web/pull/10417): Admin settings keyboard navigation
  * Bugfix [owncloud/web#10517](https://github.com/owncloud/web/pull/10517): Load thumbnail when postprocessing is finished
  * Bugfix [owncloud/web#10551](https://github.com/owncloud/web/pull/10551): Share sidebar icons
  * Bugfix [owncloud/web#10702](https://github.com/owncloud/web/pull/10702): Apply sandbox attribute to iframe in draw-io extension
  * Bugfix [owncloud/web#10706](https://github.com/owncloud/web/pull/10706): Apply sandbox attribute to iframe in app-external extension
  * Bugfix [owncloud/web#10746](https://github.com/owncloud/web/pull/10746): Versions loaded multiple times when opening sidebar
  * Bugfix [owncloud/web#10760](https://github.com/owncloud/web/pull/10760): Incoming notifications broken while notification center is open
  * Bugfix [owncloud/web#10814](https://github.com/owncloud/web/issues/10814): Vertical scroll for OcModal on small screens
  * Bugfix [owncloud/web#10900](https://github.com/owncloud/web/pull/10900): Context menu empty in tiles view
  * Bugfix [owncloud/web#10918](https://github.com/owncloud/web/issues/10918): Resource deselection on right-click
  * Bugfix [owncloud/web#10920](https://github.com/owncloud/web/pull/10920): Resources with name consist of number won't show up in trash bin
  * Bugfix [owncloud/web#10928](https://github.com/owncloud/web/pull/10928): Disable search in public link context
  * Bugfix [owncloud/web#10941](https://github.com/owncloud/web/issues/10941): Space not updating on navigation
  * Bugfix [owncloud/web#10974](https://github.com/owncloud/web/pull/10974): Local logout if IdP has no logout support
  * Change [owncloud/web#7338](https://github.com/owncloud/web/issues/7338): Remove deprecated code
  * Change [owncloud/web#9892](https://github.com/owncloud/web/issues/9892): Remove skeleton app
  * Change [owncloud/web#10102](https://github.com/owncloud/web/pull/10102): Remove deprecated extension point for adding quick actions
  * Change [owncloud/web#10122](https://github.com/owncloud/web/pull/10122): Remove homeFolder option
  * Change [owncloud/web#10210](https://github.com/owncloud/web/issues/10210): Vuex store removed
  * Change [owncloud/web#10240](https://github.com/owncloud/web/pull/10240): Remove ocs user
  * Change [owncloud/web#10330](https://github.com/owncloud/web/pull/10330): Registering app file editors
  * Change [owncloud/web#10443](https://github.com/owncloud/web/pull/10443): Add extensionPoint concept
  * Change [owncloud/web#10758](https://github.com/owncloud/web/pull/10758): Portal target removed
  * Change [owncloud/web#10786](https://github.com/owncloud/web/pull/10786): Disable opening files in embed mode
  * Enhancement [owncloud/web#5383](https://github.com/owncloud/web/issues/5383): Accessibility improvements
  * Enhancement [owncloud/web#9215](https://github.com/owncloud/web/issues/9215): Icon for .dcm files
  * Enhancement [owncloud/web#10018](https://github.com/owncloud/web/issues/10018): Tile sizes
  * Enhancement [owncloud/web#10207](https://github.com/owncloud/web/pull/10207): Enable user preferences in public links
  * Enhancement [owncloud/web#10334](https://github.com/owncloud/web/pull/10334): Move ThemeSwitcher into Account Settings
  * Enhancement [owncloud/web#10383](https://github.com/owncloud/web/issues/10383): Top loading bar increase visibility
  * Enhancement [owncloud/web#10390](https://github.com/owncloud/web/pull/10390): Integrate ToastUI editor in the text editor app
  * Enhancement [owncloud/web#10443](https://github.com/owncloud/web/pull/10443): Custom component extension type
  * Enhancement [owncloud/web#10448](https://github.com/owncloud/web/pull/10448): Epub reader app
  * Enhancement [owncloud/web#10485](https://github.com/owncloud/web/pull/10485): Highlight search term in sharing autosuggest list
  * Enhancement [owncloud/web#10519](https://github.com/owncloud/web/pull/10519): Warn user before closing browser when upload is in progress
  * Enhancement [owncloud/web#10534](https://github.com/owncloud/web/issues/10534): Full text search default
  * Enhancement [owncloud/web#10544](https://github.com/owncloud/web/pull/10544): Show locked and processing next to other status indicators
  * Enhancement [owncloud/web#10546](https://github.com/owncloud/web/pull/10546): Set emoji as space icon
  * Enhancement [owncloud/web#10586](https://github.com/owncloud/web/pull/10586): Add SSE events for locking, renaming, deleting, and restoring
  * Enhancement [owncloud/web#10611](https://github.com/owncloud/web/pull/10611): Remember left nav bar state
  * Enhancement [owncloud/web#10612](https://github.com/owncloud/web/pull/10612): Remember right side bar state
  * Enhancement [owncloud/web#10624](https://github.com/owncloud/web/pull/10624): Add details panel to trash
  * Enhancement [owncloud/web#10709](https://github.com/owncloud/web/pull/10709): Implement Server-Sent Events (SSE) for File Creation
  * Enhancement [owncloud/web#10758](https://github.com/owncloud/web/pull/10758): Search providers extension point
  * Enhancement [owncloud/web#10782](https://github.com/owncloud/web/pull/10782): Implement Server-Sent Events (SSE) for file updates
  * Enhancement [owncloud/web#10798](https://github.com/owncloud/web/pull/10798): Add SSE event for moving
  * Enhancement [owncloud/web#10801](https://github.com/owncloud/web/pull/10801): Ability to theme sharing role icons
  * Enhancement [owncloud/web#10807](https://github.com/owncloud/web/pull/10807): Add SSE event for moving
  * Enhancement [owncloud/web#10874](https://github.com/owncloud/web/pull/10874): Show loading spinner while searching or filtering users
  * Enhancement [owncloud/web#10907](https://github.com/owncloud/web/pull/10907): Display hidden resources information in files list
  * Enhancement [owncloud/web#10929](https://github.com/owncloud/web/pull/10929): Add loading spinner to admin settings spaces and groups
  * Enhancement [owncloud/web#10956](https://github.com/owncloud/web/pull/10956): Audio metadata panel
  * Enhancement [owncloud/web#10956](https://github.com/owncloud/web/pull/10956): EXIF metadata panel
  * Enhancement [owncloud/web#10976](https://github.com/owncloud/web/pull/10976): Faster page loading times
  * Enhancement [owncloud/web#11004](https://github.com/owncloud/web/pull/11004): Add enabled only filter to spaces overview
  * Enhancement [owncloud/web#11037](https://github.com/owncloud/web/pull/11037): Multiple sidebar root panels

   https://github.com/owncloud/ocis/pull/9395
   https://github.com/owncloud/web/releases/tag/v9.0.0

* Enhancement - Bump Reva to v2.20.0: [#9415](https://github.com/owncloud/ocis/pull/9415)

  *   Bugfix [cs3org/reva#4623](https://github.com/cs3org/reva/pull/4623): Consistently use spaceid and nodeid in logs
  *   Bugfix [cs3org/reva#4584](https://github.com/cs3org/reva/pull/4584): Prevent copying a file to a parent folder
  *   Bugfix [cs3org/reva#4700](https://github.com/cs3org/reva/pull/4700): Clean empty trash node path on delete
  *   Bugfix [cs3org/reva#4567](https://github.com/cs3org/reva/pull/4567): Fix error message in authprovider if user is not found
  *   Bugfix [cs3org/reva#4615](https://github.com/cs3org/reva/pull/4615): Write blob based on session id
  *   Bugfix [cs3org/reva#4557](https://github.com/cs3org/reva/pull/4557): Fix ceph build
  *   Bugfix [cs3org/reva#4711](https://github.com/cs3org/reva/pull/4711): Duplicate headers in DAV responses
  *   Bugfix [cs3org/reva#4568](https://github.com/cs3org/reva/pull/4568): Fix sharing invite on virtual drive
  *   Bugfix [cs3org/reva#4559](https://github.com/cs3org/reva/pull/4559): Fix graph drive invite
  *   Bugfix [cs3org/reva#4593](https://github.com/cs3org/reva/pull/4593): Make initiatorIDs also work on uploads
  *   Bugfix [cs3org/reva#4608](https://github.com/cs3org/reva/pull/4608): Use gateway selector in jsoncs3
  *   Bugfix [cs3org/reva#4546](https://github.com/cs3org/reva/pull/4546): Fix the mount points naming
  *   Bugfix [cs3org/reva#4678](https://github.com/cs3org/reva/pull/4678): Fix nats encoding
  *   Bugfix [cs3org/reva#4630](https://github.com/cs3org/reva/pull/4630): Fix ocm-share-id
  *   Bugfix [cs3org/reva#4518](https://github.com/cs3org/reva/pull/4518): Fix an error when lock/unlock a file
  *   Bugfix [cs3org/reva#4622](https://github.com/cs3org/reva/pull/4622): Fix public share update
  *   Bugfix [cs3org/reva#4566](https://github.com/cs3org/reva/pull/4566): Fix public link previews
  *   Bugfix [cs3org/reva#4589](https://github.com/cs3org/reva/pull/4589): Fix uploading via a public link
  *   Bugfix [cs3org/reva#4660](https://github.com/cs3org/reva/pull/4660): Fix creating documents in nested folders of public shares
  *   Bugfix [cs3org/reva#4635](https://github.com/cs3org/reva/pull/4635): Fix nil pointer when removing groups from space
  *   Bugfix [cs3org/reva#4709](https://github.com/cs3org/reva/pull/4709): Fix share update
  *   Bugfix [cs3org/reva#4661](https://github.com/cs3org/reva/pull/4661): Fix space share update for ocs
  *   Bugfix [cs3org/reva#4656](https://github.com/cs3org/reva/pull/4656): Fix space share update
  *   Bugfix [cs3org/reva#4561](https://github.com/cs3org/reva/pull/4561): Fix Stat() by Path on re-created resource
  *   Bugfix [cs3org/reva#4710](https://github.com/cs3org/reva/pull/4710): Tolerate missing user space index
  *   Bugfix [cs3org/reva#4632](https://github.com/cs3org/reva/pull/4632): Fix access to files withing a public link targeting a space root
  *   Bugfix [cs3org/reva#4603](https://github.com/cs3org/reva/pull/4603): Mask user email in output
  *   Change [cs3org/reva#4542](https://github.com/cs3org/reva/pull/4542): Drop unused service spanning stat cache
  *   Enhancement [cs3org/reva#4712](https://github.com/cs3org/reva/pull/4712): Add the error translation to the utils
  *   Enhancement [cs3org/reva#4696](https://github.com/cs3org/reva/pull/4696): Add List method to ocis and s3ng blobstore
  *   Enhancement [cs3org/reva#4693](https://github.com/cs3org/reva/pull/4693): Add mimetype for sb3 files
  *   Enhancement [cs3org/reva#4699](https://github.com/cs3org/reva/pull/4699): Add a Path method to blobstore
  *   Enhancement [cs3org/reva#4695](https://github.com/cs3org/reva/pull/4695): Add photo and image props
  *   Enhancement [cs3org/reva#4706](https://github.com/cs3org/reva/pull/4706): Add secureview flag when listing apps via http
  *   Enhancement [cs3org/reva#4585](https://github.com/cs3org/reva/pull/4585): Move more consistency checks to the usershare API
  *   Enhancement [cs3org/reva#4702](https://github.com/cs3org/reva/pull/4702): Added theme capability
  *   Enhancement [cs3org/reva#4672](https://github.com/cs3org/reva/pull/4672): Add virus filter to list uploads sessions
  *   Enhancement [cs3org/reva#4614](https://github.com/cs3org/reva/pull/4614): Bump mockery to v2.40.2
  *   Enhancement [cs3org/reva#4621](https://github.com/cs3org/reva/pull/4621): Use a memory cache for the personal space creation cache
  *   Enhancement [cs3org/reva#4556](https://github.com/cs3org/reva/pull/4556): Allow tracing requests by giving util functions a context
  *   Enhancement [cs3org/reva#4694](https://github.com/cs3org/reva/pull/4694): Expose SecureView in WebDAV permissions
  *   Enhancement [cs3org/reva#4652](https://github.com/cs3org/reva/pull/4652): Better error codes when removing a space member
  *   Enhancement [cs3org/reva#4725](https://github.com/cs3org/reva/pull/4725): Unique share mountpoint name
  *   Enhancement [cs3org/reva#4689](https://github.com/cs3org/reva/pull/4689): Extend service account permissions
  *   Enhancement [cs3org/reva#4545](https://github.com/cs3org/reva/pull/4545): Extend service account permissions
  *   Enhancement [cs3org/reva#4581](https://github.com/cs3org/reva/pull/4581): Make decomposedfs more extensible
  *   Enhancement [cs3org/reva#4564](https://github.com/cs3org/reva/pull/4564): Send file locked/unlocked events
  *   Enhancement [cs3org/reva#4730](https://github.com/cs3org/reva/pull/4730): Improve posixfs storage driver
  *   Enhancement [cs3org/reva#4587](https://github.com/cs3org/reva/pull/4587): Allow passing a initiator id
  *   Enhancement [cs3org/reva#4645](https://github.com/cs3org/reva/pull/4645): Add ItemID to LinkRemoved
  *   Enhancement [cs3org/reva#4686](https://github.com/cs3org/reva/pull/4686): Mint view only token for open in app requests
  *   Enhancement [cs3org/reva#4606](https://github.com/cs3org/reva/pull/4606): Remove resharing
  *   Enhancement [cs3org/reva#4643](https://github.com/cs3org/reva/pull/4643): Secure viewer share role
  *   Enhancement [cs3org/reva#4631](https://github.com/cs3org/reva/pull/4631): Add space-share-updated event
  *   Enhancement [cs3org/reva#4685](https://github.com/cs3org/reva/pull/4685): Support t and x in ACEs
  *   Enhancement [cs3org/reva#4625](https://github.com/cs3org/reva/pull/4625): Test async processing cornercases
  *   Enhancement [cs3org/reva#4653](https://github.com/cs3org/reva/pull/4653): Allow to resolve public shares without the ocs tokeninfo endpoint
  *   Enhancement [cs3org/reva#4657](https://github.com/cs3org/reva/pull/4657): Add ScanData to Uploadsession

   https://github.com/owncloud/ocis/pull/9415
   https://github.com/owncloud/ocis/pull/9377
   https://github.com/owncloud/ocis/pull/9330
   https://github.com/owncloud/ocis/pull/9318
   https://github.com/owncloud/ocis/pull/9269
   https://github.com/owncloud/ocis/pull/9236
   https://github.com/owncloud/ocis/pull/9188
   https://github.com/owncloud/ocis/pull/9132
   https://github.com/owncloud/ocis/pull/9041
   https://github.com/owncloud/ocis/pull/9002
   https://github.com/owncloud/ocis/pull/8917
   https://github.com/owncloud/ocis/pull/8795
   https://github.com/owncloud/ocis/pull/8701
   https://github.com/owncloud/ocis/pull/8606
   https://github.com/owncloud/ocis/pull/8937

# Changelog for [5.0.5] (2024-05-22)

The following sections list the changes for 5.0.5.

[5.0.5]: https://github.com/owncloud/ocis/compare/v5.0.4...v5.0.5

## Summary

* Enhancement - Update web to v8.0.2: [#9153](https://github.com/owncloud/ocis/pull/9153)

## Details

* Enhancement - Update web to v8.0.2: [#9153](https://github.com/owncloud/ocis/pull/9153)

   Tags: web

   We updated ownCloud Web to v8.0.2. Please refer to the changelog (linked) for
   details on the web release.

  * Bugfix [owncloud/web#10515](https://github.com/owncloud/web/issues/10515): Folder replace
  * Bugfix [owncloud/web#10598](https://github.com/owncloud/web/issues/10598): Hidden right sidebar on small screens
  * Bugfix [owncloud/web#10634](https://github.com/owncloud/web/issues/10634): Scope loss when showing search results
  * Bugfix [owncloud/web#10657](https://github.com/owncloud/web/issues/10657): Theme loading without matching theme
  * Bugfix [owncloud/web#10763](https://github.com/owncloud/web/pull/10763): Flickering loading indicator
  * Bugfix [owncloud/web#10810](https://github.com/owncloud/web/issues/10810): Download files with special chars in name
  * Bugfix [owncloud/web#10881](https://github.com/owncloud/web/pull/10881): IDP logout issues

   https://github.com/owncloud/ocis/pull/9153
   https://github.com/owncloud/web/releases/tag/v8.0.2

# Changelog for [5.0.4] (2024-05-13)

The following sections list the changes for 5.0.4.

[5.0.4]: https://github.com/owncloud/ocis/compare/v5.0.3...v5.0.4

## Summary

* Bugfix - Update reva to v2.19.7: [#9011](https://github.com/owncloud/ocis/pull/9011)
* Bugfix - Service startup of WOPI example: [#9127](https://github.com/owncloud/ocis/pull/9127)
* Bugfix - Nats reconnects: [#9139](https://github.com/owncloud/ocis/pull/9139)

## Details

* Bugfix - Update reva to v2.19.7: [#9011](https://github.com/owncloud/ocis/pull/9011)

   We updated reva to v2.19.7

  *   Enhancement [cs3org/reva#4673](https://github.com/cs3org/reva/pull/4673): Add virus filter to list uploads sessions

   https://github.com/owncloud/ocis/pull/9011

* Bugfix - Service startup of WOPI example: [#9127](https://github.com/owncloud/ocis/pull/9127)

   We fixed a bug in the service startup of the appprovider-onlyoffice in the
   ocis_wopi deployment example.

   https://github.com/owncloud/ocis/pull/9127

* Bugfix - Nats reconnects: [#9139](https://github.com/owncloud/ocis/pull/9139)

   We fixed the reconnect handling of the natjs kv registry.

   https://github.com/owncloud/ocis/pull/9139
   https://github.com/owncloud/ocis/pull/8880

# Changelog for [5.0.3] (2024-05-02)

The following sections list the changes for 5.0.3.

[5.0.3]: https://github.com/owncloud/ocis/compare/v5.0.2...v5.0.3

## Summary

* Bugfix - Update the admin user role assignment to enforce the config: [#8918](https://github.com/owncloud/ocis/pull/8918)
* Bugfix - Crash when processing crafted TIFF files: [#8981](https://github.com/owncloud/ocis/pull/8981)
* Bugfix - Update reva to v2.19.6: [#9011](https://github.com/owncloud/ocis/pull/9011)
* Bugfix - Fix infected file handling: [#9011](https://github.com/owncloud/ocis/pull/9011)

## Details

* Bugfix - Update the admin user role assignment to enforce the config: [#8918](https://github.com/owncloud/ocis/pull/8918)

   The admin user role assigment was not updated after the first assignment. We now
   read the assigned role during init and update the admin user ID accordingly if
   the role is not assigned. This is especially needed when the OCIS_ADMIN_USER_ID
   is set after the autoprovisioning of the admin user when it originates from an
   external Identity Provider.

   https://github.com/owncloud/ocis/pull/8918
   https://github.com/owncloud/ocis/pull/8897

* Bugfix - Crash when processing crafted TIFF files: [#8981](https://github.com/owncloud/ocis/pull/8981)

   Fix for a vulnerability with low severity in disintegration/imaging.

   https://github.com/owncloud/ocis/pull/8981
   https://github.com/advisories/GHSA-q7pp-wcgr-pffx

* Bugfix - Update reva to v2.19.6: [#9011](https://github.com/owncloud/ocis/pull/9011)

   We updated reva to v2.19.6

  *   Bugfix      [cs3org/reva#4654](https://github.com/cs3org/reva/pull/4654): Write blob based on session id
  *   Bugfix      [cs3org/reva#4666](https://github.com/cs3org/reva/pull/4666): Fix uploading via a public link
  *   Bugfix      [cs3org/reva#4665](https://github.com/cs3org/reva/pull/4665): Fix creating documents in nested folders of public shares
  *   Enhancement [cs3org/reva#4655](https://github.com/cs3org/reva/pull/4655): Bump mockery to v2.40.2
  *   Enhancement [cs3org/reva#4664](https://github.com/cs3org/reva/pull/4664): Add ScanData to Uploadsession

   https://github.com/owncloud/ocis/pull/9011

* Bugfix - Fix infected file handling: [#9011](https://github.com/owncloud/ocis/pull/9011)

   Reworks virus handling. Shows scandate and outcome on ocis storage-users uploads
   sessions. Avoids retrying infected files on ocis postprocessing restart.

   https://github.com/owncloud/ocis/pull/9011

# Changelog for [5.0.2] (2024-04-17)

The following sections list the changes for 5.0.2.

[5.0.2]: https://github.com/owncloud/ocis/compare/v5.0.1...v5.0.2

## Summary

* Bugfix - Fix creating new WOPI documents on public shares: [#8828](https://github.com/owncloud/ocis/pull/8828)
* Bugfix - Update reva to v2.19.5: [#8873](https://github.com/owncloud/ocis/pull/8873)

## Details

* Bugfix - Fix creating new WOPI documents on public shares: [#8828](https://github.com/owncloud/ocis/pull/8828)

   Creating a new Office document in a publicly shared folder is now possible.

   https://github.com/owncloud/ocis/issues/8691
   https://github.com/owncloud/ocis/pull/8828

* Bugfix - Update reva to v2.19.5: [#8873](https://github.com/owncloud/ocis/pull/8873)

   We updated reva to v2.19.5

  *   Bugfix [cs3org/reva#4626](https://github.com/cs3org/reva/pull/4626): Fix public share update
  *   Bugfix [cs3org/reva#4634](https://github.com/cs3org/reva/pull/4634): Fix access to files withing a public link targeting a space root

   https://github.com/owncloud/ocis/pull/8873

# Changelog for [5.0.1] (2024-04-10)

The following sections list the changes for 5.0.1.

[5.0.1]: https://github.com/owncloud/ocis/compare/v4.0.7...v5.0.1

## Summary

* Bugfix - Make IDP cookies same site strict: [#8716](https://github.com/owncloud/ocis/pull/8716)
* Bugfix - Update reva to v2.19.4: [#8781](https://github.com/owncloud/ocis/pull/8781)
* Bugfix - Fix restarting of postprocessing: [#8782](https://github.com/owncloud/ocis/pull/8782)
* Bugfix - Fix the create personal space cache: [#8799](https://github.com/owncloud/ocis/pull/8799)

## Details

* Bugfix - Make IDP cookies same site strict: [#8716](https://github.com/owncloud/ocis/pull/8716)

   To enhance the security of our application and prevent Cross-Site Request
   Forgery (CSRF) attacks, we have updated the SameSite attribute of the build in
   Identity Provider (IDP) cookies to Strict.

   This change restricts the browser from sending these cookies with any cross-site
   requests, thereby limiting the exposure of the user's session to potential
   threats.

   This update does not impact the existing functionality of the application but
   provides an additional layer of security where needed.

   This only affects cookies set by the built-in IDP. Production systems should not
   be affected.

   https://github.com/owncloud/ocis/pull/8716

* Bugfix - Update reva to v2.19.4: [#8781](https://github.com/owncloud/ocis/pull/8781)

   We updated reva to v2.19.4

  *   Bugfix [cs3org/reva#4612](https://github.com/cs3org/reva/pull/4612): Use gateway selector in jsoncs3 to scale the service

   Https://github.com/owncloud/ocis/pull/8787

   We updated reva to v2.19.3

  *   Bugfix[cs3org/reva#4607](https://github.com/cs3org/reva/pull/4607): Mask user email in output

   https://github.com/owncloud/ocis/pull/8781

* Bugfix - Fix restarting of postprocessing: [#8782](https://github.com/owncloud/ocis/pull/8782)

   When an upload is not found, the logic to restart postprocessing was bunked.
   Additionally we extended the upload sessions command to be able to restart the
   uploads without using a second command.

   NOTE: This also includes a breaking fix for the deprecated `ocis storage-users
   uploads list` command

   https://github.com/owncloud/ocis/pull/8782

* Bugfix - Fix the create personal space cache: [#8799](https://github.com/owncloud/ocis/pull/8799)

   We fixed a problem with the config for the create personal space cache which
   resulted in the cache never being used.

   https://github.com/owncloud/ocis/pull/8799

# Changelog for [4.0.7] (2024-03-27)

The following sections list the changes for 4.0.7.

[4.0.7]: https://github.com/owncloud/ocis/compare/v5.0.0...v4.0.7

## Summary

* Bugfix - Update reva to include bugfixes and improvements: [#8718](https://github.com/owncloud/ocis/pull/8718)
* Enhancement - Update to go 1.22: [#8597](https://github.com/owncloud/ocis/pull/8597)

## Details

* Bugfix - Update reva to include bugfixes and improvements: [#8718](https://github.com/owncloud/ocis/pull/8718)

   ## Changelog for reva 2.13.4

  *   Bugfix [cs3org/reva#4398](https://github.com/cs3org/reva/pull/4398): Fix ceph build
  *   Bugfix [cs3org/reva#4396](https://github.com/cs3org/reva/pull/4396): Allow an empty credentials chain in the auth middleware
  *   Bugfix [cs3org/reva#4423](https://github.com/cs3org/reva/pull/4423): Fix disconnected traces
  *   Bugfix [cs3org/reva#4590](https://github.com/cs3org/reva/pull/4590): Fix uploading via a public link
  *   Bugfix [cs3org/reva#4470](https://github.com/cs3org/reva/pull/4470): Keep failed processing status
  *   Enhancement [cs3org/reva#4397](https://github.com/cs3org/reva/pull/4397): Introduce UploadSessionLister interface

   https://github.com/owncloud/ocis/pull/8718

* Enhancement - Update to go 1.22: [#8597](https://github.com/owncloud/ocis/pull/8597)

   We have updated go to version 1.22.

   https://github.com/owncloud/ocis/pull/8597

# Changelog for [5.0.0] (2024-03-18)

The following sections list the changes for 5.0.0.

[5.0.0]: https://github.com/owncloud/ocis/compare/v4.0.6...v5.0.0

## Summary

* Bugfix - Fix wrong compile date: [#6132](https://github.com/owncloud/ocis/pull/6132)
* Bugfix - Fix the kql-bleve search: [#7290](https://github.com/owncloud/ocis/pull/7290)
* Bugfix - Bring back the USERS_LDAP_USER_SCHEMA_ID variable: [#7312](https://github.com/owncloud/ocis/issues/7312)
* Bugfix - Do not reset state of received shares when rebuilding the jsoncs3 index: [#7319](https://github.com/owncloud/ocis/issues/7319)
* Bugfix - Deprecate redundant encryptions settings for notification service: [#7345](https://github.com/owncloud/ocis/issues/7345)
* Bugfix - Check school number for duplicates before adding a school: [#7351](https://github.com/owncloud/ocis/pull/7351)
* Bugfix - Don't reload web config: [#7369](https://github.com/owncloud/ocis/pull/7369)
* Bugfix - Delete outdated userlog events: [#7410](https://github.com/owncloud/ocis/pull/7410)
* Bugfix - Set the mountpoint on auto accept: [#7460](https://github.com/owncloud/ocis/pull/7460)
* Bugfix - Fix default language fallback: [#7465](https://github.com/owncloud/ocis/issues/7465)
* Bugfix - GetUserByClaim fixed for Active Directory: [#7476](https://github.com/owncloud/ocis/pull/7476)
* Bugfix - Fix preview request 500 error when made too early: [#7502](https://github.com/owncloud/ocis/issues/7502)
* Bugfix - Fix 403 in docs pipeline: [#7509](https://github.com/owncloud/ocis/issues/7509)
* Bugfix - Fix the auth service env variable: [#7523](https://github.com/owncloud/ocis/pull/7523)
* Bugfix - Token storage config fixed: [#7528](https://github.com/owncloud/ocis/pull/7528)
* Bugfix - Set existing mountpoint on auto accept: [#7592](https://github.com/owncloud/ocis/pull/7592)
* Bugfix - Return 423 status code on tag create: [#7596](https://github.com/owncloud/ocis/pull/7596)
* Bugfix - Fix libre-graph status codes: [#7678](https://github.com/owncloud/ocis/issues/7678)
* Bugfix - Fix unlock via space API: [#7726](https://github.com/owncloud/ocis/pull/7726)
* Bugfix - Disable DEPTH infinity in PROPFIND: [#7746](https://github.com/owncloud/ocis/pull/7746)
* Bugfix - Fix the tgz mime type: [#7772](https://github.com/owncloud/ocis/pull/7772)
* Bugfix - Fix natsjs cache: [#7790](https://github.com/owncloud/ocis/pull/7790)
* Bugfix - Fix search service start: [#7795](https://github.com/owncloud/ocis/pull/7795)
* Bugfix - Fix search response: [#7815](https://github.com/owncloud/ocis/pull/7815)
* Bugfix - The race conditions in tests: [#7847](https://github.com/owncloud/ocis/pull/7847)
* Bugfix - Do not purge expired upload sessions that are still postprocessing: [#7859](https://github.com/owncloud/ocis/pull/7859)
* Bugfix - Fix the public link update: [#7862](https://github.com/owncloud/ocis/pull/7862)
* Bugfix - Fix jwt config of policies service: [#7893](https://github.com/owncloud/ocis/pull/7893)
* Bugfix - Updating logo with new theme structure: [#7930](https://github.com/owncloud/ocis/pull/7930)
* Bugfix - Password policy return code was wrong: [#7952](https://github.com/owncloud/ocis/pull/7952)
* Bugfix - Removed outdated and unused dependency from idp package: [#7957](https://github.com/owncloud/ocis/issues/7957)
* Bugfix - Update permission validation: [#7963](https://github.com/owncloud/ocis/pull/7963)
* Bugfix - Renaming a user to a string with capital letters: [#7964](https://github.com/owncloud/ocis/pull/7964)
* Bugfix - Improve OCM support: [#7973](https://github.com/owncloud/ocis/pull/7973)
* Bugfix - Permissions of a role with duplicate ID: [#7976](https://github.com/owncloud/ocis/pull/7976)
* Bugfix - Non durable streams for sse service: [#7986](https://github.com/owncloud/ocis/pull/7986)
* Bugfix - Fix empty trace ids: [#8023](https://github.com/owncloud/ocis/pull/8023)
* Bugfix - Fix search by containing special characters: [#8050](https://github.com/owncloud/ocis/pull/8050)
* Bugfix - Fix the upload postprocessing: [#8117](https://github.com/owncloud/ocis/pull/8117)
* Bugfix - Disallow to delete a file during the processing: [#8132](https://github.com/owncloud/ocis/pull/8132)
* Bugfix - Fix wrong naming in nats-js-kv registry: [#8140](https://github.com/owncloud/ocis/pull/8140)
* Bugfix - IDP CS3 backend sessions now survive a restart: [#8142](https://github.com/owncloud/ocis/pull/8142)
* Bugfix - Fix patching of language: [#8182](https://github.com/owncloud/ocis/pull/8182)
* Bugfix - Fix search service to not log expected cases as errors: [#8200](https://github.com/owncloud/ocis/pull/8200)
* Bugfix - Updating and reset logo failed: [#8211](https://github.com/owncloud/ocis/pull/8211)
* Bugfix - Cleanup graph/pkg/service/v0/driveitems.go: [#8228](https://github.com/owncloud/ocis/pull/8228)
* Bugfix - Cleanup `search/pkg/search/search.go`: [#8230](https://github.com/owncloud/ocis/pull/8230)
* Bugfix - Graph/sharedWithMe works for shares from project spaces now: [#8233](https://github.com/owncloud/ocis/pull/8233)
* Bugfix - Fix PATCH/DELETE status code for drives that don't support them: [#8235](https://github.com/owncloud/ocis/pull/8235)
* Bugfix - Fix nats authentication: [#8236](https://github.com/owncloud/ocis/pull/8236)
* Bugfix - Fix the resource name: [#8246](https://github.com/owncloud/ocis/pull/8246)
* Bugfix - Apply role constraints when creating shares via the graph API: [#8247](https://github.com/owncloud/ocis/pull/8247)
* Bugfix - Fix concurrent access to a map: [#8269](https://github.com/owncloud/ocis/pull/8269)
* Bugfix - Fix nats registry: [#8281](https://github.com/owncloud/ocis/pull/8281)
* Bugfix - Remove invalid environment variables: [#8303](https://github.com/owncloud/ocis/pull/8303)
* Bugfix - Fix concurrent shares config: [#8317](https://github.com/owncloud/ocis/pull/8317)
* Bugfix - Fix Content-Disposition header for downloads: [#8381](https://github.com/owncloud/ocis/pull/8381)
* Bugfix - Signed url verification: [#8385](https://github.com/owncloud/ocis/pull/8385)
* Bugfix - Fix an error when move: [#8396](https://github.com/owncloud/ocis/pull/8396)
* Bugfix - Fix extended env parser: [#8409](https://github.com/owncloud/ocis/pull/8409)
* Bugfix - Graph/drives/permission Expiration date update: [#8413](https://github.com/owncloud/ocis/pull/8413)
* Bugfix - Fix search error message: [#8444](https://github.com/owncloud/ocis/pull/8444)
* Bugfix - Graph/sharedWithMe align IDs with webdav response: [#8467](https://github.com/owncloud/ocis/pull/8467)
* Bugfix - Fix an error when lock/unlock a public shared file: [#8472](https://github.com/owncloud/ocis/pull/8472)
* Bugfix - Bump reva to pull in changes to fix recursive trashcan purge: [#8505](https://github.com/owncloud/ocis/pull/8505)
* Bugfix - Fix remove/update share permissions: [#8529](https://github.com/owncloud/ocis/pull/8529)
* Bugfix - Fix graph drive invite: [#8538](https://github.com/owncloud/ocis/pull/8538)
* Bugfix - We now always select the next clients when autoaccepting shares: [#8570](https://github.com/owncloud/ocis/pull/8570)
* Bugfix - Correct the default mapping of roles: [#8639](https://github.com/owncloud/ocis/pull/8639)
* Bugfix - Disable Multipart uploads: [#8667](https://github.com/owncloud/ocis/pull/8667)
* Bugfix - Fix last month search: [#31145](https://github.com/golang/go/issues/31145)
* Change - Auto-Accept Shares: [#7097](https://github.com/owncloud/ocis/pull/7097)
* Change - Change the default TUS chunk size: [#7273](https://github.com/owncloud/ocis/pull/7273)
* Change - Remove privacyURL and imprintURL from the config: [#7938](https://github.com/owncloud/ocis/pull/7938/)
* Change - Remove accessDeniedHelpUrl from the config: [#7970](https://github.com/owncloud/ocis/pull/7970)
* Change - Change the default store for presigned keys to nats-js-kv: [#8419](https://github.com/owncloud/ocis/pull/8419)
* Change - Deprecate sharing cs3 backends: [#8478](https://github.com/owncloud/ocis/pull/8478)
* Enhancement - Add the Banned Passwords List: [#4197](https://github.com/cs3org/reva/pull/4197)
* Enhancement - Introduce service accounts: [#6427](https://github.com/owncloud/ocis/pull/6427)
* Enhancement - SSE for messaging: [#6992](https://github.com/owncloud/ocis/pull/6992)
* Enhancement - Support spec violating AD FS access token issuer: [#7140](https://github.com/owncloud/ocis/pull/7140)
* Enhancement - Add OCIS_LDAP_BIND_PASSWORD as replacement for LDAP_BIND_PASSWORD: [#7176](https://github.com/owncloud/ocis/issues/7176)
* Enhancement - Keyword Query Language (KQL) search syntax: [#7212](https://github.com/owncloud/ocis/pull/7212)
* Enhancement - Introduce clientlog service: [#7217](https://github.com/owncloud/ocis/pull/7217)
* Enhancement - Proxy uses service accounts for provisioning: [#7240](https://github.com/owncloud/ocis/pull/7240)
* Enhancement - The password policies change request: [#7264](https://github.com/owncloud/ocis/pull/7264)
* Enhancement - Introduce natsjs registry: [#7272](https://github.com/owncloud/ocis/issues/7272)
* Enhancement - Add the password policies: [#7285](https://github.com/owncloud/ocis/pull/7285)
* Enhancement - Add login URL config: [#7317](https://github.com/owncloud/ocis/pull/7317)
* Enhancement - Improve SSE format: [#7325](https://github.com/owncloud/ocis/pull/7325)
* Enhancement - New value `auto` for NOTIFICATIONS_SMTP_AUTHENTICATION: [#7356](https://github.com/owncloud/ocis/issues/7356)
* Enhancement - Make sse service scalable: [#7382](https://github.com/owncloud/ocis/pull/7382)
* Enhancement - Edit wrong named enves: [#7406](https://github.com/owncloud/ocis/pull/7406)
* Enhancement - Thumbnail generation with image processors: [#7409](https://github.com/owncloud/ocis/pull/7409)
* Enhancement - Set default for Async Uploads to true: [#7416](https://github.com/owncloud/ocis/pull/7416)
* Enhancement - The default language added: [#7417](https://github.com/owncloud/ocis/pull/7417)
* Enhancement - Add "Last modified" filter Chip: [#7455](https://github.com/owncloud/ocis/pull/7455)
* Enhancement - Config for disabling Web extensions: [#7486](https://github.com/owncloud/ocis/pull/7486)
* Enhancement - Store and index metadata: [#7490](https://github.com/owncloud/ocis/pull/7490)
* Enhancement - Add support for audio files to the thumbnails service: [#7491](https://github.com/owncloud/ocis/pull/7491)
* Enhancement - Implement sharing roles: [#7524](https://github.com/owncloud/ocis/pull/7524)
* Enhancement - Add new permission to delete public link password: [#7538](https://github.com/owncloud/ocis/pull/7538)
* Enhancement - Add config to enforce passwords on all public links: [#7547](https://github.com/owncloud/ocis/pull/7547)
* Enhancement - Tika content extraction cleanup for search: [#7553](https://github.com/owncloud/ocis/pull/7553)
* Enhancement - Allow configuring storage registry with envvars: [#7554](https://github.com/owncloud/ocis/pull/7554)
* Enhancement - Add search MediaType filter: [#7602](https://github.com/owncloud/ocis/pull/7602)
* Enhancement - Add Sharing NG endpoints: [#7633](https://github.com/owncloud/ocis/pull/7633)
* Enhancement - Configs for Web embed mode: [#7670](https://github.com/owncloud/ocis/pull/7670)
* Enhancement - Support login page background configuration: [#7674](https://github.com/owncloud/ocis/issues/7674)
* Enhancement - Add new permissions: [#7700](https://github.com/owncloud/ocis/pull/7700)
* Enhancement - Add preferred language to user settings: [#7720](https://github.com/owncloud/ocis/pull/7720)
* Enhancement - Add user filter startswith and contains: [#7739](https://github.com/owncloud/ocis/pull/7739)
* Enhancement - Allow configuring additional routes: [#7741](https://github.com/owncloud/ocis/pull/7741)
* Enhancement - Default link permission config: [#7783](https://github.com/owncloud/ocis/pull/7783)
* Enhancement - Add banned password list to the default deployments: [#7784](https://github.com/owncloud/ocis/pull/7784)
* Enhancement - Update to go 1.21: [#7794](https://github.com/owncloud/ocis/pull/7794)
* Enhancement - Add Sharing NG list permissions endpoint: [#7805](https://github.com/owncloud/ocis/pull/7805)
* Enhancement - Add user list requires filter config: [#7866](https://github.com/owncloud/ocis/pull/7866)
* Enhancement - Retry antivirus postprocessing step in case of problems: [#7874](https://github.com/owncloud/ocis/pull/7874)
* Enhancement - Add validation to public share provider: [#7877](https://github.com/owncloud/ocis/pull/7877)
* Enhancement - Graphs endpoint for mounting and unmounting shares: [#7885](https://github.com/owncloud/ocis/pull/7885)
* Enhancement - Store and index metadata: [#7886](https://github.com/owncloud/ocis/pull/7886)
* Enhancement - Allow regular users to list other users: [#7887](https://github.com/owncloud/ocis/pull/7887)
* Enhancement - Add edit public share to sharing NG: [#7908](https://github.com/owncloud/ocis/pull/7908/)
* Enhancement - Add cli commands for trash-bin: [#7917](https://github.com/owncloud/ocis/pull/7917)
* Enhancement - Add validation update public share: [#7978](https://github.com/owncloud/ocis/pull/7978)
* Enhancement - Allow inmemory nats-js-kv stores: [#7979](https://github.com/owncloud/ocis/pull/7979)
* Enhancement - Disable the password policy: [#7985](https://github.com/owncloud/ocis/pull/7985)
* Enhancement - Use kv store in natsjs registry: [#7987](https://github.com/owncloud/ocis/pull/7987)
* Enhancement - Allow authentication nats connections: [#7989](https://github.com/owncloud/ocis/pull/7989)
* Enhancement - Add RED metrics to the metrics endpoint: [#7994](https://github.com/owncloud/ocis/pull/7994)
* Enhancement - Add ocm and sciencemesh services: [#7998](https://github.com/owncloud/ocis/pull/7998)
* Enhancement - Make nats-js-kv the default registry: [#8011](https://github.com/owncloud/ocis/pull/8011)
* Enhancement - Service Account roles: [#8051](https://github.com/owncloud/ocis/pull/8051)
* Enhancement - Update antivirus service: [#8062](https://github.com/owncloud/ocis/pull/8062)
* Enhancement - Remove deprecated environment variables: [#8149](https://github.com/owncloud/ocis/pull/8149)
* Enhancement - Disable the password policy: [#8152](https://github.com/owncloud/ocis/pull/8152)
* Enhancement - Allow restarting multiple uploads with one command: [#8287](https://github.com/owncloud/ocis/pull/8287)
* Enhancement - Modify the concurrency default: [#8309](https://github.com/owncloud/ocis/pull/8309)
* Enhancement - Improve ocis single binary start: [#8320](https://github.com/owncloud/ocis/pull/8320)
* Enhancement - Use environment variables in yaml config files: [#8339](https://github.com/owncloud/ocis/pull/8339)
* Enhancement - Increment filenames on upload collisions in secret filedrops: [#8340](https://github.com/owncloud/ocis/pull/8340)
* Enhancement - Allow sending multiple user ids in one sse event: [#8379](https://github.com/owncloud/ocis/pull/8379)
* Enhancement - Allow to skip service listing: [#8408](https://github.com/owncloud/ocis/pull/8408)
* Enhancement - Add a make step to validate the env var annotations: [#8436](https://github.com/owncloud/ocis/pull/8436)
* Enhancement - Drop the unnecessary grants exists check when creating shares: [#8502](https://github.com/owncloud/ocis/pull/8502)
* Enhancement - Update to go 1.22: [#8586](https://github.com/owncloud/ocis/pull/8586)
* Enhancement - Update web to v8.0.0: [#8613](https://github.com/owncloud/ocis/pull/8613)
* Enhancement - Update web to v8.0.1: [#8626](https://github.com/owncloud/ocis/pull/8626)
* Enhancement - Update reva to 2.19.2: [#8638](https://github.com/owncloud/ocis/pull/8638)

## Details

* Bugfix - Fix wrong compile date: [#6132](https://github.com/owncloud/ocis/pull/6132)

   We fixed that current date is always printed.

   https://github.com/owncloud/ocis/issues/6124
   https://github.com/owncloud/ocis/pull/6132

* Bugfix - Fix the kql-bleve search: [#7290](https://github.com/owncloud/ocis/pull/7290)

   We fixed the issue when 500 on searches that contain ":". Added the characters
   escaping according to https://blevesearch.com/docs/Query-String-Query/

   https://github.com/owncloud/ocis/issues/7282
   https://github.com/owncloud/ocis/pull/7290

* Bugfix - Bring back the USERS_LDAP_USER_SCHEMA_ID variable: [#7312](https://github.com/owncloud/ocis/issues/7312)

   We reintroduced the USERS_LDAP_USER_SCHEMA_ID variable which was accidently
   removed from the users service with the 4.0.0 release.

   https://github.com/owncloud/ocis/issues/7312
   https://github.com/owncloud/ocis-charts/issues/397

* Bugfix - Do not reset state of received shares when rebuilding the jsoncs3 index: [#7319](https://github.com/owncloud/ocis/issues/7319)

   We fixed a problem with the "ocis migrate rebuild-jsoncs3-indexes" command which
   reset the state of received shares to "pending".

   https://github.com/owncloud/ocis/issues/7319

* Bugfix - Deprecate redundant encryptions settings for notification service: [#7345](https://github.com/owncloud/ocis/issues/7345)

   The values `tls` and `ssl` for the `smtp_encryption` configuration setting are
   duplicates of `starttls` and `ssltls`. They have been marked as deprecated. A
   warning will be logged when they are still used. Please use `starttls` instead
   for `tls` and `ssltls` instead of `ssl.

   https://github.com/owncloud/ocis/issues/7345

* Bugfix - Check school number for duplicates before adding a school: [#7351](https://github.com/owncloud/ocis/pull/7351)

   We fixed an issue that allowed to create two schools with the same school number

   https://github.com/owncloud/enterprise/issues/6051
   https://github.com/owncloud/ocis/pull/7351

* Bugfix - Don't reload web config: [#7369](https://github.com/owncloud/ocis/pull/7369)

   When requesting `config.json` file from the server, web service would reload the
   file if a path is set. This will remove config entries set via Envvar. Since we
   want to have the possiblity to set configuration from both sources we removed
   the reading from file. The file will still be loaded on service startup.

   https://github.com/owncloud/ocis/pull/7369

* Bugfix - Delete outdated userlog events: [#7410](https://github.com/owncloud/ocis/pull/7410)

   Userlog will now delete events when the user has no longer access to the
   underlying resource

   https://github.com/owncloud/ocis/pull/7410

* Bugfix - Set the mountpoint on auto accept: [#7460](https://github.com/owncloud/ocis/pull/7460)

   On shares auto accept set a mountpoint with same logic as ocs handler

   https://github.com/owncloud/ocis/pull/7460

* Bugfix - Fix default language fallback: [#7465](https://github.com/owncloud/ocis/issues/7465)

   Add the default language for the webui, the settings, userlog and notification
   service.

   https://github.com/owncloud/ocis/issues/7465

* Bugfix - GetUserByClaim fixed for Active Directory: [#7476](https://github.com/owncloud/ocis/pull/7476)

   The reva ldap backend for the users and groups service did not hex escape binary
   uuids in LDAP filter correctly this could cause problems in Active Directory
   setups for services using the GetUserByClaim CS3 request with claim "userid".

   https://github.com/owncloud/ocis/issues/7469
   https://github.com/owncloud/ocis/pull/7476

* Bugfix - Fix preview request 500 error when made too early: [#7502](https://github.com/owncloud/ocis/issues/7502)

   Fix the status code and message when a thumbnail request is made too early.

   https://github.com/owncloud/ocis/issues/7502
   https://github.com/owncloud/ocis/pull/7507

* Bugfix - Fix 403 in docs pipeline: [#7509](https://github.com/owncloud/ocis/issues/7509)

   Docs pipeline was not routed through our proxies which could lead to requests
   being blacklisted

   https://github.com/owncloud/ocis/issues/7509
   https://github.com/owncloud/ocis/pull/7511

* Bugfix - Fix the auth service env variable: [#7523](https://github.com/owncloud/ocis/pull/7523)

   We the auth service env variable to the service specific name. Before it was
   configurable via `AUTH_MACHINE_JWT_SECRET` and now is configurable via
   `AUTH_SERVICE_JWT_SECRET`.

   https://github.com/owncloud/ocis/pull/7523

* Bugfix - Token storage config fixed: [#7528](https://github.com/owncloud/ocis/pull/7528)

   The token storage config in the config.json for web was missing when it was set
   to `false`.

   https://github.com/owncloud/ocis/issues/7462
   https://github.com/owncloud/ocis/pull/7528

* Bugfix - Set existing mountpoint on auto accept: [#7592](https://github.com/owncloud/ocis/pull/7592)

   When already having a share for a specific resource, auto accept would use
   custom mountpoints which lead to other errors. Now auto-accept is using the
   existing mountpoint of a share.

   https://github.com/owncloud/ocis/pull/7592

* Bugfix - Return 423 status code on tag create: [#7596](https://github.com/owncloud/ocis/pull/7596)

   When a file is locked, return 423 status code instead 500 on tag create

   https://github.com/owncloud/ocis/pull/7596

* Bugfix - Fix libre-graph status codes: [#7678](https://github.com/owncloud/ocis/issues/7678)

   Creating group: https://owncloud.dev/libre-graph-api/#/groups/CreateGroup
   changed: 200 -> 201

   Creating users: https://owncloud.dev/libre-graph-api/#/users/CreateUser changed:
   200 -> 201

   Export GDPR: https://owncloud.dev/libre-graph-api/#/user/ExportPersonalData
   changed: 201 -> 202

   https://github.com/owncloud/ocis/issues/7678
   https://github.com/owncloud/ocis/pull/7705

* Bugfix - Fix unlock via space API: [#7726](https://github.com/owncloud/ocis/pull/7726)

   We fixed a bug that caused Error 500 when user try to unlock file using fileid
   The handleSpaceUnlock has been added

   https://github.com/owncloud/ocis/issues/7708
   https://github.com/owncloud/ocis/pull/7726
   https://github.com/cs3org/reva/pull/4338

* Bugfix - Disable DEPTH infinity in PROPFIND: [#7746](https://github.com/owncloud/ocis/pull/7746)

   We fixed the Disabled DEPTH infinity in PROPFIND for: Personal
   /remote.php/dav/files/admin Public link share
   /remote.php/dav/public-files/<token> Trashbin
   /remote.php/dav/spaces/trash-bin/<personal-space-id>

   https://github.com/owncloud/ocis/issues/7359
   https://github.com/owncloud/ocis/pull/7746
   https://github.com/cs3org/reva/pull/4278

* Bugfix - Fix the tgz mime type: [#7772](https://github.com/owncloud/ocis/pull/7772)

   We have fixed a bug when the tgz mime type was not "application/gzip"

   https://github.com/owncloud/ocis/issues/7744
   https://github.com/owncloud/ocis/pull/7772

* Bugfix - Fix natsjs cache: [#7790](https://github.com/owncloud/ocis/pull/7790)

   The nats-js cache was not working. It paniced and wrote a lot of error logs.
   Both is fixed now.

   https://github.com/owncloud/ocis/pull/7790

* Bugfix - Fix search service start: [#7795](https://github.com/owncloud/ocis/pull/7795)

   The `search` service would sometimes not start correctly because config values
   are overwritten by default configuration.

   https://github.com/owncloud/ocis/pull/7795

* Bugfix - Fix search response: [#7815](https://github.com/owncloud/ocis/pull/7815)

   We fixed the search response code from 500 to 400 when the request is invalid

   https://github.com/owncloud/ocis/issues/7812
   https://github.com/owncloud/ocis/pull/7815

* Bugfix - The race conditions in tests: [#7847](https://github.com/owncloud/ocis/pull/7847)

   We fixed the race conditions in tests.

   https://github.com/owncloud/ocis/issues/7846
   https://github.com/owncloud/ocis/pull/7847

* Bugfix - Do not purge expired upload sessions that are still postprocessing: [#7859](https://github.com/owncloud/ocis/pull/7859)

   https://github.com/owncloud/ocis/pull/7859
   https://github.com/owncloud/ocis/pull/7958

* Bugfix - Fix the public link update: [#7862](https://github.com/owncloud/ocis/pull/7862)

   We fixed a bug when normal users can update the public link to delete its
   password if permission is not sent in data.

   https://github.com/owncloud/ocis/issues/7821
   https://github.com/owncloud/ocis/pull/7862

* Bugfix - Fix jwt config of policies service: [#7893](https://github.com/owncloud/ocis/pull/7893)

   Removes jwt config of policies service

   https://github.com/owncloud/ocis/pull/7893

* Bugfix - Updating logo with new theme structure: [#7930](https://github.com/owncloud/ocis/pull/7930)

   Updating and resetting the logo when using the new `theme.json` structure in Web
   has been fixed.

   https://github.com/owncloud/ocis/pull/7930

* Bugfix - Password policy return code was wrong: [#7952](https://github.com/owncloud/ocis/pull/7952)

   We fixed the status code on SharingNG update permissions for public shares.

   https://github.com/owncloud/ocis/pull/7952

* Bugfix - Removed outdated and unused dependency from idp package: [#7957](https://github.com/owncloud/ocis/issues/7957)

   We've removed the outdated and apparently unused dependency `cldr` from the
   `kpop` dependency inside the idp web ui. This resolves a security issue around
   an oudated `xmldom` package version, originating from said `kpop` library.

   https://github.com/owncloud/ocis/issues/7957
   https://github.com/owncloud/ocis/pull/7988

* Bugfix - Update permission validation: [#7963](https://github.com/owncloud/ocis/pull/7963)

   We fixed a bug where the permission validation was not working correctly.

   https://github.com/owncloud/ocis/pull/7963
   https://github.com/cs3org/reva/pull/4405

* Bugfix - Renaming a user to a string with capital letters: [#7964](https://github.com/owncloud/ocis/pull/7964)

   We fixed the issue that led to correct update but the 404 response code when
   renaming an existing user to a string with capital letters.

   https://github.com/owncloud/ocis/pull/7964

* Bugfix - Improve OCM support: [#7973](https://github.com/owncloud/ocis/pull/7973)

   We improved functionality of the OCM support.

   https://github.com/owncloud/ocis/pull/7973

* Bugfix - Permissions of a role with duplicate ID: [#7976](https://github.com/owncloud/ocis/pull/7976)

   We remove the redundant permissions of a role with duplicate ID.

   https://github.com/owncloud/ocis/issues/7931
   https://github.com/owncloud/ocis/pull/7976

* Bugfix - Non durable streams for sse service: [#7986](https://github.com/owncloud/ocis/pull/7986)

   Configure sse streams to be non-durable. This functionality is not needed for
   the sse service

   https://github.com/owncloud/ocis/pull/7986

* Bugfix - Fix empty trace ids: [#8023](https://github.com/owncloud/ocis/pull/8023)

   We changed the default tracing to produce non-empty traceids.

   https://github.com/owncloud/ocis/pull/8023
   https://github.com/owncloud/ocis/pull/8017

* Bugfix - Fix search by containing special characters: [#8050](https://github.com/owncloud/ocis/pull/8050)

   As the OData query parser interprets characters like '@' or '-' in a special
   way. Search request for users or groups needs to be quoted. We fixed the
   libregraph users and groups endpoints to handle quoted search terms correctly.

   https://github.com/owncloud/ocis/issues/7990
   https://github.com/owncloud/ocis/pull/8050
   https://github.com/owncloud/ocis/pull/8035

* Bugfix - Fix the upload postprocessing: [#8117](https://github.com/owncloud/ocis/pull/8117)

   We fixed the upload postprocessing when the destination file does not exist
   anymore.

   https://github.com/owncloud/ocis/issues/7909
   https://github.com/owncloud/ocis/pull/8117

* Bugfix - Disallow to delete a file during the processing: [#8132](https://github.com/owncloud/ocis/pull/8132)

   We want to disallow deleting a file during the processing to prevent collecting
   the orphan uploads.

   https://github.com/owncloud/ocis/issues/8127
   https://github.com/owncloud/ocis/pull/8132
   https://github.com/cs3org/reva/pull/4446

* Bugfix - Fix wrong naming in nats-js-kv registry: [#8140](https://github.com/owncloud/ocis/pull/8140)

   Registers the registry under the correct name

   https://github.com/owncloud/ocis/pull/8140

* Bugfix - IDP CS3 backend sessions now survive a restart: [#8142](https://github.com/owncloud/ocis/pull/8142)

   We now correctly reinitialize the CS3 backend session after the IDP service has
   been restarted.

   https://github.com/owncloud/ocis/pull/8142

* Bugfix - Fix patching of language: [#8182](https://github.com/owncloud/ocis/pull/8182)

   User would not be able to patch their preferred language when the ldap backend
   is set to `read-only`. This makes no sense as language is stored elsewhere.

   https://github.com/owncloud/ocis/pull/8182

* Bugfix - Fix search service to not log expected cases as errors: [#8200](https://github.com/owncloud/ocis/pull/8200)

   We changed the search service to not log cases where resources that were about
   to be indexed can no longer be found. Those are expected cases, e.g. when the
   file in question has already been deleted or renamed meanwhile.

   https://github.com/owncloud/ocis/pull/8200

* Bugfix - Updating and reset logo failed: [#8211](https://github.com/owncloud/ocis/pull/8211)

   We fixed a bug when admin tried to update or reset the logo.

   https://github.com/owncloud/ocis/issues/8101
   https://github.com/owncloud/ocis/pull/8211

* Bugfix - Cleanup graph/pkg/service/v0/driveitems.go: [#8228](https://github.com/owncloud/ocis/pull/8228)

   Main fix is using proto getters to avoid panics. But some other code
   improvements were also done

   https://github.com/owncloud/ocis/pull/8228

* Bugfix - Cleanup `search/pkg/search/search.go`: [#8230](https://github.com/owncloud/ocis/pull/8230)

   Now uses proto getters to avoid panics.

   https://github.com/owncloud/ocis/pull/8230

* Bugfix - Graph/sharedWithMe works for shares from project spaces now: [#8233](https://github.com/owncloud/ocis/pull/8233)

   We fixed a bug in the 'graph/v1beta1/me/drive/sharedWithMe' endpoint that caused
   an error response when the user received shares from project spaces.
   Additionally the endpoint now behaves more graceful in cases where the
   displayname of the owner or creator of a share or shared resource couldn't be
   resolved.

   https://github.com/owncloud/ocis/issues/8027
   https://github.com/owncloud/ocis/issues/8215
   https://github.com/owncloud/ocis/pull/8233

* Bugfix - Fix PATCH/DELETE status code for drives that don't support them: [#8235](https://github.com/owncloud/ocis/pull/8235)

   Updating and Deleting the virtual drives for shares is currently not supported.
   Instead of returning a generic 500 status we return a 405 response now.

   https://github.com/owncloud/ocis/issues/7881
   https://github.com/owncloud/ocis/pull/8235

* Bugfix - Fix nats authentication: [#8236](https://github.com/owncloud/ocis/pull/8236)

   Fixes nats authentication for registry/events/stores

   https://github.com/owncloud/ocis/pull/8236

* Bugfix - Fix the resource name: [#8246](https://github.com/owncloud/ocis/pull/8246)

   We fixed a problem where after renaming resource as sharer the receiver see a
   new name.

   https://github.com/owncloud/ocis/issues/8242
   https://github.com/owncloud/ocis/pull/8246
   https://github.com/cs3org/reva/pull/4463

* Bugfix - Apply role constraints when creating shares via the graph API: [#8247](https://github.com/owncloud/ocis/pull/8247)

   We fixed a bug in the graph API for creating and updating shares so that
   Spaceroot specific roles like 'Manager' and 'Co-owner' can no longer be assigned
   for shares on files or directories.

   https://github.com/owncloud/ocis/issues/8131
   https://github.com/owncloud/ocis/pull/8247

* Bugfix - Fix concurrent access to a map: [#8269](https://github.com/owncloud/ocis/pull/8269)

   We fixed the race condition that led to concurrent map access in a publicshare
   manager.

   https://github.com/owncloud/ocis/issues/8255
   https://github.com/owncloud/ocis/pull/8269
   https://github.com/cs3org/reva/pull/4472

* Bugfix - Fix nats registry: [#8281](https://github.com/owncloud/ocis/pull/8281)

   The nats registry would behave badly when configuring `nats-js-kv` via envvar.
   Reason is the way go-micro initializes. It took 5 developers to find the issue
   and the fix so the details cannot be shared here. Just accept that it is working
   now

   https://github.com/owncloud/ocis/pull/8281

* Bugfix - Remove invalid environment variables: [#8303](https://github.com/owncloud/ocis/pull/8303)

   We have removed two spaces related environment variables (whether project spaces
   and the share jail are enabled) and hardcoded the only allowed options. Misusing
   those variables would have resulted in invalid config.

   https://github.com/owncloud/ocis/pull/8303

* Bugfix - Fix concurrent shares config: [#8317](https://github.com/owncloud/ocis/pull/8317)

   We fixed setting the config for concurrent web requests, which did not work as
   expected before.

   https://github.com/owncloud/ocis/pull/8317

* Bugfix - Fix Content-Disposition header for downloads: [#8381](https://github.com/owncloud/ocis/pull/8381)

   We have fixed a bug that caused downloads to fail on Chromebased browsers when
   the filename contained special characters.

   https://github.com/owncloud/ocis/issues/8361
   https://github.com/owncloud/ocis/pull/8381
   https://github.com/cs3org/reva/pull/4498

* Bugfix - Signed url verification: [#8385](https://github.com/owncloud/ocis/pull/8385)

   Signed urls now expire properly

   https://github.com/owncloud/ocis/pull/8385

* Bugfix - Fix an error when move: [#8396](https://github.com/owncloud/ocis/pull/8396)

   We fixed a bug that caused Internal Server Error when move using destination id

   https://github.com/owncloud/ocis/issues/6739
   https://github.com/owncloud/ocis/pull/8396
   https://github.com/cs3org/reva/pull/4503

* Bugfix - Fix extended env parser: [#8409](https://github.com/owncloud/ocis/pull/8409)

   The extended envvar parser would be angry if there are two `os.Getenv` in the
   same line. We fixed this.

   https://github.com/owncloud/ocis/pull/8409

* Bugfix - Graph/drives/permission Expiration date update: [#8413](https://github.com/owncloud/ocis/pull/8413)

   We fixed a bug in the Update sharing permission the expiration dates can't be
   removed from link permissions.

   https://github.com/owncloud/ocis/issues/8405
   https://github.com/owncloud/ocis/pull/8413

* Bugfix - Fix search error message: [#8444](https://github.com/owncloud/ocis/pull/8444)

   We fixed an error message returned when the search request is invalid

   https://github.com/owncloud/ocis/issues/8442
   https://github.com/owncloud/ocis/pull/8444

* Bugfix - Graph/sharedWithMe align IDs with webdav response: [#8467](https://github.com/owncloud/ocis/pull/8467)

   The IDs of the driveItems returned by the 'graph/v1beta1/me/drive/sharedWithMe'
   endpoint are now aligned with the IDs returned in the PROPFIND response of the
   webdav service.

   https://github.com/owncloud/ocis/issues/8420
   https://github.com/owncloud/ocis/issues/8080
   https://github.com/owncloud/ocis/pull/8467

* Bugfix - Fix an error when lock/unlock a public shared file: [#8472](https://github.com/owncloud/ocis/pull/8472)

   We fixed a bug when anonymous user with viewer role in public link of a folder
   can lock/unlock a file inside it

   https://github.com/owncloud/ocis/issues/7785
   https://github.com/owncloud/ocis/pull/8472

* Bugfix - Bump reva to pull in changes to fix recursive trashcan purge: [#8505](https://github.com/owncloud/ocis/pull/8505)

   We have fixed a bug in the trashcan purge process that did not delete folder
   structures recursively.

   https://github.com/owncloud/ocis/issues/8473
   https://github.com/owncloud/ocis/pull/8505
   https://github.com/cs3org/reva/pull/4533

* Bugfix - Fix remove/update share permissions: [#8529](https://github.com/owncloud/ocis/pull/8529)

   This is a workaround that should prevent removing or changing the share
   permissions when the file is locked. These limitations have to be removed after
   the wopi server will be able to unlock the file properly. These limitations are
   not spread on the files inside the shared folder.

   https://github.com/owncloud/ocis/issues/8273
   https://github.com/owncloud/ocis/pull/8529
   https://github.com/cs3org/reva/pull/4534

* Bugfix - Fix graph drive invite: [#8538](https://github.com/owncloud/ocis/pull/8538)

   We fixed the issue when sharing of personal drive is allowed via graph

   https://github.com/owncloud/ocis/issues/8494
   https://github.com/owncloud/ocis/pull/8538

* Bugfix - We now always select the next clients when autoaccepting shares: [#8570](https://github.com/owncloud/ocis/pull/8570)

   https://github.com/owncloud/ocis/pull/8570

* Bugfix - Correct the default mapping of roles: [#8639](https://github.com/owncloud/ocis/pull/8639)

   The default config for the OIDC role mapping was incorrect. Lightweight users
   are now assignable.

   https://github.com/owncloud/ocis/pull/8639

* Bugfix - Disable Multipart uploads: [#8667](https://github.com/owncloud/ocis/pull/8667)

   Disables multiparts uploads as they lead to high memory consumption

   https://github.com/owncloud/ocis/pull/8667

* Bugfix - Fix last month search: [#31145](https://github.com/golang/go/issues/31145)

   We've fixed the last month search edge case when currently is 31-th.

   Https://github.com/owncloud/ocis/issues/7629
   https://github.com/owncloud/ocis/pull/7742

   https://github.com/golang/go/issues/31145
   The
   issue
   is
   related
   to
   the
   build-in
   package
   behavior

* Change - Auto-Accept Shares: [#7097](https://github.com/owncloud/ocis/pull/7097)

   Automatically accepts shares. This feature is active by default and can be
   deactivated via the environment variable `FRONTEND_AUTO_ACCEPT_SHARES`.

   https://github.com/owncloud/ocis/pull/7097

* Change - Change the default TUS chunk size: [#7273](https://github.com/owncloud/ocis/pull/7273)

   We changed the default TUS chunk size from 100MB to 10MB. You can still use the
   old value by configuring it in your deployment.

   https://github.com/owncloud/ocis/pull/7273

* Change - Remove privacyURL and imprintURL from the config: [#7938](https://github.com/owncloud/ocis/pull/7938/)

   We've removed the option privacyURL and imprintURL from the config, since other
   clients weren't able to consume these. In order to be accessible by other
   clients, not just Web, those should be configured via the theme.json file.

   https://github.com/owncloud/ocis/pull/7938/

* Change - Remove accessDeniedHelpUrl from the config: [#7970](https://github.com/owncloud/ocis/pull/7970)

   We've removed the option accessDeniedHelpUrl from the config, since other
   clients weren't able to consume it. In order to be accessible by other clients,
   not just Web, it should be configured via the theme.json file.

   https://github.com/owncloud/ocis/pull/7970

* Change - Change the default store for presigned keys to nats-js-kv: [#8419](https://github.com/owncloud/ocis/pull/8419)

   We wrapped the store service in a micro store implementation and changed the
   default to the built-in NATS instance.

   https://github.com/owncloud/ocis/pull/8419

* Change - Deprecate sharing cs3 backends: [#8478](https://github.com/owncloud/ocis/pull/8478)

   The `cs3` user and public sharing drivers have already been replaced by
   `jsoncs3`. We now mark them as deprecated in preparation to kill a lot of unused
   code in reva.

   https://github.com/owncloud/ocis/pull/8478

* Enhancement - Add the Banned Passwords List: [#4197](https://github.com/cs3org/reva/pull/4197)

   Added an option to enable a password check against a banned passwords list
   OCIS-3809

   https://github.com/cs3org/reva/pull/4197
   https://github.com/owncloud/ocis/pull/7314

* Enhancement - Introduce service accounts: [#6427](https://github.com/owncloud/ocis/pull/6427)

   Introduces service accounts to avoid impersonating users in async processes

   https://github.com/owncloud/ocis/issues/5550
   https://github.com/owncloud/ocis/pull/6427

* Enhancement - SSE for messaging: [#6992](https://github.com/owncloud/ocis/pull/6992)

   So far, sse has only been used to exchange messages between the server and the
   client. In order to be able to send more content to the client, we have moved
   the endpoint to a separate service and are now also using it for other
   notifications like:

  * notify postprocessing state changes.
  * notify file locking and unlocking.

   https://github.com/owncloud/ocis/pull/6992

* Enhancement - Support spec violating AD FS access token issuer: [#7140](https://github.com/owncloud/ocis/pull/7140)

   AD FS `/adfs/.well-known/openid-configuration` has an optional
   `access_token_issuer` which, in violation of the OpenID Connect spec, takes
   precedence over `issuer`.

   https://github.com/owncloud/ocis/pull/7140

* Enhancement - Add OCIS_LDAP_BIND_PASSWORD as replacement for LDAP_BIND_PASSWORD: [#7176](https://github.com/owncloud/ocis/issues/7176)

   The enviroment variable `OCIS_LDAP_BIND_PASSWORD` was added to be more
   consistent with all other global LDAP variables.

   `LDAP_BIND_PASSWORD` is deprecated now and scheduled for removal with the 5.0.0
   release.

   We also deprecated `LDAP_USER_SCHEMA_ID_IS_OCTETSTRING` for removal with 5.0.0.
   The replacement for it is `OCIS_LDAP_USER_SCHEMA_ID_IS_OCTETSTRING`.

   https://github.com/owncloud/ocis/issues/7176

* Enhancement - Keyword Query Language (KQL) search syntax: [#7212](https://github.com/owncloud/ocis/pull/7212)

   We've introduced support for
   [KQL](https://learn.microsoft.com/en-us/sharepoint/dev/general-development/keyword-query-language-kql-syntax-reference)
   as the default oCIS search query language.

   Simple queries:

  * `tag:golden tag:"silver"`
  * `name:file.txt name:"file.docx"`
  * `content:ahab content:"captain aha*"`

   Date/-range queries

  * `Mtime:"2023-09-05T08:42:11.23554+02:00"`
  * `Mtime>"2023-09-05T08:42:11.23554+02:00"`
  * `Mtime>="2023-09-05T08:42:11.23554+02:00"`
  * `Mtime<"2023-09-05T08:42:11.23554+02:00"`
  * `Mtime<="2023-09-05T08:42:11.23554+02:00"`
  * `Mtime:today` - range: start of today till end of today
  * `Mtime:yesterday` - range: start of yesterday till end of yesterday
  * `Mtime:"this week"` - range: start of this week till end of this week
  * `Mtime:"this month"` - range: start of this month till end of this month
  * `Mtime:"last month"` - range: start of last month till end of last month
  * `Mtime:"this year"` - range: start of this year till end of this year
  * `Mtime:"last year"` - range: start of last year till end of last year

   Conjunctive normal form queries:

  * `tag:golden AND tag:"silver`, `tag:golden OR tag:"silver`, `tag:golden NOT tag:"silver`
  * `(tag:book content:ahab*)`, `tag:(book pdf)`

   Complex queries:

  * `(name:"moby di*" OR tag:bestseller) AND tag:book NOT tag:read`

   https://github.com/owncloud/ocis/issues/7042
   https://github.com/owncloud/ocis/issues/7179
   https://github.com/owncloud/ocis/issues/7114
   https://github.com/owncloud/web/issues/9636
   https://github.com/owncloud/web/issues/9646
   https://github.com/owncloud/ocis/pull/7212
   https://github.com/owncloud/ocis/pull/7043
   https://github.com/owncloud/ocis/pull/7247
   https://github.com/owncloud/ocis/pull/7248
   https://github.com/owncloud/ocis/pull/7254
   https://github.com/owncloud/ocis/pull/7262
   https://github.com/owncloud/web/pull/9653
   https://github.com/owncloud/web/pull/9672

* Enhancement - Introduce clientlog service: [#7217](https://github.com/owncloud/ocis/pull/7217)

   Add the clientlog service which will send machine readable notifications to
   clients

   https://github.com/owncloud/ocis/pull/7217

* Enhancement - Proxy uses service accounts for provisioning: [#7240](https://github.com/owncloud/ocis/pull/7240)

   The proxy service now uses a service account for provsioning task, like role
   assignment and user auto-provisioning. This cleans up some technical debt that
   required us to mint reva tokes inside the proxy service.

   https://github.com/owncloud/ocis/issues/5550
   https://github.com/owncloud/ocis/pull/7240

* Enhancement - The password policies change request: [#7264](https://github.com/owncloud/ocis/pull/7264)

   The variables renaming OCIS-3767

   https://github.com/owncloud/ocis/pull/7264

* Enhancement - Introduce natsjs registry: [#7272](https://github.com/owncloud/ocis/issues/7272)

   Introduce a registry based on the natsjs object store

   https://github.com/owncloud/ocis/issues/7272
   https://github.com/owncloud/ocis/pull/7487

* Enhancement - Add the password policies: [#7285](https://github.com/owncloud/ocis/pull/7285)

   Add the password policies OCIS-3767

   https://github.com/owncloud/ocis/pull/7285
   https://github.com/owncloud/ocis/pull/7194
   https://github.com/cs3org/reva/pull/4147

* Enhancement - Add login URL config: [#7317](https://github.com/owncloud/ocis/pull/7317)

   Introduce a config to set the web login URL via `WEB_OPTION_LOGIN_URL`.

   https://github.com/owncloud/ocis/pull/7317

* Enhancement - Improve SSE format: [#7325](https://github.com/owncloud/ocis/pull/7325)

   Improve format of sse notifications

   https://github.com/owncloud/ocis/pull/7325

* Enhancement - New value `auto` for NOTIFICATIONS_SMTP_AUTHENTICATION: [#7356](https://github.com/owncloud/ocis/issues/7356)

   This cause the notifications service to automatically pick a suitable
   authentication method to use with the configured SMTP server. This is also the
   new default behavior. The previous default was to not use authentication at all.

   https://github.com/owncloud/ocis/issues/7356

* Enhancement - Make sse service scalable: [#7382](https://github.com/owncloud/ocis/pull/7382)

   When running multiple sse instances some events would not be reported to the
   user. This is fixed.

   https://github.com/owncloud/ocis/pull/7382

* Enhancement - Edit wrong named enves: [#7406](https://github.com/owncloud/ocis/pull/7406)

   Checked and changed the envvars specified in the task and also removed those
   that are no longer used.

   https://github.com/owncloud/ocis/pull/7406

* Enhancement - Thumbnail generation with image processors: [#7409](https://github.com/owncloud/ocis/pull/7409)

   Thumbnails can now be changed during creation, previously the images were always
   scaled to fit the given frame, but it could happen that the images were cut off
   because they could not be placed better due to the aspect ratio.

   This pr introduces the possibility of specifying how the behavior should be,
   following processors are available

  * resize
  * fit
  * fill
  * thumbnail

   The processor can be applied by adding the processor query param to the request,
   e.g. `processor=fit`, `processor=fill`, ...

   To find out more how the individual processors work please read
   https://github.com/disintegration/imaging

   If no processor is provided it behaves the same as before (resize for gif's and
   thumbnail for all other)

   https://github.com/owncloud/enterprise/issues/6057
   https://github.com/owncloud/ocis/issues/5179
   https://github.com/owncloud/web/issues/7728
   https://github.com/owncloud/ocis/pull/7409

* Enhancement - Set default for Async Uploads to true: [#7416](https://github.com/owncloud/ocis/pull/7416)

   Async Uploads are meanwhile standard and needed for multiple features. Hence we
   default them to true

   https://github.com/owncloud/ocis/pull/7416

* Enhancement - The default language added: [#7417](https://github.com/owncloud/ocis/pull/7417)

   The ability of configuration the default language has been added to the setting
   service.

   https://github.com/owncloud/enterprise/issues/5915
   https://github.com/owncloud/ocis/pull/7417

* Enhancement - Add "Last modified" filter Chip: [#7455](https://github.com/owncloud/ocis/pull/7455)

   Add "Last modified" filter Chip

   https://github.com/owncloud/ocis/issues/7431
   https://github.com/owncloud/ocis/issues/7551
   https://github.com/owncloud/ocis/pull/7455

* Enhancement - Config for disabling Web extensions: [#7486](https://github.com/owncloud/ocis/pull/7486)

   A new config for disabling specific Web extensions via their id has been added.

   https://github.com/owncloud/web/issues/8524
   https://github.com/owncloud/ocis/pull/7486

* Enhancement - Store and index metadata: [#7490](https://github.com/owncloud/ocis/pull/7490)

   Audio metadata is now extracted and stored by the search service. It is
   available for driveItems in a folder listing using the Graph API.

   https://github.com/owncloud/ocis/pull/7490

* Enhancement - Add support for audio files to the thumbnails service: [#7491](https://github.com/owncloud/ocis/pull/7491)

   The thumbnails service can now extract artwork from audio files (mp3, ogg, flac)
   and render it just like any other image.

   https://github.com/owncloud/ocis/pull/7491

* Enhancement - Implement sharing roles: [#7524](https://github.com/owncloud/ocis/pull/7524)

   Implement libre graph sharing roles

   https://github.com/owncloud/ocis/issues/7418
   https://github.com/owncloud/ocis/pull/7524

* Enhancement - Add new permission to delete public link password: [#7538](https://github.com/owncloud/ocis/pull/7538)

   Users with this new permission can now delete passwords on read-only public
   links. The permission is added to the default roles "Admin" and "Space Admin".

   https://github.com/owncloud/ocis/issues/7538
   https://github.com/owncloud/ocis/pull/7538
   https://github.com/cs3org/reva/pull/4270

* Enhancement - Add config to enforce passwords on all public links: [#7547](https://github.com/owncloud/ocis/pull/7547)

   We added the config `OCIS_SHARING_PUBLIC_SHARE_MUST_HAVE_PASSWORD` to enforce
   passwords on all public shares.

   https://github.com/owncloud/ocis/issues/7539
   https://github.com/owncloud/ocis/pull/7547

* Enhancement - Tika content extraction cleanup for search: [#7553](https://github.com/owncloud/ocis/pull/7553)

   So far it has not been possible to determine whether the content for search
   should be cleaned up of 'stop words' or not. Stop words are filling words like
   "I, you, have, am" etc and defined by the search engine.

   The behaviour can now be set with the newly introduced settings option
   `SEARCH_EXTRACTOR_TIKA_CLEAN_STOP_WORDS=false` which is enabled by default.

   In addition, the stop word cleanup is no longer as aggressive and now ignores
   numbers, urls, basically everything except the defined stop words.

   https://github.com/owncloud/ocis/issues/6674
   https://github.com/owncloud/ocis/pull/7553

* Enhancement - Allow configuring storage registry with envvars: [#7554](https://github.com/owncloud/ocis/pull/7554)

   Introduced new envvars to configure the storage registry in the gateway service

   https://github.com/owncloud/ocis/pull/7554

* Enhancement - Add search MediaType filter: [#7602](https://github.com/owncloud/ocis/pull/7602)

   Add filter MediaType filter shortcuts to search for specific document types. For
   example, a search query mediatype:documents will search for files with the
   following mimetypes:

   Application/msword
   MimeType:application/vnd.openxmlformats-officedocument.wordprocessingml.document
   MimeType:application/vnd.oasis.opendocument.text MimeType:text/plain
   MimeType:text/markdown MimeType:application/rtf
   MimeType:application/vnd.apple.pages

   Besides the document shorthand, it also contains following:

  * file
  * folder
  * document
  * spreadsheet
  * presentation
  * pdf
  * image
  * video
  * audio
  * archive

   ## File

   ## Folder

   ## Document:

   Application/msword
   application/vnd.openxmlformats-officedocument.wordprocessingml.document
   application/vnd.oasis.opendocument.text text/plain text/markdown application/rtf
   application/vnd.apple.pages

   ## Spreadsheet:

   Application/vnd.ms-excel application/vnd.oasis.opendocument.spreadsheet text/csv
   application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
   application/vnd.oasis.opendocument.spreadsheet application/vnd.apple.numbers

   ## Presentations:

   Application/vnd.ms-powerpoint
   application/vnd.openxmlformats-officedocument.presentationml.presentation
   application/vnd.oasis.opendocument.presentation application/vnd.apple.keynote

   ## PDF

   Application/pdf

   ## Image:

   Image/*

   ## Video:

   Video/*

   ## Audio:

   Audio/*

   ## Archive (zip ...):

   Application/zip application/x-tar application/x-gzip application/x-7z-compressed
   application/x-rar-compressed application/x-bzip2 application/x-bzip
   application/x-tgz

   https://github.com/owncloud/ocis/issues/7432
   https://github.com/owncloud/ocis/pull/7602

* Enhancement - Add Sharing NG endpoints: [#7633](https://github.com/owncloud/ocis/pull/7633)

   We've added new sharing ng endpoints to the graph beta api. The following
   endpoints are added:

  * /v1beta1/me/drive/sharedByMe
  * /v1beta1/me/drive/sharedWithMe
  * /v1beta1/roleManagement/permissions/roleDefinitions
  * /v1beta1/roleManagement/permissions/roleDefinitions/{roleID}
  * /v1beta1/drives/{drive-id}/items/{item-id}/createLink (create a sharing link)

   https://github.com/owncloud/ocis/issues/7436
   https://github.com/owncloud/ocis/issues/6993
   https://github.com/owncloud/ocis/pull/7633
   https://github.com/owncloud/ocis/pull/7686
   https://github.com/owncloud/ocis/pull/7684
   https://github.com/owncloud/ocis/pull/7683
   https://github.com/owncloud/ocis/pull/7239
   https://github.com/owncloud/ocis/pull/7687
   https://github.com/owncloud/ocis/pull/7751
   https://github.com/owncloud/libre-graph-api/pull/112

* Enhancement - Configs for Web embed mode: [#7670](https://github.com/owncloud/ocis/pull/7670)

   New configs for the Web embed mode have been added:

  * `enabled` Defines if embed mode is enabled.
  * `target` Defines how Web is being integrated when running in embed mode.
  * `messagesOrigin` Defines a URL under which Web can be integrated via iFrame.
  * `delegateAuthentication` Defines whether Web should require authentication to be done by the parent application.
  * `delegateAuthenticationOrigin` Defines the host to validate the message event origin against when running Web in 'embed' mode.

   https://github.com/owncloud/web/issues/9768
   https://github.com/owncloud/ocis/pull/7670

* Enhancement - Support login page background configuration: [#7674](https://github.com/owncloud/ocis/issues/7674)

   Introduce a new environment variable `IDP_LOGIN_BACKGROUND_URL` that overrides
   the default background image of the IDP login page when present.

   https://github.com/owncloud/ocis/issues/7674
   https://github.com/owncloud/ocis/pull/7900

* Enhancement - Add new permissions: [#7700](https://github.com/owncloud/ocis/pull/7700)

   Adds new permissions to admin/spaceadmin/user roles - Favorites.List allows /
   denies the Favorites Listing Request - Favorites.Write is implemented to be
   enforced on marking/unmark files as favouritesShare - Shares.Write permission
   denies / allows sharing completely for a user on all share CUD requests. (User,
   Group)

   https://github.com/owncloud/ocis/pull/7700

* Enhancement - Add preferred language to user settings: [#7720](https://github.com/owncloud/ocis/pull/7720)

   We have added the preferred language to the libre-graph api & added endpoints
   for that to ocis.

   https://github.com/owncloud/ocis/issues/5455
   https://github.com/owncloud/ocis/pull/7720
   https://github.com/owncloud/libre-graph-api/pull/130

* Enhancement - Add user filter startswith and contains: [#7739](https://github.com/owncloud/ocis/pull/7739)

   We add two new filters to the user list endpoint. The `startswith` filter allows
   to filter users by the beginning of their name. The `contains` filter allows to
   filter users by a substring of their name.

   https://github.com/owncloud/ocis/issues/5486
   https://github.com/owncloud/ocis/pull/7739

* Enhancement - Allow configuring additional routes: [#7741](https://github.com/owncloud/ocis/pull/7741)

   Allows adding custom routes to the ocis proxy. This enables custom ocis
   extensions.

   https://github.com/owncloud/ocis/pull/7741

* Enhancement - Default link permission config: [#7783](https://github.com/owncloud/ocis/pull/7783)

   A new config for default link permissions that is being announced via
   capabilities has been added. It defaults to 1 (= public link with viewer
   permissions).

   https://github.com/owncloud/web/issues/9919
   https://github.com/owncloud/ocis/pull/7783

* Enhancement - Add banned password list to the default deployments: [#7784](https://github.com/owncloud/ocis/pull/7784)

   We add banned password list to the default deployments

   https://github.com/owncloud/ocis/issues/7724
   https://github.com/owncloud/ocis/pull/7784

* Enhancement - Update to go 1.21: [#7794](https://github.com/owncloud/ocis/pull/7794)

   We updated to go 1.21.

   https://github.com/owncloud/ocis/pull/7794

* Enhancement - Add Sharing NG list permissions endpoint: [#7805](https://github.com/owncloud/ocis/pull/7805)

   We've added a new sharing ng endpoint which lists all permissions for a given
   item.

   https://github.com/owncloud/ocis/issues/6993
   https://github.com/owncloud/ocis/pull/7805

* Enhancement - Add user list requires filter config: [#7866](https://github.com/owncloud/ocis/pull/7866)

   Introduce a config to require filters in order to list users in Web via
   `WEB_OPTION_USER_LIST_REQUIRES_FILTER`.

   https://github.com/owncloud/ocis/pull/7866

* Enhancement - Retry antivirus postprocessing step in case of problems: [#7874](https://github.com/owncloud/ocis/pull/7874)

   The antivirus postprocessing step will now be retried for a configurable amount
   of times in case it can't get a result from clamav.

   https://github.com/owncloud/ocis/pull/7874

* Enhancement - Add validation to public share provider: [#7877](https://github.com/owncloud/ocis/pull/7877)

   We changed the implementation of the public share provider in reva to do the
   validation on the CS3 Api side. This makes the implementation on the graph side
   smaller.

   https://github.com/owncloud/ocis/issues/6993
   https://github.com/owncloud/ocis/pull/7877

* Enhancement - Graphs endpoint for mounting and unmounting shares: [#7885](https://github.com/owncloud/ocis/pull/7885)

   Functionality for mounting (accepting) and unmounting (rejecting) received
   shares has been added to the graph API.

   https://github.com/owncloud/ocis/pull/7885

* Enhancement - Store and index metadata: [#7886](https://github.com/owncloud/ocis/pull/7886)

   Location metadata is now extracted and stored by the search service. It is
   available for driveItems in a folder listing using the Graph API.

   https://github.com/owncloud/ocis/pull/7886

* Enhancement - Allow regular users to list other users: [#7887](https://github.com/owncloud/ocis/pull/7887)

   Regular users can search for other users and groups. The following limitations
   apply:

  * Only search queries are allowed (using the `$search=term` query parameter)
  * The search term needs to have at least 3 characters
  * for user searches the result set only contains the attributes `displayName`, `userType`, `mail` and `id`
  * for group searches the result set only contains the attributes `displayName`, `groupTypes` and `id`

   https://github.com/owncloud/ocis/issues/7782
   https://github.com/owncloud/ocis/pull/7887

* Enhancement - Add edit public share to sharing NG: [#7908](https://github.com/owncloud/ocis/pull/7908/)

   We added the ability to edit public shares to the sharing NG endpoints.

   https://github.com/owncloud/ocis/issues/6993
   https://github.com/owncloud/ocis/pull/7908/

* Enhancement - Add cli commands for trash-bin: [#7917](https://github.com/owncloud/ocis/pull/7917)

   We added the `list` and `restore` commands to the trash-bin items to the CLI

   https://github.com/owncloud/ocis/issues/7845
   https://github.com/owncloud/ocis/pull/7917
   https://github.com/cs3org/reva/pull/4392

* Enhancement - Add validation update public share: [#7978](https://github.com/owncloud/ocis/pull/7978)

   For Sharing NG, we needed validation in the implementing reva service to keep
   the client implementation simple.

   https://github.com/owncloud/ocis/pull/7978

* Enhancement - Allow inmemory nats-js-kv stores: [#7979](https://github.com/owncloud/ocis/pull/7979)

   Adds envvars to keep nats-js-kv stores in memory and not persist them on disc.

   https://github.com/owncloud/ocis/pull/7979

* Enhancement - Disable the password policy: [#7985](https://github.com/owncloud/ocis/pull/7985)

   We add the environment variable that allow to disable the password policy.

   https://github.com/owncloud/ocis/issues/7916
   https://github.com/owncloud/ocis/pull/7985
   https://github.com/cs3org/reva/pull/4409

* Enhancement - Use kv store in natsjs registry: [#7987](https://github.com/owncloud/ocis/pull/7987)

   Replaces the nats object store with the nats kv store in the natsjs registry

   https://github.com/owncloud/ocis/pull/7987

* Enhancement - Allow authentication nats connections: [#7989](https://github.com/owncloud/ocis/pull/7989)

   Allow events, store and registry implementation to pass username/password to the
   nats instance

   https://github.com/owncloud/ocis/pull/7989

* Enhancement - Add RED metrics to the metrics endpoint: [#7994](https://github.com/owncloud/ocis/pull/7994)

   We added three new metrics to the metrics endpoint to support the RED method for
   monitoring microservices.

   - Request Rate: The number of requests per second. The total count of requests
   is available under `ocis_proxy_requests_total`. - Error Rate: The number of
   failed requests per second. The total count of failed requests is available
   under `ocis_proxy_errors_total`. - Duration: The amount of time each request
   takes. The duration of all requests is available under
   `ocis_proxy_request_duration_seconds`. This is a histogram metric, so it also
   provides information about the distribution of request durations.

   The metrics are available under the following paths: `PROXY_DEBUG_ADDR/metrics`
   in a prometheus compatible format and maybe secured by `PROXY_DEBUG_TOKEN`.

   https://github.com/owncloud/ocis/pull/7994

* Enhancement - Add ocm and sciencemesh services: [#7998](https://github.com/owncloud/ocis/pull/7998)

   We added sciencemesh and ocm services to enable federation.

   https://github.com/owncloud/ocis/pull/7998
   https://github.com/owncloud/ocis/pull/7576
   https://github.com/owncloud/ocis/pull/7464
   https://github.com/owncloud/ocis/pull/7463

* Enhancement - Make nats-js-kv the default registry: [#8011](https://github.com/owncloud/ocis/pull/8011)

   The previously used default `mdns` is faulty. Deprecated it together with
   `consul`, `nats` and `etcd` implementations.

   https://github.com/owncloud/ocis/pull/8011
   https://github.com/owncloud/ocis/pull/8027

* Enhancement - Service Account roles: [#8051](https://github.com/owncloud/ocis/pull/8051)

   Use a hidden role for service accounts. It will not appear in ListRoles calls
   but internally handled by settings service

   https://github.com/owncloud/ocis/pull/8051
   https://github.com/owncloud/ocis/pull/8074

* Enhancement - Update antivirus service: [#8062](https://github.com/owncloud/ocis/pull/8062)

   We update the antivirus icap client library and optimize the antivirus scanning
   service. ANTIVIRUS_ICAP_TIMEOUT is now deprecated and
   ANTIVIRUS_ICAP_SCAN_TIMEOUT should be used instead.

   ANTIVIRUS_ICAP_SCAN_TIMEOUT supports human durations like `1s`, `1m`, `1h` and
   `1d`.

   https://github.com/owncloud/ocis/issues/6764
   https://github.com/owncloud/ocis/pull/8062

* Enhancement - Remove deprecated environment variables: [#8149](https://github.com/owncloud/ocis/pull/8149)

   We have removed all deprecated environment variables that have been marked for
   removal for 5.0.0

   https://github.com/owncloud/ocis/issues/8025
   https://github.com/owncloud/ocis/pull/8149

* Enhancement - Disable the password policy: [#8152](https://github.com/owncloud/ocis/pull/8152)

   We reworked and moved disabling the password policy logic from the reva to the
   ocis.

   https://github.com/owncloud/ocis/issues/7916
   https://github.com/owncloud/ocis/pull/8152
   https://github.com/cs3org/reva/pull/4453

* Enhancement - Allow restarting multiple uploads with one command: [#8287](https://github.com/owncloud/ocis/pull/8287)

   Allows to restart all commands in a specific state.

   https://github.com/owncloud/ocis/pull/8287

* Enhancement - Modify the concurrency default: [#8309](https://github.com/owncloud/ocis/pull/8309)

   We have changed the default MaxConcurrency value from 100 to 5 to prevent too
   frequent gc runs on low memory systems. We have also bumped reva to pull in the
   related changes from there.

   https://github.com/owncloud/ocis/issues/8257
   https://github.com/owncloud/ocis/pull/8309
   https://github.com/cs3org/reva/pull/4485

* Enhancement - Improve ocis single binary start: [#8320](https://github.com/owncloud/ocis/pull/8320)

   Removes waiting times when starting the single binary. Improves ocis single
   binary boot time from 8s to 2.5s

   https://github.com/owncloud/ocis/pull/8320

* Enhancement - Use environment variables in yaml config files: [#8339](https://github.com/owncloud/ocis/pull/8339)

   We added the ability to use environment variables in yaml config files. This
   allows to use environment variables in the config files of the ocis services
   which will be replaced by the actual value of the environment variable at
   runtime.

   Example:

   ```
   web:
     http:
       addr: ${SOME_HTTP_ADDR}
   ```

   This makes it possible to use the same config file for different environments
   without the need to change the config file itself. This is especially useful
   when using docker-compose to run the ocis services. It is a common pattern to
   create an .env file which contains the environment variables for the
   docker-compose file. Now you can use the same .env file to configure the ocis
   services.

   https://github.com/owncloud/ocis/pull/8339

* Enhancement - Increment filenames on upload collisions in secret filedrops: [#8340](https://github.com/owncloud/ocis/pull/8340)

   We have bumped reva to pull in the changes needed for automatically increment
   filenames on upload collisions in secret filedrops.

   https://github.com/owncloud/ocis/issues/8291
   https://github.com/owncloud/ocis/pull/8340

* Enhancement - Allow sending multiple user ids in one sse event: [#8379](https://github.com/owncloud/ocis/pull/8379)

   Sending multiple user ids in one sse event is now possible which reduces the
   number of sent events.

   https://github.com/owncloud/ocis/pull/8379
   https://github.com/cs3org/reva/pull/4501

* Enhancement - Allow to skip service listing: [#8408](https://github.com/owncloud/ocis/pull/8408)

   The ocis version cmd listed all services by default. This is not always
   intended, so we allow to skip the listing of the services by using the
   --skip-services flag.

   https://github.com/owncloud/ocis/issues/8070
   https://github.com/owncloud/ocis/pull/8408

* Enhancement - Add a make step to validate the env var annotations: [#8436](https://github.com/owncloud/ocis/pull/8436)

   We have added a make step `make check-env-var-annotations` to validate the
   environment variable annotations in to the environment variables.

   https://github.com/owncloud/ocis/issues/8258
   https://github.com/owncloud/ocis/pull/8436

* Enhancement - Drop the unnecessary grants exists check when creating shares: [#8502](https://github.com/owncloud/ocis/pull/8502)

   We have bumped reva to drop the unnecessary grants exists check when creating
   shares.

   https://github.com/owncloud/ocis/pull/8502

* Enhancement - Update to go 1.22: [#8586](https://github.com/owncloud/ocis/pull/8586)

   We have updated go to version 1.22.

   https://github.com/owncloud/ocis/pull/8586

* Enhancement - Update web to v8.0.0: [#8613](https://github.com/owncloud/ocis/pull/8613)

   Tags: web

   We updated ownCloud Web to v8.0.0. Please refer to the changelog (linked) for
   details on the web release.

  * Bugfix [owncloud/web#9257](https://github.com/owncloud/web/issues/9257): Filter out shares without display name
  * Bugfix [owncloud/web#9529](https://github.com/owncloud/web/pull/9529): Shared with action menu label alignment
  * Bugfix [owncloud/web#9649](https://github.com/owncloud/web/pull/9649): Add project space filter
  * Bugfix [owncloud/web#9663](https://github.com/owncloud/web/pull/9663): Respect the open-in-new-tab-config for external apps
  * Bugfix [owncloud/web#9694](https://github.com/owncloud/web/issues/9694): Special characters in username
  * Bugfix [owncloud/web#9788](https://github.com/owncloud/web/issues/9788): Create .space folder if it does not exist
  * Bugfix [owncloud/web#9799](https://github.com/owncloud/web/issues/9799): Link resolving into default app
  * Bugfix [owncloud/web#9832](https://github.com/owncloud/web/pull/9832): Copy quicklinks for webkit navigator
  * Bugfix [owncloud/web#9843](https://github.com/owncloud/web/pull/9843): Fix display path on resources
  * Bugfix [owncloud/web#9844](https://github.com/owncloud/web/pull/9844): Upload space image
  * Bugfix [owncloud/web#9861](https://github.com/owncloud/web/pull/9861): Duplicated file search request
  * Bugfix [owncloud/web#9873](https://github.com/owncloud/web/pull/9873): Tags are no longer editable for a locked file
  * Bugfix [owncloud/web#9881](https://github.com/owncloud/web/pull/9881): Prevent rendering of old/wrong set of resources in search list
  * Bugfix [owncloud/web#9915](https://github.com/owncloud/web/pull/9915): Keep both folders conflict in same-named folders
  * Bugfix [owncloud/web#9931](https://github.com/owncloud/web/pull/9931): Enabling "invite people" for password-protected folder/file
  * Bugfix [owncloud/web#10010](https://github.com/owncloud/web/issues/10010): Displaying full video in their dimensions
  * Bugfix [owncloud/web#10031](https://github.com/owncloud/web/issues/10031): Icon extension mapping
  * Bugfix [owncloud/web#10065](https://github.com/owncloud/web/pull/10065): Logout page after token expiry
  * Bugfix [owncloud/web#10083](https://github.com/owncloud/web/pull/10083): Disable expiration date for alias link (internal)
  * Bugfix [owncloud/web#10092](https://github.com/owncloud/web/pull/10092): Allow empty search query in "in-here" search
  * Bugfix [owncloud/web#10096](https://github.com/owncloud/web/pull/10096): Remove password buttons on input if disabled
  * Bugfix [owncloud/web#10118](https://github.com/owncloud/web/pull/10118): Tilesview has whitespace
  * Bugfix [owncloud/web#10149](https://github.com/owncloud/web/pull/10149): Spaces files list previews cropped
  * Bugfix [owncloud/web#10149](https://github.com/owncloud/web/pull/10149): Spaces overview tile previews zoomed
  * Bugfix [owncloud/web#10154](https://github.com/owncloud/web/pull/10154): Resolving links without drive alias
  * Bugfix [owncloud/web#10156](https://github.com/owncloud/web/pull/10156): Uploading the same files parallel
  * Bugfix [owncloud/web#10158](https://github.com/owncloud/web/pull/10158): GDPR export polling
  * Bugfix [owncloud/web#10176](https://github.com/owncloud/web/pull/10176): Turned off file extensions not always respected
  * Bugfix [owncloud/web#10179](https://github.com/owncloud/web/pull/10179): Space navigate to trash missing
  * Bugfix [owncloud/web#10182](https://github.com/owncloud/web/pull/10182): Make versions panel readonly in viewers and editors
  * Bugfix [owncloud/web#10220](https://github.com/owncloud/web/pull/10220): Loading indicator during conflict dialog
  * Bugfix [owncloud/web#10227](https://github.com/owncloud/web/issues/10227): Configurable concurrent requests
  * Bugfix [owncloud/web#10232](https://github.com/owncloud/web/pull/10232): Skip searchbar preview fetch on reload
  * Bugfix [owncloud/web#10318](https://github.com/owncloud/web/pull/10318): Scrollable account page
  * Bugfix [owncloud/web#10321](https://github.com/owncloud/web/pull/10321): Private link error messages
  * Bugfix [owncloud/web#10347](https://github.com/owncloud/web/pull/10347): Readonly user attributes have no effect on group memberships
  * Bugfix [owncloud/web#10424](https://github.com/owncloud/web/pull/10424): Restore space
  * Bugfix [owncloud/web#10473](https://github.com/owncloud/web/issues/10473): Public link file download
  * Bugfix [owncloud/web#10489](https://github.com/owncloud/web/pull/10489): Wrong share permissions when resharing off
  * Bugfix [owncloud/web#10514](https://github.com/owncloud/web/pull/10514): Indicate shares that are not manageable due to file locking
  * Change [owncloud/web#2404](https://github.com/owncloud/web/issues/2404): Theme handling
  * Change [owncloud/web#7338](https://github.com/owncloud/web/issues/7338): Remove deprecated code
  * Change [owncloud/web#9653](https://github.com/owncloud/web/pull/9653): Keyword Query Language (KQL) search syntax
  * Change [owncloud/web#9709](https://github.com/owncloud/web/issues/9709): DavProperties without namespace
  * Enhancement [owncloud/web#7317](https://github.com/owncloud/ocis/pull/7317): Make login url configurable
  * Enhancement [owncloud/web#7497](https://github.com/owncloud/ocis/issues/7497): Permission checks for shares and favorites
  * Enhancement [owncloud/web#7600](https://github.com/owncloud/web/issues/7600): Scroll to newly created folder
  * Enhancement [owncloud/web#9302](https://github.com/owncloud/web/issues/9302): Application unification
  * Enhancement [owncloud/web#9423](https://github.com/owncloud/web/pull/9423): Show local loading spinner in sharing button
  * Enhancement [owncloud/web#9441](https://github.com/owncloud/web/pull/9441): File versions tooltip with absolute date
  * Enhancement [owncloud/web#9441](https://github.com/owncloud/web/pull/9441): Disabling extensions
  * Enhancement [owncloud/web#9451](https://github.com/owncloud/web/pull/9451): Add SSE to get notifications instantly
  * Enhancement [owncloud/web#9525](https://github.com/owncloud/web/pull/9525): Tags form improved
  * Enhancement [owncloud/web#9527](https://github.com/owncloud/web/pull/9527): Don't display confirmation dialog on file deletion
  * Enhancement [owncloud/web#9531](https://github.com/owncloud/web/issues/9531): Personal shares can be shown and hidden
  * Enhancement [owncloud/web#9552](https://github.com/owncloud/web/pull/9552): Upload preparation time
  * Enhancement [owncloud/web#9561](https://github.com/owncloud/web/pull/9561): Indicate processing state
  * Enhancement [owncloud/web#9566](https://github.com/owncloud/web/pull/9566): Display locking information
  * Enhancement [owncloud/web#9584](https://github.com/owncloud/web/pull/9584): Moving share's "set expiration date" function
  * Enhancement [owncloud/web#9625](https://github.com/owncloud/web/pull/9625): Add keyboard navigation to spaces overview
  * Enhancement [owncloud/web#9627](https://github.com/owncloud/web/pull/9627): Add batch actions to spaces
  * Enhancement [owncloud/web#9671](https://github.com/owncloud/web/pull/9671): OcModal set buttons to same width
  * Enhancement [owncloud/web#9682](https://github.com/owncloud/web/pull/9682): Add password policy compatibility
  * Enhancement [owncloud/web#9691](https://github.com/owncloud/web/pull/9691): Password generator for public links
  * Enhancement [owncloud/web#9696](https://github.com/owncloud/web/pull/9696): Added app banner for mobile devices
  * Enhancement [owncloud/web#9706](https://github.com/owncloud/web/pull/9706): Unify sharing expiration date menu items
  * Enhancement [owncloud/web#9709](https://github.com/owncloud/web/issues/9709): New WebDAV implementation in web-client
  * Enhancement [owncloud/web#9727](https://github.com/owncloud/web/pull/9727): Show error if password is on a banned password list
  * Enhancement [owncloud/web#9768](https://github.com/owncloud/web/issues/9768): Embed mode
  * Enhancement [owncloud/web#9771](https://github.com/owncloud/web/pull/9771): Handle postprocessing state via Server Sent Events
  * Enhancement [owncloud/web#9794](https://github.com/owncloud/web/pull/9794): Registering search providers as extension
  * Enhancement [owncloud/web#9806](https://github.com/owncloud/web/pull/9806): Preview image presentation
  * Enhancement [owncloud/web#9809](https://github.com/owncloud/web/pull/9809): Add editors to the application menu
  * Enhancement [owncloud/web#9814](https://github.com/owncloud/web/pull/9814): Registering nav items as extension
  * Enhancement [owncloud/web#9815](https://github.com/owncloud/web/pull/9815): Add new portal into runtime to include footer
  * Enhancement [owncloud/web#9831](https://github.com/owncloud/web/pull/9831): Last modified filter chips
  * Enhancement [owncloud/web#9847](https://github.com/owncloud/web/issues/9847): Provide vendor neutral file icons
  * Enhancement [owncloud/web#9854](https://github.com/owncloud/web/pull/9854): Search query term linking
  * Enhancement [owncloud/web#9857](https://github.com/owncloud/web/pull/9857): Add permission to delete link passwords when password is enforced
  * Enhancement [owncloud/web#9858](https://github.com/owncloud/web/pull/9858): Remove settings icon from searchbar
  * Enhancement [owncloud/web#9864](https://github.com/owncloud/web/pull/9864): Search tags filter chips style aligned
  * Enhancement [owncloud/web#9884](https://github.com/owncloud/web/pull/9884): Enable dark theme on importer
  * Enhancement [owncloud/web#9890](https://github.com/owncloud/web/pull/9890): Create shortcuts
  * Enhancement [owncloud/web#9905](https://github.com/owncloud/web/pull/9905): Manage tags in details panel
  * Enhancement [owncloud/web#9906](https://github.com/owncloud/web/pull/9906): Reorganize "New" menu
  * Enhancement [owncloud/web#9912](https://github.com/owncloud/web/pull/9912): Add media type filter chip
  * Enhancement [owncloud/web#9940](https://github.com/owncloud/web/pull/9940): Display error message for upload to locked folder
  * Enhancement [owncloud/web#9966](https://github.com/owncloud/web/issues/9966): Support more audio formats with correct icon
  * Enhancement [owncloud/web#10007](https://github.com/owncloud/web/issues/10007): Additional languages
  * Enhancement [owncloud/web#10013](https://github.com/owncloud/web/issues/10013): Shared by filter
  * Enhancement [owncloud/web#10014](https://github.com/owncloud/web/issues/10014): Share search filter
  * Enhancement [owncloud/web#10024](https://github.com/owncloud/web/pull/10024): Duplicate space
  * Enhancement [owncloud/web#10037](https://github.com/owncloud/web/pull/10037): Default link permission
  * Enhancement [owncloud/web#10047](https://github.com/owncloud/web/pull/10047): Add explaining contextual helper to spaces overview
  * Enhancement [owncloud/web#10057](https://github.com/owncloud/web/pull/10057): Folder tree creation during upload
  * Enhancement [owncloud/web#10062](https://github.com/owncloud/web/pull/10062): Show webdav information in details view
  * Enhancement [owncloud/web#10099](https://github.com/owncloud/web/pull/10099): Support mandatory filter while listing users
  * Enhancement [owncloud/web#10102](https://github.com/owncloud/web/pull/10102): Registering quick actions as extension
  * Enhancement [owncloud/web#10104](https://github.com/owncloud/web/pull/10104): Create link modal
  * Enhancement [owncloud/web#10111](https://github.com/owncloud/web/pull/10111): Registering right sidebar panels as extension
  * Enhancement [owncloud/web#10111](https://github.com/owncloud/web/pull/10111): File sidebar in viewer and editor apps
  * Enhancement [owncloud/web#10224](https://github.com/owncloud/web/pull/10224): Harmonize AppSwitcher icon colors
  * Enhancement [owncloud/web#10356](https://github.com/owncloud/web/pull/10356): Preview app add reset button for images

   https://github.com/owncloud/ocis/pull/8613
   https://github.com/owncloud/web/releases/tag/v8.0.0

* Enhancement - Update web to v8.0.1: [#8626](https://github.com/owncloud/ocis/pull/8626)

   Tags: web

   We updated ownCloud Web to v8.0.1. Please refer to the changelog (linked) for
   details on the web release.

  * Bugfix [owncloud/web#10573](https://github.com/owncloud/web/pull/10573): Add link in right sidebar sharing menu, doesn't copy link to clipboard
  * Bugfix [owncloud/web#10576](https://github.com/owncloud/web/pull/10576): WebDav Url in right sidebar is missing dav in path
  * Bugfix [owncloud/web#10585](https://github.com/owncloud/web/issues/10585): Update translations

   https://github.com/owncloud/ocis/pull/8626
   https://github.com/owncloud/web/releases/tag/v8.0.1

* Enhancement - Update reva to 2.19.2: [#8638](https://github.com/owncloud/ocis/pull/8638)

   We update reva to the version 2.19.2

  *   Bugfix [cs3org/reva#4557](https://github.com/cs3org/reva/pull/4557): Fix ceph build
  *   Bugfix [cs3org/reva#4570](https://github.com/cs3org/reva/pull/4570): Fix sharing invite on virtual drive
  *   Bugfix [cs3org/reva#4559](https://github.com/cs3org/reva/pull/4559): Fix graph drive invite
  *   Bugfix [cs3org/reva#4518](https://github.com/cs3org/reva/pull/4518): Fix an error when lock/unlock a file
  *   Bugfix [cs3org/reva#4566](https://github.com/cs3org/reva/pull/4566): Fix public link previews
  *   Bugfix [cs3org/reva#4561](https://github.com/cs3org/reva/pull/4561): Fix Stat() by Path on re-created resource
  *   Enhancement [cs3org/reva#4556](https://github.com/cs3org/reva/pull/4556): Allow tracing requests by giving util functions a context
  *   Enhancement [cs3org/reva#4545](https://github.com/cs3org/reva/pull/4545): Extend service account permissions
  *   Enhancement [cs3org/reva#4564](https://github.com/cs3org/reva/pull/4564): Send file locked/unlocked events

   We update reva to the version 2.19.1

  *   Bugfix [cs3org/reva#4534](https://github.com/cs3org/reva/pull/4534): Fix remove/update share permissions
  *   Bugfix [cs3org/reva#4539](https://github.com/cs3org/reva/pull/4539): Fix a typo

   We update reva to the version 2.19.0

  *   Bugfix [cs3org/reva#4464](https://github.com/cs3org/reva/pull/4464): Don't check lock grants
  *   Bugfix [cs3org/reva#4516](https://github.com/cs3org/reva/pull/4516): The sharemanager can now reject grants with resharing permissions
  *   Bugfix [cs3org/reva#4512](https://github.com/cs3org/reva/pull/4512): Bump dependencies
  *   Bugfix [cs3org/reva#4481](https://github.com/cs3org/reva/pull/4481): Distinguish failure and node metadata reversal
  *   Bugfix [cs3org/reva#4456](https://github.com/cs3org/reva/pull/4456): Do not lose revisions when restoring the first revision
  *   Bugfix [cs3org/reva#4472](https://github.com/cs3org/reva/pull/4472): Fix concurrent access to a map
  *   Bugfix [cs3org/reva#4457](https://github.com/cs3org/reva/pull/4457): Fix concurrent map access in sharecache
  *   Bugfix [cs3org/reva#4498](https://github.com/cs3org/reva/pull/4498): Fix Content-Disposition header in dav
  *   Bugfix [cs3org/reva#4461](https://github.com/cs3org/reva/pull/4461): CORS handling for WebDAV requests fixed
  *   Bugfix [cs3org/reva#4462](https://github.com/cs3org/reva/pull/4462): Prevent setting container specific permissions on files
  *   Bugfix [cs3org/reva#4479](https://github.com/cs3org/reva/pull/4479): Fix creating documents in the app provider
  *   Bugfix [cs3org/reva#4474](https://github.com/cs3org/reva/pull/4474): Make /dav/meta consistent
  *   Bugfix [cs3org/reva#4446](https://github.com/cs3org/reva/pull/4446): Disallow to delete a file during the processing
  *   Bugfix [cs3org/reva#4517](https://github.com/cs3org/reva/pull/4517): Fix duplicated items in the sharejail root
  *   Bugfix [cs3org/reva#4473](https://github.com/cs3org/reva/pull/4473): Decomposedfs now correctly lists sessions
  *   Bugfix [cs3org/reva#4528](https://github.com/cs3org/reva/pull/4528): Respect IfNotExist option when uploading in cs3 metadata storage
  *   Bugfix [cs3org/reva#4503](https://github.com/cs3org/reva/pull/4503): Fix an error when move
  *   Bugfix [cs3org/reva#4466](https://github.com/cs3org/reva/pull/4466): Fix natsjskv store
  *   Bugfix [cs3org/reva#4533](https://github.com/cs3org/reva/pull/4533): Fix recursive trashcan purge
  *   Bugfix [cs3org/reva#4492](https://github.com/cs3org/reva/pull/4492): Fix the resource name
  *   Bugfix [cs3org/reva#4463](https://github.com/cs3org/reva/pull/4463): Fix the resource name
  *   Bugfix [cs3org/reva#4448](https://github.com/cs3org/reva/pull/4448): Fix truncating existing files
  *   Bugfix [cs3org/reva#4434](https://github.com/cs3org/reva/pull/4434): Fix the upload postprocessing
  *   Bugfix [cs3org/reva#4469](https://github.com/cs3org/reva/pull/4469): Handle interrupted uploads
  *   Bugfix [cs3org/reva#4532](https://github.com/cs3org/reva/pull/4532): Jsoncs3 cache fixes
  *   Bugfix [cs3org/reva#4449](https://github.com/cs3org/reva/pull/4449): Keep failed processing status
  *   Bugfix [cs3org/reva#4529](https://github.com/cs3org/reva/pull/4529): We aligned some OCS return codes with oc10
  *   Bugfix [cs3org/reva#4507](https://github.com/cs3org/reva/pull/4507): Make tusd CORS headers configurable
  *   Bugfix [cs3org/reva#4452](https://github.com/cs3org/reva/pull/4452): More efficient share jail
  *   Bugfix [cs3org/reva#4476](https://github.com/cs3org/reva/pull/4476): No need to unmark postprocessing when it was not started
  *   Bugfix [cs3org/reva#4454](https://github.com/cs3org/reva/pull/4454): Skip unnecessary share retrieval
  *   Bugfix [cs3org/reva#4527](https://github.com/cs3org/reva/pull/4527): Unify datagateway method handling
  *   Bugfix [cs3org/reva#4530](https://github.com/cs3org/reva/pull/4530): Drop unnecessary grant exists check
  *   Bugfix [cs3org/reva#4475](https://github.com/cs3org/reva/pull/4475): Upload session specific processing flag
  *   Enhancement [cs3org/reva#4501](https://github.com/cs3org/reva/pull/4501): Allow sending multiple user ids in one sse event
  *   Enhancement [cs3org/reva#4485](https://github.com/cs3org/reva/pull/4485): Modify the concurrency default
  *   Enhancement [cs3org/reva#4526](https://github.com/cs3org/reva/pull/4526): Configurable s3 put options
  *   Enhancement [cs3org/reva#4453](https://github.com/cs3org/reva/pull/4453): Disable the password policy
  *   Enhancement [cs3org/reva#4477](https://github.com/cs3org/reva/pull/4477): Extend ResumePostprocessing event
  *   Enhancement [cs3org/reva#4491](https://github.com/cs3org/reva/pull/4491): Add filename incrementor for secret filedrops
  *   Enhancement [cs3org/reva#4490](https://github.com/cs3org/reva/pull/4490): Lazy initialize public share manager
  *   Enhancement [cs3org/reva#4494](https://github.com/cs3org/reva/pull/4494): Start implementation of a plain posix storage driver
  *   Enhancement [cs3org/reva#4502](https://github.com/cs3org/reva/pull/4502): Add spaceindex.AddAll()

   ## Changelog for reva 2.18.0 (2023-12-22)

   The following sections list the changes in reva 2.18.0 relevant to reva users.
   The changes are ordered by importance.

  *   Bugfix [cs3org/reva#4424](https://github.com/cs3org/reva/pull/4424): Fixed panic in receivedsharecache pkg
  *   Bugfix [cs3org/reva#4425](https://github.com/cs3org/reva/pull/4425): Fix overwriting files with empty files
  *   Bugfix [cs3org/reva#4432](https://github.com/cs3org/reva/pull/4432): Fix /dav/meta endpoint for shares
  *   Bugfix [cs3org/reva#4422](https://github.com/cs3org/reva/pull/4422): Fix disconnected traces
  *   Bugfix [cs3org/reva#4429](https://github.com/cs3org/reva/pull/4429): Internal link creation
  *   Bugfix [cs3org/reva#4407](https://github.com/cs3org/reva/pull/4407): Make ocdav return correct oc:spaceid
  *   Bugfix [cs3org/reva#4410](https://github.com/cs3org/reva/pull/4410): Improve OCM support
  *   Bugfix [cs3org/reva#4402](https://github.com/cs3org/reva/pull/4402): Refactor upload session
  *   Enhancement [cs3org/reva#4421](https://github.com/cs3org/reva/pull/4421): Check permissions before adding, deleting or updating shares
  *   Enhancement [cs3org/reva#4403](https://github.com/cs3org/reva/pull/4403): Add validation to update public share
  *   Enhancement [cs3org/reva#4409](https://github.com/cs3org/reva/pull/4409): Disable the password policy
  *   Enhancement [cs3org/reva#4412](https://github.com/cs3org/reva/pull/4412): Allow authentication for nats connections
  *   Enhancement [cs3org/reva#4411](https://github.com/cs3org/reva/pull/4411): Add option to configure streams non durable
  *   Enhancement [cs3org/reva#4406](https://github.com/cs3org/reva/pull/4406): Rework cache configuration
  *   Enhancement [cs3org/reva#4414](https://github.com/cs3org/reva/pull/4414): Track more upload session metrics

   ## Changelog for reva 2.17.0 (2023-12-12)

   The following sections list the changes in reva 2.17.0 relevant to reva users.
   The changes are ordered by importance.

  *   Bugfix [cs3org/reva#4278](https://github.com/cs3org/reva/pull/4278): Disable DEPTH infinity in PROPFIND
  *   Bugfix [cs3org/reva#4318](https://github.com/cs3org/reva/pull/4318): Do not allow moves between shares
  *   Bugfix [cs3org/reva#4290](https://github.com/cs3org/reva/pull/4290): Prevent panic when trying to move a non-existent file
  *   Bugfix [cs3org/reva#4241](https://github.com/cs3org/reva/pull/4241): Allow an empty credentials chain in the auth middleware
  *   Bugfix [cs3org/reva#4216](https://github.com/cs3org/reva/pull/4216): Fix an error message
  *   Bugfix [cs3org/reva#4324](https://github.com/cs3org/reva/pull/4324): Fix capabilities decoding
  *   Bugfix [cs3org/reva#4267](https://github.com/cs3org/reva/pull/4267): Fix concurrency issue
  *   Bugfix [cs3org/reva#4362](https://github.com/cs3org/reva/pull/4362): Fix concurrent lookup
  *   Bugfix [cs3org/reva#4336](https://github.com/cs3org/reva/pull/4336): Fix definition of "file-editor" role
  *   Bugfix [cs3org/reva#4302](https://github.com/cs3org/reva/pull/4302): Fix checking of filename length
  *   Bugfix [cs3org/reva#4366](https://github.com/cs3org/reva/pull/4366): Fix CS3 status code when looking up non existing share
  *   Bugfix [cs3org/reva#4299](https://github.com/cs3org/reva/pull/4299): Fix HTTP verb of the generate-invite endpoint
  *   Bugfix [cs3org/reva#4249](https://github.com/cs3org/reva/pull/4249): GetUserByClaim not working with MSAD for claim "userid"
  *   Bugfix [cs3org/reva#4217](https://github.com/cs3org/reva/pull/4217): Fix missing case for "hide" in UpdateShares
  *   Bugfix [cs3org/reva#4140](https://github.com/cs3org/reva/pull/4140): Fix missing etag in shares jail
  *   Bugfix [cs3org/reva#4229](https://github.com/cs3org/reva/pull/4229): Fix destroying the Personal and Project spaces data
  *   Bugfix [cs3org/reva#4193](https://github.com/cs3org/reva/pull/4193): Fix overwrite a file with an empty file
  *   Bugfix [cs3org/reva#4365](https://github.com/cs3org/reva/pull/4365): Fix create public share
  *   Bugfix [cs3org/reva#4380](https://github.com/cs3org/reva/pull/4380): Fix the public link update
  *   Bugfix [cs3org/reva#4250](https://github.com/cs3org/reva/pull/4250): Fix race condition
  *   Bugfix [cs3org/reva#4345](https://github.com/cs3org/reva/pull/4345): Fix conversion of custom ocs permissions to roles
  *   Bugfix [cs3org/reva#4134](https://github.com/cs3org/reva/pull/4134): Fix share jail
  *   Bugfix [cs3org/reva#4335](https://github.com/cs3org/reva/pull/4335): Fix public shares cleanup config
  *   Bugfix [cs3org/reva#4338](https://github.com/cs3org/reva/pull/4338): Fix unlock via space API
  *   Bugfix [cs3org/reva#4341](https://github.com/cs3org/reva/pull/4341): Fix spaceID in meta endpoint response
  *   Bugfix [cs3org/reva#4351](https://github.com/cs3org/reva/pull/4351): Fix 500 when open public link
  *   Bugfix [cs3org/reva#4352](https://github.com/cs3org/reva/pull/4352): Fix the tgz mime type
  *   Bugfix [cs3org/reva#4388](https://github.com/cs3org/reva/pull/4388): Allow UpdateUserShare() to update just the expiration date
  *   Bugfix [cs3org/reva#4214](https://github.com/cs3org/reva/pull/4214): Always pass adjusted default nats options
  *   Bugfix [cs3org/reva#4291](https://github.com/cs3org/reva/pull/4291): Release lock when expired
  *   Bugfix [cs3org/reva#4386](https://github.com/cs3org/reva/pull/4386): Remove dead enable_home config
  *   Bugfix [cs3org/reva#4292](https://github.com/cs3org/reva/pull/4292): Return 403 when user is not permitted to lock
  *   Enhancement [cs3org/reva#4389](https://github.com/cs3org/reva/pull/4389): Add audio and location props
  *   Enhancement [cs3org/reva#4337](https://github.com/cs3org/reva/pull/4337): Check permissions before creating shares
  *   Enhancement [cs3org/reva#4326](https://github.com/cs3org/reva/pull/4326): Add search mediatype filter
  *   Enhancement [cs3org/reva#4367](https://github.com/cs3org/reva/pull/4367): Add GGS mime type
  *   Enhancement [cs3org/reva#4194](https://github.com/cs3org/reva/pull/4194): Add hide flag to shares
  *   Enhancement [cs3org/reva#4358](https://github.com/cs3org/reva/pull/4358): Add default permissions capability for links
  *   Enhancement [cs3org/reva#4133](https://github.com/cs3org/reva/pull/4133): Add more metadata to locks
  *   Enhancement [cs3org/reva#4353](https://github.com/cs3org/reva/pull/4353): Add support for .docxf files
  *   Enhancement [cs3org/reva#4363](https://github.com/cs3org/reva/pull/4363): Add nats-js-kv store
  *   Enhancement [cs3org/reva#4197](https://github.com/cs3org/reva/pull/4197): Add the Banned-Passwords List
  *   Enhancement [cs3org/reva#4190](https://github.com/cs3org/reva/pull/4190): Add the password policies
  *   Enhancement [cs3org/reva#4384](https://github.com/cs3org/reva/pull/4384): Add a retry postprocessing outcome and event
  *   Enhancement [cs3org/reva#4271](https://github.com/cs3org/reva/pull/4271): Add search capability
  *   Enhancement [cs3org/reva#4119](https://github.com/cs3org/reva/pull/4119): Add sse event
  *   Enhancement [cs3org/reva#4392](https://github.com/cs3org/reva/pull/4392): Add additional permissions to service accounts
  *   Enhancement [cs3org/reva#4344](https://github.com/cs3org/reva/pull/4344): Add url extension to mime type list
  *   Enhancement [cs3org/reva#4372](https://github.com/cs3org/reva/pull/4372): Add validation to the public share provider
  *   Enhancement [cs3org/reva#4244](https://github.com/cs3org/reva/pull/4244): Allow listing reveived shares by service accounts
  *   Enhancement [cs3org/reva#4129](https://github.com/cs3org/reva/pull/4129): Auto-Accept Shares through ServiceAccounts
  *   Enhancement [cs3org/reva#4374](https://github.com/cs3org/reva/pull/4374): Handle trashbin file listings concurrently
  *   Enhancement [cs3org/reva#4325](https://github.com/cs3org/reva/pull/4325): Enforce Permissions
  *   Enhancement [cs3org/reva#4368](https://github.com/cs3org/reva/pull/4368): Extract log initialization
  *   Enhancement [cs3org/reva#4375](https://github.com/cs3org/reva/pull/4375): Introduce UploadSessionLister interface
  *   Enhancement [cs3org/reva#4268](https://github.com/cs3org/reva/pull/4268): Implement sharing roles
  *   Enhancement [cs3org/reva#4160](https://github.com/cs3org/reva/pull/4160): Improve utils pkg
  *   Enhancement [cs3org/reva#4335](https://github.com/cs3org/reva/pull/4335): Add sufficient permissions check function
  *   Enhancement [cs3org/reva#4281](https://github.com/cs3org/reva/pull/4281): Port OCM changes from master
  *   Enhancement [cs3org/reva#4270](https://github.com/cs3org/reva/pull/4270): Opt out of public link password enforcement
  *   Enhancement [cs3org/reva#4181](https://github.com/cs3org/reva/pull/4181): Change the variable names for the password policy
  *   Enhancement [cs3org/reva#4256](https://github.com/cs3org/reva/pull/4256): Rename hidden share variable name
  *   Enhancement [cs3org/reva#3926](https://github.com/cs3org/reva/pull/3926): Service Accounts
  *   Enhancement [cs3org/reva#4359](https://github.com/cs3org/reva/pull/4359): Update go-ldap to v3.4.6
  *   Enhancement [cs3org/reva#4170](https://github.com/cs3org/reva/pull/4170): Update password policies
  *   Enhancement [cs3org/reva#4232](https://github.com/cs3org/reva/pull/4232): Improve error handling in utils package

   https://github.com/owncloud/ocis/pull/8638
   https://github.com/owncloud/ocis/pull/8519
   https://github.com/owncloud/ocis/pull/8502
   https://github.com/owncloud/ocis/pull/8340
   https://github.com/owncloud/ocis/pull/8381
   https://github.com/owncloud/ocis/pull/8287
   https://github.com/owncloud/ocis/pull/8278
   https://github.com/owncloud/ocis/pull/8264
   https://github.com/owncloud/ocis/pull/8100
   https://github.com/owncloud/ocis/pull/8100
   https://github.com/owncloud/ocis/pull/8038
   https://github.com/owncloud/ocis/pull/8056
   https://github.com/owncloud/ocis/pull/7949
   https://github.com/owncloud/ocis/pull/7793
   https://github.com/owncloud/ocis/pull/7978
   https://github.com/owncloud/ocis/pull/7979
   https://github.com/owncloud/ocis/pull/7963
   https://github.com/owncloud/ocis/pull/7986
   https://github.com/owncloud/ocis/pull/7721
   https://github.com/owncloud/ocis/pull/7727
   https://github.com/owncloud/ocis/pull/7752

# Changelog for [4.0.6] (2024-02-07)

The following sections list the changes for 4.0.6.

[4.0.6]: https://github.com/owncloud/ocis/compare/v4.0.5...v4.0.6

## Summary

* Bugfix - Fix RED metrics on the metrics endpoint: [#7994](https://github.com/owncloud/ocis/pull/7994)
* Bugfix - Signed url verification: [#8385](https://github.com/owncloud/ocis/pull/8385)

## Details

* Bugfix - Fix RED metrics on the metrics endpoint: [#7994](https://github.com/owncloud/ocis/pull/7994)

   We connected some metrics to the metrics endpoint to support the RED method for
   monitoring microservices.

   - Request Rate: The number of requests per second. The total count of requests
   is available under `ocis_proxy_requests_total`. - Error Rate: The number of
   failed requests per second. The total count of failed requests is available
   under `ocis_proxy_errors_total`. - Duration: The amount of time each request
   takes. The duration of all requests is available under
   `ocis_proxy_request_duration_seconds`. This is a histogram metric, so it also
   provides information about the distribution of request durations.

   The metrics are available under the following paths: `PROXY_DEBUG_ADDR/metrics`
   in a prometheus compatible format and maybe secured by `PROXY_DEBUG_TOKEN`.

   https://github.com/owncloud/ocis/pull/7994

* Bugfix - Signed url verification: [#8385](https://github.com/owncloud/ocis/pull/8385)

   Signed urls now expire properly

   https://github.com/owncloud/ocis/pull/8385

# Changelog for [4.0.5] (2023-12-21)

The following sections list the changes for 4.0.5.

[4.0.5]: https://github.com/owncloud/ocis/compare/v4.0.4...v4.0.5

## Summary

* Bugfix - Fix reva config of frontend service to avoid misleading error logs: [#7934](https://github.com/owncloud/ocis/pull/7934)
* Bugfix - Do not purge expired upload sessions that are still postprocessing: [#7941](https://github.com/owncloud/ocis/pull/7941)
* Bugfix - Fix trace ids: [#8026](https://github.com/owncloud/ocis/pull/8026)
* Enhancement - Add cli commands for trash-bin: [#7936](https://github.com/owncloud/ocis/pull/7936)

## Details

* Bugfix - Fix reva config of frontend service to avoid misleading error logs: [#7934](https://github.com/owncloud/ocis/pull/7934)

   We set an empty Credentials chain for the frontend service now. In ocis all
   non-reva token authentication is handled by the proxy. This avoids irritating
   error messages about the missing 'auth-bearer' service.

   https://github.com/owncloud/ocis/issues/6692
   https://github.com/owncloud/ocis/pull/7934
   https://github.com/owncloud/ocis/pull/7453
   https://github.com/cs3org/reva/pull/4396
   https://github.com/cs3org/reva/pull/4241

* Bugfix - Do not purge expired upload sessions that are still postprocessing: [#7941](https://github.com/owncloud/ocis/pull/7941)

   https://github.com/owncloud/ocis/pull/7941
   https://github.com/owncloud/ocis/pull/7859
   https://github.com/owncloud/ocis/pull/7958

* Bugfix - Fix trace ids: [#8026](https://github.com/owncloud/ocis/pull/8026)

   We changed the default tracing to produce non-empty traceids and fixed a problem
   where traces got disconnected further down the stack.

   https://github.com/owncloud/ocis/pull/8026

* Enhancement - Add cli commands for trash-bin: [#7936](https://github.com/owncloud/ocis/pull/7936)

   We added the `list` and `restore` commands to the trash-bin items to the CLI

   https://github.com/owncloud/ocis/issues/7845
   https://github.com/owncloud/ocis/pull/7936

# Changelog for [4.0.4] (2023-12-07)

The following sections list the changes for 4.0.4.

[4.0.4]: https://github.com/owncloud/ocis/compare/v4.0.3...v4.0.4

## Summary

* Enhancement - Update reva to improve trashbin listing: [#7858](https://github.com/owncloud/ocis/pull/7858)

## Details

* Enhancement - Update reva to improve trashbin listing: [#7858](https://github.com/owncloud/ocis/pull/7858)

   ## Changelog for reva 2.13.3

  *   Enhancement [cs3org/reva#4377](https://github.com/cs3org/reva/pull/4377): Handle trashbin file listings concurrently

   https://github.com/owncloud/ocis/pull/7858

# Changelog for [4.0.3] (2023-11-24)

The following sections list the changes for 4.0.3.

[4.0.3]: https://github.com/owncloud/ocis/compare/v4.0.2...v4.0.3

## Summary

* Bugfix - Bump reva to 2.16.2: [#7512](https://github.com/owncloud/ocis/pull/7512)
* Bugfix - Token storage config fixed: [#7546](https://github.com/owncloud/ocis/pull/7546)
* Enhancement - Support spec violating AD FS access token issuer: [#7138](https://github.com/owncloud/ocis/pull/7138)
* Enhancement - Update web to v7.1.2: [#7798](https://github.com/owncloud/ocis/pull/7798)

## Details

* Bugfix - Bump reva to 2.16.2: [#7512](https://github.com/owncloud/ocis/pull/7512)

  *   Bugfix [cs3org/reva#4251](https://github.com/cs3org/reva/pull/4251): ldap: fix GetUserByClaim for binary encoded UUIDs

   https://github.com/owncloud/ocis/issues/7469
   https://github.com/owncloud/ocis/pull/7512

* Bugfix - Token storage config fixed: [#7546](https://github.com/owncloud/ocis/pull/7546)

   The token storage config in the config.json for web was missing when it was set
   to `false`.

   https://github.com/owncloud/ocis/issues/7462
   https://github.com/owncloud/ocis/pull/7546

* Enhancement - Support spec violating AD FS access token issuer: [#7138](https://github.com/owncloud/ocis/pull/7138)

   AD FS `/adfs/.well-known/openid-configuration` has an optional
   `access_token_issuer` which, in violation of the OpenID Connect spec, takes
   precedence over `issuer`.

   https://github.com/owncloud/ocis/pull/7138

* Enhancement - Update web to v7.1.2: [#7798](https://github.com/owncloud/ocis/pull/7798)

   Tags: web

   We updated ownCloud Web to v7.1.2. Please refer to the changelog (linked) for
   details on the web release.

   ## Summary * Bugfix
   [owncloud/web#9833](https://github.com/owncloud/web/pull/9833): Resolving
   external URLs * Bugfix
   [owncloud/web#9868](https://github.com/owncloud/web/pull/9868): Respect
   "details"-query on private links * Bugfix
   [owncloud/web#9913](https://github.com/owncloud/web/pull/9913): Private link
   resolving via share jail ID

   https://github.com/owncloud/ocis/pull/7798
   https://github.com/owncloud/web/releases/tag/v7.1.2

# Changelog for [4.0.2] (2023-09-28)

The following sections list the changes for 4.0.2.

[4.0.2]: https://github.com/owncloud/ocis/compare/v4.0.1...v4.0.2

## Summary

* Bugfix - Actually pass PROXY_OIDC_SKIP_USER_INFO option to oidc client middleware: [#7220](https://github.com/owncloud/ocis/pull/7220)
* Bugfix - Disable username validation for keycloak example: [#7230](https://github.com/owncloud/ocis/pull/7230)
* Bugfix - Bring back the USERS_LDAP_USER_SCHEMA_ID variable: [#7312](https://github.com/owncloud/ocis/issues/7312)
* Bugfix - Do not reset received share state to pending: [#7319](https://github.com/owncloud/ocis/issues/7319)
* Bugfix - Bump reva to 2.16.1: [#7350](https://github.com/owncloud/ocis/pull/7350)
* Bugfix - Check school number for duplicates before adding a school: [#7351](https://github.com/owncloud/ocis/pull/7351)
* Enhancement - Add OCIS_LDAP_BIND_PASSWORD as replacement for LDAP_BIND_PASSWORD: [#7176](https://github.com/owncloud/ocis/issues/7176)

## Details

* Bugfix - Actually pass PROXY_OIDC_SKIP_USER_INFO option to oidc client middleware: [#7220](https://github.com/owncloud/ocis/pull/7220)

   https://github.com/owncloud/ocis/pull/7220

* Bugfix - Disable username validation for keycloak example: [#7230](https://github.com/owncloud/ocis/pull/7230)

   Set 'GRAPH_USERNAME_MATCH' to 'none'. To accept any username that is also valid
   for keycloak.

   https://github.com/owncloud/ocis/pull/7230

* Bugfix - Bring back the USERS_LDAP_USER_SCHEMA_ID variable: [#7312](https://github.com/owncloud/ocis/issues/7312)

   We reintroduced the USERS_LDAP_USER_SCHEMA_ID variable which was accidently
   removed from the users service with the 4.0.0 release.

   https://github.com/owncloud/ocis/issues/7312
   https://github.com/owncloud/ocis-charts/issues/397

* Bugfix - Do not reset received share state to pending: [#7319](https://github.com/owncloud/ocis/issues/7319)

   We fixed a problem where the states of received shares were reset to PENDING in
   the "ocis migrate rebuild-jsoncs3-indexes" command

   https://github.com/owncloud/ocis/issues/7319

* Bugfix - Bump reva to 2.16.1: [#7350](https://github.com/owncloud/ocis/pull/7350)

  *   Bugfix [cs3org/reva#4194](https://github.com/cs3org/reva/pull/4194): Make appctx package compatible with go v1.21
  *   Bugfix [cs3org/reva#4214](https://github.com/cs3org/reva/pull/4214): Always pass adjusted default nats options

   https://github.com/owncloud/ocis/pull/7350

* Bugfix - Check school number for duplicates before adding a school: [#7351](https://github.com/owncloud/ocis/pull/7351)

   We fixed an issue that allowed to create two schools with the same school number

   https://github.com/owncloud/enterprise/issues/6051
   https://github.com/owncloud/ocis/pull/7351

* Enhancement - Add OCIS_LDAP_BIND_PASSWORD as replacement for LDAP_BIND_PASSWORD: [#7176](https://github.com/owncloud/ocis/issues/7176)

   The enviroment variable `OCIS_LDAP_BIND_PASSWORD` was added to be more
   consistent with all other global LDAP variables.

   `LDAP_BIND_PASSWORD` is deprecated now and scheduled for removal with the 5.0.0
   release.

   We also deprecated `LDAP_USER_SCHEMA_ID_IS_OCTETSTRING` for removal with 5.0.0.
   The replacement for it is `OCIS_LDAP_USER_SCHEMA_ID_IS_OCTETSTRING`.

   https://github.com/owncloud/ocis/issues/7176

# Changelog for [4.0.1] (2023-09-01)

The following sections list the changes for 4.0.1.

[4.0.1]: https://github.com/owncloud/ocis/compare/v4.0.0...v4.0.1

## Summary

* Bugfix - Disallow sharee to search sharer files outside the share: [#7184](https://github.com/owncloud/ocis/pull/7184)

## Details

* Bugfix - Disallow sharee to search sharer files outside the share: [#7184](https://github.com/owncloud/ocis/pull/7184)

   When a file was shared with user(sharee) and the sharee searched the shared file
   the response contained unshared resources as well.

   https://github.com/owncloud/ocis/pull/7184

# Changelog for [4.0.0] (2023-08-21)

The following sections list the changes for 4.0.0.

[4.0.0]: https://github.com/owncloud/ocis/compare/v3.0.0...v4.0.0

## Summary

* Bugfix - Fix error message on 400 response for thumbnail requests: [#2064](https://github.com/owncloud/ocis/issues/2064)
* Bugfix - Handle the bad request status: [#6469](https://github.com/owncloud/ocis/pull/6469)
* Bugfix - Add missing timestamps: [#6515](https://github.com/owncloud/ocis/pull/6515)
* Bugfix - Add token to LinkAccessedEvent: [#6554](https://github.com/owncloud/ocis/pull/6554)
* Bugfix - Don't connect to ldap on startup: [#6565](https://github.com/owncloud/ocis/pull/6565)
* Bugfix - Add default store to postprocessing: [#6578](https://github.com/owncloud/ocis/pull/6578)
* Bugfix - Fix the oidc role assigner: [#6605](https://github.com/owncloud/ocis/pull/6605)
* Bugfix - Restart Postprocessing: [#6726](https://github.com/owncloud/ocis/pull/6726)
* Bugfix - Fix search shares: [#6741](https://github.com/owncloud/ocis/pull/6741)
* Bugfix - Fix the default document language for OnlyOffice: [#6878](https://github.com/owncloud/ocis/pull/6878)
* Bugfix - Fix nats registry: [#6881](https://github.com/owncloud/ocis/pull/6881)
* Bugfix - Check public auth first: [#6900](https://github.com/owncloud/ocis/pull/6900)
* Bugfix - Fix CORS issues: [#6912](https://github.com/owncloud/ocis/pull/6912)
* Bugfix - Let clients cache web and theme assets: [#6914](https://github.com/owncloud/ocis/pull/6914)
* Bugfix - Fix the search: [#6947](https://github.com/owncloud/ocis/pull/6947)
* Bugfix - Graph service did not honor the OCIS_LDAP_GROUP_SCHEMA_MEMBER setting: [#7032](https://github.com/owncloud/ocis/issues/7032)
* Bugfix - Fix the routing capability: [#9367](https://github.com/owncloud/web/issues/9367)
* Change - YAML configuration files are restricted to yaml-1.2: [#6510](https://github.com/owncloud/ocis/issues/6510)
* Enhancement - Add SSE Endpoint: [#5998](https://github.com/owncloud/ocis/pull/5998)
* Enhancement - Add postprocessing mimetype to extension helper: [#6133](https://github.com/owncloud/ocis/pull/6133)
* Enhancement - Add more metadata to the remote item: [#6300](https://github.com/owncloud/ocis/pull/6300)
* Enhancement - Add WEB_OPTION_OPEN_LINKS_WITH_DEFAULT_APP env variable: [#6328](https://github.com/owncloud/ocis/pull/6328)
* Enhancement - Fix the username validation: [#6437](https://github.com/owncloud/ocis/pull/6437)
* Enhancement - Use reva client selectors: [#6452](https://github.com/owncloud/ocis/pull/6452)
* Enhancement - Add companion URL config: [#6453](https://github.com/owncloud/ocis/pull/6453)
* Enhancement - Update go-micro kubernetes registry: [#6457](https://github.com/owncloud/ocis/pull/6457)
* Enhancement - Add imprint and privacy url config: [#6462](https://github.com/owncloud/ocis/pull/6462)
* Enhancement - Update web to v7.0.1: [#6470](https://github.com/owncloud/ocis/pull/6470)
* Enhancement - Make the app provider service name configurable: [#6482](https://github.com/owncloud/ocis/pull/6482)
* Enhancement - Fix the groupname validation: [#6490](https://github.com/owncloud/ocis/pull/6490)
* Enhancement - Add functionality to retry postprocessing: [#6500](https://github.com/owncloud/ocis/pull/6500)
* Enhancement - Fix envvar defaults: [#6516](https://github.com/owncloud/ocis/pull/6516)
* Enhancement - Add permissions to report: [#6528](https://github.com/owncloud/ocis/pull/6528)
* Enhancement - Add old & new values to audit logs: [#6537](https://github.com/owncloud/ocis/pull/6537)
* Enhancement - Allow disabling wopi chat: [#6544](https://github.com/owncloud/ocis/pull/6544)
* Enhancement - We added the storage id to the audit log for spaces: [#6548](https://github.com/owncloud/ocis/pull/6548)
* Enhancement - Add logged out url config: [#6549](https://github.com/owncloud/ocis/pull/6549)
* Enhancement - Add 'ocis decomposedfs check-treesize' command: [#6556](https://github.com/owncloud/ocis/pull/6556)
* Enhancement - Skip if the simulink is a directory: [#6574](https://github.com/owncloud/ocis/pull/6574)
* Enhancement - Thumbnails can be disabled for webdav & web now: [#6577](https://github.com/owncloud/ocis/pull/6577)
* Enhancement - Make the post logout redirect uri configurable: [#6583](https://github.com/owncloud/ocis/pull/6583)
* Enhancement - Move proxy to service tracerprovider: [#6591](https://github.com/owncloud/ocis/pull/6591)
* Enhancement - Add IDs to graph resource logging: [#6593](https://github.com/owncloud/ocis/pull/6593)
* Enhancement - Add search result content preview and term highlighting: [#6634](https://github.com/owncloud/ocis/pull/6634)
* Enhancement - Move graph to service tracerprovider: [#6695](https://github.com/owncloud/ocis/pull/6695)
* Enhancement - Provide Search filter for locations: [#6713](https://github.com/owncloud/ocis/pull/6713)
* Enhancement - Add X-Request-Id to all responses: [#6715](https://github.com/owncloud/ocis/pull/6715)
* Enhancement - Clarify license text in the dev docs: [#6755](https://github.com/owncloud/ocis/pull/6755)
* Enhancement - Add WEB_OPTION_TOKEN_STORAGE_LOCAL env variable: [#6760](https://github.com/owncloud/ocis/pull/6760)
* Enhancement - Bump Hugo: [#6787](https://github.com/owncloud/ocis/pull/6787)
* Enhancement - Bump reva to 2.16.0: [#6829](https://github.com/owncloud/ocis/pull/6829)
* Enhancement - Configure max grpc message size: [#6849](https://github.com/owncloud/ocis/pull/6849)
* Enhancement - Improve the notification logs: [#6862](https://github.com/owncloud/ocis/pull/6862)
* Enhancement - Extendable policy mimetype extension mapping: [#6869](https://github.com/owncloud/ocis/pull/6869)
* Enhancement - Evaluate policy resource information on single file shares: [#6888](https://github.com/owncloud/ocis/pull/6888)
* Enhancement - Update web to v7.1.0-rc.5: [#6944](https://github.com/owncloud/ocis/pull/6944)
* Enhancement - Add static secret to gn endpoints: [#6946](https://github.com/owncloud/ocis/pull/6946)
* Enhancement - Bump sonarcloud: [#6961](https://github.com/owncloud/ocis/pull/6961)
* Enhancement - Nats named connections: [#6979](https://github.com/owncloud/ocis/pull/6979)
* Enhancement - Add command for rebuilding the jsoncs3 share manager indexes: [#6986](https://github.com/owncloud/ocis/pull/6986)
* Enhancement - Remove deprecated environment variables: [#7099](https://github.com/owncloud/ocis/pull/7099)
* Enhancement - Update web to v7.1.0: [#7107](https://github.com/owncloud/ocis/pull/7107)

## Details

* Bugfix - Fix error message on 400 response for thumbnail requests: [#2064](https://github.com/owncloud/ocis/issues/2064)

   Fix the error message when the thumbnail request returns a '400 Bad Request'
   response.

   https://github.com/owncloud/ocis/issues/2064
   https://github.com/owncloud/ocis/pull/6911

* Bugfix - Handle the bad request status: [#6469](https://github.com/owncloud/ocis/pull/6469)

   Handle the bad request status for the CreateStorageSpace function

   https://github.com/owncloud/ocis/issues/6414
   https://github.com/owncloud/ocis/pull/6469
   https://github.com/cs3org/reva/pull/3948

* Bugfix - Add missing timestamps: [#6515](https://github.com/owncloud/ocis/pull/6515)

   We have added missing timestamps to the audit service

   https://github.com/owncloud/ocis/issues/3753
   https://github.com/owncloud/ocis/pull/6515

* Bugfix - Add token to LinkAccessedEvent: [#6554](https://github.com/owncloud/ocis/pull/6554)

   We added the link token to the LinkAccessedEvent

   https://github.com/owncloud/ocis/issues/3753
   https://github.com/owncloud/ocis/pull/6554
   https://github.com/cs3org/reva/pull/3993

* Bugfix - Don't connect to ldap on startup: [#6565](https://github.com/owncloud/ocis/pull/6565)

   This leads to misleading error messages. Instead we connect on first request

   https://github.com/owncloud/ocis/pull/6565

* Bugfix - Add default store to postprocessing: [#6578](https://github.com/owncloud/ocis/pull/6578)

   Postprocessing did not have a default store especially `database` and `table`
   are needed to talk to nats-js

   https://github.com/owncloud/ocis/pull/6578

* Bugfix - Fix the oidc role assigner: [#6605](https://github.com/owncloud/ocis/pull/6605)

   The update role method did not allow to set a role when the user already has two
   roles. This makes no sense as the user is supposed to have only one and the
   update will fix that. We still log an error level log to make the admin aware of
   that.

   https://github.com/owncloud/ocis/pull/6605
   https://github.com/owncloud/ocis/pull/6618

* Bugfix - Restart Postprocessing: [#6726](https://github.com/owncloud/ocis/pull/6726)

   In case the postprocessing service cannot find the specified upload when
   restarting postprocessing, it will now send a `RestartPostprocessing` event to
   retrigger complete postprocessing

   https://github.com/owncloud/ocis/pull/6726

* Bugfix - Fix search shares: [#6741](https://github.com/owncloud/ocis/pull/6741)

   We fixed a problem where searching shares did not yield results when the
   resource was not shared from the space root.

   https://github.com/owncloud/ocis/pull/6741

* Bugfix - Fix the default document language for OnlyOffice: [#6878](https://github.com/owncloud/ocis/pull/6878)

   Fix the default document language for OnlyOffice

   https://github.com/owncloud/enterprise/issues/5807
   https://github.com/owncloud/ocis/pull/6878

* Bugfix - Fix nats registry: [#6881](https://github.com/owncloud/ocis/pull/6881)

   Using `nats` as service registry did work, but when a service would restart and
   gets a new ip it couldn't re-register. We fixed this by using `"put"` register
   action instead of the default `"create"`

   https://github.com/owncloud/ocis/pull/6881

* Bugfix - Check public auth first: [#6900](https://github.com/owncloud/ocis/pull/6900)

   When authenticating in proxy, first check for public link authorization.

   https://github.com/owncloud/ocis/pull/6900

* Bugfix - Fix CORS issues: [#6912](https://github.com/owncloud/ocis/pull/6912)

   We fixed the CORS issues when client asking for the 'Cache-Control' header
   before load the file

   https://github.com/owncloud/ocis/issues/5108
   https://github.com/owncloud/ocis/pull/6912

* Bugfix - Let clients cache web and theme assets: [#6914](https://github.com/owncloud/ocis/pull/6914)

   We needed to remove "must-revalidate" from the cache-control header to allow
   clients to cache the web and theme assets.

   https://github.com/owncloud/ocis/pull/6914

* Bugfix - Fix the search: [#6947](https://github.com/owncloud/ocis/pull/6947)

   We fixed the issue when search using the current folder option shows the
   file/folders outside the folder if search keyword is same as current folder

   https://github.com/owncloud/ocis/issues/6935
   https://github.com/owncloud/ocis/pull/6947

* Bugfix - Graph service did not honor the OCIS_LDAP_GROUP_SCHEMA_MEMBER setting: [#7032](https://github.com/owncloud/ocis/issues/7032)

   We fixed issue when using a custom LDAP attribute for group members. The graph
   service did not honor the OCIS_LDAP_GROUP_SCHEMA_MEMBER environment variable

   https://github.com/owncloud/ocis/issues/7032

* Bugfix - Fix the routing capability: [#9367](https://github.com/owncloud/web/issues/9367)

   Fix the routing capability

   https://github.com/owncloud/web/issues/9367

* Change - YAML configuration files are restricted to yaml-1.2: [#6510](https://github.com/owncloud/ocis/issues/6510)

   For parsing YAML based configuration files we utilize the gookit/config module.
   That module has dropped support for older variants of the YAML format. It now
   only supports the YAML 1.2 syntax. If you're using yaml configuration files,
   please make sure to update your files accordingly. The most significant change
   likely is that only the string `true` and `false` (including `TRUE`,`True`,
   `FALSE` and `False`) are now parsed as booleans. `Yes`, `On` and other values
   are not longer considered valid values for booleans.

   https://github.com/owncloud/ocis/issues/6510
   https://github.com/owncloud/ocis/pull/6493

* Enhancement - Add SSE Endpoint: [#5998](https://github.com/owncloud/ocis/pull/5998)

   Add a server-sent events (sse) endpoint for the userlog service

   https://github.com/owncloud/ocis/pull/5998

* Enhancement - Add postprocessing mimetype to extension helper: [#6133](https://github.com/owncloud/ocis/pull/6133)

   Add rego helper to resolve extensions from mimetype
   `ocis.mimetype.extensions(mimetype)`. Besides that, a rego print helper is
   included also `print("PRINT MESSAGE EXAMPLE")`

   https://github.com/owncloud/ocis/pull/6133

* Enhancement - Add more metadata to the remote item: [#6300](https://github.com/owncloud/ocis/pull/6300)

   We added the drive alias, the space name and the relative path to the remote
   item. This is needed to resolve shared files directly on the source space.

   https://github.com/owncloud/ocis/pull/6300

* Enhancement - Add WEB_OPTION_OPEN_LINKS_WITH_DEFAULT_APP env variable: [#6328](https://github.com/owncloud/ocis/pull/6328)

   We introduced the open file links with default app feature in web which is
   enabled by default, this is now configurable and can be disabled by setting the
   env `WEB_OPTION_OPEN_LINKS_WITH_DEFAULT_APP` to `false`.

   https://github.com/owncloud/ocis/pull/6328

* Enhancement - Fix the username validation: [#6437](https://github.com/owncloud/ocis/pull/6437)

   Fix the username validation when an admin update the user

   https://github.com/owncloud/ocis/issues/6436
   https://github.com/owncloud/ocis/pull/6437

* Enhancement - Use reva client selectors: [#6452](https://github.com/owncloud/ocis/pull/6452)

   Use reva client selectors instead of the static clients, this introduces the
   ocis service registry in reva. The service discovery now resolves reva services
   by name and the client selectors pick a random registered service node.

   https://github.com/owncloud/ocis/pull/6452
   https://github.com/cs3org/reva/pull/3939
   https://github.com/cs3org/reva/pull/3953

* Enhancement - Add companion URL config: [#6453](https://github.com/owncloud/ocis/pull/6453)

   Introduce a config to set the Uppy Companion URL via
   `WEB_OPTION_UPLOAD_COMPANION_URL`.

   https://github.com/owncloud/ocis/pull/6453

* Enhancement - Update go-micro kubernetes registry: [#6457](https://github.com/owncloud/ocis/pull/6457)

   https://github.com/owncloud/ocis/pull/6457
   https://github.com/go-micro/plugins/pull/114
   https://github.com/go-micro/plugins/pull/113

* Enhancement - Add imprint and privacy url config: [#6462](https://github.com/owncloud/ocis/pull/6462)

   Introduce a config to set the imprint and privacy url via
   `WEB_OPTION_IMPRINT_URL` and `WEB_OPTION_PRIVACY_URL`.

   https://github.com/owncloud/ocis/pull/6462

* Enhancement - Update web to v7.0.1: [#6470](https://github.com/owncloud/ocis/pull/6470)

   Tags: web

   We updated ownCloud Web to v7.0.1. Please refer to the changelog (linked) for
   details on the web release.

   ## Summary * Bugfix
   [owncloud/web#9153](https://github.com/owncloud/web/pull/9153): Reduce space
   preloading

   https://github.com/owncloud/ocis/pull/6470
   https://github.com/owncloud/web/releases/tag/v7.0.1

* Enhancement - Make the app provider service name configurable: [#6482](https://github.com/owncloud/ocis/pull/6482)

   We needed to make the service name of the app provider configurable. This needs
   to be changed when using more than one app provider. Each of them needs be found
   by a unique service name. Possible examples are: `app-provider-collabora`,
   `app-provider-onlyoffice`, `app-provider-office365`.

   https://github.com/owncloud/ocis/pull/6482

* Enhancement - Fix the groupname validation: [#6490](https://github.com/owncloud/ocis/pull/6490)

   Fixed the ability to create a group with an empty name

   https://github.com/owncloud/ocis/issues/5050
   https://github.com/owncloud/ocis/pull/6490

* Enhancement - Add functionality to retry postprocessing: [#6500](https://github.com/owncloud/ocis/pull/6500)

   Adds a ctl command to manually retry failed postprocessing on uploads

   https://github.com/owncloud/ocis/pull/6500

* Enhancement - Fix envvar defaults: [#6516](https://github.com/owncloud/ocis/pull/6516)

   Defaults for the envvar OCIS_LDAP_DISABLE_USER_MECHANISM were not used
   consistently, correct is `attribute`.

   https://github.com/owncloud/ocis/issues/6513
   https://github.com/owncloud/ocis/pull/6516

* Enhancement - Add permissions to report: [#6528](https://github.com/owncloud/ocis/pull/6528)

   The webdav REPORT endpoint only returned permissions for personal spaces and
   shares. Now also for project spaces.

   https://github.com/owncloud/ocis/pull/6528

* Enhancement - Add old & new values to audit logs: [#6537](https://github.com/owncloud/ocis/pull/6537)

   We have added old & new values to the audit logs We have added the missing
   events for role changes

   https://github.com/owncloud/ocis/pull/6537

* Enhancement - Allow disabling wopi chat: [#6544](https://github.com/owncloud/ocis/pull/6544)

   Add a configreva for the new reva disable-chat feature

   https://github.com/owncloud/ocis/pull/6544

* Enhancement - We added the storage id to the audit log for spaces: [#6548](https://github.com/owncloud/ocis/pull/6548)

   We added the storage id to the audit log for spaces

   https://github.com/owncloud/ocis/issues/3753
   https://github.com/owncloud/ocis/pull/6548

* Enhancement - Add logged out url config: [#6549](https://github.com/owncloud/ocis/pull/6549)

   Introduce a config to set the more button url on the access denied page in web
   via `WEB_OPTION_ACCESS_DENIED_HELP_URL`.

   https://github.com/owncloud/ocis/pull/6549

* Enhancement - Add 'ocis decomposedfs check-treesize' command: [#6556](https://github.com/owncloud/ocis/pull/6556)

   We added a 'ocis decomposedfs check-treesize' command for checking (and
   reparing) the treesize metadata of a storage space.

   https://github.com/owncloud/ocis/pull/6556

* Enhancement - Skip if the simulink is a directory: [#6574](https://github.com/owncloud/ocis/pull/6574)

   Skip the error if the simulink is pointed to a directory

   https://github.com/owncloud/ocis/issues/6567
   https://github.com/owncloud/ocis/pull/6574

* Enhancement - Thumbnails can be disabled for webdav & web now: [#6577](https://github.com/owncloud/ocis/pull/6577)

   We added an env var `OCIS_DISABLE_PREVIEWS` to disable the thumbnails for web &
   webdav via a global setting. For each service this behaviour can be disabled
   using the local env vars `WEB_OPTION_DISABLE_PREVIEWS` (old) and
   `WEBDAV_DISABLE_PREVIEWS` (new).

   https://github.com/owncloud/ocis/issues/192
   https://github.com/owncloud/ocis/pull/6577

* Enhancement - Make the post logout redirect uri configurable: [#6583](https://github.com/owncloud/ocis/pull/6583)

   We added a config option to change the redirect uri after the logout action of
   the web client.

   https://github.com/owncloud/ocis/issues/6536
   https://github.com/owncloud/ocis/pull/6583

* Enhancement - Move proxy to service tracerprovider: [#6591](https://github.com/owncloud/ocis/pull/6591)

   This moves the proxy to initialise a service tracer provider at service
   initialisation time, instead of using a package global tracer provider.

   https://github.com/owncloud/ocis/pull/6591

* Enhancement - Add IDs to graph resource logging: [#6593](https://github.com/owncloud/ocis/pull/6593)

   Graph access logs were unsuable as they didn't contain IDs to match them to a
   request

   https://github.com/owncloud/ocis/pull/6593

* Enhancement - Add search result content preview and term highlighting: [#6634](https://github.com/owncloud/ocis/pull/6634)

   The search result REPORT response now contains a content preview which
   highlights the search term. The feature is only available if content extraction
   (e.g. apache tika) is configured

   https://github.com/owncloud/ocis/issues/6426
   https://github.com/owncloud/ocis/pull/6634

* Enhancement - Move graph to service tracerprovider: [#6695](https://github.com/owncloud/ocis/pull/6695)

   This moves the graph to initialise a service tracer provider at service
   initialisation time, instead of using a package global tracer provider.

   https://github.com/owncloud/ocis/pull/6695

* Enhancement - Provide Search filter for locations: [#6713](https://github.com/owncloud/ocis/pull/6713)

   The search result REPORT response now can be restricted the by the current
   folder via api (recursive) The scope needed for "current folder" (default is to
   search all available spaces) - part of the oc:pattern:"scope:<uuid> /Test"

   https://github.com/owncloud/ocis/pull/6713
   OCIS-3705

* Enhancement - Add X-Request-Id to all responses: [#6715](https://github.com/owncloud/ocis/pull/6715)

   We added the X-Request-Id to all responses to increase the debuggability of the
   platform.

   https://github.com/owncloud/ocis/pull/6715

* Enhancement - Clarify license text in the dev docs: [#6755](https://github.com/owncloud/ocis/pull/6755)

   Explain the usage of the EULA for binary builds.

   https://github.com/owncloud/ocis/pull/6755

* Enhancement - Add WEB_OPTION_TOKEN_STORAGE_LOCAL env variable: [#6760](https://github.com/owncloud/ocis/pull/6760)

   We introduced the feature to store the access token in the local storage, this
   feature is disabled by default, but can be enabled by setting the env
   `WEB_OPTION_TOKEN_STORAGE_LOCAL` to `true`.

   https://github.com/owncloud/ocis/pull/6760
   https://github.com/owncloud/ocis/pull/6771

* Enhancement - Bump Hugo: [#6787](https://github.com/owncloud/ocis/pull/6787)

   Bump hugo pkg (needed for docs generation) to `v0.115.2`

   https://github.com/owncloud/ocis/pull/6787

* Enhancement - Bump reva to 2.16.0: [#6829](https://github.com/owncloud/ocis/pull/6829)

  *   Bugfix [cs3org/reva#4086](https://github.com/cs3org/reva/pull/4086): Fix ocs status code for not enough permission response
  *   Bugfix [cs3org/reva#4078](https://github.com/cs3org/reva/pull/4078): fix the default document language for OnlyOffice
  *   Bugfix [cs3org/reva#4051](https://github.com/cs3org/reva/pull/4051): Set treesize when creating a storage space
  *   Bugfix [cs3org/reva#4089](https://github.com/cs3org/reva/pull/4089): Fix wrong import
  *   Bugfix [cs3org/reva#4082](https://github.com/cs3org/reva/pull/4082): Fix propfind permissions
  *   Bugfix [cs3org/reva#4076](https://github.com/cs3org/reva/pull/4076): Fix WebDAV permissions for space managers
  *   Bugfix [cs3org/reva#4078](https://github.com/cs3org/reva/pull/4078): fix the default document language for OnlyOffice
  *   Bugfix [cs3org/reva#4081](https://github.com/cs3org/reva/pull/4081): Propagate sizeDiff
  *   Bugfix [cs3org/reva#4051](https://github.com/cs3org/reva/pull/4051): Set treesize when creating a storage space
  *   Bugfix [cs3org/reva#4093](https://github.com/cs3org/reva/pull/4093): Fix the error handling
  *   Bugfix [cs3org/reva#4111](https://github.com/cs3org/reva/pull/4111): Return already exists error when child already exists
  *   Bugfix [cs3org/reva#4086](https://github.com/cs3org/reva/pull/4086): Fix ocs status code for not enough permission response
  *   Bugfix [cs3org/reva#4101](https://github.com/cs3org/reva/pull/4101): Make the jsoncs3 share manager indexes more robust
  *   Bugfix [cs3org/reva#4099](https://github.com/cs3org/reva/pull/4099): Fix logging upload errors
  *   Bugfix [cs3org/reva#4078](https://github.com/cs3org/reva/pull/4078): Fix the default document language for OnlyOffice
  *   Bugfix [cs3org/reva#4082](https://github.com/cs3org/reva/pull/4082): Fix propfind permissions
  *   Bugfix [cs3org/reva#4100](https://github.com/cs3org/reva/pull/4100): S3ng include md5 checksum on put
  *   Bugfix [cs3org/reva#4096](https://github.com/cs3org/reva/pull/4096): Fix the user shares list
  *   Bugfix [cs3org/reva#4076](https://github.com/cs3org/reva/pull/4076): Fix WebDAV permissions for space managers
  *   Bugfix [cs3org/reva#4117](https://github.com/cs3org/reva/pull/4117): Fix jsoncs3 atomic persistence
  *   Bugfix [cs3org/reva#4081](https://github.com/cs3org/reva/pull/4081): Propagate sizeDiff
  *   Bugfix [cs3org/reva#4091](https://github.com/cs3org/reva/pull/4091): Register WebDAV HTTP methods with chi
  *   Bugfix [cs3org/reva#4107](https://github.com/cs3org/reva/pull/4107): Return lock when requested
  *   Bugfix [cs3org/reva#4075](https://github.com/cs3org/reva/pull/4075): Revert 4065 - bypass proxy on upload
  *   Enhancement [cs3org/reva#4070](https://github.com/cs3org/reva/pull/4070): Selectable Propagators
  *   Enhancement [cs3org/reva#4074](https://github.com/cs3org/reva/pull/4074): Allow configuring the max size of grpc messages
  *   Enhancement [cs3org/reva#4085](https://github.com/cs3org/reva/pull/4085): Add registry refresh
  *   Enhancement [cs3org/reva#4090](https://github.com/cs3org/reva/pull/4090): Add Capability for sse
  *   Enhancement [cs3org/reva#4072](https://github.com/cs3org/reva/pull/4072): Allow to specify a shutdown timeout
  *   Enhancement [cs3org/reva#4083](https://github.com/cs3org/reva/pull/4083): Allow for rolling back migrations
  *   Enhancement [cs3org/reva#4014](https://github.com/cs3org/reva/pull/4014): En-/Disable DEPTH:inifinity in PROPFIND
  *   Enhancement [cs3org/reva#4089](https://github.com/cs3org/reva/pull/4089): Async propagation (experimental)
  *   Enhancement [cs3org/reva#4074](https://github.com/cs3org/reva/pull/4074): Allow configuring the max size of grpc messages
  *   Enhancement [cs3org/reva#4083](https://github.com/cs3org/reva/pull/4083): Allow for rolling back migrations
  *   Enhancement [cs3org/reva#4014](https://github.com/cs3org/reva/pull/4014): En-/Disable DEPTH:inifinity in PROPFIND
  *   Enhancement [cs3org/reva#4072](https://github.com/cs3org/reva/pull/4072): Allow to specify a shutdown timeout
  *   Enhancement [cs3org/reva#4103](https://github.com/cs3org/reva/pull/4103): Add .oform mimetype
  *   Enhancement [cs3org/reva#4098](https://github.com/cs3org/reva/pull/4098): Allow naming nats connections
  *   Enhancement [cs3org/reva#4085](https://github.com/cs3org/reva/pull/4085): Add registry refresh
  *   Enhancement [cs3org/reva#4097](https://github.com/cs3org/reva/pull/4097): Remove app ticker logs
  *   Enhancement [cs3org/reva#4090](https://github.com/cs3org/reva/pull/4090): Add Capability for sse
  *   Enhancement [cs3org/reva#4110](https://github.com/cs3org/reva/pull/4110): Tracing events propgation

   Https://github.com/owncloud/ocis/pull/6899
   https://github.com/owncloud/ocis/pull/6919
   https://github.com/owncloud/ocis/pull/6928
   https://github.com/owncloud/ocis/pull/6979

   Update reva to v2.15.0

  *   Bugfix [cs3org/reva#4004](https://github.com/cs3org/reva/pull/4004): Add path to public link POST
  *   Bugfix [cs3org/reva#3993](https://github.com/cs3org/reva/pull/3993): Add token to LinkAccessedEvent
  *   Bugfix [cs3org/reva#4007](https://github.com/cs3org/reva/pull/4007): Close archive writer properly
  *   Bugfix [cs3org/reva#3982](https://github.com/cs3org/reva/pull/3982): Fixed couple of smaller space lookup issues
  *   Bugfix [cs3org/reva#4003](https://github.com/cs3org/reva/pull/4003): Don't connect ldap on startup
  *   Bugfix [cs3org/reva#4032](https://github.com/cs3org/reva/pull/4032): Temporarily exclude ceph-iscsi when building revad-ceph image
  *   Bugfix [cs3org/reva#4042](https://github.com/cs3org/reva/pull/4042): Fix writing 0 byte msgpack metadata
  *   Bugfix [cs3org/reva#3970](https://github.com/cs3org/reva/pull/3970): Fix enforce-password issue
  *   Bugfix [cs3org/reva#4057](https://github.com/cs3org/reva/pull/4057): Properly handle not-found errors when getting a public share
  *   Bugfix [cs3org/reva#4048](https://github.com/cs3org/reva/pull/4048): Fix messagepack propagation
  *   Bugfix [cs3org/reva#4056](https://github.com/cs3org/reva/pull/4056): Fix destroys data destination when moving issue
  *   Bugfix [cs3org/reva#4012](https://github.com/cs3org/reva/pull/4012): Fix mtime if 0 size file uploaded
  *   Bugfix [cs3org/reva#4010](https://github.com/cs3org/reva/pull/4010): Omit spaceroot when archiving
  *   Bugfix [cs3org/reva#4047](https://github.com/cs3org/reva/pull/4047): Publish events synchrously
  *   Bugfix [cs3org/reva#4039](https://github.com/cs3org/reva/pull/4039): Restart Postprocessing
  *   Bugfix [cs3org/reva#3963](https://github.com/cs3org/reva/pull/3963): Treesize interger overflows
  *   Bugfix [cs3org/reva#3943](https://github.com/cs3org/reva/pull/3943): When removing metadata always use correct database and table
  *   Bugfix [cs3org/reva#3978](https://github.com/cs3org/reva/pull/3978): Decomposedfs no longer os.Stats when reading node metadata
  *   Bugfix [cs3org/reva#3959](https://github.com/cs3org/reva/pull/3959): Drop unnecessary stat
  *   Bugfix [cs3org/reva#3948](https://github.com/cs3org/reva/pull/3948): Handle the bad request status
  *   Bugfix [cs3org/reva#3955](https://github.com/cs3org/reva/pull/3955): Fix panic
  *   Bugfix [cs3org/reva#3977](https://github.com/cs3org/reva/pull/3977): Prevent direct access to trash items
  *   Bugfix [cs3org/reva#3933](https://github.com/cs3org/reva/pull/3933): Concurrently invalidate mtime cache in jsoncs3 share manager
  *   Bugfix [cs3org/reva#3985](https://github.com/cs3org/reva/pull/3985): Reduce jsoncs3 lock congestion
  *   Bugfix [cs3org/reva#3960](https://github.com/cs3org/reva/pull/3960): Add trace span details
  *   Bugfix [cs3org/reva#3951](https://github.com/cs3org/reva/pull/3951): Link context in metadata client
  *   Bugfix [cs3org/reva#3950](https://github.com/cs3org/reva/pull/3950): Use plain otel tracing in metadata client
  *   Bugfix [cs3org/reva#3975](https://github.com/cs3org/reva/pull/3975): Decomposedfs now resolves the parent without an os.Stat
  *   Change [cs3org/reva#3947](https://github.com/cs3org/reva/pull/3947): Bump golangci-lint to 1.51.2
  *   Change [cs3org/reva#3945](https://github.com/cs3org/reva/pull/3945): Revert golangci-lint back to 1.50.1
  *   Enhancement [cs3org/reva#3966](https://github.com/cs3org/reva/pull/3966): Add space metadata to ocs shares list
  *   Enhancement [cs3org/reva#3953](https://github.com/cs3org/reva/pull/3953): Client selector pool
  *   Enhancement [cs3org/reva#3941](https://github.com/cs3org/reva/pull/3941): Adding tracing for jsoncs3
  *   Enhancement [cs3org/reva#3965](https://github.com/cs3org/reva/pull/3965): ResumePostprocessing Event
  *   Enhancement [cs3org/reva#3981](https://github.com/cs3org/reva/pull/3981): We have updated the UserFeatureChangedEvent to reflect value changes
  *   Enhancement [cs3org/reva#3986](https://github.com/cs3org/reva/pull/3986): Allow disabling wopi chat
  *   Enhancement [cs3org/reva#4060](https://github.com/cs3org/reva/pull/4060): We added a go-micro based app-provider registry
  *   Enhancement [cs3org/reva#4013](https://github.com/cs3org/reva/pull/4013): Add new WebDAV permissions
  *   Enhancement [cs3org/reva#3987](https://github.com/cs3org/reva/pull/3987): Cache space indexes
  *   Enhancement [cs3org/reva#3973](https://github.com/cs3org/reva/pull/3973): More logging for metadata propagation
  *   Enhancement [cs3org/reva#4059](https://github.com/cs3org/reva/pull/4059): Improve space index performance
  *   Enhancement [cs3org/reva#3994](https://github.com/cs3org/reva/pull/3994): Load matching spaces concurrently
  *   Enhancement [cs3org/reva#4049](https://github.com/cs3org/reva/pull/4049): Do not invalidate filemetadata cache early
  *   Enhancement [cs3org/reva#4040](https://github.com/cs3org/reva/pull/4040): Allow to use external trace provider in micro service
  *   Enhancement [cs3org/reva#4019](https://github.com/cs3org/reva/pull/4019): Allow to use external trace provider
  *   Enhancement [cs3org/reva#4045](https://github.com/cs3org/reva/pull/4045): Log error message in grpc interceptor
  *   Enhancement [cs3org/reva#3989](https://github.com/cs3org/reva/pull/3989): Parallelization of jsoncs3 operations
  *   Enhancement [cs3org/reva#3809](https://github.com/cs3org/reva/pull/3809): Trace decomposedfs syscalls
  *   Enhancement [cs3org/reva#4067](https://github.com/cs3org/reva/pull/4067): Trace upload progress
  *   Enhancement [cs3org/reva#3887](https://github.com/cs3org/reva/pull/3887): Trace requests through datagateway
  *   Enhancement [cs3org/reva#4052](https://github.com/cs3org/reva/pull/4052): Update go-ldap to v3.4.5
  *   Enhancement [cs3org/reva#4065](https://github.com/cs3org/reva/pull/4065): Upload directly to dataprovider
  *   Enhancement [cs3org/reva#4046](https://github.com/cs3org/reva/pull/4046): Use correct tracer name
  *   Enhancement [cs3org/reva#3986](https://github.com/cs3org/reva/pull/3986): Allow disabling wopi chat writer properly

   https://github.com/owncloud/ocis/pull/6829
   https://github.com/owncloud/ocis/pull/6529
   https://github.com/owncloud/ocis/pull/6544
   https://github.com/owncloud/ocis/pull/6507
   https://github.com/owncloud/ocis/pull/6572
   https://github.com/owncloud/ocis/pull/6590
   https://github.com/owncloud/ocis/pull/6812

* Enhancement - Configure max grpc message size: [#6849](https://github.com/owncloud/ocis/pull/6849)

   Add a configuration option for the grpc max message size

   https://github.com/owncloud/ocis/pull/6849

* Enhancement - Improve the notification logs: [#6862](https://github.com/owncloud/ocis/pull/6862)

   Improve the notification logs when the user has no email address

   https://github.com/owncloud/ocis/issues/6855
   https://github.com/owncloud/ocis/pull/6862

* Enhancement - Extendable policy mimetype extension mapping: [#6869](https://github.com/owncloud/ocis/pull/6869)

   The extension mimetype mappings known from rego can now be extended. To do this,
   ocis must be informed where the mimetype file (apache mime.types file format) is
   located.

   `export POLICIES_ENGINE_MIMES=OCIS_CONFIG_DIR/mime.types`

   https://github.com/owncloud/ocis/pull/6869

* Enhancement - Evaluate policy resource information on single file shares: [#6888](https://github.com/owncloud/ocis/pull/6888)

   The policy environment for single file shares now also includes information
   about the resource. As a result, it is now possible to set up and check rules
   for them.

   https://github.com/owncloud/ocis/pull/6888

* Enhancement - Update web to v7.1.0-rc.5: [#6944](https://github.com/owncloud/ocis/pull/6944)

   Tags: web

   We updated ownCloud Web to v7.1.0-rc.5. Please refer to the changelog (linked)
   for details on the web release.

   ## Summary * Bugfix
   [owncloud/web#9078](https://github.com/owncloud/web/pull/9078): Favorites list
   update on removal * Bugfix
   [owncloud/web#9213](https://github.com/owncloud/web/pull/9213): Space creation
   does not block reoccurring event * Bugfix
   [owncloud/web#9247](https://github.com/owncloud/web/issues/9247): Uploading to
   folders that contain special characters * Bugfix
   [owncloud/web#9259](https://github.com/owncloud/web/issues/9259): Relative user
   quota display limited to two decimals * Bugfix
   [owncloud/web#9261](https://github.com/owncloud/web/issues/9261): Remember
   location after token invalidation * Bugfix
   [owncloud/web#9299](https://github.com/owncloud/web/pull/9299): Authenticated
   public links breaking uploads * Bugfix
   [owncloud/web#9315](https://github.com/owncloud/web/issues/9315): Switch columns
   displayed on small screens in "Shared with me" view * Bugfix
   [owncloud/web#9351](https://github.com/owncloud/web/pull/9351): Media controls
   overflow on mobile screens * Bugfix
   [owncloud/web#9389](https://github.com/owncloud/web/pull/9389): Space editors
   see empty trashbin and delete actions in space trashbin * Bugfix
   [owncloud/web#9461](https://github.com/owncloud/web/pull/9461): Merging folders
   * Bugfix [owncloud/web/#9496](https://github.com/owncloud/web/pull/9496): Logo
   not showing * Bugfix
   [owncloud/web/#9489](https://github.com/owncloud/web/pull/9489): Public drop
   zone * Bugfix [owncloud/web/#9487](https://github.com/owncloud/web/pull/9487):
   Respect supportedClouds config * Bugfix
   [owncloud/web/#9507](https://github.com/owncloud/web/pull/9507): Space
   description edit modal is cut off vertically * Bugfix
   [owncloud/web/#9501](https://github.com/owncloud/web/pull/9501): Add cloud
   importer translations * Bugfix
   [owncloud/web/#9510](https://github.com/owncloud/web/pull/9510): Double items
   after moving a file with the same name * Enhancement
   [owncloud/web#7967](https://github.com/owncloud/web/pull/7967): Add hasPriority
   property for editors per extension * Enhancement
   [owncloud/web#8422](https://github.com/owncloud/web/issues/8422): Improve
   extension app topbar * Enhancement
   [owncloud/web#8445](https://github.com/owncloud/web/issues/8445): Open
   individually shared file in dedicated view * Enhancement
   [owncloud/web#8599](https://github.com/owncloud/web/issues/8599): Shrink table
   columns * Enhancement
   [owncloud/web#8921](https://github.com/owncloud/web/pull/8921): Add whitespace
   context-menu * Enhancement
   [owncloud/web#8983](https://github.com/owncloud/web/pull/8983): Deny share
   access * Enhancement
   [owncloud/web#8984](https://github.com/owncloud/web/pull/8984): Long breadcrumb
   strategy * Enhancement
   [owncloud/web#9044](https://github.com/owncloud/web/pull/9044): Search tag
   filter * Enhancement
   [owncloud/web#9046](https://github.com/owncloud/web/pull/9046): Single file link
   open with default app * Enhancement
   [owncloud/web#9052](https://github.com/owncloud/web/pull/9052): Drag & drop on
   parent folder * Enhancement
   [owncloud/web#9055](https://github.com/owncloud/web/pull/9055): Respect archiver
   limits * Enhancement
   [owncloud/web#9056](https://github.com/owncloud/web/issues/9056): Enable
   download (archive) on spaces * Enhancement
   [owncloud/web#9059](https://github.com/owncloud/web/pull/9059): Search full-text
   filter * Enhancement
   [owncloud/web#9077](https://github.com/owncloud/web/pull/9077): Advanced search
   button * Enhancement
   [owncloud/web#9077](https://github.com/owncloud/web/pull/9077): Search
   breadcrumb * Enhancement
   [owncloud/web#9088](https://github.com/owncloud/web/pull/9088): Use app icons
   for files * Enhancement
   [owncloud/web#9140](https://github.com/owncloud/web/pull/9140): Upload file on
   paste * Enhancement
   [owncloud/web#9151](https://github.com/owncloud/web/issues/9151): Cloud import *
   Enhancement [owncloud/web#9174](https://github.com/owncloud/web/issues/9174):
   Privacy statement in account menu * Enhancement
   [owncloud/web#9178](https://github.com/owncloud/web/pull/9178): Add login button
   to top bar * Enhancement
   [owncloud/web#9195](https://github.com/owncloud/web/pull/9195): Project spaces
   list viewmode * Enhancement
   [owncloud/web#9199](https://github.com/owncloud/web/pull/9199): Add pagination
   options to admin settings * Enhancement
   [owncloud/web#9200](https://github.com/owncloud/web/pull/9200): Add batch
   actions to search result list * Enhancement
   [owncloud/web#9216](https://github.com/owncloud/web/issues/9216): Restyle
   possible sharees * Enhancement
   [owncloud/web#9226](https://github.com/owncloud/web/pull/9226): Streamline URL
   query names * Enhancement
   [owncloud/web#9263](https://github.com/owncloud/web/pull/9263): Access denied
   page update message * Enhancement
   [owncloud/web#9280](https://github.com/owncloud/web/issues/9280): Hover tooltips
   in topbar * Enhancement
   [owncloud/web#9294](https://github.com/owncloud/web/pull/9294): Search list add
   highlighted file content * Enhancement
   [owncloud/web#9299](https://github.com/owncloud/web/pull/9299): Resolve pulic
   links to their actual location * Enhancement
   [owncloud/web#9304](https://github.com/owncloud/web/pull/9304): Add search
   location filter * Enhancement
   [owncloud/web#9344](https://github.com/owncloud/web/pull/9344): Ambiguation for
   URL view mode params * Enhancement
   [owncloud/web#9346](https://github.com/owncloud/web/pull/9346): Batch actions
   redesign * Enhancement
   [owncloud/web#9348](https://github.com/owncloud/web/pull/9348): Tag comma
   separation on client side * Enhancement
   [owncloud/web#9377](https://github.com/owncloud/web/issues/9377): User
   notification for blocked pop-ups and redirects * Enhancement
   [owncloud/web#9386](https://github.com/owncloud/web/pull/9386): Allow local
   storage for auth token * Enhancement
   [owncloud/web#9394](https://github.com/owncloud/web/pull/9394): Button styling *
   Enhancement [owncloud/web#9449](https://github.com/owncloud/web/issues/9449):
   Error notifications include x-request-id * Enhancement
   [owncloud/web#9426](https://github.com/owncloud/web/pull/9426): Add error log to
   upload dialog

   https://github.com/owncloud/ocis/pull/6944
   https://github.com/owncloud/web/releases/tag/v7.1.0-rc.5

* Enhancement - Add static secret to gn endpoints: [#6946](https://github.com/owncloud/ocis/pull/6946)

   The global notifications POST and DELETE endpoints (used only for deprovision
   notifications at the moment) can now be called by adding a static secret to the
   header. Admins can still call this endpoint without knowing the secret

   https://github.com/owncloud/ocis/pull/6946

* Enhancement - Bump sonarcloud: [#6961](https://github.com/owncloud/ocis/pull/6961)

   Bump sonarcloud to `5.0` to avoid java errors

   https://github.com/owncloud/ocis/pull/6961

* Enhancement - Nats named connections: [#6979](https://github.com/owncloud/ocis/pull/6979)

   Names the nats connections for easier debugging

   https://github.com/owncloud/ocis/pull/6979

* Enhancement - Add command for rebuilding the jsoncs3 share manager indexes: [#6986](https://github.com/owncloud/ocis/pull/6986)

   We added a command for rebuilding the jsoncs3 share manager indexes.

   https://github.com/owncloud/ocis/pull/6986
   https://github.com/owncloud/ocis/pull/6971

* Enhancement - Remove deprecated environment variables: [#7099](https://github.com/owncloud/ocis/pull/7099)

   We have removed all environment variables that have been marked as deprecated
   and marked for removal for 4.0.0

   https://github.com/owncloud/ocis/pull/7099

* Enhancement - Update web to v7.1.0: [#7107](https://github.com/owncloud/ocis/pull/7107)

   Tags: web

   We updated ownCloud Web to v7.1.0. Please refer to the changelog (linked) for
   details on the web release.

   ## Summary * Bugfix
   [owncloud/web#9078](https://github.com/owncloud/web/pull/9078): Favorites list
   update on removal * Bugfix
   [owncloud/web#9213](https://github.com/owncloud/web/pull/9213): Space creation
   does not block reoccurring event * Bugfix
   [owncloud/web#9247](https://github.com/owncloud/web/issues/9247): Uploading to
   folders that contain special characters * Bugfix
   [owncloud/web#9259](https://github.com/owncloud/web/issues/9259): Relative user
   quota display limited to two decimals * Bugfix
   [owncloud/web#9261](https://github.com/owncloud/web/issues/9261): Remember
   location after token invalidation * Bugfix
   [owncloud/web#9299](https://github.com/owncloud/web/pull/9299): Authenticated
   public links breaking uploads * Bugfix
   [owncloud/web#9315](https://github.com/owncloud/web/issues/9315): Switch columns
   displayed on small screens in "Shared with me" view * Bugfix
   [owncloud/web#9351](https://github.com/owncloud/web/pull/9351): Media controls
   overflow on mobile screens * Bugfix
   [owncloud/web#9389](https://github.com/owncloud/web/pull/9389): Space editors
   see empty trashbin and delete actions in space trashbin * Bugfix
   [owncloud/web#9461](https://github.com/owncloud/web/issues/9461): Merging
   folders * Enhancement
   [owncloud/web#7967](https://github.com/owncloud/web/pull/7967): Add hasPriority
   property for editors per extension * Enhancement
   [owncloud/web#8422](https://github.com/owncloud/web/issues/8422): Improve
   extension app topbar * Enhancement
   [owncloud/web#8445](https://github.com/owncloud/web/issues/8445): Open
   individually shared file in dedicated view * Enhancement
   [owncloud/web#8599](https://github.com/owncloud/web/issues/8599): Shrink table
   columns * Enhancement
   [owncloud/web#8921](https://github.com/owncloud/web/pull/8921): Add whitespace
   context-menu * Enhancement
   [owncloud/web#8983](https://github.com/owncloud/web/pull/8983): Deny share
   access * Enhancement
   [owncloud/web#8984](https://github.com/owncloud/web/pull/8984): Long breadcrumb
   strategy * Enhancement
   [owncloud/web#9044](https://github.com/owncloud/web/pull/9044): Search tag
   filter * Enhancement
   [owncloud/web#9046](https://github.com/owncloud/web/pull/9046): Single file link
   open with default app * Enhancement
   [owncloud/web#9052](https://github.com/owncloud/web/pull/9052): Drag & drop on
   parent folder * Enhancement
   [owncloud/web#9055](https://github.com/owncloud/web/pull/9055): Respect archiver
   limits * Enhancement
   [owncloud/web#9056](https://github.com/owncloud/web/issues/9056): Enable
   download (archive) on spaces * Enhancement
   [owncloud/web#9059](https://github.com/owncloud/web/pull/9059): Search full-text
   filter * Enhancement
   [owncloud/web#9077](https://github.com/owncloud/web/pull/9077): Advanced search
   button * Enhancement
   [owncloud/web#9077](https://github.com/owncloud/web/pull/9077): Search
   breadcrumb * Enhancement
   [owncloud/web#9088](https://github.com/owncloud/web/pull/9088): Use app icons
   for files * Enhancement
   [owncloud/web#9140](https://github.com/owncloud/web/pull/9140): Upload file on
   paste * Enhancement
   [owncloud/web#9151](https://github.com/owncloud/web/issues/9151): Cloud import *
   Enhancement [owncloud/web#9174](https://github.com/owncloud/web/issues/9174):
   Privacy statement in account menu * Enhancement
   [owncloud/web#9178](https://github.com/owncloud/web/pull/9178): Add login button
   to top bar * Enhancement
   [owncloud/web#9195](https://github.com/owncloud/web/pull/9195): Project spaces
   list viewmode * Enhancement
   [owncloud/web#9199](https://github.com/owncloud/web/pull/9199): Add pagination
   options to admin settings * Enhancement
   [owncloud/web#9200](https://github.com/owncloud/web/pull/9200): Add batch
   actions to search result list * Enhancement
   [owncloud/web#9216](https://github.com/owncloud/web/issues/9216): Restyle
   possible sharees * Enhancement
   [owncloud/web#9226](https://github.com/owncloud/web/pull/9226): Streamline URL
   query names * Enhancement
   [owncloud/web#9263](https://github.com/owncloud/web/pull/9263): Access denied
   page update message * Enhancement
   [owncloud/web#9280](https://github.com/owncloud/web/issues/9280): Hover tooltips
   in topbar * Enhancement
   [owncloud/web#9294](https://github.com/owncloud/web/pull/9294): Search list add
   highlighted file content * Enhancement
   [owncloud/web#9299](https://github.com/owncloud/web/pull/9299): Resolve pulic
   links to their actual location * Enhancement
   [owncloud/web#9304](https://github.com/owncloud/web/pull/9304): Add search
   location filter * Enhancement
   [owncloud/web#9344](https://github.com/owncloud/web/pull/9344): Ambiguation for
   URL view mode params * Enhancement
   [owncloud/web#9346](https://github.com/owncloud/web/pull/9346): Batch actions
   redesign * Enhancement
   [owncloud/web#9348](https://github.com/owncloud/web/pull/9348): Tag comma
   separation on client side * Enhancement
   [owncloud/web#9377](https://github.com/owncloud/web/issues/9377): User
   notification for blocked pop-ups and redirects * Enhancement
   [owncloud/web#9386](https://github.com/owncloud/web/pull/9386): Allow local
   storage for auth token * Enhancement
   [owncloud/web#9394](https://github.com/owncloud/web/pull/9394): Button styling *
   Enhancement [owncloud/web#9436](https://github.com/owncloud/web/pull/9436): Add
   error log to upload dialog

   https://github.com/owncloud/ocis/pull/7107
   https://github.com/owncloud/web/releases/tag/v7.1.0

# Changelog for [3.0.0] (2023-06-06)

The following sections list the changes for 3.0.0.

[3.0.0]: https://github.com/owncloud/ocis/compare/v2.0.0...v3.0.0

## Summary

* Bugfix - Use UUID attribute for computing "sub" claim in lico idp: [#904](https://github.com/owncloud/ocis/issues/904)
* Bugfix - Fix default role assignment for demo users: [#3432](https://github.com/owncloud/ocis/issues/3432)
* Bugfix - Hide the existence of space when deleting/updating: [#5031](https://github.com/owncloud/ocis/issues/5031)
* Bugfix - Fix Postprocessing events: [#5269](https://github.com/owncloud/ocis/pull/5269)
* Bugfix - Return 425 on Thumbnails: [#5300](https://github.com/owncloud/ocis/pull/5300)
* Bugfix - Disassociate users from deleted school: [#5343](https://github.com/owncloud/ocis/pull/5343)
* Bugfix - Fix Search tag indexing: [#5405](https://github.com/owncloud/ocis/pull/5405)
* Bugfix - Populate expanded properties: [#5421](https://github.com/owncloud/ocis/pull/5421)
* Bugfix - Fix the empty string givenName attribute when creating user: [#5431](https://github.com/owncloud/ocis/issues/5431)
* Bugfix - Add portrait thumbnail resolutions: [#5656](https://github.com/owncloud/ocis/pull/5656)
* Bugfix - Fix so that PATCH requests for groups actually updates the group name: [#5949](https://github.com/owncloud/ocis/pull/5949)
* Bugfix - Add missing CORS config: [#5987](https://github.com/owncloud/ocis/pull/5987)
* Bugfix - Fix authenticate headers for API requests: [#5992](https://github.com/owncloud/ocis/pull/5992)
* Bugfix - Fix OIDC auth cache: [#5997](https://github.com/owncloud/ocis/pull/5997)
* Bugfix - Fix user type config for user provider: [#6027](https://github.com/owncloud/ocis/pull/6027)
* Bugfix - Fix the wrong status code when appRoleAssignments is forbidden: [#6037](https://github.com/owncloud/ocis/issues/6037)
* Bugfix - Fix Search reindexing performance regression: [#6085](https://github.com/owncloud/ocis/pull/6085)
* Bugfix - Fix userlog panic: [#6114](https://github.com/owncloud/ocis/pull/6114)
* Bugfix - Fix wrong compile date: [#6132](https://github.com/owncloud/ocis/pull/6132)
* Bugfix - Fix Logout Url config name: [#6227](https://github.com/owncloud/ocis/pull/6227)
* Bugfix - Allow selected updates on graph users: [#6233](https://github.com/owncloud/ocis/pull/6233)
* Bugfix - Add missing response to blocked requests: [#6277](https://github.com/owncloud/ocis/pull/6277)
* Bugfix - Update the default admin role: [#6310](https://github.com/owncloud/ocis/pull/6310)
* Bugfix - Trace proxy middlewares: [#6313](https://github.com/owncloud/ocis/pull/6313)
* Bugfix - Reduced default TTL of user and group caches in graph API: [#6320](https://github.com/owncloud/ocis/issues/6320)
* Bugfix - Empty exact list while searching for a sharee: [#6398](https://github.com/owncloud/ocis/pull/6398)
* Bugfix - Fix error message when disabling users: [#6435](https://github.com/owncloud/ocis/pull/6435)
* Change - Remove the settings ui: [#5463](https://github.com/owncloud/ocis/pull/5463)
* Change - Do not share versions: [#5531](https://github.com/owncloud/ocis/pull/5531)
* Change - Bump libregraph lico: [#5768](https://github.com/owncloud/ocis/pull/5768)
* Change - Updated Cache Configuration: [#5829](https://github.com/owncloud/ocis/pull/5829)
* Change - We renamed the guest role to user light: [#6456](https://github.com/owncloud/ocis/pull/6456)
* Enhancement - Rename permissions: [#3922](https://github.com/cs3org/reva/pull/3922)
* Enhancement - Open Debug endpoint for Notifications: [#5002](https://github.com/owncloud/ocis/issues/5002)
* Enhancement - Open Debug endpoint for Nats: [#5002](https://github.com/owncloud/ocis/issues/5002)
* Enhancement - Add otlp tracing exporter: [#5132](https://github.com/owncloud/ocis/pull/5132)
* Enhancement - Add global env variable extractor: [#5164](https://github.com/owncloud/ocis/pull/5164)
* Enhancement - Async Postprocessing: [#5207](https://github.com/owncloud/ocis/pull/5207)
* Enhancement - Extended search: [#5221](https://github.com/owncloud/ocis/pull/5221)
* Enhancement - Resource tags: [#5227](https://github.com/owncloud/ocis/pull/5227)
* Enhancement - Bump libre-graph-api-go: [#5309](https://github.com/owncloud/ocis/pull/5309)
* Enhancement - Drive group permissions: [#5312](https://github.com/owncloud/ocis/pull/5312)
* Enhancement - Expiration Notifications: [#5330](https://github.com/owncloud/ocis/pull/5330)
* Enhancement - Graph Drives IdentitySet displayName: [#5347](https://github.com/owncloud/ocis/pull/5347)
* Enhancement - Make the group members addition limit configurable: [#5357](https://github.com/owncloud/ocis/pull/5357)
* Enhancement - Collect global envvars: [#5367](https://github.com/owncloud/ocis/pull/5367)
* Enhancement - Add webfinger service: [#5373](https://github.com/owncloud/ocis/pull/5373)
* Enhancement - Display surname and givenName attributes: [#5388](https://github.com/owncloud/ocis/pull/5388)
* Enhancement - Add expiration to user and group shares: [#5389](https://github.com/owncloud/ocis/pull/5389)
* Enhancement - Space Management permissions: [#5441](https://github.com/owncloud/ocis/pull/5441)
* Enhancement - Better config for postprocessing service: [#5457](https://github.com/owncloud/ocis/pull/5457)
* Enhancement - Cli to purge expired trash-bin items: [#5500](https://github.com/owncloud/ocis/pull/5500)
* Enhancement - Allow username to be changed: [#5509](https://github.com/owncloud/ocis/pull/5509)
* Enhancement - Allow users to be disabled: [#5588](https://github.com/owncloud/ocis/pull/5588)
* Enhancement - Make the settings bundles part of the service config: [#5589](https://github.com/owncloud/ocis/pull/5589)
* Enhancement - Add endpoint to list permissions: [#5594](https://github.com/owncloud/ocis/pull/5594)
* Enhancement - Eventhistory service: [#5600](https://github.com/owncloud/ocis/pull/5600)
* Enhancement - Userlog Service: [#5610](https://github.com/owncloud/ocis/pull/5610)
* Enhancement - Added option to configure default quota per role: [#5616](https://github.com/owncloud/ocis/pull/5616)
* Enhancement - Add new SetProjectSpaceQuota permission: [#5660](https://github.com/owncloud/ocis/pull/5660)
* Enhancement - Make graph/education API errors more consistent: [#5682](https://github.com/owncloud/ocis/pull/5682)
* Enhancement - Add new permission for public links: [#5690](https://github.com/owncloud/ocis/pull/5690)
* Enhancement - Userlog: [#5699](https://github.com/owncloud/ocis/pull/5699)
* Enhancement - Introduce policies-service: [#5714](https://github.com/owncloud/ocis/pull/5714)
* Enhancement - Update to go 1.20 to use memlimit: [#5732](https://github.com/owncloud/ocis/pull/5732)
* Enhancement - Add endpoints to upload a custom logo: [#5735](https://github.com/owncloud/ocis/pull/5735)
* Enhancement - Add config option to enforce passwords on public links: [#5848](https://github.com/owncloud/ocis/pull/5848)
* Enhancement - Add 'ocis decomposedfs metadata' command: [#5858](https://github.com/owncloud/ocis/pull/5858)
* Enhancement - Use gotext master: [#5867](https://github.com/owncloud/ocis/pull/5867)
* Enhancement - No Notifications for own actions: [#5871](https://github.com/owncloud/ocis/pull/5871)
* Enhancement - Automate md creation: [#5901](https://github.com/owncloud/ocis/pull/5901)
* Enhancement - Notify about policies: [#5912](https://github.com/owncloud/ocis/pull/5912)
* Enhancement - Use Accept-Language Header: [#5918](https://github.com/owncloud/ocis/pull/5918)
* Enhancement - Add MessageRichParameters: [#5927](https://github.com/owncloud/ocis/pull/5927)
* Enhancement - Add more logging to av service: [#5973](https://github.com/owncloud/ocis/pull/5973)
* Enhancement - Make the LDAP base DN for new groups configurable: [#5974](https://github.com/owncloud/ocis/pull/5974)
* Enhancement - Add a capability for the Personal Data export: [#5984](https://github.com/owncloud/ocis/pull/5984)
* Enhancement - Bump go-ldap version: [#6004](https://github.com/owncloud/ocis/pull/6004)
* Enhancement - Configure GRPC in ocs: [#6022](https://github.com/owncloud/ocis/pull/6022)
* Enhancement - Web config additions: [#6032](https://github.com/owncloud/ocis/pull/6032)
* Enhancement - Notifications: [#6038](https://github.com/owncloud/ocis/pull/6038)
* Enhancement - Added possibility to assign roles based on OIDC claims: [#6048](https://github.com/owncloud/ocis/pull/6048)
* Enhancement - GDPR Export: [#6064](https://github.com/owncloud/ocis/pull/6064)
* Enhancement - Add optional services to the runtime: [#6071](https://github.com/owncloud/ocis/pull/6071)
* Enhancement - Determine the users language to translate via Transifex: [#6089](https://github.com/owncloud/ocis/pull/6089)
* Enhancement - Return Bad Request when requesting GDPR export for another user: [#6123](https://github.com/owncloud/ocis/pull/6123)
* Enhancement - Disable Notifications: [#6137](https://github.com/owncloud/ocis/pull/6137)
* Enhancement - Add the email HTML templates: [#6147](https://github.com/owncloud/ocis/pull/6147)
* Enhancement - Add debug server to idm: [#6153](https://github.com/owncloud/ocis/pull/6153)
* Enhancement - Add debug server to audit: [#6178](https://github.com/owncloud/ocis/pull/6178)
* Enhancement - Web options configuration: [#6188](https://github.com/owncloud/ocis/pull/6188)
* Enhancement - Add debug server to userlog: [#6202](https://github.com/owncloud/ocis/pull/6202)
* Enhancement - Add debug server to postprocessing: [#6203](https://github.com/owncloud/ocis/pull/6203)
* Enhancement - Add debug server to eventhistory: [#6204](https://github.com/owncloud/ocis/pull/6204)
* Enhancement - Add specific result to antivirus for debugging: [#6265](https://github.com/owncloud/ocis/pull/6265)
* Enhancement - Add Store to `postprocessing`: [#6281](https://github.com/owncloud/ocis/pull/6281)
* Enhancement - Update web to v7.0.0-rc.37: [#6294](https://github.com/owncloud/ocis/pull/6294)
* Enhancement - Remove quota from share jails api responses: [#6309](https://github.com/owncloud/ocis/pull/6309)
* Enhancement - Graph user capabilities: [#6339](https://github.com/owncloud/ocis/pull/6339)
* Enhancement - Configurable ID Cache: [#6353](https://github.com/owncloud/ocis/pull/6353)
* Enhancement - Fix err when the user share the locked file: [#6358](https://github.com/owncloud/ocis/pull/6358)
* Enhancement - Remove the email logo: [#6359](https://github.com/owncloud/ocis/issues/6359)
* Enhancement - Default LDAP write to true: [#6362](https://github.com/owncloud/ocis/pull/6362)
* Enhancement - Add fulltextsearch capabilty: [#6366](https://github.com/owncloud/ocis/pull/6366)
* Enhancement - Update web to v7.0.0-rc.38: [#6375](https://github.com/owncloud/ocis/pull/6375)
* Enhancement - Fix preview or viewing of shared animated GIFs: [#6386](https://github.com/owncloud/ocis/pull/6386)
* Enhancement - Unify CA Cert envvars: [#6392](https://github.com/owncloud/ocis/pull/6392)
* Enhancement - Fix to prevent the email X-Site scripting: [#6429](https://github.com/owncloud/ocis/pull/6429)
* Enhancement - Update web to v7.0.0: [#6438](https://github.com/owncloud/ocis/pull/6438)
* Enhancement - Update Reva to version 2.14.0: [#6448](https://github.com/owncloud/ocis/pull/6448)

## Details

* Bugfix - Use UUID attribute for computing "sub" claim in lico idp: [#904](https://github.com/owncloud/ocis/issues/904)

   By default the LDAP backend for lico uses the User DN for computing the "sub"
   claim of a user. This caused the "sub" claim to stay the same even if a user was
   deleted and recreated (and go a new UUID assgined with that). We now use the
   user's unique id (`owncloudUUID` by default) for computing the `sub` claim. So
   that user's recreated with the same name will be treated as different users by
   the IDP.

   https://github.com/owncloud/ocis/issues/904
   https://github.com/owncloud/ocis/pull/6326
   https://github.com/owncloud/ocis/pull/6338
   https://github.com/owncloud/ocis/pull/6420

* Bugfix - Fix default role assignment for demo users: [#3432](https://github.com/owncloud/ocis/issues/3432)

   The roles-assignments for demo users where duplicated with every restart of the
   settings service.

   https://github.com/owncloud/ocis/issues/3432

* Bugfix - Hide the existence of space when deleting/updating: [#5031](https://github.com/owncloud/ocis/issues/5031)

   The "code": "notAllowed" changed to "code": "itemNotFound"

   https://github.com/owncloud/ocis/issues/5031
   https://github.com/owncloud/ocis/pull/6220

* Bugfix - Fix Postprocessing events: [#5269](https://github.com/owncloud/ocis/pull/5269)

   Postprocessing service did not want to play with non-tls events. That is fixed
   now

   https://github.com/owncloud/ocis/pull/5269

* Bugfix - Return 425 on Thumbnails: [#5300](https://github.com/owncloud/ocis/pull/5300)

   Return `425` on thumbnails `GET` when file is processing. Pass `425` also
   through webdav endpoint

   https://github.com/owncloud/ocis/pull/5300

* Bugfix - Disassociate users from deleted school: [#5343](https://github.com/owncloud/ocis/pull/5343)

   When a school is deleted, users should be disassociated from it.

   https://github.com/owncloud/ocis/issues/5246
   https://github.com/owncloud/ocis/pull/5343

* Bugfix - Fix Search tag indexing: [#5405](https://github.com/owncloud/ocis/pull/5405)

   We've fixed an issue where search is not able to index tags for space resources.

   https://github.com/owncloud/ocis/pull/5405

* Bugfix - Populate expanded properties: [#5421](https://github.com/owncloud/ocis/pull/5421)

   We now return an empty array when an expanded relation has no entries. This
   makes consuming the responses a little easier.

   https://github.com/owncloud/ocis/issues/5419
   https://github.com/owncloud/ocis/pull/5421
   https://github.com/owncloud/ocis/pull/5426

* Bugfix - Fix the empty string givenName attribute when creating user: [#5431](https://github.com/owncloud/ocis/issues/5431)

   Omitempty givenName attribute when creating user

   https://github.com/owncloud/ocis/issues/5431
   https://github.com/owncloud/ocis/pull/6259

* Bugfix - Add portrait thumbnail resolutions: [#5656](https://github.com/owncloud/ocis/pull/5656)

   Add portrait-orientation resolutions to the thumbnail service's default
   configuration. This prevents portrait photos from being heavily cropped into
   landscape resolutions in the web viewer.

   https://github.com/owncloud/ocis/pull/5656

* Bugfix - Fix so that PATCH requests for groups actually updates the group name: [#5949](https://github.com/owncloud/ocis/pull/5949)

   https://github.com/owncloud/ocis/pull/5949

* Bugfix - Add missing CORS config: [#5987](https://github.com/owncloud/ocis/pull/5987)

   The graph, userlog and ocdav services had no CORS config options.

   https://github.com/owncloud/ocis/pull/5987

* Bugfix - Fix authenticate headers for API requests: [#5992](https://github.com/owncloud/ocis/pull/5992)

   We changed the www-authenticate header which should not be sent when the
   `XMLHttpRequest` header is set.

   https://github.com/owncloud/ocis/issues/5986
   https://github.com/owncloud/ocis/pull/5992

* Bugfix - Fix OIDC auth cache: [#5997](https://github.com/owncloud/ocis/pull/5997)

   We've fixed an issue rendering the OIDC auth cache useless.

   https://github.com/owncloud/ocis/pull/5997

* Bugfix - Fix user type config for user provider: [#6027](https://github.com/owncloud/ocis/pull/6027)

   We needed to provide a default value for the user type property in the user
   provider.

   https://github.com/owncloud/ocis/pull/6027

* Bugfix - Fix the wrong status code when appRoleAssignments is forbidden: [#6037](https://github.com/owncloud/ocis/issues/6037)

   Fix the wrong status code when appRoleAssignments is forbidden in the
   CreateAppRoleAssignment and DeleteAppRoleAssignment methods.

   https://github.com/owncloud/ocis/issues/6037
   https://github.com/owncloud/ocis/pull/6276

* Bugfix - Fix Search reindexing performance regression: [#6085](https://github.com/owncloud/ocis/pull/6085)

   We've fixed a regression in the search service reindexing step, causing the
   whole space to be reindexed instead of just the changed resources.

   https://github.com/owncloud/ocis/pull/6085

* Bugfix - Fix userlog panic: [#6114](https://github.com/owncloud/ocis/pull/6114)

   Userlog services paniced because of `nil` ctx. That is fixed now

   https://github.com/owncloud/ocis/pull/6114

* Bugfix - Fix wrong compile date: [#6132](https://github.com/owncloud/ocis/pull/6132)

   We fixed that current date is always printed.

   https://github.com/owncloud/ocis/issues/6124
   https://github.com/owncloud/ocis/pull/6132

* Bugfix - Fix Logout Url config name: [#6227](https://github.com/owncloud/ocis/pull/6227)

   We fixed the yaml and json name of the logout url option.

   https://github.com/owncloud/ocis/pull/6227

* Bugfix - Allow selected updates on graph users: [#6233](https://github.com/owncloud/ocis/pull/6233)

   We are now allowing a couple of update request to complete even if
   GRAPH_LDAP_SERVER_WRITE_ENABLED=false:

  *   When using a group to disable users (OCIS_LDAP_DISABLE_USER_MECHANISM=group) updates to the accountEnabled property of a user will be allowed
  *   When a distinct base dn for new groups is configured ( GRAPH_LDAP_GROUP_CREATE_BASE_DN is set to a different value than GRAPH_LDAP_GROUP_BASE_DN), allow the creation/update of local groups.

   https://github.com/owncloud/ocis/pull/6233

* Bugfix - Add missing response to blocked requests: [#6277](https://github.com/owncloud/ocis/pull/6277)

   We added the missing response body to requests which were blocked by the policy
   engine.

   https://github.com/owncloud/ocis/pull/6277

* Bugfix - Update the default admin role: [#6310](https://github.com/owncloud/ocis/pull/6310)

   The admin role was missing two permissions. We added them to make the space
   admin role a subset of the admin role. This matches better with the default user
   expectations.

   https://github.com/owncloud/ocis/pull/6310

* Bugfix - Trace proxy middlewares: [#6313](https://github.com/owncloud/ocis/pull/6313)

   We moved trace initialization to an early middleware to also trace requests made
   by other proxy middlewares.

   https://github.com/owncloud/ocis/pull/6313

* Bugfix - Reduced default TTL of user and group caches in graph API: [#6320](https://github.com/owncloud/ocis/issues/6320)

   We reduced the default TTL of the cache for user and group information on the
   /drives endpoints to 60 seconds. This fixes in issue where outdated information
   was show on the spaces list for a very long time.

   https://github.com/owncloud/ocis/issues/6320

* Bugfix - Empty exact list while searching for a sharee: [#6398](https://github.com/owncloud/ocis/pull/6398)

   We fixed a bug in the sharing api, it always returns an empty exact list while
   searching for a sharee

   https://github.com/owncloud/ocis/issues/4265
   https://github.com/owncloud/ocis/pull/6398
   https://github.com/cs3org/reva/pull/3877

* Bugfix - Fix error message when disabling users: [#6435](https://github.com/owncloud/ocis/pull/6435)

   When we disable users by adding them to a group we do not need to update the
   user entry.

   https://github.com/owncloud/ocis/pull/6435

* Change - Remove the settings ui: [#5463](https://github.com/owncloud/ocis/pull/5463)

   With ownCloud Web having transitioned to Vue 3 recently, we would have had to
   port the settings ui as well. The decision was made to discontinue the settings
   ui instead. As a result all traces of the settings ui have been removed.

   The only user facing setting that ever existed in the settings service is now
   integrated into the `account` page of ownCloud Web (click on top right user
   menu, then on your username to reach the account page).

   https://github.com/owncloud/ocis/pull/5463

* Change - Do not share versions: [#5531](https://github.com/owncloud/ocis/pull/5531)

   We changed the default behavior of shares: Share receivers have no access to
   versions. People in spaces with the "Editor" or "Manager" role can still see
   versions and work with them.

   https://github.com/owncloud/ocis/pull/5531

* Change - Bump libregraph lico: [#5768](https://github.com/owncloud/ocis/pull/5768)

   We updated lico to the latest version * Update to 0.59.4 - upstream dropped the
   kc and cookie backends

   https://github.com/owncloud/ocis/pull/5768

* Change - Updated Cache Configuration: [#5829](https://github.com/owncloud/ocis/pull/5829)

   We updated all cache related environment vars to more closely follow the go
   micro naming pattern: - `{service}_CACHE_STORE_TYPE` becomes
   `{service}_CACHE_STORE` or `{service}_PERSISTENT_STORE` -
   `{service}_CACHE_STORE_ADDRESS(ES)` becomes `{service}_CACHE_STORE_NODES` - The
   `mem` store implementation name changes to `memory` - In yaml files the cache
   `type` becomes `store` We introduced `redis-sentinel` as a store implementation.

   https://github.com/owncloud/ocis/pull/5829

* Change - We renamed the guest role to user light: [#6456](https://github.com/owncloud/ocis/pull/6456)

   We needed to rename the "Guest" role to "User Light" because the naming was
   creating confusions. The roles are not bound to a user type.

   https://github.com/owncloud/ocis/issues/6058
   https://github.com/owncloud/ocis/pull/6456

* Enhancement - Rename permissions: [#3922](https://github.com/cs3org/reva/pull/3922)

   Rename permissions to be consistent and future proof

   https://github.com/cs3org/reva/pull/3922
   https://github.com/owncloud/ocis/pull/6418

* Enhancement - Open Debug endpoint for Notifications: [#5002](https://github.com/owncloud/ocis/issues/5002)

   We added a debug server to the notifications service

   https://github.com/owncloud/ocis/issues/5002
   https://github.com/owncloud/ocis/pull/6155

* Enhancement - Open Debug endpoint for Nats: [#5002](https://github.com/owncloud/ocis/issues/5002)

   We added a debug server to nats

   https://github.com/owncloud/ocis/issues/5002
   https://github.com/owncloud/ocis/pull/6139

* Enhancement - Add otlp tracing exporter: [#5132](https://github.com/owncloud/ocis/pull/5132)

   We can now configure otlp to send traces using the otlp exporter.

   https://github.com/owncloud/ocis/pull/5132
   https://github.com/cs3org/reva/pull/3496

* Enhancement - Add global env variable extractor: [#5164](https://github.com/owncloud/ocis/pull/5164)

   We have added a little tool that will extract global env vars, that are loaded
   only through os.Getenv for documentation purposes

   https://github.com/owncloud/ocis/issues/4916
   https://github.com/owncloud/ocis/pull/5164

* Enhancement - Async Postprocessing: [#5207](https://github.com/owncloud/ocis/pull/5207)

   Provides functionality for async postprocessing. This will allow the system to
   do the postprocessing (virusscan, copying of bytes to their final destination,
   ...) asynchronous to the users request. Major change when active.

   https://github.com/owncloud/ocis/pull/5207

* Enhancement - Extended search: [#5221](https://github.com/owncloud/ocis/pull/5221)

   Provides multiple enhancement to the search implementation. * content
   extraction, search now supports apache tika to extract resource contents. *
   search engine, underlying search engine is swappable now. * event consumers, the
   number of event consumers can now be set, which improves the speed of the
   individual tasks

   https://github.com/owncloud/ocis/issues/5184
   https://github.com/owncloud/ocis/pull/5221

* Enhancement - Resource tags: [#5227](https://github.com/owncloud/ocis/pull/5227)

   We've added the ability to tag resources via the graph api. Tags can be added
   (put request) and removed (delete request) from a resource, a list of available
   tags can also be requested by sending a get request to the graph endpoint.

   https://github.com/owncloud/ocis/issues/5184
   https://github.com/owncloud/ocis/pull/5227
   https://github.com/owncloud/ocis/pull/5271

* Enhancement - Bump libre-graph-api-go: [#5309](https://github.com/owncloud/ocis/pull/5309)

   We fixed a couple of issues in libre-graph-api-go package.

  * rename drive permission grantedTo to grantedToIdentities to be ms graph spec compatible.
  * drive.name is a required property now.
  * add group property to the identitySet.

   https://github.com/owncloud/ocis/pull/5309
   https://github.com/owncloud/ocis/pull/5312

* Enhancement - Drive group permissions: [#5312](https://github.com/owncloud/ocis/pull/5312)

   We've updated the libregraph.Drive response to contain group permissions.

   https://github.com/owncloud/ocis/pull/5312

* Enhancement - Expiration Notifications: [#5330](https://github.com/owncloud/ocis/pull/5330)

   Send emails to the user informing that a share or a space membership expires.

   https://github.com/owncloud/ocis/pull/5330

* Enhancement - Graph Drives IdentitySet displayName: [#5347](https://github.com/owncloud/ocis/pull/5347)

   We've added the IdentitySet displayName property to the group and user sets for
   the graph drives endpoint. The values for groups and users get cached.

   https://github.com/owncloud/ocis/pull/5347
   https://github.com/owncloud/web/pull/8178

* Enhancement - Make the group members addition limit configurable: [#5357](https://github.com/owncloud/ocis/pull/5357)

   It's now possible to configure the limit of group members addition by PATCHing
   `/graph/v1.0/groups/{groupID}`. It still defaults to 20 as defined in the spec
   but it can be configured via `.graph.api.group_members_patch_limit` in
   `ocis.yaml` or via the `GRAPH_GROUP_MEMBERS_PATCH_LIMIT` environment variable.

   https://github.com/owncloud/ocis/issues/5262
   https://github.com/owncloud/ocis/pull/5357

* Enhancement - Collect global envvars: [#5367](https://github.com/owncloud/ocis/pull/5367)

   Compose a list of all envvars living in more than 1 service

   https://github.com/owncloud/ocis/pull/5367

* Enhancement - Add webfinger service: [#5373](https://github.com/owncloud/ocis/pull/5373)

   Adds a webfinger service to redirect ocis clients

   https://github.com/owncloud/ocis/issues/6102
   https://github.com/owncloud/ocis/pull/5373
   https://github.com/owncloud/ocis/pull/6110

* Enhancement - Display surname and givenName attributes: [#5388](https://github.com/owncloud/ocis/pull/5388)

   When querying the graph API, the surname and givenName attributes are now
   displayed for users.

   https://github.com/owncloud/ocis/issues/5386
   https://github.com/owncloud/ocis/pull/5388

* Enhancement - Add expiration to user and group shares: [#5389](https://github.com/owncloud/ocis/pull/5389)

   Added expiration to user and group shares.

   https://github.com/owncloud/ocis/pull/5389

* Enhancement - Space Management permissions: [#5441](https://github.com/owncloud/ocis/pull/5441)

   We added new space management permissions. `space-properties` will allow
   changing space properties (name, description, ...). `space-ability` will allow
   enabling and disabling spaces

   https://github.com/owncloud/ocis/pull/5441

* Enhancement - Better config for postprocessing service: [#5457](https://github.com/owncloud/ocis/pull/5457)

   The postprocessing service is now individually configurable. This is achieved by
   allowing a list of postprocessing steps that are processed in order of their
   appearance in the `POSTPROCESSING_STEPS` envvar.

   https://github.com/owncloud/ocis/pull/5457

* Enhancement - Cli to purge expired trash-bin items: [#5500](https://github.com/owncloud/ocis/pull/5500)

   Introduction of a new cli command to purge old trash-bin items. The command is
   part of the `storage-users` service and can be used as follows:

   `ocis storage-users trash-bin purge-expired`.

   The `purge-expired` command configuration is done in the `ocis`configuration or
   as usual by using environment variables.

   ENV `STORAGE_USERS_PURGE_TRASH_BIN_USER_ID` is used to obtain space trash-bin
   information and takes the system admin user as the default `OCIS_ADMIN_USER_ID`.
   It should be noted, that this is only set by default in the single binary. The
   command only considers spaces to which the user has access and delete
   permission.

   ENV `STORAGE_USERS_PURGE_TRASH_BIN_PERSONAL_DELETE_BEFORE` has a default value
   of `30 days`, which means the command will delete all files older than `30
   days`. The value is human-readable, valid values are `24h`, `60m`, `60s` etc.
   `0` is equivalent to disable and prevents the deletion of `personal space`
   trash-bin files.

   ENV `STORAGE_USERS_PURGE_TRASH_BIN_PROJECT_DELETE_BEFORE` has a default value of
   `30 days`, which means the command will delete all files older than `30 days`.
   The value is human-readable, valid values are `24h`, `60m`, `60s` etc. `0` is
   equivalent to disable and prevents the deletion of `project space` trash-bin
   files.

   Likewise, only spaces of the type `project` and `personal` are taken into
   account. Spaces of type `virtual`, for example, are ignored.

   https://github.com/owncloud/ocis/issues/5499
   https://github.com/owncloud/ocis/pull/5500

* Enhancement - Allow username to be changed: [#5509](https://github.com/owncloud/ocis/pull/5509)

   When OnPremisesSamAccountName is present in a PATCH on
   `{apiRoot}/users/{userID}` it will change the username of the user. This also
   changes the references to this user in the groups.

   https://github.com/owncloud/ocis/issues/4988
   https://github.com/owncloud/ocis/pull/5509

* Enhancement - Allow users to be disabled: [#5588](https://github.com/owncloud/ocis/pull/5588)

   By setting the `accountEnabled` property to `false` for a user via the graph
   API. Users can be disabled (i.e. they can no longer login)

   https://github.com/owncloud/ocis/pull/5588
   https://github.com/owncloud/ocis/pull/5620

* Enhancement - Make the settings bundles part of the service config: [#5589](https://github.com/owncloud/ocis/pull/5589)

   We added the settings bundles to the config. The default roles are still
   unchanged. You can now override the defaults by replacing the whole bundles list
   via json config files. The config file is loaded from a specified path which can
   be configured with `SETTINGS_BUNDLES_PATH`.

   https://github.com/owncloud/ocis/pull/5589
   https://github.com/owncloud/ocis/pull/5607

* Enhancement - Add endpoint to list permissions: [#5594](https://github.com/owncloud/ocis/pull/5594)

   We added 'https://cloud.ocis.test/api/v0/settings/permissions-list' to retrieve
   all permissions of the logged in user.

   https://github.com/owncloud/ocis/pull/5594
   https://github.com/owncloud/ocis/pull/5571

* Enhancement - Eventhistory service: [#5600](https://github.com/owncloud/ocis/pull/5600)

   Introduces the `eventhistory` service. It is a service that stores events and
   provides a grpc API to retrieve them.

   https://github.com/owncloud/ocis/pull/5600

* Enhancement - Userlog Service: [#5610](https://github.com/owncloud/ocis/pull/5610)

   Introduces userlog service. It stores eventIDs the user is interested in and
   provides an API to retrieve the events.

   https://github.com/owncloud/ocis/pull/5610

* Enhancement - Added option to configure default quota per role: [#5616](https://github.com/owncloud/ocis/pull/5616)

   Admins can assign default quotas to users with certain roles by adding the
   following config to the `proxy.yaml`. E.g.:

   ```
   role_quotas:
       d7beeea8-8ff4-406b-8fb6-ab2dd81e6b11: 2300000
   ```

   It maps a role ID to the quota in bytes.

   https://github.com/owncloud/ocis/pull/5616

* Enhancement - Add new SetProjectSpaceQuota permission: [#5660](https://github.com/owncloud/ocis/pull/5660)

   Additionally to `set-space-quota` for setting quota on personal spaces we now
   have `Drive.ReadWriteQuota.Project` for setting project spaces quota

   https://github.com/owncloud/ocis/pull/5660

* Enhancement - Make graph/education API errors more consistent: [#5682](https://github.com/owncloud/ocis/pull/5682)

   Aligned the error messages when creating schools and classes fail and changed
   the response code from 500 to 409.

   https://github.com/owncloud/ocis/issues/5660
   https://github.com/owncloud/ocis/pull/5682

* Enhancement - Add new permission for public links: [#5690](https://github.com/owncloud/ocis/pull/5690)

   Added a new permission 'PublicLink.Write' to check if a user can create or
   update public links.

   https://github.com/owncloud/ocis/pull/5690

* Enhancement - Userlog: [#5699](https://github.com/owncloud/ocis/pull/5699)

   Enhance userlog service with proper api and messages

   https://github.com/owncloud/ocis/pull/5699

* Enhancement - Introduce policies-service: [#5714](https://github.com/owncloud/ocis/pull/5714)

   Introduces policies service. The policies-service provides a new grpc api which
   can be used to return whether a requested operation is allowed or not. Open
   Policy Agent is used to determine the set of rules of what is permitted and what
   is not.

   2 further levels of authorization build on this:

  * Proxy Authorization
  * Event Authorization (needs async post-processing enabled)

   The simplest authorization layer is in the proxy, since every request is
   processed here, only simple decisions that can be processed quickly are made
   here, more complex queries such as file evaluation are explicitly excluded in
   this layer.

   The next layer is event-based as a pipeline step in asynchronous
   post-processing, since processing at this point is asynchronous, the operations
   there can also take longer and be more expensive, the bytes of a file can be
   examined here as an example.

   Since the base block is a grpc api, it is also possible to use it directly. The
   policies are written in the [rego query
   language](https://www.openpolicyagent.org/docs/latest/policy-language/).

   https://github.com/owncloud/ocis/issues/5580
   https://github.com/owncloud/ocis/pull/5714

* Enhancement - Update to go 1.20 to use memlimit: [#5732](https://github.com/owncloud/ocis/pull/5732)

   We updated to go 1.20 which allows setting GOMEMLIMIT, which we by default set
   to 0.9.

   https://github.com/owncloud/ocis/pull/5732

* Enhancement - Add endpoints to upload a custom logo: [#5735](https://github.com/owncloud/ocis/pull/5735)

   Added endpoints to upload and reset custom logos. The files are stored under the
   `WEB_ASSET_PATH` which defaults to `$OCIS_BASE_DATA_PATH/web/assets`.

   https://github.com/owncloud/ocis/pull/5735
   https://github.com/owncloud/ocis/pull/5559

* Enhancement - Add config option to enforce passwords on public links: [#5848](https://github.com/owncloud/ocis/pull/5848)

   Added a new config option to enforce passwords on public links with "Uploader,
   Editor, Contributor" roles.

   The new options are: `OCIS_SHARING_PUBLIC_WRITEABLE_SHARE_MUST_HAVE_PASSWORD`,
   `SHARING_PUBLIC_WRITEABLE_SHARE_MUST_HAVE_PASSWORD` and
   `FRONTEND_OCS_PUBLIC_WRITEABLE_SHARE_MUST_HAVE_PASSWORD`. Check the docs on how
   to properly set them.

   https://github.com/owncloud/ocis/pull/5848
   https://github.com/owncloud/ocis/pull/5785
   https://github.com/owncloud/ocis/pull/5720

* Enhancement - Add 'ocis decomposedfs metadata' command: [#5858](https://github.com/owncloud/ocis/pull/5858)

   We added a 'ocis decomposedfs metadata' command for inspecting and manipulating
   node metadata.

   https://github.com/owncloud/ocis/pull/5858

* Enhancement - Use gotext master: [#5867](https://github.com/owncloud/ocis/pull/5867)

   We needed to use forked version until our upstream changes were merged

   https://github.com/owncloud/ocis/pull/5867

* Enhancement - No Notifications for own actions: [#5871](https://github.com/owncloud/ocis/pull/5871)

   Don't send notifications on space events when the user has executed them
   herself.

   https://github.com/owncloud/ocis/pull/5871

* Enhancement - Automate md creation: [#5901](https://github.com/owncloud/ocis/pull/5901)

   Automatically create `_index.md` files from the services `README.md`

   https://github.com/owncloud/ocis/pull/5901

* Enhancement - Notify about policies: [#5912](https://github.com/owncloud/ocis/pull/5912)

   Notify the user when a file was deleted due to policies (policies service)

   https://github.com/owncloud/ocis/pull/5912

* Enhancement - Use Accept-Language Header: [#5918](https://github.com/owncloud/ocis/pull/5918)

   Use the `Accept-Language` header instead of the custom `Prefered-Language`

   https://github.com/owncloud/ocis/pull/5918

* Enhancement - Add MessageRichParameters: [#5927](https://github.com/owncloud/ocis/pull/5927)

   Adds the messageRichParameters to virus and policies notifications

   https://github.com/owncloud/ocis/pull/5927

* Enhancement - Add more logging to av service: [#5973](https://github.com/owncloud/ocis/pull/5973)

   We need more debug logging in some situations to understand the state of a virus
   scan.

   https://github.com/owncloud/ocis/pull/5973

* Enhancement - Make the LDAP base DN for new groups configurable: [#5974](https://github.com/owncloud/ocis/pull/5974)

   The LDAP backend for the Graph service introduced a new config option for
   setting the Parent DN for new groups created via the `/groups/` endpoint.
   (`GRAPH_LDAP_GROUP_CREATE_BASE_DN`)

   It defaults to the value of `GRAPH_LDAP_GROUP_BASE_DN`. If set to a different
   value the `GRAPH_LDAP_GROUP_CREATE_BASE_DN` needs to be a subordinate DN of
   `GRAPH_LDAP_GROUP_BASE_DN`.

   All existing groups with a DN outside the `GRAPH_LDAP_GROUP_CREATE_BASE_DN` tree
   will be treated as read-only groups. So it is not possible to edit these groups.

   https://github.com/owncloud/ocis/pull/5974

* Enhancement - Add a capability for the Personal Data export: [#5984](https://github.com/owncloud/ocis/pull/5984)

   Adds a capability for the personal data export endpoint

   https://github.com/owncloud/ocis/pull/5984

* Enhancement - Bump go-ldap version: [#6004](https://github.com/owncloud/ocis/pull/6004)

   Use master version of go-ldap to get rid of nasty `=` bug. See
   https://github.com/go-ldap/ldap/issues/416

   https://github.com/owncloud/ocis/pull/6004

* Enhancement - Configure GRPC in ocs: [#6022](https://github.com/owncloud/ocis/pull/6022)

   Fixes a panic in ocs when running not in single binary

   https://github.com/owncloud/ocis/pull/6022

* Enhancement - Web config additions: [#6032](https://github.com/owncloud/ocis/pull/6032)

   We've added config keys for defining additional css, scripts and translations
   for ownCloud Web.

   https://github.com/owncloud/ocis/pull/6032

* Enhancement - Notifications: [#6038](https://github.com/owncloud/ocis/pull/6038)

   Make Emails translatable via transifex The transifex translation add in to the
   email templates. The optional environment variable
   NOTIFICATIONS_TRANSLATION_PATH added to config. The optional global environment
   variable OCIS_TRANSLATION_PATH added to notifications and userlog config.

   https://github.com/owncloud/ocis/issues/6025
   https://github.com/owncloud/ocis/pull/6038

* Enhancement - Added possibility to assign roles based on OIDC claims: [#6048](https://github.com/owncloud/ocis/pull/6048)

   OCIS can now be configured to update a user's role assignment from the values of
   a claim provided via the IDPs userinfo endpoint. The claim name and the mapping
   between claim values and ocis role name can be configured via the configuration
   of the proxy service. Example:

   ```
   role_assignment:
       driver: oidc
       oidc_role_mapper:
           role_claim: ocisRoles
           role_mapping:
               - role_name: admin
                 claim_value: myAdminRole
               - role_name: spaceadmin
                 claim_value: mySpaceAdminRole
               - role_name: user
                 claim_value: myUserRole
               - role_name: guest
                 claim_value: myGuestRole
   ```

   https://github.com/owncloud/ocis/pull/6048

* Enhancement - GDPR Export: [#6064](https://github.com/owncloud/ocis/pull/6064)

   Adds an endpoint to collect all data that is related to a user

   https://github.com/owncloud/ocis/pull/6064
   https://github.com/owncloud/ocis/pull/5950

* Enhancement - Add optional services to the runtime: [#6071](https://github.com/owncloud/ocis/pull/6071)

   Make it possible to start optional services in the ocis runtime. Instead of
   using `OCIS_RUN_SERVICES` to define all services we can now use
   `OCIS_ADD_RUN_SERVICES` to add a comma separated list of additional services
   which are not started in the single process by default.

   https://github.com/owncloud/ocis/pull/6071

* Enhancement - Determine the users language to translate via Transifex: [#6089](https://github.com/owncloud/ocis/pull/6089)

   https://github.com/owncloud/ocis/issues/6087
   https://github.com/owncloud/ocis/pull/6089
   Enhance
   userlog
   service
   with
   proper
   api
   and
   messages

* Enhancement - Return Bad Request when requesting GDPR export for another user: [#6123](https://github.com/owncloud/ocis/pull/6123)

   This is an enhancement, not security related as the requested uid is never used

   https://github.com/owncloud/ocis/pull/6123

* Enhancement - Disable Notifications: [#6137](https://github.com/owncloud/ocis/pull/6137)

   Introduce new setting to disable notifications

   https://github.com/owncloud/ocis/pull/6137

* Enhancement - Add the email HTML templates: [#6147](https://github.com/owncloud/ocis/pull/6147)

   Add the email HTML templates

   https://github.com/owncloud/ocis/issues/6146
   https://github.com/owncloud/ocis/pull/6147

* Enhancement - Add debug server to idm: [#6153](https://github.com/owncloud/ocis/pull/6153)

   We added a debug server to idm.

   https://github.com/owncloud/ocis/issues/5003
   https://github.com/owncloud/ocis/pull/6153

* Enhancement - Add debug server to audit: [#6178](https://github.com/owncloud/ocis/pull/6178)

   We added a debug server to audit.

   https://github.com/owncloud/ocis/issues/5002
   https://github.com/owncloud/ocis/pull/6178

* Enhancement - Web options configuration: [#6188](https://github.com/owncloud/ocis/pull/6188)

   Hardcode web options instead of using a generic `map[string]interface{}`

   https://github.com/owncloud/ocis/pull/6188

* Enhancement - Add debug server to userlog: [#6202](https://github.com/owncloud/ocis/pull/6202)

   We added a debug server to userlog.

   https://github.com/owncloud/ocis/issues/5002
   https://github.com/owncloud/ocis/pull/6202

* Enhancement - Add debug server to postprocessing: [#6203](https://github.com/owncloud/ocis/pull/6203)

   We added a debug server to postprocessing.

   https://github.com/owncloud/ocis/issues/5002
   https://github.com/owncloud/ocis/pull/6203

* Enhancement - Add debug server to eventhistory: [#6204](https://github.com/owncloud/ocis/pull/6204)

   We added a debug server to eventhistory.

   https://github.com/owncloud/ocis/issues/5002
   https://github.com/owncloud/ocis/pull/6204

* Enhancement - Add specific result to antivirus for debugging: [#6265](https://github.com/owncloud/ocis/pull/6265)

   We added the ability to define a specific result for the virus scanner via
   env-var (ANTIVIRUS_DEBUG_SCAN_OUTCOME)

   https://github.com/owncloud/ocis/pull/6265

* Enhancement - Add Store to `postprocessing`: [#6281](https://github.com/owncloud/ocis/pull/6281)

   Add a gomicro store for the postprocessing service. Needed to run multiple
   postprocessing instances

   https://github.com/owncloud/ocis/pull/6281

* Enhancement - Update web to v7.0.0-rc.37: [#6294](https://github.com/owncloud/ocis/pull/6294)

   Tags: web

   We updated ownCloud Web to v7.0.0-rc.37. Please refer to the changelog (linked)
   for details on the web release.

  * Bugfix [owncloud/web#6423](https://github.com/owncloud/web/issues/6423): Archiver in protected public links
  * Bugfix [owncloud/web#6434](https://github.com/owncloud/web/issues/6434): Endless lazy loading indicator after sorting file table
  * Bugfix [owncloud/web#6731](https://github.com/owncloud/web/issues/6731): Layout with long breadcrumb
  * Bugfix [owncloud/web#6768](https://github.com/owncloud/web/issues/6768): Pagination after increasing items per page
  * Bugfix [owncloud/web#7513](https://github.com/owncloud/web/issues/7513): Calendar popup position in right sidebar
  * Bugfix [owncloud/web#7655](https://github.com/owncloud/web/issues/7655): Loading shares in deep nested folders
  * Bugfix [owncloud/web#7925](https://github.com/owncloud/web/pull/7925): "Paste"-action without write permissions
  * Bugfix [owncloud/web#7926](https://github.com/owncloud/web/pull/7926): Include spaces in the list info
  * Bugfix [owncloud/web#7958](https://github.com/owncloud/web/pull/7958): Prevent deletion of own account
  * Bugfix [owncloud/web#7966](https://github.com/owncloud/web/pull/7966): UI fixes for sorting and quickactions
  * Bugfix [owncloud/web#7969](https://github.com/owncloud/web/pull/7969): Space quota not displayed after creation
  * Bugfix [owncloud/web#8026](https://github.com/owncloud/web/pull/8026): Text editor appearance
  * Bugfix [owncloud/web#8040](https://github.com/owncloud/web/pull/8040): Reverting versions for read-only shares
  * Bugfix [owncloud/web#8045](https://github.com/owncloud/web/pull/8045): Resolving drives in search
  * Bugfix [owncloud/web#8054](https://github.com/owncloud/web/issues/8054): Search repeating no results message
  * Bugfix [owncloud/web#8058](https://github.com/owncloud/web/pull/8058): Current year selection in the date picker
  * Bugfix [owncloud/web#8061](https://github.com/owncloud/web/pull/8061): Omit "page"-query in breadcrumb navigation
  * Bugfix [owncloud/web#8080](https://github.com/owncloud/web/pull/8080): Left sidebar navigation item text flickers on transition
  * Bugfix [owncloud/web#8081](https://github.com/owncloud/web/issues/8081): Space member disappearing
  * Bugfix [owncloud/web#8083](https://github.com/owncloud/web/issues/8083): Re-using space images
  * Bugfix [owncloud/web#8148](https://github.com/owncloud/web/issues/8148): Show space members despite deleted entries
  * Bugfix [owncloud/web#8158](https://github.com/owncloud/web/issues/8158): Search bar input appearance
  * Bugfix [owncloud/web#8265](https://github.com/owncloud/web/pull/8265): Application menu active display on hover
  * Bugfix [owncloud/web#8276](https://github.com/owncloud/web/pull/8276): Loading additional user data
  * Bugfix [owncloud/web#8300](https://github.com/owncloud/web/pull/8300): Re-loading space members panel
  * Bugfix [owncloud/web#8326](https://github.com/owncloud/web/pull/8326): Editing users who never logged in
  * Bugfix [owncloud/web#8340](https://github.com/owncloud/web/pull/8340): Cancel custom permissions
  * Bugfix [owncloud/web#8411](https://github.com/owncloud/web/issues/8411): Drop menus with limited vertical screen space
  * Bugfix [owncloud/web#8420](https://github.com/owncloud/web/issues/8420): Token renewal in vue router hash mode
  * Bugfix [owncloud/web#8434](https://github.com/owncloud/web/issues/8434): Accessing route in admin-settings with insufficient permissions
  * Bugfix [owncloud/web#8479](https://github.com/owncloud/web/issues/8479): "Show more"-action in shares panel
  * Bugfix [owncloud/web#8480](https://github.com/owncloud/web/pull/8480): Paste action conflict dialog broken
  * Bugfix [owncloud/web#8498](https://github.com/owncloud/web/pull/8498): PDF display issue - Update CSP object-src policy
  * Bugfix [owncloud/web#8508](https://github.com/owncloud/web/pull/8508): Remove fuzzy search results
  * Bugfix [owncloud/web#8523](https://github.com/owncloud/web/issues/8523): Space image upload
  * Bugfix [owncloud/web#8549](https://github.com/owncloud/web/issues/8549): Batch context actions in admin settings
  * Bugfix [owncloud/web#8554](https://github.com/owncloud/web/pull/8554): Height of dropdown no-option
  * Bugfix [owncloud/web#8576](https://github.com/owncloud/web/pull/8576): De-duplicate event handling to prevent errors on Draw-io
  * Bugfix [owncloud/web#8585](https://github.com/owncloud/web/issues/8585): Users without role assignment
  * Bugfix [owncloud/web#8587](https://github.com/owncloud/web/issues/8587): Password enforced check for public links
  * Bugfix [owncloud/web#8592](https://github.com/owncloud/web/issues/8592): Group members sorting
  * Bugfix [owncloud/web#8694](https://github.com/owncloud/web/pull/8694): Broken re-login after logout
  * Bugfix [owncloud/web#8695](https://github.com/owncloud/web/issues/8695): Open files in external app
  * Bugfix [owncloud/web#8756](https://github.com/owncloud/web/pull/8756): Copy link to clipboard text
  * Bugfix [owncloud/web#8758](https://github.com/owncloud/web/pull/8758): Preview controls colors
  * Bugfix [owncloud/web#8776](https://github.com/owncloud/web/issues/8776): Selection reset on action click
  * Bugfix [owncloud/web#8814](https://github.com/owncloud/web/pull/8814): Share recipient container exceed
  * Bugfix [owncloud/web#8825](https://github.com/owncloud/web/pull/8825): Remove drop target in read-only folders
  * Bugfix [owncloud/web#8827](https://github.com/owncloud/web/pull/8827): Opening context menu via keyboard
  * Bugfix [owncloud/web#8834](https://github.com/owncloud/web/issues/8834): Hide upload hint in empty read-only folders
  * Bugfix [owncloud/web#8864](https://github.com/owncloud/web/pull/8864): Public link empty password stays forever
  * Bugfix [owncloud/web#8880](https://github.com/owncloud/web/issues/8880): Sidebar header after deleting resource
  * Bugfix [owncloud/web#8928](https://github.com/owncloud/web/issues/8928): Infinite login redirect
  * Bugfix [owncloud/web#8987](https://github.com/owncloud/web/pull/8987): Limit amount of concurrent tus requests
  * Bugfix [owncloud/web#8992](https://github.com/owncloud/web/pull/8992): Personal space name after language change
  * Bugfix [owncloud/web#9004](https://github.com/owncloud/web/issues/9004): Endless loading when encountering a public link error
  * Bugfix [owncloud/web#9015](https://github.com/owncloud/web/pull/9015): Prevent "virtual" spaces from being displayed in the UI
  * Change [owncloud/web#6661](https://github.com/owncloud/web/issues/6661): Streamline new tab handling in extensions
  * Change [owncloud/web#7948](https://github.com/owncloud/web/issues/7948): Update Vue to v3.2
  * Change [owncloud/web#8431](https://github.com/owncloud/web/pull/8431): Remove permission manager
  * Change [owncloud/web#8455](https://github.com/owncloud/web/pull/8455): Configurable extension autosave
  * Change [owncloud/web#8563](https://github.com/owncloud/web/pull/8563): Theme colors
  * Enhancement [owncloud/web#6183](https://github.com/owncloud/web/issues/6183): Global loading indicator
  * Enhancement [owncloud/web#7388](https://github.com/owncloud/web/pull/7388): Add tag support
  * Enhancement [owncloud/web#7721](https://github.com/owncloud/web/issues/7721): Improve performance when loading folders and share indicators
  * Enhancement [owncloud/web#7942](https://github.com/owncloud/web/pull/7942): Warn users when using unsupported browsers
  * Enhancement [owncloud/web#7965](https://github.com/owncloud/web/pull/7965): Optional Contributor role and configurable resharing permissions
  * Enhancement [owncloud/web#7968](https://github.com/owncloud/web/pull/7968): Group and user creation forms submit on enter
  * Enhancement [owncloud/web#7976](https://github.com/owncloud/web/pull/7976): Add switch to enable condensed resource table
  * Enhancement [owncloud/web#7977](https://github.com/owncloud/web/pull/7977): Introduce zoom and rotate to the preview app
  * Enhancement [owncloud/web#7983](https://github.com/owncloud/web/pull/7983): Conflict dialog UX
  * Enhancement [owncloud/web#7991](https://github.com/owncloud/web/pull/7991): Add tiles view for resource display
  * Enhancement [owncloud/web#7994](https://github.com/owncloud/web/pull/7994): Introduce full screen mode to the preview app
  * Enhancement [owncloud/web#7995](https://github.com/owncloud/web/pull/7995): Enable autoplay in the preview app
  * Enhancement [owncloud/web#8008](https://github.com/owncloud/web/issues/8008): Don't open sidebar when copying quicklink
  * Enhancement [owncloud/web#8021](https://github.com/owncloud/web/pull/8021): Access right sidebar panels via URL
  * Enhancement [owncloud/web#8051](https://github.com/owncloud/web/pull/8051): Introduce image preloading to the preview app
  * Enhancement [owncloud/web#8055](https://github.com/owncloud/web/pull/8055): Retry failed uploads on re-upload
  * Enhancement [owncloud/web#8056](https://github.com/owncloud/web/pull/8056): Increase Searchbar height
  * Enhancement [owncloud/web#8057](https://github.com/owncloud/web/pull/8057): Show text file icon for empty text files
  * Enhancement [owncloud/web#8132](https://github.com/owncloud/web/pull/8132): Update libre-graph-api to v1.0
  * Enhancement [owncloud/web#8136](https://github.com/owncloud/web/pull/8136): Make clipboard copy available to more browsers
  * Enhancement [owncloud/web#8161](https://github.com/owncloud/web/pull/8161): Space group members
  * Enhancement [owncloud/web#8161](https://github.com/owncloud/web/pull/8161): Space group shares
  * Enhancement [owncloud/web#8166](https://github.com/owncloud/web/issues/8166): Show upload speed
  * Enhancement [owncloud/web#8175](https://github.com/owncloud/web/pull/8175): Rename "user management" app
  * Enhancement [owncloud/web#8178](https://github.com/owncloud/web/pull/8178): Spaces list in admin settings
  * Enhancement [owncloud/web#8261](https://github.com/owncloud/web/pull/8261): Admin settings users section uses graph api for role assignments
  * Enhancement [owncloud/web#8279](https://github.com/owncloud/web/pull/8279): Move user group select to edit panel
  * Enhancement [owncloud/web#8280](https://github.com/owncloud/web/pull/8280): Add support for multiple clients in `theme.json`
  * Enhancement [owncloud/web#8294](https://github.com/owncloud/web/pull/8294): Move language selection to user account page
  * Enhancement [owncloud/web#8306](https://github.com/owncloud/web/pull/8306): Show selectable groups only
  * Enhancement [owncloud/web#8317](https://github.com/owncloud/web/pull/8317): Add context menu to groups
  * Enhancement [owncloud/web#8320](https://github.com/owncloud/web/pull/8320): Space member expiration
  * Enhancement [owncloud/web#8320](https://github.com/owncloud/web/pull/8320): Update SDK to v3.1.0-alpha.3
  * Enhancement [owncloud/web#8324](https://github.com/owncloud/web/pull/8324): Add context menu to users
  * Enhancement [owncloud/web#8331](https://github.com/owncloud/web/pull/8331): Admin settings users section details improvement
  * Enhancement [owncloud/web#8354](https://github.com/owncloud/web/issues/8354): Add `ItemFilter` component
  * Enhancement [owncloud/web#8356](https://github.com/owncloud/web/pull/8356): Slight improvement of key up/down performance
  * Enhancement [owncloud/web#8363](https://github.com/owncloud/web/issues/8363): Admin settings general section
  * Enhancement [owncloud/web#8375](https://github.com/owncloud/web/pull/8375): Add appearance section in general settings
  * Enhancement [owncloud/web#8377](https://github.com/owncloud/web/issues/8377): User group filter
  * Enhancement [owncloud/web#8387](https://github.com/owncloud/web/pull/8387): Batch edit quota in admin panel
  * Enhancement [owncloud/web#8398](https://github.com/owncloud/web/pull/8398): Use standardized layout for file/space action list
  * Enhancement [owncloud/web#8425](https://github.com/owncloud/web/issues/8425): Add dark ownCloud logo
  * Enhancement [owncloud/web#8432](https://github.com/owncloud/web/pull/8432): Inject customizations
  * Enhancement [owncloud/web#8433](https://github.com/owncloud/web/pull/8433): User settings login field
  * Enhancement [owncloud/web#8441](https://github.com/owncloud/web/pull/8441): Skeleton App
  * Enhancement [owncloud/web#8449](https://github.com/owncloud/web/pull/8449): Configurable top bar
  * Enhancement [owncloud/web#8450](https://github.com/owncloud/web/pull/8450): Rework notification bell
  * Enhancement [owncloud/web#8455](https://github.com/owncloud/web/pull/8455): Autosave content changes in text editor
  * Enhancement [owncloud/web#8473](https://github.com/owncloud/web/pull/8473): Update CERN links
  * Enhancement [owncloud/web#8489](https://github.com/owncloud/web/pull/8489): Respect max quota
  * Enhancement [owncloud/web#8492](https://github.com/owncloud/web/pull/8492): User role filter
  * Enhancement [owncloud/web#8503](https://github.com/owncloud/web/issues/8503): Beautify file version list
  * Enhancement [owncloud/web#8515](https://github.com/owncloud/web/pull/8515): Introduce trashbin overview
  * Enhancement [owncloud/web#8518](https://github.com/owncloud/web/pull/8518): Make notifications work with oCIS
  * Enhancement [owncloud/web#8541](https://github.com/owncloud/web/pull/8541): Public link permission `PublicLink.Write.all`
  * Enhancement [owncloud/web#8553](https://github.com/owncloud/web/pull/8553): Add and remove users from groups batch actions
  * Enhancement [owncloud/web#8554](https://github.com/owncloud/web/pull/8554): Beautify form inputs
  * Enhancement [owncloud/web#8557](https://github.com/owncloud/web/issues/8557): Rework mobile navigation
  * Enhancement [owncloud/web#8566](https://github.com/owncloud/web/pull/8566): QuickActions role configurable
  * Enhancement [owncloud/web#8612](https://github.com/owncloud/web/issues/8612): Add `Accept-Language` header to all outgoing requests
  * Enhancement [owncloud/web#8630](https://github.com/owncloud/web/pull/8630): Add logout url
  * Enhancement [owncloud/web#8652](https://github.com/owncloud/web/pull/8652): Enable guest users
  * Enhancement [owncloud/web#8711](https://github.com/owncloud/web/pull/8711): Remove placeholder, add customizable label
  * Enhancement [owncloud/web#8713](https://github.com/owncloud/web/pull/8713): Context helper read more link configurable
  * Enhancement [owncloud/web#8715](https://github.com/owncloud/web/pull/8715): Enable rename groups
  * Enhancement [owncloud/web#8730](https://github.com/owncloud/web/pull/8730): Create Space from selection
  * Enhancement [owncloud/web#8738](https://github.com/owncloud/web/issues/8738): GDPR export
  * Enhancement [owncloud/web#8762](https://github.com/owncloud/web/pull/8762): Stop bootstrapping application earlier in anonymous contexts
  * Enhancement [owncloud/web#8766](https://github.com/owncloud/web/pull/8766): Add support for read-only groups
  * Enhancement [owncloud/web#8790](https://github.com/owncloud/web/pull/8790): Custom translations
  * Enhancement [owncloud/web#8797](https://github.com/owncloud/web/pull/8797): Font family in theming
  * Enhancement [owncloud/web#8806](https://github.com/owncloud/web/pull/8806): Preview app sorting
  * Enhancement [owncloud/web#8820](https://github.com/owncloud/web/pull/8820): Adjust missing reshare permissions message
  * Enhancement [owncloud/web#8822](https://github.com/owncloud/web/pull/8822): Fix quicklink icon alignment
  * Enhancement [owncloud/web#8826](https://github.com/owncloud/web/pull/8826): Admin settings groups members panel
  * Enhancement [owncloud/web#8868](https://github.com/owncloud/web/pull/8868): Respect user read-only configuration by the server
  * Enhancement [owncloud/web#8876](https://github.com/owncloud/web/pull/8876): Update roles and permissions names, labels, texts and icons
  * Enhancement [owncloud/web#8882](https://github.com/owncloud/web/pull/8882): Layout of Share role and expiration date dropdown
  * Enhancement [owncloud/web#8883](https://github.com/owncloud/web/issues/8883): Webfinger redirect app
  * Enhancement [owncloud/web#8898](https://github.com/owncloud/web/pull/8898): Rename "Quicklink" to "link"
  * Enhancement [owncloud/web#8911](https://github.com/owncloud/web/pull/8911): Add notification setting to account page

   https://github.com/owncloud/ocis/pull/6294
   https://github.com/owncloud/web/releases/tag/v7.0.0-rc.37

* Enhancement - Remove quota from share jails api responses: [#6309](https://github.com/owncloud/ocis/pull/6309)

   We have removed the quota object from api responses for share jails, which would
   permanently show exceeded due to restrictions in the permission system.

   https://github.com/owncloud/ocis/issues/4472
   https://github.com/owncloud/ocis/pull/6309

* Enhancement - Graph user capabilities: [#6339](https://github.com/owncloud/ocis/pull/6339)

   Adds capablities to show if users are writeable in LDAP so clients can block
   their specific fields

   https://github.com/owncloud/ocis/pull/6339

* Enhancement - Configurable ID Cache: [#6353](https://github.com/owncloud/ocis/pull/6353)

   Makes the integrated idcache (used to reduce reads from disc) configurable with
   the general cache envvars

   https://github.com/owncloud/ocis/pull/6353

* Enhancement - Fix err when the user share the locked file: [#6358](https://github.com/owncloud/ocis/pull/6358)

   Fix unexpected behavior when the user try to share the locked file

   https://github.com/owncloud/ocis/issues/6197
   https://github.com/owncloud/ocis/pull/6358

* Enhancement - Remove the email logo: [#6359](https://github.com/owncloud/ocis/issues/6359)

   Remove the email logo

   https://github.com/owncloud/ocis/issues/6359
   https://github.com/owncloud/ocis/pull/6361

* Enhancement - Default LDAP write to true: [#6362](https://github.com/owncloud/ocis/pull/6362)

   Default `OCIS_LDAP_SERVER_WRITE_ENABLED` to true

   https://github.com/owncloud/ocis/pull/6362

* Enhancement - Add fulltextsearch capabilty: [#6366](https://github.com/owncloud/ocis/pull/6366)

   It needs an extra envvar `FRONTEND_FULL_TEXT_SEARCH_ENABLED`

   https://github.com/owncloud/ocis/pull/6366

* Enhancement - Update web to v7.0.0-rc.38: [#6375](https://github.com/owncloud/ocis/pull/6375)

   Tags: web

   We updated ownCloud Web to v7.0.0-rc.38. Please refer to the changelog (linked)
   for details on the web release.

  * Bugfix [owncloud/web#6423](https://github.com/owncloud/web/issues/6423): Archiver in protected public links
  * Bugfix [owncloud/web#6434](https://github.com/owncloud/web/issues/6434): Endless lazy loading indicator after sorting file table
  * Bugfix [owncloud/web#6731](https://github.com/owncloud/web/issues/6731): Layout with long breadcrumb
  * Bugfix [owncloud/web#6768](https://github.com/owncloud/web/issues/6768): Pagination after increasing items per page
  * Bugfix [owncloud/web#7513](https://github.com/owncloud/web/issues/7513): Calendar popup position in right sidebar
  * Bugfix [owncloud/web#7655](https://github.com/owncloud/web/issues/7655): Loading shares in deep nested folders
  * Bugfix [owncloud/web#7925](https://github.com/owncloud/web/pull/7925): "Paste"-action without write permissions
  * Bugfix [owncloud/web#7926](https://github.com/owncloud/web/pull/7926): Include spaces in the list info
  * Bugfix [owncloud/web#7958](https://github.com/owncloud/web/pull/7958): Prevent deletion of own account
  * Bugfix [owncloud/web#7966](https://github.com/owncloud/web/pull/7966): UI fixes for sorting and quickactions
  * Bugfix [owncloud/web#7969](https://github.com/owncloud/web/pull/7969): Space quota not displayed after creation
  * Bugfix [owncloud/web#8026](https://github.com/owncloud/web/pull/8026): Text editor appearance
  * Bugfix [owncloud/web#8040](https://github.com/owncloud/web/pull/8040): Reverting versions for read-only shares
  * Bugfix [owncloud/web#8045](https://github.com/owncloud/web/pull/8045): Resolving drives in search
  * Bugfix [owncloud/web#8054](https://github.com/owncloud/web/issues/8054): Search repeating no results message
  * Bugfix [owncloud/web#8058](https://github.com/owncloud/web/pull/8058): Current year selection in the date picker
  * Bugfix [owncloud/web#8061](https://github.com/owncloud/web/pull/8061): Omit "page"-query in breadcrumb navigation
  * Bugfix [owncloud/web#8080](https://github.com/owncloud/web/pull/8080): Left sidebar navigation item text flickers on transition
  * Bugfix [owncloud/web#8081](https://github.com/owncloud/web/issues/8081): Space member disappearing
  * Bugfix [owncloud/web#8083](https://github.com/owncloud/web/issues/8083): Re-using space images
  * Bugfix [owncloud/web#8148](https://github.com/owncloud/web/issues/8148): Show space members despite deleted entries
  * Bugfix [owncloud/web#8158](https://github.com/owncloud/web/issues/8158): Search bar input appearance
  * Bugfix [owncloud/web#8265](https://github.com/owncloud/web/pull/8265): Application menu active display on hover
  * Bugfix [owncloud/web#8276](https://github.com/owncloud/web/pull/8276): Loading additional user data
  * Bugfix [owncloud/web#8300](https://github.com/owncloud/web/pull/8300): Re-loading space members panel
  * Bugfix [owncloud/web#8326](https://github.com/owncloud/web/pull/8326): Editing users who never logged in
  * Bugfix [owncloud/web#8340](https://github.com/owncloud/web/pull/8340): Cancel custom permissions
  * Bugfix [owncloud/web#8411](https://github.com/owncloud/web/issues/8411): Drop menus with limited vertical screen space
  * Bugfix [owncloud/web#8420](https://github.com/owncloud/web/issues/8420): Token renewal in vue router hash mode
  * Bugfix [owncloud/web#8434](https://github.com/owncloud/web/issues/8434): Accessing route in admin-settings with insufficient permissions
  * Bugfix [owncloud/web#8479](https://github.com/owncloud/web/issues/8479): "Show more"-action in shares panel
  * Bugfix [owncloud/web#8480](https://github.com/owncloud/web/pull/8480): Paste action conflict dialog broken
  * Bugfix [owncloud/web#8498](https://github.com/owncloud/web/pull/8498): PDF display issue - Update CSP object-src policy
  * Bugfix [owncloud/web#8508](https://github.com/owncloud/web/pull/8508): Remove fuzzy search results
  * Bugfix [owncloud/web#8523](https://github.com/owncloud/web/issues/8523): Space image upload
  * Bugfix [owncloud/web#8549](https://github.com/owncloud/web/issues/8549): Batch context actions in admin settings
  * Bugfix [owncloud/web#8554](https://github.com/owncloud/web/pull/8554): Height of dropdown no-option
  * Bugfix [owncloud/web#8576](https://github.com/owncloud/web/pull/8576): De-duplicate event handling to prevent errors on Draw-io
  * Bugfix [owncloud/web#8585](https://github.com/owncloud/web/issues/8585): Users without role assignment
  * Bugfix [owncloud/web#8587](https://github.com/owncloud/web/issues/8587): Password enforced check for public links
  * Bugfix [owncloud/web#8592](https://github.com/owncloud/web/issues/8592): Group members sorting
  * Bugfix [owncloud/web#8694](https://github.com/owncloud/web/pull/8694): Broken re-login after logout
  * Bugfix [owncloud/web#8695](https://github.com/owncloud/web/issues/8695): Open files in external app
  * Bugfix [owncloud/web#8756](https://github.com/owncloud/web/pull/8756): Copy link to clipboard text
  * Bugfix [owncloud/web#8758](https://github.com/owncloud/web/pull/8758): Preview controls colors
  * Bugfix [owncloud/web#8776](https://github.com/owncloud/web/issues/8776): Selection reset on action click
  * Bugfix [owncloud/web#8814](https://github.com/owncloud/web/pull/8814): Share recipient container exceed
  * Bugfix [owncloud/web#8825](https://github.com/owncloud/web/pull/8825): Remove drop target in read-only folders
  * Bugfix [owncloud/web#8827](https://github.com/owncloud/web/pull/8827): Opening context menu via keyboard
  * Bugfix [owncloud/web#8834](https://github.com/owncloud/web/issues/8834): Hide upload hint in empty read-only folders
  * Bugfix [owncloud/web#8864](https://github.com/owncloud/web/pull/8864): Public link empty password stays forever
  * Bugfix [owncloud/web#8880](https://github.com/owncloud/web/issues/8880): Sidebar header after deleting resource
  * Bugfix [owncloud/web#8928](https://github.com/owncloud/web/issues/8928): Infinite login redirect
  * Bugfix [owncloud/web#8987](https://github.com/owncloud/web/pull/8987): Limit amount of concurrent tus requests
  * Bugfix [owncloud/web#8992](https://github.com/owncloud/web/pull/8992): Personal space name after language change
  * Bugfix [owncloud/web#9004](https://github.com/owncloud/web/issues/9004): Endless loading when encountering a public link error
  * Bugfix [owncloud/web#9015](https://github.com/owncloud/web/pull/9015): Prevent "virtual" spaces from being displayed in the UI
  * Bugfix [owncloud/web#9022](https://github.com/owncloud/web/issues/9022): Spaces in search results
  * Bugfix [owncloud/web#9061](https://github.com/owncloud/web/issues/9061): Resource not found and No content message at the same time
  * Change [owncloud/web#6661](https://github.com/owncloud/web/issues/6661): Streamline new tab handling in extensions
  * Change [owncloud/web#7948](https://github.com/owncloud/web/issues/7948): Update Vue to v3.2
  * Change [owncloud/web#8431](https://github.com/owncloud/web/pull/8431): Remove permission manager
  * Change [owncloud/web#8455](https://github.com/owncloud/web/pull/8455): Configurable extension autosave
  * Change [owncloud/web#8563](https://github.com/owncloud/web/pull/8563): Theme colors
  * Enhancement [owncloud/web#6183](https://github.com/owncloud/web/issues/6183): Global loading indicator
  * Enhancement [owncloud/web#7388](https://github.com/owncloud/web/pull/7388): Add tag support
  * Enhancement [owncloud/web#7721](https://github.com/owncloud/web/issues/7721): Improve performance when loading folders and share indicators
  * Enhancement [owncloud/web#7942](https://github.com/owncloud/web/pull/7942): Warn users when using unsupported browsers
  * Enhancement [owncloud/web#7965](https://github.com/owncloud/web/pull/7965): Optional Contributor role and configurable resharing permissions
  * Enhancement [owncloud/web#7968](https://github.com/owncloud/web/pull/7968): Group and user creation forms submit on enter
  * Enhancement [owncloud/web#7976](https://github.com/owncloud/web/pull/7976): Add switch to enable condensed resource table
  * Enhancement [owncloud/web#7977](https://github.com/owncloud/web/pull/7977): Introduce zoom and rotate to the preview app
  * Enhancement [owncloud/web#7983](https://github.com/owncloud/web/pull/7983): Conflict dialog UX
  * Enhancement [owncloud/web#7991](https://github.com/owncloud/web/pull/7991): Add tiles view for resource display
  * Enhancement [owncloud/web#7994](https://github.com/owncloud/web/pull/7994): Introduce full screen mode to the preview app
  * Enhancement [owncloud/web#7995](https://github.com/owncloud/web/pull/7995): Enable autoplay in the preview app
  * Enhancement [owncloud/web#8008](https://github.com/owncloud/web/issues/8008): Don't open sidebar when copying quicklink
  * Enhancement [owncloud/web#8021](https://github.com/owncloud/web/pull/8021): Access right sidebar panels via URL
  * Enhancement [owncloud/web#8051](https://github.com/owncloud/web/pull/8051): Introduce image preloading to the preview app
  * Enhancement [owncloud/web#8055](https://github.com/owncloud/web/pull/8055): Retry failed uploads on re-upload
  * Enhancement [owncloud/web#8056](https://github.com/owncloud/web/pull/8056): Increase Searchbar height
  * Enhancement [owncloud/web#8057](https://github.com/owncloud/web/pull/8057): Show text file icon for empty text files
  * Enhancement [owncloud/web#8132](https://github.com/owncloud/web/pull/8132): Update libre-graph-api to v1.0
  * Enhancement [owncloud/web#8136](https://github.com/owncloud/web/pull/8136): Make clipboard copy available to more browsers
  * Enhancement [owncloud/web#8161](https://github.com/owncloud/web/pull/8161): Space group members
  * Enhancement [owncloud/web#8161](https://github.com/owncloud/web/pull/8161): Space group shares
  * Enhancement [owncloud/web#8166](https://github.com/owncloud/web/issues/8166): Show upload speed
  * Enhancement [owncloud/web#8175](https://github.com/owncloud/web/pull/8175): Rename "user management" app
  * Enhancement [owncloud/web#8178](https://github.com/owncloud/web/pull/8178): Spaces list in admin settings
  * Enhancement [owncloud/web#8261](https://github.com/owncloud/web/pull/8261): Admin settings users section uses graph api for role assignments
  * Enhancement [owncloud/web#8279](https://github.com/owncloud/web/pull/8279): Move user group select to edit panel
  * Enhancement [owncloud/web#8280](https://github.com/owncloud/web/pull/8280): Add support for multiple clients in `theme.json`
  * Enhancement [owncloud/web#8294](https://github.com/owncloud/web/pull/8294): Move language selection to user account page
  * Enhancement [owncloud/web#8306](https://github.com/owncloud/web/pull/8306): Show selectable groups only
  * Enhancement [owncloud/web#8317](https://github.com/owncloud/web/pull/8317): Add context menu to groups
  * Enhancement [owncloud/web#8320](https://github.com/owncloud/web/pull/8320): Space member expiration
  * Enhancement [owncloud/web#8320](https://github.com/owncloud/web/pull/8320): Update SDK to v3.1.0-alpha.3
  * Enhancement [owncloud/web#8324](https://github.com/owncloud/web/pull/8324): Add context menu to users
  * Enhancement [owncloud/web#8331](https://github.com/owncloud/web/pull/8331): Admin settings users section details improvement
  * Enhancement [owncloud/web#8354](https://github.com/owncloud/web/issues/8354): Add `ItemFilter` component
  * Enhancement [owncloud/web#8356](https://github.com/owncloud/web/pull/8356): Slight improvement of key up/down performance
  * Enhancement [owncloud/web#8363](https://github.com/owncloud/web/issues/8363): Admin settings general section
  * Enhancement [owncloud/web#8375](https://github.com/owncloud/web/pull/8375): Add appearance section in general settings
  * Enhancement [owncloud/web#8377](https://github.com/owncloud/web/issues/8377): User group filter
  * Enhancement [owncloud/web#8387](https://github.com/owncloud/web/pull/8387): Batch edit quota in admin panel
  * Enhancement [owncloud/web#8398](https://github.com/owncloud/web/pull/8398): Use standardized layout for file/space action list
  * Enhancement [owncloud/web#8425](https://github.com/owncloud/web/issues/8425): Add dark ownCloud logo
  * Enhancement [owncloud/web#8432](https://github.com/owncloud/web/pull/8432): Inject customizations
  * Enhancement [owncloud/web#8433](https://github.com/owncloud/web/pull/8433): User settings login field
  * Enhancement [owncloud/web#8441](https://github.com/owncloud/web/pull/8441): Skeleton App
  * Enhancement [owncloud/web#8449](https://github.com/owncloud/web/pull/8449): Configurable top bar
  * Enhancement [owncloud/web#8450](https://github.com/owncloud/web/pull/8450): Rework notification bell
  * Enhancement [owncloud/web#8455](https://github.com/owncloud/web/pull/8455): Autosave content changes in text editor
  * Enhancement [owncloud/web#8473](https://github.com/owncloud/web/pull/8473): Update CERN links
  * Enhancement [owncloud/web#8489](https://github.com/owncloud/web/pull/8489): Respect max quota
  * Enhancement [owncloud/web#8492](https://github.com/owncloud/web/pull/8492): User role filter
  * Enhancement [owncloud/web#8503](https://github.com/owncloud/web/issues/8503): Beautify file version list
  * Enhancement [owncloud/web#8515](https://github.com/owncloud/web/pull/8515): Introduce trashbin overview
  * Enhancement [owncloud/web#8518](https://github.com/owncloud/web/pull/8518): Make notifications work with oCIS
  * Enhancement [owncloud/web#8541](https://github.com/owncloud/web/pull/8541): Public link permission `PublicLink.Write.all`
  * Enhancement [owncloud/web#8553](https://github.com/owncloud/web/pull/8553): Add and remove users from groups batch actions
  * Enhancement [owncloud/web#8554](https://github.com/owncloud/web/pull/8554): Beautify form inputs
  * Enhancement [owncloud/web#8557](https://github.com/owncloud/web/issues/8557): Rework mobile navigation
  * Enhancement [owncloud/web#8566](https://github.com/owncloud/web/pull/8566): QuickActions role configurable
  * Enhancement [owncloud/web#8612](https://github.com/owncloud/web/issues/8612): Add `Accept-Language` header to all outgoing requests
  * Enhancement [owncloud/web#8630](https://github.com/owncloud/web/pull/8630): Add logout url
  * Enhancement [owncloud/web#8652](https://github.com/owncloud/web/pull/8652): Enable guest users
  * Enhancement [owncloud/web#8711](https://github.com/owncloud/web/pull/8711): Remove placeholder, add customizable label
  * Enhancement [owncloud/web#8713](https://github.com/owncloud/web/pull/8713): Context helper read more link configurable
  * Enhancement [owncloud/web#8715](https://github.com/owncloud/web/pull/8715): Enable rename groups
  * Enhancement [owncloud/web#8730](https://github.com/owncloud/web/pull/8730): Create Space from selection
  * Enhancement [owncloud/web#8738](https://github.com/owncloud/web/issues/8738): GDPR export
  * Enhancement [owncloud/web#8762](https://github.com/owncloud/web/pull/8762): Stop bootstrapping application earlier in anonymous contexts
  * Enhancement [owncloud/web#8766](https://github.com/owncloud/web/pull/8766): Add support for read-only groups
  * Enhancement [owncloud/web#8790](https://github.com/owncloud/web/pull/8790): Custom translations
  * Enhancement [owncloud/web#8797](https://github.com/owncloud/web/pull/8797): Font family in theming
  * Enhancement [owncloud/web#8806](https://github.com/owncloud/web/pull/8806): Preview app sorting
  * Enhancement [owncloud/web#8820](https://github.com/owncloud/web/pull/8820): Adjust missing reshare permissions message
  * Enhancement [owncloud/web#8822](https://github.com/owncloud/web/pull/8822): Fix quicklink icon alignment
  * Enhancement [owncloud/web#8826](https://github.com/owncloud/web/pull/8826): Admin settings groups members panel
  * Enhancement [owncloud/web#8868](https://github.com/owncloud/web/pull/8868): Respect user read-only configuration by the server
  * Enhancement [owncloud/web#8876](https://github.com/owncloud/web/pull/8876): Update roles and permissions names, labels, texts and icons
  * Enhancement [owncloud/web#8882](https://github.com/owncloud/web/pull/8882): Layout of Share role and expiration date dropdown
  * Enhancement [owncloud/web#8883](https://github.com/owncloud/web/issues/8883): Webfinger redirect app
  * Enhancement [owncloud/web#8898](https://github.com/owncloud/web/pull/8898): Rename "Quicklink" to "link"
  * Enhancement [owncloud/web#8911](https://github.com/owncloud/web/pull/8911): Add notification setting to account page
  * Enhancement [owncloud/web#9070](https://github.com/owncloud/web/pull/9070): Disable change password capability
  * Enhancement [owncloud/web#9070](https://github.com/owncloud/web/pull/9070): Disable create user and delete user via capabilities
  * Enhancement [owncloud/web#9076](https://github.com/owncloud/web/pull/9076): Show detailed error messages while upload fails

   https://github.com/owncloud/ocis/pull/6375
   https://github.com/owncloud/web/releases/tag/v7.0.0-rc.38

* Enhancement - Fix preview or viewing of shared animated GIFs: [#6386](https://github.com/owncloud/ocis/pull/6386)

   Fix preview or viewing of shared animated GIFs

   https://github.com/owncloud/ocis/issues/5418
   https://github.com/owncloud/ocis/pull/6386

* Enhancement - Unify CA Cert envvars: [#6392](https://github.com/owncloud/ocis/pull/6392)

   Introduce a global `OCIS_EVENTS_TLS_ROOT_CA_CERTIFICATE` to avoid needing to
   configure all `{SERVICENAME}_EVENTS_TLS_ROOT_CA_CERTIFICATE` envvars

   https://github.com/owncloud/ocis/pull/6392

* Enhancement - Fix to prevent the email X-Site scripting: [#6429](https://github.com/owncloud/ocis/pull/6429)

   Fix to prevent the email notification X-Site scripting

   https://github.com/owncloud/ocis/issues/6411
   https://github.com/owncloud/ocis/pull/6429

* Enhancement - Update web to v7.0.0: [#6438](https://github.com/owncloud/ocis/pull/6438)

   Tags: web

   We updated ownCloud Web to v7.0.0. Please refer to the changelog (linked) for
   details on the web release.

   ## Breaking changes * BREAKING CHANGE for developers and admins in
   [owncloud/web#7948](https://github.com/owncloud/web/issues/7948): we've updated
   Vue.js to version 3. Existing apps that have not been updated to Vue.js version
   3 will not be compatible anymore. * BREAKING CHANGE for admins in
   [owncloud/web#8563](https://github.com/owncloud/web/pull/8563): we've introduced
   contrast colors in our theming. In case you have created a custom `theme.json`
   it needs to be adjusted accordingly: `-contrast` color values need to be added
   to all `swatches`, e.g. to `swatch-brand-contrast`. See
   https://owncloud.dev/clients/web/theming/#colors

   ## Summary * Bugfix
   [owncloud/web#6423](https://github.com/owncloud/web/issues/6423): Archiver in
   protected public links * Bugfix
   [owncloud/web#6434](https://github.com/owncloud/web/issues/6434): Endless lazy
   loading indicator after sorting file table * Bugfix
   [owncloud/web#6731](https://github.com/owncloud/web/issues/6731): Layout with
   long breadcrumb * Bugfix
   [owncloud/web#6768](https://github.com/owncloud/web/issues/6768): Pagination
   after increasing items per page * Bugfix
   [owncloud/web#7513](https://github.com/owncloud/web/issues/7513): Calendar popup
   position in right sidebar * Bugfix
   [owncloud/web#7655](https://github.com/owncloud/web/issues/7655): Loading shares
   in deep nested folders * Bugfix
   [owncloud/web#7925](https://github.com/owncloud/web/pull/7925): "Paste"-action
   without write permissions * Bugfix
   [owncloud/web#7926](https://github.com/owncloud/web/pull/7926): Include spaces
   in the list info * Bugfix
   [owncloud/web#7958](https://github.com/owncloud/web/pull/7958): Prevent deletion
   of own account * Bugfix
   [owncloud/web#7966](https://github.com/owncloud/web/pull/7966): UI fixes for
   sorting and quickactions * Bugfix
   [owncloud/web#7969](https://github.com/owncloud/web/pull/7969): Space quota not
   displayed after creation * Bugfix
   [owncloud/web#8026](https://github.com/owncloud/web/pull/8026): Text editor
   appearance * Bugfix
   [owncloud/web#8040](https://github.com/owncloud/web/pull/8040): Reverting
   versions for read-only shares * Bugfix
   [owncloud/web#8045](https://github.com/owncloud/web/pull/8045): Resolving drives
   in search * Bugfix
   [owncloud/web#8054](https://github.com/owncloud/web/issues/8054): Search
   repeating no results message * Bugfix
   [owncloud/web#8058](https://github.com/owncloud/web/pull/8058): Current year
   selection in the date picker * Bugfix
   [owncloud/web#8061](https://github.com/owncloud/web/pull/8061): Omit
   "page"-query in breadcrumb navigation * Bugfix
   [owncloud/web#8080](https://github.com/owncloud/web/pull/8080): Left sidebar
   navigation item text flickers on transition * Bugfix
   [owncloud/web#8081](https://github.com/owncloud/web/issues/8081): Space member
   disappearing * Bugfix
   [owncloud/web#8083](https://github.com/owncloud/web/issues/8083): Re-using space
   images * Bugfix
   [owncloud/web#8148](https://github.com/owncloud/web/issues/8148): Show space
   members despite deleted entries * Bugfix
   [owncloud/web#8158](https://github.com/owncloud/web/issues/8158): Search bar
   input appearance * Bugfix
   [owncloud/web#8265](https://github.com/owncloud/web/pull/8265): Application menu
   active display on hover * Bugfix
   [owncloud/web#8276](https://github.com/owncloud/web/pull/8276): Loading
   additional user data * Bugfix
   [owncloud/web#8300](https://github.com/owncloud/web/pull/8300): Re-loading space
   members panel * Bugfix
   [owncloud/web#8326](https://github.com/owncloud/web/pull/8326): Editing users
   who never logged in * Bugfix
   [owncloud/web#8340](https://github.com/owncloud/web/pull/8340): Cancel custom
   permissions * Bugfix
   [owncloud/web#8411](https://github.com/owncloud/web/issues/8411): Drop menus
   with limited vertical screen space * Bugfix
   [owncloud/web#8420](https://github.com/owncloud/web/issues/8420): Token renewal
   in vue router hash mode * Bugfix
   [owncloud/web#8434](https://github.com/owncloud/web/issues/8434): Accessing
   route in admin-settings with insufficient permissions * Bugfix
   [owncloud/web#8479](https://github.com/owncloud/web/issues/8479): "Show
   more"-action in shares panel * Bugfix
   [owncloud/web#8480](https://github.com/owncloud/web/pull/8480): Paste action
   conflict dialog broken * Bugfix
   [owncloud/web#8498](https://github.com/owncloud/web/pull/8498): PDF display
   issue - Update CSP object-src policy * Bugfix
   [owncloud/web#8508](https://github.com/owncloud/web/pull/8508): Remove fuzzy
   search results * Bugfix
   [owncloud/web#8523](https://github.com/owncloud/web/issues/8523): Space image
   upload * Bugfix
   [owncloud/web#8549](https://github.com/owncloud/web/issues/8549): Batch context
   actions in admin settings * Bugfix
   [owncloud/web#8554](https://github.com/owncloud/web/pull/8554): Height of
   dropdown no-option * Bugfix
   [owncloud/web#8576](https://github.com/owncloud/web/pull/8576): De-duplicate
   event handling to prevent errors on Draw-io * Bugfix
   [owncloud/web#8585](https://github.com/owncloud/web/issues/8585): Users without
   role assignment * Bugfix
   [owncloud/web#8587](https://github.com/owncloud/web/issues/8587): Password
   enforced check for public links * Bugfix
   [owncloud/web#8592](https://github.com/owncloud/web/issues/8592): Group members
   sorting * Bugfix [owncloud/web#8694](https://github.com/owncloud/web/pull/8694):
   Broken re-login after logout * Bugfix
   [owncloud/web#8695](https://github.com/owncloud/web/issues/8695): Open files in
   external app * Bugfix
   [owncloud/web#8756](https://github.com/owncloud/web/pull/8756): Copy link to
   clipboard text * Bugfix
   [owncloud/web#8758](https://github.com/owncloud/web/pull/8758): Preview controls
   colors * Bugfix
   [owncloud/web#8776](https://github.com/owncloud/web/issues/8776): Selection
   reset on action click * Bugfix
   [owncloud/web#8814](https://github.com/owncloud/web/pull/8814): Share recipient
   container exceed * Bugfix
   [owncloud/web#8825](https://github.com/owncloud/web/pull/8825): Remove drop
   target in read-only folders * Bugfix
   [owncloud/web#8827](https://github.com/owncloud/web/pull/8827): Opening context
   menu via keyboard * Bugfix
   [owncloud/web#8834](https://github.com/owncloud/web/issues/8834): Hide upload
   hint in empty read-only folders * Bugfix
   [owncloud/web#8864](https://github.com/owncloud/web/pull/8864): Public link
   empty password stays forever * Bugfix
   [owncloud/web#8880](https://github.com/owncloud/web/issues/8880): Sidebar header
   after deleting resource * Bugfix
   [owncloud/web#8928](https://github.com/owncloud/web/issues/8928): Infinite login
   redirect * Bugfix
   [owncloud/web#8987](https://github.com/owncloud/web/pull/8987): Limit amount of
   concurrent tus requests * Bugfix
   [owncloud/web#8992](https://github.com/owncloud/web/pull/8992): Personal space
   name after language change * Bugfix
   [owncloud/web#9004](https://github.com/owncloud/web/issues/9004): Endless
   loading when encountering a public link error * Bugfix
   [owncloud/web#9009](https://github.com/owncloud/web/pull/9009): Public link file
   previews * Bugfix
   [owncloud/web#9014](https://github.com/owncloud/web/issues/9014): Empty file
   list after deleting resources * Bugfix
   [owncloud/web#9015](https://github.com/owncloud/web/pull/9015): Prevent
   "virtual" spaces from being displayed in the UI * Bugfix
   [owncloud/web#9020](https://github.com/owncloud/web/issues/9020): Sidebar for
   spaces on "Shared via link"-page * Bugfix
   [owncloud/web#9022](https://github.com/owncloud/web/issues/9022): Spaces in
   search results * Bugfix
   [owncloud/web#9030](https://github.com/owncloud/web/issues/9030): Share
   indicator loading after pasting resources * Bugfix
   [owncloud/web#9050](https://github.com/owncloud/web/issues/9050): Preview app
   mime type detection * Bugfix
   [owncloud/web#9061](https://github.com/owncloud/web/issues/9061): Resource not
   found and No content message at the same time * Bugfix
   [owncloud/web#9080](https://github.com/owncloud/web/issues/9080): Incorrect
   pause state in upload info * Bugfix
   [owncloud/web#9131](https://github.com/owncloud/web/pull/9131): Select all
   checkbox * Bugfix
   [owncloud/web#9144](https://github.com/owncloud/web/pull/9144): Notifications
   link overflow * Change
   [owncloud/web#6661](https://github.com/owncloud/web/issues/6661): Streamline new
   tab handling in extensions * Change
   [owncloud/web#7948](https://github.com/owncloud/web/issues/7948): Update Vue to
   v3.2 * Change [owncloud/web#8431](https://github.com/owncloud/web/pull/8431):
   Remove permission manager * Change
   [owncloud/web#8455](https://github.com/owncloud/web/pull/8455): Configurable
   extension autosave * Change
   [owncloud/web#8563](https://github.com/owncloud/web/pull/8563): Theme colors *
   Enhancement [owncloud/web#6183](https://github.com/owncloud/web/issues/6183):
   Global loading indicator * Enhancement
   [owncloud/web#7388](https://github.com/owncloud/web/pull/7388): Add tag support
   * Enhancement [owncloud/web#7721](https://github.com/owncloud/web/issues/7721):
   Improve performance when loading folders and share indicators * Enhancement
   [owncloud/web#7942](https://github.com/owncloud/web/pull/7942): Warn users when
   using unsupported browsers * Enhancement
   [owncloud/web#7965](https://github.com/owncloud/web/pull/7965): Optional
   Contributor role and configurable resharing permissions * Enhancement
   [owncloud/web#7968](https://github.com/owncloud/web/pull/7968): Group and user
   creation forms submit on enter * Enhancement
   [owncloud/web#7976](https://github.com/owncloud/web/pull/7976): Add switch to
   enable condensed resource table * Enhancement
   [owncloud/web#7977](https://github.com/owncloud/web/pull/7977): Introduce zoom
   and rotate to the preview app * Enhancement
   [owncloud/web#7983](https://github.com/owncloud/web/pull/7983): Conflict dialog
   UX * Enhancement [owncloud/web#7991](https://github.com/owncloud/web/pull/7991):
   Add tiles view for resource display * Enhancement
   [owncloud/web#7994](https://github.com/owncloud/web/pull/7994): Introduce full
   screen mode to the preview app * Enhancement
   [owncloud/web#7995](https://github.com/owncloud/web/pull/7995): Enable autoplay
   in the preview app * Enhancement
   [owncloud/web#8008](https://github.com/owncloud/web/issues/8008): Don't open
   sidebar when copying quicklink * Enhancement
   [owncloud/web#8021](https://github.com/owncloud/web/pull/8021): Access right
   sidebar panels via URL * Enhancement
   [owncloud/web#8051](https://github.com/owncloud/web/pull/8051): Introduce image
   preloading to the preview app * Enhancement
   [owncloud/web#8055](https://github.com/owncloud/web/pull/8055): Retry failed
   uploads on re-upload * Enhancement
   [owncloud/web#8056](https://github.com/owncloud/web/pull/8056): Increase
   Searchbar height * Enhancement
   [owncloud/web#8057](https://github.com/owncloud/web/pull/8057): Show text file
   icon for empty text files * Enhancement
   [owncloud/web#8132](https://github.com/owncloud/web/pull/8132): Update
   libre-graph-api to v1.0 * Enhancement
   [owncloud/web#8136](https://github.com/owncloud/web/pull/8136): Make clipboard
   copy available to more browsers * Enhancement
   [owncloud/web#8161](https://github.com/owncloud/web/pull/8161): Space group
   members * Enhancement
   [owncloud/web#8161](https://github.com/owncloud/web/pull/8161): Space group
   shares * Enhancement
   [owncloud/web#8166](https://github.com/owncloud/web/issues/8166): Show upload
   speed * Enhancement
   [owncloud/web#8175](https://github.com/owncloud/web/pull/8175): Rename "user
   management" app * Enhancement
   [owncloud/web#8178](https://github.com/owncloud/web/pull/8178): Spaces list in
   admin settings * Enhancement
   [owncloud/web#8261](https://github.com/owncloud/web/pull/8261): Admin settings
   users section uses graph api for role assignments * Enhancement
   [owncloud/web#8279](https://github.com/owncloud/web/pull/8279): Move user group
   select to edit panel * Enhancement
   [owncloud/web#8280](https://github.com/owncloud/web/pull/8280): Add support for
   multiple clients in `theme.json` * Enhancement
   [owncloud/web#8294](https://github.com/owncloud/web/pull/8294): Move language
   selection to user account page * Enhancement
   [owncloud/web#8306](https://github.com/owncloud/web/pull/8306): Show selectable
   groups only * Enhancement
   [owncloud/web#8317](https://github.com/owncloud/web/pull/8317): Add context menu
   to groups * Enhancement
   [owncloud/web#8320](https://github.com/owncloud/web/pull/8320): Space member
   expiration * Enhancement
   [owncloud/web#8320](https://github.com/owncloud/web/pull/8320): Update SDK to
   v3.1.0-alpha.3 * Enhancement
   [owncloud/web#8324](https://github.com/owncloud/web/pull/8324): Add context menu
   to users * Enhancement
   [owncloud/web#8331](https://github.com/owncloud/web/pull/8331): Admin settings
   users section details improvement * Enhancement
   [owncloud/web#8354](https://github.com/owncloud/web/issues/8354): Add
   `ItemFilter` component * Enhancement
   [owncloud/web#8356](https://github.com/owncloud/web/pull/8356): Slight
   improvement of key up/down performance * Enhancement
   [owncloud/web#8363](https://github.com/owncloud/web/issues/8363): Admin settings
   general section * Enhancement
   [owncloud/web#8375](https://github.com/owncloud/web/pull/8375): Add appearance
   section in general settings * Enhancement
   [owncloud/web#8377](https://github.com/owncloud/web/issues/8377): User group
   filter * Enhancement
   [owncloud/web#8387](https://github.com/owncloud/web/pull/8387): Batch edit quota
   in admin panel * Enhancement
   [owncloud/web#8398](https://github.com/owncloud/web/pull/8398): Use standardized
   layout for file/space action list * Enhancement
   [owncloud/web#8425](https://github.com/owncloud/web/issues/8425): Add dark
   ownCloud logo * Enhancement
   [owncloud/web#8432](https://github.com/owncloud/web/pull/8432): Inject
   customizations * Enhancement
   [owncloud/web#8433](https://github.com/owncloud/web/pull/8433): User settings
   login field * Enhancement
   [owncloud/web#8441](https://github.com/owncloud/web/pull/8441): Skeleton App *
   Enhancement [owncloud/web#8449](https://github.com/owncloud/web/pull/8449):
   Configurable top bar * Enhancement
   [owncloud/web#8450](https://github.com/owncloud/web/pull/8450): Rework
   notification bell * Enhancement
   [owncloud/web#8455](https://github.com/owncloud/web/pull/8455): Autosave content
   changes in text editor * Enhancement
   [owncloud/web#8473](https://github.com/owncloud/web/pull/8473): Update CERN
   links * Enhancement
   [owncloud/web#8489](https://github.com/owncloud/web/pull/8489): Respect max
   quota * Enhancement
   [owncloud/web#8492](https://github.com/owncloud/web/pull/8492): User role filter
   * Enhancement [owncloud/web#8503](https://github.com/owncloud/web/issues/8503):
   Beautify file version list * Enhancement
   [owncloud/web#8515](https://github.com/owncloud/web/pull/8515): Introduce
   trashbin overview * Enhancement
   [owncloud/web#8518](https://github.com/owncloud/web/pull/8518): Make
   notifications work with oCIS * Enhancement
   [owncloud/web#8541](https://github.com/owncloud/web/pull/8541): Public link
   permission `PublicLink.Write.all` * Enhancement
   [owncloud/web#8553](https://github.com/owncloud/web/pull/8553): Add and remove
   users from groups batch actions * Enhancement
   [owncloud/web#8554](https://github.com/owncloud/web/pull/8554): Beautify form
   inputs * Enhancement
   [owncloud/web#8557](https://github.com/owncloud/web/issues/8557): Rework mobile
   navigation * Enhancement
   [owncloud/web#8566](https://github.com/owncloud/web/pull/8566): QuickActions
   role configurable * Enhancement
   [owncloud/web#8612](https://github.com/owncloud/web/issues/8612): Add
   `Accept-Language` header to all outgoing requests * Enhancement
   [owncloud/web#8630](https://github.com/owncloud/web/pull/8630): Add logout url *
   Enhancement [owncloud/web#8652](https://github.com/owncloud/web/pull/8652):
   Enable guest users * Enhancement
   [owncloud/web#8711](https://github.com/owncloud/web/pull/8711): Remove
   placeholder, add customizable label * Enhancement
   [owncloud/web#8713](https://github.com/owncloud/web/pull/8713): Context helper
   read more link configurable * Enhancement
   [owncloud/web#8715](https://github.com/owncloud/web/pull/8715): Enable rename
   groups * Enhancement
   [owncloud/web#8730](https://github.com/owncloud/web/pull/8730): Create Space
   from selection * Enhancement
   [owncloud/web#8738](https://github.com/owncloud/web/issues/8738): GDPR export *
   Enhancement [owncloud/web#8762](https://github.com/owncloud/web/pull/8762): Stop
   bootstrapping application earlier in anonymous contexts * Enhancement
   [owncloud/web#8766](https://github.com/owncloud/web/pull/8766): Add support for
   read-only groups * Enhancement
   [owncloud/web#8790](https://github.com/owncloud/web/pull/8790): Custom
   translations * Enhancement
   [owncloud/web#8797](https://github.com/owncloud/web/pull/8797): Font family in
   theming * Enhancement
   [owncloud/web#8806](https://github.com/owncloud/web/pull/8806): Preview app
   sorting * Enhancement
   [owncloud/web#8820](https://github.com/owncloud/web/pull/8820): Adjust missing
   reshare permissions message * Enhancement
   [owncloud/web#8822](https://github.com/owncloud/web/pull/8822): Fix quicklink
   icon alignment * Enhancement
   [owncloud/web#8826](https://github.com/owncloud/web/pull/8826): Admin settings
   groups members panel * Enhancement
   [owncloud/web#8868](https://github.com/owncloud/web/pull/8868): Respect user
   read-only configuration by the server * Enhancement
   [owncloud/web#8876](https://github.com/owncloud/web/pull/8876): Update roles and
   permissions names, labels, texts and icons * Enhancement
   [owncloud/web#8882](https://github.com/owncloud/web/pull/8882): Layout of Share
   role and expiration date dropdown * Enhancement
   [owncloud/web#8883](https://github.com/owncloud/web/issues/8883): Webfinger
   redirect app * Enhancement
   [owncloud/web#8898](https://github.com/owncloud/web/pull/8898): Rename
   "Quicklink" to "link" * Enhancement
   [owncloud/web#8911](https://github.com/owncloud/web/pull/8911): Add notification
   setting to account page * Enhancement
   [owncloud/web#9048](https://github.com/owncloud/web/issues/9048): Support
   pagination in admin settings app * Enhancement
   [owncloud/web#9070](https://github.com/owncloud/web/pull/9070): Disable change
   password capability * Enhancement
   [owncloud/web#9070](https://github.com/owncloud/web/pull/9070): Disable create
   user and delete user via capabilities * Enhancement
   [owncloud/web#9076](https://github.com/owncloud/web/pull/9076): Show detailed
   error messages while upload fails

   https://github.com/owncloud/ocis/pull/6438
   https://github.com/owncloud/web/releases/tag/v7.0.0

* Enhancement - Update Reva to version 2.14.0: [#6448](https://github.com/owncloud/ocis/pull/6448)

   Changelog for reva 2.14.0 (2023-06-05) =======================================

  *   Bugfix [cs3org/reva#3919](https://github.com/cs3org/reva/pull/3919): We added missing timestamps to events
  *   Bugfix [cs3org/reva#3911](https://github.com/cs3org/reva/pull/3911): Clean IDCache properly
  *   Bugfix [cs3org/reva#3896](https://github.com/cs3org/reva/pull/3896): Do not lose old revisions when overwriting a file during copy
  *   Bugfix [cs3org/reva#3918](https://github.com/cs3org/reva/pull/3918): Dont enumerate users
  *   Bugfix [cs3org/reva#3902](https://github.com/cs3org/reva/pull/3902): Do not try to use the cache for empty node
  *   Bugfix [cs3org/reva#3877](https://github.com/cs3org/reva/pull/3877): Empty exact list while searching for a sharee
  *   Bugfix [cs3org/reva#3906](https://github.com/cs3org/reva/pull/3906): Fix preflight requests
  *   Bugfix [cs3org/reva#3934](https://github.com/cs3org/reva/pull/3934): Fix the space editor permissions
  *   Bugfix [cs3org/reva#3899](https://github.com/cs3org/reva/pull/3899): Harden uploads
  *   Bugfix [cs3org/reva#3917](https://github.com/cs3org/reva/pull/3917): Prevent last space manager from leaving
  *   Bugfix [cs3org/reva#3866](https://github.com/cs3org/reva/pull/3866): Fix public link lookup performance
  *   Bugfix [cs3org/reva#3904](https://github.com/cs3org/reva/pull/3904): Improve performance of directory listings
  *   Enhancement [cs3org/reva#3893](https://github.com/cs3org/reva/pull/3893): Cleanup Space Delete permissions
  *   Enhancement [cs3org/reva#3894](https://github.com/cs3org/reva/pull/3894): Fix err when the user share the locked file
  *   Enhancement [cs3org/reva#3913](https://github.com/cs3org/reva/pull/3913): Introduce FullTextSearch Capability
  *   Enhancement [cs3org/reva#3898](https://github.com/cs3org/reva/pull/3898): Add Graph User capabilities
  *   Enhancement [cs3org/reva#3496](https://github.com/cs3org/reva/pull/3496): Add otlp tracing exporter
  *   Enhancement [cs3org/reva#3922](https://github.com/cs3org/reva/pull/3922): Rename permissions

   Changelog for reva 2.13.3 (2023-05-17) =======================================

  *   Bugfix [cs3org/reva#3890](https://github.com/cs3org/reva/pull/3890): Bring back public link sharing of project space roots
  *   Bugfix [cs3org/reva#3888](https://github.com/cs3org/reva/pull/3888): We fixed a bug that unnecessarily fetched all members of a group
  *   Bugfix [cs3org/reva#3886](https://github.com/cs3org/reva/pull/3886): Decomposedfs no longer deadlocks when cache is disabled
  *   Bugfix [cs3org/reva#3892](https://github.com/cs3org/reva/pull/3892): Fix public links
  *   Bugfix [cs3org/reva#3876](https://github.com/cs3org/reva/pull/3876): Remove go-micro/store/redis specific workaround
  *   Bugfix [cs3org/reva#3889](https://github.com/cs3org/reva/pull/3889): Update space root mtime when changing space metadata
  *   Bugfix [cs3org/reva#3836](https://github.com/cs3org/reva/pull/3836): Fix spaceID in the decomposedFS
  *   Bugfix [cs3org/reva#3867](https://github.com/cs3org/reva/pull/3867): Restore last version after positive result
  *   Bugfix [cs3org/reva#3849](https://github.com/cs3org/reva/pull/3849): Prevent sharing space roots and personal spaces
  *   Enhancement [cs3org/reva#3865](https://github.com/cs3org/reva/pull/3865): Remove unneccessary code from gateway
  *   Enhancement [cs3org/reva#3895](https://github.com/cs3org/reva/pull/3895): Add missing expiry date to shares

   Changelog for reva 2.13.2 (2023-05-08) =======================================

  *   Bugfix [cs3org/reva#3845](https://github.com/cs3org/reva/pull/3845): Fix propagation
  *   Bugfix [cs3org/reva#3856](https://github.com/cs3org/reva/pull/3856): Fix response code
  *   Bugfix [cs3org/reva#3857](https://github.com/cs3org/reva/pull/3857): Fix trashbin purge

   Changelog for reva 2.13.1 (2023-05-03) =======================================

  *   Bugfix [cs3org/reva#3843](https://github.com/cs3org/reva/pull/3843): Allow scope check to impersonate space owners

   Changelog for reva 2.13.0 (2023-05-02) =======================================

  *   Bugfix [cs3org/reva#3570](https://github.com/cs3org/reva/pull/3570): Return 425 on HEAD
  *   Bugfix [cs3org/reva#3830](https://github.com/cs3org/reva/pull/3830): Be more robust when logging errors
  *   Bugfix [cs3org/reva#3815](https://github.com/cs3org/reva/pull/3815): Bump micro redis store
  *   Bugfix [cs3org/reva#3596](https://github.com/cs3org/reva/pull/3596): Cache CreateHome calls
  *   Bugfix [cs3org/reva#3823](https://github.com/cs3org/reva/pull/3823): Deny correctly in decomposedfs
  *   Bugfix [cs3org/reva#3826](https://github.com/cs3org/reva/pull/3826): Add by group index to decomposedfs
  *   Bugfix [cs3org/reva#3618](https://github.com/cs3org/reva/pull/3618): Drain body on failed put
  *   Bugfix [cs3org/reva#3685](https://github.com/cs3org/reva/pull/3685): Send fileid on copy
  *   Bugfix [cs3org/reva#3688](https://github.com/cs3org/reva/pull/3688): Return 425 on GET
  *   Bugfix [cs3org/reva#3755](https://github.com/cs3org/reva/pull/3755): Fix app provider language validation
  *   Bugfix [cs3org/reva#3800](https://github.com/cs3org/reva/pull/3800): Fix building for freebsd
  *   Bugfix [cs3org/reva#3700](https://github.com/cs3org/reva/pull/3700): Fix caching
  *   Bugfix [cs3org/reva#3535](https://github.com/cs3org/reva/pull/3535): Fix ceph driver storage fs implementation
  *   Bugfix [cs3org/reva#3764](https://github.com/cs3org/reva/pull/3764): Fix missing CORS config in ocdav service
  *   Bugfix [cs3org/reva#3710](https://github.com/cs3org/reva/pull/3710): Fix error when try to delete space without permission
  *   Bugfix [cs3org/reva#3822](https://github.com/cs3org/reva/pull/3822): Fix deleting spaces
  *   Bugfix [cs3org/reva#3718](https://github.com/cs3org/reva/pull/3718): Fix revad-eos docker image which was failing to build
  *   Bugfix [cs3org/reva#3559](https://github.com/cs3org/reva/pull/3559): Fix build on freebsd
  *   Bugfix [cs3org/reva#3696](https://github.com/cs3org/reva/pull/3696): Fix ldap filters when checking for enabled users
  *   Bugfix [cs3org/reva#3767](https://github.com/cs3org/reva/pull/3767): Decode binary UUID when looking up a users group memberships
  *   Bugfix [cs3org/reva#3741](https://github.com/cs3org/reva/pull/3741): Fix listing shares to multiple groups
  *   Bugfix [cs3org/reva#3834](https://github.com/cs3org/reva/pull/3834): Return correct error during MKCOL
  *   Bugfix [cs3org/reva#3841](https://github.com/cs3org/reva/pull/3841): Fix nil pointer and improve logging
  *   Bugfix [cs3org/reva#3831](https://github.com/cs3org/reva/pull/3831): Ignore 'null' mtime on tus upload
  *   Bugfix [cs3org/reva#3758](https://github.com/cs3org/reva/pull/3758): Fix public links with enforced password
  *   Bugfix [cs3org/reva#3814](https://github.com/cs3org/reva/pull/3814): Fix stat cache access
  *   Bugfix [cs3org/reva#3650](https://github.com/cs3org/reva/pull/3650): FreeBSD xattr support
  *   Bugfix [cs3org/reva#3827](https://github.com/cs3org/reva/pull/3827): Initialize user cache for decomposedfs
  *   Bugfix [cs3org/reva#3818](https://github.com/cs3org/reva/pull/3818): Invalidate cache when deleting space
  *   Bugfix [cs3org/reva#3812](https://github.com/cs3org/reva/pull/3812): Filemetadata Cache now deletes keys without listing them first
  *   Bugfix [cs3org/reva#3817](https://github.com/cs3org/reva/pull/3817): Pipeline cache deletes
  *   Bugfix [cs3org/reva#3711](https://github.com/cs3org/reva/pull/3711): Replace ini metadata backend by messagepack backend
  *   Bugfix [cs3org/reva#3828](https://github.com/cs3org/reva/pull/3828): Send quota when listing spaces in decomposedfs
  *   Bugfix [cs3org/reva#3681](https://github.com/cs3org/reva/pull/3681): Fix etag of "empty" shares jail
  *   Bugfix [cs3org/reva#3748](https://github.com/cs3org/reva/pull/3748): Prevent service from panicking
  *   Bugfix [cs3org/reva#3816](https://github.com/cs3org/reva/pull/3816): Write Metadata once
  *   Change [cs3org/reva#3641](https://github.com/cs3org/reva/pull/3641): Hide file versions for share receivers
  *   Change [cs3org/reva#3820](https://github.com/cs3org/reva/pull/3820): Streamline stores
  *   Enhancement [cs3org/reva#3732](https://github.com/cs3org/reva/pull/3732): Make method for detecting the metadata backend public
  *   Enhancement [cs3org/reva#3789](https://github.com/cs3org/reva/pull/3789): Add capabilities indicating if user attributes are read-only
  *   Enhancement [cs3org/reva#3792](https://github.com/cs3org/reva/pull/3792): Add a prometheus gauge to keep track of active uploads and downloads
  *   Enhancement [cs3org/reva#3637](https://github.com/cs3org/reva/pull/3637): Add an ID to each events
  *   Enhancement [cs3org/reva#3704](https://github.com/cs3org/reva/pull/3704): Add more information to events
  *   Enhancement [cs3org/reva#3744](https://github.com/cs3org/reva/pull/3744): Add LDAP user type attribute
  *   Enhancement [cs3org/reva#3806](https://github.com/cs3org/reva/pull/3806): Decomposedfs now supports filtering spaces by owner
  *   Enhancement [cs3org/reva#3730](https://github.com/cs3org/reva/pull/3730): Antivirus
  *   Enhancement [cs3org/reva#3531](https://github.com/cs3org/reva/pull/3531): Async Postprocessing
  *   Enhancement [cs3org/reva#3571](https://github.com/cs3org/reva/pull/3571): Async Upload Improvements
  *   Enhancement [cs3org/reva#3801](https://github.com/cs3org/reva/pull/3801): Cache node ids
  *   Enhancement [cs3org/reva#3690](https://github.com/cs3org/reva/pull/3690): Check set project space quota permission
  *   Enhancement [cs3org/reva#3686](https://github.com/cs3org/reva/pull/3686): User disabling functionality
  *   Enhancement [cs3org/reva#3505](https://github.com/cs3org/reva/pull/3505): Fix eosgrpc package
  *   Enhancement [cs3org/reva#3575](https://github.com/cs3org/reva/pull/3575): Fix skip group grant index cleanup
  *   Enhancement [cs3org/reva#3564](https://github.com/cs3org/reva/pull/3564): Fix tag pkg
  *   Enhancement [cs3org/reva#3756](https://github.com/cs3org/reva/pull/3756): Prepare for GDPR export
  *   Enhancement [cs3org/reva#3612](https://github.com/cs3org/reva/pull/3612): Group feature changed event added
  *   Enhancement [cs3org/reva#3729](https://github.com/cs3org/reva/pull/3729): Improve decomposedfs performance, esp. with network fs/cache
  *   Enhancement [cs3org/reva#3697](https://github.com/cs3org/reva/pull/3697): Improve the ini file metadata backend
  *   Enhancement [cs3org/reva#3819](https://github.com/cs3org/reva/pull/3819): Allow creating internal links without permission
  *   Enhancement [cs3org/reva#3740](https://github.com/cs3org/reva/pull/3740): Limit concurrency in decomposedfs
  *   Enhancement [cs3org/reva#3569](https://github.com/cs3org/reva/pull/3569): Always list shares jail when listing spaces
  *   Enhancement [cs3org/reva#3788](https://github.com/cs3org/reva/pull/3788): Make resharing configurable
  *   Enhancement [cs3org/reva#3674](https://github.com/cs3org/reva/pull/3674): Introduce ini file based metadata backend
  *   Enhancement [cs3org/reva#3728](https://github.com/cs3org/reva/pull/3728): Automatically migrate file metadata from xattrs to messagepack
  *   Enhancement [cs3org/reva#3807](https://github.com/cs3org/reva/pull/3807): Name Validation
  *   Enhancement [cs3org/reva#3574](https://github.com/cs3org/reva/pull/3574): Opaque space group
  *   Enhancement [cs3org/reva#3598](https://github.com/cs3org/reva/pull/3598): Pass estream to Storage Providers
  *   Enhancement [cs3org/reva#3763](https://github.com/cs3org/reva/pull/3763): Add a capability for personal data export
  *   Enhancement [cs3org/reva#3577](https://github.com/cs3org/reva/pull/3577): Prepare for SSE
  *   Enhancement [cs3org/reva#3731](https://github.com/cs3org/reva/pull/3731): Add config option to enforce passwords on public links
  *   Enhancement [cs3org/reva#3693](https://github.com/cs3org/reva/pull/3693): Enforce the PublicLink.Write permission
  *   Enhancement [cs3org/reva#3497](https://github.com/cs3org/reva/pull/3497): Introduce owncloud 10 publiclink manager
  *   Enhancement [cs3org/reva#3714](https://github.com/cs3org/reva/pull/3714): Add global max quota option and quota for CreateHome
  *   Enhancement [cs3org/reva#3759](https://github.com/cs3org/reva/pull/3759): Set correct share type when listing shares
  *   Enhancement [cs3org/reva#3594](https://github.com/cs3org/reva/pull/3594): Add expiration to user and group shares
  *   Enhancement [cs3org/reva#3580](https://github.com/cs3org/reva/pull/3580): Share expired event
  *   Enhancement [cs3org/reva#3620](https://github.com/cs3org/reva/pull/3620): Allow a new ShareType `SpaceMembershipGroup`
  *   Enhancement [cs3org/reva#3609](https://github.com/cs3org/reva/pull/3609): Space Management Permissions
  *   Enhancement [cs3org/reva#3655](https://github.com/cs3org/reva/pull/3655): Add expiration date to space memberships
  *   Enhancement [cs3org/reva#3697](https://github.com/cs3org/reva/pull/3697): Add support for redis sentinel caches
  *   Enhancement [cs3org/reva#3552](https://github.com/cs3org/reva/pull/3552): Suppress tusd logs
  *   Enhancement [cs3org/reva#3555](https://github.com/cs3org/reva/pull/3555): Tags
  *   Enhancement [cs3org/reva#3785](https://github.com/cs3org/reva/pull/3785): Increase unit test coverage in the ocdav service
  *   Enhancement [cs3org/reva#3739](https://github.com/cs3org/reva/pull/3739): Try to rename uploaded files to their final position
  *   Enhancement [cs3org/reva#3610](https://github.com/cs3org/reva/pull/3610): Walk and log chi routes

   https://github.com/owncloud/ocis/pull/6448
   https://github.com/owncloud/ocis/pull/6447
   https://github.com/owncloud/ocis/pull/6381
   https://github.com/owncloud/ocis/pull/6305
   https://github.com/owncloud/ocis/pull/6339
   https://github.com/owncloud/ocis/pull/6205
   https://github.com/owncloud/ocis/pull/6186

# Changelog for [2.0.0] (2022-11-30)

The following sections list the changes for 2.0.0.

[2.0.0]: https://github.com/owncloud/ocis/compare/v1.20.0...v2.0.0

## Summary

* Bugfix - Substring search for sharees: [#547](https://github.com/owncloud/ocis/issues/547)
* Bugfix - Return proper errors when ocs/cloud/users is using the cs3 backend: [#3483](https://github.com/owncloud/ocis/issues/3483)
* Bugfix - Thumbnails for `/dav/xxx?preview=1` requests: [#3567](https://github.com/owncloud/ocis/pull/3567)
* Bugfix - URL encode the webdav url in the graph API: [#3597](https://github.com/owncloud/ocis/pull/3597)
* Bugfix - Idp: Check if CA certificate if present: [#3623](https://github.com/owncloud/ocis/issues/3623)
* Bugfix - Fix DN parsing issues and sizelimit handling in libregraph/idm: [#3631](https://github.com/owncloud/ocis/issues/3631)
* Bugfix - Fix the webdav URL of drive roots: [#3706](https://github.com/owncloud/ocis/issues/3706)
* Bugfix - Check permissions when deleting Space: [#3709](https://github.com/owncloud/ocis/pull/3709)
* Bugfix - Remove runtime kill and run commands: [#3740](https://github.com/owncloud/ocis/pull/3740)
* Bugfix - Make IDP secrets configurable via environment variables: [#3744](https://github.com/owncloud/ocis/pull/3744)
* Bugfix - Store user passwords hashed in idm: [#3778](https://github.com/owncloud/ocis/issues/3778)
* Bugfix - Fix version number in status page: [#3788](https://github.com/owncloud/ocis/issues/3788)
* Bugfix - Fix Thumbnails for IDs without a trailing path: [#3791](https://github.com/owncloud/ocis/pull/3791)
* Bugfix - Fix the `ocis search` command: [#3796](https://github.com/owncloud/ocis/pull/3796)
* Bugfix - Remove unused transfer secret from app provider: [#3798](https://github.com/owncloud/ocis/pull/3798)
* Bugfix - Fix the idm and settings extensions' admin user id configuration option: [#3799](https://github.com/owncloud/ocis/pull/3799)
* Bugfix - Rename search env variable for the grpc server address: [#3800](https://github.com/owncloud/ocis/pull/3800)
* Bugfix - Fix multiple storage-users env variables: [#3802](https://github.com/owncloud/ocis/pull/3802)
* Bugfix - Save Katherine: [#3823](https://github.com/owncloud/ocis/issues/3823)
* Bugfix - Enable debug server by default: [#3827](https://github.com/owncloud/ocis/pull/3827)
* Bugfix - Remove legacy accounts proxy routes: [#3831](https://github.com/owncloud/ocis/pull/3831)
* Bugfix - Set default name for public link via capabilities: [#3834](https://github.com/owncloud/ocis/pull/3834)
* Bugfix - Fix search index getting out of sync: [#3851](https://github.com/owncloud/ocis/pull/3851)
* Bugfix - Inconsistency env var naming for LDAP filter configuration: [#3890](https://github.com/owncloud/ocis/issues/3890)
* Bugfix - Allow empty environment variables: [#3892](https://github.com/owncloud/ocis/pull/3892)
* Bugfix - Fix user autoprovisioning: [#3893](https://github.com/owncloud/ocis/issues/3893)
* Bugfix - Fix LDAP insecure options: [#3897](https://github.com/owncloud/ocis/pull/3897)
* Bugfix - Rework default role provisioning: [#3900](https://github.com/owncloud/ocis/issues/3900)
* Bugfix - Fix configuration validation for extensions' server commands: [#3911](https://github.com/owncloud/ocis/pull/3911)
* Bugfix - Fix graph endpoint: [#3925](https://github.com/owncloud/ocis/issues/3925)
* Bugfix - Fix version info: [#3953](https://github.com/owncloud/ocis/pull/3953)
* Bugfix - Remove unused OCS storage configuration: [#3955](https://github.com/owncloud/ocis/pull/3955)
* Bugfix - Make ocdav service behave properly: [#3957](https://github.com/owncloud/ocis/pull/3957)
* Bugfix - Make IDP only wait for certs when using LDAP: [#3965](https://github.com/owncloud/ocis/pull/3965)
* Bugfix - Remove unused configuration options: [#3973](https://github.com/owncloud/ocis/pull/3973)
* Bugfix - CSP rules for silent token refresh in iframe: [#4031](https://github.com/owncloud/ocis/pull/4031)
* Bugfix - Logging in on the wrong account when an email address is not unique: [#4039](https://github.com/owncloud/ocis/issues/4039)
* Bugfix - Remove static ocs user backend config: [#4077](https://github.com/owncloud/ocis/pull/4077)
* Bugfix - Fix make sensitive config values in the proxy's debug server: [#4086](https://github.com/owncloud/ocis/pull/4086)
* Bugfix - Fix startup error logging: [#4093](https://github.com/owncloud/ocis/pull/4093)
* Bugfix - Polish search: [#4094](https://github.com/owncloud/ocis/pull/4094)
* Bugfix - Fix logging levels: [#4102](https://github.com/owncloud/ocis/pull/4102)
* Bugfix - Escape DN attribute value: [#4117](https://github.com/owncloud/ocis/pull/4117)
* Bugfix - Fix `OCIS_RUN_SERVICES`: [#4133](https://github.com/owncloud/ocis/pull/4133)
* Bugfix - Space Creators can hand over spaces: [#4244](https://github.com/owncloud/ocis/pull/4244)
* Bugfix - Fix handling of invalid LDAP users and groups: [#4274](https://github.com/owncloud/ocis/issues/4274)
* Bugfix - Fix search in received shares: [#4308](https://github.com/owncloud/ocis/issues/4308)
* Bugfix - Fix unrestricted quota on the graphAPI: [#4363](https://github.com/owncloud/ocis/pull/4363)
* Bugfix - Autocreate IDP private key also if file exists but is empty: [#4394](https://github.com/owncloud/ocis/pull/4394)
* Bugfix - Show help for some commands when unconfigured: [#4405](https://github.com/owncloud/ocis/pull/4405)
* Bugfix - Rename extensions to services (leftover occurrences): [#4407](https://github.com/owncloud/ocis/pull/4407)
* Bugfix - Fix configuration of mimetypes for the app registry: [#4411](https://github.com/owncloud/ocis/pull/4411)
* Bugfix - Disable default expiration for public links: [#4445](https://github.com/owncloud/ocis/issues/4445)
* Bugfix - Fix permissions in REPORT: [#4520](https://github.com/owncloud/ocis/pull/4520)
* Bugfix - Render webdav permissions as string in search report: [#4575](https://github.com/owncloud/ocis/issues/4575)
* Bugfix - Graph service now forwards trace context: [#4582](https://github.com/owncloud/ocis/pull/4582)
* Bugfix - Fix sharing jsoncs3 driver options: [#4593](https://github.com/owncloud/ocis/pull/4593)
* Bugfix - Fix the OIDC provider cache: [#4600](https://github.com/owncloud/ocis/pull/4600)
* Bugfix - Change the default value for PROXY_OIDC_INSECURE to false: [#4601](https://github.com/owncloud/ocis/pull/4601)
* Bugfix - Fix authentication for autoprovisioned users: [#4616](https://github.com/owncloud/ocis/issues/4616)
* Bugfix - Fix wopi access to public shares: [#4631](https://github.com/owncloud/ocis/pull/4631)
* Bugfix - Fix unfindable entities from shares/publicshares: [#4651](https://github.com/owncloud/ocis/pull/4651)
* Bugfix - Fix notifications service settings: [#4652](https://github.com/owncloud/ocis/pull/4652)
* Bugfix - Bring back the settings UI in Web: [#4691](https://github.com/owncloud/ocis/pull/4691)
* Bugfix - Don't run auth-bearer service by default: [#4692](https://github.com/owncloud/ocis/issues/4692)
* Bugfix - Mail notifications for group shares: [#4714](https://github.com/owncloud/ocis/pull/4714)
* Bugfix - Make tokeninfo endpoint unprotected: [#4715](https://github.com/owncloud/ocis/pull/4715)
* Bugfix - Fix cache stat table config: [#4732](https://github.com/owncloud/ocis/pull/4732)
* Bugfix - Trigger a rescan of spaces in the search index when items have changed: [#4777](https://github.com/owncloud/ocis/pull/4777)
* Bugfix - Disable cache for selected static web assets: [#4809](https://github.com/owncloud/ocis/pull/4809)
* Bugfix - Remove the storage-users event configuration: [#4825](https://github.com/owncloud/ocis/pull/4825)
* Bugfix - Fix the shareroot path in REPORT responses: [#4859](https://github.com/owncloud/ocis/pull/4859)
* Bugfix - Disable federation capabilities: [#4864](https://github.com/owncloud/ocis/pull/4864)
* Bugfix - Fix permission check in settings service: [#4890](https://github.com/owncloud/ocis/pull/4890)
* Bugfix - Fix CORS in frontend service: [#4948](https://github.com/owncloud/ocis/pull/4948)
* Bugfix - Fix notifications Web UI url: [#4998](https://github.com/owncloud/ocis/pull/4998)
* Bugfix - Do not reindex a space twice at the same time: [#5001](https://github.com/owncloud/ocis/pull/5001)
* Bugfix - Find spaces by their name: [#5044](https://github.com/owncloud/ocis/pull/5044)
* Bugfix - Initial role assignment with external IDM: [#5045](https://github.com/owncloud/ocis/issues/5045)
* Bugfix - Lower IDP token lifespans: [#5077](https://github.com/owncloud/ocis/pull/5077)
* Bugfix - Adjust cache related configuration options: [#5087](https://github.com/owncloud/ocis/pull/5087)
* Bugfix - Make storage users mount ids unique by default: [#5091](https://github.com/owncloud/ocis/pull/5091)
* Bugfix - Update reva to version 2.12.0: [#5092](https://github.com/owncloud/ocis/pull/5092)
* Bugfix - Decomposedfs increase filelock duration factor: [#5130](https://github.com/owncloud/ocis/pull/5130)
* Bugfix - Translations on login page: [#7550](https://github.com/owncloud/web/issues/7550)
* Bugfix - Fix search report: [#7557](https://github.com/owncloud/web/issues/7557)
* Bugfix - Fix unused config option `GRAPH_SPACES_INSECURE`: [#55555](https://github.com/owncloud/ocis/pull/55555)
* Change - Switched default configuration to use libregraph/idm: [#3331](https://github.com/owncloud/ocis/pull/3331)
* Change - Introduce `ocis init` and remove all default secrets: [#3551](https://github.com/owncloud/ocis/pull/3551)
* Change - Load configuration files just from one directory: [#3587](https://github.com/owncloud/ocis/pull/3587)
* Change - Reduce drives in graph /me/drives API: [#3629](https://github.com/owncloud/ocis/pull/3629)
* Change - Reduce permissions on docker image predeclared volumes: [#3641](https://github.com/owncloud/ocis/pull/3641)
* Change - Use new space ID util functions: [#3648](https://github.com/owncloud/ocis/pull/3648)
* Change - Rename MetadataUserID: [#3671](https://github.com/owncloud/ocis/pull/3671)
* Change - Split MachineAuth from SystemUser: [#3672](https://github.com/owncloud/ocis/pull/3672)
* Change - Rename serviceUser to systemUser: [#3673](https://github.com/owncloud/ocis/pull/3673)
* Change - Update ocis packages and imports to V2: [#3678](https://github.com/owncloud/ocis/pull/3678)
* Change - The `glauth` and `accounts` services are removed: [#3685](https://github.com/owncloud/ocis/pull/3685)
* Change - Prevent access to disabled space: [#3779](https://github.com/owncloud/ocis/pull/3779)
* Change - Rename "uploads purge" command to "uploads clean": [#4403](https://github.com/owncloud/ocis/pull/4403)
* Change - Enable private links by default: [#4599](https://github.com/owncloud/ocis/pull/4599/)
* Change - Use the spaceID on the cs3 resource: [#4748](https://github.com/owncloud/ocis/pull/4748)
* Change - Build service frontends with pnpm instead of yarn: [#4878](https://github.com/owncloud/ocis/pull/4878)
* Enhancement - Disable the color logging in docker compose examples: [#871](https://github.com/owncloud/ocis/issues/871)
* Enhancement - Product field in OCS version: [#2918](https://github.com/owncloud/ocis/pull/2918)
* Enhancement - Add /me/changePassword endpoint to GraphAPI: [#3063](https://github.com/owncloud/ocis/issues/3063)
* Enhancement - Update IdP UI: [#3493](https://github.com/owncloud/ocis/issues/3493)
* Enhancement - Update reva to v2.3.1: [#3552](https://github.com/owncloud/ocis/pull/3552)
* Enhancement - Update linkshare capabilities: [#3579](https://github.com/owncloud/ocis/pull/3579)
* Enhancement - Wrap metadata storage with dedicated reva gateway: [#3602](https://github.com/owncloud/ocis/pull/3602)
* Enhancement - Align service naming: [#3606](https://github.com/owncloud/ocis/pull/3606)
* Enhancement - Added `share_jail` and `projects` feature flags in spaces capability: [#3626](https://github.com/owncloud/ocis/pull/3626)
* Enhancement - Add initial version of the search extensions: [#3635](https://github.com/owncloud/ocis/pull/3635)
* Enhancement - Don't setup demo role assignments on default: [#3661](https://github.com/owncloud/ocis/issues/3661)
* Enhancement - Restrict admins from self-removal: [#3713](https://github.com/owncloud/ocis/issues/3713)
* Enhancement - Update reva to version 2.4.1: [#3746](https://github.com/owncloud/ocis/pull/3746)
* Enhancement - Add description tags to the thumbnails config structs: [#3752](https://github.com/owncloud/ocis/pull/3752)
* Enhancement - Add acting user to the audit log: [#3753](https://github.com/owncloud/ocis/issues/3753)
* Enhancement - Add descriptions to webdav configuration: [#3755](https://github.com/owncloud/ocis/pull/3755)
* Enhancement - Add descriptions for graph-explorer config: [#3759](https://github.com/owncloud/ocis/pull/3759)
* Enhancement - Add config option to provide TLS certificate: [#3818](https://github.com/owncloud/ocis/issues/3818)
* Enhancement - Introduce service registry cache: [#3833](https://github.com/owncloud/ocis/pull/3833)
* Enhancement - Improve validation of OIDC access tokens: [#3841](https://github.com/owncloud/ocis/issues/3841)
* Enhancement - Reintroduce user autoprovisioning in proxy: [#3860](https://github.com/owncloud/ocis/pull/3860)
* Enhancement - Allow resharing: [#3904](https://github.com/owncloud/ocis/pull/3904)
* Enhancement - Generate signing key and encryption secret: [#3909](https://github.com/owncloud/ocis/issues/3909)
* Enhancement - Add deprecation annotation: [#3917](https://github.com/owncloud/ocis/issues/3917)
* Enhancement - Update reva to version 2.5.1: [#3932](https://github.com/owncloud/ocis/pull/3932)
* Enhancement - Add audit events for created containers: [#3941](https://github.com/owncloud/ocis/pull/3941)
* Enhancement - Update reva: [#3944](https://github.com/owncloud/ocis/pull/3944)
* Enhancement - Make thumbnails service log less noisy: [#3959](https://github.com/owncloud/ocis/pull/3959)
* Enhancement - Refactor extensions to services: [#3980](https://github.com/owncloud/ocis/pull/3980)
* Enhancement - Add capability for alias links: [#3983](https://github.com/owncloud/ocis/issues/3983)
* Enhancement - New migrate command for migrating shares and public shares: [#3987](https://github.com/owncloud/ocis/pull/3987)
* Enhancement - Update ownCloud Web to v5.7.0-rc.1: [#4005](https://github.com/owncloud/ocis/pull/4005)
* Enhancement - Add FRONTEND_ENABLE_RESHARING env variable: [#4023](https://github.com/owncloud/ocis/pull/4023)
* Enhancement - Add drives field to users endpoint: [#4072](https://github.com/owncloud/ocis/pull/4072)
* Enhancement - Added command to reset administrator password: [#4084](https://github.com/owncloud/ocis/issues/4084)
* Enhancement - Update reva to version 2.7.2: [#4115](https://github.com/owncloud/ocis/pull/4115)
* Enhancement - Search service at the old webdav endpoint: [#4118](https://github.com/owncloud/ocis/pull/4118)
* Enhancement - Update ownCloud Web to v5.7.0-rc.4: [#4140](https://github.com/owncloud/ocis/pull/4140)
* Enhancement - Add number of total matches to the search result: [#4189](https://github.com/owncloud/ocis/issues/4189)
* Enhancement - Introduce "delete-all-spaces" permission: [#4196](https://github.com/owncloud/ocis/issues/4196)
* Enhancement - Improve error log for "could not get user by claim" error: [#4227](https://github.com/owncloud/ocis/pull/4227)
* Enhancement - Allow providing list of services NOT to start: [#4254](https://github.com/owncloud/ocis/pull/4254)
* Enhancement - Introduce insecure flag for smtp email notifications: [#4279](https://github.com/owncloud/ocis/pull/4279)
* Enhancement - Update reva to v2.7.4: [#4294](https://github.com/owncloud/ocis/pull/4294)
* Enhancement - Update ownCloud Web to v5.7.0-rc.8: [#4314](https://github.com/owncloud/ocis/pull/4314)
* Enhancement - OCS get share now also handle received shares: [#4322](https://github.com/owncloud/ocis/issues/4322)
* Enhancement - Fix behavior for foobar (in present tense): [#4346](https://github.com/owncloud/ocis/pull/4346)
* Enhancement - Use storageID when requesting special items: [#4356](https://github.com/owncloud/ocis/pull/4356)
* Enhancement - Expand personal drive on the graph user: [#4357](https://github.com/owncloud/ocis/pull/4357)
* Enhancement - Rewrite of the request authentication middleware: [#4374](https://github.com/owncloud/ocis/pull/4374)
* Enhancement - Add /app/open-with-web endpoint: [#4376](https://github.com/owncloud/ocis/pull/4376)
* Enhancement - Added language option to the app provider: [#4399](https://github.com/owncloud/ocis/pull/4399)
* Enhancement - Refactor the proxy service: [#4401](https://github.com/owncloud/ocis/issues/4401)
* Enhancement - Add previewFileMimeTypes to web default config: [#4414](https://github.com/owncloud/ocis/pull/4414)
* Enhancement - Update ownCloud Web to v5.7.0-rc.10: [#4439](https://github.com/owncloud/ocis/pull/4439)
* Enhancement - Add configuration options for mail authentication and encryption: [#4443](https://github.com/owncloud/ocis/pull/4443)
* Enhancement - Update reva to v2.8.0: [#4444](https://github.com/owncloud/ocis/pull/4444)
* Enhancement - Add missing unprotected paths: [#4454](https://github.com/owncloud/ocis/pull/4454)
* Enhancement - Automatically orientate photos when generating thumbnails: [#4477](https://github.com/owncloud/ocis/issues/4477)
* Enhancement - Improve login screen design: [#4500](https://github.com/owncloud/ocis/pull/4500)
* Enhancement - Update ownCloud Web to v5.7.0: [#4508](https://github.com/owncloud/ocis/pull/4508)
* Enhancement - Update Reva to version 2.10.0: [#4522](https://github.com/owncloud/ocis/pull/4522)
* Enhancement - Add Email templating: [#4564](https://github.com/owncloud/ocis/pull/4564)
* Enhancement - Allow to configure applications in Web: [#4578](https://github.com/owncloud/ocis/pull/4578)
* Enhancement - Add webURL to space root: [#4588](https://github.com/owncloud/ocis/pull/4588)
* Enhancement - Update reva to version 2.11.0: [#4588](https://github.com/owncloud/ocis/pull/4588)
* Enhancement - Allow to configuring the reva cache store: [#4627](https://github.com/owncloud/ocis/pull/4627)
* Enhancement - Add thumbnails support for tiff and bmp files: [#4634](https://github.com/owncloud/ocis/pull/4634)
* Enhancement - Add support for REPORT requests to /dav/spaces URLs: [#4661](https://github.com/owncloud/ocis/pull/4661)
* Enhancement - Make it possible to configure a WOPI folderurl: [#4716](https://github.com/owncloud/ocis/pull/4716)
* Enhancement - Add curl to the oCIS OCI image: [#4751](https://github.com/owncloud/ocis/pull/4751)
* Enhancement - Report parent id: [#4757](https://github.com/owncloud/ocis/pull/4757)
* Enhancement - Secure the nats connection with TLS: [#4781](https://github.com/owncloud/ocis/pull/4781)
* Enhancement - Allow to setup TLS for grpc services: [#4798](https://github.com/owncloud/ocis/pull/4798)
* Enhancement - We added e-mail subject templating: [#4799](https://github.com/owncloud/ocis/pull/4799)
* Enhancement - Logging improvements: [#4815](https://github.com/owncloud/ocis/pull/4815)
* Enhancement - Prohibit users from setting or listing other user's values: [#4897](https://github.com/owncloud/ocis/pull/4897)
* Enhancement - Deny access to resources: [#4903](https://github.com/owncloud/ocis/pull/4903)
* Enhancement - Validate space names: [#4955](https://github.com/owncloud/ocis/pull/4955)
* Enhancement - Configurable max lock cycles: [#4965](https://github.com/owncloud/ocis/pull/4965)
* Enhancement - Rename AUTH_BASIC_AUTH_PROVIDER envvar: [#4966](https://github.com/owncloud/ocis/pull/4966)
* Enhancement - Default to tls 1.2: [#4969](https://github.com/owncloud/ocis/pull/4969)
* Enhancement - Add the "hidden" state to the search index: [#5018](https://github.com/owncloud/ocis/pull/5018)
* Enhancement - Remove windows from ci & release makefile: [#5026](https://github.com/owncloud/ocis/pull/5026)
* Enhancement - Add tracing to search: [#5113](https://github.com/owncloud/ocis/pull/5113)
* Enhancement - Update ownCloud Web to v6.0.0: [#5153](https://github.com/owncloud/ocis/pull/5153)
* Enhancement - Add capability for public link single file edit: [#6787](https://github.com/owncloud/web/pull/6787)
* Enhancement - Update ownCloud Web to v5.5.0-rc.8: [#6854](https://github.com/owncloud/web/pull/6854)
* Enhancement - Update ownCloud Web to v5.5.0-rc.9: [#6854](https://github.com/owncloud/web/pull/6854)
* Enhancement - Update ownCloud Web to v5.5.0-rc.6: [#6854](https://github.com/owncloud/web/pull/6854)
* Enhancement - Optional events in graph service: [#55555](https://github.com/owncloud/ocis/pull/55555)

## Details

* Bugfix - Substring search for sharees: [#547](https://github.com/owncloud/ocis/issues/547)

   We fixed searching for sharees to be no longer case-sensitive. With this we
   introduced two new settings for the users and groups services:
   "group_substring_filter_type" for the group services and
   "user_substring_filter_type" for the users service. They allow to set the type
   of LDAP filter that is used for substring user searches. Possible values are:
   "initial", "final" and "any" to do either prefix, suffix or full substring
   searches. Both settings default to "initial".

   Also a new option "search_min_length" was added for the "frontend" service. It
   allows to configure the minimum number of characters to enter before a search
   for Sharees is started. This setting is e.g. evaluated by the web ui via the
   capabilities endpoint.

   https://github.com/owncloud/ocis/issues/547

* Bugfix - Return proper errors when ocs/cloud/users is using the cs3 backend: [#3483](https://github.com/owncloud/ocis/issues/3483)

   The ocs API was just exiting with a fatal error on any update request, when
   configured for the cs3 backend. Now it returns a proper error.

   https://github.com/owncloud/ocis/issues/3483

* Bugfix - Thumbnails for `/dav/xxx?preview=1` requests: [#3567](https://github.com/owncloud/ocis/pull/3567)

   We've added the thumbnail rendering for `/dav/xxx?preview=1`,
   `/remote.php/webdav/{relative path}?preview=1` and `/webdav/{relative
   path}?preview=1` requests, which was previously not supported because of missing
   routes. It now returns the same thumbnails as for
   `/remote.php/dav/xxx?preview=1`.

   https://github.com/owncloud/ocis/pull/3567

* Bugfix - URL encode the webdav url in the graph API: [#3597](https://github.com/owncloud/ocis/pull/3597)

   Fixed the webdav URL in the drives responses. Without encoding the URL could be
   broken by files with spaces in the file name.

   https://github.com/owncloud/ocis/issues/3538
   https://github.com/owncloud/ocis/pull/3597

* Bugfix - Idp: Check if CA certificate if present: [#3623](https://github.com/owncloud/ocis/issues/3623)

   Upon first start with the default configuration the idm service creates a server
   certificate, that might not be finished before the idp service is starting. Add
   a check to idp similar to what the user, group, and auth-providers implement.

   https://github.com/owncloud/ocis/issues/3623

* Bugfix - Fix DN parsing issues and sizelimit handling in libregraph/idm: [#3631](https://github.com/owncloud/ocis/issues/3631)

   We fixed a couple on issues in libregraph/idm related to correctly parsing LDAP
   DNs for usernames contain characters that require escaping.

   Also libregraph/idm was not properly returning "Size limit exceeded" errors when
   the result set exceeded the requested size.

   https://github.com/owncloud/ocis/issues/3631
   https://github.com/owncloud/ocis/issues/4039
   https://github.com/owncloud/ocis/issues/4078

* Bugfix - Fix the webdav URL of drive roots: [#3706](https://github.com/owncloud/ocis/issues/3706)

   Fixed the webdav URL of drive roots in the graph API.

   https://github.com/owncloud/ocis/issues/3706
   https://github.com/owncloud/ocis/pull/3916

* Bugfix - Check permissions when deleting Space: [#3709](https://github.com/owncloud/ocis/pull/3709)

   Check for manager permissions when deleting spaces. Do not allow deleting spaces
   via dav service

   https://github.com/owncloud/ocis/pull/3709

* Bugfix - Remove runtime kill and run commands: [#3740](https://github.com/owncloud/ocis/pull/3740)

   We've removed the kill and run commands from the oCIS runtime. If these dynamic
   capabilities are needed, one should switch to a full fledged supervisor and
   start oCIS as individual services.

   If one wants to start a only a subset of services, this is still possible by
   setting OCIS_RUN_EXTENSIONS.

   https://github.com/owncloud/ocis/pull/3740

* Bugfix - Make IDP secrets configurable via environment variables: [#3744](https://github.com/owncloud/ocis/pull/3744)

   We've fixed the configuration options of the IDP to make the IDP secrets again
   configurable via environment variables.

   https://github.com/owncloud/ocis/pull/3744

* Bugfix - Store user passwords hashed in idm: [#3778](https://github.com/owncloud/ocis/issues/3778)

   Support for hashing user passwords was added to libregraph/idm. The graph API
   will now set userpasswords using the LDAP Modify Extended Operation (RFC3062).
   In the default configuration passwords will be hashed using the argon2id
   algorithm.

   https://github.com/owncloud/ocis/issues/3778
   https://github.com/owncloud/ocis/pull/4053

* Bugfix - Fix version number in status page: [#3788](https://github.com/owncloud/ocis/issues/3788)

   We needed to undo the version number changes on the status page to keep
   compatibility for legacy clients. We added a new field `productversion` for the
   actual version of the product.

   https://github.com/owncloud/ocis/issues/3788
   https://github.com/owncloud/ocis/pull/3805

* Bugfix - Fix Thumbnails for IDs without a trailing path: [#3791](https://github.com/owncloud/ocis/pull/3791)

   The routes in the chi router were not matching thumbnail requests without a
   trailing path.

   https://github.com/owncloud/ocis/pull/3791

* Bugfix - Fix the `ocis search` command: [#3796](https://github.com/owncloud/ocis/pull/3796)

   We've fixed the behavior for `ocis search`, which didn't show further help when
   not all secrets have been configured. It also was not possible to start the
   search service standalone from the oCIS binary without configuring all oCIS
   secrets, even they were not needed by the search service.

   https://github.com/owncloud/ocis/pull/3796

* Bugfix - Remove unused transfer secret from app provider: [#3798](https://github.com/owncloud/ocis/pull/3798)

   We've fixed the startup of the app provider by removing the startup dependency
   on a configured transfer secret, which was not used. This only happened if you
   start the app provider without runtime (eg. `ocis app-provider server`) and
   didn't have configured all oCIS secrets.

   https://github.com/owncloud/ocis/pull/3798

* Bugfix - Fix the idm and settings extensions' admin user id configuration option: [#3799](https://github.com/owncloud/ocis/pull/3799)

   We've fixed the admin user id configuration of the settings and idm extensions.
   The have previously only been configurable via the oCIS shared configuration and
   therefore have been undocumented for the extensions. This config option is now
   part of both extensions' configuration and can now also be used when the
   extensions are compiled standalone.

   https://github.com/owncloud/ocis/pull/3799

* Bugfix - Rename search env variable for the grpc server address: [#3800](https://github.com/owncloud/ocis/pull/3800)

   We've fixed the gprc server address configuration environment variable by
   renaming it from `ACCOUNTS_GRPC_ADDR` to `SEARCH_GRPC_ADDR`

   https://github.com/owncloud/ocis/pull/3800

* Bugfix - Fix multiple storage-users env variables: [#3802](https://github.com/owncloud/ocis/pull/3802)

   We've fixed multiple environment variable configuration options for the
   storage-users extension:

  * `STORAGE_USERS_GRPC_ADDR` was used to configure both the address of the http and grpc server. This resulted in a failing startup of the storage-users extension if this config option is set, because the service tries to double-bind the configured port (one time for each of the http and grpc server). You can now configure the grpc server's address with the environment variable `STORAGE_USERS_GRPC_ADDR` and the http server's address with the environment variable `STORAGE_USERS_HTTP_ADDR`
  * `STORAGE_USERS_S3NG_USERS_PROVIDER_ENDPOINT` was used to configure the permissions service endpoint for the S3NG driver and was therefore renamed to `STORAGE_USERS_S3NG_PERMISSIONS_ENDPOINT`
  * It's now possible to configure the permissions service endpoint for all  storage drivers with the environment variable `STORAGE_USERS_PERMISSION_ENDPOINT`, which was previously only used by the S3NG driver.

   https://github.com/owncloud/ocis/pull/3802

* Bugfix - Save Katherine: [#3823](https://github.com/owncloud/ocis/issues/3823)

   SpaceManager user katherine was removed with the demo user switch. Now she comes
   back

   https://github.com/owncloud/ocis/issues/3823
   https://github.com/owncloud/ocis/pull/3824

* Bugfix - Enable debug server by default: [#3827](https://github.com/owncloud/ocis/pull/3827)

   We've fixed the behavior for the audit, idm, nats and notifications extensions,
   that did not start their debug server by default.

   https://github.com/owncloud/ocis/pull/3827

* Bugfix - Remove legacy accounts proxy routes: [#3831](https://github.com/owncloud/ocis/pull/3831)

   We've removed the legacy accounts routes from the proxy default config. There
   were no longer used since the switch to IDM as the default user backend. Also
   accounts is no longer part of the oCIS binary and therefore should not be part
   of the proxy default route config.

   https://github.com/owncloud/ocis/pull/3831

* Bugfix - Set default name for public link via capabilities: [#3834](https://github.com/owncloud/ocis/pull/3834)

   We have now added a default name for public link shares which is communicated
   via the capabilities.

   https://github.com/owncloud/ocis/issues/1237
   https://github.com/owncloud/ocis/pull/3834

* Bugfix - Fix search index getting out of sync: [#3851](https://github.com/owncloud/ocis/pull/3851)

   We fixed a problem where the search index got out of sync with child elements of
   a parent containing special characters.

   https://github.com/owncloud/ocis/pull/3851

* Bugfix - Inconsistency env var naming for LDAP filter configuration: [#3890](https://github.com/owncloud/ocis/issues/3890)

   There was a naming inconsistency for the environment variables used to define
   LDAP filters for user and groups queries. Some services used `LDAP_USER_FILTER`
   while others used `LDAP_USERFILTER`. This is now changed to use
   `LDAP_USER_FILTER` and `LDAP_GROUP_FILTER`.

   Note: If your oCIS setup is using an LDAP configuration that has any of the
   `*_LDAP_USERFILTER` or `*_LDAP_GROUPFILTER` environment variables set, please
   update the configuration to use the new unified names `*_LDAP_USER_FILTER`
   respectively `*_LDAP_GROUP_FILTER` instead.

   https://github.com/owncloud/ocis/issues/3890

* Bugfix - Allow empty environment variables: [#3892](https://github.com/owncloud/ocis/pull/3892)

   We've fixed the behavior for empty environment variables, that previously would
   not have overwritten default values. Therefore it had the same effect like not
   setting the environment variable. We now check if the environment variable is
   set at all and if so, we also allow to override a default value with an empty
   value.

   https://github.com/owncloud/ocis/pull/3892

* Bugfix - Fix user autoprovisioning: [#3893](https://github.com/owncloud/ocis/issues/3893)

   We've fixed the autoprovsioning feature that was introduced in beta2. Due to a
   bug the role assignment of the privileged user that is used to create accounts
   wasn't propagated correctly to the `graph` service.

   https://github.com/owncloud/ocis/issues/3893

* Bugfix - Fix LDAP insecure options: [#3897](https://github.com/owncloud/ocis/pull/3897)

   We've fixed multiple LDAP insecure options:

  * The Graph LDAP insecure option default was set to `true` and now defaults to `false`. This is possible after #3888, since the Graph also now uses the LDAP CAcert by default.
  * The Graph LDAP insecure option was configurable by the environment variable `OCIS_INSECURE`, which was replaced by the dedicated `LDAP_INSECURE` variable. This variable is also used by all other services using LDAP.
  * The IDP insecure option for the user backend now also picks up configuration from `LDAP_INSECURE`.

   https://github.com/owncloud/ocis/pull/3897

* Bugfix - Rework default role provisioning: [#3900](https://github.com/owncloud/ocis/issues/3900)

   We fixed a race condition in the default role assignment code that could lead to
   users loosing privileges. When authenticating before the settings service was
   fully running.

   https://github.com/owncloud/ocis/issues/3900

* Bugfix - Fix configuration validation for extensions' server commands: [#3911](https://github.com/owncloud/ocis/pull/3911)

   We've fixed the configuration validation for the extensions' server commands.
   Before this fix error messages have occurred when trying to start individual
   services without certain oCIS fullstack configuration values.

   We now no longer do the common oCIS configuration validation for extensions'
   server commands and now rely only on the extensions' validation function.

   https://github.com/owncloud/ocis/pull/3911

* Bugfix - Fix graph endpoint: [#3925](https://github.com/owncloud/ocis/issues/3925)

   We have added the memberOf slice to the /users endpoint and the member slice to
   the /group endpoint

   https://github.com/owncloud/ocis/issues/3925

* Bugfix - Fix version info: [#3953](https://github.com/owncloud/ocis/pull/3953)

   We've fixed the version info that is displayed when you run:

   - `ocis version` - `ocis <extension name> version`

   Since #2918, these commands returned an empty version only.

   https://github.com/owncloud/ocis/pull/3953

* Bugfix - Remove unused OCS storage configuration: [#3955](https://github.com/owncloud/ocis/pull/3955)

   We've removed the unused OCS configuration option `OCS_STORAGE_USERS_DRIVER`.

   https://github.com/owncloud/ocis/pull/3955

* Bugfix - Make ocdav service behave properly: [#3957](https://github.com/owncloud/ocis/pull/3957)

   The ocdav service now properly passes the tracing config and shuts down when
   receiving a kill signal.

   https://github.com/owncloud/ocis/pull/3957

* Bugfix - Make IDP only wait for certs when using LDAP: [#3965](https://github.com/owncloud/ocis/pull/3965)

   When configuring cs3 as the backend the IDP no longer waits for an LDAP
   certificate to appear.

   https://github.com/owncloud/ocis/pull/3965

* Bugfix - Remove unused configuration options: [#3973](https://github.com/owncloud/ocis/pull/3973)

   We've removed multiple unused configuration options:

   - `STORAGE_SYSTEM_DATAPROVIDER_INSECURE`, see also cs3org/reva#2993 -
   `STORAGE_USERS_DATAPROVIDER_INSECURE`, see also cs3org/reva#2993 -
   `STORAGE_SYSTEM_TEMP_FOLDER`, see also cs3org/reva#2993 -
   `STORAGE_USERS_TEMP_FOLDER`, see also cs3org/reva#2993 -
   `WEB_UI_CONFIG_VERSION`, see also owncloud/web#7130 -
   `GATEWAY_COMMIT_SHARE_TO_STORAGE_REF`, see also cs3org/reva#3017

   https://github.com/owncloud/ocis/pull/3973

* Bugfix - CSP rules for silent token refresh in iframe: [#4031](https://github.com/owncloud/ocis/pull/4031)

   When renewing the access token silently web needs to be opened in an iframe.
   This was previously blocked by a restrictive iframe CSP rule in the `Secure`
   middleware and has now been fixed by allow `self` for iframes.

   https://github.com/owncloud/web/issues/7030
   https://github.com/owncloud/ocis/pull/4031

* Bugfix - Logging in on the wrong account when an email address is not unique: [#4039](https://github.com/owncloud/ocis/issues/4039)

   The default configuration to use the same logon attribute for all services.
   Also, if the configured logon attribute is not unique access to ocis is denied.

   https://github.com/owncloud/ocis/issues/4039

* Bugfix - Remove static ocs user backend config: [#4077](https://github.com/owncloud/ocis/pull/4077)

   We've remove the `OCS_ACCOUNT_BACKEND_TYPE` configuration option. It was
   intended to allow configuration of different user backends for the ocs service.
   Right now the ocs service only has a "cs3" backend. Therefor it's a static entry
   and not configurable.

   https://github.com/owncloud/ocis/pull/4077

* Bugfix - Fix make sensitive config values in the proxy's debug server: [#4086](https://github.com/owncloud/ocis/pull/4086)

   We've fixed a security issue of the proxy's debug server config report endpoint.
   Previously sensitive configuration values haven't been masked. We now mask these
   values.

   https://github.com/owncloud/ocis/pull/4086

* Bugfix - Fix startup error logging: [#4093](https://github.com/owncloud/ocis/pull/4093)

   We've fixed the startup error logging, so that users will the reason for a
   failed startup even on "error" log level. Previously they would only see it on
   "info" log level. Also in a lot of cases the reason for the failed shutdown was
   omitted.

   https://github.com/owncloud/ocis/pull/4093

* Bugfix - Polish search: [#4094](https://github.com/owncloud/ocis/pull/4094)

   We improved the feedback when providing invalid search queries and added support
   for limiting the number of results returned.

   https://github.com/owncloud/ocis/pull/4094

* Bugfix - Fix logging levels: [#4102](https://github.com/owncloud/ocis/pull/4102)

   We've fixed the configuration of logging levels. Previously it was not possible
   to configure a service with a more or less verbose log level then all other
   services when running in the supervised / runtime mode `ocis server`.

   For example `OCIS_LOG_LEVEL=error PROXY_LOG_LEVEL=debug ocis server` did not
   configure error logging for all services except the proxy, which should be on
   debug logging. This is now fixed and working properly.

   Also we fixed the format of go-micro logs to always default to error level.
   Previously this was only ensured in the supervised / runtime mode.

   https://github.com/owncloud/ocis/issues/4089
   https://github.com/owncloud/ocis/pull/4102

* Bugfix - Escape DN attribute value: [#4117](https://github.com/owncloud/ocis/pull/4117)

   Escaped the DN attribute value on creating users and groups.

   https://github.com/owncloud/ocis/pull/4117

* Bugfix - Fix `OCIS_RUN_SERVICES`: [#4133](https://github.com/owncloud/ocis/pull/4133)

   `OCIS_RUN_SERVICES` was introduced as successor to `OCIS_RUN_EXTENSIONS` because
   we wanted to call oCIS "core" extensions services. We kept `OCIS_RUN_EXTENSIONS`
   for backwards compatibility reasons.

   It turned out, that setting `OCIS_RUN_SERVICES` has no effect since introduced.
   `OCIS_RUN_EXTENSIONS`. `OCIS_RUN_EXTENSIONS` was working fine all the time.

   We now fixed `OCIS_RUN_SERVICES`, so that you can use it as a equivalent
   replacement for `OCIS_RUN_EXTENSIONS`

   https://github.com/owncloud/ocis/pull/4133

* Bugfix - Space Creators can hand over spaces: [#4244](https://github.com/owncloud/ocis/pull/4244)

   Set no owner on non personal spaces to be able to pass the space manager role to
   a new user.

   https://github.com/owncloud/ocis/pull/4244

* Bugfix - Fix handling of invalid LDAP users and groups: [#4274](https://github.com/owncloud/ocis/issues/4274)

   We fixed an issue where ocis would exit with a panic when LDAP users or groups
   where missing required attributes (e.g. the id)

   https://github.com/owncloud/ocis/issues/4274

* Bugfix - Fix search in received shares: [#4308](https://github.com/owncloud/ocis/issues/4308)

   We fixed a problem where items in received shares were not found.

   https://github.com/owncloud/ocis/issues/4308

* Bugfix - Fix unrestricted quota on the graphAPI: [#4363](https://github.com/owncloud/ocis/pull/4363)

   Unrestricted quota needs to show 0 on the API. It is not good for clients when
   the property is missing.

   https://github.com/owncloud/ocis/pull/4363

* Bugfix - Autocreate IDP private key also if file exists but is empty: [#4394](https://github.com/owncloud/ocis/pull/4394)

   We've fixed the behavior for the IDP private key generation so that a private
   key is also generated when the file already exists but is empty.

   https://github.com/owncloud/ocis/pull/4394

* Bugfix - Show help for some commands when unconfigured: [#4405](https://github.com/owncloud/ocis/pull/4405)

   We've fixed some commands to show the help also when oCIS is not yet configured.
   Previously the help was not displayed to the user but instead a configuration
   validation error.

   https://github.com/owncloud/ocis/pull/4405

* Bugfix - Rename extensions to services (leftover occurrences): [#4407](https://github.com/owncloud/ocis/pull/4407)

   We've already renamed extensions to services in previous PRs and this PR
   performs this rename for leftover occurrences.

   https://github.com/owncloud/ocis/pull/4407

* Bugfix - Fix configuration of mimetypes for the app registry: [#4411](https://github.com/owncloud/ocis/pull/4411)

   We've fixed the configuration option for mimetypes in the app registry.
   Previously the default config would always be merged over the user provided
   configuration. Now the default mimetype configuration is only used if the user
   does not provide any mimetype configuration (like it is already done in the
   proxy with the routes configuration).

   https://github.com/owncloud/ocis/pull/4411

* Bugfix - Disable default expiration for public links: [#4445](https://github.com/owncloud/ocis/issues/4445)

   The default expiration for public links was enabled in the capabilities without
   providing a (then required) default amount of days for clients to pick a
   reasonable expiration date upon link creation. This has been fixed by disabling
   the default expiration for public links in the capabilities. With this
   configuration clients will no longer set a default expiration date upon link
   creation.

   https://github.com/owncloud/ocis/issues/4445
   https://github.com/owncloud/ocis/pull/4475

* Bugfix - Fix permissions in REPORT: [#4520](https://github.com/owncloud/ocis/pull/4520)

   The REPORT endpoint wouldn't return any permissions on personal spaces Now it
   does. Also bumps reva

   https://github.com/owncloud/ocis/pull/4520

* Bugfix - Render webdav permissions as string in search report: [#4575](https://github.com/owncloud/ocis/issues/4575)

   We now correctly render the `oc:permissions` of resources as a string.

   https://github.com/owncloud/ocis/issues/4575
   https://github.com/owncloud/ocis/pull/4579

* Bugfix - Graph service now forwards trace context: [#4582](https://github.com/owncloud/ocis/pull/4582)

   https://github.com/owncloud/ocis/pull/4582

* Bugfix - Fix sharing jsoncs3 driver options: [#4593](https://github.com/owncloud/ocis/pull/4593)

   We've fixed the environment variable config options of the jsoncs3 driver that
   previously used the same environment variables as the cs3 driver. Now the
   jsoncs3 driver has it's own configuration environment variables.

   If you used the jsoncs3 sharing driver and explicitly set
   `SHARING_PUBLIC_CS3_SYSTEM_USER_ID`, this PR is a breaking change for your
   deployment. To workaround you may set the value you had configured in
   `SHARING_PUBLIC_CS3_SYSTEM_USER_ID` to both
   `SHARING_PUBLIC_JSONCS3_SYSTEM_USER_ID` and
   `SHARING_PUBLIC_JSONCS3_SYSTEM_USER_IDP`.

   https://github.com/owncloud/ocis/pull/4593

* Bugfix - Fix the OIDC provider cache: [#4600](https://github.com/owncloud/ocis/pull/4600)

   We've fixed the OIDC provider cache. It never had a cache hit before this fix.
   Under some circumstances it could cause a painfully slow OCIS if the IDP
   well-known endpoint takes some time to respond.

   https://github.com/owncloud/ocis/pull/4600

* Bugfix - Change the default value for PROXY_OIDC_INSECURE to false: [#4601](https://github.com/owncloud/ocis/pull/4601)

   We've changed the default value for PROXY_OIDC_INSECURE to `false`. Previously
   the default values was `true` which is not acceptable since default values need
   to be secure.

   https://github.com/owncloud/ocis/pull/4601

* Bugfix - Fix authentication for autoprovisioned users: [#4616](https://github.com/owncloud/ocis/issues/4616)

   We've fixed an issue in the proxy, which made the first http request of an
   autoprovisioned user fail.

   https://github.com/owncloud/ocis/issues/4616

* Bugfix - Fix wopi access to public shares: [#4631](https://github.com/owncloud/ocis/pull/4631)

   I've added a request check to the public share authenticator middleware to allow
   wopi to access public shares.

   https://github.com/owncloud/ocis/issues/4382
   https://github.com/owncloud/ocis/pull/4631

* Bugfix - Fix unfindable entities from shares/publicshares: [#4651](https://github.com/owncloud/ocis/pull/4651)

   We fixed a problem where directories or empty files weren't findable because
   they were to the search index improperly when created through a share or
   publicshare.

   https://github.com/owncloud/ocis/issues/4489
   https://github.com/owncloud/ocis/pull/4651

* Bugfix - Fix notifications service settings: [#4652](https://github.com/owncloud/ocis/pull/4652)

   We've fixed two notifications service setting: -
   `NOTIFICATIONS_MACHINE_AUTH_API_KEY` was previously not picked up (only
   `OCIS_MACHINE_AUTH_API_KEY` was loaded) - If you used a email sender address in
   the format of the default value of `NOTIFICATIONS_SMTP_SENDER` no email could be
   send.

   https://github.com/owncloud/ocis/pull/4652

* Bugfix - Bring back the settings UI in Web: [#4691](https://github.com/owncloud/ocis/pull/4691)

   We've fixed the oC Web configuration in oCIS so that the settings UI will be
   shown again in Web.

   https://github.com/owncloud/ocis/pull/4691

* Bugfix - Don't run auth-bearer service by default: [#4692](https://github.com/owncloud/ocis/issues/4692)

   We no longer start the auth-bearer service by default. This service is currently
   unused and not required to run ocis. The equivalent functionality to verify
   OpenID connect tokens and to mint reva tokes for OIDC authenticated clients is
   currently implemented inside the oidc-auth middleware of the proxy.

   https://github.com/owncloud/ocis/issues/4692

* Bugfix - Mail notifications for group shares: [#4714](https://github.com/owncloud/ocis/pull/4714)

   We fixed multiple issues in the notifications service, which broke notification
   mails new shares with groups.

   https://github.com/owncloud/ocis/issues/4703
   https://github.com/owncloud/ocis/issues/4688
   https://github.com/owncloud/ocis/pull/4714

* Bugfix - Make tokeninfo endpoint unprotected: [#4715](https://github.com/owncloud/ocis/pull/4715)

   Make the tokeninfo endpoint unprotected as it is supposed to be available to the
   public.

   https://github.com/owncloud/ocis/pull/4715

* Bugfix - Fix cache stat table config: [#4732](https://github.com/owncloud/ocis/pull/4732)

   We have aligned the cache table config for the gateway and the dataprovider to
   make them actually use the same cache instance.

   https://github.com/owncloud/ocis/pull/4732

* Bugfix - Trigger a rescan of spaces in the search index when items have changed: [#4777](https://github.com/owncloud/ocis/pull/4777)

   The search service now scans spaces when items have been changed. This fixes the
   problem that mtime and treesize propagation was not reflected in the search
   index properly.

   https://github.com/owncloud/ocis/issues/4410
   https://github.com/owncloud/ocis/pull/4777

* Bugfix - Disable cache for selected static web assets: [#4809](https://github.com/owncloud/ocis/pull/4809)

   We've disabled caching for some static web assets. Files like the web
   index.html, oidc-callback.html or similar contain paths to timestamped resources
   and should not be cached.

   https://github.com/owncloud/ocis/pull/4809

* Bugfix - Remove the storage-users event configuration: [#4825](https://github.com/owncloud/ocis/pull/4825)

   We've removed the events configuration from the storage-users section because it
   is not needed.

   https://github.com/owncloud/ocis/pull/4825

* Bugfix - Fix the shareroot path in REPORT responses: [#4859](https://github.com/owncloud/ocis/pull/4859)

   Fixed the shareroot path in REPORT responses. Before this change the attribute
   leaked part of the folder tree of the sharer.

   https://github.com/owncloud/ocis/issues/4796
   https://github.com/owncloud/ocis/pull/4859

* Bugfix - Disable federation capabilities: [#4864](https://github.com/owncloud/ocis/pull/4864)

   We disabled the federation support in the capabilities because it is currently
   not supported.

   https://github.com/owncloud/ocis/pull/4864

* Bugfix - Fix permission check in settings service: [#4890](https://github.com/owncloud/ocis/pull/4890)

   Added a check of the stored roles as a fallback if no roles are contained in the
   context.

   https://github.com/owncloud/ocis/pull/4890

* Bugfix - Fix CORS in frontend service: [#4948](https://github.com/owncloud/ocis/pull/4948)

   We now pass CORS config to the frontend reva service middleware.

   https://github.com/owncloud/ocis/issues/1340
   https://github.com/owncloud/ocis/pull/4948

* Bugfix - Fix notifications Web UI url: [#4998](https://github.com/owncloud/ocis/pull/4998)

   We've fixed the configuration of the notification service's Web UI url that
   appears in emails.

   Previously it was only configurable via the global "OCIS_URL" and is now also
   configurable via "NOTIFICATIONS_WEB_UI_URL".

   https://github.com/owncloud/ocis/pull/4998

* Bugfix - Do not reindex a space twice at the same time: [#5001](https://github.com/owncloud/ocis/pull/5001)

   We fixed a problem where the search service reindexed a space while another
   reindex process was still in progress.

   https://github.com/owncloud/ocis/pull/5001

* Bugfix - Find spaces by their name: [#5044](https://github.com/owncloud/ocis/pull/5044)

   We've fixed finding spaces by their name in the search service.

   https://github.com/owncloud/ocis/issues/4506
   https://github.com/owncloud/ocis/pull/5044

* Bugfix - Initial role assignment with external IDM: [#5045](https://github.com/owncloud/ocis/issues/5045)

   We've the initial user role assignment when using an external LDAP server.

   https://github.com/owncloud/ocis/issues/5045

* Bugfix - Lower IDP token lifespans: [#5077](https://github.com/owncloud/ocis/pull/5077)

   We've lowered the IDP token lifespans to more reasonable durations.

   https://github.com/owncloud/ocis/pull/5077

* Bugfix - Adjust cache related configuration options: [#5087](https://github.com/owncloud/ocis/pull/5087)

   We've adjusted cache related configuration options of the gateway and
   storage-users service to the other services.

   https://github.com/owncloud/ocis/pull/5087

* Bugfix - Make storage users mount ids unique by default: [#5091](https://github.com/owncloud/ocis/pull/5091)

   The mount ID of the storage users provider needs to be unique by default. We
   made this value configurable and added it to ocis init to be sure that we have a
   random uuid v4. This is important for federated instances.

   > **Warning** >BREAKING Change: In order to make every ocis storage provider ID
   unique by default, we needed to use a random uuidv4 during ocis init. Existing
   installations need to set this value explicitly or ocis will terminate after the
   upgrade. > To upgrade from 2.0.0-rc.1 to 2.0.0-rc.2, 2.0.0 or later you need to
   set `GATEWAY_STORAGE_USERS_MOUNT_ID` and `STORAGE_USERS_MOUNT_ID` to the same
   random uuidv4. > >You can also add >``` >storage_users: > mount_id:
   some-random-uuid >gateway: > storage_registry: > storage_users_mount_id:
   some-random-uuid >``` >to the ocis.yaml file which was created during
   initialisation > >Changing the ID of the storage-users provider will change all
   >- WebDAV Urls >- FileIDs >- SpaceIDs >- Bookmarks >- and will make all existing
   shares invalid. > >The Android, Web and iOS clients will continue to work
   without interruptions. The Desktop Client sync connections need to be deleted
   and recreated. >Sorry for the inconvenience 😅 > >WORKAROUND - Not
   Recommended: You can avoid this by setting
   >`GATEWAY_STORAGE_USERS_MOUNT_ID=1284d238-aa92-42ce-bdc4-0b0000009157` and
   >`STORAGE_USERS_MOUNT_ID=1284d238-aa92-42ce-bdc4-0b0000009157` >But this will
   cause problems later when two ocis instances want to federate.

   https://github.com/owncloud/ocis/pull/5091

* Bugfix - Update reva to version 2.12.0: [#5092](https://github.com/owncloud/ocis/pull/5092)

   Changelog for reva 2.12.0 (2022-11-25)  2 ✘  14:57:56 
   =======================================

  *   Bugfix [cs3org/reva#3436](https://github.com/cs3org/reva/pull/3436): Allow updating to internal link
  *   Bugfix [cs3org/reva#3473](https://github.com/cs3org/reva/pull/3473): Decomposedfs fix revision download
  *   Bugfix [cs3org/reva#3482](https://github.com/cs3org/reva/pull/3482): Decomposedfs propagate sizediff
  *   Bugfix [cs3org/reva#3449](https://github.com/cs3org/reva/pull/3449): Don't leak space information on update drive
  *   Bugfix [cs3org/reva#3470](https://github.com/cs3org/reva/pull/3470): Add missing events for managing spaces
  *   Bugfix [cs3org/reva#3472](https://github.com/cs3org/reva/pull/3472): Fix an oCDAV error message
  *   Bugfix [cs3org/reva#3452](https://github.com/cs3org/reva/pull/3452): Fix access to spaces shared via public link
  *   Bugfix [cs3org/reva#3440](https://github.com/cs3org/reva/pull/3440): Set proper names and paths for space roots
  *   Bugfix [cs3org/reva#3437](https://github.com/cs3org/reva/pull/3437): Refactor delete error handling
  *   Bugfix [cs3org/reva#3432](https://github.com/cs3org/reva/pull/3432): Remove share jail fix
  *   Bugfix [cs3org/reva#3458](https://github.com/cs3org/reva/pull/3458): Set the Oc-Fileid header when copying items
  *   Enhancement [cs3org/reva#3441](https://github.com/cs3org/reva/pull/3441): Cover ocdav with more unit tests
  *   Enhancement [cs3org/reva#3493](https://github.com/cs3org/reva/pull/3493): Configurable filelock duration factor in decomposedfs
  *   Enhancement [cs3org/reva#3397](https://github.com/cs3org/reva/pull/3397): Reduce lock contention issues

   https://github.com/owncloud/ocis/pull/5092
   https://github.com/owncloud/ocis/pull/5131

* Bugfix - Decomposedfs increase filelock duration factor: [#5130](https://github.com/owncloud/ocis/pull/5130)

   We made the file lock duration per lock cycle for decomposedfs configurable and
   increased it to make locks work on top of NFS.

   https://github.com/owncloud/ocis/issues/5024
   https://github.com/owncloud/ocis/pull/5130

* Bugfix - Translations on login page: [#7550](https://github.com/owncloud/web/issues/7550)

   We've fixed several translations on the login page. Also, the browser language
   is now being used properly to determine the language.

   https://github.com/owncloud/web/issues/7550
   https://github.com/owncloud/ocis/pull/4504

* Bugfix - Fix search report: [#7557](https://github.com/owncloud/web/issues/7557)

   There were multiple issues with REPORT search responses from webdav. Also we
   want it to be consistent with PROPFIND responses. * the `remote.php` prefix was
   missing from the href (added even though not necessary) * the ids were formatted
   wrong, they should look different for shares and spaces. * the name of the
   resource was missing * the shareid was missing (for shares) * the prop
   `shareroot` (containing the name of the share root) was missing * the
   permissions prop was empty

   https://github.com/owncloud/web/issues/7557
   https://github.com/owncloud/ocis/pull/4485

* Bugfix - Fix unused config option `GRAPH_SPACES_INSECURE`: [#55555](https://github.com/owncloud/ocis/pull/55555)

   We've removed the unused config option `GRAPH_SPACES_INSECURE` from the GRAPH
   service.

   https://github.com/owncloud/ocis/pull/55555

* Change - Switched default configuration to use libregraph/idm: [#3331](https://github.com/owncloud/ocis/pull/3331)

   We switched the default configuration of oCIS to use the "idm" service (based on
   libregraph/idm) as the standard source for user and group information. The
   accounts and glauth services are no longer enabled by default and will be
   removed with an upcoming release.

   https://github.com/owncloud/ocis/pull/3331
   https://github.com/owncloud/ocis/pull/3633

* Change - Introduce `ocis init` and remove all default secrets: [#3551](https://github.com/owncloud/ocis/pull/3551)

   We've removed all default secrets and the hardcoded UUID of the user `admin`.
   This means you can't start oCIS any longer without setting these via environment
   variable or configuration file.

   In order to make this easy for you, we introduced a new command: `ocis init`.
   You can run this command before starting oCIS with `ocis server` and it will
   bootstrap you a configuration file for a secure oCIS instance.

   https://github.com/owncloud/ocis/issues/3524
   https://github.com/owncloud/ocis/pull/3551
   https://github.com/owncloud/ocis/pull/3743

* Change - Load configuration files just from one directory: [#3587](https://github.com/owncloud/ocis/pull/3587)

   We've changed the configuration file loading behavior and are now only loading
   configuration files from ONE single directory. This directory can be set on
   compile time or via an environment variable on startup (`OCIS_CONFIG_DIR`).

   We are using following configuration default paths:

   - Docker images: `/etc/ocis/` - Binary releases: `$HOME/.ocis/config/`

   https://github.com/owncloud/ocis/pull/3587

* Change - Reduce drives in graph /me/drives API: [#3629](https://github.com/owncloud/ocis/pull/3629)

   Reduced the drives in the graph `/me/drives` API to only the drives the user has
   access to. The endpoint `/drives` will list all drives when the user has the
   permission.

   https://github.com/owncloud/ocis/pull/3629

* Change - Reduce permissions on docker image predeclared volumes: [#3641](https://github.com/owncloud/ocis/pull/3641)

   We've lowered the permissions on the predeclared volumes of the oCIS docker
   image from 777 to 750.

   This change doesn't affect you, unless you use the docker image with the non
   default uid/guid to start oCIS (default is 1000:1000).

   https://github.com/owncloud/ocis/pull/3641

* Change - Use new space ID util functions: [#3648](https://github.com/owncloud/ocis/pull/3648)

   Changed code to use the new space ID util functions so that everything works
   with the new spaces ID format.

   https://github.com/owncloud/ocis/pull/3648
   https://github.com/owncloud/ocis/pull/3669

* Change - Rename MetadataUserID: [#3671](https://github.com/owncloud/ocis/pull/3671)

   MetadataUserID is renamed to SystemUserID including yaml tags and env vars

   https://github.com/owncloud/ocis/pull/3671

* Change - Split MachineAuth from SystemUser: [#3672](https://github.com/owncloud/ocis/pull/3672)

   We now have two different APIKeys: MachineAuth for the machine-auth service and
   SystemUser for the system user used e.g. by settings service

   https://github.com/owncloud/ocis/pull/3672

* Change - Rename serviceUser to systemUser: [#3673](https://github.com/owncloud/ocis/pull/3673)

   We renamed serviceUser to systemUser in all configs and vars including yaml-tags
   and env vars

   https://github.com/owncloud/ocis/pull/3673

* Change - Update ocis packages and imports to V2: [#3678](https://github.com/owncloud/ocis/pull/3678)

   This needs to be done in preparation for the major version bump in ocis.

   https://github.com/owncloud/ocis/pull/3678

* Change - The `glauth` and `accounts` services are removed: [#3685](https://github.com/owncloud/ocis/pull/3685)

   After switching the default configuration to libregraph/idm we could remove the
   glauth and accounts services from the source code (they were already disabled by
   default with the previous release)

   https://github.com/owncloud/ocis/pull/3685

* Change - Prevent access to disabled space: [#3779](https://github.com/owncloud/ocis/pull/3779)

   Previously managers where allowed to edit the space even when it is disabled
   This is no longer possible

   https://github.com/owncloud/ocis/pull/3779

* Change - Rename "uploads purge" command to "uploads clean": [#4403](https://github.com/owncloud/ocis/pull/4403)

   We've renamed the storage-users service's "uploads purge" command to "upload
   clean".

   https://github.com/owncloud/ocis/pull/4403

* Change - Enable private links by default: [#4599](https://github.com/owncloud/ocis/pull/4599/)

   Enable private links by default in the capabilities.

   https://github.com/owncloud/ocis/pull/4599/

* Change - Use the spaceID on the cs3 resource: [#4748](https://github.com/owncloud/ocis/pull/4748)

   We cleaned up the CS3Api to use a proper attribute for the space id.

   https://github.com/owncloud/ocis/pull/4748

* Change - Build service frontends with pnpm instead of yarn: [#4878](https://github.com/owncloud/ocis/pull/4878)

   We changed the Node.js packager from Yarn to pnpm to make it more consistent
   with the main Web repo. pnpm offers better package isolation and prevents a
   whole class of errors. This is only relevant for developers.

   https://github.com/owncloud/ocis/pull/4878
   https://github.com/owncloud/web/pull/7835

* Enhancement - Disable the color logging in docker compose examples: [#871](https://github.com/owncloud/ocis/issues/871)

   Disabled the color logging in the example docker compose deployments. Although
   colored logs are helpful during the development process they may be undesired in
   other situations like production deployments, where the logs aren't consumed by
   humans directly but instead by a log aggregator.

   https://github.com/owncloud/ocis/issues/871
   https://github.com/owncloud/ocis/pull/3935

* Enhancement - Product field in OCS version: [#2918](https://github.com/owncloud/ocis/pull/2918)

   We've added a new field to the OCS Version, which is supposed to announce the
   product name. The web ui as a client will make use of it to make the backend
   product and version available (e.g. for easier bug reports).

   https://github.com/owncloud/ocis/pull/2918

* Enhancement - Add /me/changePassword endpoint to GraphAPI: [#3063](https://github.com/owncloud/ocis/issues/3063)

   When using the builtin user management, allow users to update their own password
   via the graph/v1.0/me/changePassword endpoint.

   https://github.com/owncloud/ocis/issues/3063
   https://github.com/owncloud/ocis/pull/3705

* Enhancement - Update IdP UI: [#3493](https://github.com/owncloud/ocis/issues/3493)

   Updated our fork of the lico IdP UI. This also updated the used npm
   dependencies. The design didn't change.

   https://github.com/owncloud/ocis/issues/3493
   https://github.com/owncloud/ocis/pull/4074

* Enhancement - Update reva to v2.3.1: [#3552](https://github.com/owncloud/ocis/pull/3552)

   Updated reva to version 2.3.1. This update includes

  * Bugfix [cs3org/reva#2827](https://github.com/cs3org/reva/pull/2827): Check permissions when deleting spaces
  * Bugfix [cs3org/reva#2830](https://github.com/cs3org/reva/pull/2830): Correctly render response when accepting merged shares
  * Bugfix [cs3org/reva#2831](https://github.com/cs3org/reva/pull/2831): Fix uploads to owncloudsql storage when no mtime is provided
  * Enhancement [cs3org/reva#2833](https://github.com/cs3org/reva/pull/2833): Make status.php values configurable
  * Enhancement [cs3org/reva#2832](https://github.com/cs3org/reva/pull/2832): Add version option for ocdav go-micro service

   Updated reva to version 2.3.0. This update includes:

  * Bugfix [cs3org/reva#2693](https://github.com/cs3org/reva/pull/2693): Support editnew actions from MS Office
  * Bugfix [cs3org/reva#2588](https://github.com/cs3org/reva/pull/2588): Dockerfile.revad-ceph to use the right base image
  * Bugfix [cs3org/reva#2499](https://github.com/cs3org/reva/pull/2499): Removed check DenyGrant in resource permission
  * Bugfix [cs3org/reva#2285](https://github.com/cs3org/reva/pull/2285): Accept new userid idp format
  * Bugfix [cs3org/reva#2802](https://github.com/cs3org/reva/pull/2802): Bugfix the resource id handling for space shares
  * Bugfix [cs3org/reva#2800](https://github.com/cs3org/reva/pull/2800): Bugfix spaceid parsing in spaces trashbin API
  * Bugfix [cs3org/reva#2608](https://github.com/cs3org/reva/pull/2608): Respect the tracing_service_name config variable
  * Bugfix [cs3org/reva#2742](https://github.com/cs3org/reva/pull/2742): Use exact match in login filter
  * Bugfix [cs3org/reva#2759](https://github.com/cs3org/reva/pull/2759): Made uid, gid claims parsing more robust in OIDC auth provider
  * Bugfix [cs3org/reva#2788](https://github.com/cs3org/reva/pull/2788): Return the correct file IDs on public link resources
  * Bugfix [cs3org/reva#2322](https://github.com/cs3org/reva/pull/2322): Use RFC3339 for parsing dates
  * Bugfix [cs3org/reva#2784](https://github.com/cs3org/reva/pull/2784): Disable storageprovider cache for the share jail
  * Bugfix [cs3org/reva#2555](https://github.com/cs3org/reva/pull/2555): Bugfix site accounts endpoints
  * Bugfix [cs3org/reva#2675](https://github.com/cs3org/reva/pull/2675): Updates Makefile according to latest go standards
  * Bugfix [cs3org/reva#2572](https://github.com/cs3org/reva/pull/2572): Wait for nats server on middleware start
  * Change [cs3org/reva#2735](https://github.com/cs3org/reva/pull/2735): Avoid user enumeration
  * Change [cs3org/reva#2737](https://github.com/cs3org/reva/pull/2737): Bump go-cs3api
  * Change [cs3org/reva#2763](https://github.com/cs3org/reva/pull/2763): Change the oCIS and S3NG  storage driver blob store layout
  * Change [cs3org/reva#2596](https://github.com/cs3org/reva/pull/2596): Remove hash from public link urls
  * Change [cs3org/reva#2785](https://github.com/cs3org/reva/pull/2785): Implement workaround for chi.RegisterMethod
  * Change [cs3org/reva#2559](https://github.com/cs3org/reva/pull/2559): Do not encode webDAV ids to base64
  * Change [cs3org/reva#2740](https://github.com/cs3org/reva/pull/2740): Rename oc10 share manager driver
  * Change [cs3org/reva#2561](https://github.com/cs3org/reva/pull/2561): Merge oidcmapping auth manager into oidc
  * Enhancement [cs3org/reva#2698](https://github.com/cs3org/reva/pull/2698): Make capabilities endpoint public, authenticate users is present
  * Enhancement [cs3org/reva#2515](https://github.com/cs3org/reva/pull/2515): Enabling tracing by default if not explicitly disabled
  * Enhancement [cs3org/reva#2686](https://github.com/cs3org/reva/pull/2686): Features for favorites xattrs in EOS, cache for scope expansion
  * Enhancement [cs3org/reva#2494](https://github.com/cs3org/reva/pull/2494): Use sys ACLs for file permissions
  * Enhancement [cs3org/reva#2522](https://github.com/cs3org/reva/pull/2522): Introduce events
  * Enhancement [cs3org/reva#2811](https://github.com/cs3org/reva/pull/2811): Add event for created directories
  * Enhancement [cs3org/reva#2798](https://github.com/cs3org/reva/pull/2798): Add additional fields to events to enable search
  * Enhancement [cs3org/reva#2790](https://github.com/cs3org/reva/pull/2790): Fake providerids so API stays stable after beta
  * Enhancement [cs3org/reva#2685](https://github.com/cs3org/reva/pull/2685): Enable federated account access
  * Enhancement [cs3org/reva#1787](https://github.com/cs3org/reva/pull/1787): Add support for HTTP TPC
  * Enhancement [cs3org/reva#2799](https://github.com/cs3org/reva/pull/2799): Add flag to enable unrestricted listing of spaces
  * Enhancement [cs3org/reva#2560](https://github.com/cs3org/reva/pull/2560): Mentix PromSD extensions
  * Enhancement [cs3org/reva#2741](https://github.com/cs3org/reva/pull/2741): Meta path for user
  * Enhancement [cs3org/reva#2613](https://github.com/cs3org/reva/pull/2613): Externalize custom mime types configuration for storage providers
  * Enhancement [cs3org/reva#2163](https://github.com/cs3org/reva/pull/2163): Nextcloud-based share manager for pkg/ocm/share
  * Enhancement [cs3org/reva#2696](https://github.com/cs3org/reva/pull/2696): Preferences driver refactor and cbox sql implementation
  * Enhancement [cs3org/reva#2052](https://github.com/cs3org/reva/pull/2052): New CS3API datatx methods
  * Enhancement [cs3org/reva#2743](https://github.com/cs3org/reva/pull/2743): Add capability for public link single file edit
  * Enhancement [cs3org/reva#2738](https://github.com/cs3org/reva/pull/2738): Site accounts site-global settings
  * Enhancement [cs3org/reva#2672](https://github.com/cs3org/reva/pull/2672): Further Site Accounts improvements
  * Enhancement [cs3org/reva#2549](https://github.com/cs3org/reva/pull/2549): Site accounts improvements
  * Enhancement [cs3org/reva#2795](https://github.com/cs3org/reva/pull/2795): Add feature flags "projects" and "share_jail" to spaces capability
  * Enhancement [cs3org/reva#2514](https://github.com/cs3org/reva/pull/2514): Reuse ocs role objects in other drivers
  * Enhancement [cs3org/reva#2781](https://github.com/cs3org/reva/pull/2781): In memory user provider
  * Enhancement [cs3org/reva#2752](https://github.com/cs3org/reva/pull/2752): Refactor the rest user and group provider drivers

   https://github.com/owncloud/ocis/issues/3621
   https://github.com/owncloud/ocis/pull/3552
   https://github.com/owncloud/ocis/pull/3570
   https://github.com/owncloud/ocis/pull/3601
   https://github.com/owncloud/ocis/pull/3602
   https://github.com/owncloud/ocis/pull/3605
   https://github.com/owncloud/ocis/pull/3611
   https://github.com/owncloud/ocis/pull/3637
   https://github.com/owncloud/ocis/pull/3652
   https://github.com/owncloud/ocis/pull/3681

* Enhancement - Update linkshare capabilities: [#3579](https://github.com/owncloud/ocis/pull/3579)

   We have updated the capabilities regarding password enforcement and expiration
   dates of public links. They were previously hardcoded in a way that didn't
   reflect the actual backend functionality anymore.

   https://github.com/owncloud/ocis/pull/3579

* Enhancement - Wrap metadata storage with dedicated reva gateway: [#3602](https://github.com/owncloud/ocis/pull/3602)

   We wrapped the metadata storage in a minimal reva instance with a dedicated
   gateway, including static storage registry, static auth registry, in memory
   userprovider, machine authprovider and demo permissions service. This allows us
   to preconfigure the service user for the ocis settings service, share and public
   share providers.

   https://github.com/owncloud/ocis/pull/3602
   https://github.com/owncloud/ocis/pull/3647

* Enhancement - Align service naming: [#3606](https://github.com/owncloud/ocis/pull/3606)

   We now reflect the configured service names when listing them in the ocis
   runtime

   https://github.com/owncloud/ocis/issues/3603
   https://github.com/owncloud/ocis/pull/3606

* Enhancement - Added `share_jail` and `projects` feature flags in spaces capability: [#3626](https://github.com/owncloud/ocis/pull/3626)

   We've added feature flags to the `spaces` capability to indicate to clients
   which features are supposed to be shown to users.

   https://github.com/owncloud/ocis/pull/3626

* Enhancement - Add initial version of the search extensions: [#3635](https://github.com/owncloud/ocis/pull/3635)

   It is now possible to search for files and directories by their name using the
   web UI. Therefor new search extension indexes files in a persistent local index.

   https://github.com/owncloud/ocis/pull/3635

* Enhancement - Don't setup demo role assignments on default: [#3661](https://github.com/owncloud/ocis/issues/3661)

   Added a configuration option to explicitly tell the settings service to generate
   the default role assignments.

   https://github.com/owncloud/ocis/issues/3661
   https://github.com/owncloud/ocis/pull/3956

* Enhancement - Restrict admins from self-removal: [#3713](https://github.com/owncloud/ocis/issues/3713)

   Admin users are no longer allowed to remove their own account or to edit their
   own role assignments. By this restriction we try to prevent situation where no
   administrative users is available in the system anymore

   https://github.com/owncloud/ocis/issues/3713

* Enhancement - Update reva to version 2.4.1: [#3746](https://github.com/owncloud/ocis/pull/3746)

   Changelog for reva 2.4.1 (2022-05-24) =======================================

   The following sections list the changes in reva 2.4.1 relevant to reva users.
   The changes are ordered by importance.

   Summary -------

  * Bugfix [cs3org/reva#2891](https://github.com/cs3org/reva/pull/2891): Add missing http status code

   Changelog for reva 2.4.0 (2022-05-24) =======================================

   The following sections list the changes in reva 2.4.0 relevant to reva users.
   The changes are ordered by importance.

   Summary -------

  * Bugfix [cs3org/reva#2854](https://github.com/cs3org/reva/pull/2854): Handle non uuid space and nodeid in decomposedfs
  * Bugfix [cs3org/reva#2853](https://github.com/cs3org/reva/pull/2853): Filter CS3 share manager listing
  * Bugfix [cs3org/reva#2868](https://github.com/cs3org/reva/pull/2868): Actually remove blobs when purging
  * Bugfix [cs3org/reva#2882](https://github.com/cs3org/reva/pull/2882): Fix FileUploaded event being emitted too early
  * Bugfix [cs3org/reva#2848](https://github.com/cs3org/reva/pull/2848): Fix storage id in the references in the ItemTrashed events
  * Bugfix [cs3org/reva#2852](https://github.com/cs3org/reva/pull/2852): Fix rcbox dependency on reva 1.18
  * Bugfix [cs3org/reva#3505](https://github.com/cs3org/reva/pull/3505): Fix creating a new file with wopi
  * Bugfix [cs3org/reva#2885](https://github.com/cs3org/reva/pull/2885): Move stat out of usershareprovider
  * Bugfix [cs3org/reva#2883](https://github.com/cs3org/reva/pull/2883): Fix role consideration when updating a share
  * Bugfix [cs3org/reva#2864](https://github.com/cs3org/reva/pull/2864): Fix Grant Space IDs
  * Bugfix [cs3org/reva#2870](https://github.com/cs3org/reva/pull/2870): Update quota calculation
  * Bugfix [cs3org/reva#2876](https://github.com/cs3org/reva/pull/2876): Fix version number in status page
  * Bugfix [cs3org/reva#2829](https://github.com/cs3org/reva/pull/2829): Don't include versions in quota
  * Change [cs3org/reva#2856](https://github.com/cs3org/reva/pull/2856): Do not allow to edit disabled spaces
  * Enhancement [cs3org/reva#3741](https://github.com/cs3org/reva/pull/3741): Add download endpoint to ocdav versions API
  * Enhancement [cs3org/reva#2884](https://github.com/cs3org/reva/pull/2884): Show mounted shares in virtual share jail root
  * Enhancement [cs3org/reva#2792](https://github.com/cs3org/reva/pull/2792): Use storageproviderid for spaces routing

   https://github.com/owncloud/ocis/pull/3746
   https://github.com/owncloud/ocis/pull/3771
   https://github.com/owncloud/ocis/pull/3778
   https://github.com/owncloud/ocis/pull/3842
   https://github.com/owncloud/ocis/pull/3854
   https://github.com/owncloud/ocis/pull/3858
   https://github.com/owncloud/ocis/pull/3867

* Enhancement - Add description tags to the thumbnails config structs: [#3752](https://github.com/owncloud/ocis/pull/3752)

   Added description tags to the config structs in the thumbnails service so they
   will be included in the config documentation.

  **Important** If you ran `ocis init` with the `v2.0.0-alpha*` version then you have to manually add the `transfer_secret` to the ocis.yaml.

   Just open the `ocis.yaml` config file and look for the thumbnails section. Then
   add a random `transfer_secret` so that it looks like this:

   ```
   thumbnails:
     thumbnail:
       transfer_secret: <put random value here>
   ```

   https://github.com/owncloud/ocis/pull/3752

* Enhancement - Add acting user to the audit log: [#3753](https://github.com/owncloud/ocis/issues/3753)

   Added the acting user to the events in the audit log.

   https://github.com/owncloud/ocis/issues/3753
   https://github.com/owncloud/ocis/pull/3992

* Enhancement - Add descriptions to webdav configuration: [#3755](https://github.com/owncloud/ocis/pull/3755)

   Added descriptions to webdav config structs to include them in the config
   documentation.

   https://github.com/owncloud/ocis/pull/3755

* Enhancement - Add descriptions for graph-explorer config: [#3759](https://github.com/owncloud/ocis/pull/3759)

   Added descriptions tags to the graph-explorer config tags so that they will be
   included in the documentation.

   https://github.com/owncloud/ocis/pull/3759

* Enhancement - Add config option to provide TLS certificate: [#3818](https://github.com/owncloud/ocis/issues/3818)

   Added a config option to the graph service to provide a TLS certificate to be
   used to verify the LDAP server certificate.

   https://github.com/owncloud/ocis/issues/3818
   https://github.com/owncloud/ocis/pull/3888

* Enhancement - Introduce service registry cache: [#3833](https://github.com/owncloud/ocis/pull/3833)

   We've improved the service registry / service discovery by setting up registry
   caching (TTL 20s), so that not every requests has to do a lookup on the
   registry.

   https://github.com/owncloud/ocis/pull/3833

* Enhancement - Improve validation of OIDC access tokens: [#3841](https://github.com/owncloud/ocis/issues/3841)

   Previously OIDC access tokes were only validated by requesting the userinfo from
   the IDP. It is now possible to enable additional verification if the IDP issues
   access tokens in JWT format. In that case the oCIS proxy service will now verify
   the signature of the token using the public keys provided by jwks_uri endpoint
   of the IDP. It will also verify if the issuer claim (iss) matches the expected
   values.

   The new validation is enabled by setting `PROXY_OIDC_ACCESS_TOKEN_VERIFY_METHOD`
   to "jwt". Which is also the default. Setting it to "none" will disable the
   feature.

   https://github.com/owncloud/ocis/issues/3841
   https://github.com/owncloud/ocis/pull/4227

* Enhancement - Reintroduce user autoprovisioning in proxy: [#3860](https://github.com/owncloud/ocis/pull/3860)

   With the removal of the accounts service autoprovisioning of users upon first
   login was no longer possible. We added this feature back for the cs3 user
   backend in the proxy. Leveraging the libregraph users API for creating the
   users.

   https://github.com/owncloud/ocis/pull/3860

* Enhancement - Allow resharing: [#3904](https://github.com/owncloud/ocis/pull/3904)

   This will allow resharing files

   https://github.com/owncloud/ocis/pull/3904

* Enhancement - Generate signing key and encryption secret: [#3909](https://github.com/owncloud/ocis/issues/3909)

   The idp service now automatically generates a signing key and encryption secret
   when they don't exist. This will enable service restarts without invalidating
   existing sessions.

   https://github.com/owncloud/ocis/issues/3909
   https://github.com/owncloud/ocis/pull/4022

* Enhancement - Add deprecation annotation: [#3917](https://github.com/owncloud/ocis/issues/3917)

   We have added the ability to annotate variables in case of deprecations:

   Example:

   `services/nats/pkg/config/config.go`

   ```
   Host string `yaml:"host" env:"NATS_HOST_ADDRESS,NATS_NATS_HOST" desc:"Bind address." deprecationVersion:"1.6.2" removalVersion:"1.7.5" deprecationInfo:"the name is ugly" deprecationReplacement:"NATS_HOST_ADDRESS"`
   ```

   https://github.com/owncloud/ocis/issues/3917
   https://github.com/owncloud/ocis/pull/5143

* Enhancement - Update reva to version 2.5.1: [#3932](https://github.com/owncloud/ocis/pull/3932)

   Changelog for reva 2.5.1 (2022-06-08) =======================================

   The following sections list the changes in reva 2.5.1 relevant to reva users.
   The changes are ordered by importance.

   Summary -------

  * Bugfix [cs3org/reva#2931](https://github.com/cs3org/reva/pull/2931): Allow listing share jail space
  * Bugfix [cs3org/reva#2918](https://github.com/cs3org/reva/pull/2918): Fix propfinds with depth 0

   Changelog for reva 2.5.0 (2022-06-07) =======================================

   The following sections list the changes in reva 2.5.0 relevant to reva users.
   The changes are ordered by importance.

   Summary -------

  * Bugfix [cs3org/reva#2909](https://github.com/cs3org/reva/pull/2909): The decomposedfs now checks the GetPath permission
  * Bugfix [cs3org/reva#2899](https://github.com/cs3org/reva/pull/2899): Empty meta requests should return body
  * Bugfix [cs3org/reva#2928](https://github.com/cs3org/reva/pull/2928): Fix mkcol response code
  * Bugfix [cs3org/reva#2907](https://github.com/cs3org/reva/pull/2907): Correct share jail child aggregation
  * Bugfix [cs3org/reva#2895](https://github.com/cs3org/reva/pull/2895): Fix unlimited quota in spaces
  * Bugfix [cs3org/reva#2905](https://github.com/cs3org/reva/pull/2905): Check user permissions before updating/removing public shares
  * Bugfix [cs3org/reva#2904](https://github.com/cs3org/reva/pull/2904): Share jail now works properly when accessed as a space
  * Bugfix [cs3org/reva#2903](https://github.com/cs3org/reva/pull/2903): User owncloudsql now uses the correct userid
  * Change [cs3org/reva#2920](https://github.com/cs3org/reva/pull/2920): Clean up the propfind code
  * Change [cs3org/reva#2913](https://github.com/cs3org/reva/pull/2913): Rename ocs parameter "space_ref"
  * Enhancement [cs3org/reva#2919](https://github.com/cs3org/reva/pull/2919): EOS Spaces implementation
  * Enhancement [cs3org/reva#2888](https://github.com/cs3org/reva/pull/2888): Introduce spaces field mask
  * Enhancement [cs3org/reva#2922](https://github.com/cs3org/reva/pull/2922): Refactor webdav error handling

   https://github.com/owncloud/ocis/pull/3932
   https://github.com/owncloud/ocis/pull/3928
   https://github.com/owncloud/ocis/pull/3922

* Enhancement - Add audit events for created containers: [#3941](https://github.com/owncloud/ocis/pull/3941)

   Handle the event `ContainerCreated` in the audit service.

   https://github.com/owncloud/ocis/pull/3941

* Enhancement - Update reva: [#3944](https://github.com/owncloud/ocis/pull/3944)

   Changelog for reva 2.6.1 (2022-06-27) =======================================

   The following sections list the changes in reva 2.6.1 relevant to reva users.
   The changes are ordered by importance.

   Summary -------

  * Bugfix [cs3org/reva#2998](https://github.com/cs3org/reva/pull/2998): Fix 0-byte-uploads
  * Enhancement [cs3org/reva#3983](https://github.com/cs3org/reva/pull/3983): Add capability for alias links
  * Enhancement [cs3org/reva#3000](https://github.com/cs3org/reva/pull/3000): Make less stat requests
  * Enhancement [cs3org/reva#3003](https://github.com/cs3org/reva/pull/3003): Distinguish GRPC FAILED_PRECONDITION and ABORTED codes
  * Enhancement [cs3org/reva#3005](https://github.com/cs3org/reva/pull/3005): Remove unused HomeMapping variable

   Changelog for reva 2.6.0 (2022-06-21) =======================================

   The following sections list the changes in reva 2.6.0 relevant to reva users.
   The changes are ordered by importance.

  * Bugfix [cs3org/reva#2985](https://github.com/cs3org/reva/pull/2985): Make stat requests route based on storage providerid
  * Bugfix [cs3org/reva#2987](https://github.com/cs3org/reva/pull/2987): Let archiver handle all error codes
  * Bugfix [cs3org/reva#2994](https://github.com/cs3org/reva/pull/2994): Bugfix errors when loading shares
  * Bugfix [cs3org/reva#2996](https://github.com/cs3org/reva/pull/2996): Do not close share dump channels
  * Bugfix [cs3org/reva#2993](https://github.com/cs3org/reva/pull/2993): Remove unused configuration
  * Bugfix [cs3org/reva#2950](https://github.com/cs3org/reva/pull/2950): Bugfix sharing with space ref
  * Bugfix [cs3org/reva#2991](https://github.com/cs3org/reva/pull/2991): Make sharesstorageprovider get accepted share
  * Change [cs3org/reva#2877](https://github.com/cs3org/reva/pull/2877): Enable resharing
  * Change [cs3org/reva#2984](https://github.com/cs3org/reva/pull/2984): Update CS3Apis
  * Enhancement [cs3org/reva#3753](https://github.com/cs3org/reva/pull/3753): Add executant to the events
  * Enhancement [cs3org/reva#2820](https://github.com/cs3org/reva/pull/2820): Instrument GRPC and HTTP requests with OTel
  * Enhancement [cs3org/reva#2975](https://github.com/cs3org/reva/pull/2975): Leverage shares space storageid and type when listing shares
  * Enhancement [cs3org/reva#3882](https://github.com/cs3org/reva/pull/3882): Explicitly return on ocdav move requests with body
  * Enhancement [cs3org/reva#2932](https://github.com/cs3org/reva/pull/2932): Stat accepted shares mountpoints, configure existing share updates
  * Enhancement [cs3org/reva#2944](https://github.com/cs3org/reva/pull/2944): Improve owncloudsql connection management
  * Enhancement [cs3org/reva#2962](https://github.com/cs3org/reva/pull/2962): Per service TracerProvider
  * Enhancement [cs3org/reva#2911](https://github.com/cs3org/reva/pull/2911): Allow for dumping and loading shares
  * Enhancement [cs3org/reva#2938](https://github.com/cs3org/reva/pull/2938): Sharpen tooling

   https://github.com/owncloud/ocis/pull/3944
   https://github.com/owncloud/ocis/pull/3975
   https://github.com/owncloud/ocis/pull/3982
   https://github.com/owncloud/ocis/pull/4000
   https://github.com/owncloud/ocis/pull/4006

* Enhancement - Make thumbnails service log less noisy: [#3959](https://github.com/owncloud/ocis/pull/3959)

   Reduced the log severity when no thumbnail was found from warn to debug. This
   reduces the spam in the logs.

   https://github.com/owncloud/ocis/pull/3959

* Enhancement - Refactor extensions to services: [#3980](https://github.com/owncloud/ocis/pull/3980)

   We have decided to name all extensions, we maintain and provide with ocis,
   services from here on to avoid confusion between external extensions and code we
   provide and maintain.

   https://github.com/owncloud/ocis/pull/3980

* Enhancement - Add capability for alias links: [#3983](https://github.com/owncloud/ocis/issues/3983)

   For better UX clients need a way to discover if alias links are supported by the
   server. We added a capability under "files_sharing/public/alias"

   https://github.com/owncloud/ocis/issues/3983
   https://github.com/owncloud/ocis/pull/3991

* Enhancement - New migrate command for migrating shares and public shares: [#3987](https://github.com/owncloud/ocis/pull/3987)

   We added a new `migrate` subcommand which can be used to migrate shares and
   public shares between different share and publicshare managers.

   https://github.com/owncloud/ocis/pull/3987
   https://github.com/owncloud/ocis/pull/4019

* Enhancement - Update ownCloud Web to v5.7.0-rc.1: [#4005](https://github.com/owncloud/ocis/pull/4005)

   Tags: web

   We updated ownCloud Web to v5.7.0-rc.1. Please refer to the changelog (linked)
   for details on the web release.

  * Enhancement [owncloud/web#7119](https://github.com/owncloud/web/pull/7119): Copy/Move conflict dialog
  * Enhancement [owncloud/web#7122](https://github.com/owncloud/web/pull/7122): Enable Drag&Drop and keyboard shortcuts for all views
  * Enhancement [owncloud/web#7053](https://github.com/owncloud/web/pull/7053): Personal space id in URL
  * Enhancement [owncloud/web#6933](https://github.com/owncloud/web/pull/6933): Customize additional mimeTypes for preview app
  * Enhancement [owncloud/web#7078](https://github.com/owncloud/web/pull/7078): Add Hotkeys to ResourceTable
  * Enhancement [owncloud/web#7120](https://github.com/owncloud/web/pull/7120): Use tus chunksize from backend
  * Enhancement [owncloud/web#6749](https://github.com/owncloud/web/pull/6749): Update ODS to v13.2.0-rc.1
  * Enhancement [owncloud/web#7111](https://github.com/owncloud/web/pull/7111): Upload data during creation
  * Enhancement [owncloud/web#7109](https://github.com/owncloud/web/pull/7109): Clickable folder links in upload overlay
  * Enhancement [owncloud/web#7123](https://github.com/owncloud/web/pull/7123): Indeterminate progress bar in upload overlay
  * Enhancement [owncloud/web#7088](https://github.com/owncloud/web/pull/7088): Upload time estimation
  * Enhancement [owncloud/web#7125](https://github.com/owncloud/web/pull/7125): Wording improvements
  * Enhancement [owncloud/web#7140](https://github.com/owncloud/web/pull/7140): Separate direct and indirect link shares in sidebar
  * Bugfix [owncloud/web#7156](https://github.com/owncloud/web/pull/7156): Folder link targets
  * Bugfix [owncloud/web#7108](https://github.com/owncloud/web/pull/7108): Reload of an updated space-image and/or -readme
  * Bugfix [owncloud/web#6846](https://github.com/owncloud/web/pull/6846): Upload meta data serialization
  * Bugfix [owncloud/web#7100](https://github.com/owncloud/web/pull/7100): Complete-state of the upload overlay
  * Bugfix [owncloud/web#7104](https://github.com/owncloud/web/pull/7104): Parent folder name on public links
  * Bugfix [owncloud/web#7173](https://github.com/owncloud/web/pull/7173): Re-introduce dynamic app name in document title
  * Bugfix [owncloud/web#7166](https://github.com/owncloud/web/pull/7166): External apps fixes

   https://github.com/owncloud/ocis/pull/4005
   https://github.com/owncloud/web/pull/7158
   https://github.com/owncloud/ocis/pull/3990
   https://github.com/owncloud/web/pull/6854
   https://github.com/owncloud/web/releases/tag/v5.7.0-rc.1

* Enhancement - Add FRONTEND_ENABLE_RESHARING env variable: [#4023](https://github.com/owncloud/ocis/pull/4023)

   We introduced resharing which was enabled by default, this is now configurable
   and can be enabled by setting the env `FRONTEND_ENABLE_RESHARING` to `true`. By
   default resharing is now disabled.

   https://github.com/owncloud/ocis/pull/4023

* Enhancement - Add drives field to users endpoint: [#4072](https://github.com/owncloud/ocis/pull/4072)

   We have added `$expand=drives` to the `/users/{id}/` endpoint using the user
   filter implemented in reva.

   https://github.com/owncloud/ocis/pull/4072
   https://github.com/cs3org/reva/pull/3046
   https://github.com/owncloud/ocis/pull/4323

* Enhancement - Added command to reset administrator password: [#4084](https://github.com/owncloud/ocis/issues/4084)

   The new command `ocis idm resetpassword` allows to reset the administrator
   password when ocis is not running. So it is possible to recover setups where the
   admin password was lost.

   https://github.com/owncloud/ocis/issues/4084
   https://github.com/owncloud/ocis/pull/4365

* Enhancement - Update reva to version 2.7.2: [#4115](https://github.com/owncloud/ocis/pull/4115)

   Changelog for reva 2.7.2 (2022-07-18) =======================================

  * Bugfix [cs3org/reva#3079](https://github.com/cs3org/reva/pull/3079): Allow empty permissions
  * Bugfix [cs3org/reva#3084](https://github.com/cs3org/reva/pull/3084): Spaces related permissions and providerID cleanup
  * Bugfix [cs3org/reva#3083](https://github.com/cs3org/reva/pull/3083): Add space id to ItemTrashed event

   Changelog for reva 2.7.1 (2022-07-15) =======================================

  * Bugfix [cs3org/reva#3080](https://github.com/cs3org/reva/pull/3080): Make dataproviders return more headers
  * Enhancement [cs3org/reva#3046](https://github.com/cs3org/reva/pull/3046): Add user filter

   Changelog for reva 2.7.0 (2022-07-15) =======================================

  * Bugfix [cs3org/reva#3075](https://github.com/cs3org/reva/pull/3075): Check permissions of the move operation destination
  * Bugfix [cs3org/reva#3036](https://github.com/cs3org/reva/pull/3036):
  * Bugfix revad with EOS docker image
  * Bugfix [cs3org/reva#3037](https://github.com/cs3org/reva/pull/3037): Add uid- and gidNumber to LDAP queries
  * Bugfix [cs3org/reva#4061](https://github.com/cs3org/reva/pull/4061): Forbid resharing with higher permissions
  * Bugfix [cs3org/reva#3017](https://github.com/cs3org/reva/pull/3017): Removed unused gateway config "commit_share_to_storage_ref"
  * Bugfix [cs3org/reva#3031](https://github.com/cs3org/reva/pull/3031): Return proper response code when detecting recursive copy/move operations
  * Bugfix [cs3org/reva#3071](https://github.com/cs3org/reva/pull/3071): Make CS3 sharing drivers parse legacy resource id
  * Bugfix [cs3org/reva#3035](https://github.com/cs3org/reva/pull/3035): Prevent cross space move
  * Bugfix [cs3org/reva#3074](https://github.com/cs3org/reva/pull/3074): Send storage provider and space id to wopi server
  * Bugfix [cs3org/reva#3022](https://github.com/cs3org/reva/pull/3022): Improve the sharing internals
  * Bugfix [cs3org/reva#2977](https://github.com/cs3org/reva/pull/2977): Test valid filename on spaces tus upload
  * Change [cs3org/reva#3006](https://github.com/cs3org/reva/pull/3006): Use spaceID on the cs3api
  * Enhancement [cs3org/reva#3043](https://github.com/cs3org/reva/pull/3043): Introduce LookupCtx for index interface
  * Enhancement [cs3org/reva#3009](https://github.com/cs3org/reva/pull/3009): Prevent recursive copy/move operations
  * Enhancement [cs3org/reva#2977](https://github.com/cs3org/reva/pull/2977): Skip space lookup on space propfind

   https://github.com/owncloud/ocis/pull/4115
   https://github.com/owncloud/ocis/pull/4201
   https://github.com/owncloud/ocis/pull/4203
   https://github.com/owncloud/ocis/pull/4025
   https://github.com/owncloud/ocis/pull/4211

* Enhancement - Search service at the old webdav endpoint: [#4118](https://github.com/owncloud/ocis/pull/4118)

   We made the search service available for legacy clients at the old webdav
   endpoint.

   https://github.com/owncloud/ocis/pull/4118

* Enhancement - Update ownCloud Web to v5.7.0-rc.4: [#4140](https://github.com/owncloud/ocis/pull/4140)

   Tags: web

   We updated ownCloud Web to v5.7.0-rc.4. Please refer to the changelog (linked)
   for details on the web release.

  * Bugfix [owncloud/web#7230](https://github.com/owncloud/web/pull/7230): Context menu misplaced when triggered by keyboard navigation
  * Bugfix [owncloud/web#7214](https://github.com/owncloud/web/pull/7214): Prevent error when pasting with empty clipboard
  * Bugfix [owncloud/web#7173](https://github.com/owncloud/web/pull/7173): Re-introduce dynamic app name in document title
  * Bugfix [owncloud/web#7166](https://github.com/owncloud/web/pull/7166): External apps fixes
  * Bugfix [owncloud/web#7248](https://github.com/owncloud/web/pull/7248): Hide empty trash bin modal on error
  * Bugfix [owncloud/web#4677](https://github.com/owncloud/web/issues/4677): Logout deleted user on page reload
  * Bugfix [owncloud/web#7216](https://github.com/owncloud/web/pull/7216): Filename hovers over the image in the preview app
  * Bugfix [owncloud/web#7228](https://github.com/owncloud/web/pull/7228): Shared with others page apps not working with oc10 as backend
  * Bugfix [owncloud/web#7197](https://github.com/owncloud/web/pull/7197): Create space and access user management permission
  * Bugfix [owncloud/web#6921](https://github.com/owncloud/web/pull/6921): Space sidebar sharing indicators
  * Bugfix [owncloud/web#7030](https://github.com/owncloud/web/issues/7030): Access token renewal
  * Enhancement [owncloud/web#7217](https://github.com/owncloud/web/pull/7217): Add app top bar component
  * Enhancement [owncloud/web#7153](https://github.com/owncloud/web/pull/7153): Add Keyboard navigation/selection
  * Enhancement [owncloud/web#7030](https://github.com/owncloud/web/issues/7030): Loading context blocks application bootstrap
  * Enhancement [owncloud/web#7206](https://github.com/owncloud/web/pull/7206): Add change own password dialog to the account info page
  * Enhancement [owncloud/web#7086](https://github.com/owncloud/web/pull/7086): Re-sharing for ocis
  * Enhancement [owncloud/web#7201](https://github.com/owncloud/web/pull/7201): Added a toolbar to pdf-viewer app
  * Enhancement [owncloud/web#7139](https://github.com/owncloud/web/pull/7139): Reposition notifications
  * Enhancement [owncloud/web#7030](https://github.com/owncloud/web/issues/7030): Resolve bookmarked public links with password protection
  * Enhancement [owncloud/web#7038](https://github.com/owncloud/web/issues/7038): Improve performance of share indicators
  * Enhancement [owncloud/web#6661](https://github.com/owncloud/web/issues/6661): Option to block file extensions from text-editor app
  * Enhancement [owncloud/web#7139](https://github.com/owncloud/web/pull/7139): Update ODS to v14.0.0-alpha.4
  * Enhancement [owncloud/web#7176](https://github.com/owncloud/web/pull/7176): Introduce group assignments

   https://github.com/owncloud/ocis/pull/4140
   https://github.com/owncloud/web/releases/tag/v5.7.0-rc.4

* Enhancement - Add number of total matches to the search result: [#4189](https://github.com/owncloud/ocis/issues/4189)

   The search service now returns the number of total matches alongside the
   results.

   https://github.com/owncloud/ocis/issues/4189

* Enhancement - Introduce "delete-all-spaces" permission: [#4196](https://github.com/owncloud/ocis/issues/4196)

   This is assigned to the Admin role by default and allows to cleanup orphaned
   spaces (e.g. where the owner as been deleted)

   https://github.com/owncloud/ocis/issues/4196

* Enhancement - Improve error log for "could not get user by claim" error: [#4227](https://github.com/owncloud/ocis/pull/4227)

   We've improved the error log for "could not get user by claim" error where
   previously only the "nil" error has been logged. Now we're logging the message
   from the transport.

   https://github.com/owncloud/ocis/pull/4227

* Enhancement - Allow providing list of services NOT to start: [#4254](https://github.com/owncloud/ocis/pull/4254)

   Until now if one wanted to use a custom version of a service, one needed to
   provide `OCIS_RUN_SERVICES` which is a list of all services to start. Now one
   can provide `OCIS_EXCLUDE_RUN_SERVICES` which is a list of only services not to
   start

   https://github.com/owncloud/ocis/pull/4254

* Enhancement - Introduce insecure flag for smtp email notifications: [#4279](https://github.com/owncloud/ocis/pull/4279)

   We've introduced the `NOTIFICATIONS_SMTP_INSECURE` configuration option, that
   let's you skip certificate verification for smtp email servers.

   https://github.com/owncloud/ocis/pull/4279

* Enhancement - Update reva to v2.7.4: [#4294](https://github.com/owncloud/ocis/pull/4294)

   Updated reva to version 2.7.4 This update includes:

  *  Bugfix [cs3org/reva#3141](https://github.com/cs3org/reva/pull/3141): Check ListGrants permission when listing shares

   Updated reva to version 2.7.3 This update includes:

  *  Bugfix [cs3org/reva#3109](https://github.com/cs3org/reva/pull/3109): Bugfix missing check in MustCheckNodePermissions
  *  Bugfix [cs3org/reva#3086](https://github.com/cs3org/reva/pull/3086): Bugfix crash in ldap authprovider
  *  Bugfix [cs3org/reva#3094](https://github.com/cs3org/reva/pull/3094): Allow removing password from public links
  *  Bugfix [cs3org/reva#3096](https://github.com/cs3org/reva/pull/3096): Bugfix user filter
  *  Bugfix [cs3org/reva#3091](https://github.com/cs3org/reva/pull/3091): Project spaces need no real owner
  *  Bugfix [cs3org/reva#3088](https://github.com/cs3org/reva/pull/3088): Use correct sublogger
  *  Enhancement [cs3org/reva#3123](https://github.com/cs3org/reva/pull/3123): Allow stating links that have no permissions
  *  Enhancement [cs3org/reva#3087](https://github.com/cs3org/reva/pull/3087): Allow to set LDAP substring filter type
  *  Enhancement [cs3org/reva#3098](https://github.com/cs3org/reva/pull/3098): App provider http endpoint uses Form instead of Query
  *  Enhancement [cs3org/reva#3133](https://github.com/cs3org/reva/pull/3133): Admins can set quota on all spaces
  *  Enhancement [cs3org/reva#3117](https://github.com/cs3org/reva/pull/3117): Update go-ldap to v3.4.4
  *  Enhancement [cs3org/reva#3095](https://github.com/cs3org/reva/pull/3095): Upload expiration and cleanup

   Https://github.com/owncloud/ocis/pull/4272
   https://github.com/cs3org/reva/pull/3096
   https://github.com/cs3org/reva/pull/4315

   https://github.com/owncloud/ocis/pull/4294
   https://github.com/owncloud/ocis/pull/4330
   https://github.com/owncloud/ocis/pull/4369

* Enhancement - Update ownCloud Web to v5.7.0-rc.8: [#4314](https://github.com/owncloud/ocis/pull/4314)

   Tags: web

   We updated ownCloud Web to v5.7.0-rc.9. Please refer to the changelog (linked)
   for details on the web release.

  * Bugfix [owncloud/web#7080](https://github.com/owncloud/web/issues/7080): Add Droparea again
  * Bugfix [owncloud/web#7357](https://github.com/owncloud/web/pull/7357): Batch deleting multiple files
  * Bugfix [owncloud/web#7379](https://github.com/owncloud/web/pull/7379): Decline share not possible
  * Bugfix [owncloud/web#7322](https://github.com/owncloud/web/pull/7322): Files pagination scroll to top
  * Bugfix [owncloud/web#7348](https://github.com/owncloud/web/pull/7348): Left sidebar active navigation item has wrong cursor
  * Bugfix [owncloud/web#7355](https://github.com/owncloud/web/pull/7355): Link indicator on "Shared via link"-page
  * Bugfix [owncloud/web#7325](https://github.com/owncloud/web/pull/7325): Loading state in views
  * Bugfix [owncloud/web#7344](https://github.com/owncloud/web/pull/7344): Missing file icon in details panel
  * Bugfix [owncloud/web#7321](https://github.com/owncloud/web/pull/7321): Missing scroll bar in user management app
  * Bugfix [owncloud/web#7334](https://github.com/owncloud/web/pull/7334): No redirect after disabling space
  * Bugfix [owncloud/web#3071](https://github.com/owncloud/web/issues/3071): Don't leak oidc callback url into browser history
  * Bugfix [owncloud/web#7379](https://github.com/owncloud/web/pull/7379): Open file on shared space resource not possible
  * Bugfix [owncloud/web#7268](https://github.com/owncloud/web/issues/7268): Personal shares leaked into project space
  * Bugfix [owncloud/web#7359](https://github.com/owncloud/web/pull/7359): Fix infinite loading spinner on invalid preview links
  * Bugfix [owncloud/web#7272](https://github.com/owncloud/web/issues/7272): Print backend version
  * Bugfix [owncloud/web#7424](https://github.com/owncloud/web/pull/7424): Quicklinks not shown
  * Bugfix [owncloud/web#7379](https://github.com/owncloud/web/pull/7379): Rename shared space resource not possible
  * Bugfix [owncloud/web#7210](https://github.com/owncloud/web/pull/7210): Repair navigation highlighter
  * Bugfix [owncloud/web#7393](https://github.com/owncloud/web/pull/7393): Selected item bottom glue
  * Bugfix [owncloud/web#7308](https://github.com/owncloud/web/pull/7308): "Shared with others" and "Shared via Link" resource links not working
  * Bugfix [owncloud/web#7400](https://github.com/owncloud/web/issues/7400): Respect space quota permission
  * Bugfix [owncloud/web#7349](https://github.com/owncloud/web/pull/7349): Missing quick actions in spaces file list
  * Bugfix [owncloud/web#7396](https://github.com/owncloud/web/pull/7396): Add storage ID when navigating to a shared parent directory
  * Bugfix [owncloud/web#7394](https://github.com/owncloud/web/pull/7394): Suppress active panel error log
  * Bugfix [owncloud/web#7038](https://github.com/owncloud/web/issues/7038): File list render performance
  * Bugfix [owncloud/web#7240](https://github.com/owncloud/web/issues/7240): Access token renewal during upload
  * Bugfix [owncloud/web#7376](https://github.com/owncloud/web/pull/7376): Tooltips not shown on disabled create and upload button
  * Bugfix [owncloud/web#7297](https://github.com/owncloud/web/pull/7297): Upload overlay progress bar spacing
  * Bugfix [owncloud/web#7332](https://github.com/owncloud/web/pull/7332): Users list not loading if user has no role
  * Bugfix [owncloud/web#7313](https://github.com/owncloud/web/pull/7313): Versions of shared files not visible
  * Enhancement [owncloud/web#7404](https://github.com/owncloud/web/pull/7404): Adjust helper texts
  * Enhancement [owncloud/web#7350](https://github.com/owncloud/web/pull/7350): Change file loading mechanism in `preview` app
  * Enhancement [owncloud/web#7356](https://github.com/owncloud/web/pull/7356): Declined shares are now easily accessible
  * Enhancement [owncloud/web#7365](https://github.com/owncloud/web/pull/7365): Drop menu styling in right sidebar
  * Enhancement [owncloud/web#7252](https://github.com/owncloud/web/pull/7252): Redesign shared with list
  * Enhancement [owncloud/web#7371](https://github.com/owncloud/web/pull/7371): Use fixed width for the right sidebar
  * Enhancement [owncloud/web#7267](https://github.com/owncloud/web/pull/7267): Search all files announce limit
  * Enhancement [owncloud/web#7364](https://github.com/owncloud/web/pull/7364): Sharing panel show label instead of description for links
  * Enhancement [owncloud/web#7355](https://github.com/owncloud/web/pull/7355): Update ODS to v14.0.0-alpha.12
  * Enhancement [owncloud/web#7375](https://github.com/owncloud/web/pull/7375): User management app saved dialog

   https://github.com/owncloud/ocis/pull/4314
   https://github.com/owncloud/web/releases/tag/v5.7.0-rc.8

* Enhancement - OCS get share now also handle received shares: [#4322](https://github.com/owncloud/ocis/issues/4322)

   Requesting a specific share can now also correctly map the path to the
   mountpoint if the requested share is a received share.

   https://github.com/owncloud/ocis/issues/4322
   https://github.com/owncloud/ocis/pull/4539

* Enhancement - Fix behavior for foobar (in present tense): [#4346](https://github.com/owncloud/ocis/pull/4346)

   We've added the configuration option `PROXY_OIDC_REWRITE_WELLKNOWN` to rewrite
   the `/.well-known/openid-configuration` endpoint. If active, it serves the
   `/.well-known/openid-configuration` response of the original IDP configured in
   `OCIS_OIDC_ISSUER` / `PROXY_OIDC_ISSUER`. This is needed so that the Desktop
   Client, Android Client and iOS Client can discover the OIDC identity provider.

   Previously this rewrite needed to be performed with an external proxy as NGINX
   or Traefik if an external IDP was used.

   https://github.com/owncloud/ocis/issues/2819
   https://github.com/owncloud/ocis/issues/3280
   https://github.com/owncloud/ocis/pull/4346

* Enhancement - Use storageID when requesting special items: [#4356](https://github.com/owncloud/ocis/pull/4356)

   We need to use the storageID when requesting the special items of a space to
   spare a registry lookup and improve the performance

   https://github.com/owncloud/ocis/pull/4356

* Enhancement - Expand personal drive on the graph user: [#4357](https://github.com/owncloud/ocis/pull/4357)

   We can now list the personal drive on the users endpoint via the graph API. A
   user can add an `$expand=drive` query to list the personal drive of the
   requested user.

   https://github.com/owncloud/ocis/pull/4357

* Enhancement - Rewrite of the request authentication middleware: [#4374](https://github.com/owncloud/ocis/pull/4374)

   There were some flaws in the authentication middleware which were resolved by
   this rewrite. This rewrite also introduced the need to manually mark certain
   paths as "unprotected" if requests to these paths must not be authenticated.

   https://github.com/owncloud/ocis/pull/4374

* Enhancement - Add /app/open-with-web endpoint: [#4376](https://github.com/owncloud/ocis/pull/4376)

   We've added an /app/open-with-web endpoint to the app provider, so that clients
   that are no browser or have only limited browser access can also open apps with
   the help of a Web URL.

   https://github.com/owncloud/ocis/pull/4376
   https://github.com/cs3org/reva/pull/3143

* Enhancement - Added language option to the app provider: [#4399](https://github.com/owncloud/ocis/pull/4399)

   We've added a language option to the app provider which will in the end be
   passed to the app a user opens so that the web ui is displayed in the users
   language.

   https://github.com/owncloud/ocis/issues/4367
   https://github.com/owncloud/ocis/pull/4399
   https://github.com/cs3org/reva/pull/3156

* Enhancement - Refactor the proxy service: [#4401](https://github.com/owncloud/ocis/issues/4401)

   The routes of the proxy service now have a "unprotected" flag. This is used by
   the authentication middleware to determine if the request needs to be blocked
   when missing authentication or not.

   https://github.com/owncloud/ocis/issues/4401
   https://github.com/owncloud/ocis/issues/4497
   https://github.com/owncloud/ocis/pull/4461
   https://github.com/owncloud/ocis/pull/4498
   https://github.com/owncloud/ocis/pull/4514

* Enhancement - Add previewFileMimeTypes to web default config: [#4414](https://github.com/owncloud/ocis/pull/4414)

   We've added previewFileMimeTypes to the web default config, so web can determine
   which preview types are supported by the backend.

   https://github.com/owncloud/ocis/pull/4414

* Enhancement - Update ownCloud Web to v5.7.0-rc.10: [#4439](https://github.com/owncloud/ocis/pull/4439)

   Tags: web

   We updated ownCloud Web to v5.7.0-rc.10. Please refer to the changelog (linked)
   for details on the web release.

  * Bugfix [owncloud/web#7443](https://github.com/owncloud/web/pull/7443): Datetime formatting
  * Bugfix [owncloud/web#7437](https://github.com/owncloud/web/pull/7437): Default to user context
  * Bugfix [owncloud/web#7473](https://github.com/owncloud/web/pull/7473): Dragging a file causes no selection
  * Bugfix [owncloud/web#7469](https://github.com/owncloud/web/pull/7469): File size not updated while restoring file version
  * Bugfix [owncloud/web#7443](https://github.com/owncloud/web/pull/7443): File size formatting
  * Bugfix [owncloud/web#7474](https://github.com/owncloud/web/pull/7474): Load only supported thumbnails (configurable)
  * Bugfix [owncloud/web#7309](https://github.com/owncloud/web/pull/7309): SidebarNavItem icon flickering
  * Bugfix [owncloud/web#7425](https://github.com/owncloud/web/pull/7425): Open Folder in project space context menu
  * Bugfix [owncloud/web#7486](https://github.com/owncloud/web/issues/7486): Prevent unnecessary PROPFIND request during upload
  * Bugfix [owncloud/web#7415](https://github.com/owncloud/web/pull/7415): Re-fetch quota
  * Bugfix [owncloud/web#7478](https://github.com/owncloud/web/issues/7478): "Shared via"-indicator for links
  * Bugfix [owncloud/web#7480](https://github.com/owncloud/web/issues/7480): Missing space image in sidebar
  * Bugfix [owncloud/web#7436](https://github.com/owncloud/web/issues/7436): Hide share actions for space viewers/editors
  * Bugfix [owncloud/web#7445](https://github.com/owncloud/web/pull/7445): User management app close side bar throws error
  * Enhancement [owncloud/web#7309](https://github.com/owncloud/web/pull/7309): Keyboard shortcut indicators in ContextMenu
  * Enhancement [owncloud/web#7309](https://github.com/owncloud/web/pull/7309): Lowlight cut resources
  * Enhancement [owncloud/web#7133](https://github.com/owncloud/web/pull/7133): Permissionless (internal) link shares
  * Enhancement [owncloud/web#7309](https://github.com/owncloud/web/pull/7309): Replace locationpicker with clipboard actions
  * Enhancement [owncloud/web#7363](https://github.com/owncloud/web/pull/7363): Streamline UI sizings
  * Enhancement [owncloud/web#7355](https://github.com/owncloud/web/pull/7355): Update ODS to v14.0.0-alpha.16
  * Enhancement [owncloud/web#7476](https://github.com/owncloud/web/pull/7476): Users table on small screen
  * Enhancement [owncloud/web#7182](https://github.com/owncloud/web/pull/7182): User management app edit quota

   https://github.com/owncloud/ocis/pull/4439
   https://github.com/owncloud/web/releases/tag/v5.7.0-rc.10

* Enhancement - Add configuration options for mail authentication and encryption: [#4443](https://github.com/owncloud/ocis/pull/4443)

   We've added configuration options to configure the authentication and encryption
   for sending mails in the notifications service.

   Furthermore there is now a distinguished configuration option for the username
   to use for authentication against the mail server. This allows you to customize
   the sender address to your liking. For example sender addresses like `my oCIS
   instance <ocis@owncloud.test>` are now possible, too.

   https://github.com/owncloud/ocis/pull/4443

* Enhancement - Update reva to v2.8.0: [#4444](https://github.com/owncloud/ocis/pull/4444)

   Updated reva to version 2.8.0. This update includes:

  * Bugfix [cs3org/reva#3158](https://github.com/cs3org/reva/pull/3158): Add name to the propfind response
  * Bugfix [cs3org/reva#3157](https://github.com/cs3org/reva/pull/3157): Fix locking response codes
  * Bugfix [cs3org/reva#3152](https://github.com/cs3org/reva/pull/3152): Disable caching of not found stat responses
  * Bugfix [cs3org/reva#4251](https://github.com/cs3org/reva/pull/4251): Disable caching
  * Enhancement [cs3org/reva#3154](https://github.com/cs3org/reva/pull/3154): Dataproviders now return file metadata
  * Enhancement [cs3org/reva#3143](https://github.com/cs3org/reva/pull/3143): Add /app/open-with-web endpoint
  * Enhancement [cs3org/reva#3156](https://github.com/cs3org/reva/pull/3156): Added language option to the app provider
  * Enhancement [cs3org/reva#3148](https://github.com/cs3org/reva/pull/3148): Add new jsoncs3 share manager

   https://github.com/owncloud/ocis/pull/4444

* Enhancement - Add missing unprotected paths: [#4454](https://github.com/owncloud/ocis/pull/4454)

   Added missing unprotected paths for the text-editor, preview, pdf-viewer,
   draw-io and index.html to the authentication middleware.

   https://github.com/owncloud/ocis/pull/4454
   https://github.com/owncloud/ocis/pull/4458

* Enhancement - Automatically orientate photos when generating thumbnails: [#4477](https://github.com/owncloud/ocis/issues/4477)

   The thumbnailer now makes use of the exif orientation information to
   automatically orientate pictures before generating thumbnails.

   https://github.com/owncloud/ocis/issues/4477
   https://github.com/owncloud/ocis/pull/4513

* Enhancement - Improve login screen design: [#4500](https://github.com/owncloud/ocis/pull/4500)

   We've improved the design of the login screen to match with the general design
   used in Web.

   https://github.com/owncloud/web/issues/7552
   https://github.com/owncloud/ocis/pull/4500

* Enhancement - Update ownCloud Web to v5.7.0: [#4508](https://github.com/owncloud/ocis/pull/4508)

   Tags: web

   We updated ownCloud Web to v5.7.0. Please refer to the changelog (linked) for
   details on the web release.

  * Bugfix [owncloud/web#7522](https://github.com/owncloud/web/pull/7522): Allow uploads outside of user's home despite quota being exceeded
  * Bugfix [owncloud/web#7622](https://github.com/owncloud/web/issues/7622): Expiration date picker with long language codes
  * Bugfix [owncloud/web#7516](https://github.com/owncloud/web/pull/7516): File name in text editor
  * Bugfix [owncloud/web#7498](https://github.com/owncloud/web/issues/7498): Fix right sidebar content on small screens
  * Bugfix [owncloud/web#7455](https://github.com/owncloud/web/issues/7455): Improve keyboard shortcuts copy/cut files
  * Bugfix [owncloud/web#7510](https://github.com/owncloud/web/issues/7510): Paste action (keyboard) not working in project spaces
  * Bugfix [owncloud/web#7526](https://github.com/owncloud/web/issues/7526): Left sidebar when switching apps
  * Bugfix [owncloud/web#7582](https://github.com/owncloud/web/issues/7582): Merge share with group and group member into one
  * Bugfix [owncloud/web#7534](https://github.com/owncloud/web/issues/7534): Redirect after removing self from space members
  * Bugfix [owncloud/web#7560](https://github.com/owncloud/web/pull/7560): Search share representation
  * Bugfix [owncloud/web#7519](https://github.com/owncloud/web/issues/7519): Sidebar for current folder
  * Bugfix [owncloud/web#7453](https://github.com/owncloud/web/issues/7453): Stuck After Session Expired
  * Bugfix [owncloud/web#7595](https://github.com/owncloud/web/pull/7595): Typo when reading public links capabilities
  * Enhancement [owncloud/web#7570](https://github.com/owncloud/web/pull/7570): Adjust spacing of the files list options menu
  * Enhancement [owncloud/web#7540](https://github.com/owncloud/web/issues/7540): Left sidebar hover effect
  * Enhancement [owncloud/web#7555](https://github.com/owncloud/web/pull/7555): Propose unique file name while creating a new file
  * Enhancement [owncloud/web#7038](https://github.com/owncloud/web/issues/7038): Reduce pagination options
  * Enhancement [owncloud/web#6173](https://github.com/owncloud/web/pull/6173): Remember the UI that was last selected via the application switcher
  * Enhancement [owncloud/web#7584](https://github.com/owncloud/web/pull/7584): Remove clickOutside directive
  * Enhancement [owncloud/web#7485](https://github.com/owncloud/web/pull/7485): Add resource name to the WebDAV properties
  * Enhancement [owncloud/web#7559](https://github.com/owncloud/web/pull/7559): Don't open right sidebar from private links
  * Enhancement [owncloud/web#7586](https://github.com/owncloud/web/pull/7586): Search improvements
  * Enhancement [owncloud/web#7605](https://github.com/owncloud/web/pull/7605): Simplify mime type checking
  * Enhancement [owncloud/web#7626](https://github.com/owncloud/web/pull/7626): Update ODS to v14.0.0-alpha.18
  * Enhancement [owncloud/web#7177](https://github.com/owncloud/web/issues/7177): Update Uppy to v3.0.1
  * Enhancement [owncloud/web#7182](https://github.com/owncloud/web/pull/7182): User management app edit quota

   https://github.com/owncloud/ocis/pull/4508
   https://github.com/owncloud/ocis/pull/4547
   https://github.com/owncloud/ocis/pull/4550
   https://github.com/owncloud/web/releases/tag/v5.7.0

* Enhancement - Update Reva to version 2.10.0: [#4522](https://github.com/owncloud/ocis/pull/4522)

   Changelog for reva 2.10.0 (2022-09-09) =======================================

  * Bugfix [cs3org/reva#3210](https://github.com/cs3org/reva/pull/3210): Jsoncs3 mtime fix
  * Enhancement [cs3org/reva#3213](https://github.com/cs3org/reva/pull/3213): Allow for dumping the public shares from the cs3 publicshare manager
  * Enhancement [cs3org/reva#3199](https://github.com/cs3org/reva/pull/3199): Add support for cs3 storage backends to the json publicshare manager

   Changelog for reva 2.9.0 (2022-09-08) =======================================

  * Bugfix [cs3org/reva#3206](https://github.com/cs3org/reva/pull/3206): Add spaceid when listing share jail mount points
  * Bugfix [cs3org/reva#3194](https://github.com/cs3org/reva/pull/3194): Adds the rootinfo to storage spaces
  * Bugfix [cs3org/reva#3201](https://github.com/cs3org/reva/pull/3201): Fix shareid on PROPFIND
  * Bugfix [cs3org/reva#3176](https://github.com/cs3org/reva/pull/3176): Forbid duplicate shares
  * Bugfix [cs3org/reva#3208](https://github.com/cs3org/reva/pull/3208): Prevent panic in time conversion
  * Bugfix [cs3org/reva#3207](https://github.com/cs3org/reva/pull/3207): Align ocs status code for permission error on publiclink update
  * Enhancement [cs3org/reva#3193](https://github.com/cs3org/reva/pull/3193): Add shareid to PROPFIND
  * Enhancement [cs3org/reva#3180](https://github.com/cs3org/reva/pull/3180): Add canDeleteAllHomeSpaces permission
  * Enhancement [cs3org/reva#3203](https://github.com/cs3org/reva/pull/3203): Added "delete-all-spaces" permission
  * Enhancement [cs3org/reva#3200](https://github.com/cs3org/reva/pull/3200): OCS get share now also handle received shares
  * Enhancement [cs3org/reva#3185](https://github.com/cs3org/reva/pull/3185): Improve ldap authprovider's error reporting
  * Enhancement [cs3org/reva#3179](https://github.com/cs3org/reva/pull/3179): Improve tokeninfo endpoint
  * Enhancement [cs3org/reva#3171](https://github.com/cs3org/reva/pull/3171): Cs3 to jsoncs3 share manager migration
  * Enhancement [cs3org/reva#3204](https://github.com/cs3org/reva/pull/3204): Make the function flockFile private
  * Enhancement [cs3org/reva#3192](https://github.com/cs3org/reva/pull/3192): Enable space members to update shares

   https://github.com/owncloud/ocis/pull/4522
   https://github.com/owncloud/ocis/pull/4534
   https://github.com/owncloud/ocis/pull/4548
   https://github.com/owncloud/ocis/pull/4558

* Enhancement - Add Email templating: [#4564](https://github.com/owncloud/ocis/pull/4564)

   We have added email templating to ocis. Which are send on the SpaceShared and
   ShareCreated event.

   https://github.com/owncloud/ocis/issues/4303
   https://github.com/owncloud/ocis/pull/4564
   https://github.com/cs3org/reva/pull/3252

* Enhancement - Allow to configure applications in Web: [#4578](https://github.com/owncloud/ocis/pull/4578)

   We've added the possibility to configure applications in the Web configuration.

   https://github.com/owncloud/ocis/pull/4578

* Enhancement - Add webURL to space root: [#4588](https://github.com/owncloud/ocis/pull/4588)

   Add the web url to the space root on the graphAPI.

   https://github.com/owncloud/ocis/pull/4588

* Enhancement - Update reva to version 2.11.0: [#4588](https://github.com/owncloud/ocis/pull/4588)

   Changelog for reva 2.11.0 (2022-11-03) =======================================

  *   Bugfix  [cs3org/reva#3282](https://github.com/cs3org/reva/pull/3282):  Use Displayname in wopi apps
  *   Bugfix  [cs3org/reva#3430](https://github.com/cs3org/reva/pull/3430):  Add missing error check in decomposedfs
  *   Bugfix  [cs3org/reva#3298](https://github.com/cs3org/reva/pull/3298):  Make date only expiry dates valid for the whole day
  *   Bugfix  [cs3org/reva#3394](https://github.com/cs3org/reva/pull/3394):  Avoid AppProvider panic
  *   Bugfix  [cs3org/reva#3267](https://github.com/cs3org/reva/pull/3267):  Reduced default cache sizes for smaller memory footprint
  *   Bugfix  [cs3org/reva#3338](https://github.com/cs3org/reva/pull/3338):  Fix malformed uid string in cache
  *   Bugfix  [cs3org/reva#3255](https://github.com/cs3org/reva/pull/3255):  Properly escape oc:name in propfind response
  *   Bugfix  [cs3org/reva#3324](https://github.com/cs3org/reva/pull/3324):  Correct base URL for download URL and href when listing file public links
  *   Bugfix  [cs3org/reva#3278](https://github.com/cs3org/reva/pull/3278):  Fix public share view mode during app open
  *   Bugfix  [cs3org/reva#3377](https://github.com/cs3org/reva/pull/3377):  Fix possible race conditions
  *   Bugfix  [cs3org/reva#3274](https://github.com/cs3org/reva/pull/3274):  Fix "uploader" role permissions
  *   Bugfix  [cs3org/reva#3241](https://github.com/cs3org/reva/pull/3241):  Fix uploading empty files into shares
  *   Bugfix  [cs3org/reva#3251](https://github.com/cs3org/reva/pull/3251):  Make listing xattrs more robust
  *   Bugfix  [cs3org/reva#3287](https://github.com/cs3org/reva/pull/3287):  Return OCS forbidden error when a share already exists
  *   Bugfix  [cs3org/reva#3218](https://github.com/cs3org/reva/pull/3218):  Improve performance when listing received shares
  *   Bugfix  [cs3org/reva#3251](https://github.com/cs3org/reva/pull/3251):  Lock source on move
  *   Bugfix  [cs3org/reva#3238](https://github.com/cs3org/reva/pull/3238):  Return relative used quota amount as a percent value
  *   Bugfix  [cs3org/reva#3279](https://github.com/cs3org/reva/pull/3279):  Polish OCS error responses
  *   Bugfix  [cs3org/reva#3307](https://github.com/cs3org/reva/pull/3307):  Refresh lock in decomposedFS needs to overwrite
  *   Bugfix  [cs3org/reva#3368](https://github.com/cs3org/reva/pull/3368):  Return 404 when no permission to space
  *   Bugfix  [cs3org/reva#3341](https://github.com/cs3org/reva/pull/3341):  Validate s3ng downloads
  *   Bugfix  [cs3org/reva#3284](https://github.com/cs3org/reva/pull/3284):  Prevent nil pointer when requesting user
  *   Bugfix  [cs3org/reva#3257](https://github.com/cs3org/reva/pull/3257):  Fix wopi access to publicly shared files
  *   Change  [cs3org/reva#3267](https://github.com/cs3org/reva/pull/3267):  Decomposedfs no longer stores the idp
  *   Change  [cs3org/reva#3381](https://github.com/cs3org/reva/pull/3381):  Changed Name of the Shares Jail
  *   Enhancement  [cs3org/reva#3381](https://github.com/cs3org/reva/pull/3381):  Add capability for sharing by role
  *   Enhancement  [cs3org/reva#3320](https://github.com/cs3org/reva/pull/3320):  Add the parentID to the ocs and dav responses
  *   Enhancement  [cs3org/reva#3239](https://github.com/cs3org/reva/pull/3239):  Add privatelink to PROPFIND response
  *   Enhancement  [cs3org/reva#3340](https://github.com/cs3org/reva/pull/3340):  Add SpaceOwner to some event
  *   Enhancement  [cs3org/reva#3252](https://github.com/cs3org/reva/pull/3252):  Add SpaceShared event
  *   Enhancement  [cs3org/reva#3297](https://github.com/cs3org/reva/pull/3297):  Update dependencies
  *   Enhancement  [cs3org/reva#3429](https://github.com/cs3org/reva/pull/3429):  Make max lock cycles configurable
  *   Enhancement  [cs3org/reva#3011](https://github.com/cs3org/reva/pull/3011):  Expose capability to deny access in OCS API
  *   Enhancement  [cs3org/reva#3224](https://github.com/cs3org/reva/pull/3224):  Make the jsoncs3 share manager cache ttl configurable
  *   Enhancement  [cs3org/reva#3290](https://github.com/cs3org/reva/pull/3290):  Harden file system accesses
  *   Enhancement  [cs3org/reva#3332](https://github.com/cs3org/reva/pull/3332):  Allow to enable TLS for grpc service
  *   Enhancement  [cs3org/reva#3223](https://github.com/cs3org/reva/pull/3223):  Improve CreateShare grpc error reporting
  *   Enhancement  [cs3org/reva#3376](https://github.com/cs3org/reva/pull/3376):  Improve logging
  *   Enhancement  [cs3org/reva#3250](https://github.com/cs3org/reva/pull/3250):  Allow sharing the gateway caches
  *   Enhancement  [cs3org/reva#3240](https://github.com/cs3org/reva/pull/3240):  We now only encode &, < and > in PROPFIND PCDATA
  *   Enhancement  [cs3org/reva#3334](https://github.com/cs3org/reva/pull/3334):  Secure the nats connection with TLS
  *   Enhancement  [cs3org/reva#3300](https://github.com/cs3org/reva/pull/3300):  Do not leak existence of resources
  *   Enhancement  [cs3org/reva#3233](https://github.com/cs3org/reva/pull/3233):  Allow to override default broker for go-micro base ocdav service
  *   Enhancement  [cs3org/reva#3258](https://github.com/cs3org/reva/pull/3258):  Allow ocdav to share the registry instance with other services
  *   Enhancement  [cs3org/reva#3225](https://github.com/cs3org/reva/pull/3225):  Render file parent id for ocs shares
  *   Enhancement  [cs3org/reva#3222](https://github.com/cs3org/reva/pull/3222):  Support Prefer: return=minimal in PROPFIND
  *   Enhancement  [cs3org/reva#3395](https://github.com/cs3org/reva/pull/3395):  Reduce lock contention issues
  *   Enhancement  [cs3org/reva#3286](https://github.com/cs3org/reva/pull/3286):  Make Refresh Lock operation WOPI compliant
  *   Enhancement  [cs3org/reva#3229](https://github.com/cs3org/reva/pull/3229):  Request counting middleware
  *   Enhancement  [cs3org/reva#3312](https://github.com/cs3org/reva/pull/3312):  Implemented new share filters
  *   Enhancement  [cs3org/reva#3308](https://github.com/cs3org/reva/pull/3308):  Update the ttlcache library
  *   Enhancement  [cs3org/reva#3291](https://github.com/cs3org/reva/pull/3291):  The wopi app driver supports more options

   https://github.com/owncloud/ocis/pull/4588
   https://github.com/owncloud/ocis/pull/4716
   https://github.com/owncloud/ocis/pull/4719
   https://github.com/owncloud/ocis/pull/4750
   https://github.com/owncloud/ocis/pull/4833
   https://github.com/owncloud/ocis/pull/4867
   https://github.com/owncloud/ocis/pull/4903
   https://github.com/owncloud/ocis/pull/4908
   https://github.com/owncloud/ocis/pull/4915
   https://github.com/owncloud/ocis/pull/4964

* Enhancement - Allow to configuring the reva cache store: [#4627](https://github.com/owncloud/ocis/pull/4627)

   We have added the possibility to configure the cache store implementation for
   the users storage.

   https://github.com/owncloud/ocis/pull/4627

* Enhancement - Add thumbnails support for tiff and bmp files: [#4634](https://github.com/owncloud/ocis/pull/4634)

   Support generating thumbnails for tiff and bmp files in the thumbnails service.

   https://github.com/owncloud/ocis/pull/4634

* Enhancement - Add support for REPORT requests to /dav/spaces URLs: [#4661](https://github.com/owncloud/ocis/pull/4661)

   We added support for /dav/spaces REPORT requests which allow for searching
   specific spaces.

   https://github.com/owncloud/ocis/issues/4034
   https://github.com/owncloud/ocis/pull/4661

* Enhancement - Make it possible to configure a WOPI folderurl: [#4716](https://github.com/owncloud/ocis/pull/4716)

   The wopi folder URL is used to jump back from an application to the containing
   folder in the files list.

   https://github.com/owncloud/ocis/pull/4716

* Enhancement - Add curl to the oCIS OCI image: [#4751](https://github.com/owncloud/ocis/pull/4751)

   We've added curl to the oCIS OCI image published on Dockerhub. This can be used
   for eg. healthchecks with the services' health endpoint.

   https://github.com/owncloud/ocis/pull/4751

* Enhancement - Report parent id: [#4757](https://github.com/owncloud/ocis/pull/4757)

   We now index and return the parent id of a resource in search REPORTs.

   https://github.com/owncloud/ocis/issues/4727
   https://github.com/owncloud/ocis/pull/4757

* Enhancement - Secure the nats connection with TLS: [#4781](https://github.com/owncloud/ocis/pull/4781)

   Encrypted the connection to the event broker using TLS. Per default TLS is not
   enabled but can be enabled by setting either `OCIS_EVENTS_ENABLE_TLS=true` or
   the respective service configs:

   - `AUDIT_EVENTS_ENABLE_TLS=true` - `GRAPH_EVENTS_ENABLE_TLS=true` -
   `NATS_EVENTS_ENABLE_TLS=true` - `NOTIFICATIONS_EVENTS_ENABLE_TLS=true` -
   `SEARCH_EVENTS_ENABLE_TLS=true` - `SHARING_EVENTS_ENABLE_TLS=true` -
   `STORAGE_USERS_EVENTS_ENABLE_TLS=true`

   https://github.com/owncloud/ocis/pull/4781
   https://github.com/owncloud/ocis/pull/4800
   https://github.com/owncloud/ocis/pull/4867

* Enhancement - Allow to setup TLS for grpc services: [#4798](https://github.com/owncloud/ocis/pull/4798)

   We added config options to allow enabling TLS encryption for all reva and
   go-micro backed grpc services.

   https://github.com/owncloud/ocis/pull/4798
   https://github.com/owncloud/ocis/pull/4901

* Enhancement - We added e-mail subject templating: [#4799](https://github.com/owncloud/ocis/pull/4799)

   We have added e-mail subject templating.

   https://github.com/owncloud/ocis/pull/4799

* Enhancement - Logging improvements: [#4815](https://github.com/owncloud/ocis/pull/4815)

   We improved the logging of several http services. If possible and present, we
   now log the `X-Request-Id`.

   https://github.com/owncloud/ocis/pull/4815
   https://github.com/owncloud/ocis/pull/4974

* Enhancement - Prohibit users from setting or listing other user's values: [#4897](https://github.com/owncloud/ocis/pull/4897)

   Added checks that users can only set and list their own settings.

   https://github.com/owncloud/ocis/pull/4897

* Enhancement - Deny access to resources: [#4903](https://github.com/owncloud/ocis/pull/4903)

   We added an experimental feature to deny access to a certain resource. This
   feature is disabled by default and considered as EXPERIMENTAL. You can enable it
   by setting FRONTEND_OCS_ENABLE_DENIALS to `true`. It announces an available deny
   access permission via WebDAV on each resource. By convention it is only possible
   to deny access on folders. The clients can check the presence of the feature by
   the capability `deny_access` in the `files_sharing` section.

   https://github.com/owncloud/ocis/pull/4903

* Enhancement - Validate space names: [#4955](https://github.com/owncloud/ocis/pull/4955)

   We now return `BAD REQUEST` when space names are - too long (max 255 characters)
   - containing evil characters (`/`, `\`, `.`, `\\`, `:`, `?`, `*`, `"`, `>`, `<`,
   `|`)

   Additionally leading and trailing spaces will be removed silently.

   https://github.com/owncloud/ocis/pull/4955

* Enhancement - Configurable max lock cycles: [#4965](https://github.com/owncloud/ocis/pull/4965)

   Adds config option for max lock cycles. Also bumps reva

   https://github.com/owncloud/ocis/pull/4965

* Enhancement - Rename AUTH_BASIC_AUTH_PROVIDER envvar: [#4966](https://github.com/owncloud/ocis/pull/4966)

   Rename the `AUTH_BASIC_AUTH_PROVIDER` envvar to `AUTH_BASIC_AUTH_MANAGER`

   https://github.com/owncloud/ocis/pull/4966
   https://github.com/owncloud/ocis/pull/4981

* Enhancement - Default to tls 1.2: [#4969](https://github.com/owncloud/ocis/pull/4969)

   https://github.com/owncloud/ocis/pull/4969

* Enhancement - Add the "hidden" state to the search index: [#5018](https://github.com/owncloud/ocis/pull/5018)

   We changed the search service to store the "hidden" state of entries in the
   search index. That will allow for filtering/searching hidden files in the
   future.

   https://github.com/owncloud/ocis/pull/5018

* Enhancement - Remove windows from ci & release makefile: [#5026](https://github.com/owncloud/ocis/pull/5026)

   We have removed windows from the ci & release makefile

   https://github.com/owncloud/ocis/issues/5011
   https://github.com/owncloud/ocis/pull/5026

* Enhancement - Add tracing to search: [#5113](https://github.com/owncloud/ocis/pull/5113)

   We added tracing to search and its indexer

   https://github.com/owncloud/ocis/issues/5063
   https://github.com/owncloud/ocis/pull/5113

* Enhancement - Update ownCloud Web to v6.0.0: [#5153](https://github.com/owncloud/ocis/pull/5153)

   Tags: web

   We updated ownCloud Web to v6.0.0. Please refer to the changelog (linked) for
   details on the web release.

   ### Breaking changes * BREAKING CHANGE for users in
   [owncloud/web#6648](https://github.com/owncloud/web/issues/6648): breaks
   existing bookmarks - they won't resolve anymore. * BREAKING CHANGE for
   developers in [owncloud/web#6648](https://github.com/owncloud/web/issues/6648):
   the appDefaults composables from web-pkg now work with drive aliases,
   concatenated with relative item paths, instead of webdav paths. If you use the
   appDefaults composables in your application it's likely that your code needs to
   be adapted.

   ### Changes * Bugfix
   [owncloud/web#7419](https://github.com/owncloud/web/issues/7419): Add language
   param opening external app * Bugfix
   [owncloud/web#7731](https://github.com/owncloud/web/pull/7731): "Copy
   Quicklink"-translations * Bugfix
   [owncloud/web#7830](https://github.com/owncloud/web/pull/7830): "Cut" and "Copy"
   actions for current folder * Bugfix
   [owncloud/web#7652](https://github.com/owncloud/web/pull/7652): Disable
   copy/move overwrite on self * Bugfix
   [owncloud/web#7739](https://github.com/owncloud/web/pull/7739): Disable shares
   loading on public and trash locations * Bugfix
   [owncloud/web#7740](https://github.com/owncloud/web/pull/7740): Disappearing
   quicklink in sidebar * Bugfix
   [owncloud/web#7946](https://github.com/owncloud/web/issues/7946): Prevent shares
   from disappearing after sharing with groups * Bugfix
   [owncloud/web#7820](https://github.com/owncloud/web/pull/7820): Edit new created
   user in user management * Bugfix
   [owncloud/web#7936](https://github.com/owncloud/web/pull/7936): Editing text
   files on public pages * Bugfix
   [owncloud/web#7861](https://github.com/owncloud/web/pull/7861): Handle non 2xx
   external app responses * Bugfix
   [owncloud/web#7734](https://github.com/owncloud/web/pull/7734): File name
   reactivity * Bugfix
   [owncloud/web#7975](https://github.com/owncloud/web/pull/7975): Prevent file
   upload when folder creation failed * Bugfix
   [owncloud/web#7724](https://github.com/owncloud/web/pull/7724): Folder conflict
   dialog * Bugfix
   [owncloud/web#7603](https://github.com/owncloud/web/issues/7603): Hide search
   bar in public link context * Bugfix
   [owncloud/web#7889](https://github.com/owncloud/web/pull/7889): Hide share
   indicators on public page * Bugfix
   [owncloud/web#7903](https://github.com/owncloud/web/issues/7903): "Keep
   both"-conflict option * Bugfix
   [owncloud/web#7697](https://github.com/owncloud/web/issues/7697): Link indicator
   on "Shared with me"-page * Bugfix
   [owncloud/web#8007](https://github.com/owncloud/web/pull/8007): Missing password
   form on public drop page * Bugfix
   [owncloud/web#7652](https://github.com/owncloud/web/pull/7652): Inhibit move
   files between spaces * Bugfix
   [owncloud/web#7985](https://github.com/owncloud/web/pull/7985): Prevent retrying
   uploads with status code 5xx * Bugfix
   [owncloud/web#7811](https://github.com/owncloud/web/pull/7811): Do not load
   files from cache in public links * Bugfix
   [owncloud/web#7941](https://github.com/owncloud/web/pull/7941): Add origin check
   to Draw.io events * Bugfix
   [owncloud/web#7916](https://github.com/owncloud/web/pull/7916): Prefer alias
   links over private links * Bugfix
   [owncloud/web#7640](https://github.com/owncloud/web/pull/7640): "Private
   link"-button alignment * Bugfix
   [owncloud/web#8006](https://github.com/owncloud/web/pull/8006): Public link
   loading on role change * Bugfix
   [owncloud/web#7962](https://github.com/owncloud/web/issues/7962): Quota check
   when replacing files * Bugfix
   [owncloud/web#7748](https://github.com/owncloud/web/pull/7748): Reload file list
   after last share removal * Bugfix
   [owncloud/web#7699](https://github.com/owncloud/web/issues/7699): Remove the
   "close sidebar"-calls on delete * Bugfix
   [owncloud/web#7504](https://github.com/owncloud/web/pull/7504): Resolve upload
   existing folder * Bugfix
   [owncloud/web#7771](https://github.com/owncloud/web/pull/7771): Routing for
   re-shares * Bugfix
   [owncloud/web#7675](https://github.com/owncloud/web/pull/7675): Search bar on
   small screens * Bugfix
   [owncloud/web#7662](https://github.com/owncloud/web/pull/7662): Sidebar for
   received shares in search file list * Bugfix
   [owncloud/web#7873](https://github.com/owncloud/web/pull/7873): Share editing
   after selecting a space * Bugfix
   [owncloud/web#7657](https://github.com/owncloud/web/issues/7657): Share
   permissions for re-shares * Bugfix
   [owncloud/web#7506](https://github.com/owncloud/web/issues/7506): Shares loading
   * Bugfix [owncloud/web#7632](https://github.com/owncloud/web/pull/7632): Sidebar
   toggle icon * Bugfix
   [owncloud/web#7781](https://github.com/owncloud/web/issues/7781): Sidebar
   without highlighted resource * Bugfix
   [owncloud/web#7756](https://github.com/owncloud/web/pull/7756): Try to obtain
   refresh token before the error case * Bugfix
   [owncloud/web#7768](https://github.com/owncloud/web/pull/7768): Hide actions in
   space trash bins * Bugfix
   [owncloud/web#7651](https://github.com/owncloud/web/pull/7651): Spaces on
   "Shared via link"-page * Bugfix
   [owncloud/web#7521](https://github.com/owncloud/web/issues/7521): Spaces
   reactivity on update * Bugfix
   [owncloud/web#7960](https://github.com/owncloud/web/issues/7960): Display error
   messages in text editor * Bugfix
   [owncloud/web#8030](https://github.com/owncloud/web/pull/8030): Saving a file
   multiple times with the text editor * Bugfix
   [owncloud/web#7778](https://github.com/owncloud/web/issues/7778): Trash bin
   sidebar * Bugfix
   [owncloud/web#7956](https://github.com/owncloud/web/issues/7956): Introduce
   "upload finalizing"-state in upload overlay * Bugfix
   [owncloud/web#7630](https://github.com/owncloud/web/pull/7630): Upload modify
   time * Bugfix [owncloud/web#8011](https://github.com/owncloud/web/issues/8011):
   Prevent unnecessary request when saving a user * Bugfix
   [owncloud/web#7989](https://github.com/owncloud/web/pull/7989): Versions on the
   "Shared with me"-page * Change
   [owncloud/web#6648](https://github.com/owncloud/web/issues/6648): Drive aliases
   in URLs * Change [owncloud/web#7935](https://github.com/owncloud/web/pull/7935):
   Remove mediaSource and v-image-source * Enhancement
   [owncloud/web#7635](https://github.com/owncloud/web/pull/7635): Add restore
   conflict dialog * Enhancement
   [owncloud/web#7901](https://github.com/owncloud/web/pull/7901): Add search field
   for space members * Enhancement
   [owncloud/web#4675](https://github.com/owncloud/web/issues/4675): Add
   `X-Request-ID` header to all outgoing requests * Enhancement
   [owncloud/web#7904](https://github.com/owncloud/web/pull/7904): Batch actions
   for two or more items only * Enhancement
   [owncloud/web#7892](https://github.com/owncloud/web/pull/7892): Respect the new
   sharing denials capability (experimental) * Enhancement
   [owncloud/web#7709](https://github.com/owncloud/web/pull/7709): Edit custom
   permissions wording * Enhancement
   [owncloud/web#7373](https://github.com/owncloud/web/issues/7373): Align dark
   mode colors with given design * Enhancement
   [owncloud/web#7190](https://github.com/owncloud/web/pull/7190): Deny subfolders
   inside share * Enhancement
   [owncloud/web#7684](https://github.com/owncloud/web/pull/7684): Design polishing
   * Enhancement [owncloud/web#7865](https://github.com/owncloud/web/pull/7865):
   Disable share renaming * Enhancement
   [owncloud/web#7725](https://github.com/owncloud/web/pull/7725): Enable renaming
   on received shares * Enhancement
   [owncloud/web#7747](https://github.com/owncloud/web/pull/7747): Friendlier
   logout screen * Enhancement
   [owncloud/web#6247](https://github.com/owncloud/web/issues/6247): Id based
   routing * Enhancement
   [owncloud/web#7803](https://github.com/owncloud/web/issues/7803): Internal link
   on unaccepted share * Enhancement
   [owncloud/web#7304](https://github.com/owncloud/web/issues/7304): Resolve
   internal links * Enhancement
   [owncloud/web#7569](https://github.com/owncloud/web/pull/7569): Make keybindings
   global * Enhancement
   [owncloud/web#7894](https://github.com/owncloud/web/pull/7894): Optimize email
   validation in the user management app * Enhancement
   [owncloud/web#7707](https://github.com/owncloud/web/issues/7707): Resolve
   private links * Enhancement
   [owncloud/web#7234](https://github.com/owncloud/web/issues/7234): Auth context
   in route meta props * Enhancement
   [owncloud/web#7821](https://github.com/owncloud/web/pull/7821): Improve search
   experience * Enhancement
   [owncloud/web#7801](https://github.com/owncloud/web/pull/7801): Make search
   results sortable * Enhancement
   [owncloud/web#8028](https://github.com/owncloud/web/pull/8028): Update ODS to
   v14.0.1 * Enhancement
   [owncloud/web#7890](https://github.com/owncloud/web/pull/7890): Validate space
   names * Enhancement
   [owncloud/web#7430](https://github.com/owncloud/web/pull/7430): Webdav support
   in web-client package * Enhancement
   [owncloud/web#7900](https://github.com/owncloud/web/issues/7900): XHR upload
   timeout

   https://github.com/owncloud/ocis/pull/5153
   https://github.com/owncloud/web/releases/tag/v6.0.0

* Enhancement - Add capability for public link single file edit: [#6787](https://github.com/owncloud/web/pull/6787)

   It is now possible to share a single file by link with edit permissions.
   Therefore we need a public share capability to enable that feature in the
   clients. At the same time, we improved the WebDAV permissions for public links.

   https://github.com/owncloud/web/pull/6787
   https://github.com/owncloud/ocis/pull/3538

* Enhancement - Update ownCloud Web to v5.5.0-rc.8: [#6854](https://github.com/owncloud/web/pull/6854)

   Tags: web

   We updated ownCloud Web to v5.5.0-rc.8. Please refer to the changelog (linked)
   for details on the web release.

   https://github.com/owncloud/web/pull/6854
   https://github.com/owncloud/ocis/pull/3844
   https://github.com/owncloud/ocis/pull/3862
   https://github.com/owncloud/web/releases/tag/v5.5.0-rc.8

* Enhancement - Update ownCloud Web to v5.5.0-rc.9: [#6854](https://github.com/owncloud/web/pull/6854)

   Tags: web

   We updated ownCloud Web to v5.5.0-rc.9. Please refer to the changelog (linked)
   for details on the web release.

   Summary -------

  * Bugfix [owncloud/web#6939](https://github.com/owncloud/web/pull/6939): Not logged out if backend is ownCloud 10
  * Bugfix [owncloud/web#7061](https://github.com/owncloud/web/pull/7061): Prevent rename button from getting covered
  * Bugfix [owncloud/web#7032](https://github.com/owncloud/web/pull/7032): Show message when upload size exceeds quota
  * Bugfix [owncloud/web#7036](https://github.com/owncloud/web/pull/7036): Drag and drop upload when a file is selected
  * Enhancement [owncloud/web#7022](https://github.com/owncloud/web/pull/7022): Add config option for hoverable quick actions
  * Enhancement [owncloud/web#6555](https://github.com/owncloud/web/issues/6555): Consistent dropdown menus
  * Enhancement [owncloud/web#6994](https://github.com/owncloud/web/pull/6994): Copy/Move conflict dialog
  * Enhancement [owncloud/web#6750](https://github.com/owncloud/web/pull/6750): Make contexthelpers opt-out
  * Enhancement [owncloud/web#7038](https://github.com/owncloud/web/issues/7038): Rendering of share-indicators in ResourceTable
  * Enhancement [owncloud/web#6776](https://github.com/owncloud/web/issues/6776): Prevent the resource name in the sidebar from being truncated
  * Enhancement [owncloud/web#7067](https://github.com/owncloud/web/pull/7067): Upload progress & overlay improvements

   https://github.com/owncloud/web/pull/6854
   https://github.com/owncloud/ocis/pull/3927
   https://github.com/owncloud/web/releases/tag/v5.5.0-rc.9

* Enhancement - Update ownCloud Web to v5.5.0-rc.6: [#6854](https://github.com/owncloud/web/pull/6854)

   Tags: web

   We updated ownCloud Web to v5.5.0-rc.6. Please refer to the changelog (linked)
   for details on the web release.

   https://github.com/owncloud/web/pull/6854
   https://github.com/owncloud/ocis/pull/3664
   https://github.com/owncloud/ocis/pull/3680
   https://github.com/owncloud/ocis/pull/3727
   https://github.com/owncloud/ocis/pull/3747
   https://github.com/owncloud/ocis/pull/3797
   https://github.com/owncloud/web/releases/tag/v5.5.0-rc.6

* Enhancement - Optional events in graph service: [#55555](https://github.com/owncloud/ocis/pull/55555)

   We've changed the graph service so that you also can start it without any event
   bus. Therefore you need to set `GRAPH_EVENTS_ENDPOINT` to an empty string. The
   graph API will not emit any events in this case.

   https://github.com/owncloud/ocis/pull/55555

# Changelog for [1.20.0] (2022-04-13)

The following sections list the changes for 1.20.0.

[1.20.0]: https://github.com/owncloud/ocis/compare/v1.19.0...v1.20.0

## Summary

* Bugfix - Ensure the same data on /ocs/v?.php/config like oC10: [#3113](https://github.com/owncloud/ocis/pull/3113)
* Bugfix - Use the default server download protocol if spaces are not supported: [#3386](https://github.com/owncloud/ocis/pull/3386)
* Bugfix - Add `owncloudsql` driver to authprovider config: [#3435](https://github.com/owncloud/ocis/pull/3435)
* Bugfix - Corrected documentation: [#3439](https://github.com/owncloud/ocis/pull/3439)
* Change - Fix keys with underscores in the config files: [#3412](https://github.com/owncloud/ocis/pull/3412)
* Change - Don't create demo users by default: [#3474](https://github.com/owncloud/ocis/pull/3474)
* Enhancement - Add sorting to GraphAPI users and groups: [#3360](https://github.com/owncloud/ocis/issues/3360)
* Enhancement - Use embeddable ocdav go micro service: [#3397](https://github.com/owncloud/ocis/pull/3397)
* Enhancement - Update reva to v2.2.0: [#3397](https://github.com/owncloud/ocis/pull/3397)
* Enhancement - Make config dir configurable: [#3440](https://github.com/owncloud/ocis/pull/3440)
* Enhancement - Replace deprecated String.prototype.substr(): [#3448](https://github.com/owncloud/ocis/pull/3448)
* Enhancement - Alias links: [#3454](https://github.com/owncloud/ocis/pull/3454)
* Enhancement - Implement audit events for user and groups: [#3467](https://github.com/owncloud/ocis/pull/3467)
* Enhancement - Unify LDAP config settings across services: [#3476](https://github.com/owncloud/ocis/pull/3476)
* Enhancement - Update ownCloud Web to v5.4.0: [#6709](https://github.com/owncloud/web/pull/6709)

## Details

* Bugfix - Ensure the same data on /ocs/v?.php/config like oC10: [#3113](https://github.com/owncloud/ocis/pull/3113)

   We've fixed the returned values on the /ocs/v?.php/config endpoints, so that
   they now return the same values as an oC10 would do.

   https://github.com/owncloud/ocis/pull/3113

* Bugfix - Use the default server download protocol if spaces are not supported: [#3386](https://github.com/owncloud/ocis/pull/3386)

   https://github.com/owncloud/ocis/pull/3386

* Bugfix - Add `owncloudsql` driver to authprovider config: [#3435](https://github.com/owncloud/ocis/pull/3435)

   https://github.com/owncloud/ocis/pull/3435

* Bugfix - Corrected documentation: [#3439](https://github.com/owncloud/ocis/pull/3439)

   - ocis-pkg log File Option

   https://github.com/owncloud/ocis/pull/3439

* Change - Fix keys with underscores in the config files: [#3412](https://github.com/owncloud/ocis/pull/3412)

   We've fixed some config keys in configuration files that previously didn't
   contain underscores but should.

   Please check the documentation on https://owncloud.dev for latest configuration
   documentation.

   https://github.com/owncloud/ocis/pull/3412

* Change - Don't create demo users by default: [#3474](https://github.com/owncloud/ocis/pull/3474)

   As we are coming closer to the first beta, we need to disable the creation of
   the demo users by default.

   https://github.com/owncloud/ocis/issues/3181
   https://github.com/owncloud/ocis/pull/3474

* Enhancement - Add sorting to GraphAPI users and groups: [#3360](https://github.com/owncloud/ocis/issues/3360)

   The GraphAPI endpoints for users and groups support ordering now. User can be
   ordered by displayName, onPremisesSamAccountName and mail. Groups can be ordered
   by displayName.

   Example: https://localhost:9200/graph/v1.0/groups?$orderby=displayName asc

   https://github.com/owncloud/ocis/issues/3360

* Enhancement - Use embeddable ocdav go micro service: [#3397](https://github.com/owncloud/ocis/pull/3397)

   We now use the reva `pgk/micro/ocdav` package that implements a go micro
   compatible version of the ocdav service.

   https://github.com/owncloud/ocis/pull/3397

* Enhancement - Update reva to v2.2.0: [#3397](https://github.com/owncloud/ocis/pull/3397)

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

* Enhancement - Make config dir configurable: [#3440](https://github.com/owncloud/ocis/pull/3440)

   We have added an `OCIS_CONFIG_DIR` environment variable the will take precedence
   over the default `/etc/ocis`, `~/.ocis` and `.config` locations. When it is set
   the default locations will be ignored and only the configuration files in that
   directory will be read.

   https://github.com/owncloud/ocis/pull/3440

* Enhancement - Replace deprecated String.prototype.substr(): [#3448](https://github.com/owncloud/ocis/pull/3448)

   We've replaced all occurrences of the deprecated String.prototype.substr()
   function with String.prototype.slice() which works similarly but isn't
   deprecated.

   https://github.com/owncloud/ocis/pull/3448

* Enhancement - Alias links: [#3454](https://github.com/owncloud/ocis/pull/3454)

   Bumps reva and configures ocs token endpoint to be unprotected

   https://github.com/owncloud/ocis/pull/3454

* Enhancement - Implement audit events for user and groups: [#3467](https://github.com/owncloud/ocis/pull/3467)

   Added audit events for users and groups. This will log: * User creation * User
   deletion * User property change (currently only email) * Group creation * Group
   deletion * Group member add * Group member remove

   https://github.com/owncloud/ocis/pull/3467

* Enhancement - Unify LDAP config settings across services: [#3476](https://github.com/owncloud/ocis/pull/3476)

   The storage services where updated to adapt for the recent changes of the LDAP
   settings in reva.

   Also we allow now to use a new set of top-level LDAP environment variables that
   are shared between all LDAP-using services in ocis (graph, idp,
   storage-auth-basic, storage-userprovider, storage-groupprovider, idm). This
   should simplify the most LDAP based configurations considerably.

   Here is a list of the new environment variables: LDAP_URI LDAP_INSECURE
   LDAP_CACERT LDAP_BIND_DN LDAP_BIND_PASSWORD LDAP_LOGIN_ATTRIBUTES
   LDAP_USER_BASE_DN LDAP_USER_SCOPE LDAP_USER_FILTER LDAP_USER_OBJECTCLASS
   LDAP_USER_SCHEMA_MAIL LDAP_USER_SCHEMA_DISPLAY_NAME LDAP_USER_SCHEMA_USERNAME
   LDAP_USER_SCHEMA_ID LDAP_USER_SCHEMA_ID_IS_OCTETSTRING LDAP_GROUP_BASE_DN
   LDAP_GROUP_SCOPE LDAP_GROUP_FILTER LDAP_GROUP_OBJECTCLASS
   LDAP_GROUP_SCHEMA_GROUPNAME LDAP_GROUP_SCHEMA_ID
   LDAP_GROUP_SCHEMA_ID_IS_OCTETSTRING

   Where need these can be overwritten by service specific variables. E.g. it is
   possible to use STORAGE_LDAP_URI to override the top-level LDAP_URI variable.

   https://github.com/owncloud/ocis/issues/3150
   https://github.com/owncloud/ocis/pull/3476

* Enhancement - Update ownCloud Web to v5.4.0: [#6709](https://github.com/owncloud/web/pull/6709)

   Tags: web

   We updated ownCloud Web to v5.4.0. Please refer to the changelog (linked) for
   details on the web release.

   https://github.com/owncloud/web/pull/6709
   https://github.com/owncloud/ocis/pull/3437
   https://github.com/owncloud/ocis/pull/3487
   https://github.com/owncloud/ocis/pull/3509
   https://github.com/owncloud/web/releases/tag/v5.4.0

# Changelog for [1.19.0] (2022-03-29)

The following sections list the changes for 1.19.0.

[1.19.0]: https://github.com/owncloud/ocis/compare/v1.19.1...v1.19.0

## Summary

* Bugfix - Fix request validation on GraphAPI User updates: [#3167](https://github.com/owncloud/ocis/issues/3167)
* Bugfix - Network configuration in individual_services example: [#3238](https://github.com/owncloud/ocis/pull/3238)
* Bugfix - Improve gif thumbnails: [#3305](https://github.com/owncloud/ocis/pull/3305)
* Bugfix - Replace public mountpoint fileid with grant fileid: [#3349](https://github.com/owncloud/ocis/pull/3349)
* Bugfix - Fix error handling in GraphAPI GetUsers call: [#3357](https://github.com/owncloud/ocis/pull/3357)
* Change - Switch NATS backend: [#3192](https://github.com/owncloud/ocis/pull/3192)
* Change - Settings service now stores its data via metadata service: [#3232](https://github.com/owncloud/ocis/pull/3232)
* Change - Add remote item to mountpoint and fix spaceID: [#3365](https://github.com/owncloud/ocis/pull/3365)
* Change - Drop json config file support: [#3366](https://github.com/owncloud/ocis/pull/3366)
* Enhancement - Include etags in drives listing: [#3267](https://github.com/owncloud/ocis/pull/3267)
* Enhancement - Improve thumbnails API: [#3272](https://github.com/owncloud/ocis/pull/3272)
* Enhancement - Add space aliases: [#3283](https://github.com/owncloud/ocis/pull/3283)
* Enhancement - Log sharing events in audit service: [#3301](https://github.com/owncloud/ocis/pull/3301)
* Enhancement - Add password reset link to login page: [#3329](https://github.com/owncloud/ocis/pull/3329)
* Enhancement - Update reva to v2.1.0: [#3330](https://github.com/owncloud/ocis/pull/3330)
* Enhancement - Audit logger will now log file events: [#3332](https://github.com/owncloud/ocis/pull/3332)
* Enhancement - Update ownCloud Web to v5.3.0: [#6561](https://github.com/owncloud/web/pull/6561)

## Details

* Bugfix - Fix request validation on GraphAPI User updates: [#3167](https://github.com/owncloud/ocis/issues/3167)

   Fix PATCH on graph/v1.0/users when no 'mail' attribute is present in the request
   body

   https://github.com/owncloud/ocis/issues/3167

* Bugfix - Network configuration in individual_services example: [#3238](https://github.com/owncloud/ocis/pull/3238)

   Tidy up the deployments/examples/ocis_individual_services example so that the
   instructions work.

   https://github.com/owncloud/ocis/pull/3238

* Bugfix - Improve gif thumbnails: [#3305](https://github.com/owncloud/ocis/pull/3305)

   Improved the gif thumbnail generation for gifs with different disposal
   strategies.

   https://github.com/owncloud/ocis/pull/3305

* Bugfix - Replace public mountpoint fileid with grant fileid: [#3349](https://github.com/owncloud/ocis/pull/3349)

   We now show the same resource id for resources when accessing them via a public
   links as when using a logged in user. This allows the web ui to start a WOPI
   session with the correct resource id.

   https://github.com/owncloud/ocis/pull/3349

* Bugfix - Fix error handling in GraphAPI GetUsers call: [#3357](https://github.com/owncloud/ocis/pull/3357)

   A missing return statement caused GetUsers to return misleading results when the
   identity backend returned an error.

   https://github.com/owncloud/ocis/pull/3357

* Change - Switch NATS backend: [#3192](https://github.com/owncloud/ocis/pull/3192)

   We've switched the NATS backend from Streaming to JetStream, since NATS
   Streaming is depreciated.

   https://github.com/owncloud/ocis/pull/3192
   https://github.com/cs3org/reva/pull/2574

* Change - Settings service now stores its data via metadata service: [#3232](https://github.com/owncloud/ocis/pull/3232)

   Instead of writing files to disk it will use metadata service to do so

   https://github.com/owncloud/ocis/pull/3232

* Change - Add remote item to mountpoint and fix spaceID: [#3365](https://github.com/owncloud/ocis/pull/3365)

   A mountpoint represents the mounted share on the share receivers side. The
   original resource is located where the grant has been set. This item is now
   shown as libregraph remoteItem on the mountpoint. While adding this, we fixed
   the spaceID for mountpoints.

   https://github.com/owncloud/ocis/pull/3365

* Change - Drop json config file support: [#3366](https://github.com/owncloud/ocis/pull/3366)

   We've remove the support to configure oCIS and it's service with a json file.
   From now on we only support yaml configuration files, since they have the
   possibility to add comments.

   https://github.com/owncloud/ocis/pull/3366

* Enhancement - Include etags in drives listing: [#3267](https://github.com/owncloud/ocis/pull/3267)

   Added etags in the response of list drives.

   https://github.com/owncloud/ocis/pull/3267

* Enhancement - Improve thumbnails API: [#3272](https://github.com/owncloud/ocis/pull/3272)

   Changed the thumbnails API to no longer transfer images via GRPC. GRPC has a
   limited message size and isn't very efficient with large binary data. The new
   API transports the images over HTTP.

   https://github.com/owncloud/ocis/pull/3272

* Enhancement - Add space aliases: [#3283](https://github.com/owncloud/ocis/pull/3283)

   Space aliases can be used to resolve spaceIDs in a client.

   https://github.com/owncloud/ocis/pull/3283

* Enhancement - Log sharing events in audit service: [#3301](https://github.com/owncloud/ocis/pull/3301)

   Contains sharing related events. See full list in audit/pkg/types/events.go

   https://github.com/owncloud/ocis/pull/3301

* Enhancement - Add password reset link to login page: [#3329](https://github.com/owncloud/ocis/pull/3329)

   Added a configurable password reset link to the login page. It can be set via
   `IDP_PASSWORD_RESET_URI`. If the option is not set the link will not be shown.

   https://github.com/owncloud/ocis/pull/3329

* Enhancement - Update reva to v2.1.0: [#3330](https://github.com/owncloud/ocis/pull/3330)

   Updated reva to version 2.1.0. This update includes:

  * Fix [cs3org/reva#2636](https://github.com/cs3org/reva/pull/2636): Delay reconnect log for events
  * Fix [cs3org/reva#2645](https://github.com/cs3org/reva/pull/2645): Avoid warning about missing .flock files
  * Fix [cs3org/reva#2625](https://github.com/cs3org/reva/pull/2625): Fix locking on public links and the decomposed filesystem
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

* Enhancement - Audit logger will now log file events: [#3332](https://github.com/owncloud/ocis/pull/3332)

   See full list of supported events in `audit/pkg/types/types.go`

   https://github.com/owncloud/ocis/pull/3332

* Enhancement - Update ownCloud Web to v5.3.0: [#6561](https://github.com/owncloud/web/pull/6561)

   Tags: web

   We updated ownCloud Web to v5.3.0. Please refer to the changelog (linked) for
   details on the web release.

   https://github.com/owncloud/web/pull/6561
   https://github.com/owncloud/ocis/pull/3291
   https://github.com/owncloud/ocis/pull/3375
   https://github.com/owncloud/web/releases/tag/v5.3.0

# Changelog for [1.19.1] (2022-03-29)

The following sections list the changes for 1.19.1.

[1.19.1]: https://github.com/owncloud/ocis/compare/v1.18.0...v1.19.1

## Summary

* Bugfix - Return correct special item urls: [#3419](https://github.com/owncloud/ocis/pull/3419)

## Details

* Bugfix - Return correct special item urls: [#3419](https://github.com/owncloud/ocis/pull/3419)

   URLs for Special items (space image, readme) were broken.

   https://github.com/owncloud/ocis/pull/3419

# Changelog for [1.18.0] (2022-03-03)

The following sections list the changes for 1.18.0.

[1.18.0]: https://github.com/owncloud/ocis/compare/v1.17.0...v1.18.0

## Summary

* Bugfix - Align storage metadata GPRC bind port with other variable names: [#3169](https://github.com/owncloud/ocis/pull/3169)
* Bugfix - Make events settings configurable: [#3214](https://github.com/owncloud/ocis/pull/3214)
* Bugfix - Capabilities for password protected public links: [#3229](https://github.com/owncloud/ocis/pull/3229)
* Change - Unify file IDs: [#3185](https://github.com/owncloud/ocis/pull/3185)
* Enhancement - Re-Enabling web cache control: [#3109](https://github.com/owncloud/ocis/pull/3109)
* Enhancement - Add SPA conform fileserver for web: [#3109](https://github.com/owncloud/ocis/pull/3109)
* Enhancement - Add sorting to list Spaces: [#3200](https://github.com/owncloud/ocis/issues/3200)
* Enhancement - Change NATS port: [#3210](https://github.com/owncloud/ocis/pull/3210)
* Enhancement - Implement notifications service: [#3217](https://github.com/owncloud/ocis/pull/3217)
* Enhancement - Thumbnails in spaces: [#3219](https://github.com/owncloud/ocis/pull/3219)
* Enhancement - Update reva to v2.0.0: [#3231](https://github.com/owncloud/ocis/pull/3231)
* Enhancement - Update ownCloud Web to v5.2.0: [#6506](https://github.com/owncloud/web/pull/6506)

## Details

* Bugfix - Align storage metadata GPRC bind port with other variable names: [#3169](https://github.com/owncloud/ocis/pull/3169)

   Changed STORAGE_METADATA_GRPC_PROVIDER_ADDR to STORAGE_METADATA_GRPC_ADDR so it
   aligns with standard environment variable naming conventions used in oCIS.

   https://github.com/owncloud/ocis/pull/3169

* Bugfix - Make events settings configurable: [#3214](https://github.com/owncloud/ocis/pull/3214)

   We've fixed the hardcoded events settings to be configurable.

   https://github.com/owncloud/ocis/pull/3214

* Bugfix - Capabilities for password protected public links: [#3229](https://github.com/owncloud/ocis/pull/3229)

   Allow password protected public links to request capabilities.

   https://github.com/owncloud/web/issues/5863
   https://github.com/owncloud/ocis/pull/3229
   https://github.com/owncloud/web/pull/6471

* Change - Unify file IDs: [#3185](https://github.com/owncloud/ocis/pull/3185)

   We changed the file IDs to be consistent across all our APIs (WebDAV,
   LibreGraph, OCS). We removed the base64 encoding. Now they are formatted like
   <storageID>!<opaqueID>. They are using a reserved character ``!`` as a URL safe
   separator.

   https://github.com/owncloud/ocis/pull/3185

* Enhancement - Re-Enabling web cache control: [#3109](https://github.com/owncloud/ocis/pull/3109)

   We've re-enable browser caching headers (`Expires` and `Last-Modified`) for the
   web service, this was disabled due to a problem in the fileserver used before.
   Since we're now using our own fileserver implementation this works again and is
   enabled by default.

   https://github.com/owncloud/ocis/pull/3109

* Enhancement - Add SPA conform fileserver for web: [#3109](https://github.com/owncloud/ocis/pull/3109)

   We've added an SPA conform fileserver to the web service. It enables web to use
   vue's history mode and behaves like nginx try_files.

   https://github.com/owncloud/ocis/pull/3109

* Enhancement - Add sorting to list Spaces: [#3200](https://github.com/owncloud/ocis/issues/3200)

   We added the OData query param "orderBy" for listing spaces. We can now order by
   Space Name and LastModifiedDateTime.

   Example 1:
   https://localhost:9200/graph/v1.0/me/drives/?$orderby=lastModifiedDateTime desc
   Example 2: https://localhost:9200/graph/v1.0/me/drives/?$orderby=name asc

   https://github.com/owncloud/ocis/issues/3200
   https://github.com/owncloud/ocis/pull/3201
   https://github.com/owncloud/ocis/pull/3218

* Enhancement - Change NATS port: [#3210](https://github.com/owncloud/ocis/pull/3210)

   Currently only a certain range of ports is allowed for ocis application. Use a
   supported port for nats server

   https://github.com/owncloud/ocis/pull/3210

* Enhancement - Implement notifications service: [#3217](https://github.com/owncloud/ocis/pull/3217)

   Implemented the minimal version of the notifications service to be able to
   notify a user when they received a share.

   https://github.com/owncloud/ocis/pull/3217

* Enhancement - Thumbnails in spaces: [#3219](https://github.com/owncloud/ocis/pull/3219)

   Added support for thumbnails in spaces.

   https://github.com/owncloud/ocis/pull/3219
   https://github.com/owncloud/ocis/pull/3235

* Enhancement - Update reva to v2.0.0: [#3231](https://github.com/owncloud/ocis/pull/3231)

   We updated reva to the version 2.0.0.

  * Fix [cs3org/reva#2457](https://github.com/cs3org/reva/pull/2457) :  Do not swallow error
  * Fix [cs3org/reva#2422](https://github.com/cs3org/reva/pull/2422) :  Handle non existing spaces correctly
  * Fix [cs3org/reva#2327](https://github.com/cs3org/reva/pull/2327) :  Enable changelog on edge branch
  * Fix [cs3org/reva#2370](https://github.com/cs3org/reva/pull/2370) :  Fixes for apps in public shares, project spaces for EOS driver
  * Fix [cs3org/reva#2464](https://github.com/cs3org/reva/pull/2464) :  Pass spacegrants when adding member to space
  * Fix [cs3org/reva#2430](https://github.com/cs3org/reva/pull/2430) :  Fix aggregated child folder id
  * Fix [cs3org/reva#2348](https://github.com/cs3org/reva/pull/2348) :  Make archiver handle spaces protocol
  * Fix [cs3org/reva#2452](https://github.com/cs3org/reva/pull/2452) :  Fix create space error message
  * Fix [cs3org/reva#2445](https://github.com/cs3org/reva/pull/2445) :  Don't handle ids containing "/" in decomposedfs
  * Fix [cs3org/reva#2285](https://github.com/cs3org/reva/pull/2285) :  Accept new userid idp format
  * Fix [cs3org/reva#2503](https://github.com/cs3org/reva/pull/2503) :  Remove the protection from /v?.php/config endpoints
  * Fix [cs3org/reva#2462](https://github.com/cs3org/reva/pull/2462) :  Public shares path needs to be set
  * Fix [cs3org/reva#2427](https://github.com/cs3org/reva/pull/2427) :  Fix registry caching
  * Fix [cs3org/reva#2298](https://github.com/cs3org/reva/pull/2298) :  Remove share refs from trashbin
  * Fix [cs3org/reva#2433](https://github.com/cs3org/reva/pull/2433) :  Fix shares provider filter
  * Fix [cs3org/reva#2351](https://github.com/cs3org/reva/pull/2351) :  Fix Statcache removing
  * Fix [cs3org/reva#2374](https://github.com/cs3org/reva/pull/2374) :  Fix webdav copy of zero byte files
  * Fix [cs3org/reva#2336](https://github.com/cs3org/reva/pull/2336) :  Handle sending all permissions when creating public links
  * Fix [cs3org/reva#2440](https://github.com/cs3org/reva/pull/2440) :  Add ArbitraryMetadataKeys to statcache key
  * Fix [cs3org/reva#2582](https://github.com/cs3org/reva/pull/2582) :  Keep lock structs in a local map protected by a mutex
  * Fix [cs3org/reva#2372](https://github.com/cs3org/reva/pull/2372) :  Make owncloudsql work with the spaces registry
  * Fix [cs3org/reva#2416](https://github.com/cs3org/reva/pull/2416) :  The registry now returns complete space structs
  * Fix [cs3org/reva#3066](https://github.com/cs3org/reva/pull/3066) :  Fix propfind listing for files
  * Fix [cs3org/reva#2428](https://github.com/cs3org/reva/pull/2428) :  Remove unused home provider from config
  * Fix [cs3org/reva#2334](https://github.com/cs3org/reva/pull/2334) :  Revert fix decomposedfs upload
  * Fix [cs3org/reva#2415](https://github.com/cs3org/reva/pull/2415) :  Services should never return transport level errors
  * Fix [cs3org/reva#2419](https://github.com/cs3org/reva/pull/2419) :  List project spaces for share recipients
  * Fix [cs3org/reva#2501](https://github.com/cs3org/reva/pull/2501) :  Fix spaces stat
  * Fix [cs3org/reva#2432](https://github.com/cs3org/reva/pull/2432) :  Use space reference when listing containers
  * Fix [cs3org/reva#2572](https://github.com/cs3org/reva/pull/2572) :  Wait for nats server on middleware start
  * Fix [cs3org/reva#2454](https://github.com/cs3org/reva/pull/2454) :  Fix webdav paths in PROPFINDS
  * Chg [cs3org/reva#2329](https://github.com/cs3org/reva/pull/2329) :  Activate the statcache
  * Chg [cs3org/reva#2596](https://github.com/cs3org/reva/pull/2596) :  Remove hash from public link urls
  * Chg [cs3org/reva#2495](https://github.com/cs3org/reva/pull/2495) :  Remove the ownCloud storage driver
  * Chg [cs3org/reva#2527](https://github.com/cs3org/reva/pull/2527) :  Store space attributes in decomposedFS
  * Chg [cs3org/reva#2581](https://github.com/cs3org/reva/pull/2581) :  Update hard-coded status values
  * Chg [cs3org/reva#2524](https://github.com/cs3org/reva/pull/2524) :  Use description during space creation
  * Chg [cs3org/reva#2554](https://github.com/cs3org/reva/pull/2554) :  Shard nodes per space in decomposedfs
  * Chg [cs3org/reva#2576](https://github.com/cs3org/reva/pull/2576) :  Harden xattrs errors
  * Chg [cs3org/reva#2436](https://github.com/cs3org/reva/pull/2436) :  Replace template in GroupFilter for UserProvider with a simple string
  * Chg [cs3org/reva#2429](https://github.com/cs3org/reva/pull/2429) :  Make archiver id based
  * Chg [cs3org/reva#2340](https://github.com/cs3org/reva/pull/2340) :  Allow multiple space configurations per provider
  * Chg [cs3org/reva#2396](https://github.com/cs3org/reva/pull/2396) :  The ocdav handler is now spaces aware
  * Chg [cs3org/reva#2349](https://github.com/cs3org/reva/pull/2349) :  Require `ListRecycle` when listing trashbin
  * Chg [cs3org/reva#2353](https://github.com/cs3org/reva/pull/2353) :  Reduce log output
  * Chg [cs3org/reva#2542](https://github.com/cs3org/reva/pull/2542) :  Do not encode webDAV ids to base64
  * Chg [cs3org/reva#2519](https://github.com/cs3org/reva/pull/2519) :  Remove the auto creation of the .space folder
  * Chg [cs3org/reva#2394](https://github.com/cs3org/reva/pull/2394) :  Remove logic from gateway
  * Chg [cs3org/reva#2023](https://github.com/cs3org/reva/pull/2023) :  Add a sharestorageprovider
  * Chg [cs3org/reva#2234](https://github.com/cs3org/reva/pull/2234) :  Add a spaces registry
  * Chg [cs3org/reva#2339](https://github.com/cs3org/reva/pull/2339) :  Fix static registry regressions
  * Chg [cs3org/reva#2370](https://github.com/cs3org/reva/pull/2370) :  Fix static registry regressions
  * Chg [cs3org/reva#2354](https://github.com/cs3org/reva/pull/2354) :  Return not found when updating non existent space
  * Chg [cs3org/reva#2589](https://github.com/cs3org/reva/pull/2589) :  Remove deprecated linter modules
  * Chg [cs3org/reva#2016](https://github.com/cs3org/reva/pull/2016) :  Move wrapping and unwrapping of paths to the storage gateway
  * Enh [cs3org/reva#2591](https://github.com/cs3org/reva/pull/2591) :  Set up App Locks with basic locks
  * Enh [cs3org/reva#1209](https://github.com/cs3org/reva/pull/1209) :  Reva CephFS module v0.2.1
  * Enh [cs3org/reva#2511](https://github.com/cs3org/reva/pull/2511) :  Error handling cleanup in decomposed FS
  * Enh [cs3org/reva#2516](https://github.com/cs3org/reva/pull/2516) :  Cleaned up some code
  * Enh [cs3org/reva#2512](https://github.com/cs3org/reva/pull/2512) :  Consolidate xattr setter and getter
  * Enh [cs3org/reva#2341](https://github.com/cs3org/reva/pull/2341) :  Use CS3 permissions API
  * Enh [cs3org/reva#2343](https://github.com/cs3org/reva/pull/2343) :  Allow multiple space type fileters on decomposedfs
  * Enh [cs3org/reva#2460](https://github.com/cs3org/reva/pull/2460) :  Add locking support to decomposedfs
  * Enh [cs3org/reva#2540](https://github.com/cs3org/reva/pull/2540) :  Refactored the xattrs package in the decomposedfs
  * Enh [cs3org/reva#2463](https://github.com/cs3org/reva/pull/2463) :  Do not log whole nodes
  * Enh [cs3org/reva#2350](https://github.com/cs3org/reva/pull/2350) :  Add file locking methods to the storage and filesystem interfaces
  * Enh [cs3org/reva#2379](https://github.com/cs3org/reva/pull/2379) :  Add new file url of the app provider to the ocs capabilities
  * Enh [cs3org/reva#2369](https://github.com/cs3org/reva/pull/2369) :  Implement TouchFile from the CS3apis
  * Enh [cs3org/reva#2385](https://github.com/cs3org/reva/pull/2385) :  Allow to create new files with the app provider on public links
  * Enh [cs3org/reva#2397](https://github.com/cs3org/reva/pull/2397) :  Product field in OCS version
  * Enh [cs3org/reva#2393](https://github.com/cs3org/reva/pull/2393) :  Update tus/tusd to version 1.8.0
  * Enh [cs3org/reva#2522](https://github.com/cs3org/reva/pull/2522) :  Introduce events
  * Enh [cs3org/reva#2528](https://github.com/cs3org/reva/pull/2528) :  Use an exclusive write lock when writing multiple attributes
  * Enh [cs3org/reva#2595](https://github.com/cs3org/reva/pull/2595) :  Add integration test for the groupprovider
  * Enh [cs3org/reva#2439](https://github.com/cs3org/reva/pull/2439) :  Ignore handled errors when creating spaces
  * Enh [cs3org/reva#2500](https://github.com/cs3org/reva/pull/2500) :  Invalidate listproviders cache
  * Enh [cs3org/reva#2345](https://github.com/cs3org/reva/pull/2345) :  Don't assume that the LDAP groupid in reva matches the name
  * Enh [cs3org/reva#2525](https://github.com/cs3org/reva/pull/2525) :  Allow using AD UUID as userId values
  * Enh [cs3org/reva#2584](https://github.com/cs3org/reva/pull/2584) :  Allow running userprovider integration tests for the LDAP driver
  * Enh [cs3org/reva#2585](https://github.com/cs3org/reva/pull/2585) :  Add metadata storage layer and indexer
  * Enh [cs3org/reva#2163](https://github.com/cs3org/reva/pull/2163) :  Nextcloud-based share manager for pkg/ocm/share
  * Enh [cs3org/reva#2278](https://github.com/cs3org/reva/pull/2278) :  OIDC driver changes for lightweight users
  * Enh [cs3org/reva#2315](https://github.com/cs3org/reva/pull/2315) :  Add new attributes to public link propfinds
  * Enh [cs3org/reva#2431](https://github.com/cs3org/reva/pull/2431) :  Delete shares when purging spaces
  * Enh [cs3org/reva#2434](https://github.com/cs3org/reva/pull/2434) :  Refactor ocdav into smaller chunks
  * Enh [cs3org/reva#2524](https://github.com/cs3org/reva/pull/2524) :  Add checks when removing space members
  * Enh [cs3org/reva#2457](https://github.com/cs3org/reva/pull/2457) :  Restore spaces that were previously deleted
  * Enh [cs3org/reva#2498](https://github.com/cs3org/reva/pull/2498) :  Include grants in list storage spaces response
  * Enh [cs3org/reva#2344](https://github.com/cs3org/reva/pull/2344) :  Allow listing all storage spaces
  * Enh [cs3org/reva#2547](https://github.com/cs3org/reva/pull/2547) :  Add an if-match check to the storage provider
  * Enh [cs3org/reva#2486](https://github.com/cs3org/reva/pull/2486) :  Update cs3apis to include lock api changes
  * Enh [cs3org/reva#2526](https://github.com/cs3org/reva/pull/2526) :  Upgrade ginkgo to v2

   https://github.com/owncloud/ocis/pull/3231
   https://github.com/owncloud/ocis/pull/3258

* Enhancement - Update ownCloud Web to v5.2.0: [#6506](https://github.com/owncloud/web/pull/6506)

   Tags: web

   We updated ownCloud Web to v5.2.0. Please refer to the changelog (linked) for
   details on the web release.

   https://github.com/owncloud/web/pull/6506
   https://github.com/owncloud/ocis/pull/3202
   https://github.com/owncloud/web/releases/tag/v5.2.0

# Changelog for [1.17.0] (2022-02-16)

The following sections list the changes for 1.17.0.

[1.17.0]: https://github.com/owncloud/ocis/compare/v1.16.0...v1.17.0

## Summary

* Bugfix - Fix configuration for space membership endpoint: [#2893](https://github.com/owncloud/ocis/pull/2893)
* Bugfix - Add `ocis storage-auth-machine` subcommand: [#2910](https://github.com/owncloud/ocis/pull/2910)
* Bugfix - Fix the default tracing provider: [#2952](https://github.com/owncloud/ocis/pull/2952)
* Bugfix - Fix retry handling for LDAP connections: [#2974](https://github.com/owncloud/ocis/issues/2974)
* Bugfix - Remove group memberships when deleting a user: [#3027](https://github.com/owncloud/ocis/issues/3027)
* Bugfix - Make the default grpc client use the registry settings: [#3041](https://github.com/owncloud/ocis/pull/3041)
* Bugfix - Use same jwt secret for accounts as for metadata storage: [#3081](https://github.com/owncloud/ocis/pull/3081)
* Change - Unify configuration and commands: [#2818](https://github.com/owncloud/ocis/pull/2818)
* Change - Update libre-graph-api to v0.3.0: [#2858](https://github.com/owncloud/ocis/pull/2858)
* Change - Return not found when updating non existent space: [#2869](https://github.com/cs3org/reva/pull/2869)
* Change - Update the graph api: [#2885](https://github.com/owncloud/ocis/pull/2885)
* Change - Change log level default from debug to error: [#3071](https://github.com/owncloud/ocis/pull/3071)
* Change - Remove the ownCloud storage driver: [#3072](https://github.com/owncloud/ocis/pull/3072)
* Change - Functionality to restore spaces: [#3092](https://github.com/owncloud/ocis/pull/3092)
* Change - Extended Space Properties: [#3141](https://github.com/owncloud/ocis/pull/3141)
* Enhancement - Support signature auth in the public share auth middleware: [#2831](https://github.com/owncloud/ocis/pull/2831)
* Enhancement - Update REVA to v1.16.1-0.20220215130802-df1264deff58: [#2878](https://github.com/owncloud/ocis/pull/2878)
* Enhancement - Add new file url of the app provider to the ocs capabilities: [#2884](https://github.com/owncloud/ocis/pull/2884)
* Enhancement - Update ownCloud Web to v5.0.0: [#2895](https://github.com/owncloud/ocis/pull/2895)
* Enhancement - Add spaces capability: [#2931](https://github.com/owncloud/ocis/pull/2931)
* Enhancement - Add filter by driveType and id to /me/drives: [#2946](https://github.com/owncloud/ocis/pull/2946)
* Enhancement - Introduce User and Group Management capabilities on GraphAPI: [#2947](https://github.com/owncloud/ocis/pull/2947)
* Enhancement - Update REVA to v1.16.1-0.20220112085026-07451f6cd806: [#2953](https://github.com/owncloud/ocis/pull/2953)
* Enhancement - Add endpoint to retrieve a single space: [#2978](https://github.com/owncloud/ocis/pull/2978)
* Enhancement - Add graph endpoint to delete and purge spaces: [#2979](https://github.com/owncloud/ocis/pull/2979)
* Enhancement - Add permissions to graph drives: [#3095](https://github.com/owncloud/ocis/pull/3095)
* Enhancement - Consul as supported service registry: [#3133](https://github.com/owncloud/ocis/pull/3133)
* Enhancement - Provide Description when creating a space: [#3167](https://github.com/owncloud/ocis/pull/3167)

## Details

* Bugfix - Fix configuration for space membership endpoint: [#2893](https://github.com/owncloud/ocis/pull/2893)

   Added a missing config value to the ocs config related to the space membership
   endpoint.

   https://github.com/owncloud/ocis/pull/2893

* Bugfix - Add `ocis storage-auth-machine` subcommand: [#2910](https://github.com/owncloud/ocis/pull/2910)

   We added the ocis subcommand to start the machine auth provider.

   https://github.com/owncloud/ocis/pull/2910

* Bugfix - Fix the default tracing provider: [#2952](https://github.com/owncloud/ocis/pull/2952)

   We've fixed the default tracing provider which was no longer configured after
   [owncloud/ocis#2818](https://github.com/owncloud/ocis/pull/2818).

   https://github.com/owncloud/ocis/pull/2952
   https://github.com/owncloud/ocis/pull/2818

* Bugfix - Fix retry handling for LDAP connections: [#2974](https://github.com/owncloud/ocis/issues/2974)

   We've fixed the handling of network issues (e.g. connection loss) during LDAP
   Write Operations to correctly retry the request.

   https://github.com/owncloud/ocis/issues/2974

* Bugfix - Remove group memberships when deleting a user: [#3027](https://github.com/owncloud/ocis/issues/3027)

   The LDAP backend in the graph API now takes care of removing a user's group
   membership when deleting the user.

   https://github.com/owncloud/ocis/issues/3027

* Bugfix - Make the default grpc client use the registry settings: [#3041](https://github.com/owncloud/ocis/pull/3041)

   We've fixed the default grpc client to use the registry settings. Previously it
   always used mdns.

   https://github.com/owncloud/ocis/pull/3041

* Bugfix - Use same jwt secret for accounts as for metadata storage: [#3081](https://github.com/owncloud/ocis/pull/3081)

   We've the metadata storage uses the same jwt secret as all other REVA services.
   Therefore the accounts service needs to use the same secret.

   Secrets are documented here:
   https://owncloud.dev/ocis/deployment/#change-default-secrets

   https://github.com/owncloud/ocis/pull/3081

* Change - Unify configuration and commands: [#2818](https://github.com/owncloud/ocis/pull/2818)

   We've unified the configuration and commands of all non storage services. This
   also includes the change, that environment variables are now defined on the
   config struct as tags instead in a separate mapping.

   https://github.com/owncloud/ocis/pull/2818

* Change - Update libre-graph-api to v0.3.0: [#2858](https://github.com/owncloud/ocis/pull/2858)

   This updates the libre-graph-api to use the latest spec and types.

   https://github.com/owncloud/ocis/pull/2858

* Change - Return not found when updating non existent space: [#2869](https://github.com/cs3org/reva/pull/2869)

   If a spaceid of a space which is updated doesn't exist, handle it as a not found
   error.

   https://github.com/cs3org/reva/pull/2869

* Change - Update the graph api: [#2885](https://github.com/owncloud/ocis/pull/2885)

   GraphApi has been updated to version 0.4.1 and the existing dependency was
   removed

   https://github.com/owncloud/ocis/pull/2885

* Change - Change log level default from debug to error: [#3071](https://github.com/owncloud/ocis/pull/3071)

   We've changed the default log level for all services from "info" to "error".

   https://github.com/owncloud/ocis/pull/3071

* Change - Remove the ownCloud storage driver: [#3072](https://github.com/owncloud/ocis/pull/3072)

   We've removed the ownCloud storage driver because it was no longer maintained
   after the ownCloud SQL storage driver was added.

   If you have been using the ownCloud storage driver, please switch to the
   ownCloud SQL storage driver which brings you more features and is under active
   maintenance.

   https://github.com/owncloud/ocis/pull/3072

* Change - Functionality to restore spaces: [#3092](https://github.com/owncloud/ocis/pull/3092)

   Disabled spaces can now be restored via the graph api. An information was added
   to the root item of each space when it is deleted

   https://github.com/owncloud/ocis/pull/3092

* Change - Extended Space Properties: [#3141](https://github.com/owncloud/ocis/pull/3141)

   We can now set and modify short description, space image and space readme. Only
   managers can set the short description. Editors can change the space image and
   readme id.

   https://github.com/owncloud/ocis/pull/3141

* Enhancement - Support signature auth in the public share auth middleware: [#2831](https://github.com/owncloud/ocis/pull/2831)

   Enabled public share requests to be authenticated using the public share
   signature.

   https://github.com/owncloud/ocis/pull/2831

* Enhancement - Update REVA to v1.16.1-0.20220215130802-df1264deff58: [#2878](https://github.com/owncloud/ocis/pull/2878)

   Updated REVA to v1.16.1-0.20220215130802-df1264deff58 This update includes:

  * Enh [cs3org/reva#2524](https://github.com/cs3org/reva/pull/2524): Remove space members
  * Fix [cs3org/reva#2541](https://github.com/cs3org/reva/pull/2541): fix xattr error types, remove error wrapper
  * Chg [cs3org/reva#2540](https://github.com/cs3org/reva/pull/2540): decomposedfs: refactor xattrs package errors
  * Enh [cs3org/reva#2533](https://github.com/cs3org/reva/pull/2533): Use space description on creation
  * Enh [cs3org/reva#2527](https://github.com/cs3org/reva/pull/2527): Add space props
  * Enh [cs3org/reva#2522](https://github.com/cs3org/reva/pull/2522): Events
  * Chg [cs3org/reva#2512](https://github.com/cs3org/reva/pull/2512): Consolidate all metadata Get's and Set's to central functions.
  * Chg [cs3org/reva#2511](https://github.com/cs3org/reva/pull/2511): Some error cleanup steps in the decomposed FS
  * Enh [cs3org/reva#2460](https://github.com/cs3org/reva/pull/2460): decomposedfs: add locking support
  * Chg [cs3org/reva#2519](https://github.com/cs3org/reva/pull/2519): remove creation of .space folder
  * Fix [cs3org/reva#2506](https://github.com/cs3org/reva/pull/2506): fix propfind listing for files
  * Chg [cs3org/reva#2503](https://github.com/cs3org/reva/pull/2503): unprotected ocs config endpoint
  * Enh [cs3org/reva#2458](https://github.com/cs3org/reva/pull/2458): Restoring Spaces
  * Enh [cs3org/reva#2498](https://github.com/cs3org/reva/pull/2498): add grants to list-spaces
  * Fix [cs3org/reva#2500](https://github.com/cs3org/reva/pull/2500): invalidate cache when modifying or deleting a space
  * Fix [cs3org/reva#2501](https://github.com/cs3org/reva/pull/2501): fix spaces stat requests
  * Enh [cs3org/reva#2472](https://github.com/cs3org/reva/pull/2472): Make owncloudsql spaces aware
  * Enh [cs3org/reva#2464](https://github.com/cs3org/reva/pull/2464): Space grants
  * Fix [cs3org/reva#2463](https://github.com/cs3org/reva/pull/2463): Do not log nodes
  * Enh [cs3org/reva#2437](https://github.com/cs3org/reva/pull/2437): Make gateway dumb again
  * Enh [cs3org/reva#2459](https://github.com/cs3org/reva/pull/2459): prevent purging of enabled spaces
  * Fix [cs3org/reva#2457](https://github.com/cs3org/reva/pull/2457): decomposedfs: do not swallow errors when creating nodes
  * Fix [cs3org/reva#2454](https://github.com/cs3org/reva/pull/2454): fix path construction in webdav propfind
  * Fix [cs3org/reva#2452](https://github.com/cs3org/reva/pull/2452): fix create space error message
  * Enh [cs3org/reva#2431](https://github.com/cs3org/reva/pull/2431): Purge spaces
  * Fix [cs3org/reva#2445](https://github.com/cs3org/reva/pull/2445): Fix publiclinks and decomposedfs
  * Chg [cs3org/reva#2439](https://github.com/cs3org/reva/pull/2439): ignore handled errors when creating spaces
  * Enh [cs3org/reva#2436](https://github.com/cs3org/reva/pull/2436): Adjust "groupfilter" to be able to search by member name
  * Fix [cs3org/reva#2434](https://github.com/cs3org/reva/pull/2434): Start splitting up ocdav
  * Fix [cs3org/reva#2433](https://github.com/cs3org/reva/pull/2433): fix shares provider filter
  * Chg [cs3org/reva#2432](https://github.com/cs3org/reva/pull/2432): use space reference when listing containers
  * Fix [cs3org/reva#2430](https://github.com/cs3org/reva/pull/2430): fix aggregated child folder id
  * Enh [cs3org/reva#2429](https://github.com/cs3org/reva/pull/2429): make archiver id based
  * Fix [cs3org/reva#2427](https://github.com/cs3org/reva/pull/2427): fix registry caching
  * Fix [cs3org/reva#2422](https://github.com/cs3org/reva/pull/2422): handle space does not exist
  * Fix [cs3org/reva#2419](https://github.com/cs3org/reva/pull/2419): Spaces fixes
  * Chg [cs3org/reva#2415](https://github.com/cs3org/reva/pull/2415): services should never return transport level errors
  * Chg [cs3org/reva#2396](https://github.com/cs3org/reva/pull/2396): Ocdav spaces aware
  * Fix [cs3org/reva#2348](https://github.com/cs3org/reva/pull/2348): fix-archiver
  * Chg [cs3org/reva#2344](https://github.com/cs3org/reva/pull/2344): allow listing all storage spaces
  * Chg [cs3org/reva#2345](https://github.com/cs3org/reva/pull/2345): Switch LDAP test to use entryUUID as unique id for groups
  * Chg [cs3org/reva#2343](https://github.com/cs3org/reva/pull/2343): allow multiple space type filters on decomposedfs
  * Enh [cs3org/reva#2329](https://github.com/cs3org/reva/pull/2329): Activate Statcache
  * Enh [cs3org/reva#2340](https://github.com/cs3org/reva/pull/2340): Space registry multiple spaces per provider
  * Chg [cs3org/reva#2336](https://github.com/cs3org/reva/pull/2336): handle sending all permissions when creating public links
  * Fix [cs3org/reva#2330](https://github.com/cs3org/reva/pull/2330): fix decomposedfs upload
  * Enh [cs3org/reva#2234](https://github.com/cs3org/reva/pull/2234): Spaces registry
  * Enh [cs3org/reva#2217](https://github.com/cs3org/reva/pull/2217): New OIDC ESCAPE auth driver.
  * Enh [cs3org/reva#2250](https://github.com/cs3org/reva/pull/2250): Implement space membership endpoints
  * Fix [cs3org/reva#1941](https://github.com/cs3org/reva/pull/1941): fix tus with transfer token only
  * Fix [cs3org/reva#2309](https://github.com/cs3org/reva/pull/2309): Bugfix: Remove early finish for zero byte file uploads
  * Fix [cs3org/reva#2303](https://github.com/cs3org/reva/pull/2303): Fix content disposition
  * Fix [cs3org/reva#2314](https://github.com/cs3org/reva/pull/2314): OIDC: fallback to "email" if IDP doesn't provide "preferred_username" claim
  * Enh [cs3org/reva#2256](https://github.com/cs3org/reva/pull/2256): Return user type in the response of the ocs GET user call
  * Enh [cs3org/reva#2310](https://github.com/cs3org/reva/pull/2310): Implement setting arbitrary metadata for the public storage provider
  * Fix [cs3org/reva#2305](https://github.com/cs3org/reva/pull/2305): Make sure /app/new takes target as absolute path
  * Fix [cs3org/reva#2297](https://github.com/cs3org/reva/pull/2297): Fix public link paths for file shares

   https://github.com/owncloud/ocis/pull/2878
   https://github.com/owncloud/ocis/pull/2901
   https://github.com/owncloud/ocis/pull/2997
   https://github.com/owncloud/ocis/pull/3116
   https://github.com/owncloud/ocis/pull/3130
   https://github.com/owncloud/ocis/pull/3175
   https://github.com/owncloud/ocis/pull/3182

* Enhancement - Add new file url of the app provider to the ocs capabilities: [#2884](https://github.com/owncloud/ocis/pull/2884)

   We've added the new file capability of the app provider to the ocs capabilities,
   so that clients can discover this url analogous to the app list and file open
   urls.

   https://github.com/owncloud/ocis/pull/2884
   https://github.com/owncloud/ocis/pull/2907
   https://github.com/cs3org/reva/pull/2379
   https://github.com/owncloud/web/pull/5890#issuecomment-993905242

* Enhancement - Update ownCloud Web to v5.0.0: [#2895](https://github.com/owncloud/ocis/pull/2895)

   Tags: web

   We updated ownCloud Web to v5.0.0. Please refer to the changelog (linked) for
   details on the web release.

   https://github.com/owncloud/ocis/pull/2895
   https://github.com/owncloud/ocis/pull/3157
   https://github.com/owncloud/web/releases/tag/v4.8.0
   https://github.com/owncloud/web/releases/tag/v5.0.0

* Enhancement - Add spaces capability: [#2931](https://github.com/owncloud/ocis/pull/2931)

   We've added the spaces capability with version 0.0.1 and enabled defaulting to
   true.

   https://github.com/owncloud/ocis/pull/2931
   https://github.com/cs3org/reva/pull/2015
   https://github.com/owncloud/ocis/pull/2965

* Enhancement - Add filter by driveType and id to /me/drives: [#2946](https://github.com/owncloud/ocis/pull/2946)

   We added two possible filter terms (driveType, id) to the /me/drives endpoint on
   the graph api. These can be used with the odata query parameter "$filter". We
   only support the "eq" operator for now.

   https://github.com/owncloud/ocis/pull/2946

* Enhancement - Introduce User and Group Management capabilities on GraphAPI: [#2947](https://github.com/owncloud/ocis/pull/2947)

   The GraphAPI LDAP Backend is now able to add/modify and delete Users and Groups

   https://github.com/owncloud/ocis/pull/2947
   https://github.com/owncloud/ocis/pull/2996

* Enhancement - Update REVA to v1.16.1-0.20220112085026-07451f6cd806: [#2953](https://github.com/owncloud/ocis/pull/2953)

   Update REVA to v1.16.1-0.20220112085026-07451f6cd806

   https://github.com/owncloud/ocis/pull/2953

* Enhancement - Add endpoint to retrieve a single space: [#2978](https://github.com/owncloud/ocis/pull/2978)

   We added the endpoint ``/drives/{driveID}`` to get a single space by id from the
   server.

   https://github.com/owncloud/ocis/pull/2978

* Enhancement - Add graph endpoint to delete and purge spaces: [#2979](https://github.com/owncloud/ocis/pull/2979)

   Added a new graph endpoint to delete and purge spaces.

   https://github.com/owncloud/ocis/pull/2979
   https://github.com/owncloud/ocis/pull/3000

* Enhancement - Add permissions to graph drives: [#3095](https://github.com/owncloud/ocis/pull/3095)

   Added permissions to graph drives when listing drives.

   https://github.com/owncloud/ocis/pull/3095

* Enhancement - Consul as supported service registry: [#3133](https://github.com/owncloud/ocis/pull/3133)

   We have added Consul as an supported service registry. You can now use it to let
   oCIS services discover each other.

   https://github.com/owncloud/ocis/pull/3133

* Enhancement - Provide Description when creating a space: [#3167](https://github.com/owncloud/ocis/pull/3167)

   We added the possibility to send a short description when creating a space.

   https://github.com/owncloud/ocis/pull/3167

# Changelog for [1.16.0] (2021-12-10)

The following sections list the changes for 1.16.0.

[1.16.0]: https://github.com/owncloud/ocis/compare/v1.15.0...v1.16.0

## Summary

* Bugfix - Fix claim selector based routing for basic auth: [#2779](https://github.com/owncloud/ocis/pull/2779)
* Bugfix - Fix using s3ng as the metadata storage backend: [#2807](https://github.com/owncloud/ocis/pull/2807)
* Bugfix - Disallow creation of a group with empty name via the OCS api: [#2825](https://github.com/owncloud/ocis/pull/2825)
* Bugfix - Use the CS3api up- and download workflow for the accounts service: [#2837](https://github.com/owncloud/ocis/pull/2837)
* Change - OIDC: fallback if IDP doesn't provide "preferred_username" claim: [#2644](https://github.com/owncloud/ocis/issues/2644)
* Change - Restructure Configuration Parsing: [#2708](https://github.com/owncloud/ocis/pull/2708)
* Change - Rename `APP_PROVIDER_BASIC_*` environment variables: [#2812](https://github.com/owncloud/ocis/pull/2812)
* Enhancement - Cleanup ocis-pkg config: [#2813](https://github.com/owncloud/ocis/pull/2813)
* Enhancement - Correct shutdown of services under runtime: [#2843](https://github.com/owncloud/ocis/pull/2843)
* Enhancement - Update ownCloud Web to v4.6.1: [#2846](https://github.com/owncloud/ocis/pull/2846)
* Enhancement - Update REVA to v1.17.0: [#2849](https://github.com/owncloud/ocis/pull/2849)

## Details

* Bugfix - Fix claim selector based routing for basic auth: [#2779](https://github.com/owncloud/ocis/pull/2779)

   We've fixed the claim selector based routing for requests using basic auth.
   Previously requests using basic auth have always been routed to the
   DefaultPolicy when using the claim selector despite the set cookie because the
   basic auth middleware fakes some OIDC claims.

   Now the cookie is checked before routing to the DefaultPolicy and therefore set
   cookie will also be respected for requests using basic auth.

   https://github.com/owncloud/ocis/pull/2779

* Bugfix - Fix using s3ng as the metadata storage backend: [#2807](https://github.com/owncloud/ocis/pull/2807)

   It is now possible to use s3ng as the metadata storage backend.

   https://github.com/owncloud/ocis/issues/2668
   https://github.com/owncloud/ocis/pull/2807

* Bugfix - Disallow creation of a group with empty name via the OCS api: [#2825](https://github.com/owncloud/ocis/pull/2825)

   We've fixed the behavior for group creation on the OCS api, where it was
   possible to create a group with an empty name. This was is not possible on oC10
   and is therefore also forbidden on oCIS to keep compatibility. This PR forbids
   the creation and also ensures the correct status code for both OCS v1 and OCS v2
   apis.

   https://github.com/owncloud/ocis/issues/2823
   https://github.com/owncloud/ocis/pull/2825

* Bugfix - Use the CS3api up- and download workflow for the accounts service: [#2837](https://github.com/owncloud/ocis/pull/2837)

   We've fixed the interaction of the accounts service with the metadata storage
   after bypassing the InitiateUpload and InitiateDownload have been removed from
   various storage drivers. The accounts service now uses the proper CS3apis
   workflow for up- and downloads.

   https://github.com/owncloud/ocis/pull/2837
   https://github.com/cs3org/reva/pull/2309

* Change - OIDC: fallback if IDP doesn't provide "preferred_username" claim: [#2644](https://github.com/owncloud/ocis/issues/2644)

   Some IDPs don't add the "preferred_username" claim. Fallback to the "email"
   claim in that case

   https://github.com/owncloud/ocis/issues/2644

* Change - Restructure Configuration Parsing: [#2708](https://github.com/owncloud/ocis/pull/2708)

   Tags: ocis

   CLI flags are no longer needed for subcommands, as we rely solely on env
   variables and config files. This greatly simplifies configuration and
   deployment.

   https://github.com/owncloud/ocis/pull/2708

* Change - Rename `APP_PROVIDER_BASIC_*` environment variables: [#2812](https://github.com/owncloud/ocis/pull/2812)

   We've renamed the `APP_PROVIDER_BASIC_*` to `APP_PROVIDER_*` since the `_BASIC_`
   part is a copy and paste error. Now all app provider environment variables are
   consistently starting with `APP_PROVIDER_*`.

   https://github.com/owncloud/ocis/pull/2812
   https://github.com/owncloud/ocis/pull/2811

* Enhancement - Cleanup ocis-pkg config: [#2813](https://github.com/owncloud/ocis/pull/2813)

   Certain values were of no use when configuring the ocis runtime.

   https://github.com/owncloud/ocis/pull/2813

* Enhancement - Correct shutdown of services under runtime: [#2843](https://github.com/owncloud/ocis/pull/2843)

   Supervised goroutines now shut themselves down on context cancellation
   propagation.

   https://github.com/owncloud/ocis/pull/2843

* Enhancement - Update ownCloud Web to v4.6.1: [#2846](https://github.com/owncloud/ocis/pull/2846)

   Tags: web

   We updated ownCloud Web to v4.6.1. Please refer to the changelog (linked) for
   details on the web release.

   https://github.com/owncloud/ocis/pull/2846
   https://github.com/owncloud/web/releases/tag/v4.6.1

* Enhancement - Update REVA to v1.17.0: [#2849](https://github.com/owncloud/ocis/pull/2849)

   Updated REVA to v1.17.0 This update includes:

  * Fix [cs3org/reva#2305](https://github.com/cs3org/reva/pull/2305): Make sure /app/new takes `target` as absolute path
  * Fix [cs3org/reva#2303](https://github.com/cs3org/reva/pull/2303): Fix content disposition header for public links files
  * Fix [cs3org/reva#2316](https://github.com/cs3org/reva/pull/2316): Fix the share types in propfinds
  * Fix [cs3org/reva#2803](https://github.com/cs3org/reva/pull/2310): Fix app provider for editor public links
  * Fix [cs3org/reva#2298](https://github.com/cs3org/reva/pull/2298): Remove share refs from trashbin
  * Fix [cs3org/reva#2309](https://github.com/cs3org/reva/pull/2309): Remove early finish for zero byte file uploads
  * Fix [cs3org/reva#1941](https://github.com/cs3org/reva/pull/1941): Fix TUS uploads with transfer token only
  * Chg [cs3org/reva#2210](https://github.com/cs3org/reva/pull/2210): Fix app provider new file creation and improved error codes
  * Enh [cs3org/reva#2217](https://github.com/cs3org/reva/pull/2217): OIDC auth driver for ESCAPE IAM
  * Enh [cs3org/reva#2256](https://github.com/cs3org/reva/pull/2256): Return user type in the response of the ocs GET user call
  * Enh [cs3org/reva#2315](https://github.com/cs3org/reva/pull/2315): Add new attributes to public link propfinds
  * Enh [cs3org/reva#2740](https://github.com/cs3org/reva/pull/2250): Implement space membership endpoints
  * Enh [cs3org/reva#2252](https://github.com/cs3org/reva/pull/2252): Add the xattr sys.acl to SysACL (eosgrpc)
  * Enh [cs3org/reva#2314](https://github.com/cs3org/reva/pull/2314): OIDC: fallback if IDP doesn't provide "preferred_username" claim

   https://github.com/owncloud/ocis/pull/2849
   https://github.com/owncloud/ocis/pull/2835
   https://github.com/owncloud/ocis/pull/2837

# Changelog for [1.15.0] (2021-11-19)

The following sections list the changes for 1.15.0.

[1.15.0]: https://github.com/owncloud/ocis/compare/v1.14.0...v1.15.0

## Summary

* Bugfix - Don't allow empty password: [#197](https://github.com/owncloud/product/issues/197)
* Bugfix - Don't announce resharing via capabilities: [#2690](https://github.com/owncloud/ocis/pull/2690)
* Bugfix - Fix oCIS startup ony systems with IPv6: [#2698](https://github.com/owncloud/ocis/pull/2698)
* Bugfix - Fix error logging when there is no thumbnail for a file: [#2702](https://github.com/owncloud/ocis/pull/2702)
* Bugfix - Fix basic auth config: [#2719](https://github.com/owncloud/ocis/pull/2719)
* Bugfix - Fix opening images in media viewer for some usernames: [#2738](https://github.com/owncloud/ocis/pull/2738)
* Bugfix - Fix basic auth with custom user claim: [#2755](https://github.com/owncloud/ocis/pull/2755)
* Change - Make all insecure options configurable and change the default to false: [#2700](https://github.com/owncloud/ocis/issues/2700)
* Change - Update ownCloud Web to v4.5.0: [#2780](https://github.com/owncloud/ocis/pull/2780)
* Enhancement - Add API to list all spaces: [#2692](https://github.com/owncloud/ocis/pull/2692)
* Enhancement - Update REVA to v1.16.0: [#2737](https://github.com/owncloud/ocis/pull/2737)

## Details

* Bugfix - Don't allow empty password: [#197](https://github.com/owncloud/product/issues/197)

   It was allowed to create users with empty or spaces-only password. This is fixed

   https://github.com/owncloud/product/issues/197

* Bugfix - Don't announce resharing via capabilities: [#2690](https://github.com/owncloud/ocis/pull/2690)

   OCIS / Reva is not capable of resharing, yet. We've set the resharing capability
   to false, so that clients have a chance to react accordingly.

   https://github.com/owncloud/ocis/pull/2690

* Bugfix - Fix oCIS startup ony systems with IPv6: [#2698](https://github.com/owncloud/ocis/pull/2698)

   We've fixed failing startup of oCIS on systems with IPv6 addresses.

   https://github.com/owncloud/ocis/issues/2300
   https://github.com/owncloud/ocis/pull/2698

* Bugfix - Fix error logging when there is no thumbnail for a file: [#2702](https://github.com/owncloud/ocis/pull/2702)

   We've fixed the behavior of the logging when there is no thumbnail for a file
   (because the filetype is not supported for thumbnail generation). Previously the
   WebDAV service always issues an error log in this case. Now, we don't log this
   event any more.

   https://github.com/owncloud/ocis/pull/2702

* Bugfix - Fix basic auth config: [#2719](https://github.com/owncloud/ocis/pull/2719)

   Users could authenticate using basic auth even though `PROXY_ENABLE_BASIC_AUTH`
   was set to false.

   https://github.com/owncloud/ocis/issues/2466
   https://github.com/owncloud/ocis/pull/2719

* Bugfix - Fix opening images in media viewer for some usernames: [#2738](https://github.com/owncloud/ocis/pull/2738)

   We've fixed the opening of images in the media viewer for user names containing
   special characters (eg. `@`) which will be URL-escaped. Before this fix users
   could not see the image in the media viewer. Now the user name is correctly
   escaped and the user can view the image in the media viewer.

   https://github.com/owncloud/ocis/pull/2738

* Bugfix - Fix basic auth with custom user claim: [#2755](https://github.com/owncloud/ocis/pull/2755)

   We've fixed authentication with basic if oCIS is configured to use a
   non-standard claim as user claim (`PROXY_USER_OIDC_CLAIM`). Prior to this bugfix
   the authentication always failed and is now working.

   https://github.com/owncloud/ocis/pull/2755

* Change - Make all insecure options configurable and change the default to false: [#2700](https://github.com/owncloud/ocis/issues/2700)

   We had several hard-coded 'insecure' flags. These options are now configurable
   and default to false. Also we changed all other 'insecure' flags with a previous
   default of true to false.

   In development environments using self signed certs (the default) you now need
   to set these flags:

   ```
   PROXY_OIDC_INSECURE=true
   STORAGE_FRONTEND_APPPROVIDER_INSECURE=true
   STORAGE_FRONTEND_ARCHIVER_INSECURE=true
   STORAGE_FRONTEND_OCDAV_INSECURE=true
   STORAGE_HOME_DATAPROVIDER_INSECURE=true
   STORAGE_METADATA_DATAPROVIDER_INSECURE=true
   STORAGE_OIDC_INSECURE=true
   STORAGE_USERS_DATAPROVIDER_INSECURE=true
   THUMBNAILS_CS3SOURCE_INSECURE=true
   THUMBNAILS_WEBDAVSOURCE_INSECURE=true
   ```

   As an alternative you also can set a single flag, which configures all options
   together:

   ```
   OCIS_INSECURE=true
   ```

   https://github.com/owncloud/ocis/issues/2700
   https://github.com/owncloud/ocis/pull/2745

* Change - Update ownCloud Web to v4.5.0: [#2780](https://github.com/owncloud/ocis/pull/2780)

   Tags: web

   We updated ownCloud Web to v4.5.0. Please refer to the changelog (linked) for
   details on the web release.

   https://github.com/owncloud/ocis/pull/2780
   https://github.com/owncloud/web/releases/tag/v4.5.0

* Enhancement - Add API to list all spaces: [#2692](https://github.com/owncloud/ocis/pull/2692)

   Added a graph endpoint to enable users with the `list-all-spaces` permission to
   list all spaces.

   https://github.com/owncloud/ocis/pull/2692

* Enhancement - Update REVA to v1.16.0: [#2737](https://github.com/owncloud/ocis/pull/2737)

   Updated REVA to v1.16.0 This update includes:

  * Fix [cs3org/reva#2245](https://github.com/cs3org/reva/pull/2245): Don't announce search-files capability
  * Fix [cs3org/reva#2247](https://github.com/cs3org/reva/pull/2247): Merge user ACLs from EOS to sys ACLs
  * Fix [cs3org/reva#2279](https://github.com/cs3org/reva/pull/2279): Return the inode of the version folder for files when listing in EOS
  * Fix [cs3org/reva#2294](https://github.com/cs3org/reva/pull/2294): Fix HTTP return code when path is invalid
  * Fix [cs3org/reva#2231](https://github.com/cs3org/reva/pull/2231): Fix share permission on a single file in sql share driver (cbox pkg)
  * Fix [cs3org/reva#2230](https://github.com/cs3org/reva/pull/2230): Fix open by default app and expose default app
  * Fix [cs3org/reva#2265](https://github.com/cs3org/reva/pull/2265): Fix nil pointer exception when resolving members of a group (rest driver)
  * Fix [cs3org/reva#1214](https://github.com/cs3org/reva/pull/1214): Fix restoring versions
  * Fix [cs3org/reva#2254](https://github.com/cs3org/reva/pull/2254): Fix spaces propfind
  * Fix [cs3org/reva#2260](https://github.com/cs3org/reva/pull/2260): Fix unset quota xattr on darwin
  * Fix [cs3org/reva#5776](https://github.com/cs3org/reva/pull/5776): Enforce permissions in public share apps
  * Fix [cs3org/reva#2767](https://github.com/cs3org/reva/pull/2767): Fix status code for WebDAV mkcol requests where an ancestor is missing
  * Fix [cs3org/reva#2287](https://github.com/cs3org/reva/pull/2287): Add public link access via mount-ID:token/relative-path to the scope
  * Fix [cs3org/reva#2244](https://github.com/cs3org/reva/pull/2244): Fix the permissions response for shared files in the cbox sql driver
  * Enh [cs3org/reva#2219](https://github.com/cs3org/reva/pull/2219): Add virtual view tests
  * Enh [cs3org/reva#2230](https://github.com/cs3org/reva/pull/2230): Add priority to app providers
  * Enh [cs3org/reva#2258](https://github.com/cs3org/reva/pull/2258): Improved error messages from the AppProviders
  * Enh [cs3org/reva#2119](https://github.com/cs3org/reva/pull/2119): Add authprovider owncloudsql
  * Enh [cs3org/reva#2211](https://github.com/cs3org/reva/pull/2211): Enhance the cbox share sql driver to store accepted group shares
  * Enh [cs3org/reva#2212](https://github.com/cs3org/reva/pull/2212): Filter root path according to the agent that makes the request
  * Enh [cs3org/reva#2237](https://github.com/cs3org/reva/pull/2237): Skip get user call in eosfs in case previous ones also failed
  * Enh [cs3org/reva#2266](https://github.com/cs3org/reva/pull/2266): Callback for the EOS UID cache to retry fetch for failed keys
  * Enh [cs3org/reva#2215](https://github.com/cs3org/reva/pull/2215): Aggregate resource info properties for virtual views
  * Enh [cs3org/reva#2271](https://github.com/cs3org/reva/pull/2271): Revamp the favorite manager and add the cbox sql driver
  * Enh [cs3org/reva#2248](https://github.com/cs3org/reva/pull/2248): Cache whether a user home was created or not
  * Enh [cs3org/reva#2282](https://github.com/cs3org/reva/pull/2282): Return a proper NOT_FOUND error when a user or group is not found
  * Enh [cs3org/reva#2268](https://github.com/cs3org/reva/pull/2268): Add the reverseproxy http service
  * Enh [cs3org/reva#2207](https://github.com/cs3org/reva/pull/2207): Enable users to list all spaces
  * Enh [cs3org/reva#2286](https://github.com/cs3org/reva/pull/2286): Add trace ID to middleware loggers
  * Enh [cs3org/reva#2251](https://github.com/cs3org/reva/pull/2251): Mentix service inference
  * Enh [cs3org/reva#2218](https://github.com/cs3org/reva/pull/2218): Allow filtering of mime types supported by app providers
  * Enh [cs3org/reva#2213](https://github.com/cs3org/reva/pull/2213): Add public link share type to propfind response
  * Enh [cs3org/reva#2253](https://github.com/cs3org/reva/pull/2253): Support the file editor role for public links
  * Enh [cs3org/reva#2208](https://github.com/cs3org/reva/pull/2208): Reduce redundant stat calls when statting by resource ID
  * Enh [cs3org/reva#2235](https://github.com/cs3org/reva/pull/2235): Specify a list of allowed folders/files to be archived
  * Enh [cs3org/reva#2267](https://github.com/cs3org/reva/pull/2267): Restrict the paths where share creation is allowed
  * Enh [cs3org/reva#2252](https://github.com/cs3org/reva/pull/2252): Add the xattr sys.acl to SysACL (eosgrpc)
  * Enh [cs3org/reva#2239](https://github.com/cs3org/reva/pull/2239): Update toml configs

   https://github.com/owncloud/ocis/pull/2737
   https://github.com/owncloud/ocis/pull/2726
   https://github.com/owncloud/ocis/pull/2790
   https://github.com/owncloud/ocis/pull/2797

# Changelog for [1.14.0] (2021-10-27)

The following sections list the changes for 1.14.0.

[1.14.0]: https://github.com/owncloud/ocis/compare/v1.13.0...v1.14.0

## Summary

* Security - Don't expose services by default: [#2612](https://github.com/owncloud/ocis/issues/2612)
* Bugfix - Create parent directories for idp configuration: [#2667](https://github.com/owncloud/ocis/issues/2667)
* Change - New default data paths and easier configuration of the data path: [#2590](https://github.com/owncloud/ocis/pull/2590)
* Change - Configurable default quota: [#2621](https://github.com/owncloud/ocis/issues/2621)
* Change - Split spaces webdav url and graph url in base and path: [#2660](https://github.com/owncloud/ocis/pull/2660)
* Change - Update ownCloud Web to v4.4.0: [#2681](https://github.com/owncloud/ocis/pull/2681)
* Enhancement - Replace fileb0x with go-embed: [#1199](https://github.com/owncloud/ocis/issues/1199)
* Enhancement - Start up a new machine auth provider in the storage service: [#2528](https://github.com/owncloud/ocis/pull/2528)
* Enhancement - Add a middleware to authenticate public share requests: [#2536](https://github.com/owncloud/ocis/pull/2536)
* Enhancement - Lower TUS max chunk size: [#2584](https://github.com/owncloud/ocis/pull/2584)
* Enhancement - Upgrade to go-micro v4.1.0: [#2616](https://github.com/owncloud/ocis/pull/2616)
* Enhancement - Report quota states: [#2628](https://github.com/owncloud/ocis/pull/2628)
* Enhancement - Broaden bufbuild/Buf usage: [#2630](https://github.com/owncloud/ocis/pull/2630)
* Enhancement - Add sharees additional info parameter config to ocs: [#2637](https://github.com/owncloud/ocis/pull/2637)
* Enhancement - Enforce permission on update space quota: [#2650](https://github.com/owncloud/ocis/pull/2650)
* Enhancement - Update lico to v0.51.1: [#2654](https://github.com/owncloud/ocis/pull/2654)
* Enhancement - Add user setting capability: [#2655](https://github.com/owncloud/ocis/pull/2655)
* Enhancement - Update reva to v1.15: [#2658](https://github.com/owncloud/ocis/pull/2658)
* Enhancement - Review and correct http header: [#2666](https://github.com/owncloud/ocis/pull/2666)

## Details

* Security - Don't expose services by default: [#2612](https://github.com/owncloud/ocis/issues/2612)

   We've changed the bind behaviour for all non public facing services. Before this
   PR all services would listen on all interfaces. After this PR, all services
   listen on 127.0.0.1 only, except the proxy which is listening on 0.0.0.0:9200.

   https://github.com/owncloud/ocis/issues/2612

* Bugfix - Create parent directories for idp configuration: [#2667](https://github.com/owncloud/ocis/issues/2667)

   The parent directories of the identifier-registration.yaml config file might not
   exist when starting idp. Create them, when that is the case.

   https://github.com/owncloud/ocis/issues/2667

* Change - New default data paths and easier configuration of the data path: [#2590](https://github.com/owncloud/ocis/pull/2590)

   We've changed the default data path for our release artifacts: - oCIS docker
   images will now store all data in `/var/lib/ocis` instead in `/var/tmp/ocis` -
   binary releases will now store all data in `~/.ocis` instead of `/var/tmp/ocis`

   Also if you're a developer and you run oCIS from source, it will store all data
   in `~/.ocis` from now on.

   You can now easily change the data path for all extensions by setting the
   environment variable `OCIS_BASE_DATA_PATH`.

   If you want to package oCIS, you also can set the default data path at compile
   time, eg. by passing `-X
   "github.com/owncloud/ocis/ocis-pkg/config/defaults.BaseDataPathType=path" -X
   "github.com/owncloud/ocis/ocis-pkg/config/defaults.BaseDataPathValue=/var/lib/ocis"`
   to your go build step.

   https://github.com/owncloud/ocis/pull/2590

* Change - Configurable default quota: [#2621](https://github.com/owncloud/ocis/issues/2621)

   When creating a new space a (configurable) default quota will be used (instead
   the hardcoded one). One can set the EnvVar `GRAPH_SPACES_DEFAULT_QUOTA` to
   configure it

   https://github.com/owncloud/ocis/issues/2621
   https://jira.owncloud.com/browse/OCIS-2070

* Change - Split spaces webdav url and graph url in base and path: [#2660](https://github.com/owncloud/ocis/pull/2660)

   We've fixed the behavior for the spaces webdav url and graph explorer graph url
   settings, so that they respect the environment variable `OCIS_URL`. Previously
   oCIS admins needed to set these URLs manually to make spaces and the graph
   explorer work.

   https://github.com/owncloud/ocis/issues/2659
   https://github.com/owncloud/ocis/pull/2660

* Change - Update ownCloud Web to v4.4.0: [#2681](https://github.com/owncloud/ocis/pull/2681)

   Tags: web

   We updated ownCloud Web to v4.4.0. Please refer to the changelog (linked) for
   details on the web release.

   https://github.com/owncloud/ocis/pull/2681
   https://github.com/owncloud/web/releases/tag/v4.4.0

* Enhancement - Replace fileb0x with go-embed: [#1199](https://github.com/owncloud/ocis/issues/1199)

   Go-embed already brings the functionality we need but with less code. We decided
   to use it instead of 3rd party fileb0x

   https://github.com/owncloud/ocis/issues/1199
   https://github.com/owncloud/ocis/pull/2631
   https://github.com/owncloud/ocis/pull/2649

* Enhancement - Start up a new machine auth provider in the storage service: [#2528](https://github.com/owncloud/ocis/pull/2528)

   This PR also adds the config to skip encoding user groups in reva tokens

   https://github.com/owncloud/ocis/pull/2528
   https://github.com/owncloud/ocis/pull/2529

* Enhancement - Add a middleware to authenticate public share requests: [#2536](https://github.com/owncloud/ocis/pull/2536)

   Added a new middleware to authenticate public share requests. This makes it
   possible to use APIs which require an authenticated context with public shares.

   https://github.com/owncloud/ocis/issues/2479
   https://github.com/owncloud/ocis/pull/2536
   https://github.com/owncloud/ocis/pull/2652

* Enhancement - Lower TUS max chunk size: [#2584](https://github.com/owncloud/ocis/pull/2584)

   We've lowered the TUS max chunk size from infinite to 0.1GB so that chunking
   actually happens.

   https://github.com/owncloud/ocis/pull/2584
   https://github.com/cs3org/reva/pull/2136

* Enhancement - Upgrade to go-micro v4.1.0: [#2616](https://github.com/owncloud/ocis/pull/2616)

   We've upgraded to go-micro v4.1.0

   https://github.com/owncloud/ocis/pull/2616

* Enhancement - Report quota states: [#2628](https://github.com/owncloud/ocis/pull/2628)

   When listing the available spaces via the GraphAPI we now return quota states to
   make it easier for the clients to add visual indicators.

   https://github.com/owncloud/ocis/pull/2628

* Enhancement - Broaden bufbuild/Buf usage: [#2630](https://github.com/owncloud/ocis/pull/2630)

   We've switched the usage of bufbuild/Buf from a protoc replacement only to also
   using it to configure the outputs and pinning dependencies.

   https://github.com/owncloud/ocis/pull/2630
   https://github.com/owncloud/ocis/pull/2616

* Enhancement - Add sharees additional info parameter config to ocs: [#2637](https://github.com/owncloud/ocis/pull/2637)

   https://github.com/owncloud/ocis/pull/2637

* Enhancement - Enforce permission on update space quota: [#2650](https://github.com/owncloud/ocis/pull/2650)

   Added a check that only users with the `set-space-quota` permission can update
   the space quota.

   https://github.com/owncloud/ocis/pull/2650

* Enhancement - Update lico to v0.51.1: [#2654](https://github.com/owncloud/ocis/pull/2654)

   Updated lico to v0.51.1 This update includes: * Apply LibreGraph naming treewide
   * move to go1.17 * Update 3rd party Go dependencies

   https://github.com/owncloud/ocis/pull/2654

* Enhancement - Add user setting capability: [#2655](https://github.com/owncloud/ocis/pull/2655)

   We've added a capability to communicate the existence of a user settings service
   to clients.

   https://github.com/owncloud/web/issues/5926
   https://github.com/owncloud/ocis/pull/2655

* Enhancement - Update reva to v1.15: [#2658](https://github.com/owncloud/ocis/pull/2658)

   Updated reva to v1.15 This update includes:

  * Fix [cs3org/reva#2168](https://github.com/cs3org/reva/pull/2168): Override provider if was previously registered
  * Fix [cs3org/reva#2173](https://github.com/cs3org/reva/pull/2173): Fix archiver max size reached error
  * Fix [cs3org/reva#2167](https://github.com/cs3org/reva/pull/2167): Handle nil quota in decomposedfs
  * Fix [cs3org/reva#2153](https://github.com/cs3org/reva/pull/2153): Restrict EOS project spaces sharing permissions to admins and writers
  * Fix [cs3org/reva#2179](https://github.com/cs3org/reva/pull/2179): Fix the returned permissions for webdav uploads
  * Chg [cs3org/reva#2479](https://github.com/cs3org/reva/pull/2479): Make apps able to work with public shares
  * Enh [cs3org/reva#2174](https://github.com/cs3org/reva/pull/2174): Inherit ACLs for files from parent directories
  * Enh [cs3org/reva#2152](https://github.com/cs3org/reva/pull/2152): Add a reference parameter to the getQuota request
  * Enh [cs3org/reva#2171](https://github.com/cs3org/reva/pull/2171): Add optional claim parameter to machine auth
  * Enh [cs3org/reva#2135](https://github.com/cs3org/reva/pull/2135): Nextcloud test improvements
  * Enh [cs3org/reva#2180](https://github.com/cs3org/reva/pull/2180): Remove OCDAV options namespace parameter
  * Enh [cs3org/reva#2170](https://github.com/cs3org/reva/pull/2170): Handle propfind requests for existing files
  * Enh [cs3org/reva#2165](https://github.com/cs3org/reva/pull/2165): Allow access to recycle bin for arbitrary paths outside homes
  * Enh [cs3org/reva#2189](https://github.com/cs3org/reva/pull/2189): Add user settings capability
  * Enh [cs3org/reva#2162](https://github.com/cs3org/reva/pull/2162): Implement the UpdateStorageSpace method
  * Enh [cs3org/reva#2117](https://github.com/cs3org/reva/pull/2117): Add ocs cache warmup strategy for first request from the user

   https://github.com/owncloud/ocis/pull/2658
   https://github.com/owncloud/ocis/pull/2536
   https://github.com/owncloud/ocis/pull/2650
   https://github.com/owncloud/ocis/pull/2680

* Enhancement - Review and correct http header: [#2666](https://github.com/owncloud/ocis/pull/2666)

   Reviewed and corrected the necessary http headers. Made CORS configurable.

   https://github.com/owncloud/ocis/pull/2666

# Changelog for [1.13.0] (2021-10-13)

The following sections list the changes for 1.13.0.

[1.13.0]: https://github.com/owncloud/ocis/compare/v1.12.0...v1.13.0

## Summary

* Bugfix - Use proper url path decode on the username: [#2511](https://github.com/owncloud/ocis/pull/2511)
* Bugfix - Remove notifications placeholder: [#2514](https://github.com/owncloud/ocis/pull/2514)
* Bugfix - Fix the account resolver middleware: [#2557](https://github.com/owncloud/ocis/pull/2557)
* Bugfix - Race condition in config parsing: [#2574](https://github.com/owncloud/ocis/pull/2574)
* Bugfix - Fix version information for extensions: [#2575](https://github.com/owncloud/ocis/pull/2575)
* Bugfix - Remove asset path configuration option from proxy: [#2576](https://github.com/owncloud/ocis/pull/2576)
* Bugfix - Add the gatewaysvc to all shared configuration in REVA services: [#2597](https://github.com/owncloud/ocis/pull/2597)
* Change - Make the drives create method odata compliant: [#2531](https://github.com/owncloud/ocis/pull/2531)
* Change - Unify Envvar names configuring REVA gateway address: [#2587](https://github.com/owncloud/ocis/pull/2587)
* Change - Update ownCloud Web to v4.3.0: [#2589](https://github.com/owncloud/ocis/pull/2589)
* Change - Configure users and metadata storage separately: [#2598](https://github.com/owncloud/ocis/pull/2598)
* Enhancement - TLS config options for ldap in reva: [#2492](https://github.com/owncloud/ocis/pull/2492)
* Enhancement - Redirect invalid links to oC Web: [#2493](https://github.com/owncloud/ocis/pull/2493)
* Enhancement - Add option to skip generation of demo users and groups: [#2495](https://github.com/owncloud/ocis/pull/2495)
* Enhancement - Allow overriding the cookie based route by claim: [#2508](https://github.com/owncloud/ocis/pull/2508)
* Enhancement - Expose the reva archiver in OCIS: [#2509](https://github.com/owncloud/ocis/pull/2509)
* Enhancement - Set reva JWT token expiration time to 24 hours by default: [#2527](https://github.com/owncloud/ocis/pull/2527)
* Enhancement - Use reva's Authenticate method instead of spawning token managers: [#2528](https://github.com/owncloud/ocis/pull/2528)
* Enhancement - Add maximum files and size to archiver capabilities: [#2544](https://github.com/owncloud/ocis/pull/2544)
* Enhancement - Make mimetype allow list configurable for app provider: [#2553](https://github.com/owncloud/ocis/pull/2553)
* Enhancement - Reduced repository size: [#2579](https://github.com/owncloud/ocis/pull/2579)
* Enhancement - Add allow_creation parameter to mime type config: [#2591](https://github.com/owncloud/ocis/pull/2591)
* Enhancement - Favorites capability: [#2599](https://github.com/owncloud/ocis/pull/2599)
* Enhancement - Updated MimeTypes configuration for AppRegistry: [#2603](https://github.com/owncloud/ocis/pull/2603)
* Enhancement - Upgrade to GO 1.17: [#2605](https://github.com/owncloud/ocis/pull/2605)
* Enhancement - Return the newly created space: [#2610](https://github.com/owncloud/ocis/pull/2610)
* Enhancement - Update reva to v1.14.0: [#2615](https://github.com/owncloud/ocis/pull/2615)

## Details

* Bugfix - Use proper url path decode on the username: [#2511](https://github.com/owncloud/ocis/pull/2511)

   We now properly decode the username when reading it from a url parameter

   https://github.com/owncloud/ocis/pull/2511

* Bugfix - Remove notifications placeholder: [#2514](https://github.com/owncloud/ocis/pull/2514)

   Since Reva was communicating its notification capabilities incorrectly, oCIS
   relied on a hardcoded string to overwrite them. This has been fixed in
   [reva#1819](https://github.com/cs3org/reva/pull/1819) so we now removed the
   hardcoded string and don't modify Reva's notification capabilities anymore in
   order to fix clients having to poll a (non-existent) notifications endpoint.

   https://github.com/owncloud/ocis/pull/2514

* Bugfix - Fix the account resolver middleware: [#2557](https://github.com/owncloud/ocis/pull/2557)

   The accounts resolver middleware put an empty token into the request when the
   user was already present. Added a step to get the token for the user.

   https://github.com/owncloud/ocis/pull/2557

* Bugfix - Race condition in config parsing: [#2574](https://github.com/owncloud/ocis/pull/2574)

   There was a race condition in the config parsing when configuring the storage
   services caused by services overwriting a pointer to a config value. We fixed it
   by setting sane defaults.

   https://github.com/owncloud/ocis/pull/2574

* Bugfix - Fix version information for extensions: [#2575](https://github.com/owncloud/ocis/pull/2575)

   We've fixed the behavior for `ocis version` which previously always showed
   `0.0.0` as version for extensions. Now the real version of the extensions are
   shown.

   https://github.com/owncloud/ocis/pull/2575

* Bugfix - Remove asset path configuration option from proxy: [#2576](https://github.com/owncloud/ocis/pull/2576)

   We've remove the asset path configuration option (`--asset-path` or
   `PROXY_ASSET_PATH`) since it didn't do anything at all.

   https://github.com/owncloud/ocis/pull/2576

* Bugfix - Add the gatewaysvc to all shared configuration in REVA services: [#2597](https://github.com/owncloud/ocis/pull/2597)

   We've fixed the configuration for REVA services which didn't have a gatewaysvc
   in their shared configuration. This could lead to default gatewaysvc addresses
   in the auth middleware. Now it is set everywhere.

   https://github.com/owncloud/ocis/pull/2597

* Change - Make the drives create method odata compliant: [#2531](https://github.com/owncloud/ocis/pull/2531)

   When creating a space on the graph API we now use the POST Body to provide the
   parameters.

   https://github.com/owncloud/ocis/pull/2531
   https://github.com/owncloud/ocis/pull/2535
   https://www.odata.org/getting-started/basic-tutorial/#modifyData

* Change - Unify Envvar names configuring REVA gateway address: [#2587](https://github.com/owncloud/ocis/pull/2587)

   We've renamed all envvars configuring REVA gateway address to `REVA_GATEWAY`,
   additionally we renamed the cli parameters to `--reva-gateway-addr` and adjusted
   the description

   https://github.com/owncloud/ocis/issues/2091
   https://github.com/owncloud/ocis/pull/2587

* Change - Update ownCloud Web to v4.3.0: [#2589](https://github.com/owncloud/ocis/pull/2589)

   Tags: web

   We updated ownCloud Web to v4.3.0. Please refer to the changelog (linked) for
   details on the web release.

   https://github.com/owncloud/ocis/pull/2589
   https://github.com/owncloud/web/releases/tag/v4.3.0

* Change - Configure users and metadata storage separately: [#2598](https://github.com/owncloud/ocis/pull/2598)

   We've fixed the configuration behaviour of the user and metadata service writing
   in the same directory when using oCIS storage.

   Therefore we needed to separate the configuration of the users and metadata
   storage so that they now can be configured totally separate.

   https://github.com/owncloud/ocis/pull/2598

* Enhancement - TLS config options for ldap in reva: [#2492](https://github.com/owncloud/ocis/pull/2492)

   We added the new config options "ldap-cacert" and "ldap-insecure" to the auth-,
   users- and groups-provider services to be able to do proper TLS configuration
   for the LDAP clients. "ldap-cacert" is by default configured to add the bundled
   glauth LDAP servers certificate to the trusted set for the LDAP clients.
   "ldap-insecure" is set to "false" by default and can be used to disable
   certificate checks (only advisable for development and test environments).

   https://github.com/owncloud/ocis/pull/2492

* Enhancement - Redirect invalid links to oC Web: [#2493](https://github.com/owncloud/ocis/pull/2493)

   Invalid links (eg. https://foo.bar/index.php/apps/pdfviewer) will be redirect to
   ownCloud Web instead of displaying a blank page with a "not found" message.

   https://github.com/owncloud/ocis/pull/2493
   https://github.com/owncloud/ocis/pull/2512

* Enhancement - Add option to skip generation of demo users and groups: [#2495](https://github.com/owncloud/ocis/pull/2495)

   We've added a new environment variable to decide whether we should generate the
   demo users and groups or not. This environment variable is set to `true` by
   default, so the demo users and groups will get generated by default as long as
   oCIS is in its "technical preview" stage.

   In any case, there are still some users and groups automatically generated: for
   users: Reva IOP, Kopano IDP, admin; for groups: sysusers and users.

   https://github.com/owncloud/ocis/pull/2495

* Enhancement - Allow overriding the cookie based route by claim: [#2508](https://github.com/owncloud/ocis/pull/2508)

   When determining the routing policy we now let the claim override the cookie so
   that users are routed to the correct backend after login.

   https://github.com/owncloud/ocis/pull/2508

* Enhancement - Expose the reva archiver in OCIS: [#2509](https://github.com/owncloud/ocis/pull/2509)

   The reva archiver can now be accessed through the storage frontend service

   https://github.com/owncloud/ocis/pull/2509

* Enhancement - Set reva JWT token expiration time to 24 hours by default: [#2527](https://github.com/owncloud/ocis/pull/2527)

   https://github.com/owncloud/ocis/pull/2527

* Enhancement - Use reva's Authenticate method instead of spawning token managers: [#2528](https://github.com/owncloud/ocis/pull/2528)

   When using the CS3 proxy backend, we previously obtained the user from reva's
   userprovider service and minted the token ourselves. This required maintaining a
   shared JWT secret between ocis and reva, as well duplication of logic. This PR
   delegates this logic by using the `Authenticate` method provided by the reva
   gateway service to obtain this token, making it an arbitrary, indestructible
   entry. Currently, the changes have been made to the proxy service but will be
   extended to others as well.

   https://github.com/owncloud/ocis/pull/2528

* Enhancement - Add maximum files and size to archiver capabilities: [#2544](https://github.com/owncloud/ocis/pull/2544)

   We added the maximum files count and maximum archive size of the archiver to the
   capabilities endpoint. Clients can use this to generate warnings before the
   actual archive creation fails.

   https://github.com/owncloud/ocis/issues/2537
   https://github.com/owncloud/ocis/pull/2544
   https://github.com/cs3org/reva/pull/2105

* Enhancement - Make mimetype allow list configurable for app provider: [#2553](https://github.com/owncloud/ocis/pull/2553)

   We've added a configuration option to configure the mimetype allow list
   introduced in cs3org/reva#2095. This also makes it possible to set one
   application per mime type as a default.

   https://github.com/owncloud/ocis/issues/2563
   https://github.com/owncloud/ocis/pull/2553
   https://github.com/cs3org/reva/pull/2095

* Enhancement - Reduced repository size: [#2579](https://github.com/owncloud/ocis/pull/2579)

   We removed leftover artifacts from the migration to a single repository.

   https://github.com/owncloud/ocis/pull/2579

* Enhancement - Add allow_creation parameter to mime type config: [#2591](https://github.com/owncloud/ocis/pull/2591)

   https://github.com/owncloud/ocis/pull/2591

* Enhancement - Favorites capability: [#2599](https://github.com/owncloud/ocis/pull/2599)

   We've added a capability for the storage frontend which can be used to announce
   to clients whether or not favorites are supported. By default this is disabled
   because the listing of favorites doesn't survive service restarts at the moment.

   https://github.com/owncloud/ocis/pull/2599

* Enhancement - Updated MimeTypes configuration for AppRegistry: [#2603](https://github.com/owncloud/ocis/pull/2603)

   We updated the type of the mime types config to a list, to keep the order of
   mime types from the config.

   https://github.com/owncloud/ocis/pull/2603

* Enhancement - Upgrade to GO 1.17: [#2605](https://github.com/owncloud/ocis/pull/2605)

   We've upgraded the used GO version from 1.16 to 1.17.

   https://github.com/owncloud/ocis/pull/2605

* Enhancement - Return the newly created space: [#2610](https://github.com/owncloud/ocis/pull/2610)

   Changed the response of the CreateSpace method to include the newly created
   space.

   https://github.com/owncloud/ocis/pull/2610
   https://github.com/cs3org/reva/pull/2158

* Enhancement - Update reva to v1.14.0: [#2615](https://github.com/owncloud/ocis/pull/2615)

   This update includes:

  * Bugfix [cs3org/reva#2103](https://github.com/cs3org/reva/pull/2103): AppProvider: propagate back errors reported by WOPI
  * Bugfix [cs3org/reva#2149](https://github.com/cs3org/reva/pull/2149): Remove excess info from the http list app providers endpoint
  * Bugfix [cs3org/reva#2114](https://github.com/cs3org/reva/pull/2114): Add as default app while registering and skip unset mimetypes
  * Bugfix [cs3org/reva#2095](https://github.com/cs3org/reva/pull/2095): Fix app open when multiple app providers are present
  * Bugfix [cs3org/reva#2135](https://github.com/cs3org/reva/pull/2135): Make TUS capabilities configurable
  * Bugfix [cs3org/reva#2076](https://github.com/cs3org/reva/pull/2076): Fix chi routing
  * Bugfix [cs3org/reva#2077](https://github.com/cs3org/reva/pull/2077): Fix concurrent registration of mimetypes
  * Bugfix [cs3org/reva#2154](https://github.com/cs3org/reva/pull/2154): Return OK when trying to delete a non existing reference
  * Bugfix [cs3org/reva#2078](https://github.com/cs3org/reva/pull/2078): Fix nil pointer exception in stat
  * Bugfix [cs3org/reva#2073](https://github.com/cs3org/reva/pull/2073): Fix opening a readonly filetype with WOPI
  * Bugfix [cs3org/reva#2140](https://github.com/cs3org/reva/pull/2140): Map GRPC error codes to REVA errors
  * Bugfix [cs3org/reva#2147](https://github.com/cs3org/reva/pull/2147): Follow up of #2138: this is the new expected format
  * Bugfix [cs3org/reva#2116](https://github.com/cs3org/reva/pull/2116): Differentiate share types when retrieving received shares in sql driver
  * Bugfix [cs3org/reva#2074](https://github.com/cs3org/reva/pull/2074): Fix Stat() for EOS storage provider
  * Bugfix [cs3org/reva#2151](https://github.com/cs3org/reva/pull/2151): Fix return code for webdav uploads when the token expired
  * Change [cs3org/reva#2121](https://github.com/cs3org/reva/pull/2121): Sharemanager API change
  * Enhancement [cs3org/reva#2090](https://github.com/cs3org/reva/pull/2090): Return space name during list storage spaces
  * Enhancement [cs3org/reva#2138](https://github.com/cs3org/reva/pull/2138): Default AppProvider on top of the providers list
  * Enhancement [cs3org/reva#2137](https://github.com/cs3org/reva/pull/2137): Revamp app registry and add parameter to control file creation
  * Enhancement [cs3org/reva#145](https://github.com/cs3org/reva/pull/2137): UI improvements for the AppProviders
  * Enhancement [cs3org/reva#2088](https://github.com/cs3org/reva/pull/2088): Add archiver and app provider to ocs capabilities
  * Enhancement [cs3org/reva#2537](https://github.com/cs3org/reva/pull/2537): Add maximum files and size to archiver capabilities
  * Enhancement [cs3org/reva#2100](https://github.com/cs3org/reva/pull/2100): Add support for resource id to the archiver
  * Enhancement [cs3org/reva#2158](https://github.com/cs3org/reva/pull/2158): Augment the Id of new spaces
  * Enhancement [cs3org/reva#2085](https://github.com/cs3org/reva/pull/2085): Make encoding user groups in access tokens configurable
  * Enhancement [cs3org/reva#146](https://github.com/cs3org/reva/pull/146): Filter the denial shares (permission = 0) out of
  * Enhancement [cs3org/reva#2141](https://github.com/cs3org/reva/pull/2141): Use golang v1.17
  * Enhancement [cs3org/reva#2053](https://github.com/cs3org/reva/pull/2053): Safer defaults for TLS verification on LDAP connections
  * Enhancement [cs3org/reva#2115](https://github.com/cs3org/reva/pull/2115): Reduce code duplication in LDAP related drivers
  * Enhancement [cs3org/reva#1989](https://github.com/cs3org/reva/pull/1989): Add redirects from OC10 URL formats
  * Enhancement [cs3org/reva#2479](https://github.com/cs3org/reva/pull/2479): Limit publicshare and resourceinfo scope content
  * Enhancement [cs3org/reva#2071](https://github.com/cs3org/reva/pull/2071): Implement listing favorites via the dav report API
  * Enhancement [cs3org/reva#2091](https://github.com/cs3org/reva/pull/2091): Nextcloud share managers
  * Enhancement [cs3org/reva#2070](https://github.com/cs3org/reva/pull/2070): More unit tests for the Nextcloud storage provider
  * Enhancement [cs3org/reva#2087](https://github.com/cs3org/reva/pull/2087): More unit tests for the Nextcloud auth and user managers
  * Enhancement [cs3org/reva#2075](https://github.com/cs3org/reva/pull/2075): Make owncloudsql leverage existing filecache index
  * Enhancement [cs3org/reva#2050](https://github.com/cs3org/reva/pull/2050): Add a share types filter to the OCS API
  * Enhancement [cs3org/reva#2134](https://github.com/cs3org/reva/pull/2134): Use space Type from request
  * Enhancement [cs3org/reva#2132](https://github.com/cs3org/reva/pull/2132): Align local tests with drone setup
  * Enhancement [cs3org/reva#2095](https://github.com/cs3org/reva/pull/2095): Whitelisting for apps
  * Enhancement [cs3org/reva#2155](https://github.com/cs3org/reva/pull/2155): Pass an extra query parameter to WOPI /openinapp with a

   https://github.com/owncloud/ocis/pull/2615
   https://github.com/owncloud/ocis/pull/2566
   https://github.com/owncloud/ocis/pull/2520

# Changelog for [1.12.0] (2021-09-14)

The following sections list the changes for 1.12.0.

[1.12.0]: https://github.com/owncloud/ocis/compare/v1.11.0...v1.12.0

## Summary

* Bugfix - Set English as default language in the dropdown in the settings page: [#2465](https://github.com/owncloud/ocis/pull/2465)
* Bugfix - Remove non working proxy route and fix cs3 users example: [#2474](https://github.com/owncloud/ocis/pull/2474)
* Change - Remove OnlyOffice extension: [#2433](https://github.com/owncloud/ocis/pull/2433)
* Change - Remove OnlyOffice extension: [#2433](https://github.com/owncloud/ocis/pull/2433)
* Change - Update ownCloud Web to v4.2.0: [#2501](https://github.com/owncloud/ocis/pull/2501)
* Enhancement - Add app provider and app provider registry: [#2204](https://github.com/owncloud/ocis/pull/2204)
* Enhancement - Update go-chi/chi to version 5.0.3: [#2429](https://github.com/owncloud/ocis/pull/2429)
* Enhancement - Upgrade go micro to v3.6.0: [#2451](https://github.com/owncloud/ocis/pull/2451)
* Enhancement - Add set space quota permission: [#2459](https://github.com/owncloud/ocis/pull/2459)
* Enhancement - Add the create space permission: [#2461](https://github.com/owncloud/ocis/pull/2461)
* Enhancement - Create a Space using the Graph API: [#2471](https://github.com/owncloud/ocis/pull/2471)
* Enhancement - Update reva to v1.13.0: [#2477](https://github.com/owncloud/ocis/pull/2477)

## Details

* Bugfix - Set English as default language in the dropdown in the settings page: [#2465](https://github.com/owncloud/ocis/pull/2465)

   The language dropdown didn't have a default language selected, and it was
   showing an empty value. Now it shows English instead.

   https://github.com/owncloud/ocis/pull/2465

* Bugfix - Remove non working proxy route and fix cs3 users example: [#2474](https://github.com/owncloud/ocis/pull/2474)

   We removed a non working route from the proxy default configuration and fixed
   the cs3 users deployment example since it still used the accounts service. It
   now only uses the configured LDAP.

   https://github.com/owncloud/ocis/pull/2474

* Change - Remove OnlyOffice extension: [#2433](https://github.com/owncloud/ocis/pull/2433)

   Tags: OnlyOffice

   We've removed the OnlyOffice extension in oCIS. OnlyOffice has their own web
   extension for OC10 backend now with [a dedicated
   guide](https://owncloud.dev/clients/web/deployments/oc10-app/#onlyoffice). In
   oCIS, we will follow up with a guide on how to start a WOPI server providing
   OnlyOffice soon.

   https://github.com/owncloud/ocis/pull/2433

* Change - Remove OnlyOffice extension: [#2433](https://github.com/owncloud/ocis/pull/2433)

   Tags: OnlyOffice

   We've removed the OnlyOffice extension in oCIS. OnlyOffice has their own web
   extension for OC10 backend now with [a dedicated
   guide](https://owncloud.dev/clients/web/deployments/oc10-app/#onlyoffice). In
   oCIS, we will follow up with a guide on how to start a WOPI server providing
   OnlyOffice soon.

   https://github.com/owncloud/ocis/pull/2433

* Change - Update ownCloud Web to v4.2.0: [#2501](https://github.com/owncloud/ocis/pull/2501)

   Tags: web

   We updated ownCloud Web to v4.2.0. Please refer to the changelog (linked) for
   details on the web release.

   https://github.com/owncloud/ocis/pull/2501
   https://github.com/owncloud/web/releases/tag/v4.2.0

* Enhancement - Add app provider and app provider registry: [#2204](https://github.com/owncloud/ocis/pull/2204)

   We added the app provider and app provider registry. Now the CS3org WOPI server
   can be registered and OpenInApp requests can be done.

   https://github.com/owncloud/ocis/pull/2204
   https://github.com/cs3org/reva/pull/1785

* Enhancement - Update go-chi/chi to version 5.0.3: [#2429](https://github.com/owncloud/ocis/pull/2429)

   Updated go-chi/chi to the latest release

   https://github.com/owncloud/ocis/pull/2429

* Enhancement - Upgrade go micro to v3.6.0: [#2451](https://github.com/owncloud/ocis/pull/2451)

   Go micro and all go micro plugins are now on v3.6.0

   https://github.com/owncloud/ocis/pull/2451

* Enhancement - Add set space quota permission: [#2459](https://github.com/owncloud/ocis/pull/2459)

   In preparation for the upcoming spaces features a `SetSpaceQuota` permission was
   added.

   https://github.com/owncloud/ocis/pull/2459

* Enhancement - Add the create space permission: [#2461](https://github.com/owncloud/ocis/pull/2461)

   In preparation for the upcoming spaces features a `Create Space` permission was
   added.

   https://github.com/owncloud/ocis/pull/2461

* Enhancement - Create a Space using the Graph API: [#2471](https://github.com/owncloud/ocis/pull/2471)

   Spaces can now be created on `POST /drives/{drive-name}`. Only users with the
   `create-space` permissions can perform this operation.

   Allowed body form values are:

   - `quota` (bytes) maximum amount of bytes stored in the space. - `maxQuotaFiles`
   (integer) maximum amount of files supported by the space.

   https://github.com/owncloud/ocis/pull/2471

* Enhancement - Update reva to v1.13.0: [#2477](https://github.com/owncloud/ocis/pull/2477)

   This update includes:

  * Bugfix [cs3org/reva#2054](https://github.com/cs3org/reva/pull/2054): Fix the response after deleting a share
  * Bugfix [cs3org/reva#2026](https://github.com/cs3org/reva/pull/2026): Fix moving of a shared file
  * Bugfix [cs3org/reva#1605](https://github.com/cs3org/reva/pull/1605): Allow to expose full paths in OCS API
  * Bugfix [cs3org/reva#2033](https://github.com/cs3org/reva/pull/2033): Fix the storage id of shares
  * Bugfix [cs3org/reva#1991](https://github.com/cs3org/reva/pull/1991): Remove share references when declining shares
  * Enhancement [cs3org/reva#1994](https://github.com/cs3org/reva/pull/1994): Add owncloudsql driver for the userprovider
  * Enhancement [cs3org/reva#2065](https://github.com/cs3org/reva/pull/2065): New sharing role Manager
  * Enhancement [cs3org/reva#2015](https://github.com/cs3org/reva/pull/2015): Add spaces to the list of capabilities
  * Enhancement [cs3org/reva#2041](https://github.com/cs3org/reva/pull/2041): Create operations for Spaces
  * Enhancement [cs3org/reva#2029](https://github.com/cs3org/reva/pull/2029): Tracing agent configuration

   https://github.com/owncloud/ocis/pull/2477

# Changelog for [1.11.0] (2021-08-24)

The following sections list the changes for 1.11.0.

[1.11.0]: https://github.com/owncloud/ocis/compare/v1.10.0...v1.11.0

## Summary

* Bugfix - Specify primary user type for all accounts: [#2364](https://github.com/owncloud/ocis/pull/2364)
* Bugfix - Fix naming of the user- and groupprovider services: [#2388](https://github.com/owncloud/ocis/pull/2388)
* Change - Update ownCloud Web to v4.1.0: [#2426](https://github.com/owncloud/ocis/pull/2426)
* Enhancement - Use non root user for the owncloud/ocis docker image: [#2380](https://github.com/owncloud/ocis/pull/2380)
* Enhancement - Replace unmaintained jwt library: [#2386](https://github.com/owncloud/ocis/pull/2386)
* Enhancement - Update bleve to version 2.1.0: [#2391](https://github.com/owncloud/ocis/pull/2391)
* Enhancement - Update github.com/coreos/go-oidc to v3.0.0: [#2393](https://github.com/owncloud/ocis/pull/2393)
* Enhancement - Update reva to v1.12: [#2423](https://github.com/owncloud/ocis/pull/2423)

## Details

* Bugfix - Specify primary user type for all accounts: [#2364](https://github.com/owncloud/ocis/pull/2364)

   https://github.com/owncloud/ocis/pull/2364

* Bugfix - Fix naming of the user- and groupprovider services: [#2388](https://github.com/owncloud/ocis/pull/2388)

   The services are called "storage-userprovider" and "storage-groupprovider". The
   'ocis help' output was misleading.

   https://github.com/owncloud/ocis/pull/2388

* Change - Update ownCloud Web to v4.1.0: [#2426](https://github.com/owncloud/ocis/pull/2426)

   Tags: web

   We updated ownCloud Web to v4.1.0. Please refer to the changelog (linked) for
   details on the web release.

   https://github.com/owncloud/ocis/pull/2426
   https://github.com/owncloud/web/releases/tag/v4.1.0

* Enhancement - Use non root user for the owncloud/ocis docker image: [#2380](https://github.com/owncloud/ocis/pull/2380)

   The owncloud/ocis docker image now uses a non root user and enables you to set a
   different user with the docker `--user` parameter. The default user has the UID
   1000 is part of a group with the GID 1000.

   This is a breaking change for existing docker deployments. The permission on the
   files and folders in persistent volumes need to be changed to the UID and GID
   used for oCIS (default 1000:1000 if not changed by the user).

   https://github.com/owncloud/ocis/pull/2380

* Enhancement - Replace unmaintained jwt library: [#2386](https://github.com/owncloud/ocis/pull/2386)

   The old library
   [github.com/dgrijalva/jwt-go](https://github.com/dgrijalva/jwt-go) is
   unmaintained and was replaced by the community maintained fork
   [github.com/golang-jwt/jwt](https://github.com/golang-jwt/jwt).

   https://github.com/owncloud/ocis/pull/2386

* Enhancement - Update bleve to version 2.1.0: [#2391](https://github.com/owncloud/ocis/pull/2391)

   Updated bleve to the current version.

   https://github.com/owncloud/ocis/pull/2391

* Enhancement - Update github.com/coreos/go-oidc to v3.0.0: [#2393](https://github.com/owncloud/ocis/pull/2393)

   Updated the github.com/coreos/go-oidc library to the version 3.0.0.

   https://github.com/owncloud/ocis/pull/2393

* Enhancement - Update reva to v1.12: [#2423](https://github.com/owncloud/ocis/pull/2423)

  * Enhancement cs3org/reva#1803: Introduce new webdav spaces endpoint
  * Bugfix cs3org/reva#1819: Disable notifications
  * Enhancement cs3org/reva#1861: Add support for runtime plugins
  * Bugfix cs3org/reva#1913: Logic to restore files to readonly nodes
  * Enhancement cs3org/reva#1946: Add share manager that connects to oc10 databases
  * Bugfix cs3org/reva#1954: Fix response format of the sharees API
  * Bugfix cs3org/reva#1956: Fix trashbin listing with depth 0
  * Bugfix cs3org/reva#1957: Fix etag propagation on deletes
  * Bugfix cs3org/reva#1960: Return the updated share after updating
  * Bugfix cs3org/reva#1965 cs3org/reva#1967: Fix the file target of user and group shares
  * Bugfix cs3org/reva#1980: Propagate the etag after restoring a file version
  * Enhancement cs3org/reva#1984: Replace OpenCensus with OpenTelemetry
  * Bugfix cs3org/reva#1985: Add quota stubs
  * Bugfix cs3org/reva#1987: Fix windows build
  * Bugfix cs3org/reva#1990: Increase oc10 compatibility of owncloudsql
  * Bugfix cs3org/reva#1992: Check if symlink exists instead of spamming the console
  * Bugfix cs3org/reva#1993: fix owncloudsql GetMD

   https://github.com/owncloud/ocis/pull/2423

# Changelog for [1.10.0] (2021-08-06)

The following sections list the changes for 1.10.0.

[1.10.0]: https://github.com/owncloud/ocis/compare/v1.9.0...v1.10.0

## Summary

* Bugfix - Forward basic auth to OpenID connect token authentication endpoint: [#2095](https://github.com/owncloud/ocis/issues/2095)
* Bugfix - Log all requests in the proxy access log: [#2301](https://github.com/owncloud/ocis/pull/2301)
* Bugfix - Update glauth to 20210729125545-b9aecdfcac31: [#2336](https://github.com/owncloud/ocis/pull/2336)
* Bugfix - Improve IDP Login Accessibility: [#5376](https://github.com/owncloud/web/issues/5376)
* Change - Update ownCloud Web to v4.0.0: [#2353](https://github.com/owncloud/ocis/pull/2353)
* Enhancement - Proxy: Add claims policy selector: [#2248](https://github.com/owncloud/ocis/pull/2248)
* Enhancement - Refactor graph API: [#2277](https://github.com/owncloud/ocis/pull/2277)
* Enhancement - Add ocs cache warmup config and warn on protobuf ns conflicts: [#2328](https://github.com/owncloud/ocis/pull/2328)
* Enhancement - Use only one go.mod file for project dependencies: [#2344](https://github.com/owncloud/ocis/pull/2344)
* Enhancement - Update REVA: [#2355](https://github.com/owncloud/ocis/pull/2355)

## Details

* Bugfix - Forward basic auth to OpenID connect token authentication endpoint: [#2095](https://github.com/owncloud/ocis/issues/2095)

   When using `PROXY_ENABLE_BASIC_AUTH=true` we now forward request to the idp
   instead of trying to authenticate the request ourself.

   https://github.com/owncloud/ocis/issues/2095
   https://github.com/owncloud/ocis/issues/2094

* Bugfix - Log all requests in the proxy access log: [#2301](https://github.com/owncloud/ocis/pull/2301)

   We now use a dedicated middleware to log all requests, regardless of routing
   selector outcome. While the log now includes the remote address, the selected
   routing policy is only logged when log level is set to debug because the request
   context cannot be changed in the `directorSelectionDirector`, as per the
   `ReverseProxy.Director` documentation.

   https://github.com/owncloud/ocis/pull/2301

* Bugfix - Update glauth to 20210729125545-b9aecdfcac31: [#2336](https://github.com/owncloud/ocis/pull/2336)

  * Fixes the backend config not being passed correctly in ocis
  * Fixes a mutex being copied, leading to concurrent writes
  * Fixes UTF8 chars in filters
  * Fixes case insensitive strings

   https://github.com/owncloud/ocis/pull/2336
   https://github.com/glauth/glauth/pull/198
   https://github.com/glauth/glauth/pull/194

* Bugfix - Improve IDP Login Accessibility: [#5376](https://github.com/owncloud/web/issues/5376)

   We have addressed the feedback from the `a11y` audit and improved the IDP login
   screen accordingly.

   https://github.com/owncloud/web/issues/5376
   https://github.com/owncloud/web/issues/5377

* Change - Update ownCloud Web to v4.0.0: [#2353](https://github.com/owncloud/ocis/pull/2353)

   Tags: web

   We updated ownCloud Web to v4.0.0. Please refer to the changelog (linked) for
   details on the web release.

   https://github.com/owncloud/ocis/pull/2353
   https://github.com/owncloud/web/releases/tag/v4.0.0

* Enhancement - Proxy: Add claims policy selector: [#2248](https://github.com/owncloud/ocis/pull/2248)

   Using the proxy config file, it is now possible to let let the IdP determine the
   routing policy by sending an `ocis.routing.policy` claim. Its value will be used
   to determine the set of routes for the logged in user.

   https://github.com/owncloud/ocis/pull/2248

* Enhancement - Refactor graph API: [#2277](https://github.com/owncloud/ocis/pull/2277)

   We refactored the `/graph/v1.0/` endpoint which now relies on the internal
   access token fer authentication, getting rid of any LDAP or OIDC code to
   authenticate requests. This allows using the graph api when using basic auth or
   any other auth mechanism provided by the CS3 auth providers / reva gateway /
   ocis proxy.

   https://github.com/owncloud/ocis/pull/2277

* Enhancement - Add ocs cache warmup config and warn on protobuf ns conflicts: [#2328](https://github.com/owncloud/ocis/pull/2328)

   https://github.com/owncloud/ocis/pull/2328

* Enhancement - Use only one go.mod file for project dependencies: [#2344](https://github.com/owncloud/ocis/pull/2344)

   We now use one single go.mod file at the root of the repository rather than one
   per core extension.

   https://github.com/owncloud/ocis/pull/2344

* Enhancement - Update REVA: [#2355](https://github.com/owncloud/ocis/pull/2355)

   Update REVA from v1.10.1-0.20210730095301-fcb7a30a44a6 to
   v1.11.1-0.20210809134415-3fe79c870fb5 * Fix cs3org/reva#1978: Fix owner type is
   optional * Fix cs3org/reva#1965: fix value of file_target in shares * Fix
   cs3org/reva#1960: fix updating shares in the memory share manager * Fix
   cs3org/reva#1956: fix trashbin listing with depth 0 * Fix cs3org/reva#1957: fix
   etag propagation on deletes * Enh cs3org/reva#1861: [WIP] Runtime plugins * Fix
   cs3org/reva#1954: fix response format of the sharees API * Fix cs3org/reva#1819:
   Remove notifications key from ocs response * Enh cs3org/reva#1946: Add a share
   manager that connects to oc10 databases * Fix cs3org/reva#1899: Fix chunked
   uploads for new versions * Fix cs3org/reva#1906: Fix copy over existing resource
   * Fix cs3org/reva#1891: Delete Shared Resources as Receiver * Fix
   cs3org/reva#1907: Error when creating folder with existing name * Fix
   cs3org/reva#1937: Do not overwrite more specific matches when finding storage
   providers * Fix cs3org/reva#1939: Fix the share jail permissions in the
   decomposedfs * Fix cs3org/reva#1932: Numerous fixes to the owncloudsql storage
   driver * Fix cs3org/reva#1912: Fix response when listing versions of another
   user * Fix cs3org/reva#1910: Get user groups recursively in the cbox rest user
   driver * Fix cs3org/reva#1904: Set Content-Length to 0 when swallowing body in
   the datagateway * Fix cs3org/reva#1911: Fix version order in propfind responses
   * Fix cs3org/reva#1926: Trash Bin in oCIS Storage Operations * Fix
   cs3org/reva#1901: Fix response code when folder doesnt exist on upload * Enh
   cs3org/reva#1785: Extend app registry with AddProvider method and mimetype
   filters * Enh cs3org/reva#1938: Add methods to get and put context values * Enh
   cs3org/reva#1798: Add support for a deny-all permission on references * Enh
   cs3org/reva#1916: Generate updated protobuf bindings for EOS GRPC * Enh
   cs3org/reva#1887: Add "a" and "l" filter for grappa queries * Enh
   cs3org/reva#1919: Run gofmt before building * Enh cs3org/reva#1927: Implement
   RollbackToVersion for eosgrpc (needs a newer EOS MGM) * Enh cs3org/reva#1944:
   Implement listing supported mime types in app registry * Enh cs3org/reva#1870:
   Be defensive about wrongly quoted etags * Enh cs3org/reva#1940: Reduce memory
   usage when uploading with S3ng storage * Enh cs3org/reva#1888: Refactoring of
   the webdav code * Enh cs3org/reva#1900: Check for illegal names while uploading
   or moving files * Enh cs3org/reva#1925: Refactor listing and statting across
   providers for virtual views * Fix cs3org/reva#1883: Pass directories with
   trailing slashes to eosclient.GenerateToken * Fix cs3org/reva#1878: Improve the
   webdav error handling in the trashbin * Fix cs3org/reva#1884: Do not send body
   on failed range request * Enh cs3org/reva#1744: Add support for lightweight user
   types * Fix cs3org/reva#1904: Set Content-Length to 0 when swallowing body in
   the datagateway * Fix cs3org/reva#1899: Bugfix: Fix chunked uploads for new
   versions * Enh cs3org/reva#1888: Refactoring of the webdav code * Enh
   cs3org/reva#1887: Add "a" and "l" filter for grappa queries

   https://github.com/owncloud/ocis/pull/2355
   https://github.com/owncloud/ocis/pull/2295
   https://github.com/owncloud/ocis/pull/2314

# Changelog for [1.9.0] (2021-07-13)

The following sections list the changes for 1.9.0.

[1.9.0]: https://github.com/owncloud/ocis/compare/v1.8.0...v1.9.0

## Summary

* Bugfix - Panic when service fails to start: [#2252](https://github.com/owncloud/ocis/pull/2252)
* Bugfix - Dont use port 80 as debug for GroupsProvider: [#2271](https://github.com/owncloud/ocis/pull/2271)
* Change - Update ownCloud Web to v3.4.0: [#2276](https://github.com/owncloud/ocis/pull/2276)
* Change - Update WEB to v3.4.1: [#2283](https://github.com/owncloud/ocis/pull/2283)
* Enhancement - Remove unnecessary Service.Init(): [#1705](https://github.com/owncloud/ocis/pull/1705)
* Enhancement - Update REVA to v1.9.1-0.20210628143859-9d29c36c0c3f: [#2227](https://github.com/owncloud/ocis/pull/2227)
* Enhancement - Runtime support for cherry picking extensions: [#2229](https://github.com/owncloud/ocis/pull/2229)
* Enhancement - Add readonly mode for storagehome and storageusers: [#2230](https://github.com/owncloud/ocis/pull/2230)
* Enhancement - Update REVA to v1.9.1: [#2280](https://github.com/owncloud/ocis/pull/2280)

## Details

* Bugfix - Panic when service fails to start: [#2252](https://github.com/owncloud/ocis/pull/2252)

   Tags: runtime

   When attempting to run a service through the runtime that is currently running
   and fails to start, a race condition still redirect os Interrupt signals to a
   closed channel.

   https://github.com/owncloud/ocis/pull/2252

* Bugfix - Dont use port 80 as debug for GroupsProvider: [#2271](https://github.com/owncloud/ocis/pull/2271)

   A copy/paste error where the configuration for the groupsprovider's debug
   address was not present leaves go-micro to start the debug service in port 80 by
   default.

   https://github.com/owncloud/ocis/pull/2271

* Change - Update ownCloud Web to v3.4.0: [#2276](https://github.com/owncloud/ocis/pull/2276)

   Tags: web

   We updated ownCloud Web to v3.4.0. Please refer to the changelog (linked) for
   details on the web release.

   https://github.com/owncloud/ocis/pull/2276
   https://github.com/owncloud/web/releases/tag/v3.4.0

* Change - Update WEB to v3.4.1: [#2283](https://github.com/owncloud/ocis/pull/2283)

  * Fix [5501](https://github.com/owncloud/web/pull/5501): loading previews in the right sidebar
  * Fix [5493](https://github.com/owncloud/web/pull/5493): view options position

   https://github.com/owncloud/ocis/pull/2283
   https://github.com/owncloud/web/releases/tag/v3.4.1

* Enhancement - Remove unnecessary Service.Init(): [#1705](https://github.com/owncloud/ocis/pull/1705)

   As it turns out oCIS already calls this method. Invoking it twice would end in
   accidentally resetting values.

   https://github.com/owncloud/ocis/pull/1705

* Enhancement - Update REVA to v1.9.1-0.20210628143859-9d29c36c0c3f: [#2227](https://github.com/owncloud/ocis/pull/2227)

   https://github.com/owncloud/ocis/pull/2227

* Enhancement - Runtime support for cherry picking extensions: [#2229](https://github.com/owncloud/ocis/pull/2229)

   Support for running certain extensions supervised via cli flags. Example usage:

   ```
   > ocis server --extensions="proxy, idp, storage-metadata, accounts"
   ```

   https://github.com/owncloud/ocis/pull/2229

* Enhancement - Add readonly mode for storagehome and storageusers: [#2230](https://github.com/owncloud/ocis/pull/2230)

   To enable the readonly mode use `STORAGE_HOME_READ_ONLY=true` and
   `STORAGE_USERS_READ_ONLY=true`. Alternative: use `OCIS_STORAGE_READ_ONLY=true`

   https://github.com/owncloud/ocis/pull/2230

* Enhancement - Update REVA to v1.9.1: [#2280](https://github.com/owncloud/ocis/pull/2280)

  * Fix cs3org/reva#1843: Correct Dockerfile path for the reva CLI and alpine3.13 as builder
  * Fix cs3org/reva#1835: Cleanup owncloudsql driver
  * Fix cs3org/reva#1868: Minor fixes to the grpc/http plugin: checksum, url escaping
  * Fix cs3org/reva#1885: Fix template in eoshomewrapper to use context user rather than resource
  * Fix cs3org/reva#1833: Properly handle name collisions for deletes in the owncloud driver
  * Fix cs3org/reva#1874: Use the original file mtime during upload
  * Fix cs3org/reva#1854: Add the uid/gid to the url for eos
  * Fix cs3org/reva#1848: Fill in missing gid/uid number with nobody
  * Fix cs3org/reva#1831: Make the ocm-provider endpoint in the ocmd service unprotected
  * Fix cs3org/reva#1808: Use empty array in OCS Notifications endpoints
  * Fix cs3org/reva#1825: Raise max grpc message size
  * Fix cs3org/reva#1828: Send a proper XML header with error messages
  * Chg cs3org/reva#1828: Remove the oidc provider in order to upgrad mattn/go-sqlite3 to v1.14.7
  * Enh cs3org/reva#1834: Add API key to Mentix GOCDB connector
  * Enh cs3org/reva#1855: Minor optimization in parsing EOS ACLs
  * Enh cs3org/reva#1873: Update the EOS image tag to be for revad-eos image
  * Enh cs3org/reva#1802: Introduce list spaces
  * Enh cs3org/reva#1849: Add readonly interceptor
  * Enh cs3org/reva#1875: Simplify resource comparison
  * Enh cs3org/reva#1827: Support trashbin sub paths in the recycle API

   https://github.com/owncloud/ocis/pull/2280

# Changelog for [1.8.0] (2021-06-28)

The following sections list the changes for 1.8.0.

[1.8.0]: https://github.com/owncloud/ocis/compare/v1.7.0...v1.8.0

## Summary

* Bugfix - External storage registration used wrong config: [#2120](https://github.com/owncloud/ocis/pull/2120)
* Bugfix - Remove authentication from /status.php completely: [#2188](https://github.com/owncloud/ocis/pull/2188)
* Bugfix - Make webdav namespace configurable across services: [#2198](https://github.com/owncloud/ocis/pull/2198)
* Change - Update ownCloud Web to v3.3.0: [#2187](https://github.com/owncloud/ocis/pull/2187)
* Enhancement - Properly configure graph-explorer client registration: [#2118](https://github.com/owncloud/ocis/pull/2118)
* Enhancement - Use system default location to store TLS artefacts: [#2129](https://github.com/owncloud/ocis/pull/2129)
* Enhancement - Update REVA to v1.9: [#2205](https://github.com/owncloud/ocis/pull/2205)

## Details

* Bugfix - External storage registration used wrong config: [#2120](https://github.com/owncloud/ocis/pull/2120)

   The go-micro registry-singleton ignores the ocis configuration and defaults to
   mdns

   https://github.com/owncloud/ocis/pull/2120

* Bugfix - Remove authentication from /status.php completely: [#2188](https://github.com/owncloud/ocis/pull/2188)

   Despite requests without Authentication header being successful, requests with
   an invalid bearer token in the Authentication header were rejected in the proxy
   with an 401 unauthenticated. Now the Authentication header is completely ignored
   for the /status.php route.

   https://github.com/owncloud/client/issues/8538
   https://github.com/owncloud/ocis/pull/2188

* Bugfix - Make webdav namespace configurable across services: [#2198](https://github.com/owncloud/ocis/pull/2198)

   The WebDAV namespace is used across various services, but it was previously
   hardcoded in some of the services. This PR uses the same environment variable to
   set the config correctly across the services.

   https://github.com/owncloud/ocis/pull/2198

* Change - Update ownCloud Web to v3.3.0: [#2187](https://github.com/owncloud/ocis/pull/2187)

   Tags: web

   We updated ownCloud Web to v3.3.0. Please refer to the changelog (linked) for
   details on the web release.

   https://github.com/owncloud/ocis/pull/2187
   https://github.com/owncloud/web/releases/tag/v3.3.0

* Enhancement - Properly configure graph-explorer client registration: [#2118](https://github.com/owncloud/ocis/pull/2118)

   The client registration in the `identifier-registration.yaml` for the
   graph-explorer didn't contain `redirect_uris` nor `origins`. Both were added to
   prevent exploitation.

   https://github.com/owncloud/ocis/pull/2118

* Enhancement - Use system default location to store TLS artefacts: [#2129](https://github.com/owncloud/ocis/pull/2129)

   This used to default to the current location of the binary, which is not ideal
   after a first run as it leaves traces behind. It now uses the system's location
   for artefacts with the help of https://golang.org/pkg/os/#UserConfigDir.

   https://github.com/owncloud/ocis/pull/2129

* Enhancement - Update REVA to v1.9: [#2205](https://github.com/owncloud/ocis/pull/2205)

   This update includes * [set Content-Type
   correctly](https://github.com/cs3org/reva/pull/1750) * [Return file checksum
   available from the metadata for the EOS
   driver](https://github.com/cs3org/reva/pull/1755) * [Sort share entries
   alphabetically](https://github.com/cs3org/reva/pull/1772) * [Initial work on the
   owncloudsql driver](https://github.com/cs3org/reva/pull/1710) * [Add user ID
   cache warmup to EOS storage driver](https://github.com/cs3org/reva/pull/1774) *
   [Use UidNumber and GidNumber fields in User
   objects](https://github.com/cs3org/reva/pull/1573) * [EOS GRPC
   interface](https://github.com/cs3org/reva/pull/1471) * [switch
   references](https://github.com/cs3org/reva/pull/1721) * [remove user's uuid from
   trashbin file key](https://github.com/cs3org/reva/pull/1793) * [fix restore
   behavior of the trashbin API](https://github.com/cs3org/reva/pull/1795) *
   [eosfs: add arbitrary metadata
   support](https://github.com/cs3org/reva/pull/1811)

   https://github.com/owncloud/ocis/pull/2205
   https://github.com/owncloud/ocis/pull/2210

# Changelog for [1.7.0] (2021-06-04)

The following sections list the changes for 1.7.0.

[1.7.0]: https://github.com/owncloud/ocis/compare/v1.6.0...v1.7.0

## Summary

* Bugfix - Change the groups index to be case sensitive: [#2109](https://github.com/owncloud/ocis/pull/2109)
* Change - Update ownCloud Web to v3.2.0: [#2096](https://github.com/owncloud/ocis/pull/2096)
* Enhancement - Enable the s3ng storage driver: [#1886](https://github.com/owncloud/ocis/pull/1886)
* Enhancement - Announce user profile picture capability: [#2036](https://github.com/owncloud/ocis/pull/2036)
* Enhancement - Color contrasts on IDP/OIDC login pages: [#2088](https://github.com/owncloud/ocis/pull/2088)
* Enhancement - Update reva to v1.7.1-0.20210531093513-b74a2b156af6: [#2104](https://github.com/owncloud/ocis/pull/2104)

## Details

* Bugfix - Change the groups index to be case sensitive: [#2109](https://github.com/owncloud/ocis/pull/2109)

   Groups are considered to be case-sensitive. The index must handle them
   case-sensitive too otherwise we will have non-deterministic behavior while
   editing or deleting groups.

   https://github.com/owncloud/ocis/pull/2109

* Change - Update ownCloud Web to v3.2.0: [#2096](https://github.com/owncloud/ocis/pull/2096)

   Tags: web

   We updated ownCloud Web to v3.2.0. Please refer to the changelog (linked) for
   details on the web release.

   https://github.com/owncloud/ocis/pull/2096
   https://github.com/owncloud/web/releases/tag/v3.2.0

* Enhancement - Enable the s3ng storage driver: [#1886](https://github.com/owncloud/ocis/pull/1886)

   We made it possible to use the new s3ng storage driver by adding according
   commandline flags and environment variables.

   https://github.com/owncloud/ocis/pull/1886

* Enhancement - Announce user profile picture capability: [#2036](https://github.com/owncloud/ocis/pull/2036)

   Added a new capability (through https://github.com/cs3org/reva/pull/1694) to
   prevent the web frontend from fetching (nonexistent) user avatar profile
   pictures which added latency & console errors.

   https://github.com/owncloud/ocis/pull/2036

* Enhancement - Color contrasts on IDP/OIDC login pages: [#2088](https://github.com/owncloud/ocis/pull/2088)

   We have updated the color contrasts on the IDP pages in order to improve
   accessibility.

   https://github.com/owncloud/ocis/pull/2088

* Enhancement - Update reva to v1.7.1-0.20210531093513-b74a2b156af6: [#2104](https://github.com/owncloud/ocis/pull/2104)

   This reva update includes: * [fix move in the owncloud storage
   driver](https://github.com/cs3org/reva/pull/1696) * [add checksum header to the
   tus preflight response](https://github.com/cs3org/reva/pull/1702) * [Add
   reliability calculations support to
   Mentix](https://github.com/cs3org/reva/pull/1649) * [fix response format when
   accepting shares](https://github.com/cs3org/reva/pull/1724) * [Datatx
   createtransfershare](https://github.com/cs3org/reva/pull/1725)

   https://github.com/owncloud/ocis/issues/2102
   https://github.com/owncloud/ocis/pull/2104

# Changelog for [1.6.0] (2021-05-12)

The following sections list the changes for 1.6.0.

[1.6.0]: https://github.com/owncloud/ocis/compare/v1.5.0...v1.6.0

## Summary

* Bugfix - Fix STORAGE_METADATA_ROOT default value override: [#1956](https://github.com/owncloud/ocis/pull/1956)
* Bugfix - Stop the supervisor if a service fails to start: [#1963](https://github.com/owncloud/ocis/pull/1963)
* Change - Update ownCloud Web to v3.1.0: [#2045](https://github.com/owncloud/ocis/pull/2045)
* Enhancement - User Deprovisioning for the OCS API: [#1962](https://github.com/owncloud/ocis/pull/1962)
* Enhancement - Use oc-select: [#1979](https://github.com/owncloud/ocis/pull/1979)
* Enhancement - Support thumbnails for txt files: [#1988](https://github.com/owncloud/ocis/pull/1988)
* Enhancement - Introduce login form with h1 tag for screen readers only: [#1991](https://github.com/owncloud/ocis/pull/1991)
* Enhancement - Added dictionary files: [#2003](https://github.com/owncloud/ocis/pull/2003)
* Enhancement - Update reva to v1.7.1-0.20210430154404-69bd21f2cc97: [#2010](https://github.com/owncloud/ocis/pull/2010)
* Enhancement - Set SameSite settings to Strict for Web: [#2019](https://github.com/owncloud/ocis/pull/2019)
* Enhancement - Update reva to v1.7.1-0.20210507160327-e2c3841d0dbc: [#2044](https://github.com/owncloud/ocis/pull/2044)

## Details

* Bugfix - Fix STORAGE_METADATA_ROOT default value override: [#1956](https://github.com/owncloud/ocis/pull/1956)

   The way the value was being set ensured that it was NOT being overridden where
   it should have been. This patch ensures the correct loading order of values.

   https://github.com/owncloud/ocis/pull/1956

* Bugfix - Stop the supervisor if a service fails to start: [#1963](https://github.com/owncloud/ocis/pull/1963)

   Steps to make the supervisor fail:

   `PROXY_HTTP_ADDR=0.0.0.0:9144 bin/ocis server`

   https://github.com/owncloud/ocis/pull/1963

* Change - Update ownCloud Web to v3.1.0: [#2045](https://github.com/owncloud/ocis/pull/2045)

   Tags: web

   We updated ownCloud Web to v3.1.0. Please refer to the changelog (linked) for
   details on the web release.

   https://github.com/owncloud/ocis/pull/2045
   https://github.com/owncloud/web/releases/tag/v3.1.0

* Enhancement - User Deprovisioning for the OCS API: [#1962](https://github.com/owncloud/ocis/pull/1962)

   Use the CS3 API and Reva to deprovision users completely.

   Two new environment variables introduced:

   ```
   OCS_IDM_ADDRESS
   OCS_STORAGE_USERS_DRIVER
   ```

   `OCS_IDM_ADDRESS` is also an alias for `OCIS_URL`; allows the OCS service to
   mint jwt tokens for the authenticated user that will be read by the reva
   authentication middleware.

   `OCS_STORAGE_USERS_DRIVER` determines how a user is deprovisioned. This kind of
   behavior is needed since every storage driver deals with deleting differently.

   https://github.com/owncloud/ocis/pull/1962

* Enhancement - Use oc-select: [#1979](https://github.com/owncloud/ocis/pull/1979)

   Replace oc-drop with oc select in settings

   https://github.com/owncloud/ocis/pull/1979

* Enhancement - Support thumbnails for txt files: [#1988](https://github.com/owncloud/ocis/pull/1988)

   Implemented support for thumbnails for txt files in the thumbnails service.

   https://github.com/owncloud/ocis/pull/1988

* Enhancement - Introduce login form with h1 tag for screen readers only: [#1991](https://github.com/owncloud/ocis/pull/1991)

   https://github.com/owncloud/ocis/pull/1991

* Enhancement - Added dictionary files: [#2003](https://github.com/owncloud/ocis/pull/2003)

   Added the dictionary.js file for package settings and accounts which contains
   strings that should be synced to transifex but not exist in the UI directly.

   https://github.com/owncloud/ocis/pull/2003

* Enhancement - Update reva to v1.7.1-0.20210430154404-69bd21f2cc97: [#2010](https://github.com/owncloud/ocis/pull/2010)

  * Fix recycle to different locations (https://github.com/cs3org/reva/pull/1541)
  * Fix user share as grantee in json backend (https://github.com/cs3org/reva/pull/1650)
  * Introduce named services (https://github.com/cs3org/reva/pull/1509)
  * Improve json marshalling of share protobuf messages (https://github.com/cs3org/reva/pull/1655)
  * Cache resources from share getter methods in OCS (https://github.com/cs3org/reva/pull/1643)
  * Fix public file shares (https://github.com/cs3org/reva/pull/1666)

   https://github.com/owncloud/ocis/pull/2010

* Enhancement - Set SameSite settings to Strict for Web: [#2019](https://github.com/owncloud/ocis/pull/2019)

   Changed SameSite settings to Strict for Web to prevent warnings in Firefox

   https://github.com/owncloud/ocis/pull/2019

* Enhancement - Update reva to v1.7.1-0.20210507160327-e2c3841d0dbc: [#2044](https://github.com/owncloud/ocis/pull/2044)

  * Add user profile picture to capabilities (https://github.com/cs3org/reva/pull/1694)
  * Mint scope-based access tokens for RBAC (https://github.com/cs3org/reva/pull/1669)
  * Add cache warmup strategy for OCS resource infos (https://github.com/cs3org/reva/pull/1664)
  * Filter shares based on type in OCS (https://github.com/cs3org/reva/pull/1683)

   https://github.com/owncloud/ocis/pull/2044

# Changelog for [1.5.0] (2021-04-21)

The following sections list the changes for 1.5.0.

[1.5.0]: https://github.com/owncloud/ocis/compare/v1.4.0...v1.5.0

## Summary

* Bugfix - Fixes "unaligned 64-bit atomic operation" panic on 32-bit ARM: [#1888](https://github.com/owncloud/ocis/pull/1888)
* Change - Make Protobuf package names unique: [#1875](https://github.com/owncloud/ocis/pull/1875)
* Change - Update ownCloud Web to v3.0.0: [#1938](https://github.com/owncloud/ocis/pull/1938)
* Enhancement - Update reva to v1.6.1-0.20210414111318-a4b5148cbfb2: [#1872](https://github.com/owncloud/ocis/pull/1872)
* Enhancement - Change default path for thumbnails: [#1892](https://github.com/owncloud/ocis/pull/1892)
* Enhancement - Add config for public share SQL driver: [#1916](https://github.com/owncloud/ocis/pull/1916)
* Enhancement - Add option to reading registry rules from json file: [#1917](https://github.com/owncloud/ocis/pull/1917)
* Enhancement - Remove dead runtime code: [#1923](https://github.com/owncloud/ocis/pull/1923)
* Enhancement - Parse config on supervised mode with run subcommand: [#1931](https://github.com/owncloud/ocis/pull/1931)
* Enhancement - Update ODS in accounts & settings extension: [#1934](https://github.com/owncloud/ocis/pull/1934)

## Details

* Bugfix - Fixes "unaligned 64-bit atomic operation" panic on 32-bit ARM: [#1888](https://github.com/owncloud/ocis/pull/1888)

   Sync/cache had uint64s that were not 64-bit aligned causing panics on 32-bit
   systems during atomic access

   https://github.com/owncloud/ocis/issues/1887
   https://github.com/owncloud/ocis/pull/1888

* Change - Make Protobuf package names unique: [#1875](https://github.com/owncloud/ocis/pull/1875)

   Introduce unique `package` and `go_package` names for our Protobuf definitions

   https://github.com/owncloud/ocis/pull/1875

* Change - Update ownCloud Web to v3.0.0: [#1938](https://github.com/owncloud/ocis/pull/1938)

   Tags: web

   We updated ownCloud Web to v3.0.0. Please refer to the changelog (linked) for
   details on the web release.

   https://github.com/owncloud/ocis/pull/1938
   https://github.com/owncloud/web/releases/tag/v3.0.0

* Enhancement - Update reva to v1.6.1-0.20210414111318-a4b5148cbfb2: [#1872](https://github.com/owncloud/ocis/pull/1872)

  * enforce quota (https://github.com/cs3org/reva/pull/1557)
  * Make additional info attribute configurable (https://github.com/cs3org/reva/pull/1588)
  * check ENOTDIR for readlink (https://github.com/cs3org/reva/pull/1597)
  * Add wrappers for EOS and EOS Home storage drivers (https://github.com/cs3org/reva/pull/1624)
  * eos: fixes for enabling file sharing (https://github.com/cs3org/reva/pull/1619)
  * implement checksums in the owncloud storage driver (https://github.com/cs3org/reva/pull/1629)

   https://github.com/owncloud/ocis/pull/1872

* Enhancement - Change default path for thumbnails: [#1892](https://github.com/owncloud/ocis/pull/1892)

   Changes the default path for thumbnails from `<os tmp dir>/ocis-thumbnails` to
   `/var/tmp/ocis/thumbnails`

   https://github.com/owncloud/ocis/issues/1891
   https://github.com/owncloud/ocis/pull/1892

* Enhancement - Add config for public share SQL driver: [#1916](https://github.com/owncloud/ocis/pull/1916)

   https://github.com/owncloud/ocis/pull/1916

* Enhancement - Add option to reading registry rules from json file: [#1917](https://github.com/owncloud/ocis/pull/1917)

   https://github.com/owncloud/ocis/pull/1917

* Enhancement - Remove dead runtime code: [#1923](https://github.com/owncloud/ocis/pull/1923)

   When moving from the old runtime to the new one there were lots of files left
   behind that are essentially dead code and should be removed. The original code
   lives here github.com/refs/pman/ if someone finds it interesting to read.

   https://github.com/owncloud/ocis/pull/1923

* Enhancement - Parse config on supervised mode with run subcommand: [#1931](https://github.com/owncloud/ocis/pull/1931)

   Currently it is not possible to parse a single config file from an extension
   when running on supervised mode.

   https://github.com/owncloud/ocis/pull/1931

* Enhancement - Update ODS in accounts & settings extension: [#1934](https://github.com/owncloud/ocis/pull/1934)

   The accounts and settings extensions were updated to reflect the latest changes
   in the ownCloud design system. In addition, a couple of quick wins in terms of
   accessibility are included.

   https://github.com/owncloud/ocis/pull/1934

# Changelog for [1.4.0] (2021-03-30)

The following sections list the changes for 1.4.0.

[1.4.0]: https://github.com/owncloud/ocis/compare/v1.3.0...v1.4.0

## Summary

* Bugfix - Fix thumbnail generation for jpegs: [#1785](https://github.com/owncloud/ocis/pull/1785)
* Change - Update ownCloud Web to v2.1.0: [#1870](https://github.com/owncloud/ocis/pull/1870)
* Enhancement - Update reva to v1.6.1-0.20210326165326-e8a00d9b2368: [#1683](https://github.com/owncloud/ocis/pull/1683)
* Enhancement - Clarify expected failures: [#1790](https://github.com/owncloud/ocis/pull/1790)
* Enhancement - Generate thumbnails for .gif files: [#1791](https://github.com/owncloud/ocis/pull/1791)
* Enhancement - Add focus to input elements on login page: [#1792](https://github.com/owncloud/ocis/pull/1792)
* Enhancement - Improve accessibility to input elements on login page: [#1794](https://github.com/owncloud/ocis/pull/1794)
* Enhancement - Replace special character in login page title with a regular minus: [#1813](https://github.com/owncloud/ocis/pull/1813)
* Enhancement - File Logging: [#1816](https://github.com/owncloud/ocis/pull/1816)
* Enhancement - Tracing Refactor: [#1819](https://github.com/owncloud/ocis/pull/1819)
* Enhancement - Runtime Hostname and Port are now configurable: [#1822](https://github.com/owncloud/ocis/pull/1822)
* Enhancement - Add new build targets: [#1824](https://github.com/owncloud/ocis/pull/1824)

## Details

* Bugfix - Fix thumbnail generation for jpegs: [#1785](https://github.com/owncloud/ocis/pull/1785)

   Images with the extension `.jpeg` were not properly supported.

   https://github.com/owncloud/ocis/issues/1490
   https://github.com/owncloud/ocis/pull/1785

* Change - Update ownCloud Web to v2.1.0: [#1870](https://github.com/owncloud/ocis/pull/1870)

   Tags: web

   We updated ownCloud Web to v2.1.0. Please refer to the changelog (linked) for
   details on the web release.

   https://github.com/owncloud/ocis/pull/1870
   https://github.com/owncloud/web/releases/tag/v2.1.0

* Enhancement - Update reva to v1.6.1-0.20210326165326-e8a00d9b2368: [#1683](https://github.com/owncloud/ocis/pull/1683)

  * quota querying and tree accounting [cs3org/reva#1405](https://github.com/cs3org/reva/pull/1405)
  * Fix webdav file versions endpoint bugs [cs3org/reva#1526](https://github.com/cs3org/reva/pull/1526)
  * Fix etag changing only once a second [cs3org/reva#1576](https://github.com/cs3org/reva/pull/1576)
  * Trashbin API parity [cs3org/reva#1552](https://github.com/cs3org/reva/pull/1552)
  * Signature authentication for public links [cs3org/reva#1590](https://github.com/cs3org/reva/pull/1590)

   https://github.com/owncloud/ocis/pull/1683
   https://github.com/cs3org/reva/pull/1405
   https://github.com/owncloud/ocis/pull/1861

* Enhancement - Clarify expected failures: [#1790](https://github.com/owncloud/ocis/pull/1790)

   Some features, while covered by the ownCloud 10 acceptance tests, will not be
   implemented for now: - blacklisted / ignored files, because ocis does not need
   to blacklist `.htaccess` files - `OC-LazyOps` support was [removed from the
   clients](https://github.com/owncloud/client/pull/8398). We are thinking about [a
   state machine for uploads to properly solve that scenario and also list the
   state of files in progress in the web
   ui](https://github.com/owncloud/ocis/issues/214). The expected failures files
   now have a dedicated _Won't fix_ section for these items.

   https://github.com/owncloud/ocis/issues/214
   https://github.com/owncloud/ocis/pull/1790
   https://github.com/owncloud/client/pull/8398

* Enhancement - Generate thumbnails for .gif files: [#1791](https://github.com/owncloud/ocis/pull/1791)

   Added support for gifs to the thumbnails service.

   https://github.com/owncloud/ocis/pull/1791

* Enhancement - Add focus to input elements on login page: [#1792](https://github.com/owncloud/ocis/pull/1792)

   https://github.com/owncloud/web/issues/4322
   https://github.com/owncloud/ocis/pull/1792

* Enhancement - Improve accessibility to input elements on login page: [#1794](https://github.com/owncloud/ocis/pull/1794)

   https://github.com/owncloud/web/issues/4319
   https://github.com/owncloud/ocis/pull/1794
   https://github.com/owncloud/ocis/pull/1811

* Enhancement - Replace special character in login page title with a regular minus: [#1813](https://github.com/owncloud/ocis/pull/1813)

   https://github.com/owncloud/ocis/pull/1813

* Enhancement - File Logging: [#1816](https://github.com/owncloud/ocis/pull/1816)

   When running supervised, support for configuring all logs to a single log file:
   `OCIS_LOG_FILE=/Users/foo/bar/ocis.log MICRO_REGISTRY=etcd bin/ocis server`

   Supports directing log from single extensions to a log file:
   `PROXY_LOG_FILE=/Users/foo/bar/proxy.log MICRO_REGISTRY=etcd bin/ocis proxy`

   https://github.com/owncloud/ocis/pull/1816

* Enhancement - Tracing Refactor: [#1819](https://github.com/owncloud/ocis/pull/1819)

   Centralize tracing handling per extension.

   https://github.com/owncloud/ocis/pull/1819

* Enhancement - Runtime Hostname and Port are now configurable: [#1822](https://github.com/owncloud/ocis/pull/1822)

   Without any configuration the ocis runtime will start on `localhost:9250` unless
   specified otherwise. Usage:

   - `OCIS_RUNTIME_PORT=6061 bin/ocis server` - overrides the oCIS runtime and
   starts on port 6061 - `OCIS_RUNTIME_PORT=6061 bin/ocis list` - lists running
   extensions for the runtime on `localhost:6061`

   All subcommands are updated and expected to work with the following environment
   variables:

   ```
   OCIS_RUNTIME_HOST
   OCIS_RUNTIME_PORT
   ```

   https://github.com/owncloud/ocis/pull/1822

* Enhancement - Add new build targets: [#1824](https://github.com/owncloud/ocis/pull/1824)

   Make build target `build` used to build a binary twice, the second occurrence
   having symbols for debugging. We split this step in two and added `build-all`
   and `build-debug` targets.

   - `build-all` now behaves as the previous `build` target, it will generate 2
   binaries, one for debug. - `build-debug` will build a single binary for
   debugging.

   https://github.com/owncloud/ocis/pull/1824

# Changelog for [1.3.0] (2021-03-09)

The following sections list the changes for 1.3.0.

[1.3.0]: https://github.com/owncloud/ocis/compare/v1.2.0...v1.3.0

## Summary

* Bugfix - Fix accounts initialization: [#1696](https://github.com/owncloud/ocis/pull/1696)
* Bugfix - Fix the ttl of the authentication middleware cache: [#1699](https://github.com/owncloud/ocis/pull/1699)
* Bugfix - Add missing gateway config: [#1716](https://github.com/owncloud/ocis/pull/1716)
* Bugfix - Purposely delay accounts service startup: [#1734](https://github.com/owncloud/ocis/pull/1734)
* Change - Update ownCloud Web to v2.0.1: [#1683](https://github.com/owncloud/ocis/pull/1683)
* Change - Update ownCloud Web to v2.0.2: [#1776](https://github.com/owncloud/ocis/pull/1776)
* Enhancement - Update go-micro to v3.5.1-0.20210217182006-0f0ace1a44a9: [#1670](https://github.com/owncloud/ocis/pull/1670)
* Enhancement - Update reva to v1.6.1-0.20210223065028-53f39499762e: [#1683](https://github.com/owncloud/ocis/pull/1683)
* Enhancement - Add initial nats and kubernetes registry support: [#1697](https://github.com/owncloud/ocis/pull/1697)
* Enhancement - Remove the JWT from the log: [#1758](https://github.com/owncloud/ocis/pull/1758)

## Details

* Bugfix - Fix accounts initialization: [#1696](https://github.com/owncloud/ocis/pull/1696)

   Originally the accounts service relies on both the `settings` and
   `storage-metadata` to be up and running at the moment it starts. This is an
   antipattern as it will cause the entire service to panic if the dependants are
   not present.

   We inverted this dependency and moved the default initialization data (i.e:
   creating roles, permissions, settings bundles) and instead of notifying the
   settings service that the account has to provide with such options, the settings
   is instead initialized with the options the accounts rely on. Essentially saving
   bandwidth as there is no longer a gRPC call to the settings service.

   For the `storage-metadata` a retry mechanism was added that retries by default
   20 times to fetch the `com.owncloud.storage.metadata` from the service registry
   every `500` milliseconds. If this retry expires the accounts panics, as its
   dependency on the `storage-metadata` service cannot be resolved.

   We also introduced a client wrapper that acts as middleware between a client and
   a server. For more information on how it works further read
   [here](https://github.com/sony/gobreaker)

   https://github.com/owncloud/ocis/pull/1696

* Bugfix - Fix the ttl of the authentication middleware cache: [#1699](https://github.com/owncloud/ocis/pull/1699)

   The authentication cache ttl was multiplied with `time.Second` multiple times.
   This resulted in a ttl that was not intended.

   https://github.com/owncloud/ocis/pull/1699

* Bugfix - Add missing gateway config: [#1716](https://github.com/owncloud/ocis/pull/1716)

   The auth provider `ldap` and `oidc` drivers now need to be able talk to the reva
   gateway. We added the `gatewayscv` to the config that is passed to reva.

   https://github.com/owncloud/ocis/pull/1716

* Bugfix - Purposely delay accounts service startup: [#1734](https://github.com/owncloud/ocis/pull/1734)

   As it turns out the race condition between `accounts <-> storage-metadata` still
   remains. This PR is a hotfix, and it should be followed up with a proper fix.
   Either:

   - block the accounts' initialization until the storage metadata is ready (using
   the registry) or - allow the accounts service to initialize and use a message
   broker to signal the accounts the metadata storage is ready to receive requests.

   https://github.com/owncloud/ocis/pull/1734

* Change - Update ownCloud Web to v2.0.1: [#1683](https://github.com/owncloud/ocis/pull/1683)

   Tags: web

   We updated ownCloud Web to v2.0.1. Please refer to the changelog (linked) for
   details on the web release.

   https://github.com/owncloud/ocis/pull/1683
   https://github.com/owncloud/web/releases/tag/v2.0.1

* Change - Update ownCloud Web to v2.0.2: [#1776](https://github.com/owncloud/ocis/pull/1776)

   Tags: web

   We updated ownCloud Web to v2.0.2. Please refer to the changelog (linked) for
   details on the web release.

   https://github.com/owncloud/ocis/pull/1776
   https://github.com/owncloud/web/releases/tag/v2.0.2

* Enhancement - Update go-micro to v3.5.1-0.20210217182006-0f0ace1a44a9: [#1670](https://github.com/owncloud/ocis/pull/1670)

   - We updated from go micro v2 (v2.9.1) go-micro v3 (v3.5.1 edge). - oCIS runtime
   is now aware of `MICRO_LOG_LEVEL` and is set to `error` by default. This
   decision was made because ownCloud, as framework builders, want to log
   everything oCIS related and hide everything unrelated by default. It can be
   re-enabled by setting it to a log level other than `error`. i.e:
   `MICRO_LOG_LEVEL=info`. - Updated `protoc-gen-micro` to the [latest
   version](https://github.com/asim/go-micro/tree/master/cmd/protoc-gen-micro). -
   We're using Prometheus wrappers from go-micro.

   https://github.com/owncloud/ocis/pull/1670
   https://github.com/asim/go-micro/pull/2126

* Enhancement - Update reva to v1.6.1-0.20210223065028-53f39499762e: [#1683](https://github.com/owncloud/ocis/pull/1683)

  * quota querying and tree accounting [cs3org/reva#1405](https://github.com/cs3org/reva/pull/1405)

   https://github.com/owncloud/ocis/pull/1683
   https://github.com/cs3org/reva/pull/1405

* Enhancement - Add initial nats and kubernetes registry support: [#1697](https://github.com/owncloud/ocis/pull/1697)

   We added initial support to use nats and kubernetes as a service registry using
   `MICRO_REGISTRY=nats` and `MICRO_REGISTRY=kubernetes` respectively. Multiple
   nodes can be given with `MICRO_REGISTRY_ADDRESS=1.2.3.4,5.6.7.8,9.10.11.12`.

   https://github.com/owncloud/ocis/pull/1697

* Enhancement - Remove the JWT from the log: [#1758](https://github.com/owncloud/ocis/pull/1758)

   We were logging the JWT in some places. Secrets should not be exposed in logs so
   it got removed.

   https://github.com/owncloud/ocis/pull/1758

# Changelog for [1.2.0] (2021-02-17)

The following sections list the changes for 1.2.0.

[1.2.0]: https://github.com/owncloud/ocis/compare/v1.1.0...v1.2.0

## Summary

* Bugfix - Check if roles are present in user object before looking those up: [#1388](https://github.com/owncloud/ocis/pull/1388)
* Bugfix - Fix etcd address configuration: [#1546](https://github.com/owncloud/ocis/pull/1546)
* Bugfix - Fix thumbnail generation when using different idp: [#1624](https://github.com/owncloud/ocis/issues/1624)
* Bugfix - Remove unimplemented config file option for oCIS root command: [#1636](https://github.com/owncloud/ocis/pull/1636)
* Change - Move runtime code on refs/pman over to owncloud/ocis/ocis: [#1483](https://github.com/owncloud/ocis/pull/1483)
* Change - Initial release of graph and graph explorer: [#1594](https://github.com/owncloud/ocis/pull/1594)
* Change - Update ownCloud Web to v2.0.0: [#1661](https://github.com/owncloud/ocis/pull/1661)
* Enhancement - Introduce ADR: [#1042](https://github.com/owncloud/ocis/pull/1042)
* Enhancement - Functionality to map home directory to different storage providers: [#1186](https://github.com/owncloud/ocis/pull/1186)
* Enhancement - Use a default protocol parameter instead of explicitly disabling tus: [#1331](https://github.com/cs3org/reva/pull/1331)
* Enhancement - Switch to opencontainers annotation scheme: [#1381](https://github.com/owncloud/ocis/pull/1381)
* Enhancement - Update reva to v1.5.2-0.20210125114636-0c10b333ee69: [#1482](https://github.com/owncloud/ocis/pull/1482)
* Enhancement - Migrate ocis-graph to ocis monorepo: [#1594](https://github.com/owncloud/ocis/pull/1594)
* Enhancement - Migrate ocis-graph-explorer to ocis monorepo: [#1596](https://github.com/owncloud/ocis/pull/1596)
* Enhancement - Make use of new design-system oc-table: [#1597](https://github.com/owncloud/ocis/pull/1597)
* Enhancement - Enable group sharing and add config for sharing SQL driver: [#1626](https://github.com/owncloud/ocis/pull/1626)

## Details

* Bugfix - Check if roles are present in user object before looking those up: [#1388](https://github.com/owncloud/ocis/pull/1388)

   https://github.com/owncloud/ocis/pull/1388

* Bugfix - Fix etcd address configuration: [#1546](https://github.com/owncloud/ocis/pull/1546)

   The etcd server address in `MICRO_REGISTRY_ADDRESS` was not picked up when etcd
   was set as service discovery registry `MICRO_REGISTRY=etcd`. Therefore etcd was
   only working if available on localhost / 127.0.0.1.

   https://github.com/owncloud/ocis/pull/1546

* Bugfix - Fix thumbnail generation when using different idp: [#1624](https://github.com/owncloud/ocis/issues/1624)

   The thumbnail service was relying on a konnectd specific field in the access
   token. This logic was now replaced by a service parameter for the username.

   https://github.com/owncloud/ocis/issues/1624
   https://github.com/owncloud/ocis/pull/1628

* Bugfix - Remove unimplemented config file option for oCIS root command: [#1636](https://github.com/owncloud/ocis/pull/1636)

   https://github.com/owncloud/ocis/pull/1636

* Change - Move runtime code on refs/pman over to owncloud/ocis/ocis: [#1483](https://github.com/owncloud/ocis/pull/1483)

   Tags: ocis, runtime

   Currently, the runtime is under the private account of an oCIS developer. For
   future-proofing we don't want oCIS mission critical components to depend on
   external repositories, so we're including refs/pman module as an oCIS package
   instead.

   https://github.com/owncloud/ocis/pull/1483

* Change - Initial release of graph and graph explorer: [#1594](https://github.com/owncloud/ocis/pull/1594)

   Tags: graph, graph-explorer

   We brought initial basic Graph and Graph-Explorer support for the ownCloud
   Infinite Scale project.

   https://github.com/owncloud/ocis/pull/1594
   https://github.com/owncloud/ocis-graph-explorer/pull/3

* Change - Update ownCloud Web to v2.0.0: [#1661](https://github.com/owncloud/ocis/pull/1661)

   Tags: web

   We updated ownCloud Web to v2.0.0. Please refer to the changelog (linked) for
   details on the web release.

   https://github.com/owncloud/ocis/pull/1661
   https://github.com/owncloud/web/releases/tag/v2.0.0

* Enhancement - Introduce ADR: [#1042](https://github.com/owncloud/ocis/pull/1042)

   We will keep track of [Architectural Decision Records using
   Markdown](https://adr.github.io/madr/) in `/docs/adr`.

   https://github.com/owncloud/ocis/pull/1042

* Enhancement - Functionality to map home directory to different storage providers: [#1186](https://github.com/owncloud/ocis/pull/1186)

   We added a parameter in reva that allows us to redirect /home requests to
   different storage providers based on a mapping derived from the user attributes,
   which was previously not possible since we hardcode the /home path for all
   users. For example, having its value as `/home/{{substr 0 1 .Username}}` can be
   used to redirect home requests for different users to different storage
   providers.

   https://github.com/owncloud/ocis/pull/1186
   https://github.com/cs3org/reva/pull/1142

* Enhancement - Use a default protocol parameter instead of explicitly disabling tus: [#1331](https://github.com/cs3org/reva/pull/1331)

   https://github.com/cs3org/reva/pull/1331
   https://github.com/owncloud/ocis/pull/1374

* Enhancement - Switch to opencontainers annotation scheme: [#1381](https://github.com/owncloud/ocis/pull/1381)

   Switch docker image annotation scheme to org.opencontainers standard because
   org.label-schema is depreciated.

   https://github.com/owncloud/ocis/pull/1381

* Enhancement - Update reva to v1.5.2-0.20210125114636-0c10b333ee69: [#1482](https://github.com/owncloud/ocis/pull/1482)

  * initial checksum support for ocis [cs3org/reva#1400](https://github.com/cs3org/reva/pull/1400)
  * Use updated etag of home directory even if it is cached [cs3org/reva#1416](https://github.com/cs3org/reva/pull/#1416)
  * Indicate in EOS containers that TUS is not supported [cs3org/reva#1415](https://github.com/cs3org/reva/pull/#1415)
  * Get status code from recycle response [cs3org/reva#1408](https://github.com/cs3org/reva/pull/#1408)

   https://github.com/owncloud/ocis/pull/1482
   https://github.com/cs3org/reva/pull/1400
   https://github.com/cs3org/reva/pull/1416
   https://github.com/cs3org/reva/pull/1415
   https://github.com/cs3org/reva/pull/1408

* Enhancement - Migrate ocis-graph to ocis monorepo: [#1594](https://github.com/owncloud/ocis/pull/1594)

   Tags: ocis, ocis-graph

   Ocis-graph was not migrated during the monorepo conversion.

   https://github.com/owncloud/ocis/pull/1594

* Enhancement - Migrate ocis-graph-explorer to ocis monorepo: [#1596](https://github.com/owncloud/ocis/pull/1596)

   Tags: ocis, ocis-graph-explorer

   Ocis-graph-explorer was not migrated during the monorepo conversion.

   https://github.com/owncloud/ocis/pull/1596

* Enhancement - Make use of new design-system oc-table: [#1597](https://github.com/owncloud/ocis/pull/1597)

   Tags: ui, accounts

   The design-system table component has changed the way it's used. We updated
   accounts-ui to use the new 'oc-table-simple' component.

   https://github.com/owncloud/ocis/pull/1597

* Enhancement - Enable group sharing and add config for sharing SQL driver: [#1626](https://github.com/owncloud/ocis/pull/1626)

   This PR adds config to support sharing with groups. It also introduces a
   breaking change for the CS3APIs definitions since grantees can now refer to both
   users as well as groups. Since we store the grantee information in a json file,
   `/var/tmp/ocis/storage/shares.json`, its previous version needs to be removed as
   we won't be able to unmarshal data corresponding to the previous definitions.

   https://github.com/owncloud/ocis/pull/1626
   https://github.com/cs3org/reva/pull/1453

# Changelog for [1.1.0] (2021-01-22)

The following sections list the changes for 1.1.0.

[1.1.0]: https://github.com/owncloud/ocis/compare/v1.0.0...v1.1.0

## Summary

* Change - Disable pretty logging by default: [#1133](https://github.com/owncloud/ocis/pull/1133)
* Change - Update ownCloud Web to v1.0.1: [#1191](https://github.com/owncloud/ocis/pull/1191)
* Change - Generate cryptographically secure state token: [#1203](https://github.com/owncloud/ocis/pull/1203)
* Change - Move k6 to cdperf: [#1358](https://github.com/owncloud/ocis/pull/1358)
* Change - Update go version: [#1364](https://github.com/owncloud/ocis/pull/1364)
* Change - Add "expose" information to docker images: [#1366](https://github.com/owncloud/ocis/pull/1366)
* Change - Add "volume" declaration to docker images: [#1375](https://github.com/owncloud/ocis/pull/1375)
* Enhancement - Add OCIS_URL env var: [#1148](https://github.com/owncloud/ocis/pull/1148)
* Enhancement - Update reva to v1.4.1-0.20210111080247-f2b63bfd6825: [#1194](https://github.com/owncloud/ocis/pull/1194)
* Enhancement - Add named locks and refactor cache: [#1212](https://github.com/owncloud/ocis/pull/1212)
* Enhancement - Use sync.cache for roles cache: [#1367](https://github.com/owncloud/ocis/pull/1367)
* Enhancement - Update reva to v1.5.1: [#1372](https://github.com/owncloud/ocis/pull/1372)

## Details

* Change - Disable pretty logging by default: [#1133](https://github.com/owncloud/ocis/pull/1133)

   Tags: ocis

   Disable pretty logging default for performance reasons.

   https://github.com/owncloud/ocis/pull/1133

* Change - Update ownCloud Web to v1.0.1: [#1191](https://github.com/owncloud/ocis/pull/1191)

   Tags: web

   We updated ownCloud Web to v1.0.1. Please refer to the changelog (linked) for
   details on the web release.

   https://github.com/owncloud/ocis/pull/1191
   https://github.com/owncloud/web/releases/tag/v1.0.1

* Change - Generate cryptographically secure state token: [#1203](https://github.com/owncloud/ocis/pull/1203)

   Replaced Math.random with a cryptographically secure way to generate the oidc
   state token using the javascript crypto api.

   https://github.com/owncloud/ocis/pull/1203
   https://developer.mozilla.org/en-US/docs/Web/API/Crypto/getRandomValues
   https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Math/random

* Change - Move k6 to cdperf: [#1358](https://github.com/owncloud/ocis/pull/1358)

   Tags: performance, testing, k6

   The ownCloud performance tests can not only be used to test oCIS. This is why we
   have decided to move the k6 tests to https://github.com/owncloud/cdperf

   https://github.com/owncloud/ocis/pull/1358

* Change - Update go version: [#1364](https://github.com/owncloud/ocis/pull/1364)

   Tags: go

   Update go from 1.13 to 1.15

   https://github.com/owncloud/ocis/pull/1364

* Change - Add "expose" information to docker images: [#1366](https://github.com/owncloud/ocis/pull/1366)

   Tags: docker

   Add "expose" information to docker images. Docker users will now see that we
   offer services on port 9200.

   https://github.com/owncloud/ocis/pull/1366

* Change - Add "volume" declaration to docker images: [#1375](https://github.com/owncloud/ocis/pull/1375)

   Tags: docker

   Add "volume" declaration to docker images. This makes it easier for Docker users
   to see where oCIS stores data.

   https://github.com/owncloud/ocis/pull/1375

* Enhancement - Add OCIS_URL env var: [#1148](https://github.com/owncloud/ocis/pull/1148)

   Tags: ocis

   We introduced a new environment variable `OCIS_URL` that expects a URL including
   protocol, host and optionally port to simplify configuring all the different
   services. These existing environment variables still take precedence, but will
   also fall back to `OCIS_URL`: `STORAGE_LDAP_IDP`, `STORAGE_OIDC_ISSUER`,
   `PROXY_OIDC_ISSUER`, `STORAGE_FRONTEND_PUBLIC_URL`, `KONNECTD_ISS`,
   `WEB_OIDC_AUTHORITY`, and `WEB_UI_CONFIG_SERVER`.

   Some environment variables are now built dynamically if they are not set: -
   `STORAGE_DATAGATEWAY_PUBLIC_URL` defaults to
   `<STORAGE_FRONTEND_PUBLIC_URL>/data`, also falling back to `OCIS_URL` -
   `WEB_OIDC_METADATA_URL` defaults to
   `<WEB_OIDC_AUTHORITY>/.well-known/openid-configuration`, also falling back to
   `OCIS_URL`

   Furthermore, the built in konnectd will generate an
   `identifier-registration.yaml` that uses the `KONNECTD_ISS` in the allowed
   `redirect_uris` and `origins`. It simplifies the default
   `https://localhost:9200` and remote deployment with `OCIS_URL` which is
   evaluated as a fallback if `KONNECTD_ISS` is not set.

   An oCIS server can now be started on a remote machine as easy as
   `OCIS_URL=https://cloud.ocis.test PROXY_HTTP_ADDR=0.0.0.0:443 ocis server`.

   Note that the `OCIS_DOMAIN` environment variable is not used by oCIS, but by the
   docker containers.

   https://github.com/owncloud/ocis/pull/1148

* Enhancement - Update reva to v1.4.1-0.20210111080247-f2b63bfd6825: [#1194](https://github.com/owncloud/ocis/pull/1194)

  * Enhancement: calculate and expose actual file permission set [cs3org/reva#1368](https://github.com/cs3org/reva/pull/1368)
  * initial range request support [cs3org/reva#1326](https://github.com/cs3org/reva/pull/1388)

   https://github.com/owncloud/ocis/pull/1194
   https://github.com/cs3org/reva/pull/1368
   https://github.com/cs3org/reva/pull/1388

* Enhancement - Add named locks and refactor cache: [#1212](https://github.com/owncloud/ocis/pull/1212)

   Tags: ocis-pkg, accounts

   We had the case that we needed kind of a named locking mechanism which enables
   us to lock only under certain conditions. It's used in the indexer package where
   we do not need to lock everything, instead just lock the requested parts and
   differentiate between reads and writes.

   This made it possible to entirely remove locks from the accounts service and
   move them to the ocis-pkg indexer. Another part of this refactor was to make the
   cache atomic and write tests for it.

   - remove locking from accounts service - add sync package with named mutex - add
   named locking to indexer - move cache to sync package

   https://github.com/owncloud/ocis/issues/966
   https://github.com/owncloud/ocis/pull/1212

* Enhancement - Use sync.cache for roles cache: [#1367](https://github.com/owncloud/ocis/pull/1367)

   Tags: ocis-pkg

   Update ocis-pkg/roles cache to use ocis-pkg/sync cache

   https://github.com/owncloud/ocis/pull/1367

* Enhancement - Update reva to v1.5.1: [#1372](https://github.com/owncloud/ocis/pull/1372)

   Summary -------

  * Fix #1401: Use the user in request for deciding the layout for non-home DAV requests
  * Fix #1413: Re-include the '.git' dir in the Docker images to pass the version tag
  * Fix #1399: Fix ocis trash-bin purge
  * Enh #1397: Bump the Copyright date to 2021
  * Enh #1398: Support site authorization status in Mentix
  * Enh #1393: Allow setting favorites, mtime and a temporary etag
  * Enh #1403: Support remote cloud gathering metrics

   Details -------

  * Bugfix #1401: Use the user in request for deciding the layout for non-home DAV requests

   For the incoming /dav/files/userID requests, we have different namespaces
   depending on whether the request is for the logged-in user's namespace or not.
   Since in the storage drivers, we specify the layout depending only on the user
   whose resources are to be accessed, this fails when a user wants to access
   another user's namespace when the storage provider depends on the logged in
   user's namespace. This PR fixes that.

   For example, consider the following case. The owncloud fs uses a layout {{substr
   0 1 .Id.OpaqueId}}/{{.Id.OpaqueId}}. The user einstein sends a request to access
   a resource shared with him, say /dav/files/marie/abcd, which should be allowed.
   However, based on the way we applied the layout, there's no way in which this
   can be translated to /m/marie/.

   Https://github.com/cs3org/reva/pull/1401

  * Bugfix #1413: Re-include the '.git' dir in the Docker images to pass the version tag

   And git SHA to the release tool.

   Https://github.com/cs3org/reva/pull/1413

  * Bugfix #1399: Fix ocis trash-bin purge

   Fixes the empty trash-bin functionality for ocis-storage

   Https://github.com/owncloud/product/issues/254
   https://github.com/cs3org/reva/pull/1399

  * Enhancement #1397: Bump the Copyright date to 2021

   Https://github.com/cs3org/reva/pull/1397

  * Enhancement #1398: Support site authorization status in Mentix

   This enhancement adds support for a site authorization status to Mentix. This
   way, sites registered via a web app can now be excluded until authorized
   manually by an administrator.

   Furthermore, Mentix now sets the scheme for Prometheus targets. This allows us
   to also support monitoring of sites that do not support the default HTTPS
   scheme.

   Https://github.com/cs3org/reva/pull/1398

  * Enhancement #1393: Allow setting favorites, mtime and a temporary etag

   We now let the oCIS driver persist favorites, set temporary etags and the mtime
   as arbitrary metadata.

   Https://github.com/owncloud/ocis/issues/567
   https://github.com/cs3org/reva/issues/1394
   https://github.com/cs3org/reva/pull/1393

  * Enhancement #1403: Support remote cloud gathering metrics

   The current metrics package can only gather metrics either from json files. With
   this feature, the metrics can be gathered polling the http endpoints exposed by
   the owncloud/nextcloud sciencemesh apps.

   Https://github.com/cs3org/reva/pull/1403

   https://github.com/owncloud/ocis/pull/1372

# Changelog for [1.0.0] (2020-12-17)

The following sections list the changes for 1.0.0.

## Summary

* Bugfix - Fix path of files shared with me in ocs api: [#204](https://github.com/owncloud/product/issues/204)
* Bugfix - Add missing env vars to docker compose: [#392](https://github.com/owncloud/ocis/pull/392)
* Bugfix - Build docker images with alpine:latest instead of alpine:edge: [#416](https://github.com/owncloud/ocis/pull/416)
* Bugfix - Don't enforce empty external apps slice: [#473](https://github.com/owncloud/ocis/pull/473)
* Bugfix - Fix director selection in proxy: [#521](https://github.com/owncloud/ocis/pull/521)
* Bugfix - Fix button layout after phoenix update: [#625](https://github.com/owncloud/ocis/pull/625)
* Bugfix - Don't create account if id/mail/username already taken: [#709](https://github.com/owncloud/ocis/pull/709)
* Bugfix - Use micro default client: [#718](https://github.com/owncloud/ocis/pull/718)
* Bugfix - Mint token with uid and gid: [#737](https://github.com/owncloud/ocis/pull/737)
* Bugfix - Lower Bound was not working for the cs3 api index implementation: [#741](https://github.com/owncloud/ocis/pull/741)
* Bugfix - Fix id or username query handling: [#745](https://github.com/owncloud/ocis/pull/745)
* Bugfix - Allow consent-prompt with switch-account: [#788](https://github.com/owncloud/ocis/pull/788)
* Bugfix - Accounts config sometimes being overwritten: [#808](https://github.com/owncloud/ocis/pull/808)
* Bugfix - Fix konnectd build: [#809](https://github.com/owncloud/ocis/pull/809)
* Bugfix - Make settings service start without go coroutines: [#835](https://github.com/owncloud/ocis/pull/835)
* Bugfix - Fix choose account dialogue: [#846](https://github.com/owncloud/ocis/pull/846)
* Bugfix - Enable scrolling in accounts list: [#909](https://github.com/owncloud/ocis/pull/909)
* Bugfix - Serve index.html for directories: [#912](https://github.com/owncloud/ocis/pull/912)
* Bugfix - Disable public link expiration by default: [#987](https://github.com/owncloud/ocis/issues/987)
* Bugfix - Fix minor ui bugs: [#1043](https://github.com/owncloud/ocis/issues/1043)
* Bugfix - Permission checks for settings write access: [#1092](https://github.com/owncloud/ocis/pull/1092)
* Change - Initial release of basic version: [#2](https://github.com/owncloud/ocis/issues/2)
* Change - Start ocis-accounts with the ocis server command: [#25](https://github.com/owncloud/product/issues/25)
* Change - Add cli-commands to manage accounts: [#115](https://github.com/owncloud/product/issues/115)
* Change - Start ocis-proxy with the ocis server command: [#119](https://github.com/owncloud/ocis/issues/119)
* Change - Account management permissions for Admin role: [#124](https://github.com/owncloud/product/issues/124)
* Change - Add the thumbnails command: [#156](https://github.com/owncloud/ocis/issues/156)
* Change - Integrate import command from ocis-migration: [#249](https://github.com/owncloud/ocis/pull/249)
* Change - Switch over to a new custom-built runtime: [#287](https://github.com/owncloud/ocis/pull/287)
* Change - Make ocis-settings available: [#287](https://github.com/owncloud/ocis/pull/287)
* Change - Update reva config: [#336](https://github.com/owncloud/ocis/pull/336)
* Change - Use bcrypt to hash the user passwords: [#510](https://github.com/owncloud/ocis/issues/510)
* Change - Improve reva service descriptions: [#536](https://github.com/owncloud/ocis/pull/536)
* Change - Choose disk or cs3 storage for accounts and groups: [#623](https://github.com/owncloud/ocis/pull/623)
* Change - Update phoenix to v0.18.0: [#651](https://github.com/owncloud/ocis/pull/651)
* Change - Accounts UI shows message when no permissions: [#656](https://github.com/owncloud/ocis/pull/656)
* Change - Settings and accounts appear in the user menu: [#656](https://github.com/owncloud/ocis/pull/656)
* Change - Update phoenix to v0.20.0: [#674](https://github.com/owncloud/ocis/pull/674)
* Change - Unify Configuration Parsing: [#675](https://github.com/owncloud/ocis/pull/675)
* Change - Default apps in ownCloud Web: [#688](https://github.com/owncloud/ocis/pull/688)
* Change - Bring oC theme: [#698](https://github.com/owncloud/ocis/pull/698)
* Change - Filesystem based index: [#709](https://github.com/owncloud/ocis/pull/709)
* Change - Remove username field in OCS: [#709](https://github.com/owncloud/ocis/pull/709)
* Change - Update phoenix to v0.21.0: [#728](https://github.com/owncloud/ocis/pull/728)
* Change - Clarify storage driver env vars: [#729](https://github.com/owncloud/ocis/pull/729)
* Change - Rebuild index command for accounts: [#748](https://github.com/owncloud/ocis/pull/748)
* Change - Properly style konnectd consent page: [#754](https://github.com/owncloud/ocis/pull/754)
* Change - Update phoenix to v0.22.0: [#757](https://github.com/owncloud/ocis/pull/757)
* Change - Update phoenix to v0.23.0: [#785](https://github.com/owncloud/ocis/pull/785)
* Change - Move the indexer package from ocis/accounts to ocis/ocis-pkg: [#794](https://github.com/owncloud/ocis/pull/794)
* Change - Enable OpenID dynamic client registration: [#811](https://github.com/owncloud/ocis/issues/811)
* Change - Update phoenix to v0.24.0: [#817](https://github.com/owncloud/ocis/pull/817)
* Change - Move ocis default config to root level: [#842](https://github.com/owncloud/ocis/pull/842)
* Change - Update phoenix to v0.25.0: [#868](https://github.com/owncloud/ocis/pull/868)
* Change - Theme welcome and choose account pages: [#887](https://github.com/owncloud/ocis/pull/887)
* Change - Replace the library which scales the images: [#910](https://github.com/owncloud/ocis/pull/910)
* Change - Update phoenix to v0.26.0: [#935](https://github.com/owncloud/ocis/pull/935)
* Change - Update phoenix to v0.27.0: [#943](https://github.com/owncloud/ocis/pull/943)
* Change - Cache password validation: [#958](https://github.com/owncloud/ocis/pull/958)
* Change - Proxy allow insecure upstreams: [#1007](https://github.com/owncloud/ocis/pull/1007)
* Change - CS3 can be used as accounts-backend: [#1020](https://github.com/owncloud/ocis/pull/1020)
* Change - Update phoenix to v0.28.0: [#1027](https://github.com/owncloud/ocis/pull/1027)
* Change - Update phoenix to v0.29.0: [#1034](https://github.com/owncloud/ocis/pull/1034)
* Change - Make all paths configurable and default to a common temp dir: [#1080](https://github.com/owncloud/ocis/pull/1080)
* Change - Update reva to v1.4.1-0.20201209113234-e791b5599a89: [#1089](https://github.com/owncloud/ocis/pull/1089)
* Change - Update ownCloud Web to v1.0.0-beta3: [#1105](https://github.com/owncloud/ocis/pull/1105)
* Change - Update ownCloud Web to v1.0.0-beta4: [#1110](https://github.com/owncloud/ocis/pull/1110)
* Enhancement - Simplify tracing config: [#92](https://github.com/owncloud/product/issues/92)
* Enhancement - Document how to run OCIS on top of EOS: [#172](https://github.com/owncloud/ocis/pull/172)
* Enhancement - Add a command to list the versions of running instances: [#226](https://github.com/owncloud/product/issues/226)
* Enhancement - Add the accounts service: [#244](https://github.com/owncloud/product/issues/244)
* Enhancement - Add the glauth service: [#244](https://github.com/owncloud/product/issues/244)
* Enhancement - Add the konnectd service: [#244](https://github.com/owncloud/product/issues/244)
* Enhancement - Add the ocis-phoenix service: [#244](https://github.com/owncloud/product/issues/244)
* Enhancement - Add the ocis-pkg package: [#244](https://github.com/owncloud/product/issues/244)
* Enhancement - Add the ocs service: [#244](https://github.com/owncloud/product/issues/244)
* Enhancement - Add the proxy service: [#244](https://github.com/owncloud/product/issues/244)
* Enhancement - Add the settings service: [#244](https://github.com/owncloud/product/issues/244)
* Enhancement - Add the storage service: [#244](https://github.com/owncloud/product/issues/244)
* Enhancement - Add the store service: [#244](https://github.com/owncloud/product/issues/244)
* Enhancement - Add the thumbnails service: [#244](https://github.com/owncloud/product/issues/244)
* Enhancement - Add the webdav service: [#244](https://github.com/owncloud/product/issues/244)
* Enhancement - Launch a storage to store ocis-metadata: [#602](https://github.com/owncloud/ocis/pull/602)
* Enhancement - Add basic auth option: [#627](https://github.com/owncloud/ocis/pull/627)
* Enhancement - Add glauth fallback backend: [#649](https://github.com/owncloud/ocis/pull/649)
* Enhancement - Update reva to dd3a8c0f38: [#725](https://github.com/owncloud/ocis/pull/725)
* Enhancement - Update konnectd to v0.33.8: [#744](https://github.com/owncloud/ocis/pull/744)
* Enhancement - Update reva to cdb3d6688da5: [#748](https://github.com/owncloud/ocis/pull/748)
* Enhancement - Update glauth to dev 4f029234b2308: [#786](https://github.com/owncloud/ocis/pull/786)
* Enhancement - Update reva to v1.4.1-0.20201123062044-b2c4af4e897d: [#823](https://github.com/owncloud/ocis/pull/823)
* Enhancement - Update glauth to dev fd3ac7e4bbdc93578655d9a08d8e23f105aaa5b2: [#834](https://github.com/owncloud/ocis/pull/834)
* Enhancement - Better adopt Go-Micro: [#840](https://github.com/owncloud/ocis/pull/840)
* Enhancement - Tidy dependencies: [#845](https://github.com/owncloud/ocis/pull/845)
* Enhancement - Create OnlyOffice extension: [#857](https://github.com/owncloud/ocis/pull/857)
* Enhancement - Cache userinfo in proxy: [#877](https://github.com/owncloud/ocis/pull/877)
* Enhancement - Add permission check when assigning and removing roles: [#879](https://github.com/owncloud/ocis/issues/879)
* Enhancement - Show basic-auth warning only once: [#886](https://github.com/owncloud/ocis/pull/886)
* Enhancement - Create a proxy access-log: [#889](https://github.com/owncloud/ocis/pull/889)
* Enhancement - Add a version command to ocis: [#915](https://github.com/owncloud/ocis/pull/915)
* Enhancement - Add k6: [#941](https://github.com/owncloud/ocis/pull/941)
* Enhancement - Update reva to v1.4.1-0.20201127111856-e6a6212c1b7b: [#971](https://github.com/owncloud/ocis/pull/971)
* Enhancement - Update reva to v1.4.1-0.20201130061320-ac85e68e0600: [#980](https://github.com/owncloud/ocis/pull/980)
* Enhancement - Add www-authenticate based on user agent: [#1009](https://github.com/owncloud/ocis/pull/1009)
* Enhancement - Add tracing to the accounts service: [#1016](https://github.com/owncloud/ocis/issues/1016)
* Enhancement - Runtime Cleanup: [#1066](https://github.com/owncloud/ocis/pull/1066)
* Enhancement - Update reva to 063b3db9162b: [#1091](https://github.com/owncloud/ocis/pull/1091)
* Enhancement - Update OCIS Runtime: [#1108](https://github.com/owncloud/ocis/pull/1108)
* Enhancement - Update reva to v1.4.1-0.20201125144025-57da0c27434c: [#1320](https://github.com/cs3org/reva/pull/1320)

## Details

* Bugfix - Fix path of files shared with me in ocs api: [#204](https://github.com/owncloud/product/issues/204)

   The path of files shared with me using the ocs api was pointing to an incorrect
   location.

   https://github.com/owncloud/product/issues/204
   https://github.com/owncloud/ocis/pull/994

* Bugfix - Add missing env vars to docker compose: [#392](https://github.com/owncloud/ocis/pull/392)

   Tags: docker

   Without setting `REVA_FRONTEND_URL` and `REVA_DATAGATEWAY_URL` uploads would
   default to localhost and fail if `OCIS_DOMAIN` was used to run ocis on a remote
   host.

   https://github.com/owncloud/ocis/pull/392

* Bugfix - Build docker images with alpine:latest instead of alpine:edge: [#416](https://github.com/owncloud/ocis/pull/416)

   Tags: docker

   ARM builds were failing when built on alpine:edge, so we switched to
   alpine:latest instead.

   https://github.com/owncloud/ocis/pull/416

* Bugfix - Don't enforce empty external apps slice: [#473](https://github.com/owncloud/ocis/pull/473)

   Tags: web

   The command for ocis-phoenix enforced an empty external apps configuration. This
   was removed, as it was blocking a new set of default external apps in
   ocis-phoenix.

   https://github.com/owncloud/ocis/pull/473

* Bugfix - Fix director selection in proxy: [#521](https://github.com/owncloud/ocis/pull/521)

   Tags: proxy

   We fixed a bug in ocis-proxy where simultaneous requests could be executed on
   the wrong backend.

   https://github.com/owncloud/ocis/pull/521
   https://github.com/owncloud/ocis-proxy/pull/99

* Bugfix - Fix button layout after phoenix update: [#625](https://github.com/owncloud/ocis/pull/625)

   Tags: accounts

   With the phoenix update to v0.17.0 a new ODS version was released which has a
   breaking change for buttons regarding their layout. We adjusted the button
   layout in the accounts UI accordingly.

   https://github.com/owncloud/ocis/pull/625

* Bugfix - Don't create account if id/mail/username already taken: [#709](https://github.com/owncloud/ocis/pull/709)

   Tags: accounts

   We don't allow anymore to create a new account if the provided id/mail/username
   is already taken.

   https://github.com/owncloud/ocis/pull/709

* Bugfix - Use micro default client: [#718](https://github.com/owncloud/ocis/pull/718)

   Tags: glauth

   We found a file descriptor leak in the glauth connections to the accounts
   service. Fixed it by using the micro default client.

   https://github.com/owncloud/ocis/pull/718

* Bugfix - Mint token with uid and gid: [#737](https://github.com/owncloud/ocis/pull/737)

   Tags: accounts

   The eos driver expects the uid and gid from the opaque map of a user. While the
   proxy does mint tokens correctly, the accounts service wasn't.

   https://github.com/owncloud/ocis/pull/737

* Bugfix - Lower Bound was not working for the cs3 api index implementation: [#741](https://github.com/owncloud/ocis/pull/741)

   Tags: accounts

   Lower bound working on the cs3 index implementation

   https://github.com/owncloud/ocis/pull/741

* Bugfix - Fix id or username query handling: [#745](https://github.com/owncloud/ocis/pull/745)

   Tags: accounts

   The code was stopping execution when encountering an error while loading an
   account by id. But for or queries we can continue execution.

   https://github.com/owncloud/ocis/pull/745

* Bugfix - Allow consent-prompt with switch-account: [#788](https://github.com/owncloud/ocis/pull/788)

   Multiple prompt values are allowed and this change fixes the check for
   select_account if it was used together with other prompt values. Where
   select_account previously was ignored, it is now processed as required, fixing
   the use case when a RP wants to trigger select_account first while at the same
   time wants also to request interactive consent.

   https://github.com/owncloud/ocis/pull/788

* Bugfix - Accounts config sometimes being overwritten: [#808](https://github.com/owncloud/ocis/pull/808)

   Tags: accounts

   Sometimes when running the accounts extensions flags were not being taken into
   consideration.

   https://github.com/owncloud/ocis/pull/808

* Bugfix - Fix konnectd build: [#809](https://github.com/owncloud/ocis/pull/809)

   Tags: konnectd

   We fixed the default config for konnectd and updated the Makefile to include the
   `yarn install`and `yarn build` steps if the static assets are missing.

   https://github.com/owncloud/ocis/pull/809

* Bugfix - Make settings service start without go coroutines: [#835](https://github.com/owncloud/ocis/pull/835)

   The go routines cause a race condition that sometimes causes the tests to fail.
   The ListRoles request would not return all permissions.

   https://github.com/owncloud/ocis/pull/835

* Bugfix - Fix choose account dialogue: [#846](https://github.com/owncloud/ocis/pull/846)

   Tags: konnectd

   We've fixed the choose account dialogue in konnectd bug that the user hasn't
   been logged in after selecting account.

   https://github.com/owncloud/ocis/pull/846

* Bugfix - Enable scrolling in accounts list: [#909](https://github.com/owncloud/ocis/pull/909)

   Tags: accounts

   We've fixed the accounts list to enable scrolling.

   https://github.com/owncloud/ocis/pull/909

* Bugfix - Serve index.html for directories: [#912](https://github.com/owncloud/ocis/pull/912)

   The static middleware in ocis-pkg now serves index.html instead of returning 404
   on paths with a trailing `/`.

   https://github.com/owncloud/ocis-pkg/issues/63
   https://github.com/owncloud/ocis/pull/912

* Bugfix - Disable public link expiration by default: [#987](https://github.com/owncloud/ocis/issues/987)

   Tags: storage

   The public link expiration was enabled by default and didn't have a default
   expiration span by default, which resulted in already expired public links
   coming from the public link quick action. We fixed this by disabling the public
   link expiration by default.

   https://github.com/owncloud/ocis/issues/987
   https://github.com/owncloud/ocis/pull/1035

* Bugfix - Fix minor ui bugs: [#1043](https://github.com/owncloud/ocis/issues/1043)

   - the ui haven't updated the language of the items in the settings view menu.
   Now we listen to the selected language and update the ui - deduplicate
   resetMenuItems call

   https://github.com/owncloud/ocis/issues/1043
   https://github.com/owncloud/ocis/pull/1044

* Bugfix - Permission checks for settings write access: [#1092](https://github.com/owncloud/ocis/pull/1092)

   Tags: settings

   There were several endpoints with write access to the settings service that were
   not protected by permission checks. We introduced a generic settings management
   permission to fix this for now. Will be more fine grained later on.

   https://github.com/owncloud/ocis/pull/1092

* Change - Initial release of basic version: [#2](https://github.com/owncloud/ocis/issues/2)

   Just prepared an initial basic version which simply embeds the minimum of
   required services in the context of the ownCloud Infinite Scale project.

   https://github.com/owncloud/ocis/issues/2

* Change - Start ocis-accounts with the ocis server command: [#25](https://github.com/owncloud/product/issues/25)

   Tags: accounts

   Starts ocis-accounts in single binary mode (./ocis server). This service stores
   the user-account information.

   https://github.com/owncloud/product/issues/25
   https://github.com/owncloud/ocis/pull/239/files

* Change - Add cli-commands to manage accounts: [#115](https://github.com/owncloud/product/issues/115)

   Tags: accounts

   COMMANDS:

  * list, ls        List existing accounts
  * add, create     Create a new account
  * update          Make changes to an existing account
  * remove, rm      Removes an existing account
  * inspect         Show detailed data on an existing account
  * help, h         Shows a list of commands or help for one command

   https://github.com/owncloud/product/issues/115

* Change - Start ocis-proxy with the ocis server command: [#119](https://github.com/owncloud/ocis/issues/119)

   Tags: proxy

   Starts the proxy in single binary mode (./ocis server) on port 9200. The proxy
   serves as a single-entry point for all http-clients.

   https://github.com/owncloud/ocis/issues/119
   https://github.com/owncloud/ocis/issues/136

* Change - Account management permissions for Admin role: [#124](https://github.com/owncloud/product/issues/124)

   Tags: accounts, settings

   We created an `AccountManagement` permission and added it to the default admin
   role. There are permission checks in place to protected http endpoints in
   ocis-accounts against requests without the permission. All existing default
   users (einstein, marie, richard) have the default user role now (doesn't have
   the `AccountManagement` permission). Additionally, there is a new default Admin
   user with credentials `moss:vista`.

   Known issue: for users without the `AccountManagement` permission, the accounts
   UI extension is still available in the ocis-web app switcher, but the requests
   for loading the users will fail (as expected). We are working on a way to hide
   the accounts UI extension if the user doesn't have the `AccountManagement`
   permission.

   https://github.com/owncloud/product/issues/124
   https://github.com/owncloud/ocis-settings/pull/59
   https://github.com/owncloud/ocis-settings/pull/66
   https://github.com/owncloud/ocis-settings/pull/67
   https://github.com/owncloud/ocis-settings/pull/69
   https://github.com/owncloud/ocis-proxy/pull/95
   https://github.com/owncloud/ocis-pkg/pull/59
   https://github.com/owncloud/ocis-accounts/pull/95
   https://github.com/owncloud/ocis-accounts/pull/100
   https://github.com/owncloud/ocis-accounts/pull/102

* Change - Add the thumbnails command: [#156](https://github.com/owncloud/ocis/issues/156)

   Tags: thumbnails

   Added the thumbnails command so that the thumbnails service can get started via
   ocis.

   https://github.com/owncloud/ocis/issues/156

* Change - Integrate import command from ocis-migration: [#249](https://github.com/owncloud/ocis/pull/249)

   Tags: migration

   https://github.com/owncloud/ocis/pull/249
   https://github.com/owncloud/ocis-migration

* Change - Switch over to a new custom-built runtime: [#287](https://github.com/owncloud/ocis/pull/287)

   We moved away from using the go-micro runtime and are now using [our own
   runtime](https://github.com/refs/pman). This allows us to spawn service
   processes even when they are using different versions of go-micro. On top of
   that we now have the commands `ocis list`, `ocis kill` and `ocis run` available
   for service runtime management.

   https://github.com/owncloud/ocis/pull/287

* Change - Make ocis-settings available: [#287](https://github.com/owncloud/ocis/pull/287)

   Tags: settings

   This version delivers `settings` as a new service. It is part of the array of
   services in the `server` command.

   https://github.com/owncloud/ocis/pull/287

* Change - Update reva config: [#336](https://github.com/owncloud/ocis/pull/336)

  * EOS homes are not configured with an enable-flag anymore, but with a dedicated storage driver.
  * We're using it now and adapted default configs of storages

   https://github.com/owncloud/ocis/pull/336
   https://github.com/owncloud/ocis/pull/337
   https://github.com/owncloud/ocis/pull/338
   https://github.com/owncloud/ocis-reva/pull/891

* Change - Use bcrypt to hash the user passwords: [#510](https://github.com/owncloud/ocis/issues/510)

   Change the hashing algorithm from SHA-512 to bcrypt since the latter is better
   suitable for password hashing. This is a breaking change. Existing deployments
   need to regenerate the accounts folder.

   https://github.com/owncloud/ocis/issues/510

* Change - Improve reva service descriptions: [#536](https://github.com/owncloud/ocis/pull/536)

   Tags: docs

   The descriptions make it clearer that the services actually represent a mount
   point in the combined storage. Each mount point can have a different driver.

   https://github.com/owncloud/ocis/pull/536

* Change - Choose disk or cs3 storage for accounts and groups: [#623](https://github.com/owncloud/ocis/pull/623)

   Tags: accounts

   The accounts service now has an abstraction layer for the storage. In addition
   to the local disk implementation we implemented a cs3 storage, which is the new
   default for the accounts service.

   https://github.com/owncloud/ocis/pull/623

* Change - Update phoenix to v0.18.0: [#651](https://github.com/owncloud/ocis/pull/651)

   Tags: web

   We updated phoenix to v0.18.0. Please refer to the changelog (linked) for
   details on the phoenix release. With the ODS release brought in by phoenix we
   now have proper oc-checkbox and oc-radio components for the settings and
   accounts UI.

   https://github.com/owncloud/ocis/pull/651
   https://github.com/owncloud/phoenix/releases/tag/v0.18.0
   https://github.com/owncloud/owncloud-design-system/releases/tag/v1.12.1

* Change - Accounts UI shows message when no permissions: [#656](https://github.com/owncloud/ocis/pull/656)

   We improved the UX of the accounts UI by showing a message information the user
   about missing permissions when the accounts or roles fail to load. This was
   showing an indeterminate progress bar before.

   https://github.com/owncloud/ocis/pull/656

* Change - Settings and accounts appear in the user menu: [#656](https://github.com/owncloud/ocis/pull/656)

   We moved settings and accounts to the user menu.

   https://github.com/owncloud/ocis/pull/656

* Change - Update phoenix to v0.20.0: [#674](https://github.com/owncloud/ocis/pull/674)

   Tags: web

   We updated phoenix to v0.20.0. Please refer to the changelog (linked) for
   details on the phoenix release.

   https://github.com/owncloud/ocis/pull/674
   https://github.com/owncloud/phoenix/releases/tag/v0.20.0

* Change - Unify Configuration Parsing: [#675](https://github.com/owncloud/ocis/pull/675)

   Tags: ocis

   - responsibility for config parsing should be on the subcommand - if there is a
   config file in the environment location, env var should take precedence -
   general rule of thumb: the more explicit the config file is that would be picked
   up. Order from less to more explicit: - config location (/etc/ocis) -
   environment variable - cli flag

   https://github.com/owncloud/ocis/pull/675

* Change - Default apps in ownCloud Web: [#688](https://github.com/owncloud/ocis/pull/688)

   Tags: web

   We changed the default apps for ownCloud Web to be only files and media-viewer.
   Markdown-editor and draw-io have been removed as defaults.

   https://github.com/owncloud/ocis/pull/688

* Change - Bring oC theme: [#698](https://github.com/owncloud/ocis/pull/698)

   Tags: konnectd

   We've styled our konnectd login page to reflect ownCloud theme.

   https://github.com/owncloud/ocis/pull/698

* Change - Filesystem based index: [#709](https://github.com/owncloud/ocis/pull/709)

   Tags: accounts, storage

   We replaced `bleve` with a new filesystem based index implementation. There is
   an `indexer` which is capable of orchestrating different index types to build
   indices on documents by field. You can choose from the index types `unique`,
   `non-unique` or `autoincrement`. Indices can be utilized to run search queries
   (full matches or globbing) on document fields. The accounts service is using
   this index internally to run the search queries coming in via `ListAccounts` and
   `ListGroups` and to generate UIDs for new accounts as well as GIDs for new
   groups.

   The accounts service can be configured to store the index on the local FS / a
   NFS (`disk` implementation of the index) or to use an arbitrary storage ( `cs3`
   implementation of the index). `cs3` is the new default, which is configured to
   use the `metadata` storage.

   https://github.com/owncloud/ocis/pull/709

* Change - Remove username field in OCS: [#709](https://github.com/owncloud/ocis/pull/709)

   Tags: ocs

   We use the incoming userid as both the `id` and the
   `on_premises_sam_account_name` for new accounts in the accounts service. The
   userid in OCS requests is in fact the username, not our internal account id. We
   need to enforce the userid as our internal account id though, because the
   account id is part of various `path` formats.

   https://github.com/owncloud/ocis/pull/709
   https://github.com/owncloud/ocis/pull/816

* Change - Update phoenix to v0.21.0: [#728](https://github.com/owncloud/ocis/pull/728)

   Tags: web

   We updated phoenix to v0.21.0. Please refer to the changelog (linked) for
   details on the phoenix release.

   https://github.com/owncloud/ocis/pull/728
   https://github.com/owncloud/phoenix/releases/tag/v0.21.0

* Change - Clarify storage driver env vars: [#729](https://github.com/owncloud/ocis/pull/729)

   After renaming ocsi-reva to storage and combining the storage and data providers
   some env vars were confusingly named `STORAGE_STORAGE_...`. We are changing the
   prefix for driver related env vars to `STORAGE_DRIVER_...`. This makes changing
   the storage driver using eg.: `STORAGE_HOME_DRIVER=eos` and setting driver
   options using `STORAGE_DRIVER_EOS_LAYOUT=...` less confusing.

   https://github.com/owncloud/ocis/pull/729

* Change - Rebuild index command for accounts: [#748](https://github.com/owncloud/ocis/pull/748)

   Tags: accounts

   The index for the accounts service can now be rebuilt by running the cli command
   `./bin/ocis accounts rebuild`. It deletes all configured indices and rebuilds
   them from the documents found on storage. For this we also introduced a
   `LoadAccounts` and `LoadGroups` function on storage for loading all existing
   documents.

   https://github.com/owncloud/ocis/pull/748

* Change - Properly style konnectd consent page: [#754](https://github.com/owncloud/ocis/pull/754)

   Tags: konnectd

   After bringing our theme into konnectd, we've had to adjust the styles of the
   consent page so the text is visible and button reflects our theme.

   https://github.com/owncloud/ocis/pull/754

* Change - Update phoenix to v0.22.0: [#757](https://github.com/owncloud/ocis/pull/757)

   Tags: web

   We updated phoenix to v0.22.0. Please refer to the changelog (linked) for
   details on the phoenix release.

   https://github.com/owncloud/ocis/pull/757
   https://github.com/owncloud/phoenix/releases/tag/v0.22.0

* Change - Update phoenix to v0.23.0: [#785](https://github.com/owncloud/ocis/pull/785)

   Tags: web

   We updated phoenix to v0.23.0. Please refer to the changelog (linked) for
   details on the phoenix release.

   https://github.com/owncloud/ocis/pull/785
   https://github.com/owncloud/phoenix/releases/tag/v0.23.0

* Change - Move the indexer package from ocis/accounts to ocis/ocis-pkg: [#794](https://github.com/owncloud/ocis/pull/794)

   We are making that change for semantic reasons. So consumers of any index don't
   necessarily need to know of the accounts service.

   https://github.com/owncloud/ocis/pull/794

* Change - Enable OpenID dynamic client registration: [#811](https://github.com/owncloud/ocis/issues/811)

   Enable OpenID dynamic client registration

   https://github.com/owncloud/ocis/issues/811
   https://github.com/owncloud/ocis/pull/813

* Change - Update phoenix to v0.24.0: [#817](https://github.com/owncloud/ocis/pull/817)

   Tags: web

   We updated phoenix to v0.24.0. Please refer to the changelog (linked) for
   details on the phoenix release.

   https://github.com/owncloud/ocis/pull/817
   https://github.com/owncloud/phoenix/releases/tag/v0.24.0

* Change - Move ocis default config to root level: [#842](https://github.com/owncloud/ocis/pull/842)

   Tags: ocis

   We moved the tracing config to the `root` flagset so that they are parsed on all
   commands. We also introduced a `JWTSecret` flag in the root flagset, in order to
   apply a common default JWTSecret to all services that have one.

   https://github.com/owncloud/ocis/pull/842
   https://github.com/owncloud/ocis/pull/843

* Change - Update phoenix to v0.25.0: [#868](https://github.com/owncloud/ocis/pull/868)

   Tags: web

   We updated phoenix to v0.25.0. Please refer to the changelog (linked) for
   details on the phoenix release.

   https://github.com/owncloud/ocis/pull/868
   https://github.com/owncloud/phoenix/releases/tag/v0.25.0

* Change - Theme welcome and choose account pages: [#887](https://github.com/owncloud/ocis/pull/887)

   Tags: konnectd

   We've themed the konnectd pages Welcome and Choose account. All text has a white
   color now to be easily readable on the dark background.

   https://github.com/owncloud/ocis/pull/887

* Change - Replace the library which scales the images: [#910](https://github.com/owncloud/ocis/pull/910)

   The library went out of support. Also did some refactoring of the thumbnails
   service code.

   https://github.com/owncloud/ocis/pull/910

* Change - Update phoenix to v0.26.0: [#935](https://github.com/owncloud/ocis/pull/935)

   Tags: web

   We updated phoenix to v0.26.0. Please refer to the changelog (linked) for
   details on the phoenix release.

   https://github.com/owncloud/ocis/pull/935
   https://github.com/owncloud/phoenix/releases/tag/v0.26.0

* Change - Update phoenix to v0.27.0: [#943](https://github.com/owncloud/ocis/pull/943)

   Tags: web

   We updated phoenix to v0.27.0. Please refer to the changelog (linked) for
   details on the phoenix release.

   https://github.com/owncloud/ocis/pull/943
   https://github.com/owncloud/phoenix/releases/tag/v0.27.0

* Change - Cache password validation: [#958](https://github.com/owncloud/ocis/pull/958)

   Tags: accounts

   The password validity check for requests like `login eq '%s' and password eq
   '%s'` is now cached for 10 minutes. This improves the performance for basic auth
   requests.

   https://github.com/owncloud/ocis/pull/958

* Change - Proxy allow insecure upstreams: [#1007](https://github.com/owncloud/ocis/pull/1007)

   Tags: proxy

   We can now configure the proxy if insecure upstream servers are allowed. This
   was added since you need to disable certificate checks fore some situations like
   testing.

   https://github.com/owncloud/ocis/pull/1007

* Change - CS3 can be used as accounts-backend: [#1020](https://github.com/owncloud/ocis/pull/1020)

   Tags: proxy

   PROXY_ACCOUNT_BACKEND_TYPE=cs3 PROXY_ACCOUNT_BACKEND_TYPE=accounts (default)

   By using a backend which implements the CS3 user-api (currently provided by
   reva/storage) it is possible to bypass the ocis-accounts service and for example
   use ldap directly.

   https://github.com/owncloud/ocis/pull/1020

* Change - Update phoenix to v0.28.0: [#1027](https://github.com/owncloud/ocis/pull/1027)

   Tags: web

   We updated phoenix to v0.28.0. Please refer to the changelog (linked) for
   details on the phoenix release.

   https://github.com/owncloud/ocis/pull/1027
   https://github.com/owncloud/phoenix/releases/tag/v0.28.0

* Change - Update phoenix to v0.29.0: [#1034](https://github.com/owncloud/ocis/pull/1034)

   Tags: web

   We updated phoenix to v0.29.0. Please refer to the changelog (linked) for
   details on the phoenix release.

   https://github.com/owncloud/ocis/pull/1034
   https://github.com/owncloud/phoenix/releases/tag/v0.29.0

* Change - Make all paths configurable and default to a common temp dir: [#1080](https://github.com/owncloud/ocis/pull/1080)

   Aligned all services to use a dir following`/var/tmp/ocis/<service>/...` by
   default. Also made some missing temp paths configurable via env vars and config
   flags.

   https://github.com/owncloud/ocis/pull/1080

* Change - Update reva to v1.4.1-0.20201209113234-e791b5599a89: [#1089](https://github.com/owncloud/ocis/pull/1089)

   Updated reva to v1.4.1-0.20201209113234-e791b5599a89

   https://github.com/owncloud/ocis/pull/1089

* Change - Update ownCloud Web to v1.0.0-beta3: [#1105](https://github.com/owncloud/ocis/pull/1105)

   Tags: web

   We updated ownCloud Web to v1.0.0-beta3. Please refer to the changelog (linked)
   for details on the web release.

   https://github.com/owncloud/ocis/pull/1105
   https://github.com/owncloud/phoenix/releases/tag/v1.0.0-beta3

* Change - Update ownCloud Web to v1.0.0-beta4: [#1110](https://github.com/owncloud/ocis/pull/1110)

   Tags: web

   We updated ownCloud Web to v1.0.0-beta4. Please refer to the changelog (linked)
   for details on the web release.

   https://github.com/owncloud/ocis/pull/1110
   https://github.com/owncloud/phoenix/releases/tag/v1.0.0-beta4

* Enhancement - Simplify tracing config: [#92](https://github.com/owncloud/product/issues/92)

   We now apply the oCIS tracing config to all services which have tracing. With
   this it is possible to set one tracing config for all services at the same time.

   https://github.com/owncloud/product/issues/92
   https://github.com/owncloud/ocis/pull/329
   https://github.com/owncloud/ocis/pull/409

* Enhancement - Document how to run OCIS on top of EOS: [#172](https://github.com/owncloud/ocis/pull/172)

   Tags: eos

   We have added rules to the Makefile that use the official [eos docker
   images](https://gitlab.cern.ch/eos/eos-docker) to boot an eos cluster and
   configure OCIS to use it.

   https://github.com/owncloud/ocis/pull/172

* Enhancement - Add a command to list the versions of running instances: [#226](https://github.com/owncloud/product/issues/226)

   Tags: accounts

   Added a micro command to list the versions of running accounts services.

   https://github.com/owncloud/product/issues/226

* Enhancement - Add the accounts service: [#244](https://github.com/owncloud/product/issues/244)

   Tags: accounts

  * Bugfix - Initialize roleService client in GRPC server: [#114](https://github.com/owncloud/ocis-accounts/pull/114)
  * Bugfix - Cleanup separated indices in memory: [#224](https://github.com/owncloud/product/issues/224)
  * Change - Set user role on builtin users: [#102](https://github.com/owncloud/ocis-accounts/pull/102)
  * Change - Add new builtin admin user: [#102](https://github.com/owncloud/ocis-accounts/pull/102)
  * Change - We make use of the roles cache to enforce permission checks: [#100](https://github.com/owncloud/ocis-accounts/pull/100)
  * Change - We make use of the roles manager to enforce permission checks: [#108](https://github.com/owncloud/ocis-accounts/pull/108)
  * Enhancement - Add create account form: [#148](https://github.com/owncloud/product/issues/148)
  * Enhancement - Add delete accounts action: [#148](https://github.com/owncloud/product/issues/148)
  * Enhancement - Add enable/disable capabilities to the WebUI: [#118](https://github.com/owncloud/product/issues/118)
  * Enhancement - Improve visual appearance of accounts UI: [#222](https://github.com/owncloud/product/issues/222)
  * Bugfix - Adapting to new settings API for fetching roles: [#96](https://github.com/owncloud/ocis-accounts/pull/96)
  * Change - Create account api-call implicitly adds "default-user" role: [#173](https://github.com/owncloud/product/issues/173)
  * Change - Add role selection to accounts UI: [#103](https://github.com/owncloud/product/issues/103)
  * Bugfix - Atomic Requests: [#82](https://github.com/owncloud/ocis-accounts/pull/82)
  * Bugfix - Unescape value for prefix query: [#76](https://github.com/owncloud/ocis-accounts/pull/76)
  * Change - Adapt to new ocis-settings data model: [#87](https://github.com/owncloud/ocis-accounts/pull/87)
  * Change - Add permissions for language to default roles: [#88](https://github.com/owncloud/ocis-accounts/pull/88)
  * Bugfix - Add write mutexes: [#71](https://github.com/owncloud/ocis-accounts/pull/71)
  * Bugfix - Fix the accountId and groupId mismatch in DeleteGroup Method: [#60](https://github.com/owncloud/ocis-accounts/pull/60)
  * Bugfix - Fix index mapping: [#73](https://github.com/owncloud/ocis-accounts/issues/73)
  * Bugfix - Use NewNumericRangeInclusiveQuery for numeric literals: [#28](https://github.com/owncloud/ocis-glauth/issues/28)
  * Bugfix - Prevent segfault when no password is set: [#65](https://github.com/owncloud/ocis-accounts/pull/65)
  * Bugfix - Update account return value not used: [#70](https://github.com/owncloud/ocis-accounts/pull/70)
  * Bugfix - Build docker images with alpine:latest instead of alpine:edge: [#64](https://github.com/owncloud/ocis-accounts/pull/64)
  * Change - Align structure of this extension with other extensions: [#51](https://github.com/owncloud/ocis-accounts/pull/51)
  * Change - Change api errors: [#11](https://github.com/owncloud/ocis-accounts/issues/11)
  * Change - Enable accounts on creation: [#43](https://github.com/owncloud/ocis-accounts/issues/43)
  * Change - Fix index update on create/update: [#57](https://github.com/owncloud/ocis-accounts/issues/57)
  * Change - Pass around the correct logger throughout the code: [#41](https://github.com/owncloud/ocis-accounts/issues/41)
  * Change - Remove timezone setting: [#33](https://github.com/owncloud/ocis-accounts/pull/33)
  * Change - Tighten screws on usernames and email addresses: [#65](https://github.com/owncloud/ocis-accounts/pull/65)
  * Enhancement - Add early version of cli tools for user-management: [#69](https://github.com/owncloud/ocis-accounts/pull/69)
  * Enhancement - Update accounts API: [#30](https://github.com/owncloud/ocis-accounts/pull/30)
  * Enhancement - Add simple user listing UI: [#51](https://github.com/owncloud/ocis-accounts/pull/51)
  * Enhancement - Logging is configurable: [#24](https://github.com/owncloud/ocis-accounts/pull/24)
  * Change - Initial release of basic version: [#1](https://github.com/owncloud/ocis-accounts/issues/1)
  * Enhancement - Configuration: [#15](https://github.com/owncloud/ocis-accounts/pull/15)

   https://github.com/owncloud/product/issues/244

* Enhancement - Add the glauth service: [#244](https://github.com/owncloud/product/issues/244)

   Tags: glauth

  * Bugfix - Return invalid credentials when user was not found: [#30](https://github.com/owncloud/ocis-glauth/pull/30)
  * Bugfix - Query numeric attribute values without quotes: [#28](https://github.com/owncloud/ocis-glauth/issues/28)
  * Bugfix - Use searchBaseDN if already a user/group name: [#214](https://github.com/owncloud/product/issues/214)
  * Bugfix - Fix LDAP substring startswith filters: [#31](https://github.com/owncloud/ocis-glauth/pull/31)
  * Enhancement - Add build information to the metrics: [#226](https://github.com/owncloud/product/issues/226)
  * Enhancement - Reenable configuring backends: [#600](https://github.com/owncloud/ocis/pull/600)
  * Bugfix - Ignore case when comparing objectclass values: [#26](https://github.com/owncloud/ocis-glauth/pull/26)
  * Bugfix - Build docker images with alpine:latest instead of alpine:edge: [#24](https://github.com/owncloud/ocis-glauth/pull/24)
  * Enhancement - Handle ownCloudUUID attribute: [#27](https://github.com/owncloud/ocis-glauth/pull/27)
  * Enhancement - Implement group queries: [#22](https://github.com/owncloud/ocis-glauth/issues/22)
  * Enhancement - Configuration: [#11](https://github.com/owncloud/ocis-glauth/pull/11)
  * Enhancement - Improve default settings: [#12](https://github.com/owncloud/ocis-glauth/pull/12)
  * Enhancement - Generate temporary ldap certificates if LDAPS is enabled: [#12](https://github.com/owncloud/ocis-glauth/pull/12)
  * Enhancement - Provide additional tls-endpoint: [#12](https://github.com/owncloud/ocis-glauth/pull/12)
  * Change - Use physicist demo users: [#5](https://github.com/owncloud/ocis-glauth/issues/5)
  * Change - Default to config based user backend: [#6](https://github.com/owncloud/ocis-glauth/pull/6)

   https://github.com/owncloud/product/issues/244

* Enhancement - Add the konnectd service: [#244](https://github.com/owncloud/product/issues/244)

   Tags: konnectd

  * Enhancement - Add version command: [#226](https://github.com/owncloud/product/issues/226)
  * Bugfix - Add silent redirect url: [#69](https://github.com/owncloud/ocis-konnectd/issues/69)
  * Bugfix - Build docker images with alpine:latest instead of alpine:edge: [#71](https://github.com/owncloud/ocis-konnectd/pull/71)
  * Bugfix - Include the assets for #62: [#64](https://github.com/owncloud/ocis-konnectd/pull/64)
  * Bugfix - Redirect to the provided uri: [#26](https://github.com/owncloud/ocis-konnectd/issues/26)
  * Change - Add a trailing slash to trusted redirect uris: [#26](https://github.com/owncloud/ocis-konnectd/issues/26)
  * Change - Improve client identifiers for end users: [#62](https://github.com/owncloud/ocis-konnectd/pull/62)
  * Enhancement - Use upstream version of konnect library: [#14](https://github.com/owncloud/product/issues/14)
  * Enhancement - Change default config for single-binary: [#55](https://github.com/owncloud/ocis-konnectd/pull/55)
  * Bugfix - Generate a random CSP-Nonce in the webapp: [#17](https://github.com/owncloud/ocis-konnectd/issues/17)
  * Change - Dummy index.html is not required anymore by upstream: [#25](https://github.com/owncloud/ocis-konnectd/issues/25)
  * Change - Initial release of basic version: [#1](https://github.com/owncloud/ocis-konnectd/issues/1)
  * Change - Use glauth as ldap backend, default to running behind ocis-proxy: [#52](https://github.com/owncloud/ocis-konnectd/pull/52)

   https://github.com/owncloud/product/issues/244

* Enhancement - Add the ocis-phoenix service: [#244](https://github.com/owncloud/product/issues/244)

   Tags: web

  * Bugfix - Fix external app URLs: [#218](https://github.com/owncloud/product/issues/218)
  * Change - Remove pdf-viewer from default apps: [#85](https://github.com/owncloud/ocis-phoenix/pull/85)
  * Change - Enable Settings and Accounts apps by default: [#80](https://github.com/owncloud/ocis-phoenix/pull/80)
  * Bugfix - Exit when assets or config are not found: [#76](https://github.com/owncloud/ocis-phoenix/pull/76)
  * Bugfix - Build docker images with alpine:latest instead of alpine:edge: [#73](https://github.com/owncloud/ocis-phoenix/pull/73)
  * Change - Hide searchbar by default: [#116](https://github.com/owncloud/product/issues/116)
  * Bugfix - Allow silent refresh of access token: [#69](https://github.com/owncloud/ocis-konnectd/issues/69)
  * Change - Update Phoenix: [#60](https://github.com/owncloud/ocis-phoenix/pull/60)
  * Enhancement - Configuration: [#57](https://github.com/owncloud/ocis-phoenix/pull/57)
  * Bugfix - Config file value not being read: [#45](https://github.com/owncloud/ocis-phoenix/pull/45)
  * Change - Default to running behind ocis-proxy: [#55](https://github.com/owncloud/ocis-phoenix/pull/55)

   https://github.com/owncloud/product/issues/244

* Enhancement - Add the ocis-pkg package: [#244](https://github.com/owncloud/product/issues/244)

   Tags: ocis-pkg

  * Change - Unwrap roleIDs from access-token into metadata context: [#59](https://github.com/owncloud/ocis-pkg/pull/59)
  * Change - Provide cache for roles: [#59](https://github.com/owncloud/ocis-pkg/pull/59)
  * Change - Roles manager: [#60](https://github.com/owncloud/ocis-pkg/pull/60)
  * Change - Use go-micro's metadata context for account id: [#56](https://github.com/owncloud/ocis-pkg/pull/56)
  * Bugfix - Remove redigo 2.0.0+incompatible dependency: [#33](https://github.com/owncloud/ocis-graph/pull/33)
  * Change - Add middleware for x-access-token dismantling: [#46](https://github.com/owncloud/ocis-pkg/pull/46)
  * Enhancement - Add `ocis.id` and numeric id claims: [#50](https://github.com/owncloud/ocis-pkg/pull/50)
  * Bugfix - Pass flags to micro service: [#44](https://github.com/owncloud/ocis-pkg/pull/44)
  * Change - Add header to cors handler: [#41](https://github.com/owncloud/ocis-pkg/issues/41)
  * Enhancement - Tracing middleware: [#35](https://github.com/owncloud/ocis-pkg/pull/35/)
  * Enhancement - Allow http services to register handlers: [#33](https://github.com/owncloud/ocis-pkg/pull/33)
  * Change - Upgrade the micro libraries: [#22](https://github.com/owncloud/ocis-pkg/pull/22)
  * Bugfix - Fix Module Path: [#25](https://github.com/owncloud/ocis-pkg/pull/25)
  * Bugfix - Change import paths to ocis-pkg/v2: [#27](https://github.com/owncloud/ocis-pkg/pull/27)
  * Bugfix - Fix serving static assets: [#14](https://github.com/owncloud/ocis-pkg/pull/14)
  * Change - Add TLS support for http services: [#19](https://github.com/owncloud/ocis-pkg/issues/19)
  * Enhancement - Introduce OpenID Connect middleware: [#8](https://github.com/owncloud/ocis-pkg/issues/8)
  * Change - Add root path to static middleware: [#9](https://github.com/owncloud/ocis-pkg/issues/9)
  * Change - Better log level handling within micro: [#2](https://github.com/owncloud/ocis-pkg/issues/2)

   https://github.com/owncloud/product/issues/244

* Enhancement - Add the ocs service: [#244](https://github.com/owncloud/product/issues/244)

   Tags: ocs

  * Bugfix - Match the user response to the OC10 format: [#181](https://github.com/owncloud/product/issues/181)
  * Enhancement - Add version command: [#226](https://github.com/owncloud/product/issues/226)
  * Bugfix - Add the top level response structure to json responses: [#181](https://github.com/owncloud/product/issues/181)
  * Enhancement - Update ocis-accounts: [#42](https://github.com/owncloud/ocis-ocs/pull/42)
  * Bugfix - Mimic oc10 user enabled as string in provisioning api: [#39](https://github.com/owncloud/ocis-ocs/pull/39)
  * Bugfix - Use opaque ID of a user for signing keys: [#436](https://github.com/owncloud/ocis/issues/436)
  * Enhancement - Add option to create user with uidnumber and gidnumber: [#34](https://github.com/owncloud/ocis-ocs/pull/34)
  * Bugfix - Fix file descriptor leak: [#79](https://github.com/owncloud/ocis-accounts/issues/79)
  * Enhancement - Add Group management for OCS Provisioning API: [#25](https://github.com/owncloud/ocis-ocs/pull/25)
  * Enhancement - Basic Support for the User Provisioning API: [#23](https://github.com/owncloud/ocis-ocs/pull/23)
  * Bugfix - Build docker images with alpine:latest instead of alpine:edge: [#20](https://github.com/owncloud/ocis-ocs/pull/20)
  * Change - Initial release of basic version: [#1](https://github.com/owncloud/ocis-ocs/issues/1)
  * Change - Upgrade micro libraries: [#11](https://github.com/owncloud/ocis-ocs/issues/11)
  * Enhancement - Configuration: [#14](https://github.com/owncloud/ocis-ocs/pull/14)
  * Enhancement - Support signing key: [#18](https://github.com/owncloud/ocis-ocs/pull/18)

   https://github.com/owncloud/product/issues/244

* Enhancement - Add the proxy service: [#244](https://github.com/owncloud/product/issues/244)

   Tags: proxy

  * Bugfix - Fix director selection: [#99](https://github.com/owncloud/ocis-proxy/pull/99)
  * Bugfix - Add settings API and app endpoints to example config: [#93](https://github.com/owncloud/ocis-proxy/pull/93)
  * Change - Remove accounts caching: [#100](https://github.com/owncloud/ocis-proxy/pull/100)
  * Enhancement - Add autoprovision accounts flag: [#219](https://github.com/owncloud/product/issues/219)
  * Enhancement - Add hello API and app endpoints to example config and builtin config: [#96](https://github.com/owncloud/ocis-proxy/pull/96)
  * Enhancement - Add roleIDs to the access token: [#95](https://github.com/owncloud/ocis-proxy/pull/95)
  * Enhancement - Add version command: [#226](https://github.com/owncloud/product/issues/226)
  * Enhancement - Add numeric uid and gid to the access token: [#89](https://github.com/owncloud/ocis-proxy/pull/89)
  * Enhancement - Add configuration options for the pre-signed url middleware: [#91](https://github.com/owncloud/ocis-proxy/issues/91)
  * Bugfix - Enable new accounts by default: [#79](https://github.com/owncloud/ocis-proxy/pull/79)
  * Bugfix - Lookup user by id for presigned URLs: [#85](https://github.com/owncloud/ocis-proxy/pull/85)
  * Bugfix - Build docker images with alpine:latest instead of alpine:edge: [#78](https://github.com/owncloud/ocis-proxy/pull/78)
  * Change - Add settings and ocs group routes: [#81](https://github.com/owncloud/ocis-proxy/pull/81)
  * Change - Add route for user provisioning API in ocis-ocs: [#80](https://github.com/owncloud/ocis-proxy/pull/80)
  * Bugfix - Provide token configuration from config: [#69](https://github.com/owncloud/ocis-proxy/pull/69)
  * Bugfix - Provide token configuration from config: [#76](https://github.com/owncloud/ocis-proxy/pull/76)
  * Change - Add OIDC config flags: [#66](https://github.com/owncloud/ocis-proxy/pull/66)
  * Change - Mint new username property in the reva token: [#62](https://github.com/owncloud/ocis-proxy/pull/62)
  * Enhancement - Add Accounts UI routes: [#65](https://github.com/owncloud/ocis-proxy/pull/65)
  * Enhancement - Add option to disable TLS: [#71](https://github.com/owncloud/ocis-proxy/issues/71)
  * Enhancement - Only send create home request if an account has been migrated: [#52](https://github.com/owncloud/ocis-proxy/issues/52)
  * Enhancement - Create a root span on proxy that propagates down to consumers: [#64](https://github.com/owncloud/ocis-proxy/pull/64)
  * Enhancement - Support signed URLs: [#73](https://github.com/owncloud/ocis-proxy/issues/73)
  * Bugfix - Accounts service response was ignored: [#43](https://github.com/owncloud/ocis-proxy/pull/43)
  * Bugfix - Fix x-access-token in header: [#41](https://github.com/owncloud/ocis-proxy/pull/41)
  * Change - Point /data endpoint to reva frontend: [#45](https://github.com/owncloud/ocis-proxy/pull/45)
  * Change - Send autocreate home request to reva gateway: [#51](https://github.com/owncloud/ocis-proxy/pull/51)
  * Change - Update to new accounts API: [#39](https://github.com/owncloud/ocis-proxy/issues/39)
  * Enhancement - Retrieve Account UUID From User Claims: [#36](https://github.com/owncloud/ocis-proxy/pull/36)
  * Enhancement - Create account if it doesn't exist in ocis-accounts: [#55](https://github.com/owncloud/ocis-proxy/issues/55)
  * Enhancement - Disable keep-alive on server-side OIDC requests: [#268](https://github.com/owncloud/ocis/issues/268)
  * Enhancement - Make jwt secret configurable: [#41](https://github.com/owncloud/ocis-proxy/pull/41)
  * Enhancement - Respect account_enabled flag: [#53](https://github.com/owncloud/ocis-proxy/issues/53)
  * Change - Update ocis-pkg: [#30](https://github.com/owncloud/ocis-proxy/pull/30)
  * Change - Insecure http-requests are now redirected to https: [#29](https://github.com/owncloud/ocis-proxy/pull/29)
  * Enhancement - Configurable OpenID Connect client: [#27](https://github.com/owncloud/ocis-proxy/pull/27)
  * Enhancement - Add policy selectors: [#4](https://github.com/owncloud/ocis-proxy/issues/4)
  * Bugfix - Set TLS-Certificate correctly: [#25](https://github.com/owncloud/ocis-proxy/pull/25)
  * Change - Route requests based on regex or query parameters: [#21](https://github.com/owncloud/ocis-proxy/issues/21)
  * Enhancement - Proxy client urls in default configuration: [#19](https://github.com/owncloud/ocis-proxy/issues/19)
  * Enhancement - Make TLS-Cert configurable: [#14](https://github.com/owncloud/ocis-proxy/pull/14)
  * Enhancement - Load Proxy Policies at Runtime: [#17](https://github.com/owncloud/ocis-proxy/issues/17)

   https://github.com/owncloud/product/issues/244

* Enhancement - Add the settings service: [#244](https://github.com/owncloud/product/issues/244)

   Tags: settings

  * Bugfix - Fix loading and saving system scoped values: [#66](https://github.com/owncloud/ocis-settings/pull/66)
  * Bugfix - Complete input validation: [#66](https://github.com/owncloud/ocis-settings/pull/66)
  * Change - Add filter option for bundle ids in ListBundles and ListRoles: [#59](https://github.com/owncloud/ocis-settings/pull/59)
  * Change - Reuse roleIDs from the metadata context: [#69](https://github.com/owncloud/ocis-settings/pull/69)
  * Change - Update ocis-pkg/v2: [#72](https://github.com/owncloud/ocis-settings/pull/72)
  * Enhancement - Add version command: [#226](https://github.com/owncloud/product/issues/226)
  * Bugfix - Fix fetching bundles in settings UI: [#61](https://github.com/owncloud/ocis-settings/pull/61)
  * Change - Filter settings by permissions: [#99](https://github.com/owncloud/product/issues/99)
  * Change - Add role service: [#110](https://github.com/owncloud/product/issues/110)
  * Change - Rename endpoints and message types: [#36](https://github.com/owncloud/ocis-settings/issues/36)
  * Change - Use UUIDs instead of alphanumeric identifiers: [#46](https://github.com/owncloud/ocis-settings/pull/46)
  * Bugfix - Adjust UUID validation to be more tolerant: [#41](https://github.com/owncloud/ocis-settings/issues/41)
  * Bugfix - Fix runtime error when type asserting on nil value: [#38](https://github.com/owncloud/ocis-settings/pull/38)
  * Bugfix - Fix multiple submits on string and number form elements: [#745](https://github.com/owncloud/owncloud-design-system/issues/745)
  * Bugfix - Build docker images with alpine:latest instead of alpine:edge: [#39](https://github.com/owncloud/ocis-settings/pull/39)
  * Change - Dynamically add navItems for extensions with settings bundles: [#25](https://github.com/owncloud/ocis-settings/pull/25)
  * Change - Introduce input validation: [#22](https://github.com/owncloud/ocis-settings/pull/22)
  * Change - Use account uuid from x-access-token: [#14](https://github.com/owncloud/ocis-settings/pull/14)
  * Change - Use server config variable from ocis-web: [#34](https://github.com/owncloud/ocis-settings/pull/34)
  * Enhancement - Remove paths from Makefile: [#33](https://github.com/owncloud/ocis-settings/pull/33)
  * Enhancement - Extend the docs: [#11](https://github.com/owncloud/ocis-settings/issues/11)
  * Enhancement - Update ocis-pkg/v2: [#42](https://github.com/owncloud/ocis-settings/pull/42)

   https://github.com/owncloud/product/issues/244

* Enhancement - Add the storage service: [#244](https://github.com/owncloud/product/issues/244)

   Tags: storage, reva

  * Enhancement - Enable ocis driver treetime accounting: [#620](https://github.com/owncloud/ocis/pull/620)
  * Enhancement - Launch a storage to store ocis-metadata: [#602](https://github.com/owncloud/ocis/pull/602)

   In the future accounts, settings etc. should be stored in a dedicated metadata
   storage. The services should talk to this storage directly, bypassing
   reva-gateway.

   Https://github.com/owncloud/ocis/pull/602

  * Enhancement - Update reva to v1.2.2-0.20200924071957-e6676516e61e: [#601](https://github.com/owncloud/ocis/pull/601)

   - Update reva to v1.2.2-0.20200924071957-e6676516e61e - eos client: Handle eos
   EPERM as permission denied
   [(reva/#1183)](https://github.com/cs3org/reva/pull/1183) - ocis driver: synctime
   based etag propagation [(reva/#1180)](https://github.com/cs3org/reva/pull/1180)
   - ocis driver: fix litmus
   [(reva/#1179)](https://github.com/cs3org/reva/pull/1179) - ocis driver: fix move
   [(reva/#1177)](https://github.com/cs3org/reva/pull/1177) - ocs service: cache
   displaynames [(reva/#1161)](https://github.com/cs3org/reva/pull/1161)

   Https://github.com/owncloud/ocis-reva/issues/262
   https://github.com/owncloud/ocis-reva/issues/357
   https://github.com/owncloud/ocis-reva/issues/301
   https://github.com/owncloud/ocis-reva/issues/302
   https://github.com/owncloud/ocis/pull/601

  * Bugfix - Fix default configuration for accessing shares: [#205](https://github.com/owncloud/product/issues/205)

   The storage provider mounted at `/home` should always have EnableHome set to
   `true`. The other storage providers should have it set to `false`.

   Https://github.com/owncloud/product/issues/205
   https://github.com/owncloud/ocis-reva/pull/461

  * Enhancement - Allow configuring arbitrary storage registry rules: [#193](https://github.com/owncloud/product/issues/193)

   We added a new config flag `storage-registry-rule` that can be given multiple
   times for the gateway to specify arbitrary storage registry rules. You can also
   use a comma separated list of rules in the `REVA_STORAGE_REGISTRY_RULES`
   environment variable.

   Https://github.com/owncloud/product/issues/193
   https://github.com/owncloud/ocis-reva/pull/461

  * Enhancement - Update reva to v1.2.1-0.20200826162318-c0f54e1f37ea: [#454](https://github.com/owncloud/ocis-reva/pull/454)

   - Update reva to v1.2.1-0.20200826162318-c0f54e1f37ea - Do not swallow 'not
   found' errors in Stat [(reva/#1124)](https://github.com/cs3org/reva/pull/1124) -
   Rewire dav files to the home storage
   [(reva/#1125)](https://github.com/cs3org/reva/pull/1125) - Do not restore
   recycle entry on purge [(reva/#1099)](https://github.com/cs3org/reva/pull/1099)
   - Allow listing the trashbin
   [(reva/#1091)](https://github.com/cs3org/reva/pull/1091) - Restore and delete
   trash items via ocs [(reva/#1103)](https://github.com/cs3org/reva/pull/1103) -
   Ensure ignoring public stray shares
   [(reva/#1090)](https://github.com/cs3org/reva/pull/1090) - Ensure ignoring stray
   shares [(reva/#1064)](https://github.com/cs3org/reva/pull/1064) - Minor fixes in
   reva cmd, gateway uploads and smtpclient
   [(reva/#1082)](https://github.com/cs3org/reva/pull/1082) - Owncloud driver -
   propagate mtime on RemoveGrant
   [(reva/#1115)](https://github.com/cs3org/reva/pull/1115) - Handle redirection
   prefixes when extracting destination from URL
   [(reva/#1111)](https://github.com/cs3org/reva/pull/1111) - Add UID and GID in
   ldap auth driver [(reva/#1101)](https://github.com/cs3org/reva/pull/1101) - Add
   calens check to verify changelog entries in CI
   [(reva/#1077)](https://github.com/cs3org/reva/pull/1077) - Refactor Reva CLI
   with prompts [(reva/#1072)](https://github.com/cs3org/reva/pull/1072j) - Get
   file info using fxids from EOS
   [(reva/#1079)](https://github.com/cs3org/reva/pull/1079) - Update LDAP user
   driver [(reva/#1088)](https://github.com/cs3org/reva/pull/1088) - System
   information metrics cleanup
   [(reva/#1114)](https://github.com/cs3org/reva/pull/1114) - System information
   included in Prometheus metrics
   [(reva/#1071)](https://github.com/cs3org/reva/pull/1071) - Add logic for
   resolving storage references over webdav
   [(reva/#1094)](https://github.com/cs3org/reva/pull/1094)

   Https://github.com/owncloud/ocis-reva/pull/454

  * Enhancement - Update reva to v1.2.1-0.20200911111727-51649e37df2d: [#466](https://github.com/owncloud/ocis-reva/pull/466)

   - Update reva to v1.2.1-0.20200911111727-51649e37df2d - Added new OCIS storage
   driver ocis [(reva/#1155)](https://github.com/cs3org/reva/pull/1155) - App
   provider: fallback to env. variable if 'iopsecret' unset
   [(reva/#1146)](https://github.com/cs3org/reva/pull/1146) - Add switch to
   database [(reva/#1135)](https://github.com/cs3org/reva/pull/1135) - Add the
   ocdav HTTP svc to the standalone config
   [(reva/#1128)](https://github.com/cs3org/reva/pull/1128)

   Https://github.com/owncloud/ocis-reva/pull/466

  * Enhancement - Separate user and auth providers, add config for rest user: [#412](https://github.com/owncloud/ocis-reva/pull/412)

   Previously, the auth and user provider services used to have the same driver,
   which restricted using separate drivers and configs for both. This PR separates
   the two and adds the config for the rest user driver and the gatewaysvc
   parameter to EOS fs.

   Https://github.com/owncloud/ocis-reva/pull/412
   https://github.com/cs3org/reva/pull/995

  * Enhancement - Update reva to v1.1.1-0.20200819100654-dcbf0c8ea187: [#447](https://github.com/owncloud/ocis-reva/pull/447)

   - Update reva to v1.1.1-0.20200819100654-dcbf0c8ea187 - fix restoring and
   deleting trash items via ocs
   [(reva/#1103)](https://github.com/cs3org/reva/pull/1103) - Add UID and GID in
   ldap auth driver [(reva/#1101)](https://github.com/cs3org/reva/pull/1101) -
   Allow listing the trashbin
   [(reva/#1091)](https://github.com/cs3org/reva/pull/1091) - Ignore Stray Public
   Shares [(reva/#1090)](https://github.com/cs3org/reva/pull/1090) - Implement
   GetUserByClaim for LDAP user driver
   [(reva/#1088)](https://github.com/cs3org/reva/pull/1088) - eosclient: get file
   info by fxid [(reva/#1079)](https://github.com/cs3org/reva/pull/1079) - Ensure
   stray shares get ignored
   [(reva/#1064)](https://github.com/cs3org/reva/pull/1064) - Improve timestamp
   precision while logging [(reva/#1059)](https://github.com/cs3org/reva/pull/1059)
   - Ocfs lookup userid (update)
   [(reva/#1052)](https://github.com/cs3org/reva/pull/1052) - Disallow sharing the
   shares directory [(reva/#1051)](https://github.com/cs3org/reva/pull/1051) -
   Local storage provider: Fixed resolution of fileid
   [(reva/#1046)](https://github.com/cs3org/reva/pull/1046) - List public shares
   only created by the current user
   [(reva/#1042)](https://github.com/cs3org/reva/pull/1042)

   Https://github.com/owncloud/ocis-reva/pull/447

  * Bugfix - Update LDAP filters: [#399](https://github.com/owncloud/ocis-reva/pull/399)

   With the separation of use and find filters we can now use a filter that taken
   into account a users uuid as well as his username. This is necessary to make
   sharing work with the new account service which assigns accounts an immutable
   account id that is different from the username. Furthermore, the separate find
   filters now allows searching users by their displayname or email as well.

   ```
     userfilter =
      "(&(objectclass=posixAccount)(|(ownclouduuid={{.OpaqueId}})(cn={{.OpaqueId}})))"
      findfilter =
      "(&(objectclass=posixAccount)(|(cn={{query}}*)(displayname={{query}}*)(mail={{query}}*)))"
      ```

   Https://github.com/owncloud/ocis-reva/pull/399
   https://github.com/cs3org/reva/pull/996

  * Change - Environment updates for the username userid split: [#420](https://github.com/owncloud/ocis-reva/pull/420)

   We updated the owncloud storage driver in reva to properly look up users by
   userid or username using the userprovider instead of taking the path segment as
   is. This requires the user service address as well as changing the default
   layout to the userid instead of the username. The latter is not considered a
   stable and persistent identifier.

   Https://github.com/owncloud/ocis-reva/pull/420
   https://github.com/cs3org/reva/pull/1033

  * Enhancement - Update storage documentation: [#384](https://github.com/owncloud/ocis-reva/pull/384)

   We added details to the documentation about storage requirements known from
   ownCloud 10, the local storage driver and the ownCloud storage driver.

   Https://github.com/owncloud/ocis-reva/pull/384
   https://github.com/owncloud/ocis-reva/pull/390

  * Enhancement - Update reva to v0.1.1-0.20200724135750-b46288b375d6: [#399](https://github.com/owncloud/ocis-reva/pull/399)

   - Update reva to v0.1.1-0.20200724135750-b46288b375d6 - Split LDAP user filters
   (reva/#996) - meshdirectory: Add invite forward API to provider links
   (reva/#1000) - OCM: Pass the link to the meshdirectory service in token mail
   (reva/#1002) - Update github.com/go-ldap/ldap to v3 (reva/#1004)

   Https://github.com/owncloud/ocis-reva/pull/399
   https://github.com/cs3org/reva/pull/996 https://github.com/cs3org/reva/pull/1000
   https://github.com/cs3org/reva/pull/1002
   https://github.com/cs3org/reva/pull/1004

  * Enhancement - Update reva to v0.1.1-0.20200728071211-c948977dd3a0: [#407](https://github.com/owncloud/ocis-reva/pull/407)

   - Update reva to v0.1.1-0.20200728071211-c948977dd3a0 - Use proper logging for
   ldap auth requests (reva/#1008) - Update github.com/eventials/go-tus to
   v0.0.0-20200718001131-45c7ec8f5d59 (reva/#1007) - Check if SMTP credentials are
   nil (reva/#1006)

   Https://github.com/owncloud/ocis-reva/pull/407
   https://github.com/cs3org/reva/pull/1008
   https://github.com/cs3org/reva/pull/1007
   https://github.com/cs3org/reva/pull/1006

  * Bugfix - Build docker images with alpine:latest instead of alpine:edge: [#393](https://github.com/owncloud/ocis-reva/pull/393)

   ARM builds were failing when built on alpine:edge, so we switched to
   alpine:latest instead.

   Https://github.com/owncloud/ocis-reva/pull/393

  * Enhancement - Update reva to v0.1.1-0.20200710143425-cf38a45220c5: [#371](https://github.com/owncloud/ocis-reva/pull/371)

   - Update reva to v0.1.1-0.20200710143425-cf38a45220c5 (#371) - Add wopi open
   (reva/#920) - Added a CS3API compliant data exporter to Mentix (reva/#955) -
   Read SMTP password from env if not set in config (reva/#953) - OCS share fix
   including file info after update (reva/#958) - Add flag to smtpclient for for
   unauthenticated SMTP (reva/#963)

   Https://github.com/owncloud/ocis-reva/pull/371
   https://github.com/cs3org/reva/pull/920 https://github.com/cs3org/reva/pull/953
   https://github.com/cs3org/reva/pull/955 https://github.com/cs3org/reva/pull/958
   https://github.com/cs3org/reva/pull/963

  * Enhancement - Update reva to v0.1.1-0.20200722125752-6dea7936f9d1: [#392](https://github.com/owncloud/ocis-reva/pull/392)

   - Update reva to v0.1.1-0.20200722125752-6dea7936f9d1 - Added signing key
   capability (reva/#986) - Add functionality to create webdav references for OCM
   shares (reva/#974) - Added a site locations exporter to Mentix (reva/#972) - Add
   option to config to allow requests to hosts with unverified certificates
   (reva/#969)

   Https://github.com/owncloud/ocis-reva/pull/392
   https://github.com/cs3org/reva/pull/986 https://github.com/cs3org/reva/pull/974
   https://github.com/cs3org/reva/pull/972 https://github.com/cs3org/reva/pull/969

  * Enhancement - Make frontend prefixes configurable: [#363](https://github.com/owncloud/ocis-reva/pull/363)

   We introduce three new environment variables and preconfigure them the following
   way:

  * `REVA_FRONTEND_DATAGATEWAY_PREFIX="data"`
  * `REVA_FRONTEND_OCDAV_PREFIX=""`
  * `REVA_FRONTEND_OCS_PREFIX="ocs"`

   This restores the reva defaults that were changed upstream.

   Https://github.com/owncloud/ocis-reva/pull/363
   https://github.com/cs3org/reva/pull/936/files#diff-51bf4fb310f7362f5c4306581132fc3bR63

  * Enhancement - Update reva to v0.1.1-0.20200701152626-2f6cc60e2f66: [#341](https://github.com/owncloud/ocis-reva/pull/341)

   - Update reva to v0.1.1-0.20200701152626-2f6cc60e2f66 (#341) - Added country
   information to Mentix (reva/#924) - Refactor metrics package to implement reader
   interface (reva/#934) - Fix OCS public link share update values logic (#252,
   #288, reva/#930)

   Https://github.com/owncloud/ocis-reva/issues/252
   https://github.com/owncloud/ocis-reva/issues/288
   https://github.com/owncloud/ocis-reva/pull/341
   https://github.com/cs3org/reva/pull/924 https://github.com/cs3org/reva/pull/934
   https://github.com/cs3org/reva/pull/930

  * Enhancement - Update reva to v0.1.1-0.20200709064551-91eed007038f: [#362](https://github.com/owncloud/ocis-reva/pull/362)

   - Update reva to v0.1.1-0.20200709064551-91eed007038f (#362) - Fix config for
   uploads when data server is not exposed (reva/#936) - Update OCM partners
   endpoints (reva/#937) - Update Ailleron endpoint (reva/#938) - OCS: Fix
   initialization of shares json file (reva/#940) - OCS: Fix returned public link
   URL (#336, reva/#945) - OCS: Share wrap resource id correctly (#344, reva/#951)
   - OCS: Implement share handling for accepting and listing shares (#11,
   reva/#929) - ocm: dynamically lookup IPs for provider check (reva/#946) - ocm:
   add functionality to mail OCM invite tokens (reva/#944) - Change percentagused
   to percentageused (reva/#903) - Fix file-descriptor leak (reva/#954)

   Https://github.com/owncloud/ocis-reva/issues/344
   https://github.com/owncloud/ocis-reva/issues/336
   https://github.com/owncloud/ocis-reva/issues/11
   https://github.com/owncloud/ocis-reva/pull/362
   https://github.com/cs3org/reva/pull/936 https://github.com/cs3org/reva/pull/937
   https://github.com/cs3org/reva/pull/938 https://github.com/cs3org/reva/pull/940
   https://github.com/cs3org/reva/pull/951 https://github.com/cs3org/reva/pull/945
   https://github.com/cs3org/reva/pull/929 https://github.com/cs3org/reva/pull/946
   https://github.com/cs3org/reva/pull/944 https://github.com/cs3org/reva/pull/903
   https://github.com/cs3org/reva/pull/954

  * Enhancement - Add new config options for the http client: [#330](https://github.com/owncloud/ocis-reva/pull/330)

   The internal certificates are checked for validity after
   https://github.com/cs3org/reva/pull/914, which causes the acceptance tests to
   fail. This change sets new hardcoded defaults.

   Https://github.com/owncloud/ocis-reva/pull/330

  * Enhancement - Allow datagateway transfers to take 24h: [#323](https://github.com/owncloud/ocis-reva/pull/323)

   - Increase transfer token life time to 24h (PR #323)

   Https://github.com/owncloud/ocis-reva/pull/323

  * Enhancement - Update reva to v0.1.1-0.20200630075923-39a90d431566: [#320](https://github.com/owncloud/ocis-reva/pull/320)

   - Update reva to v0.1.1-0.20200630075923-39a90d431566 (#320) - Return special
   value for public link password (#294, reva/#904) - Fix public stat and
   listcontainer response to contain the correct prefix (#310, reva/#902)

   Https://github.com/owncloud/ocis-reva/issues/310
   https://github.com/owncloud/ocis-reva/issues/294
   https://github.com/owncloud/ocis-reva/pull/320
   https://github.com/cs3org/reva/pull/902 https://github.com/cs3org/reva/pull/904

  * Enhancement - Update reva to v0.1.1-0.20200701152626-2f6cc60e2f66: [#328](https://github.com/owncloud/ocis-reva/pull/328)

   - Update reva to v0.1.1-0.20200701152626-2f6cc60e2f66 (#328) - Use sync.Map on
   pool package (reva/#909) - Use mutex instead of sync.Map (reva/#915) - Use
   gatewayProviders instead of storageProviders on conn pool (reva/#916) - Add
   logic to ls and stat to process arbitrary metadata keys (reva/#905) -
   Preliminary implementation of Set/UnsetArbitraryMetadata (reva/#912) - Make
   datagateway forward headers (reva/#913, reva/#926) - Add option to cmd upload to
   disable tus (reva/#911) - OCS Share Allow date-only expiration for public shares
   (#288, reva/#918) - OCS Share Remove array from OCS Share update response (#252,
   reva/#919) - OCS Share Implement GET request for single shares (#249, reva/#921)

   Https://github.com/owncloud/ocis-reva/issues/288
   https://github.com/owncloud/ocis-reva/issues/252
   https://github.com/owncloud/ocis-reva/issues/249
   https://github.com/owncloud/ocis-reva/pull/328
   https://github.com/cs3org/reva/pull/909 https://github.com/cs3org/reva/pull/915
   https://github.com/cs3org/reva/pull/916 https://github.com/cs3org/reva/pull/905
   https://github.com/cs3org/reva/pull/912 https://github.com/cs3org/reva/pull/913
   https://github.com/cs3org/reva/pull/926 https://github.com/cs3org/reva/pull/911
   https://github.com/cs3org/reva/pull/918 https://github.com/cs3org/reva/pull/919
   https://github.com/cs3org/reva/pull/921

  * Enhancement - Update reva to v0.1.1-0.20200629131207-04298ea1c088: [#309](https://github.com/owncloud/ocis-reva/pull/309)

   - Update reva to v0.1.1-0.20200629094927-e33d65230abc (#309) - Fix public link
   file share (#278, reva/#895, reva/#900) - Delete public share (reva/#899) -
   Updated reva to v0.1.1-0.20200629131207-04298ea1c088 (#313)

   Https://github.com/owncloud/ocis-reva/issues/278
   https://github.com/owncloud/ocis-reva/pull/309
   https://github.com/cs3org/reva/pull/895 https://github.com/cs3org/reva/pull/899
   https://github.com/cs3org/reva/pull/900
   https://github.com/owncloud/ocis-reva/pull/313

  * Enhancement - Update reva to v0.1.1-0.20200626111234-e21c32db9614: [#261](https://github.com/owncloud/ocis-reva/pull/261)

   - Updated reva to v0.1.1-0.20200626111234-e21c32db9614 (#304) - TUS upload
   support through datagateway (#261, reva/#878, reva/#888) - Added support for
   differing metrics path for Prometheus to Mentix (reva/#875) - More data exported
   by Mentix (reva/#881) - Implementation of file operations in public folder
   shares (#49, #293, reva/#877) - Make httpclient trust local certificates for now
   (reva/#880) - EOS homes are not configured with an enable-flag anymore, but with
   a dedicated storage driver. We're using it now and adapted default configs of
   storages (reva/#891, #304)

   Https://github.com/owncloud/ocis-reva/issues/49
   https://github.com/owncloud/ocis-reva/issues/293
   https://github.com/owncloud/ocis-reva/issues/261
   https://github.com/owncloud/ocis-reva/pull/261
   https://github.com/cs3org/reva/pull/875 https://github.com/cs3org/reva/pull/877
   https://github.com/cs3org/reva/pull/878 https://github.com/cs3org/reva/pull/881
   https://github.com/cs3org/reva/pull/880 https://github.com/cs3org/reva/pull/888
   https://github.com/owncloud/ocis-reva/pull/304
   https://github.com/cs3org/reva/pull/891

  * Enhancement - Update reva to v0.1.1-0.20200624063447-db5e6635d5f0: [#279](https://github.com/owncloud/ocis-reva/pull/279)

   - Updated reva to v0.1.1-0.20200624063447-db5e6635d5f0 (#279) - Local storage:
   URL-encode file ids to ease integration with other microservices like WOPI
   (reva/#799) - Mentix fixes (reva/#803, reva/#817) - OCDAV: fix returned
   timestamp format (#116, reva/#805) - OCM: add default prefix (#814) - add the
   content-length header to the responses (reva/#816) - Deps: clean (reva/#818) -
   Fix trashbin listing (#112, #253, #254, reva/#819) - Make the json publicshare
   driver configurable (reva/#820) - TUS: Return metadata headers after direct
   upload (ocis/#216, reva/#813) - Set mtime to storage after simple upload (#174,
   reva/#823, reva/#841) - Configure grpc client to allow for insecure conns and
   skip server certificate verification (reva/#825) - Deployment: simplify config
   with more default values (reva/#826, reva/#837, reva/#843, reva/#848, reva/#842)
   - Separate local fs into home and with home disabled (reva/#829) - Register
   reflection after other services (reva/#831) - Refactor EOS fs (reva/#830) - Add
   ocs-share-permissions to the propfind response (#47, reva/#836) - OCS: Properly
   read permissions when creating public link (reva/#852) - localfs: make normalize
   return associated error (reva/#850) - EOS grpc driver (reva/#664) - OCS: Add
   support for legacy public link arg publicUpload (reva/#853) - Add cache layer to
   user REST package (reva/#849) - Meshdirectory: pass query params to selected
   provider (reva/#863) - Pass etag in quotes from the fs layer (#269, reva/#866,
   reva/#867) - OCM: use refactored cs3apis provider definition (reva/#864)

   Https://github.com/owncloud/ocis-reva/issues/116
   https://github.com/owncloud/ocis-reva/issues/112
   https://github.com/owncloud/ocis-reva/issues/253
   https://github.com/owncloud/ocis-reva/issues/254
   https://github.com/owncloud/ocis/issues/216
   https://github.com/owncloud/ocis-reva/issues/174
   https://github.com/owncloud/ocis-reva/issues/47
   https://github.com/owncloud/ocis-reva/issues/269
   https://github.com/owncloud/ocis-reva/pull/279
   https://github.com/owncloud/cs3org/reva/pull/799
   https://github.com/owncloud/cs3org/reva/pull/803
   https://github.com/owncloud/cs3org/reva/pull/817
   https://github.com/owncloud/cs3org/reva/pull/805
   https://github.com/owncloud/cs3org/reva/pull/814
   https://github.com/owncloud/cs3org/reva/pull/816
   https://github.com/owncloud/cs3org/reva/pull/818
   https://github.com/owncloud/cs3org/reva/pull/819
   https://github.com/owncloud/cs3org/reva/pull/820
   https://github.com/owncloud/cs3org/reva/pull/823
   https://github.com/owncloud/cs3org/reva/pull/841
   https://github.com/owncloud/cs3org/reva/pull/813
   https://github.com/owncloud/cs3org/reva/pull/825
   https://github.com/owncloud/cs3org/reva/pull/826
   https://github.com/owncloud/cs3org/reva/pull/837
   https://github.com/owncloud/cs3org/reva/pull/843
   https://github.com/owncloud/cs3org/reva/pull/848
   https://github.com/owncloud/cs3org/reva/pull/842
   https://github.com/owncloud/cs3org/reva/pull/829
   https://github.com/owncloud/cs3org/reva/pull/831
   https://github.com/owncloud/cs3org/reva/pull/830
   https://github.com/owncloud/cs3org/reva/pull/836
   https://github.com/owncloud/cs3org/reva/pull/852
   https://github.com/owncloud/cs3org/reva/pull/850
   https://github.com/owncloud/cs3org/reva/pull/664
   https://github.com/owncloud/cs3org/reva/pull/853
   https://github.com/owncloud/cs3org/reva/pull/849
   https://github.com/owncloud/cs3org/reva/pull/863
   https://github.com/owncloud/cs3org/reva/pull/866
   https://github.com/owncloud/cs3org/reva/pull/867
   https://github.com/owncloud/cs3org/reva/pull/864

  * Enhancement - Add TUS global capability: [#177](https://github.com/owncloud/ocis-reva/issues/177)

   The TUS global capabilities from Reva are now exposed.

   The advertised max chunk size can be configured using the
   "--upload-max-chunk-size" CLI switch or "REVA_FRONTEND_UPLOAD_MAX_CHUNK_SIZE"
   environment variable. The advertised http method override can be configured
   using the "--upload-http-method-override" CLI switch or
   "REVA_FRONTEND_UPLOAD_HTTP_METHOD_OVERRIDE" environment variable.

   Https://github.com/owncloud/ocis-reva/issues/177
   https://github.com/owncloud/ocis-reva/pull/228

  * Enhancement - Update reva to v0.1.1-0.20200603071553-e05a87521618: [#244](https://github.com/owncloud/ocis-reva/issues/244)

   - Updated reva to v0.1.1-0.20200603071553-e05a87521618 (#244) - Add option to
   disable TUS on OC layer (#177, reva/#791) - Dataprovider now supports method
   override (#177, reva/#792) - OCS fixes for create public link (reva/#798)

   Https://github.com/owncloud/ocis-reva/issues/244
   https://github.com/owncloud/ocis-reva/issues/177
   https://github.com/cs3org/reva/pull/791 https://github.com/cs3org/reva/pull/792
   https://github.com/cs3org/reva/pull/798

  * Enhancement - Add public shares service: [#49](https://github.com/owncloud/ocis-reva/issues/49)

   Added Public Shares service with CRUD operations and File Public Shares Manager

   Https://github.com/owncloud/ocis-reva/issues/49
   https://github.com/owncloud/ocis-reva/pull/232

  * Enhancement - Update reva to v0.1.1-0.20200529120551-4f2d9c85d3c9: [#49](https://github.com/owncloud/ocis-reva/issues/49)

   - Updated reva to v0.1.1-0.20200529120551 (#232) - Public Shares CRUD, File
   Public Shares Manager (#49, #232, reva/#681, reva/#788) - Disable
   HTTP-KeepAlives to reduce fd count (ocis/#268, reva/#787) - Fix trashbin listing
   (#229, reva/#782) - Create PUT wrapper for TUS uploads (reva/#770) - Add
   security access headers for ocdav requests (#66, reva/#780) - Add option to
   revad cmd to specify logging level (reva/#772) - New metrics package (reva/#740)
   - Remove implicit data member from memory store (reva/#774) - Added TUS global
   capabilities (#177, reva/#775) - Fix PROPFIND with Depth 1 for cross-storage
   operations (reva/#779)

   Https://github.com/owncloud/ocis-reva/issues/49
   https://github.com/owncloud/ocis-reva/issues/229
   https://github.com/owncloud/ocis-reva/issues/66
   https://github.com/owncloud/ocis-reva/issues/177
   https://github.com/owncloud/ocis/issues/268
   https://github.com/owncloud/ocis-reva/pull/232
   https://github.com/cs3org/reva/pull/787 https://github.com/cs3org/reva/pull/681
   https://github.com/cs3org/reva/pull/788 https://github.com/cs3org/reva/pull/782
   https://github.com/cs3org/reva/pull/770 https://github.com/cs3org/reva/pull/780
   https://github.com/cs3org/reva/pull/772 https://github.com/cs3org/reva/pull/740
   https://github.com/cs3org/reva/pull/774 https://github.com/cs3org/reva/pull/775
   https://github.com/cs3org/reva/pull/779

  * Enhancement - Update reva to v0.1.1-0.20200520150229: [#161](https://github.com/owncloud/ocis-reva/pull/161)

   - Update reva to v0.1.1-0.20200520150229 (#161, #180, #192, #207, #221) - Return
   arbitrary metadata with stat, upload without TUS (reva/#766) - Stat file before
   returning datagateway URL when initiating download (reva/#765) - REST driver for
   user package (reva/#747) - Sharing behavior now consistent with the old backend
   (#20, #26, #43, #44, #46, #94 ,reva/#748) - Mentix service (reva/#755) -
   meshdirectory: add mentix driver for gocdb sites integration (reva/#754) - Add
   functionality to commit to storage for OCM shares (reva/#760) - Add option in
   config to disable tus (reva/#759) - ocdav: fix custom property XML parsing in
   PROPPATCH handler (#203, reva/#743) - ocdav: fix PROPPATCH response for removed
   properties (#186, reva/#742) - ocdav: implement PROPFIND infinity depth (#212,
   reva/#758) - Local fs: Allow setting of arbitrary metadata, minor bug fixes
   (reva/#764) - Local fs: metadata handling and share persistence (reva/#732) -
   Local fs: return file owner info in stat (reva/#750) - Fixed regression when
   uploading empty files to OCFS or EOS with PUT and TUS (#188, reva/#734) - On
   delete move the file versions to the trashbin (#94, reva/#731) - Fix OCFS move
   operation (#182, reva/#729) - Fix OCFS custom property / xattr removal
   (reva/#728) - Retry trashbin in case of timestamp collision (reva/#730) -
   Disable chunking v1 by default (reva/#678) - Implement ocs to http status code
   mapping (#26, reva/#696, reva/#707, reva/#711) - Handle the case if directory
   already exists (reva/#695) - Added TUS upload support (reva/#674, reva/#725,
   reva/#717) - Always return file sizes in Webdav PROPFIND (reva/#712) - Use
   default mime type when none was detected (reva/#713) - Fixed Webdav shallow COPY
   (reva/#714) - Fixed arbitrary namespace usage for custom properties in PROPFIND
   (#57, reva/#720) - Implement returning Webdav custom properties from xattr (#57,
   reva/#721) - Minor fix in OCM share pkg (reva/#718)

   Https://github.com/owncloud/ocis-reva/issues/20
   https://github.com/owncloud/ocis-reva/issues/26
   https://github.com/owncloud/ocis-reva/issues/43
   https://github.com/owncloud/ocis-reva/issues/44
   https://github.com/owncloud/ocis-reva/issues/46
   https://github.com/owncloud/ocis-reva/issues/94
   https://github.com/owncloud/ocis-reva/issues/26
   https://github.com/owncloud/ocis-reva/issues/67
   https://github.com/owncloud/ocis-reva/issues/57
   https://github.com/owncloud/ocis-reva/issues/94
   https://github.com/owncloud/ocis-reva/issues/188
   https://github.com/owncloud/ocis-reva/issues/182
   https://github.com/owncloud/ocis-reva/issues/212
   https://github.com/owncloud/ocis-reva/issues/186
   https://github.com/owncloud/ocis-reva/issues/203
   https://github.com/owncloud/ocis-reva/pull/161
   https://github.com/owncloud/ocis-reva/pull/180
   https://github.com/owncloud/ocis-reva/pull/192
   https://github.com/owncloud/ocis-reva/pull/207
   https://github.com/owncloud/ocis-reva/pull/221
   https://github.com/cs3org/reva/pull/766 https://github.com/cs3org/reva/pull/765
   https://github.com/cs3org/reva/pull/755 https://github.com/cs3org/reva/pull/754
   https://github.com/cs3org/reva/pull/747 https://github.com/cs3org/reva/pull/748
   https://github.com/cs3org/reva/pull/760 https://github.com/cs3org/reva/pull/759
   https://github.com/cs3org/reva/pull/678 https://github.com/cs3org/reva/pull/696
   https://github.com/cs3org/reva/pull/707 https://github.com/cs3org/reva/pull/711
   https://github.com/cs3org/reva/pull/695 https://github.com/cs3org/reva/pull/674
   https://github.com/cs3org/reva/pull/725 https://github.com/cs3org/reva/pull/717
   https://github.com/cs3org/reva/pull/712 https://github.com/cs3org/reva/pull/713
   https://github.com/cs3org/reva/pull/720 https://github.com/cs3org/reva/pull/718
   https://github.com/cs3org/reva/pull/731 https://github.com/cs3org/reva/pull/734
   https://github.com/cs3org/reva/pull/729 https://github.com/cs3org/reva/pull/728
   https://github.com/cs3org/reva/pull/730 https://github.com/cs3org/reva/pull/758
   https://github.com/cs3org/reva/pull/742 https://github.com/cs3org/reva/pull/764
   https://github.com/cs3org/reva/pull/743 https://github.com/cs3org/reva/pull/732
   https://github.com/cs3org/reva/pull/750

  * Bugfix - Stop advertising unsupported chunking v2: [#145](https://github.com/owncloud/ocis-reva/pull/145)

   Removed "chunking" attribute in the DAV capabilities. Please note that chunking
   v2 is advertised as "chunking 1.0" while chunking v1 is the attribute
   "bigfilechunking" which is already false.

   Https://github.com/owncloud/ocis-reva/pull/145

  * Enhancement - Allow configuring the gateway for dataproviders: [#136](https://github.com/owncloud/ocis-reva/pull/136)

   This allows using basic or bearer auth when directly talking to dataproviders.

   Https://github.com/owncloud/ocis-reva/pull/136

  * Enhancement - Use a configured logger on reva runtime: [#153](https://github.com/owncloud/ocis-reva/pull/153)

   For consistency reasons we need a configured logger that is inline with an ocis
   logger, so the log cascade can be easily parsed by a human.

   Https://github.com/owncloud/ocis-reva/pull/153

  * Bugfix - Fix eos user sharing config: [#127](https://github.com/owncloud/ocis-reva/pull/127)

   We have added missing config options for the user sharing manager and added a
   dedicated eos storage command with pre configured settings for the eos-docker
   container. It configures a `Shares` folder in a users home when using eos as the
   storage driver.

   Https://github.com/owncloud/ocis-reva/pull/127

  * Enhancement - Update reva to v1.1.0-20200414133413: [#127](https://github.com/owncloud/ocis-reva/pull/127)

   Adds initial public sharing and ocm implementation.

   Https://github.com/owncloud/ocis-reva/pull/127

  * Bugfix - Fix eos config: [#125](https://github.com/owncloud/ocis-reva/pull/125)

   We have added missing config options for the home layout to the config struct
   that is passed to eos.

   Https://github.com/owncloud/ocis-reva/pull/125

  * Bugfix - Set correct flag type in the flagsets: [#75](https://github.com/owncloud/ocis-reva/issues/75)

   While upgrading to the micro/cli version 2 there where two instances of
   `StringFlag` which had not been changed to `StringSliceFlag`. This caused
   `ocis-reva users` and `ocis-reva storage-root` to fail on startup.

   Https://github.com/owncloud/ocis-reva/issues/75
   https://github.com/owncloud/ocis-reva/pull/76

  * Bugfix - We fixed a typo in the `REVA_LDAP_SCHEMA_MAIL` environment variable: [#113](https://github.com/owncloud/ocis-reva/pull/113)

   It was misspelled as `REVA_LDAP_SCHEMA_Mail`.

   Https://github.com/owncloud/ocis-reva/pull/113

  * Bugfix - Allow different namespaces for /webdav and /dav/files: [#68](https://github.com/owncloud/ocis-reva/pull/68)

   After fbf131c the path for the "new" webdav path does not contain a username
   `/remote.php/dav/files/textfile0.txt`. It used to be
   `/remote.php/dav/files/oc/einstein/textfile0.txt` So it lost `oc/einstein`.

   This PR allows setting up different namespaces for `/webav` and `/dav/files`:

   `/webdav` is jailed into `/home` - which uses the home storage driver and uses
   the logged in user to construct the path `/dav/files` is jailed into `/oc` -
   which uses the owncloud storage driver and expects a username as the first path
   segment

   This mimics oc10

   The `WEBDAV_NAMESPACE_JAIL` environment variable is split into -
   `WEBDAV_NAMESPACE` and - `DAV_FILES_NAMESPACE` accordingly.

   Https://github.com/owncloud/ocis-reva/pull/68 related:

  * Change - Use /home as default namespace: [#68](https://github.com/owncloud/ocis-reva/pull/68)

   Currently, cross storage etag propagation is not yet implemented, which prevents
   the desktop client from detecting changes via the PROPFIND to /. / is managed by
   the root storage provider which is independent of the home and oc storage
   providers. If a file changes in /home/foo, the etag change will only be
   propagated to the root of the home storage provider.

   This change jails users into the `/home` namespace, and allows configuring the
   namespace to use for the two webdav endpoints using the new environment variable
   `WEBDAV_NAMESPACE_JAIL` which affects both endpoints `/dav/files` and `/webdav`.

   This will allow us to focus on getting a single storage driver like eos or
   owncloud tested and better resembles what owncloud 10 does.

   To get back the global namespace, which ultimately is the goal, just set the
   above environment variable to `/`.

   Https://github.com/owncloud/ocis-reva/pull/68

  * Change - Initial release of basic version: [#1](https://github.com/owncloud/ocis-reva/issues/1)

   Just prepared an initial basic version to start a reva server and start
   integrating with the go-micro base dextension framework of ownCloud Infinite
   Scale.

   Https://github.com/owncloud/ocis-reva/issues/1

  * Change - Start multiple services with dedicated commands: [#6](https://github.com/owncloud/ocis-reva/issues/6)

   The initial version would only allow us to use a set of reva configurations to
   start multiple services. We use a more opinionated set of commands to start
   dedicated services that allows us to configure them individually. It allows us
   to switch eg. the user backend to LDAP and fully use it on the cli.

   Https://github.com/owncloud/ocis-reva/issues/6

  * Change - Storage providers now default to exposing data servers: [#89](https://github.com/owncloud/ocis-reva/issues/89)

   The flags that let reva storage providers announce that they expose a data
   server now defaults to true:

   `REVA_STORAGE_HOME_EXPOSE_DATA_SERVER=1` `REVA_STORAGE_OC_EXPOSE_DATA_SERVER=1`

   Https://github.com/owncloud/ocis-reva/issues/89

  * Change - Default to running behind ocis-proxy: [#113](https://github.com/owncloud/ocis-reva/pull/113)

   We changed the default configuration to integrate better with ocis.

   - We use ocis-glauth as the default ldap server on port 9125 with base
   `dc=example,dc=org`. - We use a dedicated technical `reva` user to make ldap
   binds - Clients are supposed to use the ocis-proxy endpoint
   `https://localhost:9200` - We removed unneeded ocis configuration from the
   frontend which no longer serves an oidc provider. - We changed the default user
   OpaqueID attribute from `sub` to `preferred_username`. The latter is a claim
   populated by konnectd that can also be used by the reva ldap user manager to
   look up users by their OpaqueId

   Https://github.com/owncloud/ocis-reva/pull/113

  * Enhancement - Expose owncloud storage driver config in flagset: [#87](https://github.com/owncloud/ocis-reva/issues/87)

   Three new flags are now available:

   - scan files on startup to generate missing fileids default: `true` env var:
   `REVA_STORAGE_OWNCLOUD_SCAN` cli option: `--storage-owncloud-scan`

   - autocreate home path for new users default: `true` env var:
   `REVA_STORAGE_OWNCLOUD_AUTOCREATE` cli option: `--storage-owncloud-autocreate`

   - the address of the redis server default: `:6379` env var:
   `REVA_STORAGE_OWNCLOUD_REDIS_ADDR` cli option: `--storage-owncloud-redis`

   Https://github.com/owncloud/ocis-reva/issues/87

  * Enhancement - Update reva to v0.0.2-0.20200212114015-0dbce24f7e8b: [#91](https://github.com/owncloud/ocis-reva/pull/91)

   Reva has seen a lot of changes that allow us to - reduce the configuration
   overhead - use the autocreate home folder option - use the home folder path
   layout option - no longer start the root storage

   Https://github.com/owncloud/ocis-reva/pull/91 related:

  * Enhancement - Allow configuring user sharing driver: [#115](https://github.com/owncloud/ocis-reva/pull/115)

   We now default to `json` which persists shares in the sharing manager in a json
   file instead of an in memory db.

   Https://github.com/owncloud/ocis-reva/pull/115

   https://github.com/owncloud/product/issues/244

* Enhancement - Add the store service: [#244](https://github.com/owncloud/product/issues/244)

   Tags: store

  * Enhancement - Add version command: [#226](https://github.com/owncloud/product/issues/226)
  * Bugfix - Removed code from other service: [#7](https://github.com/owncloud/ocis-store/pull/7)
  * Bugfix - Build docker images with alpine:latest instead of alpine:edge: [#5](https://github.com/owncloud/ocis-store/pull/5)
  * Change - Initial release of basic version: [#1](https://github.com/owncloud/ocis-store/pull/1)

   https://github.com/owncloud/product/issues/244

* Enhancement - Add the thumbnails service: [#244](https://github.com/owncloud/product/issues/244)

   Tags: thumbnails

  * Enhancement - Add version command: [#226](https://github.com/owncloud/product/issues/226)
  * Bugfix - Build docker images with alpine:latest instead of alpine:edge: [#35](https://github.com/owncloud/ocis-thumbnails/pull/35)
  * Enhancement - Serve the metrics endpoint: [#37](https://github.com/owncloud/ocis-thumbnails/issues/37)
  * Change - Add more default resolutions: [#23](https://github.com/owncloud/ocis-thumbnails/issues/23)
  * Change - Refactor code to remove code smells: [#21](https://github.com/owncloud/ocis-thumbnails/issues/21)
  * Change - Use micro service error api: [#31](https://github.com/owncloud/ocis-thumbnails/issues/31)
  * Enhancement - Limit users to access own thumbnails: [#5](https://github.com/owncloud/ocis-thumbnails/issues/5)
  * Bugfix - Fix usage of context.Context: [#18](https://github.com/owncloud/ocis-thumbnails/issues/18)
  * Bugfix - Fix execution when passing program flags: [#15](https://github.com/owncloud/ocis-thumbnails/issues/15)
  * Change - Initial release of basic version: [#1](https://github.com/owncloud/ocis-thumbnails/issues/1)
  * Change - Use predefined resolutions for thumbnail generation: [#7](https://github.com/owncloud/ocis-thumbnails/issues/7)
  * Change - Implement the first working version: [#3](https://github.com/owncloud/ocis-thumbnails/pull/3)

   https://github.com/owncloud/product/issues/244

* Enhancement - Add the webdav service: [#244](https://github.com/owncloud/product/issues/244)

   Tags: webdav

  * Enhancement - Add version command: [#226](https://github.com/owncloud/product/issues/226)
  * Bugfix - Build docker images with alpine:latest instead of alpine:edge: [#22](https://github.com/owncloud/ocis-webdav/pull/22)
  * Change Change status not found on missing thumbnail: [#20](https://github.com/owncloud/ocis-webdav/issues/20)
  * Change - Initial release of basic version: [#1](https://github.com/owncloud/ocis-webdav/issues/1)
  * Change - Update ocis-pkg to version 2.2.0: [#16](https://github.com/owncloud/ocis-webdav/issues/16)
  * Enhancement - Configuration: [#14](https://github.com/owncloud/ocis-webdav/pull/14)
  * Enhancement - Implement preview API: [#13](https://github.com/owncloud/ocis-webdav/pull/13)

   https://github.com/owncloud/product/issues/244

* Enhancement - Launch a storage to store ocis-metadata: [#602](https://github.com/owncloud/ocis/pull/602)

   Tags: metadata, accounts, settings

   In the future accounts, settings etc. should be stored in a dedicated metadata
   storage. The services should talk to this storage directly, bypassing
   reva-gateway.

   https://github.com/owncloud/ocis/pull/602

* Enhancement - Add basic auth option: [#627](https://github.com/owncloud/ocis/pull/627)

   We added a new `enable-basic-auth` option and `PROXY_ENABLE_BASIC_AUTH`
   environment variable that can be set to `true` to make the proxy verify the
   basic auth header with the accounts service. This should only be used for
   testing and development and is disabled by default.

   https://github.com/owncloud/product/issues/198
   https://github.com/owncloud/ocis/pull/627

* Enhancement - Add glauth fallback backend: [#649](https://github.com/owncloud/ocis/pull/649)

   We introduced the `fallback-datastore` config option and the corresponding
   options to allow configuring a simple chain of two handlers.

   Simple, because it is intended for bind and single result search queries.
   Merging large sets of results is currently out of scope. For now, the
   implementation will only search the fallback backend if the default backend
   returns an error or the number of results is 0. This is sufficient to allow an
   IdP to authenticate users from ocis as well as owncloud 10 as described in the
   [bridge scenario](https://owncloud.github.io/ocis/deployment/bridge/).

   https://github.com/owncloud/ocis-glauth/issues/18
   https://github.com/owncloud/ocis/pull/649

* Enhancement - Update reva to dd3a8c0f38: [#725](https://github.com/owncloud/ocis/pull/725)

  * fixes etag propagation in the ocis driver

   https://github.com/owncloud/ocis/pull/725
   https://github.com/cs3org/reva/pull/1264

* Enhancement - Update konnectd to v0.33.8: [#744](https://github.com/owncloud/ocis/pull/744)

   This update adds options which allow the configuration of oidc-token expiration
   parameters: KONNECTD_ACCESS_TOKEN_EXPIRATION, KONNECTD_ID_TOKEN_EXPIRATION and
   KONNECTD_REFRESH_TOKEN_EXPIRATION.

   Other changes from upstream:

   - Generate random endsession state for external authority - Update dependencies
   in Dockerfile - Set prompt=None to avoid loops with external authority - Update
   Jenkins reporting plugin from checkstyle to recordIssues - Remove extra kty key
   from JWKS top level document - Fix regression which encodes URL fragments twice
   - Avoid generating fragment/query URLs with wrong order - Return state for oidc
   endsession response redirects - Use server provided username to avoid case
   mismatch - Use signed-out-uri if set as fallback for goodbye redirect on saml
   slo - Add checks to ensure post_logout_redirect_uri is not empty - Fix SAML2
   logout request parsing - Cure panic when no state is found in saml esr - Use
   SAML IdP Issuer value from meta data entityID - Allow configuration of
   expiration of oidc access, id and refresh tokens - Implement trampolin for
   external OIDC authority end session - Update ca-certificates version

   https://github.com/owncloud/ocis/pull/744

* Enhancement - Update reva to cdb3d6688da5: [#748](https://github.com/owncloud/ocis/pull/748)

  * let the gateway filter invalid references

   https://github.com/owncloud/ocis/pull/748
   https://github.com/cs3org/reva/pull/1274

* Enhancement - Update glauth to dev 4f029234b2308: [#786](https://github.com/owncloud/ocis/pull/786)

   Includes a bugfix, don't mix graph and provisioning api.

   https://github.com/owncloud/ocis/pull/786

* Enhancement - Update reva to v1.4.1-0.20201123062044-b2c4af4e897d: [#823](https://github.com/owncloud/ocis/pull/823)

  * Refactor the uploading files workflow from various clients [cs3org/reva#1285](https://github.com/cs3org/reva/pull/1285), [cs3org/reva#1314](https://github.com/cs3org/reva/pull/1314)
  * [OCS] filter share with me requests [cs3org/reva#1302](https://github.com/cs3org/reva/pull/1302)
  * Fix listing shares for nonexistent path [cs3org/reva#1316](https://github.com/cs3org/reva/pull/1316)
  * prevent nil pointer when listing shares [cs3org/reva#1317](https://github.com/cs3org/reva/pull/1317)
  * Sharee retrieves the information about a share -but gets response containing all the shares [owncloud/ocis-reva#260](https://github.com/owncloud/ocis-reva/issues/260)
  * Deleting a public link after renaming a file [owncloud/ocis-reva#311](https://github.com/owncloud/ocis-reva/issues/311)
  * Avoid log spam [cs3org/reva#1323](https://github.com/cs3org/reva/pull/1323), [cs3org/reva#1324](https://github.com/cs3org/reva/pull/1324)
  * Fix trashbin [cs3org/reva#1326](https://github.com/cs3org/reva/pull/1326)

   https://github.com/owncloud/ocis-reva/issues/260
   https://github.com/owncloud/ocis-reva/issues/311
   https://github.com/owncloud/ocis/pull/823
   https://github.com/cs3org/reva/pull/1285
   https://github.com/cs3org/reva/pull/1302
   https://github.com/cs3org/reva/pull/1314
   https://github.com/cs3org/reva/pull/1316
   https://github.com/cs3org/reva/pull/1317
   https://github.com/cs3org/reva/pull/1323
   https://github.com/cs3org/reva/pull/1324
   https://github.com/cs3org/reva/pull/1326

* Enhancement - Update glauth to dev fd3ac7e4bbdc93578655d9a08d8e23f105aaa5b2: [#834](https://github.com/owncloud/ocis/pull/834)

   We updated glauth to dev commit fd3ac7e4bbdc93578655d9a08d8e23f105aaa5b2, which
   allows to skip certificate checks for the owncloud backend.

   https://github.com/owncloud/ocis/pull/834

* Enhancement - Better adopt Go-Micro: [#840](https://github.com/owncloud/ocis/pull/840)

   Tags: ocis

   There are a few building blocks that we were relying on default behavior, such
   as `micro.Registry` and the go-micro client. In order for oCIS to work in any
   environment and not relying in black magic configuration or running daemons we
   need to be able to:

   - Provide with a configurable go-micro registry. - Use our own go-micro client
   adjusted to our own needs (i.e: custom timeout, custom dial timeout, custom
   transport...)

   This PR is relying on 2 env variables from Micro: `MICRO_REGISTRY` and
   `MICRO_REGISTRY_ADDRESS`. The latter does not make sense to provide if the
   registry is not `etcd`.

   The current implementation only accounts for `mdns` and `etcd` registries,
   defaulting to `mdns` when not explicitly defined to use `etcd`.

   https://github.com/owncloud/ocis/pull/840

* Enhancement - Tidy dependencies: [#845](https://github.com/owncloud/ocis/pull/845)

   Methodology:

   ```
   go-modules() {
     find . \( -name vendor -o -name '[._].*' -o -name node_modules \) -prune -o -name go.mod -print | sed 's:/go.mod$::'
   }
   ```

   ```
   for m in $(go-modules); do (cd $m && go mod tidy); done
   ```

   https://github.com/owncloud/ocis/pull/845

* Enhancement - Create OnlyOffice extension: [#857](https://github.com/owncloud/ocis/pull/857)

   Tags: OnlyOffice

   We've created an OnlyOffice extension which enables users to create and edit
   docx documents and open spreadsheets and presentations.

   https://github.com/owncloud/ocis/pull/857

* Enhancement - Cache userinfo in proxy: [#877](https://github.com/owncloud/ocis/pull/877)

   Tags: proxy

   We introduced caching for the userinfo response. The token expiration is used
   for cache invalidation if available. Otherwise we fall back to a preconfigured
   TTL (default 10 seconds).

   https://github.com/owncloud/ocis/pull/877

* Enhancement - Add permission check when assigning and removing roles: [#879](https://github.com/owncloud/ocis/issues/879)

   Everyone could add and remove roles from users. Added a new permission and a
   check so that only users with the role management permissions can assign and
   unassign roles.

   https://github.com/owncloud/ocis/issues/879

* Enhancement - Show basic-auth warning only once: [#886](https://github.com/owncloud/ocis/pull/886)

   Show basic-auth warning only on startup instead on every request.

   https://github.com/owncloud/ocis/pull/886

* Enhancement - Create a proxy access-log: [#889](https://github.com/owncloud/ocis/pull/889)

   Logs client access at the proxy

   https://github.com/owncloud/ocis/pull/889

* Enhancement - Add a version command to ocis: [#915](https://github.com/owncloud/ocis/pull/915)

   The version command was only implemented in the extensions. This adds the
   version command to ocis to list all services in the ocis namespace.

   https://github.com/owncloud/ocis/pull/915

* Enhancement - Add k6: [#941](https://github.com/owncloud/ocis/pull/941)

   Tags: tests

   Add k6 as a performance testing framework

   https://github.com/owncloud/ocis/pull/941
   https://github.com/owncloud/ocis/pull/983

* Enhancement - Update reva to v1.4.1-0.20201127111856-e6a6212c1b7b: [#971](https://github.com/owncloud/ocis/pull/971)

   Tags: reva

  * Fix capabilities response for multiple client versions #1331 [cs3org/reva#1331](https://github.com/cs3org/reva/pull/1331)
  * Fix home storage redirect for remote.php/dav/files [cs3org/reva#1342](https://github.com/cs3org/reva/pull/1342)

   https://github.com/owncloud/ocis/pull/971
   https://github.com/cs3org/reva/pull/1331
   https://github.com/cs3org/reva/pull/1342

* Enhancement - Update reva to v1.4.1-0.20201130061320-ac85e68e0600: [#980](https://github.com/owncloud/ocis/pull/980)

  * Fix move operation in ocis storage driver [csorg/reva#1343](https://github.com/cs3org/reva/pull/1343)

   https://github.com/owncloud/ocis/issues/975
   https://github.com/owncloud/ocis/pull/980
   https://github.com/cs3org/reva/pull/1343

* Enhancement - Add www-authenticate based on user agent: [#1009](https://github.com/owncloud/ocis/pull/1009)

   Tags: reva, proxy

   We now comply with HTTP spec by adding Www-Authenticate headers on every `401`
   request. Furthermore, we not only take care of such a thing at the Proxy but
   also Reva will take care of it. In addition, we now are able to lock-in a set of
   User-Agent to specific challenges.

   Admins can use this feature by configuring oCIS + Reva following this approach:

   ```
   STORAGE_FRONTEND_MIDDLEWARE_AUTH_CREDENTIALS_BY_USER_AGENT="mirall:basic, Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:83.0) Gecko/20100101 Firefox/83.0:bearer" \
   PROXY_MIDDLEWARE_AUTH_CREDENTIALS_BY_USER_AGENT="mirall:basic, Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:83.0) Gecko/20100101 Firefox/83.0:bearer" \
   PROXY_ENABLE_BASIC_AUTH=true \
   go run cmd/ocis/main.go server
   ```

   We introduced two new environment variables:

   `STORAGE_FRONTEND_MIDDLEWARE_AUTH_CREDENTIALS_BY_USER_AGENT` as well as
   `PROXY_MIDDLEWARE_AUTH_CREDENTIALS_BY_USER_AGENT`, The reason they have the same
   value is not to rely on the os env on a distributed environment, so in
   redundancy we trust. They both configure the same on the backend storage and
   oCIS Proxy.

   https://github.com/owncloud/ocis/pull/1009

* Enhancement - Add tracing to the accounts service: [#1016](https://github.com/owncloud/ocis/issues/1016)

   Added tracing to the accounts service.

   https://github.com/owncloud/ocis/issues/1016

* Enhancement - Runtime Cleanup: [#1066](https://github.com/owncloud/ocis/pull/1066)

   Small runtime cleanup prior to Tech Preview release

   https://github.com/owncloud/ocis/pull/1066

* Enhancement - Update reva to 063b3db9162b: [#1091](https://github.com/owncloud/ocis/pull/1091)

   - bring public link removal changes to OCIS. - fix subcommand name collision
   from renaming phoenix -> web.

   https://github.com/owncloud/ocis/issues/1098
   https://github.com/owncloud/ocis/pull/1091

* Enhancement - Update OCIS Runtime: [#1108](https://github.com/owncloud/ocis/pull/1108)

   - enhances the overall behavior of our runtime - runtime `db` file configurable
   - two new env variables to deal with the runtime - `RUNTIME_DB_FILE` and
   `RUNTIME_KEEP_ALIVE` - `RUNTIME_KEEP_ALIVE` defaults to `false` to provide
   backwards compatibility - if `RUNTIME_KEEP_ALIVE` is set to `true`, if a
   supervised process terminates the runtime will attempt to start with the same
   environment provided.

   https://github.com/owncloud/ocis/pull/1108

* Enhancement - Update reva to v1.4.1-0.20201125144025-57da0c27434c: [#1320](https://github.com/cs3org/reva/pull/1320)

   Mostly to bring fixes to pressing changes.

   https://github.com/cs3org/reva/pull/1320
   https://github.com/cs3org/reva/pull/1338
