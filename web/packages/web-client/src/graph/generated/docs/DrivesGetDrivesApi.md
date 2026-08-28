# DrivesGetDrivesApi

All URIs are relative to *https://ocis.ocis.rolling.owncloud.works/graph*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**listAllDrives**](DrivesGetDrivesApi.md#listalldrives) | **GET** /v1.0/drives | Get all available drives |
| [**listAllDrivesBeta**](DrivesGetDrivesApi.md#listalldrivesbeta) | **GET** /v1beta1/drives | Alias for \&#39;/v1.0/drives\&#39;, the difference is that grantedtoV2 is used and roles contain unified roles instead of cs3 roles |



## listAllDrives

> CollectionOfDrives1 listAllDrives($orderby, $filter)

Get all available drives

### Example

```ts
import {
  Configuration,
  DrivesGetDrivesApi,
} from '';
import type { ListAllDrivesRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // To configure HTTP basic authorization: basicAuth
    username: "YOUR USERNAME",
    password: "YOUR PASSWORD",
  });
  const api = new DrivesGetDrivesApi(config);

  const body = {
    // string | The $orderby system query option allows clients to request resources in either ascending order using asc or descending order using desc. (optional)
    $orderby: lastModifiedDateTime desc,
    // string | Filter items by property values (optional)
    $filter: driveType eq 'project',
  } satisfies ListAllDrivesRequest;

  try {
    const data = await api.listAllDrives(body);
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
| **$orderby** | `string` | The $orderby system query option allows clients to request resources in either ascending order using asc or descending order using desc. | [Optional] [Defaults to `undefined`] |
| **$filter** | `string` | Filter items by property values | [Optional] [Defaults to `undefined`] |

### Return type

[**CollectionOfDrives1**](CollectionOfDrives1.md)

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Retrieved spaces |  -  |
| **0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## listAllDrivesBeta

> CollectionOfDrives1 listAllDrivesBeta($orderby, $filter)

Alias for \&#39;/v1.0/drives\&#39;, the difference is that grantedtoV2 is used and roles contain unified roles instead of cs3 roles

### Example

```ts
import {
  Configuration,
  DrivesGetDrivesApi,
} from '';
import type { ListAllDrivesBetaRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // To configure HTTP basic authorization: basicAuth
    username: "YOUR USERNAME",
    password: "YOUR PASSWORD",
  });
  const api = new DrivesGetDrivesApi(config);

  const body = {
    // string | The $orderby system query option allows clients to request resources in either ascending order using asc or descending order using desc. (optional)
    $orderby: lastModifiedDateTime desc,
    // string | Filter items by property values (optional)
    $filter: driveType eq 'project',
  } satisfies ListAllDrivesBetaRequest;

  try {
    const data = await api.listAllDrivesBeta(body);
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
| **$orderby** | `string` | The $orderby system query option allows clients to request resources in either ascending order using asc or descending order using desc. | [Optional] [Defaults to `undefined`] |
| **$filter** | `string` | Filter items by property values | [Optional] [Defaults to `undefined`] |

### Return type

[**CollectionOfDrives1**](CollectionOfDrives1.md)

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Retrieved spaces |  -  |
| **0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

