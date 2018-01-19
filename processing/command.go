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
