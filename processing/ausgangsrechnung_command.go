package processing

import (
	"time"

	"github.com/ahojsenn/kontrol/account"
	"github.com/ahojsenn/kontrol/booking"
	"github.com/ahojsenn/kontrol/owner"
	"github.com/ahojsenn/kontrol/util"
	"github.com/ahojsenn/kontrol/accountSystem"
)

type BookAusgangsrechnungCommand struct {
	Booking    booking.Booking
	Repository accountSystem.AccountSystem
}

func (this BookAusgangsrechnungCommand) run() {


	// book from bankaccount...
	bank := this.Repository.BankAccount()
	a := booking.Booking{
		Amount:      util.Net(this.Booking.Amount),
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
		skr1400, _ := this.Repository.Get(accountSystem.SKR03_1400.Id)
		skr1400.Book(this.Booking)
		return
	} else {
		// haben umsatzerloese
		umsatzerloese, _ := this.Repository.Get(accountSystem.SKR03_Umsatzerloese.Id)
		b := booking.Booking{
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
		umsatzsteuernKonto ,_ := this.Repository.Get(accountSystem.SKR03_Umsatzsteuer.Id)
		c :=
			booking.Booking{
				Amount:      this.Booking.Amount - util.Net(this.Booking.Amount),
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

		if benefited.Type == owner.StakeholderTypePartner {

			// book partner share
			b := booking.Booking{
				Amount:      this.Booking.Net[benefited] * owner.PartnerShare,
				Type:        booking.Nettoanteil,
				Text:        this.Booking.Text + "#NetShare#" + benefited.Id,
				Month:       this.Booking.Month,
				Year:        this.Booking.Year,
				FileCreated: this.Booking.FileCreated,
				BankCreated: this.Booking.BankCreated}
			a, _ := this.Repository.Get(benefited.Id)
			a.Book(b)

			// book kommitment share
			kommitmentShare := booking.Booking{
				Amount:      this.Booking.Net[benefited] * owner.KommmitmentShare,
				Type:        booking.Kommitmentanteil,
				Text:        this.Booking.Text + "#Kommitment#" + benefited.Id,
				Month:       this.Booking.Month,
				Year:        this.Booking.Year,
				FileCreated: this.Booking.FileCreated,
				BankCreated: this.Booking.BankCreated}

			kommitmentAccount, _ := this.Repository.Get(owner.StakeholderKM.Id)
			kommitmentAccount.Book(kommitmentShare)
		}

		if benefited.Type == owner.StakeholderTypeExtern {

			// book kommitment share
			kommitmentShare := booking.Booking{
				Amount:      this.Booking.Net[benefited] * owner.KommmitmentExternShare,
				Type:        booking.Kommitmentanteil,
				Text:        this.Booking.Text + "#Kommitment#" + benefited.Id,
				Month:       this.Booking.Month,
				Year:        this.Booking.Year,
				FileCreated: this.Booking.FileCreated,
				BankCreated: this.Booking.BankCreated}
			kommitmentAccount, _ := this.Repository.Get(owner.StakeholderKM.Id)
			kommitmentAccount.Book(kommitmentShare)
		}

		// Book the rest. This can happen e.g. due to
		// non person related things on the invoice like travel expenses or similar

		if benefited.Type == owner.StakeholderTypeOthers {

			// book kommitment share
			kommitmentShare := booking.Booking{
				Amount:      this.Booking.Net[benefited] * owner.KommmitmentOthersShare,
				Type:        booking.Kommitmentanteil,
				Text:        this.Booking.Text + "#Kommitment#Rest#" + benefited.Id,
				Month:       this.Booking.Month,
				Year:        this.Booking.Year,
				FileCreated: this.Booking.FileCreated,
				BankCreated: this.Booking.BankCreated}
			kommitmentAccount, _ := this.Repository.Get(owner.StakeholderKM.Id)
			kommitmentAccount.Book(kommitmentShare)
		}

		if benefited.Type == owner.StakeholderTypeEmployee {

			// book kommitment share
			kommitmentShare := booking.Booking{
				Amount:      this.Booking.Net[benefited] * owner.KommmitmentEmployeeShare,
				Type:        booking.Kommitmentanteil,
				Text:        this.Booking.Text,
				Month:       this.Booking.Month,
				Year:        this.Booking.Year,
				FileCreated: this.Booking.FileCreated,
				BankCreated: this.Booking.BankCreated,
				CostCenter:  benefited.Id}
			kommitmentAccount, _ := this.Repository.Get(owner.StakeholderKM.Id)
			kommitmentAccount.Book(kommitmentShare)
		}

		// Die Vertriebsprovision bekommt entweder ein Partner, oder wird dem K-Account gut geschrieben.
		// Letzters, wenn der Vertriebserfolg einem Angestellten zuzuordnen ist. In diesem Fall wird die
		// Kostenstelle auf die Id des Angestellten gesetzt, so dass die Gutschrift diesem zugeordnet
		// werden kann.
		if benefited.Type != owner.StakeholderTypeOthers { // DON'T givr 5% for travel expenses and co...
			var provisionAccount *account.Account
			var costcenter string
			stakeholderRepository := owner.StakeholderRepository{}
			if stakeholderRepository.TypeOf(this.Booking.Responsible) == owner.StakeholderTypeEmployee {
				provisionAccount, _ = this.Repository.Get(owner.StakeholderKM.Id)
				costcenter = this.Booking.Responsible
			} else {
				provisionAccount, _ = this.Repository.Get(this.Booking.Responsible)
			}
			b := booking.Booking{
				Amount:      this.Booking.Net[benefited] * owner.PartnerProvision,
				Type:        booking.Vertriebsprovision,
				Text:        this.Booking.Text + "#Provision#" + benefited.Id,
				Month:       this.Booking.Month,
				Year:        this.Booking.Year,
				FileCreated: this.Booking.FileCreated,
				BankCreated: this.Booking.BankCreated,
				CostCenter:  costcenter}
			provisionAccount.Book(b)
		}
	}
}

// Eine Buchung kann mehrere Nettopositionen enthalten, den je einem Stakeholder zugeschrieben wird.
// Diese Funktion liefert ein Array mit Stateholder, deren Nettoanteil in der Buchung != 0 ist.
func (this BookAusgangsrechnungCommand) stakeholderWithNetPositions() []owner.Stakeholder {
	var result []owner.Stakeholder

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
