---
title: "ocis.messages.accounts.v0"
url: /grpc_apis/ocis_messages_accounts_v0
date: 2022-03-29T07:14:36Z
weight: 50
geekdocRepo: https://github.com/owncloud/ocis
---

{{< toc >}}



## ocis/messages/accounts/v0/accounts.proto

### Account

Account follows the properties of the ms graph api user resource.
See https://docs.microsoft.com/en-us/graph/api/resources/user?view=graph-rest-1.0#properties

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  | The unique identifier for the user. Key. Not nullable. Non reassignable. Read-only. |
| account_enabled | [bool](#bool) |  | `true` if the account is enabled; otherwise, `false`. This property is required when a user is created. Supports $filter. |
| is_resource_account | [bool](#bool) |  | A resource account is also known as a /disabled user object/ in Azure AD, and can be used to represent resources in general.<br>In Exchange it might be used to represent conference rooms, for example, and allow them to have a phone number.<br>You could give printers or machines with a sync client resource accounts as well.<br>A resource account can be homed in Microsoft 365 or on premises using Skype for Business Server 2019.<br>`true` if the user is a resource account; otherwise, `false`. Null value should be considered false. |
| creation_type | [string](#string) |  | Indicates whether the account was created as<br>- a regular school or work account ("" / emptystring),<br>- a local account, fully managed by oCIS (LocalAccount), includes synced accounts or<br>- an external account (Invitation),<br>- self-service sign-up using email verification (EmailVerified). Read-only. |
| identities | [Identities](#identities) | repeated | Represents the identities that can be used to sign in to this account.<br>An identity can be provided by oCIS (also known as a local account), by organizations, or by social identity providers such as Facebook, Google, and Microsoft, and is tied to an account.<br>May contain multiple items with the same signInType value. Supports $filter. |
| display_name | [string](#string) |  | The name displayed in the address book for the account.<br>This is usually the combination of the user's first name, middle initial and last name.<br>This property is required when a user is created and it cannot be cleared during updates.<br>Supports $filter and $orderby.<br>posixaccount MUST cn |
| preferred_name | [string](#string) |  | The username<br>posixaccount MUST uid |
| uid_number | [int64](#int64) |  | used for exposing the user using ldap<br>posixaccount MUST uidnumber |
| gid_number | [int64](#int64) |  | used for exposing the user using ldap<br>posixaccount MUST gidnumber |
| mail | [string](#string) |  | The SMTP address for the user, for example, "jeff@contoso.onmicrosoft.com". Read-Only. Supports $filter.<br>inetorgperson MAY mail |
| description | [string](#string) |  | A description, useful for resource accounts<br>posixaccount MAY description |
| password_profile | [PasswordProfile](#passwordprofile) |  | Specifies the password profile for the user.<br>The profile contains the user’s password. This property is required when a user is created.<br>The password in the profile must satisfy minimum requirements as specified by the passwordPolicies property.<br>By default, a strong password is required.<br>posixaccount MAY authPassword |
| memberOf | [Group](#group) | repeated | The groups, directory roles and administrative units that the user is a member of. Read-only. Nullable.<br>should we only respond with repeated strings of ids? no clients should a proper filter mask! |
| created_date_time | [google.protobuf.Timestamp](#googleprotobuftimestamp) |  | The created date of the account object. |
| deleted_date_time | [google.protobuf.Timestamp](#googleprotobuftimestamp) |  | The date and time the user was deleted. Returned only on $select. |
| on_premises_sync_enabled | [bool](#bool) |  | `true` if this object is synced from an on-premises directory;<br>`false` if this object was originally synced from an on-premises directory but is no longer synced;<br>null if this object has never been synced from an on-premises directory (default). Read-only |
| on_premises_immutable_id | [string](#string) |  | This property is used to associate an on-premises LDAP user to the oCIS account object.<br>This property must be specified when creating a new user account in the Graph if you are using a federated domain for the user’s userPrincipalName (UPN) property.<br>Important: The $ and _ characters cannot be used when specifying this property. Supports $filter. |
| on_premises_security_identifier | [string](#string) |  | Contains the on-premises security identifier (SID) for the user that was synchronized from on-premises to the cloud. Read-only. |
| on_premises_distinguished_name | [string](#string) |  | Contains the on-premises LDAP `distinguished name` or `DN`.<br>The property is only populated for customers who are synchronizing their on-premises directory to ocis-accounts. Read-only. |
| on_premises_sam_account_name | [string](#string) |  | Contains the on-premises `samAccountName` synchronized from the on-premises directory.<br>The property is only populated for customers who are synchronizing their on-premises directory to ocis-accounts. Read-only. |
| on_premises_domain_name | [string](#string) |  | Contains the on-premises `domainFQDN`, also called `dnsDomainName` synchronized from the on-premises directory<br>The property is only populated for customers who are synchronizing their on-premises directory to ocis-accounts. Read-only. |
| on_premises_user_principal_name | [string](#string) |  | Contains the on-premises userPrincipalName synchronized from the on-premises directory.<br>The property is only populated for customers who are synchronizing their on-premises directory to ocis-accounts. Read-only. |
| on_premises_last_sync_date_time | [google.protobuf.Timestamp](#googleprotobuftimestamp) |  | Indicates the last time at which the object was synced with the on-premises directory; Read-only. |
| on_premises_provisioning_errors | [OnPremisesProvisioningError](#onpremisesprovisioningerror) | repeated | Errors when using synchronization during provisioning. |
| external_user_state | [string](#string) |  | For an external user invited to the tenant using the invitation API, this property represents the invited user's invitation status.<br>For invited users, the state can be `PendingAcceptance` or `Accepted`, or "" / emptystring for all other users.<br>Returned only on $select. Supports $filter with the supported values. For example: $filter=externalUserState eq 'PendingAcceptance'. |
| external_user_state_change_date_time | [google.protobuf.Timestamp](#googleprotobuftimestamp) |  | Shows the timestamp for the latest change to the externalUserState property. Returned only on $select. |
| refresh_tokens_valid_from_date_time | [google.protobuf.Timestamp](#googleprotobuftimestamp) |  | Any refresh tokens or sessions tokens (session cookies) issued before this time are invalid, and applications will get<br>an error when using an invalid refresh or sessions token to acquire a delegated access token (to access APIs such as Microsoft Graph).<br>If this happens, the application will need to acquire a new refresh token by making a request to the authorize endpoint.<br>Returned only on $select. Read-only. Use invalidateAllRefreshTokens to reset. |
| sign_in_sessions_valid_from_date_time | [google.protobuf.Timestamp](#googleprotobuftimestamp) |  | Any refresh tokens or sessions tokens (session cookies) issued before this time are invalid, and applications will get<br>an error when using an invalid refresh or sessions token to acquire a delegated access token (to access APIs such as Microsoft Graph).<br>If this happens, the application will need to acquire a new refresh token by making a request to the authorize endpoint.<br>Read-only. Use revokeSignInSessions to reset. |

### Group



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  | The unique identifier for the group.<br>Returned by default. Inherited from directoryObject. Key. Not nullable. Read-only. |
| display_name | [string](#string) |  | The display name for the group. This property is required when a group is created and cannot be cleared during updates.<br>Returned by default. Supports $filter and $orderby.<br>groupofnames MUST cn<br><br>groupofnames MUST/MAY member |
| members | [Account](#account) | repeated | Users, contacts, and groups that are members of this group. HTTP Methods: GET (supported for all groups), POST (supported for security groups and mail-enabled security groups), DELETE (supported only for security groups) Read-only. Nullable. |
| owners | [Account](#account) | repeated | groupofnames MAY businessCategory<br>groupofnames MAY o<br>groupofnames MAY ou<br>groupofnames MAY owner, SINGLE-VALUE but there might be multiple owners |
| description | [string](#string) |  | An optional description for the group. Returned by default.<br>groupofnames MAY description |
| gid_number | [int64](#int64) |  | used for exposing the user using ldap<br>posixgroup MUST gidnumber<br><br>posixgroup MAY authPassword<br>posixgroup MAY userPassword<br>posixgroup MAY memberUid -> groupofnames member<br>posixgroup MAY description  -> groupofnames |
| created_date_time | [google.protobuf.Timestamp](#googleprotobuftimestamp) |  | Timestamp of when the group was created. The value cannot be modified and is automatically populated when the group is created<br>Returned by default. Read-only. |
| deleted_date_time | [google.protobuf.Timestamp](#googleprotobuftimestamp) |  | For some Azure Active Directory objects (user, group, application), if the object is deleted, it is first logically deleted, and this property is updated with the date and time when the object was deleted. Otherwise this property is null. If the object is restored, this property is updated to null.<br>Returned by default. Read-only. |
| expiration_date_time | [google.protobuf.Timestamp](#googleprotobuftimestamp) |  | Timestamp of when the group is set to expire. The value cannot be modified and is automatically populated when the group is created.<br>Returned by default. Read-only. |
| hide_from_address_lists | [bool](#bool) |  | True if the group is not displayed in certain parts of the Outlook user interface:<br>in the Address Book, in address lists for selecting message recipients, and in the Browse Groups dialog for searching groups; false otherwise. Default value is false.<br>Returned only on $select. |
| visibility | [string](#string) |  | Specifies the visibility of an Office 365 group. Possible values are: Private, Public, or Hiddenmembership; blank values are treated as public. See group visibility options to learn more.<br>Visibility can be set only when a group is created; it is not editable.<br>Returned by default. |
| on_premises_sync_enabled | [bool](#bool) |  | `true` if this group is synced from an on-premises directory;<br>`false` if this group was originally synced from an on-premises directory but is no longer synced;<br>null if this object has never been synced from an on-premises directory (default).<br>Returned by default. Read-only. Supports $filter. |
| on_premises_immutable_id | [string](#string) |  | This property is used to associate an on-premises LDAP user to the oCIS account object.<br>This property must be specified when creating a new user account in the Graph if you are using a federated domain for the user’s userPrincipalName (UPN) property.<br>Important: The $ and _ characters cannot be used when specifying this property. Supports $filter. |
| on_premises_security_identifier | [string](#string) |  | Contains the on-premises security identifier (SID) for the group that was synchronized from on-premises to the cloud. Returned by default. Read-only. |
| on_premises_distinguished_name | [string](#string) |  | Contains the on-premises LDAP `distinguished name` or `DN`.<br>The property is only populated for customers who are synchronizing their on-premises directory to ocis-accounts. Read-only. |
| on_premises_sam_account_name | [string](#string) |  | Contains the on-premises `samAccountName` synchronized from the on-premises directory.<br>The property is only populated for customers who are synchronizing their on-premises directory to ocis-accounts. Returned by default. Read-only. |
| on_premises_domain_name | [string](#string) |  | Contains the on-premises domain FQDN, also called dnsDomainName synchronized from the on-premises directory. The property is only populated for customers who are synchronizing their on-premises directory to Azure Active Directory via Azure AD Connect.<br>Returned by default. Read-only. |
| on_premises_net_bios_name | [string](#string) |  | Contains the on-premises netBios name synchronized from the on-premises directory. The property is only populated for customers who are synchronizing their on-premises directory to Azure Active Directory via Azure AD Connect.<br>Returned by default. Read-only. |
| on_premises_last_sync_date_time | [string](#string) |  | Indicates the last time at which the group was synced with the on-premises directory.<br>Returned by default. Read-only. Supports $filter. |
| on_premises_provisioning_errors | [OnPremisesProvisioningError](#onpremisesprovisioningerror) | repeated | Errors when using synchronization during provisioning. |

### Identities

Identities Represents an identity used to sign in to a user account.
An identity can be provided by oCIS, by organizations, or by social identity providers such as Facebook, Google, or Microsoft, that are tied to a user account.
This enables the user to sign in to the user account with any of those associated identities.
They are also used to keep a history of old usernames.

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| sign_in_type | [string](#string) |  | Specifies the user sign-in types in your directory, such as `emailAddress`, `userName` or `federated`.<br>Here, federated represents a unique identifier for a user from an issuer, that can be in any format chosen by the issuer.<br>Additional validation is enforced on *issuer_assigned_id* when the sign-in type is set to `emailAddress` or `userName`.<br>This property can also be set to any custom string. |
| issuer | [string](#string) |  | Specifies the issuer of the identity, for example facebook.com.<br>For local accounts (where signInType is not federated), this property is<br>the local B2C tenant default domain name, for example contoso.onmicrosoft.com.<br>For external users from other Azure AD organization, this will be the domain of<br>the federated organization, for example contoso.com.<br>Supports $filter. 512 character limit. |
| issuer_assigned_id | [string](#string) |  | Specifies the unique identifier assigned to the user by the issuer. The combination of *issuer* and *issuerAssignedId* must be unique within the organization. Represents the sign-in name for the user, when signInType is set to emailAddress or userName (also known as local accounts).<br>When *signInType* is set to:<br>* `emailAddress`, (or starts with `emailAddress` like `emailAddress1`) `issuerAssignedId` must be a valid email address<br>* `userName`, issuer_assigned_id must be a valid local part of an email address<br>Supports $filter. 512 character limit. |

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
| password | [string](#string) |  | The password for the user. This property is required when a user is created.<br>It can be updated, but the user will be required to change the password on the next login.<br>The password must satisfy minimum requirements as specified by the user’s passwordPolicies property. By default, a strong password is required. |
| last_password_change_date_time | [google.protobuf.Timestamp](#googleprotobuftimestamp) |  | The time when this account last changed their password. |
| password_policies | [string](#string) | repeated | Specifies password policies for the user.<br>This value is an enumeration with one possible value being “DisableStrongPassword”, which allows weaker passwords than the default policy to be specified.<br>“DisablePasswordExpiration” can also be specified. |
| force_change_password_next_sign_in | [bool](#bool) |  | `true` if the user must change her password on the next login; otherwise false. |
| force_change_password_next_sign_in_with_mfa | [bool](#bool) |  | If `true`, at next sign-in, the user must perform a multi-factor authentication (MFA) before being forced to change their password. The behavior is identical to forceChangePasswordNextSignIn except that the user is required to first perform a multi-factor authentication before password change. After a password change, this property will be automatically reset to false. If not set, default is false. |


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

