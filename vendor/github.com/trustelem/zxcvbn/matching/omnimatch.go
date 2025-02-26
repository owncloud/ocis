package matching

import (
	"github.com/trustelem/zxcvbn/adjacency"
	"github.com/trustelem/zxcvbn/frequency"
	"github.com/trustelem/zxcvbn/match"
	"regexp"
)

func Omnimatch(password string, userInputs []string) (matches []*match.Match) {
	dictMatcher := defaultRankedDictionnaries.withDict("user_inputs", buildRankedDict(userInputs))

	matchers := []match.Matcher{
		dictMatcher,
		reverseDictionnaryMatch{dm: dictMatcher},
		l33tMatch{dm: dictMatcher, table: l33tTable},
		spatialMatch{graphs: defaultGraphs},
		repeatMatch{},
		sequenceMatch{},
		regexpMatch{regexes: defaultRegexpMatch},
		dateMatch{},
	}

	for _, m := range matchers {
		matches = append(matches, m.Matches(password)...)
	}
	match.Sort(matches)
	return matches
}

var (
	defaultRankedDictionnaries = loadDefaultDictionnaries()
	defaultGraphs              = loadDefaultAdjacencyGraphs()
	defaultRegexpMatch         = []struct {
		Name   string
		Regexp *regexp.Regexp
	}{
		{
			Name:   "recent_year",
			Regexp: regexp.MustCompile(`19\d\d|200\d|201\d`),
		},
	}
	l33tTable = map[string][]string{
		"a": {"4", "@"},
		"b": {"8"},
		"c": {"(", "{", "[", "<"},
		"e": {"3"},
		"g": {"6", "9"},
		"i": {"1", "!", "|"},
		"l": {"1", "|", "7"},
		"o": {"0"},
		"s": {"$", "5"},
		"t": {"+", "7"},
		"x": {"%"},
		"z": {"2"},
	}
)

func loadDefaultDictionnaries() dictionaryMatch {
	rd := make(map[string]rankedDictionnary)
	for n, list := range frequency.FrequencyLists {
		rd[n] = buildRankedDict(list)
	}
	return dictionaryMatch{
		rankedDictionaries: rd,
	}
}

func loadDefaultAdjacencyGraphs() []*adjacency.Graph {
	return []*adjacency.Graph{
		adjacency.Graphs["qwerty"],
		adjacency.Graphs["dvorak"],
		adjacency.Graphs["keypad"],
		adjacency.Graphs["mac_keypad"],
	}

}
