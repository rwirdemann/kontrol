package account

import (
	"fmt"
	"github.com/ahojsenn/kontrol/util"
	"log"
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
	Superaccount string
}

type Account struct {
	Description     AccountDescription
	Bookings  []booking.Booking `json:",omitempty"`
	Nbookings int
	Soll	float64
	Haben 	float64
	KommitmenschNettoFaktura   float64 	// is used for the net faktura of people
//	YearS	float64 					// this ist used for the sum of all bookings in the current year
	Saldo	float64
}

func (a *Account) SumOfBookingType  (btype string ) float64 {
	sum := 0.0
	for _, b := range a.Bookings {
		switch btype {
		case b.Type,"":
			sum += b.Amount
		}
	}
	return sum
}


func NewAccount(a AccountDescription) *Account {
	return &Account{Description: a }
}

func (a *Account) Book(b booking.Booking) {
	if a == nil {
		log.Println ("in Book, got nil account in row ", b.RowNr,b )
		return
	}
	b.Id = util.GetNewBookingId()
	a.Bookings = append(a.Bookings, b)
	a.UpdateSaldo()  // this might be expensive, but this way the Salden should always be accurate
	// log.Println("in Book: ", a)
}

func (a *Account) UpdateSaldo() {
	saldo 	:= 0.0
	soll 	:= 0.0
	haben 	:= 0.0
	kommitmenschNettoFaktura := 0.0
	for _, b := range a.Bookings {
		saldo += b.Amount

		switch b.Type {
		case booking.CC_Employeeaanteil:
			kommitmenschNettoFaktura += b.Amount / EmployeeShare
		case booking.CC_PartnerNettoFaktura:
			kommitmenschNettoFaktura += b.Amount
		}

		// soll und haben richtig verbuchen
		switch a.Description.Type {
		case KontenartAufwand,  KontenartAktiv:
			if b.Amount > 0.0 { soll += b.Amount } else { haben += b.Amount }
		case KontenartErtrag,  KontenartPassiv, KontenartProject:
			if b.Amount < 0.0 { soll += b.Amount } else { haben += b.Amount }
		default:
		}
	}
	a.Soll = soll
	a.Haben = haben
	a.Saldo = saldo
	a.Nbookings = len(a.Bookings)
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
