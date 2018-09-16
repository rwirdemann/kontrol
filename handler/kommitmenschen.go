package handler

import (
	"fmt"
	"github.com/ahojsenn/kontrol/valueMagnets"
	"net/http"
	"sort"

	"github.com/ahojsenn/kontrol/account"
	"github.com/ahojsenn/kontrol/accountSystem"
	"github.com/ahojsenn/kontrol/util"
)

func MakeGetKommitmenschenAccountsHandler(repository accountSystem.AccountSystem) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var accounts []account.Account
		accounts = repository.GetBilanzAccounts(valueMagnets.StakeholderTypeCompany)
		accounts = append(accounts, repository.GetBilanzAccounts(valueMagnets.StakeholderTypeExtern)... )
		accounts = append(accounts, repository.GetBilanzAccounts(valueMagnets.StakeholderTypeEmployee)... )
		accounts = append(accounts, repository.GetBilanzAccounts(valueMagnets.StakeholderTypePartner)... )
		sort.Sort(account.ByType(accounts))

		// wrap response with "Accounts" element
		response := struct {
			Accounts []account.Account
		}{
			accounts,
		}
		json := util.Json(response)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, json)
	}
}

