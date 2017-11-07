package kontrol

import (
	"fmt"
	"sort"
)

type Account struct {
	Owner    Stakeholder
	Bookings []Booking
}

var Accounts map[string]*Account

func init() {
	Accounts = make(map[string]*Account)
	for _, sh := range AllStakeholder {
		Accounts[sh.Id] = &Account{Owner: sh}
	}
}

func (a *Account) Book(booking Booking) {
	a.Bookings = append(a.Bookings, booking)
}

func (a Account) Saldo() float64 {
	saldo := 0.0
	for _, b := range a.Bookings {
		saldo += b.Amount
	}
	return saldo
}

func (a Account) Print() {
	sort.Sort(ByMonth(a.Bookings))
	for _, b := range a.Bookings {
		b.Print(a.Owner)
	}
	fmt.Println("-------------------------------------------------------------------------------------------")
	fmt.Printf("[Saldo: \t\t\t\t\t\t\t\t\t%10.2f]\n", a.Saldo())
}
