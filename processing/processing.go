package processing

import (
	"log"

	"bitbucket.org/rwirdemann/kontrol/account"
	"bitbucket.org/rwirdemann/kontrol/booking"
	"bitbucket.org/rwirdemann/kontrol/owner"
	"bitbucket.org/rwirdemann/kontrol/util"
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
	case "GWSteuer":
		b.Amount = b.Amount * -1
	}

	// Interne Stunden werden nicht auf dem Bankkonto verbucht. Sie sind da nie eingegangen, sondern werden durch
	// Einnahmen bestritten
	if b.BookOnBankAccount() &&
		b.Type != "Gehalt" &&
		b.Type != "SV-Beitrag" &&
		b.Type != "LNSteuer" &&
		b.Type != "GWSteuer" &&
		b.Type != "GV" {
		repository.BankAccount().Book(b)
	}

	// Assign booking to one or more virtual stakeholder accounts
	var command Command
	switch booking.Typ {
	case "GV":
		command = BookPartnerEntnahmeCommand{Repository: repository, Booking: booking}
		command.run()
	case "AR":
		command = BookAusgangsrechnungCommand{Repository: repository, Booking: booking}
		command.run()
	case "ER":
		command = BookEingangsrechnungCommand{Repository: repository, Booking: booking}
		command.run()
	case "IS":
		bookInternalHours(repository, booking)
	case "SV-Beitrag":
		command = BookSVBeitragCommand{Repository: repository, Booking: booking}
		command.run()
	case "GWSteuer":
		command = BookGWSteuerCommand{Repository: repository, Booking: booking}
		command.run()
	case "Gehalt":
		command = BookGehaltCommand{Repository: repository, Booking: booking}
		command.run()
	case "LNSteuer":
		command = BookLNSteuerCommand{Repository: repository, Booking: booking}
		command.run()
	default:
		log.Printf("could not process booking type '%s'", booking.Typ)
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
