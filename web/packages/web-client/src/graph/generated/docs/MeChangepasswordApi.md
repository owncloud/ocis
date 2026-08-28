# MeChangepasswordApi

All URIs are relative to *https://ocis.ocis.rolling.owncloud.works/graph*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**changeOwnPassword**](MeChangepasswordApi.md#changeownpassword) | **POST** /v1.0/me/changePassword | Change your own password |



## changeOwnPassword

> changeOwnPassword(passwordChange)

Change your own password

### Example

```ts
import {
  Configuration,
  MeChangepasswordApi,
} from '';
import type { ChangeOwnPasswordRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // To configure HTTP basic authorization: basicAuth
    username: "YOUR USERNAME",
    password: "YOUR PASSWORD",
  });
  const api = new MeChangepasswordApi(config);

  const body = {
    // PasswordChange | Password change request
    passwordChange: ...,
  } satisfies ChangeOwnPasswordRequest;

  try {
    const data = await api.changeOwnPassword(body);
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
| **passwordChange** | [PasswordChange](PasswordChange.md) | Password change request | |

### Return type

`void` (Empty response body)

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **204** | Success |  -  |
| **0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

