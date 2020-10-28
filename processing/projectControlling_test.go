package processing

import (
	"testing"
	"time"

	"github.com/ahojsenn/kontrol/accountSystem"
	"github.com/ahojsenn/kontrol/booking"
	"github.com/ahojsenn/kontrol/util"
	"github.com/ahojsenn/kontrol/valueMagnets"
)

func TestGenerateProjectControlling(t *testing.T) {

	var as accountSystem.AccountSystem
	util.Global.FinancialYear = 2017
	as = accountSystem.NewDefaultAccountSystem()

	net := make(map[string]float64)
	shrepo := valueMagnets.Stakeholder{}

	// given the following booking of 1190
	net[shrepo.Get("AN").Id] = 500.0
	net[shrepo.Get("JM").Id] = 500.0
	net[shrepo.Get("RR").Id] = 190.0

	bkng := booking.NewBooking(13, "AR", "", "", "BW,JM,AN,blupp", "Project-Z", net, 17225.25, "Rechnung 1234", 1, 2017, time.Time{})

	// when: the position is processed
	Process(as, *bkng)
	BookRevenueToEmployeeCostCenter{AccSystem: as, Booking: *bkng}.run()
	GenerateProjectControlling(as)

	acc, _ := as.Get("_PROJEKT-Project-Z")
	//	log.Println("in TestMultipleCostCenters", acc)
	util.AssertFloatEquals(t, 14475, acc.Saldo)
}
