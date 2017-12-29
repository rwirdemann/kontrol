package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"bitbucket.org/rwirdemann/kontrol/account"
)

func main() {
	baseUrl := "http://localhost:8991/kontrol"

	accountFlag := flag.String("account", "", "fetches account")
	flag.Parse()

	if *accountFlag != "" {
		url := fmt.Sprintf("%s/accounts/%s", baseUrl, *accountFlag)
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

			var a account.Account
			if err := json.Unmarshal(contents, &a); err != nil {
				log.Fatal(err)
			}

			a.Print()
		}
	}
}
