package processing

import (
	"errors"
	"testing"

	"bitbucket.org/rwirdemann/kontrol/account"

	"bitbucket.org/rwirdemann/kontrol/owner"
	"bitbucket.org/rwirdemann/kontrol/util"
)

var repository account.Repository

func setUp() {
	repository = account.NewDefaultRepository()
}

func TestPartnerNettoAnteil(t *testing.T) {
	setUp()

	// given: a booking
	net := make(map[owner.Stakeholder]float64)
	net[owner.StakeholderRW] = 10800.0
	net[owner.StakeholderJM] = 3675.0
	p := account.NewBooking("AR", "JM", net, 17225.25, "Rechnung 1234", 1, 2017)

	// when: the position is processed
	Process(repository, *p)

	// then ralf 1 booking: his own net share
	accountRalf, _ := repository.Get(owner.StakeholderRW.Id)
	bookingsRalf := accountRalf.Bookings
	util.AssertEquals(t, 1, len(bookingsRalf))
	bRalf, _ := findBookingByText(bookingsRalf, "Rechnung 1234#NetShare#RW")
	util.AssertFloatEquals(t, 10800.0*owner.PartnerShare, bRalf.Amount)
	util.AssertEquals(t, 1, bRalf.Month)
	util.AssertEquals(t, 2017, bRalf.Year)
	util.AssertEquals(t, account.Nettoanteil, bRalf.DestType)

	// and hannes got 3 bookings: his own net share and 2 provisions
	accountHannes, _ := repository.Get(owner.StakeholderJM.Id)
	bookingsHannes := accountHannes.Bookings
	util.AssertEquals(t, 3, len(bookingsHannes))

	// net share
	b, _ := findBookingByText(bookingsHannes, "Rechnung 1234#NetShare#JM")
	util.AssertFloatEquals(t, 3675.0*owner.PartnerShare, b.Amount)
	util.AssertEquals(t, 1, b.Month)
	util.AssertEquals(t, 2017, b.Year)

	// provision from ralf
	provisionRalf, _ := findBookingByText(bookingsHannes, "Rechnung 1234#Provision#RW")
	util.AssertFloatEquals(t, 10800.0*owner.PartnerProvision, provisionRalf.Amount)
	util.AssertEquals(t, account.Vertriebsprovision, provisionRalf.DestType)

	// // provision from hannes
	provisionHannes, _ := findBookingByText(bookingsHannes, "Rechnung 1234#Provision#JM")
	util.AssertFloatEquals(t, 3675.0*owner.PartnerProvision, provisionHannes.Amount)
	util.AssertEquals(t, account.Vertriebsprovision, provisionHannes.DestType)

	// kommitment got 25% from ralfs net booking
	accountKommitment, _ := repository.Get(owner.StakeholderKM.Id)
	bookingsKommitment := accountKommitment.Bookings
	util.AssertEquals(t, 2, len(bookingsKommitment))
	kommitmentRalf, _ := findBookingByText(bookingsKommitment, "Rechnung 1234#Kommitment#RW")
	util.AssertFloatEquals(t, 10800.0*owner.KommmitmentShare, kommitmentRalf.Amount)
	util.AssertEquals(t, account.Kommitmentanteil, kommitmentRalf.DestType)

	// and kommitment got 25% from hannes net booking
	kommitmentHannes, _ := findBookingByText(bookingsKommitment, "Rechnung 1234#Kommitment#JM")
	util.AssertFloatEquals(t, 3675.0*owner.KommmitmentShare, kommitmentHannes.Amount)
	util.AssertEquals(t, account.Kommitmentanteil, kommitmentHannes.DestType)
}

func findBookingByText(bookings []account.Booking, text string) (*account.Booking, error) {
	for _, b := range bookings {
		if b.Text == text {
			return &b, nil
		}
	}
	return nil, errors.New("booking with test '" + text + " not found")
}

func TestExternAngestellterNettoAnteil(t *testing.T) {
	setUp()

	// given: a booking
	net := map[owner.Stakeholder]float64{
		owner.StakeholderBW: 10800.0,
	}
	p := account.NewBooking("AR", "JM", net, 12852.0, "Rechnung 1234", 1, 2017)

	// when: the position is processed
	Process(repository, *p)

	// and hannes got his provision
	accountHannes, _ := repository.Get(owner.StakeholderJM.Id)
	provision := accountHannes.Bookings[0]
	util.AssertFloatEquals(t, 10800.0*owner.PartnerProvision, provision.Amount)
	util.AssertEquals(t, account.Vertriebsprovision, provision.DestType)

	// and kommitment got 95%
	util.AssertEquals(t, 1, len(accountHannes.Bookings))
	accountKommitment, _ := repository.Get(owner.StakeholderKM.Id)
	kommitment := accountKommitment.Bookings[0]
	util.AssertFloatEquals(t, 10800.0*owner.KommmitmentEmployeeShare, kommitment.Amount)
	util.AssertEquals(t, account.Kommitmentanteil, kommitment.DestType)

	// 100% is booked to employee account to see how much money is made by this employee
	accountBen, _ := repository.Get(owner.StakeholderBW.Id)
	util.AssertEquals(t, 1, len(accountBen.Bookings))
	bookingBen := accountBen.Bookings[0]
	util.AssertFloatEquals(t, 10800.0, bookingBen.Amount)
}

func TestExternNettoAnteil(t *testing.T) {
	setUp()

	// given: a booking
	net := map[owner.Stakeholder]float64{
		owner.StakeholderEX: 10800.0,
	}
	p := account.NewBooking("AR", "JM", net, 12852.0, "Rechnung 1234", 1, 2017)

	// when: the position is processed
	Process(repository, *p)

	// and hannes got his provision
	accountHannes, _ := repository.Get(owner.StakeholderJM.Id)
	provision := accountHannes.Bookings[0]
	util.AssertFloatEquals(t, 10800.0*owner.PartnerProvision, provision.Amount)
	util.AssertEquals(t, account.Vertriebsprovision, provision.DestType)

	// and kommitment got 95%
	util.AssertEquals(t, 1, len(accountHannes.Bookings))
	accountKommitment, _ := repository.Get(owner.StakeholderKM.Id)
	kommitment := accountKommitment.Bookings[0]
	util.AssertFloatEquals(t, 10800.0*owner.KommmitmentExternShare, kommitment.Amount)
	util.AssertEquals(t, account.Kommitmentanteil, kommitment.DestType)
}

func TestEingangsrechnung(t *testing.T) {
	setUp()

	// given: a booking
	p := account.NewBooking("ER", "k", nil, 12852.0, "Eingangsrechnung 1234", 1, 2017)

	// when: the position is processed
	Process(repository, *p)

	// the invoice is booked to the kommitment account
	accountKommitment, _ := repository.Get(owner.StakeholderKM.Id)
	util.AssertEquals(t, 1, len(accountKommitment.Bookings))
	kommitment := accountKommitment.Bookings[0]
	util.AssertFloatEquals(t, util.Net(-12852.0), kommitment.Amount)
	util.AssertEquals(t, account.Eingangsrechnung, kommitment.DestType)
}

func TestPartnerWithdrawals(t *testing.T) {
	setUp()

	extras := account.CsvBookingExtras{SourceType: "GV", CostCenter: "RW"}
	extras.Net = make(map[owner.Stakeholder]float64)
	b := account.NewBooking("GV", "RW", nil, 6000, "", 1, 2017)
	Process(repository, *b)
	accountRalf, _ := repository.Get(owner.StakeholderRW.Id)
	util.AssertEquals(t, 1, len(accountRalf.Bookings))
	bRalf := accountRalf.Bookings[0]
	util.AssertFloatEquals(t, -6000, bRalf.Amount)
	util.AssertEquals(t, account.Entnahme, bRalf.DestType)
}



// Interne Stunden
// - werden nicht auf das Bankkonto gebucht
// - 100% werden auf das Partner-Konto gebucht
// - 100% werden gegen das Kommitment-Konto gebucht
func TestInterneStunden(t *testing.T) {
	setUp()

	// given: a internal hours booking
	p := account.NewBooking("IS", "AN", nil, 8250.0, "Interne Stunden 2017", 12, 2017)

	// when: the position is processed
	Process(repository, *p)

	// the booking is booked to anke's account
	a1, _ := repository.Get(owner.StakeholderAN.Id)
	util.AssertEquals(t, 1, len(a1.Bookings))
	b1 := a1.Bookings[0]
	util.AssertFloatEquals(t, 8250.00, b1.Amount)
	util.AssertEquals(t, account.InterneStunden, b1.DestType)

	// the booking is booked against kommitment account
	a2, _ := repository.Get(owner.StakeholderKM.Id)
	b2 := a2.Bookings[0]
	util.AssertFloatEquals(t, -8250.00, b2.Amount)
	util.AssertEquals(t, account.InterneStunden, b1.DestType)

	// internal hours are not booked on bank account
	util.AssertEquals(t, 0, len(repository.BankAccount().Bookings))
}

func TestBookEingangsrechnungToBankAccount(t *testing.T) {
	setUp()
	b := account.NewBooking("ER", "K", nil, 6000, "Eingangsrechnung", 1, 2017)

	Process(repository, *b)

	util.AssertEquals(t, 1, len(repository.BankAccount().Bookings))
	actual := repository.BankAccount().Bookings[0]
	util.AssertFloatEquals(t, util.Net(-6000), actual.Amount)
	util.AssertEquals(t, "Eingangsrechnung", actual.Text)
	util.AssertEquals(t, "ER", actual.DestType)
}

func TestBookAusgangsrechnungToBankAccount(t *testing.T) {
	setUp()
	b := account.NewBooking("AR", "K", nil, 6000, "Ausgangsrechnung", 1, 2017)

	Process(repository, *b)

	util.AssertEquals(t, 1, len(repository.BankAccount().Bookings))
	actual := repository.BankAccount().Bookings[0]
	util.AssertFloatEquals(t, util.Net(6000), actual.Amount)
	util.AssertEquals(t, "Ausgangsrechnung", actual.Text)
	util.AssertEquals(t, "AR", actual.DestType)
}

// 100% werden auf das Bankkonto gebucht
// 100% werden gegen das Kommitment-Konto gebucht
func TestProcessSVBeitrag(t *testing.T) {
	setUp()
	b := account.NewBooking("SV-Beitrag", "BEN", nil, 1385.10, "KKH, Ben", 5, 2017)

	Process(repository, *b)

	// Buchung wurde gegen Kommitment-Konto gebucht
	a, _ := repository.Get(owner.StakeholderKM.Id)
	b1 := a.Bookings[0]
	util.AssertFloatEquals(t, -1385.10, b1.Amount)
	util.AssertEquals(t, account.SVBeitrag, b1.DestType)

	// Buchung wurde aufs Bankkonto gebucht
	util.AssertEquals(t, 1, len(repository.BankAccount().Bookings))
	actual := repository.BankAccount().Bookings[0]
	util.AssertFloatEquals(t, -1385.10, actual.Amount)
	util.AssertEquals(t, "KKH, Ben", actual.Text)
	util.AssertEquals(t, "SV-Beitrag", actual.DestType)

}

// 100% werden auf das Bankkonto gebucht
// 100% werden gegen das Kommitment-Konto gebucht. Diese Regel ist nich unscharf:
// eigentlich müssen die 100% aufgeteilt werden auf: 70% auf Partner, 25% auf
// Kommitment und 5% auf Dealbringer
func TestProcessGWSteuer(t *testing.T) {
	setUp()

	b := account.NewBooking("GWSteuer", "K", nil, 2385.10, "STEUERKASSE HAMBURG STEUERNR 048/638/01147 GEW.ST 4VJ.17", 9, 2017)

	Process(repository, *b)

	// Buchung wurde gegen Kommitment-Konto gebucht
	a, _ := repository.Get(owner.StakeholderKM.Id)
	b1 := a.Bookings[0]
	assertBooking(t, b1, -2385.10, "STEUERKASSE HAMBURG STEUERNR 048/638/01147 GEW.ST 4VJ.17", account.GWSteuer)

	// Buchung wurde aufs Bankkonto gebucht
	util.AssertEquals(t, 1, len(repository.BankAccount().Bookings))
	actual := repository.BankAccount().Bookings[0]
	assertBooking(t, actual, -2385.10, "STEUERKASSE HAMBURG STEUERNR 048/638/01147 GEW.ST 4VJ.17", "GWSteuer")
}

func assertBooking(t *testing.T, b account.Booking, amount float64, text string, destType string) {
	util.AssertFloatEquals(t, amount, b.Amount)
	util.AssertEquals(t, text, b.Text)
	util.AssertEquals(t, destType, b.DestType)
}
