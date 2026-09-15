# EducationClassApi

All URIs are relative to *https://ocis.ocis.rolling.owncloud.works/graph*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**addUserToClass**](EducationClassApi.md#addusertoclass) | **POST** /v1.0/education/classes/{class-id}/members/$ref | Assign a user to a class |
| [**createClass**](EducationClassApi.md#createclass) | **POST** /v1.0/education/classes | Add new education class |
| [**deleteClass**](EducationClassApi.md#deleteclass) | **DELETE** /v1.0/education/classes/{class-id} | Delete education class |
| [**deleteUserFromClass**](EducationClassApi.md#deleteuserfromclass) | **DELETE** /v1.0/education/classes/{class-id}/members/{user-id}/$ref | Unassign user from a class |
| [**getClass**](EducationClassApi.md#getclass) | **GET** /v1.0/education/classes/{class-id} | Get class by key |
| [**listClassMembers**](EducationClassApi.md#listclassmembers) | **GET** /v1.0/education/classes/{class-id}/members | Get the educationClass resources owned by an educationSchool |
| [**listClasses**](EducationClassApi.md#listclasses) | **GET** /v1.0/education/classes | list education classes |
| [**updateClass**](EducationClassApi.md#updateclass) | **PATCH** /v1.0/education/classes/{class-id} | Update properties of a education class |



## addUserToClass

> addUserToClass(classId, classMemberReference)

Assign a user to a class

### Example

```ts
import {
  Configuration,
  EducationClassApi,
} from '';
import type { AddUserToClassRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new EducationClassApi(config);

  const body = {
    // string | key: id or externalId of class
    classId: 86948e45-96a6-43df-b83d-46e92afd30de,
    // ClassMemberReference | educationUser to be added as member
    classMemberReference: ...,
  } satisfies AddUserToClassRequest;

  try {
    const data = await api.addUserToClass(body);
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
| **classId** | `string` | key: id or externalId of class | [Defaults to `undefined`] |
| **classMemberReference** | [ClassMemberReference](ClassMemberReference.md) | educationUser to be added as member | |

### Return type

`void` (Empty response body)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **204** | Success |  -  |
| **0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## createClass

> EducationClass createClass(educationClass)

Add new education class

### Example

```ts
import {
  Configuration,
  EducationClassApi,
} from '';
import type { CreateClassRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new EducationClassApi(config);

  const body = {
    // EducationClass | New entity
    educationClass: ...,
  } satisfies CreateClassRequest;

  try {
    const data = await api.createClass(body);
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
| **educationClass** | [EducationClass](EducationClass.md) | New entity | |

### Return type

[**EducationClass**](EducationClass.md)

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


## deleteClass

> deleteClass(classId)

Delete education class

### Example

```ts
import {
  Configuration,
  EducationClassApi,
} from '';
import type { DeleteClassRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new EducationClassApi(config);

  const body = {
    // string | key: id or externalId of class
    classId: 86948e45-96a6-43df-b83d-46e92afd30de,
  } satisfies DeleteClassRequest;

  try {
    const data = await api.deleteClass(body);
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
| **classId** | `string` | key: id or externalId of class | [Defaults to `undefined`] |

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


## deleteUserFromClass

> deleteUserFromClass(classId, userId)

Unassign user from a class

### Example

```ts
import {
  Configuration,
  EducationClassApi,
} from '';
import type { DeleteUserFromClassRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new EducationClassApi(config);

  const body = {
    // string | key: id or externalId of class
    classId: classId_example,
    // string | key: id or username of the user to unassign from class
    userId: 90eedea1-dea1-90ee-a1de-ee90a1deee90,
  } satisfies DeleteUserFromClassRequest;

  try {
    const data = await api.deleteUserFromClass(body);
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
| **classId** | `string` | key: id or externalId of class | [Defaults to `undefined`] |
| **userId** | `string` | key: id or username of the user to unassign from class | [Defaults to `undefined`] |

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


## getClass

> EducationClass getClass(classId)

Get class by key

### Example

```ts
import {
  Configuration,
  EducationClassApi,
} from '';
import type { GetClassRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new EducationClassApi(config);

  const body = {
    // string | key: id or externalId of class
    classId: 86948e45-96a6-43df-b83d-46e92afd30de,
  } satisfies GetClassRequest;

  try {
    const data = await api.getClass(body);
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
| **classId** | `string` | key: id or externalId of class | [Defaults to `undefined`] |

### Return type

[**EducationClass**](EducationClass.md)

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


## listClassMembers

> CollectionOfEducationUser listClassMembers(classId)

Get the educationClass resources owned by an educationSchool

### Example

```ts
import {
  Configuration,
  EducationClassApi,
} from '';
import type { ListClassMembersRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new EducationClassApi(config);

  const body = {
    // string | key: id or externalId of class
    classId: 86948e45-96a6-43df-b83d-46e92afd30de,
  } satisfies ListClassMembersRequest;

  try {
    const data = await api.listClassMembers(body);
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
| **classId** | `string` | key: id or externalId of class | [Defaults to `undefined`] |

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
| **200** | Retrieved class members |  -  |
| **0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## listClasses

> CollectionOfClass listClasses()

list education classes

### Example

```ts
import {
  Configuration,
  EducationClassApi,
} from '';
import type { ListClassesRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new EducationClassApi(config);

  try {
    const data = await api.listClasses();
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

[**CollectionOfClass**](CollectionOfClass.md)

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


## updateClass

> EducationClass updateClass(classId, educationClass)

Update properties of a education class

### Example

```ts
import {
  Configuration,
  EducationClassApi,
} from '';
import type { UpdateClassRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new EducationClassApi(config);

  const body = {
    // string | key: id or externalId of class
    classId: 86948e45-96a6-43df-b83d-46e92afd30de,
    // EducationClass | New property values
    educationClass: {"displayName":"Musik"},
  } satisfies UpdateClassRequest;

  try {
    const data = await api.updateClass(body);
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
| **classId** | `string` | key: id or externalId of class | [Defaults to `undefined`] |
| **educationClass** | [EducationClass](EducationClass.md) | New property values | |

### Return type

[**EducationClass**](EducationClass.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | New property values |  -  |
| **204** | Success |  -  |
| **0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

