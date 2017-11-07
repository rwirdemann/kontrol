package rest

import (
	"fmt"
	"net/http"
	"sort"

	"bitbucket.org/rwirdemann/kontrol/kontrol"
	"bitbucket.org/rwirdemann/notux/util"
	"github.com/gorilla/mux"
)

func StartService() {
	r := mux.NewRouter()
	r.HandleFunc("/kontrol", index)
	r.HandleFunc("/kontrol/accounts", accounts)

	// todo: should be account instead of booking
	r.HandleFunc("/kontrol/accounts/{id}/bookings", bookings)

	fmt.Printf("http://localhost:8991/kontrol/accounts/RW/bookings")
	http.ListenAndServe(":8991", r)
}

func index(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

type AccountsResponse struct {
	Accounts []kontrol.Account
}

func accounts(w http.ResponseWriter, r *http.Request) {

	// convert account map to array
	accounts := make([]kontrol.Account, 0, len(kontrol.Accounts))
	for _, a := range kontrol.Accounts {
		a.UpdateSaldo()
		accounts = append(accounts, *a)
	}

	json := util.Json(AccountsResponse{Accounts: accounts})
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, json)
}

func bookings(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountId := vars["id"]
	account := kontrol.Accounts[accountId]

	if account != nil {
		w.WriteHeader(http.StatusOK)
		sort.Sort(kontrol.ByMonth(account.Bookings))
		json := util.Json(account)
		fmt.Fprintf(w, json)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}
