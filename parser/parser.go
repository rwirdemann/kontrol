package parser

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"bitbucket.org/rwirdemann/kontrol/account"
	"bitbucket.org/rwirdemann/kontrol/owner"
)

// Beschreibt, dass die netto (Rechnungs-)Position in Spalte X der CSV-Datei dem Stakeholder Y gehört
var netBookings = []struct {
	Owner  owner.Stakeholder
	Column int
}{
	{Owner: owner.StakeholderRW, Column: 21},
	{Owner: owner.StakeholderAN, Column: 20},
	{Owner: owner.StakeholderJM, Column: 22},
	{Owner: owner.StakeholderBW, Column: 19},
	{Owner: owner.StakeholderEX, Column: 23},
}

func Import(file string) []account.Booking {
	var positions []account.Booking

	if f, err := openCsvFile(file); err == nil {
		r := csv.NewReader(bufio.NewReader(f))
		for {
			record, err := r.Read()
			if err == io.EOF {
				break
			}
			if record[0] == "GV" || record[0] == "AR" || record[0] == "ER" || record[0] == "IS" {
				typ := record[0]
				cs := record[1]
				subject := strings.Replace(record[2], "\n", ",", -1)
				amount := parseAmount(record[3])
				year, month := parseMonth(record[4])
				extras := account.CsvBookingExtras{Typ: typ, CostCenter: cs}
				extras.Net = make(map[owner.Stakeholder]float64)
				for _, p := range netBookings {
					extras.Net[p.Owner] = parseAmount(record[p.Column])
				}
				position := account.Booking{Extras: extras, Text: subject, Amount: amount, Year: year, Month: month}
				positions = append(positions, position)
			} else {
				fmt.Printf("unknown booking type %s: %s\n", record[0], record[3])
			}
		}
	} else {
		panic(err)
	}

	return positions
}

func parseAmount(amount string) float64 {
	amount = strings.Trim(amount, " ")
	if amount == "" {
		return 0
	}

	idx := strings.Index(amount, " ")
	var s string
	if idx != -1 {
		s = strings.Replace(strings.Replace(amount[0:idx], ".", "", -1), ",", ".", -1)
	} else {
		s = strings.Replace(strings.Replace(amount, ".", "", -1), ",", ".", -1)
	}

	if a, err := strconv.ParseFloat(s, 64); err == nil {
		return a
	} else {
		panic(err)
	}
}

func parseMonth(yearMonth string) (int, int) {
	s := strings.Split(yearMonth, "-")
	y, _ := strconv.Atoi(s[0])
	m, _ := strconv.Atoi(s[1])
	return y, m
}

func openCsvFile(fileName string) (*os.File, error) {

	// Open file from current directory
	if file, err := os.Open(fileName); err == nil {
		return file, nil
	}

	// Open file from GOPATH
	gopath := os.Getenv("GOPATH")
	if gopath != "" {
		if file, err := os.Open(gopath + "/src/bitbucket.org/rwirdemann/kontrol/" + fileName); err == nil {
			return file, nil
		}
	}

	return nil, errors.New("could not open " + fileName)
}
