package zxcvbn

import (
	"fmt"
	"math"
)

type EstimatedTimes struct {
	CrackTimesSeconds map[string]float64 `json:"crack_times_seconds"`
	CrashTimesDisplay map[string]string  `json:"crack_times_display"`
	Score             int
}

func estimateAttackTimes(guesses float64) (t EstimatedTimes) {
	// crack_times_seconds
	t.CrackTimesSeconds = make(map[string]float64)
	t.CrackTimesSeconds["online_throttling_100_per_hour"] = guesses / (100 / 3600)
	t.CrackTimesSeconds["online_no_throttling_10_per_second"] = guesses / 10
	t.CrackTimesSeconds["offline_slow_hashing_1e4_per_second"] = guesses / 1e4
	t.CrackTimesSeconds["offline_fast_hashing_1e10_per_second"] = guesses / 1e10

	t.CrashTimesDisplay = make(map[string]string)

	for scenario, seconds := range t.CrackTimesSeconds {
		t.CrashTimesDisplay[scenario] = displayTime(seconds)
	}

	t.Score = guessesToScore(guesses)
	return
}

func guessesToScore(guesses float64) int {
	const DELTA = 5
	if guesses < 1e3+DELTA {
		// risky password: "too guessable"
		return 0
	}
	if guesses < 1e6+DELTA {
		// modest protection from throttled online attacks: "very guessable"
		return 1
	}
	if guesses < 1e8+DELTA {
		// modest protection from unthrottled online attacks: "somewhat guessable"
		return 2
	}
	if guesses < 1e10+DELTA {
		// modest protection from offline attacks: "safely unguessable"
		// assuming a salted, slow hash function like bcrypt, scrypt, PBKDF2, argon, etc
		return 3
	}
	// strong protection from offline attacks under same scenario: "very unguessable"
	return 4
}

func displayTime(seconds float64) string {
	minute := float64(60)
	hour := minute * 60
	day := hour * 24
	month := day * 31
	year := month * 12
	century := year * 100

	if seconds < 1 {
		return "less than a second"
	}
	if seconds < minute {
		return strCount(seconds, "second")
	} else if seconds < hour {
		return strCount(seconds/minute, "minute")
	} else if seconds < day {
		return strCount(seconds/hour, "hour")
	} else if seconds < month {
		return strCount(seconds/day, "day")
	} else if seconds < year {
		return strCount(seconds/month, "month")
	} else if seconds < century {
		return strCount(seconds/year, "year")
	} else {
		return "centuries"
	}
}

func strCount(count float64, base string) string {
	c := int(round(count, 0.5, 0))
	str := fmt.Sprintf("%d %s", c, base)
	if c > 1 {
		str += "s"
	}
	return str
}

func round(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}
