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
	p := booking.NewBooking("AR", "JM", net, 12852.0, "Rechnung 1234", 1, 2017, time.Time{}, time.Time{})

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
	setUp()

	p := booking.NewBooking("Gehalt", "BW", nil, 3869.65, "Gehalt Ben", 1, 2017, time.Time{}, time.Time{})

	Process(repository, *p)

	// 100% Brutto gegen Bankkonto
	assert.Equal(t, -3869.65, accountBank.Bookings[0].Amount)
	assert.Equal(t, "Gehalt Ben", accountBank.Bookings[0].Text)
	assert.Equal(t, "Gehalt", accountBank.Bookings[0].Type)

	// 100% Brutto gegen Kommitment
	assert.Equal(t, -3869.65, accountKommitment.Bookings[0].Amount)
	assert.Equal(t, booking.Gehalt, accountKommitment.Bookings[0].Type)

	// Kommitment-Buchung ist der Kostenstelle "BW" zugeordnet
	assert.Equal(t, "BW", accountKommitment.Bookings[0].CostCenter)
}

func TestExternNettoAnteil(t *testing.T) {
	setUp()

	// given: a booking
	net := map[owner.Stakeholder]float64{
		owner.StakeholderEX: 10800.0,
	}
	p := booking.NewBooking("AR", "JM", net, 12852.0, "Rechnung 1234", 1, 2017, time.Time{}, time.Time{})

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

	p := booking.NewBooking("ER", "K", nil, 12852.0, "Eingangsrechnung 1234", 1, 2017, time.Time{}, time.Time{})

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

func TestPartnerEntnahme(t *testing.T) {
	setUp()

	extras := booking.CsvBookingExtras{Typ: "GV", Responsible: "RW"}
	extras.Net = make(map[owner.Stakeholder]float64)
	b := booking.NewBooking("GV", "RW", nil, 6000, "", 1, 2017, time.Time{}, time.Time{})

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

// Interne Stunden
// - werden nicht auf das Bankkonto gebucht
// - 100% werden auf das Partner-Konto gebucht
// - 100% werden gegen das Kommitment-Konto gebucht
func TestInterneStunden(t *testing.T) {
	setUp()

	// given: a internal hours booking
	p := booking.NewBooking("IS", "AN", nil, 8250.0, "Interne Stunden 2017", 12, 2017, time.Time{}, time.Time{})

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
	b := booking.NewBooking("AR", "K", nil, 6000, "Ausgangsrechnung", 1, 2017, time.Time{}, time.Time{})

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
	b := booking.NewBooking("SV-Beitrag", "BW", nil, 1385.10, "KKH, Ben", 5, 2017, time.Time{}, time.Time{})

	Process(repository, *b)

	// Buchung wurde gegen Kommitment-Konto gebucht
	a, _ := repository.Get(owner.StakeholderKM.Id)
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
	b := booking.NewBooking("LNSteuer", "BW", nil, 1511.45, "Lohnsteuer Ben", 5, 2017, time.Time{}, time.Time{})

	Process(repository, *b)

	// Buchung wurde gegen Kommitment-Konto gebucht
	assertBooking(t, accountKommitment.Bookings[0], -1511.45, "Lohnsteuer Ben", "LNSteuer")
	assert.Equal(t, "BW", accountKommitment.Bookings[0].CostCenter)

	// Buchung wurde aufs Bankkonto gebucht
	assertBooking(t, accountBank.Bookings[0], -1511.45, "Lohnsteuer Ben", "LNSteuer")
}

// 100% werden auf das Bankkonto gebucht
// 100% werden gegen das Kommitment-Konto gebucht. Diese Regel ist nich unscharf:
// eigentlich müssen die 100% aufgeteilt werden auf: 70% auf Partner, 25% auf
// Kommitment und 5% auf Dealbringer
func TestProcessGWSteuer(t *testing.T) {
	setUp()

	b := booking.NewBooking("GWSteuer", "K", nil, 2385.10, "STEUERKASSE HAMBURG STEUERNR 048/638/01147 GEW.ST 4VJ.17", 9, 2017, time.Time{}, time.Time{})

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

func assertBooking(t *testing.T, b booking.Booking, amount float64, text string, destType string) {
	util.AssertFloatEquals(t, amount, b.Amount)
	util.AssertEquals(t, text, b.Text)
	util.AssertEquals(t, destType, b.Type)
}
