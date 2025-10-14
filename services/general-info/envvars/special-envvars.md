---
title: "Special Envvars"
date: 2025-07-04T00:00:00+01:00
weight: 14
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/services/general-info/envvars
geekdocFilePath: special-envvars.md
---

{{< toc >}}

Handling these envvars properly is very important for the automated doc process for both the developer and the admin docs!

## Special Scope Envvars

Variables with special scope are only related to a deployment method such as `OCIS_RUN_SERVICES`. These variables cannot be gathered automatically, rarely change, can be viewed and must be maintained manually via the [admin documentation](https://doc.owncloud.com/ocis/next/deployment/services/env-vars-special-scope.html).

## Extended Envvars

Environment variables with extended scope are not included in a service. They are rarely added or removed, though their code location can change during development.

They are variables that must be present before the core or services start up because they depend on information such as the path to configuration files. Therefore, they are not bound to services like other environment variables. Extended environment variables are identified via `os.Getenv` and are usually defined via a subfolder of `ocis-pkg`. The real envvar name cannot be automatically assigned and must be defined manually via the code in the generated `extended_vars.yaml` file.

While generating the ocis Developer Documentation, the `extended_vars.yaml` file located in `/docs/helpers/` might be updated which needs special care to finalize the process. When merged, an `extended_configvars.adoc` file is internally generated from it. This file, along with the others, is moved to the `doc` branch. The process runs automatically, and no additional effort is required. The admin docs picks this .adoc file for further processing.

The file for the master (`docs`) branch is located at:\
`https://github.com/owncloud/ocis/tree/docs/services/_includes/adoc/extended_configvars.adoc`\
respectivle in any `docs-stable-x.y` branch.

### General Info

The process behind updating the `extended_vars.yaml` is non-destructive. This means, that the process checks the code if values found match those in the file already present. If differences occur, only **new** content blocks are added, independent if it is new or moved code. The file is recreated when  deleted - try to avoid this and maintain the changed one.

This also means, that if generating the docs result in a change in the `extended_vars.yaml` file, manual action **must** be taken and the final changes need to be committed/pushed/merged. If this is not done, the `extended_configvars.adoc` will contain invalid and/or corrupt data.

It can happen that extended envvars are found but do not need to be published as they are for internal use only. Those envvars can be defined to be ignored for further processing.

**IMPORTANT:**

- **First Time Identification**\
  Once an extended envvar has been identified, it is added to the `extended_vars.yaml` file found, but never changed or touched by the process anymore. There is one exception with respect to single/double quote usage. While you can (and will) manually define a text like: `"'/var/lib/ocis'"`, quotes are transformed by the process in the .yaml file to: `'''/var/lib/ocis'''`. There is no need to change this back, as the final step transforms this correctly for the adoc table.

- **Item Naming**\
  An extended envvar may not have the right naming. It may appear as `name: _registryEnv`. In case, this envvar needs to be named properly like `name: MICRO_REGISTRY` which can only be done in close alignment with development.

- **Automatic Data Population**:\
  `rawname`, `path` and `foundincode` are automatically filled by the program.\
  **IMPORTANT**: DO NOT EDIT THESE VALUES MANUALLY - except the line number in the `path` key.

- **Manual Data Population**:\
  The following keys can and must be updated manually: `name`, `type`, `default_value`, `description`, `do_ignore`\
  For the `path` key, **only** the line number of the value may be changed, see fixing values below.

- **Item Uniqueness**\
  The identification, if an envvar is already present in the yaml file, is made via the `rawname` and the `path` identifier which includes the line number. **If there is a change in the source file shifting line numbers, new items will get added and old ones do not get touched.** Though technically ok, this can cause confusion to identify which items are correctly present or just added additionally just be cause code location has changed. If there are multiple occurrences of the same `rawname` value, check which one contains relevant data and set `do_ignore` to `false` and all others to `true`. When there are two identical blocks with different source references, mostly the one containing a proper `default_value` is the active one. Populate the false block with the envvar data to be used.

- **Sort Ordering**\
  Do not change the sort order of extended envvar blocks as they are automatically reordered alphabetically.

- **Mandatory Key Values**\
  Because extended envvars do not have the same structural setup as "normal" envvars (like type, description or defaults), this info needs to be provided manually once - for each valid block. Any change of this info will be noticed during the next CI run, the corresponding `extended_configvars.adoc` file will be generated, changes moved to the docs branch and published in the next admin docs build. See the following example with all keys listed and populated:
    ```yaml
    rawname: registryAddressEnv
    path: ocis-pkg/registry/registry.go:44
    foundincode: true
    name: MICRO_REGISTRY_ADDRESS
    type: string
    default_value: ""
    description: The bind address of the internal go micro framework. Only change on
        supervision of ownCloud Support.
    do_ignore: false
    ```

### Fixing Changed Items

If there is a change in `extended_vars.yaml` which you can identify via git when running e.g. `make -C docs docs-generate`, read the [General Info](#general-info) section first and follow the items listed below afterwards.

- **Fixing Items**\
  If an item has been identified as additionally added such as there was a change in the code location only, it is mostly sufficient to just fix the line number in the `path` key of the existing/correct one and double check by removing the newly added ones. Then, re-run `make -C docs docs-generate`. If the fix was correct, no new items of the same will re-appear.

- **Remove Orphaned Items**\
  To get rid of items with wrong line numbers, check `rawname` the `path` and correct the _existing ones_, especially the one containing the description and which is marked `do_ignore` false. Only items that have a real line number match need to be present, orphaned items can safely be removed.

  You can double-check valid items by creating a dummy branch, delete the `extended_vars.yaml` and run `make -C docs docs-generate` to regenerate the file having only items with valid path references. With that info, you can remove orphaned items from the live file. Note to be careful on judging only on `foundincode` set to false indicating an item not existing anymore. Fix all items first, when rerunning `make -C docs docs-generate`, this may change back to true!

  I an envvar has been removed from the code completely, you can also remove the respective entry block from the file. 

- When all is set, create a PR and merge it.
