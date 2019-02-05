package handler

import (
	"fmt"
	"github.com/ahojsenn/kontrol/account"
	"github.com/ahojsenn/kontrol/accountSystem"
	"github.com/ahojsenn/kontrol/util"
	"github.com/ahojsenn/kontrol/valueMagnets"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"sort"
	"strings"
)

func MakeGetKommitmenschenAccountsHandler(as accountSystem.AccountSystem) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var accounts []account.Account
		//accounts = append(accounts, as.CloneAccountsOfType(valueMagnets.StakeholderTypeCompany)... )
		// company
		for _,acc := range as.CloneAccountsOfType(valueMagnets.StakeholderTypeCompany) {
			accounts = append (accounts, acc)
			// add underraccounts
			ua := as.CloneAccountsOfType(valueMagnets.StakeholderTypeKUA)
			sort.Sort(account.ByName(ua))
			for _,acc2 := range ua {
				if acc2.Description.Superaccount ==  acc.Description.Id {
					accounts = append (accounts, acc2)
				}
			}
		}

		// Others
		// accounts = append(accounts, as.CloneAccountsOfType(valueMagnets.StakeholderTypeOthers)... )
		for _,acc := range as.CloneAccountsOfType(valueMagnets.StakeholderTypeOthers) {
			accounts = append (accounts, acc)
			// add underraccounts
			ua := as.CloneAccountsOfType(valueMagnets.StakeholderTypeKUA)
			sort.Sort(account.ByName(ua))
			for _,acc2 := range ua {
				if acc2.Description.Superaccount ==  acc.Description.Id {
					accounts = append (accounts, acc2)
				}
			}
		}

		// external accounts
		for _,acc := range as.CloneAccountsOfType(valueMagnets.StakeholderTypeExtern) {
			accounts = append (accounts, acc)
			// add underraccounts
			ua := as.CloneAccountsOfType(valueMagnets.StakeholderTypeKUA)
			sort.Sort(account.ByName(ua))
			for _,acc2 := range ua {
				if acc2.Description.Superaccount ==  acc.Description.Id {
					accounts = append (accounts, acc2)
				}
			}
		}

		// for the  rest of all employees
		employees := as.CloneAccountsOfType(valueMagnets.StakeholderTypeEmployee)
		sort.Sort(account.ByName(employees))
		for _,acc := range employees {
			accounts = append (accounts, acc)
			// add underraccounts
			ua := as.CloneAccountsOfType(valueMagnets.StakeholderTypeKUA)
			sort.Sort(account.ByName(ua))
			for _,acc2 := range ua {
				if acc2.Description.Superaccount ==  acc.Description.Id {
					accounts = append (accounts, acc2)
				}
			}
		}
		// for the  rest of all partners
		partners := as.CloneAccountsOfType(valueMagnets.StakeholderTypePartner )
		sort.Sort(account.ByName(partners))
		for _,acc := range partners {
			accounts = append (accounts, acc)
			// add underraccounts
			ua := as.CloneAccountsOfType(valueMagnets.StakeholderTypeKUA)
			sort.Sort(account.ByName(ua))
			for _,acc2 := range ua {
				if acc2.Description.Superaccount ==  acc.Description.Id {
					accounts = append (accounts, acc2)
				}
			}
		}

		vars := mux.Vars(r)
		accountId, ok := vars["id"]
		log.Println("in MakeGetKommitmenschenAccountsHandler:", accountId, ok)
		if ok {
			accounts = FilterAccountsByStakeholder(accounts, accountId)
		} else {
			accounts = accounts
		}

		// update all saldos
		for _,acc := range accounts {
			acc.UpdateSaldo()
		}

		// wrap response with "Accounts" element
		response := struct {
			Accounts []account.Account
		}{
			accounts,
		}
		json := util.Json(response)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, json)
	}
}


func FilterAccountsByStakeholder(accounts []account.Account, stakeholder string) []account.Account {
	var filtered []account.Account
	for _, b := range accounts {
		if strings.Contains(b.Description.Id, stakeholder)   {
			filtered = append(filtered, b)
		}
	}
	return filtered
}
