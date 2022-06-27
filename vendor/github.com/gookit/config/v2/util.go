package config

import (
	"os"
	"reflect"
	"strings"

	"github.com/gookit/goutil/envutil"
	"github.com/mitchellh/mapstructure"
)

// ParseEnvVarStringHookFunc returns a DecodeHookFunc that parse ENV var
func ParseEnvVarStringHookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}

		str := envutil.ParseEnvValue(data.(string))
		return str, nil
	}
}

/*************************************************************
 * Deprecated methods
 *************************************************************/

// SetDecoder add/set a format decoder
//
// Deprecated: please use driver instead
func SetDecoder(format string, decoder Decoder) {
	dc.SetDecoder(format, decoder)
}

// SetDecoder set decoder
//
// Deprecated: please use driver instead
func (c *Config) SetDecoder(format string, decoder Decoder) {
	format = fixFormat(format)
	c.decoders[format] = decoder
}

// SetDecoders set decoders
//
// Deprecated: please use driver instead
func (c *Config) SetDecoders(decoders map[string]Decoder) {
	for format, decoder := range decoders {
		c.SetDecoder(format, decoder)
	}
}

// SetEncoder set a encoder for the format
//
// Deprecated: please use driver instead
func SetEncoder(format string, encoder Encoder) {
	dc.SetEncoder(format, encoder)
}

// SetEncoder set a encoder for the format
//
// Deprecated: please use driver instead
func (c *Config) SetEncoder(format string, encoder Encoder) {
	format = fixFormat(format)
	c.encoders[format] = encoder
}

// SetEncoders set encoders
//
// Deprecated: please use driver instead
func (c *Config) SetEncoders(encoders map[string]Encoder) {
	for format, encoder := range encoders {
		c.SetEncoder(format, encoder)
	}
}

/*************************************************************
 * helper methods/functions
 *************************************************************/

// GetEnv get os ENV value by name
func GetEnv(name string, defVal ...string) (val string) {
	return Getenv(name, defVal...)
}

// Getenv get os ENV value by name. like os.Getenv, but support default value
// Notice:
// - Key is not case-sensitive when getting
func Getenv(name string, defVal ...string) (val string) {
	if val = os.Getenv(name); val != "" {
		return
	}

	if len(defVal) > 0 {
		val = defVal[0]
	}
	return
}

func parseVarNameAndType(key string) (string, string) {
	typ := "string"
	key = strings.Trim(key, "-")

	// can set var type: int, uint, bool
	if strings.IndexByte(key, ':') > 0 {
		list := strings.SplitN(key, ":", 2)
		key, typ = list[0], list[1]

		if _, ok := validTypes[typ]; !ok {
			typ = "string"
		}
	}
	return key, typ
}

// format key
func formatKey(key, sep string) string {
	return strings.Trim(strings.TrimSpace(key), sep)
}

// resolve fix inc/conf/yaml format
func fixFormat(f string) string {
	if f == Yml {
		f = Yaml
	}

	if f == "inc" {
		f = Ini
	}

	// eg nginx config file.
	if f == "conf" {
		f = Hcl
	}
	return f
}
