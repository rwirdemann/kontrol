package processing

import (
	"errors"
	"testing"

	"bitbucket.org/rwirdemann/kontrol/account"

	"bitbucket.org/rwirdemann/kontrol/domain"
	"bitbucket.org/rwirdemann/kontrol/util"
)

var repository account.Repository

func setUp() {
	repository = account.NewDefaultRepository()
	for _, sh := range domain.AllStakeholder {
		if sh.Type != domain.StakeholderTypeExtern {
			repository.Add(domain.NewAccount(sh))
		}
	}
}

func TestPartnerNettoAnteil(t *testing.T) {
	setUp()

	// given: a booking
	extras := domain.CsvBookingExtras{Typ: "AR", CostCenter: "JM"}
	extras.Net = make(map[domain.Stakeholder]float64)
	extras.Net[domain.StakeholderRW] = 10800.0
	extras.Net[domain.StakeholderJM] = 3675.0
	p := domain.Booking{Extras: extras, Amount: 17225.25, Text: "Rechnung 1234", Month: 1, Year: 2017}

	// when: the position is processed
	Process(repository, p)

	// then ralf 1 booking: his own net share
	accountRalf, _ := repository.Get(domain.StakeholderRW.Id)
	bookingsRalf := accountRalf.Bookings
	util.AssertEquals(t, 1, len(bookingsRalf))
	bRalf, _ := findBookingByText(bookingsRalf, "Rechnung 1234#NetShare#RW")
	util.AssertFloatEquals(t, 10800.0*domain.PartnerShare, bRalf.Amount)
	util.AssertEquals(t, 1, bRalf.Month)
	util.AssertEquals(t, 2017, bRalf.Year)
	util.AssertEquals(t, domain.Nettoanteil, bRalf.Typ)

	// and hannes got 3 bookings: his own net share and 2 provisions
	accountHannes, _ := repository.Get(domain.StakeholderJM.Id)
	bookingsHannes := accountHannes.Bookings
	util.AssertEquals(t, 3, len(bookingsHannes))

	// net share
	b, _ := findBookingByText(bookingsHannes, "Rechnung 1234#NetShare#JM")
	util.AssertFloatEquals(t, 3675.0*domain.PartnerShare, b.Amount)
	util.AssertEquals(t, 1, b.Month)
	util.AssertEquals(t, 2017, b.Year)

	// provision from ralf
	provisionRalf, _ := findBookingByText(bookingsHannes, "Rechnung 1234#Provision#RW")
	util.AssertFloatEquals(t, 10800.0*domain.PartnerProvision, provisionRalf.Amount)
	util.AssertEquals(t, domain.Vertriebsprovision, provisionRalf.Typ)

	// // provision from hannes
	provisionHannes, _ := findBookingByText(bookingsHannes, "Rechnung 1234#Provision#JM")
	util.AssertFloatEquals(t, 3675.0*domain.PartnerProvision, provisionHannes.Amount)
	util.AssertEquals(t, domain.Vertriebsprovision, provisionHannes.Typ)

	// kommitment got 25% from ralfs net booking
	accountKommitment, _ := repository.Get(domain.StakeholderKM.Id)
	bookingsKommitment := accountKommitment.Bookings
	util.AssertEquals(t, 2, len(bookingsKommitment))
	kommitmentRalf, _ := findBookingByText(bookingsKommitment, "Rechnung 1234#Kommitment#RW")
	util.AssertFloatEquals(t, 10800.0*domain.KommmitmentShare, kommitmentRalf.Amount)
	util.AssertEquals(t, domain.Kommitmentanteil, kommitmentRalf.Typ)

	// and kommitment got 25% from hannes net booking
	kommitmentHannes, _ := findBookingByText(bookingsKommitment, "Rechnung 1234#Kommitment#JM")
	util.AssertFloatEquals(t, 3675.0*domain.KommmitmentShare, kommitmentHannes.Amount)
	util.AssertEquals(t, domain.Kommitmentanteil, kommitmentHannes.Typ)
}

func findBookingByText(bookings []domain.Booking, text string) (*domain.Booking, error) {
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
	extras := domain.CsvBookingExtras{Typ: "AR", CostCenter: "JM"}
	extras.Net = make(map[domain.Stakeholder]float64)
	extras.Net[domain.StakeholderBW] = 10800.0
	p := domain.Booking{Extras: extras, Amount: 12852.0, Text: "Rechnung 1234", Month: 1, Year: 2017}

	// when: the position is processed
	Process(repository, p)

	// and hannes got his provision
	accountHannes, _ := repository.Get(domain.StakeholderJM.Id)
	provision := accountHannes.Bookings[0]
	util.AssertFloatEquals(t, 10800.0*domain.PartnerProvision, provision.Amount)
	util.AssertEquals(t, domain.Vertriebsprovision, provision.Typ)

	// and kommitment got 95%
	util.AssertEquals(t, 1, len(accountHannes.Bookings))
	accountKommitment, _ := repository.Get(domain.StakeholderKM.Id)
	kommitment := accountKommitment.Bookings[0]
	util.AssertFloatEquals(t, 10800.0*domain.KommmitmentEmployeeShare, kommitment.Amount)
	util.AssertEquals(t, domain.Kommitmentanteil, kommitment.Typ)

	// 100% is booked to employee account to see how much money is made by this employee
	accountBen, _ := repository.Get(domain.StakeholderBW.Id)
	util.AssertEquals(t, 1, len(accountBen.Bookings))
	bookingBen := accountBen.Bookings[0]
	util.AssertFloatEquals(t, 10800.0, bookingBen.Amount)
}

func TestExternNettoAnteil(t *testing.T) {
	setUp()

	// given: a booking
	extras := domain.CsvBookingExtras{Typ: "AR", CostCenter: "JM"}
	extras.Net = make(map[domain.Stakeholder]float64)
	extras.Net[domain.StakeholderEX] = 10800.0
	p := domain.Booking{Extras: extras, Amount: 12852.0, Text: "Rechnung 1234", Month: 1, Year: 2017}

	// when: the position is processed
	Process(repository, p)

	// and hannes got his provision
	accountHannes, _ := repository.Get(domain.StakeholderJM.Id)
	provision := accountHannes.Bookings[0]
	util.AssertFloatEquals(t, 10800.0*domain.PartnerProvision, provision.Amount)
	util.AssertEquals(t, domain.Vertriebsprovision, provision.Typ)

	// and kommitment got 95%
	util.AssertEquals(t, 1, len(accountHannes.Bookings))
	accountKommitment, _ := repository.Get(domain.StakeholderKM.Id)
	kommitment := accountKommitment.Bookings[0]
	util.AssertFloatEquals(t, 10800.0*domain.KommmitmentExternShare, kommitment.Amount)
	util.AssertEquals(t, domain.Kommitmentanteil, kommitment.Typ)
}

func TestEingangsrechnung(t *testing.T) {
	setUp()

	// given: a booking
	extras := domain.CsvBookingExtras{Typ: "ER", CostCenter: "K"}
	p := domain.Booking{Extras: extras, Amount: 12852.0, Text: "Eingangsrechnung 1234", Month: 1, Year: 2017}

	// when: the position is processed
	Process(repository, p)

	// the invoice is booked to the kommitment account
	accountKommitment, _ := repository.Get(domain.StakeholderKM.Id)
	util.AssertEquals(t, 1, len(accountKommitment.Bookings))
	kommitment := accountKommitment.Bookings[0]
	util.AssertFloatEquals(t, util.Net(-12852.0), kommitment.Amount)
	util.AssertEquals(t, domain.Eingangsrechnung, kommitment.Typ)
}

func TestPartnerWithdrawals(t *testing.T) {
	setUp()

	extras := domain.CsvBookingExtras{Typ: "GV", CostCenter: "RW"}
	extras.Net = make(map[domain.Stakeholder]float64)
	b := domain.Booking{Extras: extras, Amount: 6000}
	Process(repository, b)
	accountRalf, _ := repository.Get(domain.StakeholderRW.Id)
	util.AssertEquals(t, 1, len(accountRalf.Bookings))
	bRalf := accountRalf.Bookings[0]
	util.AssertFloatEquals(t, -6000, bRalf.Amount)
	util.AssertEquals(t, domain.Entnahme, bRalf.Typ)
}

func TestInterneStunden(t *testing.T) {
	setUp()

	// given: a internal hours booking
	extras := domain.CsvBookingExtras{Typ: "IS", CostCenter: "AN"}
	p := domain.Booking{Extras: extras, Amount: 8250.00, Text: "Internet Stunden 2017", Month: 12, Year: 2017}

	// when: the position is processed
	Process(repository, p)

	// the booking is booked to anke's account
	accountAnke, _ := repository.Get(domain.StakeholderAN.Id)
	util.AssertEquals(t, 1, len(accountAnke.Bookings))
	booking := accountAnke.Bookings[0]
	util.AssertFloatEquals(t, 8250.00, booking.Amount)
	util.AssertEquals(t, domain.InterneStunden, booking.Typ)
}
