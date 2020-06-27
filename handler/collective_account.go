package handler

import (
	"fmt"
	"github.com/ahojsenn/kontrol/account"
	"github.com/ahojsenn/kontrol/accountSystem"
	"github.com/ahojsenn/kontrol/booking"
	"github.com/ahojsenn/kontrol/util"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func MakeGetCollectiveAccountHandler(repository accountSystem.AccountSystem) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ad := account.AccountDescription{Id: "all", Name: "Hauptbuch", Type: account.KontenartVerrechnung}
		resultAccount := &account.Account{Description: ad}
		vars := mux.Vars(r)
		year, urlHasAYear := vars["year"]
		month, urlHasAMonth := vars["month"]
		fmt.Sprintf(resultAccount.Description.Type, "Verrechnungskonto")

		if urlHasAMonth {
			fmt.Sprintf(resultAccount.Description.Id, "allBookings_%s-%s", year, month)
			fmt.Sprintf(resultAccount.Description.Name, "allBookings_%s-%s", year, month)
			resultAccount.Bookings = filterBookingsByYearAndMonth(repository.GetCollectiveAccount_allYears(), year, month)
		} else if urlHasAYear {
			fmt.Sprintf(resultAccount.Description.Id, "allBookings_%s", year)
			fmt.Sprintf(resultAccount.Description.Name, "allBookings_%s", year)
			resultAccount.Bookings = filterBookingsByYear(repository.GetCollectiveAccount_allYears(), year)
		} else {
			fmt.Sprintf( resultAccount.Description.Id,  "allBookings")
			fmt.Sprintf( resultAccount.Description.Name,  "allBookings")
			ca  := repository.GetCollectiveAccount_allYears()
			resultAccount.Bookings = ca.Bookings
		}

		// kopie erstellen!

		resultAccount.Nbookings = len(resultAccount.Bookings)

		w.Header().Set("Content-Type", "application/json")
		//sort.Sort(booking.ByMonth(resultAccount.Bookings))
		//sort.Sort(booking.ByRowNr(resultAccount.Bookings))
		json := util.Json(resultAccount)
		fmt.Fprintf(w, json)
	}
}


func filterBookingsByYear(account *account.Account, year string) []booking.Booking {
	var bookings  []booking.Booking

	iyear, err := strconv.Atoi(year)
	if err != nil {
		fmt.Println(err)
	} else {
		for _, b := range account.Bookings {
			if b.Year == iyear {
				bookings = append (bookings, b)
			}
		}
	}
	return bookings
}

func filterBookingsByYearAndMonth(account *account.Account, year, month string) []booking.Booking {
	var bookings []booking.Booking

	iyear, err := strconv.Atoi(year)
	imonth, err1 := strconv.Atoi(month)
	if err != nil {
		fmt.Println(err)
	} else if err1 != nil {
		fmt.Println(err1)
    } else {
		for _, b := range account.Bookings {
			if b.Year == iyear && b.Month == imonth{
				bookings = append (bookings, b)
			}
		}
	}
	return bookings
}
