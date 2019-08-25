package processing

import (
	"fmt"
	"github.com/ahojsenn/kontrol/accountSystem"
	"github.com/ahojsenn/kontrol/booking"
	"github.com/ahojsenn/kontrol/valueMagnets"
	"log"
	"strconv"
	"strings"
)

// for a given booking
// generate a new pair of bookings,
// soll in UK_Kosten of the booking responsible Stakeholder
// haben in the valueMagnets.StakeholderKM.Id accoung (company account)
type BookCostToCostCenter struct {
	Booking    booking.Booking
	AccSystem  accountSystem.AccountSystem
}

func (this BookCostToCostCenter) run() {

	// Sollbuchung
	bkresp := this.Booking.CostCenter
	if bkresp == "" {
		log.Println("in BookCostToCostCenter, cc empty in row ", this.Booking.RowNr)
		log.Println("    , setting it to 'K' ")
		bkresp = valueMagnets.StakeholderKM.Id
	}

	this.Booking.Type = booking.Kosten

	// book from kommitment company account to k subacc costs
	acc_k_main,_ := this.AccSystem.Get(valueMagnets.StakeholderKM.Id)
	acc_k_subcosts,_ := this.AccSystem.GetSubacc(valueMagnets.StakeholderKM.Id, accountSystem.UK_Kosten.Id)
	bookFromTo(this.Booking, acc_k_main, acc_k_subcosts)

	// if costcenter is != Company (k) then distribute costs
	// 1. from acc_k_subcosts to employees main account
	// 2. and from there to employees subacc.
	if (bkresp != valueMagnets.StakeholderKM.Id) {
		this.Booking.Type = booking.Kosten+"_aufKommitmensch"
		habenAccount,_ := this.AccSystem.Get(bkresp)
		//log.Println("in BookCostToCostCenter, rowNr:", this.Booking.RowNr)
		bookFromTo(this.Booking, acc_k_subcosts, habenAccount)

		sollAccount,_ := this.AccSystem.GetSubacc(bkresp, accountSystem.UK_Kosten.Id)
		bookFromTo(this.Booking, habenAccount, sollAccount)
	}
}



type BookToValuemagnetsByShares struct{
	Booking    booking.Booking
	AccSystem  accountSystem.AccountSystem
	SubAcc 	   string
}

func (c BookToValuemagnetsByShares) run() {

	// book from k-main account to k-subacc type
	// var k_subAcc *account.Account
	//k_mainAcc,_ := c.AccSystem.Get(valueMagnets.StakeholderKM.Id)
	//k_subAcc,_ := c.AccSystem.GetSubacc(valueMagnets.StakeholderKM.Id, c.SubAcc)
	//bookFromTo(c.Booking, k_mainAcc, k_subAcc)

	// loop through all Stakeholders
	shrepo := valueMagnets.Stakeholder{}
	for _,sh := range shrepo.GetAllOfType (valueMagnets.StakeholderTypePartner) {
		k_subAcc,_ := c.AccSystem.GetSubacc(valueMagnets.StakeholderKM.Id, c.SubAcc)
		fairshares,_  := strconv.ParseFloat(sh.Fairshares, 64)
		amount := c.Booking.Amount*fairshares

		// book from k-subacc to stakeholders main acc
		sh_mainAcc,_ := c.AccSystem.Get(sh.Id)
		b1 := booking.CloneBooking(c.Booking, amount, c.Booking.Type, c.Booking.CostCenter, c.Booking.Soll, c.Booking.Haben, c.Booking.Project)
		b1.Text += " "+sh.Id+" Anteil von "+fmt.Sprintf("%.2f",c.Booking.Amount)+"€"
		bookFromTo(b1, k_subAcc, sh_mainAcc)

		// and from stakeholders mainaccount to subacc
		sh_subAcc,_ := c.AccSystem.GetSubacc(sh.Id, c.SubAcc)
		bookFromTo(b1, sh_mainAcc,sh_subAcc)
	}
}







func isEmployee  (id string) bool {
	shrepo := valueMagnets.Stakeholder{}
	if  strings.Compare (shrepo.TypeOf(id), valueMagnets.StakeholderTypeEmployee ) ==  0 {
		return true
	}
	return false
}