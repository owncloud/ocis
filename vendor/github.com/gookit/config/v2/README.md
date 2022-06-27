# Config

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/gookit/config?style=flat-square)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/d6ac163ee63649ec92c1566e42f09c11)](https://app.codacy.com/app/inhere/config)
[![GoDoc](https://godoc.org/github.com/gookit/config?status.svg)](https://pkg.go.dev/github.com/gookit/config)
[![Build Status](https://travis-ci.org/gookit/config.svg?branch=master)](https://travis-ci.org/gookit/config)
[![Actions Status](https://github.com/gookit/config/workflows/Unit-Tests/badge.svg)](https://github.com/gookit/config/actions)
[![Coverage Status](https://coveralls.io/repos/github/gookit/config/badge.svg?branch=master)](https://coveralls.io/github/gookit/config?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/gookit/config)](https://goreportcard.com/report/github.com/gookit/config)

Golang's application config manage tool library. 

> **[中文说明](README.zh-CN.md)**

## Features

- Support multi format: `JSON`(default), `JSON5`, `INI`, `YAML`, `TOML`, `HCL`, `ENV`, `Flags`
  - `JSON` content support comments. will auto clear comments
  - Other drivers are used on demand, not used will not be loaded into the application.
    - Possibility to add custom driver for your specific format
- Support multi-file and multi-data loading
- Support loading configuration from os ENV
- Support for loading configuration data from remote URLs
- Support for setting configuration data from command line arguments(`flags`)
- Support listen and fire events on config data changed. 
  - allow events: `set.value`, `set.data`, `load.data`, `clean.data`
- Support data overlay and merge, automatically load by key when loading multiple copies of data
- Support for binding all or part of the configuration data to the structure
- Support get sub value by path, like `map.key` `arr.2`
- Support parse ENV name and allow with default value. like `envKey: ${SHELL|/bin/bash}` -> `envKey: /bin/zsh`
- Generic api `Get` `Int` `Uint` `Int64` `Float` `String` `Bool` `Ints` `IntMap` `Strings` `StringMap` ...
- Complete unit test(code coverage > 95%)

> Provide a sub-package `dotenv` that supports importing data from files (eg `.env`) to ENV

## Only use INI

> If you just want to use INI for simple config management, recommended use [gookit/ini](https://github.com/gookit/ini)

## GoDoc

- [godoc for github](https://pkg.go.dev/github.com/gookit/config)

## Usage

**Install**:

```bash
go get github.com/gookit/config/v2
```

Here using the yaml format as an example(`testdata/yml_other.yml`):

```yaml
name: app2
debug: false
baseKey: value2
shell: ${SHELL}
envKey1: ${NotExist|defValue}

map1:
    key: val2
    key2: val20

arr1:
    - val1
    - val21
```

### Load data

> examples code please see [_examples/yaml.go](_examples/yaml.go):

```go
package main

import (
    "github.com/gookit/config/v2"
    "github.com/gookit/config/v2/yaml"
)

// go run ./examples/yaml.go
func main() {
	config.WithOptions(config.ParseEnv)

	// add driver for support yaml content
	config.AddDriver(yaml.Driver)

	err := config.LoadFiles("testdata/yml_base.yml")
	if err != nil {
		panic(err)
	}

	// load more files
	err = config.LoadFiles("testdata/yml_other.yml")
	// can also load multi at once
	// err := config.LoadFiles("testdata/yml_base.yml", "testdata/yml_other.yml")
	if err != nil {
		panic(err)
	}

	// fmt.Printf("config data: \n %#v\n", config.Data())
}
```

## Bind Structure

> Note: The default binding mapping tag of a structure is `mapstructure`, which can be changed by setting the decoder's option `options.DecoderConfig.TagName`

```go
user := struct {
    Age  int  `mapstructure:"age"`
    Key  string `mapstructure:"key"`
    UserName  string `mapstructure:"user_name"`
    Tags []int  `mapstructure:"tags"`
}{}

err = config.BindStruct("user", &user)

fmt.Println(user.UserName) // inhere
```

**Change struct tag name**

```go
config.WithOptions(func(opt *Options) {
    options.DecoderConfig.TagName = "config"
})
```

Can use empty string for bind all config data to a struct:

```go
config.BindStruct("", &myConf)
```

> `config.MapOnExists` like `BindStruct`，but map binding only if key exists

### Direct read data

- Get integer

```go
age := config.Int("age")
fmt.Print(age) // 100
```

- Get bool

```go
val := config.Bool("debug")
fmt.Print(val) // true
```

- Get string

```go
name := config.String("name")
fmt.Print(name) // inhere
```

- Get strings(slice)

```go
arr1 := config.Strings("arr1")
fmt.Printf("%#v", arr1) // []string{"val1", "val21"}
```

- Get string map

```go
val := config.StringMap("map1")
fmt.Printf("%#v",val) // map[string]string{"key":"val2", "key2":"val20"}
```

- Value contains ENV var

```go
value := config.String("shell")
fmt.Print(value) // "/bin/zsh"
```

- Get value by key path

```go
// from array
value := config.String("arr1.0")
fmt.Print(value) // "val1"

// from map
value := config.String("map1.key")
fmt.Print(value) // "val2"
```

- Setting new value

```go
// set value
config.Set("name", "new name")
name = config.String("name")
fmt.Print(name) // "new name"
```

## Load from flags

> Support simple flags parameter parsing, loading

```go
// flags like: --name inhere --env dev --age 99 --debug

// load flag info
keys := []string{"name", "env", "age:int" "debug:bool"}
err := config.LoadFlags(keys)

// read
config.String("name") // "inhere"
config.String("env") // "dev"
config.Int("age") // 99
config.Bool("debug") // true
```

## Load from ENV

```go
// os env: APP_NAME=config APP_DEBUG=true
// load ENV info
config.LoadOSEnv([]string{"APP_NAME", "APP_DEBUG"}, true)

// read
config.Bool("app_debug") // true
config.String("app_name") // "config"
```

## New config instance

You can create custom config instance

```go
// create new instance, will auto register JSON driver
myConf := config.New("my-conf")

// create empty instance
myConf := config.NewEmpty("my-conf")

// create and with some options
myConf := config.NewWithOptions("my-conf", config.ParseEnv, config.ReadOnly)
```

## Listen config change

Now, you can add a hook func for listen config data change. then, you can do something like: write data to file

Add hook func on create config:

```go
	hookFn := func(event string, c *Config) {
		fmt.Println("fire the:", event)
	}

	c := NewWithOptions("test", WithHookFunc(hookFn))
	// for global config
	config.WithOptions(WithHookFunc(hookFn))
```

After that, when calling `LoadXXX, Set, SetData, ClearData` etc methods, it will output:

```text
fire the: load.data
fire the: set.value
fire the: set.data
fire the: clean.data
```

## Dump config data

> Can use `config.DumpTo()` export the configuration data to the specified `writer`, such as: buffer,file

**Dump to JSON file**

```go
buf := new(bytes.Buffer)

_, err := config.DumpTo(buf, config.JSON)
ioutil.WriteFile("my-config.json", buf.Bytes(), 0755)
```

**Dump pretty JSON**

You can set the default var `JSONMarshalIndent` or custom a new JSON driver. 

```go
config.JSONMarshalIndent = "    "
```

**Dump to YAML file**

```go
_, err := config.DumpTo(buf, config.YAML)
ioutil.WriteFile("my-config.yaml", buf.Bytes(), 0755)
```

## Available options

```go
// Options config options
type Options struct {
	// parse env value. like: "${EnvName}" "${EnvName|default}"
	ParseEnv bool
	// config is readonly. default is False
	Readonly bool
	// enable config data cache. default is False
	EnableCache bool
	// parse key, allow find value by key path. default is True eg: 'key.sub' will find `map[key]sub`
	ParseKey bool
	// tag name for binding data to struct
	// Deprecated
	// please set tag name by DecoderConfig
	TagName string
	// the delimiter char for split key path, if `FindByPath=true`. default is '.'
	Delimiter byte
	// default write format
	DumpFormat string
	// default input format
	ReadFormat string
	// DecoderConfig setting for binding data to struct
	DecoderConfig *mapstructure.DecoderConfig
	// HookFunc on data changed.
	HookFunc HookFunc
}
```

## API Methods Refer

### Load Config

- `LoadOSEnv(keys []string)` Load from os ENV
- `LoadData(dataSource ...interface{}) (err error)` Load from struts or maps
- `LoadFlags(keys []string) (err error)` Load from CLI flags
- `LoadExists(sourceFiles ...string) (err error)` 
- `LoadFiles(sourceFiles ...string) (err error)`
- `LoadRemote(format, url string) (err error)`
- `LoadSources(format string, src []byte, more ...[]byte) (err error)`
- `LoadStrings(format string, str string, more ...string) (err error)`
- `LoadFilesByFormat(format string, sourceFiles ...string) (err error)`
- `LoadExistsByFormat(format string, sourceFiles ...string) error`

### Getting Values

- `Bool(key string, defVal ...bool) bool`
- `Int(key string, defVal ...int) int`
- `Uint(key string, defVal ...uint) uint`
- `Int64(key string, defVal ...int64) int64`
- `Ints(key string) (arr []int)`
- `IntMap(key string) (mp map[string]int)`
- `Float(key string, defVal ...float64) float64`
- `String(key string, defVal ...string) string`
- `Strings(key string) (arr []string)`
- `StringMap(key string) (mp map[string]string)`
- `Get(key string, findByPath ...bool) (value interface{})`

**Mapping data to struct:**

- `BindStruct(key string, dst interface{}) error`
- `MapOnExists(key string, dst interface{}) error`

### Setting Values

- `Set(key string, val interface{}, setByPath ...bool) (err error)`

### Useful Methods

- `Getenv(name string, defVal ...string) (val string)`
- `AddDriver(driver Driver)`
- `Data() map[string]interface{}`
- `SetData(data map[string]interface{})` set data to override the Config.Data
- `Exists(key string, findByPath ...bool) bool`
- `DumpTo(out io.Writer, format string) (n int64, err error)`

## Run Tests

```bash
go test -cover
// contains all sub-folder
go test -cover ./...
```

## Projects using config

Check out these projects, which use https://github.com/gookit/config :

- https://github.com/JanDeDobbeleer/oh-my-posh A prompt theme engine for any shell.
- [+ See More](https://pkg.go.dev/github.com/gookit/config?tab=importedby)

## Gookit packages

- [gookit/ini](https://github.com/gookit/ini) Go config management, use INI files
- [gookit/rux](https://github.com/gookit/rux) Simple and fast request router for golang HTTP 
- [gookit/gcli](https://github.com/gookit/gcli) build CLI application, tool library, running CLI commands
- [gookit/event](https://github.com/gookit/event) Lightweight event manager and dispatcher implements by Go
- [gookit/cache](https://github.com/gookit/cache) Generic cache use and cache manager for golang. support File, Memory, Redis, Memcached.
- [gookit/config](https://github.com/gookit/config) Go config management. support JSON, YAML, TOML, INI, HCL, ENV and Flags
- [gookit/color](https://github.com/gookit/color) A command-line color library with true color support, universal API methods and Windows support
- [gookit/filter](https://github.com/gookit/filter) Provide filtering, sanitizing, and conversion of golang data
- [gookit/validate](https://github.com/gookit/validate) Use for data validation and filtering. support Map, Struct, Form data
- [gookit/goutil](https://github.com/gookit/goutil) Some utils for the Go: string, array/slice, map, format, cli, env, filesystem, test and more
- More, please see https://github.com/gookit

## See also

- Ini parse [gookit/ini/parser](https://github.com/gookit/ini/tree/master/parser)
- Json5 parse [json5](https://github.com/yosuke-furukawa/json5)
- Yaml parse [go-yaml](https://github.com/go-yaml/yaml)
- Toml parse [go toml](https://github.com/BurntSushi/toml)
- Data merge [mergo](https://github.com/imdario/mergo)
- Map structure [mapstructure](https://github.com/mitchellh/mapstructure)

## License

**MIT**
