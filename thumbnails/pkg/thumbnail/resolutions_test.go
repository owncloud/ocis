package thumbnail

import (
	"image"
	"testing"
)

func TestInitWithEmptyArray(t *testing.T) {
	rs, err := ParseResolutions([]string{})
	if err != nil {
		t.Errorf("Init with an empty array should not fail. Error: %s.\n", err.Error())
	}
	if len(rs) != 0 {
		t.Error("Init with an empty array should return an empty Resolutions instance.\n")
	}
}

func TestInitWithNil(t *testing.T) {
	rs, err := ParseResolutions(nil)
	if err != nil {
		t.Errorf("Init with nil parameter should not fail. Error: %s.\n", err.Error())
	}
	if len(rs) != 0 {
		t.Error("Init with nil parameter should return an empty Resolutions instance.\n")
	}
}

func TestInitWithInvalidValuesInArray(t *testing.T) {
	_, err := ParseResolutions([]string{"invalid"})
	if err == nil {
		t.Error("Init with invalid parameter should fail.\n")
	}
}

func TestInit(t *testing.T) {
	rs, err := ParseResolutions([]string{"16x16"})
	if err != nil {
		t.Errorf("Init with valid parameter should not fail. Error: %s.\n", err.Error())
	}
	if len(rs) != 1 {
		t.Errorf("resolutions has size %d, expected size %d.\n", len(rs), 1)
	}
}

func TestInitWithMultipleResolutions(t *testing.T) {
	rStrs := []string{"16x16", "32x32", "64x64", "128x128"}
	rs, err := ParseResolutions(rStrs)
	if err != nil {
		t.Errorf("Init with valid parameter should not fail. Error: %s.\n", err.Error())
	}
	if len(rs) != len(rStrs) {
		t.Errorf("resolutions has size %d, expected size %d.\n", len(rs), len(rStrs))
	}
}

func TestClosestMatchWithEmptyResolutions(t *testing.T) {
	rs, _ := ParseResolutions(nil)
	want := image.Rect(0, 0, 24, 24)
	imgSize := image.Rect(0, 0, 24, 24)

	r := rs.ClosestMatch(want, imgSize)
	if r.Dx() != want.Dx() || r.Dy() != want.Dy() {
		t.Errorf("ClosestMatch from empty resolutions should return the given resolution")
	}
}

func TestClosestMatch(t *testing.T) {
	rs, _ := ParseResolutions([]string{"16x16", "24x24", "32x32", "64x64", "128x128"})

	testData := [][]image.Rectangle{
		{image.Rect(0, 0, 17, 17), image.Rect(0, 0, 1920, 1080), image.Rect(0, 0, 24, 24)},
		{image.Rect(0, 0, 12, 17), image.Rect(0, 0, 1080, 1920), image.Rect(0, 0, 24, 24)},
		{image.Rect(0, 0, 24, 24), image.Rect(0, 0, 1920, 1080), image.Rect(0, 0, 24, 24)},
		{image.Rect(0, 0, 20, 20), image.Rect(0, 0, 1920, 1080), image.Rect(0, 0, 24, 24)},
		{image.Rect(0, 0, 20, 80), image.Rect(0, 0, 1080, 1920), image.Rect(0, 0, 128, 128)},
		{image.Rect(0, 0, 80, 20), image.Rect(0, 0, 1920, 1080), image.Rect(0, 0, 128, 128)},
		{image.Rect(0, 0, 48, 48), image.Rect(0, 0, 1920, 1080), image.Rect(0, 0, 64, 64)},
		{image.Rect(0, 0, 1024, 1024), image.Rect(0, 0, 1920, 1080), image.Rect(0, 0, 128, 128)},
		{image.Rect(0, 0, 1920, 1080), image.Rect(0, 0, 256, 36), image.Rect(0, 0, 256, 36)},
	}

	for _, row := range testData {
		given := row[0]
		imgSize := row[1]
		expected := row[2]

		match := rs.ClosestMatch(given, imgSize)

		if match != expected {
			t.Errorf("Expected resolution %dx%d got %dx%d", expected.Dx(), expected.Dy(), match.Dx(), match.Dy())
		}
	}
}

func TestParseWithEmptyString(t *testing.T) {
	_, err := ParseResolution("")
	if err == nil {
		t.Error("Parse with empty string should return an error.")
	}
}

func TestParseWithInvalidWidth(t *testing.T) {
	_, err := ParseResolution("invalidx42")
	if err == nil {
		t.Error("Parse with invalid width should return an error.")
	}
}

func TestParseWithInvalidHeight(t *testing.T) {
	_, err := ParseResolution("42xinvalid")
	if err == nil {
		t.Error("Parse with invalid height should return an error.")
	}
}

func TestParseResolution(t *testing.T) {
	rStr := "42x23"
	r, _ := ParseResolution(rStr)
	if r.Dx() != 42 || r.Dy() != 23 {
		t.Errorf("Expected resolution %s got %s", rStr, r.String())
	}
}

func TestMapRatio(t *testing.T) {
	testData := [][]image.Rectangle{
		{image.Rect(0, 0, 1920, 1080), image.Rect(0, 0, 32, 32), image.Rect(0, 0, 32, 18)},
		{image.Rect(0, 0, 1080, 1920), image.Rect(0, 0, 32, 32), image.Rect(0, 0, 18, 32)},
		{image.Rect(0, 0, 1024, 735), image.Rect(0, 0, 32, 32), image.Rect(0, 0, 32, 22)},
	}
	for _, row := range testData {
		given := row[0]
		other := row[1]
		expected := row[2]
		mapped := mapRatio(given, other)
		if mapped.Dx() != expected.Dx() || mapped.Dy() != expected.Dy() {
			t.Errorf("Expected %dx%d got %dx%d", expected.Dx(), expected.Dy(), mapped.Dx(), mapped.Dy())
		}
	}
}
