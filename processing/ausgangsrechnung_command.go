package processing

import (
	"fmt"
	"log"
	"time"

	"github.com/ahojsenn/kontrol/account"
	"github.com/ahojsenn/kontrol/booking"
	"github.com/ahojsenn/kontrol/valueMagnets"
	"github.com/ahojsenn/kontrol/util"
	"github.com/ahojsenn/kontrol/accountSystem"
	"math"
)

const (
	PartnerShare             = 0.70
	EmployeeShare			 = 0.70
	KommmitmentShare         = 0.25
	KommmitmentExternShare   = 0.95
	KommmitmentOthersShare   = 1.00
	KommmitmentEmployeeShare = 0.95
	PartnerProvision         = 0.05
)

type BookAusgangsrechnungCommand struct {
	Booking    booking.Booking
	AccSystem  accountSystem.AccountSystem
}

func (this BookAusgangsrechnungCommand) run() {

	log.Println("BookAusgangsrechnungCommand")


	// if booking with empty timestamp in position "BankCreated"
	// the book it to open positions SKR03_1400
	if this.isOpenPosition() == true {
		skr1400, _ := this.AccSystem.Get(accountSystem.SKR03_1400.Id)
		skr1400.Book(this.Booking)

		this.Booking.Text = "Achtung OPOS "+this.Booking.Text
	} else {
		// book from bankaccount...
		bank,_ := this.AccSystem.Get(accountSystem.SKR03_1200.Id)
		a := booking.Booking{
			RowNr: 		 this.Booking.RowNr,
			Amount:      -this.Booking.Amount,
			Type:        booking.Erloese,
			Text:        this.Booking.Text,
			Month:       this.Booking.Month,
			Year:        this.Booking.Year,
			FileCreated: this.Booking.FileCreated,
			BankCreated: this.Booking.BankCreated,
		}
		bank.Book(a)
	}



	// haben umsatzerloese
	umsatzerloese, _ := this.AccSystem.Get(accountSystem.SKR03_Umsatzerloese.Id)
	b := booking.Booking{
		RowNr: 		 this.Booking.RowNr,
		CostCenter:  this.Booking.CostCenter,
		Amount:      util.Net(this.Booking.Amount),
		Type:        booking.Erloese,
		Text:        this.Booking.Text,
		Month:       this.Booking.Month,
		Year:        this.Booking.Year,
		FileCreated: this.Booking.FileCreated,
		BankCreated: this.Booking.BankCreated,
	}
	umsatzerloese.Book(b)

	// haben Steuern
	umsatzsteuernKonto ,_ := this.AccSystem.Get(accountSystem.SKR03_Umsatzsteuer.Id)
	c :=
		booking.Booking{
			RowNr: 		 this.Booking.RowNr,
			Amount:      this.Booking.Amount - util.Net(this.Booking.Amount),
			CostCenter:  this.Booking.CostCenter,
			Type:        booking.Erloese,
			Text:        this.Booking.Text,
			Month:       this.Booking.Month,
			Year:        this.Booking.Year,
			FileCreated: this.Booking.FileCreated,
			BankCreated: this.Booking.BankCreated,
		}
	umsatzsteuernKonto.Book(c)


	// hier kommt nun die ganze Verteilung unter den kommitmenschen
	// 1. get rid of external share
	// 2. 5% Vertriebsprovision (Achtung die ist zu hoch, wenn der db bei externen nicht 20% erreicht)
	// 3. 70% Employeeshare, falls Leistungserbringer employee ist
	// 4. 15% mal Rendite auf shareholder, bämm geht nicht bei rein externen...
	// DAS ISTUNFERTIG DENK DENK DENK...

	benefitees := this.stakeholderWithNetPositions()
	for _, benefited := range benefitees {


		if benefited.Type == valueMagnets.StakeholderTypePartner {
			// book partner share
			b := booking.Booking{
				RowNr: 		 this.Booking.RowNr,
				Amount:      math.Round(this.Booking.Net[benefited] * PartnerShare*1000000)/1000000,
				Type:        booking.Nettoanteil,
				CostCenter:  benefited.Id,
				Text:        this.Booking.Text + "#NetShare#" + benefited.Id,
				Month:       this.Booking.Month,
				Year:        this.Booking.Year,
				FileCreated: this.Booking.FileCreated,
				BankCreated: this.Booking.BankCreated}
			a, _ := this.AccSystem.Get(benefited.Id)
			a.Book(b)

			// book kommitment share
			kommitmentShare := booking.Booking{
				RowNr:       this.Booking.RowNr,
				Amount:      this.Booking.Net[benefited] * KommmitmentShare,
				Type:        booking.Kommitmentanteil,
				CostCenter:  valueMagnets.StakeholderKM.Id,
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
				Amount:      this.Booking.Net[benefited] * KommmitmentExternShare,
				Type:        booking.Kommitmentanteil,
				CostCenter:  valueMagnets.StakeholderKM.Id,
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
				Amount:      this.Booking.Net[benefited] * KommmitmentOthersShare,
				Type:        booking.Kommitmentanteil,
				CostCenter:  benefited.Id,
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
				Amount:      this.Booking.Net[benefited] * (KommmitmentEmployeeShare-EmployeeShare),
				Type:        booking.Kommitmentanteil,
				Text:        this.Booking.Text,
				Month:       this.Booking.Month,
				Year:        this.Booking.Year,
				FileCreated: this.Booking.FileCreated,
				BankCreated: this.Booking.BankCreated,
				CostCenter:  benefited.Id}
			kommitmentAccount, _ := this.AccSystem.Get(valueMagnets.StakeholderKM.Id)
			kommitmentAccount.Book(kommitmentShare)

			// book employee share
			employeeshare := booking.Booking{
				RowNr: 		 this.Booking.RowNr,
				Amount:      this.Booking.Net[benefited] * EmployeeShare,
				Type:        booking.Employeeaanteil,
				Text:        fmt.Sprintf("%f", EmployeeShare)+"*Netto: " + this.Booking.Text,
				Month:       this.Booking.Month,
				Year:        this.Booking.Year,
				FileCreated: this.Booking.FileCreated,
				BankCreated: this.Booking.BankCreated,
				CostCenter:  benefited.Id}
			employeeaccount, _ := this.AccSystem.Get(benefited.Id)
			employeeaccount.Book(employeeshare)

		}

		// Die Vertriebsprovision bekommt der Dealbringer
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
				Amount:      this.Booking.Net[benefited] * PartnerProvision,
				Type:        booking.Vertriebsprovision,
				Text:        this.Booking.Text + "#Provision#" + benefited.Id,
				Month:       this.Booking.Month,
				Year:        this.Booking.Year,
				FileCreated: this.Booking.FileCreated,
				BankCreated: this.Booking.BankCreated,
				CostCenter:  this.Booking.Responsible}

			provisionAccount.Book(b)
		}
	}
}

// Eine Buchung kann mehrere Nettopositionen enthalten, den je einem Stakeholder zugeschrieben wird.
// Diese Funktion liefert ein Array mit Stateholdern, deren Nettoanteil in der Buchung != 0 ist.
func (this BookAusgangsrechnungCommand) stakeholderWithNetPositions() []valueMagnets.Stakeholder {
	var result []valueMagnets.Stakeholder
	for k, v := range this.Booking.Net {
		if v != 0 {
			result = append(result, k)
		}
	}
	return result
}

// is this an Open Position?
func (this BookAusgangsrechnungCommand) isOpenPosition() bool {
	emptyTime := time.Time{}
	return this.Booking.BankCreated == emptyTime
}


