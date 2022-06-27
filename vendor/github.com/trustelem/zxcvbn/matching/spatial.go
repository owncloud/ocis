package matching

import (
	"strings"

	"github.com/trustelem/zxcvbn/adjacency"
	"github.com/trustelem/zxcvbn/match"
)

type spatialMatch struct {
	graphs []*adjacency.Graph
}

func (s spatialMatch) Matches(password string) (matches []*match.Match) {
	for _, graph := range s.graphs {
		if graph.Graph != nil {
			matches = append(matches, spatialMatchHelper(password, graph)...)
		}
	}
	match.Sort(matches)
	return matches
}

var shiftedChars = map[string]map[byte]bool{
	"qwerty": stringToSet(`~!@#$%^&*()_+QWERTYUIOP{}|ASDFGHJKL:"ZXCVBNM<>?`),
	"dvorak": stringToSet(`~!@#$%^&*()_+QWERTYUIOP{}|ASDFGHJKL:"ZXCVBNM<>?`),
}

func stringToSet(s string) map[byte]bool {
	set := make(map[byte]bool)
	for i := 0; i < len(s); i++ {
		set[s[i]] = true
	}
	return set
}

func spatialMatchHelper(password string, graph *adjacency.Graph) (matches []*match.Match) {
	shifted := shiftedChars[graph.Name]

	i := 0
	for i < len(password)-1 {
		j := i + 1
		lastDirection := -99
		turns := 0
		shiftedCount := 0
		if shifted[password[i]] {
			shiftedCount = 1
		}

		for {
			prevChar := password[j-1]
			found := false
			foundDirection := -1
			curDirection := -1
			adjacents := graph.Graph[string(prevChar)]
			// Consider growing pattern by one character if j hasn't gone over the edge
			if j < len(password) {
				curChar := password[j]
				for _, adj := range adjacents {
					curDirection++

					if idx := strings.Index(adj, string(curChar)); idx != -1 {
						found = true
						foundDirection = curDirection

						if idx == 1 {
							// index 1 in the adjacency means the key is shifted, 0 means unshifted: A vs a, % vs 5, etc.
							// for example, 'q' is adjacent to the entry '2@'. @ is shifted w/ index 1, 2 is unshifted.
							shiftedCount++
						}

						if lastDirection != foundDirection {
							// adding a turn is correct even in the initial case when last_direction is null:
							// every spatial pattern starts with a turn.
							turns++
							lastDirection = foundDirection
						}
						break
					}
				}
			}

			// if the current pattern continued, extend j and try to grow again
			if found {
				j++
			} else {
				// otherwise push the pattern discovered so far, if any...
				if j-i > 2 {
					// don't consider length 1 or 2 chains.
					matchSpc := &match.Match{
						Pattern:      "spatial",
						I:            i,
						J:            j - 1,
						Token:        password[i:j],
						Graph:        graph.Name,
						Turns:        turns,
						ShiftedCount: shiftedCount,
					}
					matches = append(matches, matchSpc)
				}
				//. . . and then start a new search from the rest of the password
				i = j
				break
			}
		}

	}
	return matches
}
