package repositories

type Repositories struct {
	Users                 UsersRepository
	Sessions              SessionsRepository
	PlaidItems            PlaidItemsRepository
	TransactionCategories TransactionCategoriesRepository
	Budgeting             BudgetingRepository
	Transactions          TransactionsRepository
}

func GetNilRepositories() *Repositories {
	return &Repositories{
		Users:                 &NilUsersRepository{},
		Sessions:              &NilSessionsRepository{},
		PlaidItems:            &NilPlaidItemsRepository{},
		TransactionCategories: &NilTransactionCategoriesRepository{},
		Budgeting:             &NilBudgetingRepository{},
		Transactions:          &NilTransactionsRepository{},
	}
}
