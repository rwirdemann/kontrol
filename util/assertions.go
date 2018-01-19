package util

import (
	"testing"
)

func AssertEquals(t *testing.T, expect interface{}, actual interface{}) {
	if expect != actual {
		t.Errorf("wanted: %v, got: %v", expect, actual)
		t.FailNow()
	}
}

func AssertTrue(t *testing.T, b bool) {
	if !b {
		t.Errorf("b = %v is not true", b)
		t.FailNow()
	}
}

func AssertFalse(t *testing.T, b bool) {
	if b {
		t.Errorf("b = %v is true, expected was false", b)
		t.FailNow()
	}
}

func AssertFloatEquals(t *testing.T, expect float64, actual float64) {
	const EPSILON float64 = 0.00000001

	if (actual-expect) >= EPSILON || (expect-actual) >= EPSILON {
		t.Errorf("wanted: %v, got: %v", expect, actual)
		t.FailNow()
	}
}
