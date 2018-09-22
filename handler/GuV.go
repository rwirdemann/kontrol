package handler

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/ahojsenn/kontrol/account"
	"github.com/ahojsenn/kontrol/accountSystem"
	"github.com/ahojsenn/kontrol/util"
)

func MakeGetGuVAccountsHandler(as accountSystem.AccountSystem) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var accounts []account.Account
		accounts = as.GetBilanzAccounts(accountSystem.KontenartErtrag)
		accounts = append(accounts, as.GetBilanzAccounts(accountSystem.KontenartAufwand)... )
		// acc ,_ := as.Get("SKR03_ErgebnisNachSteuern")
		// accounts = append(accounts, *acc )

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

