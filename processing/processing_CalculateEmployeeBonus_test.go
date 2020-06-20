package processing

import (
	"github.com/ahojsenn/kontrol/accountSystem"
	"github.com/ahojsenn/kontrol/booking"
	"github.com/ahojsenn/kontrol/valueMagnets"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCalculateEmplyeeBonnusses (t *testing.T) {

	its2018 := time.Date(2018, 1, 23, 0, 0, 0, 0, time.UTC)

	valueMagnets.KommimtmentYear{}.Init(2018)
	valueMagnets.StakeholderRepository = append(valueMagnets.StakeholderRepository,
		valueMagnets.Stakeholder{"AB", "Anna Blume", "Employee", "1.0", "", 0})
	valueMagnets.StakeholderRepository = append(valueMagnets.StakeholderRepository,
		valueMagnets.Stakeholder{"JM", "Johannes Mainusch", "Kommanditist", "1.0", "0.5", 0})
	valueMagnets.StakeholderRepository = append(valueMagnets.StakeholderRepository,
		valueMagnets.Stakeholder{"K", "Kompanie", "Company", "0", "0", 0})

	as := accountSystem.NewDefaultAccountSystem()
	stakeholder := valueMagnets.Stakeholder{}
	net := make(map[valueMagnets.Stakeholder]float64)
	net[stakeholder.Get("AB")] = 100.0
	net[stakeholder.Get("JM")] = 100.0

	hauptbuch := as.GetCollectiveAccount_thisYear(2018)
//	b1 := *booking.NewBooking(13, "AR", "", "", "K", "Project-X", net, 1190, "Anna+Johannes", 1, 2018, its2018)
	Process(as, *booking.NewBooking(13, "AR", "", "", "K", "Project-X", net, 238, "Anna+Johannes", 1, 2018, its2018))
	Process(as, *booking.NewBooking(13, "AR", "", "", "AB", "Project-X", net, 238, "Anna+Johannes", 1, 2018, its2018))
	Process(as, *booking.NewBooking(13, "AR", "", "", "AB", "Project-X", net, 238, "Anna+Johannes", 1, 2018, its2018))
	Process(as, *booking.NewBooking(13, "ER", "", "", "AB", "Project-X", net, 11.9, "Anna+Johannes", 1, 2018, its2018))

	// nun verteilen
	for _, p := range hauptbuch.Bookings {
		Process(as, p)
	}

	// now distribution of costs & profits
	Kostenerteilung(as)
	ErloesverteilungAnEmployees(as)
	CalculateEmployeeBonus(as)
	CalculateEmployeeBonus(as)
	CalculateEmployeeBonus(as)  // calling this twice should not double the bonus...

	// 70% of 100€
	gehaelterAccount, _ := as.Get("4100_4199")
	annasAccount_Erloese, _ := as.GetSubacc("AB", accountSystem.UK_AnteileAuserloesen.Id)

	assert.Equal(t, 210.0, annasAccount_Erloese.Saldo )
	assert.Equal(t, -220.0, gehaelterAccount.Saldo )

	return
}

