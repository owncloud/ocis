// Copyright 2015 Benjamin BALET. All rights reserved.
// Use of this source code is governed by the BSD license
// license that can be found in the LICENSE file.

// Package stopwords allows you to customize the list of stopwords
package stopwords

import (
	"io/ioutil"
	"strings"

	"golang.org/x/text/language"
)

// LoadStopWordsFromFile loads a list of stop words from a file
// filePath is the full path to the file to be loaded
// langCode is a BCP 47 or ISO 639-1 language code (e.g. "en" for English).
// sep is the string separator (e.g. "\n" for newline)
func LoadStopWordsFromFile(filePath string, langCode string, sep string) {
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
			panic(err)
	}
	LoadStopWordsFromString(string(b), langCode, sep)
}

// LoadStopWordsFromString loads a list of stop words from a string
// filePath is the full path to the file to be loaded
// langCode is a BCP 47 or ISO 639-1 language code (e.g. "en" for English).
// sep is the string separator (e.g. "\n" for newline)
func LoadStopWordsFromString(wordsList string, langCode string, sep string) {

	//Parse language
	tag := language.Make(langCode)
	base, _ := tag.Base()
	langCode = base.String()

	words := strings.Split(strings.ToLower(wordsList), sep)

	switch langCode {
	case "ar":
		arabic = make(map[string]string)
		for _, word := range words {
				arabic[word] = ""
		}
	case "bg":
		bulgarian = make(map[string]string)
		for _, word := range words {
				bulgarian[word] = ""
		}
	case "cs":
		czech = make(map[string]string)
		for _, word := range words {
				czech[word] = ""
		}
	case "da":
		danish = make(map[string]string)
		for _, word := range words {
				danish[word] = ""
		}
	case "de":
		german = make(map[string]string)
		for _, word := range words {
				german[word] = ""
		}
	case "el":
		greek = make(map[string]string)
		for _, word := range words {
				greek[word] = ""
		}
	case "en":
		english = make(map[string]string)
		for _, word := range words {
				english[word] = ""
		}
	case "es":
		spanish = make(map[string]string)
		for _, word := range words {
				spanish[word] = ""
		}
	case "fa":
		persian = make(map[string]string)
		for _, word := range words {
				persian[word] = ""
		}
	case "fr":
		french = make(map[string]string)
		for _, word := range words {
				french[word] = ""
		}
	case "fi":
		finnish = make(map[string]string)
		for _, word := range words {
				finnish[word] = ""
		}
	case "hu":
		hungarian = make(map[string]string)
		for _, word := range words {
				hungarian[word] = ""
		}
	case "id":
		indonesian = make(map[string]string)
		for _, word := range words {
				indonesian[word] = ""
		}
	case "it":
		italian = make(map[string]string)
		for _, word := range words {
				italian[word] = ""
		}
	case "ja":
		japanese = make(map[string]string)
		for _, word := range words {
				japanese[word] = ""
		}
	case "km":
		khmer = make(map[string]string)
		for _, word := range words {
				khmer[word] = ""
		}
	case "lv":
		latvian = make(map[string]string)
		for _, word := range words {
				latvian[word] = ""
		}
	case "nl":
		dutch = make(map[string]string)
		for _, word := range words {
				dutch[word] = ""
		}
	case "no":
		norwegian = make(map[string]string)
		for _, word := range words {
				norwegian[word] = ""
		}
	case "pl":
		polish = make(map[string]string)
		for _, word := range words {
				polish[word] = ""
		}
	case "pt":
		portuguese = make(map[string]string)
		for _, word := range words {
				portuguese[word] = ""
		}
	case "ro":
		romanian = make(map[string]string)
		for _, word := range words {
				romanian[word] = ""
		}
	case "ru":
		russian = make(map[string]string)
		for _, word := range words {
				russian[word] = ""
		}
	case "sk":
		slovak = make(map[string]string)
		for _, word := range words {
				slovak[word] = ""
		}
	case "sv":
		swedish = make(map[string]string)
		for _, word := range words {
				swedish[word] = ""
		}
	case "th":
		thai = make(map[string]string)
		for _, word := range words {
				thai[word] = ""
		}
	case "tr":
		turkish = make(map[string]string)
		for _, word := range words {
				turkish[word] = ""
		}
	}
}
