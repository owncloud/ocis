**stopwords** is a go package that removes stop words from a text content.
If instructed to do so, it will remove HTML tags and parse HTML entities.
The objective is to prepare a text in view to be used by natural processing algos
or text comparison algorithms such as SimHash.

[![GoDoc](https://godoc.org/github.com/bbalet/stopwords?status.svg)](https://godoc.org/github.com/bbalet/stopwords)
[![Build Status](https://api.travis-ci.org/bbalet/stopwords.png)](https://travis-ci.org/bbalet/stopwords)
[![codecov.io](https://codecov.io/github/bbalet/stopwords/coverage.svg?branch=master)](https://codecov.io/github/bbalet/stopwords?branch=master)
[![Go Report Card](https://goreportcard.com/badge/bbalet/stopwords)](https://goreportcard.com/report/bbalet/stopwords)

[![Join the chat at https://gitter.im/bbalet/stopwords](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/bbalet/stopwords?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

It uses a curated list of the most frequent words used in these languages:

 * Arabic
 * Bulgarian
 * Czech
 * Danish
 * English
 * Finnish
 * French
 * German
 * Hungarian
 * Italian
 * Japanese
 * Khmer
 * Latvian
 * Norwegian
 * Persian
 * Polish
 * Portuguese
 * Romanian
 * Russian
 * Slovak
 * Spanish
 * Swedish
 * Thai
 * Turkish

If the function is used with an unsupported language, it doesn't fail, but will apply english filter to the content.

## How to use this package?

You can find an example here https:github.com/bbalet/gorelated where **stopwords**
package is used in conjunction with SimHash algorithm in order to find a list of
related content for a static website generator:

    import (
	      "github.com/bbalet/stopwords"
    )

    //Example with 2 strings containing P html tags
    //"la", "un", etc. are (stop) words without lexical value in French
    string1 := []byte("<p>la fin d'un bel après-midi d'été</p>")
    string2 := []byte("<p>cet été, nous avons eu un bel après-midi</p>")

    //Return a string where HTML tags and French stop words has been removed
    cleanContent := stopwords.CleanString(string1, "fr", true)

    //Get two (Sim) hash representing the content of each string
    hash1 := stopwords.Simhash(string1, "fr", true)
    hash2 := stopwords.Simhash(string2, "fr", true)

  	//Hamming distance between the two strings (diffference between contents)
  	distance := stopwords.CompareSimhash(hash1, hash2)

    //Clean the content of string1 and string2, compute the Levenshtein Distance
    stopwords.LevenshteinDistance(string1, string2, "fr", true)

Where *fr* is the ISO 639-1 code for French (it accepts a BCP 47 tag as well).
https://en.wikipedia.org/wiki/List_of_ISO_639-1_codes

## How to load a custom list of stop words from a file/string?

This package comes with a predefined list of stopwords.
However, two functions allow you to use your own list of words:

    stopwords.LoadStopWordsFromFile(filePath, langCode, separator)
    stopwords.LoadStopWordsFromString(wordsList, langCode, separator)

They will overwrite the predefined words for a given language.
You can find an example with the file `stopwords.txt`

## How to overwrite the word segmenter?

If you don't want to strip the Unicode Characters of the 'Number, Decimal Digit'
Category, call the function `DontStripDigits` before using the package :

    stopwords.DontStripDigits()

If you want to use your own segmenter, you can overwrite the regular expression:

    stopwords.OverwriteWordSegmenter(`[\pL]+`)

## Limitations

Please note that this library doesn't break words. If you want to break words prior using stopwords, you need to use another library that provides a binding to ICU library.

These curated lists contain the most used words in various topics, they were not built with a corpus limited to any given specialized topic.

## Credits

Most of the lists were built by IR Multilingual Resources at UniNE
http://members.unine.ch/jacques.savoy/clef/index.html

## License

**stopwords** is released under the BSD license.
