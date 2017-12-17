package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"bitbucket.org/rwirdemann/kontrol/domain"

	"bitbucket.org/rwirdemann/kontrol/processing"

	"log"

	"bitbucket.org/rwirdemann/kontrol/handler"
	"bitbucket.org/rwirdemann/kontrol/parser"
	"github.com/gorilla/mux"
	"github.com/howeyc/fsnotify"
	"github.com/rs/cors"
)

var (
	FileName   = "2017-Buchungen-KG - Buchungen 2017.csv"
	githash    string
	buildstamp string
)

const port = 8991

func main() {
	version := flag.Bool("version", false, "prints current kontrol version")
	flag.Parse()
	if *version {
		fmt.Printf("Build: %s Git: %s\n", buildstamp, githash)
		os.Exit(0)
	}

	watchBookingFile()
	importAndProcessBookings()

	r := mux.NewRouter()
	r.HandleFunc("/kontrol/version", handler.MakeVersionHandler(githash, buildstamp))
	r.HandleFunc("/kontrol/accounts", handler.Accounts)
	r.HandleFunc("/kontrol/accounts/{id}", handler.Account)

	fmt.Printf("http://localhost:8991/kontrol/accounts\n")

	// cors.Default() setup the middleware with default options being all origins accepted with simple
	// methods (GET, POST)
	handler := cors.Default().Handler(r)

	http.ListenAndServe(":"+strconv.Itoa(port), handler)
}

func importAndProcessBookings() {
	domain.ResetAccounts()
	bookings := parser.Import(FileName)
	for _, p := range bookings {
		processing.Process(p)
	}
}

func watchBookingFile() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			select {
			case <-watcher.Event:
				log.Printf("booking reimport start: %s\n", time.Now())
				importAndProcessBookings()
				log.Printf("booking reimport end: %s\n", time.Now())
			case err := <-watcher.Error:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Watch(FileName)
	if err != nil {
		log.Fatal(err)
	}
}
