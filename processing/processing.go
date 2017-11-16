package processing

import (
	"bitbucket.org/rwirdemann/kontrol/kontrol"
	"bitbucket.org/rwirdemann/kontrol/util"
)

func Process(booking kontrol.Booking) {

	if booking.Extras.Typ == "GV" {
		b := kontrol.Booking{
			Amount: -1 * booking.Amount,
			Typ:    kontrol.Entnahme,
			Text:   "GV Entnahme",
			Month:  booking.Month,
			Year:   booking.Year}
		account := kontrol.Accounts[booking.Extras.CostCenter]
		account.Book(b)
	}

	if booking.Extras.Typ == "AR" {
		benefitees := stakeholderWithNetPositions(booking)
		for _, benefited := range benefitees {

			if benefited.Type == kontrol.StakeholderTypePartner {

				// book partner share
				b := kontrol.Booking{
					Amount: booking.Extras.Net[benefited] * kontrol.PartnerShare,
					Typ:    kontrol.Nettoanteil,
					Text:   booking.Text + "#NetShare#" + benefited.Id,
					Month:  booking.Month,
					Year:   booking.Year}
				account := kontrol.Accounts[benefited.Id]
				account.Book(b)

				// book kommitment share
				kommitmentShare := kontrol.Booking{
					Amount: booking.Extras.Net[benefited] * kontrol.KommmitmentShare,
					Typ:    kontrol.Kommitmentanteil,
					Text:   booking.Text + "#Kommitment#" + benefited.Id,
					Month:  booking.Month,
					Year:   booking.Year}
				kommitmentAccount := kontrol.Accounts[kontrol.SH_KM.Id]
				kommitmentAccount.Book(kommitmentShare)
			}

			if benefited.Type == kontrol.StakeholderTypeExtern {

				// book kommitment share
				kommitmentShare := kontrol.Booking{
					Amount: booking.Extras.Net[benefited] * kontrol.KommmitmentExternShare,
					Typ:    kontrol.Kommitmentanteil,
					Text:   booking.Text + "#Kommitment#" + benefited.Id,
					Month:  booking.Month,
					Year:   booking.Year}
				kommitmentAccount := kontrol.Accounts[kontrol.SH_KM.Id]
				kommitmentAccount.Book(kommitmentShare)
			}

			if benefited.Type == kontrol.StakeholderTypeEmployee {

				// 100% net is booked to employee account to see how much money is made by him
				b := kontrol.Booking{
					Amount: booking.Extras.Net[benefited],
					Typ:    kontrol.Nettoanteil,
					Text:   booking.Text + "#NetShare#" + benefited.Id,
					Month:  booking.Month,
					Year:   booking.Year}
				account := kontrol.Accounts[benefited.Id]
				account.Book(b)

				// book kommitment share
				kommitmentShare := kontrol.Booking{
					Amount: booking.Extras.Net[benefited] * kontrol.KommmitmentEmployeeShare,
					Typ:    kontrol.Kommitmentanteil,
					Text:   booking.Text + "#Kommitment#" + benefited.Id,
					Month:  booking.Month,
					Year:   booking.Year}
				kommitmentAccount := kontrol.Accounts[kontrol.SH_KM.Id]
				kommitmentAccount.Book(kommitmentShare)
			}

			// book cost center provision
			account := kontrol.Accounts[booking.Extras.CostCenter]
			b := kontrol.Booking{
				Amount: booking.Extras.Net[benefited] * kontrol.PartnerProvision,
				Typ:    kontrol.Vertriebsprovision,
				Text:   booking.Text + "#Provision#" + benefited.Id,
				Month:  booking.Month,
				Year:   booking.Year}
			account.Book(b)
		}
	}

	if booking.Extras.Typ == "ER" {
		kommitmentShare := kontrol.Booking{
			Amount: util.Net(booking.Amount) * -1,
			Typ:    kontrol.Eingangsrechnung,
			Text:   booking.Text,
			Month:  booking.Month,
			Year:   booking.Year}
		kommitmentAccount := kontrol.Accounts[kontrol.SH_KM.Id]
		kommitmentAccount.Book(kommitmentShare)
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
