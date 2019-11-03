package processing

import (
	"github.com/ahojsenn/kontrol/account"
	"github.com/ahojsenn/kontrol/booking"
)

// book from source to target account
func bookFromTo (b booking.Booking, source, target *account.Account)  bool {
	// https://www.rechnungswesen-portal.de/Fachinfo/Eigenkapital/Erfolgskonten.html
	// Aufwandskonto: Aufwendungen werden immer im Soll gebucht, da sie das Eigenkapital mindern
	// Ertragskonto:  Erträge werden im Haben gebucht, da sie das Eigenkapital erhöhen.
	if target.Description.Type == account.KontenartAufwand  {
		// Aufwandskonten werden auf der Habenseite im Soll bebucht --> *= -1.0
		// Aktivkonten werden auf der Habenseite im Soll bebucht --> *= -1.0
		//b.Amount *= -1.0
	}
	b.Text += ": "+b.Soll + "-->"+ b.Haben
	b.Amount *= -1.0 // von soll (negativ)
	source.Book (b)
	b.Amount *= -1.0 // nach haben (positiv)
	b.Text += " -->" + source.Description.Id
	target.Book(b)
	return true
}
