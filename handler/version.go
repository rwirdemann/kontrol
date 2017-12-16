package handler

import (
	"fmt"
	"net/http"
)

func MakeVersionHandler(githash string, buildstamp string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("Githash: %s Buildstamp: %s", githash, buildstamp)))
	}
}
