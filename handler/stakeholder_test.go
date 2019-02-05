package handler

import (
	"encoding/json"
	"github.com/ahojsenn/kontrol/accountSystem"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func initRouter() *mux.Router {
	acs := accountSystem.EmptyDefaultAccountSystem()
	return NewRouter("githash", "buildtime", acs)
}

func TestGetAllStakeholders(t *testing.T) {

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

	ankesName := m["Stakeholder"].([]interface{})[0].(map[string]interface{})["Name"]
	assert.Equal(t, ankesName, "k: Anke Nehrenberg" )

	johannesName := m["Stakeholder"].([]interface{})[2].(map[string]interface{})["Name"]
	assert.Equal(t, johannesName, "k: Johannes Mainusch" )

}
