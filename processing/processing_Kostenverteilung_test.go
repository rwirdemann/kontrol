package processing

import (
	"github.com/ahojsenn/kontrol/accountSystem"
	"github.com/ahojsenn/kontrol/booking"
	"github.com/ahojsenn/kontrol/valueMagnets"
	"github.com/stretchr/testify/assert"
	"log"
	"math"
	"testing"
	"time"
)

func TestKostenerteilung (t *testing.T) {
	as := accountSystem.NewDefaultAccountSystem()
	as.ClearBookings()

	its2018 := time.Date(2018, 1, 23, 0, 0, 0, 0, time.UTC)

	net := make(map[valueMagnets.Stakeholder]float64)
	shrepo := valueMagnets.Stakeholder{}
	net[shrepo.Get("AN")] = 119.0
	net[shrepo.Get("JM")] = 119.0
	net[shrepo.Get("BW")] = 119.0

	hauptbuch := as.GetCollectiveAccount_thisYear(2018)
	// Anke und Johannes haben Nettoeinnahmen von 100€
	// Ben hat Kosten von 10€ netto
	// Bebn hat Einkünfte von 100€ netto, davon werden 70% angerechnet = 70€
	// ==> Ben sollte nun  70€ Jahreseinkommen haben...
	bkgs := &(hauptbuch.Bookings)

	*bkgs =  append (*bkgs,  *booking.NewBooking(13, "AR", "", "", "K", "Project-X", net, 1190, "Anke+Johannes", 1, 2018, its2018) )
	*bkgs =  append (*bkgs,  *booking.NewBooking(13, "ER", "", "", "JM", "Project-X", net, 11.9, "H-costs", 1, 2018, its2018) )
	*bkgs =  append (*bkgs,  *booking.NewBooking(13, "ER", "", "", "JM", "Project-X", net, 11.9, "H-costs", 1, 2018, its2018) )
	*bkgs =  append (*bkgs,  *booking.NewBooking(13, "ER", "", "", "BW", "Project-X", net, 11.9, "H-costs", 1, 2018, its2018) )

	log.Println("in TestKostenerteilung: ", len(as.GetCollectiveAccount_thisYear(2018).Bookings) )

	// nun verteilen
	for _, p := range as.GetCollectiveAccount_thisYear(2018).Bookings {
		Process(as, p)
	}
	Kostenerteilung(as)

	// 33% von 200€ k-anteil + 50% von 800€ - Gewerbesteuer
	assert.Equal(t, -10.0, math.Round(StakeholderYearlyIncome(as, "BW")) )
	assert.Equal(t, -20.0, math.Round(StakeholderYearlyIncome(as, "JM")) )
	assert.Equal(t, 0.0, math.Round(StakeholderYearlyIncome(as, "AN")) )

	return

}