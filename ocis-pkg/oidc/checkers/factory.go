package checkers

import (
	"strconv"
	"strings"
)

// Factory takes care of creating new claims checkers
type Factory struct {
}

// NewFactory creates a new factory
func NewFactory() *Factory {
	return &Factory{}
}

// GetChecker gets a new claims checker based on the provided name, and
// configured with the provided parameter string.
// If the name is unknown or there are problems with the parameters, a
// NoopChecker will be returned instead.
func (f *Factory) GetChecker(name, paramString string) Checker {
	params := f.parseParamString(paramString)
	switch name {
	case "Bool":
		key, keyok := params["key"]
		value, valueok := params["value"]
		if keyok && valueok {
			if boolValue, err := strconv.ParseBool(value); err == nil {
				return NewBooleanChecker(key, boolValue)
			}
		}
	case "Regexp":
		key, keyok := params["key"]
		pattern, patternok := params["value"]
		if keyok && patternok {
			return NewRegexpChecker(key, pattern)
		}
	case "Acr":
		value, valueok := params["value"]
		if valueok {
			return NewAcrChecker(value)
		}
	}
	return NewNoopChecker()
}

func (f *Factory) parseParamString(paramString string) map[string]string {
	params := make(map[string]string)
	paramList := strings.Split(paramString, ";")
	for _, keyvalue := range paramList {
		p := strings.SplitN(keyvalue, "=", 2)
		if len(p) == 2 {
			params[p[0]] = p[1]
		}
	}
	return params
}
