# EducationUserApi

All URIs are relative to *https://ocis.ocis.rolling.owncloud.works/graph*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**createEducationUser**](EducationUserApi.md#createeducationuser) | **POST** /v1.0/education/users | Add new education user |
| [**deleteEducationUser**](EducationUserApi.md#deleteeducationuser) | **DELETE** /v1.0/education/users/{user-id} | Delete educationUser |
| [**getEducationUser**](EducationUserApi.md#geteducationuser) | **GET** /v1.0/education/users/{user-id} | Get properties of educationUser |
| [**listEducationUsers**](EducationUserApi.md#listeducationusers) | **GET** /v1.0/education/users | Get entities from education users |
| [**updateEducationUser**](EducationUserApi.md#updateeducationuser) | **PATCH** /v1.0/education/users/{user-id} | Update properties of educationUser |



## createEducationUser

> EducationUser createEducationUser(educationUser)

Add new education user

### Example

```ts
import {
  Configuration,
  EducationUserApi,
} from '';
import type { CreateEducationUserRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new EducationUserApi(config);

  const body = {
    // EducationUser | New entity
    educationUser: ...,
  } satisfies CreateEducationUserRequest;

  try {
    const data = await api.createEducationUser(body);
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
| **educationUser** | [EducationUser](EducationUser.md) | New entity | |

### Return type

[**EducationUser**](EducationUser.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **201** | Created entity |  -  |
| **0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## deleteEducationUser

> deleteEducationUser(userId)

Delete educationUser

### Example

```ts
import {
  Configuration,
  EducationUserApi,
} from '';
import type { DeleteEducationUserRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new EducationUserApi(config);

  const body = {
    // string | key: id or username of user
    userId: 90eedea1-dea1-90ee-a1de-ee90a1deee90,
  } satisfies DeleteEducationUserRequest;

  try {
    const data = await api.deleteEducationUser(body);
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
| **userId** | `string` | key: id or username of user | [Defaults to `undefined`] |

### Return type

`void` (Empty response body)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **204** | Success |  -  |
| **0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## getEducationUser

> EducationUser getEducationUser(userId, $expand)

Get properties of educationUser

### Example

```ts
import {
  Configuration,
  EducationUserApi,
} from '';
import type { GetEducationUserRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new EducationUserApi(config);

  const body = {
    // string | key: id or username of user
    userId: 90eedea1-dea1-90ee-a1de-ee90a1deee90,
    // Set<'memberOf'> | Expand related entities (optional)
    $expand: ...,
  } satisfies GetEducationUserRequest;

  try {
    const data = await api.getEducationUser(body);
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
| **userId** | `string` | key: id or username of user | [Defaults to `undefined`] |
| **$expand** | `memberOf` | Expand related entities | [Optional] [Enum: memberOf] |

### Return type

[**EducationUser**](EducationUser.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Retrieved entity |  -  |
| **0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## listEducationUsers

> CollectionOfEducationUser listEducationUsers($orderby, $expand)

Get entities from education users

### Example

```ts
import {
  Configuration,
  EducationUserApi,
} from '';
import type { ListEducationUsersRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new EducationUserApi(config);

  const body = {
    // Set<'displayName' | 'displayName desc' | 'mail' | 'mail desc' | 'onPremisesSamAccountName' | 'onPremisesSamAccountName desc'> | Order items by property values (optional)
    $orderby: ...,
    // Set<'memberOf'> | Expand related entities (optional)
    $expand: ...,
  } satisfies ListEducationUsersRequest;

  try {
    const data = await api.listEducationUsers(body);
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
| **$orderby** | `displayName`, `displayName desc`, `mail`, `mail desc`, `onPremisesSamAccountName`, `onPremisesSamAccountName desc` | Order items by property values | [Optional] [Enum: displayName, displayName desc, mail, mail desc, onPremisesSamAccountName, onPremisesSamAccountName desc] |
| **$expand** | `memberOf` | Expand related entities | [Optional] [Enum: memberOf] |

### Return type

[**CollectionOfEducationUser**](CollectionOfEducationUser.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Retrieved entities |  -  |
| **0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## updateEducationUser

> EducationUser updateEducationUser(userId, educationUser)

Update properties of educationUser

### Example

```ts
import {
  Configuration,
  EducationUserApi,
} from '';
import type { UpdateEducationUserRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new EducationUserApi(config);

  const body = {
    // string | key: id or username of user
    userId: 90eedea1-dea1-90ee-a1de-ee90a1deee90,
    // EducationUser | New property values
    educationUser: {"mail":"max.mustermann@new.domain"},
  } satisfies UpdateEducationUserRequest;

  try {
    const data = await api.updateEducationUser(body);
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
| **userId** | `string` | key: id or username of user | [Defaults to `undefined`] |
| **educationUser** | [EducationUser](EducationUser.md) | New property values | |

### Return type

[**EducationUser**](EducationUser.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Success |  -  |
| **204** | Success |  -  |
| **0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

