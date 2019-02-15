package handler

import (
	"fmt"
	"github.com/ahojsenn/kontrol/util"
	"github.com/ahojsenn/kontrol/valueMagnets"
	"net/http"
)

// delivers all stakeholders in the company
// including employees, externals, kommanditisten and the company itsself


func MakeGetStakeholderHandler() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		var stakeholder valueMagnets.Stakeholder

		// wrap response with "Accounts" element
		response := struct {
			Stakeholder []valueMagnets.Stakeholder
		}{
			stakeholder.All(util.Global.FinancialYear),
		}
		json := util.Json(response)
		fmt.Fprintf(w, json)
	}
}

