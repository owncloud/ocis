---
title: "LibreGraph"
date: 2018-05-02T00:00:00+00:00
weight: 20
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/apis/http/graph
geekdocFilePath: _index.md
geekdocCollapseSection: true
---

The LibreGraph API is a REST Api which is inspired by the [Microsoft Graph API](https://developer.microsoft.com/en-us/graph). It tries to stay compliant with the Microsoft Graph API and aims to be the Next Generation Api in Infinite Scale where we want to support most of the features of the platform.
The [API specification](https://github.com/owncloud/libre-graph-api) is available in the OpenApi 3 standard and there are generated client and server [SDKs](https://github.com/owncloud/libre-graph-api#clients) available. You can browse the API with the [Swagger UI](https://owncloud.dev/libre-graph-api/).

## Calling the LibreGraph API

```sh
{HTTP method} https://ocis.url/graph/{version}/{resource}?{query-parameters}
```

The request component consists of:

| Component          | Description                                                             |
|--------------------|-------------------------------------------------------------------------|
| {HTTP method}      | The HTTP method which is used in the request.                           |
| {version}          | The version of the LibreGraph API used by the client.                   |
| {resource}         | The LibreGraph Resource which the client is referencing in the request. |
| {query-parameters} | Optional parameters for the request to customize the response.          |

### HTTP methods

| Method | Description                   |
|--------|-------------------------------|
| GET    | Read data from a resource.    |
| POST   | Create a new resource.        |
| PATCH  | Update an existing resource.  |
| PUT    | Replace an existing resource. |
| DELETE | Delete an existing resource.  |

The methods `GET` and `DELETE` need no request body. The methods `POST`, `PATCH` and `PUT` require a request body, normally in JSON format to provide the needed values.

### Version

Infinite Scale currently provides the version `v1.0`.

### Resource

A resource could be an entity or a complex type and is usually defined by properties. Entities are always recognizable by an `Id` property. The URL contains the resource which you are interacting with e.g. `/me/drives` or `/groups/{group-id}`.

Each resource could possibly require different permissions. Usually you need permissions on a higher level for creating or updating an existing resource than for reading.

### Query parameters

Query parameters can be OData system query options, or other strings that a method accepts to customize its response.

You can use optional OData system query options to include more or fewer properties than the default response, filter the response for items that match a custom query, or provide additional parameters for a method.

For example, adding the following filter parameter restricts the drives returned to only those with the driveType property of `project`.

```shell
GET https://ocis.url/graph/v1.0/drives?$filter=driveType eq 'project'
```
For more information about OData query options please check the [API specification](https://github.com/owncloud/libre-graph-api) and the provided examples.

## Resources

{{< toc-tree >}}
