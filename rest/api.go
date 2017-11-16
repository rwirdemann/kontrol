package rest

import (
	"fmt"
	"net/http"
	"sort"

	"bitbucket.org/rwirdemann/kontrol/kontrol"
	"github.com/gorilla/mux"
	"github.com/arschles/go-bindata-html-template"
	"bitbucket.org/rwirdemann/kontrol/html"
	"strconv"
	"bitbucket.org/rwirdemann/kontrol/util"
)

const port = 8991

func StartService() {
	r := mux.NewRouter()
	r.HandleFunc("/kontrol", index)
	r.HandleFunc("/kontrol/accounts", accounts)
	r.HandleFunc("/kontrol/accounts/{id}", account)

	fmt.Printf("Visit http://%s:8991/kontrol...", util.GetHostname())
	http.ListenAndServe(":"+strconv.Itoa(port), r)
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
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, json)
}

func account(w http.ResponseWriter, r *http.Request) {
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

// Data struct for index.html
type Index struct {
	Hostname string
	Port     int
}

func index(w http.ResponseWriter, r *http.Request) {
	hostname := util.GetHostname()
	index := Index{Hostname: hostname, Port: port}
	t, _ := template.New("index", html.Asset).Parse("html/index.html")
	t.Execute(w, struct{ Index }{index})
}
