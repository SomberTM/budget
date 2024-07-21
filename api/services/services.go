package services

type Services struct {
	Plaid     PlaidService
	Users     UsersService
	Budgeting BudgetingService
}

func GetNilServices() *Services {
	return &Services{
		Plaid:     &NilPlaidService{},
		Users:     &NilUsersService{},
		Budgeting: &NilBudgetingService{},
	}
}
