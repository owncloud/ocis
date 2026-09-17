# MeDriveApi

All URIs are relative to *https://ocis.ocis.rolling.owncloud.works/graph*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**getHome**](MeDriveApi.md#gethome) | **GET** /v1.0/me/drive | Get personal space for user |
| [**listSharedByMe**](MeDriveApi.md#listsharedbyme) | **GET** /v1beta1/me/drive/sharedByMe | Get a list of driveItem objects shared by the current user. |
| [**listSharedWithMe**](MeDriveApi.md#listsharedwithme) | **GET** /v1beta1/me/drive/sharedWithMe | Get a list of driveItem objects shared with the owner of a drive. |



## getHome

> Drive getHome()

Get personal space for user

### Example

```ts
import {
  Configuration,
  MeDriveApi,
} from '';
import type { GetHomeRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // To configure HTTP basic authorization: basicAuth
    username: "YOUR USERNAME",
    password: "YOUR PASSWORD",
  });
  const api = new MeDriveApi(config);

  try {
    const data = await api.getHome();
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters

This endpoint does not need any parameter.

### Return type

[**Drive**](Drive.md)

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Retrieved personal space |  -  |
| **0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## listSharedByMe

> CollectionOfDriveItems1 listSharedByMe()

Get a list of driveItem objects shared by the current user.

The &#x60;driveItems&#x60; returned from the &#x60;sharedByMe&#x60; method always include the &#x60;permissions&#x60; relation that indicates they are shared items. 

### Example

```ts
import {
  Configuration,
  MeDriveApi,
} from '';
import type { ListSharedByMeRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // To configure HTTP basic authorization: basicAuth
    username: "YOUR USERNAME",
    password: "YOUR PASSWORD",
  });
  const api = new MeDriveApi(config);

  try {
    const data = await api.listSharedByMe();
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters

This endpoint does not need any parameter.

### Return type

[**CollectionOfDriveItems1**](CollectionOfDriveItems1.md)

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | OK |  -  |
| **0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## listSharedWithMe

> CollectionOfDriveItems1 listSharedWithMe()

Get a list of driveItem objects shared with the owner of a drive.

The &#x60;driveItems&#x60; returned from the &#x60;sharedWithMe&#x60; method always include the &#x60;remoteItem&#x60; facet that indicates they are items from a different drive. 

### Example

```ts
import {
  Configuration,
  MeDriveApi,
} from '';
import type { ListSharedWithMeRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // To configure HTTP basic authorization: basicAuth
    username: "YOUR USERNAME",
    password: "YOUR PASSWORD",
  });
  const api = new MeDriveApi(config);

  try {
    const data = await api.listSharedWithMe();
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters

This endpoint does not need any parameter.

### Return type

[**CollectionOfDriveItems1**](CollectionOfDriveItems1.md)

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | OK |  -  |
| **0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

