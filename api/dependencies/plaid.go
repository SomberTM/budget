package dependencies

import "github.com/plaid/plaid-go/v27/plaid"

type Plaid interface {
	SetClientId(clientId string)
	SetClientSecret(clientSecret string)
	Init()
	GetApiService() *plaid.PlaidApiService
}

type PlaidSandbox struct {
	client       *plaid.APIClient
	clientId     string
	clientSecret string
}

func NewPlaidSandbox() *PlaidSandbox {
	return &PlaidSandbox{}
}

func (s *PlaidSandbox) SetClientId(clientId string) {
	s.clientId = clientId
}

func (s *PlaidSandbox) SetClientSecret(clientSecret string) {
	s.clientSecret = clientSecret
}

func (s *PlaidSandbox) Init() {
	configuration := plaid.NewConfiguration()
	configuration.AddDefaultHeader("PLAID-CLIENT-ID", s.clientId)
	configuration.AddDefaultHeader("PLAID-SECRET", s.clientSecret)
	configuration.UseEnvironment(plaid.Sandbox)
	s.client = plaid.NewAPIClient(configuration)
}

func (s *PlaidSandbox) GetApiService() *plaid.PlaidApiService {
	return s.client.PlaidApi
}

type PlaidProduction struct {
	client       *plaid.APIClient
	clientId     string
	clientSecret string
}

func NewPlaidProduction() *PlaidProduction {
	return &PlaidProduction{}
}

func (s *PlaidProduction) SetClientId(clientId string) {
	s.clientId = clientId
}

func (s *PlaidProduction) SetClientSecret(clientSecret string) {
	s.clientSecret = clientSecret
}

func (s *PlaidProduction) Init() {
	configuration := plaid.NewConfiguration()
	configuration.AddDefaultHeader("PLAID-CLIENT-ID", s.clientId)
	configuration.AddDefaultHeader("PLAID-SECRET", s.clientSecret)
	configuration.UseEnvironment(plaid.Production)
	s.client = plaid.NewAPIClient(configuration)
}

func (s *PlaidProduction) GetApiService() *plaid.PlaidApiService {
	return s.client.PlaidApi
}

type NilPlaid struct{}

func (s *NilPlaid) SetClientId(clientId string)           {}
func (s *NilPlaid) SetClientSecret(clientSecret string)   {}
func (s *NilPlaid) Init()                                 {}
func (s *NilPlaid) GetApiService() *plaid.PlaidApiService { return nil }
