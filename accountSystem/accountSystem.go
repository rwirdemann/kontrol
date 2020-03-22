package accountSystem

import (
	"fmt"
	"github.com/ahojsenn/kontrol/account"
	"github.com/ahojsenn/kontrol/booking"
	"github.com/ahojsenn/kontrol/util"
	"github.com/ahojsenn/kontrol/valueMagnets"
	"log"
	"sort"
	"strconv"
)

type AccountSystem interface {
	GetCollectiveAccount() *account.Account
	Add(a *account.Account)
	All() []account.Account
	Get(id string) (*account.Account, bool)
	GetSubacc(id string, subacctype string) (*account.Account, bool)
	CloneAccountsOfType (typ string) []account.Account
	GetSKR03(id string) *account.Account
	GetByType (typ string) map[string]*account.Account
	ClearBookings()
	GetAllAccountsOfStakeholder (sh valueMagnets.Stakeholder) []account.Account
}

type DefaultAccountSystem struct {
	collectiveAccount *account.Account
	accounts          map[string]*account.Account
}

// Bilanzkonten
var SKR03_Anlagen25_35 = account.AccountDescription{Id: "SKR03_Anlagen25_35", Name: "01_Immaterielle Vermögensgegenstände", Type: account.KontenartAktiv}
var SKR03_Anlagen = account.AccountDescription{Id: "SKR03_Anlagen", Name: "021_Sachanlagen", Type: account.KontenartAktiv}
var SKR03_FinanzAnlagen = account.AccountDescription{Id: "SKR03_FinanzAnlagen", Name: "022_Finanzanlagen", Type: account.KontenartAktiv}
var SKR03_Vorraete = account.AccountDescription{Id: "SKR03_Vorraete", Name: "03_Vorräte", Type: account.KontenartAktiv}
var SKR03_1400 = account.AccountDescription{Id: "1400", Name: "04_Forderungen aus L&L + sonstige Vermögensg.", Type: account.KontenartAktiv}
var SKR03_Kautionen = account.AccountDescription{Id: "SKR03_Kautionen", Name: "05_Anzahlungen & Kautionen", Type: account.KontenartAktiv}
var SKR03_Vorsteuer = account.AccountDescription{Id: "SKR03_Vorsteuer", Name: "06_SKR03_1570-1579_Steuer:  Vorsteuer", Type: account.KontenartAktiv}
var SKR03_Umsatzsteuer = account.AccountDescription{Id: "SKR03_Umsatzsteuer", Name: "07_SKR03_1770_Steuer: Umsatzsteuer", Type: account.KontenartAktiv}
var SKR03_1200 = account.AccountDescription{Id: "1200", Name: "08_SKR03_1200_Bank", Type: account.KontenartAktiv}

var SKR03_sonstVerb = account.AccountDescription{Id: "SKR03_sonstVerb", Name: "11_sonstige Verbindlichkeiten", Type: account.KontenartPassiv}
var SKR03_900_Haftkapital = account.AccountDescription{Id: "SKR03_900_Haftkapital", Name: "12_SKR03_900_Haftkapital", Type: account.KontenartPassiv}
var SKR03_Eigenkapital_880 = account.AccountDescription{Id: "Eigenkapital", Name: "13_SKR03_880_Eigenkapital", Type: account.KontenartPassiv}
var SKR03_Rueckstellungen = account.AccountDescription{Id: "Rückstellung", Name: "14_SKR03_956-977_Rückstellung", Type: account.KontenartPassiv}
var SKR03_920_Gesellschafterdarlehen = account.AccountDescription{Id: "15_SKR03_920_Gesellschafterdarlehen", Name: "14_SKR03_920_Gesellschafterdarlehen", Type: account.KontenartPassiv}
var SKR03_1600 = account.AccountDescription{Id: "1600", Name: "16_SKR03_1600_OPOS-Lieferant", Type: account.KontenartPassiv}
var SKR03_1900 = account.AccountDescription{Id: "1900", Name: "17_SKR03_1900_Privatentnahmen", Type: account.KontenartPassiv}

// Erfolgskonten
var SKR03_Umsatzerloese = account.AccountDescription{Id: "SKR03_Umsatzerloese", Name: "1 SKR03_Umsatzerloese 8100-8402", Type: account.KontenartErtrag}
var SKR03_sonstigeErloese = account.AccountDescription{Id: "SKR03_sonstigeErloese", Name: "2 SKR03_Umsatzerloese 23**-88**", Type: account.KontenartErtrag}
var SKR03_4100_4199 = account.AccountDescription{Id: "4100_4199", Name: "3 Löhne und Gehälter 4100-4199", Type: account.KontenartAufwand}
//var SKR03_AnlagenabgaengeSachanlagen2310 = account.AccountDescription{Id: "SKR03_AnlagenabgängeSachanlagen-Buchverlust", Name: "4 AnlagenabgängeSachanlagen 2310 (Restb. Buchverlust)", Type: account.KontenartAufwand}
//var SKR03_AnlagenabgaengeSachanlagen2315 = account.AccountDescription{Id: "SKR03_AnlagenabgängeSachanlagen-Buchgewinn", Name: "4 AnlagenabgängeSachanlagen 2315 (Restb. Buchgewinn)", Type: account.KontenartErtrag}
var SKR03_Abschreibungen = account.AccountDescription{Id: "SKR03_Abschreibungen", Name: "4 Abschreibungen auf Anlagen 4822-4855", Type: account.KontenartAufwand}
var SKR03_sonstigeAufwendungen = account.AccountDescription{Id: "SKR03_sonstigeAufwendungen", Name: "5 sonstige Aufwendungen", Type: account.KontenartAufwand}
var SKR03_Steuern = account.AccountDescription{Id: "SKR03_Steuern", Name: "6 SKR03_Steuern 4320 Gewerbesteuer", Type: account.KontenartAufwand}
var ErgebnisNachSteuern = account.AccountDescription{Id: "SKR03_ErgebnisNachSteuern 10000", Name: "7 ErgebnisNachSteuern", Type: account.KontenartVerrechnung}

// Verrechnungskonten
var SKR03_Saldenvortrag = account.AccountDescription{Id: "SKR03_Saldenvortrag", Name: "Saldenvortrag 9000", Type: account.KontenartVerrechnung}
var SKR03_9790_Restanteil = account.AccountDescription{Id: "SKR03_9790_Restanteil", Name: "SKR03_9790_Restanteil", Type: account.KontenartVerrechnung}
var SummeAktiva 	= account.AccountDescription{Id: "SummeAktiva", Name: "V: SummeAktiva", Type: account.KontenartVerrechnung}
var SummePassiva 	= account.AccountDescription{Id: "SummePassiva", Name: "V: SummePassiva", Type: account.KontenartVerrechnung}
var AlleKLRBuchungen = account.AccountDescription{Id: "AlleKLRBuchungen", Name: "V: AlleKLRBuchungen", Type: account.KontenartKLR}
var k_ErloeseVerteilung = account.AccountDescription{Id: "k_ErloeseVerteilung", Name: "V: k_ErloeseVerteilung", Type: account.KontenartVerrechnung}


// Unterkonten für kommitmenschen
var UK_Kosten 				= account.AccountDescription{Id: "_0-Kosten", Name: "", Type: "Aktiv"}
var UK_AnteileausFairshare 	= account.AccountDescription{Id: "_1-AnteilausFairshare", Name: "", Type: "Aktiv"}
var UK_Vertriebsprovision 	= account.AccountDescription{Id: "_2-Vertriebsprovision", Name: "", Type: "Aktiv"}
var UK_AnteilMitmachen 		= account.AccountDescription{Id: "_3-Anteil-Mitmachen", Name: "", Type: "Aktiv"}
var UK_AnteileAuserloesen 	= account.AccountDescription{Id: "_4-Anteilauserloesen", Name: "", Type: "Aktiv"}
var UK_Erloese 				= account.AccountDescription{Id: "_5-Erloese", Name: "", Type: "Aktiv"}

var UK_Entnahmen 			= account.AccountDescription{Id: "_6-Entnahmen", Name: "", Type: "Passiv"}
var UK_VeraenderungAnlagen 	= account.AccountDescription{Id: "_7-VeränderungAnlagen", Name: "", Type: "Passiv"}
var UK_AnteilAnAnlagen 		= account.AccountDescription{Id: "_8-AnteilAnAnlagen", Name: "", Type: "Passiv"}
var UK_Darlehen 			= account.AccountDescription{Id: "_9-Darlehen", Name: "", Type: "Passiv"}
var UK_LiquidityReserve 	= account.AccountDescription{Id: "_A-Liquiditätsreserve", Name: "", Type: "Passiv"}
var UK_Verfuegungsrahmen	= account.AccountDescription{Id: "_B-Jahresüberschuss-Bonus", Name: "", Type: "Passiv"}
var Hauptkonto 				= account.AccountDescription{Id: "Hauptkonto", Name: "", Type: "Hauptkonto"}


var UK  = [...]account.AccountDescription {
	UK_Kosten,
	UK_AnteileausFairshare,
	UK_AnteilMitmachen,
	UK_Vertriebsprovision,
	UK_AnteileAuserloesen,
	UK_Erloese,
	UK_VeraenderungAnlagen,
	UK_AnteilAnAnlagen,
	UK_Darlehen,
	UK_Entnahmen,
	UK_LiquidityReserve,
	UK_Verfuegungsrahmen,
}


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
//		SKR03_AnlagenabgaengeSachanlagen2310,
//		SKR03_AnlagenabgaengeSachanlagen2315,
		SKR03_sonstigeAufwendungen,
		SKR03_Anlagen,
		SKR03_Anlagen25_35,
		SKR03_Abschreibungen,
		SKR03_FinanzAnlagen,
		SKR03_Vorraete,
		SKR03_Kautionen,
		SKR03_Umsatzerloese,
		SKR03_sonstigeErloese,
		SKR03_Steuern,
		SKR03_Vorsteuer,
		SKR03_Umsatzsteuer,
		SKR03_Saldenvortrag,
		SKR03_sonstVerb,
		ErgebnisNachSteuern,
		SummeAktiva,
		SummePassiva,
		SKR03_9790_Restanteil,
		AlleKLRBuchungen,
		k_ErloeseVerteilung,
	}
}

func EmptyDefaultAccountSystem() AccountSystem {
	o := account.AccountDescription{Id: "all", Name: "Hauptbuch", Type: account.KontenartVerrechnung}
	return &DefaultAccountSystem{collectiveAccount: &account.Account{Description: o}, accounts: make(map[string]*account.Account)}
}

func NewDefaultAccountSystem() AccountSystem {

	ad := account.AccountDescription{Id: "all", Name: "Hauptbuch", Type: account.KontenartVerrechnung}
	accountSystem := DefaultAccountSystem{collectiveAccount: &account.Account{Description: ad}, accounts: make(map[string]*account.Account)}

	// generate accounts according to the AccountList
	accountlist := Accountlist{}
	for _, a := range accountlist.All() {
		accountSystem.Add(account.NewAccount(a))
	}

	// generate accounts for all stakeholders
	stakeholder := valueMagnets.Stakeholder{}

	for _, sh := range stakeholder.All(util.Global.FinancialYear) {
		ad := account.AccountDescription{Id: sh.Id, Name: sh.Name, Type: Hauptkonto.Id}
		accountSystem.Add(account.NewAccount(ad))

		// create subaccounts
		for _, uk := range UK {
			//			sa := account.AccountDescription{Id: sh.Id+ukname, Name: sh.Name+ukname, Type: valueMagnets.StakeholderTypeKUA}
			sa := account.AccountDescription{Id: sh.Id+uk.Id, Name: sh.Name+uk.Id, Type: uk.Type}
			sa.Superaccount = ad.Id
			accountSystem.Add(account.NewAccount(sa))
		}
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
		clone := *a
		//clone.Bookings = nil
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

func (r *DefaultAccountSystem) CloneAccountsOfStakeholder(sh valueMagnets.Stakeholder) []account.Account {
	var accounts []account.Account

	acc,_ := r.Get(sh.Id)
	log.Println("in CloneAccountsOfStakeholder", acc)
	accounts = append(accounts, *acc)

	// find subaccounts
	for _, uk := range UK {
		sa := account.AccountDescription{Id: sh.Id+uk.Id, Name: sh.Name+uk.Id, Type: uk.Type}
		log.Println("in CloneAccountsOfStakeholder", sa)
	}

	return accounts
}


func (r *DefaultAccountSystem) GetSubacc(id string, subacctype string) (*account.Account, bool) {
	if a, ok := r.accounts[id+subacctype]; ok {
		return a, true
	}
	log.Println("in accountSystem.GetSubacc, could not find account ", id+subacctype)
	return nil, false
}


func (as *DefaultAccountSystem) GetByType(typ string) map[string]*account.Account {
	filtered  := make (map[string]*account.Account)

	for _, account := range as.accounts {
		if account.Description.Type == typ {
			clone := account
			filtered[account.Description.Name] = clone
		}
	}
	return filtered
}


func (as *DefaultAccountSystem) CloneAccountsOfType (typ string) []account.Account {
	var filtered  []account.Account
	for _, account := range as.accounts {
		if account.Description.Type == typ {
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


// get all accounts and subaccounts of a given Stakeholder or empty
func  (as *DefaultAccountSystem) GetAllAccountsOfStakeholder (sh valueMagnets.Stakeholder) []account.Account {
	var accountlist []account.Account
	var stakeholder valueMagnets.Stakeholder

	// check if stakeholder is in sztakeholder list
	for _,s := range stakeholder.All(util.Global.FinancialYear){
		if s == sh {
			// now add all the beloging accounts...
			for _, account := range as.accounts {
				if account.Description.Superaccount == sh.Id || account.Description.Id == sh.Id {
					clone := *account
					clone.Bookings = nil
					accountlist = append(accountlist, clone)
				}
			}
		}
	}

	sort.Slice(accountlist, func(i, j int) bool { return accountlist[i].Description.Name < accountlist[j].Description.Name })
	//	log.Println("in GetAllAccountsOfStakeholder",sh, accountlist)
	return accountlist
}




// find the right account for the SKR03konto string
func (r *DefaultAccountSystem) GetSKR03(SKR03konto string) *account.Account {
	var account *account.Account
	switch  {
	case IsInRange(SKR03konto, 25, 35): // Anlage buchen
		account = r.accounts[SKR03_Anlagen25_35.Id]
	case IsInRange(SKR03konto, 300, 490): // Anlage buchen
		account = r.accounts[SKR03_Anlagen.Id]
	case IsInRange(SKR03konto, 500, 599): // Anlage buchen
		account = r.accounts[SKR03_FinanzAnlagen.Id]
	case IsInRange(SKR03konto, 880, 899): // Eigenkapital bilden
		account = r.accounts[SKR03_Eigenkapital_880.Id]
	case IsInRange(SKR03konto, 900, 919):
		account = r.accounts[SKR03_900_Haftkapital.Id]
	case IsInRange(SKR03konto, 920, 929):
		account = r.accounts[SKR03_920_Gesellschafterdarlehen.Id]
	case IsInRange(SKR03konto, 930, 979): // Rückstellung bilden
		account = r.accounts[SKR03_Rueckstellungen.Id]
	case IsInRange(SKR03konto, 1200, 1250): // Bank buchen
		account = r.accounts[SKR03_1200.Id]
	case IsInRange(SKR03konto, 1518, 1518):
		account = r.accounts[SKR03_Vorraete.Id]
	case IsInRange(SKR03konto, 1525, 1525):
		account = r.accounts[SKR03_Kautionen.Id]
	case IsInRange(SKR03konto, 1548, 1587):
		account = r.accounts[SKR03_Vorsteuer.Id]
	case SKR03konto == "1400", SKR03konto == "1595":
		account = r.accounts[SKR03_1400.Id]
	case SKR03konto == "731", IsInRange(SKR03konto, 1600, 1699):
		account = r.accounts[SKR03_1600.Id]
	case IsInRange(SKR03konto, 1700, 1754),
		IsInRange(SKR03konto, 1758, 1799):
		account =  r.accounts[SKR03_sonstVerb.Id]  // bspw. 1755
	case IsInRange(SKR03konto, 1769, 1791):
		account = r.accounts[SKR03_Umsatzsteuer.Id]
	case SKR03konto == "1900":
		account = r.accounts[SKR03_1900.Id]
//	case IsInRange(SKR03konto, 2310, 2313):  --> ist nun in sonstigen Aufwendungen
//		account = r.accounts[SKR03_AnlagenabgaengeSachanlagen2310.Id]
	case SKR03konto == "4120":
		account = r.accounts[SKR03_4100_4199.Id]
	case IsInRange(SKR03konto, 1755, 1756),
		IsInRange(SKR03konto, 4130, 4140):
		account = r.accounts[SKR03_4100_4199.Id]  // Löhne & Gehälter
	case SKR03konto == "4320",
		IsInRange(SKR03konto, 2200, 2289):
		account = r.accounts[SKR03_Steuern.Id]
	case IsInRange(SKR03konto, 4822, 4855):
		account = r.accounts[SKR03_Abschreibungen.Id]
	case IsInRange(SKR03konto, 2000, 2199),
		 IsInRange(SKR03konto, 2300, 2309),
		 IsInRange(SKR03konto, 2310, 2313),
		 IsInRange(SKR03konto, 2320, 2350),
		 IsInRange(SKR03konto, 2380, 2409),
		 IsInRange(SKR03konto, 4200, 4306),
		 IsInRange(SKR03konto, 4360, 4500),
		 IsInRange(SKR03konto, 4520, 4810),
 		 IsInRange(SKR03konto, 4886, 4887),
		 IsInRange(SKR03konto, 4900, 4980):
		account = r.accounts[SKR03_sonstigeAufwendungen.Id]
	case IsInRange(SKR03konto, 8000, 8402),
		IsInRange(SKR03konto, 8700, 8799):
		account = r.accounts[SKR03_Umsatzerloese.Id]
	case IsInRange(SKR03konto, 2315, 2318),
		IsInRange(SKR03konto, 2700, 2749),
		IsInRange(SKR03konto, 2510, 2520),
		IsInRange(SKR03konto, 8600, 8699),
		IsInRange(SKR03konto, 8800, 8853):
		account = r.accounts[SKR03_sonstigeErloese.Id]
	case SKR03konto == "9000":
		account = r.accounts[SKR03_Saldenvortrag.Id]
	case SKR03konto == "9790":
		account = r.accounts[SKR03_9790_Restanteil.Id]
	case SKR03konto == "10000":
		account = r.accounts[ErgebnisNachSteuern.Id]
	default:
		log.Printf("GetSKR03: could not process booking type '%s'", SKR03konto)
		panic(fmt.Sprintf("GetSKR03: SKR03Bucket/Stakeholder/Konto '%s' not found", account.Description))
	}
	return account
}


// check if an SKR03 account number is in range
func IsInRange(num string, low, high int) bool {
	n, err := strconv.Atoi(num)
	if err != nil {
		fmt.Println("Error in IsInRange", num, low, high)
		panic(err)
	}
	return n >= low && n <= high
}
