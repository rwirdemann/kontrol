package parser

import (
	"bufio"
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strconv"
	"strings"

	"bitbucket.org/rwirdemann/kontrol/kontrol"
)

func Parse(file string) []kontrol.Position {
	var positions []kontrol.Position

	if f, err := openCsvFile(file); err == nil {
		r := csv.NewReader(bufio.NewReader(f))
		for {
			record, err := r.Read()
			if err == io.EOF {
				break
			}
			if record[0] == "ER" || record[0] == "AR" {
				typ := record[0]
				cs := record[1]
				subject := record[2]
				amount := parseAmount(record[3])
				year, month := parseMonth(record[4])

				position := kontrol.Position{Typ: typ, CostCenter: cs, Subject: subject, Amount: amount, Year: year, Month: month}

				position.Net = make(map[string]float64)
				for _, p := range kontrol.NetPositions {
					position.Net[p.Stakeholder] = parseAmount(record[p.Column])
				}

				positions = append(positions, position)
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
