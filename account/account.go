package account

import (
	"fmt"
	"sort"
	"strings"

	"github.com/ahojsenn/kontrol/booking"
	"github.com/ahojsenn/kontrol/owner"
)

type Account struct {
	Owner     owner.Stakeholder
	Bookings  []booking.Booking `json:",omitempty"`
	Costs     float64
	Advances  float64
	Reserves  float64
	Provision float64
	Revenue   float64
	Taxes     float64
	Internals float64
	Saldo     float64
}

func NewAccount(o owner.Stakeholder) *Account {
	return &Account{Owner: o}
}

func (a *Account) Book(booking booking.Booking) {
	a.Bookings = append(a.Bookings, booking)
}

func (a *Account) UpdateSaldo() {
	costs := 0.0
	provision := 0.0
	revenue := 0.0
	taxes := 0.0
	saldo := 0.0
	internals := 0.0
	advances := 0.0
	for _, b := range a.Bookings {
		saldo += b.Amount
		if b.Type == "Entnahme" {
			advances += b.Amount
		}
		if b.Type == "Vertriebsprovision" {
			provision += b.Amount
		}
		if b.Type == "Vertriebsprovision" || b.Type == "Nettoanteil" || b.Type == "Kommitmentanteil" {
			revenue += b.Amount
		}
		if b.Type == "GWSteuer" {
			taxes += b.Amount
		}
		if b.Type == "Eingangsrechnung" || b.Type == "Gehalt" || b.Type == "SV-Beitrag" {
			costs += b.Amount
		}
		if b.Type == "Interne Stunden" {
			internals += b.Amount
		}
	}
	a.Saldo = saldo
	a.Advances = advances
	a.Costs = costs
	a.Revenue = revenue
	a.Taxes = taxes
	a.Provision = provision
	a.Internals = internals
}

func (a Account) Print() {
	sort.Sort(booking.ByMonth(a.Bookings))
	for _, b := range a.Bookings {
		b.Print(a.Owner)
	}
	fmt.Println("-------------------------------------------------------------------------------------------")
	fmt.Printf("[Saldo: \t\t\t\t\t\t\t\t\t%10.2f]\n", a.Saldo)
}

func (a Account) CSV() string {
	result := "Konto;Monat;Jahr;Mitarbeiter;Typ;Buchungstext;Betrag\n"
	sort.Sort(booking.ByMonth(a.Bookings))
	for _, b := range a.Bookings {
		result = result + b.CSV(a.Owner)
	}
	return result
}

func (a Account) FilterBookingsByCostcenter(costcenter string) *Account {
	var filtered []booking.Booking
	for _, b := range a.Bookings {
		if b.CostCenter == costcenter {
			filtered = append(filtered, b)
		}
	}
	a.Bookings = filtered
	return &a
}

type ByOwner []Account

func (a ByOwner) Len() int           { return len(a) }
func (a ByOwner) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByOwner) Less(i, j int) bool { return strings.Compare(a[i].Owner.Name, a[j].Owner.Name) < 0 }
