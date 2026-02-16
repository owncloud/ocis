// Original work Copyright (c) 2016 Jonas Obrist (https://github.com/ojii/gettext.go)
// Modified work Copyright (c) 2018 DeineAgentur UG https://www.deineagentur.com
// Modified work Copyright (c) 2018-present gotext maintainers (https://github.com/leonelquinteros/gotext)
//
// Licensed under the 3-Clause BSD License. See LICENSE in the project root for license information.

package plurals

type math interface {
	calc(n uint32) uint32
}

type mod struct {
	value uint32
}

func (m mod) calc(n uint32) uint32 {
	return n % m.value
}
