package booking

import (
	"fmt"
	"time"

	"github.com/ahojsenn/kontrol/util"

	"golang.org/x/text/language"

	"golang.org/x/text/message"
)

var ValidBookingTypes = [...]string{
	"ER",
	"AR",
	"GV",
	"GV-Vorjahr",
	"Reisekosten",
	"RKE",
	"RK-Erstattung",
	"IS",
	"SV-Beitrag",
	"GWSteuer",
	"Gehalt",
	"LNSteuer",
	"Rückstellung",
	"Anfangsbestand",
	"SKR03",
	"openingBalance",
	"closingBalance",
	"UstVZ",
}

// Zusatzinformationen einer Buchung, deren Quelle die CSV-Datei ist, und die für die weitere
// Bearbeitung erforderlich sind.
type CsvBookingExtras struct {
	// "ER", "AR", "GV", "IS", "SV-Beitrag", "CC_GWSteuer"
	Typ string

	// Wer für die Buchung verantwortlich ist. Unterschiedliche Bedeutung für unterschiedliche Buchungsarten:
	// - "ER": Wer die Kosten verursacht hat
	// - "AR": Wer den Auftrag gebracht hat.
	// - "GV": Wer die CC_Entnahme getätigt hat
	// - "IS": Wer die internen Stunden geleistet hat
	// - "SV-Beitrag": Für wen der SV-Beitrag gezhalt wurde
	Responsible string

	// Verteilung der netto Rechnungspositionen auf Stakeholder
	// Net map[valueMagnets.Stakeholder]float64
	Net map[string]float64
}

// Aus einer Buchung in der Quelldatei wird eine oder mehrere virtuelle Buchungen erstellt. Dies ist die Liste
// möglicher Werte für den Type einer virtuellen Buchung
const (
	Erloese                  = "Erloese"
	CC_Vertriebsprovision    = "CC_Vertriebsprovision"
	CC_Nettoanteil           = "CC_Nettoanteil"
	CC_Employeeaanteil       = "CC_Employeeanteil"
	CC_Kommitmentanteil      = "CC_Kommitmentanteil"
	CC_KommitmentanteilEX    = "CC_KommitmentanteilEX"
	CC_KommitmentanteilREST  = "CC_KommitmentanteilREST"
	CC_Entnahme              = "CC_Entnahme"
	CC_AnteilAusFairshares   = "CC_AnteilAusFairshares"
	CC_AnteilAusFaktura      = "CC_AnteilAusFaktura"
	CC_Fakturasumme          = "CC_Fakturasumme"
	CC_RevDistribution_1     = "CC_Erlösverteilung Schritt 1"
	CC_KommitmenschDarlehen  = "CC_KommitmenschDarlehen"
	Eingangsrechnung         = "Eingangsrechnung"
	CC_SVBeitrag             = "CC_SV-Beitrag"
	CC_GWSteuer              = "CC_GWSteuer"
	CC_Gehalt                = "CC_Gehalt"
	CC_J_Bonus               = "CC_Jahresüberschuss/Bonus"
	CC_LNSteuer              = "CC_LNSteuer"
	CC_PartnerNettoFaktura   = "CC_PartnerNettofaktura"
	CC_LiquidityReserve      = "CC_Liquiditätsreserve"
	CC_Kostenrueckerstattung = "CC_Kostenrueckerstattung"
	CC_Anlagenzugang         = "CC_Anlagenzugang"
	SKR03                    = "SKR03"
	UstVZ                    = "UstVZ"
	Ust                      = "Ust"
	Kosten                   = "Kosten"
)

type Booking struct {
	Id          int
	RowNr       int
	Type        string // siehe const-Block hier drüber für gültige Werte
	Soll        string
	Haben       string
	CostCenter  string
	Project     string
	Amount      float64
	Text        string
	Year        int
	Month       int
	FileCreated time.Time
	BankCreated time.Time

	CsvBookingExtras `json:"Net"`
}

func NewBooking(
	rownr int,
	csvType string,
	soll string,
	haben string,
	costCenter string,
	project string,
	//net map[valueMagnets.Stakeholder]float64,
	net map[string]float64,
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
		Id:               util.GetNewBookingId(),
		RowNr:            rownr,
		CsvBookingExtras: extraCsv,
		Type:             csvType,
		Soll:             soll,
		Haben:            haben,
		CostCenter:       costCenter,
		Project:          project,
		Amount:           amount,
		Text:             text,
		Month:            month,
		Year:             year,
		BankCreated:      bankCreated,
	}
}

func CloneBooking(b Booking, amount float64, typ string, costcenter string, soll string, haben string, project string) Booking {
	return Booking{
		Id:          util.GetNewBookingId(),
		Amount:      amount,
		Type:        typ,
		RowNr:       b.RowNr,
		Text:        b.Text,
		Month:       b.Month,
		Year:        b.Year,
		FileCreated: b.FileCreated,
		BankCreated: b.BankCreated,
		CostCenter:  costcenter,
		Soll:        soll,
		Haben:       haben,
		Project:     project,
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

// is this an Open Position?
func (b *Booking) IsOpenPosition() bool {
	emptyTime := time.Time{}
	return b.BankCreated == emptyTime
}

// ist this booking beyond the budget date?
func (b *Booking) IsBeyondBudgetDate() bool {
	return b.BankCreated.After(util.Global.BalanceDate)
}

type ByMonth []Booking

func (a ByMonth) Len() int           { return len(a) }
func (a ByMonth) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByMonth) Less(i, j int) bool { return a[i].Month < a[j].Month }

type ByRowNr []Booking

func (a ByRowNr) Len() int           { return len(a) }
func (a ByRowNr) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByRowNr) Less(i, j int) bool { return a[i].RowNr < a[j].RowNr }
