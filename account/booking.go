package account

import (
	"fmt"

	"bitbucket.org/rwirdemann/kontrol/owner"
)

var ValidBookingTypes = [...]string{"ER", "AR", "GV", "IS", "SV-Beitrag"}

// Zusatzinformationen einer Booking, deren Quelle die CSV-Datei ist, und die für die weitere
// Bearbeitung erforderlich sind.
type CsvBookingExtras struct {
	SourceType string                        // siehe ValidBookingTypes for valid values
	CostCenter string                        // JM, AN, K, usw.
	Net        map[owner.Stakeholder]float64 // Verteilung der netto Rechnungspositionen auf Stakeholder
}

// Aus einer Bankbuchung wird eine oder mehrere virtuelle Buchungen erstellt. Dis ist die Liste
// mögliche Werte für den DestType einer virtuellen Buchung
const (
	Vertriebsprovision = "Vertriebsprovision"
	Nettoanteil        = "Nettoanteil"
	Kommitmentanteil   = "Kommitmentanteil"
	Entnahme           = "Entnahme"
	Eingangsrechnung   = "Eingangsrechnung"
	InterneStunden     = "Interne Stunden"
)

type Booking struct {
	DestType string // siehe const-Block hier drüber für gültige Werte
	Amount   float64
	Text     string
	Year     int
	Month    int

	Extras CsvBookingExtras `json:"-"`
}

func (b Booking) Print(owner owner.Stakeholder) {
	text := b.Text
	if len(text) > 37 {
		text = text[:37] + "..."
	}

	fmt.Printf("[%s: %2d-%d %-22s %-40s \t %9.2f]\n", owner.Id, b.Month, b.Year, b.DestType, text, b.Amount)
}

type ByMonth []Booking

func (a ByMonth) Len() int           { return len(a) }
func (a ByMonth) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByMonth) Less(i, j int) bool { return a[i].Month < a[j].Month }
