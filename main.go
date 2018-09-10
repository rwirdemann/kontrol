package main

import (
		"flag"
	"fmt"
		"net/http"
	"os"
	"time"

		"github.com/ahojsenn/kontrol/processing"

	"log"

	"github.com/ahojsenn/kontrol/handler"
	"github.com/ahojsenn/kontrol/parser"
	"github.com/howeyc/fsnotify"
	"github.com/rs/cors"
	"github.com/ahojsenn/kontrol/util"
	"github.com/ahojsenn/kontrol/accountSystem"
)

const DefaultBookingFile = "2017-Buchungen-KG - Buchungen 2017.csv"

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
	httpPort := flag.String("httpPort", "8991", "http server port")
	httpsPort := flag.String("httpsPort", "8992", "https server port")
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
	log.Println("in main, util.Global.FinancialYear:", util.Global.FinancialYear)

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

func importAndProcessBookings(repository accountSystem.AccountSystem, year int) {
	repository.ClearBookings()
	log.Printf("importAndProcessBookings: %d\n", year)
	hauptbuch := repository.GetCollectiveAccount()
	parser.Import(fileName, year, &(hauptbuch.Bookings))
	log.Println("in main, import done")
	for _, p := range hauptbuch.Bookings {
		processing.Process(repository, p)
	}
	// now calculate GuV
	processing.GuV(repository)
	processing.Bilanz(repository)
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
