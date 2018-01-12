package util

import (
	"testing"
	"bitbucket.org/rwirdemann/kontrol/account"
)

func AssertEquals(t *testing.T, expect interface{}, actual interface{}) {
	if expect != actual {
		t.Errorf("wanted: %v, got: %v", expect, actual)
		t.FailNow()
	}
}

func AssertBooking(t *testing.T, b account.Booking, amount float64, text string, destType string) {
	AssertFloatEquals(t, amount, b.Amount)
	AssertEquals(t, text, b.Text)
	AssertEquals(t, destType, b.DestType)
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

func AssertFloatEquals(t *testing.T, expect float64, actual float64) {
	const EPSILON float64 = 0.00000001

	if (actual-expect) >= EPSILON || (expect-actual) >= EPSILON {
		t.Errorf("wanted: %v, got: %v", expect, actual)
		t.FailNow()
	}
}
