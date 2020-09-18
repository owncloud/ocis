package resolution

import (
	"fmt"
	"strconv"
	"strings"
)

// Parse parses a resolution string in the form <width>x<height> and returns a resolution instance.
func Parse(s string) (Resolution, error) {
	parts := strings.Split(s, "x")
	if len(parts) != 2 {
		return Resolution{}, fmt.Errorf("failed to parse resolution: %s. Expected format <width>x<height>", s)
	}
	width, err := strconv.Atoi(parts[0])
	if err != nil {
		return Resolution{}, fmt.Errorf("width: %s has an invalid value. Expected an integer", parts[0])
	}
	height, err := strconv.Atoi(parts[1])
	if err != nil {
		return Resolution{}, fmt.Errorf("height: %s has an invalid value. Expected an integer", parts[1])
	}
	return Resolution{Width: width, Height: height}, nil
}

// Resolution defines represents the width and height of a thumbnail.
type Resolution struct {
	Width  int
	Height int
}

// String returns the resolution in the format:
//
// <width>x<height>
func (r Resolution) String() string {
	return strconv.Itoa(r.Width) + "x" + strconv.Itoa(r.Height)
}
