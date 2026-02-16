// Copyright 2015 Benjamin BALET. All rights reserved.
// Use of this source code is governed by the BSD license
// license that can be found in the LICENSE file.

// stopwords package removes most frequent words from a text content.
// It can be used to improve the accuracy of SimHash algo for example.
// It uses a list of most frequent words used in various languages :
//
// arabic, bulgarian, czech, danish, english, finnish, french, german,
// hungarian, italian, japanese, latvian, norwegian, persian, polish,
// portuguese, romanian, russian, slovak, spanish, swedish, turkish

// Package stopwords contains various algorithms of text comparison (Simhash, Levenshtein)
package stopwords

import (
	"bytes"
	"html"
	"regexp"

	"golang.org/x/text/language"
	"golang.org/x/text/unicode/norm"
)

var (
	remTags      = regexp.MustCompile(`<[^>]*>`)
	oneSpace     = regexp.MustCompile(`\s{2,}`)
	wordSegmenter = regexp.MustCompile(`[\pL\p{Mc}\p{Mn}-_']+`)
)

// DontStripDigits changes the behaviour of the default word segmenter
// by including 'Number, Decimal Digit' Unicode Category as words
func DontStripDigits() {
	wordSegmenter = regexp.MustCompile(`[\pL\p{Mc}\p{Mn}\p{Nd}-_']+`)
}

// OverwriteWordSegmenter allows you to overwrite the default word segmenter
// with your own regular expression
func OverwriteWordSegmenter(expression string) {
	wordSegmenter = regexp.MustCompile(expression)
}

// CleanString removes useless spaces and stop words from string content.
// BCP 47 or ISO 639-1 language code (if unknown, we'll apply english filters).
// If cleanHTML is TRUE, remove HTML tags from content and unescape HTML entities.
func CleanString(content string, langCode string, cleanHTML bool) string {
	return string(Clean([]byte(content), langCode, cleanHTML))
}

// Clean removes useless spaces and stop words from a byte slice.
// BCP 47 or ISO 639-1 language code (if unknown, we'll apply english filters).
// If cleanHTML is TRUE, remove HTML tags from content and unescape HTML entities.
func Clean(content []byte, langCode string, cleanHTML bool) []byte {
	//Remove HTML tags
	if cleanHTML {
		content = remTags.ReplaceAll(content, []byte(" "))
		content = []byte(html.UnescapeString(string(content)))
	}

	//Parse language
	tag := language.Make(langCode)
	base, _ := tag.Base()
	langCode = base.String()

	//Remove stop words by using a list of most frequent words
	switch langCode {
	case "ar":
		content = removeStopWords(content, arabic)
	case "bg":
		content = removeStopWords(content, bulgarian)
	case "cs":
		content = removeStopWords(content, czech)
	case "da":
		content = removeStopWords(content, danish)
	case "de":
		content = removeStopWords(content, german)
	case "el":
		content = removeStopWords(content, greek)
	case "en":
		content = removeStopWords(content, english)
	case "es":
		content = removeStopWords(content, spanish)
	case "fa":
		content = removeStopWords(content, persian)
	case "fr":
		content = removeStopWords(content, french)
	case "fi":
		content = removeStopWords(content, finnish)
	case "hu":
		content = removeStopWords(content, hungarian)
	case "id":
		content = removeStopWords(content, indonesian)
	case "it":
		content = removeStopWords(content, italian)
	case "ja":
		content = removeStopWords(content, japanese)
	case "km":
		content = removeStopWords(content, khmer)
	case "lv":
		content = removeStopWords(content, latvian)
	case "nl":
		content = removeStopWords(content, dutch)
	case "no":
		content = removeStopWords(content, norwegian)
	case "pl":
		content = removeStopWords(content, polish)
	case "pt":
		content = removeStopWords(content, portuguese)
	case "ro":
		content = removeStopWords(content, romanian)
	case "ru":
		content = removeStopWords(content, russian)
	case "sk":
		content = removeStopWords(content, slovak)
	case "sv":
		content = removeStopWords(content, swedish)
	case "th":
		content = removeStopWords(content, thai)
	case "tr":
		content = removeStopWords(content, turkish)
	}

	//Remove duplicated space characters
	content = oneSpace.ReplaceAll(content, []byte(" "))

	return content
}

// removeStopWords iterates through a list of words and removes stop words.
func removeStopWords(content []byte, dict map[string]string) []byte {
	var result []byte
	content = norm.NFC.Bytes(content)
	content = bytes.ToLower(content)
	words := wordSegmenter.FindAll(content, -1)
	for _, w := range words {
		//log.Println(w)
		if _, ok := dict[string(w)]; ok {
			result = append(result, ' ')
		} else {
			result = append(result, []byte(w)...)
			result = append(result, ' ')
		}
	}
	return result
}
