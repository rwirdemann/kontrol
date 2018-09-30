package account

import (
	"fmt"
	"sort"
	"strings"

	"github.com/ahojsenn/kontrol/booking"
)

const (
	KontenartAktiv            = "Aktivkonto"
	KontenartPassiv           = "Passivkonto"
	KontenartAufwand          = "Aufwandskonto"
	KontenartErtrag           = "Ertragskonto"
	KontenartVerrechnung      = "Verrechnungskonto"
	KontenartKLR      		  = "VerrechnungskontoKLR"
)

type AccountDescription struct {
	Id   string `json:",omitempty"`
	Name string
	Type string
}

type Account struct {
	Description     AccountDescription
	Bookings  []booking.Booking `json:",omitempty"`
	Costs     float64
	Advances  float64
	Reserves  float64
	Provision float64
	Rest      float64
	Revenue   float64
	Taxes     float64
	Internals float64
	Saldo     float64
}

func NewAccount(a AccountDescription) *Account {
	return &Account{Description: a}
}

func (a *Account) Book(b booking.Booking) {
	a.Bookings = append(a.Bookings, b)
}

func (a *Account) UpdateSaldo() {
	provision := 0.0
	revenue := 0.0
	saldo := 0.0
	internals := 0.0
	advances := 0.0
	rest := 0.0
	costs := 0.0
	for _, b := range a.Bookings {
		saldo += b.Amount

		switch b.Type {
		case booking.Nettoanteil, booking.Kommitmentanteil:
			revenue += b.Amount
		case booking.Entnahme:
			advances += b.Amount
		case booking.Vertriebsprovision:
			provision += b.Amount
		case booking.InterneStunden:
			internals += b.Amount
		case booking.Erloese, booking.Employeeaanteil:
			revenue += b.Amount
		case booking.Kosten, booking.LNSteuer, booking.GWSteuer, booking.SVBeitrag, booking.Gehalt, booking.Eingangsrechnung:
			costs += b.Amount
		case booking.SKR03:
			// hier genauer gucken...
		default:
			rest += b.Amount
		}
	}
	a.Saldo = saldo
	a.Advances = advances
	a.Revenue = revenue
	a.Provision = provision
	a.Internals = internals
	a.Rest = rest
	a.Costs = costs
}

func (a Account) Print() {
	sort.Sort(booking.ByMonth(a.Bookings))
	for _, b := range a.Bookings {
		b.Print(a.Description.Id)
	}
	fmt.Println("-------------------------------------------------------------------------------------------")
	fmt.Printf("[Saldo: \t\t\t\t\t\t\t\t\t%10.2f]\n", a.Saldo)
}

func (a Account) CSV() string {
	result := "Konto;Monat;Jahr;Mitarbeiter;Typ;Buchungstext;Betrag\n"
	sort.Sort(booking.ByMonth(a.Bookings))
	for _, b := range a.Bookings {
		result = result + b.CSV(a.Description.Id)
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
func (a ByOwner) Less(i, j int) bool { return strings.Compare(a[i].Description.Name, a[j].Description.Name) < 0 }

type ByType []Account
func (a ByType) Len() int           { return len(a) }
func (a ByType) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByType) Less(i, j int) bool { return strings.Compare(a[i].Description.Type, a[j].Description.Type) < 0 }

