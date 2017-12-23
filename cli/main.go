package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"bitbucket.org/rwirdemann/kontrol/domain"
)

func main() {
	baseUrl := "http://localhost:8991/kontrol"

	account := flag.String("account", "", "fetches account")
	flag.Parse()

	if *account != "" {
		url := fmt.Sprintf("%s/accounts/%s", baseUrl, *account)
		response, err := http.Get(url)
		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		} else {
			defer response.Body.Close()
			contents, err := ioutil.ReadAll(response.Body)
			if err != nil {
				fmt.Printf("%s", err)
				os.Exit(1)
			}

			var account domain.Account
			if err := json.Unmarshal(contents, &account); err != nil {
				log.Fatal(err)
			}

			account.Print()
		}

	}
}
