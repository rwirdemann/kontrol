package handler

import (
	"bitbucket.org/rwirdemann/kontrol/account"
	"github.com/gorilla/mux"
	"bitbucket.org/rwirdemann/kontrol/middleware"
	"io/ioutil"
	"log"
)

func NewRouter(githash string, buildstamp string, repository account.Repository) *mux.Router {
	r := mux.NewRouter()
	key := keycloakRSAPub()
	r.HandleFunc("/kontrol/version", middleware.JWTMiddleware(key, MakeVersionHandler(githash, buildstamp)))
	r.HandleFunc("/kontrol/bankaccount", middleware.JWTMiddleware(key, MakeGetBankAccountHandler(repository)))
	r.HandleFunc("/kontrol/accounts", MakeGetAccountsHandler(repository))
	r.HandleFunc("/kontrol/accounts/{id}", MakeGetAccountHandler(repository))
	return r
}

const defaultRSAPublicKey = `-----BEGIN RSA PRIVATE KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAl+hgn1UL8XjGZA6PHJfA16+adPHcgFBReXqiw0ykWivX46Ass0Psw60Q6nxhgOU4fa0QNVuxeTrvxpgiPu2MhrSSJxsNW2BiYETfAdjCpbwqmR1JhpyL1UR1MOdtkIe+Ucy6tYmpL4lt9gkREgDpv8pfQXAk5tYlaGnhPZM/53kRV3N1cFYlYC65ykY+JDkJdT74gFKbekOtYQiJPfmeuBOtBNZ1FiSf7T9k1bhtks6Q9ZbDjr8y9ax5OHCZEJLhTzQTIN4YdpV5nFR3eBZ5/kOS/E60JUjTWKgMYuSeoRi5drcYxQXPlM/gmCgY9igKfNa4gEeGx1cq1LSxV01tQQIDAQAB
-----END RSA PRIVATE KEY-----`

func keycloakRSAPub()[]byte  {
	if k, err := ioutil.ReadFile("keycloak_rsa.pub"); err == nil {
		return k
	}

	log.Printf("create 'keycloak_rsa.pub' in your main application folder")
	return []byte(defaultRSAPublicKey)
}