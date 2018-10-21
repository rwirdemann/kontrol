package processing

import (
	"testing"
	"time"

	"github.com/ahojsenn/kontrol/account"

	"github.com/ahojsenn/kontrol/accountSystem"
	"github.com/ahojsenn/kontrol/booking"
	"github.com/ahojsenn/kontrol/util"
	"github.com/ahojsenn/kontrol/valueMagnets"
	"github.com/stretchr/testify/assert"
)

var accSystem accountSystem.AccountSystem
var accountBank *account.Account
var accountHannes *account.Account
var accountRalf *account.Account
var accountKommitment *account.Account

func setUp() {
	accSystem = accountSystem.NewDefaultAccountSystem()
	accountBank = accSystem.GetCollectiveAccount()
	accountHannes, _ = accSystem.Get(valueMagnets.StakeholderRepository{}.Get("JM").Id)
	accountRalf, _ = accSystem.Get(valueMagnets.StakeholderRepository{}.Get("RW").Id)
	accountKommitment, _ = accSystem.Get(valueMagnets.StakeholderRepository{}.Get("K").Id)
	util.Global.BalanceDate = time.Date(2018, 1, 24, 0, 0, 0, 0, time.UTC)

}


// CC_Gehalt Angestellter
// - 100% Brutto gegen Bankkonto
// - 100% Brutto gegen Kommitmentkonto
// - Kostenstelle: Kürzel des Angestellten
func TestGehaltAngestellter(t *testing.T) {
	repository := accountSystem.NewDefaultAccountSystem()

	// given: a booking
	its2018 := time.Date(2018, 1, 23, 0, 0, 0, 0, time.UTC)
	p := booking.NewBooking(13, "Gehalt", "", "", "BW", nil, 3869.65, "Gehalt Ben", 1, 2017, its2018)

	// when: the position is processed
	Process(repository, *p)

	// 100% Brutto gegen SKR03_4100_4199
	account2, _ := repository.Get(accountSystem.SKR03_4100_4199.Id)
	util.AssertEquals(t, 1, len(account2.Bookings))
	assert.Equal(t, -3869.65, account2.Bookings[0].Amount)
	assert.Equal(t, booking.CC_Gehalt, account2.Bookings[0].Type)


	// 100% Brutto gegen Bankkonto
	accountBank,_ := repository.Get(accountSystem.SKR03_1200.Id)
	assert.Equal(t, 3869.65, accountBank.Bookings[0].Amount)
	assert.Equal(t, "Gehalt Ben", accountBank.Bookings[0].Text)
	assert.Equal(t, "CC_Gehalt", accountBank.Bookings[0].Type)


	// Kommitment-Buchung ist der Kostenstelle "BW" zugeordnet
	assert.Equal(t, "BW", account2.Bookings[0].CostCenter)
}

// Eingangsrechnungen
// - 100% werden netto vom Bankkonto gebucht
// - 100% des Nettobetrags werden gegen das SKR03_sonstigeAufwendungen gebucht
func TestEingangsrechnung(t *testing.T) {
	setUp()

	// given: BOOKING ER
	// Eingangsrechnung 12852.0€ von Bank an SKR03_sonstigeAufwendungen
	its2018 := time.Date(2018, 1, 23, 0, 0, 0, 0, time.UTC)
	p := booking.NewBooking(13, "ER", "", "", "K", nil, 12852.0, "Eingangsrechnung 1234", 1, 2017, its2018)

	// when: the position is processed
	Process(accSystem, *p)

	// Soll Buchung wurde von SKR03_sonstigeAufwendungen gebucht, Achtung Passivkonto, da werden Soll auf die Haben Seite gebucht
	a, _ := accSystem.Get(accountSystem.SKR03_sonstigeAufwendungen.Id)
	assert.Equal(t, 1, len(a.Bookings))
	bk := a.Bookings[0]
	assert.Equal(t, util.Net(-12852.0), bk.Amount)
	assert.Equal(t, bk.Type, booking.Kosten)
	assert.Equal(t, "K", bk.CostCenter)

	//  Haben wurde auf das Bankkonto gebucht, Achtung Bank ist Aktivkonto, da werden Soll Eintrage im Haben gebucht
	habenAccount,_ := accSystem.Get(accountSystem.SKR03_1200.Id)
	assert.Equal(t, 1, len(habenAccount.Bookings))
	actual := habenAccount.Bookings[0]
	assert.Equal(t,12852.0, actual.Amount)
	assert.Equal(t, "Eingangsrechnung 1234", actual.Text)
	assert.Equal(t, booking.Kosten, actual.Type)


}

func TestRueckstellungAufloesen(t *testing.T) {
	setUp()
	// Rückstellungen können gegen das kommitment Konto aufgelöst werden

	// given a Buchung Eingangsrechnung gegen Rücksttellung
	its2018 := time.Date(2018, 1, 23, 0, 0, 0, 0, time.UTC)
	p := booking.NewBooking(13,"SKR03", "965", "4957", "K", nil, 12852.0, "Auflösung Rückstellungsdifferenz", 1, 2017, its2018)

	// when: the position is processed
	Process(accSystem, *p)

	// the booking is booked from Rückstellung account
	a1, _ := accSystem.Get(accountSystem.SKR03_Rueckstellungen.Id)
	util.AssertEquals(t, 1, len(a1.Bookings))
	b1 := a1.Bookings[0]
	util.AssertFloatEquals(t, -12852.0, b1.Amount)
	util.AssertEquals(t, booking.SKR03, b1.Type)

	// the booking is not booked to the bankaccout
	util.AssertEquals(t, 0, len(accSystem.GetCollectiveAccount().Bookings))

	// the booking is  booked on SKR03_sonstigeAufwendungen account
	a2, _ := accSystem.Get(accountSystem.SKR03_sonstigeAufwendungen.Id)
	assert.Equal(t, 1, len(a2.Bookings))
	bk := a2.Bookings[0]
	assert.Equal(t, 12852.0, bk.Amount)
	assert.Equal(t, bk.Type, booking.SKR03)
	assert.Equal(t, "K", bk.CostCenter)
}


func TestAnfangsbestandRueckstellung(t *testing.T) {
	setUp()
	// Rückstellungen können gegen das kommitment Konto aufgelöst werden

	// given a Buchung Eingangsrechnung gegen Rücksttellung
	its2018 := time.Date(2018, 1, 23, 0, 0, 0, 0, time.UTC)
	p := booking.NewBooking(13,"SKR03", "9000", "956", "K", nil, 12852.0, "Anfangsbestand GWteuewrRückst.", 1, 2017, its2018)

	// when: the position is processed
	Process(accSystem, *p)

	// the booking is booked from Saldenvortrag 9000 account
	a1, _ := accSystem.Get(accountSystem.SKR03_Saldenvortrag.Id)
	util.AssertEquals(t, 1, len(a1.Bookings))
	b1 := a1.Bookings[0]
	util.AssertFloatEquals(t, -12852.0, b1.Amount)
	util.AssertEquals(t, booking.SKR03, b1.Type)

	// the booking is not booked to the bankaccout
	util.AssertEquals(t, 0, len(accSystem.GetCollectiveAccount().Bookings))

	// the booking is  booked on Rückstellung account
	a2, _ := accSystem.Get(accountSystem.SKR03_Rueckstellungen.Id)
	assert.Equal(t, 1, len(a2.Bookings))
	bk := a2.Bookings[0]
	assert.Equal(t, 12852.0, bk.Amount)
	assert.Equal(t, bk.Type, booking.SKR03)
	assert.Equal(t, "K", bk.CostCenter)
}

func TestPartnerEntnahme(t *testing.T) {
	setUp()

	extras := booking.CsvBookingExtras{Typ: "GV", Responsible: "RW"}
	extras.Net = make(map[valueMagnets.Stakeholder]float64)
	its2018 := time.Date(2018, 1, 23, 0, 0, 0, 0, time.UTC)
	b := booking.NewBooking(13,"GV", "", "", "RW", nil, 6000, "", 1, 2017, its2018)

	Process(accSystem, *b)

	acc1,_ := accSystem.Get(accountSystem.SKR03_1900.Id)
	bRalf := acc1.Bookings[0]
	util.AssertFloatEquals(t, -6000, bRalf.Amount)
	util.AssertEquals(t, booking.CC_Entnahme, bRalf.Type)

	// Buchung wurde gegen das Bankkonto gebucht
	acc,_ := accSystem.Get(accountSystem.SKR03_1200.Id)
	util.AssertEquals(t, 1, len(acc.Bookings))
	actual := acc.Bookings[0]
	util.AssertFloatEquals(t, 6000, actual.Amount)
	util.AssertEquals(t, booking.CC_Entnahme, actual.Type)
}

// Rückstellungen
// - werden nicht auf das Bankkonto gebucht
// - 100% werden auf das Rückstellung-Konto gebucht
// - 100% werden gegen das Kommitment-Konto gebucht
func TestRueckstellung(t *testing.T) {
	setUp()

	// given: a Rückstellung booking
	its2018 := time.Date(2018, 1, 23, 0, 0, 0, 0, time.UTC)
	p := booking.NewBooking(13,"SKR03", "4120", "965", "BW", nil, 4711.0, "Bonus Rückstellung", 12, 2017, its2018)

	// when: the position is processed
	Process(accSystem, *p)

	// the booking is booked to Rückstellung account
	a1, _ := accSystem.Get(accountSystem.SKR03_Rueckstellungen.Id)
	util.AssertEquals(t, 1, len(a1.Bookings))
	b1 := a1.Bookings[0]
	util.AssertFloatEquals(t, 4711.00, b1.Amount)
	util.AssertEquals(t, booking.SKR03, b1.Type)

	// the booking is booked against kommitment account
	a2, _ := accSystem.Get(accountSystem.SKR03_4100_4199.Id)
	b2 := a2.Bookings[0]
	util.AssertFloatEquals(t, -4711.00, b2.Amount)
	util.AssertEquals(t, booking.SKR03, b1.Type)

	// Rückstellungen are not booked on bank account
	util.AssertEquals(t, 0, len(accSystem.GetCollectiveAccount().Bookings))
}

// Interne Stunden
// - werden nicht auf das Bankkonto gebucht
// - 100% werden auf das Rückstellung-Konto gebucht
// - 100% werden gegen das Kommitment-Konto gebucht
func TestInterneStunden(t *testing.T) {
	setUp()

	// given: a internal hours booking
	its2018 := time.Date(2018, 1, 23, 0, 0, 0, 0, time.UTC)
	p := booking.NewBooking(13,"IS", "", "", "AN", nil, 8250.0, "Interne Stunden 2017", 12, 2017, its2018)

	// when: the position is processed
	Process(accSystem, *p)

	// the booking is booked to anke's account
	a1, _ := accSystem.Get(valueMagnets.StakeholderRepository{}.Get("AN").Id)
	util.AssertEquals(t, 1, len(a1.Bookings))
	b1 := a1.Bookings[0]
	util.AssertFloatEquals(t, 8250.00, b1.Amount)
	util.AssertEquals(t, booking.CC_InterneStunden, b1.Type)

	// the booking is booked against kommitment account
	a2, _ := accSystem.Get(valueMagnets.StakeholderKM.Id)
	b2 := a2.Bookings[0]
	util.AssertFloatEquals(t, -8250.00, b2.Amount)
	util.AssertEquals(t, booking.CC_InterneStunden, b1.Type)

	// internal hours are not booked on bank account
	util.AssertEquals(t, 0, len(accSystem.GetCollectiveAccount().Bookings))
}

func TestBookAusgangsrechnungToBankAccount(t *testing.T) {
	setUp()
	its2018 := time.Date(2018, 1, 23, 0, 0, 0, 0, time.UTC)
	b := booking.NewBooking(13,"AR", "", "", "K", nil, 6000, "Ausgangsrechnung", 1, 2017, its2018)

	Process(accSystem, *b)
	acc, _ := accSystem.Get(accountSystem.SKR03_1200.Id)
	util.AssertEquals(t, 1, len(acc.Bookings))
	actual := acc.Bookings[0]
	util.AssertFloatEquals(t, -6000, actual.Amount)
	util.AssertEquals(t, "Ausgangsrechnung", actual.Text)
	util.AssertEquals(t, "Erloese", actual.Type)
}

// 100% werden auf das Bankkonto gebucht
// 100% werden gegen das Kommitment-Konto gebucht
func TestProcessSVBeitrag(t *testing.T) {
	setUp()
	its2018 := time.Date(2018, 1, 23, 0, 0, 0, 0, time.UTC)
	b := booking.NewBooking(13,"SV-Beitrag", "", "", "BW", nil, 1385.10, "KKH, Ben", 5, 2017, its2018)

	Process(accSystem, *b)

	// Buchung wurde gegen Kommitment-Konto gebucht
	a, _ := accSystem.Get(accountSystem.SKR03_4100_4199.Id)
	b1 := a.Bookings[0]
	assert.Equal(t, -1385.10, b1.Amount)
	assert.Equal(t, booking.CC_SVBeitrag, b1.Type)
	assert.Equal(t, "BW", b1.CostCenter)

	// Buchung wurde aufs Bankkonto gebucht
	acc,_ := accSystem.Get(accountSystem.SKR03_1200.Id)
	assert.Equal(t, 1, len(acc.Bookings))
	actual := acc.Bookings[0]
	assert.Equal(t, 1385.10, actual.Amount)
	assert.Equal(t, "KKH, Ben", actual.Text)
	assert.Equal(t, booking.CC_SVBeitrag, actual.Type)
}

// 100% werden auf das Bankkonto gebucht
// 100% werden gegen das Kommitment-Konto gebucht
// Kostenstelle: Angestellter, für den Lohnsteuer gezahlt wurde
func TestProcessLNSteuer(t *testing.T) {
	setUp()
	its2018 := time.Date(2018, 1, 23, 0, 0, 0, 0, time.UTC)
	b := booking.NewBooking(13,"LNSteuer", "", "", "BW", nil, 1511.45, "Lohnsteuer Ben", 5, 2017, its2018)

	Process(accSystem, *b)

	// Buchung wurde gegen Kommitment-Konto gebucht
	account2, _ := accSystem.Get(accountSystem.SKR03_4100_4199.Id)
	assertBooking(t, account2.Bookings[0], -1511.45, "Lohnsteuer Ben", "CC_LNSteuer")
	assert.Equal(t, "BW", account2.Bookings[0].CostCenter)

	// Buchung wurde aufs Bankkonto gebucht
	bacc, _ := accSystem.Get(accountSystem.SKR03_1200.Id)
	assertBooking(t, bacc.Bookings[0], 1511.45, "Lohnsteuer Ben", "CC_LNSteuer")
}

// 100% werden auf das Bankkonto gebucht
// 100% werden gegen das Kommitment-Konto gebucht. Diese Regel ist nicht unscharf:
// eigentlich müssen die 100% aufgeteilt werden auf: 70% auf Partner, 25% auf
// Kommitment und 5% auf Dealbringer
func TestProcessGWSteuer(t *testing.T) {
	setUp()

	its2018 := time.Date(2018, 1, 23, 0, 0, 0, 0, time.UTC)
	b := booking.NewBooking(13,"GWSteuer", "", "", "K", nil, 2385.10, "STEUERKASSE HAMBURG STEUERNR 048/638/01147 GEW.ST 4VJ.17", 9, 2017, its2018)

	Process(accSystem, *b)


	// Buchung wurde gegen Gewerbesteuer Konto gebucht
	a, _ := accSystem.Get(accountSystem.SKR03_Steuern.Id)
	b1 := a.Bookings[0]
	assertBooking(t, b1, -2385.10, "STEUERKASSE HAMBURG STEUERNR 048/638/01147 GEW.ST 4VJ.17", booking.CC_GWSteuer)

	// Buchung wurde aufs Bankkonto gebucht
	acc,_ := accSystem.Get(accountSystem.SKR03_1200.Id)
	util.AssertEquals(t, 1, len(acc.Bookings))
	actual := acc.Bookings[0]
	assertBooking(t, actual, 2385.10, "STEUERKASSE HAMBURG STEUERNR 048/638/01147 GEW.ST 4VJ.17", booking.CC_GWSteuer)
}



// 100% werden auf das Bankkonto gebucht
// 100% werden gegen das JahresüberschussVJ gebucht
func TestProcessGV_Vorjahr(t *testing.T) {
	setUp()
	b := booking.NewBooking(13,"GV-Vorjahr", "", "", "JM", nil, 77777, "Rest Anteil Johannes", 5, 2017, time.Time{})

	Process(accSystem, *b)

	// Buchung wurde gegen JahresüberschussVJ gebucht
	a, _ := accSystem.Get(accountSystem.SKR03_920_Gesellschafterdarlehen.Id)
	b1 := a.Bookings[0]
	assert.Equal(t, -77777.0, b1.Amount)
	assert.Equal(t, "GV-Vorjahr", b1.Type)
	assert.Equal(t, "JM", b1.CostCenter)

	// Buchung wurde aufs Bankkonto gebucht
	habenAccount,_ := accSystem.Get(accountSystem.SKR03_1200.Id)
	assert.Equal(t, 1, len(habenAccount.Bookings))
	actual := habenAccount.Bookings[0]
	assert.Equal(t, 77777.0, actual.Amount)
	assert.Equal(t, "Rest Anteil Johannes", actual.Text)
	assert.Equal(t, "GV-Vorjahr", actual.Type)
}

// test whether there is a not yet payed invoice
func TestProcessOPOS_SKR1600(t *testing.T) {
	setUp()

	// given: a internal hours booking
	bkDate,_ := time.Parse("2006 01 02 15 04 05",  "2017 11 11 11 11 11"  )
	tomorrow := bkDate.AddDate(+1, 0, +1)
	p := booking.NewBooking(13,"ER", "", "", "K", nil, 8250.0, "Interne Stunden 2017", 11, 2017, tomorrow)

	// when: the position is processed
	Process(accSystem, *p)
	ErloesverteilungAnValueMagnets(accSystem)

	// the booking is booked to SRK1600 account
	account1600, _ := accSystem.Get(accountSystem.SKR03_1600.Id)
	bookings1600 := account1600.Bookings
	assert.Equal(t, 1, len(bookings1600))

	// the booking is booked to partners via costCenter booking
	accountK, _ := accSystem.Get(valueMagnets.StakeholderKM.Id)
	bookingsK := accountK.Bookings
	assert.Equal(t, 1, len(bookingsK))

}

// Teste TestBonusRückstellungAngestellterSKR03
func TestBonusRueckstellungAngestellterSKR03(t *testing.T) {
	accSystem = accountSystem.NewDefaultAccountSystem()

	// given: a internal hours booking
	now := time.Now().AddDate(0, 0, 0)
	p := booking.NewBooking(13,"SKR03", "4120", "965", "BW", nil, 1337.42, "CC_Gehalt Januar 2017", 12, 2017, now)

	// when: the position is processed
	Process(accSystem, *p)

	// soll account
	a, _ := accSystem.Get(accountSystem.SKR03_4100_4199.Id)
	assert.Equal(t, 1, len(a.Bookings))
	assert.Equal(t, -1337.42, a.Bookings[0].Amount)
	assert.Equal(t, booking.SKR03, a.Bookings[0].Type)
	assert.Equal(t, "BW", a.Bookings[0].CostCenter)

	// booking is on Rückstellungsaccount
	rueckstellungen, _ := accSystem.Get(accountSystem.SKR03_Rueckstellungen.Id)
	assert.Equal(t, 1, len(rueckstellungen.Bookings))
	assert.Equal(t, 1337.42, rueckstellungen.Bookings[0].Amount)
	assert.Equal(t, booking.SKR03, rueckstellungen.Bookings[0].Type)
}

// Test Abschreibungen auf Anlagen
func TestAbschreibungenAufAnlagen(t *testing.T) {
	accSystem = accountSystem.NewDefaultAccountSystem()

	// given: Abschreibung
	now := time.Now().AddDate(0, 0, 0)
	p := booking.NewBooking(13,"SKR03", "4830", "25", "", nil, 1337.23, "Abschreibung Sachanlage", 12, 2017, now)

	// when: the position is processed
	Process(accSystem, *p)

	// soll account
	a, _ := accSystem.Get(accountSystem.SKR03_Abschreibungen.Id)
	assert.Equal(t, 1, len(a.Bookings))
	assert.Equal(t, -1337.23 , a.Bookings[0].Amount )
	assert.Equal(t, booking.SKR03, a.Bookings[0].Type)

	// booking is not on bankaccount
	ba := accSystem.GetCollectiveAccount()
	assert.Equal(t, 0, len(ba.Bookings))

	// booking is posiv von haben account
	ha, _ := accSystem.Get(accountSystem.SKR03_Anlagen25.Id)
	assert.Equal(t, 1, len(ha.Bookings))

}

// TestUstVZ
func TestUstVZ(t *testing.T) {
	accSystem = accountSystem.NewDefaultAccountSystem()

	// given: Abschreibung
	now,_ := time.Parse("2006 01 02 15 04 05",  "2017 11 11 11 11 11"  )
	p := booking.NewBooking(13,"UstVZ", "", "", "",nil, 1337.23, "UST", 12, 2017, now)

	// when: the position is processed
	Process(accSystem, *p)

	// soll account
	a, _ := accSystem.Get(accountSystem.SKR03_Umsatzsteuer.Id)
	assert.Equal(t, 1, len(a.Bookings))
	assert.Equal(t, -1337.23 , a.Bookings[0].Amount )
	assert.Equal(t, booking.UstVZ, a.Bookings[0].Type)

	// booking is  on bankaccount
	habenAccount,_ := accSystem.Get(accountSystem.SKR03_1200.Id)
	assert.Equal(t, 1, len(habenAccount.Bookings))
	assert.Equal(t, 1337.23 , habenAccount.Bookings[0].Amount )

}

func TestErloesverteilungAnValueMagnetsSimple(t *testing.T) {
	as := accountSystem.NewDefaultAccountSystem()

	// given: BOOKING ER
	// Eingangsrechnung 12852.0€ von Bank an SKR03_sonstigeAufwendungen
	its2018 := time.Date(2018, 1, 23, 0, 0, 0, 0, time.UTC)
	net := make(map[valueMagnets.Stakeholder]float64)
	net[valueMagnets.StakeholderRepository{}.Get("BW")] = 1000.0
	p3 := booking.NewBooking(13, "AR", "", "", "BW", net,  1190, "ARGSSLL", 1, 2017, its2018)

	// when: the position is processed
	Process(as, *p3)
	ErloesverteilungAnValueMagnets(as)

	// booking of 25% K-share is on CostCenter K
	b,_ := as.Get("K")
	b.UpdateSaldo()
	assert.Equal(t, 1, len(b.Bookings))
	assert.Equal(t, 0.0, b.Advances)
	assert.Equal(t, 250.0, b.Saldo)

}


func TestErloesverteilungAnValueMagnets(t *testing.T) {
	as := accountSystem.NewDefaultAccountSystem()

	// given: BOOKING ER
	// Eingangsrechnung 12852.0€ von Bank an SKR03_sonstigeAufwendungen
	its2018 := time.Date(2018, 1, 23, 0, 0, 0, 0, time.UTC)
	p1 := booking.NewBooking(13, "ER", "", "", "K", nil, 119, "hugo 1234", 1, 2017, its2018)
	p2 := booking.NewBooking(13, "ER", "", "", "BW", nil, 11900, "gugo. blupp", 1, 2017, its2018)
	net := make(map[valueMagnets.Stakeholder]float64)
	net[valueMagnets.StakeholderRepository{}.Get("BW")] = 1000.0
	p3 := booking.NewBooking(13, "AR", "", "", "BW", net,  1190, "ARGSSLL", 1, 2017, its2018)
	p4 := booking.NewBooking(13, "GV", "", "", "JM", nil, 5000, "ARGSSLL", 1, 2017, its2018)
	p5 := booking.NewBooking(13, "SKR03", "965", "4957", "K", nil, 42, "SKR03test", 1, 2017, its2018)

	// when: the position is processed
	Process(as, *p1)
	Process(as, *p2)
	Process(as, *p3)
	Process(as, *p4)
	Process(as, *p5)
	ErloesverteilungAnValueMagnets(as)


	// booking ist on CostCenter K
	b,_ := as.Get("K")
	b.UpdateSaldo()
	assert.Equal(t, 3, len(b.Bookings))
	assert.Equal(t, 0.0, b.Advances)
	assert.Equal(t, 192.0, b.Saldo)

	// Booking is on CostCenter BW
	a, _ := as.Get("BW")
	a.UpdateSaldo()
	assert.Equal(t, 3, len(a.Bookings))
	assert.Equal(t, booking.CC_Employeeaanteil, a.Bookings[1].Type)
	assert.Equal(t, 0.0, a.Internals)
	assert.Equal(t, 0.0, a.Advances)
	assert.Equal(t, -9250.0, a.Saldo)

	// Booking is on CostCenter JM
	c, _ := as.Get("JM")
	c.UpdateSaldo()
	assert.Equal(t, 1, len(c.Bookings))
	assert.Equal(t, booking.CC_Entnahme, c.Bookings[0].Type)
	assert.Equal(t, 0.0, c.Revenue)
	assert.Equal(t, 0.0, c.Internals)
	assert.Equal(t, -5000.0, c.Advances)
	assert.Equal(t, -5000.0, c.Saldo)


}


func assertBooking(t *testing.T, b booking.Booking, amount float64, text string, destType string) {
	util.AssertFloatEquals(t, amount, b.Amount)
	util.AssertEquals(t, text, b.Text)
	util.AssertEquals(t, destType, b.Type)
}
