package matching

import (
	"github.com/dlclark/regexp2"
	"github.com/trustelem/zxcvbn/match"
	"github.com/trustelem/zxcvbn/scoring"
)

type repeatMatch struct{}

var greedy = regexp2.MustCompile(`(.+)\1+`, 0)
var lazy = regexp2.MustCompile(`(.+?)\1+`, 0)
var lazyAnchored = regexp2.MustCompile(`^(.+?)\1+$`, 0)

func runeToStringIndex(index int, password string) int {
	runes := 0
	for i := range password {
		if runes == index {
			return i
		}
		runes++
	}
	//shouldn't really get here
	return len(password)
}

func (repeatMatch) Matches(password string) []*match.Match {
	var matches []*match.Match

	lastIndex := 0
	for lastIndex < len(password) {
		greedyMatch, err := greedy.FindStringMatchStartingAt(password, lastIndex)
		if err != nil || greedyMatch == nil {
			break
		}
		lazyMatch, _ := lazy.FindStringMatchStartingAt(password, lastIndex)

		var rmatch *regexp2.Match
		var baseToken string
		if greedyMatch.Captures[0].Length > lazyMatch.Captures[0].Length {
			// greedy beats lazy for 'aabaab'
			//   greedy: [aabaab, aab]
			//   lazy:   [aa,     a]
			rmatch = greedyMatch
			// greedy's repeated string might itself be repeated, eg.
			// aabaab in aabaabaabaab.
			// run an anchored lazy match on greedy's repeated string
			// to find the shortest repeated string
			if m, err := lazyAnchored.FindStringMatch(rmatch.Captures[0].String()); err == nil {
				baseToken = m.GroupByNumber(1).String()
			}
		} else {
			// lazy beats greedy for 'aaaaa'
			//   greedy: [aaaa,  aa]
			//   lazy:   [aaaaa, a]
			rmatch = lazyMatch
			baseToken = rmatch.GroupByNumber(1).String()
		}
		// FindStringMatchStartingAt takes an index into the string (basically an offset
		// into a byte array). rmatch indices will be rune offsets and so need to be converted
		// to string offsets
		i := runeToStringIndex(rmatch.Index, password)
		j := runeToStringIndex(rmatch.Index+rmatch.Captures[0].Length-1, password)

		// recursively match and score the base string
		baseAnalysis := scoring.MostGuessableMatchSequence(
			baseToken,
			Omnimatch(baseToken, nil),
			false,
		)
		matches = append(matches, &match.Match{
			Pattern:     "repeat",
			I:           i,
			J:           j,
			Token:       rmatch.Captures[0].String(),
			BaseToken:   baseToken,
			BaseGuesses: baseAnalysis.Guesses,
			BaseMatches: baseAnalysis.Sequence,
			RepeatCount: rmatch.Captures[0].Length / len(baseToken),
		})
		lastIndex = j + 1

	}
	return matches
}
