package processing

import (
	"github.com/ahojsenn/kontrol/account"
	"github.com/ahojsenn/kontrol/booking"
	"github.com/ahojsenn/kontrol/owner"
	"github.com/ahojsenn/kontrol/util"
)

type Command interface {
	run()
}

type BookRueckstellungCommand struct {
	Booking    booking.Booking
	Repository account.Repository
}

func (c BookRueckstellungCommand) run() {

	// Bankbuchung
	bankBooking := c.Booking
	bankBooking.Type = c.Booking.Typ
	bankBooking.Amount = bankBooking.Amount * -1
	c.Repository.BankAccount().Book(bankBooking)

	// Buchung Kommitment-Konto
	kBooking := booking.CloneBooking(c.Booking, c.Booking.Amount*-1, booking.Rueckstellung, c.Booking.Responsible)
	kommitmentAccount, _ := c.Repository.Get(owner.StakeholderKM.Id)
	kommitmentAccount.Book(kBooking)
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
	kBooking := booking.CloneBooking(c.Booking, c.Booking.Amount*-1, booking.Gehalt, c.Booking.Responsible)
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
	kBooking := booking.CloneBooking(c.Booking, c.Booking.Amount*-1, booking.SVBeitrag, c.Booking.Responsible)
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
	kBooking := booking.CloneBooking(c.Booking, c.Booking.Amount*-1, booking.LNSteuer, c.Booking.Responsible)
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
	kBooking := booking.CloneBooking(c.Booking, c.Booking.Amount*-1, booking.GWSteuer, c.Booking.Responsible)
	kommitmentAccount, _ := c.Repository.Get(owner.StakeholderKM.Id)
	kommitmentAccount.Book(kBooking)
}

type BookPartnerEntnahmeCommand struct {
	Booking    booking.Booking
	Repository account.Repository
}

func (c BookPartnerEntnahmeCommand) run() {

	// Bankbuchung
	bankBooking := c.Booking
	bankBooking.Type = c.Booking.Typ
	bankBooking.Amount = bankBooking.Amount * -1
	c.Repository.BankAccount().Book(bankBooking)

	// Buchung gegen Kommanditstenkonto
	b := booking.CloneBooking(c.Booking, c.Booking.Amount*-1, booking.Entnahme, c.Booking.Responsible)
	a, _ := c.Repository.Get(c.Booking.Responsible)
	a.Book(b)
}

type BookEingangsrechnungCommand struct {
	Booking    booking.Booking
	Repository account.Repository
}

func (c BookEingangsrechnungCommand) run() {

	// Buchung Kommitment-Konto
	b := booking.CloneBooking(c.Booking, util.Net(c.Booking.Amount)*-1, booking.Eingangsrechnung, c.Booking.Responsible)
	kommitmentAccount, _ := c.Repository.Get(owner.StakeholderKM.Id)
	kommitmentAccount.Book(b)
}
