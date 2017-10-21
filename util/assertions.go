package util

import "testing"

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

func AssertNotNil(t *testing.T, actual interface{}) {
	if actual == nil {
		t.Errorf("actial = %v is nil", actual)
		t.FailNow()
	}
}

func AssertNil(t *testing.T, actual interface{}) {
	if actual != nil {
		t.Errorf("actial = %v is not nil", actual)
		t.FailNow()
	}
}
