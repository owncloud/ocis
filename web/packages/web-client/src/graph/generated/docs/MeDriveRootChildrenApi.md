# MeDriveRootChildrenApi

All URIs are relative to *https://ocis.ocis.rolling.owncloud.works/graph*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**homeGetChildren**](MeDriveRootChildrenApi.md#homegetchildren) | **GET** /v1.0/me/drive/root/children | Get children from drive |



## homeGetChildren

> CollectionOfDriveItems homeGetChildren()

Get children from drive

### Example

```ts
import {
  Configuration,
  MeDriveRootChildrenApi,
} from '';
import type { HomeGetChildrenRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // To configure HTTP basic authorization: basicAuth
    username: "YOUR USERNAME",
    password: "YOUR PASSWORD",
  });
  const api = new MeDriveRootChildrenApi(config);

  try {
    const data = await api.homeGetChildren();
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

[**CollectionOfDriveItems**](CollectionOfDriveItems.md)

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Retrieved resource list |  -  |
| **0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

