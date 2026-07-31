# DrivesGetDrivesApi

All URIs are relative to *https://ocis.ocis.rolling.owncloud.works/graph*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**listAllDrives**](#listalldrives) | **GET** /v1.0/drives | Get all available drives|
|[**listAllDrivesBeta**](#listalldrivesbeta) | **GET** /v1beta1/drives | Alias for \&#39;/v1.0/drives\&#39;, the difference is that grantedtoV2 is used and roles contain unified roles instead of cs3 roles|

# **listAllDrives**
> CollectionOfDrives1 listAllDrives()


### Example

```typescript
import {
    DrivesGetDrivesApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DrivesGetDrivesApi(configuration);

let $orderby: string; //The $orderby system query option allows clients to request resources in either ascending order using asc or descending order using desc. (optional) (default to undefined)
let $filter: string; //Filter items by property values (optional) (default to undefined)

const { status, data } = await apiInstance.listAllDrives(
    $orderby,
    $filter
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **$orderby** | [**string**] | The $orderby system query option allows clients to request resources in either ascending order using asc or descending order using desc. | (optional) defaults to undefined|
| **$filter** | [**string**] | Filter items by property values | (optional) defaults to undefined|


### Return type

**CollectionOfDrives1**

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Retrieved spaces |  -  |
|**0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **listAllDrivesBeta**
> CollectionOfDrives1 listAllDrivesBeta()


### Example

```typescript
import {
    DrivesGetDrivesApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DrivesGetDrivesApi(configuration);

let $orderby: string; //The $orderby system query option allows clients to request resources in either ascending order using asc or descending order using desc. (optional) (default to undefined)
let $filter: string; //Filter items by property values (optional) (default to undefined)

const { status, data } = await apiInstance.listAllDrivesBeta(
    $orderby,
    $filter
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **$orderby** | [**string**] | The $orderby system query option allows clients to request resources in either ascending order using asc or descending order using desc. | (optional) defaults to undefined|
| **$filter** | [**string**] | Filter items by property values | (optional) defaults to undefined|


### Return type

**CollectionOfDrives1**

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Retrieved spaces |  -  |
|**0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

