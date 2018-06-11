package processing

import (
	"log"

	"github.com/ahojsenn/kontrol/account"
	"github.com/ahojsenn/kontrol/booking"
	"github.com/ahojsenn/kontrol/util"
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
		b.Type != "Rückstellung" &&
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
		command = BookInterneStundenCommand{Repository: repository, Booking: booking}
		command.run()
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
	case "Rückstellung":
		command = BookRueckstellungCommand{Repository: repository, Booking: booking}
		command.run()
	default:
		log.Printf("could not process booking type '%s'", booking.Typ)
	}
}
