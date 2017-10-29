package processing

import (
	"testing"

	"bitbucket.org/rwirdemann/kontrol/kontrol"
	"bitbucket.org/rwirdemann/kontrol/util"
)

func setUp() {
	for _, a := range kontrol.Accounts {
		a.Bookings = make([]kontrol.Booking, 0)
	}
}

func TestProcessing(t *testing.T) {
	setUp()

	// given: a positon
	p := kontrol.Position{Typ: "AR", CostCenter: "JM", Amount: 17225.25, Subject: "Rechnung 1234", Month: 1, Year: 2017}
	p.Net = make(map[string]float64)
	p.Net[kontrol.SA_RW] = 10800.0
	p.Net[kontrol.SA_JM] = 3675.0

	// when: the position is processed
	Process(p)

	// then ralf got his net booking
	util.AssertEquals(t, 1, len(kontrol.Accounts[kontrol.SA_RW].Bookings))
	bRalf := kontrol.Accounts[kontrol.SA_RW].Bookings[0]
	util.AssertFloatEquals(t, 10800.0*kontrol.PartnerShare, bRalf.Amount)
	util.AssertEquals(t, "Rechnung 1234", bRalf.Text)
	util.AssertEquals(t, 1, bRalf.Month)
	util.AssertEquals(t, 2017, bRalf.Year)

	// and hannes got his net booking
	util.AssertEquals(t, 1, len(kontrol.Accounts[kontrol.SA_JM].Bookings))
	bHannes := kontrol.Accounts[kontrol.SA_JM].Bookings[0]
	util.AssertFloatEquals(t, 3675*kontrol.PartnerShare, bHannes.Amount)
	util.AssertEquals(t, "Rechnung 1234", bHannes.Text)
	util.AssertEquals(t, 1, bHannes.Month)
	util.AssertEquals(t, 2017, bHannes.Year)

	// TODO: book 5 percent to cost center
}

func TestPartnerWithdrawals(t *testing.T) {
	setUp()

	p := kontrol.Position{Typ: "GV", CostCenter: "RW", Amount: 6000}
	Process(p)
	util.AssertEquals(t, 1, len(kontrol.Accounts[kontrol.SA_RW].Bookings))
	bRalf := kontrol.Accounts[kontrol.SA_RW].Bookings[0]
	util.AssertFloatEquals(t, -6000, bRalf.Amount)
}
