package handler

import (
	"fmt"
	"net/http"
	"sort"

	"bitbucket.org/rwirdemann/kontrol/account"
	"bitbucket.org/rwirdemann/kontrol/util"
	"github.com/gorilla/mux"
	"bitbucket.org/rwirdemann/kontrol/booking"
)

func MakeGetAccountsHandler(repository account.Repository) http.HandlerFunc {
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

func MakeGetAccountHandler(repository account.Repository) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		accountId := vars["id"]
		if a, ok := repository.Get(accountId); ok {

			// Check to filter account by costcenter
			costcenter := r.URL.Query().Get("cs")
			if costcenter != "" {
				a = a.FilterBookingsByCostcenter(costcenter)
			}

			a.UpdateSaldo()
			w.Header().Set("Content-Type", "application/json")
			sort.Sort(booking.ByMonth(a.Bookings))
			json := util.Json(a)
			fmt.Fprintf(w, json)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}
}
