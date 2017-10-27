package processing

import (
	"testing"

	"bitbucket.org/rwirdemann/kontrol/kontrol"
	"bitbucket.org/rwirdemann/kontrol/util"
)

func TestProcessing(t *testing.T) {

	// given: a positon
	p := kontrol.Position{Typ: "AR", CostCenter: "JM", Amount: 17225.25, Subject: "Rechnung 1234"}
	p.Net = make(map[string]float64)
	p.Net[kontrol.SA_RW] = 10800
	p.Net[kontrol.SA_JM] = 3675

	// when: the position is processed
	Process(p)

	// then ralf got his net booking
	util.AssertEquals(t, 1, len(kontrol.Accounts[kontrol.SA_RW].Bookings))
	bRalf := kontrol.Accounts[kontrol.SA_RW].Bookings[0]
	util.AssertFloatEquals(t, 10800, bRalf.Amount)
	util.AssertEquals(t, "Rechnung 1234", bRalf.Text)

	// and hannes got his net booking
	util.AssertEquals(t, 1, len(kontrol.Accounts[kontrol.SA_JM].Bookings))
	bHannes := kontrol.Accounts[kontrol.SA_JM].Bookings[0]
	util.AssertFloatEquals(t, 3675, bHannes.Amount)
	util.AssertEquals(t, "Rechnung 1234", bHannes.Text)

	// TODO: book only 70% of net value
	// TODO: book 5 percent to cost center
}
