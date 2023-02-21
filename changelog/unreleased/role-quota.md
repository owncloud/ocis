Enhancement: Added option to configure default quota per role

Admins can assign default quotas to users with certain roles by adding the following config to the `proxy.yaml`.
E.g.:
```
role_quotas:
    d7beeea8-8ff4-406b-8fb6-ab2dd81e6b11: 2300000
```

It maps a role ID to the quota in bytes.

https://github.com/owncloud/ocis/pull/5616
