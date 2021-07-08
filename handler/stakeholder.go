package handler

import (
	"fmt"
	"github.com/ahojsenn/kontrol/util"
	"github.com/ahojsenn/kontrol/valueMagnets"
	"net/http"
	"strconv"
	"log"
	"github.com/gorilla/mux"
)

// delivers all stakeholders in the company
// including employees, externals, kommanditisten and the company itsself

func MakeGetStakeholderHandler() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		var stakeholder valueMagnets.Stakeholder
		vars := mux.Vars(r)
		yearString, urlHasAYear := vars["year"]

		year := util.Global.FinancialYear
		if urlHasAYear { 
			year,_ = strconv.Atoi(yearString)
		}

		response := struct {
			Stakeholder []valueMagnets.Stakeholder
		}{
			stakeholder.All(year),
		}
	
		log.Println("MakeGetStakeholderHandler: ", urlHasAYear, year, stakeholder.All(year))


		json := util.Json(response)
		fmt.Fprintf(w, json)
	}
}

