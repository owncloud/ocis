# Docs Helpers

`docs/helpers` contains a small go program creating docs by extracting information from the code. It has three main responsibilities
- Generate docs for envvars in config structs
- Extract envvars that are not mentioned in config structs (aka "Rogue" envvars)
- Generate docs for rogue envvars

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
	// Ignore this envvar when creating docs?
	Ignore bool `yaml:"do_ignore"`

	// For simplicity ignored for now:
	// DependendServices []Service `yaml:"dependend_services"`
}
```
This yaml file can later be manually edited to add descriptions, default values, etc.

IMPORTANT: `RawName`, `Path` and `FoundInCode` are automatically filled by the program. DO NOT EDIT THESE VALUES MANUALLY

## Generate Rogue Envvar docs

It picks up the `yaml` file generated in `Extract Rogue Envvars` step and renders it to a adoc file using a go template.

Template for this can be found at `docs/templates/ADOC_global.tmpl`
