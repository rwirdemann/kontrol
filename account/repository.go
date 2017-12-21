package account

import "bitbucket.org/rwirdemann/kontrol/domain"

type Repository interface {
	Add(a domain.Account)
}

type DefaultRepository struct {
	Accounts map[string]*domain.Account
}

var Repository Repository

func init() {

}

func Add(a domain.Account) {
}
