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

func TestErloesverteilungAnEmployees (t *testing.T) {
	DEBUG := false
	as := accountSystem.NewDefaultAccountSystem()
	as.ClearBookings()
	its2018 := time.Date(2018, 1, 23, 0, 0, 0, 0, time.UTC)

	net := make(map[valueMagnets.Stakeholder]float64)
	shrepo := valueMagnets.Stakeholder{}
	net[shrepo.Get("AN")] = 100.0
	net[shrepo.Get("JM")] = 100.0
	net[shrepo.Get("BW")] = 100.0

	hauptbuch := as.GetCollectiveAccount_thisYear()
	// Anke und Johannes haben Nettoeinnahmen von 100€
	// Ben hat Kosten von 10€ netto
	// ==> Ben sollte nun  70€ + 70€ + 5% von 300€ netto€ = 155€ Jahreseinkommen haben...
	bkgs := &(hauptbuch.Bookings)
	*bkgs = append(*bkgs, *booking.NewBooking(13, "AR", "", "", "K", "Project-X", net, 357, "Anke+Johannes+Ben+K", 1, 2018, its2018))
	*bkgs = append(*bkgs, *booking.NewBooking(13, "AR", "", "", "BW", "Project-X", net, 357, "Anke+Johannes+Ben+K", 1, 2018, its2018))
	*bkgs = append(*bkgs, *booking.NewBooking(13, "ER", "", "", "JM", "Project-X", net, 11.9, "H-costs", 1, 2018, its2018))
	*bkgs = append(*bkgs, *booking.NewBooking(13, "ER", "", "", "BW", "Project-X", net, 11.9, "H-costs", 1, 2018, its2018))

	if (DEBUG) {log.Println("in TestStakeholderYearlyIncome: ", len(as.GetCollectiveAccount_thisYear().Bookings))}

	// nun verteilen
	for _, p := range as.GetCollectiveAccount_thisYear().Bookings {
		if (DEBUG) {log.Println("   processing ", p)}
		Process(as, p)
	}
	ErloesverteilungAnEmployees(as)
	// write out all stakehoders accounts
	for _, sh := range shrepo.All(2018) {
		acc,_ := as.GetSubacc("BW", accountSystem.UK_AnteileAuserloesen.Id )
		if (DEBUG) {log.Println("in TestErloesverteilungAnEmployees: ", sh.Id, acc)}
	}
	// ==> Ben sollte nun  70€ + 70€ + 5% von 300€ netto€ = 155€ Jahreseinkommen haben...
	assert.Equal(t, 155.0, math.Round(StakeholderYearlyIncome(as, "BW")) )
	assert.Equal(t, 0.0, math.Round(StakeholderYearlyIncome(as, "AN")) ) // because Anke is no employee...


	// test if there are booked revenues on k-account
	k_mainaccount, _ := as.Get("K")
//	for _,b := range k_mainaccount.Bookings {
//		log.Println("     b:", b.Amount, b.Text)
//	}
	k_subacc_Erloese, _ := as.GetSubacc("K", accountSystem.UK_AnteileAuserloesen.Id)
	assert.Equal(t, -600.0, k_mainaccount.Saldo )
	assert.Equal(t, -0.0, k_subacc_Erloese.Saldo )


	return
}

