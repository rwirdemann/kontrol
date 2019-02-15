package handler

import (
	"encoding/json"
	"github.com/ahojsenn/kontrol/processing"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ahojsenn/kontrol/account"
	"github.com/ahojsenn/kontrol/accountSystem"
	"github.com/ahojsenn/kontrol/booking"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

var router *mux.Router

var as accountSystem.AccountSystem

func init() {
	repository := accountSystem.EmptyDefaultAccountSystem()
	repository.Add(account.NewAccount(account.AccountDescription{Id: "JM", Name: "k: Ralf", Type: "partner"}))
	repository.Add(account.NewAccount(account.AccountDescription{Id: "RW", Name: "k: Ralf", Type: "partner"}))
	repository.Add(account.NewAccount(account.AccountDescription{Id: "BW", Name: "k: Ben", Type: "partner"}))
	repository.Add(account.NewAccount(account.AccountDescription{Id: "AN", Name: "k: Anke Nehrenberg", Type: "partner"}))
	repository.Add(account.NewAccount(account.AccountDescription{Id: "SKR03_ErgebnisNachSteuern", Name: "SKR03_ErgebnisNachSteuern", Type: "Verrechnungskonto"}))

	k := account.NewAccount(account.AccountDescription{Id: "K", Name: "k: Kommitment", Type: "company"})
	ar := booking.NewBooking(13,"AR", "800", "1337", "JM", "Project-X",nil, 2000, "Rechnung WLW", 1, 2018, time.Time{})
	ar.CostCenter = "BW"
	k.Book(*ar)
	ar2 := booking.NewBooking(13,"AR", "800", "1337", "JM", "Project-X",nil, 2400, "Rechnung JH", 1, 2018, time.Time{})
	ar2.CostCenter = "RW"
	k.Book(*ar2)
	repository.Add(k)

	ar2 = booking.NewBooking(13,"SKR03", "4900", "977", "JM", "Project-X",	nil, 2400, "Rückstellung Schornsteinfegerrechnung", 1, 2018, time.Time{})
	ar2.CostCenter = "RW"
	k.Book(*ar2)
	repository.Add(k)

	as = repository

	router = NewRouter("githash", "buildtime", repository)
}

func TestGetAllAccounts(t *testing.T) {
	req, _ := http.NewRequest("GET", "/kontrol/accounts", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	b := []byte(rr.Body.String())
	var f interface{}
	err := json.Unmarshal(b, &f)
	if err != nil {
		log.Fatal(err)
	}
	m := f.(map[string]interface{})

	acc0 := m["Accounts"].([]interface{})[0].
	(map[string]interface{})["Description"].
	(map[string]interface{})["Id"]
	assert.Equal(t, "SKR03_ErgebnisNachSteuern", acc0)

	acc2 := m["Accounts"].([]interface{})[2].
	(map[string]interface{})["Description"].
	(map[string]interface{})["Id"]
	assert.Equal(t, "BW", acc2)

	//	expected := "{\"Accounts\":[{\"Description\":{\"Id\":\"SKR03_ErgebnisNachSteuern\",\"Name\":\"SKR03_ErgebnisNachSteuern\",\"Type\":\"Verrechnungskonto\",\"Superaccount\":\"\"},\"KommitmenschNettoFaktura\":0,\"AnteilAusFaktura\":0,\"AnteilAusFairshares\":0,\"KommitmenschDarlehen\":0,\"Nbkngs\":0,\"Costs\":0,\"Advances\":0,\"Reserves\":0,\"Rest\":0,\"Revenue\":0,\"Salesprv\":0,\"Taxes\":0,\"Internals\":0,\"YearS\":0,\"Saldo\":0},{\"Description\":{\"Id\":\"AN\",\"Name\":\"k: Anke Nehrenberg\",\"Type\":\"partner\",\"Superaccount\":\"\"},\"KommitmenschNettoFaktura\":0,\"AnteilAusFaktura\":0,\"AnteilAusFairshares\":0,\"KommitmenschDarlehen\":0,\"Nbkngs\":0,\"Costs\":0,\"Advances\":0,\"Reserves\":0,\"Rest\":0,\"Revenue\":0,\"Salesprv\":0,\"Taxes\":0,\"Internals\":0,\"YearS\":0,\"Saldo\":0},{\"Description\":{\"Id\":\"BW\",\"Name\":\"k: Ben\",\"Type\":\"partner\",\"Superaccount\":\"\"},\"KommitmenschNettoFaktura\":0,\"AnteilAusFaktura\":0,\"AnteilAusFairshares\":0,\"KommitmenschDarlehen\":0,\"Nbkngs\":0,\"Costs\":0,\"Advances\":0,\"Reserves\":0,\"Rest\":0,\"Revenue\":0,\"Salesprv\":0,\"Taxes\":0,\"Internals\":0,\"YearS\":0,\"Saldo\":0},{\"Description\":{\"Id\":\"K\",\"Name\":\"k: Kommitment\",\"Type\":\"company\",\"Superaccount\":\"\"},\"KommitmenschNettoFaktura\":0,\"AnteilAusFaktura\":0,\"AnteilAusFairshares\":0,\"KommitmenschDarlehen\":0,\"Nbkngs\":3,\"Costs\":0,\"Advances\":0,\"Reserves\":0,\"Rest\":4400,\"Revenue\":0,\"Salesprv\":0,\"Taxes\":0,\"Internals\":0,\"YearS\":0,\"Saldo\":6800},{\"Description\":{\"Id\":\"JM\",\"Name\":\"k: Ralf\",\"Type\":\"partner\",\"Superaccount\":\"\"},\"KommitmenschNettoFaktura\":0,\"AnteilAusFaktura\":0,\"AnteilAusFairshares\":0,\"KommitmenschDarlehen\":0,\"Nbkngs\":0,\"Costs\":0,\"Advances\":0,\"Reserves\":0,\"Rest\":0,\"Revenue\":0,\"Salesprv\":0,\"Taxes\":0,\"Internals\":0,\"YearS\":0,\"Saldo\":0},{\"Description\":{\"Id\":\"RW\",\"Name\":\"k: Ralf\",\"Type\":\"partner\",\"Superaccount\":\"\"},\"KommitmenschNettoFaktura\":0,\"AnteilAusFaktura\":0,\"AnteilAusFairshares\":0,\"KommitmenschDarlehen\":0,\"Nbkngs\":0,\"Costs\":0,\"Advances\":0,\"Reserves\":0,\"Rest\":0,\"Revenue\":0,\"Salesprv\":0,\"Taxes\":0,\"Internals\":0,\"YearS\":0,\"Saldo\":0}]}"
}

func TestGetCollectiveAccount(t *testing.T) {
	req, _ := http.NewRequest("GET", "/kontrol/collectiveaccount", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestGetAccountFilterByCostcenter(t *testing.T) {
	req, _ := http.NewRequest("GET", "/kontrol/accounts/K?cs=BW", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	b := []byte(rr.Body.String())
	var f interface{}
	err := json.Unmarshal(b, &f)
	if err != nil {
		log.Fatal(err)
	}
	m := f.(map[string]interface{})

	acc0 := m["Description"].(map[string]interface{})["Id"]
	assert.Equal(t, "K", acc0)

	acc2 := m["Bookings"].([]interface{})[0].
	(map[string]interface{})["RowNr"]
	assert.Equal(t, float64(13), acc2)

	//expected := "{\"Description\":{\"Id\":\"K\",\"Name\":\"k: Kommitment\",\"Type\":\"company\",\"Superaccount\":\"\"},\"Bookings\":[{\"RowNr\":13,\"Type\":\"AR\",\"Soll\":\"800\",\"Haben\":\"1337\",\"CostCenter\":\"BW\",\"Project\":\"Project-X\",\"Amount\":2000,\"Text\":\"Rechnung WLW\",\"Year\":2018,\"Month\":1,\"FileCreated\":\"0001-01-01T00:00:00Z\",\"BankCreated\":\"0001-01-01T00:00:00Z\"}],\"KommitmenschNettoFaktura\":0,\"AnteilAusFaktura\":0,\"AnteilAusFairshares\":0,\"KommitmenschDarlehen\":0,\"Nbkngs\":1,\"Costs\":0,\"Advances\":0,\"Reserves\":0,\"Rest\":2000,\"Revenue\":0,\"Salesprv\":0,\"Taxes\":0,\"Internals\":0,\"YearS\":0,\"Saldo\":2000}"

}

func TestGetAccountsBilanzkonten (t *testing.T) {
	as.Add(account.NewAccount(account.AccountDescription{Id: "1400", Name: "SKR03_1400_OPOS-Kunde", Type: "Aktivkonto"}))

	req, _ := http.NewRequest("GET", "/kontrol/bilanz", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	b := []byte(rr.Body.String())
	var f interface{}
	err := json.Unmarshal(b, &f)
	if err != nil {
		log.Fatal(err)
	}
	m := f.(map[string]interface{})

	acc0 := m["Accounts"].([]interface{})[0].
	(map[string]interface{})["Description"].
	(map[string]interface{})["Id"]
	assert.Equal(t, "1400", acc0)

	acc2 := m["Accounts"].([]interface{})[0].
	(map[string]interface{})["Description"].
	(map[string]interface{})["Name"]
	assert.Equal(t, "SKR03_1400_OPOS-Kunde", acc2)

	// expected := "{\"Accounts\":[{\"Description\":{\"Id\":\"1400\",\"Name\":\"SKR03_1400_OPOS-Kunde\",\"Type\":\"Aktivkonto\",\"Superaccount\":\"\"},\"KommitmenschNettoFaktura\":0,\"AnteilAusFaktura\":0,\"AnteilAusFairshares\":0,\"KommitmenschDarlehen\":0,\"Nbkngs\":0,\"Costs\":0,\"Advances\":0,\"Reserves\":0,\"Rest\":0,\"Revenue\":0,\"Salesprv\":0,\"Taxes\":0,\"Internals\":0,\"YearS\":0,\"Saldo\":0}]}"
}

func TestGetAccountsGuV(t *testing.T) {
	as.Add(account.NewAccount(account.AccountDescription{Id: "8100", Name: "Ertrag", Type: "Ertragskonto"}))

	req, _ := http.NewRequest("GET", "/kontrol/GuV", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, http.StatusOK, rr.Code)
	b := []byte(rr.Body.String())
	var f interface{}
	err := json.Unmarshal(b, &f)
	if err != nil {
		log.Fatal(err)
	}
	m := f.(map[string]interface{})

	acc0 := m["Accounts"].([]interface{})[0].
	(map[string]interface{})["Description"].
	(map[string]interface{})["Id"]
	assert.Equal(t, "8100", acc0)

	acc2 := m["Accounts"].([]interface{})[0].
	(map[string]interface{})["Description"].
	(map[string]interface{})["Name"]
	assert.Equal(t, "Ertrag", acc2)
	//expected := "{\"Accounts\":[{\"Description\":{\"Id\":\"8100\",\"Name\":\"Ertrag\",\"Type\":\"Ertragskonto\",\"Superaccount\":\"\"},\"KommitmenschNettoFaktura\":0,\"AnteilAusFaktura\":0,\"AnteilAusFairshares\":0,\"KommitmenschDarlehen\":0,\"Nbkngs\":0,\"Costs\":0,\"Advances\":0,\"Reserves\":0,\"Rest\":0,\"Revenue\":0,\"Salesprv\":0,\"Taxes\":0,\"Internals\":0,\"YearS\":0,\"Saldo\":0}]}"

}

func TestGetAccountsProjects(t *testing.T) {
	ac := account.NewAccount(account.AccountDescription{Id: "TestProjekt", Name: "k: KLR Kommitment", Type: account.KontenartProject})
	bk := booking.NewBooking(13,"AR", "", "", "JM", "Project-X",nil, 2000, "Rechnung WLW", 1, 2018, time.Time{})
	ac.Book(*bk)
	//( and a booking which is not supposed to be picked up...
	ac2 := account.NewAccount(account.AccountDescription{Id: "somethingelse", Name: "k: KLR Kommitment", Type: account.KontenartErtrag})
	bk2 := booking.NewBooking(13,"SKR03", "4900", "977", "JM", "Project-X",	nil, 2400, "Rückstellung Schornsteinfegerrechnung", 1, 2018, time.Time{})
	ac2.Book(*bk2)

	as.Add(ac)


	req, _ := http.NewRequest("GET", "/kontrol/projects", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	processing.GenerateProjectControlling(as)
	assert.Equal(t, http.StatusOK, rr.Code)

	b := []byte(rr.Body.String())
	var f interface{}
	err := json.Unmarshal(b, &f)
	if err != nil {
		log.Fatal(err)
	}
	m := f.(map[string]interface{})

	acc0 := m["Accounts"].([]interface{})[0].
	(map[string]interface{})["Description"].
	(map[string]interface{})["Id"]
	assert.Equal(t, "TestProjekt", acc0)

	acc2 := m["Accounts"].([]interface{})[0].
	(map[string]interface{})["Bookings"].
	([]interface{})[0].
	(map[string]interface{})["RowNr"]
	assert.Equal(t, float64(13), acc2)



	// only Bilanzkonten are wanted here [Type Passivkonto or Aktivkonto]
//	expected := "{\"Accounts\":[{\"Description\":{\"Id\":\"TestProjekt\",\"Name\":\"k: KLR Kommitment\",\"Type\":\"VerrechnungskontoProjekt\",\"Superaccount\":\"\"},\"Bookings\":[{\"RowNr\":13,\"Type\":\"AR\",\"Soll\":\"\",\"Haben\":\"\",\"CostCenter\":\"JM\",\"Project\":\"Project-X\",\"Amount\":2000,\"Text\":\"Rechnung WLW\",\"Year\":2018,\"Month\":1,\"FileCreated\":\"0001-01-01T00:00:00Z\",\"BankCreated\":\"0001-01-01T00:00:00Z\"}],\"KommitmenschNettoFaktura\":0,\"AnteilAusFaktura\":0,\"AnteilAusFairshares\":0,\"KommitmenschDarlehen\":0,\"Nbkngs\":1,\"Costs\":0,\"Advances\":0,\"Reserves\":0,\"Rest\":2000,\"Revenue\":0,\"Salesprv\":0,\"Taxes\":0,\"Internals\":0,\"YearS\":0,\"Saldo\":2000}]}"

}
