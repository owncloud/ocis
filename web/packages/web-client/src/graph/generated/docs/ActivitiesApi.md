# ActivitiesApi

All URIs are relative to *https://ocis.ocis.rolling.owncloud.works/graph*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**getActivities**](#getactivities) | **GET** /v1beta1/extensions/org.libregraph/activities | Get activities|

# **getActivities**
> CollectionOfActivities getActivities()


### Example

```typescript
import {
    ActivitiesApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new ActivitiesApi(configuration);

let kql: string; // (optional) (default to undefined)

const { status, data } = await apiInstance.getActivities(
    kql
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **kql** | [**string**] |  | (optional) defaults to undefined|


### Return type

**CollectionOfActivities**

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Found activities |  -  |
|**0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

