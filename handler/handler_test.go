package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"bitbucket.org/rwirdemann/kontrol/account"
	"bitbucket.org/rwirdemann/kontrol/domain"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

var router *mux.Router

func init() {
	repository := account.NewDefaultRepository()
	repository.Add(domain.NewAccount(domain.StakeholderAN))
	router = NewRouter("githash", "buildtime", repository)
}

func TestGetAllAccounts(t *testing.T) {
	req, _ := http.NewRequest("GET", "/kontrol/accounts", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	expected := `{"Accounts":[{"Owner":{"Id":"AN","Name":"Anke Nehrenberg","Type":"partner"},"Bookings":null,"Saldo":0}]}`
	assert.Equal(t, expected, rr.Body.String())
}
