package accountSystem

import (
	"github.com/ahojsenn/kontrol/account"
	"github.com/ahojsenn/kontrol/util"
	"github.com/ahojsenn/kontrol/valueMagnets"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestNewDefaultAccountSystem(t *testing.T) {
	util.Global.FinancialYear = 1337
	accountSystem := NewDefaultAccountSystem()
	assert.Equal(t, "all_1337", accountSystem.GetCollectiveAccount_thisYear().Description.Id)
	assert.Equal(t, "all", accountSystem.GetCollectiveAccount_allYears().Description.Id)


	log.Println("in TestNewDefaultAccountSystem", accountSystem)

	ac ,exists := accountSystem.Get("SKR03_Anlagen")
	log.Println("in TestNewDefaultAccountSystem", exists, ac)
	util.AssertTrue(t, exists)

	_ ,exists = accountSystem.Get("JM")
	log.Println("in TestNewDefaultAccountSystem", exists, ac)
	util.AssertTrue(t, exists)
}
/*
func TestEmptyDefaultAccountSystem (t *testing.T) {
	util.Global.FinancialYear = 1337
	accountSystem := EmptyDefaultAccountSystem()

	assert.Equal(t, accountSystem.GetCollectiveAccount_thisYear().Description.Id, "all")
	assert.Equal(t, accountSystem.GetCollectiveAccount_allYears().Description.Id, "all_1337")

	_ ,exists := accountSystem.Get("SKR03_Anlagen")
	util.AssertFalse(t, exists)
}
*/
func TestAdd (t *testing.T) {
	accountSystem := EmptyDefaultAccountSystem()
	assert.Equal(t, accountSystem.GetCollectiveAccount_thisYear().Description.Id, "all")

	newAccount := account.NewAccount(account.AccountDescription{Id: "K", Name: "k: Kommitment", Type: "company"})
	accountSystem.Add (newAccount)

	a,_ := accountSystem.Get("K")
	assert.Equal(t, a.Description.Name,  "k: Kommitment")
}

func TestGetAllAccountsOfStakeholder(t *testing.T) {
	accountSystem := NewDefaultAccountSystem()
	shrepo := valueMagnets.Stakeholder{}
	stakeholder := shrepo.Get("JM")

	al := accountSystem.GetAllAccountsOfStakeholder (stakeholder)
	assert.Equal(t, al[0].Description.Name,  "k: Johannes Mainusch")
	assert.Equal(t, al[3].Description.Name,  "k: Johannes Mainusch_2-Vertriebsprovision")
}
