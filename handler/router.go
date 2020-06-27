package handler

import (
	"github.com/ahojsenn/kontrol/accountSystem"
	"github.com/ahojsenn/kontrol/jwt-notimpl"
	"github.com/gorilla/mux"
	"io/ioutil"
)




func NewRouter( githash string, buildstamp string, as accountSystem.AccountSystem ) *mux.Router {
	r := mux.NewRouter()
	key := keycloakRSAPub()

	r.HandleFunc("/kontrol/version", jwt_notimpl.JWTMiddleware(key, MakeVersionHandler(githash, buildstamp)))
	r.HandleFunc("/kontrol/stakeholder", MakeGetStakeholderHandler())
	r.HandleFunc("/kontrol/errors", MakeGetErrorHandler(as))
	r.HandleFunc("/kontrol/collectiveaccount",  MakeGetCollectiveAccountHandler(as))
	r.HandleFunc("/kontrol/collectiveaccount/{year}", MakeGetCollectiveAccountHandler(as))
	r.HandleFunc("/kontrol/collectiveaccount/{year}/{month}", MakeGetCollectiveAccountHandler(as))
	r.HandleFunc("/kontrol/bilanz", MakeGetBilanzAccountsHandler(as))
	r.HandleFunc("/kontrol/GuV", MakeGetGuVAccountsHandler(as))
	r.HandleFunc("/kontrol/kommitmenschenaccounts", MakeGetKommitmenschenAccountsHandler(as))
	r.HandleFunc("/kontrol/kommitmenschenaccounts/{id}", MakeGetKommitmenschenAccountsHandler(as))
	r.HandleFunc("/kontrol/accounts", MakeGetAccountsHandler(as))
	r.HandleFunc("/kontrol/accounts/{id}", MakeGetAccountHandler(as))
	r.HandleFunc("/kontrol/projects", MakeGetProjectsHandler(as))

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
