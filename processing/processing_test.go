package processing

import (
	"testing"
	"time"

	"github.com/ahojsenn/kontrol/account"

	"github.com/ahojsenn/kontrol/booking"
	"github.com/ahojsenn/kontrol/owner"
	"github.com/ahojsenn/kontrol/util"
	"github.com/stretchr/testify/assert"
)

var repository account.Repository
var accountBank *account.Account
var accountHannes *account.Account
var accountRalf *account.Account
var accountKommitment *account.Account

func setUp() {
	repository = account.NewDefaultRepository()
	accountBank = repository.BankAccount()
	accountHannes, _ = repository.Get(owner.StakeholderJM.Id)
	accountRalf, _ = repository.Get(owner.StakeholderRW.Id)
	accountKommitment, _ = repository.Get(owner.StakeholderKM.Id)
}

// Ausgangsrechnung Angestellter
// - 5% Provision für Dealbringer
// - 95% für Kommitmentment, Kostenstelle Angestellter
func TestAusgangsrechnungAngestellter(t *testing.T) {
	setUp()

	// Ben hat auf einer Buchung nett 10.800 EUR erwirtschaftet
	net := map[owner.Stakeholder]float64{
		owner.StakeholderBW: 10800.0,
	}
	its2018 := time.Date(2018, 1, 23, 0, 0, 0, 0, time.UTC)
	p := booking.NewBooking("AR", "800", "1337", "JM", net, 12852.0, "Rechnung 1234", 1, 2017, time.Time{}, its2018)

	Process(repository, *p)

	// Johannes kriegt 5% Provision
	provision := accountHannes.Bookings[0]
	util.AssertFloatEquals(t, 10800.0*owner.PartnerProvision, provision.Amount)
	util.AssertEquals(t, booking.Vertriebsprovision, provision.Type)

	// Kommitment kriegt 95% der Nettorechnung
	util.AssertEquals(t, 1, len(accountHannes.Bookings))
	kommitment := accountKommitment.Bookings[0]
	util.AssertFloatEquals(t, 10800.0*owner.KommmitmentEmployeeShare, kommitment.Amount)
	util.AssertEquals(t, booking.Kommitmentanteil, kommitment.Type)

	// Kommitment-Buchung ist der Kostenstelle "BW" zugeordnet
	assert.Equal(t, "BW", kommitment.CostCenter)
}

// Gehalt Angestellter
// - 100% Brutto gegen Bankkonto
// - 100% Brutto gegen Kommitmentkonto
// - Kostenstelle: Kürzel des Angestellten
func TestGehaltAngestellter(t *testing.T) {
	repository := account.NewDefaultRepository()

	// given: a booking
	its2018 := time.Date(2018, 1, 23, 0, 0, 0, 0, time.UTC)
	p := booking.NewBooking("Gehalt", "800", "1337", "BW", nil, 3869.65, "Gehalt Ben", 1, 2017, time.Time{}, its2018)

	// when: the position is processed
	Process(repository, *p)

	// 100% Brutto gegen Bankkonto
	accountBank := repository.BankAccount()
	assert.Equal(t, -3869.65, accountBank.Bookings[0].Amount)
	assert.Equal(t, "Gehalt Ben", accountBank.Bookings[0].Text)
	assert.Equal(t, "Gehalt", accountBank.Bookings[0].Type)

	// 100% Brutto gegen SKR03_4100_4199
	account2, _ := repository.Get(owner.SKR03_4100_4199.Id)
	util.AssertEquals(t, 1, len(account2.Bookings))
	assert.Equal(t, -3869.65, account2.Bookings[0].Amount)
	assert.Equal(t, booking.Gehalt, account2.Bookings[0].Type)

	// Kommitment-Buchung ist der Kostenstelle "BW" zugeordnet
	assert.Equal(t, "BW", account2.Bookings[0].CostCenter)
}

func TestExternNettoAnteil(t *testing.T) {
	setUp()

	// given: a booking
	net := map[owner.Stakeholder]float64{
		owner.StakeholderEX: 10800.0,
	}
	its2018 := time.Date(2018, 1, 23, 0, 0, 0, 0, time.UTC)
	p := booking.NewBooking("AR", "800", "1337", "JM", net, 12852.0, "Rechnung 1234", 1, 2017, time.Time{}, its2018)

	// when: the position is processed
	Process(repository, *p)

	// and hannes got his provision
	accountHannes, _ := repository.Get(owner.StakeholderJM.Id)
	provision := accountHannes.Bookings[0]
	util.AssertFloatEquals(t, 10800.0*owner.PartnerProvision, provision.Amount)
	util.AssertEquals(t, booking.Vertriebsprovision, provision.Type)

	// and kommitment got 95%
	util.AssertEquals(t, 1, len(accountHannes.Bookings))
	accountKommitment, _ := repository.Get(owner.StakeholderKM.Id)
	kommitment := accountKommitment.Bookings[0]
	util.AssertFloatEquals(t, 10800.0*owner.KommmitmentExternShare, kommitment.Amount)
	util.AssertEquals(t, booking.Kommitmentanteil, kommitment.Type)
}

// Eingangsrechnungen
// - 100% werden netto gegen das Bankkonto gebucht
// - 100% des Nettobetrags werden gegen das Kommitment-Konto gebucht
func TestEingangsrechnung(t *testing.T) {
	setUp()

	its2018 := time.Date(2018, 1, 23, 0, 0, 0, 0, time.UTC)
	p := booking.NewBooking("ER", "800", "1337", "K", nil, 12852.0, "Eingangsrechnung 1234", 1, 2017, time.Time{}, its2018)

	Process(repository, *p)

	// Buchung wurde gegen das Kommitment-Konto gebucht
	accountKommitment, _ := repository.Get(owner.StakeholderKM.Id)
	assert.Equal(t, 1, len(accountKommitment.Bookings))
	kommitment := accountKommitment.Bookings[0]
	assert.Equal(t, util.Net(-12852.0), kommitment.Amount)
	assert.Equal(t, kommitment.Type, booking.Eingangsrechnung)
	assert.Equal(t, "K", kommitment.CostCenter)

	// Buchung wurde gegen das Bankkonto gebucht
	assert.Equal(t, 1, len(repository.BankAccount().Bookings))
	actual := repository.BankAccount().Bookings[0]
	assert.Equal(t, util.Net(-12852.0), actual.Amount)
	assert.Equal(t, "Eingangsrechnung 1234", actual.Text)
	assert.Equal(t, "ER", actual.Type)
}

func TestEingangsrechnungGegenRückstellung(t *testing.T) {
	setUp()
	// Eingangserechnungen können auch gegen Rückstellungen gebucht werden

	// given a Buchung Eingangsrechnung gegen Rücksttellung
	its2018 := time.Date(2018, 1, 23, 0, 0, 0, 0, time.UTC)
	p := booking.NewBooking("ERgegenRückstellung", "800", "1337", "Rückstellung", nil, 12852.0, "ER 1234", 1, 2017, time.Time{}, its2018)

	// when: the position is processed
	Process(repository, *p)

	// the booking is booked from Rückstellung account
	a1, _ := repository.Get(owner.StakeholderRueckstellung.Id)
	util.AssertEquals(t, 1, len(a1.Bookings))
	b1 := a1.Bookings[0]
	util.AssertFloatEquals(t, util.Net(-12852.0), b1.Amount)
	util.AssertEquals(t, booking.Eingangsrechnung, b1.Type)

	// the booking is not booked to K Accout
	c1, _ := repository.Get(owner.StakeholderKM.Id)
	util.AssertEquals(t, 0, len(c1.Bookings))

	// the booking is  booked on bank account
	assert.Equal(t, 1, len(repository.BankAccount().Bookings))
	actual := repository.BankAccount().Bookings[0]
	assert.Equal(t, util.Net(-12852.0), actual.Amount)
	assert.Equal(t, "ER 1234", actual.Text)
	assert.Equal(t, "ERgegenRückstellung", actual.Type)

}

func TestRückstellungAuflösen(t *testing.T) {
	setUp()
	// Rückstellungen können gegen das kommitment Konto aufgelöst werden

	// given a Buchung Eingangsrechnung gegen Rücksttellung
	its2018 := time.Date(2018, 1, 23, 0, 0, 0, 0, time.UTC)
	p := booking.NewBooking("RückstellungAuflösen", "800", "1337", "K", nil, -12852.0, "Auflösung Rückstellungsdifferenz", 1, 2017, time.Time{}, its2018)

	// when: the position is processed
	Process(repository, *p)

	// the booking is booked from Rückstellung account
	a1, _ := repository.Get(owner.StakeholderRueckstellung.Id)
	util.AssertEquals(t, 1, len(a1.Bookings))
	b1 := a1.Bookings[0]
	util.AssertFloatEquals(t, -12852.0, b1.Amount)
	util.AssertEquals(t, booking.Eingangsrechnung, b1.Type)

	// the booking is not booked to the bankaccout
	util.AssertEquals(t, 0, len(repository.BankAccount().Bookings))

	// the booking is  booked on kommitment account
	accountKommitment, _ := repository.Get(owner.StakeholderKM.Id)
	assert.Equal(t, 1, len(accountKommitment.Bookings))
	kommitment := accountKommitment.Bookings[0]
	assert.Equal(t, 12852.0, kommitment.Amount)
	assert.Equal(t, kommitment.Type, booking.RueckstellungAuflösen)
	assert.Equal(t, "K", kommitment.CostCenter)
}

func TestPartnerEntnahme(t *testing.T) {
	setUp()

	extras := booking.CsvBookingExtras{Typ: "GV", Responsible: "RW"}
	extras.Net = make(map[owner.Stakeholder]float64)
	its2018 := time.Date(2018, 1, 23, 0, 0, 0, 0, time.UTC)
	b := booking.NewBooking("GV", "800", "1337", "RW", nil, 6000, "", 1, 2017, time.Time{}, its2018)

	Process(repository, *b)

	bRalf := accountRalf.Bookings[0]
	util.AssertFloatEquals(t, -6000, bRalf.Amount)
	util.AssertEquals(t, booking.Entnahme, bRalf.Type)

	// Buchung wurde gegen das Bankkonto gebucht
	util.AssertEquals(t, 1, len(repository.BankAccount().Bookings))
	actual := repository.BankAccount().Bookings[0]
	util.AssertFloatEquals(t, -6000, actual.Amount)
	util.AssertEquals(t, "GV", actual.Type)
}

// Rückstellungen
// - werden nicht auf das Bankkonto gebucht
// - 100% werden auf das Rückstellung-Konto gebucht
// - 100% werden gegen das Kommitment-Konto gebucht
func TestRückstellung(t *testing.T) {
	setUp()

	// given: a Rückstellung booking
	its2018 := time.Date(2018, 1, 23, 0, 0, 0, 0, time.UTC)
	p := booking.NewBooking("Rückstellung", "800", "1337", "BW", nil, 4711.0, "Bonus Rückstellung", 12, 2017, time.Time{}, its2018)

	// when: the position is processed
	Process(repository, *p)

	// the booking is booked to Rückstellung account
	a1, _ := repository.Get(owner.StakeholderRueckstellung.Id)
	util.AssertEquals(t, 1, len(a1.Bookings))
	b1 := a1.Bookings[0]
	util.AssertFloatEquals(t, 4711.00, b1.Amount)
	util.AssertEquals(t, booking.Rueckstellung, b1.Type)

	// the booking is booked against kommitment account
	a2, _ := repository.Get(owner.StakeholderKM.Id)
	b2 := a2.Bookings[0]
	util.AssertFloatEquals(t, -4711.00, b2.Amount)
	util.AssertEquals(t, booking.Rueckstellung, b1.Type)

	// Rückstellungen are not booked on bank account
	util.AssertEquals(t, 0, len(repository.BankAccount().Bookings))
}

// Interne Stunden
// - werden nicht auf das Bankkonto gebucht
// - 100% werden auf das Rückstellung-Konto gebucht
// - 100% werden gegen das Kommitment-Konto gebucht
func TestInterneStunden(t *testing.T) {
	setUp()

	// given: a internal hours booking
	its2018 := time.Date(2018, 1, 23, 0, 0, 0, 0, time.UTC)
	p := booking.NewBooking("IS", "800", "1337", "AN", nil, 8250.0, "Interne Stunden 2017", 12, 2017, time.Time{}, its2018)

	// when: the position is processed
	Process(repository, *p)

	// the booking is booked to anke's account
	a1, _ := repository.Get(owner.StakeholderAN.Id)
	util.AssertEquals(t, 1, len(a1.Bookings))
	b1 := a1.Bookings[0]
	util.AssertFloatEquals(t, 8250.00, b1.Amount)
	util.AssertEquals(t, booking.InterneStunden, b1.Type)

	// the booking is booked against kommitment account
	a2, _ := repository.Get(owner.StakeholderKM.Id)
	b2 := a2.Bookings[0]
	util.AssertFloatEquals(t, -8250.00, b2.Amount)
	util.AssertEquals(t, booking.InterneStunden, b1.Type)

	// internal hours are not booked on bank account
	util.AssertEquals(t, 0, len(repository.BankAccount().Bookings))
}

func TestBookAusgangsrechnungToBankAccount(t *testing.T) {
	setUp()
	its2018 := time.Date(2018, 1, 23, 0, 0, 0, 0, time.UTC)
	b := booking.NewBooking("AR", "800", "1337", "K", nil, 6000, "Ausgangsrechnung", 1, 2017, time.Time{}, its2018)

	Process(repository, *b)

	util.AssertEquals(t, 1, len(repository.BankAccount().Bookings))
	actual := repository.BankAccount().Bookings[0]
	util.AssertFloatEquals(t, util.Net(6000), actual.Amount)
	util.AssertEquals(t, "Ausgangsrechnung", actual.Text)
	util.AssertEquals(t, "AR", actual.Type)
}

// 100% werden auf das Bankkonto gebucht
// 100% werden gegen das Kommitment-Konto gebucht
func TestProcessSVBeitrag(t *testing.T) {
	setUp()
	its2018 := time.Date(2018, 1, 23, 0, 0, 0, 0, time.UTC)
	b := booking.NewBooking("SV-Beitrag", "800", "1337", "BW", nil, 1385.10, "KKH, Ben", 5, 2017, time.Time{}, its2018)

	Process(repository, *b)

	// Buchung wurde gegen Kommitment-Konto gebucht
	a, _ := repository.Get(owner.SKR03_4100_4199.Id)
	b1 := a.Bookings[0]
	assert.Equal(t, -1385.10, b1.Amount)
	assert.Equal(t, booking.SVBeitrag, b1.Type)
	assert.Equal(t, "BW", b1.CostCenter)

	// Buchung wurde aufs Bankkonto gebucht
	assert.Equal(t, 1, len(repository.BankAccount().Bookings))
	actual := repository.BankAccount().Bookings[0]
	assert.Equal(t, -1385.10, actual.Amount)
	assert.Equal(t, "KKH, Ben", actual.Text)
	assert.Equal(t, "SV-Beitrag", actual.Type)
}

// 100% werden auf das Bankkonto gebucht
// 100% werden gegen das Kommitment-Konto gebucht
// Kostenstelle: Angestellter, für den Lohnsteuer gezahlt wurde
func TestProcessLNSteuer(t *testing.T) {
	setUp()
	its2018 := time.Date(2018, 1, 23, 0, 0, 0, 0, time.UTC)
	b := booking.NewBooking("LNSteuer", "800", "1337", "BW", nil, 1511.45, "Lohnsteuer Ben", 5, 2017, time.Time{}, its2018)

	Process(repository, *b)

	// Buchung wurde gegen Kommitment-Konto gebucht
	account2, _ := repository.Get(owner.SKR03_4100_4199.Id)
	assertBooking(t, account2.Bookings[0], -1511.45, "Lohnsteuer Ben", "LNSteuer")
	assert.Equal(t, "BW", account2.Bookings[0].CostCenter)

	// Buchung wurde aufs Bankkonto gebucht
	assertBooking(t, accountBank.Bookings[0], -1511.45, "Lohnsteuer Ben", "LNSteuer")
}

// 100% werden auf das Bankkonto gebucht
// 100% werden gegen das Kommitment-Konto gebucht. Diese Regel ist nicht unscharf:
// eigentlich müssen die 100% aufgeteilt werden auf: 70% auf Partner, 25% auf
// Kommitment und 5% auf Dealbringer
func TestProcessGWSteuer(t *testing.T) {
	setUp()

	its2018 := time.Date(2018, 1, 23, 0, 0, 0, 0, time.UTC)
	b := booking.NewBooking("GWSteuer", "800", "1337", "K", nil, 2385.10, "STEUERKASSE HAMBURG STEUERNR 048/638/01147 GEW.ST 4VJ.17", 9, 2017, time.Time{}, its2018)

	Process(repository, *b)

	// Buchung wurde gegen Kommitment-Konto gebucht
	a, _ := repository.Get(owner.StakeholderKM.Id)
	b1 := a.Bookings[0]
	assertBooking(t, b1, -2385.10, "STEUERKASSE HAMBURG STEUERNR 048/638/01147 GEW.ST 4VJ.17", booking.GWSteuer)

	// Buchung wurde aufs Bankkonto gebucht
	util.AssertEquals(t, 1, len(repository.BankAccount().Bookings))
	actual := repository.BankAccount().Bookings[0]
	assertBooking(t, actual, -2385.10, "STEUERKASSE HAMBURG STEUERNR 048/638/01147 GEW.ST 4VJ.17", "GWSteuer")
}

// 100% werden auf das Bankkonto gebucht
// 100% werden gegen das Rückstellung gebucht.
func TestProcessGWSteuer_gegenRückstellung(t *testing.T) {
	setUp()

	its2018 := time.Date(2018, 1, 23, 0, 0, 0, 0, time.UTC)
	b := booking.NewBooking("GWSteuer", "800", "1337", "Rückstellung", nil, 2385.10, "STEUERKASSE HAMBURG STEUERNR 048/638/01147 GEW.ST 4VJ.17", 9, 2017, time.Time{}, its2018)

	Process(repository, *b)

	// Buchung wurde gegen Kommitment-Konto gebucht
	a, _ := repository.Get(owner.StakeholderRueckstellung.Id)
	b1 := a.Bookings[0]
	assertBooking(t, b1, -2385.10, "STEUERKASSE HAMBURG STEUERNR 048/638/01147 GEW.ST 4VJ.17", booking.GWSteuer)

	// Buchung wurde aufs Bankkonto gebucht
	util.AssertEquals(t, 1, len(repository.BankAccount().Bookings))
	actual := repository.BankAccount().Bookings[0]
	assertBooking(t, actual, -2385.10, "STEUERKASSE HAMBURG STEUERNR 048/638/01147 GEW.ST 4VJ.17", "GWSteuer")
}

func assertBooking(t *testing.T, b booking.Booking, amount float64, text string, destType string) {
	util.AssertFloatEquals(t, amount, b.Amount)
	util.AssertEquals(t, text, b.Text)
	util.AssertEquals(t, destType, b.Type)
}

// 100% werden als Anfangsbestand auf ein Konto gebucht, bspw. Rückstellung
func TestProcessAnfangsbestand(t *testing.T) {
	setUp()

	// given a booking with Anfangsbestand
	b := booking.NewBooking("Anfangsbestand", "800", "1337", "Rückstellung", nil, 42.23, "Anfangsbestand aus Vorjahr", 9, 2017, time.Time{}, time.Time{})

	// when: the position is processed
	Process(repository, *b)

	// the booking is booked to Rückstellung account
	a1, _ := repository.Get(owner.StakeholderRueckstellung.Id)
	util.AssertEquals(t, 1, len(a1.Bookings))
	b1 := a1.Bookings[0]
	util.AssertFloatEquals(t, 42.23, b1.Amount)
	util.AssertEquals(t, booking.Anfangsbestand, b1.Type)

	// Anfangsbestand is not booked on bank account
	util.AssertEquals(t, 0, len(repository.BankAccount().Bookings))
}

// 100% werden als Anfangsbestand auf ein Konto gebucht, bspw. Rückstellung
func TestProcessAnfangsbestand_JahresüberschusssVJ(t *testing.T) {
	setUp()

	// given a booking with Anfangsbestand
	b := booking.NewBooking("Anfangsbestand", "800", "1337", "JahresüberschussVJ", nil, 10042.23, "Anfangsbestand aus Vorjahr", 9, 2017, time.Time{}, time.Time{})

	// when: the position is processed
	Process(repository, *b)

	// the booking is booked to Rückstellung account
	a1, _ := repository.Get(owner.KontoJUSVJ.Id)
	util.AssertEquals(t, 1, len(a1.Bookings))
	b1 := a1.Bookings[0]
	util.AssertFloatEquals(t, 10042.23, b1.Amount)
	util.AssertEquals(t, booking.Anfangsbestand, b1.Type)

	// Anfangsbestand is not booked on bank account
	util.AssertEquals(t, 0, len(repository.BankAccount().Bookings))
}

// 100% werden auf das Bankkonto gebucht
// 100% werden gegen das JahresüberschussVJ gebucht
func TestProcessGV_Vorjahr(t *testing.T) {
	setUp()
	b := booking.NewBooking("GV-Vorjahr", "800", "1337", "JM", nil, 77777, "Rest Anteil Johannes", 5, 2017, time.Time{}, time.Time{})

	Process(repository, *b)

	// Buchung wurde gegen JahresüberschussVJ gebucht
	a, _ := repository.Get(owner.KontoJUSVJ.Id)
	b1 := a.Bookings[0]
	assert.Equal(t, -77777.0, b1.Amount)
	assert.Equal(t, booking.GVVorjahr, b1.Type)
	assert.Equal(t, "JM", b1.CostCenter)

	// Buchung wurde aufs Bankkonto gebucht
	assert.Equal(t, 1, len(repository.BankAccount().Bookings))
	actual := repository.BankAccount().Bookings[0]
	assert.Equal(t, -77777.0, actual.Amount)
	assert.Equal(t, "Rest Anteil Johannes", actual.Text)
	assert.Equal(t, "GV-Vorjahr", actual.Type)
}

// test whether there is a not yet payed invoice
func TestProcessOPOS_SKR1600(t *testing.T) {
	setUp()

	// given: a internal hours booking
	tomorrow := time.Now().AddDate(0, 0, +1)
	p := booking.NewBooking("ER", "800", "1337", "K", nil, 8250.0, "Interne Stunden 2017", 12, 2017, tomorrow, tomorrow)

	// when: the position is processed
	Process(repository, *p)

	// the booking is booked to SRK1600 account
	// then the booking is on SKR03_1400
	account1600, _ := repository.Get(owner.SKR03_1600.Id)
	bookings1600 := account1600.Bookings
	assert.Equal(t, 1, len(bookings1600))

	// the booking is not yet booked to partners
	accountK, _ := repository.Get(owner.StakeholderKM.Id)
	bookingsK := accountK.Bookings
	assert.Equal(t, 0, len(bookingsK))

}
