// Copyright 2015 Benjamin BALET. All rights reserved.
// Use of this source code is governed by the BSD license
// license that can be found in the LICENSE file.

// Package stopwords implements the Levenshtein Distance algorithm to evaluate the diference
// between 2 strings
package stopwords

// LevenshteinDistance compute the LevenshteinDistance between 2 strings
// it removes useless spaces and stop words from a byte slice.
// BCP 47 or ISO 639-1 language code (if unknown, we'll apply english filters).
// If cleanHTML is TRUE, remove HTML tags from content and unescape HTML entities.
func LevenshteinDistance(contentA []byte, contentB []byte, langCode string, cleanHTML bool) int {
	stringA := Clean(contentA, langCode, cleanHTML)
	stringB := Clean(contentB, langCode, cleanHTML)
	distance := levenshteinAlgo(&stringA, &stringB)
	return distance
}

// levenshteinAlgo compute the LevenshteinDistance between 2 strings
func levenshteinAlgo(a, b *[]byte) int {
	la := len(*a)
	lb := len(*b)
	d := make([]int, la+1)
	var lastdiag, olddiag, temp int

	for i := 1; i <= la; i++ {
		d[i] = i
	}
	for i := 1; i <= lb; i++ {
		d[0] = i
		lastdiag = i - 1
		for j := 1; j <= la; j++ {
			olddiag = d[j]
			min := d[j] + 1
			if (d[j-1] + 1) < min {
				min = d[j-1] + 1
			}
			if (*a)[j-1] == (*b)[i-1] {
				temp = 0
			} else {
				temp = 1
			}
			if (lastdiag + temp) < min {
				min = lastdiag + temp
			}
			d[j] = min
			lastdiag = olddiag
		}
	}
	return d[la]
}
