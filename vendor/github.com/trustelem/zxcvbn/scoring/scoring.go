package scoring

import (
	"math"
	"sort"

	"github.com/trustelem/zxcvbn/internal/mathutils"
	"github.com/trustelem/zxcvbn/match"
)

type Result struct {
	Password string
	Guesses  float64
	Sequence []*match.Match
}

// ------------------------------------------------------------------------------
// search --- most guessable match sequence -------------------------------------
// ------------------------------------------------------------------------------

// MostGuessableMatchSequence takes a sequence of overlapping matches, returns the non-overlapping sequence with
// minimum guesses. the following is a O(l_max * (n + m)) dynamic programming algorithm
// for a length-n password with m candidate matches. l_max is the maximum optimal
// sequence length spanning each prefix of the password. In practice it rarely exceeds 5 and the
// search terminates rapidly.
//
// the optimal "minimum guesses" sequence is here defined to be the sequence that
// minimizes the following function:
//
//    g = l! * Product(m.guesses for m in sequence) + D^(l - 1)
//
// where l is the length of the sequence.
//
// the factorial term is the number of ways to order l patterns.
//
// the D^(l-1) term is another length penalty, roughly capturing the idea that an
// attacker will try lower-length sequences first before trying length-l sequences.
//
// for example, consider a sequence that is date-repeat-dictionary.
//  - an attacker would need to try other date-repeat-dictionary combinations,
//    hence the product term.
//  - an attacker would need to try repeat-date-dictionary, dictionary-repeat-date,
//    ..., hence the factorial term.
//  - an attacker would also likely try length-1 (dictionary) and length-2 (dictionary-date)
//    sequences before length-3. assuming at minimum D guesses per pattern type,
//    D^(l-1) approximates Sum(D^i for i in [1..l-1]
//
func MostGuessableMatchSequence(password string, matches []*match.Match, excludeAdditive bool) (result Result) {
	n := len(password)
	validIndexes := make([]bool, n)
	for i := range password {
		validIndexes[i] = true
	}

	// partition matches into sublists according to ending index j
	matchesByJ := make([][]*match.Match, n)
	for _, m := range matches {
		if m.J < 0 || m.J > n {
			// panic(fmt.Sprintf("Invalid match %#v", m))
			continue
		}
		matchesByJ[m.J] = append(matchesByJ[m.J], m)
	}
	// small detail: for deterministic output, sort each sublist by i.
	for _, matches := range matchesByJ {
		sort.SliceStable(matches, func(a, b int) bool {
			return matches[a].I < matches[b].I
		})
	}

	var optimal struct {
		// optimal.m[k][l] holds final match in the best length-l match sequence covering the
		// password prefix up to k, inclusive.
		// if there is no length-l sequence that scores better (fewer guesses) than
		// a shorter match sequence spanning the same prefix, optimal.m[k][l] is undefined.
		m []map[int]*match.Match

		// same structure as optimal.m -- holds the product term Prod(m.guesses for m in sequence).
		// optimal.pi allows for fast (non-looping) updates to the minimization function.
		pi []map[int]float64

		// same structure as optimal.m -- holds the overall metric.
		g []map[int]float64
	}

	optimal.m = make([]map[int]*match.Match, n)
	for i := 0; i < n; i++ {
		optimal.m[i] = make(map[int]*match.Match)
	}

	optimal.pi = make([]map[int]float64, n)
	for i := 0; i < n; i++ {
		optimal.pi[i] = make(map[int]float64)
	}

	optimal.g = make([]map[int]float64, n)
	for i := 0; i < n; i++ {
		optimal.g[i] = make(map[int]float64)
	}

	// helper: considers whether a length-l sequence ending at match m is better (fewer guesses)
	// than previously encountered sequences, updating state if so.
	update := func(m *match.Match, l int) {
		k := m.J
		pi := EstimateGuesses(m, password)
		if l > 1 {
			// we're considering a length-l sequence ending with match m:
			// obtain the product term in the minimization function by multiplying m's guesses
			// by the product of the length-(l-1) sequence ending just before m, at m.i - 1.
			pi *= optimal.pi[m.I-1][l-1]
		}
		// calculate the minimization func
		g := mathutils.Factorial(l) * pi
		if !excludeAdditive {
			g += math.Pow(MinGuessesBeforeGrowingSequence, float64(l-1))
		}
		// update state if new best.
		// first see if any competing sequences covering this prefix, with l or fewer matches,
		// fare better than this sequence. if so, skip it and return.
		if k >= 0 && k < len(optimal.g) {
			for competingL, competingG := range optimal.g[k] {
				if competingL <= l && competingG <= g {
					return
				}
			}
		}
		// this sequence might be part of the final optimal sequence.
		optimal.g[k][l] = g
		optimal.m[k][l] = m
		optimal.pi[k][l] = pi
	}

	// helper: evaluate bruteforce matches ending at k.
	bruteforceUpdate := func(k int) {
		// see if a single bruteforce match spanning the k-prefix is optimal.
		m := makeBruteforceMatch(0, k, password)
		update(m, 1)
		previous := 0
		for i := 1; i <= k; i++ {
			if !validIndexes[i] {
				continue
			}
			prev := previous
			previous = i
			// generate k bruteforce matches, spanning from (i=1, j=k) up to (i=k, j=k).
			// see if adding these new matches to any of the sequences in optimal[i-1]
			// leads to new bests.
			m := makeBruteforceMatch(i, k, password)
			for l := 0; l < n; l++ {
				lastM, ok := optimal.m[prev][l]
				if !ok {
					continue
				}
				// corner: an optimal sequence will never have two adjacent bruteforce matches.
				// it is strictly better to have a single bruteforce match spanning the same region:
				// same contribution to the guess product with a lower length.
				// --> safe to skip those cases.
				if lastM.Pattern == "bruteforce" {
					continue
				}
				// try adding m to this length-l sequence.
				update(m, l+1)
			}
		}
	}

	// helper: step backwards through optimal.m starting at the end,
	// constructing the final optimal match sequence.
	unwind := func(n int) []*match.Match {
		var optimalMatchSequence []*match.Match
		k := n - 1
		// find the final best sequence length and score
		l := -1            // = undefined
		var g float64 = -1 // = Infinity
		if k >= 0 && k < len(optimal.g) {
			first := true
			for candidateL, candidateG := range optimal.g[k] {
				if first || candidateG < g || candidateG == g && candidateL > l {
					l = candidateL
					g = candidateG
					first = false
				}
			}

			for k >= 0 {
				m, ok := optimal.m[k][l]
				l--
				if !ok {
					//we're counting down through the potential keys (l). It's possible that
					//the keys are non-contiguous so we need to skip values of l which aren't
					//valid keys
					continue
				}
				optimalMatchSequence = append([]*match.Match{m}, optimalMatchSequence...)
				k = m.I - 1
			}
		}

		return optimalMatchSequence
	}

	for k := 0; k < n; k++ {
		for _, m := range matchesByJ[k] {
			if m.I > 0 {
				for l := 0; l < n; l++ {
					if optimal.m[m.I-1][l] != nil {
						update(m, l+1)
					}
				}
			} else {
				update(m, 1)
			}
		}
		bruteforceUpdate(k)
	}

	optimalMatchSequence := unwind(n)
	optimalL := len(optimalMatchSequence)

	var guesses float64
	// corner: empty password
	if len(password) == 0 {
		guesses = 1
	} else {
		guesses = optimal.g[n-1][optimalL]
	}

	// final result object
	result.Password = password
	result.Guesses = guesses
	result.Sequence = optimalMatchSequence
	return
}

// helper: make bruteforce match objects spanning i to j, inclusive.
func makeBruteforceMatch(i int, j int, password string) *match.Match {
	return &match.Match{
		Pattern: "bruteforce",
		Token:   password[i : j+1],
		I:       i,
		J:       j,
	}
}
