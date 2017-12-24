package processing

import (
	"bitbucket.org/rwirdemann/kontrol/account"
	"bitbucket.org/rwirdemann/kontrol/owner"
	"bitbucket.org/rwirdemann/kontrol/util"
)

func Process(repository account.Repository, booking account.Booking) {
	switch booking.Extras.Typ {
	case "GV":
		bookPartnerWithdrawal(repository, booking)
	case "AR":
		bookOutgoingInvoice(repository, booking)
	case "ER":
		bookIncomingInvoice(repository, booking)
	case "IS":
		bookInternalHours(repository, booking)
	}
}

func bookPartnerWithdrawal(repository account.Repository, booking account.Booking) {
	if booking.Extras.Typ == "GV" {
		b := account.Booking{
			Amount: -1 * booking.Amount,
			Typ:    account.Entnahme,
			Text:   "GV Entnahme",
			Month:  booking.Month,
			Year:   booking.Year}
		account, _ := repository.Get(booking.Extras.CostCenter)
		account.Book(b)
	}
}

func bookOutgoingInvoice(repository account.Repository, booking account.Booking) {
	benefitees := stakeholderWithNetPositions(booking)
	for _, benefited := range benefitees {

		if benefited.Type == owner.StakeholderTypePartner {

			// book partner share
			b := account.Booking{
				Amount: booking.Extras.Net[benefited] * owner.PartnerShare,
				Typ:    account.Nettoanteil,
				Text:   booking.Text + "#NetShare#" + benefited.Id,
				Month:  booking.Month,
				Year:   booking.Year}
			a, _ := repository.Get(benefited.Id)
			a.Book(b)

			// book kommitment share
			kommitmentShare := account.Booking{
				Amount: booking.Extras.Net[benefited] * owner.KommmitmentShare,
				Typ:    account.Kommitmentanteil,
				Text:   booking.Text + "#Kommitment#" + benefited.Id,
				Month:  booking.Month,
				Year:   booking.Year}

			kommitmentAccount, _ := repository.Get(owner.StakeholderKM.Id)
			kommitmentAccount.Book(kommitmentShare)
		}

		if benefited.Type == owner.StakeholderTypeExtern {

			// book kommitment share
			kommitmentShare := account.Booking{
				Amount: booking.Extras.Net[benefited] * owner.KommmitmentExternShare,
				Typ:    account.Kommitmentanteil,
				Text:   booking.Text + "#Kommitment#" + benefited.Id,
				Month:  booking.Month,
				Year:   booking.Year}
			kommitmentAccount, _ := repository.Get(owner.StakeholderKM.Id)
			kommitmentAccount.Book(kommitmentShare)
		}

		if benefited.Type == owner.StakeholderTypeEmployee {

			// 100% net is booked to employee account to see how much money is made by him
			b := account.Booking{
				Amount: booking.Extras.Net[benefited],
				Typ:    account.Nettoanteil,
				Text:   booking.Text + "#NetShare#" + benefited.Id,
				Month:  booking.Month,
				Year:   booking.Year}
			a, _ := repository.Get(benefited.Id)
			a.Book(b)

			// book kommitment share
			kommitmentShare := account.Booking{
				Amount: booking.Extras.Net[benefited] * owner.KommmitmentEmployeeShare,
				Typ:    account.Kommitmentanteil,
				Text:   booking.Text + "#Kommitment#" + benefited.Id,
				Month:  booking.Month,
				Year:   booking.Year}
			kommitmentAccount, _ := repository.Get(owner.StakeholderKM.Id)
			kommitmentAccount.Book(kommitmentShare)
		}

		// book cost center provision
		a, _ := repository.Get(booking.Extras.CostCenter)
		b := account.Booking{
			Amount: booking.Extras.Net[benefited] * owner.PartnerProvision,
			Typ:    account.Vertriebsprovision,
			Text:   booking.Text + "#Provision#" + benefited.Id,
			Month:  booking.Month,
			Year:   booking.Year}
		a.Book(b)
	}
}

func bookIncomingInvoice(repository account.Repository, booking account.Booking) {
	kommitmentShare := account.Booking{
		Amount: util.Net(booking.Amount) * -1,
		Typ:    account.Eingangsrechnung,
		Text:   booking.Text,
		Month:  booking.Month,
		Year:   booking.Year}
	kommitmentAccount, _ := repository.Get(owner.StakeholderKM.Id)
	kommitmentAccount.Book(kommitmentShare)
}

// Interne Stunden werden direkt netto verbucht
func bookInternalHours(repository account.Repository, booking account.Booking) {
	// Buchung aufs Partner-Konto
	b := account.Booking{
		Amount: booking.Amount,
		Typ:    account.InterneStunden,
		Text:   booking.Text,
		Month:  booking.Month,
		Year:   booking.Year}
	a, _ := repository.Get(booking.Extras.CostCenter)
	a.Book(b)

	// Gegenbuchung Kommitment-Konto
	counterBooking := account.Booking{
		Amount: booking.Amount * -1,
		Typ:    account.InterneStunden,
		Text:   booking.Text,
		Month:  booking.Month,
		Year:   booking.Year}
	kommitmentAccount, _ := repository.Get(owner.StakeholderKM.Id)
	kommitmentAccount.Book(counterBooking)
}

// Eine Buchung kann mehrere Nettopositionen enthalten, den je einem Stakeholder zugeschrieben wird.
// Diese Funktion liefert ein Array mit Stateholder, deren Nettoanteil in der Buchung > 0 ist.
func stakeholderWithNetPositions(booking account.Booking) []owner.Stakeholder {
	var result []owner.Stakeholder

	for k, v := range booking.Extras.Net {
		if v > 0 {
			result = append(result, k)
		}
	}
	return result
}
