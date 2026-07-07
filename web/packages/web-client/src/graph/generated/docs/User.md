# User

Represents an Active Directory user object.

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **string** | Read-only. | [optional] [readonly] [default to undefined]
**accountEnabled** | **boolean** | Set to \&quot;true\&quot; when the account is enabled. | [optional] [default to undefined]
**appRoleAssignments** | [**Array&lt;AppRoleAssignment&gt;**](AppRoleAssignment.md) | The apps and app roles which this user has been assigned. | [optional] [readonly] [default to undefined]
**displayName** | **string** | The name displayed in the address book for the user. This value is usually the combination of the user\&#39;s first name, middle initial, and last name. This property is required when a user is created and it cannot be cleared during updates. Returned by default. Supports $orderby. | [default to undefined]
**drives** | [**Array&lt;Drive&gt;**](Drive.md) | A collection of drives available for this user. Read-only. | [optional] [readonly] [default to undefined]
**drive** | [**Drive**](Drive.md) |  | [optional] [default to undefined]
**identities** | [**Array&lt;ObjectIdentity&gt;**](ObjectIdentity.md) | Identities associated with this account. | [optional] [default to undefined]
**mail** | **string** | The SMTP address for the user, for example, \&#39;jeff@contoso.onowncloud.com\&#39;. Returned by default. | [optional] [default to undefined]
**memberOf** | [**Array&lt;Group&gt;**](Group.md) | Groups that this user is a member of. HTTP Methods: GET (supported for all groups). Read-only. Nullable. Supports $expand. | [optional] [readonly] [default to undefined]
**onPremisesSamAccountName** | **string** | Contains the on-premises SAM account name synchronized from the on-premises directory. | [default to undefined]
**passwordProfile** | [**PasswordProfile**](PasswordProfile.md) |  | [optional] [default to undefined]
**surname** | **string** | The user\&#39;s surname (family name or last name). Returned by default. | [optional] [default to undefined]
**givenName** | **string** | The user\&#39;s givenName. Returned by default. | [optional] [default to undefined]
**userType** | **string** | The user&#x60;s type. This can be either \&quot;Member\&quot; for regular user, \&quot;Guest\&quot; for guest users or \&quot;Federated\&quot; for users imported from a federated instance. | [optional] [readonly] [default to undefined]
**preferredLanguage** | **string** | Represents the users language setting, ISO-639-1 Code | [optional] [default to undefined]
**signInActivity** | [**SignInActivity**](SignInActivity.md) |  | [optional] [default to undefined]
**externalID** | **string** | A unique identifier assigned to the user by the organization. | [optional] [default to undefined]
**crossInstanceReference** | **string** | A unique reference to the user. This is used to query the user from a different oCIS instance connected to the same identity provider. | [optional] [default to undefined]
**instances** | [**Array&lt;Instance&gt;**](Instance.md) | oCIS instances that the user is either a member or a guest of. | [optional] [default to undefined]

## Example

```typescript
import { User } from './api';

const instance: User = {
    id,
    accountEnabled,
    appRoleAssignments,
    displayName,
    drives,
    drive,
    identities,
    mail,
    memberOf,
    onPremisesSamAccountName,
    passwordProfile,
    surname,
    givenName,
    userType,
    preferredLanguage,
    signInActivity,
    externalID,
    crossInstanceReference,
    instances,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
