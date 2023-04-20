// Copyright 2015 Benjamin BALET. All rights reserved.
// Use of this source code is governed by the BSD license
// license that can be found in the LICENSE file.

// Package stopwords implements Charikar's simhash algorithm to generate a 64-bit
// fingerprint of a given document.
package stopwords

import (
	"bytes"
	"hash/fnv"
	"html"

	"golang.org/x/text/language"
	"golang.org/x/text/unicode/norm"
)

// Each item of the array is the hash of a word of the content.
type vector [64]int

// Internal struct: 64-bit hash and weight of a word
type feature struct {
	Sum    uint64
	Weight int
}

// Simhash returns a 64-bit simhash representing the content of the string
// removes useless spaces and stop words from a byte slice.
// BCP 47 or ISO 639-1 language code (if unknown, we'll apply english filters).
// If cleanHTML is TRUE, remove HTML tags from content and unescape HTML entities.
func Simhash(content []byte, langCode string, cleanHTML bool) uint64 {
	//Remove HTML tags
	if cleanHTML {
		content = remTags.ReplaceAll(content, []byte(" "))
		content = []byte(html.UnescapeString(string(content)))
	}

	//Parse language
	tag := language.Make(langCode)
	base, _ := tag.Base()
	langCode = base.String()
	var hash uint64

	//Remove stop words by using a list of most frequent words
	switch langCode {
	case "ar":
		hash = removeStopWordsAndHash(content, arabic)
	case "bg":
		hash = removeStopWordsAndHash(content, bulgarian)
	case "cs":
		hash = removeStopWordsAndHash(content, czech)
	case "da":
		hash = removeStopWordsAndHash(content, danish)
	case "de":
		hash = removeStopWordsAndHash(content, german)
	case "el":
		hash = removeStopWordsAndHash(content, greek)
	case "en":
		hash = removeStopWordsAndHash(content, english)
	case "es":
		hash = removeStopWordsAndHash(content, spanish)
	case "fa":
		hash = removeStopWordsAndHash(content, persian)
	case "fr":
		hash = removeStopWordsAndHash(content, french)
	case "fi":
		hash = removeStopWordsAndHash(content, finnish)
	case "hu":
		hash = removeStopWordsAndHash(content, hungarian)
	case "it":
		hash = removeStopWordsAndHash(content, italian)
	case "ja":
		hash = removeStopWordsAndHash(content, japanese)
	case "km":
		hash = removeStopWordsAndHash(content, khmer)
	case "lv":
		hash = removeStopWordsAndHash(content, latvian)
	case "nl":
		hash = removeStopWordsAndHash(content, dutch)
	case "no":
		hash = removeStopWordsAndHash(content, norwegian)
	case "pl":
		hash = removeStopWordsAndHash(content, polish)
	case "pt":
		hash = removeStopWordsAndHash(content, portuguese)
	case "ro":
		hash = removeStopWordsAndHash(content, romanian)
	case "ru":
		hash = removeStopWordsAndHash(content, russian)
	case "sk":
		hash = removeStopWordsAndHash(content, slovak)
	case "sv":
		hash = removeStopWordsAndHash(content, swedish)
	case "th":
		hash = removeStopWordsAndHash(content, thai)
	case "tr":
		hash = removeStopWordsAndHash(content, turkish)
	}

	return hash
}

// removeStopWords iterates through a list of words and removes stop words.
func removeStopWordsAndHash(content []byte, dict map[string]string) uint64 {
	var v vector
	var i int

	content = norm.NFC.Bytes(content)
	content = bytes.ToLower(content)
	words := wordSegmenter.FindAll(content, -1)

	for _, w := range words {
		if _, ok := dict[string(w)]; !ok {
			aFeature := newFeature(w)
			sum := aFeature.Sum
			weight := aFeature.Weight
			for i := uint8(0); i < 64; i++ {
				bit := ((sum >> i) & 1)
				if bit == 1 {
					v[i] += weight
				} else {
					v[i] -= weight
				}
			}
			i++
		}
	}

	// compute and return the fingerprint of the content
	// The fingerprint f of a given 64-dimension vector v is defined as follows:
	//   f[j] = 1 if v[j] >= 0
	//   f[j] = 0 if v[j] < 0
	var fingerprint uint64
	for j := uint8(0); j < 64; j++ {
		if v[j] >= 0 {
			fingerprint |= (1 << j)
		}
	}
	return fingerprint
}

// Returns a new feature representing the given byte slice, using a weight of 1
func newFeature(f []byte) feature {
	h := fnv.New64()
	h.Write(f)
	return feature{h.Sum64(), 1}
}

// CompareSimhash calculates the Hamming distance between two 64-bit integers
// using the Kernighan method.
func CompareSimhash(a uint64, b uint64) uint8 {
	v := a ^ b
	var c uint8
	for c = 0; v != 0; c++ {
		v &= v - 1
	}
	return c
}
