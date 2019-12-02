package handler

import (
	"github.com/ahojsenn/kontrol/accountSystem"
	"github.com/ahojsenn/kontrol/jwt-notimpl"
	"github.com/ahojsenn/kontrol/util"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

type reimportFunc func (as accountSystem.AccountSystem, year int, month string)
type queryParser struct {
	reimport reimportFunc
	as accountSystem.AccountSystem
}

func (qp *queryParser) parseQ(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		qp.reimportDataIfneeded(w, r)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func (qp *queryParser) reimportDataIfneeded (w http.ResponseWriter, r *http.Request) {

	// find financialMonth if present in query and if changed, reimport data
	financialMonth := r.URL.Query().Get("financialMonth")

	if financialMonth != util.Global.FinancialMonth &&
		r.URL.Path != "/kontrol/version" &&
		r.URL.Path != "/kontrol/stakeholder" &&
		r.URL.Path != "/kontrol/errors" {
		util.Global.FinancialMonth  = financialMonth
		log.Println("MakeParseQueryParams: triggerring reprocessing of data...\n\n", util.Global.FinancialYear, financialMonth, r.URL)
		//accountSystem.ImportAndProcessBookings(util.Global.AccountSystem, util.Global.FinancialYear, util.Global.FinancialMonth)
		qp.reimport (qp.as, util.Global.FinancialYear, util.Global.FinancialMonth)
	}
}


func NewRouter( githash string, buildstamp string, as accountSystem.AccountSystem, reimport reimportFunc ) *mux.Router {
	r := mux.NewRouter()
	key := keycloakRSAPub()

	r.HandleFunc("/kontrol/version", jwt_notimpl.JWTMiddleware(key, MakeVersionHandler(githash, buildstamp)))
	r.HandleFunc("/kontrol/stakeholder", MakeGetStakeholderHandler())
	r.HandleFunc("/kontrol/errors", MakeGetErrorHandler(as))

	r.HandleFunc("/kontrol/collectiveaccount", jwt_notimpl.JWTMiddleware(key, MakeGetCollectiveAccountHandler(as)))
	r.HandleFunc("/kontrol/bilanz", MakeGetBilanzAccountsHandler(as))
	r.HandleFunc("/kontrol/GuV", MakeGetGuVAccountsHandler(as))
	r.HandleFunc("/kontrol/kommitmenschenaccounts", MakeGetKommitmenschenAccountsHandler(as))
	r.HandleFunc("/kontrol/kommitmenschenaccounts/{id}", MakeGetKommitmenschenAccountsHandler(as))
	r.HandleFunc("/kontrol/accounts", MakeGetAccountsHandler(as))
	r.HandleFunc("/kontrol/accounts/{id}", MakeGetAccountHandler(as))
	r.HandleFunc("/kontrol/projects", MakeGetProjectsHandler(as))

	qp := queryParser{}
	qp.reimport = reimport
	qp.as = as
	r.Use( qp.parseQ )

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
