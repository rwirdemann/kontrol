package handler

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/ahojsenn/kontrol/account"
	"github.com/ahojsenn/kontrol/accountSystem"
	"github.com/ahojsenn/kontrol/util"
)

func MakeGetProjectsHandler(as accountSystem.AccountSystem) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var accounts []account.Account
		// accounts = as.CloneAccountsOfType(account.KontenartKLR)
		for _, acc := range as.GetByType(account.KontenartProject) {
			acc.UpdateSaldo()
			accounts = append(accounts, *acc)
		}

		sort.Sort(account.ByOwner(accounts))

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

