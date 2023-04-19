package matching

import (
	"regexp"

	"github.com/trustelem/zxcvbn/match"
)

type sequenceMatch struct{}

const maxDelta = 5

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

var reLower = regexp.MustCompile(`^[a-z]+$`)
var reUpper = regexp.MustCompile(`^[A-Z]+$`)
var reDigits = regexp.MustCompile(`^\d+$`)

func (sequenceMatch) Matches(password string) []*match.Match {
	matches := []*match.Match{}
	if len(password) == 1 {
		return matches
	}

	update := func(i, j, delta int) {
		absDelta := abs(delta)
		if j-i > 1 || absDelta == 1 {
			if absDelta > 0 && absDelta <= maxDelta {
				token := password[i : j+1]
				// conservatively stick with roman alphabet size.
				// (this could be improved)
				seqName := "unicode"
				seqSpace := 26
				if reLower.MatchString(token) {
					seqName = "lower"
				} else if reUpper.MatchString(token) {
					seqName = "upper"
				} else if reDigits.MatchString(token) {
					seqName = "digits"
					seqSpace = 10
				}
				matches = append(matches, &match.Match{
					Pattern:       "sequence",
					I:             i,
					J:             j,
					Token:         password[i : j+1],
					SequenceName:  seqName,
					SequenceSpace: seqSpace,
					Ascending:     delta > 0,
				})
			}
		}
	}

	i := 0
	lastDelta := 0 // null
	for k := 1; k <= len(password)-1; k++ {
		delta := int(password[k]) - int(password[k-1])
		if k == 1 {
			lastDelta = delta
		}
		if delta == lastDelta {
			continue
		}
		j := k - 1
		update(i, j, lastDelta)
		i = j
		lastDelta = delta
	}

	update(i, len(password)-1, lastDelta)
	return matches
}
