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
	extras := account.CsvBookingExtras{SourceType: "AR", CostCenter: "JM"}
	extras.Net = make(map[owner.Stakeholder]float64)
	extras.Net[owner.StakeholderRW] = 10800.0
	extras.Net[owner.StakeholderJM] = 3675.0
	p := account.Booking{Extras: extras, Amount: 17225.25, Text: "Rechnung 1234", Month: 1, Year: 2017}

	// when: the position is processed
	Process(repository, p)

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
	extras := account.CsvBookingExtras{SourceType: "AR", CostCenter: "JM"}
	extras.Net = make(map[owner.Stakeholder]float64)
	extras.Net[owner.StakeholderBW] = 10800.0
	p := account.Booking{Extras: extras, Amount: 12852.0, Text: "Rechnung 1234", Month: 1, Year: 2017}

	// when: the position is processed
	Process(repository, p)

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
	extras := account.CsvBookingExtras{SourceType: "AR", CostCenter: "JM"}
	extras.Net = make(map[owner.Stakeholder]float64)
	extras.Net[owner.StakeholderEX] = 10800.0
	p := account.Booking{Extras: extras, Amount: 12852.0, Text: "Rechnung 1234", Month: 1, Year: 2017}

	// when: the position is processed
	Process(repository, p)

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
	extras := account.CsvBookingExtras{SourceType: "ER", CostCenter: "K"}
	p := account.Booking{Extras: extras, Amount: 12852.0, Text: "Eingangsrechnung 1234", Month: 1, Year: 2017}

	// when: the position is processed
	Process(repository, p)

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
	b := account.Booking{Extras: extras, Amount: 6000}
	Process(repository, b)
	accountRalf, _ := repository.Get(owner.StakeholderRW.Id)
	util.AssertEquals(t, 1, len(accountRalf.Bookings))
	bRalf := accountRalf.Bookings[0]
	util.AssertFloatEquals(t, -6000, bRalf.Amount)
	util.AssertEquals(t, account.Entnahme, bRalf.DestType)
}

func TestInterneStunden(t *testing.T) {
	setUp()

	// given: a internal hours booking
	extras := account.CsvBookingExtras{SourceType: "IS", CostCenter: "AN"}
	p := account.Booking{Extras: extras, Amount: 8250.00, Text: "Internet Stunden 2017", Month: 12, Year: 2017}

	// when: the position is processed
	Process(repository, p)

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
	util.AssertEquals(t, 0, len(repository.CollectiveAccount().Bookings))
}

func TestSVBeitrag(t *testing.T) {
	setUp()

	// given: a sv-beitrag booking
	extras := account.CsvBookingExtras{SourceType: "SV-Beitrag", CostCenter: "BEN"}
	b := account.Booking{Extras: extras, Amount: 1385.10, Text: "KKH, Ben"}

	// when: the booking is processed
	Process(repository, b)

	// the booking is booked against kommitment account
	a, _ := repository.Get(owner.StakeholderKM.Id)
	b1 := a.Bookings[0]
	util.AssertFloatEquals(t, -1385.10, b1.Amount)
	util.AssertEquals(t, account.SVBeitrag, b1.DestType)
}

func TestBookEingangsrechnungToBankAccount(t *testing.T) {
	setUp()
	extras := account.CsvBookingExtras{SourceType: "ER", CostCenter: "K"}
	extras.Net = make(map[owner.Stakeholder]float64)
	b := account.Booking{Extras: extras, Amount: 6000, Text: "Eingangsrechnung"}

	Process(repository, b)

	util.AssertEquals(t, 1, len(repository.CollectiveAccount().Bookings))
	actual := repository.CollectiveAccount().Bookings[0]
	util.AssertFloatEquals(t, util.Net(-6000), actual.Amount)
	util.AssertEquals(t, "Eingangsrechnung", actual.Text)
	util.AssertEquals(t, "ER", actual.DestType)
}

func TestBookAusgangsrechnungToBankAccount(t *testing.T) {
	setUp()
	extras := account.CsvBookingExtras{SourceType: "AR", CostCenter: "K"}
	extras.Net = make(map[owner.Stakeholder]float64)
	b := account.Booking{Extras: extras, Amount: 6000, Text: "Ausgangsrechnung"}

	Process(repository, b)

	util.AssertEquals(t, 1, len(repository.CollectiveAccount().Bookings))
	actual := repository.CollectiveAccount().Bookings[0]
	util.AssertFloatEquals(t, util.Net(6000), actual.Amount)
	util.AssertEquals(t, "Ausgangsrechnung", actual.Text)
	util.AssertEquals(t, "AR", actual.DestType)
}

func TestBookSVBeitragToBankAccount(t *testing.T) {
	setUp()
	extras := account.CsvBookingExtras{SourceType: "SV-Beitrag", CostCenter: "BEN"}
	b := account.Booking{Extras: extras, Amount: 1385.10, Text: "KKH, Ben"}

	Process(repository, b)

	util.AssertEquals(t, 1, len(repository.CollectiveAccount().Bookings))
	actual := repository.CollectiveAccount().Bookings[0]
	util.AssertFloatEquals(t, -1385.10, actual.Amount)
	util.AssertEquals(t, "KKH, Ben", actual.Text)
	util.AssertEquals(t, "SV-Beitrag", actual.DestType)
}
