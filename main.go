package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"bitbucket.org/rwirdemann/kontrol/domain"
	"bitbucket.org/rwirdemann/kontrol/handler"

	"bitbucket.org/rwirdemann/kontrol/processing"

	"log"

	"bitbucket.org/rwirdemann/kontrol/parser"
	"github.com/howeyc/fsnotify"
)

var (
	FileName = "2017-Buchungen-KG - Buchungen 2017.csv"
	Version  string
	Build    string
)

func main() {
	version := flag.Bool("version", false, "prints current kontrol version")
	flag.Parse()
	if *version {
		fmt.Printf("Build: %s Git: %s\n", Build, Version)
		os.Exit(0)
	}

	watchBookingFile()
	importAndProcessBookings()
	handler.StartService()
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
