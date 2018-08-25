package accountSystem

import (
	"testing"
	"github.com/ahojsenn/kontrol/util"
		"github.com/ahojsenn/kontrol/account"
			)

func TestNewDefaultAccountSystem(t *testing.T) {
	accountSystem := NewDefaultAccountSystem()
	util.AssertEquals(t, accountSystem.GetCollectiveAccount().Description.Id, "all")

	_ ,exists := accountSystem.Get("SKR03_Anlagen")
	util.AssertTrue(t, exists)

	_ ,exists = accountSystem.Get("JM")
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
