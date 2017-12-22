package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"bitbucket.org/rwirdemann/kontrol/account"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

var router *mux.Router

func init() {
	repository := account.DefaultRepository{}
	router = NewRouter("githash", "buildtime", &repository)
}

func TestGetAllAccounts(t *testing.T) {

	req, _ := http.NewRequest("GET", "/kontrol/accounts", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	expected := `{"Accounts":[{"Owner":{"Id":"AN","Name":"Anke Nehrenberg","Type":"partner"},"Bookings":null,"Saldo":0},{"Owner":{"Id":"BW","Name":"Ben Wiedenmann","Type":"employee"},"Bookings":null,"Saldo":0},{"Owner":{"Id":"JM","Name":"Johannes Mainusch","Type":"partner"},"Bookings":null,"Saldo":0},{"Owner":{"Id":"K","Name":"Kommitment","Type":"company"},"Bookings":null,"Saldo":0},{"Owner":{"Id":"RW","Name":"Ralf Wirdemann","Type":"partner"},"Bookings":null,"Saldo":0}]}`
	assert.Equal(t, expected, rr.Body.String())
}
