package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ahojsenn/kontrol/account"
	"github.com/ahojsenn/kontrol/booking"
		"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/ahojsenn/kontrol/accountSystem"
)

var router *mux.Router

func init() {
	repository := accountSystem.EmptyDefaultAccountSystem(2016)
	repository.Add(account.NewAccount(account.AccountDescription{Id: "AN", Name: "k: Anke Nehrenberg", Type: "partner"}))

	k := account.NewAccount(account.AccountDescription{Id: "K", Name: "k: Kommitment", Type: "company"})
	ar := booking.NewBooking("AR", "800", "1337", "JM", nil, 2000, "Rechnung WLW", 1, 2018, time.Time{})
	ar.CostCenter = "BW"
	k.Book(*ar)
	ar2 := booking.NewBooking("AR", "800", "1337", "JM", nil, 2400, "Rechnung JH", 1, 2018, time.Time{})
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
	expected := "{\"Accounts\":[{\"Description\":{\"Id\":\"AN\",\"Name\":\"k: Anke Nehrenberg\",\"Type\":\"partner\"},\"Costs\":0,\"Advances\":0,\"Reserves\":0,\"Provision\":0,\"Revenue\":0,\"Taxes\":0,\"Internals\":0,\"Saldo\":0},{\"Description\":{\"Id\":\"K\",\"Name\":\"k: Kommitment\",\"Type\":\"company\"},\"Costs\":0,\"Advances\":0,\"Reserves\":0,\"Provision\":0,\"Revenue\":0,\"Taxes\":0,\"Internals\":0,\"Saldo\":4400}]}"
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
	expected := "{\"Description\":{\"Id\":\"K\",\"Name\":\"k: Kommitment\",\"Type\":\"company\"},\"Bookings\":[{\"Type\":\"\",\"Soll\":\"800\",\"Haben\":\"1337\",\"CostCenter\":\"BW\",\"Amount\":2000,\"Text\":\"Rechnung WLW\",\"Year\":2018,\"Month\":1,\"FileCreated\":\"0001-01-01T00:00:00Z\",\"BankCreated\":\"0001-01-01T00:00:00Z\"}],\"Costs\":0,\"Advances\":0,\"Reserves\":0,\"Provision\":0,\"Revenue\":0,\"Taxes\":0,\"Internals\":0,\"Saldo\":2000}"
	assert.Equal(t, expected, rr.Body.String())
}
