package rest

import (
	"fmt"
	"net/http"

	"bitbucket.org/rwirdemann/kontrol/kontrol"
	"bitbucket.org/rwirdemann/notux/util"
	"github.com/gorilla/mux"
)

func StartService() {
	r := mux.NewRouter()
	r.HandleFunc("/kontrol", index)
	r.HandleFunc("/kontrol/accounts/{id}/bookings", bookings)

	fmt.Printf("http://localhost:8991/kontrol/accounts/RW/bookings")
	http.ListenAndServe(":8991", r)
}

func index(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

type ResponseWrapper struct {
	Bookings []kontrol.Booking
	Owner    string
	Saldo    string
}

func bookings(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountId := vars["id"]
	account := kontrol.Accounts[accountId]

	if account != nil {
		w.WriteHeader(http.StatusOK)
		response := ResponseWrapper{Owner: account.Owner.Name, Bookings: account.Bookings, Saldo: fmt.Sprintf("%.2f €", account.Saldo())}
		json := util.Json(response)
		fmt.Fprintf(w, json)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}
