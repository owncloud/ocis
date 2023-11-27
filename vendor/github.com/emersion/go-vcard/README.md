# go-vcard

[![GoDoc](https://godoc.org/github.com/emersion/go-vcard?status.svg)](https://godoc.org/github.com/emersion/go-vcard)
[![Build Status](https://travis-ci.org/emersion/go-vcard.svg?branch=master)](https://travis-ci.org/emersion/go-vcard)
[![codecov](https://codecov.io/gh/emersion/go-vcard/branch/master/graph/badge.svg)](https://codecov.io/gh/emersion/go-vcard)

A Go library to parse and format [vCard](https://tools.ietf.org/html/rfc6350).

## Usage

```go
f, err := os.Open("cards.vcf")
if err != nil {
	log.Fatal(err)
}
defer f.Close()

dec := vcard.NewDecoder(f)
for {
	card, err := dec.Decode()
	if err == io.EOF {
		break
	} else if err != nil {
		log.Fatal(err)
	}

	log.Println(card.PreferredValue(vcard.FieldFormattedName))
}
```

## License

MIT
