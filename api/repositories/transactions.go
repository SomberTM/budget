package repositories

import (
	"budget/api/dependencies"
	"budget/api/models"
	"context"
	"database/sql"

	"github.com/plaid/plaid-go/v27/plaid"
)

type TransactionsRepository interface {
	GetTransactionsForUser(ctx context.Context, userId string) ([]plaid.Transaction, error)
	GetTransactionCursorForUser(ctx context.Context, userId string) (models.TransactionCursor, error)
	UpsertTransactionCursorForUser(ctx context.Context, userId string, newCursor string) error
	AddTransactionsForUser(ctx context.Context, userId string, transactions []plaid.Transaction) error
	ModifyTransactionsForUser(ctx context.Context, userId string, transactions []plaid.Transaction) error
	DeleteTransactionsForUser(ctx context.Context, userId string, transactionIds []string) error
}

type DatabaseTransactionsRepository struct {
	db dependencies.Database
}

func NewDatabaseTransactionsRepository(db dependencies.Database) *DatabaseBudgetingRepository {
	return &DatabaseBudgetingRepository{db: db}
}

func (r *DatabaseTransactionsRepository) GetTransactionsForUser(ctx context.Context, userId string) ([]plaid.Transaction, error) {
	transactions := []plaid.Transaction{}
	err := r.db.GetConnection().SelectContext(ctx, &transactions, "SELECT * FROM transactions WHERE user_id = $1", userId)
	if err != nil {
		return []plaid.Transaction{}, err
	}

	return transactions, nil
}

func (r *DatabaseTransactionsRepository) GetTransactionCursorForUser(ctx context.Context, userId string) (models.TransactionCursor, error) {
	cursor := models.TransactionCursor{}
	err := r.db.GetConnection().GetContext(ctx, &cursor, "SELECT * FROM transaction_cursors WHERE user_id = $1 LIMIT 1", userId)
	if err != nil {
		return models.TransactionCursor{}, err
	}

	return cursor, nil
}

func (r *DatabaseTransactionsRepository) UpsertTransactionCursorForUser(ctx context.Context, userId string, newCursor string) error {
	_, err := r.GetTransactionCursorForUser(ctx, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			_, err = r.db.GetConnection().ExecContext(ctx, "INSERT INTO transaction_cursors (id, user_id, cursor) VALUES (default, $1, $2)", userId, newCursor)
			if err != nil {
				return err
			}

			return nil
		}

		return err
	}

	_, err = r.db.GetConnection().ExecContext(ctx, "UPDATE transaction_cursors SET cursor = $1 WHERE user_id = $2", newCursor, userId)
	if err != nil {
		return err
	}

	return nil
}

func (r *DatabaseTransactionsRepository) AddTransactionsForUser(ctx context.Context, userId string, transactions []plaid.Transaction) error {
	db := r.db.GetConnection()

	_, err := db.NamedExecContext(ctx, `INSERT INTO transactions 
		(id, user_id, account_id, amount, iso_currency_code, unofficial_currency_code, check_number, "date", "datetime", "location", "name", merchant_name, original_description, payment_meta, pending, pending_transaction_id, account_owner, transaction_id, logo_url, website, authorized_date, authorized_datetime, transaction_category_detailed, transaction_code, personal_finance_category_icon_url, counterparties, merchant_entity_id)
		VALUES (default, :user_id, :account_id, :amount, :iso_currency_code, :unofficial_currency_code, :check_number, :"date", :"datetime", :"location", :"name", :merchant_name, :original_description, :payment_meta, :pending, :pending_transaction_id, :account_owner, :transaction_id, :logo_url, :website, :authorized_date, :authorized_datetime, :transaction_category_detailed, :transaction_code, :personal_finance_category_icon_url, :counterparties, :merchant_entity_id)`,
		transactions)
	if err != nil {
		return err
	}

	return nil
}

func (r *DatabaseTransactionsRepository) ModifyTransactionsForUser(ctx context.Context, userId string, transactions []plaid.Transaction) error {
	return nil
}

func (r *DatabaseTransactionsRepository) DeleteTransactionsForUser(ctx context.Context, userId string, transactionIds []string) error {
	return nil
}

type NilTransactionsRepository struct{}

func (r *NilTransactionsRepository) GetTransactionsForUser(ctx context.Context, userId string) ([]plaid.Transaction, error) {
	return []plaid.Transaction{}, nil
}
func (r *NilTransactionsRepository) GetTransactionCursorForUser(ctx context.Context, userId string) (models.TransactionCursor, error) {
	return models.TransactionCursor{}, nil
}
func (r *NilTransactionsRepository) UpsertTransactionCursorForUser(ctx context.Context, userId string, newCursor string) error {
	return nil
}
func (r *NilTransactionsRepository) AddTransactionsForUser(ctx context.Context, userId string, transactions []plaid.Transaction) error {
	return nil
}
func (r *NilTransactionsRepository) ModifyTransactionsForUser(ctx context.Context, userId string, transactions []plaid.Transaction) error {
	return nil
}
func (r *NilTransactionsRepository) DeleteTransactionsForUser(ctx context.Context, userId string, transactionIds []string) error {
	return nil
}
