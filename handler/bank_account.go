package handler

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/ahojsenn/kontrol/account"
	"github.com/ahojsenn/kontrol/booking"
	"github.com/ahojsenn/kontrol/util"
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
