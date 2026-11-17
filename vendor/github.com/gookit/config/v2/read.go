package config

import (
	"maps"
	"math"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/gookit/goutil/envutil"
	"github.com/gookit/goutil/maputil"
	"github.com/gookit/goutil/mathutil"
	"github.com/gookit/goutil/strutil"
)

// Exists key exists check
func Exists(key string, findByPath ...bool) bool { return dc.Exists(key, findByPath...) }

// Exists key exists check
func (c *Config) Exists(key string, findByPath ...bool) (ok bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	sep := c.opts.Delimiter
	if key = formatKey(key, string(sep)); key == "" {
		return
	}

	if _, ok = c.data[key]; ok {
		return
	}

	// disable find by path.
	if len(findByPath) > 0 && !findByPath[0] {
		return
	}

	// has sub key? eg. "lang.dir"
	if strings.IndexByte(key, sep) == -1 {
		return
	}

	_, ok = maputil.GetByPathKeys(c.data, strings.Split(key, string(sep)))
	return
}

/*************************************************************
 * read config data
 *************************************************************/

// Data return all config data
func Data() map[string]any { return dc.Data() }

// Data get all config data.
//
// Note: will don't apply any options, like ParseEnv
func (c *Config) Data() map[string]any {
	return c.data
}

// Sub return a map config data by key
func Sub(key string) map[string]any { return dc.Sub(key) }

// Sub get a map config data by key
//
// Note: will don't apply any options, like ParseEnv
func (c *Config) Sub(key string) map[string]any {
	if mp, ok := c.GetValue(key); ok {
		if mmp, ok := mp.(map[string]any); ok {
			return mmp
		}
	}
	return nil
}

// Keys return all config data
func Keys() []string { return dc.Keys() }

// Keys get all config data
func (c *Config) Keys() []string {
	c.lock.RLock()
	defer c.lock.RUnlock()

	keys := make([]string, 0, len(c.data))
	for key := range c.data {
		keys = append(keys, key)
	}
	return keys
}

// Get config value by key string, support get sub-value by key path(eg. 'map.key'),
func Get(key string, findByPath ...bool) any { return dc.Get(key, findByPath...) }

// Get config value by key, findByPath default is true.
func (c *Config) Get(key string, findByPath ...bool) any {
	val, _ := c.GetValue(key, findByPath...)
	return val
}

// GetValue get value by given key string. findByPath default is true.
func GetValue(key string, findByPath ...bool) (any, bool) {
	return dc.GetValue(key, findByPath...)
}

// GetValue get value by given key string. findByPath default is true.
//
// Return:
//   - ok is true, find value from config
//   - ok is false, not found or error
func (c *Config) GetValue(key string, findByPath ...bool) (value any, ok bool) {
	sep := c.opts.Delimiter
	if key = formatKey(key, string(sep)); key == "" {
		c.addError(ErrKeyIsEmpty)
		return
	}

	// if not is readonly
	if !c.opts.Readonly {
		c.lock.RLock()
		defer c.lock.RUnlock()
	}

	// is top key
	if value, ok = c.data[key]; ok {
		return
	}

	// disable find by path.
	if len(findByPath) > 0 && !findByPath[0] {
		// c.addError(ErrNotFound)
		return
	}

	// has sub key? eg. "lang.dir"
	if strings.IndexByte(key, sep) == -1 {
		// c.addError(ErrNotFound)
		return
	}

	// find by path, the bracket index syntax (eg. "top.arr[1]") is handled there
	keys := strings.Split(key, string(sep))
	if val, found := maputil.GetByPathKeys(c.data, keys); found {
		return val, true
	}

	c.recordGetPathError(key, keys)
	return nil, false
}

// recordGetPathError reports a path that runs into a value which cannot have sub-keys,
// eg. "name.sub" when "name" is a string. A missing key stays silent, as before.
func (c *Config) recordGetPathError(key string, keys []string) {
	for i := len(keys) - 1; i >= 1; i-- {
		if val, ok := maputil.GetByPathKeys(c.data, keys[:i]); ok {
			if !canHaveSubValue(val) {
				c.addErrorf("cannot get value of the key '%s'", key)
			}
			return
		}
	}
}

// canHaveSubValue reports whether the value can be walked into by a sub key.
func canHaveSubValue(val any) bool {
	if val == nil {
		return false
	}

	switch reflect.ValueOf(val).Kind() {
	case reflect.Map, reflect.Slice, reflect.Array:
		return true
	}
	return false
}

/*************************************************************
 * read config (basic data type)
 *************************************************************/

// String get a string by key
func String(key string, defVal ...string) string { return dc.String(key, defVal...) }

// String get a string by key, if not found return default value
func (c *Config) String(key string, defVal ...string) string {
	value, ok := c.getString(key)

	if !ok && len(defVal) > 0 { // give default value
		value = defVal[0]
	}
	return value
}

// MustString get a string by key, will panic on empty or not exists
func MustString(key string) string { return dc.MustString(key) }

// MustString get a string by key, will panic on empty or not exists
func (c *Config) MustString(key string) string {
	value, ok := c.getString(key)
	if !ok {
		panic("config: string value not found, key: " + key)
	}
	return value
}

func (c *Config) getString(key string) (value string, ok bool) {
	enableCache := c.opts.EnableCache
	var gen uint64

	// find from cache
	if enableCache {
		if value, ok = c.cachedString(key); ok {
			return
		}
		gen = c.cacheGeneration()
	}

	val, ok := c.GetValue(key)
	if !ok {
		return
	}

	switch typVal := val.(type) {
	// from json `int` always is float64
	case string:
		value = typVal
		if c.opts.ParseEnv {
			value = envutil.ParseEnvValue(value)
		}
	default:
		var err error
		value, err = strutil.AnyToString(val, false)
		if err != nil {
			return "", false
		}
	}

	// add cache
	if enableCache {
		c.cacheString(gen, key, value)
	}
	return
}

// Int get an int by key
func Int(key string, defVal ...int) int { return dc.Int(key, defVal...) }

// Int get a int value, if not found return default value
func (c *Config) Int(key string, defVal ...int) (value int) {
	i64, exist := c.tryInt64(key)

	if exist {
		// the int64 value may be out of the int range on 32-bit platforms
		if i64 > math.MaxInt || i64 < math.MinInt {
			c.addErrorf("the int value '%d' of key '%s' is out of range", i64, key)
		} else {
			value = int(i64)
		}
	} else if len(defVal) > 0 {
		value = defVal[0]
	}
	return
}

// Uint get a uint value, if not found return default value
func Uint(key string, defVal ...uint) uint { return dc.Uint(key, defVal...) }

// Uint get a int value, if not found return default value
func (c *Config) Uint(key string, defVal ...uint) (value uint) {
	u64, exist := c.tryUint(key)

	if exist {
		value = uint(u64)
	} else if len(defVal) > 0 {
		value = defVal[0]
	}
	return
}

// Int64 get a int value, if not found return default value
func Int64(key string, defVal ...int64) int64 { return dc.Int64(key, defVal...) }

// Int64 get a int value, if not found return default value
func (c *Config) Int64(key string, defVal ...int64) (value int64) {
	value, exist := c.tryInt64(key)

	if !exist && len(defVal) > 0 {
		value = defVal[0]
	}
	return
}

// try to get an int64 value by given key
func (c *Config) tryInt64(key string) (value int64, ok bool) {
	strVal, ok := c.getString(key)
	if !ok {
		return
	}
	value, err := strconv.ParseInt(strVal, 10, 64)
	if err != nil {
		c.addError(err)
		value = 0 // reset the value on error, don't return the clamped max/min value
	}
	return
}

// try to get an uint value by given key
func (c *Config) tryUint(key string) (value uint64, ok bool) {
	strVal, ok := c.getString(key)
	if !ok {
		return
	}

	// the bitSize is the platform uint width, so a value out of range records an error
	value, err := strconv.ParseUint(strVal, 10, strconv.IntSize)
	if err != nil {
		c.addError(err)
		value = 0 // reset the value on error, don't return the clamped max value
	}
	return
}

// Duration get a time.Duration type value. if not found return default value
func Duration(key string, defVal ...time.Duration) time.Duration { return dc.Duration(key, defVal...) }

// Duration get a time.Duration type value. if not found return default value.
//
// The stored value may be a Go duration string (e.g. "300s", "1h30m", "20m").
// A bare integer is treated as nanoseconds for backwards compatibility.
func (c *Config) Duration(key string, defVal ...time.Duration) time.Duration {
	str, exist := c.getString(key)
	if !exist {
		if len(defVal) > 0 {
			return defVal[0]
		}
		return 0
	}

	// Prefer a Go duration string, e.g. "300s", "1h30m".
	if dur, err := time.ParseDuration(str); err == nil {
		return dur
	}

	// Backwards compatible: a bare integer is nanoseconds.
	if n, err := strconv.ParseInt(str, 10, 64); err == nil {
		return time.Duration(n)
	}

	if len(defVal) > 0 {
		return defVal[0]
	}
	return 0
}

// Float get a float64 value, if not found return default value
func Float(key string, defVal ...float64) float64 { return dc.Float(key, defVal...) }

// Float get a float64 by key
func (c *Config) Float(key string, defVal ...float64) (value float64) {
	str, ok := c.getString(key)
	if !ok {
		if len(defVal) > 0 {
			value = defVal[0]
		}
		return
	}

	value, err := strconv.ParseFloat(str, 64)
	if err != nil {
		c.addError(err)
		value = 0 // reset the value on error, don't return the clamped +Inf/-Inf
	}
	return
}

// Bool get a bool value, if not found return default value
func Bool(key string, defVal ...bool) bool { return dc.Bool(key, defVal...) }

// Bool looks up a value for a key in this section and attempts to parse that value as a boolean,
// along with a boolean result similar to a map lookup.
//
// of following(case insensitive):
//   - true
//   - yes
//   - false
//   - no
//   - 1
//   - 0
//
// The `ok` boolean will be false in the event that the value could not be parsed as a bool
func (c *Config) Bool(key string, defVal ...bool) (value bool) {
	rawVal, ok := c.getString(key)
	if !ok {
		if len(defVal) > 0 {
			return defVal[0]
		}
		return
	}

	lowerCase := strings.ToLower(rawVal)
	switch lowerCase {
	case "", "0", "false", "no":
		value = false
	case "1", "true", "yes":
		value = true
	default:
		c.addErrorf("the value '%s' cannot be convert to bool", lowerCase)
	}
	return
}

/*************************************************************
 * read config (complex data type)
 *************************************************************/

// Ints get config data as an int slice/array
func Ints(key string) []int { return dc.Ints(key) }

// Ints get config data as an int slice/array
//
// Note: when ParseEnv is enabled, a string elem is parsed by ParseEnvValue before
// the conversion, so a "${PORT}" elem reads the env value.
func (c *Config) Ints(key string) (arr []int) {
	rawVal, ok := c.GetValue(key)
	if !ok {
		return
	}

	switch typeData := rawVal.(type) {
	case []int:
		arr = typeData
	case []string: // from Set
		for _, v := range typeData {
			iv, err := mathutil.ToInt(c.parseEnvValue(v))
			if err != nil {
				c.addError(err)
				arr = arr[0:0] // reset
				return
			}

			arr = append(arr, iv)
		}
	case []any:
		for _, v := range typeData {
			// iv, err := strconv.Atoi(fmt.Sprintf("%v", v))
			iv, err := mathutil.ToInt(c.parseEnvValue(v))
			if err != nil {
				c.addError(err)
				arr = arr[0:0] // reset
				return
			}

			arr = append(arr, iv)
		}
	default:
		c.addErrorf("value cannot be convert to []int, key is '%s'", key)
	}
	return
}

// IntMap get config data as a map[string]int
func IntMap(key string) map[string]int { return dc.IntMap(key) }

// IntMap get config data as a map[string]int
//
// Note: when ParseEnv is enabled, a string value is parsed by ParseEnvValue before
// the conversion, so a "${PORT}" value reads the env value.
func (c *Config) IntMap(key string) (mp map[string]int) {
	rawVal, ok := c.GetValue(key)
	if !ok {
		return
	}

	switch typeData := rawVal.(type) {
	case map[string]int: // from Set
		mp = typeData
	case map[string]string: // from Set
		mp = make(map[string]int, len(typeData))
		for k, v := range typeData {
			iv, err := mathutil.ToInt(c.parseEnvValue(v))
			if err != nil {
				c.addError(err)
				mp = map[string]int{} // reset
				return
			}
			mp[k] = iv
		}
	case map[string]any: // decode from json,toml
		mp = make(map[string]int)
		for k, v := range typeData {
			// iv, err := strconv.Atoi(fmt.Sprintf("%v", v))
			iv, err := mathutil.ToInt(c.parseEnvValue(v))
			if err != nil {
				c.addError(err)
				mp = map[string]int{} // reset
				return
			}
			mp[k] = iv
		}
	case map[any]any: // if decode from yaml
		mp = make(map[string]int)
		for k, v := range typeData {
			// iv, err := strconv.Atoi(fmt.Sprintf( "%v", v))
			iv, err := mathutil.ToInt(c.parseEnvValue(v))
			if err != nil {
				c.addError(err)
				mp = map[string]int{} // reset
				return
			}

			// sk := fmt.Sprintf("%v", k)
			sk, _ := strutil.AnyToString(k, false)
			mp[sk] = iv
		}
	default:
		c.addErrorf("value cannot be convert to map[string]int, key is '%s'", key)
	}
	return
}

// Strings get strings by key
func Strings(key string) []string { return dc.Strings(key) }

// Strings get config data as a string slice/array
//
// Note: when ParseEnv is enabled, a string elem is parsed by ParseEnvValue and a copy
// of the []string value is returned. Otherwise the []string value of the config data
// may be returned directly, do not change the returned slice in that case.
func (c *Config) Strings(key string) (arr []string) {
	enableCache := c.opts.EnableCache
	var gen uint64

	// find from cache
	if enableCache {
		if cached, ok := c.cachedStringArr(key); ok {
			return cached
		}
		gen = c.cacheGeneration()
	}

	rawVal, ok := c.GetValue(key)
	if !ok {
		return
	}

	switch typeData := rawVal.(type) {
	case []string:
		arr = typeData
		if c.opts.ParseEnv {
			arr = slices.Clone(typeData)
			for i, value := range arr {
				arr[i] = envutil.ParseEnvValue(value)
			}
		}
	case []any:
		for _, v := range typeData {
			arr = append(arr, c.parseEnvString(v))
		}
	default:
		c.addErrorf("value cannot be convert to []string, key is '%s'", key)
		return
	}

	// add cache
	if enableCache {
		c.cacheStringArr(gen, key, arr)
	}
	return
}

// parseEnvString convert a slice elem value to a string. a string elem is parsed by
// ParseEnv when it is enabled, the other types use strutil.MustString as before.
// parseEnvValue parses the ${VAR} references of a string value when ParseEnv is
// enabled. Any other value, and any value with ParseEnv off, is returned as-is.
func (c *Config) parseEnvValue(value any) any {
	if str, ok := value.(string); ok && c.opts.ParseEnv {
		return envutil.ParseEnvValue(str)
	}
	return value
}

func (c *Config) parseEnvString(value any) string {
	return strutil.MustString(c.parseEnvValue(value))
}

// StringsBySplit get []string by split a string value.
func StringsBySplit(key, sep string) []string { return dc.StringsBySplit(key, sep) }

// StringsBySplit get []string by split a string value.
func (c *Config) StringsBySplit(key, sep string) (ss []string) {
	if str, ok := c.getString(key); ok {
		ss = strutil.Split(str, sep)
	}
	return
}

// StringMap get config data as a map[string]string
func StringMap(key string) map[string]string { return dc.StringMap(key) }

// StringMap get config data as a map[string]string
//
// Note: when ParseEnv is enabled, a map[string]string value from Set or LoadData is
// parsed by ParseEnvValue and a copy of it is returned. Otherwise the map[string]string
// value of the config data is returned directly, do not change the returned map in that case.
func (c *Config) StringMap(key string) (mp map[string]string) {
	enableCache := c.opts.EnableCache
	var gen uint64

	// find from cache
	if enableCache {
		if cached, ok := c.cachedStringMap(key); ok {
			return cached
		}
		gen = c.cacheGeneration()
	}

	rawVal, ok := c.GetValue(key)
	if !ok {
		return
	}

	switch typeData := rawVal.(type) {
	case map[string]string: // from Set
		mp = typeData
		if c.opts.ParseEnv {
			mp = maps.Clone(typeData)
			for k, v := range mp {
				mp[k] = envutil.ParseEnvValue(v)
			}
		}
	case map[string]any: // decode from json,toml,yaml.v3
		mp = make(map[string]string, len(typeData))

		for k, v := range typeData {
			switch tv := v.(type) {
			case string:
				if c.opts.ParseEnv {
					mp[k] = envutil.ParseEnvValue(tv)
				} else {
					mp[k] = tv
				}
			default:
				mp[k] = strutil.QuietString(v)
			}
		}
	case map[any]any: // decode from yaml v2
		mp = make(map[string]string, len(typeData))

		for k, v := range typeData {
			sk := strutil.QuietString(k)

			switch typVal := v.(type) {
			case string:
				if c.opts.ParseEnv {
					mp[sk] = envutil.ParseEnvValue(typVal)
				} else {
					mp[sk] = typVal
				}
			default:
				mp[sk] = strutil.QuietString(v)
			}
		}
	default:
		c.addErrorf("value cannot be convert to map[string]string, key is %q", key)
		return
	}

	// add cache
	if enableCache {
		c.cacheStringMap(gen, key, mp)
	}
	return
}

// SubDataMap get sub config data as maputil.Map
func SubDataMap(key string) maputil.Map { return dc.SubDataMap(key) }

// SubDataMap get sub config data as maputil.Map
//
// TIP: will not enable parse Env and more
func (c *Config) SubDataMap(key string) maputil.Map {
	if mp, ok := c.GetValue(key); ok {
		if mmp, ok := mp.(map[string]any); ok {
			return mmp
		}
	}

	// keep is not nil
	return maputil.Map{}
}
