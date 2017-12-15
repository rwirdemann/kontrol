package handler

import (
	"fmt"
	"net/http"
	"sort"

	"strconv"

	"bitbucket.org/rwirdemann/kontrol/domain"
	"bitbucket.org/rwirdemann/kontrol/util"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

const port = 8991

func StartService() {
	r := mux.NewRouter()
	r.HandleFunc("/kontrol/accounts", accounts)
	r.HandleFunc("/kontrol/accounts/{id}", account)

	fmt.Printf("Visit http://%s:8991/domain....", util.GetHostname())

	// cors.Default() setup the middleware with default options being all origins accepted with simple
	// methods (GET, POST)
	handler := cors.Default().Handler(r)

	http.ListenAndServe(":"+strconv.Itoa(port), handler)
}

type AccountsResponse struct {
	Accounts []domain.Account
}

func accounts(w http.ResponseWriter, r *http.Request) {

	// convert account map to array
	accounts := make([]domain.Account, 0, len(domain.Accounts))
	for _, a := range domain.Accounts {
		a.UpdateSaldo()
		accounts = append(accounts, *a)
	}

	sort.Sort(domain.ByOwner(accounts))
	json := util.Json(AccountsResponse{Accounts: accounts})
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, json)
}

func account(w http.ResponseWriter, r *http.Request) {
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
