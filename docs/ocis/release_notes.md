---
title: "Release Notes"
date: 2020-12-16T20:35:00+01:00
weight: 0
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis
geekdocFilePath: release_notes.md
---

## ownCloud Infinite Scale 1.10.0 Technology Preview

Version 1.10.0 brings new features, usability improvements and bug fixes. ownCloud Web 4.0.0 now supports ONLYOFFICE document editors and can search/filter files and folders. Furthermore it brings a new context menu for file actions that can be accessed via right click and comes with a big bunch of other notable improvements and fixes.

The most prominent changes in ownCloud Infinite Scale 1.10.0 and ownCloud Web 4.0.0 comprise:

- ownCloud Web now supports ONLYOFFICE document editors when used with ownCloud Classic Server. See the [documentation](https://owncloud.dev/clients/web/deployments/oc10-app/#onlyoffice) for more information on requirements and configuration.
- ownCloud Web now supports global search and filtering for the current folder via the search bar. Both will work when ownCloud Web is used with ownCloud Classic. The Infinite Scale capabilities are currently limited to filtering the current folder. [#5415](https://github.com/owncloud/web/pull/5415)
- A context menu for a file/folder which contains related actions has been introduced to ownCloud Web (in addition to the actions in the right sidebar). [#5160](https://github.com/owncloud/web/issues/5160)
- The context menu for a file/folder in ownCloud Web can be opened via right click and using the "..." menu. [#5102](https://github.com/owncloud/web/issues/5102)
- As a first step of a larger redesign of the sharing dialog in ownCloud Web, the autocomplete and share recipient selection have been redesigned. [#5554](https://github.com/owncloud/web/pull/5554)
- The right sidebar navigation in ownCloud Web has been redesigned. Moving away from structuring all functionality on a single view using accordions, each section now has their own, dedicated view. [#5549](https://github.com/owncloud/web/pull/5549)
- The maximum number of sharing autocomplete suggestions in ownCloud Web can now be configured. See [the documentation](https://owncloud.dev/clients/web/getting-started/#options) for more information. [#5506](https://github.com/owncloud/web/pull/5506)
- ownCloud Web works now with ownCloud Classic when OpenID Connect authentication is used. [#5536](https://github.com/owncloud/web/pull/5536)
- ownCloud Web now respects the server-side capability for user avatars. [#5178](https://github.com/owncloud/web/pull/5178)
- The login page has been optimized in regards of accessibility. [#5376](https://github.com/owncloud/web/issues/5376)
- The Infinite Scale backend is being further hardened by fixing known issues, improving error handling and stabilizing existing features. 

You can also read the full [ownCloud Infinite Scale changelog](https://github.com/owncloud/ocis/releases/tag/v1.10.0) and [ownCloud Web changelog](https://github.com/owncloud/web/releases/tag/v4.0.0) for further details on what has changed.

### Breaking changes
{{< hint warning >}}
We are currently in a Tech Preview state and breaking changes may occur at any time. For more information see our [release roadmap]({{< ref "./release_roadmap" >}})
{{< /hint >}}

## ownCloud Infinite Scale 1.9.0 Technology Preview

Version 1.9.0 is a feature and maintenance release. More features have been added and the platform was matured further. ownCloud Web 3.4.1 brings usability improvements and new features. The right sidebar now shows details about the selected resource and offers previews for images. View options for the file list and a feedback button have been added.

The most prominent changes in ownCloud Infinite Scale 1.9.0 and ownCloud Web 3.4.1 comprise:

- The right sidebar in ownCloud Web now shows details about the selected file/folder (e.g., size, owner, sharing status, modification time). [#5161](https://github.com/owncloud/web/issues/5161)
- The right sidebar in ownCloud Web now shows previews for images. [#5501](https://github.com/owncloud/web/pull/5501)
- View options for the file list have been introduced in ownCloud Web. Currently this allows to change the number of files/folders per page and to show/hide hidden files. [#5408]https://github.com/owncloud/web/pull/5408 [#5470](https://github.com/owncloud/web/pull/5470)
- A feedback button has been added to the top bar. It guides to user to an ownCloud Web feedback survey. If undesired, this feature [can be disabled in the ownCloud Web configuration](https://owncloud.dev/clients/web/getting-started/#options). [#5468](https://github.com/owncloud/web/pull/5468)
- Received shares can now be accepted/declined as batches in the "Shared with me" view. [#5374](https://github.com/owncloud/web/pull/5374)
- The oCIS backend now supports to enable extensions by name. [#2229](https://github.com/owncloud/ocis/pull/2229)
- Storage drivers can be set to read only. [#2230](https://github.com/owncloud/ocis/pull/2230)
- Micro service init has been improved for faster startup. [#1705](https://github.com/owncloud/ocis/pull/1705)

You can also read the full [ownCloud Infinite Scale changelog](https://github.com/owncloud/ocis/releases/tag/v1.9.0) and [ownCloud Web changelog](https://github.com/owncloud/web/releases/tag/v3.4.1) for further details on what has changed.

### Breaking changes
{{< hint warning >}}
We are currently in a Tech Preview state and breaking changes may occur at any time. For more information see our [release roadmap]({{< ref "./release_roadmap" >}})
{{< /hint >}}

## ownCloud Infinite Scale 1.8.0 Technology Preview

Version 1.8.0 is a maintenance and bug fix release. ownCloud Web 3.3.0 has received further performance and major accessibility improvements.

The most prominent changes in ownCloud Infinite Scale 1.8.0 and ownCloud Web 3.3.0 comprise:

- ownCloud Web is now fully translatable on Transifex [#5042](https://github.com/owncloud/web/pull/5042)
- ownCloud Web now supports keyboard navigation [#4937](https://github.com/owncloud/web/pull/4937) [#5013](https://github.com/owncloud/web/pull/5013) [#5027](https://github.com/owncloud/web/pull/5027) [#5147](https://github.com/owncloud/web/pull/5147)
- ownCloud Web now supports screenreaders [#5182](https://github.com/owncloud/web/pull/5182) [#5166](https://github.com/owncloud/web/pull/5166) [#5058](https://github.com/owncloud/web/pull/5058) [#5046](https://github.com/owncloud/web/pull/5046) [#5010](https://github.com/owncloud/web/pull/5010)
- ownCloud Web has received many performance improvements (image cache, fixes to avoid duplicate resource loading, asynchronous image loading) [#5194](https://github.com/owncloud/web/pull/5194)
- The file lists in ownCloud Web are now paginated to control loading times [#5224](https://github.com/owncloud/web/pull/5224) [#5309](https://github.com/owncloud/web/pull/5309)
- ownCloud Web now supports TypeScript [#5194](https://github.com/owncloud/web/pull/5194)

You can also read the full [ownCloud Infinite Scale changelog](https://github.com/owncloud/ocis/releases/tag/v1.8.0) and [ownCloud Web changelog](https://github.com/owncloud/web/releases/tag/v3.3.0) for further details on what has changed.

### Breaking changes
{{< hint warning >}}
We are currently in a Tech Preview state and breaking changes may occur at any time. For more information see our [release roadmap]({{< ref "./release_roadmap" >}})
{{< /hint >}}

## ownCloud Infinite Scale 1.7.0 Technology Preview

Version 1.7.0 is a maintenance and bug fix release. ownCloud Web 3.2.0 has received further performance improvements and minor usability tweaks.

The most prominent changes in ownCloud Infinite Scale 1.7.0 and ownCloud Web 3.2.0 comprise:

- The S3 storage driver can now be used for testing using the configuration values in the [documentation](https://owncloud.dev/extensions/storage/configuration/#s3ng-driver) [#1886](https://github.com/owncloud/ocis/pull/1886)
- A confirmation dialog for public link deletion has been added [#5125](https://github.com/owncloud/web/pull/5125)
- To improve performance, the file types which are being rendered as previews can now be specified using an [allow list in config.json](https://owncloud.dev/clients/web/getting-started/#options) [#5159](https://github.com/owncloud/web/pull/5159)
- A warning has been added when a user tries to leave the page while an operation is in progress (e.g., an upload) [#2590](https://github.com/owncloud/web/issues/2590)

You can also read the full [ownCloud Infinite Scale changelog](https://github.com/owncloud/ocis/releases/tag/v1.7.0) and [ownCloud Web changelog](https://github.com/owncloud/web/releases/tag/v3.2.0) for further details on what has changed.

### Breaking changes
{{< hint warning >}}
We are currently in a Tech Preview state and breaking changes may occur at any time. For more information see our [release roadmap]({{< ref "./release_roadmap" >}})
{{< /hint >}}

## ownCloud Infinite Scale 1.6.0 Technology Preview

To get the full potential out of the microservice architecture, version 1.6.0 introduces a dynamic service registry to ownCloud Infinite Scale. The dynamic service registry facilitates the configuration and contributes to the scalability of the platform. ownCloud Web 3.1.0 has received further improvements for accessibility like keyboard navigation and it comes with performance improvements by loading certain elements asynchronously.

The most prominent changes in ownCloud Infinite Scale 1.6.0 and ownCloud Web 3.1.0 comprise:

- Introducing a dynamic service registry: The dynamic service registry takes care of dynamically assigning network addresses between the oCIS services and enables the services to find and work with each other automatically. It replaces the previous hardcoded service configuration which simplifies the initial setup and makes distributed, scale-out environments a lot easier to handle. [#1509](https://github.com/cs3org/reva/pull/1509)
- User avatars are now fetched asynchronously, enabling a non-blocking loading of the file list and improving user experience [#1295](https://github.com/owncloud/owncloud-design-system/pull/1295)
- Further accessibility and keyboard navigation improvements have been added [#1979](https://github.com/owncloud/ocis/pull/1979) [#1991](https://github.com/owncloud/ocis/pull/1991) [#4942](https://github.com/owncloud/web/pull/4942) [#4965](https://github.com/owncloud/web/pull/4965) [#4991](https://github.com/owncloud/web/pull/4991)
- The OCS user deprovisioning endpoint has been added, enabling a full user deprovisioning including storage. [#1962](https://github.com/owncloud/ocis/pull/1962)
- Text files (.txt) now have previews (thumbnails) [#1988](https://github.com/owncloud/ocis/pull/1988)
- The translations in the Settings and Accounts extensions have been improved [#2003](https://github.com/owncloud/ocis/pull/2003)

You can also read the full [ownCloud Infinite Scale changelog](https://github.com/owncloud/ocis/releases/tag/v1.6.0) and [ownCloud Web changelog](https://github.com/owncloud/web/releases/tag/v3.1.0) for further details on what has changed.

### Breaking changes
{{< hint warning >}}
We are currently in a Tech Preview state and breaking changes may occur at any time. For more information see our [release roadmap]({{< ref "./release_roadmap" >}})
{{< /hint >}}

#### Changed oCIS JSON share driver storage format

Related: [#1655](https://github.com/cs3org/reva/pull/1655)

The storage format of the oCIS JSON share driver has changed. You will be affected if you plan to update from a previous version of oCIS to oCIS 1.6.0, you have shared files or folders with users or groups and you are using the oCIS JSON share driver, which is currently the default share driver.

Implications:
- manual action required

Our recommended update strategy to oCIS 1.6.0 is:
1. let users note all their shares with users and groups they set up in oCIS
1. stop oCIS
1. move / delete the JSON share driver storage file `/var/tmp/ocis/storage/shares.json`
1. update to oCIS 1.6.0
1. let users recreate their shares

#### Fixed / changed oCIS metadata storage driver filesystem path
Related: [#1956](https://github.com/owncloud/ocis/pull/1956)

The filesystem path of the oCIS metadata storage driver has changed (been fixed). You will be affected if you plan to update from a previous version of oCIS to oCIS 1.6.0 and are using the oCIS storage driver for metadata storage.

Implications:
- manual action required

Our recommended update strategy to oCIS 1.6.0 is:
1. let users backup all their data stored in oCIS
1. stop oCIS
1. prune all oCIS data in `/var/tmp/ocis`
1. update to oCIS 1.6.0
1. recreate user accounts (can be skipped if an external IDP is used)
1. let users upload all their data again
1. let users recreate their shares

If you want to use oCIS 1.6.0 without following our recommended update strategy, you can also keep the pre 1.6.0 behaviour by setting this environment variable:
`export STORAGE_METADATA_ROOT=/var/tmp/ocis/storage/users`
This may lead to faulty behaviour since both the metadata and user storage driver will be storing their data in the same filesystem path.

## ownCloud Infinite Scale 1.5.0 Technology Preview

Version 1.5.0 is a maintenance release for the Infinite Scale backend with a number of bug fixes and smaller improvements. For ownCloud Web it brings further accessibility improvements and a whole bunch of new features. The web interface can now be branded and there is a new, dedicated view in the left sidebar to list all link shares of a user.

The most prominent changes in ownCloud Infinite Scale 1.5.0 and ownCloud Web 3.0.0 comprise:

- Config file based theming for ownCloud Web (see https://owncloud.dev/clients/web/theming/ for more information) [#4822](https://github.com/owncloud/web/pull/4822)
- A dedicated view for "Shared by link" has been added [#4881](https://github.com/owncloud/web/pull/4881)
- The file list table has been replaced and is now more performant and accessible [#4627](https://github.com/owncloud/web/pull/4627)
- Many further accessibility improvements have been added, e.g., around the app switcher, sidebar, sharing list and focus management
- User storage quotas will now be enforced [#1557](https://github.com/cs3org/reva/pull/1557)
- The "owncloud" storage driver now supports file integrity checking with checksums [#1629](https://github.com/cs3org/reva/pull/1629)

You can also read the full [ownCloud Infinite Scale changelog](https://github.com/owncloud/ocis/releases/tag/v1.5.0) and [ownCloud Web changelog](https://github.com/owncloud/web/releases/tag/v3.0.0) for further details on what has changed.

### Breaking changes
{{< hint warning >}}
We are currently in a Tech Preview state and breaking changes may occur at any time. For more information see our [release roadmap]({{< ref "./release_roadmap" >}})
{{< /hint >}}

## ownCloud Infinite Scale 1.4.0 Technology Preview

Version 1.4.0 brings new features, bug fixes and further improvements. The accessibility of ownCloud Web has greatly improved, paving the way for WCAG 2.1 compliance. The Infinite Scale platform has received major improvements regarding memory consumption. The user storage quota feature has been implemented and folder sizes are now properly calculated. It is now possible to write log messages to log files and to specify configuration values using a config file.

The most prominent changes in ownCloud Infinite Scale 1.4.0 and ownCloud Web 2.1.0 comprise:

- ownCloud Web is now able to use pre-signed url downloads for password protected shares [#38376](https://github.com/owncloud/core/pull/38376)
- Reduced the memory consumption of the runtime drastically (by a factor of 24) [#1762](https://github.com/owncloud/ocis/pull/1762)
- Initial quota support to impose storage space restrictions for users (query / set) [#1405](https://github.com/cs3org/reva/pull/1405)
- Folder sizes are now calculated correctly (tree size accounting) [#1405](https://github.com/cs3org/reva/pull/1405)
- Added the possibility to write the log to a file with the option to write separated log files by service [#1816](https://github.com/owncloud/ocis/pull/1816)
- Added the possibility to specify configuration values for the entire platform in a single config file [#1762](https://github.com/owncloud/ocis/pull/1762)
- Added GIF and JPEG file types for thumbnail generation (allows to display thumbnails and use the media viewer for GIF/JPEG images) [#1791](https://github.com/owncloud/ocis/pull/1791)
- Fixes for the trash bin feature [#1552](https://github.com/cs3org/reva/pull/1552)

You can also read the full [ownCloud Infinite Scale changelog](https://github.com/owncloud/ocis/releases/tag/v1.4.0) and [ownCloud Web changelog](https://github.com/owncloud/web/releases/tag/v2.1.0) for further details on what has changed.

### Breaking changes
{{< hint warning >}}
We are currently in a Tech Preview state and breaking changes may occur at any time. For more information see our [release roadmap]({{< ref "./release_roadmap" >}})
{{< /hint >}}

#### Changed oCIS storage driver file layout

Related: [#1452](https://github.com/cs3org/reva/pull/1452)

Despite a breaking change in the oCIS storage driver file layout, data is not automatically migrated. You will be affected if you plan to update from a previous version of oCIS to oCIS 1.4.0 and are using the oCIS storage driver, which is currently the default storage driver.

Implications:
- manual action required

Our recommended update strategy to oCIS 1.4.0 is:
1. let users backup all their data stored in oCIS
1. stop oCIS
1. prune all oCIS data in `/var/tmp/ocis`
1. update to oCIS 1.4.0
1. recreate user accounts (can be skipped if an external IDP is used)
1. let users upload all their data again
1. let users recreate their shares

If you already updated to oCIS 1.4.0 without our recommended update strategy you will see no data in oCIS anymore, even after a downgrade to your previous version of oCIS. But be assured that your data is still there.

You have to follow these steps to be able to access your data again in oCIS:
1. stop oCIS
1. navigate to `/var/tmp/ocis/storage/users/nodes/root/`
1. in this directory you will find directories with UUID as names. These are the home folders of the oCIS users. Find the ones with content your oCIS users uploaded to oCIS.
1. create an temporary directory eg. `/tmp/dereferenced-ocis-storage`
1. copy the data from oCIS to the temporary directory while dereferencing symlinks. On Linux you can do this by running `cp --recursive --dereference /var/tmp/ocis/storage/users/nodes/root/ /tmp/dereferenced-ocis-storage`
1. you now have a backup of all users data in `/tmp/dereferenced-ocis-storage` and can follow our recommended update strategy above


## ownCloud Infinite Scale 1.3.0 Technology Preview
Version 1.3.0 is a regular maintenance and bugfix release. It provides the latest improvements to users and administrators.

### Changes in Reva

[Reva](https://github.com/cs3org/Reva) is one of the fundamental components of oCIS. It has these significant changes:

- Align href URL encoding with oc10 [#1425](https://github.com/cs3org/Reva/pull/1425)
- Fix public link webdav permissions [#1461](https://github.com/cs3org/Reva/pull/1461)
- Purge non-empty dirs from trash-bin [#1429](https://github.com/cs3org/Reva/pull/1429)
- Checksum support [#1400](https://github.com/cs3org/Reva/pull/1400)
- Set quota when creating home directory in EOS [#1477](https://github.com/cs3org/Reva/pull/1477)
- Add functionality to share resources with groups [#1453](https://github.com/cs3org/Reva/pull/1453)
- Add s3ng storage driver, storing blobs in a s3-compatible blobstore [#1428](https://github.com/cs3org/Reva/pull/1428)

### Changes in oCIS

These are the major changes in oCIS:

- Update ownCloud Web to v2.0.2: [#1776](https://github.com/owncloud/ocis/pull/1776)
- Enhancement - Update go-micro to v3.5.1-0.20210217182006-0f0ace1a44a9: [#1670](https://github.com/owncloud/ocis/pull/1670)
- Enhancement - Update reva to v1.6.1-0.20210223065028-53f39499762e: [#1683](https://github.com/owncloud/ocis/pull/1683)
- Enhancement - Add initial nats and kubernetes registry support: [#1697](https://github.com/owncloud/ocis/pull/1697)

More details about this release can be found in the full [ownCloud Infinite Scale changelog](https://github.com/owncloud/ocis/releases/tag/v1.3.0) and [ownCloud Web changelog](https://github.com/owncloud/web/releases/tag/v2.0.2).

### Breaking changes
{{< hint warning >}}
We are currently in a Tech Preview state and breaking changes may occur at any time. For more information see our [release roadmap]({{< ref "./release_roadmap" >}})
{{< /hint >}}

## ownCloud Infinite Scale 1.2.0 Technology Preview
Version 1.2.0 brings more functionality and stability to ownCloud Infinite Scale. ownCloud Web now loads a lot faster and is prepared for the introduction of accessibility features. An initial implementation for S3 storage support is available and file integrity checking has been introduced.

The most prominent changes in ownCloud Infinite Scale 1.2.0 and ownCloud Web 2.0.0 comprise:

- The initial loading time for ownCloud Web has been reduced by handling dependencies more efficiently (the bundle size of ownCloud Web has been drastically reduced) [#4584](https://github.com/owncloud/web/pull/4584)
- Preparations for accessibility features have been implemented to work towards WCAG 2.1 compliance [#4594](https://github.com/owncloud/web/pull/4594)
- Initial S3 storage support is available [#1429](https://github.com/cs3org/reva/issues/1429)
- File integrity checking has been introduced: When uploading files, Infinite Scale now makes sure that the file integrity is protected between server and clients by comparing checksums [#1400](https://github.com/cs3org/reva/issues/1400)
- Public link passwords are now stored as hashes to improve security [#1462](https://github.com/cs3org/reva/issues/1462)

You can also read the full [ownCloud Infinite Scale changelog](https://github.com/owncloud/ocis/releases/tag/v1.2.0) and [ownCloud Web changelog](https://github.com/owncloud/web/releases/tag/v2.0.0) for further details on what has changed.

### Breaking changes
{{< hint warning >}}
We are currently in a Tech Preview state and breaking changes may occur at any time. For more information see our [release roadmap]({{< ref "./release_roadmap" >}})
{{< /hint >}}

#### Fix IDP service user
Related: [#1390](https://github.com/owncloud/ocis/pull/1390), [#1569](https://github.com/owncloud/ocis/issues/1569)

After upgrading oCIS from a previous version to oCIS 1.2.0 you will not be able to login in ownCloud Web

Implications:
- manual action required

Migration steps:
- Stop oCIS
- Open following file `/var/tmp/ocis/storage/metadata/nodes/root/accounts/820ba2a1-3f54-4538-80a4-2d73007e30bf`
- Change password to `$2y$12$ywfGLDPsSlBTVZU0g.2GZOPO8Wap3rVOpm8e3192VlytNdGWH7x72`
- Change onPremisesSamAccountName to `idp`
- Change preferredName to `idp`
- Save the changed file
- Start oCIS
- You now are able to lock back in again.

Please have a look at [how to secure an oCIS instance]({{< ref "./deployment#secure-an-ocis-instance" >}}) since you seem to run it with default secrets.

#### Reset shares
Related: [#1626](https://github.com/owncloud/ocis/pull/1626)

After upgrading oCIS from a previous version to oCIS 1.2.0 you will will not be able to use previous shares or create new shares.

Implications:
- manual action required
- loss of shares (manual resharing is needed, files will not be lost)

Migration steps:
- Stop oCIS
- Delete following file `/var/tmp/ocis/storage/shares.json`
- Start oCIS
- Recreate shares manually

## ownCloud Infinite Scale 1.1.0 Technology Preview

Version 1.1.0 is a hardening and patch release. It ships with the latest version of ownCloud Web and brings a couple of minor improvements. The minor version increase is needed due to non-backwards compatible changes in configuration. The documentation has been updated to reflect the changes. Please note that this version is still a Technology Preview and not suited for production use.

The most prominent changes in ownCloud Infinite Scale 1.1.0 and ownCloud Web 1.0.1 comprise:

- Performance and stability improvements for installations with multiple concurrent users
- Simplified configuration by introducing the new environment variable OCIS_URL
- Beta release of [ownCloud performance scripts](https://github.com/owncloud/cdperf)
- Update ownCloud web to [v1.0.1](https://github.com/owncloud/web/releases/tag/v1.0.1)
- Update reva to [v1.5.1](https://github.com/cs3org/reva/releases/tag/v1.5.1)

You can also read the full [ownCloud Infinite Scale changelog](https://github.com/owncloud/ocis/releases/tag/v1.1.0) for further details on what has changed.

## ownCloud Infinite Scale 1.0.0 Technology Preview

We are pleased to announce the availability of ownCloud Infinite Scale 1.0.0 Technology Preview which is released as the first public version of the new Infinite Scale platform.

### Microservice architecture

ownCloud Infinite Scale is following the microservices architectural pattern. It is implemented as a set of microservices which are independent of each other. They are coupled with well-defined APIs. This architecture fosters a lot of benefits that we were aiming for with the new design for oCIS:

- Every service is independent, comparably small and brings it's own webserver, backend/APIs and frontend components
- Each service runs as a separate service on the system, increasing security and stability
- Scalability:  High performance demands can be fulfilled by scaling and distributing of services
- Testability: Each service can be tested on its own due to well-defined APIs and functionality
- Protocol-driven development using protobuf
- High-performance communication between services through gRPC
- Multi-platform support powered by Golang - only minimal dependency on platform packages
- Cloud-native deployment, update, monitoring, logging, tracing and orchestration strategies

### Key figures

- The all-new ownCloud Web frontend is shipped as part of the platform
- OpenID Connect is the future-proof technology choice for authentication
- An Identity Provider is bundled to ease deployment and operations. It can be replaced with an external OpenID IdP, if desired
- Automatically built and fully maintained Docker containers are available
- Flexible configuration through environment variables, config files or command-line flags
- Database-less architecture - metadata and data are kept together in the storage as a single source of truth
- Native storage capabilities are used where like native versioning and trashbin
- Public APIs like WebDAV and OCS have been kept compatible with ownCloud 10
- A secure and flexible framework to create extensions

#### Supported platforms

- Linux-amd64
- Darwin-amd64
- Experimental: Windows, ARM (e.g., Raspberry Pi, Termux on Android)

#### Client support

All official ownCloud Clients support the Infinite Scale server with the following versions:
- Desktop >= 2.7
- Android >= 2.15
- iOS >= 1.2

### Architecture components

ownCloud Infinite Scale is built as a modular framework in which components can be scaled individually. It consists of

- a user management service
- a settings service
- a frontend service
- a storage backend service
- a built-in IdP
- an application gateway/proxy

These components can be deployed in a multi-tier deployment architecture. See the [documentation]({{< ref "./" >}}) for an overview of the services.

### Operation modes

#### Standalone mode (with oCIS storage driver)

In standalone mode oCIS uses its built-in orchestrator to start all necessary services. This allows you to run oCIS on a single node without any outside dependencies like docker-compose, kubernetes or even a webserver. It will start an OpenID IdP and create a self-signed certificate. You can start right away by navigating to <https://localhost:9200>.

#### Single services scaleout

oCIS allows you to scale individual services using well-known orchestration frameworks like docker-compose, dockerSwarm and kubernetes.

#### Bridge mode with ownCloud 10 backend

For the product transition phase, ownCloud Infinite Scale comes with an operation mode ("bridge mode") that allows a hybrid deployment, between both server generations to operate the new web frontend with ownCloud 10 and Infinite Scale in parallel. This setup allows the ownCloud Web frontend to operate with both server generations and provides the foundation to migrate users gradually to the new backend.

**Requirements for the bridge mode**
- ownCloud Server >= 10.6
- [Open ID Connect](https://marketplace.owncloud.com/apps/openidconnect) is used for user authentication
- The [Graph API](https://marketplace.owncloud.com/apps/graphapi) app is installed on ownCloud Server
- The latest client versions are rolled-out to users (required for OpenID Connect support). See the [documentation](https://doc.owncloud.com/server/admin_manual/configuration/user/oidc/#owncloud-desktop-and-mobile-clients) for more information.

See the [documentation]({{< ref "./deployment/bridge" >}}) on how to deploy Infinite Scale in bridge mode.

{{< hint "warning" >}}
**Technology Preview**

ownCloud Infinite Scale is currently in Technology Preview. The bridge mode should only be used in non-production environments.
{{< /hint >}}

### What to expect?

This is the first promoted public release of ownCloud Infinite Scale, released as "Technical Preview". Infinite Scale is not yet ready for production installations. Technical audiences will be able to get a good understanding of the potential of ownCloud's new platform.

Version 1.0.0 comes with the base functionality for sync and share with a much higher performance-, stability- and security-level compared to all available platforms. Based on ten years of experience in enterprise sync and share and a long standing collaboration with the biggest global science organizations this new platform will exceed what content collaboration is today.

### How to get started?

One of the most important objectives for oCIS was to ease the setup of a working instance dramatically. Since oCIS is built with Google's powerful Go language it supports the single-file-deployment: Installing oCIS 1.0.0 is as easy as downloading a single file, applying execution permission to it and get started. No more fiddling around with complicated LAMP stacks.

#### Deployment Options

Given the architecture of oCIS, there are various deployment options based on the users requirements. In our experience setting up the LAMP stack for ownCloud 10 was difficult for many users. Therefore a big emphasis was put on easy yet functional deployment strategies.

{{< tabs "deployments" >}}
{{< tab "Single binary" >}}
#### Delivery as single binary

The single binary is the best option to test the new ownCloud Infinite Scale 1.0.0 Technical Preview release on a local machine. Follow these instructions to get the platform running in the most simple way:

1. Download the binary

    **Linux**
    `curl https://download.owncloud.com/ocis/ocis/1.0.0/ocis-1.0.0-linux-amd64 -o ocis`

    **MacOS**
    `curl https://download.owncloud.com/ocis/ocis/1.0.0/ocis-1.0.0-darwin-amd64 -o ocis`

2. Make it executable

    `chmod +x ocis`

3. Run it

    `./ocis server`

4. Navigate to <https://localhost:9200> and log in to ownCloud Web (admin:admin)

Production environments will need a more sophisticated setup, see <{{< ref "./deployment" >}}> for more information.

{{< /tab >}}
{{< tab "Docker" >}}
#### Containerized Setup

For more sophisticated setups we recommend using one of our docker setup examples. See the [documentation](<{{< ref "./deployment/ocis_traefik" >}}>) for a setup with [Traefik](https://traefik.io/traefik/) as a reverse proxy which also includes automated SSL certificate provisioning using Letsencrypt tools.

{{< /tab >}}
{{< /tabs >}}

### ownCloud Web Features
{{< tabs "web-features" >}}
{{< tab "Framework" >}}
#### Framework
- User avatars (compatible with oC 10 API)
- Alerts for information/errors
- Notifications (bell icon, compatible with oC 10 API)
- Extension points
- Available extensions
  - Media Viewer (images and videos)
  - Draw.io

{{< /tab >}}
{{< tab "Files" >}}
#### Files
- Listing and browsing the hierarchy
- Sorting by columns (name/size/updated)
- Breadcrumb
- Thumbnail previews for images (compatible with oC 10 API and Thumbnails service API)
- Upload (file/folder), using the TUS protocol for reliable uploads
- Download (file)
- Rename
- Copy
- Move
- Delete
- Indicators for resources shared with people (including subfiles and subfolders)
- Indicators for resources shared by link (including subfiles and subfolders)
- Quick actions
  - Add people
  - Create public link on-the-fly and copy it to the clipboard
- Favorites (view + add/remove)
- Shared with me (view)
- Shared with others (view)
- Deleted files
- Versions (list/restore/download/delete)
- File/folder search

{{< /tab >}}
{{< tab "Sharing" >}}
#### Sharing with People (user/group shares)
- Adding people to a resource
  - Adding multiple people at once (users and groups)
  - Autocomplete search to find users
  - Roles: Viewer / Editor (folder) / Advanced permissions (granular permissions)
  - Expiration date
- Listing people who have access to a resource
  - People can be listed when a resource is directly shared and when it's indirectly shared via a parent folder
  - When listing people of an indirectly shared resource, there is a "via" indicator that guides to the directly shared parent
  - Every person can recognize the owner of a resource
  - Every person can recognize their role
  - The owner of a resource can recognize persons that added other people (reshare indicator)
  - Editing persons
  - Removing persons

{{< /tab >}}
{{< tab "Links" >}}
#### Sharing with Links
- Private links (copy)
- Public links
  - Adding public links on single files and folders
    - Roles: Viewer / Editor (folder) / Contributor (folder) / Uploader (folder)
    - Password-protection
    - Expiration date
  - Listing public links
    - Public links can be listed when a resource is directly shared and when it's indirectly shared via a parent folder
    - When listing public links of an indirectly shared resource, there is a "via" indicator that guides to the directly shared parent
    - Copying existing public links
    - Editing existing public links
    - Removing existing public links
  - Viewing public links

{{< /tab >}}
{{< tab "User Profile" >}}
#### User Profile
- Display basic profile information (user name, display name, e-mail, group memberships)
- "Edit" button guides to ownCloud 10 user settings (when used with oC 10)

{{< /tab >}}
{{< tab "User Settings" >}}

##### Basic user settings
- Language of the web interface

{{< /tab >}}
{{< /tabs >}}

### oCIS Backend Features

{{< tabs "backend-features" >}}
{{< tab "Storage" >}}

#### Storage

The default oCIS storage driver deconstructs a filesystem to be able to efficiently look up files by fileid as well as path. It stores all folders and files by a uuid and persists share and other metadata using extended attributes. This allows using the linux VFS cache using stat syscalls instead of a database or key/value store. The driver implements trash, versions and sharing. It not only serves as the current default storage driver, but also as a blueprint for future storage driver implementations.

{{< /tab >}}
{{< tab "IDM" >}}
#### User and group management
- Functionality available via API and frontend ("Accounts" extension)
- User listing (API/FE)
- User creation (API/FE)
- User deletion (API/FE)
- User activation/blocking (API/FE)
- Role assignment for users (API/FE)
- User editing (API)
- Multi-select in the frontend (delete & block/activate)
- Group creation (API)
- Add/remove users to/from groups (API)
- Group deletion (API)
- Create/read/update/delete users and groups (CLI)

{{< /tab >}}
{{< tab "Settings" >}}

##### Settings

The settings service provides APIs for other services for registering a set of settings as `Bundle`. It also provides a pluggable extension for ownCloud Web which provides dynamically built web forms, so that users can customize their own settings. Some well known settings are directly used by ownCloud Web for adapted user experience, e.g. the UI language. Services can query the users' chosen settings for customized backend and frontend operations as needed.

##### Roles & Permissions System

Infinite Scale follows a role-based access control model. Based on permissions for actions which are provided by the system and by extensions, roles can be composed. Ultimately, these roles can be assigned to users to define what users are permitted to do. This model allows a segregation of duties for administration and allows granular control of how different types of users (e.g., Guests) can use the platform.

- Currently available permissions: Manage accounts (gives access to the internal user management), manage roles (allows assigning roles to users)
- The current roles are exemplary default roles which are used for demonstration purposes
  - "Admin": Has the permissions to "manage accounts" and to "manage roles"
  - "User": Does not have any dedicated permission
  - "Guest": Does not have any dedicated permission
- Currently a user can only have one role
- Users with the role "Admin" can assign/unassign roles to/from other users (as part of the permission to "manage roles")

{{< /tab >}}
{{< tab "APIs" >}}
#### APIs

- WebDAV
- OCS

{{< /tab >}}
{{< /tabs >}}

### Known issues

- There are feature differences depending on the operation mode, e.g., no user management with ownCloud Web and oC 10 backend
- Public links do not yet respect the given role (a recipient has full permissions no matter which role has been set)
- Resharing does not yet work as expected
  - Share recipients can create public links with higher permissions than they originally had
  - Share recipients can add other people but they will not be able to access the data
- Sharing indicators in the file list will only be shown after opening the right sidebar for a resource
- Users can't change their password yet
- Folder sizes will not be calculated
- Cleanups are not yet available (e.g., shares of a deleted user will not be removed)
- Sharing from the desktop client does not work yet
- There are no notifications yet
- There can be issues with access tokens not being refreshed correctly, leading to interruptions, e.g., during uploads
- Deleting non-empty folders from the trash bin does not work
- Emptying the whole trash bin does not work

For feedback and bug reports, please use the [public issue tracker](https://github.com/owncloud/ocis/issues).
