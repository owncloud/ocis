
[![GoDoc](https://godoc.org/github.com/trustelem/zxcvbn?status.svg)](https://godoc.org/github.com/trustelem/zxcvbn)
[![Build
Status](https://travis-ci.org/trustelem/zxcvbn.svg?branch=master)](https://travis-ci.org/trustelem/zxcvbn)
[![Coverage Status](https://coveralls.io/repos/github/trustelem/zxcvbn/badge.svg?branch=master)](https://coveralls.io/github/trustelem/zxcvbn?branch=master)

This is a go port of [zxcvbn](https://github.com/dropbox/zxcvbn), a password strength estimator inspired by password crackers. Through pattern matching and conservative estimation, it recognizes and weighs 30k common passwords, common names and surnames according to US census data, popular English words from Wikipedia and US television and movies, and other common patterns like dates, repeats (aaa), sequences (abcd), keyboard patterns (qwertyuiop), and l33t speak.

This port aims to be fully compatible (i.e. give the same results for a given password using the same set of dictionnaries) with the upstream coffeescript libray from Dropbox: all unit tests from the upstream library have been ported (and even more tests have been added) to ensure that this holds.

------------------------------------------------------------------------

Current status:
- this library should be 100% compatible (score, sequence and number of guesses) with [release 4.4.2](https://github.com/dropbox/zxcvbn/releases/tag/v4.4.2) of the coffeescript library.
- feedback messages are missing