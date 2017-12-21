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

	// And: Body contains 1 product
	//	expected := `[{"Id":"1","Name":"Schuhe","Description":"","Category":"","Price":0}]`
	//	assert.Equal(t, expected, rr.Body.String())
}
