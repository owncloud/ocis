---
title: "Deprecating Variables"
date: 2022-11-29T15:41:00+01:00
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/development
geekdocFilePath: deprecating-variables.md
---

{{< toc >}}

## Deprecating Environment Variables

Sometimes it is necessary to deprecate environment to align their naming with conventions. We therefore added annotations to automate the documentation process. It is necessary to know when the variable is going to be deprecated, when it is going to be removed and why.

### Example

```golang
// Notifications defines the config options for the notifications service.
type Notifications struct {
RevaGateway string `yaml:"reva_gateway" env:"OCIS_REVA_GATEWAY;REVA_GATEWAY" desc:"CS3 gateway used to look up user metadata" deprecationVersion:"3.0" removalVersion:"4.0.0" deprecationInfo:"REVA_GATEWAY changing name for consistency" deprecationReplacement:"OCIS_REVA_GATEWAY"`
...
}
```

There are four different annotation variables that need to be filled:

| Annotation |Description| Format|
|---|---|---|
| deprecationVersion | The version the variable will be deprecated | semver (e.g. 3.0)|
| removalVersion| The version the variable will be removed from the codebase. Note that according to semver, a removal **MUST NOT** be made in a minor or patch version change, but only in a major release | semver (e.g. 4.0.0) |
| deprecationInfo | Information why the variable is deprecated, must start with the name of the variable in order to avoid confusion, when there are multiple options in the `env:`-field | string (e.g. NATS_NATS_HOST is confusing) |
| deprecationReplacement | The name of the variable that is going to replace the deprecated one.| string (e.g. NATS_HOST_ADDRESS) |

### What happens next?

Once a variable has been finally removed, the annotations must be removed again from the code, since they do not serve any purpose from this point.
