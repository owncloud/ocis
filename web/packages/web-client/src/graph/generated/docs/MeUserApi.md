# MeUserApi

All URIs are relative to *https://ocis.ocis.rolling.owncloud.works/graph*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**getOwnUser**](#getownuser) | **GET** /v1.0/me | Get current user|
|[**updateOwnUser**](#updateownuser) | **PATCH** /v1.0/me | Update the current user|

# **getOwnUser**
> User getOwnUser()


### Example

```typescript
import {
    MeUserApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new MeUserApi(configuration);

let $expand: Set<'memberOf'>; //Expand related entities (optional) (default to undefined)

const { status, data } = await apiInstance.getOwnUser(
    $expand
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **$expand** | **Array<&#39;memberOf&#39;>** | Expand related entities | (optional) defaults to undefined|


### Return type

**User**

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Retrieved entity |  -  |
|**0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **updateOwnUser**
> User updateOwnUser()


### Example

```typescript
import {
    MeUserApi,
    Configuration,
    UserUpdate
} from './api';

const configuration = new Configuration();
const apiInstance = new MeUserApi(configuration);

let userUpdate: UserUpdate; //New user values (optional)

const { status, data } = await apiInstance.updateOwnUser(
    userUpdate
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **userUpdate** | **UserUpdate**| New user values | |


### Return type

**User**

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Success |  -  |
|**0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

