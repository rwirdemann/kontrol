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


type BookFromKtoKommitmensch struct{
	Booking    booking.Booking
	AccSystem  accountSystem.AccountSystem
	SubAcc 	   string
}

func (c BookFromKtoKommitmensch) run() {
	amount := c.Booking.Amount

	// Sollbuchung
	bkresp := c.Booking.CostCenter
	if bkresp == "" {
		log.Println("in BookToCostCenter, cc empty in row ", c.Booking.RowNr)
		log.Println("    , setting it to 'K' ")
		bkresp = valueMagnets.StakeholderKM.Id
	}
	sollAccount,_ := c.AccSystem.GetSubacc(bkresp, c.SubAcc)
	b1 := booking.CloneBooking(c.Booking, amount, c.Booking.Type, c.Booking.CostCenter, c.Booking.Soll, c.Booking.Haben, c.Booking.Project)
	sollAccount.Book(b1)

	// Habenbuchung
	habenAccount,_ := c.AccSystem.GetSubacc(valueMagnets.StakeholderKM.Id, c.SubAcc)
	b2 := booking.CloneBooking(c.Booking, -amount, c.Booking.Type, c.Booking.CostCenter, c.Booking.Soll, c.Booking.Haben, c.Booking.Project)
	habenAccount.Book(b2)
}


// for a given booking
// generate a new pair of bookings,
// soll in UK_Kosten of the booking responsible Stakeholder
// haben in the valueMagnets.StakeholderKM.Id accoung (company account)
type BookCostToCostCenter struct {
	Booking    booking.Booking
	AccSystem  accountSystem.AccountSystem
}

func (c BookCostToCostCenter) run() {

	amount := c.Booking.Amount

	// set booking Type
	var bkt string = "hier steht der Buchungstyp"
	switch c.Booking.Type {
	case booking.Eingangsrechnung, booking.SKR03:
		bkt = booking.Kosten
	default:
		bkt = c.Booking.Type
	}

	// Sollbuchung
	bkresp := c.Booking.CostCenter
	if bkresp == "" {
		log.Println("in BookCostToCostCenter, cc empty in row ", c.Booking.RowNr,  c.Booking)
		log.Println("    , setting it to 'K' ")
		bkresp = valueMagnets.StakeholderKM.Id
	}
	sollAccount,_ := c.AccSystem.GetSubacc(bkresp, accountSystem.UK_Kosten)
	b1 := booking.CloneBooking(c.Booking, amount, bkt, c.Booking.CostCenter, c.Booking.Soll, c.Booking.Haben, c.Booking.Project)
	sollAccount.Book(b1)

	// Habenbuchung
	habenAccount,_ := c.AccSystem.Get(valueMagnets.StakeholderKM.Id)
	b2 := booking.CloneBooking(c.Booking, -amount, bkt, c.Booking.CostCenter, c.Booking.Soll, c.Booking.Haben, c.Booking.Project)
	b2.Text = b2.Text + "-Gegenbuchung"
	habenAccount.Book(b2)
}


type BookRevenueToEmployeeCostCenter struct {
	Booking    booking.Booking
	AccSystem  accountSystem.AccountSystem

}

func (this BookRevenueToEmployeeCostCenter) run() {

	// hier kommt nun die ganze Verteilung unter den kommitmenschen
	// 1. erst employees auszahlen und vertriebsprovisionen, den Rest auf kommitment
	// 2. später den kommitment topf unter den kommanditisten aufteilen


	benefitees := this.stakeholderWithNetPositions()

	for _, benefited := range benefitees {

		// book Partners revenue to kommitment for now
		if benefited.Type == valueMagnets.StakeholderTypePartner {
			kommitmentShare := booking.Booking{
				RowNr:       this.Booking.RowNr,
				Amount:      this.Booking.Net[benefited] * (1.00 - account.PartnerProvision),
				Type:        booking.CC_PartnerNettoFaktura,
				CostCenter:  benefited.Id,
				Text:        this.Booking.Text + "#Kommitment#" + benefited.Id,
				Month:       this.Booking.Month,
				Year:        this.Booking.Year,
				FileCreated: this.Booking.FileCreated,
				BankCreated: this.Booking.BankCreated}

			kommitmentAccount,_ := this.AccSystem.GetSubacc(valueMagnets.StakeholderKM.Id, accountSystem.UK_AnteileAuserloesen)
			kommitmentAccount.Book(kommitmentShare)

			// Gegenbuchung
			sollAccount,_ := this.AccSystem.Get(valueMagnets.StakeholderKM.Id)
			kommitmentShare.Amount *= -1
			sollAccount.Book (kommitmentShare)
		}

		// book Externals revenue to kommitment
		if benefited.Type == valueMagnets.StakeholderTypeExtern {
			bkt := booking.CC_KommitmentanteilEX
			kommitmentShare := booking.Booking{
				RowNr:       this.Booking.RowNr,
				Amount:      this.Booking.Net[benefited] * account.KommmitmentExternShare,
				Type:        bkt,
				CostCenter:  this.Booking.CostCenter,
				Text:        this.Booking.Text + "#Kommitment#" + benefited.Id,
				Month:       this.Booking.Month,
				Year:        this.Booking.Year,
				FileCreated: this.Booking.FileCreated,
				BankCreated: this.Booking.BankCreated}

			kommitmentAccount,_ := this.AccSystem.GetSubacc(valueMagnets.StakeholderKM.Id, accountSystem.UK_AnteileAuserloesen)
			kommitmentAccount.Book(kommitmentShare)

			// Gegenbuchung
			sollAccount,_ := this.AccSystem.Get(valueMagnets.StakeholderKM.Id)
			kommitmentShare.Amount *= -1
			sollAccount.Book (kommitmentShare)
		}

		// Book the rest (position) to kommitment. This can happen e.g. due to
		// non person related things on the invoice like travel expenses or similar
		if benefited.Type == valueMagnets.StakeholderTypeOthers {
			kommitmentShare := booking.Booking{
				RowNr: 		 this.Booking.RowNr,
				Amount:      this.Booking.Net[benefited] * account.KommmitmentOthersShare,
				Type:        booking.CC_Kommitmentanteil,
				CostCenter:  this.Booking.CostCenter,
				Text:        this.Booking.Text + "#Kommitment#Rest#" + benefited.Id,
				Month:       this.Booking.Month,
				Year:        this.Booking.Year,
				FileCreated: this.Booking.FileCreated,
				BankCreated: this.Booking.BankCreated}

			kommitmentAccount,_ := this.AccSystem.GetSubacc(valueMagnets.StakeholderKM.Id, accountSystem.UK_AnteileAuserloesen)
			kommitmentAccount.Book(kommitmentShare)

			// Gegenbuchung
			sollAccount,_ := this.AccSystem.Get(valueMagnets.StakeholderKM.Id)
			kommitmentShare.Amount *= -1
			sollAccount.Book (kommitmentShare)
		}

		// now book kommitment part of employee's revenue
		if benefited.Type == valueMagnets.StakeholderTypeEmployee {
			kommitmentShare := booking.Booking{
				RowNr: 		 this.Booking.RowNr,
				Amount:      this.Booking.Net[benefited] * (account.KommmitmentEmployeeShare-account.EmployeeShare),
				Type:        booking.CC_Kommitmentanteil,
				Text:        this.Booking.Text,
				Month:       this.Booking.Month,
				Year:        this.Booking.Year,
				FileCreated: this.Booking.FileCreated,
				BankCreated: this.Booking.BankCreated,
				CostCenter:  this.Booking.CostCenter}


			kommitmentAccount,_ := this.AccSystem.GetSubacc(valueMagnets.StakeholderKM.Id, accountSystem.UK_AnteileAuserloesen)
			kommitmentAccount.Book(kommitmentShare)

			// Gegenbuchung
			sollAccount,_ := this.AccSystem.Get(valueMagnets.StakeholderKM.Id)
			kommitmentShare.Amount *= -1
			sollAccount.Book (kommitmentShare)


			// book employee share of employees revenue
			employeeshare := booking.Booking{
				RowNr: 		 this.Booking.RowNr,
				Amount:      this.Booking.Net[benefited] * account.EmployeeShare,
				Type:        booking.CC_Employeeaanteil,
				Text:        fmt.Sprintf("%f", account.EmployeeShare)+"*Netto: " + this.Booking.Text,
				Month:       this.Booking.Month,
				Year:        this.Booking.Year,
				FileCreated: this.Booking.FileCreated,
				BankCreated: this.Booking.BankCreated,
				CostCenter:  this.Booking.CostCenter}

			employeeaccount,_ := this.AccSystem.GetSubacc(benefited.Id, accountSystem.UK_AnteileAuserloesen)
			employeeaccount.Book(employeeshare)

			// Gegenbuchung, das sind Kosten für Kommitment
			sollAccount,_ = this.AccSystem.Get(valueMagnets.StakeholderKM.Id)
			employeeshare.Amount *= -1
			sollAccount.Book (employeeshare)
		}


		// Die CC_Vertriebsprovision bekommt der Dealbringer
		if benefited.Type != valueMagnets.StakeholderTypeOthers { // Don't give 5% for travel expenses and co...
			var provisionAccount *account.Account
			var sollAccount *account.Account
			/*
			// nur wenn auf der Buchung der Stakeholder ein employee ist...
			// log.Println("in BookRevenueToEmployeeCostCenter.run ==>", this.Booking.CostCenter,
			//	valueMagnets.StakeholderRepository{}.TypeOf(this.Booking.CostCenter))
			if isEmployee(this.Booking.CostCenter) {
				//&& provisionAccount.Description.Type != valueMagnets.StakeholderTypePartner
				// then provision goes to kommitment
				provisionAccount, _ = this.AccSystem.GetSubacc(this.Booking.Responsible, accountSystem.UK_Vertriebsprovision)
				// sollAccount, _ = this.AccSystem.Get(this.Booking.Responsible)
				} else {
				provisionAccount, _ = this.AccSystem.GetSubacc(valueMagnets.StakeholderKM.Id, accountSystem.UK_Vertriebsprovision)
				// sollAccount,_ = this.AccSystem.Get(valueMagnets.StakeholderKM.Id)
				}
			// das sind Kosten für Kommitment
			*/

			provisionAccount, _ = this.AccSystem.GetSubacc(this.Booking.Responsible, accountSystem.UK_Vertriebsprovision)
			b := booking.Booking{
				RowNr: 		 this.Booking.RowNr,
				Amount:      this.Booking.Net[benefited] * account.PartnerProvision,
				Type:        booking.CC_Vertriebsprovision,
				Text:        this.Booking.Text + "#Provision#" + benefited.Id,
				Month:       this.Booking.Month,
				Year:        this.Booking.Year,
				FileCreated: this.Booking.FileCreated,
				BankCreated: this.Booking.BankCreated,
				CostCenter:  this.Booking.CostCenter}
			provisionAccount.Book(b)

			// Gegenbuchung
			sollAccount,_ = this.AccSystem.Get(valueMagnets.StakeholderKM.Id)
			b.Amount *= -1
			sollAccount.Book (b)
		}
	}
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

func isEmployee  (id string) bool {
	shrepo := valueMagnets.Stakeholder{}
	if  strings.Compare (shrepo.TypeOf(id), valueMagnets.StakeholderTypeEmployee ) ==  0 {
		return true
	}
	return false
}