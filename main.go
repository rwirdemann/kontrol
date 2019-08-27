package main

import (
	"flag"
	"fmt"
	"github.com/ahojsenn/kontrol/valueMagnets"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/ahojsenn/kontrol/processing"

	"log"

	"github.com/ahojsenn/kontrol/accountSystem"
	"github.com/ahojsenn/kontrol/handler"
	"github.com/ahojsenn/kontrol/parser"
	"github.com/ahojsenn/kontrol/util"
	"github.com/howeyc/fsnotify"
	"github.com/rs/cors"
)

const DefaultBookingFile = "Buchungen-KG.csv"

var (
	fileName   	string
	githash    	string
	buildstamp 	string
)


func main() {
	environment := util.GetEnv()
	version := flag.Bool("version", false, "prints current kontrol version")
	file := flag.String("file", DefaultBookingFile, "booking file")
	year := flag.Int("year", 2018, "year to control")
	httpPort := flag.String("httpPort", "20171", "http server port")
	httpsPort := flag.String("httpsPort", "20172", "https server port")
	certFile := flag.String("certFile", environment.CertFile, "https certificate")
	keyFile := flag.String("keyFile", environment.KeyFile, "https key")
	flag.Parse()

	if *version {
		fmt.Printf("Build: %s Git: %s\n", buildstamp, githash)
		os.Exit(0)
	}
	fileName = *file



	// set FinancialYear
	util.Global.FinancialYear =  *year
	bd, e := time.Parse("2006 01 02 15 04 05", strconv.Itoa(*year) + " 12 31 23 59 59"  )
	if e != nil {
		fmt.Println(e)
	}
	util.Global.BalanceDate = bd
	log.Println("in main, util.Global.FinancialYear:", util.Global.FinancialYear,
		"\n    BalanceDate=",util.Global.BalanceDate)

	// set LiquidityNeed
	util.Global.LiquidityNeed =  valueMagnets.KommimtmentYear{}.Liqui(util.Global.FinancialYear)
	log.Println("in main: util.Global.LiquidityNeed=",util.Global.LiquidityNeed)


	accountSystem := accountSystem.NewDefaultAccountSystem()
	log.Println("in main, created accountsystem for ", util.Global.FinancialYear)
	watchBookingFile(accountSystem, *year)
	importAndProcessBookings(accountSystem, *year)

	handler := cors.AllowAll().Handler(handler.NewRouter(githash, buildstamp, accountSystem))

	go func() {
		fmt.Printf("listing on http://localhost:%s...\n", *httpPort)
		log.Fatal(http.ListenAndServe(":"+*httpPort, handler))
	}()
	log.Println("started http server... ")

	// start HTTPS
	log.Println("starting https server	 \n  try https://localhost:" + *httpsPort + "/kontrol/accounts")
	log.Fatal(http.ListenAndServeTLS(":"+*httpsPort, *certFile, *keyFile, handler))
}

func importAndProcessBookings(as accountSystem.AccountSystem, year int) {
	as.ClearBookings()
	log.Printf("importAndProcessBookings: %d\n", year)
	hauptbuch := as.GetCollectiveAccount()
	parser.Import(fileName, year, &(hauptbuch.Bookings))
	log.Println("in main, import done")
	for _, p := range hauptbuch.Bookings {
		processing.Process(as, p)
	}

	// distribute revenues and costs to valueMagnets
	// in this step only employees revenues will be booked to employee cost centers
	// partners reneue will bi primarily booked to company account for this step
	processing.ErloesverteilungAnStakeholder(as)
	// now employee bonusses are calculated and booked
	processing.CalculateEmployeeBonus(as)

	// now (after employee bonusses are booked) calculate GuV and Bilanz
	processing.GuV(as)
	processing.Bilanz(as)

	// distribution profit among partners
	processing.DistributeKTopf(as)
	// calculate liquidity needs per partner
	processing.BookLiquidityNeedToPartners(as, valueMagnets.KommimtmentYear{}.Liqui(util.Global.FinancialYear))
	processing.BookAmountAtDisposition(as)

	// procject Controlling
	processing.GenerateProjectControlling(as)
}

func watchBookingFile(repository accountSystem.AccountSystem, year int) {
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
				log.Printf("booking reimport start: %s\n", time.Now())
				importAndProcessBookings(repository, year)
				log.Printf("booking reimport end: %s\n", time.Now())
			case err := <-watcher.Error:
				log.Println("error:", err)
			}

		}
	}()

	err = watcher.Watch(fileName)
	if err != nil {
		log.Fatal(err)
	}

}
