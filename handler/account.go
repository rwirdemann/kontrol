package handler

import (
	"fmt"
	"net/http"
	"sort"

	"bitbucket.org/rwirdemann/kontrol/domain"
	"bitbucket.org/rwirdemann/kontrol/util"
	"github.com/gorilla/mux"
)

func Accounts(w http.ResponseWriter, r *http.Request) {

	// convert account map to array
	accounts := make([]domain.Account, 0, len(domain.Accounts))
	for _, a := range domain.Accounts {
		a.UpdateSaldo()
		accounts = append(accounts, *a)
	}

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

func Account(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountId := vars["id"]
	account := domain.Accounts[accountId]
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
