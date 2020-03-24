package resolution

import (
	"testing"
)

func TestInitWithEmptyArray(t *testing.T) {
	rs, err := Init([]string{})
	if err != nil {
		t.Errorf("Init with an empty array should not fail. Error: %s.\n", err.Error())
	}
	if len(rs) != 0 {
		t.Error("Init with an empty array should return an empty Resolutions instance.\n")
	}
}

func TestInitWithNil(t *testing.T) {
	rs, err := Init(nil)
	if err != nil {
		t.Errorf("Init with nil parameter should not fail. Error: %s.\n", err.Error())
	}
	if len(rs) != 0 {
		t.Error("Init with nil parameter should return an empty Resolutions instance.\n")
	}
}

func TestInitWithInvalidValuesInArray(t *testing.T) {
	_, err := Init([]string{"invalid"})
	if err == nil {
		t.Error("Init with invalid parameter should fail.\n")
	}
}

func TestInit(t *testing.T) {
	rs, err := Init([]string{"16x16"})
	if err != nil {
		t.Errorf("Init with valid parameter should not fail. Error: %s.\n", err.Error())
	}
	if len(rs) != 1 {
		t.Errorf("resolutions has size %d, expected size %d.\n", len(rs), 1)
	}
}

func TestInitWithMultipleResolutions(t *testing.T) {
	rStrs := []string{"16x16", "32x32", "64x64", "128x128"}
	rs, err := Init(rStrs)
	if err != nil {
		t.Errorf("Init with valid parameter should not fail. Error: %s.\n", err.Error())
	}
	if len(rs) != len(rStrs) {
		t.Errorf("resolutions has size %d, expected size %d.\n", len(rs), len(rStrs))
	}
}

func TestInitWithMultipleResolutionsShouldBeSorted(t *testing.T) {
	rStrs := []string{"32x32", "64x64", "16x16", "128x128"}
	rs, err := Init(rStrs)
	if err != nil {
		t.Errorf("Init with valid parameter should not fail. Error: %s.\n", err.Error())
	}

	for i := 0; i < len(rs)-1; i++ {
		current := rs[i]
		currentSize := current.Width * current.Height
		next := rs[i]
		nextSize := next.Width * next.Height

		if currentSize > nextSize {
			t.Error("Resolutions are not sorted.")
		}

	}
}
func TestClosestMatchWithEmptyResolutions(t *testing.T) {
	rs, _ := Init(nil)
	width := 24
	height := 24

	r := rs.ClosestMatch(width, height)
	if r.Width != width || r.Height != height {
		t.Errorf("ClosestMatch from empty resolutions should return the given resolution")
	}
}

func TestClosestMatch(t *testing.T) {
	rs, _ := Init([]string{"16x16", "24x24", "32x32", "64x64", "128x128"})
	table := [][]int{
		// width, height, expectedWidth, expectedHeight
		[]int{17, 17, 24, 24},
		[]int{12, 17, 24, 24},
		[]int{24, 24, 24, 24},
		[]int{20, 20, 24, 24},
		[]int{20, 80, 128, 128},
		[]int{80, 20, 128, 128},
		[]int{48, 48, 64, 64},
		[]int{1024, 1024, 128, 128},
	}

	for _, row := range table {
		width := row[0]
		height := row[1]
		expectedWidth := row[2]
		expectedHeight := row[3]

		match := rs.ClosestMatch(width, height)

		if match.Width != expectedWidth || match.Height != expectedHeight {
			t.Errorf("Expected resolution %dx%d got %s", expectedWidth, expectedHeight, match.String())
		}
	}
}
