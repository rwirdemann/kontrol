package processing

import (
	"errors"
	"testing"

	"bitbucket.org/rwirdemann/kontrol/kontrol"
	"bitbucket.org/rwirdemann/kontrol/util"
)

func setUp() {
	for _, a := range kontrol.Accounts {
		a.Bookings = make([]kontrol.Booking, 0)
	}
}

func TestPartnerNettoAnteil(t *testing.T) {
	setUp()

	// given: a booking
	extras := kontrol.CsvBookingExtras{Typ: "AR", CostCenter: "JM"}
	extras.Net = make(map[kontrol.Stakeholder]float64)
	extras.Net[kontrol.StakeholderRW] = 10800.0
	extras.Net[kontrol.StakeholderJM] = 3675.0
	p := kontrol.Booking{Extras: extras, Amount: 17225.25, Text: "Rechnung 1234", Month: 1, Year: 2017}

	// when: the position is processed
	Process(p)

	// then ralf 1 booking: his own net share
	bookingsRalf := kontrol.Accounts[kontrol.StakeholderRW.Id].Bookings
	util.AssertEquals(t, 1, len(kontrol.Accounts[kontrol.StakeholderRW.Id].Bookings))
	bRalf, _ := findBookingByText(bookingsRalf, "Rechnung 1234#NetShare#RW")
	util.AssertFloatEquals(t, 10800.0*kontrol.PartnerShare, bRalf.Amount)
	util.AssertEquals(t, 1, bRalf.Month)
	util.AssertEquals(t, 2017, bRalf.Year)
	util.AssertEquals(t, kontrol.Nettoanteil, bRalf.Typ)

	// and hannes got 3 bookings: his own net share and 2 provisions
	bookingsHannes := kontrol.Accounts[kontrol.StakeholderJM.Id].Bookings
	util.AssertEquals(t, 3, len(bookingsHannes))

	// net share
	b, _ := findBookingByText(bookingsHannes, "Rechnung 1234#NetShare#JM")
	util.AssertFloatEquals(t, 3675.0*kontrol.PartnerShare, b.Amount)
	util.AssertEquals(t, 1, b.Month)
	util.AssertEquals(t, 2017, b.Year)

	// provision from ralf
	provisionRalf, _ := findBookingByText(bookingsHannes, "Rechnung 1234#Provision#RW")
	util.AssertFloatEquals(t, 10800.0*kontrol.PartnerProvision, provisionRalf.Amount)
	util.AssertEquals(t, kontrol.Vertriebsprovision, provisionRalf.Typ)

	// // provision from hannes
	provisionHannes, _ := findBookingByText(bookingsHannes, "Rechnung 1234#Provision#JM")
	util.AssertFloatEquals(t, 3675.0*kontrol.PartnerProvision, provisionHannes.Amount)
	util.AssertEquals(t, kontrol.Vertriebsprovision, provisionHannes.Typ)

	// kommitment got 25% from ralfs net booking
	bookingsKommitment := kontrol.Accounts[kontrol.StakeholderKM.Id].Bookings
	util.AssertEquals(t, 2, len(bookingsKommitment))
	kommitmentRalf, _ := findBookingByText(bookingsKommitment, "Rechnung 1234#Kommitment#RW")
	util.AssertFloatEquals(t, 10800.0*kontrol.KommmitmentShare, kommitmentRalf.Amount)
	util.AssertEquals(t, kontrol.Kommitmentanteil, kommitmentRalf.Typ)

	// and kommitment got 25% from hannes net booking
	kommitmentHannes, _ := findBookingByText(bookingsKommitment, "Rechnung 1234#Kommitment#JM")
	util.AssertFloatEquals(t, 3675.0*kontrol.KommmitmentShare, kommitmentHannes.Amount)
	util.AssertEquals(t, kontrol.Kommitmentanteil, kommitmentHannes.Typ)
}

func findBookingByText(bookings []kontrol.Booking, text string) (*kontrol.Booking, error) {
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
	extras := kontrol.CsvBookingExtras{Typ: "AR", CostCenter: "JM"}
	extras.Net = make(map[kontrol.Stakeholder]float64)
	extras.Net[kontrol.StakeholderBW] = 10800.0
	p := kontrol.Booking{Extras: extras, Amount: 12852.0, Text: "Rechnung 1234", Month: 1, Year: 2017}

	// when: the position is processed
	Process(p)

	// and hannes got his provision
	provision := kontrol.Accounts[kontrol.StakeholderJM.Id].Bookings[0]
	util.AssertFloatEquals(t, 10800.0*kontrol.PartnerProvision, provision.Amount)
	util.AssertEquals(t, kontrol.Vertriebsprovision, provision.Typ)

	// and kommitment got 95%
	util.AssertEquals(t, 1, len(kontrol.Accounts[kontrol.StakeholderKM.Id].Bookings))
	kommitment := kontrol.Accounts[kontrol.StakeholderKM.Id].Bookings[0]
	util.AssertFloatEquals(t, 10800.0*kontrol.KommmitmentEmployeeShare, kommitment.Amount)
	util.AssertEquals(t, kontrol.Kommitmentanteil, kommitment.Typ)

	// 100% is booked to employee account to see how much money is made by this employee
	util.AssertEquals(t, 1, len(kontrol.Accounts[kontrol.StakeholderBW.Id].Bookings))
	bookingBen := kontrol.Accounts[kontrol.StakeholderBW.Id].Bookings[0]
	util.AssertFloatEquals(t, 10800.0, bookingBen.Amount)
}

func TestExternNettoAnteil(t *testing.T) {
	setUp()

	// given: a booking
	extras := kontrol.CsvBookingExtras{Typ: "AR", CostCenter: "JM"}
	extras.Net = make(map[kontrol.Stakeholder]float64)
	extras.Net[kontrol.StakeholderEX] = 10800.0
	p := kontrol.Booking{Extras: extras, Amount: 12852.0, Text: "Rechnung 1234", Month: 1, Year: 2017}

	// when: the position is processed
	Process(p)

	// and hannes got his provision
	provision := kontrol.Accounts[kontrol.StakeholderJM.Id].Bookings[0]
	util.AssertFloatEquals(t, 10800.0*kontrol.PartnerProvision, provision.Amount)
	util.AssertEquals(t, kontrol.Vertriebsprovision, provision.Typ)

	// and kommitment got 95%
	util.AssertEquals(t, 1, len(kontrol.Accounts[kontrol.StakeholderKM.Id].Bookings))
	kommitment := kontrol.Accounts[kontrol.StakeholderKM.Id].Bookings[0]
	util.AssertFloatEquals(t, 10800.0*kontrol.KommmitmentExternShare, kommitment.Amount)
	util.AssertEquals(t, kontrol.Kommitmentanteil, kommitment.Typ)
}

func TestEingangsrechnung(t *testing.T) {
	setUp()

	// given: a booking
	extras := kontrol.CsvBookingExtras{Typ: "ER", CostCenter: "K"}
	p := kontrol.Booking{Extras: extras, Amount: 12852.0, Text: "Eingangsrechnung 1234", Month: 1, Year: 2017}

	// when: the position is processed
	Process(p)

	// the invoice is booked to the kommitment account
	util.AssertEquals(t, 1, len(kontrol.Accounts[kontrol.StakeholderKM.Id].Bookings))
	kommitment := kontrol.Accounts[kontrol.StakeholderKM.Id].Bookings[0]
	util.AssertFloatEquals(t, util.Net(-12852.0), kommitment.Amount)
	util.AssertEquals(t, kontrol.Eingangsrechnung, kommitment.Typ)
}

func TestPartnerWithdrawals(t *testing.T) {
	setUp()

	extras := kontrol.CsvBookingExtras{Typ: "GV", CostCenter: "RW"}
	extras.Net = make(map[kontrol.Stakeholder]float64)
	b := kontrol.Booking{Extras: extras, Amount: 6000}
	Process(b)
	util.AssertEquals(t, 1, len(kontrol.Accounts[kontrol.StakeholderRW.Id].Bookings))
	bRalf := kontrol.Accounts[kontrol.StakeholderRW.Id].Bookings[0]
	util.AssertFloatEquals(t, -6000, bRalf.Amount)
	util.AssertEquals(t, kontrol.Entnahme, bRalf.Typ)
}
