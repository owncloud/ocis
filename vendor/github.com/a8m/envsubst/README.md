# envsubst
[![GoDoc][godoc-img]][godoc-url]
[![License][license-image]][license-url]
[![Build status][travis-image]][travis-url]
[![Github All Releases][releases-image]][releases]

> Environment variables substitution for Go. see docs [below](#docs)

#### Installation:

##### From binaries
Latest stable `envsubst` [prebuilt binaries for 64-bit Linux, or Mac OS X][releases] are available via Github releases.

###### Linux and MacOS
```console
curl -L https://github.com/a8m/envsubst/releases/download/v1.2.0/envsubst-`uname -s`-`uname -m` -o envsubst
chmod +x envsubst
sudo mv envsubst /usr/local/bin
```

###### Windows
Download the latest prebuilt binary from [releases page][releases], or if you have curl installed:
```console
curl -L https://github.com/a8m/envsubst/releases/download/v1.2.0/envsubst.exe
```

##### With go
You can install via `go get` (provided you have installed go):
```console
go get github.com/a8m/envsubst/cmd/envsubst
```


#### Using via cli
```sh
envsubst < input.tmpl > output.text
echo 'welcome $HOME ${USER:=a8m}' | envsubst
envsubst -help
```

#### Imposing restrictions
There are three command line flags with which you can cause the substitution to stop with an error code, should the restriction associated with the flag not be met. This can be handy if you want to avoid creating e.g. configuration files with unset or empty parameters.
Setting a `-fail-fast` flag in conjunction with either no-unset or no-empty or both will result in a faster feedback loop, this can be especially useful when running through a large file or byte array input, otherwise a list of errors is returned.

The flags and their restrictions are: 

|__Option__     | __Meaning__    | __Type__ | __Default__  |
| ------------| -------------- | ------------ | ------------ |
|`-i`  | input file  | ```string | stdin``` | `stdin`
|`-o`  | output file | ```string | stdout``` |  `stdout` 
|`-no-digit`  | do not replace variables starting with a digit, e.g. $1 and ${1} | `flag` |  `false` 
|`-no-unset`  | fail if a variable is not set | `flag` |  `false` 
|`-no-empty`  | fail if a variable is set but empty | `flag` | `false`
|`-fail-fast`  | fails at first occurence of an error, if `-no-empty` or `-no-unset` flags were **not** specified this is ignored | `flag` | `false`

These flags can be combined to form tighter restrictions. 

#### Using `envsubst` programmatically ?
You can take a look on [`_example/main`](https://github.com/a8m/envsubst/blob/master/_example/main.go) or see the example below.
```go
package main

import (
	"fmt"
	"github.com/a8m/envsubst"
)

func main() {
    input := "welcom $HOME"
    str, err := envsubst.String(input)
    // ...
    buf, err := envsubst.Bytes([]byte(input))
    // ...
    buf, err := envsubst.ReadFile("filename")
}
```
### Docs
> api docs here: [![GoDoc][godoc-img]][godoc-url]

|__Expression__     | __Meaning__    |
| ----------------- | -------------- |
|`${var}`           | Value of var (same as `$var`)
|`${var-$DEFAULT}`  | If var not set, evaluate expression as $DEFAULT
|`${var:-$DEFAULT}` | If var not set or is empty, evaluate expression as $DEFAULT
|`${var=$DEFAULT}`  | If var not set, evaluate expression as $DEFAULT
|`${var:=$DEFAULT}` | If var not set or is empty, evaluate expression as $DEFAULT
|`${var+$OTHER}`    | If var set, evaluate expression as $OTHER, otherwise as empty string
|`${var:+$OTHER}`   | If var set, evaluate expression as $OTHER, otherwise as empty string
|`$$var`            | Escape expressions. Result will be `$var`. 

<sub>Most of the rows in this table were taken from [here](http://www.tldp.org/LDP/abs/html/refcards.html#AEN22728)</sub>

### See also

* `os.ExpandEnv(s string) string` - only supports `$var` and `${var}` notations

#### License
MIT

[releases]: https://github.com/a8m/envsubst/releases
[releases-image]: https://img.shields.io/github/downloads/a8m/envsubst/total.svg?style=for-the-badge
[godoc-url]: https://godoc.org/github.com/a8m/envsubst
[godoc-img]: https://img.shields.io/badge/godoc-reference-blue.svg?style=for-the-badge
[license-image]: https://img.shields.io/badge/license-MIT-blue.svg?style=for-the-badge
[license-url]: LICENSE
[travis-image]: https://img.shields.io/travis/a8m/envsubst.svg?style=for-the-badge
[travis-url]: https://travis-ci.org/a8m/envsubst
