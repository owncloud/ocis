package matching

import (
	"strconv"

	"github.com/dlclark/regexp2"
	"github.com/trustelem/zxcvbn/internal/mathutils"
	"github.com/trustelem/zxcvbn/match"
	"github.com/trustelem/zxcvbn/scoring"
)

const dateMaxYear = 2050
const dateMinYear = 1000

var dateSplits = map[int][]struct{ k, l int }{
	4: { // for length-4 strings, eg 1191 or 9111, two ways to split:
		{1, 2}, // 1 1 91 (2nd split starts at index 1, 3rd at index 2)
		{2, 3}, // 91 1 1
	},
	5: {
		{1, 3}, // 1 11 91
		{2, 3}, // 11 1 91
	},
	6: {
		{1, 2}, // 1 1 1991
		{2, 4}, // 11 11 91
		{4, 5}, // 1991 1 1
	},
	7: {
		{1, 3}, // 1 11 1991
		{2, 3}, // 11 1 1991
		{4, 5}, // 1991 1 11
		{4, 6}, // 1991 11 1
	},
	8: {
		{2, 4}, // 11 11 1991
		{4, 6}, // 1991 11 11
	},
}

var maybeDateNoSeparator = regexp2.MustCompile(
	`^\d{4,8}$`, 0)
var maybeDateWithSeparator = regexp2.MustCompile(
	`^(\d{1,4})([\s/\\_.-])(\d{1,2})\2(\d{1,4})$`, 0)

// a "date" is recognized as:
//   any 3-tuple that starts or ends with a 2- or 4-digit year,
//   with 2 or 0 separator chars (1.1.91 or 1191),
//   maybe zero-padded (01-01-91 vs 1-1-91),
//   a month between 1 and 12,
//   a day between 1 and 31.
//
// note: this isn't true date parsing in that "feb 31st" is allowed,
// this doesn't check for leap years, etc.
//
// recipe:
// start with regex to find maybe-dates, then attempt to map the integers
// onto month-day-year to filter the maybe-dates into dates.
// finally, remove matches that are substrings of other matches to reduce noise.
//
// note: instead of using a lazy or greedy regex to find many dates over the full string,
// this uses a ^...$ regex against every substring of the password -- less performant but leads
// to every possible date match.

type dateMatchCandidate struct {
	Day   int
	Month int
	Year  int
}

type dateMatch struct{}

func (dm dateMatch) Matches(password string) []*match.Match {
	matches := []*match.Match{}

	// dates without separators are between length 4 '1191' and 8 '11111991'
	for i := 0; i <= len(password)-4; i++ {
		for j := i + 3; j <= i+7; j++ {
			if j >= len(password) {
				break
			}
			token := password[i : j+1]
			if m, err := maybeDateNoSeparator.MatchString(token); !m || err != nil {
				continue
			}
			var candidates []*dateMatchCandidate
			for _, s := range dateSplits[len(token)] {
				s1, s2, s3 := token[0:s.k], token[s.k:s.l], token[s.l:]
				if dmy := mapIntsToDMY(s1, s2, s3); dmy != nil {
					candidates = append(candidates, dmy)
				}
			}
			if len(candidates) == 0 {
				continue
			}
			// at this point: different possible dmy mappings for the same i,j substring.
			// match the candidate date that likely takes the fewest guesses: a year closest to 2000.
			// (scoring.REFERENCE_YEAR).
			//
			// ie, considering '111504', prefer 11-15-04 to 1-1-1504
			// (interpreting '04' as 2004)
			bestCandidate := candidates[0]
			minDistance := dateMatchMetric(candidates[0])
			for _, candidate := range candidates[1:] {
				distance := dateMatchMetric(candidate)
				if distance < minDistance {
					bestCandidate = candidate
					minDistance = distance
				}

			}
			matches = append(matches, &match.Match{
				Pattern:   "date",
				Token:     token,
				I:         i,
				J:         j,
				Separator: "",
				Year:      bestCandidate.Year,
				Month:     bestCandidate.Month,
				Day:       bestCandidate.Day,
			})
		}
	}

	// dates with separators are between length 6 '1/1/91' and 10 '11/11/1991'
	for i := 0; i <= len(password)-6; i++ {
		for j := i + 5; j <= i+9; j++ {
			if j >= len(password) {
				break
			}
			token := password[i : j+1]
			m, err := maybeDateWithSeparator.FindStringMatch(token)
			if m == nil || err != nil {
				continue
			}

			dmy := mapIntsToDMY(
				m.GroupByNumber(1).String(),
				m.GroupByNumber(3).String(),
				m.GroupByNumber(4).String(),
			)
			if dmy != nil {
				matches = append(matches, &match.Match{
					Pattern:   "date",
					Token:     token,
					I:         i,
					J:         j,
					Separator: m.GroupByNumber(2).String(),
					Year:      dmy.Year,
					Month:     dmy.Month,
					Day:       dmy.Day,
				})
			}
		}
	}

	// matches now contains all valid date strings in a way that is tricky to capture
	// with regexes only. while thorough, it will contain some unintuitive noise:
	//
	// '2015_06_04', in addition to matching 2015_06_04, will also contain
	// 5(!) other date matches: 15_06_04, 5_06_04, ..., even 2015 (matched as 5/1/2020)
	//
	// to reduce noise, remove date matches that are strict substrings of others
	var filteredMatches []*match.Match
	for _, m := range matches {
		isSubmatch := false
		for _, o := range matches {
			if m == o {
				continue
			}
			if o.I <= m.I && o.J >= m.J {
				isSubmatch = true
				break
			}
		}
		if !isSubmatch {
			filteredMatches = append(filteredMatches, m)
		}
	}
	match.Sort(filteredMatches)
	return filteredMatches
}

func dateMatchMetric(c *dateMatchCandidate) int {
	return mathutils.Abs(c.Year - scoring.ReferenceYear)
}

func mapIntsToDMY(s1, s2, s3 string) *dateMatchCandidate {
	// given a 3-tuple, discard if:
	//   middle int is over 31 (for all dmy formats, years are never allowed in the middle)
	//   middle int is zero
	//   any int is over the max allowable year
	//   any int is over two digits but under the min allowable year
	//   2 ints are over 31, the max allowable day
	//   2 ints are zero
	//   all ints are over 12, the max allowable month
	i1, _ := strconv.Atoi(s1)
	i2, _ := strconv.Atoi(s2)
	i3, _ := strconv.Atoi(s3)
	if i2 > 31 || i2 <= 0 {
		return nil
	}
	over12 := 0
	over31 := 0
	under1 := 0
	for _, i := range [3]int{i1, i2, i3} {
		if (i > 99 && i < dateMinYear) || i > dateMaxYear {
			return nil
		}
		if i > 31 {
			over31++
		}
		if i > 12 {
			over12++
		}
		if i <= 0 {
			under1++
		}
	}
	if over31 >= 2 || over12 == 3 || under1 >= 2 {
		return nil
	}

	// first look for a four digit year: yyyy + daymonth or daymonth + yyyy
	possibleYearSplits := [][3]int{
		{i3, i1, i2}, // year last
		{i1, i2, i3}, // year first
	}
	for _, split := range possibleYearSplits {
		y := split[0]
		if dateMinYear <= y && y <= dateMaxYear {
			// for a candidate that includes a four-digit year,
			// when the remaining ints don't match to a day and month,
			// it is not a date.
			return mapIntsToDM(split[1], split[2], y)
		}
	}

	// given no four-digit year, two digit years are the most flexible int to match, so
	// try to parse a day-month out of ints[0..1] or ints[1..0]
	for _, split := range possibleYearSplits {
		y := split[0]
		dm := mapIntsToDM(split[1], split[2], y)
		if dm != nil {
			dm.Year = twoToFourDigitYear(dm.Year)
			return dm
		}
	}
	return nil
}

func mapIntsToDM(i1, i2 int, year int) *dateMatchCandidate {
	if i1 <= 31 && i2 <= 12 {
		return &dateMatchCandidate{
			Day:   i1,
			Month: i2,
			Year:  year,
		}
	}
	if i2 <= 31 && i1 <= 12 {
		return &dateMatchCandidate{
			Day:   i2,
			Month: i1,
			Year:  year,
		}
	}
	return nil
}

func twoToFourDigitYear(year int) int {
	if year > 99 {
		return year
	} else if year > 50 {
		// 87 -> 1987
		return year + 1900
	} else {
		// 15 -> 2015
		return year + 2000
	}
}
