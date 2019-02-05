package handler

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/ahojsenn/kontrol/account"
	"github.com/ahojsenn/kontrol/booking"
	"github.com/ahojsenn/kontrol/util"
	"github.com/gorilla/mux"
	"github.com/ahojsenn/kontrol/accountSystem"
)

func MakeGetAccountsHandler(repository accountSystem.AccountSystem) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accounts := repository.All()
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

func MakeGetAccountHandler(as accountSystem.AccountSystem) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		// accept only requests from localhost
		//...

		vars := mux.Vars(r)
		accountId := vars["id"]
		if a, ok := as.Get(accountId); ok {

			// Check to filter account by costcenter
			costcenter := r.URL.Query().Get("cs")
			if costcenter != "" {
				a = a.FilterBookingsByCostcenter(costcenter)
			}

			a.UpdateSaldo()
			w.Header().Set("Content-Type", "application/json")
			sort.Sort(booking.ByRowNr(a.Bookings))
			json := util.Json(a)
			fmt.Fprintf(w, json)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}
}
