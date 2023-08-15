package kql

import (
	"fmt"

	"github.com/owncloud/ocis/v2/services/search/pkg/kql/grammar"
)

func Parse(s string) ([]*grammar.Token, error) {
	v, err := grammar.Parse("", []byte(s))
	if err != nil {
		return nil, err
	}

	grammarInterfaceTokens, ok := v.([]interface{})
	if !ok {
		return nil, fmt.Errorf("unable to convert '%v'", grammarInterfaceTokens)
	}

	var grammarTokens []*grammar.Token
	for _, t := range grammarInterfaceTokens {
		grammarToken, ok := t.(*grammar.Token)
		if !ok {
			return nil, fmt.Errorf("unable to convert '%v'", t)
		}

		grammarTokens = append(grammarTokens, grammarToken)
	}

	return grammarTokens, nil
}
