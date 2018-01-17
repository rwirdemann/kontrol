package account

import (
	"fmt"
	"sort"
	"strings"

	"bitbucket.org/rwirdemann/kontrol/owner"
	"bitbucket.org/rwirdemann/kontrol/booking"
)

type Account struct {
	Owner    owner.Stakeholder
	Bookings []booking.Booking `json:",omitempty"`
	Saldo    float64
}

func NewAccount(o owner.Stakeholder) *Account {
	return &Account{Owner: o}
}

func (a *Account) Book(booking booking.Booking) {
	a.Bookings = append(a.Bookings, booking)
}

func (a *Account) UpdateSaldo() {
	saldo := 0.0
	for _, b := range a.Bookings {
		saldo += b.Amount
	}
	a.Saldo = saldo
}

func (a Account) Print() {
	sort.Sort(booking.ByMonth(a.Bookings))
	for _, b := range a.Bookings {
		b.Print(a.Owner)
	}
	fmt.Println("-------------------------------------------------------------------------------------------")
	fmt.Printf("[Saldo: \t\t\t\t\t\t\t\t\t%10.2f]\n", a.Saldo)
}

type ByOwner []Account

func (a ByOwner) Len() int           { return len(a) }
func (a ByOwner) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByOwner) Less(i, j int) bool { return strings.Compare(a[i].Owner.Name, a[j].Owner.Name) < 0 }
