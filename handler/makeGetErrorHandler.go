package handler

import (
	"fmt"
	"github.com/ahojsenn/kontrol/accountSystem"
	"github.com/ahojsenn/kontrol/util"
	"net/http"
)

func MakeGetErrorHandler(as accountSystem.AccountSystem) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//util.Global.Errors = append(util.Global.Errors, "testerror from your friendly Server")

		// wrap response with "Accounts" element
		response := struct {
			Errors []string
		}{
			util.Global.Errors,
		}
		json := util.Json(response)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, json)
	}
}


