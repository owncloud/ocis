# MeDrivesApi

All URIs are relative to *https://ocis.ocis.rolling.owncloud.works/graph*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**listMyDrives**](MeDrivesApi.md#listmydrives) | **GET** /v1.0/me/drives | Get all drives where the current user is a regular member of |
| [**listMyDrivesBeta**](MeDrivesApi.md#listmydrivesbeta) | **GET** /v1beta1/me/drives | Alias for \&#39;/v1.0/drives\&#39;, the difference is that grantedtoV2 is used and roles contain unified roles instead of cs3 roles |



## listMyDrives

> CollectionOfDrives listMyDrives($orderby, $filter)

Get all drives where the current user is a regular member of

### Example

```ts
import {
  Configuration,
  MeDrivesApi,
} from '';
import type { ListMyDrivesRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // To configure HTTP basic authorization: basicAuth
    username: "YOUR USERNAME",
    password: "YOUR PASSWORD",
  });
  const api = new MeDrivesApi(config);

  const body = {
    // string | The $orderby system query option allows clients to request resources in either ascending order using asc or descending order using desc. (optional)
    $orderby: lastModifiedDateTime desc,
    // string | Filter items by property values (optional)
    $filter: driveType eq 'project',
  } satisfies ListMyDrivesRequest;

  try {
    const data = await api.listMyDrives(body);
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

[**CollectionOfDrives**](CollectionOfDrives.md)

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


## listMyDrivesBeta

> CollectionOfDrives listMyDrivesBeta($orderby, $filter)

Alias for \&#39;/v1.0/drives\&#39;, the difference is that grantedtoV2 is used and roles contain unified roles instead of cs3 roles

### Example

```ts
import {
  Configuration,
  MeDrivesApi,
} from '';
import type { ListMyDrivesBetaRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // To configure HTTP basic authorization: basicAuth
    username: "YOUR USERNAME",
    password: "YOUR PASSWORD",
  });
  const api = new MeDrivesApi(config);

  const body = {
    // string | The $orderby system query option allows clients to request resources in either ascending order using asc or descending order using desc. (optional)
    $orderby: lastModifiedDateTime desc,
    // string | Filter items by property values (optional)
    $filter: driveType eq 'project',
  } satisfies ListMyDrivesBetaRequest;

  try {
    const data = await api.listMyDrivesBeta(body);
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

[**CollectionOfDrives**](CollectionOfDrives.md)

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

