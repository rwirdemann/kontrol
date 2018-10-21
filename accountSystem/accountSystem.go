package accountSystem

import (
	"fmt"
	"log"

	"github.com/ahojsenn/kontrol/valueMagnets"
	"github.com/ahojsenn/kontrol/account"
	"github.com/ahojsenn/kontrol/booking"
	"github.com/ahojsenn/kontrol/util"
)




type AccountSystem interface {
	GetCollectiveAccount() *account.Account
	Add(a *account.Account)
	All() []account.Account
	Get(id string) (*account.Account, bool)
	GetBilanzAccounts (typ string) []account.Account
	GetSKR03(id string) *account.Account
	ClearBookings()
}

type DefaultAccountSystem struct {
	collectiveAccount *account.Account
	accounts          map[string]*account.Account
}

const SKR03 = "SKR03"

// Bilanzkonten
var SKR03_Anlagen25 = account.AccountDescription{Id: "SKR03_Anlagen25", Name: "01_SKR03_25_Anlagen Ähnl.R&W", Type: account.KontenartAktiv}
var SKR03_Anlagen = account.AccountDescription{Id: "SKR03_Anlagen", Name: "02_SKR03_410-480_Anlagen", Type: account.KontenartAktiv}
var SKR03_1400 = account.AccountDescription{Id: "1400", Name: "03_SKR03_1400_OPOS-Kunde", Type: account.KontenartAktiv}
var SKR03_Kautionen = account.AccountDescription{Id: "SKR03_Kautionen", Name: "04_SKR03_1525_Kautionen", Type: account.KontenartAktiv}
var SKR03_Vorsteuer = account.AccountDescription{Id: "SKR03_Vorsteuer", Name: "06_05_SKR03_05_1570-1579_Steuer: Vorsteuer", Type: account.KontenartAktiv}
var SKR03_Umsatzsteuer = account.AccountDescription{Id: "SKR03_Umsatzsteuer", Name: "07_SKR03_6_1770_Steuer: Umsatzsteuer", Type: account.KontenartAktiv}
var SKR03_1200 = account.AccountDescription{Id: "1200", Name: "08_SKR03_1200_Bank", Type: account.KontenartAktiv}

var SKR03_900_Haftkapital = account.AccountDescription{Id: "SKR03_900_Haftkapital", Name: "08_SKR03_900_Haftkapital", Type: account.KontenartPassiv}
var SKR03_Eigenkapital_880 = account.AccountDescription{Id: "Eigenkapital", Name: "09_SKR03_880_Eigenkapital", Type: account.KontenartPassiv}
var SKR03_Rueckstellungen = account.AccountDescription{Id: "Rückstellung", Name: "10_SKR03_10_956-977_Rückstellung", Type: account.KontenartPassiv}
var SKR03_920_Gesellschafterdarlehen = account.AccountDescription{Id: "SKR03_920_Gesellschafterdarlehen", Name: "11_SKR03_920_Gesellschafterdarlehen", Type: account.KontenartPassiv}
var SKR03_1600 = account.AccountDescription{Id: "1600", Name: "12_SKR03_1600_OPOS-Lieferant", Type: account.KontenartPassiv}
var SKR03_1900 = account.AccountDescription{Id: "1900", Name: "13_SKR03_1900_Privatentnahmen", Type: account.KontenartPassiv}

// Erfolgskonten
var SKR03_Umsatzerloese = account.AccountDescription{Id: "SKR03_Umsatzerloese", Name: "1 SKR03_Umsatzerloese 8100-8402", Type: account.KontenartErtrag}
var SKR03_4100_4199 = account.AccountDescription{Id: "4100_4199", Name: "3 Löhne und Gehälter 4100-4199", Type: account.KontenartAufwand}
var SKR03_AnlagenabgaengeSachanlagen2310 = account.AccountDescription{Id: "SKR03_AnlagenabgängeSachanlagen", Name: "4 AnlagenabgängeSachanlagen 2310", Type: account.KontenartAufwand}
var SKR03_Abschreibungen = account.AccountDescription{Id: "SKR03_Abschreibungen", Name: "4 Abschreibungen auf Anlagen 4822-4855", Type: account.KontenartAufwand}
var SKR03_sonstigeAufwendungen = account.AccountDescription{Id: "SKR03_sonstigeAufwendungen", Name: "5 sonstige Aufwendungen", Type: account.KontenartAufwand}
var SKR03_Steuern = account.AccountDescription{Id: "SKR03_Steuern", Name: "6 SKR03_Steuern 4320", Type: account.KontenartAufwand}
var ErgebnisNachSteuern = account.AccountDescription{Id: "SKR03_ErgebnisNachSteuern", Name: "7 ErgebnisNachSteuern", Type: account.KontenartVerrechnung}
// Verrechnungskonten
var SKR03_Saldenvortrag = account.AccountDescription{Id: "SKR03_Saldenvortrag", Name: "Saldenvortrag 9000", Type: account.KontenartVerrechnung}
var SummeAktiva 	= account.AccountDescription{Id: "SummeAktiva", Name: "V: SummeAktiva", Type: account.KontenartVerrechnung}
var SummePassiva 	= account.AccountDescription{Id: "SummePassiva", Name: "V: SummePassiva", Type: account.KontenartVerrechnung}
var AlleKLRBuchungen = account.AccountDescription{Id: "AlleKLRBuchungen", Name: "V: AlleKLRBuchungen", Type: account.KontenartKLR}





type Accountlist struct {
}

func (this Accountlist) All() []account.AccountDescription {
	return []account.AccountDescription{
		SKR03_Rueckstellungen,
		SKR03_Eigenkapital_880,
		SKR03_900_Haftkapital,
		SKR03_920_Gesellschafterdarlehen,
		SKR03_1200,
		SKR03_1400,
		SKR03_1600,
		SKR03_1900,
		SKR03_4100_4199,
		SKR03_AnlagenabgaengeSachanlagen2310,
		SKR03_sonstigeAufwendungen,
		SKR03_Anlagen,
		SKR03_Anlagen25,
		SKR03_Abschreibungen,
		SKR03_Kautionen,
		SKR03_Umsatzerloese,
		SKR03_Steuern,
		SKR03_Vorsteuer,
		SKR03_Umsatzsteuer,
		SKR03_Saldenvortrag,
		ErgebnisNachSteuern,
		SummeAktiva,
		SummePassiva,
		AlleKLRBuchungen,
	}
}

func EmptyDefaultAccountSystem() AccountSystem {
	o := account.AccountDescription{Id: "all", Name: "Hauptbuch", Type: account.KontenartVerrechnung}
	return &DefaultAccountSystem{collectiveAccount: &account.Account{Description: o}, accounts: make(map[string]*account.Account)}
}

func NewDefaultAccountSystem() AccountSystem {
	year := util.Global.FinancialYear

	ad := account.AccountDescription{Id: "all", Name: "Hauptbuch", Type: account.KontenartVerrechnung}
	accountSystem := DefaultAccountSystem{collectiveAccount: &account.Account{Description: ad}, accounts: make(map[string]*account.Account)}

	// generate accounts according to the AccountList
	accountlist := Accountlist{}
	for _, a := range accountlist.All() {
		accountSystem.Add(account.NewAccount(a))
	}

	// generate accounts for all stakeholders
	stakeholderRepository := valueMagnets.StakeholderRepository{}
	for _, sh := range stakeholderRepository.All(year) {
			ad := account.AccountDescription{Id: sh.Id, Name: sh.Name, Type: sh.Type}
			accountSystem.Add(account.NewAccount(ad))
	}
	return &accountSystem
}

func (r *DefaultAccountSystem) GetCollectiveAccount() *account.Account {
	return r.collectiveAccount
}

func (r *DefaultAccountSystem) Add(a *account.Account) {
	r.accounts[a.Description.Id] = a
}

func (r *DefaultAccountSystem) All() []account.Account {
	result := make([]account.Account, 0, len(r.accounts))
	for _, a := range r.accounts {
		a.UpdateSaldo()
		clone := *a
		clone.Bookings = nil
		result = append(result, clone)
	}
	return result
}

func (r *DefaultAccountSystem) Get(id string) (*account.Account, bool) {
	if a, ok := r.accounts[id]; ok {
		return a, true
	}
	return nil, false
}

func (as *DefaultAccountSystem) GetBilanzAccounts (typ string) []account.Account {
	var filtered  []account.Account
	for _, account := range as.accounts {
		if account.Description.Type == typ {
			account.UpdateSaldo()
			clone := *account
			clone.Bookings = nil
			filtered = append(filtered, clone)
		}
	}
	return filtered
}

func (r *DefaultAccountSystem) ClearBookings() {
	r.collectiveAccount.Bookings = []booking.Booking{}
	for _, account := range r.accounts {
		account.Bookings = []booking.Booking{}
	}
}


func (r *DefaultAccountSystem) GetSKR03(SKR03konto string) *account.Account {
	var account *account.Account
	switch SKR03konto {
	case "25": // Anlage buchen
		account = r.accounts[SKR03_Anlagen25.Id]
	case "410": // Anlage buchen
		account = r.accounts[SKR03_Anlagen.Id]
	case "480": // Anlage buchen
		account = r.accounts[SKR03_Anlagen.Id]
	case "880": // Eigenkapital bilden
		account = r.accounts[SKR03_Eigenkapital_880.Id]
	case "900": // Eigenkapital bilden
		account = r.accounts[SKR03_900_Haftkapital.Id]
	case "920": // Rückstellung bilden
		account = r.accounts[SKR03_920_Gesellschafterdarlehen.Id]
	case "956","965","970","977": // Rückstellung bilden
		account = r.accounts[SKR03_Rueckstellungen.Id]
	case "1200": // Bank buchen
		account = r.accounts[SKR03_1200.Id]
	case "1525":
		account = r.accounts[SKR03_Kautionen.Id]
	case "1548":
		account = r.accounts[SKR03_Vorsteuer.Id]
	case "1400":
		account = r.accounts[SKR03_1400.Id]
	case "731", "1600":
		account = r.accounts[SKR03_1600.Id]
	case "1770":
		account = r.accounts[SKR03_Umsatzsteuer.Id]
	case "1900":
		account = r.accounts[SKR03_1900.Id]
	case "2310":
		account = r.accounts[SKR03_AnlagenabgaengeSachanlagen2310.Id]
	case "4120":
		account = r.accounts[SKR03_4100_4199.Id]
	case "4130", "4138", "4140":
		account = r.accounts[SKR03_4100_4199.Id]
	case "4320":
		account = r.accounts[SKR03_Steuern.Id]
	case "4822", "4830", "4855":
		account = r.accounts[SKR03_Abschreibungen.Id]
	case "4200", "4210", "4250", "4260",
		"4360", "4380", "4390", "4806", "4595", "4600", "4640", "4650", "4654", "4655", "4663", "4664", "4666", "4670", "4672", "4673", "4674", "4676", "4780", "2300", "4900", "4909", "4910", "4920", "4921", "4925", "4930", "4940", "4945", "4949", "4950", "4955", "4957", "4960", "4964", "4970", "4980":
		account = r.accounts[SKR03_sonstigeAufwendungen.Id]
	case "8100":
		account = r.accounts[SKR03_Umsatzerloese.Id]
	case "9000":
		account = r.accounts[SKR03_Saldenvortrag.Id]
	default:
		log.Printf("GetSKR03: could not process booking type '%s'", SKR03konto)
		panic(fmt.Sprintf("GetSKR03: SKR03Bucket/Stakeholder/Konto '%s' not found", account.Description))
	}
	return account
}
