package processing

import (
	"github.com/ahojsenn/kontrol/account"
	"github.com/ahojsenn/kontrol/accountSystem"
	"github.com/ahojsenn/kontrol/booking"
	"github.com/ahojsenn/kontrol/util"
	"github.com/ahojsenn/kontrol/valueMagnets"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func setCCTestUp() {
	accSystem = accountSystem.NewDefaultAccountSystem()
	accountBank = accSystem.GetCollectiveAccount()
	accountHannes, _ = accSystem.Get(valueMagnets.StakeholderRepository{}.Get("JM").Id)
	accountRalf, _ = accSystem.Get(valueMagnets.StakeholderRepository{}.Get("RW").Id)
	accountKommitment, _ = accSystem.Get(valueMagnets.StakeholderRepository{}.Get("K").Id)
}


func TestBookCostToCostCenter(t *testing.T) {

	// given an accountingSystem and a booking
	accsystem := accountSystem.NewDefaultAccountSystem()
	bkng := booking.Booking{}
	bkng.Type = "SKR03"
	bkng.Amount =  1337.23
	bkng.Year = 2017
	bkng.Month = 7
	bkng.Haben = "4200"
	bkng.Soll = "1200"  //
	bkng.CostCenter = "JM"
	bkng.Text = "This is a test"

	// when you book it
	BookSKR03Command{AccSystem: accsystem, Booking: bkng}.run()
	BookCostToCostCenter{AccSystem: accsystem, Booking: bkng}.run()

	// there is money on the costcenter JM
	account, _ := accsystem.Get("JM")
	util.AssertEquals(t, 1, len(account.Bookings))
	assert.Equal(t, 1337.23, account.Bookings[0].Amount)

	// there is something on the other account too
	account2, _ := accsystem.Get(accountSystem.AlleKLRBuchungen.Id)
	util.AssertEquals(t, 1, len(account2.Bookings))
	assert.Equal(t, -1337.23, account2.Bookings[0].Amount)

}

func TestBookRevenueToEmployeeCostCenter(t *testing.T) {

	// given an accountingSystem and a booking
	accsystem := accountSystem.NewDefaultAccountSystem()
	net := make(map[valueMagnets.Stakeholder]float64)
	net[valueMagnets.StakeholderRepository{}.Get("BW")] = 1000
	net[valueMagnets.StakeholderRepository{}.Get("JM")] = 3675.0
	bkng := *booking.NewBooking(13, "AR", "", "", "JM", net, 17225.25, "Rechnung 1234", 1, 2017, time.Time{})


	// when you book it
	BookAusgangsrechnungCommand{AccSystem: accsystem, Booking: bkng}.run()
	BookRevenueToEmployeeCostCenter{AccSystem: accsystem, Booking: bkng}.run()

	// there is money on the costcenter JM
	acc, _ := accsystem.Get("BW")
	util.AssertEquals(t, 1, len(acc.Bookings))
	assert.Equal(t, 1000*  account.EmployeeShare, acc.Bookings[0].Amount)


}



func TestExternNettoAnteil(t *testing.T) {
	setCCTestUp()

	// given: a booking
	net := map[valueMagnets.Stakeholder]float64{
		valueMagnets.StakeholderEX: 10800.0,
	}
	its2018 := time.Date(2018, 1, 23, 0, 0, 0, 0, time.UTC)
	p := booking.NewBooking(13, "AR", "", "", "JM", net, 12852.0, "Rechnung 1234", 1, 2017, its2018)

	// when: the position is processed
	Process(accSystem, *p)
	BookRevenueToEmployeeCostCenter{AccSystem: accSystem, Booking: *p}.run()

	// and hannes got his provision
	accountHannes, _ := accSystem.Get(valueMagnets.StakeholderRepository{}.Get("JM").Id)
	provision := accountHannes.Bookings[0]
	util.AssertFloatEquals(t, 10800.0*account.PartnerProvision, provision.Amount)
	util.AssertEquals(t, booking.CC_Vertriebsprovision, provision.Type)

	// and kommitment got 95%
	util.AssertEquals(t, 1, len(accountHannes.Bookings))
	acc, _ := accSystem.Get(valueMagnets.StakeholderKM.Id)
	bk := acc.Bookings[0]
	util.AssertFloatEquals(t, 10800.0*account.KommmitmentExternShare, bk.Amount)
	util.AssertEquals(t, booking.CC_KommitmentanteilEX, bk.Type)
}
