package handler

import (
	"bitbucket.org/rwirdemann/kontrol/account"
	"github.com/gorilla/mux"
	"bitbucket.org/rwirdemann/kontrol/middleware"
)

func NewRouter(githash string, buildstamp string, repository account.Repository) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/kontrol/version", middleware.JWTMiddleware(MakeVersionHandler(githash, buildstamp)))
	r.HandleFunc("/kontrol/bankaccount", middleware.JWTMiddleware(MakeGetBankAccountHandler(repository)))
	r.HandleFunc("/kontrol/accounts", MakeGetAccountsHandler(repository))
	r.HandleFunc("/kontrol/accounts/{id}", MakeGetAccountHandler(repository))
	return r
}
