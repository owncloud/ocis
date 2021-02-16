---
title: "GRPC API"
date: 2018-05-02T00:00:00+00:00
weight: 50
geekdocRepo: https://github.com/owncloud/ocis-thumbnails
geekdocEditPath: edit/master/docs
geekdocFilePath: grpc.md
---

{{< toc >}}

## accounts.proto

### Account

Account follows the properties of the ms graph api user resuorce.
See https://docs.microsoft.com/en-us/graph/api/resources/user?view=graph-rest-1.0#properties

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  | The unique identifier for the user. Key. Not nullable. Non reassignable. Read-only. |
| account_enabled | [bool](#bool) |  | true* if the account is enabled; otherwise, *false*. This property is required when a user is created. Supports $filter. |
| is_resource_account | [bool](#bool) |  | A resource account is also known as a /disabled user object/ in Azure AD, and can be used to represent resources in general. In Exchange it might be used to represent conference rooms, for example, and allow them to have a phone number. You could give printers or machines with a sync client resource accounts as well. A resource account can be homed in Microsoft 365 or on premises using Skype for Business Server 2019. *true* if the user is a resource account; otherwise, *false*. Null value should be considered false. |
| creation_type | [string](#string) |  | Indicates whether the account was created as - a regular school or work account ("" / emptystring), - a local account, fully managed by oCIS (LocalAccount), includes synced accounts or - an external account (Invitation), - self-service sign-up using email verification (EmailVerified). Read-only. |
| identities | [Identities](#identities) | repeated | Represents the identities that can be used to sign in to this account. An identity can be provided by oCIS (also known as a local account), by organizations, or by social identity providers such as Facebook, Google, and Microsoft, and is tied to an account. May contain multiple items with the same signInType value. Supports $filter. |
| display_name | [string](#string) |  | The name displayed in the address book for the account. This is usually the combination of the user's first name, middle initial and last name. This property is required when a user is created and it cannot be cleared during updates. Supports $filter and $orderby. posixaccount MUST cn |
| preferred_name | [string](#string) |  | The username posixaccount MUST uid |
| uid_number | [int64](#int64) |  | TODO rename to on_premise_? or move to extension? see https://docs.microsoft.com/en-us/graph/extensibility-open-users used for exposing the user using ldap posixaccount MUST uidnumber |
| gid_number | [int64](#int64) |  | used for exposing the user using ldap posixaccount MUST gidnumber |
| mail | [string](#string) |  | The SMTP address for the user, for example, "jeff@contoso.onmicrosoft.com". Read-Only. Supports $filter. inetorgperson MAY mail |
| description | [string](#string) |  | A description, useful for resource accounts posixaccount MAY description |
| password_profile | [PasswordProfile](#passwordprofile) |  | Specifies the password profile for the user. The profile contains the user’s password. This property is required when a user is created. The password in the profile must satisfy minimum requirements as specified by the passwordPolicies property. By default, a strong password is required. posixaccount MAY authPassword |
| memberOf | [Group](#group) | repeated | The groups, directory roles and administrative units that the user is a member of. Read-only. Nullable. should we only respond with repeated strings of ids? no clients should a proper filter mask! |
| created_date_time | [google.protobuf.Timestamp](#googleprotobuftimestamp) |  | The created date of the account object. |
| deleted_date_time | [google.protobuf.Timestamp](#googleprotobuftimestamp) |  | The date and time the user was deleted. Returned only on $select. |
| on_premises_sync_enabled | [bool](#bool) |  | true* if this object is synced from an on-premises directory; *false* if this object was originally synced from an on-premises directory but is no longer synced; null if this object has never been synced from an on-premises directory (default). Read-only |
| on_premises_immutable_id | [string](#string) |  | This property is used to associate an on-premises LDAP user to the oCIS account object. This property must be specified when creating a new user account in the Graph if you are using a federated domain for the user’s userPrincipalName (UPN) property. Important: The $ and _ characters cannot be used when specifying this property. Supports $filter. |
| on_premises_security_identifier | [string](#string) |  | Contains the on-premises security identifier (SID) for the user that was synchronized from on-premises to the cloud. Read-only. |
| on_premises_distinguished_name | [string](#string) |  | Contains the on-premises LDAP `distinguished name` or `DN`. The property is only populated for customers who are synchronizing their on-premises directory to ocis-accounts. Read-only. |
| on_premises_sam_account_name | [string](#string) |  | Contains the on-premises `samAccountName` synchronized from the on-premises directory. The property is only populated for customers who are synchronizing their on-premises directory to ocis-accounts. Read-only. |
| on_premises_domain_name | [string](#string) |  | Contains the on-premises `domainFQDN`, also called `dnsDomainName` synchronized from the on-premises directory The property is only populated for customers who are synchronizing their on-premises directory to ocis-accounts. Read-only. |
| on_premises_user_principal_name | [string](#string) |  | Contains the on-premises userPrincipalName synchronized from the on-premises directory. The property is only populated for customers who are synchronizing their on-premises directory to ocis-accounts. Read-only. |
| on_premises_last_sync_date_time | [google.protobuf.Timestamp](#googleprotobuftimestamp) |  | Indicates the last time at which the object was synced with the on-premises directory; Read-only. |
| on_premises_provisioning_errors | [OnPremisesProvisioningError](#onpremisesprovisioningerror) | repeated | Errors when using synchronization during provisioning. |
| external_user_state | [string](#string) |  | For an external user invited to the tenant using the invitation API, this property represents the invited user's invitation status. For invited users, the state can be `PendingAcceptance` or `Accepted`, or "" / emptystring for all other users. Returned only on $select. Supports $filter with the supported values. For example: $filter=externalUserState eq 'PendingAcceptance'. |
| external_user_state_change_date_time | [google.protobuf.Timestamp](#googleprotobuftimestamp) |  | Shows the timestamp for the latest change to the externalUserState property. Returned only on $select. |
| refresh_tokens_valid_from_date_time | [google.protobuf.Timestamp](#googleprotobuftimestamp) |  | Any refresh tokens or sessions tokens (session cookies) issued before this time are invalid, and applications will get an error when using an invalid refresh or sessions token to acquire a delegated access token (to access APIs such as Microsoft Graph). If this happens, the application will need to acquire a new refresh token by making a request to the authorize endpoint. Returned only on $select. Read-only. Use invalidateAllRefreshTokens to reset. |
| sign_in_sessions_valid_from_date_time | [google.protobuf.Timestamp](#googleprotobuftimestamp) |  | Any refresh tokens or sessions tokens (session cookies) issued before this time are invalid, and applications will get an error when using an invalid refresh or sessions token to acquire a delegated access token (to access APIs such as Microsoft Graph). If this happens, the application will need to acquire a new refresh token by making a request to the authorize endpoint. Read-only. Use revokeSignInSessions to reset. |

### AddMemberRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| group_id | [string](#string) |  | The id of the group to add a member to |
| account_id | [string](#string) |  | The account id to add |

### CreateAccountRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| account | [Account](#account) |  | The account resource to create |

### CreateGroupRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| group | [Group](#group) |  | The account resource to create |

### DeleteAccountRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |

### DeleteGroupRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |

### GetAccountRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |

### GetGroupRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |

### Group



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  | The unique identifier for the group. Returned by default. Inherited from directoryObject. Key. Not nullable. Read-only. |
| display_name | [string](#string) |  | The display name for the group. This property is required when a group is created and cannot be cleared during updates. Returned by default. Supports $filter and $orderby. groupofnames MUST cn

groupofnames MUST/MAY member |
| members | [Account](#account) | repeated | Users, contacts, and groups that are members of this group. HTTP Methods: GET (supported for all groups), POST (supported for security groups and mail-enabled security groups), DELETE (supported only for security groups) Read-only. Nullable. TODO accounts (users) only for now, we can add groups with the dedicated message using oneof construct later |
| owners | [Account](#account) | repeated | groupofnames MAY businessCategory groupofnames MAY o groupofnames MAY ou groupofnames MAY owner, SINGLE-VALUE but there might be multiple owners |
| description | [string](#string) |  | An optional description for the group. Returned by default. groupofnames MAY description |
| gid_number | [int64](#int64) |  | used for exposing the user using ldap posixgroup MUST gidnumber

posixgroup MAY authPassword posixgroup MAY userPassword posixgroup MAY memberUid -> groupofnames member posixgroup MAY description -> groupofnames |
| created_date_time | [google.protobuf.Timestamp](#googleprotobuftimestamp) |  | Timestamp of when the group was created. The value cannot be modified and is automatically populated when the group is created Returned by default. Read-only. |
| deleted_date_time | [google.protobuf.Timestamp](#googleprotobuftimestamp) |  | For some Azure Active Directory objects (user, group, application), if the object is deleted, it is first logically deleted, and this property is updated with the date and time when the object was deleted. Otherwise this property is null. If the object is restored, this property is updated to null. Returned by default. Read-only. |
| expiration_date_time | [google.protobuf.Timestamp](#googleprotobuftimestamp) |  | Timestamp of when the group is set to expire. The value cannot be modified and is automatically populated when the group is created. Returned by default. Read-only. |
| hide_from_address_lists | [bool](#bool) |  | True if the group is not displayed in certain parts of the Outlook user interface: in the Address Book, in address lists for selecting message recipients, and in the Browse Groups dialog for searching groups; false otherwise. Default value is false. Returned only on $select. |
| visibility | [string](#string) |  | Specifies the visibility of an Office 365 group. Possible values are: Private, Public, or Hiddenmembership; blank values are treated as public. See group visibility options to learn more. Visibility can be set only when a group is created; it is not editable. Returned by default. |
| on_premises_sync_enabled | [bool](#bool) |  | true* if this group is synced from an on-premises directory; *false* if this group was originally synced from an on-premises directory but is no longer synced; null if this object has never been synced from an on-premises directory (default). Returned by default. Read-only. Supports $filter. |
| on_premises_immutable_id | [string](#string) |  | This property is used to associate an on-premises LDAP user to the oCIS account object. This property must be specified when creating a new user account in the Graph if you are using a federated domain for the user’s userPrincipalName (UPN) property. Important: The $ and _ characters cannot be used when specifying this property. Supports $filter. |
| on_premises_security_identifier | [string](#string) |  | Contains the on-premises security identifier (SID) for the group that was synchronized from on-premises to the cloud. Returned by default. Read-only. |
| on_premises_distinguished_name | [string](#string) |  | Contains the on-premises LDAP `distinguished name` or `DN`. The property is only populated for customers who are synchronizing their on-premises directory to ocis-accounts. Read-only. |
| on_premises_sam_account_name | [string](#string) |  | Contains the on-premises `samAccountName` synchronized from the on-premises directory. The property is only populated for customers who are synchronizing their on-premises directory to ocis-accounts. Returned by default. Read-only. |
| on_premises_domain_name | [string](#string) |  | Contains the on-premises domain FQDN, also called dnsDomainName synchronized from the on-premises directory. The property is only populated for customers who are synchronizing their on-premises directory to Azure Active Directory via Azure AD Connect. Returned by default. Read-only. |
| on_premises_net_bios_name | [string](#string) |  | Contains the on-premises netBios name synchronized from the on-premises directory. The property is only populated for customers who are synchronizing their on-premises directory to Azure Active Directory via Azure AD Connect. Returned by default. Read-only. |
| on_premises_last_sync_date_time | [string](#string) |  | Indicates the last time at which the group was synced with the on-premises directory. Returned by default. Read-only. Supports $filter. |
| on_premises_provisioning_errors | [OnPremisesProvisioningError](#onpremisesprovisioningerror) | repeated | Errors when using synchronization during provisioning. |

### Identities

Identities Represents an identity used to sign in to a user account.
An identity can be provided by oCIS, by organizations, or by social identity providers such as Facebook, Google, or Microsoft, that are tied to a user account.
This enables the user to sign in to the user account with any of those associated identities.
They are also used to keep a history of old usernames.

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| sign_in_type | [string](#string) |  | Specifies the user sign-in types in your directory, such as `emailAddress`, `userName` or `federated`. Here, federated represents a unique identifier for a user from an issuer, that can be in any format chosen by the issuer. Additional validation is enforced on *issuer_assigned_id* when the sign-in type is set to `emailAddress` or `userName`. This property can also be set to any custom string. |
| issuer | [string](#string) |  | Specifies the issuer of the identity, for example facebook.com. For local accounts (where signInType is not federated), this property is the local B2C tenant default domain name, for example contoso.onmicrosoft.com. For external users from other Azure AD organization, this will be the domain of the federated organization, for example contoso.com. Supports $filter. 512 character limit. |
| issuer_assigned_id | [string](#string) |  | Specifies the unique identifier assigned to the user by the issuer. The combination of *issuer* and *issuerAssignedId* must be unique within the organization. Represents the sign-in name for the user, when signInType is set to emailAddress or userName (also known as local accounts). When *signInType* is set to: * `emailAddress`, (or starts with `emailAddress` like `emailAddress1`) *issuerAssignedId* must be a valid email address * `userName`, issuer_assigned_id must be a valid local part of an email address Supports $filter. 512 character limit. |

### ListAccountsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| page_size | [int32](#int32) |  | Optional. The maximum number of accounts to return in the response |
| page_token | [string](#string) |  | Optional. A pagination token returned from a previous call to `Get` that indicates from where search should continue |
| field_mask | [google.protobuf.FieldMask](#googleprotobuffieldmask) |  | Optional. Used to specify a subset of fields that should be returned by a get operation or modified by an update operation. |
| query | [string](#string) |  | Optional. Search criteria used to select the accounts to return. If no search criteria is specified then all accounts will be returned

TODO update query language Query expressions can be used to restrict results based upon the account properties where the operators `=`, `NOT`, `AND` and `OR` can be used along with the suffix wildcard symbol `*`.

The string properties in a query expression should use escaped quotes for values that include whitespace to prevent unexpected behavior.

Some example queries are:

* Query `display_name=Th*` returns accounts whose display_name starts with "Th" * Query `email=foo@example.com` returns accounts with `email` set to `foo@example.com` * Query `display_name=\\"Test String\\"` returns accounts with display names that include both "Test" and "String" |

### ListAccountsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| accounts | [Account](#account) | repeated | The field name should match the noun "accounts" in the method name. There will be a maximum number of items returned based on the page_size field in the request |
| next_page_token | [string](#string) |  | Token to retrieve the next page of results, or empty if there are no more results in the list |

### ListGroupsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| page_size | [int32](#int32) |  | Optional. The maximum number of groups to return in the response |
| page_token | [string](#string) |  | Optional. A pagination token returned from a previous call to `Get` that indicates from where search should continue |
| field_mask | [google.protobuf.FieldMask](#googleprotobuffieldmask) |  | Optional. Used to specify a subset of fields that should be returned by a get operation or modified by an update operation. |
| query | [string](#string) |  | Optional. Search criteria used to select the groups to return. If no search criteria is specified then all groups will be returned

TODO update query language Query expressions can be used to restrict results based upon the account properties where the operators `=`, `NOT`, `AND` and `OR` can be used along with the suffix wildcard symbol `*`.

The string properties in a query expression should use escaped quotes for values that include whitespace to prevent unexpected behavior.

Some example queries are:

* Query `display_name=Th*` returns accounts whose display_name starts with "Th" * Query `display_name=\\"Test String\\"` returns groups with display names that include both "Test" and "String" |

### ListGroupsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| groups | [Group](#group) | repeated | The field name should match the noun "group" in the method name. There will be a maximum number of items returned based on the page_size field in the request |
| next_page_token | [string](#string) |  | Token to retrieve the next page of results, or empty if there are no more results in the list |

### ListMembersRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| page_size | [int32](#int32) |  |  |
| page_token | [string](#string) |  | Optional. A pagination token returned from a previous call to `Get` that indicates from where search should continue |
| field_mask | [google.protobuf.FieldMask](#googleprotobuffieldmask) |  | Optional. Used to specify a subset of fields that should be returned by a get operation or modified by an update operation. |
| query | [string](#string) |  | Optional. Search criteria used to select the groups to return. If no search criteria is specified then all groups will be returned

TODO update query language Query expressions can be used to restrict results based upon the account properties where the operators `=`, `NOT`, `AND` and `OR` can be used along with the suffix wildcard symbol `*`.

The string properties in a query expression should use escaped quotes for values that include whitespace to prevent unexpected behavior.

Some example queries are:

* Query `display_name=Th*` returns accounts whose display_name starts with "Th" * Query `display_name=\\"Test String\\"` returns groups with display names that include both "Test" and "String" |
| id | [string](#string) |  | The id of the group to list members from |

### ListMembersResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| members | [Account](#account) | repeated | The field name should match the noun "members" in the method name. There will be a maximum number of items returned based on the page_size field in the request |
| next_page_token | [string](#string) |  | Token to retrieve the next page of results, or empty if there are no more results in the list |

### OnPremisesProvisioningError



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| category | [string](#string) |  | Category of the provisioning error. Note: Currently, there is only one possible value. Possible value: PropertyConflict - indicates a property value is not unique. Other objects contain the same value for the property. |
| occurred_date_time | [google.protobuf.Timestamp](#googleprotobuftimestamp) |  | The date and time at which the error occurred. |
| property_causing_error | [string](#string) |  | Name of the directory property causing the error. Current possible values: UserPrincipalName or ProxyAddress |
| value | [string](#string) |  | Value of the property causing the error. |

### PasswordProfile



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| password | [string](#string) |  | The password for the user. This property is required when a user is created. It can be updated, but the user will be required to change the password on the next login. The password must satisfy minimum requirements as specified by the user’s passwordPolicies property. By default, a strong password is required. |
| last_password_change_date_time | [google.protobuf.Timestamp](#googleprotobuftimestamp) |  | The time when this account last changed their password. |
| password_policies | [string](#string) | repeated | Specifies password policies for the user. This value is an enumeration with one possible value being “DisableStrongPassword”, which allows weaker passwords than the default policy to be specified. “DisablePasswordExpiration” can also be specified. |
| force_change_password_next_sign_in | [bool](#bool) |  | true* if the user must change her password on the next login; otherwise false. |
| force_change_password_next_sign_in_with_mfa | [bool](#bool) |  | If *true*, at next sign-in, the user must perform a multi-factor authentication (MFA) before being forced to change their password. The behavior is identical to forceChangePasswordNextSignIn except that the user is required to first perform a multi-factor authentication before password change. After a password change, this property will be automatically reset to false. If not set, default is false. |

### RebuildIndexRequest




### RebuildIndexResponse




### RemoveMemberRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| group_id | [string](#string) |  | The id of the group to remove a member from |
| account_id | [string](#string) |  | The account id to remove |

### UpdateAccountRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| account | [Account](#account) |  | The account resource which replaces the resource on the server |
| update_mask | [google.protobuf.FieldMask](#googleprotobuffieldmask) |  | The update mask applies to the resource. For the `FieldMask` definition, see https://developers.google.com/protocol-buffers/docs/reference/google.protobuf#fieldmask |

### UpdateGroupRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| group | [Group](#group) |  | The group resource which replaces the resource on the server |
| update_mask | [google.protobuf.FieldMask](#googleprotobuffieldmask) |  | The update mask applies to the resource. For the `FieldMask` definition, see https://developers.google.com/protocol-buffers/docs/reference/google.protobuf#fieldmask |


### AccountsService

Follow recommended Methods for rpc APIs https://cloud.google.com/apis/design/resources?hl=de#methods
https://cloud.google.com/apis/design/standard_methods?hl=de#list
https://cloud.google.com/apis/design/naming_convention?hl=de

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| ListAccounts | [ListAccountsRequest](#listaccountsrequest) | [ListAccountsResponse](#listaccountsresponse) | Lists accounts |
| GetAccount | [GetAccountRequest](#getaccountrequest) | [Account](#account) | Gets an account |
| CreateAccount | [CreateAccountRequest](#createaccountrequest) | [Account](#account) | Creates an account |
| UpdateAccount | [UpdateAccountRequest](#updateaccountrequest) | [Account](#account) | Updates an account |
| DeleteAccount | [DeleteAccountRequest](#deleteaccountrequest) | [.google.protobuf.Empty](#googleprotobufempty) | Deletes an account |

### GroupsService



| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| ListGroups | [ListGroupsRequest](#listgroupsrequest) | [ListGroupsResponse](#listgroupsresponse) | Lists groups |
| GetGroup | [GetGroupRequest](#getgrouprequest) | [Group](#group) | Gets an groups |
| CreateGroup | [CreateGroupRequest](#creategrouprequest) | [Group](#group) | Creates a group |
| UpdateGroup | [UpdateGroupRequest](#updategrouprequest) | [Group](#group) | Updates a group |
| DeleteGroup | [DeleteGroupRequest](#deletegrouprequest) | [.google.protobuf.Empty](#googleprotobufempty) | Deletes a group |
| AddMember | [AddMemberRequest](#addmemberrequest) | [Group](#group) | group:addmember https://docs.microsoft.com/en-us/graph/api/group-post-members?view=graph-rest-1.0&tabs=http |
| RemoveMember | [RemoveMemberRequest](#removememberrequest) | [Group](#group) | group:removemember https://docs.microsoft.com/en-us/graph/api/group-delete-members?view=graph-rest-1.0 |
| ListMembers | [ListMembersRequest](#listmembersrequest) | [ListMembersResponse](#listmembersresponse) | group:listmembers https://docs.microsoft.com/en-us/graph/api/group-list-members?view=graph-rest-1.0 |

### IndexService



| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| RebuildIndex | [RebuildIndexRequest](#rebuildindexrequest) | [RebuildIndexResponse](#rebuildindexresponse) |  |

## Scalar Value Types

| .proto Type | Notes | C++ | Java |
| ----------- | ----- | --- | ---- |
| {{< div id="double" content="double" >}} |  | double | double |
| {{< div id="float" content="float" >}} |  | float | float |
| {{< div id="int32" content="int32" >}} | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int |
| {{< div id="int64" content="int64" >}} | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long |
| {{< div id="uint32" content="uint32" >}} | Uses variable-length encoding. | uint32 | int |
| {{< div id="uint64" content="uint64" >}} | Uses variable-length encoding. | uint64 | long |
| {{< div id="sint32" content="sint32" >}} | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int |
| {{< div id="sint64" content="sint64" >}} | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long |
| {{< div id="fixed32" content="fixed32" >}} | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int |
| {{< div id="fixed64" content="fixed64" >}} | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long |
| {{< div id="sfixed32" content="sfixed32" >}} | Always four bytes. | int32 | int |
| {{< div id="sfixed64" content="sfixed64" >}} | Always eight bytes. | int64 | long |
| {{< div id="bool" content="bool" >}} |  | bool | boolean |
| {{< div id="string" content="string" >}} | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String |
| {{< div id="bytes" content="bytes" >}} | May contain any arbitrary sequence of bytes. | string | ByteString |
