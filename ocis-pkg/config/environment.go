package config

import (
	"fmt"
	"reflect"

	gofig "github.com/gookit/config/v2"
	"github.com/owncloud/ocis/ocis-pkg/shared"
)

// BindEnv takes a config `c` and a EnvBinding and binds the values from the environment to the address location in cfg.
func BindEnv(c *gofig.Config, bindings []shared.EnvBinding) error {
	return bindEnv(c, bindings)
}

func bindEnv(c *gofig.Config, bindings []shared.EnvBinding) error {
	for i := range bindings {
		for j := range bindings[i].EnvVars {
			// we need to guard against v != "" because this is the condition that checks that the value is set from the environment.
			// the `ok` guard is not enough, apparently.
			if v, ok := c.GetValue(bindings[i].EnvVars[j]); ok && v != "" {

				// get the destination type from destination
				switch reflect.ValueOf(bindings[i].Destination).Type().String() {
				case "*bool":
					r := c.Bool(bindings[i].EnvVars[j])
					*bindings[i].Destination.(*bool) = r
				case "*string":
					r := c.String(bindings[i].EnvVars[j])
					*bindings[i].Destination.(*string) = r
				case "*int":
					r := c.Int(bindings[i].EnvVars[j])
					*bindings[i].Destination.(*int) = r
				case "*float64":
					// defaults to float64
					r := c.Float(bindings[i].EnvVars[j])
					*bindings[i].Destination.(*float64) = r
				default:
					// it is unlikely we will ever get here. Let this serve more as a runtime check for when debugging.
					return fmt.Errorf("invalid type for env var: `%v`", bindings[i].EnvVars[j])
				}
			}
		}
	}

	return nil
}
