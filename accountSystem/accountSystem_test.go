package accountSystem

import (
	"testing"
	"github.com/ahojsenn/kontrol/util"
		"github.com/ahojsenn/kontrol/account"
			)

func TestNewDefaultAccountSystem(t *testing.T) {
	accountSystem := NewDefaultAccountSystem()
	util.AssertEquals(t, accountSystem.BankAccount().Description.Id, "GLS")

	_ ,exists := accountSystem.Get("SKR03_Anlagen")
	util.AssertTrue(t, exists)

	_ ,exists = accountSystem.Get("JM")
	util.AssertTrue(t, exists)
}

func TestEmptyDefaultAccountSystem (t *testing.T) {
	accountSystem := EmptyDefaultAccountSystem()
	util.AssertEquals(t, accountSystem.BankAccount().Description.Id, "GLS")

	_ ,exists := accountSystem.Get("SKR03_Anlagen")
	util.AssertFalse(t, exists)
}

func TestAdd (t *testing.T) {
	accountSystem := EmptyDefaultAccountSystem()
	util.AssertEquals(t, accountSystem.BankAccount().Description.Id, "GLS")

	newAccount := account.NewAccount(account.AccountDescription{Id: "K", Name: "k: Kommitment", Type: "company"})
	accountSystem.Add (newAccount)

	a,_ := accountSystem.Get("K")
	util.AssertEquals(t, a.Description.Name,  "k: Kommitment")
}

func TestDetermineSollOrHaben (t *testing.T) {
	accountSystem := NewDefaultAccountSystem()
	var amount float64
	a,_ := accountSystem.Get("1600")
	amount = accountSystem.DetermineSollOrHaben(100.0, a, "haben")
	util.AssertEquals(t, amount,  -100.0)

	a,_ = accountSystem.Get("4100_4199")
	amount = accountSystem.DetermineSollOrHaben(100.0, a, "soll")
	util.AssertEquals(t, amount,  -100.0)
}

