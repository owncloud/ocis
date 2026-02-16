package scoring

import (
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/trustelem/zxcvbn/adjacency"
	"github.com/trustelem/zxcvbn/internal/mathutils"
	"github.com/trustelem/zxcvbn/match"
)

const (
	BruteforceCardinality           = 10
	MinGuessesBeforeGrowingSequence = 10000
	MinSubmatchGuessesSingleChar    = 10
	MinSubmatchGuessesMultiChar     = 50
)

func EstimateGuesses(m *match.Match, password string) float64 {
	if m.Guesses > 0 {
		// a match's guess estimate doesn't change. cache it.
		return m.Guesses
	}
	minGuesses := float64(1)
	if len(m.Token) < len(password) {
		if len(m.Token) == 1 {
			minGuesses = MinSubmatchGuessesSingleChar
		} else {
			minGuesses = MinSubmatchGuessesMultiChar
		}
	}
	var guesses float64
	switch m.Pattern {
	case "bruteforce":
		guesses = BruteforceGuesses(m)
	case "dictionary":
		guesses = DictionaryGuesses(m)
	case "spatial":
		guesses = SpatialGuesses(m)
	case "repeat":
		guesses = RepeatGuesses(m)
	case "sequence":
		guesses = SequenceGuesses(m)
	case "regex":
		guesses = RegexGuesses(m)
	case "date":
		guesses = DateGuesses(m)
	default:
		// panic("unknown pattern " + m.Pattern)
	}
	m.Guesses = guesses
	if m.Guesses < minGuesses {
		m.Guesses = minGuesses
	}
	return m.Guesses
}

func BruteforceGuesses(m *match.Match) float64 {
	runeCount := utf8.RuneCountInString(m.Token)
	guesses := math.Pow(BruteforceCardinality, float64(runeCount))
	/* if guesses == Number.POSITIVE_INFINITY {
		guesses = math.MaxInt;
	}*/
	// small detail: make bruteforce matches at minimum one guess bigger than smallest allowed
	// submatch guesses, such that non-bruteforce submatches over the same [i..j] take precedence.
	minGuesses := float64(0)
	if runeCount == 1 {
		minGuesses = MinSubmatchGuessesSingleChar + 1
	} else {
		minGuesses = MinSubmatchGuessesMultiChar + 1
	}

	if guesses < minGuesses {
		return minGuesses
	}
	return guesses
}

func DictionaryGuesses(m *match.Match) float64 {
	m.BaseGuesses = float64(m.Rank)
	m.UppercaseVariations = UppercaseVariations(m.Token)
	m.L33tVariations = L33tVariations(m)
	reversedVariations := 1
	if m.Reversed {
		reversedVariations = 2
	}
	return float64(m.BaseGuesses) * float64(m.UppercaseVariations) * float64(m.L33tVariations) * float64(reversedVariations)
}

var reStartUpper = regexp.MustCompile(`^[A-Z][^A-Z]+$`)
var reEndUpper = regexp.MustCompile(`^[^A-Z]+[A-Z]$`)
var reAllUpper = regexp.MustCompile(`^[^a-z]+$`)
var reAllLower = regexp.MustCompile(`^[^A-Z]+$`)

func UppercaseVariations(w string) float64 {
	if reAllLower.MatchString(w) || strings.ToLower(w) == w {
		return 1
	}
	// a capitalized word is the most common capitalization scheme,
	// so it only doubles the search space (uncapitalized + capitalized).
	// allcaps and end-capitalized are common enough too, underestimate as 2x factor to be safe.
	if reStartUpper.MatchString(w) ||
		reEndUpper.MatchString(w) ||
		reAllUpper.MatchString(w) {
		return 2
	}
	// otherwise calculate the number of ways to capitalize U+L uppercase+lowercase letters
	// with U uppercase letters or less. or, if there's more uppercase than lower (for eg. PASSwORD),
	// the number of ways to lowercase U+L letters with L lowercase letters or less.
	u := 0
	l := 0
	for _, c := range w {
		if c >= 'A' && c <= 'Z' {
			u++
			continue
		}
		if c >= 'a' && c <= 'z' {
			l++
		}
	}
	variations := float64(0)
	for i := 1; i <= u && i <= l; i++ {
		variations += mathutils.NCk(u+l, i)
	}
	return variations
}

func L33tVariations(m *match.Match) float64 {
	if !m.L33t {
		return 1
	}
	// lower-case match.token before calculating: capitalization shouldn't affect l33t calc.
	chrs := strings.ToLower(m.Token)
	variations := float64(1)
	for subbed, unsubbed := range m.Sub {
		s := 0 // num of subbed chars
		u := 0 // num of unsubbed chars
		for _, c := range chrs {
			if string(c) == subbed {
				s++
			}
			if string(c) == unsubbed {
				u++
			}
		}
		if s == 0 || u == 0 {
			// for this sub, password is either fully subbed (444) or fully unsubbed (aaa)
			// treat that as doubling the space (attacker needs to try fully subbed chars in addition to
			// unsubbed.)
			variations *= 2
		} else {
			// this case is similar to capitalization:
			// with aa44a, U = 3, S = 2, attacker needs to try unsubbed + one sub + two subs
			p := mathutils.Min(u, s)
			possibilities := float64(0)
			for i := 1; i <= p; i++ {
				possibilities += mathutils.NCk(u+s, i)
			}
			variations *= possibilities
		}
	}
	return variations
}

func SpatialGuesses(m *match.Match) float64 {
	s := float64(0)
	d := float64(0)
	switch m.Graph {
	case "qwerty", "dvorak":
		s = float64(len(adjacency.Graphs["qwerty"].Graph))
		d = adjacency.Graphs["qwerty"].AverageDegree
	default:
		s = float64(len(adjacency.Graphs["keypad"].Graph))
		d = adjacency.Graphs["keypad"].AverageDegree
	}
	guesses := float64(0)
	runeCount := utf8.RuneCountInString(m.Token)
	l := runeCount
	t := m.Turns
	// estimate the number of possible patterns w/ length L or less with t turns or less.
	for i := 2; i <= l; i++ {
		possibleTurns := mathutils.Min(t, i-1)
		for j := 1; j <= possibleTurns; j++ {
			guesses += mathutils.NCk(i-1, j-1) * s * math.Pow(d, float64(j))
		}
	}
	// add extra guesses for shifted keys. (% instead of 5, A instead of a.)
	// math is similar to extra guesses of l33t substitutions in dictionary matches.
	if m.ShiftedCount > 0 {
		s := m.ShiftedCount
		u := runeCount - m.ShiftedCount // unshifted count
		if s == 0 || u == 0 {
			guesses *= 2
		} else {
			shiftedVariations := float64(0)
			for i := 1; i <= mathutils.Min(s, u); i++ {
				shiftedVariations += mathutils.NCk(s+u, i)
			}
			guesses *= shiftedVariations
		}
	}
	return (guesses)
}

func RepeatGuesses(m *match.Match) float64 {
	return float64(m.BaseGuesses) * float64(m.RepeatCount)
}

func SequenceGuesses(m *match.Match) float64 {
	firstChr := m.Token[0]
	// lower guesses for obvious starting points
	baseGuesses := 0
	switch firstChr {
	case 'a', 'A', 'z', 'Z', '0', '1', '9':
		baseGuesses = 4
	default:
		if firstChr >= '0' && firstChr <= '9' {
			baseGuesses = 10 // digits
		} else {
			// could give a higher base for uppercase,
			// assigning 26 to both upper and lower sequences is more conservative.
			baseGuesses = 26
		}
	}
	if !m.Ascending {
		// need to try a descending sequence in addition to every ascending sequence ->
		// 2x guesses
		baseGuesses *= 2
	}
	return float64(baseGuesses * len(m.Token))
}

func RegexGuesses(m *match.Match) float64 {
	switch m.RegexName {
	case "alpha_lower":
		return math.Pow(26, float64(len(m.Token)))
	case "alpha_upper":
		return math.Pow(26, float64(len(m.Token)))
	case "alpha":
		return math.Pow(52, float64(len(m.Token)))
	case "alphanumeric":
		return math.Pow(62, float64(len(m.Token)))
	case "digits":
		return math.Pow(10, float64(len(m.Token)))
	case "symbols":
		return math.Pow(33, float64(len(m.Token)))
	case "recent_year":
		// conservative estimate of year space: num years from REFERENCE_YEAR.
		// if year is close to REFERENCE_YEAR, estimate a year space of MIN_YEAR_SPACE.
		year, _ := strconv.Atoi(m.Token)
		yearSpace := mathutils.Abs(year - ReferenceYear)
		yearSpace = mathutils.Max(yearSpace, MinYearSpace)
		return float64(yearSpace)
	default:
		return 0
	}
}

const MinYearSpace = 20

var ReferenceYear = time.Now().Year()

func DateGuesses(m *match.Match) float64 {
	// base guesses: (year distance from ReferenceYear) * num_days * num_years
	yearSpace := mathutils.Max(mathutils.Abs(m.Year-ReferenceYear), MinYearSpace)
	guesses := yearSpace * 365
	// add factor of 4 for separator selection (one of ~4 choices)
	if m.Separator != "" {
		guesses *= 4
	}
	return float64(guesses)
}
