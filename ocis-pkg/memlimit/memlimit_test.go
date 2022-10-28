package memlimit

import (
	"testing"
)

func TestSetDifferentEnabled(t *testing.T) {

	err := Set(false, 0.9, "")
	if err != nil {
		t.Fatalf("Set failed on the first time")
	}

	err = Set(true, 0.9, "")
	if err == nil {
		t.Fatalf("Set did not fail on the second time")
	}

	err = Set(false, 0.9, "")
	if err != nil {
		t.Fatalf("Set failed on the third time")
	}
}

func TestSetDifferentRatio(t *testing.T) {

	err := Set(true, 0.1, "")
	if err != nil {
		t.Fatalf("Set failed on the first time")
	}

	err = Set(true, 0.8, "")
	if err == nil {
		t.Fatalf("Set did not fail on the second time")
	}

	err = Set(true, 0.1, "")
	if err != nil {
		t.Fatalf("Set failed on the third time")
	}
}

func TestSetDifferentAmount(t *testing.T) {

	err := Set(true, 0, "1G")
	if err != nil {
		t.Fatalf("Set failed on the first time")
	}

	err = Set(true, 0, "2G")
	if err == nil {
		t.Fatalf("Set did not fail on the second time")
	}

	err = Set(true, 0, "1G")
	if err != nil {
		t.Fatalf("Set failed on the third time")
	}
}

func TestSetAmountWhenRatio(t *testing.T) {

	err := Set(true, 0.5, "")
	if err != nil {
		t.Fatalf("Set failed on the first time")
	}

	err = Set(true, 0, "2G")
	if err == nil {
		t.Fatalf("Set did not fail on the second time")
	}
}

func TestSetRatioWhenAmount(t *testing.T) {

	err := Set(true, 0, "1G")
	if err != nil {
		t.Fatalf("Set failed on the first time")
	}

	err = Set(true, 0.5, "")
	if err == nil {
		t.Fatalf("Set did not fail on the second time")
	}
}

func TestSetInvalidAmount(t *testing.T) {

	err := Set(true, 0, "1foobar")
	if err == nil {
		t.Fatalf("Set did not fail for a non existent data unit")
	}
}

func TestSetInvalidRatio(t *testing.T) {

	err := Set(true, -1, "")
	if err == nil {
		t.Fatalf("Set did not fail for an invalid ratio")
	}

	err = Set(true, 0, "")
	if err == nil {
		t.Fatalf("Set did not fail for an invalid ratio")
	}

	err = Set(true, 2, "")
	if err == nil {
		t.Fatalf("Set did not fail for an invalid ratio")
	}
}
