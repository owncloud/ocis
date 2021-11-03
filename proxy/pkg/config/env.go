package config

import (
	"fmt"
	"reflect"

	gofig "github.com/gookit/config/v2"
)

type mapping struct {
	env         []string    // name of the env var.
	destination interface{} // memory address of the original config value to modify.
}

// GetEnv fetches a list of known env variables for this extension. It is to be used by gookit, as it provides a list
// with all the environment variables an extension supports.
func GetEnv() []string {
	var r = make([]string, len(structMappings(&Config{})))
	for i := range structMappings(&Config{}) {
		r = append(r, structMappings(&Config{})[i].env...)
	}

	return r
}

// UnmapEnv loads values from the gooconf.Config argument and sets them in the expected destination.
func (c *Config) UnmapEnv(gooconf *gofig.Config) error {
	vals := structMappings(c)
	for i := range vals {
		for j := range vals[i].env {
			// we need to guard against v != "" because this is the condition that checks that the value is set from the environment.
			// the `ok` guard is not enough, apparently.
			if v, ok := gooconf.GetValue(vals[i].env[j]); ok && v != "" {

				// get the destination type from destination
				switch reflect.ValueOf(vals[i].destination).Type().String() {
				case "*bool":
					r := gooconf.Bool(vals[i].env[j])
					*vals[i].destination.(*bool) = r
				case "*string":
					r := gooconf.String(vals[i].env[j])
					*vals[i].destination.(*string) = r
				case "*int":
					r := gooconf.Int(vals[i].env[j])
					*vals[i].destination.(*int) = r
				case "*float64":
					// defaults to float64
					r := gooconf.Float(vals[i].env[j])
					*vals[i].destination.(*float64) = r
				default:
					// it is unlikely we will ever get here. Let this serve more as a runtime check for when debugging.
					return fmt.Errorf("invalid type for env var: `%v`", vals[i].env[j])
				}
			}
		}
	}

	return nil
}
