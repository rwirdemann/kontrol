package processing

import (
	"github.com/ahojsenn/kontrol/booking"
	"github.com/ahojsenn/kontrol/account"
	"github.com/ahojsenn/kontrol/accountSystem"
)

func GenerateProjectControlling  (as accountSystem.AccountSystem) {

	// concat KontenartErtrag and KontenartAufwand into accList
	accList := as.GetByType(account.KontenartErtrag)
	for k, v := range as.GetByType(account.KontenartAufwand) {
		accList[k] = v
	}


	for _, acc := range accList {
		for _, bk := range acc.Bookings {

			// handle empty projects
			if bk.Project == "" {
				bk.Project = "_PROJEKT-emptyProject"
			}
			// check if there is an project account bk.Projects
			acc, exists := as.Get("_PROJEKT-"+bk.Project)
			if !exists {
				// create a new projects account
				acc = account.NewAccount(account.AccountDescription{Id: "_PROJEKT-"+bk.Project, Name: bk.Project, Type: account.KontenartProject})
				as.Add(acc)
			}

			// subtract "ER" from "AR"
			sign := +1.0
			if (bk.Type == "ER") {
				sign = -1.0
			}

			// now create a booking in the appropriate projectAccount
			clonedBooking := booking.Booking{
				RowNr:       bk.RowNr,
				Amount:      sign*bk.Amount,
				Project:     bk.Project,
				Type:        bk.Type,
				CostCenter:  bk.CostCenter,
				Month:		 bk.Month,
				Year:        bk.Year,
				Text:        bk.Text,
				FileCreated: bk.FileCreated,
				BankCreated: bk.BankCreated,
			}
			// and book it to tha account
			acc.Book(clonedBooking)
		}
	}
}


