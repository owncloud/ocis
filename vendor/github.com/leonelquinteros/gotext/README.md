[![GitHub release](https://img.shields.io/github/release/leonelquinteros/gotext.svg)](https://github.com/leonelquinteros/gotext)
[![MIT license](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
![Gotext build](https://github.com/leonelquinteros/gotext/workflows/Gotext%20build/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/leonelquinteros/gotext)](https://goreportcard.com/report/github.com/leonelquinteros/gotext)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/leonelquinteros/gotext)](https://pkg.go.dev/github.com/leonelquinteros/gotext)


# Gotext

[GNU gettext utilities](https://www.gnu.org/software/gettext) for Go.

`gotext` is a native Go implementation of the GNU Gettext utilities. It provides a thread-safe, flexible, and powerful way to handle internationalization (i18n) and localization (l10n) in your Go applications.

---

**📖 [Read the Full Documentation](https://leonelquinteros.github.io/gotext/)**

---

### Table of Contents
- [Features](#features)
- [Installation](#installation)
- [Getting Started](#getting-started)
- [CLI Tool (xgotext)](#cli-tool-xgotext)
- [Advanced Usage](#advanced-usage)
  - [Locale Object](#using-locale-object)
  - [Plural Forms](#use-plural-forms-of-translations)
  - [Dynamic Variables](#using-dynamic-variables-on-translations)
- [Directory Structure](#locales-directories-structure)
- [Contributing](#contributing)
- [License](#license)

---

## Features
- Native Go implementation of Gettext (the core library has no external dependencies).
- Full support for **PO and MO files**.
- **Pluralization rules** support via GNU Gettext plural-form expressions.
- **Message context** support (`msgctxt`).
- Thread-safe for concurrent use.
- Works with UTF-8 by default.
- Integrated CLI tool (`xgotext`) for string extraction.
- Serializable objects for caching.
- Plain-string results that plug directly into Go's `text/template` and `html/template` pipelines.

---

## Installation
```bash
go get github.com/leonelquinteros/gotext
```

---

## Getting Started
For a quick start, use the package-level API:

```go
package main

import (
    "fmt"
    "github.com/leonelquinteros/gotext"
)

func main() {
    // Configure package: locales path, language, and domain
    gotext.Configure("/path/to/locales", "en_US", "default")

    // Simple translation
    fmt.Println(gotext.Get("Hello, world!"))

    // Translation with variables
    fmt.Println(gotext.Get("Hello, %s!", "Gopher"))
}
```

For more details, see the [Getting Started Guide](docs/GETTING_STARTED.md).

---

## CLI Tool (xgotext)
`gotext` includes a command-line tool to extract translatable strings from your Go source code.

**Install xgotext:**
```bash
go install github.com/leonelquinteros/gotext/cli/xgotext@latest
```

**Extract strings:**
```bash
xgotext -in . -out locales/templates
```
The command generates `.pot` templates under `locales/templates`; merge those templates into language-specific `.po` files.

See the [xgotext Documentation](docs/xgotext.md) for full usage details.

---

## Advanced Usage

### Using Locale object
For managing multiple languages or domains independently:

```go
l := gotext.NewLocale("/path/to/locales", "es_UY")
l.AddDomain("default")
fmt.Println(l.Get("Translate this"))
```

### Use plural forms of translations
`gotext` handles complex pluralization rules defined in PO headers:

```go
// GetN(singular, plural, quantity, args...)
fmt.Println(gotext.GetN("I have one apple.", "I have %d apples.", 5, 5))
```
See the [Plural Forms Guide](docs/PLURALS.md) for more examples.

### Using dynamic variables on translations
Supports standard `fmt` package syntax:
```go
name := "John"
fmt.Println(gotext.Get("Hi, my name is %s", name))
```

### More advanced scenarios
- **Embedded translation files** (`embed.FS`) and named `%(name)s` formatting: see the [Advanced Usage Guide](docs/ADVANCED.md).
- **Languages fallback**: pass multiple languages to `Configure` as a colon-separated list (`gotext.Configure("/path/to/locales", "es_UY:es:en_US", "default")`).
- **Serialization/caching** of parsed locales: `Locale.MarshalBinary` / `UnmarshalBinary`.

---

## Locales directories structure
The package expects a standard Gettext directory structure:
```
/path/to/locales
  /en_US
    /LC_MESSAGES
      default.po
  /es_ES
    /LC_MESSAGES
      default.po
```
Translation files can live either inside an `LC_MESSAGES` subdirectory or directly under the language code (`/es_ES/default.po`); both layouts are supported.

It supports automatic language simplification (e.g., falling back from `en_UK` to `en`).

---

## Contributing
We welcome contributions of all kinds! 
- [Contributing Guidelines](CONTRIBUTING.md)
- [Code of Conduct](CODE_OF_CONDUCT.md)

---

## License
`gotext` is released under the [MIT License](LICENSE).
