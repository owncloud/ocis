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
   * [Backporting](#backporting)

## Introduction

`docs/helpers` contains small go programs creating docs by extracting information from the code. The `main.go` program is manually started with `make docs-generate` or via the CI. It calls the other required programs and has these main responsibilities:

- Generate docs for envvars in config structs including deprecations if there are any.
- Extract and generate docs for `extended` envvars that are not mentioned in config structs (aka "rogue" envvars).
- Extract and generate docs for `global` envvars which occur in multiple services.
- Create `docs/service/<service-name>/_index.md` from `service/<service-name>/README.md` files while keeping the existing `_index.md` if the README.md has not been created so far.

## Output Generated

- The generated yaml files can be found at: `docs/services/_includes` when running locally respectively in the `docs branch` after the CI has finished.
- The generated adoc files can be found at: `docs/services/_includes/adoc` when running locally respectively in the `docs branch` after the CI has finished.
- The file name for global envvars is named: `global_configvars.adoc`.
- The file name for extended envvars is named: `extended_configvars.adoc`.

## Admin Doc Process

Whenever a build from the [ocis admin](https://github.com/owncloud/docs-ocis) or any other related documentation is triggered, the files generated here are included into the build process and added in a proper manner defined by the admin documentation. The updated documentation will then show up on the [web](https://doc.owncloud.com/ocis/next/).

## Branching

The following is valid for envvars and yaml files related to the doc process:

* When filing a pull request in the ocis master branch relating to docs, CI runs `make docs-generate` and copies the result into the `docs` branch of ocis. This branch is then taken as base for owncloud.dev and as reference for the [admin docs](https://doc.owncloud.com/ocis/next/).
* When running `make docs-generate` locally, the same output is created as above but it stays in the same branch where the make command was issued.

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

## Extended Envvars

### General Extended Envvars Info

"Extended" envvars are variables that need to be present *before* the core or services are starting up as they depend on the info provided like path for config files etc. Therefore they are _not_ bound to services like other envvars.

It can happen that extended envvars are found but do not need to be published as they are for internal use only. Those envvars can be defined to be ignored for further processing.

IMPORTANT:

- Once an extended envvar has been identified, it is added to the `extended_vars.yaml` file found in this folder but never changed or touched by the process anymore. There is one exception with respect to single/double quote usage. While you can (and will) manually define a text like: `"'/var/lib/ocis'"`, quotes are transformed by the process in the .yaml file to: `'''/var/lib/ocis'''`. There is no need to change this back, as the final step transforms this correctly for the adoc table.

- Because extended envvars do not have the same structural setup as "normal" envvars (like type, description or defaults), this info needs to be provided manually once - even if found multiple times. Any change of this info will be noticed during the next CI run, the corresponding adoc file generated, changes transported to the docs branch and published in the next admin docs build.

- The identification if an envvar is in the yaml file already present is made via the `rawname` and the `path` identifier which includes the line number. If there is a change in the source file shifting line numbers, new items will get added and the old ones not touched. Though technically ok, this can cause confusion to identify which items have a correct path reference. To get rid of items with wrong line numbers, correct the existing ones, especially the one containing the description and which is marked to be shown. Only items that have a real line number match need to be present, orphaned items can safely be deleted. You can double-check valid items by creating a dummy branch, delete the `extended_vars.yaml` and run `make docs-generate` to regenerate the file having only items with valid path references.

- Do not change the sort order of extended envvar blocks as they are automatically reordered alphabetically.

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

IMPORTANT: `RawName`, `Path` and `FoundInCode` are automatically filled by the program. DO NOT EDIT THESE VALUES MANUALLY.

### Generate Extended Envvar Docs

The process further picks up the `yaml` file generated in the `Extract Rogue Envvars` step and renders it to an adoc file (a table is created) using a go template. The template file for this step can be found at `docs/templates/ADOC_extended.tmpl`.

## Backporting

The ocis repo contains branches which are necessary for the documentation. The `docs` branch is related to changes in master, necessary for owncloud.dev and the admin docs referencing master content when it comes to envvars and yaml files.

When a new stable ocis release (branch) is published, like `stable-2.0`, an additional branch (including CI) is set up manually by the dev team for referencing docs content like `docs-stable-2.0` - related to envvars and yaml files only - and added to the CI.

In case it is necessary to transport a change from master to a stable branch like `docs-stable-2.0`, you must backport the original changes that will create that file to the `stable-2.0` branch. The CI will then take care of creating the results in the target `docs-stable-2.0`.

Cases for a backport can be a typo in an envvar description you want to have fixed in a stable branch too or a file  was created after the stable branch was set up but needs to be available in that branch.
