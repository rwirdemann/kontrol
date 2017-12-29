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
	flag.Parse()

	var a account.Account
	switch {
	case *bankFlag:
		get(fmt.Sprintf("%s/bankaccount", baseUrl), &a)
		a.Print()
	case *accountFlag != "":
		get(fmt.Sprintf("%s/accounts/%s", baseUrl, *accountFlag), &a)
		a.Print()
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
