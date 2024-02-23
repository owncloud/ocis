---
title: "26. Application based user settings"
date: 2024-02-09T17:30:00+01:00
weight: 26
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/adr
geekdocFilePath: 0026-application-based-user-settings.md
---

* Status: draft
* Deciders: [@butonic](https://github.com/butonic), [@fschade](https://github.com/fschade), [@kulmann](https://github.com/kulmann)
* Date: 2024-02-09

## Context and Problem Statement

To share user settings across devices applications want to store user specific settings on the server. The ePUB app wants to remember which page the user is on. The iOS app wants to rember search queries. The Caldav app needs a space to store data.

## Decision Drivers <!-- optional -->

## Considered Options

* OCS provisioning API
* settings service
* libregraph API

## Decision Outcome

Chosen option: *???*

### Positive Consequences:

* TODO

### Negative Consequences:

* TODO

## Pros and Cons of the Options <!-- optional -->

### OCS provisioning API

Nextcloud added a `/ocs/v2.php/apps/provisioning_api/api/v1/config/users/{appId}/{configKey}` endpoint

* Bad, legacy API we want to get rid of

### settings service

- Bad, yet another API. Always uses POST requests.

### libregraph API

The MS Graph API has [a special approot driveItem](https://learn.microsoft.com/en-us/graph/api/drive-get-specialfolder?view=graph-rest-1.0&tabs=http) that apps can use to store arbitrary files. See also: 
[Using an App Folder to store user content without access to all files](https://learn.microsoft.com/en-us/onedrive/developer/rest-api/concepts/special-folders-appfolder?view=odsp-graph-online) and a blog post with the section [Store data in the application’s personal folder](https://blog.mastykarz.nl/easiest-store-user-settings-microsoft-365-app/#store-data-in-the-applications-personal-folder).

It basically uses the `/me/drive/special/approot:/{filename}` endpoint to
```http
PUT https://graph.microsoft.com/v1.0/me/drive/special/approot:/settings.json:/content
content-type: text/plain
authorization: Bearer abc

{"key": "value"}
```
or
```http
GET https://graph.microsoft.com/v1.0/me/drive/special/approot:/settings.json:/content
authorization: Bearer abc
```

On single page apps you need two requests:
```http
GET https://graph.microsoft.com/v1.0/me/drive/special/approot:/settings.json?select=@microsoft.graph.downloadUrl
authorization: Bearer abc
```
followed by
```http
GET <url from the response['@microsoft.graph.downloadUrl'] property>
```

Currently, applications have no dedicated tokens that we could use to derive the `appid` from. All apps should have an `appid` and [be discoverable under](https://learn.microsoft.com/en-us/graph/api/application-list?view=graph-rest-1.0&tabs=http)
```http
GET /applications
```

In any case for libregraph we could introduce a `LIBRE_GRAPH_APPID` header to make these requests possible rather soon.

Then we can decide if we want to store these files in the users personal drive, or if we create a space for every app that then uses the userid as a folder that contains all the files for the user.

- Good, because clients can remain in libregraph API land
- Bad, we currently have no application tokens



## Links <!-- optional -->
