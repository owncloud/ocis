# DrivesPermissionsApi

All URIs are relative to *https://ocis.ocis.rolling.owncloud.works/graph*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**createLink**](DrivesPermissionsApi.md#createlink) | **POST** /v1beta1/drives/{drive-id}/items/{item-id}/createLink | Create a sharing link for a DriveItem |
| [**deletePermission**](DrivesPermissionsApi.md#deletepermission) | **DELETE** /v1beta1/drives/{drive-id}/items/{item-id}/permissions/{perm-id} | Remove access to a DriveItem |
| [**getPermission**](DrivesPermissionsApi.md#getpermission) | **GET** /v1beta1/drives/{drive-id}/items/{item-id}/permissions/{perm-id} | Get sharing permission for a file or folder |
| [**invite**](DrivesPermissionsApi.md#invite) | **POST** /v1beta1/drives/{drive-id}/items/{item-id}/invite | Send a sharing invitation |
| [**listPermissions**](DrivesPermissionsApi.md#listpermissions) | **GET** /v1beta1/drives/{drive-id}/items/{item-id}/permissions | List the effective sharing permissions on a driveItem. |
| [**setPermissionPassword**](DrivesPermissionsApi.md#setpermissionpassword) | **POST** /v1beta1/drives/{drive-id}/items/{item-id}/permissions/{perm-id}/setPassword | Set sharing link password |
| [**updatePermission**](DrivesPermissionsApi.md#updatepermission) | **PATCH** /v1beta1/drives/{drive-id}/items/{item-id}/permissions/{perm-id} | Update sharing permission |



## createLink

> Permission createLink(driveId, itemId, driveItemCreateLink)

Create a sharing link for a DriveItem

You can use the createLink action to share a driveItem via a sharing link.  The response will be a permission object with the link facet containing the created link details.  ## Link types  For now, The following values are allowed for the type parameter.  | Value          | Display name      | Description                                                     | | -------------- | ----------------- | --------------------------------------------------------------- | | view           | View              | Creates a read-only link to the driveItem.                      | | upload         | Upload            | Creates a read-write link to the folder driveItem.              | | edit           | Edit              | Creates a read-write link to the driveItem.                     | | createOnly     | File Drop         | Creates an upload-only link to the folder driveItem.            | | blocksDownload | Secure View       | Creates a read-only link that blocks download to the driveItem. | 

### Example

```ts
import {
  Configuration,
  DrivesPermissionsApi,
} from '';
import type { CreateLinkRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // To configure HTTP basic authorization: basicAuth
    username: "YOUR USERNAME",
    password: "YOUR PASSWORD",
  });
  const api = new DrivesPermissionsApi(config);

  const body = {
    // string | key: id of drive
    driveId: driveId_example,
    // string | key: id of item
    itemId: itemId_example,
    // DriveItemCreateLink | In the request body, provide a JSON object with the following parameters. (optional)
    driveItemCreateLink: {"type":"view"},
  } satisfies CreateLinkRequest;

  try {
    const data = await api.createLink(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **driveId** | `string` | key: id of drive | [Defaults to `undefined`] |
| **itemId** | `string` | key: id of item | [Defaults to `undefined`] |
| **driveItemCreateLink** | [DriveItemCreateLink](DriveItemCreateLink.md) | In the request body, provide a JSON object with the following parameters. | [Optional] |

### Return type

[**Permission**](Permission.md)

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Response |  -  |
| **207** | Partial success response TODO |  -  |
| **0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## deletePermission

> deletePermission(driveId, itemId, permId)

Remove access to a DriveItem

Remove access to a DriveItem.  Only sharing permissions that are not inherited can be deleted. The &#x60;inheritedFrom&#x60; property must be &#x60;null&#x60;. 

### Example

```ts
import {
  Configuration,
  DrivesPermissionsApi,
} from '';
import type { DeletePermissionRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // To configure HTTP basic authorization: basicAuth
    username: "YOUR USERNAME",
    password: "YOUR PASSWORD",
  });
  const api = new DrivesPermissionsApi(config);

  const body = {
    // string | key: id of drive
    driveId: driveId_example,
    // string | key: id of item
    itemId: itemId_example,
    // string | key: id of permission
    permId: permId_example,
  } satisfies DeletePermissionRequest;

  try {
    const data = await api.deletePermission(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **driveId** | `string` | key: id of drive | [Defaults to `undefined`] |
| **itemId** | `string` | key: id of item | [Defaults to `undefined`] |
| **permId** | `string` | key: id of permission | [Defaults to `undefined`] |

### Return type

`void` (Empty response body)

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **204** | Success |  -  |
| **0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## getPermission

> Permission getPermission(driveId, itemId, permId)

Get sharing permission for a file or folder

Return the effective sharing permission for a particular permission resource. 

### Example

```ts
import {
  Configuration,
  DrivesPermissionsApi,
} from '';
import type { GetPermissionRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // To configure HTTP basic authorization: basicAuth
    username: "YOUR USERNAME",
    password: "YOUR PASSWORD",
  });
  const api = new DrivesPermissionsApi(config);

  const body = {
    // string | key: id of drive
    driveId: driveId_example,
    // string | key: id of item
    itemId: itemId_example,
    // string | key: id of permission
    permId: permId_example,
  } satisfies GetPermissionRequest;

  try {
    const data = await api.getPermission(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **driveId** | `string` | key: id of drive | [Defaults to `undefined`] |
| **itemId** | `string` | key: id of item | [Defaults to `undefined`] |
| **permId** | `string` | key: id of permission | [Defaults to `undefined`] |

### Return type

[**Permission**](Permission.md)

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Retrieved resource |  -  |
| **0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## invite

> CollectionOfPermissions invite(driveId, itemId, driveItemInvite)

Send a sharing invitation

Sends a sharing invitation for a &#x60;driveItem&#x60;. A sharing invitation provides permissions to the recipients and optionally sends them an email with a sharing link.  The response will be a permission object with the grantedToV2 property containing the created grant details.  ## Roles property values For now, roles are only identified by a uuid. There are no hardcoded aliases like &#x60;read&#x60; or &#x60;write&#x60; because role actions can be completely customized. 

### Example

```ts
import {
  Configuration,
  DrivesPermissionsApi,
} from '';
import type { InviteRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // To configure HTTP basic authorization: basicAuth
    username: "YOUR USERNAME",
    password: "YOUR PASSWORD",
  });
  const api = new DrivesPermissionsApi(config);

  const body = {
    // string | key: id of drive
    driveId: driveId_example,
    // string | key: id of item
    itemId: itemId_example,
    // DriveItemInvite | In the request body, provide a JSON object with the following parameters. To create a custom role submit a list of actions instead of roles. (optional)
    driveItemInvite: {"recipients":[{"@libre.graph.recipient.type":"user","objectId":"4c510ada-c86b-4815-8820-42cdf82c3d51"}],"roles":["b1e2218d-eef8-4d4c-b82d-0f1a1b48f3b5"]},
  } satisfies InviteRequest;

  try {
    const data = await api.invite(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **driveId** | `string` | key: id of drive | [Defaults to `undefined`] |
| **itemId** | `string` | key: id of item | [Defaults to `undefined`] |
| **driveItemInvite** | [DriveItemInvite](DriveItemInvite.md) | In the request body, provide a JSON object with the following parameters. To create a custom role submit a list of actions instead of roles. | [Optional] |

### Return type

[**CollectionOfPermissions**](CollectionOfPermissions.md)

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Response |  -  |
| **207** | Partial success response TODO |  -  |
| **400** | Bad request |  -  |
| **0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## listPermissions

> CollectionOfPermissionsWithAllowedValues listPermissions(driveId, itemId, $filter, $select)

List the effective sharing permissions on a driveItem.

The permissions collection includes potentially sensitive information and may not be available for every caller.  * For the owner of the item, all sharing permissions will be returned. This includes co-owners. * For a non-owner caller, only the sharing permissions that apply to the caller are returned. * Sharing permission properties that contain secrets (e.g. &#x60;webUrl&#x60;) are only returned for callers that are able to create the sharing permission.  All permission objects have an &#x60;id&#x60;. A permission representing * a link has the &#x60;link&#x60; facet filled with details. * a share has the &#x60;roles&#x60; property set and the &#x60;grantedToV2&#x60; property filled with the grant recipient details. 

### Example

```ts
import {
  Configuration,
  DrivesPermissionsApi,
} from '';
import type { ListPermissionsRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // To configure HTTP basic authorization: basicAuth
    username: "YOUR USERNAME",
    password: "YOUR PASSWORD",
  });
  const api = new DrivesPermissionsApi(config);

  const body = {
    // string | key: id of drive
    driveId: driveId_example,
    // string | key: id of item
    itemId: itemId_example,
    // string | Filter items by property values. By default all permissions are returned and the avalable sharing roles are limited to normal users. To get a list of sharing roles applicable to federated users use the example $select query and combine it with $filter to omit the list of permissions. (optional)
    $filter: @libre.graph.permissions.roles.allowedValues/rolePermissions/any(p:contains(p/condition, '@Subject.UserType=="Federated"')),
    // Set<'@libre.graph.permissions.actions.allowedValues' | '@libre.graph.permissions.roles.allowedValues' | 'value'> | Select properties to be returned. By default all properties are returned. Select the roles property to fetch the available sharing roles without resolving all the permissions. Combine this with the $filter parameter to fetch the actions applicable to federated users. (optional)
    $select: ...,
  } satisfies ListPermissionsRequest;

  try {
    const data = await api.listPermissions(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **driveId** | `string` | key: id of drive | [Defaults to `undefined`] |
| **itemId** | `string` | key: id of item | [Defaults to `undefined`] |
| **$filter** | `string` | Filter items by property values. By default all permissions are returned and the avalable sharing roles are limited to normal users. To get a list of sharing roles applicable to federated users use the example $select query and combine it with $filter to omit the list of permissions. | [Optional] [Defaults to `undefined`] |
| **$select** | `@libre.graph.permissions.actions.allowedValues`, `@libre.graph.permissions.roles.allowedValues`, `value` | Select properties to be returned. By default all properties are returned. Select the roles property to fetch the available sharing roles without resolving all the permissions. Combine this with the $filter parameter to fetch the actions applicable to federated users. | [Optional] [Enum: @libre.graph.permissions.actions.allowedValues, @libre.graph.permissions.roles.allowedValues, value] |

### Return type

[**CollectionOfPermissionsWithAllowedValues**](CollectionOfPermissionsWithAllowedValues.md)

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Retrieved resource |  -  |
| **0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## setPermissionPassword

> Permission setPermissionPassword(driveId, itemId, permId, sharingLinkPassword)

Set sharing link password

Set the password of a sharing permission.  Only the &#x60;password&#x60; property can be modified this way. 

### Example

```ts
import {
  Configuration,
  DrivesPermissionsApi,
} from '';
import type { SetPermissionPasswordRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // To configure HTTP basic authorization: basicAuth
    username: "YOUR USERNAME",
    password: "YOUR PASSWORD",
  });
  const api = new DrivesPermissionsApi(config);

  const body = {
    // string | key: id of drive
    driveId: driveId_example,
    // string | key: id of item
    itemId: itemId_example,
    // string | key: id of permission
    permId: permId_example,
    // SharingLinkPassword | New password value
    sharingLinkPassword: {"password":"TestPassword123!"},
  } satisfies SetPermissionPasswordRequest;

  try {
    const data = await api.setPermissionPassword(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **driveId** | `string` | key: id of drive | [Defaults to `undefined`] |
| **itemId** | `string` | key: id of item | [Defaults to `undefined`] |
| **permId** | `string` | key: id of permission | [Defaults to `undefined`] |
| **sharingLinkPassword** | [SharingLinkPassword](SharingLinkPassword.md) | New password value | |

### Return type

[**Permission**](Permission.md)

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Updated permission |  -  |
| **0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## updatePermission

> Permission updatePermission(driveId, itemId, permId, permission)

Update sharing permission

Update the properties of a sharing permission by patching the permission resource.  Only the &#x60;roles&#x60;, &#x60;expirationDateTime&#x60; and &#x60;password&#x60; properties can be modified this way. 

### Example

```ts
import {
  Configuration,
  DrivesPermissionsApi,
} from '';
import type { UpdatePermissionRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // To configure HTTP basic authorization: basicAuth
    username: "YOUR USERNAME",
    password: "YOUR PASSWORD",
  });
  const api = new DrivesPermissionsApi(config);

  const body = {
    // string | key: id of drive
    driveId: driveId_example,
    // string | key: id of item
    itemId: itemId_example,
    // string | key: id of permission
    permId: permId_example,
    // Permission | New property values
    permission: {"link":{"type":"edit"}},
  } satisfies UpdatePermissionRequest;

  try {
    const data = await api.updatePermission(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **driveId** | `string` | key: id of drive | [Defaults to `undefined`] |
| **itemId** | `string` | key: id of item | [Defaults to `undefined`] |
| **permId** | `string` | key: id of permission | [Defaults to `undefined`] |
| **permission** | [Permission](Permission.md) | New property values | |

### Return type

[**Permission**](Permission.md)

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Updated permission |  -  |
| **0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

