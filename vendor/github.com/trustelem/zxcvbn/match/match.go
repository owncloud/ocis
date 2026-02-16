package match

import (
	"encoding/json"
)

type Match struct {
	Pattern string `json:"pattern"`
	I       int    `json:"i"`
	J       int    `json:"j"`
	Token   string `json:"token"`

	// Dictionary
	Reversed            bool              `json:"reversed,omitempty"`
	UppercaseVariations float64           `json:"uppercase_variations,omitempty"`
	L33tVariations      float64           `json:"l33t_variations,omitempty"`
	MatchedWord         string            `json:"matched_word,omitempty"`
	Rank                int               `json:"rank,omitempty"`
	DictionaryName      string            `json:"dictionary_name,omitempty"`
	L33t                bool              `json:"l33t,omitempty"`
	Sub                 map[string]string `json:"sub,omitempty"`

	// Sequence
	Graph         string `json:"graph,omitempty"`
	SequenceName  string `json:"sequence_name,omitempty"`
	SequenceSpace int    `json:"sequence_space,omitempty"`
	Ascending     bool   `json:"ascending,omitempty"`
	Turns         int    `json:"turns,omitempty"`
	ShiftedCount  int    `json:"shifted_count,omitempty"`

	// Repeat
	BaseToken   string   `json:"base_token,omitempty"`
	BaseGuesses float64  `json:"base_guesses,omitempty"`
	BaseMatches []*Match `json:"base_matches,omitempty"`
	RepeatCount int      `json:"repeat_count,omitempty"`

	// Regexp
	RegexName string `json:"regex_name,omitempty"`

	// Date
	Year      int     `json:"year,omitempty"`
	Month     int     `json:"month,omitempty"`
	Day       int     `json:"day,omitempty"`
	Separator string  `json:"separator,omitempty"`
	Entropy   float64 `json:"entropy,omitempty"`
	Guesses   float64 `json:"guesses,omitempty"`
}

type Matcher interface {
	Matches(password string) []*Match
}

// ToString returns a string representation of a sequence of matches
func ToString(matches []*Match) string {
	b, _ := json.Marshal(matches)
	return string(b)
}
