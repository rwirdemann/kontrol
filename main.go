package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

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

const port = 8991

func main() {
	version := flag.Bool("version", false, "prints current kontrol version")
	file := flag.String("file", DefaultBookingFile, "booking file")
	year := flag.Int("year", 2017, "year to control")
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
	fmt.Printf("listing on http://localhost:%d...\n", port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), handler))
}

func importAndProcessBookings(repository account.Repository, year int) {
	repository.ClearBookings()
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
			select {
			case <-watcher.Event:
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
