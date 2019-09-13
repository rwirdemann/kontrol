package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
)

// environments and HTTPS certificate locations.
type Environment struct {
	Hostname string `json:"hostname"`
	CertFile string `json:"certfile"`
	KeyFile  string `json:"keyfile"`
	KommitmentFile  string `json:"kommitmentFile"`
}

// global scope
type GlobalStruct struct {
	Filename string
	FinancialMonth string
	FinancialYear int
	BalanceDate time.Time
	LiquidityNeed float64
}

var Global GlobalStruct

func GetEnv() *Environment{
	hostname := getHostname()

	raw, err := ioutil.ReadFile("./httpsconfig.env")
	if err != nil {
		//
		raw, err = ioutil.ReadFile("../httpsconfig.env")
		if err != nil {
			fmt.Println("in GetEnv(): ", err.Error())
			os.Exit(1)
		}
	}
	var environments []Environment

	json.Unmarshal(raw, &environments)
	for i := range environments {
		if environments[i].Hostname == hostname {
			return &environments[i]
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

