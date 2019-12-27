package processing

import (
	"github.com/ahojsenn/kontrol/accountSystem"
	"github.com/ahojsenn/kontrol/booking"
	"github.com/ahojsenn/kontrol/parser"
	"github.com/ahojsenn/kontrol/util"
	"log"
	"testing"
	"time"
)


// 100% werden auf das Bankkonto gebucht
// 100% werden gegen das Kommitment-Konto gebucht. Diese Regel ist nicht unscharf:
// eigentlich müssen die 100% aufgeteilt werden auf: 70% auf Partner, 25% auf
// Kommitment und 5% auf Dealbringer
func TestProcessGWSteuer(t *testing.T) {
	setUp()

	its2018 := time.Date(2018, 1, 23, 0, 0, 0, 0, time.UTC)
	b := booking.NewBooking(13,"GWSteuer", "", "", "K", "Project-X",nil, 2385.10, "STEUERKASSE HAMBURG STEUERNR 048/638/01147 GEW.ST 4VJ.17", 9, 2017, its2018)

	Process(accSystem, *b)


	// Buchung wurde gegen Gewerbesteuer Konto gebucht
	a, _ := accSystem.Get(accountSystem.SKR03_Steuern.Id)
	b1 := a.Bookings[0]
	util.AssertFloatEquals(t, b1.Amount,-2385.10 )

	// Buchung wurde aufs Bankkonto gebucht
	acc,_ := accSystem.Get(accountSystem.SKR03_1200.Id)
	util.AssertEquals(t, 1, len(acc.Bookings))
	actual := acc.Bookings[0]
	util.AssertFloatEquals(t, actual.Amount,2385.10 )
}

func TestGuV(t *testing.T) {
	util.Global.FinancialYear = 2018
	util.Global.Filename = "/Users/docjoe/mystuff/development/kontrol/Buchungen-KG.csv"
	as := accountSystem.NewDefaultAccountSystem()
	hauptbuch := as.GetCollectiveAccount()
	log.Println("in TestGuV: ", util.Global.FinancialYear, len(hauptbuch.Bookings), util.Global.Filename)

	parser.Import(util.Global.Filename, util.Global.FinancialYear, "*",&(hauptbuch.Bookings))

	for _, p := range hauptbuch.Bookings {
		Process(as, p)
	}

	// distribute revenues and costs to valueMagnets
	// in this step only employees revenues will be booked to employee cost centers
	// partners reneue will bi primarily booked to company account for this step
	ErloesverteilungAnEmployees(as)
	// now employee bonusses are calculated and booked
	CalculateEmployeeBonus(as)

	// now (after employee bonusses are booked) calculate GuV and Bilanz
	jahresüberschuss := GuV(as)
	util.AssertEquals(t, jahresüberschuss, GuV(as))
	util.AssertEquals(t, jahresüberschuss, GuV(as))
	util.AssertEquals(t, jahresüberschuss, GuV(as))
}



