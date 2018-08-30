package booking

import (
	"testing"
	"time"

	"github.com/ahojsenn/kontrol/util"
)

// Eine aus der CSV-Datei imporierte Buchung muss nicht zwingend über das Bankkonto in System gekommen sein.
// Ein Beispiel sind interne Stunden, die nie als Zahlung eingegangen sind.
func TestBookOnBankAccount(t *testing.T) {
	bankCreated, _ := time.Parse(time.RFC822, "01 Jan 17 10:17 UTC")

	is := NewBooking(13,"IS", "800", "1337", "K", nil, 0, "Booking", 1, 2017, bankCreated)
	util.AssertFalse(t, is.BookOnBankAccount())

	gv := NewBooking(13,"GV", "800", "1337", "RW", nil, 0, "Booking", 1, 2017, bankCreated)
	util.AssertTrue(t, gv.BookOnBankAccount())

	start, _ := time.Parse(time.RFC822, "01 Jan 17 10:12 UTC")
	end, _ := time.Parse(time.RFC822, "01 Jan 17 10:18 UTC")
	util.AssertTrue(t, inTimeSpan(start, end, gv.BankCreated))
	util.AssertTrue(t, inTimeSpan(start, end, is.BankCreated))

}

func inTimeSpan(start, end, check time.Time) bool {
	return check.After(start) && check.Before(end)
}
