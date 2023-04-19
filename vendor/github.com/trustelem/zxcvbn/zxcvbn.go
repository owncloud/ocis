package zxcvbn

import (
	"github.com/trustelem/zxcvbn/match"
	"time"
	"unicode/utf8"

	"github.com/trustelem/zxcvbn/matching"
	"github.com/trustelem/zxcvbn/scoring"
)

type Result struct {
	Guesses  float64
	Sequence []*match.Match
	Score    int
	CalcTime float64
}

func PasswordStrength(password string, userInputs []string) Result {
	start := time.Now()
	var result Result
	if !utf8.ValidString(password) {
		// Do not evaluate passwords containing invalid utf8
		// => those will be reported as weak passwords
		return result
	}
	matches := matching.Omnimatch(password, userInputs)
	seq := scoring.MostGuessableMatchSequence(password, matches, false)
	end := time.Now()
	calcTime := end.Nanosecond() - start.Nanosecond()
	result.CalcTime = round(float64(calcTime)*time.Nanosecond.Seconds(), .5, 3)
	result.Sequence = seq.Sequence
	result.Guesses = seq.Guesses
	result.Score = guessesToScore(seq.Guesses)
	return result
}
