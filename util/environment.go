package util

import (
	"log"
	"io/ioutil"
	"fmt"
	"encoding/json"
	"os"
)

// environments and HTTPS certificate locations.
type Environment struct {
	Hostname string `json:"hostname"`
	CertFile string `json:"certfile"`
	KeyFile  string `json:"keyfile"`
}

func (this Environment) Get() *Environment {
	log.Println("getEnvironment: ")
	raw, err := ioutil.ReadFile("./httpsconfig.env")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	var environments []Environment
	hostname := getHostname()
	json.Unmarshal(raw, &environments)
	for i := range environments {
		if environments[i].Hostname == hostname {
			// Found hostname
			return &environments[i]
			break
		}
	}
	log.Fatal("there is no environment configured for '", hostname, "' in ./httpsconfig.env")
	return nil
}

func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	return hostname
}

