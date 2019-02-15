package processing

import (
	"time"
	"github.com/ahojsenn/kontrol/booking"
	"github.com/ahojsenn/kontrol/account"
	"github.com/ahojsenn/kontrol/accountSystem"
	"github.com/ahojsenn/kontrol/valueMagnets"
	"github.com/stretchr/testify/suite"
	"testing"
)

type AusgangsrechnungTestSuite struct {
	suite.Suite
	repository        accountSystem.AccountSystem
	accountBank       *account.Account
	accountHannes     *account.Account
	accountRalf       *account.Account
	accountBen        *account.Account
	accountKommitment *account.Account
}

func (suite *AusgangsrechnungTestSuite) SetupTest() {
	suite.repository = accountSystem.NewDefaultAccountSystem()
	suite.accountBank = suite.repository.GetCollectiveAccount()
	shrepo := valueMagnets.Stakeholder{}
	suite.accountRalf, _ = suite.repository.Get(  shrepo.Get("RW").Id )
	suite.accountHannes, _ = suite.repository.Get(  shrepo.Get("JM").Id )
	suite.accountBen, _ = suite.repository.Get(  shrepo.Get("BW").Id )
	suite.accountKommitment, _ = suite.repository.Get( shrepo.Get("K").Id)
}

func TestAusgangsRechnungTestSuite(t *testing.T) {
	suite.Run(t, new(AusgangsrechnungTestSuite))
}



//
// Tests für Offene Posten
//
func (suite *AusgangsrechnungTestSuite) TestOffeneRechnung() {

	// given: a booking with empty timestamp in position "BankCreated"
	net := make(map[valueMagnets.Stakeholder]float64)
	shrepo := valueMagnets.Stakeholder{}
	net[shrepo.Get("RW")] = 10800.0
	net[shrepo.Get("JM")] = 3675.0
	p := booking.NewBooking(13, "AR", "", "", "JM", "Project-X", net, 17225.25, "Rechnung 1234", 1, 2017, time.Time{})

	// when: the position is processed
	Process(suite.repository, *p)

	// then the booking is on SKR03_1400
	account1400, _ := suite.repository.Get(accountSystem.SKR03_1400.Id)
	bookings1400 := account1400.Bookings
	suite.Equal(1, len(bookings1400))
}




func (suite *AusgangsrechnungTestSuite) assertBooking(amount float64, _type string, b *booking.Booking) {
	suite.Equal(amount, b.Amount)
	suite.Equal(_type, b.Type)
}

func Round(x, unit float64) float64 {
	return float64(int64(x/unit+0.5)) * unit
}