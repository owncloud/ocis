package resolutions

import "testing"

func TestParseWithEmptyString(t *testing.T) {
	_, err := Parse("")
	if err == nil {
		t.Error("Parse with empty string should return an error.")
	}
}

func TestParseWithInvalidWidth(t *testing.T) {
	_, err := Parse("invalidx42")
	if err == nil {
		t.Error("Parse with invalid width should return an error.")
	}
}

func TestParseWithInvalidHeight(t *testing.T) {
	_, err := Parse("42xinvalid")
	if err == nil {
		t.Error("Parse with invalid height should return an error.")
	}
}

func TestParse(t *testing.T) {
	rStr := "42x23"
	r, _ := Parse(rStr)
	if r.Width != 42 || r.Height != 23 {
		t.Errorf("Expected resolution %s got %s", rStr, r.String())
	}
}

func TestString(t *testing.T) {
	r := Resolution{Width: 42, Height: 23}
	expected := "42x23"
	if r.String() != expected {
		t.Errorf("Expected string %s got %s", expected, r.String())
	}
}
