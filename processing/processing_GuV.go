package processing

import (
	"github.com/ahojsenn/kontrol/account"
	"github.com/ahojsenn/kontrol/accountSystem"
	"github.com/ahojsenn/kontrol/booking"
	"github.com/ahojsenn/kontrol/util"
	"github.com/ahojsenn/kontrol/valueMagnets"
	"log"
	"strconv"
	"time"
)

func GuV (as accountSystem.AccountSystem) float64 {
	jahresueberschussVs, gwsGezahlt := ermitteleJahresueberschussGWsteuer(as)
	gwsRück := calculateGewerbesteuerRueckstellung (jahresueberschussVs, gwsGezahlt)
	bucheGewerbesteuer (as, gwsRück)
	jahresueberschuss := jahresueberschussVs  - gwsRück
	bucheJahresueberschuss (as, jahresueberschuss  )
	log.Printf("	Gewinn nach Steuer: %+9.2f€\n", jahresueberschuss)
	return jahresueberschuss
}


func ermitteleJahresueberschussGWsteuer (as accountSystem.AccountSystem) (float64, float64) {
	var jahresueberschuss, gwsteuer, ertrag, aufwand float64

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
	log.Printf("	Ertrag:  %+9.2f€\n", ertrag)
	log.Printf("	Aufwand [vor Gewerbesteuer Korrektur]: %+9.2f€\n", aufwand)
	log.Printf("	Gewerbesteuer gebucht: %+9.2f€\n", gwsteuer)
	jahresueberschuss = ertrag + aufwand
	log.Printf("	Gewinn vor GWRücks: %+9.2f€\n", jahresueberschuss)
	return jahresueberschuss, gwsteuer
}


func calculateGewerbesteuerRueckstellung (jahresueberschuss float64, gwsteuer float64) float64 {
	gwsRück := 0.0
	// calculate Gewerbesteuer
	// only do that for the current year!
	if  !util.Global.JahresAbschluss_done { // JM overwrite for now until jahresabschluss
		log.Println("in calculateGewerbesteuerRueckstellung, Jahresabschluss not done yet ==> I will estimate the Gewerbesteuer.")
		gwsRück = berechneGewerbesteuer(jahresueberschuss-gwsteuer) + gwsteuer
	}
	log.Printf("	GWsteuerrückstellung: %7.2f€", gwsRück)
	return gwsRück
}


//
func bucheGewerbesteuer (as accountSystem.AccountSystem, gwsRück float64 )  {
	if (gwsRück != 0.0) {
		now := time.Now().AddDate(0, 0, 0)
		gws := booking.NewBooking(0,booking.CC_GWSteuer, "4320", "956", "K", "",nil,  gwsRück, ("in kontrol ermittelte Gewerbesteuer-Rückstellung "+strconv.Itoa(util.Global.FinancialYear)), int(now.Month()), now.Year(), now)
		// bookFromTo( *gws, gwsKonto, gwsGegenKonto)
		Process(as, *gws)

		// costenter bookings
		gws.Amount *= -1  // but why??? --> Costcenterbookings are negative...
		BookCostToCostCenter{AccSystem: as, Booking: *gws}.run()

		// ermittelte GWSteuer Rückstellung von jahresueberschuss abziehen
		log.Printf("	Gewerbesteuer-Rückstellung  %7.2f€", gwsRück)
	} else {
		log.Printf("in bucheGewerbesteuer, gwsRück == 0 ==> nothing booked...")
	}
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

