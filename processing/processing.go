package processing

import (
	"github.com/ahojsenn/kontrol/account"
	"github.com/ahojsenn/kontrol/booking"
		"github.com/ahojsenn/kontrol/accountSystem"
	"github.com/ahojsenn/kontrol/valueMagnets"
	"github.com/ahojsenn/kontrol/util"
	"log"
	"strconv"
	"time"
)

type Command interface {
	run()
}

// Implementiert den Kommitment-Verteilungsalgorithmus
func Process(accsystem accountSystem.AccountSystem, booking booking.Booking) {

	// Assign booking to one or more virtual stakeholder accounts
	var command Command

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
	case "SKR03":
		command = BookSKR03Command{AccSystem: accsystem, Booking: booking}
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


	now := time.Now().AddDate(0, 0, 0)
	// Jarhesüberschuss ist nun ermittelt

	// Buchung auf Verrechnungskonto Jahresüberschuss
	jue,_ := as.Get(accountSystem.ErgebnisNachSteuern.Id)
	soll := booking.NewBooking(0,"Jahresüberschuss "+strconv.Itoa(util.Global.FinancialYear), "", "", "", nil,  jahresueberschuss, "Buchung Jahresüberschuss", int(now.Month()), now.Year(), now)
	jue.Book(*soll)

	// und Buchung auf Verbindlichkeitenkonto
	verb,_ := as.Get(accountSystem.SKR03_920_Gesellschafterdarlehen.Id)
	haben := booking.NewBooking(0,"Jahresüberschuss "+strconv.Itoa(util.Global.FinancialYear), "", "", valueMagnets.StakeholderKM.Id, nil,  jahresueberschuss, "Buchung Jahresüberschuss", int(now.Month()), now.Year(), now)
	verb.Book(*haben)

}

func Bilanz (as accountSystem.AccountSystem) {

	var konto *account.Account
	var bk *booking.Booking
	now := time.Now().AddDate(0, 0, 0)


	// Aktiva
	for rownr, account := range as.All() {
		if account.Description.Type == accountSystem.KontenartAktiv {
			// Buchung auf SummeAktiva
			konto,_ = as.Get(accountSystem.SummeAktiva.Id)
			bk = booking.NewBooking(
				rownr,
				account.Description.Name+strconv.Itoa(util.Global.FinancialYear),
				"",
				"",
				"",
				nil,
				account.Saldo,
				"SummeAktiva "+strconv.Itoa(util.Global.FinancialYear),
				int(now.Month()),
				now.Year(),
				now)
			konto.Book(*bk)
		}
	}


	// Passiva
	for rownr, account := range as.All() {
		if account.Description.Type == accountSystem.KontenartPassiv {
			// Buchung auf SummePassiva
			konto,_ = as.Get(accountSystem.SummePassiva.Id)
			bk = booking.NewBooking(
				rownr,
				account.Description.Name+strconv.Itoa(util.Global.FinancialYear),
				"",
				"",
				"",
				nil,
				account.Saldo,
				"SummePassiva "+strconv.Itoa(util.Global.FinancialYear),
				int(now.Month()),
				now.Year(),
				now)
			konto.Book(*bk)
		}
	}




}






func ErloesverteilungAnProfitcenter (as accountSystem.AccountSystem) {

	// Kosten werden auf Profitcenter zugeordnet

	// Angestellten Boni werden verteils

	// zunächst werden 15% (oder 20%) davon werden an die Shareholder (gemäß  fairshare) verteilt

	// etwa 5% (1/4 des jeweiligen Deckungsbeitrags) der  Sales comission wird beglichen,

	// interne Stunden werden beglichen

	// Rest wird ermittelt und nach Umsatzanteil verteilt


}
