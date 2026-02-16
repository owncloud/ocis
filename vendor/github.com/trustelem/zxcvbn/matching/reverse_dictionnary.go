package matching

import (
	"github.com/trustelem/zxcvbn/match"
)

type reverseDictionnaryMatch struct {
	dm dictionaryMatch
}

func (rdm reverseDictionnaryMatch) Matches(password string) []*match.Match {
	reversedPassword := reverse(password)
	matches := rdm.dm.Matches(reversedPassword)
	for _, m := range matches {
		m.Token = reverse(m.Token)
		m.Reversed = true
		m.I, m.J = len(password)-1-m.J, len(password)-1-m.I
	}
	match.Sort(matches)
	return matches
}

func reverse(input string) string {
	// Get Unicode code points.
	n := 0
	rune := make([]rune, len(input))
	for _, r := range input {
		rune[n] = r
		n++
	}
	rune = rune[0:n]
	// Reverse
	for i := 0; i < n/2; i++ {
		rune[i], rune[n-1-i] = rune[n-1-i], rune[i]
	}
	// Convert back to UTF-8.
	return string(rune)
}
