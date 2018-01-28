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
	costcenterFlag := flag.String("costcenter", "", "fetches kommitmentaccount filtered by costcenter")
	bankFlag := flag.Bool("bank", false, "fetches bank account")
	vSaldoFlag := flag.Bool("vsaldo", false, "saldo sum virtual accounts")
	checkFlag := flag.Bool("check", false, "checks virtual accounts against bank account saldo")
	flag.Parse()

	var a account.Account
	switch {
	case *bankFlag:
		bankAccount(baseUrl).Print()
	case *accountFlag != "":
		if *costcenterFlag != "" {
			get(fmt.Sprintf("%s/accounts/%s?cs=%s", baseUrl, *accountFlag, *costcenterFlag), &a)
		} else {
			get(fmt.Sprintf("%s/accounts/%s", baseUrl, *accountFlag), &a)
		}
		a.Print()
	case *vSaldoFlag:
		saldo := virtualAccountsSaldo(baseUrl)
		fmt.Println("-------------------------------------------------------------------------------------------")
		fmt.Printf("[Saldo vAccounts: \t\t\t\t\t\t\t\t%10.2f]\n", saldo)
	case *checkFlag:
		banksaldo := bankAccount(baseUrl).Saldo
		saldo := virtualAccountsSaldo(baseUrl)
		fmt.Printf("Saldo Bank Accoount:\t%10.2f\n"+
			"Saldo vAccounts:\t\t%10.2f\n"+
			"Diff:\t\t\t\t\t%10.2f\n", banksaldo, saldo, banksaldo-saldo)
	}
}

func virtualAccountsSaldo(baseUrl string) float64 {
	response := struct {
		Accounts []account.Account
	}{}
	get(fmt.Sprintf("%s/accounts", baseUrl), &response)
	saldo := 0.0
	for _, a := range response.Accounts {
		saldo += a.Saldo
	}
	return saldo
}

func bankAccount(baseUrl string) *account.Account {
	var a account.Account
	get(fmt.Sprintf("%s/bankaccount", baseUrl), &a)
	return &a
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
