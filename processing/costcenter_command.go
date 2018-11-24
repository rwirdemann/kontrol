package processing

import (
	"fmt"
	"github.com/ahojsenn/kontrol/account"
	"github.com/ahojsenn/kontrol/accountSystem"
	"github.com/ahojsenn/kontrol/booking"
	"github.com/ahojsenn/kontrol/valueMagnets"
	"log"
	"time"
)


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
		log.Println("in BookToCostCenter, cc empty in row ", c.Booking.RowNr)
		log.Println("    , setting it to 'K' ")
		bkresp = valueMagnets.StakeholderKM.Id
	}
	sollAccount,_ := c.AccSystem.Get(bkresp)
	b1 := booking.CloneBooking(c.Booking, amount, bkt, c.Booking.CostCenter, c.Booking.Soll, c.Booking.Haben, c.Booking.Project)
	sollAccount.Book(b1)

	// Habenbuchung
	habenAccount,_ := c.AccSystem.Get(accountSystem.AlleKLRBuchungen.Id)
	b2 := booking.CloneBooking(c.Booking, -amount, bkt, c.Booking.CostCenter, c.Booking.Soll, c.Booking.Haben, c.Booking.Project)
	habenAccount.Book(b2)
}


type BookRevenueToEmployeeCostCenter struct {
	Booking    booking.Booking
	AccSystem  accountSystem.AccountSystem

}

func (this BookRevenueToEmployeeCostCenter) run() {

	// hier kommt nun die ganze Verteilung unter den kommitmenschen
	// 1. erst employees auszahlen und vertriebsprovisionen, den Rest auf kommitment
	// 2. später den kommitmenttopf unter den kommanditisten aufteilen


	benefitees := this.stakeholderWithNetPositions()

	for _, benefited := range benefitees {

		if benefited.Type == valueMagnets.StakeholderTypePartner {

			// book kommitment share
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

			kommitmentAccount, _ := this.AccSystem.Get(valueMagnets.StakeholderKM.Id)
			kommitmentAccount.Book(kommitmentShare)
		}

		if benefited.Type == valueMagnets.StakeholderTypeExtern {

			// book Extern share
			kommitmentShare := booking.Booking{
				RowNr:       this.Booking.RowNr,
				Amount:      this.Booking.Net[benefited] * account.KommmitmentExternShare,
				Type:        booking.CC_KommitmentanteilEX,
				CostCenter:  this.Booking.CostCenter,
				Text:        this.Booking.Text + "#Kommitment#" + benefited.Id,
				Month:       this.Booking.Month,
				Year:        this.Booking.Year,
				FileCreated: this.Booking.FileCreated,
				BankCreated: this.Booking.BankCreated}
			kommitmentAccount, _ := this.AccSystem.Get(valueMagnets.StakeholderKM.Id)
			kommitmentAccount.Book(kommitmentShare)
		}

		// Book the rest. This can happen e.g. due to
		// non person related things on the invoice like travel expenses or similar

		if benefited.Type == valueMagnets.StakeholderTypeOthers {

			// book kommitment share
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
			kommitmentAccount, _ := this.AccSystem.Get(valueMagnets.StakeholderKM.Id)
			kommitmentAccount.Book(kommitmentShare)
		}

		if benefited.Type == valueMagnets.StakeholderTypeEmployee {
			// book kommitment share
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
			kommitmentAccount, _ := this.AccSystem.Get(valueMagnets.StakeholderKM.Id)
			kommitmentAccount.Book(kommitmentShare)

			// book employee share
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
			employeeaccount, _ := this.AccSystem.Get(benefited.Id)
			employeeaccount.Book(employeeshare)

		}


		// Die CC_Vertriebsprovision bekommt der Dealbringer
		if benefited.Type != valueMagnets.StakeholderTypeOthers { // Don't give 5% for travel expenses and co...
			var provisionAccount *account.Account

			// Vertriebsprovisionen gehen nur an employees und partner, ansonsten fallen die an Kommitment
			provisionAccount, _ = this.AccSystem.Get(this.Booking.Responsible)
			if ( provisionAccount.Description.Type != valueMagnets.StakeholderTypeEmployee &&
				provisionAccount.Description.Type != valueMagnets.StakeholderTypePartner ) {
				// then provision goes to kommitment
				provisionAccount, _ = this.AccSystem.Get(valueMagnets.StakeholderKM.Id)
			}
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