package matching

import (
	"bytes"
	// "github.com/trustelem/zxcvbn/entropy"
	"github.com/trustelem/zxcvbn/match"
	"sort"
	"strings"
)

type l33tMatch struct {
	dm    dictionaryMatch
	table map[string][]string
}

func (lm l33tMatch) Matches(password string) []*match.Match {
	matches := []*match.Match{}

	substitutions := relevantSubtable(password, lm.table)

	for _, sub := range enumerateLeetSubs(substitutions) {
		if len(sub) == 0 {
			break
		}
		subbedPassword := translate(password, sub)
		for _, m := range lm.dm.Matches(subbedPassword) {
			token := password[m.I : m.J+1]
			if len(token) <= 1 {
				// filter single-character l33t matches to reduce noise.
				// otherwise '1' matches 'i', '4' matches 'a', both very common English words
				continue
			}

			if strings.ToLower(token) == m.MatchedWord {
				continue // only return the matches that return an actual substitution
			}
			m.Sub = make(map[string]string)
			for subbed, chr := range sub {
				if strings.Contains(token, subbed) {
					m.Sub[subbed] = chr
				}
			}
			m.L33t = true
			m.Token = token
			matches = append(matches, m)
		}

	}

	match.Sort(matches)
	return matches
}

func translate(password string, sub map[string]string) string {
	var res string
	for _, s := range password {
		if v, ok := sub[string(s)]; ok {
			res = res + v
		} else {
			res = res + string(s)
		}
	}
	return res
}

type kv struct {
	k string
	v string
}

func dedup(subs [][]kv) [][]kv {
	var res [][]kv
	var b bytes.Buffer
	members := make(map[string]bool)
	for _, sub := range subs {
		sort.SliceStable(sub, func(i, j int) bool {
			return sub[i].k < sub[j].k
		})
		b.Reset()
		for _, x := range sub {
			b.WriteString(x.k)
			b.WriteString(",")
			b.WriteString(x.v)
		}
		key := b.String()
		if !members[key] {
			res = append(res, sub)
			members[key] = true
		}
	}
	return res
}

// enumerateLeetSubs returns the list of possible 1337 replacement dictionaries for a given password
func enumerateLeetSubs(table map[string][]string) []map[string]string {
	var keys []string
	for k := range table {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var subs = [][]kv{[]kv{}}

	var helper func(keys []string)
	helper = func(keys []string) {
		if len(keys) == 0 {
			return
		}
		firstKey := keys[0]
		restKeys := keys[1:]
		var nextSubs [][]kv
		for _, l33tChr := range table[firstKey] {
			for _, sub := range subs {
				dupL33tIndex := -1
				for i := 0; i < len(sub); i++ {
					if sub[i].k == l33tChr {
						dupL33tIndex = i
						break
					}
				}
				if dupL33tIndex == -1 {
					subExtension := append(sub, kv{k: l33tChr, v: firstKey})
					nextSubs = append(nextSubs, subExtension)
				} else {
					subAlternative := make([]kv, 0, len(sub))
					subAlternative = append(subAlternative, sub[0:dupL33tIndex]...)
					subAlternative = append(subAlternative, sub[dupL33tIndex+1:]...)
					subAlternative = append(subAlternative, kv{k: l33tChr, v: firstKey})
					// subAlternative := make([]kv, 0, len(sub))
					// subAlternative = append(subAlternative, sub)
					// subAlternative[dupL33tIndex] = {k:l33tChr,v:firstKey}
					nextSubs = append(nextSubs, sub)
					nextSubs = append(nextSubs, subAlternative)
				}
			}
		}
		subs = dedup(nextSubs)
		helper(restKeys)
	}

	helper(keys)
	var subDicts []map[string]string
	for _, sub := range subs {
		subDict := make(map[string]string)

		for _, x := range sub {
			subDict[x.k] = x.v
		}
		subDicts = append(subDicts, subDict)
	}
	return subDicts
}

func relevantSubtable(password string, table map[string][]string) map[string][]string {
	passwordChars := make(map[rune]bool)
	for _, chr := range password {
		passwordChars[chr] = true
	}

	relevantSubs := make(map[string][]string)
	for key, values := range table {
		for _, value := range values {
			if passwordChars[rune(value[0])] {
				relevantSubs[key] = append(relevantSubs[key], value)
			}
		}
	}
	return relevantSubs
}
