package account

import (
	"fmt"
	"log"

	"github.com/ahojsenn/kontrol/booking"
	"github.com/ahojsenn/kontrol/owner"
)

type Repository interface {
	BankAccount() *Account
	Add(a *Account)
	All() []Account
	Get(id string) (*Account, bool)
	GetSKR03(id string) *Account
	ClearBookings()
}

type DefaultRepository struct {
	collectiveAccount *Account
	accounts          map[string]*Account
}

func EmptyDefaultRepository() Repository {
	o := owner.Stakeholder{Id: "GLS", Name: "Kommitment GmbH & Co. KG", Type: owner.StakeholderTypeBank}
	return &DefaultRepository{collectiveAccount: &Account{Owner: o}, accounts: make(map[string]*Account)}
}

func NewDefaultRepository() Repository {
	o := owner.Stakeholder{Id: "GLS", Name: "Kommitment GmbH & Co. KG", Type: owner.StakeholderTypeBank}
	r := DefaultRepository{collectiveAccount: &Account{Owner: o}, accounts: make(map[string]*Account)}
	stakeholderRepository := owner.StakeholderRepository{}
	for _, sh := range stakeholderRepository.All() {
		if sh.Type != owner.StakeholderTypeExtern &&
			sh.Type != owner.StakeholderTypeOthers {
			r.Add(NewAccount(sh))
		}
	}
	return &r
}

func (r *DefaultRepository) BankAccount() *Account {
	return r.collectiveAccount
}

func (r *DefaultRepository) Add(a *Account) {
	r.accounts[a.Owner.Id] = a
}

func (r *DefaultRepository) All() []Account {
	result := make([]Account, 0, len(r.accounts))
	for _, a := range r.accounts {
		a.UpdateSaldo()
		clone := *a
		clone.Bookings = nil
		result = append(result, clone)
	}
	return result
}

func (r *DefaultRepository) Get(id string) (*Account, bool) {
	if a, ok := r.accounts[id]; ok {
		return a, true
	}
	return nil, false
}

func (r *DefaultRepository) ClearBookings() {
	r.collectiveAccount.Bookings = []booking.Booking{}
	for _, account := range r.accounts {
		account.Bookings = []booking.Booking{}
	}
}

func (r *DefaultRepository) GetSKR03(SKR03konto string) *Account {
	var account *Account
	switch SKR03konto {
	case "410": // Anlage buchen
		account = r.accounts[owner.SKR03_Anlagen.Id]
	case "965": // Rückstellung bilden
		account = r.accounts[owner.StakeholderRueckstellung.Id]
	case "4120":
		account = r.accounts[owner.SKR03_4100_4199.Id]
	case "4130", "4138", "4140":
		account = r.accounts[owner.SKR03_4100_4199.Id]
	case "1200":
		account = r.BankAccount()
	default:
		log.Printf("GetSKR03: could not process booking type '%s'", SKR03konto)
		panic(fmt.Sprintf("GetSKR03: SKR03Bucket/Stakeholder/Konto '%s' not found", account))
	}
	return account
}
