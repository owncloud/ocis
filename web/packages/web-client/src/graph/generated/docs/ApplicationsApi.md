# ApplicationsApi

All URIs are relative to *https://ocis.ocis.rolling.owncloud.works/graph*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**getApplication**](#getapplication) | **GET** /v1.0/applications/{application-id} | Get application by id|
|[**listApplications**](#listapplications) | **GET** /v1.0/applications | Get all applications|

# **getApplication**
> Application getApplication()


### Example

```typescript
import {
    ApplicationsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new ApplicationsApi(configuration);

let applicationId: string; //key: id of application (default to undefined)

const { status, data } = await apiInstance.getApplication(
    applicationId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **applicationId** | [**string**] | key: id of application | defaults to undefined|


### Return type

**Application**

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |
|**0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **listApplications**
> CollectionOfApplications listApplications()


### Example

```typescript
import {
    ApplicationsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new ApplicationsApi(configuration);

const { status, data } = await apiInstance.listApplications();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**CollectionOfApplications**

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Retrieved entities |  -  |
|**0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

