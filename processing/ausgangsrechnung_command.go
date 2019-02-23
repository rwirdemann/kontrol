package processing

import (
	"time"

	"github.com/ahojsenn/kontrol/accountSystem"
	"github.com/ahojsenn/kontrol/booking"
	"github.com/ahojsenn/kontrol/util"
)


type BookAusgangsrechnungCommand struct {
	Booking    booking.Booking
	AccSystem  accountSystem.AccountSystem
}

func (this BookAusgangsrechnungCommand) run() {

	// if booking with empty timestamp in position "BankCreated"
	// then book it to open positions SKR03_1400
	//
	// the same
	// if booking is beyond the current financial year / balance date, then book to SKR 1400
	// "Forderungen aus Lieferung und Leistung"
	// and from that to bank... || this.isBeyondBudgetDate()

	sollAccId := ""
	if this.Booking.IsOpenPosition() || this.Booking.IsBeyondBudgetDate() {
		sollAccId = accountSystem.SKR03_1400.Id
		this.Booking.Text = "Achtung OPOS "+this.Booking.Text
	} else {
		// book from bankaccount...
		sollAccId = accountSystem.SKR03_1200.Id
		}

	sollkto,_ := this.AccSystem.Get(sollAccId)
	a := booking.Booking{
		RowNr: 		 this.Booking.RowNr,
		Amount:      -this.Booking.Amount,
		Project:     this.Booking.Project,
		Type:        booking.Erloese,
		Text:        this.Booking.Text,
		Month:       this.Booking.Month,
		Year:        this.Booking.Year,
		FileCreated: this.Booking.FileCreated,
		BankCreated: this.Booking.BankCreated,
		CsvBookingExtras: 		 this.Booking.CsvBookingExtras,
	}
	sollkto.Book(a)



	// haben umsatzerloese
	umsatzerloese, _ := this.AccSystem.Get(accountSystem.SKR03_Umsatzerloese.Id)
	b := booking.Booking{
		RowNr: 		 this.Booking.RowNr,
		CostCenter:  this.Booking.CostCenter,
		Project:     this.Booking.Project,
		Amount:      util.Net(this.Booking.Amount),
		Type:        booking.Erloese,
		Text:        this.Booking.Text,
		Month:       this.Booking.Month,
		Year:        this.Booking.Year,
		FileCreated: this.Booking.FileCreated,
		BankCreated: this.Booking.BankCreated,
		CsvBookingExtras: 		 this.Booking.CsvBookingExtras,
	}
	umsatzerloese.Book(b)

	// haben Steuern
	umsatzsteuernKonto ,_ := this.AccSystem.Get(accountSystem.SKR03_Umsatzsteuer.Id)
	c :=
		booking.Booking{
			RowNr: 		 this.Booking.RowNr,
			Amount:      this.Booking.Amount - util.Net(this.Booking.Amount),
			CostCenter:  this.Booking.CostCenter,
			Project:     this.Booking.Project,
			Type:        booking.Erloese,
			Text:        this.Booking.Text,
			Month:       this.Booking.Month,
			Year:        this.Booking.Year,
			FileCreated: this.Booking.FileCreated,
			BankCreated: this.Booking.BankCreated,
			CsvBookingExtras: 		 this.Booking.CsvBookingExtras,
		}
	umsatzsteuernKonto.Book(c)

}

// is this an Open Position?
func (this BookAusgangsrechnungCommand) isOpenPosition() bool {
	emptyTime := time.Time{}
	return this.Booking.BankCreated == emptyTime
}

// ist this booking beyond the budget date?
func (this BookAusgangsrechnungCommand) isBeyondBudgetDate () bool {
	return this.Booking.BankCreated.After(util.Global.BalanceDate)
}


