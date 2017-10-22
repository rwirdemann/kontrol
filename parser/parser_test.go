package parser

import (
	"testing"

	"bitbucket.org/rwirdemann/kontrol/kontrol"

	"bitbucket.org/rwirdemann/kontrol/util"
)

func TestImport(t *testing.T) {
	positions := Import("bookings.csv")
	util.AssertEquals(t, 3, len(positions))
	assertPosition(t, positions[0], "ER", "K", "Ulrike Klode", 2142, 2017, 2, 0, 0, 0)
	assertPosition(t, positions[1], "AR", "AN", "moebel.de", 4730.25, 2017, 1, 0, 0, 3975)
	assertPosition(t, positions[2], "AR", "JM", "RN_20170131-picue", 17225.25, 2017, 1, 10800, 3675, 0)
}

func assertPosition(t *testing.T, p kontrol.Position, typ string, costCenter string, subject string,
	amount float64, year int, month int,
	nettoRW float64, nettoJM float64, nettoAN float64) {
	util.AssertEquals(t, typ, p.Typ)
	util.AssertEquals(t, costCenter, p.CostCenter)
	util.AssertEquals(t, subject, p.Subject)
	util.AssertEquals(t, amount, p.Amount)
	util.AssertEquals(t, year, p.Year)
	util.AssertEquals(t, month, p.Month)

	util.AssertEquals(t, nettoRW, p.Net[kontrol.SA_RW])
	util.AssertEquals(t, nettoJM, p.Net[kontrol.SA_JM])
	util.AssertEquals(t, nettoAN, p.Net[kontrol.SA_AN])
}
