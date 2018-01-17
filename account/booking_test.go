package account

import (
	"testing"
	"bitbucket.org/rwirdemann/kontrol/util"
)

// Eine aus der CSV-Datei imporierte Buchung muss nicht zwingend über das Bankkonto in System gekommen sein.
// Ein Beispiel sind interne Stunden, die nie als Zahlung eingegangen sind.
func TestBookOnBankAccount(t *testing.T) {
	is := NewBooking("IS", "K", nil, 0, "Booking", 1, 2017)
	util.AssertFalse(t, is.BookOnBankAccount())

	gv := NewBooking("GV", "RW", nil, 0, "Booking", 1, 2017)
	util.AssertTrue(t, gv.BookOnBankAccount())
}
