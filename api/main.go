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

	err = db.Migrate(nil)
	if err != nil {
		log.Fatal("Failed to migrate database")
	}
	env.Database = db

	plaid := dependencies.NewPlaidSandbox()
	// plaid := dependencies.NewPlaidProduction()
	plaid.SetClientId(plaidClientId)
	plaid.SetClientSecret(plaidSecret)
	plaid.Init()
	env.Plaid = plaid

	repos := repositories.GetNilRepositories()
	repos.Users = repositories.NewDatabaseUsersRepository(db)
	repos.PlaidItems = repositories.NewDatabasePlaidItemsRepository(db)
	repos.Sessions = repositories.NewDatabaseSessionsRepository(db)
	repos.TransactionCategories = repositories.NewDatabaseTransactionCategoriesRepository(db)
	repos.Budgeting = repositories.NewDatabaseBudgetingRepository(db)
	repos.Transactions = repositories.NewDatabaseTransactionsRepository(db)
	env.Repositories = repos

	servs := services.GetNilServices()
	servs.Plaid = services.NewPlaidFreeService(plaid.GetApiService(), repos.PlaidItems, repos.Transactions)
	servs.Users = services.NewSessionStrategyUsersService(repos.Users, repos.Sessions)
	servs.Budgeting = services.NewPrimaryBudgetingService(plaid.GetApiService(), repos.Budgeting, servs.Plaid)
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
		authorized.GET("/ubudgeting/budgets", environmentInjector(controllers.MyBudgets))
		authorized.POST("/ubudgeting/budgets", environmentInjector(controllers.CreateBudget))
		authorized.GET("/ubudgeting/budgets/:budgetId/breakdown", environmentInjector(controllers.BudgetBreakdown))
		authorized.POST("/ubudgeting/budgets/:budgetId/definitions", environmentInjector(controllers.CreateBudgetDefinition))
	}

	authorization := router.Group("/auth")
	{
		authorization.POST("/login", environmentInjector(controllers.Login))
		authorization.POST("/logout", environmentInjector(controllers.Logout))
	}

	budgeting := router.Group("/budgeting")
	{
		budgeting.GET("/categories", environmentInjector(controllers.TransactionCategories))

	}

	router.Run("localhost:8080")
}
