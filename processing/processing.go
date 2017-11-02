package processing

import (
	"bitbucket.org/rwirdemann/kontrol/kontrol"
)

func Process(booking kontrol.Booking) {

	if booking.Extras.Typ == "GV" {
		b := kontrol.Booking{
			Amount: -1 * booking.Amount,
			Text:   "GV Entnahme",
			Month:  booking.Month,
			Year:   booking.Year}
		account := kontrol.Accounts[booking.Extras.CostCenter]
		account.Book(b)
	}

	if booking.Extras.Typ == "AR" {
		benefitees := stakeholderWithNetPositions(booking)
		for _, benefited := range benefitees {

			// todo fix this: Externe bekommen keinen Sharem, Angestellte bekommen einen anderen Share
			if benefited.Type == kontrol.STAKEHOLDER_TYPE_PARTNER {
				b := kontrol.Booking{
					Amount: booking.Extras.Net[benefited] * kontrol.PartnerShare,
					Typ:    kontrol.Nettoanteil,
					Text:   booking.Text,
					Month:  booking.Month,
					Year:   booking.Year}
				account := kontrol.Accounts[benefited.Id]
				account.Bookings = append(kontrol.Accounts[benefited.Id].Bookings, b)
			}
		}
	}
}

// Eine Buchung kann mehrere Nettopositionen enthalten, den je einem Stakeholder zugeschrieben wird.
// Diese Funktion liefert ein Array mit Stateholder, deren Nettoanteil in der Buchung > 0 ist.
func stakeholderWithNetPositions(booking kontrol.Booking) []kontrol.Stakeholder {
	var result []kontrol.Stakeholder

	for k, v := range booking.Extras.Net {
		if v > 0 {
			result = append(result, k)
		}
	}
	return result
}
