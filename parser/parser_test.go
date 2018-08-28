package parser

import (
		"testing"
	"time"

	"github.com/ahojsenn/kontrol/booking"
	"github.com/ahojsenn/kontrol/owner"
	"github.com/stretchr/testify/assert"
	)

func TestImport(t *testing.T) {
	var positions []booking.Booking
	Import("bookings.csv", 2017, &positions)
	assert.Equal(t, 11, len(positions))
	assertPosition(t, positions[0], "ER", "800", "1337", "K", "Ulrike Klode", 2142, 2017, 2, 0, 0, 0)
	assertPosition(t, positions[1], "AR", "800", "1337", "AN", "moebel.de", 4730.25, 2017, 1, 0, 0, 3975)
	assertPosition(t, positions[2], "AR", "800", "1337", "JM", "RN_20170131-picue", 17225.25, 2017, 1, 10800, 3675, 0)
	assertPosition(t, positions[3], "GV", "800", "1337", "RW", "Ralf", 6000, 2017, 1, 0, 0, 0)
	assertPosition(t, positions[4], "IS", "800", "1337", "AN", "165", 8250, 2017, 12, 0, 0, 0)
	assertPosition(t, positions[5], "SV-Beitrag", "800", "1337", "BW", "KKH, Ben", 1385.1, 2017, 7, 0, 0, 0)
	assertPosition(t, positions[6], "GWSteuer", "800", "1337", "K", "STEUERKASSE HAMBURG STEUERNR 048/638/01147 GEW.ST 4VJ.17", 5160, 2017, 11, 0, 0, 0)
	assertPosition(t, positions[7], "Gehalt", "800", "1337", "BW", "Gehalt Ben 2017-07", 3869.65, 2017, 7, 0, 0, 0)
	assertPosition(t, positions[8], "AR", "800", "1337", "BW", "wlw 2017-11-07-10053", 7461.3, 2017, 10, 0, 0, 0)
	assertPosition(t, positions[9], "LNSteuer", "800", "1337", "BW", "STEUERKASSE HAMBURGSTEUERNR 046/638/02084 LOHNST DEZ.17 1.511,45EUR UMS.ST NOV.17 10.843,11EUR", 1511.45, 2017, 12, 0, 0, 0)
	assertPosition(t, positions[10], "Anfangsbestand", "800", "1337", "Rückstellung", "Anfangsbestand aus Vorjahr", 42.23, 2017, 2, 0, 0, 0)
}

func assertPosition(t *testing.T, p booking.Booking, typ string, soll string, haben string,
	costCenter string, subject string,
	amount float64, year int, month int,
	nettoRW float64, nettoJM float64, nettoAN float64) {
	assert.Equal(t, typ, p.Typ)
	assert.Equal(t, soll, p.Soll)
	assert.Equal(t, haben, p.Haben)
	assert.Equal(t, costCenter, p.Responsible)
	assert.Equal(t, subject, p.Text)
	assert.Equal(t, amount, p.Amount)
	assert.Equal(t, year, p.Year)
	assert.Equal(t, month, p.Month)

	assert.Equal(t, nettoRW, p.Net[ owner.StakeholderRepository{}.Get("RW") ])
	assert.Equal(t, nettoJM, p.Net[ owner.StakeholderRepository{}.Get("JM") ])
	assert.Equal(t, nettoAN, p.Net[ owner.StakeholderRepository{}.Get("AN") ])
}

func TestParseFileCreated(t *testing.T) {
	assert.Equal(t, time.Date(2017, 1, 1, 0, 0, 0, 0, time.UTC), parseFileCreated("1.1.2017"))
	assert.Equal(t, time.Date(2017, 12, 1, 0, 0, 0, 0, time.UTC), parseFileCreated("1.12.2017"))
	assert.Equal(t, time.Date(2017, 7, 31, 0, 0, 0, 0, time.UTC), parseFileCreated("31.7.2017"))
	assert.Equal(t, time.Time{}, parseFileCreated(""))
}
