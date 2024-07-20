package services

type Services struct {
	Plaid PlaidService
	Users UsersService
}

func GetNilServices() *Services {
	return &Services{
		Plaid: &NilPlaidService{},
		Users: &NilUsersService{},
	}
}
