package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"
	"io/ioutil"
	"encoding/json"

	"bitbucket.org/rwirdemann/kontrol/account"

	"bitbucket.org/rwirdemann/kontrol/processing"

	"log"

	"bitbucket.org/rwirdemann/kontrol/handler"
	"bitbucket.org/rwirdemann/kontrol/parser"
	"github.com/howeyc/fsnotify"
	"github.com/rs/cors"
)

const DefaultBookingFile = "2017-Buchungen-KG - Buchungen 2017.csv"

var (
	fileName   string
	githash    string
	buildstamp string
)

// environments and HTTPS certificate locations.
type Environment struct {
	Hostname string  `json:"hostname"`
	CertFile string  `json:"certfile"`
	KeyFile  string  `json:"keyfile"`
}

func getEnvironment() *Environment {
	log.Println ("getEnvironment: ")
	raw, err := ioutil.ReadFile("./httpsconfig.env")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	var environments []Environment
	hostname := getHostname()
	json.Unmarshal(raw, &environments)
	for i := range environments {
		if environments[i].Hostname == hostname {
			// Found hostname
			return &environments[i]
			break
		}
	}
	log.Fatal ("there is no environment configured for '", hostname, "' in ./httpsconfig.env")
	return nil
}

func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	return hostname
}

func main() {
	environment := getEnvironment()
	version 		:= flag.Bool("version", false, "prints current kontrol version")
	file 				:= flag.String("file", DefaultBookingFile, "booking file")
	year 				:= flag.Int("year", 2018, "year to control")
	httpPort		:= flag.String("httpPort", "8991", "http server port")
	httpsPort		:= flag.String("httpsPort", "8992", "https server port")
	certFile 		:= flag.String("certFile", environment.CertFile, "https certificate")
	keyFile		 	:= flag.String("keyFile", environment.KeyFile, "https key")
	flag.Parse()

	if *version {
		fmt.Printf("Build: %s Git: %s\n", buildstamp, githash)
		os.Exit(0)
	}
	fileName = *file

	repository := account.NewDefaultRepository()
	watchBookingFile(repository, *year)
	importAndProcessBookings(repository, *year)

	handler := cors.AllowAll().Handler(handler.NewRouter(githash, buildstamp, repository))
	go func() {
		fmt.Printf("listing on http://localhost:%s...\n", *httpPort)
		log.Fatal(http.ListenAndServe(":"+ *httpPort, handler))
		} ()
		log.Println("started http server... ")
		// start HTTPS
		log.Println("starting https server	 \n  try https://localhost:"+ *httpsPort+"/kontrol/accounts")
		log.Fatal(http.ListenAndServeTLS(":"+ *httpsPort, *certFile, *keyFile, handler))
	}

	func importAndProcessBookings(repository account.Repository, year int) {
		repository.ClearBookings()
		log.Printf("importAndProcessBookings: %d\n", year)
		bookings := parser.Import(fileName, year)
		for _, p := range bookings {
			processing.Process(repository, p)
		}
	}

	func watchBookingFile(repository account.Repository, year int) {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Fatal(err)
		}

		go func() {
			for {
				// there is no nicer war to create an 'empty NewTimer'
				// https://github.com/golang/go/issues/12721
				timer := time.NewTimer(0)
				<- timer.C
				select {
				case event := <-watcher.Event:
					// reset the timer if there are more events within the 3 seconds
					// i.e. when the file ist still loading
					log.Println("event:", event)
					timer.Reset(3*time.Second)
				case err := <-watcher.Error:
					log.Println("error:", err)
				case <-timer.C:
					// now wait until no further event has been written for one second...
					// to prevent the process from reading the file while it is still
					// being written...
					log.Printf("booking reimport start: %s\n", time.Now())
					importAndProcessBookings(repository, year)
					log.Printf("booking reimport end: %s\n", time.Now())
				}
				timer.Stop()
			}
			}()

			err = watcher.Watch(fileName)
			if err != nil {
				log.Fatal(err)
			}
		}
