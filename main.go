package main

import (
	"fmt"
	"sort"

	"bitbucket.org/rwirdemann/kontrol/kontrol"

	"bitbucket.org/rwirdemann/kontrol/processing"

	"bitbucket.org/rwirdemann/kontrol/parser"
)

func main() {
	positions := parser.Import("buchungen-2017.csv")
	for _, p := range positions {
		processing.Process(p)
	}

	for owner, account := range kontrol.Accounts {
		if owner == kontrol.SA_RW {
			sort.Sort(kontrol.ByMonth(account.Bookings))
			for _, b := range account.Bookings {
				b.Print(owner)
			}
			fmt.Println("----------------------------------------------------------------")
			fmt.Printf("Saldo: %10.2f", account.Saldo())
		}
	}
}
