package processing

import (
	"testing"
	"github.com/ahojsenn/kontrol/accountSystem"
	"github.com/ahojsenn/kontrol/booking"
	"github.com/ahojsenn/kontrol/util"
		"github.com/stretchr/testify/assert"
	)

func TestSKR03Command(t *testing.T) {

	// iven an accountingSystem and a booking
	accsystem := accountSystem.NewDefaultAccountSystem()
	bkng := booking.Booking{}
	bkng.Type = "SKR03"
	bkng.Amount =  1337.23
	bkng.Year = 2017
	bkng.Month = 7
	bkng.Haben = "1200"
	bkng.Soll = "4200"
	bkng.CostCenter = "JM"
	bkng.Text = "This is a test"

	// when you book it
	command := BookSKR03Command{AccSystem: accsystem, Booking: bkng}
	command.run()

	// there is money on the bank account
	account, _ := accsystem.Get(accountSystem.SKR03_1200.Id)
	util.AssertEquals(t, 1, len(account.Bookings))
	assert.Equal(t, 1337.23, account.Bookings[0].Amount)

	// there is something on the other account too
	account2, _ := accsystem.Get(accountSystem.SKR03_sonstigeAufwendungen.Id)
	util.AssertEquals(t, 1, len(account2.Bookings))
	assert.Equal(t, -1337.23, account2.Bookings[0].Amount)

}

