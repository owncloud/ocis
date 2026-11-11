# DrivesRootApi

All URIs are relative to *https://ocis.ocis.rolling.owncloud.works/graph*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**createDriveItem**](DrivesRootApi.md#createdriveitem) | **POST** /v1beta1/drives/{drive-id}/root/children | Create a drive item |
| [**createLinkSpaceRoot**](DrivesRootApi.md#createlinkspaceroot) | **POST** /v1beta1/drives/{drive-id}/root/createLink | Create a sharing link for the root item of a Drive |
| [**deletePermissionSpaceRoot**](DrivesRootApi.md#deletepermissionspaceroot) | **DELETE** /v1beta1/drives/{drive-id}/root/permissions/{perm-id} | Remove access to a Drive |
| [**getPermissionSpaceRoot**](DrivesRootApi.md#getpermissionspaceroot) | **GET** /v1beta1/drives/{drive-id}/root/permissions/{perm-id} | Get a single sharing permission for the root item of a drive |
| [**getRoot**](DrivesRootApi.md#getroot) | **GET** /v1.0/drives/{drive-id}/root | Get root from arbitrary space |
| [**inviteSpaceRoot**](DrivesRootApi.md#invitespaceroot) | **POST** /v1beta1/drives/{drive-id}/root/invite | Send a sharing invitation |
| [**listPermissionsSpaceRoot**](DrivesRootApi.md#listpermissionsspaceroot) | **GET** /v1beta1/drives/{drive-id}/root/permissions | List the effective permissions on the root item of a drive. |
| [**setPermissionPasswordSpaceRoot**](DrivesRootApi.md#setpermissionpasswordspaceroot) | **POST** /v1beta1/drives/{drive-id}/root/permissions/{perm-id}/setPassword | Set sharing link password for the root item of a drive |
| [**updatePermissionSpaceRoot**](DrivesRootApi.md#updatepermissionspaceroot) | **PATCH** /v1beta1/drives/{drive-id}/root/permissions/{perm-id} | Update sharing permission |



## createDriveItem

> DriveItem createDriveItem(driveId, driveItem)

Create a drive item

You can use the root childrens endpoint to mount a remoteItem in the share jail. The &#x60;@client.synchronize&#x60; property of the &#x60;driveItem&#x60; in the [sharedWithMe](#/me.drive/ListSharedWithMe) endpoint will change to true. 

### Example

```ts
import {
  Configuration,
  DrivesRootApi,
} from '';
import type { CreateDriveItemRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // To configure HTTP basic authorization: basicAuth
    username: "YOUR USERNAME",
    password: "YOUR PASSWORD",
  });
  const api = new DrivesRootApi(config);

  const body = {
    // string | key: id of drive
    driveId: a0ca6a90-a365-4782-871e-d44447bbc668$a0ca6a90-a365-4782-871e-d44447bbc668,
    // DriveItem | In the request body, provide a JSON object with the following parameters. For mounting a share the necessary remoteItem id and permission id can be taken from the [sharedWithMe](#/me.drive/ListSharedWithMe) endpoint. (optional)
    driveItem: {"name":"Einsteins project share","remoteItem":{"id":"a-storage-provider-id$a-space-id!a-node-id"}},
  } satisfies CreateDriveItemRequest;

  try {
    const data = await api.createDriveItem(body);
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
| **driveItem** | [DriveItem](DriveItem.md) | In the request body, provide a JSON object with the following parameters. For mounting a share the necessary remoteItem id and permission id can be taken from the [sharedWithMe](#/me.drive/ListSharedWithMe) endpoint. | [Optional] |

### Return type

[**DriveItem**](DriveItem.md)

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Response |  -  |
| **0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## createLinkSpaceRoot

> Permission createLinkSpaceRoot(driveId, driveItemCreateLink)

Create a sharing link for the root item of a Drive

You can use the createLink action to share a driveItem via a sharing link.  The response will be a permission object with the link facet containing the created link details.  ## Link types  For now, The following values are allowed for the type parameter.  | Value          | Display name      | Description                                                     | | -------------- | ----------------- | --------------------------------------------------------------- | | view           | View              | Creates a read-only link to the driveItem.                      | | upload         | Upload            | Creates a read-write link to the folder driveItem.              | | edit           | Edit              | Creates a read-write link to the driveItem.                     | | createOnly     | File Drop         | Creates an upload-only link to the folder driveItem.            | | blocksDownload | Secure View       | Creates a read-only link that blocks download to the driveItem. | 

### Example

```ts
import {
  Configuration,
  DrivesRootApi,
} from '';
import type { CreateLinkSpaceRootRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // To configure HTTP basic authorization: basicAuth
    username: "YOUR USERNAME",
    password: "YOUR PASSWORD",
  });
  const api = new DrivesRootApi(config);

  const body = {
    // string | key: id of drive
    driveId: driveId_example,
    // DriveItemCreateLink | In the request body, provide a JSON object with the following parameters. (optional)
    driveItemCreateLink: {"type":"view"},
  } satisfies CreateLinkSpaceRootRequest;

  try {
    const data = await api.createLinkSpaceRoot(body);
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


## deletePermissionSpaceRoot

> deletePermissionSpaceRoot(driveId, permId)

Remove access to a Drive

Remove access to the root item of a drive.  Only sharing permissions that are not inherited can be deleted. The &#x60;inheritedFrom&#x60; property must be &#x60;null&#x60;. 

### Example

```ts
import {
  Configuration,
  DrivesRootApi,
} from '';
import type { DeletePermissionSpaceRootRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // To configure HTTP basic authorization: basicAuth
    username: "YOUR USERNAME",
    password: "YOUR PASSWORD",
  });
  const api = new DrivesRootApi(config);

  const body = {
    // string | key: id of drive
    driveId: driveId_example,
    // string | key: id of permission
    permId: permId_example,
  } satisfies DeletePermissionSpaceRootRequest;

  try {
    const data = await api.deletePermissionSpaceRoot(body);
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


## getPermissionSpaceRoot

> Permission getPermissionSpaceRoot(driveId, permId)

Get a single sharing permission for the root item of a drive

Return the effective sharing permission for a particular permission resource. 

### Example

```ts
import {
  Configuration,
  DrivesRootApi,
} from '';
import type { GetPermissionSpaceRootRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // To configure HTTP basic authorization: basicAuth
    username: "YOUR USERNAME",
    password: "YOUR PASSWORD",
  });
  const api = new DrivesRootApi(config);

  const body = {
    // string | key: id of drive
    driveId: driveId_example,
    // string | key: id of permission
    permId: permId_example,
  } satisfies GetPermissionSpaceRootRequest;

  try {
    const data = await api.getPermissionSpaceRoot(body);
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


## getRoot

> DriveItem getRoot(driveId)

Get root from arbitrary space

### Example

```ts
import {
  Configuration,
  DrivesRootApi,
} from '';
import type { GetRootRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // To configure HTTP basic authorization: basicAuth
    username: "YOUR USERNAME",
    password: "YOUR PASSWORD",
  });
  const api = new DrivesRootApi(config);

  const body = {
    // string | key: id of drive
    driveId: driveId_example,
  } satisfies GetRootRequest;

  try {
    const data = await api.getRoot(body);
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

### Return type

[**DriveItem**](DriveItem.md)

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


## inviteSpaceRoot

> CollectionOfPermissions inviteSpaceRoot(driveId, driveItemInvite)

Send a sharing invitation

Sends a sharing invitation for the root of a &#x60;drive&#x60;. A sharing invitation provides permissions to the recipients and optionally sends them an email with a sharing link.  The response will be a permission object with the grantedToV2 property containing the created grant details.  ## Roles property values For now, roles are only identified by a uuid. There are no hardcoded aliases like &#x60;read&#x60; or &#x60;write&#x60; because role actions can be completely customized. 

### Example

```ts
import {
  Configuration,
  DrivesRootApi,
} from '';
import type { InviteSpaceRootRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // To configure HTTP basic authorization: basicAuth
    username: "YOUR USERNAME",
    password: "YOUR PASSWORD",
  });
  const api = new DrivesRootApi(config);

  const body = {
    // string | key: id of drive
    driveId: driveId_example,
    // DriveItemInvite | In the request body, provide a JSON object with the following parameters. To create a custom role submit a list of actions instead of roles. (optional)
    driveItemInvite: {"recipients":[{"@libre.graph.recipient.type":"user","objectId":"4c510ada-c86b-4815-8820-42cdf82c3d51"}],"roles":["b1e2218d-eef8-4d4c-b82d-0f1a1b48f3b5"]},
  } satisfies InviteSpaceRootRequest;

  try {
    const data = await api.inviteSpaceRoot(body);
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


## listPermissionsSpaceRoot

> CollectionOfPermissionsWithAllowedValues listPermissionsSpaceRoot(driveId, $filter, $select)

List the effective permissions on the root item of a drive.

The permissions collection includes potentially sensitive information and may not be available for every caller.  * For the owner of the item, all sharing permissions will be returned. This includes co-owners. * For a non-owner caller, only the sharing permissions that apply to the caller are returned. * Sharing permission properties that contain secrets (e.g. &#x60;webUrl&#x60;) are only returned for callers that are able to create the sharing permission.  All permission objects have an &#x60;id&#x60;. A permission representing * a link has the &#x60;link&#x60; facet filled with details. * a share has the &#x60;roles&#x60; property set and the &#x60;grantedToV2&#x60; property filled with the grant recipient details. 

### Example

```ts
import {
  Configuration,
  DrivesRootApi,
} from '';
import type { ListPermissionsSpaceRootRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // To configure HTTP basic authorization: basicAuth
    username: "YOUR USERNAME",
    password: "YOUR PASSWORD",
  });
  const api = new DrivesRootApi(config);

  const body = {
    // string | key: id of drive
    driveId: driveId_example,
    // string | Filter items by property values. By default all permissions are returned and the avalable sharing roles are limited to normal users. To get a list of sharing roles applicable to federated users use the example $select query and combine it with $filter to omit the list of permissions. (optional)
    $filter: @libre.graph.permissions.roles.allowedValues/rolePermissions/any(p:contains(p/condition, '@Subject.UserType=="Federated"')),
    // Set<'@libre.graph.permissions.actions.allowedValues' | '@libre.graph.permissions.roles.allowedValues' | 'value'> | Select properties to be returned. By default all properties are returned. Select the roles property to fetch the available sharing roles without resolving all the permissions. Combine this with the $filter parameter to fetch the actions applicable to federated users. (optional)
    $select: ...,
  } satisfies ListPermissionsSpaceRootRequest;

  try {
    const data = await api.listPermissionsSpaceRoot(body);
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


## setPermissionPasswordSpaceRoot

> Permission setPermissionPasswordSpaceRoot(driveId, permId, sharingLinkPassword)

Set sharing link password for the root item of a drive

Set the password of a sharing permission.  Only the &#x60;password&#x60; property can be modified this way. 

### Example

```ts
import {
  Configuration,
  DrivesRootApi,
} from '';
import type { SetPermissionPasswordSpaceRootRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // To configure HTTP basic authorization: basicAuth
    username: "YOUR USERNAME",
    password: "YOUR PASSWORD",
  });
  const api = new DrivesRootApi(config);

  const body = {
    // string | key: id of drive
    driveId: driveId_example,
    // string | key: id of permission
    permId: permId_example,
    // SharingLinkPassword | New password value
    sharingLinkPassword: {"password":"TestPassword123!"},
  } satisfies SetPermissionPasswordSpaceRootRequest;

  try {
    const data = await api.setPermissionPasswordSpaceRoot(body);
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


## updatePermissionSpaceRoot

> Permission updatePermissionSpaceRoot(driveId, permId, permission)

Update sharing permission

Update the properties of a sharing permission by patching the permission resource.  Only the &#x60;roles&#x60;, &#x60;expirationDateTime&#x60; and &#x60;password&#x60; properties can be modified this way. 

### Example

```ts
import {
  Configuration,
  DrivesRootApi,
} from '';
import type { UpdatePermissionSpaceRootRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // To configure HTTP basic authorization: basicAuth
    username: "YOUR USERNAME",
    password: "YOUR PASSWORD",
  });
  const api = new DrivesRootApi(config);

  const body = {
    // string | key: id of drive
    driveId: driveId_example,
    // string | key: id of permission
    permId: permId_example,
    // Permission | New property values
    permission: {"link":{"type":"edit"}},
  } satisfies UpdatePermissionSpaceRootRequest;

  try {
    const data = await api.updatePermissionSpaceRoot(body);
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

