package processing

import (
	"github.com/ahojsenn/kontrol/account"
	"github.com/ahojsenn/kontrol/accountSystem"
	"log"
)

// book all employees cost and revenues
// 20191228:
// this is performed before GuV thus

func ErloesverteilungAnEmployees (as accountSystem.AccountSystem) {
	sumOfRevs := 0.0
	//log.Println("in ErloesverteilungAnEmployees, got as with n bookings, n=", len(as.GetCollectiveAccount_thisYear().Bookings ))
	// precondition: alle bookings from the general ledger have been processed and are now booked to GuC accounts
	// loop though all GuV relevant accounts,
	// i.e. Ertrag
	// for employees only --> BookRevenueToEmployeeCostCenter
	// to employees cost center accounts

	// loop through all accounts in accountSystem,
	for _, acc := range as.All() {
		a, _ := as.Get(acc.Description.Id)
		// filter out non GuV accounts
		isRevenueAccount := acc.Description.Type == account.KontenartErtrag
		if !isRevenueAccount {
			continue
		}
		// process bookings on GuV accounts
		// log.Println("	account:", a.Description.Id, len(a.Bookings) )
		for _, bk := range a.Bookings {
			BookRevenueToEmployeeCostCenter{AccSystem: as, Booking: bk}.run()
			sumOfRevs += bk.Amount
		}
	}
	log.Printf("ErloesverteilungAnEmployees, distributed %7.2f€", sumOfRevs)
	return
}

