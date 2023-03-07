package structs

import "testing"

type example struct {
	Attribute1 string
	Attribute2 string
}

func TestCopyOrZeroValue(t *testing.T) {
	var e *example

	zv := CopyOrZeroValue(e)

	if zv == nil {
		t.Error("CopyOrZeroValue returned nil")
	}

	if zv.Attribute1 != "" || zv.Attribute2 != "" {
		t.Error("CopyOrZeroValue didn't return zero value")
	}

	e2 := &example{Attribute1: "One", Attribute2: "Two"}

	cp := CopyOrZeroValue(e2)

	if cp == nil {
		t.Error("CopyOrZeroValue returned nil")
	}

	if cp == e2 {
		t.Error("CopyOrZeroValue returned reference with same address")
	}

	if cp.Attribute1 != e2.Attribute1 || cp.Attribute2 != e2.Attribute2 {
		t.Error("CopyOrZeroValue didn't correctly copy attributes")
	}
}
