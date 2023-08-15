package grammar

import (
	"fmt"
)

type TokenName string

const (
	TagToken      TokenName = "tag"
	NameToken               = "name"
	ContentToken            = "content"
	FallbackToken           = "fallback"
)

type Token struct {
	Name TokenName
	Val  string
}

func NewToken(name TokenName, tail interface{}) (*Token, error) {
	token := &Token{
		Name: name,
	}

	switch name {
	case TagToken, NameToken, ContentToken, FallbackToken:
		token.Val = stringFromChars(tail)
	default:
		return nil, fmt.Errorf("unknown tokentype: '%s'", name)
	}

	return token, nil
}

func stringFromChars(chars interface{}) string {
	str := ""
	r := chars.([]interface{})
	for _, i := range r {
		j := i.([]uint8)
		str += string(j[0])
	}
	return str
}
