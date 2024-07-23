package repositories_test

import (
	"budget/api/dependencies"
	"budget/api/environment"
	"budget/api/models"
	"budget/api/repositories"
	"context"
	"os"
	"testing"

	"github.com/plaid/plaid-go/v27/plaid"
)

var env *environment.Environment = nil

var (
	dataSourceName string
	migrationFile  string           = "../data_test.sql"
	testUserId     string           = "006f052d-2f4c-4110-a01f-7dc536667994"
	testItem       models.PlaidItem = models.PlaidItem{
		Id:          "2773ba5f-9c14-48f4-8b6d-08d87eda4149",
		UserId:      testUserId,
		ItemId:      "",
		AccessToken: "",
	}
)

func init() {
	dataSourceName = os.Getenv("CONNECTION_STRING")

	env = environment.GetNilEnvironment()

	db := dependencies.NewPostgreSql()
	db.Init(dataSourceName)
	db.Migrate(&migrationFile)

	c := db.GetConnection()
	c.MustExec("INSERT INTO users (id, user_name, password_hash) VALUES ($1, $2, $3)", testUserId, "sombertm", "123456")
	c.MustExec("INSERT INTO plaid_items (id, user_id, item_id, access_token) VALUES ($1, $2, $3, $4)", testItem.Id, testItem.UserId, testItem.ItemId, testItem.AccessToken)

	env.Database = db

	repos := repositories.GetNilRepositories()
	repos.Transactions = repositories.NewDatabaseTransactionsRepository(db)

	env.Repositories = repos
}

func nullableStringFromLiteral(str string) plaid.NullableString {
	return *plaid.NewNullableString(&str)
}

var transactions []models.Transaction

func TestEnsureTransactionConversionFromPlaid(t *testing.T) {
	transactions = models.NewTransactionsForItem(
		[]plaid.Transaction{
			{
				AccountId:              "BxBXxLj1m4HMXBm9WZZmCWVbPjX16EHwv99vp",
				AccountOwner:           plaid.NullableString{},
				Amount:                 72.1,
				IsoCurrencyCode:        nullableStringFromLiteral("USD"),
				UnofficialCurrencyCode: plaid.NullableString{},
				Category: []string{
					"Shops",
					"Supermarkets and Groceries",
				},
				CategoryId:  nullableStringFromLiteral("19046000"),
				CheckNumber: plaid.NullableString{},
				Counterparties: &[]plaid.TransactionCounterparty{
					{
						Name:            "Walmart",
						Type:            "merchant",
						LogoUrl:         nullableStringFromLiteral("https://plaid-merchant-logos.plaid.com/walmart_1100.png"),
						Website:         nullableStringFromLiteral("walmart.com"),
						EntityId:        nullableStringFromLiteral("O5W5j4dN9OR3E6ypQmjdkWZZRoXEzVMz2ByWM"),
						ConfidenceLevel: nullableStringFromLiteral("VERY_HIGH"),
					},
				},
				Date:               "2023-09-24",
				Datetime:           plaid.NullableTime{},
				AuthorizedDate:     nullableStringFromLiteral("2023-09-22"),
				AuthorizedDatetime: plaid.NullableTime{},
				Location: plaid.Location{
					Address:     nullableStringFromLiteral("13425 Community Rd"),
					City:        nullableStringFromLiteral("Poway"),
					Region:      nullableStringFromLiteral("CA"),
					PostalCode:  nullableStringFromLiteral("92064"),
					Country:     nullableStringFromLiteral("US"),
					StoreNumber: nullableStringFromLiteral("1700"),
				},
				Name:                           "PURCHASE WM SUPERCENTER #1700",
				MerchantName:                   nullableStringFromLiteral("Walmart"),
				MerchantEntityId:               nullableStringFromLiteral("O5W5j4dN9OR3E6ypQmjdkWZZRoXEzVMz2ByWM"),
				LogoUrl:                        nullableStringFromLiteral("https://plaid-merchant-logos.plaid.com/walmart_1100.png"),
				Website:                        nullableStringFromLiteral("walmart.com"),
				PaymentMeta:                    plaid.PaymentMeta{},
				PaymentChannel:                 "in store",
				Pending:                        false,
				PendingTransactionId:           nullableStringFromLiteral("no86Eox18VHMvaOVL7gPUM9ap3aR1LsAVZ5nc"),
				PersonalFinanceCategory:        plaid.NullablePersonalFinanceCategory{},
				PersonalFinanceCategoryIconUrl: nullableStringFromLiteral("https://plaid-category-icons.plaid.com/PFC_GENERAL_MERCHANDISE.png").Get(),
				TransactionId:                  "lPNjeW1nR6CDn5okmGQ6hEpMo4lLNoSrzqDje",
				TransactionCode:                plaid.NullableTransactionCode{},
				TransactionType:                nullableStringFromLiteral("place").Get(),
			},
			{
				AccountId:              "ajgkajglakjgklajsglkalk",
				AccountOwner:           plaid.NullableString{},
				Amount:                 72.1,
				IsoCurrencyCode:        nullableStringFromLiteral("USD"),
				UnofficialCurrencyCode: plaid.NullableString{},
				CategoryId:             nullableStringFromLiteral("19046000"),
				CheckNumber:            plaid.NullableString{},
				Counterparties: &[]plaid.TransactionCounterparty{
					{
						Name:            "Walmart",
						Type:            "merchant",
						LogoUrl:         nullableStringFromLiteral("https://plaid-merchant-logos.plaid.com/walmart_1100.png"),
						Website:         nullableStringFromLiteral("walmart.com"),
						EntityId:        nullableStringFromLiteral("O5W5j4dN9OR3E6ypQmjdkWZZRoXEzVMz2ByWM"),
						ConfidenceLevel: nullableStringFromLiteral("VERY_HIGH"),
					},
				},
				Date:               "2023-09-29",
				Datetime:           plaid.NullableTime{},
				AuthorizedDate:     nullableStringFromLiteral("2023-09-28"),
				AuthorizedDatetime: plaid.NullableTime{},
				Location: plaid.Location{
					Address:     nullableStringFromLiteral("13425 Community Rd"),
					City:        nullableStringFromLiteral("Poway"),
					Region:      nullableStringFromLiteral("CA"),
					PostalCode:  nullableStringFromLiteral("92064"),
					Country:     nullableStringFromLiteral("US"),
					StoreNumber: nullableStringFromLiteral("1700"),
				},
				Name:                           "PURCHASE WM SUPERCENTER #1700",
				MerchantName:                   nullableStringFromLiteral("Walmart"),
				MerchantEntityId:               nullableStringFromLiteral("O5W5j4dN9OR3E6ypQmjdkWZZRoXEzVMz2ByWM"),
				LogoUrl:                        nullableStringFromLiteral("https://plaid-merchant-logos.plaid.com/walmart_1100.png"),
				Website:                        nullableStringFromLiteral("walmart.com"),
				PaymentMeta:                    plaid.PaymentMeta{},
				PaymentChannel:                 "in store",
				Pending:                        false,
				PendingTransactionId:           nullableStringFromLiteral("no86Eox18VHMvaOVL7gPUM9ap3aR1LsAVZ5nc"),
				PersonalFinanceCategory:        plaid.NullablePersonalFinanceCategory{},
				PersonalFinanceCategoryIconUrl: nullableStringFromLiteral("https://plaid-category-icons.plaid.com/PFC_GENERAL_MERCHANDISE.png").Get(),
				TransactionId:                  "Ranakghjahhjktahkjghajk",
				TransactionCode:                plaid.NullableTransactionCode{},
				TransactionType:                nullableStringFromLiteral("place").Get(),
			},
		}, testItem)

	t.Logf("Converted %d plaid transaction(s) successfully", len(transactions))
}

func TestEnvironment(t *testing.T) {
	if env == nil || env.Database == nil || env.Repositories == nil || env.Repositories.Transactions == nil {
		t.Fatalf("Environment not initialized properly for the following tests")
	}
}

func countTransactions(t *testing.T) int {
	var c int
	row := env.GetConnection().QueryRow("SELECT COUNT(*) FROM transactions")
	if err := row.Scan(&c); err != nil {
		t.Fatalf("Error counting transactions: %v", err)
	}

	return c
}

func TestAddTransactions(t *testing.T) {
	err := env.Repositories.Transactions.AddTransactions(context.Background(), transactions)
	if err != nil {
		t.Fatalf("Expected transactions to be inserted: %v", err)
	}

	expected := len(transactions)
	actual := countTransactions(t)

	if expected != actual {
		t.Fatalf("Not all transactions inserted. Expected (%d), Actual (%d)", expected, actual)
	}

	t.Logf("All %d transaction(s) inserted", actual)
}

func TestModifyTransactions(t *testing.T) {
	newAmount := 1000000.0

	for i := 0; i < len(transactions); i++ {
		transactions[i].Amount = newAmount
	}

	err := env.Repositories.Transactions.ModifyTransactions(context.Background(), transactions)
	if err != nil {
		t.Fatalf("Error modifying transactions: %v", err)
	}

	var amounts []float64
	err = env.GetConnection().Select(&amounts, "SELECT amount FROM transactions")
	if err != nil {
		t.Fatalf("Failed to select modified amounts: %v", err)
	}

	for _, a := range amounts {
		if a != newAmount {
			t.Fatalf("Expected transaction amount to be updated. Expected (%f), Actual (%f)", newAmount, a)
		}
	}

	t.Logf("%d transaction(s) updated with amount %f", len(transactions), newAmount)
}

func TestRemoveTransactions(t *testing.T) {
	transactionIds := make([]string, 0, len(transactions))
	for _, t := range transactions {
		transactionIds = append(transactionIds, t.TransactionId)
	}

	err := env.Repositories.Transactions.DeleteTransactions(context.Background(), transactionIds)
	if err != nil {
		t.Fatalf("Error deleting transactions: %v", err)
	}

	expected := 0
	actual := countTransactions(t)

	if expected != actual {
		t.Fatalf("Not all transactions deleted. Expected (%d), Actual (%d)", expected, actual)
	}

	t.Logf("All transaction(s) deleted")
}
