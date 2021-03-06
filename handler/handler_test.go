package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"bitbucket.org/rwirdemann/kontrol/account"
	"bitbucket.org/rwirdemann/kontrol/booking"
	"bitbucket.org/rwirdemann/kontrol/owner"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

var router *mux.Router

func init() {
	repository := account.EmptyDefaultRepository()
	repository.Add(account.NewAccount(owner.StakeholderAN))

	k := account.NewAccount(owner.StakeholderKM)
	ar := booking.NewBooking("AR", "JM", nil, 2000, "Rechnung WLW", 1, 2018, time.Time{})
	ar.CostCenter = "BW"
	k.Book(*ar)
	ar2 := booking.NewBooking("AR", "JM", nil, 2400, "Rechnung JH", 1, 2018, time.Time{})
	ar2.CostCenter = "RW"
	k.Book(*ar2)
	repository.Add(k)

	router = NewRouter("githash", "buildtime", repository)
}

func TestGetAllAccounts(t *testing.T) {
	req, _ := http.NewRequest("GET", "/kontrol/accounts", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	expected := "{\"Accounts\":[{\"Owner\":{\"Id\":\"AN\",\"Name\":\"Anke Nehrenberg\",\"Type\":\"partner\"},\"Saldo\":0},{\"Owner\":{\"Id\":\"K\",\"Name\":\"Kommitment\",\"Type\":\"company\"},\"Saldo\":4400}]}"
	assert.Equal(t, expected, rr.Body.String())
}

func TestGetBankAccount(t *testing.T) {
	req, _ := http.NewRequest("GET", "/kontrol/bankaccount", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestGetAccountFilterByCostcenter(t *testing.T) {
	req, _ := http.NewRequest("GET", "/kontrol/accounts/K?cs=BW", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	expected := "{\"Owner\":{\"Id\":\"K\",\"Name\":\"Kommitment\",\"Type\":\"company\"},\"Bookings\":[{\"Type\":\"\",\"CostCenter\":\"BW\",\"Amount\":2000,\"Text\":\"Rechnung WLW\",\"Year\":2018,\"Month\":1,\"FileCreated\":\"0001-01-01T00:00:00Z\"}],\"Saldo\":2000}"
	assert.Equal(t, expected, rr.Body.String())
}
