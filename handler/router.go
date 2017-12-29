package handler

import (
	"bitbucket.org/rwirdemann/kontrol/account"
	"github.com/gorilla/mux"
)

func NewRouter(githash string, buildstamp string, repository account.Repository) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/kontrol/version", MakeVersionHandler(githash, buildstamp))
	r.HandleFunc("/kontrol/bankaccount", MakeGetBankAccountHandler(repository))
	r.HandleFunc("/kontrol/accounts", MakeGetAccountsHandler(repository))
	r.HandleFunc("/kontrol/accounts/{id}", MakeGetAccountHandler(repository))
	return r
}
