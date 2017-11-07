package kontrol

import (
	"fmt"
)

// Zusatzinformationen einer Buchung, deren Quelle die CSV-Datei ist, und die für die weitere
// Bearbeitung erforderlich sind.
type CsvBookingExtras struct {
	Typ        string                  // ER, AR, GV
	CostCenter string                  // JM, AN, K, usw.
	Net        map[Stakeholder]float64 // Verteilung der netto Rechnungspositionen auf Stakeholder
}

const (
	Vertriebsprovision = "Vertriebsprovision"
	Nettoanteil        = "Nettoanteil"
	Kommitmentanteil   = "Kommitmentanteil"
	Entnahme           = "Entnahme"
)

type Booking struct {
	Typ    string // Vertriebsprovision, Nettoanteil, Kommitmentanteil
	Amount float64
	Text   string
	Year   int
	Month  int

	Extras CsvBookingExtras `json:"-"`
}

func (b Booking) Print(owner Stakeholder) {
	text := b.Text
	if len(text) > 37 {
		text = text[:37] + "..."
	}

	fmt.Printf("[%s: %2d-%d %-22s %-40s \t %9.2f]\n", owner.Id, b.Month, b.Year, b.Typ, text, b.Amount)
}

type ByMonth []Booking

func (a ByMonth) Len() int           { return len(a) }
func (a ByMonth) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByMonth) Less(i, j int) bool { return a[i].Month < a[j].Month }
