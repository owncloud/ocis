// Original work Copyright (c) 2016 Jonas Obrist (https://github.com/ojii/gettext.go)
// Modified work Copyright (c) 2018 DeineAgentur UG https://www.deineagentur.com
// Modified work Copyright (c) 2018-present gotext maintainers (https://github.com/leonelquinteros/gotext)
//
// Licensed under the 3-Clause BSD License. See LICENSE in the project root for license information.

package plurals

// Expression is a plurals expression. Eval evaluates the expression for
// a given n value. Use plurals.Compile to generate Expression instances.
type Expression interface {
	Eval(n uint32) int
}

type constValue struct {
	value int
}

func (c constValue) Eval(n uint32) int {
	return c.value
}

type variableValue struct{}

func (variableValue) Eval(n uint32) int {
	return int(n)
}

type mathValue struct {
	value math
}

func (m mathValue) Eval(n uint32) int {
	if m.value == nil {
		return -1
	}
	return int(m.value.calc(n))
}

type testValue struct {
	condition test
}

func (t testValue) Eval(n uint32) int {
	if t.condition == nil {
		return -1
	}
	if t.condition.test(n) {
		return 1
	}
	return 0
}

type test interface {
	test(n uint32) bool
}

type ternary struct {
	test      test
	trueExpr  Expression
	falseExpr Expression
}

func (t ternary) Eval(n uint32) int {
	if t.test == nil {
		return -1
	}
	if t.test.test(n) {
		if t.trueExpr == nil {
			return -1
		}
		return t.trueExpr.Eval(n)
	}
	if t.falseExpr == nil {
		return -1
	}
	return t.falseExpr.Eval(n)
}
