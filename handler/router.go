package handler

import (
	"io/ioutil"

		"github.com/ahojsenn/kontrol/middleware"
	"github.com/gorilla/mux"
	"github.com/ahojsenn/kontrol/accountSystem"
)

func NewRouter(githash string, buildstamp string, repository accountSystem.AccountSystem) *mux.Router {
	r := mux.NewRouter()
	key := keycloakRSAPub()
	r.HandleFunc("/kontrol/version", middleware.JWTMiddleware(key, MakeVersionHandler(githash, buildstamp)))
	r.HandleFunc("/kontrol/collectiveaccount", middleware.JWTMiddleware(key, MakeGetCollectiveAccountHandler(repository)))
	r.HandleFunc("/kontrol/bilanz", MakeGetBilanzAccountsHandler(repository))
	r.HandleFunc("/kontrol/GuV", MakeGetGuVAccountsHandler(repository))
	r.HandleFunc("/kontrol/kommitmenschen", MakeGetKommitmenschenAccountsHandler(repository))
	r.HandleFunc("/kontrol/accounts", MakeGetAccountsHandler(repository))
	r.HandleFunc("/kontrol/accounts/{id}", MakeGetAccountHandler(repository))
	r.HandleFunc("/kontrol/projects", MakeGetProjectsHandler(repository))
	return r
}

const defaultRSAPublicKey = `-----BEGIN RSA PRIVATE KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAl+hgn1UL8XjGZA6PHJfA16+adPHcgFBReXqiw0ykWivX46Ass0Psw60Q6nxhgOU4fa0QNVuxeTrvxpgiPu2MhrSSJxsNW2BiYETfAdjCpbwqmR1JhpyL1UR1MOdtkIe+Ucy6tYmpL4lt9gkREgDpv8pfQXAk5tYlaGnhPZM/53kRV3N1cFYlYC65ykY+JDkJdT74gFKbekOtYQiJPfmeuBOtBNZ1FiSf7T9k1bhtks6Q9ZbDjr8y9ax5OHCZEJLhTzQTIN4YdpV5nFR3eBZ5/kOS/E60JUjTWKgMYuSeoRi5drcYxQXPlM/gmCgY9igKfNa4gEeGx1cq1LSxV01tQQIDAQAB
-----END RSA PRIVATE KEY-----`

func keycloakRSAPub() []byte {
	if k, err := ioutil.ReadFile("keycloak_rsa.pub"); err == nil {
		return k
	}

	return []byte(defaultRSAPublicKey)
}
