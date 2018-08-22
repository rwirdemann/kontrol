package accountSystem

import (
	"fmt"
	"log"

	"github.com/ahojsenn/kontrol/owner"
	"github.com/ahojsenn/kontrol/account"
	"github.com/ahojsenn/kontrol/booking"
	"github.com/ahojsenn/kontrol/util"
)


const (
	KontenartAktiv            = "Aktivkonto"
	KontenartPassiv           = "Passivkonto"
	KontenartAufwand          = "Aufwandskonto"
	KontenartErtrag           = "Ertragskonto"
	KontenartVerrechnung      = "Verrechnungskonto"
)

type AccountSystem interface {
	BankAccount() *account.Account
	Add(a *account.Account)
	All() []account.Account
	Get(id string) (*account.Account, bool)
	GetSKR03(id string) *account.Account
	ClearBookings()
}

type DefaultAccountSystem struct {
	collectiveAccount *account.Account
	accounts          map[string]*account.Account
}

const SKR03 = "SKR03"

var SKR03_Rueckstellungen = account.AccountDescription{Id: "Rückstellung", Name: "Rückstellung 956-977", Type: KontenartPassiv}
var SKR03_Eigenkapital_880 = account.AccountDescription{Id: "Eigenkapital", Name: "Eigenkapital 880", Type: KontenartPassiv}
var SKR03_KontoJUSVJ = account.AccountDescription{Id: "JahresüberschussVJ", Name: "JahresüberschussVJ Gesellschafterdarlehen 920", Type: KontenartPassiv}
var SKR03_1400 = account.AccountDescription{Id: "1400", Name: "OPOS-Kunde 1400", Type: KontenartAktiv}
var SKR03_1600 = account.AccountDescription{Id: "1600", Name: "OPOS-Lieferant 1600", Type: KontenartPassiv}
var SKR03_Anlagen = account.AccountDescription{Id: "SKR03_Anlagen", Name: "Zugang Anlagen", Type: KontenartAktiv}
var SKR03_Anlagen25 = account.AccountDescription{Id: "SKR03_Anlagen25", Name: "Zugang Anlagen Ähnl.R&W 25", Type: KontenartAktiv}
var SKR03_Kautionen = account.AccountDescription{Id: "SKR03_Kautionen", Name: "SKR03_Kautionen 1525", Type: KontenartAktiv}
var SKR03_Vorsteuer = account.AccountDescription{Id: "SKR03_Vorsteuer", Name: "Steuer: Vorsteuer 1570-1579", Type: KontenartAktiv}
var SKR03_Umsatzsteuer = account.AccountDescription{Id: "SKR03_Umsatzsteuer", Name: "Steuer: Umsatzsteuer 1770", Type: KontenartAktiv}
// Erfolgskonten
var SKR03_Umsatzerloese = account.AccountDescription{Id: "SKR03_Umsatzerloese", Name: "1 SKR03_Umsatzerloese 8100-8402", Type: KontenartErtrag}
var SKR03_4100_4199 = account.AccountDescription{Id: "4100_4199", Name: "3 Löhne und Gehälter 4100-4199", Type: KontenartAufwand}
var SKR03_Abschreibungen = account.AccountDescription{Id: "SKR03_Abschreibungen", Name: "4 Abschreibungen auf Anlagen 4822-4855", Type: KontenartAufwand}
var SKR03_sonstigeAufwendungen = account.AccountDescription{Id: "SKR03_sonstigeAufwendungen", Name: "5 sonstige Aufwendungen", Type: KontenartAufwand}
var SKR03_Steuern = account.AccountDescription{Id: "SKR03_Steuern", Name: "6 SKR03_Steuern 4320", Type: KontenartAufwand}
var ErgebnisNachSteuern = account.AccountDescription{Id: "SKR03_ErgebnisNachSteuern", Name: "7 ErgebnisNachSteuern", Type: KontenartVerrechnung}
// Verrechnungskonten
var SKR03_Saldenvortrag = account.AccountDescription{Id: "SKR03_Saldenvortrag", Name: "Saldenvortrag 9000", Type: KontenartVerrechnung}

type Accountlist struct {
}

func (this Accountlist) All() []account.AccountDescription {
	return []account.AccountDescription{
		SKR03_Rueckstellungen,
		SKR03_Eigenkapital_880,
		SKR03_KontoJUSVJ,
		SKR03_1400,
		SKR03_1600,
		SKR03_4100_4199,
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
	}
}

func EmptyDefaultAccountSystem() AccountSystem {
	o := account.AccountDescription{Id: "GLS", Name: "Kommitment GmbH & Co. KG", Type: KontenartAktiv}
	return &DefaultAccountSystem{collectiveAccount: &account.Account{Description: o}, accounts: make(map[string]*account.Account)}
}

func NewDefaultAccountSystem() AccountSystem {
	year := util.Global.FinancialYear

	ad := account.AccountDescription{Id: "GLS", Name: "Kommitment GmbH & Co. KG", Type: KontenartAktiv}
	accountSystem := DefaultAccountSystem{collectiveAccount: &account.Account{Description: ad}, accounts: make(map[string]*account.Account)}

	// generate accounts according to the AccountList
	accountlist := Accountlist{}
	for _, a := range accountlist.All() {
		accountSystem.Add(account.NewAccount(a))
	}

	// generate accounts for all stakeholders
	stakeholderRepository := owner.StakeholderRepository{}
	for _, sh := range stakeholderRepository.All(year) {
		if sh.Type != owner.StakeholderTypeOthers {
			ad := account.AccountDescription{Id: sh.Id, Name: sh.Name, Type: sh.Type}
			accountSystem.Add(account.NewAccount(ad))
		}
	}
	return &accountSystem
}

func (r *DefaultAccountSystem) BankAccount() *account.Account {
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
	case "920": // Rückstellung bilden
		account = r.accounts[SKR03_KontoJUSVJ.Id]
	case "956","965","970","977": // Rückstellung bilden
		account = r.accounts[SKR03_Rueckstellungen.Id]
	case "1525":
		account = r.accounts[SKR03_Kautionen.Id]
	case "4120":
		account = r.accounts[SKR03_4100_4199.Id]
	case "4130", "4138", "4140":
		account = r.accounts[SKR03_4100_4199.Id]
	case "1200":
		account = r.BankAccount()
	case "4320":
		account = r.accounts[SKR03_Steuern.Id]
	case "4822", "4830", "4855":
		account = r.accounts[SKR03_Abschreibungen.Id]
	case "4200", "4210", "4250", "4260",
		"4360", "4380", "4390", "4806", "4595", "4600", "4640", "4650", "4654", "4655", "4663", "4664", "4666", "4670", "4672", "4673", "4674", "4676", "4780", "2300", "4900", "4909", "4910", "4920", "4921", "4925", "4930", "4940", "4945", "4949", "4950", "4955", "4957", "4960", "4964", "4970", "4980":
		account = r.accounts[SKR03_sonstigeAufwendungen.Id]
	case "9000":
		account = r.accounts[SKR03_Saldenvortrag.Id]
	default:
		log.Printf("GetSKR03: could not process booking type '%s'", SKR03konto)
		panic(fmt.Sprintf("GetSKR03: SKR03Bucket/Stakeholder/Konto '%s' not found", account.Description))
	}
	return account
}
