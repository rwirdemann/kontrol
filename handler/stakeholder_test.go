package handler

import (
	"encoding/json"
	"github.com/ahojsenn/kontrol/accountSystem"
	"github.com/ahojsenn/kontrol/util"
	"github.com/ahojsenn/kontrol/valueMagnets"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)



func initRouter() *mux.Router {
	acs := accountSystem.EmptyDefaultAccountSystem()
	reimport := func (as accountSystem.AccountSystem, year int, month string) {
		//
	}
	return NewRouter("githash", "buildtime", acs, reimport)
}

func TestGetAllStakeholders(t *testing.T) {

	var stakeholder valueMagnets.Stakeholder
	sh := stakeholder.All(util.Global.FinancialYear)
	sh[2].YearlySaldo = 42.0

	req, _ := http.NewRequest("GET", "/kontrol/stakeholder", nil)
	rr := httptest.NewRecorder()

	router := initRouter()
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

	name := m["Stakeholder"].([]interface{})[0].(map[string]interface{})["Name"]
	assert.Equal(t, "k:  kommitment" , name)

	johannesName := m["Stakeholder"].([]interface{})[3].(map[string]interface{})["Name"]
	assert.Equal(t, johannesName, "k: Johannes Mainusch" )

	sh = stakeholder.All(util.Global.FinancialYear)
	johannesYearlySaldo := m["Stakeholder"].([]interface{})[2].(map[string]interface{})["YearlySaldo"]
	assert.Equal(t, 42.0 , johannesYearlySaldo )

}
