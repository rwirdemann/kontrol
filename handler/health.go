package handler

import (
	"fmt"
	"net/http"

	"github.com/ahojsenn/kontrol/accountSystem"
	"github.com/ahojsenn/kontrol/util"
)

// MakeGetHealthHandler ...
func MakeGetHealthHandler(as accountSystem.AccountSystem) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		response := struct {
			Health string
		}{
			"OK",
		}
		json := util.Json(response)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, json)
	}
}
