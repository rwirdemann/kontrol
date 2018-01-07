package processing

import (
	"log"

	"bitbucket.org/rwirdemann/kontrol/account"
	"bitbucket.org/rwirdemann/kontrol/owner"
	"bitbucket.org/rwirdemann/kontrol/util"
)

func Process(repository account.Repository, booking account.Booking) {

	// Book booking to bank account
	b := booking
	b.DestType = booking.Extras.SourceType
	switch b.DestType {
	case "ER":
		b.Amount = util.Net(b.Amount) * -1
	case "AR":
		b.Amount = util.Net(b.Amount)
	case "GV", "SV-Beitrag":
		b.Amount = b.Amount * -1
	}

	// Interne Stunden werden nicht auf dem Bankkonto verbucht
	if b.DestType != "IS" {
		repository.CollectiveAccount().Book(b)
	}

	// Assign booking to one or more virtual stakeholder accounts
	switch booking.Extras.SourceType {
	case "GV":
		bookPartnerWithdrawal(repository, booking)
	case "AR":
		bookOutgoingInvoice(repository, booking)
	case "ER":
		bookIncomingInvoice(repository, booking)
	case "IS":
		bookInternalHours(repository, booking)
	case "SV-Beitrag":
		bookSVBeitrag(repository, booking)
	default:
		log.Printf("could not process booking type '%s'", booking.Extras.SourceType)
	}
}

func bookPartnerWithdrawal(repository account.Repository, booking account.Booking) {
	if booking.Extras.SourceType == "GV" {
		b := account.Booking{
			Amount:   -1 * booking.Amount,
			DestType: account.Entnahme,
			Text:     "GV Entnahme",
			Month:    booking.Month,
			Year:     booking.Year}
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
				Amount:   booking.Extras.Net[benefited] * owner.PartnerShare,
				DestType: account.Nettoanteil,
				Text:     booking.Text + "#NetShare#" + benefited.Id,
				Month:    booking.Month,
				Year:     booking.Year}
			a, _ := repository.Get(benefited.Id)
			a.Book(b)

			// book kommitment share
			kommitmentShare := account.Booking{
				Amount:   booking.Extras.Net[benefited] * owner.KommmitmentShare,
				DestType: account.Kommitmentanteil,
				Text:     booking.Text + "#Kommitment#" + benefited.Id,
				Month:    booking.Month,
				Year:     booking.Year}

			kommitmentAccount, _ := repository.Get(owner.StakeholderKM.Id)
			kommitmentAccount.Book(kommitmentShare)
		}

		if benefited.Type == owner.StakeholderTypeExtern {

			// book kommitment share
			kommitmentShare := account.Booking{
				Amount:   booking.Extras.Net[benefited] * owner.KommmitmentExternShare,
				DestType: account.Kommitmentanteil,
				Text:     booking.Text + "#Kommitment#" + benefited.Id,
				Month:    booking.Month,
				Year:     booking.Year}
			kommitmentAccount, _ := repository.Get(owner.StakeholderKM.Id)
			kommitmentAccount.Book(kommitmentShare)
		}

		if benefited.Type == owner.StakeholderTypeEmployee {

			// 100% net is booked to employee account to see how much money is made by him
			b := account.Booking{
				Amount:   booking.Extras.Net[benefited],
				DestType: account.Nettoanteil,
				Text:     booking.Text + "#NetShare#" + benefited.Id,
				Month:    booking.Month,
				Year:     booking.Year}
			a, _ := repository.Get(benefited.Id)
			a.Book(b)

			// book kommitment share
			kommitmentShare := account.Booking{
				Amount:   booking.Extras.Net[benefited] * owner.KommmitmentEmployeeShare,
				DestType: account.Kommitmentanteil,
				Text:     booking.Text + "#Kommitment#" + benefited.Id,
				Month:    booking.Month,
				Year:     booking.Year}
			kommitmentAccount, _ := repository.Get(owner.StakeholderKM.Id)
			kommitmentAccount.Book(kommitmentShare)
		}

		// book cost center provision
		a, _ := repository.Get(booking.Extras.CostCenter)
		b := account.Booking{
			Amount:   booking.Extras.Net[benefited] * owner.PartnerProvision,
			DestType: account.Vertriebsprovision,
			Text:     booking.Text + "#Provision#" + benefited.Id,
			Month:    booking.Month,
			Year:     booking.Year}
		a.Book(b)
	}
}

func bookIncomingInvoice(repository account.Repository, booking account.Booking) {
	kommitmentShare := account.Booking{
		Amount:   util.Net(booking.Amount) * -1,
		DestType: account.Eingangsrechnung,
		Text:     booking.Text,
		Month:    booking.Month,
		Year:     booking.Year}
	kommitmentAccount, _ := repository.Get(owner.StakeholderKM.Id)
	kommitmentAccount.Book(kommitmentShare)
}

// Interne Stunden werden direkt netto verbucht
func bookInternalHours(repository account.Repository, booking account.Booking) {
	// Buchung aufs Partner-Konto
	b := account.Booking{
		Amount:   booking.Amount,
		DestType: account.InterneStunden,
		Text:     booking.Text,
		Month:    booking.Month,
		Year:     booking.Year}
	a, _ := repository.Get(booking.Extras.CostCenter)
	a.Book(b)

	// Gegenbuchung Kommitment-Konto
	counterBooking := account.Booking{
		Amount:   booking.Amount * -1,
		DestType: account.InterneStunden,
		Text:     booking.Text,
		Month:    booking.Month,
		Year:     booking.Year}
	kommitmentAccount, _ := repository.Get(owner.StakeholderKM.Id)
	kommitmentAccount.Book(counterBooking)
}

// SV-Beitrag wird direkt netto gegen das Kommitment-Konto gebucht
func bookSVBeitrag(repository account.Repository, booking account.Booking) {

	// Gegenbuchung Kommitment-Konto
	counterBooking := account.Booking{
		Amount:   booking.Amount * -1,
		DestType: account.SVBeitrag,
		Text:     booking.Text,
		Month:    booking.Month,
		Year:     booking.Year}
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
