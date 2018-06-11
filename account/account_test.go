package account

import (
	"testing"

	"github.com/ahojsenn/kontrol/booking"
	"github.com/ahojsenn/kontrol/util"
)

func TestSaldo(t *testing.T) {
	a := Account{}
	a.Book(booking.Booking{Amount: 12.55})
	a.UpdateSaldo()
	util.AssertEquals(t, 12.55, a.Saldo)

	a.Book(booking.Booking{Amount: 15.57})
	a.UpdateSaldo()
	util.AssertEquals(t, 28.12, a.Saldo)
}
