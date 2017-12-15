package handler

import (
	"fmt"
	"net/http"
	"sort"

	"strconv"

	"bitbucket.org/rwirdemann/kontrol/kontrol"
	"bitbucket.org/rwirdemann/kontrol/util"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

const port = 8991

func StartService() {
	r := mux.NewRouter()
	r.HandleFunc("/kontrol/accounts", accounts)
	r.HandleFunc("/kontrol/accounts/{id}", account)

	fmt.Printf("Visit http://%s:8991/kontrol...", util.GetHostname())

	// cors.Default() setup the middleware with default options being all origins accepted with simple
	// methods (GET, POST)
	handler := cors.Default().Handler(r)

	http.ListenAndServe(":"+strconv.Itoa(port), handler)
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

	sort.Sort(kontrol.ByOwner(accounts))
	json := util.Json(AccountsResponse{Accounts: accounts})
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, json)
}

func account(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountId := vars["id"]
	account := kontrol.Accounts[accountId]
	account.UpdateSaldo()

	if account != nil {
		w.Header().Set("Content-Type", "application/json")
		sort.Sort(kontrol.ByMonth(account.Bookings))
		json := util.Json(account)
		fmt.Fprintf(w, json)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}
