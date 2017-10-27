package main

import (
	"bitbucket.org/rwirdemann/kontrol/kontrol"

	"bitbucket.org/rwirdemann/kontrol/processing"

	"bitbucket.org/rwirdemann/kontrol/parser"
)

func main() {
	positions := parser.Import("buchungen-2017.csv")
	for _, p := range positions {
		processing.Process(p)
	}

	for k, v := range kontrol.Accounts {
		if k == kontrol.SA_RW {
			for _, b := range v.Bookings {
				b.Print(k)
			}
		}
	}
}
