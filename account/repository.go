package account

import "bitbucket.org/rwirdemann/kontrol/domain"

type Repository interface {
	Add(a *domain.Account)
	All() []domain.Account
	Get(id string) (*domain.Account, bool)
	ClearBookings()
}

type DefaultRepository struct {
	accounts map[string]*domain.Account
}

func NewDefaultRepository() Repository {
	return &DefaultRepository{accounts: make(map[string]*domain.Account)}
}

func (r *DefaultRepository) Add(a *domain.Account) {
	r.accounts[a.Owner.Id] = a
}

func (r *DefaultRepository) All() []domain.Account {
	result := make([]domain.Account, 0, len(r.accounts))
	for _, a := range r.accounts {
		a.UpdateSaldo()
		result = append(result, *a)
	}
	return result
}

func (r *DefaultRepository) Get(id string) (*domain.Account, bool) {
	if a, ok := r.accounts[id]; ok {
		return a, true
	}
	return nil, false
}

func (r *DefaultRepository) ClearBookings() {
	for _, account := range r.accounts {
		account.Bookings = []domain.Booking{}
	}
}
