package parser

import (
	"testing"

	"bitbucket.org/rwirdemann/kontrol/owner"
	"bitbucket.org/rwirdemann/kontrol/util"
	"bitbucket.org/rwirdemann/kontrol/booking"
)

func TestImport(t *testing.T) {
	positions := Import("bookings.csv")
	util.AssertEquals(t, 7, len(positions))
	assertPosition(t, positions[0], "ER", "K", "Ulrike Klode", 2142, 2017, 2, 0, 0, 0)
	assertPosition(t, positions[1], "AR", "AN", "moebel.de", 4730.25, 2017, 1, 0, 0, 3975)
	assertPosition(t, positions[2], "AR", "JM", "RN_20170131-picue", 17225.25, 2017, 1, 10800, 3675, 0)
	assertPosition(t, positions[3], "GV", "RW", "Ralf", 6000, 2017, 1, 0, 0, 0)
	assertPosition(t, positions[4], "IS", "AN", "165", 8250, 2017, 12, 0, 0, 0)
	assertPosition(t, positions[5], "SV-Beitrag", "BW", "KKH, Ben", 1385.1, 2017, 7, 0, 0, 0)
	assertPosition(t, positions[6], "GWSteuer", "K", "STEUERKASSE HAMBURG STEUERNR 048/638/01147 GEW.ST 4VJ.17", 5160, 2017, 11, 0, 0, 0)
}

func assertPosition(t *testing.T, p booking.Booking, typ string, costCenter string, subject string,
	amount float64, year int, month int,
	nettoRW float64, nettoJM float64, nettoAN float64) {
	util.AssertEquals(t, typ, p.Typ)
	util.AssertEquals(t, costCenter, p.Responsible)
	util.AssertEquals(t, subject, p.Text)
	util.AssertEquals(t, amount, p.Amount)
	util.AssertEquals(t, year, p.Year)
	util.AssertEquals(t, month, p.Month)

	util.AssertEquals(t, nettoRW, p.Net[owner.StakeholderRW])
	util.AssertEquals(t, nettoJM, p.Net[owner.StakeholderJM])
	util.AssertEquals(t, nettoAN, p.Net[owner.StakeholderAN])
}
