package booking

import (
	"fmt"

	"bitbucket.org/rwirdemann/kontrol/owner"
)

var ValidBookingTypes = [...]string{"ER", "AR", "GV", "IS", "SV-Beitrag", "GWSteuer"}

// Zusatzinformationen einer Booking, deren Quelle die CSV-Datei ist, und die für die weitere
// Bearbeitung erforderlich sind.
type CsvBookingExtras struct {
	CSVType    string                        // siehe ValidBookingTypes for valid values
	CostCenter string                        // JM, AN, K, usw.
	Net        map[owner.Stakeholder]float64 // Verteilung der netto Rechnungspositionen auf Stakeholder
}

// Aus einer Bankbuchung wird eine oder mehrere virtuelle Buchungen erstellt. Dis ist die Liste
// mögliche Werte für den Type einer virtuellen Buchung
const (
	Vertriebsprovision = "Vertriebsprovision"
	Nettoanteil        = "Nettoanteil"
	Kommitmentanteil   = "Kommitmentanteil"
	Entnahme           = "Entnahme"
	Eingangsrechnung   = "Eingangsrechnung"
	InterneStunden     = "Interne Stunden"
	SVBeitrag          = "SV-Beitrag"
	GWSteuer           = "GWSteuer"
)

type Booking struct {
	Type   string // siehe const-Block hier drüber für gültige Werte
	Amount float64
	Text   string
	Year   int
	Month  int

	CsvBookingExtras `json:"-"`
}

func NewBooking(
	csvType string,
	costCenter string,
	net map[owner.Stakeholder]float64,
	amount float64,
	text string,
	month int,
	year int) *Booking {

	return &Booking{
		CsvBookingExtras: CsvBookingExtras{
			CSVType:    csvType,
			CostCenter: costCenter,
			Net:        net,
		},
		Amount: amount,
		Text:   text,
		Month:  month,
		Year:   year,
	}
}

func (b Booking) Print(owner owner.Stakeholder) {
	text := b.Text
	if len(text) > 37 {
		text = text[:37] + "..."
	}

	fmt.Printf("[%s: %2d-%d %-22s %-40s \t %9.2f]\n", owner.Id, b.Month, b.Year, b.Type, text, b.Amount)
}

func (b *Booking) BookOnBankAccount() bool {
	if b.CSVType == "IS" {
		return false
	}
	return true
}

type ByMonth []Booking

func (a ByMonth) Len() int           { return len(a) }
func (a ByMonth) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByMonth) Less(i, j int) bool { return a[i].Month < a[j].Month }
