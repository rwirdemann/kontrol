package processing

import (
	"bitbucket.org/rwirdemann/kontrol/kontrol"
)

func Process(position kontrol.Position) {

	if position.Typ == "GV" {
		b := kontrol.Booking{Amount: -1 * position.Amount, Text: "GV Entnahme", Month: position.Month, Year: position.Year}
		account := kontrol.Accounts[position.CostCenter]
		account.Book(b)
	}

	if position.Typ == "AR" {
		benefitees := stakeholderWithNetPositions(position)
		for _, benefited := range benefitees {

			// todo: externe kriegen ExternShare statt PartnerShare

			b := kontrol.Booking{Amount: position.Net[benefited] * kontrol.PartnerShare, Text: position.Subject,
				Month: position.Month, Year: position.Year}
			account := kontrol.Accounts[benefited]
			account.Bookings = append(kontrol.Accounts[benefited].Bookings, b)
		}
	}
}

func stakeholderWithNetPositions(position kontrol.Position) []string {
	var result []string

	for k, v := range position.Net {
		if v > 0 {
			result = append(result, k)
		}
	}
	return result
}
