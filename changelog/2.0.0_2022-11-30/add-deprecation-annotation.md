Enhancement: Add deprecation annotation

We have added the ability to annotate variables in case of deprecations:

Example:

`services/nats/pkg/config/config.go`

```
Host string `yaml:"host" env:"NATS_HOST_ADDRESS,NATS_NATS_HOST" desc:"Bind address." deprecationVersion:"1.6.2" removalVersion:"1.7.5" deprecationInfo:"the name is ugly" deprecationReplacement:"NATS_HOST_ADDRESS"`
```

https://github.com/owncloud/ocis/issues/3917
https://github.com/owncloud/ocis/pull/5143
