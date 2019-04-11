package main

import (
	"flag"
	"fmt"
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
	year := flag.Int("year", 2017, "year to control")
	liquidityNeed := flag.Float64("liquidityNeed", 300000.0, "needed liquidity for this year")
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
	util.Global.LiquidityNeed =  *liquidityNeed


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

	// verteile Erlöse
	processing.ErloesverteilungAnStakeholder(as)
	processing.CalculateEmployeeBonus(as)

	// now calculate GuV and Bilanz
	processing.GuV(as)
	processing.Bilanz(as)

	processing.DistributeKTopf(as)
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
