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

//		for _, acc := range as.All() {
//			acc.UpdateSaldo()
//			// log.Println ("in MakeGetGuVAccountsHandler", acc)
//		}

		var accounts []account.Account
		accounts = as.CloneAccountsOfType(account.KontenartErtrag)
		accounts = append(accounts, as.CloneAccountsOfType(account.KontenartAufwand)... )

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

