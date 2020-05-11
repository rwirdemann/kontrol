package processing

import (
	"github.com/ahojsenn/kontrol/account"
	"github.com/ahojsenn/kontrol/accountSystem"
	"github.com/ahojsenn/kontrol/booking"
	"github.com/ahojsenn/kontrol/util"
	"github.com/ahojsenn/kontrol/valueMagnets"
	"log"
	"strconv"
	"strings"
	"time"
)

func GuV (as accountSystem.AccountSystem) float64 {

	jahresueberschuss := ermitteleJahresueberschuss(as)
	bucheJahresueberschuss (as, jahresueberschuss  )
	log.Printf("	Gewinn nach Steuer: %+9.2f€\n", jahresueberschuss)
	return jahresueberschuss
}


func ermitteleJahresueberschuss (as accountSystem.AccountSystem) (float64) {
	var gewinn_vorSteuer, gwsteuer_gezahlt, ertrag, aufwand, ausgleichsbuchung float64

	for _, acc := range as.All() {
		if (acc.Description.Id == accountSystem.SKR03_Steuern.Id) {
			// Ausnahme machen für dien Gewerbesteuerkorrekturbuchumg am Jahresende
			// die wird zum Ausgleich der Gewerbesteuerberechnung
			// zwischen Steuerberater und kommitment gebraucht
			// Kennzeichen ist
			for _,booking := range acc.Bookings {
				if strings.Contains(booking.Text, "Ausgleichsbuchung") { // true
					ausgleichsbuchung += booking.Amount
				}
			}
			gwsteuer_gezahlt += acc.Saldo
		}
		switch  {
		case acc.Description.Type == account.KontenartAufwand:
			aufwand += acc.Saldo
		case acc.Description.Type == account.KontenartErtrag:
			ertrag += acc.Saldo
		default:
		}
	}
	// ausgleichsbuchung muss aus zahlungen und Aufwand wieder raus...
	gwsteuer_gezahlt -= ausgleichsbuchung
	aufwand -= ausgleichsbuchung

	log.Printf("	Ertrag:  %+9.2f€\n", ertrag)
	log.Printf("	Aufwand [vor Gewerbesteuer Korrektur]: %+9.2f€\n", aufwand -gwsteuer_gezahlt)
	gewinn_vorSteuer = ertrag + aufwand - gwsteuer_gezahlt
	log.Printf("	Gewinn vor Steuer: %+9.2f€\n", gewinn_vorSteuer)

	log.Printf("	Gewerbesteuer gezahlt(ohne ausgleichsbuchung): %+9.2f€\n", gwsteuer_gezahlt)
	// die ausgleichsbuchung sollte bei der Gewerbesteuer Ermittlung nicht berücksichtig werden.
	gewerbeSteuerSchätzungKontrol := -1.0*berechneGewerbesteuer (gewinn_vorSteuer)
	log.Printf("	Gewerbesteuer geschätzt (kontrol): %+9.2f€\n", gewerbeSteuerSchätzungKontrol)
	gwsRück := gwsteuer_gezahlt - gewerbeSteuerSchätzungKontrol
	log.Printf("	Gewerbesteuer Rückstellung (kontrol): %+9.2f€\n", gwsRück)
	bucheGewerbesteuer (as, gwsRück)
	log.Printf("	GWSteuer Rückstellung aus kontrol: %+9.2f€\n", gwsRück)
	log.Printf("	GWSteuer Ausgleichsbuchung: %+9.2f€\n", ausgleichsbuchung)
	gwSteuer := gwsteuer_gezahlt +ausgleichsbuchung-gwsRück
	log.Printf("	Gewerbesteuergezahlt-Ausgleichsbuchung-GWSteuer Rückst.: %+9.2f€\n", gwSteuer)

	return gewinn_vorSteuer + gwSteuer
}

//
func bucheJahresueberschuss (as accountSystem.AccountSystem, jahresueberschuss float64 ) float64 {

	now := time.Now().AddDate(0, 0, 0)

	// Buchung auf Verrechnungskonto Jahresüberschuss
	jue,okay := as.Get(accountSystem.ErgebnisNachSteuern.Id)
	if !okay {
		log.Panic("in GuV, there is no accountSystem.ErgebnisNachSteuern.Id")
	}
	soll := booking.NewBooking(0,"Jahresüberschuss "+strconv.Itoa(util.Global.FinancialYear), "", "", "", "",nil,  -jahresueberschuss, "Buchung Jahresüberschuss", int(now.Month()), now.Year(), now)
	jue.Book(*soll)

	// und Buchung auf Verbindlichkeitenkonto
	verb,okay := as.Get(accountSystem.SKR03_920_Gesellschafterdarlehen.Id)
	if !okay {
		log.Panic("in GuV, there is no accountSystem.SKR03_920_Gesellschafterdarlehen.Id")
	}
	haben := booking.NewBooking(0,"Jahresüberschuss "+strconv.Itoa(util.Global.FinancialYear), "", "", valueMagnets.StakeholderKM.Id, "", nil,  jahresueberschuss, "Buchung Jahresüberschuss", int(now.Month()), now.Year(), now)
	verb.Book(*haben)

	return jahresueberschuss

}

