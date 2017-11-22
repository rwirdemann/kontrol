package main

import (
	"bitbucket.org/rwirdemann/kontrol/rest"

	"bitbucket.org/rwirdemann/kontrol/processing"

	"bitbucket.org/rwirdemann/kontrol/parser"
	"log"
	"github.com/howeyc/fsnotify"
)

var FileName = "2017-Buchungen-KG - Buchungen 2017.csv"

func main() {
	bookings := parser.Import(FileName)
	//watchBookingFile()
	for _, p := range bookings {
		processing.Process(p)
	}

	rest.StartService()
}

func watchBookingFile() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	//done := make(chan bool)

	// Process events
	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				log.Println("event:", ev)
			case err := <-watcher.Error:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Watch("main.go")
	if err != nil {
		log.Fatal(err)
	}

	// Hang so program doesn't exit
	//<-done

	/* ... do stuff ... */
	watcher.Close()
}
