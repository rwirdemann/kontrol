package main

import (
	"time"

	"bitbucket.org/rwirdemann/kontrol/domain"
	"bitbucket.org/rwirdemann/kontrol/handler"

	"bitbucket.org/rwirdemann/kontrol/processing"

	"log"

	"bitbucket.org/rwirdemann/kontrol/parser"
	"github.com/howeyc/fsnotify"
)

var FileName = "2017-Buchungen-KG - Buchungen 2017.csv"

func main() {
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
