package booking

import (
	"fmt"
	"time"

	"golang.org/x/text/language"

	"github.com/ahojsenn/kontrol/valueMagnets"
	"golang.org/x/text/message"
)

var ValidBookingTypes = [...]string{
	"ER",
	"AR",
	"GV",
	"GV-Vorjahr",
	"IS",
	"SV-Beitrag",
	"GWSteuer",
	"Gehalt",
	"LNSteuer",
	"Rückstellung",
	"Anfangsbestand",
 	"SKR03",
	"UstVZ",
}

// Zusatzinformationen einer Buchung, deren Quelle die CSV-Datei ist, und die für die weitere
// Bearbeitung erforderlich sind.
type CsvBookingExtras struct {
	// "ER", "AR", "GV", "IS", "SV-Beitrag", "GWSteuer"
	Typ string

	// Wer für die Buchung verantwortlich ist. Unterschiedliche Bedeutung für unterschiedliche Buchungsarten:
	// - "ER": Wer die Kosten verursacht hat
	// - "AR": Wer den Auftrag gebracht hat.
	// - "GV": Wer die Entnahme getätigt hat
	// - "IS": Wer die internen Stunden geleistet hat
	// - "SV-Beitrag": Für wen der SV-Beitrag gezhalt wurde
	Responsible string

	// Verteilung der netto Rechnungspositionen auf Stakeholder
	Net map[valueMagnets.Stakeholder]float64
}

// Aus einer Buchung in der Quelldatei wird eine oder mehrere virtuelle Buchungen erstellt. Dies ist die Liste
// möglicher Werte für den Type einer virtuellen Buchung
const (
	Erloese      		  = "Erloese"
	Vertriebsprovision    = "Vertriebsprovision"
	Nettoanteil           = "Nettoanteil"
	Employeeaanteil		  = "Employeeanteil"
	Kommitmentanteil      = "Kommitmentanteil"
	Entnahme              = "Entnahme"
	Eingangsrechnung      = "Eingangsrechnung"
	InterneStunden        = "Interne Stunden"
	SVBeitrag             = "SV-Beitrag"
	GWSteuer              = "GWSteuer"
	Gehalt                = "Gehalt"
	LNSteuer              = "LNSteuer"
	Rueckstellung         = "Rueckstellung"
	Anfangsbestand        = "Anfangsbestand"
	GVVorjahr             = "GVVorjahr"
	SKR03                 = "SKR03"
	UstVZ				  = "UstVZ"
	Kosten				  = "Kosten"
)

type Booking struct {
	RowNr 		int
	Type        string // siehe const-Block hier drüber für gültige Werte
	Soll        string
	Haben       string
	CostCenter  string
	Amount      float64
	Text        string
	Year        int
	Month       int
	FileCreated time.Time
	BankCreated time.Time

	CsvBookingExtras `json:"-"`
}

func NewBooking(
	rownr int,
	csvType string,
	soll string,
	haben string,
	costCenter string,
	net map[valueMagnets.Stakeholder]float64,
	amount float64,
	text string,
	month int,
	year int,
	bankCreated time.Time) *Booking {

		extraCsv := CsvBookingExtras{
			Typ:         csvType,
			Responsible: costCenter,
			Net:         net,
		}


	return &Booking{
		RowNr:		 rownr,
		CsvBookingExtras: extraCsv,
		Type:        csvType,
		Soll:        soll,
		Haben:       haben,
		CostCenter:  costCenter,
		Amount:      amount,
		Text:        text,
		Month:       month,
		Year:        year,
		BankCreated: bankCreated,
	}
}

func Ausgangsrechnung(
	rownr int,
	dealbringer string,
	net map[valueMagnets.Stakeholder]float64,
	amount float64,
	text string,
	month int,
	year int,
	bankCreated time.Time) *Booking {

	return NewBooking(rownr,"AR", "", "", dealbringer, net, 17225.25, "Rechnung 1234", 1, 2017, bankCreated)
}

func CloneBooking(b Booking, amount float64, typ string, costcenter string, soll string, haben string) Booking {
	return Booking{
		Amount:      amount,
		Type:        typ,
		RowNr:       b.RowNr,
		Text:        b.Text,
		Month:       b.Month,
		Year:        b.Year,
		FileCreated: b.FileCreated,
		BankCreated: b.BankCreated,
		CostCenter:  costcenter,
		Soll: soll,
		Haben: haben,
	}
}

func (b Booking) Print(id string) {
	text := b.Text
	if len(text) > 37 {
		text = text[:37] + "..."
	}

	fmt.Printf("[%s: %2d-%d %2s %-22s %-40s \t %9.2f]\n", id, b.Month, b.Year, b.CostCenter, b.Type, text, b.Amount)
}

func (b Booking) CSV(id string) string {
	p := message.NewPrinter(language.German)
	text := b.Text
	if len(text) > 37 {
		text = text[:37] + "..."
	}
	amount := p.Sprintf("%9.2f", b.Amount)
	return fmt.Sprintf("%s;%2d;%d;%s;%s;%s;%s\n", id, b.Month, b.Year, b.CostCenter, b.Type, text, amount)
}

func (b *Booking) BookOnBankAccount() bool {
	if b.Typ == "IS" {
		return false
	}
	return true
}

type ByMonth []Booking

func (a ByMonth) Len() int           { return len(a) }
func (a ByMonth) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByMonth) Less(i, j int) bool { return a[i].Month < a[j].Month }

type ByRowNr []Booking
func (a ByRowNr) Len() int           { return len(a) }
func (a ByRowNr) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByRowNr) Less(i, j int) bool { return a[i].RowNr < a[j].RowNr }
