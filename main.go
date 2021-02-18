package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/ahojsenn/kontrol/parser"
	"github.com/ahojsenn/kontrol/processing"
	"github.com/ahojsenn/kontrol/valueMagnets"

	"log"

	"github.com/ahojsenn/kontrol/accountSystem"
	"github.com/ahojsenn/kontrol/handler"
	"github.com/ahojsenn/kontrol/util"
	"github.com/howeyc/fsnotify"
	"github.com/rs/cors"
)

const DefaultBookingFile = "Buchungen-KG.csv"

var (
	githash    string
	buildstamp string
)

func main() {
	environment := util.GetEnv()
	version := flag.Bool("version", false, "prints current kontrol version")
	file := flag.String("file", DefaultBookingFile, "booking file")
	year := flag.Int("year", 2019, "year to control")
	month := flag.String("month", "*", "month to control")
	httpPort := flag.String("httpPort", "20171", "http server port")
	httpsPort := flag.String("httpsPort", "20172", "https server port")
	certFile := flag.String("certFile", environment.CertFile, "https certificate")
	keyFile := flag.String("keyFile", environment.KeyFile, "https key")
	flag.Parse()

	if *version {
		fmt.Printf("Build: %s Git: %s\n", buildstamp, githash)
		os.Exit(0)
	}
	util.Global.Filename = *file

	log.SetFlags(0)

	// set FinancialYear & month
	util.Global.FinancialYear = *year
	util.Global.FinancialMonth = *month
	bd, e := time.Parse("2006 01 02 15 04 05", strconv.Itoa(*year)+" 12 31 23 59 59")
	if e != nil {
		fmt.Println(e)
	}
	util.Global.BalanceDate = bd
	fmt.Println("\n\n#############################################################")
	log.Println("in main, util.Global.FinancialYear:", util.Global.FinancialYear,
		"\n    BalanceDate=", util.Global.BalanceDate)
	fmt.Println("#############################################################")

	// set LiquidityNeed
	util.Global.LiquidityNeed = valueMagnets.KommimtmentYear{}.Liqui(util.Global.FinancialYear)

	as := accountSystem.NewDefaultAccountSystem()
	ImportAndProcessBookings(as, *year)

	watchBookingFile(as, *year, *month)
	//ImportAndProcessBookings(as, *year, *month)

	h1 := cors.AllowAll().Handler(handler.NewRouter(githash, buildstamp, as))

	go func() {
		fmt.Printf("listing on http://localhost:%s...\n", *httpPort)
		log.Fatal(http.ListenAndServe(":"+*httpPort, h1))
	}()
	log.Println("started http server... ")

	// start HTTPS
	log.Println("starting https server	 \n  try https://localhost:" + *httpsPort + "/kontrol/accounts")
	log.Fatal(http.ListenAndServeTLS(":"+*httpsPort, *certFile, *keyFile, h1))
}

func ImportAndProcessBookings(as accountSystem.AccountSystem, year int) {
	log.Println("in ImportAndProcessBookings...")
	util.Global.Errors = nil

	as.ClearBookings()
	//	hauptbuch_allYears := as.GetCollectiveAccount_allYears()

	hauptbuch_thisYear := as.GetCollectiveAccount_thisYear()

	parser.Import(util.Global.Filename, year, as)
	log.Println("in ImportAndProcessBookings, import done...")

	// process all bookings from the general ledger
	for _, bk := range hauptbuch_thisYear.Bookings {
		processing.Process(as, bk)
	}

	// distribute revenues and costs to valueMagnets
	// in this step only employees revenues will be booked to employee cost centers
	// partners reneue will bi primarily booked to company account for this step
	processing.Kostenerteilung(as)
	processing.ErloesverteilungAnEmployees(as)
	// now employee bonusses are calculated and booked
	processing.CalculateEmployeeBonus(as)

	// now (after employee bonusses are booked) calculate GuV and Bilanz
	processing.GuV(as)
	processing.Bilanz(as)

	processing.ErloesverteilungAnKommanditisten(as)
	// distribution profit among partners
	processing.DistributeKTopf(as)
	// calculate liquidity needs per partner
	processing.BookLiquidityNeedToPartners(as, valueMagnets.KommimtmentYear{}.Liqui(util.Global.FinancialYear))
	processing.BookAmountAtDisposition(as)

	// procject Controlling
	processing.GenerateProjectControlling(as)
}

func watchBookingFile(repository accountSystem.AccountSystem, year int, month string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {

			select {
			case event := <-watcher.Event:
				// reset the timer if there are more events within the 3 seconds
				// i.e. when the file ist still loading
				log.Println("event:", event)
				select {
				case <-time.After(3 * time.Second):
					fmt.Println("timeout 3 sec")
				}
				// there might be more than one server of this kind running on this server
				// so wait for year mod 10 (seconds year %5)*5 seconds
				waitFor := (util.Global.FinancialYear % 5) * 5.0
				time.Sleep(time.Duration(waitFor)* time.Second)
				log.Printf("booking reimport start: %s %i\n", time.Now(), waitFor)
				ImportAndProcessBookings(repository, year)
				log.Printf("booking reimport end: %s\n", time.Now())
			case err := <-watcher.Error:
				log.Println("error:", err)
			}

		}
	}()

	err = watcher.Watch(util.Global.Filename)
	if err != nil {
		log.Fatal(err)
	}

}
