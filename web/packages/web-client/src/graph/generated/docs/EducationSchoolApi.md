# EducationSchoolApi

All URIs are relative to *https://ocis.ocis.rolling.owncloud.works/graph*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**addClassToSchool**](EducationSchoolApi.md#addclasstoschool) | **POST** /v1.0/education/schools/{school-id}/classes/$ref | Assign a class to a school |
| [**addUserToSchool**](EducationSchoolApi.md#addusertoschool) | **POST** /v1.0/education/schools/{school-id}/users/$ref | Assign a user to a school |
| [**createSchool**](EducationSchoolApi.md#createschool) | **POST** /v1.0/education/schools | Add new school |
| [**deleteClassFromSchool**](EducationSchoolApi.md#deleteclassfromschool) | **DELETE** /v1.0/education/schools/{school-id}/classes/{class-id}/$ref | Unassign class from a school |
| [**deleteSchool**](EducationSchoolApi.md#deleteschool) | **DELETE** /v1.0/education/schools/{school-id} | Delete school |
| [**deleteUserFromSchool**](EducationSchoolApi.md#deleteuserfromschool) | **DELETE** /v1.0/education/schools/{school-id}/users/{user-id}/$ref | Unassign user from a school |
| [**getSchool**](EducationSchoolApi.md#getschool) | **GET** /v1.0/education/schools/{school-id} | Get the properties of a specific school |
| [**listSchoolClasses**](EducationSchoolApi.md#listschoolclasses) | **GET** /v1.0/education/schools/{school-id}/classes | Get the educationClass resources owned by an educationSchool |
| [**listSchoolUsers**](EducationSchoolApi.md#listschoolusers) | **GET** /v1.0/education/schools/{school-id}/users | Get the educationUser resources associated with an educationSchool |
| [**listSchools**](EducationSchoolApi.md#listschools) | **GET** /v1.0/education/schools | Get a list of schools and their properties |
| [**updateSchool**](EducationSchoolApi.md#updateschool) | **PATCH** /v1.0/education/schools/{school-id} | Update properties of a school |



## addClassToSchool

> addClassToSchool(schoolId, classReference)

Assign a class to a school

### Example

```ts
import {
  Configuration,
  EducationSchoolApi,
} from '';
import type { AddClassToSchoolRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new EducationSchoolApi(config);

  const body = {
    // string | key: id or schoolNumber of school
    schoolId: 43b879c4-14c6-4e0a-9b3f-b1b33c5a4bd4,
    // ClassReference | educationClass to be added as member
    classReference: ...,
  } satisfies AddClassToSchoolRequest;

  try {
    const data = await api.addClassToSchool(body);
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
| **schoolId** | `string` | key: id or schoolNumber of school | [Defaults to `undefined`] |
| **classReference** | [ClassReference](ClassReference.md) | educationClass to be added as member | |

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


## addUserToSchool

> addUserToSchool(schoolId, educationUserReference)

Assign a user to a school

### Example

```ts
import {
  Configuration,
  EducationSchoolApi,
} from '';
import type { AddUserToSchoolRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new EducationSchoolApi(config);

  const body = {
    // string | key: id or schoolNumber of school
    schoolId: 43b879c4-14c6-4e0a-9b3f-b1b33c5a4bd4,
    // EducationUserReference | educationUser to be added as member
    educationUserReference: ...,
  } satisfies AddUserToSchoolRequest;

  try {
    const data = await api.addUserToSchool(body);
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
| **schoolId** | `string` | key: id or schoolNumber of school | [Defaults to `undefined`] |
| **educationUserReference** | [EducationUserReference](EducationUserReference.md) | educationUser to be added as member | |

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


## createSchool

> EducationSchool createSchool(educationSchool)

Add new school

### Example

```ts
import {
  Configuration,
  EducationSchoolApi,
} from '';
import type { CreateSchoolRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new EducationSchoolApi(config);

  const body = {
    // EducationSchool | New school
    educationSchool: ...,
  } satisfies CreateSchoolRequest;

  try {
    const data = await api.createSchool(body);
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
| **educationSchool** | [EducationSchool](EducationSchool.md) | New school | |

### Return type

[**EducationSchool**](EducationSchool.md)

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


## deleteClassFromSchool

> deleteClassFromSchool(schoolId, classId)

Unassign class from a school

### Example

```ts
import {
  Configuration,
  EducationSchoolApi,
} from '';
import type { DeleteClassFromSchoolRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new EducationSchoolApi(config);

  const body = {
    // string | key: id or schoolNumber of school
    schoolId: 43b879c4-14c6-4e0a-9b3f-b1b33c5a4bd4,
    // string | key: id or externalId of the class to unassign from school
    classId: 7e84a069-f374-479b-817d-71590117d443,
  } satisfies DeleteClassFromSchoolRequest;

  try {
    const data = await api.deleteClassFromSchool(body);
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
| **schoolId** | `string` | key: id or schoolNumber of school | [Defaults to `undefined`] |
| **classId** | `string` | key: id or externalId of the class to unassign from school | [Defaults to `undefined`] |

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


## deleteSchool

> deleteSchool(schoolId)

Delete school

Deletes a school. A school can only be delete if it has the terminationDate property set. And if that termination Date is in the past.

### Example

```ts
import {
  Configuration,
  EducationSchoolApi,
} from '';
import type { DeleteSchoolRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new EducationSchoolApi(config);

  const body = {
    // string | key: id or schoolNumber of school
    schoolId: 43b879c4-14c6-4e0a-9b3f-b1b33c5a4bd4,
  } satisfies DeleteSchoolRequest;

  try {
    const data = await api.deleteSchool(body);
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
| **schoolId** | `string` | key: id or schoolNumber of school | [Defaults to `undefined`] |

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


## deleteUserFromSchool

> deleteUserFromSchool(schoolId, userId)

Unassign user from a school

### Example

```ts
import {
  Configuration,
  EducationSchoolApi,
} from '';
import type { DeleteUserFromSchoolRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new EducationSchoolApi(config);

  const body = {
    // string | key: id or schoolNumber of school
    schoolId: 43b879c4-14c6-4e0a-9b3f-b1b33c5a4bd4,
    // string | key: id or username of the user to unassign from school
    userId: 90eedea1-dea1-90ee-a1de-ee90a1deee90,
  } satisfies DeleteUserFromSchoolRequest;

  try {
    const data = await api.deleteUserFromSchool(body);
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
| **schoolId** | `string` | key: id or schoolNumber of school | [Defaults to `undefined`] |
| **userId** | `string` | key: id or username of the user to unassign from school | [Defaults to `undefined`] |

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


## getSchool

> EducationSchool getSchool(schoolId)

Get the properties of a specific school

### Example

```ts
import {
  Configuration,
  EducationSchoolApi,
} from '';
import type { GetSchoolRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new EducationSchoolApi(config);

  const body = {
    // string | key: id or schoolNumber of school
    schoolId: 43b879c4-14c6-4e0a-9b3f-b1b33c5a4bd4,
  } satisfies GetSchoolRequest;

  try {
    const data = await api.getSchool(body);
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
| **schoolId** | `string` | key: id or schoolNumber of school | [Defaults to `undefined`] |

### Return type

[**EducationSchool**](EducationSchool.md)

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


## listSchoolClasses

> CollectionOfEducationClass listSchoolClasses(schoolId)

Get the educationClass resources owned by an educationSchool

### Example

```ts
import {
  Configuration,
  EducationSchoolApi,
} from '';
import type { ListSchoolClassesRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new EducationSchoolApi(config);

  const body = {
    // string | key: id or schoolNumber of school
    schoolId: 43b879c4-14c6-4e0a-9b3f-b1b33c5a4bd4,
  } satisfies ListSchoolClassesRequest;

  try {
    const data = await api.listSchoolClasses(body);
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
| **schoolId** | `string` | key: id or schoolNumber of school | [Defaults to `undefined`] |

### Return type

[**CollectionOfEducationClass**](CollectionOfEducationClass.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Retrieved classes |  -  |
| **0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## listSchoolUsers

> CollectionOfEducationUser listSchoolUsers(schoolId)

Get the educationUser resources associated with an educationSchool

### Example

```ts
import {
  Configuration,
  EducationSchoolApi,
} from '';
import type { ListSchoolUsersRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new EducationSchoolApi(config);

  const body = {
    // string | key: id or schoolNumber of school
    schoolId: 43b879c4-14c6-4e0a-9b3f-b1b33c5a4bd4,
  } satisfies ListSchoolUsersRequest;

  try {
    const data = await api.listSchoolUsers(body);
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
| **schoolId** | `string` | key: id or schoolNumber of school | [Defaults to `undefined`] |

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
| **200** | Retrieved educationUser |  -  |
| **0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## listSchools

> CollectionOfSchools listSchools()

Get a list of schools and their properties

### Example

```ts
import {
  Configuration,
  EducationSchoolApi,
} from '';
import type { ListSchoolsRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new EducationSchoolApi(config);

  try {
    const data = await api.listSchools();
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

[**CollectionOfSchools**](CollectionOfSchools.md)

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


## updateSchool

> EducationSchool updateSchool(schoolId, educationSchool)

Update properties of a school

### Example

```ts
import {
  Configuration,
  EducationSchoolApi,
} from '';
import type { UpdateSchoolRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new EducationSchoolApi(config);

  const body = {
    // string | key: id or schoolNumber of school
    schoolId: 43b879c4-14c6-4e0a-9b3f-b1b33c5a4bd4,
    // EducationSchool | New property values
    educationSchool: ...,
  } satisfies UpdateSchoolRequest;

  try {
    const data = await api.updateSchool(body);
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
| **schoolId** | `string` | key: id or schoolNumber of school | [Defaults to `undefined`] |
| **educationSchool** | [EducationSchool](EducationSchool.md) | New property values | |

### Return type

[**EducationSchool**](EducationSchool.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Success |  -  |
| **0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

