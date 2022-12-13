# Docs Helpers

`docs/helpers` contains a small go program creating docs by extracting information from the code. It is manually started with `make docs-generate` or via the CI and has three main responsibilities:

- Generate docs for envvars in config structs
- Extract envvars that are not mentioned in config structs (aka "Rogue" envvars)
- Generate docs for rogue envvars

Output:

- The generated yaml files can be found at: `docs/services/_includes` when running locally respectively in the docs branch after the CI has finished.
- The generated adoc files can be found at: `docs/services/_includes/adoc` when running locally respectively in the docs branch after the CI has finished.
- The file name for rouge envvars is named: `global_configvars.adoc`.

Admin doc process:

Whenever a build from the admin documentation is triggered, the files generated here are included into the build process and added in a proper manner defined by the admin documentation.

Genreal info:

"Rouge" envvars are variables that need to be present *before* the core or services are starting up as they depend on the info provided like path for config files etc. Therefore they are _not_ bound to services like other envvars do. 

It can happen that rouge envvars are found but do not need to be published as they are for internal use only. Those rouge envvars can be defined to be ignored for further processing.

IMPORTANT:

- Once a rouge envvar has been identified, it is added to the `global_vars.yaml` file but never changed or touched by the process. There is one exception with respect to single/double quote usage. While you manually can (and will) define a text like: `"'/var/lib/ocis'"`, quotes are transformed by the process in the .yaml file to: `'''/var/lib/ocis'''`. There is no need to change this back, as the final step transforms this correctly to the adoc table.

- Because rouge envvars do not have the same structural setup as "normal" envvars like type, description or defaults, these infos need to be provided manually one time - even if found multiple times. Any change on this info will be used on the next CI run and published on the next admin docs build.

- Do not change the sort order of rouge envvar blocks as they are automatically reordered alphabetically.

## Generate Envvar docs for config structs

Generates docs from a template file, mainly extracting `"env"` and `"desc"` tags from the config structs.

Templates can be found in `docs/helpers` folder. (Same as this `README`.) Check `.tmpl` files

## Extract Rogue Envvars

It `grep`s over the code, looking for `os.Getenv` and parses these contents to a yaml file along with the following information:
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

## Generate Rogue Envvar docs

It picks up the `yaml` file generated in `Extract Rogue Envvars` step and renders it to a adoc file (table) using a go template.

The adoc template file for this step can be found at `docs/templates/ADOC_global.tmpl`.
