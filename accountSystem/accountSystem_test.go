package accountSystem

import (
	"github.com/ahojsenn/kontrol/account"
	"github.com/ahojsenn/kontrol/util"
	"github.com/ahojsenn/kontrol/valueMagnets"
	"log"
	"testing"
)

func TestNewDefaultAccountSystem(t *testing.T) {
	accountSystem := NewDefaultAccountSystem()
	util.AssertEquals(t, accountSystem.GetCollectiveAccount().Description.Id, "all")


	log.Println("in TestNewDefaultAccountSystem", accountSystem)

	ac ,exists := accountSystem.Get("SKR03_Anlagen")
	log.Println("in TestNewDefaultAccountSystem", exists, ac)
	util.AssertTrue(t, exists)

	_ ,exists = accountSystem.Get("JM")
	log.Println("in TestNewDefaultAccountSystem", exists, ac)
	util.AssertTrue(t, exists)
}

func TestEmptyDefaultAccountSystem (t *testing.T) {
	accountSystem := EmptyDefaultAccountSystem()
	util.AssertEquals(t, accountSystem.GetCollectiveAccount().Description.Id, "all")

	_ ,exists := accountSystem.Get("SKR03_Anlagen")
	util.AssertFalse(t, exists)
}

func TestAdd (t *testing.T) {
	accountSystem := EmptyDefaultAccountSystem()
	util.AssertEquals(t, accountSystem.GetCollectiveAccount().Description.Id, "all")

	newAccount := account.NewAccount(account.AccountDescription{Id: "K", Name: "k: Kommitment", Type: "company"})
	accountSystem.Add (newAccount)

	a,_ := accountSystem.Get("K")
	util.AssertEquals(t, a.Description.Name,  "k: Kommitment")
}

func TestGetAllAccountsOfStakeholder(t *testing.T) {
	accountSystem := NewDefaultAccountSystem()
	shrepo := valueMagnets.Stakeholder{}
	stakeholder := shrepo.Get("JM")

	al := accountSystem.GetAllAccountsOfStakeholder (stakeholder)
	util.AssertEquals(t, al[0].Description.Name,  "k: Johannes Mainusch")
	util.AssertEquals(t, al[4].Description.Name,  "k: Johannes Mainusch_3-Vertriebsprovision")
}
