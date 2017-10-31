package main

import (
	"bitbucket.org/rwirdemann/kontrol/kontrol"

	"bitbucket.org/rwirdemann/kontrol/processing"

	"bitbucket.org/rwirdemann/kontrol/parser"
)

func main() {
	bookings := parser.Import("buchungen-2017.csv")
	for _, p := range bookings {
		processing.Process(p)
	}

	kontrol.Accounts[kontrol.SA_JM].Print()
}
