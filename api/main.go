package main

import (
	"budget/api/controllers"
	"budget/api/dependencies"
	"budget/api/environment"
	"budget/api/middleware"
	"budget/api/repositories"
	"budget/api/services"
	"flag"
	"log"

	"github.com/gin-gonic/gin"

	_ "github.com/lib/pq"
)

var (
	router *gin.Engine
	env    *environment.Environment
)

var (
	driverName     string
	dataSourceName string

	plaidClientId string
	plaidSecret   string

	sqlScriptPath string
)

func init() {
	flag.StringVar(&driverName, "driver-name", "postgres", "name of the driver to use")
	flag.StringVar(&dataSourceName, "data-source-name", "", "data source name or connection string")
	flag.StringVar(&plaidClientId, "plaid-client-id", "", "plaid client id")
	flag.StringVar(&plaidSecret, "plaid-secret", "", "plaid secret")
	flag.StringVar(&sqlScriptPath, "sql-script", "", "sql script to run before starting api")
	flag.Parse()

	env = environment.GetNilEnvironment()

	db := dependencies.NewPostgreSql()
	err := db.Init(dataSourceName)
	if err != nil {
		log.Fatal("Failed to initialize database")
	}

	err = db.Migrate()
	if err != nil {
		log.Fatal("Failed to migrate database")
	}
	env.Database = db

	plaid := dependencies.NewPlaidSandbox()
	plaid.SetClientId(plaidClientId)
	plaid.SetClientSecret(plaidSecret)
	plaid.Init()
	env.Plaid = plaid

	repos := repositories.GetNilRepositories()
	repos.Users = repositories.NewDatabaseUsersRepository(db)
	repos.PlaidItems = repositories.NewDatabasePlaidItemsRepository(db)
	repos.Sessions = repositories.NewDatabaseSessionsRepository(db)
	env.Repositories = repos

	servs := services.GetNilServices()
	servs.Plaid = services.NewPlaidFreeService(plaid.GetApiService(), repos.PlaidItems)
	servs.Users = services.NewSessionStrategyUsersService(repos.Users, repos.Sessions)
	env.Services = servs

	router = gin.Default()
}

func environmentInjector(handler func(e *environment.Environment, c *gin.Context)) func(c *gin.Context) {
	return func(c *gin.Context) { handler(env, c) }
}

func main() {
	authorized := router.Group("/")
	authorized.Use(environmentInjector(middleware.RequireLoggedInUser))
	{
		authorized.GET("/me", controllers.Me)
		authorized.GET("/link/create", environmentInjector(controllers.CreateLinkToken))
		authorized.POST("/link/exchange", environmentInjector(controllers.ExchangePublicToken))

		authorized.GET("/user-accounts", environmentInjector(controllers.GetAccounts))
		// authorized.GET("/transactions", controllers.GetTransactions)
	}

	authorization := router.Group("/auth")
	{
		authorization.POST("/login", environmentInjector(controllers.Login))
		authorization.POST("/logout", environmentInjector(controllers.Logout))
	}

	router.Run("localhost:8080")
}
