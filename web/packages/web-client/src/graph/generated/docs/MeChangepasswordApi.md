# MeChangepasswordApi

All URIs are relative to *https://ocis.ocis.rolling.owncloud.works/graph*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**changeOwnPassword**](#changeownpassword) | **POST** /v1.0/me/changePassword | Change your own password|

# **changeOwnPassword**
> changeOwnPassword(passwordChange)


### Example

```typescript
import {
    MeChangepasswordApi,
    Configuration,
    PasswordChange
} from './api';

const configuration = new Configuration();
const apiInstance = new MeChangepasswordApi(configuration);

let passwordChange: PasswordChange; //Password change request

const { status, data } = await apiInstance.changeOwnPassword(
    passwordChange
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **passwordChange** | **PasswordChange**| Password change request | |


### Return type

void (empty response body)

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**204** | Success |  -  |
|**0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

