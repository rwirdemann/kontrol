package handler

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/ahojsenn/kontrol/account"
	"github.com/ahojsenn/kontrol/accountSystem"
	"github.com/ahojsenn/kontrol/util"
)

func MakeGetBilanzAccountsHandler(repository accountSystem.AccountSystem) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var accounts []account.Account
		accounts = repository.GetBilanzAccounts(account.KontenartPassiv)
		accounts = append(accounts, repository.GetBilanzAccounts(account.KontenartAktiv)... )
		sort.Sort(account.ByName(accounts))

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

