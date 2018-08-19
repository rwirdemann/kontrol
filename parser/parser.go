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
	"time"

	"github.com/ahojsenn/kontrol/booking"
	"github.com/ahojsenn/kontrol/owner"
		)

// Beschreibt, dass die netto (Rechnungs-)Position in Spalte X der CSV-Datei dem Stakeholder Y gehört
var netBookings = []struct {
	Owner  string
	Column int
}{
	{Owner: "BW", Column: 21},
	{Owner: "AN", Column: 22},
	{Owner: "RW", Column: 23},
	{Owner: "JM", Column: 24},
	{Owner: "KR", Column: 25},
	{Owner: "EX", Column: 26},
	{Owner: "RR", Column: 27},
}

func Import(file string, aYear int) []booking.Booking {
	var positions []booking.Booking

	if f, err := openCsvFile(file); err == nil {
		r := csv.NewReader(bufio.NewReader(f))
		for {
			record, err := r.Read()
			if err == io.EOF {
				break
			}

			if isHeader(record[0]) {
				continue
			}

			if isValidBookingType(record[0]) {
				typ := record[0]
				soll := record[1]
				haben := record[2]
				cs := record[3]
				subject := strings.Replace(record[4], "\n", ",", -1)
				amount := parseAmount(record[5])
				year, month := parseMonth(record[6])
				bankCreated := parseFileCreated(record[7])
				if year == aYear {
					m := make(map[owner.Stakeholder]float64)
					for _, p := range netBookings {
						// die Owner Zuordnung muss hier anders sein...
						stakeholder := owner.StakeholderRepository{}.Get(p.Owner)
						m[stakeholder] = parseAmount(record[p.Column])
					}
					position := booking.NewBooking(typ, soll, haben, cs, m, amount, subject, month, year, bankCreated)
					positions = append(positions, *position)
				}
			} else {
				fmt.Printf("unknown booking type '%s'\n", record[0])
			}
		}
	} else {
		panic(err)
	}

	return positions
}

func isHeader(s string) bool {
	return strings.Contains(s, ":")
}

func isValidBookingType(s string) bool {
	for _, t := range booking.ValidBookingTypes {
		if s == t {
			return true
		}
	}
	return false
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
		fmt.Printf("parsing error '%s'\n", err)
		return 0
	}
}

func parseMonth(yearMonth string) (int, int) {
	if len(yearMonth) < 2 {
		return 0, 0
	}
	s := strings.Split(yearMonth, "-")
	y, _ := strconv.Atoi(s[0])
	m, _ := strconv.Atoi(s[1])
	return y, m
}

func parseFileCreated(fileCreated string) time.Time {
	s := strings.Split(fileCreated, ".")
	if len(s) != 3 {
		return time.Time{}
	}

	day, _ := strconv.Atoi(s[0])
	month, _ := strconv.Atoi(s[1])
	year, _ := strconv.Atoi(s[2])
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

func openCsvFile(fileName string) (*os.File, error) {

	// Open file from current directory
	if file, err := os.Open(fileName); err == nil {
		return file, nil
	}

	// Open file from GOPATH
	gopath := os.Getenv("GOPATH")
	if gopath != "" {
		if file, err := os.Open(gopath + "/src/github.com/ahojsenn/kontrol/" + fileName); err == nil {
			return file, nil
		}
	}

	return nil, errors.New("could not open " + fileName)
}
