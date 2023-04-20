package matching

import (
	"github.com/trustelem/zxcvbn/match"
	"regexp"
)

type regexpMatch struct {
	regexes []struct {
		Name   string
		Regexp *regexp.Regexp
	}
}

func (r regexpMatch) Matches(password string) []*match.Match {
	var matches []*match.Match
	for _, rx := range r.regexes {
		for _, indexes := range rx.Regexp.FindAllStringIndex(password, -1) {
			token := password[indexes[0]:indexes[1]]
			matches = append(matches, &match.Match{
				Pattern:   "regex",
				Token:     token,
				I:         indexes[0],
				J:         indexes[1] - 1,
				RegexName: rx.Name,
			})
		}
	}
	match.Sort(matches)
	return matches
}
