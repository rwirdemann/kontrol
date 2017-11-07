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

func TestPartnerNettoAnteil(t *testing.T) {
	setUp()

	// given: a booking
	extras := kontrol.CsvBookingExtras{Typ: "AR", CostCenter: "JM"}
	extras.Net = make(map[kontrol.Stakeholder]float64)
	extras.Net[kontrol.SH_RW] = 10800.0
	extras.Net[kontrol.SH_JM] = 3675.0
	p := kontrol.Booking{Extras: extras, Amount: 17225.25, Text: "Rechnung 1234", Month: 1, Year: 2017}

	// when: the position is processed
	Process(p)

	// then ralf got his net booking
	util.AssertEquals(t, 1, len(kontrol.Accounts[kontrol.SH_RW.Id].Bookings))
	bRalf := kontrol.Accounts[kontrol.SH_RW.Id].Bookings[0]
	util.AssertFloatEquals(t, 10800.0*kontrol.PartnerShare, bRalf.Amount)
	util.AssertEquals(t, "Rechnung 1234#RW", bRalf.Text)
	util.AssertEquals(t, 1, bRalf.Month)
	util.AssertEquals(t, 2017, bRalf.Year)
	util.AssertEquals(t, kontrol.Nettoanteil, bRalf.Typ)

	// and hannes got his net booking
	util.AssertEquals(t, 2, len(kontrol.Accounts[kontrol.SH_JM.Id].Bookings))
	bHannes := kontrol.Accounts[kontrol.SH_JM.Id].Bookings[0]
	util.AssertFloatEquals(t, 3675*kontrol.PartnerShare, bHannes.Amount)
	util.AssertEquals(t, "Rechnung 1234#JM", bHannes.Text)
	util.AssertEquals(t, 1, bHannes.Month)
	util.AssertEquals(t, 2017, bHannes.Year)

	// and hannes got his provision
	provision := kontrol.Accounts[kontrol.SH_JM.Id].Bookings[1]
	util.AssertFloatEquals(t, 14475.0*kontrol.PartnerProvision, provision.Amount)
	util.AssertEquals(t, kontrol.Vertriebsprovision, provision.Typ)

	// and kommitment got 75%
	util.AssertEquals(t, 2, len(kontrol.Accounts[kontrol.SH_KM.Id].Bookings))
	bKommitmentRalf := kontrol.Accounts[kontrol.SH_KM.Id].Bookings[0]
	util.AssertFloatEquals(t, 10800.0*kontrol.KommmitmentShare, bKommitmentRalf.Amount)
	util.AssertEquals(t, kontrol.Kommitmentanteil, bKommitmentRalf.Typ)

	bKommitmentHannes := kontrol.Accounts[kontrol.SH_KM.Id].Bookings[1]
	util.AssertFloatEquals(t, 3675.0*kontrol.KommmitmentShare, bKommitmentHannes.Amount)
	util.AssertEquals(t, kontrol.Kommitmentanteil, bKommitmentRalf.Typ)
}

func TestExternNettoAnteil(t *testing.T) {
	setUp()

	// given: a booking
	extras := kontrol.CsvBookingExtras{Typ: "AR", CostCenter: "JM"}
	extras.Net = make(map[kontrol.Stakeholder]float64)
	extras.Net[kontrol.SH_EX] = 10800.0
	p := kontrol.Booking{Extras: extras, Amount: 12852.0, Text: "Rechnung 1234", Month: 1, Year: 2017}

	// when: the position is processed
	Process(p)

	// and hannes got his provision
	provision := kontrol.Accounts[kontrol.SH_JM.Id].Bookings[0]
	util.AssertFloatEquals(t, 10800.0*kontrol.PartnerProvision, provision.Amount)
	util.AssertEquals(t, kontrol.Vertriebsprovision, provision.Typ)

	// and kommitment got 95%
	util.AssertEquals(t, 1, len(kontrol.Accounts[kontrol.SH_KM.Id].Bookings))
	kommitment := kontrol.Accounts[kontrol.SH_KM.Id].Bookings[0]
	util.AssertFloatEquals(t, 10800.0*kontrol.KommmitmentExternShare, kommitment.Amount)
	util.AssertEquals(t, kontrol.Kommitmentanteil, kommitment.Typ)
}

func TestPartnerWithdrawals(t *testing.T) {
	setUp()

	extras := kontrol.CsvBookingExtras{Typ: "GV", CostCenter: "RW"}
	extras.Net = make(map[kontrol.Stakeholder]float64)
	b := kontrol.Booking{Extras: extras, Amount: 6000}
	Process(b)
	util.AssertEquals(t, 1, len(kontrol.Accounts[kontrol.SH_RW.Id].Bookings))
	bRalf := kontrol.Accounts[kontrol.SH_RW.Id].Bookings[0]
	util.AssertFloatEquals(t, -6000, bRalf.Amount)
	util.AssertEquals(t, kontrol.Entnahme, bRalf.Typ)
}
