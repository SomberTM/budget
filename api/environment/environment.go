package environment

import (
	"budget/api/dependencies"
	"budget/api/repositories"
	"budget/api/services"
)

// https://www.jerf.org/iri/post/2929/
type Environment struct {
	dependencies.Database
	dependencies.Plaid
	*repositories.Repositories
	*services.Services
}

func NewEnvironment(database dependencies.Database, plaid dependencies.Plaid, repos *repositories.Repositories, servs *services.Services) *Environment {
	if database == nil {
		database = &dependencies.NilDatabase{}
	}
	if plaid == nil {
		plaid = &dependencies.NilPlaid{}
	}
	if repos == nil {
		repos = repositories.GetNilRepositories()
	}
	if servs == nil {
		servs = services.GetNilServices()
	}

	return &Environment{Database: database, Plaid: plaid, Repositories: repos, Services: servs}
}

func GetNilEnvironment() *Environment {
	return NewEnvironment(nil, nil, nil, nil)
}
