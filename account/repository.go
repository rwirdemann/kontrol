package account

import "bitbucket.org/rwirdemann/kontrol/domain"

type Repository interface {
	Add(a domain.Account)
}

type DefaultRepository struct {
	Accounts map[string]*domain.Account
}

func NewDefaultRepository() Repository {
	return &DefaultRepository{}
}

func (r *DefaultRepository) Add(a domain.Account) {
}
