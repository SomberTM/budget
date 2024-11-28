package repositories

type Repositories struct {
	Users        *UsersRepository
	Sessions     SessionsRepository
	PlaidItems   *PlaidItemsRepository
	Transactions TransactionsRepository
}

func GetNilRepositories() *Repositories {
	return &Repositories{
		Users:        nil,
		Sessions:     &NilSessionsRepository{},
		PlaidItems:   nil,
		Transactions: &NilTransactionsRepository{},
	}
}
