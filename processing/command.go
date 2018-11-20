package processing

import (
	"github.com/ahojsenn/kontrol/accountSystem"
	"github.com/ahojsenn/kontrol/booking"
	"github.com/ahojsenn/kontrol/util"
	"github.com/ahojsenn/kontrol/valueMagnets"
)

type BookGehaltCommand struct {
	Booking    booking.Booking
	AccSystem  accountSystem.AccountSystem
}

func (c BookGehaltCommand) run() {

	amount := c.Booking.Amount
	// Gehaltsbuchung ist 4120 and 1200, also CC_Gehalt an Bank
	// Buchung Kommitment-Konto
	sollAccount, _ := c.AccSystem.Get(accountSystem.SKR03_4100_4199.Id)
	kBooking := booking.CloneBooking(c.Booking, -amount, booking.CC_Gehalt, c.Booking.Responsible, c.Booking.Soll, c.Booking.Haben, c.Booking.Project)
	sollAccount.Book(kBooking)

	// Bankbuchung, Haben
	bankBooking := c.Booking
	bankBooking.Type = booking.CC_Gehalt
	bankBooking.Amount =  amount
	bankBooking.Responsible = c.Booking.Responsible
	acc,_ := c.AccSystem.Get(accountSystem.SKR03_1200.Id)
	acc.Book(bankBooking)

}

type BookSVBeitragCommand struct {
	Booking    booking.Booking
	AccSystem  accountSystem.AccountSystem
}

func (c BookSVBeitragCommand) run() {

	amount := c.Booking.Amount

	// Buchung SKR03_4100_4199
	sollAccount, _ := c.AccSystem.Get(accountSystem.SKR03_4100_4199.Id)
	kBooking := booking.CloneBooking(c.Booking, -amount, booking.CC_SVBeitrag, c.Booking.Responsible, c.Booking.Soll, c.Booking.Haben, c.Booking.Project)
	sollAccount.Book(kBooking)

	// Habenbuchung
	habenAccountId := ""
	if c.Booking.IsBeyondBudgetDate() {
		habenAccountId = accountSystem.SKR03_Rueckstellungen.Id
	} else {
		habenAccountId = accountSystem.SKR03_1200.Id
	}
	habenAccount,_ := c.AccSystem.Get(habenAccountId)
	bk := c.Booking
	bk.Type = booking.CC_SVBeitrag
	bk.Amount = amount
	habenAccount.Book(bk)

}

type BookLNSteuerCommand struct {
	Booking    booking.Booking
	AccSystem  accountSystem.AccountSystem
}

func (c BookLNSteuerCommand) run() {

	amount := c.Booking.Amount

	// Buchung SKR03_4100_4199
	account, _ := c.AccSystem.Get(accountSystem.SKR03_4100_4199.Id)
	kBooking := booking.CloneBooking(c.Booking, - amount, booking.CC_LNSteuer, c.Booking.Responsible, c.Booking.Soll, c.Booking.Haben, c.Booking.Project)
	account.Book(kBooking)

	// Habenbuchung
	habenAccountId := ""
	if c.Booking.IsBeyondBudgetDate() {
		habenAccountId = accountSystem.SKR03_Rueckstellungen.Id
	} else {
		habenAccountId = accountSystem.SKR03_1200.Id
	}
	habenAccount,_ := c.AccSystem.Get(habenAccountId)
	bk := c.Booking
	bk.Type = booking.CC_LNSteuer
	bk.Amount = amount
	habenAccount.Book(bk)

}

type BookGWSteuerCommand struct {
	Booking    booking.Booking
	AccSystem  accountSystem.AccountSystem
}

func (c BookGWSteuerCommand) run() {

	// Gewerbesteuer an Bank
	amount := c.Booking.Amount

	// Buchung Kommitment-Konto oder Rückstellung oder ...
	gwAccount,_ := c.AccSystem.Get(accountSystem.SKR03_Steuern.Id)
	kBooking := booking.CloneBooking(c.Booking, - amount, booking.CC_GWSteuer, c.Booking.Responsible, c.Booking.Soll, c.Booking.Haben, c.Booking.Project)
	gwAccount.Book(kBooking)

	// Habenbuchung
	habenAccountId := ""
	if c.Booking.IsBeyondBudgetDate() {
		habenAccountId = accountSystem.SKR03_Rueckstellungen.Id
	} else {
		habenAccountId = accountSystem.SKR03_1200.Id
	}
	habenAccount,_ := c.AccSystem.Get(habenAccountId)
	habenbk := c.Booking
	habenbk.Type = booking.CC_GWSteuer
	habenbk.Amount = amount
	habenAccount.Book(habenbk)

}

type BookPartnerEntnahmeCommand struct {
	Booking    booking.Booking
	AccSystem  accountSystem.AccountSystem
}

func (c BookPartnerEntnahmeCommand) run() {

	// auflösen eines Gesellschafterdarlehens, Buchung: Privatentnahme 1900 an Bank 1200
	amount := c.Booking.Amount

	// Soll Privatentnahme
	sollAccount,_ := c.AccSystem.Get(accountSystem.SKR03_1900.Id)
	b := booking.CloneBooking(c.Booking, -amount, booking.CC_Entnahme, c.Booking.Responsible, c.Booking.Soll, c.Booking.Haben, c.Booking.Project)
	sollAccount.Book(b)

	// an Bank
	bankBooking := c.Booking
	bankBooking.Type = booking.CC_Entnahme
	bankBooking.Amount = bankBooking.Amount
	acc,_ := c.AccSystem.Get(accountSystem.SKR03_1200.Id)
	acc.Book(bankBooking)

}

type BookPartnerEntnahmeVorjahrCommand struct {
	Booking    booking.Booking
	AccSystem  accountSystem.AccountSystem
}

func (c BookPartnerEntnahmeVorjahrCommand) run() {

	// auflösen eines Gesellschafterdarlehens, Buchung: Gesellschafterdarlehen 920 an Bank 1200
	amount := c.Booking.Amount
	//c.Booking.Soll = "920"
	//c.Booking.Haben = "1200"

	// Soll Gesellschafterdarlehens
	sollAccount,_ := c.AccSystem.Get(accountSystem.SKR03_920_Gesellschafterdarlehen.Id)
	b := booking.CloneBooking(c.Booking, -amount, c.Booking.Typ, c.Booking.Responsible, c.Booking.Soll, c.Booking.Haben, c.Booking.Project)
	sollAccount.Book(b)

	// Haben Bankbuchung
	habenAccount,_ := c.AccSystem.Get(accountSystem.SKR03_1200.Id)
	b2 := booking.CloneBooking(c.Booking, amount, c.Booking.Typ, c.Booking.Responsible, c.Booking.Soll, c.Booking.Haben, c.Booking.Project)
	habenAccount.Book(b2)

}

type BookEingangsrechnungCommand struct {
	Booking    booking.Booking
	AccSystem  accountSystem.AccountSystem
}

func (c BookEingangsrechnungCommand) run() {

	// if booking with empty timestamp in position "BankCreated"
	// then book it to open positions SKR03_1600


	if c.Booking.IsOpenPosition() {
		skr1600, _ := c.AccSystem.Get(accountSystem.SKR03_1600.Id)
		skr1600.Book(c.Booking)
		return
	}

	amount := c.Booking.Amount

	// Soll Buchung UST-Konto, Erträge werden im Haben gebucht, Ausgaben im Soll
	ustAccount,_ := c.AccSystem.Get(accountSystem.SKR03_Vorsteuer.Id)
	b2 := booking.CloneBooking(c.Booking, -1.0*(amount-util.Net(amount)), booking.Ust, c.Booking.Responsible, c.Booking.Soll, c.Booking.Haben, c.Booking.Project)
	ustAccount.Book (b2)

	// Soll Buchung Kommitment-Konto
	sollAccount,_ := c.AccSystem.Get(accountSystem.SKR03_sonstigeAufwendungen.Id)
	b := booking.CloneBooking(c.Booking, -util.Net(amount), booking.Kosten, c.Booking.Responsible, c.Booking.Soll, c.Booking.Haben, c.Booking.Project)
	sollAccount.Book(b)

	// Haben Buchung Bank if not booking.IsBeyondBudgetDate()
	// otherwise in accountSystem.SKR03_1600.Id
	habenAccountId := ""
	if c.Booking.IsBeyondBudgetDate() {
		habenAccountId = accountSystem.SKR03_1600.Id
	} else {
		habenAccountId = accountSystem.SKR03_1200.Id
	}
	a :=  booking.CloneBooking(c.Booking, amount, booking.Kosten, c.Booking.Responsible, c.Booking.Soll, c.Booking.Haben, c.Booking.Project)
	habenAccount,_ := c.AccSystem.Get(habenAccountId)
	habenAccount.Book(a)

}

type BookInterneStundenCommand struct {
	Booking    booking.Booking
	AccSystem  accountSystem.AccountSystem
}

func (c BookInterneStundenCommand) run() {

	// Buchung interner Stunden auf Kommanditstenkonto
	a := booking.CloneBooking(c.Booking, c.Booking.Amount, booking.CC_InterneStunden, c.Booking.Responsible, c.Booking.Soll, c.Booking.Haben, c.Booking.Project)
	partnerAccount, _ := c.AccSystem.Get(c.Booking.Responsible)
	partnerAccount.Book(a)

	// Buchung interner Stunden von kommitment Konto auf Stakeholder
	b := booking.CloneBooking(c.Booking, c.Booking.Amount*-1, booking.CC_InterneStunden, c.Booking.Responsible, c.Booking.Soll, c.Booking.Haben, c.Booking.Project)
	kommitmentAccount, _ := c.AccSystem.Get(valueMagnets.StakeholderKM.Id)
	kommitmentAccount.Book(b)

}


type BookSKR03Command struct {
	Booking    		booking.Booking
	AccSystem  		accountSystem.AccountSystem
}

func (c BookSKR03Command) run() {

	amount := c.Booking.Amount

	// Sollbuchung
	sollAccount := c.AccSystem.GetSKR03(c.Booking.Soll)
	a := booking.CloneBooking(c.Booking, -amount, c.Booking.Typ, c.Booking.Responsible, c.Booking.Soll, c.Booking.Haben, c.Booking.Project)
	sollAccount.Book(a)

	// Habenbuchung
	habenAccount := c.AccSystem.GetSKR03(c.Booking.Haben)
	b := booking.CloneBooking(c.Booking, amount, c.Booking.Typ, c.Booking.Responsible, c.Booking.Soll, c.Booking.Haben, c.Booking.Project)
	habenAccount.Book(b)

}

type BookUstCommand struct {
	Booking    booking.Booking
	AccSystem  accountSystem.AccountSystem
}

func (c BookUstCommand) run() {

	amount := c.Booking.Amount

	// Sollbuchung
	sollAccount,_ := c.AccSystem.Get(accountSystem.SKR03_Umsatzsteuer.Id)
	a := booking.CloneBooking(c.Booking, -amount, c.Booking.Typ, c.Booking.Responsible, c.Booking.Soll, c.Booking.Haben, c.Booking.Project)
	sollAccount.Book(a)

	// Habenbuchung
	habenAccountId := ""
	if c.Booking.IsBeyondBudgetDate() {
		habenAccountId = accountSystem.SKR03_Rueckstellungen.Id
	} else {
		habenAccountId = accountSystem.SKR03_1200.Id
	}
	habenAccount,_ := c.AccSystem.Get(habenAccountId)
	b := booking.CloneBooking(c.Booking, amount, c.Booking.Typ, c.Booking.Responsible, c.Booking.Soll, c.Booking.Haben, c.Booking.Project)
	habenAccount.Book(b)

}

