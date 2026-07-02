# MeDriveRootChildrenApi

All URIs are relative to *https://ocis.ocis.rolling.owncloud.works/graph*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**homeGetChildren**](#homegetchildren) | **GET** /v1.0/me/drive/root/children | Get children from drive|

# **homeGetChildren**
> CollectionOfDriveItems homeGetChildren()


### Example

```typescript
import {
    MeDriveRootChildrenApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new MeDriveRootChildrenApi(configuration);

const { status, data } = await apiInstance.homeGetChildren();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**CollectionOfDriveItems**

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Retrieved resource list |  -  |
|**0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

