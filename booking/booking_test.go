package booking

import (
	"testing"
	"time"

	"bitbucket.org/rwirdemann/kontrol/util"
)

// Eine aus der CSV-Datei imporierte Buchung muss nicht zwingend über das Bankkonto in System gekommen sein.
// Ein Beispiel sind interne Stunden, die nie als Zahlung eingegangen sind.
func TestBookOnBankAccount(t *testing.T) {
	is := NewBooking("IS", "K", nil, 0, "Booking", 1, 2017, time.Time{}, time.Time{})
	util.AssertFalse(t, is.BookOnBankAccount())

	gv := NewBooking("GV", "RW", nil, 0, "Booking", 1, 2017, time.Time{}, time.Time{})
	util.AssertTrue(t, gv.BookOnBankAccount())
}
