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

To create the changed envvar tables, you must proceed the following steps in order:

1. Install, if not already done, the converter for adoc to markdown tables: `npm install -g downdoc`\
This is only required when converting adoc to markdown tables but highly recommended!

2. Run `make -C docs docs-generate` from the ocis root.\
A file named `env_vars.yaml` is generated. Check for validity. If issues are found, fix them in the service sources first which need to be merged before you rerun make. When the changes are fine, create a PR and merge it. Any delta information is based on an actual `env_vars.yaml` file which is pull from master by the python script described below!

3. Configure the Python script `docs/helpers/changed_envvars.py` variables for the new version.\
Note that you **must** use semver and not code names!

4. Run the python script from the ocis root such as `python3 docs/helpers/changed_envvars.py`.\
Note that the script pulls data from the master branch only, therefore the `env_vars.yaml` file MUST be up to date.\
adoc tables will be generated which are used for the admin docs and are the basis for markdown.

5. Because the script can not determine the link target in the `Service` column, you must manually adapt them in the generated adoc files with respect to the name and target. Only one entry per identical service block is required, delete the cell content for the rest for ease of readability.

6. Run `npx downdoc <filename.adoc>` for each of the newly generated `added`, `removed` and `deprecated` files.

7. Add in each markdown file on top the following sentence:\
`Note that the links provided in the service column are non functional when clicked.` including a newline.

8. Create a PR and merge it.
