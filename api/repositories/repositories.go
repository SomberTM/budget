package repositories

type Repositories struct {
	Users      UsersRepository
	Sessions   SessionsRepository
	PlaidItems PlaidItemsRepository
}

func GetNilRepositories() *Repositories {
	return &Repositories{
		Users:      &NilUsersRepository{},
		Sessions:   &NilSessionsRepository{},
		PlaidItems: &NilPlaidItemsRepository{},
	}
}
