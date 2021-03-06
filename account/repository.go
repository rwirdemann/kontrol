package account

import (
	"bitbucket.org/rwirdemann/kontrol/booking"
	"bitbucket.org/rwirdemann/kontrol/owner"
)

type Repository interface {
	BankAccount() *Account
	Add(a *Account)
	All() []Account
	Get(id string) (*Account, bool)
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
		if sh.Type != owner.StakeholderTypeExtern {
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
