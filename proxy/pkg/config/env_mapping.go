package config

import (
	"fmt"
	"reflect"
	"strings"

	gofig "github.com/gookit/config/v2"
)

// mappings holds a record of how to get an env variable's value onto a config.Config value. Field selectors are made via
// the `tagName` field. For instance having the following value:
//	type example struct {
//	  enable `mapstructure:"enable"`
//	}
//
//  e := example{enable: false}
//
// we can link the field `e.enable` with the environment variable EXTENSION_ENABLE by adding an entry in this mappings:
//	{
//		gType:   "bool",
//		envName: "EXTENSION_ENABLE",
//		tagName: "enable",
//	}
//
// so when a config is parsed the value is read from the environment, parsed and loaded onto whatever destination
// has the tagName.
var mappings = []struct {
	gType   string // expected type, used for decoding. It is the type expected from gookit.
	envName string // name of the env var
	tagName string // name of the tag to select the value from. Tag names are to be unique.
}{
	{
		gType:   "bool",
		envName: "PROXY_ENABLE_BASIC_AUTH",
		tagName: "enable_basic_auth",
	},
}

// GetEnv fetches a list of known env variables for this extension.
func GetEnv() []string {
	var r []string
	for i := range mappings {
		r = append(r, mappings[i].envName)
	}

	return r
}

func UnmapEnv(gooconf *gofig.Config, cfg *Config) error {
	for i := range mappings {
		switch mappings[i].gType {
		case "bool":
			v := gooconf.Bool(mappings[i].envName)
			if err := setField(cfg, mappings[i].tagName, v); err != nil {
				return err
			}
		case "string":
			v := gooconf.String(mappings[i].envName)
			if err := setField(cfg, mappings[i].tagName, v); err != nil {
				return err
			}
		default:
			return fmt.Errorf("invalid type for env var: `%v`", mappings[i].envName)
		}
	}

	return nil
}

// setField allows us to set a value on a struct selecting by its `mapstructure` tag.
func setField(item interface{}, fieldName string, value interface{}) error {
	v := reflect.ValueOf(item).Elem()
	if !v.CanAddr() {
		return fmt.Errorf("cannot assign to the item passed, item must be a pointer in order to assign")
	}
	fName := func(t reflect.StructTag) (string, error) {
		if jt, ok := t.Lookup("mapstructure"); ok {
			return strings.Split(jt, ",")[0], nil
		}
		return "", fmt.Errorf("tag %s provided does not define a json tag", fieldName)
	}

	fieldNames := map[string]int{}
	for i := 0; i < v.NumField(); i++ {
		typeField := v.Type().Field(i)
		tag := typeField.Tag
		jName, _ := fName(tag)
		fieldNames[jName] = i
	}

	fieldNum, ok := fieldNames[fieldName]
	if !ok {
		return fmt.Errorf("field does not exist within the provided item")
	}
	fieldVal := v.Field(fieldNum)
	fieldVal.Set(reflect.ValueOf(value))
	return nil
}
