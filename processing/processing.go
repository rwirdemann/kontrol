package processing

import (
	"github.com/ahojsenn/kontrol/booking"
		"github.com/ahojsenn/kontrol/accountSystem"
	"log"
	"time"
)

type Command interface {
	run()
}

// Implementiert den Kommitment-Verteilungsalgorithmus
func Process(accsystem accountSystem.AccountSystem, booking booking.Booking) {

	// Assign booking to one or more virtual stakeholder accounts
	var command Command

	// if Soll and Haben are filled, then book according to the account values given
	if booking.Soll != "" && booking.Haben != "" {
		// log.Println("booking ", booking.Amount, "€ from ", booking.Soll, " to ", booking.Haben)
		// find the right soll account
		command = BookSKR03Command{AccSystem: accsystem, Booking: booking}

	} else {
		// otherwise use booking.Typ
		switch booking.Typ {
		case "GV":
			command = BookPartnerEntnahmeCommand{AccSystem: accsystem, Booking: booking}
		case "GV-Vorjahr":
			command = BookPartnerEntnahmeVorjahrCommand{AccSystem: accsystem, Booking: booking}
		case "AR":
			command = BookAusgangsrechnungCommand{AccSystem: accsystem, Booking: booking}
		case "ER":
			command = BookEingangsrechnungCommand{AccSystem: accsystem, Booking: booking}
		case "IS":
			command = BookInterneStundenCommand{AccSystem: accsystem, Booking: booking}
		case "SV-Beitrag":
			command = BookSVBeitragCommand{AccSystem: accsystem, Booking: booking}
		case "GWSteuer":
			command = BookGWSteuerCommand{AccSystem: accsystem, Booking: booking}
		case "Gehalt":
			command = BookGehaltCommand{AccSystem: accsystem, Booking: booking}
		case "LNSteuer":
			command = BookLNSteuerCommand{AccSystem: accsystem, Booking: booking}
		case "UstVZ":
			command = BookUstCommand{AccSystem: accsystem, Booking: booking}
				}
	}
	command.run()
}

func GuV (as accountSystem.AccountSystem) {
	log.Println("in GuV")

	var jahresueberschuss float64

	for _, account := range as.All() {
		if account.Description.Type == accountSystem.KontenartAufwand ||  account.Description.Type == accountSystem.KontenartErtrag {
			jahresueberschuss += account.Saldo
		}
	}
	a,_ := as.Get(accountSystem.ErgebnisNachSteuern.Id)
	now := time.Now().AddDate(0, 0, 0)
	p := booking.NewBooking(0,"Jahresüberschuss", "", "", "", nil,  jahresueberschuss, "Buchung Jahresüberschuss", int(now.Month()), now.Year(), now)
	a.Book(*p)
}
