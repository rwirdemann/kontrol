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
	KontenartProject      	  = "VerrechnungskontoProjekt"
)

const (
	EmployeeShare			 = 0.70
	KommmitmentExternShare   = 0.95
	KommmitmentOthersShare   = 1.00
	KommmitmentEmployeeShare = 0.95
	PartnerProvision         = 0.05
)

type AccountDescription struct {
	Id   string `json:",omitempty"`
	Name string
	Type string
}

type Account struct {
	Description     AccountDescription
	Bookings  []booking.Booking `json:",omitempty"`
	KommitmenschNettoFaktura   float64
	AnteilAusFaktura float64
	AnteilAusFairshares float64
	KommitmenschDarlehen float64
	Costs     float64
	Advances  float64
	Reserves  float64
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
	revenue := 0.0
	saldo := 0.0
	kommitmenschNettoFaktura := 0.0
	internals := 0.0
	advances := 0.0
	rest := 0.0
	costs := 0.0
	anteilAusFairshares := 0.0
	anteilAusFaktura := 0.0
	darlehen := 0.0
	for _, b := range a.Bookings {
		saldo += b.Amount

		switch b.Type {
		case booking.CC_Nettoanteil, booking.CC_Kommitmentanteil:
			revenue += b.Amount
		case booking.CC_Entnahme:
			advances += b.Amount
		case booking.CC_InterneStunden:
			internals += b.Amount
		case booking.CC_Vertriebsprovision:
			revenue += b.Amount
		case booking.Erloese, booking.CC_KommitmentanteilEX:
			revenue += b.Amount
		case booking.CC_Employeeaanteil:
			revenue += b.Amount
			kommitmenschNettoFaktura += b.Amount/ EmployeeShare
		case booking.CC_PartnerNettoFaktura:
			revenue += b.Amount
			kommitmenschNettoFaktura += b.Amount
		case booking.Kosten, booking.CC_LNSteuer, booking.CC_GWSteuer, booking.CC_SVBeitrag, booking.CC_Gehalt:
			costs += b.Amount
		case booking.CC_AnteilAusFairshares:
			anteilAusFairshares += b.Amount
		case booking.CC_AnteilAusFaktura:
			anteilAusFaktura += b.Amount
		case booking.CC_KommitmenschDarlehen:
			darlehen += b.Amount
		case booking.SKR03:
			// hier genauer gucken...
		default:
			rest += b.Amount
		}
	}
	a.Saldo = saldo
	a.Advances = advances
	a.Revenue = revenue
	a.Internals = internals
	a.Rest = rest
	a.Costs = costs
	a.KommitmenschNettoFaktura = kommitmenschNettoFaktura
	a.AnteilAusFaktura = anteilAusFaktura
	a.AnteilAusFairshares = anteilAusFairshares
	a.KommitmenschDarlehen = darlehen
	a.KommitmenschNettoFaktura = kommitmenschNettoFaktura
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

type ByName []Account
func (a ByName) Len() int           { return len(a) }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByName) Less(i, j int) bool { return strings.Compare(a[i].Description.Name, a[j].Description.Name) < 0 }
