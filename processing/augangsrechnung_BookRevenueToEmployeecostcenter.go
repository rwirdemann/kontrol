package processing

import (
	"fmt"
	"github.com/ahojsenn/kontrol/account"
	"github.com/ahojsenn/kontrol/accountSystem"
	"github.com/ahojsenn/kontrol/booking"
	"github.com/ahojsenn/kontrol/valueMagnets"
	"log"
	"strings"
	"time"
)

type BookRevenueToEmployeeCostCenter struct {
	Booking    booking.Booking
	AccSystem  accountSystem.AccountSystem

}

func (this BookRevenueToEmployeeCostCenter) run() {

	// hier kommt nun die ganze Verteilung unter den kommitmenschen
	// 1. erst employees auszahlen und vertriebsprovisionen, den Rest auf kommitment
	// 2. später den kommitment Topf unter den kommanditisten aufteilen

	benefitees := this.stakeholderWithNetPositions()

	// book the sum of the net positions from the k-main account to k-revenues subaccount
	acc_k_main,_ := this.AccSystem.Get(valueMagnets.StakeholderKM.Id)
	acc_k_subrev,_ := this.AccSystem.GetSubacc(valueMagnets.StakeholderKM.Id, accountSystem.UK_Erloese.Id)

	// if there are no benefitees, then simply book this to from main to subacc
	if len(benefitees) == 0 {
		bookFromTo(this.Booking,  acc_k_main, acc_k_subrev)
	}


	for _, benefited := range benefitees {

		// 1. Step, book everything to k subacc revenue for each benefitee individually
		bk := booking.Booking{
			RowNr:       this.Booking.RowNr,
			Amount:      this.Booking.Net[benefited],
			Type:        booking.CC_RevDistribution_1,
			Project:     this.Booking.Project,
			CostCenter:  benefited.Id,
			Text:        this.Booking.Text + "#Kommitment#" + benefited.Id,
			Month:       this.Booking.Month,
			Year:        this.Booking.Year,
			FileCreated: this.Booking.FileCreated,
			BankCreated: this.Booking.BankCreated}

		bookFromTo(bk,  acc_k_main, acc_k_subrev)


		// Employee revenue =70% from subacc costs to employee
		if benefited.Type == valueMagnets.StakeholderTypeEmployee {

			// book employee share of employees revenue
			employeeshare := booking.Booking{
				RowNr: 		 this.Booking.RowNr,
				Amount:      this.Booking.Net[benefited] * account.EmployeeShare,
				Type:        booking.CC_Employeeaanteil,
				Project:     this.Booking.Project,
				Text:        fmt.Sprintf("%f", account.EmployeeShare)+"*Netto: " + this.Booking.Text+"#"+benefited.Id,
				Month:       this.Booking.Month,
				Year:        this.Booking.Year,
				FileCreated: this.Booking.FileCreated,
				BankCreated: this.Booking.BankCreated,
				CostCenter:  this.Booking.CostCenter}

			// from k-costs account to employees main account
			acc_k_costs,_ := this.AccSystem.GetSubacc(valueMagnets.StakeholderKM.Id, accountSystem.UK_Kosten.Id)
			employeeMainAccount,_ := this.AccSystem.Get(benefited.Id)
			bookFromTo(employeeshare,  acc_k_costs, employeeMainAccount)

			// book from employees main account to employees subaccount
			employeeSubbAccountRev,_ := this.AccSystem.GetSubacc(benefited.Id, accountSystem.UK_AnteileAuserloesen.Id)
			bookFromTo(employeeshare,  employeeMainAccount, employeeSubbAccountRev)

		}


		// Die CC_Vertriebsprovision bekommt der oder die Dealbringer
		if benefited.Type != valueMagnets.StakeholderTypeOthers { // Don't give 5% for travel expenses and co...
	
			// get all involved Dealbringer
			// split Vertriebsprovision between all involved CostCenters
			ccArr := this.getCostCenter()
			numCc := float64(len(ccArr))
			for _,cc := range ccArr {
				// log.Println("in BookRevenueToEmployeeCostCenter:",  benefited.Id, cc)
				// Buchung Vertriebsprovision
				b := booking.Booking{
					RowNr: 		 this.Booking.RowNr,
					Amount:      this.Booking.Net[benefited] * account.PartnerProvision / numCc,
					Type:        booking.CC_Vertriebsprovision,
					Project:     this.Booking.Project,
					Text:        this.Booking.Text + "#Provision#" + cc + " of " + benefited.Id,
					Month:       this.Booking.Month,
					Year:        this.Booking.Year,
					FileCreated: this.Booking.FileCreated,
					BankCreated: this.Booking.BankCreated,
					CostCenter:  cc,
				}

				// from k-subacc-provisions --> should be from costs account in case of employees
				// --> otherwise from provisions account
				// to employees main account
				var acc_k_subprov *account.Account
				if isEmployee(cc) {
					acc_k_subprov ,_ = this.AccSystem.GetSubacc(valueMagnets.StakeholderKM.Id, accountSystem.UK_Kosten.Id)
				} else  {
					acc_k_subprov ,_ = this.AccSystem.GetSubacc(valueMagnets.StakeholderKM.Id, accountSystem.UK_Vertriebsprovision.Id)
				}
				acc_cc_main,_ := this.AccSystem.Get(cc)
				bookFromTo(b,  acc_k_subprov, acc_cc_main )

				// and from that to subacc provisions
				provisionAccount, _ := this.AccSystem.GetSubacc(cc, accountSystem.UK_Vertriebsprovision.Id)
				bookFromTo(b,  acc_cc_main, provisionAccount )

			}
		}
	}
}

// Split CostCenter String by Comma and rreturn an Array of costCenters
func  (this BookRevenueToEmployeeCostCenter) getCostCenter() []string {
	var ccArr []string
	var sh valueMagnets.Stakeholder

	// check if valid costCenter
	for _,cc := range strings.Split(this.Booking.Responsible,",") {
		if !sh.IsValidStakeholder(cc) {
			log.Printf("in BookRevenueToEmployeeCostCenter.getCostCenter(), invalid cc: '%s'\n", cc)
			log.Printf("                ==> setting '%s' to '%s'\n", cc, valueMagnets.StakeholderKM.Id)
			ccArr = append(ccArr, valueMagnets.StakeholderKM.Id)
		} else {
			ccArr = append(ccArr, cc)
		}
	}

	return ccArr
}


// Eine Buchung kann mehrere Nettopositionen enthalten, den je einem Stakeholder zugeschrieben wird.
// Diese Funktion liefert ein Array mit Stateholdern, deren CC_Nettoanteil in der Buchung != 0 ist.
func (this BookRevenueToEmployeeCostCenter) stakeholderWithNetPositions() []valueMagnets.Stakeholder {
	var result []valueMagnets.Stakeholder
	for k, v := range this.Booking.Net {
		if v != 0 {
			result = append(result, k)
		}
	}
	return result
}

// is this an Open Position?
func (this BookRevenueToEmployeeCostCenter) isOpenPosition() bool {
	emptyTime := time.Time{}
	return this.Booking.BankCreated == emptyTime
}

func (this BookRevenueToEmployeeCostCenter) sumOfNetPositions () float64{
	result := 0.0
	for _, v := range this.Booking.Net {
		if v != 0 {
			result += v
		}
	}
	return result
}

