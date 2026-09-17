# MeUserApi

All URIs are relative to *https://ocis.ocis.rolling.owncloud.works/graph*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**getOwnUser**](MeUserApi.md#getownuser) | **GET** /v1.0/me | Get current user |
| [**updateOwnUser**](MeUserApi.md#updateownuser) | **PATCH** /v1.0/me | Update the current user |



## getOwnUser

> User getOwnUser($expand)

Get current user

### Example

```ts
import {
  Configuration,
  MeUserApi,
} from '';
import type { GetOwnUserRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // To configure HTTP basic authorization: basicAuth
    username: "YOUR USERNAME",
    password: "YOUR PASSWORD",
  });
  const api = new MeUserApi(config);

  const body = {
    // Set<'memberOf'> | Expand related entities (optional)
    $expand: ...,
  } satisfies GetOwnUserRequest;

  try {
    const data = await api.getOwnUser(body);
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
| **$expand** | `memberOf` | Expand related entities | [Optional] [Enum: memberOf] |

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


## updateOwnUser

> User updateOwnUser(userUpdate)

Update the current user

### Example

```ts
import {
  Configuration,
  MeUserApi,
} from '';
import type { UpdateOwnUserRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // To configure HTTP basic authorization: basicAuth
    username: "YOUR USERNAME",
    password: "YOUR PASSWORD",
  });
  const api = new MeUserApi(config);

  const body = {
    // UserUpdate | New user values (optional)
    userUpdate: {"preferredLanguage":"en"},
  } satisfies UpdateOwnUserRequest;

  try {
    const data = await api.updateOwnUser(body);
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
| **userUpdate** | [UserUpdate](UserUpdate.md) | New user values | [Optional] |

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

