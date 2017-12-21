package handler

import (
	"github.com/gorilla/mux"
)

func NewRouter(githash string, buildstamp string) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/kontrol/version", MakeVersionHandler(githash, buildstamp))
	r.HandleFunc("/kontrol/accounts", Accounts)
	r.HandleFunc("/kontrol/accounts/{id}", Account)
	return r
}
