package processing

import (
		"time"

		"github.com/ahojsenn/kontrol/booking"
	"github.com/ahojsenn/kontrol/owner"
	"github.com/ahojsenn/kontrol/util"
	"github.com/ahojsenn/kontrol/accountSystem"
	)

type BookGehaltCommand struct {
	Booking    booking.Booking
	Repository accountSystem.AccountSystem
}

func (c BookGehaltCommand) run() {

	amount := c.Booking.Amount
	// GEhaltsbuchung ist 4120 and 1200, also Gehalt an Bank
	// Buchung Kommitment-Konto
	sollAccount, _ := c.Repository.Get(accountSystem.SKR03_4100_4199.Id)
	kBooking := booking.CloneBooking(c.Booking, -amount, booking.Gehalt, c.Booking.Responsible, c.Booking.Soll, c.Booking.Haben)
	sollAccount.Book(kBooking)


	// Bankbuchung, Haben
	bankBooking := c.Booking
	bankBooking.Type = c.Booking.Typ
	bankBooking.Amount =  amount
	bankBooking.Responsible = c.Booking.Responsible
	c.Repository.BankAccount().Book(bankBooking)

}

type BookSVBeitragCommand struct {
	Booking    booking.Booking
	Repository accountSystem.AccountSystem
}

func (c BookSVBeitragCommand) run() {

	amount := c.Booking.Amount

	// Buchung SKR03_4100_4199
	sollAccount, _ := c.Repository.Get(accountSystem.SKR03_4100_4199.Id)
	kBooking := booking.CloneBooking(c.Booking, -amount, booking.SVBeitrag, c.Booking.Responsible, c.Booking.Soll, c.Booking.Haben)
	sollAccount.Book(kBooking)

	// Bankbuchung
	bankBooking := c.Booking
	bankBooking.Type = c.Booking.Typ
	bankBooking.Amount = amount
	c.Repository.BankAccount().Book(bankBooking)
}

type BookLNSteuerCommand struct {
	Booking    booking.Booking
	Repository accountSystem.AccountSystem
}

func (c BookLNSteuerCommand) run() {

	amount := c.Booking.Amount

	// Buchung SKR03_4100_4199
	account, _ := c.Repository.Get(accountSystem.SKR03_4100_4199.Id)
	kBooking := booking.CloneBooking(c.Booking, - amount, booking.LNSteuer, c.Booking.Responsible, c.Booking.Soll, c.Booking.Haben)
	account.Book(kBooking)

	// Bankbuchung
	bankBooking := c.Booking
	bankBooking.Type = c.Booking.Typ
	bankBooking.Amount = amount
	c.Repository.BankAccount().Book(bankBooking)
}

type BookGWSteuerCommand struct {
	Booking    booking.Booking
	Repository accountSystem.AccountSystem
}

func (c BookGWSteuerCommand) run() {

	// Gewerbesteuer an Bank
	amount := c.Booking.Amount

	// Buchung Kommitment-Konto oder Rückstellung oder ...
	gwAccount,_ := c.Repository.Get(accountSystem.SKR03_Steuern.Id)
	kBooking := booking.CloneBooking(c.Booking, - amount, booking.GWSteuer, c.Booking.Responsible, c.Booking.Soll, c.Booking.Haben)
	gwAccount.Book(kBooking)

	// Bankbuchung Haben
	bankBooking := c.Booking
	bankBooking.Type = c.Booking.Typ
	bankBooking.Amount = amount
	c.Repository.BankAccount().Book(bankBooking)

}

type BookPartnerEntnahmeCommand struct {
	Booking    booking.Booking
	Repository accountSystem.AccountSystem
}

func (c BookPartnerEntnahmeCommand) run() {

	// Bankbuchung
	bankBooking := c.Booking
	bankBooking.Type = c.Booking.Typ
	bankBooking.Amount = bankBooking.Amount * -1
	c.Repository.BankAccount().Book(bankBooking)

	// Buchung gegen Kommanditstenkonto
	b := booking.CloneBooking(c.Booking, c.Booking.Amount*-1, booking.Entnahme, c.Booking.Responsible, c.Booking.Soll, c.Booking.Haben)
	a, _ := c.Repository.Get(c.Booking.Responsible)
	a.Book(b)
}

type BookPartnerEntnahmeVorjahrCommand struct {
	Booking    booking.Booking
	Repository accountSystem.AccountSystem
}

func (c BookPartnerEntnahmeVorjahrCommand) run() {

	// auflösen eines Gesellschafterdarlehens, Buchung: Gesellschafterdarlehen 920 an Bank 1200
	amount := c.Booking.Amount

	// Soll Gesellschafterdarlehens
	sollAccount,_ := c.Repository.Get(accountSystem.SKR03_KontoJUSVJ.Id)
	b := booking.CloneBooking(c.Booking, -amount, booking.GVVorjahr, c.Booking.Responsible, c.Booking.Soll, c.Booking.Haben)
	sollAccount.Book(b)

	// Haben Bankbuchung
	habenAccount := c.Repository.BankAccount()
	b2 := booking.CloneBooking(c.Booking, amount, booking.GVVorjahr, c.Booking.Responsible, c.Booking.Soll, c.Booking.Haben)
	habenAccount.Book(b2)

}

type BookEingangsrechnungCommand struct {
	Booking    booking.Booking
	Repository accountSystem.AccountSystem
}

func (c BookEingangsrechnungCommand) run() {

	// if booking with empty timestamp in position "BankCreated"
	// then book it to open positions SKR03_1600
	if c.Booking.BankCreated.After(time.Now()) {
		skr1600, _ := c.Repository.Get(accountSystem.SKR03_1600.Id)
		skr1600.Book(c.Booking)
		return
	}

	amount := c.Booking.Amount

	// Soll Buchung UST-Konto, Erträge werden im Haben gebucht, Ausgaben im Soll
	ustAccount,_ := c.Repository.Get(accountSystem.SKR03_Vorsteuer.Id)
	b2 := booking.CloneBooking(c.Booking, -1.0*(amount-util.Net(amount)), booking.Eingangsrechnung, c.Booking.Responsible, c.Booking.Soll, c.Booking.Haben)
	ustAccount.Book (b2)

	// Soll Buchung Kommitment-Konto
	sollAccount,_ := c.Repository.Get(accountSystem.SKR03_sonstigeAufwendungen.Id)
	b := booking.CloneBooking(c.Booking, -util.Net(amount), booking.Eingangsrechnung, c.Booking.Responsible, c.Booking.Soll, c.Booking.Haben)
	sollAccount.Book(b)

	// Haben Buchung Bank
	habenAccount := c.Repository.BankAccount()
	a :=  booking.CloneBooking(c.Booking, amount, booking.Eingangsrechnung, c.Booking.Responsible, c.Booking.Soll, c.Booking.Haben)
	habenAccount.Book(a)


}

type BookInterneStundenCommand struct {
	Booking    booking.Booking
	Repository accountSystem.AccountSystem
}

func (c BookInterneStundenCommand) run() {

	// Buchung interner Stunden auf Kommanditstenkonto
	a := booking.CloneBooking(c.Booking, c.Booking.Amount, booking.InterneStunden, c.Booking.Responsible, c.Booking.Soll, c.Booking.Haben)
	partnerAccount, _ := c.Repository.Get(c.Booking.Responsible)
	partnerAccount.Book(a)

	// Buchung interner Stunden von kommitment Konto auf Stakeholder
	b := booking.CloneBooking(c.Booking, c.Booking.Amount*-1, booking.InterneStunden, c.Booking.Responsible, c.Booking.Soll, c.Booking.Haben)
	kommitmentAccount, _ := c.Repository.Get(owner.StakeholderKM.Id)
	kommitmentAccount.Book(b)
}

type BookRueckstellungCommand struct {
	Booking    booking.Booking
	Repository accountSystem.AccountSystem
}

func (c BookRueckstellungCommand) run() {

	// Rückstellungsbuchung
	a := booking.CloneBooking(c.Booking, c.Booking.Amount, booking.Rueckstellung, c.Booking.Responsible, c.Booking.Soll, c.Booking.Haben)
	rueckstellungsAccount, _ := c.Repository.Get(accountSystem.SKR03_Rueckstellungen.Id)
	rueckstellungsAccount.Book(a)

	// Buchung gegen kommitment Konto
	kBooking := booking.CloneBooking(c.Booking, c.Booking.Amount*-1, booking.Rueckstellung, c.Booking.Responsible, c.Booking.Soll, c.Booking.Haben)
	account, _ := c.Repository.Get(owner.StakeholderKM.Id)
	account.Book(kBooking)
}

type BookAnfangsbestandCommand struct {
	Booking    booking.Booking
	Repository accountSystem.AccountSystem
}

func (c BookAnfangsbestandCommand) run() {

	// Anfangsbestand buchen
	a := booking.CloneBooking(c.Booking, c.Booking.Amount, booking.Anfangsbestand, c.Booking.Responsible, c.Booking.Soll, c.Booking.Haben)
	zielkonto, _ := c.Repository.Get(c.Booking.Responsible)
	zielkonto.Book(a)
}


type BookSKR03Command struct {
	Booking    booking.Booking
	Repository accountSystem.AccountSystem
}

func (c BookSKR03Command) run() {

	amount := c.Booking.Amount

	// Sollbuchung
	sollAccount := c.Repository.GetSKR03(c.Booking.Soll)
	a := booking.CloneBooking(c.Booking, -amount, c.Booking.Typ, c.Booking.Responsible, c.Booking.Soll, c.Booking.Haben)
	sollAccount.Book(a)

	// Habenbuchung
	habenAccount := c.Repository.GetSKR03(c.Booking.Haben)
	b := booking.CloneBooking(c.Booking, amount, c.Booking.Typ, c.Booking.Responsible, c.Booking.Soll, c.Booking.Haben)
	habenAccount.Book(b)
}

type BookUstCommand struct {
	Booking    booking.Booking
	Repository accountSystem.AccountSystem
}

func (c BookUstCommand) run() {

	amount := c.Booking.Amount

	// Sollbuchung
	sollAccount,_ := c.Repository.Get(accountSystem.SKR03_Umsatzsteuer.Id)
	a := booking.CloneBooking(c.Booking, -amount, c.Booking.Typ, c.Booking.Responsible, c.Booking.Soll, c.Booking.Haben)
	sollAccount.Book(a)

	// Habenbuchung
	habenAccount := c.Repository.BankAccount()
	b := booking.CloneBooking(c.Booking, amount, c.Booking.Typ, c.Booking.Responsible, c.Booking.Soll, c.Booking.Haben)
	habenAccount.Book(b)
}

