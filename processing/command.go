package processing

import (
	"bitbucket.org/rwirdemann/kontrol/account"
	"bitbucket.org/rwirdemann/kontrol/booking"
	"bitbucket.org/rwirdemann/kontrol/owner"
)

type Command interface {
	run()
}

type BookGehaltCommand struct {
	Booking    booking.Booking
	Repository account.Repository
}

func (c BookGehaltCommand) run() {

	// Bankbuchung
	bankBooking := c.Booking
	bankBooking.Type = c.Booking.Typ
	bankBooking.Amount = bankBooking.Amount * -1
	c.Repository.BankAccount().Book(bankBooking)

	// Buchung Kommitment-Konto
	kBooking := booking.CloneBooking(c.Booking, -1, booking.Gehalt, c.Booking.Responsible)
	kommitmentAccount, _ := c.Repository.Get(owner.StakeholderKM.Id)
	kommitmentAccount.Book(kBooking)
}

type BookSVBeitragCommand struct {
	Booking    booking.Booking
	Repository account.Repository
}

func (c BookSVBeitragCommand) run() {

	// Bankbuchung
	bankBooking := c.Booking
	bankBooking.Type = c.Booking.Typ
	bankBooking.Amount = bankBooking.Amount * -1
	c.Repository.BankAccount().Book(bankBooking)

	// Buchung Kommitment-Konto
	kBooking := booking.CloneBooking(c.Booking, -1, booking.SVBeitrag, c.Booking.Responsible)
	kommitmentAccount, _ := c.Repository.Get(owner.StakeholderKM.Id)
	kommitmentAccount.Book(kBooking)
}

type BookLNSteuerCommand struct {
	Booking    booking.Booking
	Repository account.Repository
}

func (c BookLNSteuerCommand) run() {

	// Bankbuchung
	bankBooking := c.Booking
	bankBooking.Type = c.Booking.Typ
	bankBooking.Amount = bankBooking.Amount * -1
	c.Repository.BankAccount().Book(bankBooking)

	// Buchung Kommitment-Konto
	kBooking := booking.CloneBooking(c.Booking, -1, booking.LNSteuer, c.Booking.Responsible)
	kommitmentAccount, _ := c.Repository.Get(owner.StakeholderKM.Id)
	kommitmentAccount.Book(kBooking)
}

type BookGWSteuerCommand struct {
	Booking    booking.Booking
	Repository account.Repository
}

func (c BookGWSteuerCommand) run() {

	// Bankbuchung
	bankBooking := c.Booking
	bankBooking.Type = c.Booking.Typ
	bankBooking.Amount = bankBooking.Amount * -1
	c.Repository.BankAccount().Book(bankBooking)

	// Buchung Kommitment-Konto
	kBooking := booking.CloneBooking(c.Booking, -1, booking.GWSteuer, c.Booking.Responsible)
	kommitmentAccount, _ := c.Repository.Get(owner.StakeholderKM.Id)
	kommitmentAccount.Book(kBooking)
}

type BookAusgangsrechnungCommand struct {
	Booking    booking.Booking
	Repository account.Repository
}

func (this BookAusgangsrechnungCommand) run() {
	benefitees := stakeholderWithNetPositions(this.Booking)
	for _, benefited := range benefitees {

		if benefited.Type == owner.StakeholderTypePartner {

			// book partner share
			b := booking.Booking{
				Amount: this.Booking.Net[benefited] * owner.PartnerShare,
				Type:   booking.Nettoanteil,
				Text:   this.Booking.Text + "#NetShare#" + benefited.Id,
				Month:  this.Booking.Month,
				Year:   this.Booking.Year}
			a, _ := this.Repository.Get(benefited.Id)
			a.Book(b)

			// book kommitment share
			kommitmentShare := booking.Booking{
				Amount: this.Booking.Net[benefited] * owner.KommmitmentShare,
				Type:   booking.Kommitmentanteil,
				Text:   this.Booking.Text + "#Kommitment#" + benefited.Id,
				Month:  this.Booking.Month,
				Year:   this.Booking.Year}

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
				Year:   this.Booking.Year}
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
				CostCenter: benefited.Id}
			kommitmentAccount, _ := this.Repository.Get(owner.StakeholderKM.Id)
			kommitmentAccount.Book(kommitmentShare)
		}

		// book cost center provision
		a, _ := this.Repository.Get(this.Booking.Responsible)
		b := booking.Booking{
			Amount: this.Booking.Net[benefited] * owner.PartnerProvision,
			Type:   booking.Vertriebsprovision,
			Text:   this.Booking.Text + "#Provision#" + benefited.Id,
			Month:  this.Booking.Month,
			Year:   this.Booking.Year}
		a.Book(b)
	}
}

// Eine Buchung kann mehrere Nettopositionen enthalten, den je einem Stakeholder zugeschrieben wird.
// Diese Funktion liefert ein Array mit Stateholder, deren Nettoanteil in der Buchung > 0 ist.
func stakeholderWithNetPositions(booking booking.Booking) []owner.Stakeholder {
	var result []owner.Stakeholder

	for k, v := range booking.Net {
		if v > 0 {
			result = append(result, k)
		}
	}
	return result
}
