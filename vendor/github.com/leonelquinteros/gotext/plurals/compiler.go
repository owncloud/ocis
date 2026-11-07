// Original work Copyright (c) 2016 Jonas Obrist (https://github.com/ojii/gettext.go)
// Modified work Copyright (c) 2018 DeineAgentur UG https://www.deineagentur.com
// Modified work Copyright (c) 2018-present gotext maintainers (https://github.com/leonelquinteros/gotext)
//
// Licensed under the 3-Clause BSD License. See LICENSE in the project root for license information.

/*
Package plurals is the pluralform compiler to get the correct translation id of the plural string
*/
package plurals

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"unicode"
)

type testToken interface {
	compile(tokens []string) (test test, err error)
}

type cmpTestBuilder func(val uint32, flipped bool) test
type logicTestBuild func(left test, right test) test

func compileTernary(tokens []string, depth int) (expr Expression, err error) {
	if depth > maxParseDepth {
		return nil, errors.New("expression nesting is too deep")
	}

	main, err := splitTokens(tokens, "?")
	if err != nil {
		return nil, err
	}
	test, err := compileTestDepth(strings.Join(main.Left, ""), depth+1)
	if err != nil {
		return nil, err
	}
	actions, err := splitTernaryTokens(main.Right)
	if err != nil {
		return nil, err
	}
	trueAction, err := compileExpressionDepth(strings.Join(actions.Left, ""), depth+1)
	if err != nil {
		return nil, err
	}
	falseAction, err := compileExpressionDepth(strings.Join(actions.Right, ""), depth+1)
	if err != nil {
		return nil, err
	}
	return ternary{
		test:      test,
		trueExpr:  trueAction,
		falseExpr: falseAction,
	}, nil
}

const maxParseDepth = 1024

func splitTernaryTokens(tokens []string) (s splitted, err error) {
	depth := 0
	for index, token := range tokens {
		switch token {
		case "?":
			depth++
		case ":":
			if depth == 0 {
				return splitted{
					Left:  tokens[:index],
					Right: tokens[index+1:],
				}, nil
			}
			depth--
		}
	}
	return s, errors.New("':' not found for ternary expression")
}

var constToken constValStruct

type constValStruct struct{}

func (constValStruct) compile(tokens []string) (expr Expression, err error) {
	if len(tokens) == 0 {
		return nil, errors.New("got nothing instead of constant")
	}
	if len(tokens) != 1 {
		return nil, fmt.Errorf("invalid constant: %s", strings.Join(tokens, ""))
	}
	i, err := strconv.Atoi(tokens[0])
	if err != nil {
		return nil, err
	}
	return constValue{value: i}, nil
}

func compileLogicTestDepth(tokens []string, sep string, builder logicTestBuild, depth int) (test test, err error) {
	if depth > maxParseDepth {
		return nil, errors.New("expression nesting is too deep")
	}
	split, err := splitTokens(tokens, sep)
	if err != nil {
		return nil, err
	}
	left, err := compileTestDepth(strings.Join(split.Left, ""), depth+1)
	if err != nil {
		return nil, err
	}
	right, err := compileTestDepth(strings.Join(split.Right, ""), depth+1)
	if err != nil {
		return nil, err
	}
	return builder(left, right), nil
}

func buildOr(left test, right test) test {
	return or{left: left, right: right}
}

func buildAnd(left test, right test) test {
	return and{left: left, right: right}
}

func unwrapSingleGroup(tokens []string) ([]string, error) {
	for len(tokens) == 1 && strings.HasPrefix(tokens[0], "(") {
		inner, err := tokenize(tokens[0])
		if err != nil {
			return nil, err
		}
		if len(inner) == 1 && inner[0] == tokens[0] {
			break
		}
		tokens = inner
	}
	return tokens, nil
}

func isSimpleN(tokens []string) bool {
	tokens, err := unwrapSingleGroup(tokens)
	if err != nil || len(tokens) != 1 {
		return false
	}
	return tokens[0] == "n"
}

func parseLiteralTokens(tokens []string) (uint32, error) {
	tokens, err := unwrapSingleGroup(tokens)
	if err != nil {
		return 0, err
	}
	if len(tokens) != 1 {
		return 0, errors.New("expected a single integer")
	}
	return parseUint32(tokens[0])
}

func containsSimpleN(tokens []string) bool {
	for _, token := range tokens {
		if isSimpleN([]string{token}) {
			return true
		}
	}
	return false
}

func compileMod(tokens []string) (math math, err error) {
	split, err := splitTokens(tokens, "%")
	if err != nil {
		return math, err
	}
	if !isSimpleN(split.Left) {
		return math, errors.New("modulus operation requires 'n' as left operand")
	}
	i, err := parseLiteralTokens(split.Right)
	if err != nil {
		return math, err
	}
	if i == 0 {
		return math, errors.New("modulus operation requires a non-zero integer as right operand")
	}
	return mod{value: i}, nil
}

func subPipe(modTokens []string, actionTokens []string, builder cmpTestBuilder, flipped bool) (test test, err error) {
	modifier, err := compileMod(modTokens)
	if err != nil {
		return nil, err
	}
	if len(actionTokens) == 0 {
		return nil, errors.New("can only get modulus of integer")
	}
	i, err := parseLiteralTokens(actionTokens)
	if err != nil {
		return nil, err
	}
	action := builder(i, flipped)
	return pipe{
		modifier: modifier,
		action:   action,
	}, nil
}

func compileEquality(tokens []string, sep string, builder cmpTestBuilder) (test test, err error) {
	split, err := splitTokens(tokens, sep)
	if err != nil {
		return nil, err
	}
	left, err := unwrapSingleGroup(split.Left)
	if err != nil {
		return nil, err
	}
	right, err := unwrapSingleGroup(split.Right)
	if err != nil {
		return nil, err
	}
	if isSimpleN(left) {
		i, err := parseLiteralTokens(right)
		if err != nil {
			return nil, err
		}
		return builder(i, false), nil
	}
	if isSimpleN(right) {
		i, err := parseLiteralTokens(left)
		if err != nil {
			return nil, err
		}
		return builder(i, true), nil
	}
	if containsSimpleN(left) && slices.Contains(left, "%") {
		return subPipe(left, right, builder, false)
	}
	if containsSimpleN(right) && slices.Contains(right, "%") {
		return subPipe(right, left, builder, true)
	}
	return nil, errors.New("equality test must have 'n' as one of the two tests")
}

var eqToken eqStruct

type eqStruct struct{}

func (eqStruct) compile(tokens []string) (test test, err error) {
	return compileEquality(tokens, "==", buildEq)
}
func buildEq(val uint32, flipped bool) test {
	return equal{value: val}
}

var neqToken neqStruct

type neqStruct struct{}

func (neqStruct) compile(tokens []string) (test test, err error) {
	return compileEquality(tokens, "!=", buildNeq)
}
func buildNeq(val uint32, flipped bool) test {
	return notequal{value: val}
}

var gtToken gtStruct

type gtStruct struct{}

func (gtStruct) compile(tokens []string) (test test, err error) {
	return compileEquality(tokens, ">", buildGt)
}
func buildGt(val uint32, flipped bool) test {
	return gt{value: val, flipped: flipped}
}

var gteToken gteStruct

type gteStruct struct{}

func (gteStruct) compile(tokens []string) (test test, err error) {
	return compileEquality(tokens, ">=", buildGte)
}
func buildGte(val uint32, flipped bool) test {
	return gte{value: val, flipped: flipped}
}

var ltToken ltStruct

type ltStruct struct{}

func (ltStruct) compile(tokens []string) (test test, err error) {
	return compileEquality(tokens, "<", buildLt)
}
func buildLt(val uint32, flipped bool) test {
	return lt{value: val, flipped: flipped}
}

var lteToken lteStruct

type lteStruct struct{}

func (lteStruct) compile(tokens []string) (test test, err error) {
	return compileEquality(tokens, "<=", buildLte)
}
func buildLte(val uint32, flipped bool) test {
	return lte{value: val, flipped: flipped}
}

type testTokenDef struct {
	op    string
	token testToken
}

var precedence = []testTokenDef{
	{op: "||"},
	{op: "&&"},
	{op: "==", token: eqToken},
	{op: "!=", token: neqToken},
	{op: ">=", token: gteToken},
	{op: ">", token: gtToken},
	{op: "<=", token: lteToken},
	{op: "<", token: ltToken},
}

type splitted struct {
	Left  []string
	Right []string
}

// Split a list of tokens by a token into a splitted struct holding the tokens
// before and after the token to be split by.
func splitTokens(tokens []string, sep string) (s splitted, err error) {
	index := slices.Index(tokens, sep)
	if index == -1 {
		return s, fmt.Errorf("'%s' not found in ['%s']", sep, strings.Join(tokens, "','"))
	}
	return splitted{
		Left:  tokens[:index],
		Right: tokens[index+1:],
	}, nil
}

func trimLexicalWhitespace(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, s)
}

func lexicalToken(s string, index int) (token string, next int, err error) {
	if index >= len(s) {
		return "", index, errors.New("unexpected end of expression")
	}

	switch s[index] {
	case '?', ':', '%', '(', ')':
		return s[index : index+1], index + 1, nil
	case '|':
		if index+1 >= len(s) || s[index+1] != '|' {
			return "", index, errors.New("logical OR must use '||'")
		}
		return "||", index + 2, nil
	case '&':
		if index+1 >= len(s) || s[index+1] != '&' {
			return "", index, errors.New("logical AND must use '&&'")
		}
		return "&&", index + 2, nil
	case '=':
		if index+1 >= len(s) || s[index+1] != '=' {
			return "", index, errors.New("equality must use '=='")
		}
		return "==", index + 2, nil
	case '!':
		if index+1 >= len(s) || s[index+1] != '=' {
			return "", index, errors.New("inequality must use '!='")
		}
		return "!=", index + 2, nil
	case '>':
		if index+1 < len(s) && s[index+1] == '=' {
			return ">=", index + 2, nil
		}
		return ">", index + 1, nil
	case '<':
		if index+1 < len(s) && s[index+1] == '=' {
			return "<=", index + 2, nil
		}
		return "<", index + 1, nil
	case 'n':
		return "n", index + 1, nil
	default:
		if s[index] >= '0' && s[index] <= '9' {
			next = index + 1
			for next < len(s) && s[next] >= '0' && s[next] <= '9' {
				next++
			}
			return s[index:next], next, nil
		}
		return "", index, fmt.Errorf("invalid character %q in expression", s[index])
	}
}

func validateParentheses(s string) error {
	depth := 0
	for index := 0; index < len(s); {
		token, next, err := lexicalToken(s, index)
		if err != nil {
			return err
		}
		switch token {
		case "(":
			depth++
		case ")":
			depth--
			if depth < 0 {
				return errors.New("unmatched closing parenthesis")
			}
		}
		index = next
	}
	if depth != 0 {
		return errors.New("unmatched opening parenthesis")
	}
	return nil
}

func stripRedundantOuterParens(s string) string {
	if len(s) < 2 || s[0] != '(' {
		return s
	}

	leading := 0
	for leading < len(s) && s[leading] == '(' {
		leading++
	}

	openers := make([]int, 0, leading)
	outer := leading
	for index := range s {
		switch s[index] {
		case '(':
			openers = append(openers, index)
		case ')':
			if len(openers) == 0 {
				return s
			}
			opener := openers[len(openers)-1]
			openers = openers[:len(openers)-1]
			if opener < outer && index != len(s)-1-opener {
				outer = opener
			}
		}
	}
	if len(openers) != 0 {
		return s
	}
	return s[outer : len(s)-outer]
}

// Tokenizes a string into a list of strings, tokens grouped by parenthesis are
// not split! If the string starts with ( and ends in ), those are stripped.
func tokenize(s string) ([]string, error) {
	s = trimLexicalWhitespace(s)
	if s == "" {
		return []string{}, nil
	}
	if err := validateParentheses(s); err != nil {
		return nil, err
	}
	s = stripRedundantOuterParens(s)
	if s == "" {
		return []string{}, nil
	}

	ret := []string{}
	for index := 0; index < len(s); {
		if s[index] == '(' {
			start := index
			depth := 0
			groupEnd := false
			for index < len(s) {
				switch s[index] {
				case '(':
					depth++
				case ')':
					depth--
					if depth == 0 {
						index++
						ret = append(ret, s[start:index])
						groupEnd = true
					}
				}
				if groupEnd {
					break
				}
				index++
			}
			if depth != 0 {
				return nil, errors.New("unmatched opening parenthesis")
			}
			continue
		}
		token, next, err := lexicalToken(s, index)
		if err != nil {
			return nil, err
		}
		if token == ")" {
			return nil, errors.New("unmatched closing parenthesis")
		}
		ret = append(ret, token)
		index = next
	}
	return ret, nil
}

// Compile a string containing a plural form expression to a Expression object.
func Compile(s string) (expr Expression, err error) {
	s = trimLexicalWhitespace(s)
	if s == "" {
		return nil, errors.New("empty expression")
	}
	return compileExpressionDepth(s, 0)
}

// Compiles an expression (ternary, numeric leaf, or test).
func compileExpressionDepth(s string, depth int) (expr Expression, err error) {
	if depth > maxParseDepth {
		return nil, errors.New("expression nesting is too deep")
	}
	tokens, err := tokenize(s)
	if err != nil {
		return nil, err
	}
	if slices.Contains(tokens, "?") {
		return compileTernary(tokens, depth)
	}
	for _, tokenDef := range precedence {
		if !slices.Contains(tokens, tokenDef.op) {
			continue
		}
		condition, err := compileTestDepth(s, depth+1)
		if err != nil {
			return nil, err
		}
		return testValue{condition: condition}, nil
	}
	if isSimpleN(tokens) {
		return variableValue{}, nil
	}
	if slices.Contains(tokens, "%") {
		modifier, err := compileMod(tokens)
		if err != nil {
			return nil, err
		}
		return mathValue{value: modifier}, nil
	}
	return constToken.compile(tokens)
}

// Compiles a test (comparison)
func compileTestDepth(s string, depth int) (test test, err error) {
	if depth > maxParseDepth {
		return nil, errors.New("expression nesting is too deep")
	}
	tokens, err := tokenize(s)
	if err != nil {
		return nil, err
	}
	for _, tokenDef := range precedence {
		if !slices.Contains(tokens, tokenDef.op) {
			continue
		}
		switch tokenDef.op {
		case "||":
			return compileLogicTestDepth(tokens, "||", buildOr, depth)
		case "&&":
			return compileLogicTestDepth(tokens, "&&", buildAnd, depth)
		default:
			return tokenDef.token.compile(tokens)
		}
	}
	if slices.Contains(tokens, "%") {
		m, err := compileMod(tokens)
		if err != nil {
			return nil, err
		}
		return pipe{
			modifier: m,
			action:   notequal{value: 0},
		}, nil
	}
	if len(tokens) == 1 && strings.HasPrefix(tokens[0], "(") {
		inner, err := tokenize(tokens[0])
		if err != nil {
			return nil, err
		}
		if len(inner) == 1 && inner[0] == tokens[0] {
			return nil, errors.New("cannot compile parenthesized test")
		}
		return compileTestDepth(tokens[0], depth+1)
	}
	if len(tokens) == 1 && tokens[0] == "n" {
		return notequal{value: 0}, nil
	}
	if value, err := parseLiteralTokens(tokens); err == nil {
		return literalTest{value: value}, nil
	}
	return nil, errors.New("cannot compile")
}

func parseUint32(s string) (ui uint32, err error) {
	s = trimLexicalWhitespace(s)
	i, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return ui, err
	}
	return uint32(i), nil
}
