package config

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/mitchellh/mapstructure"
)

// MapStruct alias method of the 'Structure'
//
// Usage:
// 	dbInfo := &Db{}
// 	config.MapStruct("db", dbInfo)
func MapStruct(key string, dst interface{}) error { return dc.MapStruct(key, dst) }

// MapStruct alias method of the 'Structure'
func (c *Config) MapStruct(key string, dst interface{}) error {
	return c.Structure(key, dst)
}

// BindStruct alias method of the 'Structure'
func BindStruct(key string, dst interface{}) error { return dc.BindStruct(key, dst) }

// BindStruct alias method of the 'Structure'
func (c *Config) BindStruct(key string, dst interface{}) error {
	return c.Structure(key, dst)
}

// MapOnExists mapping data to the dst structure only on key exists.
func MapOnExists(key string, dst interface{}) error {
	return dc.MapOnExists(key, dst)
}

// MapOnExists mapping data to the dst structure only on key exists.
func (c *Config) MapOnExists(key string, dst interface{}) error {
	err := c.Structure(key, dst)
	if err != nil && err == errNotFound {
		return nil
	}

	return err
}

// Structure get config data and binding to the dst structure.
//
// Usage:
// 	dbInfo := Db{}
// 	config.Structure("db", &dbInfo)
func (c *Config) Structure(key string, dst interface{}) error {
	var data interface{}
	if key == "" { // binding all data
		data = c.data
	} else { // some data of the config
		var ok bool
		data, ok = c.GetValue(key)
		if !ok {
			return errNotFound
		}
	}

	var bindConf *mapstructure.DecoderConfig
	if c.opts.DecoderConfig == nil {
		bindConf = newDefaultDecoderConfig()
	} else {
		bindConf = c.opts.DecoderConfig
		// Compatible with previous settings opts.TagName
		if bindConf.TagName == "" {
			bindConf.TagName = c.opts.TagName
		}
	}

	// parse env var
	if c.opts.ParseEnv && bindConf.DecodeHook == nil {
		bindConf.DecodeHook = ParseEnvVarStringHookFunc()
	}

	bindConf.Result = dst // set result struct ptr
	decoder, err := mapstructure.NewDecoder(bindConf)
	if err != nil {
		return err
	}

	return decoder.Decode(data)
}

// ToJSON string
func (c *Config) ToJSON() string {
	buf := &bytes.Buffer{}

	_, err := c.DumpTo(buf, JSON)
	if err != nil {
		return ""
	}

	return buf.String()
}

// WriteTo a writer
func WriteTo(out io.Writer) (int64, error) { return dc.WriteTo(out) }

// WriteTo Write out config data representing the current state to a writer.
func (c *Config) WriteTo(out io.Writer) (n int64, err error) {
	return c.DumpTo(out, c.opts.DumpFormat)
}

// DumpTo a writer and use format
func DumpTo(out io.Writer, format string) (int64, error) { return dc.DumpTo(out, format) }

// DumpTo use the format(json,yaml,toml) dump config data to a writer
func (c *Config) DumpTo(out io.Writer, format string) (n int64, err error) {
	var ok bool
	var encoder Encoder

	format = fixFormat(format)
	if encoder, ok = c.encoders[format]; !ok {
		err = errors.New("not exists/register encoder for the format: " + format)
		return
	}

	// is empty
	if len(c.data) == 0 {
		return
	}

	// encode data to string
	encoded, err := encoder(c.data)
	if err != nil {
		return
	}

	// write content to out
	num, _ := fmt.Fprintln(out, string(encoded))

	return int64(num), nil
}
