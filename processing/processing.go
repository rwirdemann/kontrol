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

			// todo: externe kriegen ExternShare statt PartnerShare

			b := kontrol.Booking{
				Amount: booking.Extras.Net[benefited] * kontrol.PartnerShare,
				Text:   booking.Text,
				Month:  booking.Month,
				Year:   booking.Year}
			account := kontrol.Accounts[benefited]
			account.Bookings = append(kontrol.Accounts[benefited].Bookings, b)
		}
	}
}

func stakeholderWithNetPositions(position kontrol.Booking) []string {
	var result []string

	for k, v := range position.Extras.Net {
		if v > 0 {
			result = append(result, k)
		}
	}
	return result
}
