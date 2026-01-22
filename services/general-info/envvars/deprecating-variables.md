---
title: "Envvar Deprecation"
date: 2024-08-22T15:41:00+01:00
weight: 15
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/services/general-info/envvars
geekdocFilePath: deprecating-variables.md
---

{{< toc >}}

## Deprecating Environment Variables

Sometimes it is necessary to deprecate an environment variable to align the naming with conventions or remove it completely. We therefore added annotations to automate the *documentation* process.

The relevant annotations in the envvar struct tag are:

* `deprecationVersion`\
  The release an envvar is announced for deprecation.
* `removalVersion`\
  The version it is finally going to be removed is defined via the mandatory placeholder `%%NEXT_PRODUCTION_VERSION%%`, not an actual version number.
* `deprecationInfo`\
  The reason why it was deprecated.
* `deprecationReplacement`\
  Only if it is going to be replaced, not necessary if removed.

{{< hint warning >}}
* During the development cycle, the value for the `removalVersion` must be set to `%%NEXT_PRODUCTION_VERSION%%`. This placeholder will be replaced by the real semantic-version number during the production releasing process.
* Compared when introducing new envvars where you can use arbitrary alphabetic identifyers, the string for deprecation is fixed and cannot be altered.
{{< /hint >}}

For the documentation to show the correct value for the `removalVersion`, our docs helper scripts will automatically generate the correct version to be printed in the documentation. If `%%NEXT_PRODUCTION_VERSION%%` is found in the query, it will be replaced with `next-prod`, else the value found is used.

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
| removalVersion | The version the variable will be removed from the codebase. Consider semver rules when finally removing a deprecated ennvar | `%%NEXT_PRODUCTION_VERSION%%` |
| deprecationInfo | Information why the variable is deprecated, must start with the name of the variable in order to avoid confusion, when there are multiple options in the `env:`-field | string (e.g. NATS_NATS_HOST is confusing) |
| deprecationReplacement | The name of the variable that is going to replace the deprecated one.| string (e.g. NATS_HOST_ADDRESS) |

### What Happens Next?

To remove an environment variable, which needs to be done before the planned release has codefreeze:

* Check if the envvar is also present in the REVA code and adapt if so.
* The envvar needs to be removed from any occurrences in the ocis code.
* The envvar must also be removed from the `docs/helpers/env_vars.yaml` file. This should be done in the same PR that removes the envvar from the code.
* Notify docs. The added/deprecated/removed envvar tables need to be created/updated. The release notes will get a note about this change.
