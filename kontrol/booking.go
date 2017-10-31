package kontrol

import (
	"fmt"
)

// Beschreibt, dass die netto (Rechnungs-)Position in Spalte X der CSV-Datei dem Stakeholder Y gehört
type NetBookingColumn struct {
	Stakeholder string
	Column      int
}

// Liste aller Spalten-Stateholder Positions-Mappings
var NetBookings = []NetBookingColumn{
	NetBookingColumn{Stakeholder: SA_RW, Column: NET_COL_RW},
	NetBookingColumn{Stakeholder: SA_AN, Column: NET_COL_AN},
	NetBookingColumn{Stakeholder: SA_JM, Column: NET_COL_JM},
	NetBookingColumn{Stakeholder: SA_BW, Column: NET_COL_BW},
	NetBookingColumn{Stakeholder: SA_EX, Column: NET_COL_EX},
}

// Zusatzinformationen einer Buchung, deren Quelle die CSV-Datei ist, und die für die weitere
// Bearbeitung erforderlich sind.
type CsvBookingExtras struct {
	Typ        string             // ER, AR, GV
	CostCenter string             // JM, AN, K, usw.
	Net        map[string]float64 // Verteilung der netto Rechnungspositionen auf Stakeholder
}

type Booking struct {
	Amount float64
	Text   string
	Year   int
	Month  int

	Extras CsvBookingExtras
}

func (b Booking) Print(account string) {
	text := b.Text
	if len(text) > 37 {
		text = text[:37] + "..."
	}

	fmt.Printf("[%s: %2d-%d %-40s \t %9.2f]\n", account, b.Month, b.Year, text, b.Amount)
}

type ByMonth []Booking

func (a ByMonth) Len() int           { return len(a) }
func (a ByMonth) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByMonth) Less(i, j int) bool { return a[i].Month < a[j].Month }
