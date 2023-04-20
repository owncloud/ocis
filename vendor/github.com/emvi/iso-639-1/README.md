# ISO 639-1

[![Go Reference](https://pkg.go.dev/badge/github.com/emvi/iso-639-1?status.svg)](https://pkg.go.dev/github.com/emvi/iso-639-1?status)
[![CircleCI](https://circleci.com/gh/emvi/iso-639-1.svg?style=svg)](https://circleci.com/gh/emvi/iso-639-1)
[![Go Report Card](https://goreportcard.com/badge/github.com/emvi/iso-639-1)](https://goreportcard.com/report/github.com/emvi/iso-639-1)
<a href="https://discord.gg/fAYm4Cz"><img src="https://img.shields.io/discord/739184135649886288?logo=discord" alt="Chat on Discord"></a>

List of all ISO 639-1 language names, native names and two character codes as well as functions for convenient access.
The lists of all names and codes (`Codes`, `Names`, `NativeNames`, `Languages`) are build in the init function for quick read access. 
For full documentation please read the Godocs.

## Installation

```
go get github.com/emvi/iso-639-1
```

## Example

```
fmt.Println(iso6391.Codes)          // print all codes
fmt.Println(iso6391.Names)          // print all names
fmt.Println(iso6391.NativeNames)    // print all native names
fmt.Println(iso6391.Languages)      // print all language objects {Code, Name, NativeName}

fmt.Println(iso6391.FromCode("en"))             // prints {Code: "en", Name: "English", NativeName: "English"}
fmt.Println(iso6391.Name("en"))                 // prints "English"
fmt.Println(iso6391.NativeName("zh"))           // prints "中文"
fmt.Println(iso6391.CodeFromName("English"))    // prints "en"
fmt.Println(iso6391.ValidCode("en"))            // prints true
// ... see Godoc for more functions
```

## Contribute

[See CONTRIBUTING.md](CONTRIBUTING.md)

## License

MIT
