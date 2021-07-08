package handler

import (
	"encoding/json"
	"github.com/ahojsenn/kontrol/account"
	"github.com/ahojsenn/kontrol/accountSystem"
	"github.com/ahojsenn/kontrol/util"
	"github.com/ahojsenn/kontrol/valueMagnets"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)


func TestMakeGetKommitmenschenAccountsHandler(t *testing.T) {
	acs := accountSystem.NewDefaultAccountSystem()
	router := NewRouter("githash", "buildtime", acs)
	acs.Add(account.NewAccount(account.AccountDescription{Id: "1400", Name: "SKR03_1400_OPOS-Kunde", Type: "Aktivkonto"}))


	var stakeholder valueMagnets.Stakeholder
	sh := stakeholder.All(util.Global.FinancialYear)
	sh[2].YearlySaldo = 42.0

	req, _ := http.NewRequest("GET", "/kontrol/kommitmenschenaccounts", nil)
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

	//{"Stakeholder":[{"Id":"AN","Name":"k: Anke Nehrenberg","Type":"Partner","Arbeit":"1.00","Fairshares":"0.5"},{"Id":"RW","Name":"k: Ralf Wirdemann","Type":"Extern","Arbeit":"1.00","Fairshares":""},{"Id":"JM","Name":"k: Johannes Mainusch","Type":"Partner","Arbeit":"1.00","Fairshares":"0.5"},

	name := m["Accounts"].([]interface{})[0].(map[string]interface{})["Description"].(map[string]interface{})["Name"]
	assert.Equal(t,  "k:  kommitment" ,name)

	fSAnteil := m["Accounts"].([]interface{})[2].(map[string]interface{})["Description"].(map[string]interface{})["Name"]
	assert.Equal(t, "k:  kommitment_1-AnteilausFairshare" ,fSAnteil)

}
