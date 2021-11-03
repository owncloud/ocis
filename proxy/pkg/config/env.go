package config

import (
	"fmt"

	gofig "github.com/gookit/config/v2"
)

type mapping struct {
	goType      string      // expected type, used for decoding. It is the field dynamic type.
	env         string      // name of the env var.
	destination interface{} // memory address of the original config value to modify.
}

// GetEnv fetches a list of known env variables for this extension. It is to be used by gookit, as it provides a list
// with all the environment variables an extension supports.
func GetEnv() []string {
	var r = make([]string, len(structMappings(&Config{})))
	for i := range structMappings(&Config{}) {
		r = append(r, structMappings(&Config{})[i].env)
	}

	return r
}

// UnmapEnv loads values from the gooconf.Config argument and sets them in the expected destination.
func (c *Config) UnmapEnv(gooconf *gofig.Config) error {
	vals := structMappings(c)
	for i := range vals {
		// we need to guard against v != "" because this is the condition that checks that the value is set from the environment.
		// the `ok` guard is not enough, apparently.
		if v, ok := gooconf.GetValue(vals[i].env); ok && v != "" {
			switch vals[i].goType {
			case "bool":
				r := gooconf.Bool(vals[i].env)
				*vals[i].destination.(*bool) = r
			case "string":
				r := gooconf.String(vals[i].env)
				*vals[i].destination.(*string) = r
			case "int":
				r := gooconf.Int(vals[i].env)
				*vals[i].destination.(*int) = r
			case "float":
				// defaults to float64
				r := gooconf.Float(vals[i].env)
				*vals[i].destination.(*float64) = r
			default:
				// it is unlikely we will ever get here. Let this serve more as a runtime check for when debugging.
				return fmt.Errorf("invalid type for env var: `%v`", vals[i].env)
			}
		}
	}

	return nil
}
