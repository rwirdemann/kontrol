package processing

import (
	"log"
	"time"

	"github.com/ahojsenn/kontrol/account"
	"github.com/ahojsenn/kontrol/booking"
	"github.com/ahojsenn/kontrol/owner"
	"github.com/ahojsenn/kontrol/util"
)

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
	account2, _ := c.Repository.Get(owner.SKR03_4100_4199.Id)
	account2.Book(kBooking)
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
	account2, _ := c.Repository.Get(owner.SKR03_4100_4199.Id)
	account2.Book(kBooking)
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
	account2, _ := c.Repository.Get(owner.SKR03_4100_4199.Id)
	account2.Book(kBooking)
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

	// Buchung Kommitment-Konto oder Rückstellung oder ...
	kBooking := booking.CloneBooking(c.Booking, c.Booking.Amount*-1, booking.GWSteuer, c.Booking.Responsible)
	kommitmentAccount, _ := c.Repository.Get(c.Booking.Responsible)
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

type BookPartnerEntnahmeVorjahrCommand struct {
	Booking    booking.Booking
	Repository account.Repository
}

func (c BookPartnerEntnahmeVorjahrCommand) run() {

	// Bankbuchung
	bankBooking := c.Booking
	bankBooking.Type = c.Booking.Typ
	bankBooking.Amount = bankBooking.Amount * -1
	c.Repository.BankAccount().Book(bankBooking)

	// Buchung gegen Kommanditstenkonto
	b := booking.CloneBooking(c.Booking, c.Booking.Amount*-1, booking.GVVorjahr, c.Booking.Responsible)
	a, _ := c.Repository.Get(owner.KontoJUSVJ.Id)
	a.Book(b)
}

type BookEingangsrechnungCommand struct {
	Booking    booking.Booking
	Repository account.Repository
}

func (c BookEingangsrechnungCommand) run() {

	// if booking with empty timestamp in position "BankCreated"
	// the book it to open positions SKR03_1600
	log.Println("in BookEingangsrechnungCommand: ", c.Booking.BankCreated)
	if c.Booking.BankCreated.After(time.Now()) {
		skr1600, _ := c.Repository.Get(owner.SKR03_1600.Id)
		skr1600.Book(c.Booking)
		return
	}

	// Buchung Kommitment-Konto
	b := booking.CloneBooking(c.Booking, util.Net(c.Booking.Amount)*-1, booking.Eingangsrechnung, c.Booking.Responsible)
	kommitmentAccount, _ := c.Repository.Get(owner.StakeholderKM.Id)
	kommitmentAccount.Book(b)
}

type BookInterneStundenCommand struct {
	Booking    booking.Booking
	Repository account.Repository
}

func (c BookInterneStundenCommand) run() {

	// Buchung interner Stunden auf Kommanditstenkonto
	a := booking.CloneBooking(c.Booking, c.Booking.Amount, booking.InterneStunden, c.Booking.Responsible)
	partnerAccount, _ := c.Repository.Get(c.Booking.Responsible)
	partnerAccount.Book(a)

	// Buchung interner Stunden von kommitment Konto auf Stakeholder
	b := booking.CloneBooking(c.Booking, c.Booking.Amount*-1, booking.InterneStunden, c.Booking.Responsible)
	kommitmentAccount, _ := c.Repository.Get(owner.StakeholderKM.Id)
	kommitmentAccount.Book(b)
}

type BookRueckstellungCommand struct {
	Booking    booking.Booking
	Repository account.Repository
}

func (c BookRueckstellungCommand) run() {

	// Rückstellungsbuchung
	a := booking.CloneBooking(c.Booking, c.Booking.Amount, booking.Rueckstellung, c.Booking.Responsible)
	rueckstellungsAccount, _ := c.Repository.Get(owner.StakeholderRueckstellung.Id)
	rueckstellungsAccount.Book(a)

	// Buchung gegen kommitment Konto
	kBooking := booking.CloneBooking(c.Booking, c.Booking.Amount*-1, booking.Rueckstellung, c.Booking.Responsible)
	account, _ := c.Repository.Get(owner.StakeholderKM.Id)
	account.Book(kBooking)
}

type BookAnfangsbestandCommand struct {
	Booking    booking.Booking
	Repository account.Repository
}

func (c BookAnfangsbestandCommand) run() {

	// Anfangsbestand buchen
	a := booking.CloneBooking(c.Booking, c.Booking.Amount, booking.Anfangsbestand, c.Booking.Responsible)
	zielkonto, _ := c.Repository.Get(c.Booking.Responsible)
	zielkonto.Book(a)
}

type BookERgegenRückstellungCommand struct {
	Booking    booking.Booking
	Repository account.Repository
}

func (c BookERgegenRückstellungCommand) run() {

	// Buchung gegen Rückstellungskonto
	a := booking.CloneBooking(c.Booking, util.Net(c.Booking.Amount)*-1, booking.Eingangsrechnung, c.Booking.Responsible)
	rückstellungskonto, _ := c.Repository.Get(owner.StakeholderRueckstellung.Id)
	rückstellungskonto.Book(a)
}

type BookRückstellungAuflösenCommand struct {
	Booking    booking.Booking
	Repository account.Repository
}

func (c BookRückstellungAuflösenCommand) run() {

	// Buchung vom Rückstellungskonto
	a := booking.CloneBooking(c.Booking, c.Booking.Amount, booking.Eingangsrechnung, c.Booking.Responsible)
	rückstellungskonto, _ := c.Repository.Get(owner.StakeholderRueckstellung.Id)
	rückstellungskonto.Book(a)

	// Buchung auf das kommitment Konto
	kBooking := booking.CloneBooking(c.Booking, c.Booking.Amount*-1.0, booking.RueckstellungAuflösen, c.Booking.Responsible)
	account, _ := c.Repository.Get(owner.StakeholderKM.Id)
	account.Book(kBooking)
}
