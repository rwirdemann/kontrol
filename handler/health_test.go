package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ahojsenn/kontrol/account"
	"github.com/ahojsenn/kontrol/accountSystem"
	"github.com/ahojsenn/kontrol/booking"
	"github.com/stretchr/testify/assert"
)

//var router *mux.Router

// var as accountSystem.AccountSystem

func init() {
	repository := accountSystem.NewDefaultAccountSystem()
	repository.Add(account.NewAccount(account.AccountDescription{Id: "JM", Name: "k: Ralf", Type: "partner"}))
	repository.Add(account.NewAccount(account.AccountDescription{Id: "RW", Name: "k: Ralf", Type: "partner"}))
	repository.Add(account.NewAccount(account.AccountDescription{Id: "BW", Name: "k: Ben", Type: "partner"}))
	repository.Add(account.NewAccount(account.AccountDescription{Id: "AN", Name: "k: Anke Nehrenberg", Type: "partner"}))
	repository.Add(account.NewAccount(account.AccountDescription{Id: "SKR03_ErgebnisNachSteuern", Name: "SKR03_ErgebnisNachSteuern", Type: "Verrechnungskonto"}))

	k := account.NewAccount(account.AccountDescription{Id: "K", Name: "k: Kommitment", Type: "company"})
	ar := booking.NewBooking(13, "AR", "800", "1337", "JM", "Project-X", nil, 2000, "Rechnung WLW", 1, 2018, time.Time{})
	ar.CostCenter = "BW"
	k.Book(*ar)
	ar2 := booking.NewBooking(13, "AR", "800", "1337", "JM", "Project-X", nil, 2400, "Rechnung JH", 1, 2018, time.Time{})
	ar2.CostCenter = "RW"
	k.Book(*ar2)
	repository.Add(k)

	ar2 = booking.NewBooking(13, "SKR03", "4900", "977", "JM", "Project-X", nil, 2400, "Rückstellung Schornsteinfegerrechnung", 1, 2018, time.Time{})
	ar2.CostCenter = "RW"
	k.Book(*ar2)
	repository.Add(k)

	as = repository

	router = NewRouter("githash", "buildtime", repository)
}

func TestHealth(t *testing.T) {
	req, _ := http.NewRequest("GET", "/kontrol/health", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestHealthContent(t *testing.T) {
	req, _ := http.NewRequest("GET", "/kontrol/health", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	log.Println("in TestHealthContents: ", rr)

	b := []byte(rr.Body.String())

	var f interface{}
	err := json.Unmarshal(b, &f)
	if err != nil {
		log.Fatal(err)
	}
	m := f.(map[string]interface{})
	log.Println("in TestHealthContents: ", m)

	acc0 := m["Health"]
	assert.Equal(t, "OK", acc0)

}
