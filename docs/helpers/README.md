# Docs Helpers

   * [Introduction](#introduction)
   * [Output Generated](#output)
   * [Admin Doc Process](#admin-doc-process)
   * [Branching](#branching)
   * [Service-Dependent Output](#service-dependent-output)
      * [Generate Envvar Docs for Config Structs](#generate-envvar-docs-for-config-structs)
      * [Deprecation Process](#deprecation-process)
   * [Global Envvars](#global-envvars)
   * [Extended Envvars](#extended-envvars)
      * [General Extended Envvars Info](#general-extended-envvars-info)
      * [Extract Extended Envvars](#extract-extended-envvars)
      * [Generate Extended Envvar Docs](#generate-extended-envvar-docs)
   * [Tasks for New Releases](#tasks-for-new-releases)
   * [Backporting](#backporting)

## Introduction

`docs/helpers` contains a go program named `main.go` which creates docs by extracting information from the code using additional go programs. Individual steps (programs) can be called manually if needed. Note that not all programs are called automatically on purpose, see the [Tasks for New Releases](#tasks-for-new-releases) below. `main.go` is used by `make docs-generate` (or `make -C docs docs-generate` when running manually from the repos root) which is triggered by the CI or can be called manually. It calls the other required programs and has these main responsibilities for automatic runs:

- Generate docs for envvars in config structs including deprecations if there are any.
- Extract and generate docs for `extended` envvars that are not mentioned in config structs (aka "rogue" envvars).
- Extract and generate docs for `global` envvars which occur in multiple services.
- Create `docs/service/<service-name>/_index.md` from `service/<service-name>/README.md` files while keeping the existing `_index.md` if the service README.md has not been created so far. Also see the important note at [docs README](../README.md).

## Output Generated

- The generated yaml files can be found at: `docs/services/_includes` when running locally respectively in the `docs branch` after the CI has finished.
- The generated adoc files can be found at: `docs/services/_includes/adoc` when running locally respectively in the `docs branch` after the CI has finished.
- The file name for global envvars is named: `global_configvars.adoc`.
- The file name for extended envvars is named: `extended_configvars.adoc`.
- A file named `docs/helpers/env_vars.yaml` containing envvar changes gets updated if changes have been identified.
- A file named `docs/helpers/extended_vars.yaml` containing changes for extended envvars gets updated if changes have been identified. Note, if changes appear, **this file needs manual treatment** before committing, see [Extended Envvars](#extended-envvars) below.

## Admin Doc Process

Whenever a build from the [ocis admin](https://github.com/owncloud/docs-ocis) documentation or any other admin related documentation is triggered, files generated here in the ocis repo are included into the build process and added in a proper manner defined by the admin documentation. The updated documentation will then show up on the public [admin documentation](https://doc.owncloud.com/ocis/next/).

## Branching

The following is valid for envvars and yaml files related to the doc process:

* When filing a pull request in the ocis master branch relating to docs, CI runs `make docs-generate` and copies the result into the `docs` branch of ocis. This branch is then taken as base for owncloud.dev and as reference for the [admin docs](https://doc.owncloud.com/ocis/next/).
* When running `make docs-generate` _locally_, the same output is created as above but it stays in the same branch where the make command was issued.

In both cases, `make docs-generate` removes files in the target folder `_includes` to avoid remnants. All content is recreated.

On a side note (unrelated to the `docs` branch), [deployment examples](https://github.com/owncloud/ocis/tree/master/deployments/examples) have their own branch related to an ocis stable version to keep the state consistent, which is necessary for the admin documentation.

## Service-Dependent Output

For each service available, a file named like `<service name>_configvars.adoc` is created containing a:

* table on top defining deprecated envvars - if applicable
* table containing all envvars with their name, type, default value and description

The table with deprecations is always printed in the final adoc file even if there are none, but is rendered in the docs build process only if the `HasDeprecations` value is set. This value is automatically handed over via the adoc file. The template file can be found at `docs/templates/ADOC.tmpl`.

### Generate Envvar Docs for Config Structs

Generates docs from a template file, mainly extracting `"env"` and `"desc"` tags from the config structs.

Templates can be found in `docs/helpers` folder. (Same as this `README`.) Check `.tmpl` files

### Deprecation Process

For details on deprecation see the [deprecating-variables](https://github.com/owncloud/ocis/blob/master/docs/ocis/development/deprecating-variables.md) documentation.

## Global Envvars

Global envvars are gathered by checking if the envvar is available in more than one service. The table created is similar to the service-dependent envvar table but additionally contains a column with all service names where this envvar occurs. The output is rendered in list form where each item is clickable and automatically points to the corresponding service page. The template file can be found at `docs/templates/ADOC_global.tmpl`.

If global envvars do not appear in the list of globals, before checking if the code works, do a manual search in the ocis/services folder with `grep -rn OCIS_xxx` if the envvar in question appears at least twice. If the envvar only appears once, the helpers code works correct.

## Extended Envvars

### General Extended Envvars Info

"Extended" envvars are variables that need to be present *before* the core or services are starting up as they depend on the info provided like path for config files etc. Therefore they are _not_ bound to services like other envvars. Extended envvars are identified via `os.Getenv`, usually defined via a subfolder of `ocis-pkg`. The real envvar name name cant be automatically assigned and needs to be manually defined via the code in the `extended_vars.yaml` file.

It can happen that extended envvars are found but do not need to be published as they are for internal use only. Those envvars can be defined to be ignored for further processing.

**IMPORTANT:**

- **First Time Identification**\
  Once an extended envvar has been identified, it is added to the `extended_vars.yaml` file found, but never changed or touched by the process anymore. There is one exception with respect to single/double quote usage. While you can (and will) manually define a text like: `"'/var/lib/ocis'"`, quotes are transformed by the process in the .yaml file to: `'''/var/lib/ocis'''`. There is no need to change this back, as the final step transforms this correctly for the adoc table.

- **Item Naming**\
  An extended envvar may not have the right naming. It may appear as `name: _registryEnv`. In case, this envvar needs to be named properly like `name: MICRO_REGISTRY` which can only be done in close alignment with development.

- **Item Uniqueness**\
  The identification, if an envvar is already present in the yaml file, is made via the `rawname` and the `path` identifier which includes the line number. **If there is a change in the source file shifting line numbers, new items will get added and old ones do not get touched.** Though technically ok, this can cause confusion to identify which items are correctly present or just added additionally just be cause code location has changed. If there are multiple occurrences of the same `rawname` value, check which one contains relevant data and set `do_ignore` to `false` and all others to `true`. When there are two identical blocks with different source references, mostly the one containing a proper `default_value` is the active one. Populate the false block with the envvar data to be used.

- **Fixing Items**\
  If an item has been identified as additionally added because there was a change in the code location, it is mostly sufficient to just fix the line number in the `path` key of the existing/correct one and double check by removing the newly added ones. Then, re-run `make docs-generate`. If the fix was correct, no new items of the same will re-appear.

- **Remove Orphaned Items**\
  To get rid of items with wrong line numbers, check `rawname` the `path` and correct the _existing ones_, especially the one containing the description and which is marked `do_ignore` false. Only items that have a real line number match need to be present, orphaned items can safely be removed. You can double-check valid items by creating a dummy branch, delete the `extended_vars.yaml` and run `make docs-generate` to regenerate the file having only items with valid path references. With that info, you can remove orphaned items from the live file. Note to be careful on judging only on `foundincode` set to false indicating an item not existing anymore. Fix all items first, when rerunning `make docs-generate`, this may change back to true!

- **Sort Ordering**\
  Do not change the sort order of extended envvar blocks as they are automatically reordered alphabetically.

- **Mandatory Key Values**\
  Because extended envvars do not have the same structural setup as "normal" envvars (like type, description or defaults), this info needs to be provided manually once - for each valid block. Any change of this info will be noticed during the next CI run, the corresponding adoc file generated, changes transported to the docs branch and published in the next admin docs build. See the following example with all keys listed and populated:
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

### Extract Extended Envvars

The grep command parses the code, looking for `os.Getenv` and passes these contents to a yaml file along with the following information:
```golang
// Variable contains all information about one rogue envvar
type Variable struct {
	// These field structs are automatically filled:
	// RawName can be the name of the envvar or the name of its var
	RawName string `yaml:"rawname"`
	// Path to the envvar with linenumber
	Path string `yaml:"path"`
	// FoundInCode indicates if the variable is still found in the codebase.
	FoundInCode bool `yaml:"foundincode"`
	// Name is equal to RawName but will not be overwritten in consecutive runs
	Name string `yaml:"name"`

	// These field structs need manual filling:
	// Type of the envvar
	Type string `yaml:"type"`
	// DefaultValue of the envvar
	DefaultValue string `yaml:"default_value"`
	// Description of what this envvar does
	Description string `yaml:"description"`
	// Do not export this envvar into the generated adoc table
	Ignore bool `yaml:"do_ignore"`

	// For simplicity ignored for now:
	// DependendServices []Service `yaml:"dependend_services"`
}
```

This yaml file can later be manually edited to add descriptions, default values, etc.

**IMPORTANT**: `RawName`, `Path` and `FoundInCode` are automatically filled by the program. DO NOT EDIT THESE VALUES MANUALLY.

### Generate Extended Envvar Docs

The process further picks up the `yaml` file generated in the `Extract Rogue Envvars` step and renders it to an adoc file (a table is created) using a go template. The template file for this step can be found at `docs/templates/ADOC_extended.tmpl`.

## Tasks for New Releases

Close before a new release gets published, and **before** a new ocis admin docs branch is created, some tasks need to be done manually. These tasks cant be done automatically! By processing these manual tasks, doc related and important files get updated or created, ready to be consumed by the admin docs. The following situations can occur, see the solution provided:

1. Admin docs prepares in master for envvar delta changes. Building the admin docs report non-existent referenced delta files.\
&#8594; Run the manual tasks described below and restart the branched admin docs.

2. A new admin docs branch gets created but no `docs-stable-x.y` branch has been created in the ocis repo so far. Building the branched admin docs reports non-existent referenced files.\
&#8594; Run the manual tasks described below.\
&#8594; Take care that the `docs-stable-x.y` branch gets created in the ocis repo.\
&#8594; Alternatively you can _temporarily_ point in the branched [docs-ocis/antora.yml](https://github.com/owncloud/docs-ocis/blob/master/antora.yml) `compose_url_component` to `master`.\
&#8594; Restart building the branched admin docs.

3. A new ocis release is close to being released, the manual tasks **have not been** processed so far.\
&#8594; Run the manual tasks AND backport the results into the `docs-stable-x.y` branch.

4. A new ocis release is close to being released, the manual tasks **have been** processed some time ago.\
&#8594; Re-run the manual tasks AND backport the results into the `docs-stable-x.y` branch.

5. A new ocis release **has been released**, the manual tasks **have not been** processed so far.\
&#8594; Re-run the manual tasks AND backport the results into the `docs-stable-x.y` branch.
  
### Task List

The following can be done at any time but *latest* it must be done before a new release gets published.

1. From the ocis root, run `make -C docs docs-generate` and check if there is a change in the `extended-envars.yaml` output. In this case, process [Extended Envvars](#extended-envvars). When done, re-run the make command and check if the output of `./docs/services/_includes/adoc/extended_configvars.adoc` matches the expectations.

2. In `./docs/helpers` run: `go run . --help` This will give you an overview of available commands. 
    1. Run `go run . all` to generate and/or update required base files.\
   &#8594; Any change in `env_vars.yaml` must be check out to be available in master.
    2. Create delta files for added, removed and deprecated envvars. To do so type:\
    `go run . env-var-delta-table` and use as parameter the versions you want to compare. Example: `v5.0.0 v6.0.0`.
    3. List and check the files created in `./docs/helpers/output/env-deltas/`. These are the files defining a table in markdown and asciidoc where the former is used in the dev docs and the latter is consumed by the admin docs build process.\
    **Important**: To make the files consumable, the following attribute must be changed in the respective admin docs respective branch definition: [docs-ocis/antora.yml](https://github.com/owncloud/docs-ocis/blob/master/antora.yml) `env_var_delta_name`. Only the delta filename part needs to be used. The complete filenames will be assembled automatically.
       ```
       v5.0.0-v6.0.0 -->
       
       v5.0.0-v6.0.0-added.md
       v5.0.0-v6.0.0-deprecated.md
       v5.0.0-v6.0.0-removed.md
       ```

3. Commit all changes, create a PR and post merging, adapt, if necessary, the ocis admin docs. With any merge in the admin docs, the newly created content for that ocis admin docs branch will be consumed.    

## Backporting

The ocis repo contains branches which are necessary for the documentation. The `docs` branch is related to changes in master, necessary for owncloud.dev and the admin docs referencing master content when it comes to envvars and yaml files.

Cases for a backport can be a typo in an envvar description you want to have fixed in a stable branch too or a file  was created after the stable branch was set up but needs to be available in that branch.

When a new stable ocis release (branch) is published, like `stable-5.0`, an additional branch (including CI) is set up manually by the dev team for referencing docs content like `docs-stable-5.0` - related to envvars and yaml files only - and added to the CI.

In case it is necessary to transport a change from master to a stable branch like `docs-stable-5.0`, you must backport the original changes that will create that file to the `stable-5.0` branch. The CI will then take care of creating the results in the target `docs-stable-5.0`.

If the change is expected to have a bigger impact on documenation, you can locally run `make -C docs docs-generate` in the respective branch containing the changes or independently in the `stable-x.y` branch after merging to see if there are additional actions necessary and changed files may need to get checked in.
