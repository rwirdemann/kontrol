package processing

import (
	"time"

	"github.com/ahojsenn/kontrol/account"
	"github.com/ahojsenn/kontrol/booking"
	"github.com/ahojsenn/kontrol/profitCenter"
	"github.com/ahojsenn/kontrol/util"
	"github.com/ahojsenn/kontrol/accountSystem"
	"math"
)

const (
	PartnerShare             = 0.70
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

	// if booking with empty timestamp in position "BankCreated"
	// the book it to open positions SKR03_1400
	if this.isOpenPosition() == true {
		skr1400, _ := this.AccSystem.Get(accountSystem.SKR03_1400.Id)
		skr1400.Book(this.Booking)
		return
	} else {
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
	}

	// hier kommt nun die ganze Verteilung unter den kommitmenschen

	benefitees := this.stakeholderWithNetPositions()
	for _, benefited := range benefitees {

		if benefited.Type == profitCenter.StakeholderTypePartner {

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
				CostCenter:  profitCenter.StakeholderKM.Id,
				Text:        this.Booking.Text + "#Kommitment#" + benefited.Id,
				Month:       this.Booking.Month,
				Year:        this.Booking.Year,
				FileCreated: this.Booking.FileCreated,
				BankCreated: this.Booking.BankCreated}

			kommitmentAccount, _ := this.AccSystem.Get(profitCenter.StakeholderKM.Id)
			kommitmentAccount.Book(kommitmentShare)
		}

		if benefited.Type == profitCenter.StakeholderTypeExtern {

			// book kommitment share
			kommitmentShare := booking.Booking{
				RowNr:       this.Booking.RowNr,
				Amount:      this.Booking.Net[benefited] * KommmitmentExternShare,
				Type:        booking.Kommitmentanteil,
				CostCenter:  profitCenter.StakeholderKM.Id,
				Text:        this.Booking.Text + "#Kommitment#" + benefited.Id,
				Month:       this.Booking.Month,
				Year:        this.Booking.Year,
				FileCreated: this.Booking.FileCreated,
				BankCreated: this.Booking.BankCreated}
			kommitmentAccount, _ := this.AccSystem.Get(profitCenter.StakeholderKM.Id)
			kommitmentAccount.Book(kommitmentShare)
		}

		// Book the rest. This can happen e.g. due to
		// non person related things on the invoice like travel expenses or similar

		if benefited.Type == profitCenter.StakeholderTypeOthers {

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
			kommitmentAccount, _ := this.AccSystem.Get(profitCenter.StakeholderKM.Id)
			kommitmentAccount.Book(kommitmentShare)
		}

		if benefited.Type == profitCenter.StakeholderTypeEmployee {
			// book kommitment share
			kommitmentShare := booking.Booking{
				RowNr: 		 this.Booking.RowNr,
				Amount:      this.Booking.Net[benefited] * KommmitmentEmployeeShare,
				Type:        booking.Kommitmentanteil,
				Text:        this.Booking.Text,
				Month:       this.Booking.Month,
				Year:        this.Booking.Year,
				FileCreated: this.Booking.FileCreated,
				BankCreated: this.Booking.BankCreated,
				CostCenter:  benefited.Id}
			kommitmentAccount, _ := this.AccSystem.Get(profitCenter.StakeholderKM.Id)
			kommitmentAccount.Book(kommitmentShare)
		}

		// Die Vertriebsprovision bekommt der Dealbringer
		if benefited.Type != profitCenter.StakeholderTypeOthers { // Don't give 5% for travel expenses and co...
			var provisionAccount *account.Account

			// Vertriebsprovisionen gehen nur an employees und partner, ansonsten fallen die an Kommitment
			provisionAccount, _ = this.AccSystem.Get(this.Booking.Responsible)
			if ( provisionAccount.Description.Type != profitCenter.StakeholderTypeEmployee &&
				 provisionAccount.Description.Type != profitCenter.StakeholderTypePartner ) {
				 	// then provision goes to kommitment
				provisionAccount, _ = this.AccSystem.Get(profitCenter.StakeholderKM.Id)
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
func (this BookAusgangsrechnungCommand) stakeholderWithNetPositions() []profitCenter.Stakeholder {
	var result []profitCenter.Stakeholder
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


