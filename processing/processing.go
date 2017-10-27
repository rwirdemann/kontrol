package processing

import (
	"bitbucket.org/rwirdemann/kontrol/kontrol"
)

func Process(position kontrol.Position) {
	benefitees := stakeholderWithNetPositions(position)
	for _, benefited := range benefitees {
		b := kontrol.Booking{Amount: position.Net[benefited], Text: position.Subject}
		account := kontrol.Accounts[benefited]
		account.Bookings = append(kontrol.Accounts[benefited].Bookings, b)
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
