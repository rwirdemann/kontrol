package processing

import (
	"fmt"
	"github.com/ahojsenn/kontrol/account"
	"github.com/ahojsenn/kontrol/accountSystem"
	"github.com/ahojsenn/kontrol/booking"
	"github.com/ahojsenn/kontrol/util"
	"github.com/ahojsenn/kontrol/valueMagnets"
	"log"
	"math"
	"strconv"
	"time"
)

func GuV (as accountSystem.AccountSystem) {

	var jahresueberschuss, gwsteuer, ertrag, aufwand float64
	now := time.Now().AddDate(0, 0, 0)

	for _, acc := range as.All() {
		if (acc.Description.Id == accountSystem.SKR03_Steuern.Id) {
			gwsteuer += acc.Saldo
		}
		switch  {
			case acc.Description.Type == account.KontenartAufwand:
				aufwand += acc.Saldo
			case acc.Description.Type == account.KontenartErtrag:
				ertrag += acc.Saldo
			default:
		}
	}
	fmt.Printf("			in GuV, Ertrag:  %+9.2f€\n", math.Round(100*ertrag)/100)
	fmt.Printf("			in GuV, Aufwand: %+9.2f€\n", math.Round(100*aufwand)/100)
	jahresueberschuss = ertrag + aufwand
	fmt.Printf("			in GuV, Jahresueberschuss: %+9.2f€\n", math.Round(100*jahresueberschuss)/100)


	// calculate Gewerbesteuer
	// only do that for the current year!
	log.Println("in GuV, Gewerbesteuer gebucht:", gwsteuer)
	log.Println("in GuV, Gewinn vor Steuer:", jahresueberschuss-gwsteuer)
	if  util.Global.FinancialYear == time.Now().Year() {
		log.Println("in GuV, GWsteuer:", berechne_Gewerbesteuer(jahresueberschuss-gwsteuer))
		gwsRück := math.Round( 100* (berechne_Gewerbesteuer(jahresueberschuss-gwsteuer) + gwsteuer ) /100 )

		// ermittelte GWSteuer Rückstellung verbuchen
		gwsKonto,_ := as.Get(accountSystem.SKR03_Steuern.Id)
		gwsGegenKonto,_ := as.Get(accountSystem.SKR03_Rueckstellungen.Id)
		gws := booking.NewBooking(0,"in kontrol ermittelte Gewerbesteuer-Rückstellung "+strconv.Itoa(util.Global.FinancialYear), "4320", "956", "", "",nil,  -gwsRück, ("in kontrol ermittelte Gewerbesteuer-Rückstellung "+strconv.Itoa(util.Global.FinancialYear)), int(now.Month()), now.Year(), now)

		bookFromTo( *gws, gwsKonto, gwsGegenKonto)
		// ermittelte GWSteuer Rückstellung von jahresueberschuss abziehen
		log.Println("in GuV, Gewerbesteuer-Rückstellung", gwsRück)
		jahresueberschuss -= gwsRück

	}

	log.Println("in GuV, Gewinn nach Steuer:", jahresueberschuss)


	// Jahresüberschuss ist nun ermittelt

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

	log.Printf("in GuV, Jahresüberschuss: %6.2f€\n", math.Round(100*jahresueberschuss)/100)
}

