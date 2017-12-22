package handler

import (
	"fmt"
	"net/http"
	"sort"

	"bitbucket.org/rwirdemann/kontrol/account"
	"bitbucket.org/rwirdemann/kontrol/domain"
	"bitbucket.org/rwirdemann/kontrol/util"
	"github.com/gorilla/mux"
)

func MakeGetAccountsHandler(repository account.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accounts := repository.All()
		sort.Sort(domain.ByOwner(accounts))

		// wrap response with "Accounts" element
		response := struct {
			Accounts []domain.Account
		}{
			accounts,
		}
		json := util.Json(response)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, json)
	}
}

func MakeGetAccountHandler(repository account.Repository) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		accountId := vars["id"]
		account, _ := repository.Get(accountId)
		account.UpdateSaldo()

		if account != nil {
			w.Header().Set("Content-Type", "application/json")
			sort.Sort(domain.ByMonth(account.Bookings))
			json := util.Json(account)
			fmt.Fprintf(w, json)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}
}
