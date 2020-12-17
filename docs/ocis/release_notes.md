---
title: "Release Notes"
date: 2020-12-16T20:35:00+01:00
weight: 0
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis
geekdocFilePath: release_notes.md
---
# Release Notes

## ownCloud Infinite Scale 1.0.0 Technology Preview - Initial Release Notes

We are pleased to annouce the availability of ownCloud Infinite Scale 1.0.0 Technology Preview which is released as the first public version of the new Infinite Scale platform.

### Microservice architecture

ownCloud Infinite Scale is following the microservices architectual pattern. It is implemented as a set of microservices which are independent of each other. They are coupled with very well-defined APIs and communicate via HTTP. This architecture fosters a lot of benefits that we were going for with the new design for oCIS:

- Independent services: Every service is independant, comparably small and brings it's own webserver, backend/APIs and frontend components.
- Each service runs as a separate service on the system, increasing security and stability
- Scalability:  High performance demands can be fulfilled by scaling the amount of services
- Testability: Each service can be tested on its own due to well-defined APIs and functionality
- Protocol-driven development
- High-performance communication between services through technologies like GRPC
- Multi-platform support through utilizing Golang - only minimal dependency on platform packages.
- Cloud-native deployment and update strategies

### Key figures

- The all-new ownCloud Web frontend ships with the platform
- OpenID Connect is the technology choice for authentication
- An Identity Provider is bundled to ease deployment and operations. It can be replaced with other applications if desired.
- Up-to-date, cloud-native deployment options are available
- Flexible configuration through environment variables, yaml files or command-line switches
- Database-less architecture - metadata and data are kept together in the storage as a single source of truth
- Native storage capabilities are used
- Public ownCloud APIs like WebDAV and OCS have been kept compatible to ownCloud 10
- A secure and flexible framework to create extensions for ownCloud. It allows to integrate with ownCloud data in a very easy yet powerful way.

#### Supported platforms
- Linux-amd64
- Darwin-amd64
- Experimental: Windows, ARM (e.g., Raspberry Pi)

#### Client support
All official ownCloud Clients support the Infinite Scale server with the following versions:
- Desktop >= 2.7.0
- Android >= 2.15
- iOS >= 1.2

### Architecture

ownCloud Infinite Scale is built as a modular framework in which components can be scaled individually. It consists of

- a user management service
- a storage backend service
- Built-in IdP
- Frontend
- Application gateway/proxy

These components can be deployed in a multi-tier deployment architecture. See the [documentation](https://owncloud.github.io/ocis/) for an overview of the services.

### Operation modes
#### Standalone Full Stack Server mode (with oCIS storage driver)

#### Standalone Single service mode for scaleouts

#### Bridge mode with ownCloud 10 backend

For the product transition phase, ownCloud Infinite Scale comes with an operation mode ("bridge mode") that allows to create a hybrid deployment between both server generations to operate the new web frontend with ownCloud 10 and Infinite Scale in parallel. This setup allows to operate the ownCloud Web frontend with both server generations and provides the foundation to migrate users gradually to the new backend.

**Requirements for the bridge mode**
- ownCloud Server >= 10.6
- [Open ID Connect](https://marketplace.owncloud.com/apps/openidconnect) is used for user authentication
- The [Graph API](https://marketplace.owncloud.com/apps/graphapi) app is installed on ownCloud Server
- The latest client versions are rolled-out to users (required for OpenID Connect support). See the [ownCloud Documentation](https://doc.owncloud.com/server/admin_manual/configuration/user/oidc/#owncloud-desktop-and-mobile-clients) for more information.

{{< hint [warning] >}}
**ownCloud Infinite Scale is currently in Technology Preview. The bridge mode should only be used in non-productive environments.**
{{< /hint >}}

https://owncloud.github.io/ocis/deployment/owncloud10_with_oc_web/

[To illustrate, a little graphic that describes the various operation modes would be cool

### What to expect?

This is the first promoted public release of ownCloud Infinite Scale, released as "Technical Preview". Infinite Scale is not yet ready for production installations. Technical audience will get a good impression of the potential of ownClouds new platform.

Version 1.0.0 comes with the base functionality for sync and share on a much higher performance-, stability- and security-level compared to all available platforms. Based on ten years of experience in enterprise sync and share and a long standing collaboration with the biggest global science organizations this new platform will exceed what content collaboration is today.

TODO: Mention the base modules of oCIS

### How to get started?

One of the most important objectives for oCIS was to ease the setup of a working instance dramatically. Since oCIS is built on Google's powerful GO language it supports the single-file-deployment: Installing oCIS 1.0.0 is as easy as downloading a single file, applying execution permission to it and get started. No more fiddling around with complicated LAMP stacks.

#### Deployment Options

Given the Golang-based architecture of oCIS, there are various deployment options based on the users requirements. With our experience with the for many users difficult setup of the LAMP stack before a big emphasis was put on easy yet functional deployment strategies.

##### Delivery as single binary

The single binary is the best option to test the new ownCloud Infinite Scale 1.0.0 Technical Preview release on a local machine. Follow these instructions to get the platform running in the most simple way:

1. Download the binary

**Linux**

`curl https://download.owncloud.com/ocis/ocis/testing/ocis-testing-linux-amd64 --output ocis`

**MacOS**

`curl https://download.owncloud.com/ocis/ocis/testing/ocis-testing-darwin-amd64 --output ocis`

2. Make it executable
`chmod +x ocis`

3. Run it
`./ocis server`

4. Navigate to http://localhost:9200 and log in to ownCloud Web (admin/admin)

Infinite Scale environments on remote machines should use a more sophisticated setup. See the [documentation](https://owncloud.github.io/ocis/deployment/) for more information.

##### Containerized Setup

For more sophisticated and production setups we recommend to use one of our proposed docker setups. See the [documentation](https://owncloud.github.io/ocis/deployment/ocis_traefik/) for a setup with Traefik as a reverse proxy which also includes automated SSL certificate provisioning using Letsencrypt tools.

### ownCloud Web Features
#### Framework
- User avatars (compatible with oC 10 API)
- Alerts for information/errors
- Notifications (bell icon, compatible with oC 10 API)
- Extension points
- Available extensions
  - Media Viewer (images and videos)
  - Draw.io

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
- Quick actions (extension point)
  - Add people
  - Create public link on-the-fly and copy it to the clipboard
- Favorites (view + add/remove)
- Shared with me (view)
- Shared with others (view)
- Deleted files
- Versions (list/restore/download/delete)
- File/folder search

#### Sharing with People (user/group shares)
- Adding people to a resource
  - Adding multiple people at once (users and groups)
  - Autocomplete search to find users
  - Roles: Viewer / Editor (folder) / Advanced permissions (granular permissions)
  - Expiration date
- Listing people who have access to a resource
  - People can be listed when a resource is directly shared and when it's indirectly shared via a parent folder.
  - When listing people of an indirectly shared resource, there is a "via" indicator that guides to the directly shared parent.
  - Every person can recognize the owner of a resource.
  - Every person can recognize their role.
  - The owner of a resource can recognize persons that added other people (reshare indicator).
  - Editing persons
  - Removing persons

#### Sharing with Links
- Private links (copy)
- Public links
  - Adding public links on single files and folders
    - Roles: Viewer / Editor (folder) / Contributor (folder) / Uploader (folder)
    - Password-protection
    - Expiration date
  - Listing public links
    - Public links can be listed when a resource is directly shared and when it's indirectly shared via a parent folder.
    - When listing public links of an indirectly shared resource, there is a "via" indicator that guides to the directly shared parent.
    - Copying existing public links
    - Editing existing public links
    - Removing existing public links
  - Viewing public links

#### User Profile
- Display basic profile information (user name, display name, e-mail, group memberships)
- "Edit" button guides to ownCloud 10 user settings (when used with oC 10)

#### Basic user settings
- Language of the web interface

### oCIS Backend Features

#### Storage

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
- Create/read/update users and groups (CLI)

#### Settings
##### Settings bundles framework
- What is a settings bundle?
- What can you do with it?
- Extensions?

##### Roles & Permissions System
Infinite Scale follows a role-based access control model. Based on permissions for actions which are provided by the system and by extensions, roles can be composed. Ultimately, these roles can be assigned to users to define what users are permitted to do. This model allows to realize a segregation of duties for administration and allows to control granularly how different types of users (e.g., Guests) can use the platform.

- Currently available permissions: Manage accounts (gives access to the internal user management)
- The current roles are exemplary default roles which defined in config files
  - "Admin": Has the permission to "manage accounts"
  - "User": Does not have any dedicated permission
  - "Guest": Does not have any dedicated permission
- Currently a user can only have one role
- Users with the role "Admin" can assign/unassign roles to/from other users (as part of the permission to "manage accounts")

#### APIs
- WebDAV
- OCS

### Known issues for OCIS standalone
- There are feature differences depending on the operation  mode, e.g., no user management with ownCloud Web and oC 10 backend
- Public links do not yet respect the given role (a recipient has full permissions no matter which role has been set)
- Resharing works but does not have necessary restrictions in place
  - Share recipients can add more people or create public links with higher permissions than they originally had
  - Every person in a share can see all other people in the people list
- Sharing indicators in the file list will only be shown after opening the right sidebar for a resource
- Users can't change their password yet
- Folder sizes will not be calculated
- Cleanups are not yet available (e.g., shares of a deleted user will not be removed)
