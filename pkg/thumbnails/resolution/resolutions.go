package resolution

import (
	"fmt"
	"math"
)

// Init creates an instance of Resolutions from resolution strings.
func Init(rStrs []string) (Resolutions, error) {
	var rs Resolutions
	for _, rStr := range rStrs {
		r, err := Parse(rStr)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize resolutions: %s", err.Error())
		}
		rs = append(rs, r)
	}
	return rs, nil
}

// Resolutions represents the available thumbnail resolutions.
type Resolutions []Resolution

// ClosestMatch returns the resolution which is closest to the provided resolution.
func (r Resolutions) ClosestMatch(width, height int) Resolution {
	if len(r) == 0 {
		return Resolution{Width: width, Height: height}
	}

	isLandscape := width > height
	givenLen := math.Max(float64(width), float64(height))

	// Initialize with the first resolution
	match := r[0]
	matchLen := dimensionLength(match, isLandscape)
	minDiff := math.Abs(givenLen - float64(matchLen))

	for i := 1; i < len(r); i++ {
		r := r[i]
		rLen := dimensionLength(r, isLandscape)
		diff := math.Abs(givenLen - float64(rLen))

		if diff <= minDiff {
			minDiff = diff
			match = r
			continue
		}
	}

	return match
}

func dimensionLength(r Resolution, landscape bool) int {
	if landscape {
		return r.Width
	}
	return r.Height
}
