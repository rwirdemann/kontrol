package processing

import (
	"bitbucket.org/rwirdemann/kontrol/account"
	"bitbucket.org/rwirdemann/kontrol/booking"
	"bitbucket.org/rwirdemann/kontrol/owner"
)

type Command interface {
	run(r account.Repository, b booking.Booking)
}

type BookGehaltCommand struct{}

func (BookGehaltCommand) run(r account.Repository, b booking.Booking) {

	// Bankbuchung
	bankBooking := b
	bankBooking.Type = b.Typ
	bankBooking.Amount = bankBooking.Amount * -1
	r.BankAccount().Book(bankBooking)

	// Buchung Kommitment-Konto
	kBooking := booking.CloneBooking(b, -1, booking.Gehalt, b.Responsible)
	kommitmentAccount, _ := r.Get(owner.StakeholderKM.Id)
	kommitmentAccount.Book(kBooking)
}

type BookSVBeitragCommand struct{}

func (BookSVBeitragCommand) run(r account.Repository, b booking.Booking) {

	// Bankbuchung
	bankBooking := b
	bankBooking.Type = b.Typ
	bankBooking.Amount = bankBooking.Amount * -1
	r.BankAccount().Book(bankBooking)

	// Buchung Kommitment-Konto
	kBooking := booking.CloneBooking(b, -1, booking.SVBeitrag, b.Responsible)
	kommitmentAccount, _ := r.Get(owner.StakeholderKM.Id)
	kommitmentAccount.Book(kBooking)
}
