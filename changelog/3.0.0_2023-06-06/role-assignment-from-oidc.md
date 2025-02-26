Enhancement: Added possibility to assign roles based on OIDC claims

oCIS can now be configured to update a user's role assignment from the values of a claim provided
via the IDPs userinfo endpoint. The claim name and the mapping between claim values and ocis role
name can be configured via the configuration of the proxy service. Example:

```
role_assignment:
    driver: oidc
    oidc_role_mapper:
        role_claim: ocisRoles
        role_mapping:
            - role_name: admin
              claim_value: myAdminRole
            - role_name: spaceadmin
              claim_value: mySpaceAdminRole
            - role_name: user
              claim_value: myUserRole
            - role_name: guest
              claim_value: myGuestRole
```

https://github.com/owncloud/ocis/pull/6048
