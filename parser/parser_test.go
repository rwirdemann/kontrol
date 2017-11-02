package parser

import (
	"testing"

	"bitbucket.org/rwirdemann/kontrol/kontrol"

	"bitbucket.org/rwirdemann/kontrol/util"
)

func TestImport(t *testing.T) {
	positions := Import("bookings.csv")
	util.AssertEquals(t, 4, len(positions))
	assertPosition(t, positions[0], "ER", "K", "Ulrike Klode", 2142, 2017, 2, 0, 0, 0)
	assertPosition(t, positions[1], "AR", "AN", "moebel.de", 4730.25, 2017, 1, 0, 0, 3975)
	assertPosition(t, positions[2], "AR", "JM", "RN_20170131-picue", 17225.25, 2017, 1, 10800, 3675, 0)
	assertPosition(t, positions[3], "GV", "RW", "Ralf", 6000, 2017, 1, 0, 0, 0)
}

func assertPosition(t *testing.T, p kontrol.Booking, typ string, costCenter string, subject string,
	amount float64, year int, month int,
	nettoRW float64, nettoJM float64, nettoAN float64) {
	util.AssertEquals(t, typ, p.Extras.Typ)
	util.AssertEquals(t, costCenter, p.Extras.CostCenter)
	util.AssertEquals(t, subject, p.Text)
	util.AssertEquals(t, amount, p.Amount)
	util.AssertEquals(t, year, p.Year)
	util.AssertEquals(t, month, p.Month)

	util.AssertEquals(t, nettoRW, p.Extras.Net[kontrol.SH_RW])
	util.AssertEquals(t, nettoJM, p.Extras.Net[kontrol.SH_JM])
	util.AssertEquals(t, nettoAN, p.Extras.Net[kontrol.SH_AN])
}
