package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"

	"bitbucket.org/rwirdemann/kontrol/account"
)

func main() {
	baseUrl := "http://localhost:8991/kontrol"

	accountFlag := flag.String("account", "", "fetches given account")
	bankFlag := flag.Bool("bank", false, "fetches bank account")
	vSaldoFlag := flag.Bool("vsaldo", false, "saldo sum virtual accounts")
	flag.Parse()

	var a account.Account
	switch {
	case *bankFlag:
		get(fmt.Sprintf("%s/bankaccount", baseUrl), &a)
		a.Print()
	case *accountFlag != "":
		get(fmt.Sprintf("%s/accounts/%s", baseUrl, *accountFlag), &a)
		a.Print()
	case *vSaldoFlag:
		response := struct {
			Accounts []account.Account
		}{}
		get(fmt.Sprintf("%s/accounts", baseUrl), &response)
		saldo := 0.0
		for _, a := range response.Accounts {
			saldo += a.Saldo
		}
		fmt.Println("-------------------------------------------------------------------------------------------")
		fmt.Printf("[Saldo vAccounts: \t\t\t\t\t\t\t\t%10.2f]\n", saldo)
	}
}

func get(url string, entity interface{}) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(contents, entity)
}
