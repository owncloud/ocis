# 5. Graph API for Spaces

* Status: proposed
* Deciders: @butonic, @micbar, @dragotin, @hodyroff, @pmaier1
* Date: 2021-03-19

Technical Story: API to enable the concept of [Spaces](https://github.com/owncloud/enterprise/issues/3863)

## Context and Problem Statement

As one of the building blocks for Spaces in oCIS we plan to add an API that returns information about available spaces. This ADR discusses the API design oriented on the Microsoft Graph API.

> Note: The term "spaces" is used here in the context of "a space where files can be saved", similar to a directory. It is not to be confused with space in the sense of free file space for example.

The purpose of this new API is to give clients a very simple way to query the dynamic list of spaces, that the user has access to. Clients can provide a better user experience with that.

This API is supposed to be queried often, to give clients a condensed view of the available spaces for a user, but also their eTags and cTags. Hence the clients do not have to perform a PROPFIND for every space separately.

This API would even allow to provide (WebDAV-) endpoints depending on the kind and version of the client asking for it.

## Decision Drivers

- Make it easy to work with a dynamic list of spaces of a user for the clients.
- No longer the need to make assumptions about WebDAV- and other routes in clients.
- More meta data available about spaces for a better user experience.
- Part of the bigger spaces plan.
- Important to consider in client migration scenarios, ie. in CERN.

## Considered Options

1. [Microsoft Graph API](https://developer.microsoft.com/en-us/graph) inspired API that provides the requested information.

## Decision Outcome

This the DRAFT for the API.

### API to Get Info about Spaces

ownCloud servers provide an API to query for available spaces of an user.

Most important, the API returns the WebDAV endpoint for each space. With that, clients do not have to make assumptions about WebDAV routes any more.

See [Drive item in Microsoft Graph API](https://docs.microsoft.com/en-us/graph/api/resources/onedrive?view=graph-rest-1.0) for an overview of `drive` and `driveItem` resources. The concrete list of drives / spaces a user has access to can be obtained on multiple endpoints.

### Get "Home folder"

Retrieve information about the home space of a user. Note: The user has access to more spaces. This call only returns the home space to provide API parity with the Graph API.

API Call: `/me/drive`: Returns the information about the users home folder.

### Get All Spaces of a User

Retrieve a list of available spaces of a user. This includes all spaces the user has access to at that moment, also the home space.

API Call: `/me/drives`: Returns a list of spaces.

There is also `/drives`, returning the list of spaces the user has access to. This endpoint is used to access any space by id using `/drives/{drive-id}`.

### Common Reply

The reply to both calls is either one or a list of [Drive representation objects](https://docs.microsoft.com/de-de/graph/api/resources/drive?view=graph-rest-1.0):

```
{
  "id": "string",
  "createdDateTime": "string (timestamp)",
  "description": "string",
  "driveType": "personal | projectSpaces | shares",
  "oCDriveStatus": "accepted | pending | mandatory",
  "eTag": "string",
  "lastModifiedDateTime": "string (timestamp)",
  "name": "string",
  "owner": { "@odata.type": "microsoft.graph.identitySet" },
  "quota": { "@odata.type": "microsoft.graph.quota" },
  "root":  { "@odata.type": "microsoft.graph.driveItem" },
  "webUrl": "url",
  "oCCoOwner": [ { "@odata.type": "microsoft.graph.identitySet" } ]
}

```

The meaning of the object are in ownClouds context:

1. **id** - a unique ID identifying the space
2. **driveType** - describing the type of the space.
3. **oCDriveStatus** - telling the status (*)
4. **eTag** - the current ETag of the space
5. **owner** - an owner object to whom the space belongs
6. **quota** - quota information about this space
7. **root**  - the root driveItem object.
6. **webUrl** - The URL to make this space visible in the browser.
7. **oCCoOwner** - optional array owner objects of the co-owners of a space (*)

> Note: the **root** object is a [driveItem](https://docs.microsoft.com/de-de/graph/api/resources/driveitem?view=graph-rest-1.0) and contains information about the root of the space. Especially important is the member `webDavUrl` in the root object which provides the full URL to be used by clients to address items within the space. Also it contains the fields `eTag` and `cTag`.

The following *driveType* values are available in the first step, but might be enhanced later:

* **personal**: The users home space
* **projectSpaces**: The project spaces available for the user (*)
* **shares**: The share jail, contains all shares for the user (*)

Other space types such as backup, hidden etc. can be added later as requested.

The (*) marked types are not defined in the official MS API. They are prefixed with `oC` to avoid namespace clashes.


The following *driveStatus* values are available:

* **accepted**: The user has accepted the space and uses it
* **pending**: The user has not yet accepted the space , but can use it after having it accepted.
* **mandatory**: This is an mandatory space. Used for the personal- and shares-space. The user can not influence if it is visible or not, it is always available.
* **offline**: The space is currently not available.

### Positive Consequences

- A well understood and mature API from Microsoft adopted to our needs.
- Prerequisite for Spaces in oCIS.
- Enables further steps in client development.

### Negative Consequences

- Migration impact on existing installations. Still to be investigated.
- Requires additional webdav endpoint that allows accessing an arbitrary storage space, either
  - with an id: `/dav/spaces/<spaceid>/relative/path/to/file.ext`, or
  - with a global path: `/dav/global/<accessible>/<mount>/<point>/relative/path/to/file.ext`, e.g. `/dav/global/projects/Golive 2021/Resources/slides.odt`

### Open Topics

- What are the WebDAV pathes for Trashbin, Versions
    + option: additional entries in the reply struct
- The identitySet object used for "owner" and "coowner" require to implement the [https://docs.microsoft.com/de-de/graph/api/resources/identityset?view=graph-rest-1.0](IdentitySet) JSON object, which contains information that seems to be of limited benefit for oCIS. An alternative would be to implement a simpler identity object for oCIS and use that.
