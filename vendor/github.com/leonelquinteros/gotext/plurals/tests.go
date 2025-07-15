// Original work Copyright (c) 2016 Jonas Obrist (https://github.com/ojii/gettext.go)
// Modified work Copyright (c) 2018 DeineAgentur UG https://www.deineagentur.com
// Modified work Copyright (c) 2018-present gotext maintainers (https://github.com/leonelquinteros/gotext)
//
// Licensed under the BSD License. See License.txt in the project root for license information.

package plurals

type equal struct {
	value uint32
}

func (e equal) test(n uint32) bool {
	return n == e.value
}

type notequal struct {
	value uint32
}

func (e notequal) test(n uint32) bool {
	return n != e.value
}

type gt struct {
	value   uint32
	flipped bool
}

func (e gt) test(n uint32) bool {
	if e.flipped {
		return e.value > n
	}

	return n > e.value
}

type lt struct {
	value   uint32
	flipped bool
}

func (e lt) test(n uint32) bool {
	if e.flipped {
		return e.value < n
	}
	return n < e.value
}

type gte struct {
	value   uint32
	flipped bool
}

func (e gte) test(n uint32) bool {
	if e.flipped {
		return e.value >= n
	}
	return n >= e.value
}

type lte struct {
	value   uint32
	flipped bool
}

func (e lte) test(n uint32) bool {
	if e.flipped {
		return e.value <= n
	}
	return n <= e.value
}

type and struct {
	left  test
	right test
}

func (e and) test(n uint32) bool {
	if !e.left.test(n) {
		return false
	}
	return e.right.test(n)
}

type or struct {
	left  test
	right test
}

func (e or) test(n uint32) bool {
	if e.left.test(n) {
		return true
	}
	return e.right.test(n)
}

type pipe struct {
	modifier math
	action   test
}

func (e pipe) test(n uint32) bool {
	return e.action.test(e.modifier.calc(n))
}
