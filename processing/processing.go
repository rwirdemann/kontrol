package processing

import (
	"bitbucket.org/rwirdemann/kontrol/account"
	"bitbucket.org/rwirdemann/kontrol/domain"
	"bitbucket.org/rwirdemann/kontrol/util"
)

func Process(repository account.Repository, booking domain.Booking) {

	// GV Entnahme
	if booking.Extras.Typ == "GV" {
		b := domain.Booking{
			Amount: -1 * booking.Amount,
			Typ:    domain.Entnahme,
			Text:   "GV Entnahme",
			Month:  booking.Month,
			Year:   booking.Year}
		account, _ := repository.Get(booking.Extras.CostCenter)
		account.Book(b)
	}

	// Ausgangsrechnungen
	if booking.Extras.Typ == "AR" {
		benefitees := stakeholderWithNetPositions(booking)
		for _, benefited := range benefitees {

			if benefited.Type == domain.StakeholderTypePartner {

				// book partner share
				b := domain.Booking{
					Amount: booking.Extras.Net[benefited] * domain.PartnerShare,
					Typ:    domain.Nettoanteil,
					Text:   booking.Text + "#NetShare#" + benefited.Id,
					Month:  booking.Month,
					Year:   booking.Year}
				account, _ := repository.Get(benefited.Id)
				account.Book(b)

				// book kommitment share
				kommitmentShare := domain.Booking{
					Amount: booking.Extras.Net[benefited] * domain.KommmitmentShare,
					Typ:    domain.Kommitmentanteil,
					Text:   booking.Text + "#Kommitment#" + benefited.Id,
					Month:  booking.Month,
					Year:   booking.Year}

				kommitmentAccount, _ := repository.Get(domain.StakeholderKM.Id)
				kommitmentAccount.Book(kommitmentShare)
			}

			if benefited.Type == domain.StakeholderTypeExtern {

				// book kommitment share
				kommitmentShare := domain.Booking{
					Amount: booking.Extras.Net[benefited] * domain.KommmitmentExternShare,
					Typ:    domain.Kommitmentanteil,
					Text:   booking.Text + "#Kommitment#" + benefited.Id,
					Month:  booking.Month,
					Year:   booking.Year}
				kommitmentAccount, _ := repository.Get(domain.StakeholderKM.Id)
				kommitmentAccount.Book(kommitmentShare)
			}

			if benefited.Type == domain.StakeholderTypeEmployee {

				// 100% net is booked to employee account to see how much money is made by him
				b := domain.Booking{
					Amount: booking.Extras.Net[benefited],
					Typ:    domain.Nettoanteil,
					Text:   booking.Text + "#NetShare#" + benefited.Id,
					Month:  booking.Month,
					Year:   booking.Year}
				account, _ := repository.Get(benefited.Id)
				account.Book(b)

				// book kommitment share
				kommitmentShare := domain.Booking{
					Amount: booking.Extras.Net[benefited] * domain.KommmitmentEmployeeShare,
					Typ:    domain.Kommitmentanteil,
					Text:   booking.Text + "#Kommitment#" + benefited.Id,
					Month:  booking.Month,
					Year:   booking.Year}
				kommitmentAccount, _ := repository.Get(domain.StakeholderKM.Id)
				kommitmentAccount.Book(kommitmentShare)
			}

			// book cost center provision
			account, _ := repository.Get(booking.Extras.CostCenter)
			b := domain.Booking{
				Amount: booking.Extras.Net[benefited] * domain.PartnerProvision,
				Typ:    domain.Vertriebsprovision,
				Text:   booking.Text + "#Provision#" + benefited.Id,
				Month:  booking.Month,
				Year:   booking.Year}
			account.Book(b)
		}
	}

	if booking.Extras.Typ == "ER" {
		kommitmentShare := domain.Booking{
			Amount: util.Net(booking.Amount) * -1,
			Typ:    domain.Eingangsrechnung,
			Text:   booking.Text,
			Month:  booking.Month,
			Year:   booking.Year}
		kommitmentAccount, _ := repository.Get(domain.StakeholderKM.Id)
		kommitmentAccount.Book(kommitmentShare)
	}
}

// Eine Buchung kann mehrere Nettopositionen enthalten, den je einem Stakeholder zugeschrieben wird.
// Diese Funktion liefert ein Array mit Stateholder, deren Nettoanteil in der Buchung > 0 ist.
func stakeholderWithNetPositions(booking domain.Booking) []domain.Stakeholder {
	var result []domain.Stakeholder

	for k, v := range booking.Extras.Net {
		if v > 0 {
			result = append(result, k)
		}
	}
	return result
}
