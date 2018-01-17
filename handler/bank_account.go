package handler

import (
	"fmt"
	"net/http"
	"sort"

	"bitbucket.org/rwirdemann/kontrol/account"
	"bitbucket.org/rwirdemann/kontrol/util"
	"bitbucket.org/rwirdemann/kontrol/booking"
)

func MakeGetBankAccountHandler(repository account.Repository) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		a := repository.BankAccount()
		a.UpdateSaldo()
		w.Header().Set("Content-Type", "application/json")
		sort.Sort(booking.ByMonth(a.Bookings))
		json := util.Json(a)
		fmt.Fprintf(w, json)
	}
}
