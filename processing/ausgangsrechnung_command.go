package processing

import (
	"bitbucket.org/rwirdemann/kontrol/account"
	"bitbucket.org/rwirdemann/kontrol/booking"
	"bitbucket.org/rwirdemann/kontrol/owner"
)

type BookAusgangsrechnungCommand struct {
	Booking    booking.Booking
	Repository account.Repository
}

func (this BookAusgangsrechnungCommand) run() {
	benefitees := this.stakeholderWithNetPositions()
	for _, benefited := range benefitees {

		if benefited.Type == owner.StakeholderTypePartner {

			// book partner share
			b := booking.Booking{
				Amount: this.Booking.Net[benefited] * owner.PartnerShare,
				Type:   booking.Nettoanteil,
				Text:   this.Booking.Text + "#NetShare#" + benefited.Id,
				Month:  this.Booking.Month,
				Year:   this.Booking.Year,
				FileCreated: this.Booking.FileCreated,
				BankCreated: this.Booking.BankCreated}
			a, _ := this.Repository.Get(benefited.Id)
			a.Book(b)

			// book kommitment share
			kommitmentShare := booking.Booking{
				Amount: this.Booking.Net[benefited] * owner.KommmitmentShare,
				Type:   booking.Kommitmentanteil,
				Text:   this.Booking.Text + "#Kommitment#" + benefited.Id,
				Month:  this.Booking.Month,
				Year:   this.Booking.Year,
				FileCreated: this.Booking.FileCreated,
				BankCreated: this.Booking.BankCreated}

			kommitmentAccount, _ := this.Repository.Get(owner.StakeholderKM.Id)
			kommitmentAccount.Book(kommitmentShare)
		}

		if benefited.Type == owner.StakeholderTypeExtern {

			// book kommitment share
			kommitmentShare := booking.Booking{
				Amount: this.Booking.Net[benefited] * owner.KommmitmentExternShare,
				Type:   booking.Kommitmentanteil,
				Text:   this.Booking.Text + "#Kommitment#" + benefited.Id,
				Month:  this.Booking.Month,
				Year:   this.Booking.Year,
				FileCreated: this.Booking.FileCreated,
				BankCreated: this.Booking.BankCreated}
			kommitmentAccount, _ := this.Repository.Get(owner.StakeholderKM.Id)
			kommitmentAccount.Book(kommitmentShare)
		}

		if benefited.Type == owner.StakeholderTypeEmployee {

			// book kommitment share
			kommitmentShare := booking.Booking{
				Amount:     this.Booking.Net[benefited] * owner.KommmitmentEmployeeShare,
				Type:       booking.Kommitmentanteil,
				Text:       this.Booking.Text,
				Month:      this.Booking.Month,
				Year:       this.Booking.Year,
				FileCreated: this.Booking.FileCreated,
				BankCreated: this.Booking.BankCreated,
				CostCenter: benefited.Id}
			kommitmentAccount, _ := this.Repository.Get(owner.StakeholderKM.Id)
			kommitmentAccount.Book(kommitmentShare)
		}

		// Die Vertriebsprovision bekommt entweder ein Partner, oder wird dem K-Account gut geschrieben.
		// Letzters, wenn der Vertriebserfolg einem Angestellten zuzuordnen ist. In diesem Fall wird die
		// Kostenstelle auf die Id des Angestellten gesetzt, so dass die Gutschrift diesem zugeordnet
		// werden kann.
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
			Amount:     this.Booking.Net[benefited] * owner.PartnerProvision,
			Type:       booking.Vertriebsprovision,
			Text:       this.Booking.Text + "#Provision#" + benefited.Id,
			Month:      this.Booking.Month,
			Year:       this.Booking.Year,
			FileCreated: this.Booking.FileCreated,
			BankCreated: this.Booking.BankCreated,
			CostCenter: costcenter}
		provisionAccount.Book(b)
	}
}

// Eine Buchung kann mehrere Nettopositionen enthalten, den je einem Stakeholder zugeschrieben wird.
// Diese Funktion liefert ein Array mit Stateholder, deren Nettoanteil in der Buchung > 0 ist.
func (this BookAusgangsrechnungCommand) stakeholderWithNetPositions() []owner.Stakeholder {
	var result []owner.Stakeholder

	for k, v := range this.Booking.Net {
		if v > 0 {
			result = append(result, k)
		}
	}
	return result
}
