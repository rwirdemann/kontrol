package main

import (
	"fmt"

	"bitbucket.org/rwirdemann/kontrol/parser"
)

func main() {
	positions := parser.Import("buchungen-2017.csv")
	for _, p := range positions {
		if p.Typ == "AR" {
			fmt.Printf("%v\n", p)
		}
	}
}
