package processing

import (
	"log"
	"testing"
	"time"

	"github.com/ahojsenn/kontrol/account"
	"github.com/ahojsenn/kontrol/accountSystem"
	"github.com/ahojsenn/kontrol/booking"
	"github.com/ahojsenn/kontrol/util"
	"github.com/ahojsenn/kontrol/valueMagnets"
	"github.com/stretchr/testify/assert"
)

func setCCTestUp() {
	accSystem = accountSystem.NewDefaultAccountSystem()
	accountBank = accSystem.GetCollectiveAccount_thisYear()
	shrepo := valueMagnets.Stakeholder{}
	accountHannes, _ = accSystem.Get(shrepo.Get("JM").Id)
	accountRalf, _ = accSystem.Get(shrepo.Get("RW").Id)
	accountKommitment, _ = accSystem.Get(shrepo.Get("K").Id)
}

func TestBookCostToCostCenter(t *testing.T) {

	// given an accountingSystem and a booking
	accsystem := accountSystem.NewDefaultAccountSystem()
	bkng := booking.Booking{}
	bkng.Type = "SKR03"
	bkng.Amount = 1337.23
	bkng.Year = 2017
	bkng.Month = 7
	bkng.Haben = "4200"
	bkng.Soll = "1200" //
	bkng.CostCenter = "JM"
	bkng.Text = "This is a test"

	// when you book it
	//BookSKR03Command{AccSystem: accsystem, Booking: bkng}.run()
	BookCostToCostCenter{AccSystem: accsystem, Booking: bkng}.run()

	// there is money on the costcenter JM
	account, _ := accsystem.GetSubacc("JM", accountSystem.UK_Kosten.Id)

	util.AssertEquals(t, 1, len(account.Bookings))
	assert.Equal(t, 1337.23, account.Bookings[0].Amount)

	// this acfcount has two bookings, since it is jus a passing through accound ...
	account2, _ := accsystem.Get("JM")
	util.AssertEquals(t, 2, len(account2.Bookings))
	assert.Equal(t, 1337.23, account2.Bookings[0].Amount)

}

func TestBookRevenueToEmployeeCostCenter(t *testing.T) {

	// given an accountingSystem and a booking
	accsystem := accountSystem.NewDefaultAccountSystem()
	net := make(map[string]float64)
	shrepo := valueMagnets.Stakeholder{}
	net[shrepo.Get("BW").Id] = 1000
	net[shrepo.Get("JM").Id] = 3675.0
	bkng := *booking.NewBooking(13, "AR", "", "", "JM", "Project-X", net, 17225.25, "Rechnung 1234", 1, 2017, time.Time{})

	// when you book it
	BookAusgangsrechnungCommand{AccSystem: accsystem, Booking: bkng}.run()
	BookRevenueToEmployeeCostCenter{AccSystem: accsystem, Booking: bkng}.run()

	// there is money on the costcenter BW --> now on subaccount
	acc, _ := accsystem.GetSubacc("BW", accountSystem.UK_AnteileAuserloesen.Id)
	util.AssertEquals(t, 1, len(acc.Bookings))
	assert.Equal(t, 1000*account.EmployeeShare, acc.Bookings[0].Amount)
}

func TestExternNettoAnteil(t *testing.T) {
	setCCTestUp()

	// given: a booking
	net := map[string]float64{
		valueMagnets.StakeholderEX.Id: 10800.0,
	}
	its2018 := time.Date(2018, 1, 23, 0, 0, 0, 0, time.UTC)
	p := booking.NewBooking(13, "AR", "", "", "JM", "Project-X", net, 12852.0, "Rechnung 1234", 1, 2017, its2018)

	// when: the position is processed
	Process(accSystem, *p)
	BookRevenueToEmployeeCostCenter{AccSystem: accSystem, Booking: *p}.run()

	// and hannes got his provision <-- this is noch anymore booked here
	accountHannes, _ := accSystem.GetSubacc("JM", accountSystem.UK_Vertriebsprovision.Id)
	provision := accountHannes.Bookings[0]
	util.AssertFloatEquals(t, 10800.0*account.PartnerProvision, provision.Amount)
	util.AssertEquals(t, booking.CC_Vertriebsprovision, provision.Type)
	util.AssertEquals(t, 1, len(accountHannes.Bookings))

	// and kommitment got 95%
	acc, _ := accSystem.Get(valueMagnets.StakeholderKM.Id)
	bk := acc.Bookings[0]
	util.AssertFloatEquals(t, -10800.0, bk.Amount)
	util.AssertEquals(t, booking.CC_RevDistribution_1, bk.Type)
}

func TestStakeholderWithNetPositions(t *testing.T) {
	setCCTestUp()

	// given: a booking
	net := map[string]float64{
		valueMagnets.StakeholderEX.Id: 10800.0,
	}
	shrepo := valueMagnets.Stakeholder{}
	net[shrepo.Get("BW").Id] = 1000
	net[shrepo.Get("JM").Id] = 3675.0
	net[shrepo.Get("KR").Id] = 3675.0
	net[shrepo.Get("IK").Id] = 3675.0
	p := booking.NewBooking(13, "AR", "", "", "JM", "Project-X", net, 17225.25, "Rechnung 1234", 1, 2017, time.Time{})

	// when: the position is processed
	Process(accSystem, *p)

	// make shure alle benefitees are recognized
	benefitees := BookRevenueToEmployeeCostCenter{AccSystem: accSystem, Booking: *p}.stakeholderWithNetPositions()
	log.Println("in TestStakeholderWithNetPositions:", benefitees)
	util.AssertEquals(t, 5, len(benefitees))
}
