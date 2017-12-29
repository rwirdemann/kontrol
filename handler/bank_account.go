package handler

import (
	"fmt"
	"net/http"
	"sort"

	"bitbucket.org/rwirdemann/kontrol/account"
	"bitbucket.org/rwirdemann/kontrol/util"
)

func MakeGetBankAccountHandler(repository account.Repository) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		a := repository.CollectiveAccount()
		a.UpdateSaldo()
		w.Header().Set("Content-Type", "application/json")
		sort.Sort(account.ByMonth(a.Bookings))
		json := util.Json(a)
		fmt.Fprintf(w, json)
	}
}
