package booking

import (
	"fmt"
	"time"

	"golang.org/x/text/language"

	"github.com/ahojsenn/kontrol/owner"
	"golang.org/x/text/message"
)

var ValidBookingTypes = [...]string{"ER", "AR", "GV", "GV-Vorjahr", "IS", "SV-Beitrag", "GWSteuer", "Gehalt", "LNSteuer", "Rückstellung", "Anfangsbestand", "RückstellungAuflösen", "ERgegenRückstellung", "SKR03"}

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
	Net map[owner.Stakeholder]float64
}

// Aus einer Buchung in der Quelldatei wird eine oder mehrere virtuelle Buchungen erstellt. Dies ist die Liste
// möglicher Werte für den Type einer virtuellen Buchung
const (
	Vertriebsprovision    = "Vertriebsprovision"
	Nettoanteil           = "Nettoanteil"
	Kommitmentanteil      = "Kommitmentanteil"
	Entnahme              = "Entnahme"
	Eingangsrechnung      = "Eingangsrechnung"
	InterneStunden        = "Interne Stunden"
	SVBeitrag             = "SV-Beitrag"
	GWSteuer              = "GWSteuer"
	Gehalt                = "Gehalt"
	LNSteuer              = "LNSteuer"
	Rueckstellung         = "Rueckstellung"
	RueckstellungAuflösen = "RueckstellungAufloesen"
	Anfangsbestand        = "Anfangsbestand"
	GVVorjahr             = "GVVorjahr"
	SKR03                 = "SKR03"
)

type Booking struct {
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
	csvType string,
	soll string,
	haben string,
	dealBringer string,
	net map[owner.Stakeholder]float64,
	amount float64,
	text string,
	month int,
	year int,
	bankCreated time.Time) *Booking {

	return &Booking{
		CsvBookingExtras: CsvBookingExtras{
			Typ:         csvType,
			Responsible: dealBringer,
			Net:         net,
		},
		Soll:        soll,
		Haben:       haben,
		Amount:      amount,
		Text:        text,
		Month:       month,
		Year:        year,
		BankCreated: bankCreated,
	}
}

func Ausgangsrechnung(
	dealbringer string,
	net map[owner.Stakeholder]float64,
	amount float64,
	text string,
	month int,
	year int,
	bankCreated time.Time) *Booking {

	return NewBooking("AR", "", "", dealbringer, net, 17225.25, "Rechnung 1234", 1, 2017, bankCreated)
}

func CloneBooking(b Booking, amount float64, typ string, costcenter string) Booking {
	return Booking{
		Amount:      amount,
		Type:        typ,
		Text:        b.Text,
		Month:       b.Month,
		Year:        b.Year,
		FileCreated: b.FileCreated,
		BankCreated: b.BankCreated,
		CostCenter:  costcenter,
	}
}

func (b Booking) Print(owner owner.Stakeholder) {
	text := b.Text
	if len(text) > 37 {
		text = text[:37] + "..."
	}

	fmt.Printf("[%s: %2d-%d %2s %-22s %-40s \t %9.2f]\n", owner.Id, b.Month, b.Year, b.CostCenter, b.Type, text, b.Amount)
}

func (b Booking) CSV(owner owner.Stakeholder) string {
	p := message.NewPrinter(language.German)
	text := b.Text
	if len(text) > 37 {
		text = text[:37] + "..."
	}
	amount := p.Sprintf("%9.2f", b.Amount)
	return fmt.Sprintf("%s;%2d;%d;%s;%s;%s;%s\n", owner.Id, b.Month, b.Year, b.CostCenter, b.Type, text, amount)
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
