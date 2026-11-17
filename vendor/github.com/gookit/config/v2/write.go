package config

import (
	"errors"
	"strings"

	"github.com/gookit/goutil/maputil"
)

// some common errors definitions
var (
	ErrReadonly   = errors.New("the config instance in 'readonly' mode")
	ErrKeyIsEmpty = errors.New("the config key is cannot be empty")
	ErrNotFound   = errors.New("this key does not exist in the config")
)

// SetData for override the Config.Data
func SetData(data map[string]any) {
	dc.SetData(data)
}

// SetData for override the Config.Data
func (c *Config) SetData(data map[string]any) {
	_ = c.writeAndHook(OnSetData, func() error {
		c.data = data
		return nil
	})
}

// Set value by key. setByPath default is true
func Set(key string, val any, setByPath ...bool) error {
	return dc.Set(key, val, setByPath...)
}

// Set a value by key string. setByPath default is true
func (c *Config) Set(key string, val any, setByPath ...bool) (err error) {
	if c.opts.Readonly {
		return ErrReadonly
	}

	sep := c.opts.Delimiter
	if key = formatKey(key, string(sep)); key == "" {
		return ErrKeyIsEmpty
	}

	// tip: the lock order is hookMu -> lock, see Config.writeAndHook().
	return c.writeAndHook(OnSetValue, func() error {
		return c.setValue(key, val, sep, setByPath)
	})
}

// setValue set the value to data by key. the c.lock must be locked by the caller.
func (c *Config) setValue(key string, val any, sep byte, setByPath []bool) error {
	if c.data == nil {
		c.data = make(map[string]any)
	}

	if strings.IndexByte(key, sep) == -1 {
		c.data[key] = val
		return nil
	}

	// disable set by path.
	if len(setByPath) > 0 && !setByPath[0] {
		c.data[key] = val
		return nil
	}

	// set by path
	keys := strings.Split(key, string(sep))
	return maputil.SetByKeys(&c.data, keys, val)
}
