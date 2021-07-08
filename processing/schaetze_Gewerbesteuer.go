package processing

import (
	"github.com/ahojsenn/kontrol/accountSystem"
	"github.com/ahojsenn/kontrol/booking"
	"github.com/ahojsenn/kontrol/util"
	"math"
	"strconv"
	"time"
)

const Gewerbesteuer_Freibetrag 	= 24500.0
const Gewerbesteuer_Hebesatz	= 4.70
const Gewerbesteuer_Messbetrtag = 0.035

func berechneGewerbesteuer(gewinn_vorSteuer float64) float64 {
	if gewinn_vorSteuer <= Gewerbesteuer_Freibetrag {
		return 0.0
	}
	gewinn_vorSteuer -= Gewerbesteuer_Freibetrag
	gewinn_vorSteuer = math.Round(gewinn_vorSteuer /100)*100
	return  gewinn_vorSteuer * Gewerbesteuer_Messbetrtag * Gewerbesteuer_Hebesatz
}

//
func bucheGewerbesteuer (as accountSystem.AccountSystem, gwsRück float64 )  {
	if (gwsRück != 0.0) {
		now := time.Now().AddDate(0, 0, 0)
		gws := booking.NewBooking(0,
			booking.CC_GWSteuer, "4320", "956",
			"K", "",nil,  gwsRück,
			("in kontrol ermittelte Gewerbesteuer-Rückstellung "+strconv.Itoa(util.Global.FinancialYear)),
			int(now.Month()), now.Year(), now)
		// book to GuV
		Process(as, *gws)

		// costenter bookings
		gws.Amount *= -1  // but why??? --> Costcenterbookings are negative...
		BookCostToCostCenter{AccSystem: as, Booking: *gws}.run()

		// ermittelte GWSteuer Rückstellung von jahresueberschuss abziehen
		// log.Printf("	Gewerbesteuer-Rückstellung  %7.2f€", gwsRück)
	} else {
		//log.Printf("in bucheGewerbesteuer, gwsRück == 0 ==> nothing booked...")
	}
}

