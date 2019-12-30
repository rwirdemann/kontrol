package processing

import (
	"github.com/ahojsenn/kontrol/account"
	"github.com/ahojsenn/kontrol/accountSystem"
	"log"
)

func Kostenerteilung (as accountSystem.AccountSystem) {
	sumOfCosts := 0.0
	//log.Println("in Kostenverteilung, got as with n bookings, n=", len(as.GetCollectiveAccount().Bookings ))
	// precondition: alle bookings from the general ledger havbe been processed and are now booked to GuC accounts
	// loop though all GuV relevant accounts,
	// i.e. Aufwand...
	// for every stakeholder
	// book all costs
	// to cost center accounts

	// loop through all accounts in accountSystem,
	for _, acc := range as.All() {
		a, _ := as.Get(acc.Description.Id)
		// filter out non GuV accounts
		isCostAccount := acc.Description.Type == account.KontenartAufwand
		if !isCostAccount {
			continue
		}
		// process bookings on GuV accounts
		// log.Println("	account:", a.Description.Id, len(a.Bookings) )
		for _, bk := range a.Bookings {
			bk.Text = "autom. Kostenvert.: " + bk.Text
			BookCostToCostCenter{AccSystem: as, Booking: bk}.run()
			sumOfCosts += bk.Amount
		}
	}
	log.Printf("Kostenverteilung: distributed %7.2f € ", sumOfCosts)
	return
}



