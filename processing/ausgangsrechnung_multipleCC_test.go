package processing

import (
	"testing"
	"time"

	"github.com/ahojsenn/kontrol/accountSystem"
	"github.com/ahojsenn/kontrol/booking"
	"github.com/ahojsenn/kontrol/util"
	"github.com/ahojsenn/kontrol/valueMagnets"
)

func TestMultipleCostCenters(t *testing.T) {

	var as accountSystem.AccountSystem
	util.Global.FinancialYear = 2017
	as = accountSystem.NewDefaultAccountSystem()

	net := make(map[string]float64)
	shrepo := valueMagnets.Stakeholder{}

	// given the following booking of 1190
	net[shrepo.Get("AN").Id] = 500.0
	net[shrepo.Get("JM").Id] = 500.0
	net[shrepo.Get("Rest").Id] = 190.0

	bkng := booking.NewBooking(13, "AR", "", "", "BW,JM,AN,blupp", "Project-X", net, 1119.0, "Rechnung 1234", 1, 2017, time.Time{})

	// when: the position is processed
	Process(as, *bkng)
	BookRevenueToEmployeeCostCenter{AccSystem: as, Booking: *bkng}.run()

	// 1/4 of of 5% provision = 12,5€ goes to each of the four parties
	acc, _ := as.GetSubacc("BW", accountSystem.UK_Vertriebsprovision.Id)
	//	log.Println("in TestMultipleCostCenters", acc.Saldo)
	util.AssertFloatEquals(t, 12.5, acc.Saldo)

}
