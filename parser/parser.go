package parser

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ahojsenn/kontrol/accountSystem"
	"github.com/ahojsenn/kontrol/util"
	"github.com/ahojsenn/kontrol/valueMagnets"

	"github.com/ahojsenn/kontrol/booking"
)

type headerItem = struct {
	Description string
	Column      int
}

var header_basics = []headerItem{}
var header_stakeholder = []headerItem{}

func Import(file string, aYear int, as accountSystem.AccountSystem) {
	imported := 0
	hauptbuch_thisYear := as.GetCollectiveAccount_thisYear()
	hauptbuch_allYears := as.GetCollectiveAccount_allYears()
	hauptbuch_thisYear.Bookings = nil
	hauptbuch_allYears.Bookings = nil
	if file == "" {
		util.Global.Errors = append(util.Global.Errors, "in Import, no file provided...")
		log.Println("ERROR: in Import, no file provided...")
	}

	if f, err := openCsvFile(file); err == nil {
		r := csv.NewReader(bufio.NewReader(f))
		rownr := 0

		for {
			rownr++
			record, err := r.Read()
			if err == io.EOF {
				break
			}
			//log.Println("in Import, reading line ", rownr)

			if isHeader(record[0]) {
				log.Println("in Import, read header", record)
				var hi headerItem

				var stakeholderStartCol int
				if record[1] == "Cost1" {
					stakeholderStartCol = 12
				}  else {
					stakeholderStartCol = 11
				}	

				for i, s := range record {
					hi.Column = i
					hi.Description = s
					// break out of loop when hitting first empty description
					// log.Println("in Import, read header, col=", hi, hi.Description == "Cost1" )
					if hi.Description == "" {
						break
					}

					if hi.Column < stakeholderStartCol {
						header_basics = append(header_basics, hi)
					} else {
						// log.Println("in Import, read header, stakeholder =", hi, stakeholderStartCol, hi.Column )
						header_stakeholder = append(header_stakeholder, hi)
					}

				}
				//log.Println("in Import, read header", header_basics, header_stakeholder)
				continue
				//CONTINUE:
			}

			var typ, soll, haben, cs, project, cost1, subject string
			var amount float64
			var bankCreated time.Time
			var year, month int
			if isValidBookingType(record[0]) {
				if header_basics[1].Description == "Cost1" {
					typ = record[0]
					cost1 = sanitizeMyString(record[1])
					soll = record[2]
					haben = record[3]
					cs = strings.Replace(record[4], " ", "", -1) // suppress whitespace
					project = sanitizeMyString(record[5])
					subject = sanitizeMyString(record[6])
					amount = parseAmount(record[7], rownr)
					year, month = parseMonth(record[8])
					bankCreated = parseFileCreated(record[9])
				} else {
					typ = record[0]
					cost1 = "empty"
					soll = record[1]
					haben = record[2]
					cs = strings.Replace(record[3], " ", "", -1) // suppress whitespace
					project = sanitizeMyString(record[4])
					subject = sanitizeMyString(record[5])
					amount = parseAmount(record[6], rownr)
					year, month = parseMonth(record[7])
					bankCreated = parseFileCreated(record[8])
				}
				imported++
				//m := make(map[valueMagnets.Stakeholder]float64)
				m := make(map[string]float64)
				// now loop over columns with personal revenues of all stakeholders...
				shrepo := valueMagnets.Stakeholder{}

				// loop over columns until header column is empty
				for _, p := range header_stakeholder {
					//
					stakeholder := shrepo.Get(p.Description).Id
					m[stakeholder] = parseAmount(record[p.Column], rownr)
				}
				bkng := booking.NewBooking(rownr, typ, soll, haben, cs, project+";"+cost1, m, amount, subject, month, year, bankCreated)

				//				log.Println ("in Immport, ", imported, year, bkng)

				hauptbuch_allYears.Bookings = append(hauptbuch_allYears.Bookings, *bkng)

				if year == aYear {
					hauptbuch_thisYear.Bookings = append(hauptbuch_thisYear.Bookings, *bkng)
				} else {
					// log.Println ("in Immport, ", year, " is not in	 period ", aYear, rownr)
				}
			} else {
				err := fmt.Sprintf("unknown booking type '%s' in row '%d'\n", record[0], rownr)
				util.Global.Errors = append(util.Global.Errors, err)
				fmt.Printf(err)
			}
		}
	} else {
		fmt.Println("file not found", file)
		panic(err)
	}
	log.Printf("in Import, imported %d rows from %s", imported, file)
	return
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

func parseAmount(amount string, rownr int) float64 {
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
		e := fmt.Sprintf("in parseAmount: parsing error '%s' on amount '%s' in line %d\n", err, amount, rownr)
		util.Global.Errors = append(util.Global.Errors, e)
		fmt.Printf(e)
		return 0
	}
}

func parseMonth(yearMonth string) (int, int) {
	if len(yearMonth) < 2 {
		return 0, 0
	}
	s := strings.Split(yearMonth, "-")
	if len(s) < 2 {
		util.Global.Errors = append(util.Global.Errors, "in parseMonth, something went wrong with this entry")
		log.Fatal("in parseMonth, something went wrong with this entry", s)
	}
	y, err := strconv.Atoi(s[0])
	if err != nil {
		log.Fatal("ERROR in parseMonth, ", err)
	}
	m, err := strconv.Atoi(s[1])
	if err != nil {
		log.Fatal("ERROR in parseMonth, ", err)
	}
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

	/*
		// Open file from GOPATH
		gopath := os.Getenv("GOPATH")
		if gopath != "" {
			if file, err := os.Open(gopath + "/src/github.com/ahojsenn/kontrol/" + fileName); err == nil {
				return file, nil
			}
		}
	*/

	return nil, errors.New("could not open " + fileName)
}

func sanitizeMyString(in string) string {
	out := in
	//	out = strings.Replace(out, "/", "-", -1)
	out = strings.Replace(out, "\n", ",", -1)
	out = strings.Replace(out, "%", "Prozent", -1)
	return out
}
