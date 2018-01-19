package processing

import (
	"log"

	"bitbucket.org/rwirdemann/kontrol/account"
	"bitbucket.org/rwirdemann/kontrol/owner"
	"bitbucket.org/rwirdemann/kontrol/util"
	"bitbucket.org/rwirdemann/kontrol/booking"
)

// Implementiert den Kommitment-Verteilungsalgorithmus
func Process(repository account.Repository, booking booking.Booking) {

	// Book booking to bank account
	b := booking
	b.Type = booking.Typ
	switch b.Type {
	case "ER":
		b.Amount = util.Net(b.Amount) * -1
	case "AR":
		b.Amount = util.Net(b.Amount)
	case "GV", "GWSteuer":
		b.Amount = b.Amount * -1
	}

	// Interne Stunden werden nicht auf dem Bankkonto verbucht. Sie sind da nie eingegangen, sondern werden durch
	// Einnahmen bestritten
	if b.BookOnBankAccount() && b.Type != "Gehalt" && b.Type != "SV-Beitrag" {
		repository.BankAccount().Book(b)
	}

	// Assign booking to one or more virtual stakeholder accounts
	switch booking.Typ {
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
	case "GWSteuer":
		bookGWSteuer(repository, booking)
	case "Gehalt":
		bookGehalt(repository, booking)
	default:
		log.Printf("could not process booking type '%s'", booking.Typ)
	}
}

func bookPartnerWithdrawal(repository account.Repository, sourceBooking booking.Booking) {
	if sourceBooking.Typ == "GV" {
		b := booking.Booking{
			Amount: -1 * sourceBooking.Amount,
			Type:   booking.Entnahme,
			Text:   "GV Entnahme",
			Month:  sourceBooking.Month,
			Year:   sourceBooking.Year}
		a, _ := repository.Get(sourceBooking.Responsible)
		a.Book(b)
	}
}

func bookOutgoingInvoice(repository account.Repository, sourceBooking booking.Booking) {
	benefitees := stakeholderWithNetPositions(sourceBooking)
	for _, benefited := range benefitees {

		if benefited.Type == owner.StakeholderTypePartner {

			// book partner share
			b := booking.Booking{
				Amount: sourceBooking.Net[benefited] * owner.PartnerShare,
				Type:   booking.Nettoanteil,
				Text:   sourceBooking.Text + "#NetShare#" + benefited.Id,
				Month:  sourceBooking.Month,
				Year:   sourceBooking.Year}
			a, _ := repository.Get(benefited.Id)
			a.Book(b)

			// book kommitment share
			kommitmentShare := booking.Booking{
				Amount: sourceBooking.Net[benefited] * owner.KommmitmentShare,
				Type:   booking.Kommitmentanteil,
				Text:   sourceBooking.Text + "#Kommitment#" + benefited.Id,
				Month:  sourceBooking.Month,
				Year:   sourceBooking.Year}

			kommitmentAccount, _ := repository.Get(owner.StakeholderKM.Id)
			kommitmentAccount.Book(kommitmentShare)
		}

		if benefited.Type == owner.StakeholderTypeExtern {

			// book kommitment share
			kommitmentShare := booking.Booking{
				Amount: sourceBooking.Net[benefited] * owner.KommmitmentExternShare,
				Type:   booking.Kommitmentanteil,
				Text:   sourceBooking.Text + "#Kommitment#" + benefited.Id,
				Month:  sourceBooking.Month,
				Year:   sourceBooking.Year}
			kommitmentAccount, _ := repository.Get(owner.StakeholderKM.Id)
			kommitmentAccount.Book(kommitmentShare)
		}

		if benefited.Type == owner.StakeholderTypeEmployee {

			// book kommitment share
			kommitmentShare := booking.Booking{
				Amount:     sourceBooking.Net[benefited] * owner.KommmitmentEmployeeShare,
				Type:       booking.Kommitmentanteil,
				Text:       sourceBooking.Text,
				Month:      sourceBooking.Month,
				Year:       sourceBooking.Year,
				CostCenter: benefited.Id}
			kommitmentAccount, _ := repository.Get(owner.StakeholderKM.Id)
			kommitmentAccount.Book(kommitmentShare)
		}

		// book cost center provision
		a, _ := repository.Get(sourceBooking.Responsible)
		b := booking.Booking{
			Amount: sourceBooking.Net[benefited] * owner.PartnerProvision,
			Type:   booking.Vertriebsprovision,
			Text:   sourceBooking.Text + "#Provision#" + benefited.Id,
			Month:  sourceBooking.Month,
			Year:   sourceBooking.Year}
		a.Book(b)
	}
}

func bookIncomingInvoice(repository account.Repository, sourceBooking booking.Booking) {
	kommitmentShare := booking.Booking{
		Amount: util.Net(sourceBooking.Amount) * -1,
		Type:   booking.Eingangsrechnung,
		Text:   sourceBooking.Text,
		Month:  sourceBooking.Month,
		Year:   sourceBooking.Year}
	kommitmentAccount, _ := repository.Get(owner.StakeholderKM.Id)
	kommitmentAccount.Book(kommitmentShare)
}

// Interne Stunden werden direkt netto verbucht
func bookInternalHours(repository account.Repository, sourceBooking booking.Booking) {
	// Buchung aufs Partner-Konto
	b := booking.Booking{
		Amount: sourceBooking.Amount,
		Type:   booking.InterneStunden,
		Text:   sourceBooking.Text,
		Month:  sourceBooking.Month,
		Year:   sourceBooking.Year}
	a, _ := repository.Get(sourceBooking.Responsible)
	a.Book(b)

	// Gegenbuchung Kommitment-Konto
	counterBooking := booking.Booking{
		Amount: sourceBooking.Amount * -1,
		Type:   booking.InterneStunden,
		Text:   sourceBooking.Text,
		Month:  sourceBooking.Month,
		Year:   sourceBooking.Year}
	kommitmentAccount, _ := repository.Get(owner.StakeholderKM.Id)
	kommitmentAccount.Book(counterBooking)
}

// SV-Beitrag wird direkt netto gegen das Kommitment-Konto gebucht
func bookSVBeitrag(repository account.Repository, sourceBooking booking.Booking) {
	c := BookSVBeitragCommand{}
	c.run(repository, sourceBooking)
}

// GWSteuer wird direkt netto gegen das Kommitment-Konto gebucht
func bookGWSteuer(repository account.Repository, sourceBooking booking.Booking) {

	// Gegenbuchung Kommitment-Konto
	counterBooking := booking.Booking{
		Amount: sourceBooking.Amount * -1,
		Type:   booking.GWSteuer,
		Text:   sourceBooking.Text,
		Month:  sourceBooking.Month,
		Year:   sourceBooking.Year}
	kommitmentAccount, _ := repository.Get(owner.StakeholderKM.Id)
	kommitmentAccount.Book(counterBooking)
}

// Gehalt wird direkt netto gegen das Kommitment-Konto gebucht
func bookGehalt(repository account.Repository, sourceBooking booking.Booking) {
	c := BookGehaltCommand{}
	c.run(repository, sourceBooking)
}

// Eine Buchung kann mehrere Nettopositionen enthalten, den je einem Stakeholder zugeschrieben wird.
// Diese Funktion liefert ein Array mit Stateholder, deren Nettoanteil in der Buchung > 0 ist.
func stakeholderWithNetPositions(booking booking.Booking) []owner.Stakeholder {
	var result []owner.Stakeholder

	for k, v := range booking.Net {
		if v > 0 {
			result = append(result, k)
		}
	}
	return result
}
