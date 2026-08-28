# UserApi

All URIs are relative to *https://ocis.ocis.rolling.owncloud.works/graph*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**deleteUser**](UserApi.md#deleteuser) | **DELETE** /v1.0/users/{user-id} | Delete entity from users |
| [**exportPersonalData**](UserApi.md#exportpersonaldataoperation) | **POST** /v1.0/users/{user-id}/exportPersonalData | export personal data of a user |
| [**getUser**](UserApi.md#getuser) | **GET** /v1.0/users/{user-id} | Get entity from users by key |
| [**updateUser**](UserApi.md#updateuser) | **PATCH** /v1.0/users/{user-id} | Update entity in users |



## deleteUser

> deleteUser(userId, ifMatch)

Delete entity from users

### Example

```ts
import {
  Configuration,
  UserApi,
} from '';
import type { DeleteUserRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // To configure HTTP basic authorization: basicAuth
    username: "YOUR USERNAME",
    password: "YOUR PASSWORD",
  });
  const api = new UserApi(config);

  const body = {
    // string | key: id or name of user
    userId: userId_example,
    // string | ETag (optional)
    ifMatch: ifMatch_example,
  } satisfies DeleteUserRequest;

  try {
    const data = await api.deleteUser(body);
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
| **userId** | `string` | key: id or name of user | [Defaults to `undefined`] |
| **ifMatch** | `string` | ETag | [Optional] [Defaults to `undefined`] |

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


## exportPersonalData

> exportPersonalData(userId, exportPersonalDataRequest)

export personal data of a user

### Example

```ts
import {
  Configuration,
  UserApi,
} from '';
import type { ExportPersonalDataOperationRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // To configure HTTP basic authorization: basicAuth
    username: "YOUR USERNAME",
    password: "YOUR PASSWORD",
  });
  const api = new UserApi(config);

  const body = {
    // string | key: id or name of user
    userId: userId_example,
    // ExportPersonalDataRequest | destination the file should be created at (optional)
    exportPersonalDataRequest: ...,
  } satisfies ExportPersonalDataOperationRequest;

  try {
    const data = await api.exportPersonalData(body);
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
| **userId** | `string` | key: id or name of user | [Defaults to `undefined`] |
| **exportPersonalDataRequest** | [ExportPersonalDataRequest](ExportPersonalDataRequest.md) | destination the file should be created at | [Optional] |

### Return type

`void` (Empty response body)

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **202** | success |  -  |
| **0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## getUser

> User getUser(userId, $select, $expand)

Get entity from users by key

### Example

```ts
import {
  Configuration,
  UserApi,
} from '';
import type { GetUserRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // To configure HTTP basic authorization: basicAuth
    username: "YOUR USERNAME",
    password: "YOUR PASSWORD",
  });
  const api = new UserApi(config);

  const body = {
    // string | key: id or name of user
    userId: userId_example,
    // Set<'id' | 'displayName' | 'drive' | 'drives' | 'mail' | 'memberOf' | 'onPremisesSamAccountName' | 'surname'> | Select properties to be returned (optional)
    $select: ...,
    // Set<'drive' | 'drives' | 'memberOf' | 'appRoleAssignments'> | Expand related entities (optional)
    $expand: ...,
  } satisfies GetUserRequest;

  try {
    const data = await api.getUser(body);
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
| **userId** | `string` | key: id or name of user | [Defaults to `undefined`] |
| **$select** | `id`, `displayName`, `drive`, `drives`, `mail`, `memberOf`, `onPremisesSamAccountName`, `surname` | Select properties to be returned | [Optional] [Enum: id, displayName, drive, drives, mail, memberOf, onPremisesSamAccountName, surname] |
| **$expand** | `drive`, `drives`, `memberOf`, `appRoleAssignments` | Expand related entities | [Optional] [Enum: drive, drives, memberOf, appRoleAssignments] |

### Return type

[**User**](User.md)

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Retrieved entity |  -  |
| **0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## updateUser

> User updateUser(userId, userUpdate)

Update entity in users

### Example

```ts
import {
  Configuration,
  UserApi,
} from '';
import type { UpdateUserRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // To configure HTTP basic authorization: basicAuth
    username: "YOUR USERNAME",
    password: "YOUR PASSWORD",
  });
  const api = new UserApi(config);

  const body = {
    // string | key: id of user
    userId: userId_example,
    // UserUpdate | New property values
    userUpdate: {"displayName":"Marie Skłodowska Curie"},
  } satisfies UpdateUserRequest;

  try {
    const data = await api.updateUser(body);
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
| **userId** | `string` | key: id of user | [Defaults to `undefined`] |
| **userUpdate** | [UserUpdate](UserUpdate.md) | New property values | |

### Return type

[**User**](User.md)

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Success |  -  |
| **0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

