---
title: "Envvar Deprecation"
date: 2024-08-22T15:41:00+01:00
weight: 20
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/services/general-info
geekdocFilePath: deprecating-variables.md
---

{{< toc >}}

## Deprecating Environment Variables

Sometimes it is necessary to deprecate an environment variable to align the naming with conventions or remove it at all. We therefore added annotations to automate the *documentation* process.

The relevant annotations in the envvar struct tag are:

* `deprecationVersion`\
  The release an envvar is announced for deprecation.
* `removalVersion`\
  The version it is finally going to be removed is defined via the mandatory placeholder `%%NEXT_PRODUCTION_VERSION%%`, not an actual version number.
* `deprecationInfo`\
  The reason why it got deprecated.
* `deprecationReplacement`\
  Only if it is going to be replaced, not necessary if removed.

{{< hint warning >}}
During the development cycle, the value for the `removalVersion` must be set to `%%NEXT_PRODUCTION_VERSION%%`.  This placeholder will be removed by the real version number during the production releasing process.
{{< /hint >}}

For the documentation only to show the correct value for the `removalVersion`, our docs helper scripts will automatically generate the correct version to be printed in the documentation. If `%%NEXT_PRODUCTION_VERSION%%` is found in the query, it will be replaced with `next-prod`, else the value found is used.

### Example

```golang
// Notifications defines the config options for the notifications service.
type Notifications struct {
RevaGateway string `yaml:"reva_gateway" env:"OCIS_REVA_GATEWAY;REVA_GATEWAY" desc:"CS3 gateway used to look up user metadata" deprecationVersion:"3.0" removalVersion:"%%NEXT_PRODUCTION_VERSION%%" deprecationInfo:"REVA_GATEWAY changing name for consistency" deprecationReplacement:"OCIS_REVA_GATEWAY"`
...
}
```

There are four different annotation variables that need to be filled:

| Annotation |Description| Format|
|---|---|---|
| deprecationVersion | The version the variable will be deprecated | semver (e.g. 3.0)|
| removalVersion| The version the variable will be removed from the codebase. Note that according to semver, a removal **MUST NOT** be made in a minor or patch version change, but only in a major release | `%%NEXT_PRODUCTION_VERSION%%` |
| deprecationInfo | Information why the variable is deprecated, must start with the name of the variable in order to avoid confusion, when there are multiple options in the `env:`-field | string (e.g. NATS_NATS_HOST is confusing) |
| deprecationReplacement | The name of the variable that is going to replace the deprecated one.| string (e.g. NATS_HOST_ADDRESS) |

### What Happens Next?

Once a variable has been finally been removed, the annotations must be removed again from the code, since they don't serve any purpose.
