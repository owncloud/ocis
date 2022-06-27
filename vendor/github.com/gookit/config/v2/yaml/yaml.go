/*
Package yaml is a driver use YAML format content as config source

Usage please see example:

*/
package yaml

// see https://pkg.go.dev/gopkg.in/yaml.v2
import (
	"github.com/gookit/config/v2"
	"gopkg.in/yaml.v2"
)

// Decoder the yaml content decoder
var Decoder config.Decoder = yaml.Unmarshal

// Encoder the yaml content encoder
var Encoder config.Encoder = yaml.Marshal

// Driver for yaml
var Driver = config.NewDriver(config.Yaml, Decoder, Encoder)
