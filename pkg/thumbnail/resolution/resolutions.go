package resolution

import (
	"fmt"
	"math"
	"sort"
)

// New creates an instance of Resolutions from resolution strings.
func New(rStrs []string) (Resolutions, error) {
	var rs Resolutions
	for _, rStr := range rStrs {
		r, err := Parse(rStr)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize resolutions: %s", err.Error())
		}
		rs = append(rs, r)
	}
	sort.Slice(rs, func(i, j int) bool {
		left := rs[i]
		right := rs[j]

		leftSize := left.Width * left.Height
		rightSize := right.Width * right.Height

		return leftSize < rightSize
	})

	return rs, nil
}

// Resolutions represents the available thumbnail resolutions.
type Resolutions []Resolution

// ClosestMatch returns the resolution which is closest to the provided resolution.
// If there is no exact match the resolution will be the next higher one.
// If the given resolution is bigger than all available resolutions the biggest available one is used.
func (r Resolutions) ClosestMatch(width, height int) Resolution {
	if len(r) == 0 {
		return Resolution{Width: width, Height: height}
	}

	isLandscape := width > height
	givenLen := int(math.Max(float64(width), float64(height)))

	// Initialize with the first resolution
	var match Resolution
	minDiff := math.MaxInt32

	for _, current := range r {
		len := dimensionLength(current, isLandscape)
		diff := givenLen - len
		if diff > 0 {
			continue
		}
		absDiff := int(math.Abs(float64(diff)))
		if absDiff < minDiff {
			minDiff = absDiff
			match = current
		}
	}

	if match == (Resolution{}) {
		match = r[len(r)-1]
	}
	return match
}

func dimensionLength(r Resolution, landscape bool) int {
	if landscape {
		return r.Width
	}
	return r.Height
}
