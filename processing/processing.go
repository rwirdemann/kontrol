package processing

import (
	"log"

	"github.com/ahojsenn/kontrol/booking"
	"github.com/ahojsenn/kontrol/util"
	"github.com/ahojsenn/kontrol/accountSystem"
)

type Command interface {
	run()
}

// Implementiert den Kommitment-Verteilungsalgorithmus
func Process(repository accountSystem.AccountSystem, booking booking.Booking) {

	// Book booking to bank account
	b := booking
	b.Type = booking.Typ
	switch b.Type {
	case "ER":
		b.Amount = util.Net(b.Amount) * -1
	case "ERgegenRückstellung":
		b.Amount = util.Net(b.Amount) * -1
	case "AR":
		b.Amount = util.Net(b.Amount)
	case "GWSteuer":
		b.Amount = b.Amount * -1
	}

	// Interne Stunden werden nicht auf dem Bankkonto verbucht. Sie sind da nie eingegangen, sondern werden durch
	// Einnahmen bestritten
	// hier füge ich immer mehr buchungstypen hinzu. Die Bankbuchung sollte immer in command.go stattfinden (resp. ausgangsre*...
	// und nicht hier...
	if b.BookOnBankAccount() &&
		b.Type != "Gehalt" &&
		b.Type != "ER" &&
		b.Type != "SV-Beitrag" &&
		b.Type != "LNSteuer" &&
		b.Type != "GWSteuer" &&
		b.Type != "Rückstellung" &&
		b.Type != "Anfangsbestand" &&
		b.Type != "RückstellungAuflösen" &&
		b.Type != "GV-Vorjahr" &&
		b.Type != "GV" {
		repository.BankAccount().Book(b)
	}

	// Assign booking to one or more virtual stakeholder accounts
	var command Command

	// if Soll and Haben are filled, then book according to the account values given
	if booking.Soll != "" && booking.Haben != "" {
		// log.Println("booking ", booking.Amount, "€ from ", booking.Soll, " to ", booking.Haben)
		// find the right soll account
		log.Println(">>> soll, haben", booking.Soll, booking.Haben)
		command = BookSKR03Command{Repository: repository, Booking: booking}

	} else {
		// otherwise use booking.Typ
		switch booking.Typ {
		case "GV":
			command = BookPartnerEntnahmeCommand{Repository: repository, Booking: booking}
		case "GV-Vorjahr":
			command = BookPartnerEntnahmeVorjahrCommand{Repository: repository, Booking: booking}
		case "AR":
			command = BookAusgangsrechnungCommand{Repository: repository, Booking: booking}
		case "ER":
			command = BookEingangsrechnungCommand{Repository: repository, Booking: booking}
		case "IS":
			command = BookInterneStundenCommand{Repository: repository, Booking: booking}
		case "SV-Beitrag":
			command = BookSVBeitragCommand{Repository: repository, Booking: booking}
		case "GWSteuer":
			command = BookGWSteuerCommand{Repository: repository, Booking: booking}
		case "Gehalt":
			command = BookGehaltCommand{Repository: repository, Booking: booking}
		case "LNSteuer":
			command = BookLNSteuerCommand{Repository: repository, Booking: booking}
		case "Rückstellung":
			command = BookRueckstellungCommand{Repository: repository, Booking: booking}
		case "Anfangsbestand":
			command = BookAnfangsbestandCommand{Repository: repository, Booking: booking}
		case "ERgegenRückstellung":
			command = BookERgegenRückstellungCommand{Repository: repository, Booking: booking}
		case "RückstellungAuflösen":
			command = BookRückstellungAuflösenCommand{Repository: repository, Booking: booking}
		}
	}
	command.run()
}
