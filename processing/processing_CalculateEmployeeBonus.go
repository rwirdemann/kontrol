package processing

import (
	"fmt"
	"github.com/ahojsenn/kontrol/accountSystem"
	"github.com/ahojsenn/kontrol/booking"
	"github.com/ahojsenn/kontrol/util"
	"github.com/ahojsenn/kontrol/valueMagnets"
	"log"
	"time"
)

// loop through the employees, calculates their Bonusses and book them from employee to kommitment cost.
func CalculateEmployeeBonus (as accountSystem.AccountSystem) accountSystem.AccountSystem {
	sumOfBonusses := 0.0
	//log.Println("in CalculateEmployeeBonus: ")
	shrepo := valueMagnets.Stakeholder{}
	for _,sh := range shrepo.GetAllOfType (valueMagnets.StakeholderTypeEmployee) {
		bonus := StakeholderYearlyIncome(as, sh.Id)
		log.Printf("in CalculateEmployeeBonus: %s:  %7.2f€",sh.Id, bonus)

		// take care, this is idempotent, i.e. that the next bonus calculation overwrites the last one...
		if bonus > 0.0 {
			// only book positive bonusses of valuemagnets
			// log.Println("in CalculateEmployeeBonus: ", sh.Id, math.Round(bonus*100)/100)
			// book in into GuV
			bk := booking.Booking{
				RowNr:       0,
				Amount:      bonus,
				Soll:		 		"4120",
				Haben: 		 	"965",
				Type:        booking.CC_Gehalt,
				CostCenter:  sh.Id,
				Text:        fmt.Sprintf("in kontrol kalkulierte Bonusrückstellung für %s in %d", sh.Id, util.Global.FinancialYear),
				Month:       12,
				Year:        util.Global.FinancialYear,
				FileCreated: time.Now().AddDate(0, 0, 0),
				BankCreated: time.Now().AddDate(0, 0, 0),
			}
			// create new booking in GuV Aufwand und Ertragskonten
			Process(as, bk)

			// costenter bookings
			// book from k-ompanie mainaccount to subacc cos
			bk.Amount *= -1  // but why???
			BookCostToCostCenter{AccSystem: as, Booking: bk}.run()

			// now also store the bonus in stakeholders account
			stakeHoldersAccount, _ := as.Get(sh.Id)
			stakeHoldersAccount.Bonus = bonus

			sumOfBonusses += bk.Amount
		}
	}

	log.Printf ("CalculateEmployeeBonus: sumOfBonusses=%7.2f€", sumOfBonusses)
	return as
}



// sum up whatever the stakeholder earned in the actual year
func StakeholderYearlyIncome (as accountSystem.AccountSystem, stkhldr string) float64 {

	a1,_  := as.GetSubacc(stkhldr, accountSystem.UK_Kosten.Id)
	a2,_  := as.GetSubacc(stkhldr, accountSystem.UK_AnteileausFairshare.Id)
	a3,_  := as.GetSubacc(stkhldr, accountSystem.UK_Vertriebsprovision.Id)
	a4,_  := as.GetSubacc(stkhldr, accountSystem.UK_AnteilMitmachen.Id)
	a5,_  := as.GetSubacc(stkhldr, accountSystem.UK_AnteileAuserloesen.Id)
	a6,_  := as.GetSubacc(stkhldr, accountSystem.UK_Erloese.Id)

	log.Printf("in StakeholderYearlyIncome: %s %7.2f€ %7.2f€ %7.2f€", stkhldr, a1.Saldo, a3.Saldo, a5.Saldo)

	return	a1.Saldo +
		a2.Saldo +
		a3.Saldo +
		a4.Saldo +
		a5.Saldo +
		a6.Saldo
}



