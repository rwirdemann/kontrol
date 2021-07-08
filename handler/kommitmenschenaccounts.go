package handler

import (
	"fmt"
	"github.com/ahojsenn/kontrol/account"
	"github.com/ahojsenn/kontrol/accountSystem"
	"github.com/ahojsenn/kontrol/util"
	"github.com/ahojsenn/kontrol/valueMagnets"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

func MakeGetKommitmenschenAccountsHandler(as accountSystem.AccountSystem) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var stkhldr valueMagnets.Stakeholder
		var accounts []account.Account

		for _,sh := range stkhldr.All(util.Global.FinancialYear) {
			al := as.GetAllAccountsOfStakeholder(sh)
			accounts = append (accounts, al...)
		}

		// filter by id if provided in URL
		vars := mux.Vars(r)
		accountId, ok := vars["id"]
		if ok {
			accounts = FilterAccountsByStakeholder(accounts, accountId)
		} else {
			accounts = accounts
		}

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


func FilterAccountsByStakeholder(accounts []account.Account, stakeholder string) []account.Account {
	var filtered []account.Account
	for _, b := range accounts {
		if strings.Contains(b.Description.Id, stakeholder)   {
			filtered = append(filtered, b)
		}
	}
	return filtered
}
