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
func Process(repository accountSystem.AccountSystem, booking booking.Booking) {
	
	// Assign booking to one or more virtual stakeholder accounts
	var command Command

	// if Soll and Haben are filled, then book according to the account values given
	if booking.Soll != "" && booking.Haben != "" {
		// log.Println("booking ", booking.Amount, "€ from ", booking.Soll, " to ", booking.Haben)
		// find the right soll account
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
	p := booking.NewBooking("Jahresüberschuss", "", "", "", nil, jahresueberschuss, "Buchung Jahresüberschuss", int(now.Month()), now.Year(), now)
	a.Book(*p)
}
