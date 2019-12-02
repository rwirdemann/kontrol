package processing

import (
	"github.com/ahojsenn/kontrol/account"
	"github.com/ahojsenn/kontrol/booking"
)

// book from source to target account
func bookFromTo (b booking.Booking, source, target *account.Account)  bool {

	b.Text += ": "+b.Soll + "-->"+ b.Haben
	b.Amount *= -1.0 // von soll (negativ)
	source.Book (b)
	b.Amount *= -1.0 // nach haben (positiv)
	b.Text += " -->" + source.Description.Id
	target.Book(b)
	return true
}
