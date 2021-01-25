---
title: "Release Notes"
date: 2020-12-16T20:35:00+01:00
weight: 0
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis
geekdocFilePath: release_notes.md
---

## ownCloud Infinite Scale 1.1.0 Technology Preview

Version 1.1.0 is a hardening and patch release. It ships with the latest version of ownCloud Web and brings a couple of minor improvements. The minor version increase is needed due to non-backwards compatible changes in configuration. The documentation has been updated to reflect the changes. Please note that this version is still a Technology Preview and not suited for production use.

The most prominent changes in version 1.1.0 comprise
- Performance and stability improvements for installations with multiple concurrent users
- Simplified configuration by introducing the new environment variable OCIS_URL
- Beta release of [ownCloud performance scripts](https://github.com/owncloud/cdperf)
- Update ownCloud web to [v1.0.1](https://github.com/owncloud/web/releases/tag/v1.0.1)
- Update reva to [v1.5.1](https://github.com/cs3org/reva/releases/tag/v1.5.1)

You can also read the full [ownCloud Infinite Scale changelog](https://github.com/owncloud/ocis/blob/master/CHANGELOG.md) for further details on what has changed.

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

These components can be deployed in a multi-tier deployment architecture. See the [documentation](https://owncloud.github.io/ocis/) for an overview of the services.

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

See the [documentation](https://owncloud.github.io/ocis/deployment/owncloud10_with_oc_web/) on how to deploy Infinite Scale in bridge mode.

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

Production environments will need a more sophisticated setup, see <https://owncloud.github.io/ocis/deployment/> for more information.

{{< /tab >}}
{{< tab "Docker" >}}
#### Containerized Setup

For more sophisticated setups we recommend using one of our docker setup examples. See the [documentation](https://owncloud.github.io/ocis/deployment/ocis_traefik/) for a setup with [Traefik](https://traefik.io/traefik/) as a reverse proxy which also includes automated SSL certificate provisioning using Letsencrypt tools.

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
