---
title: Deltas Between Versions
date: 2024-02-08T00:00:00+00:00
weight: 20
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/services/general-info/envvars/env-var-deltas
geekdocFilePath: _index.md
geekdocCollapseSection: true
---

## General Information

This section provides information about `added`, `removed` and `deprecated` environment variables between two major/minor versions.

{{< hint warning >}}
* When creating a new release, this step should be completed **before** the new `stable-x.y` branch is created. Then, all changes will go cleanly into this stable branch.
* If the changes required are not part of the stable branch created for the release process, you must backport all `added`, `removed` and `deprecated` files created from the described process below into the stable branch. Backporting `env_vars.yaml`to the stable branch is not required and can be omitted.
{{< /hint >}}

To create the changed envvar tables, you must proceed the following steps in order:

1. Install, if not already done, the converter for adoc to markdown tables: `npm install -g downdoc`\
This is only required when converting adoc to markdown tables but it is highly recommended to show them in the dev docs too!

1. Run `make docs-local` from the ocis root.\
Usually, a file named `env_vars.yaml` gets changed. Check for validity. If issues are found, fix them in the service sources first which need to be merged before you rerun make. For details how to do so, see [Maintain the 'env_vars.yaml' File]({{< ref "../new-release-process.md#maintain-the-env_varsyaml-file" >}}). Any delta information is based on an actual `env_vars.yaml` file which is pulled **from master** by the python script described below!

1. Configure the Python script `docs/helpers/changed_envvars.py` variables for the new version.\
Note that you **must** use semver and not code names!

1. Run the python script from the ocis root such as `python3 docs/helpers/changed_envvars.py`.\
Note that the script pulls data from the master branch as a base reference, therefore the `env_vars.yaml` file must be kept up to date. The adoc tables generated are used for the admin documentation and form the basis for Markdown.

1. As the script cannot determine the link target (xref:) in the `Service` column, you must adapt these manually in the generated adoc files according to the file name and printed name. Envvars starting with `OCIS_` are displayed differently in the `Service` column because the file name and printed name cannot be easily identified and must be resolved differently. You have to check where they are defined, unlike the others where the name provides a clue. The final xref path must be corrected manually in all cases. Only one entry per identical service is required to generate an easy block view. Delete the cell content except for the pipe symbol (`|`) to make it easier to read. See existing files for an example.

1. Change into the directory that contains the generated adoc files and run `npx downdoc <filename.adoc>` for each of the newly generated `added`, `removed` and `deprecated` files. This will generate markdown files for the dev docs.

1. Add in each markdown file on top the following sentence:\
`Note that the links provided in the service column are non functional when clicked.`, including a newline.

1. Create a PR and merge it.
