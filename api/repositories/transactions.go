package repositories

import (
	"budget/api/dependencies"
	"budget/api/models"
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type TransactionsRepository interface {
	GetTransactionsForAccount(ctx context.Context, accountId string) ([]models.Transaction, error)
	GetTransactionsForItem(ctx context.Context, itemId string) ([]models.Transaction, error)
	GetTransactionsForUser(ctx context.Context, userId string) ([]models.Transaction, error)
	GetTransactionCursorForUser(ctx context.Context, userId string) (*models.TransactionCursor, error)
	UpsertTransactionCursorForUser(ctx context.Context, userId string, newCursor string) error
	AddTransactions(ctx context.Context, transactions []models.Transaction) error
	ModifyTransactions(ctx context.Context, transactions []models.Transaction) error
	DeleteTransactions(ctx context.Context, transactionIds []string) error
}

type DatabaseTransactionsRepository struct {
	db dependencies.Database
}

func NewDatabaseTransactionsRepository(db dependencies.Database) *DatabaseTransactionsRepository {
	return &DatabaseTransactionsRepository{db: db}
}

func (r *DatabaseTransactionsRepository) GetTransactionsForAccount(ctx context.Context, accountId string) ([]models.Transaction, error) {
	transactions := []models.Transaction{}
	err := r.db.GetConnection().SelectContext(ctx, &transactions, "SELECT * FROM transactions WHERE account_id = $1", accountId)
	if err != nil {
		return []models.Transaction{}, err
	}

	return transactions, nil
}

func (r *DatabaseTransactionsRepository) GetTransactionsForItem(ctx context.Context, itemId string) ([]models.Transaction, error) {
	transactions := []models.Transaction{}
	err := r.db.GetConnection().SelectContext(ctx, &transactions, "SELECT * FROM transactions WHERE item_id = $1", itemId)
	if err != nil {
		return []models.Transaction{}, err
	}

	return transactions, nil
}

func (r *DatabaseTransactionsRepository) GetTransactionsForUser(ctx context.Context, userId string) ([]models.Transaction, error) {
	transactions := []models.Transaction{}
	err := r.db.GetConnection().SelectContext(ctx, &transactions, "SELECT * FROM transactions WHERE user_id = $1", userId)
	if err != nil {
		return []models.Transaction{}, err
	}

	return transactions, nil
}

func (r *DatabaseTransactionsRepository) GetTransactionCursorForUser(ctx context.Context, userId string) (*models.TransactionCursor, error) {
	cursor := models.TransactionCursor{}
	err := r.db.GetConnection().GetContext(ctx, &cursor, "SELECT * FROM transaction_cursors WHERE user_id = $1 LIMIT 1", userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &cursor, nil
}

func (r *DatabaseTransactionsRepository) UpsertTransactionCursorForUser(ctx context.Context, userId string, newCursor string) error {
	db := r.db.GetConnection()

	_, err := db.ExecContext(ctx, "INSERT INTO transaction_cursors (id, user_id, cursor) VALUES (default, $1, $2) ON CONFLICT (user_id) DO UPDATE SET cursor = $2", userId, newCursor)
	if err != nil {
		return err
	}

	return nil
}

func (r *DatabaseTransactionsRepository) AddTransactions(ctx context.Context, transactions []models.Transaction) error {
	db := r.db.GetConnection()

	var err error
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	if len(transactions) > 0 {
		_, err = tx.NamedExecContext(ctx, `INSERT INTO transactions
		(id, item_id, user_id, transaction_category_detailed, account_id, transaction_id, amount, "date", data)
		VALUES (default, :item_id, :user_id, :transaction_category_detailed, :account_id, :transaction_id, :amount, :date, :data)`,
			transactions)
	}

	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (r *DatabaseTransactionsRepository) ModifyTransactions(ctx context.Context, transactions []models.Transaction) error {
	db := r.db.GetConnection()

	// there are fancier ways but we'll take an iterative approach for simplicity. wont be many modifications anyway
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	for _, t := range transactions {
		_, err = tx.NamedExec(`UPDATE transactions
			SET transaction_id = :transaction_id,
				account_id = :account_id,
				transaction_category_detailed = :transaction_category_detailed,
				amount = :amount,
				"date" = :date,
				data = :data
			WHERE transaction_id = :transaction_id
		`, t)

		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (r *DatabaseTransactionsRepository) DeleteTransactions(ctx context.Context, transactionIds []string) error {
	query, args, err := sqlx.In("DELETE FROM transactions WHERE transaction_id IN (?)", transactionIds)
	if err != nil {
		return err
	}

	db := r.db.GetConnection()
	query = db.Rebind(query)
	_, err = db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

type NilTransactionsRepository struct{}

func (r *NilTransactionsRepository) GetTransactionsForAccount(ctx context.Context, accountId string) ([]models.Transaction, error) {
	return []models.Transaction{}, nil
}
func (r *NilTransactionsRepository) GetTransactionsForItem(ctx context.Context, itemId string) ([]models.Transaction, error) {
	return []models.Transaction{}, nil
}
func (r *NilTransactionsRepository) GetTransactionsForUser(ctx context.Context, userId string) ([]models.Transaction, error) {
	return []models.Transaction{}, nil
}
func (r *NilTransactionsRepository) GetTransactionCursorForUser(ctx context.Context, userId string) (*models.TransactionCursor, error) {
	return &models.TransactionCursor{}, nil
}
func (r *NilTransactionsRepository) UpsertTransactionCursorForUser(ctx context.Context, userId string, newCursor string) error {
	return nil
}
func (r *NilTransactionsRepository) AddTransactions(ctx context.Context, transactions []models.Transaction) error {
	return nil
}
func (r *NilTransactionsRepository) ModifyTransactions(ctx context.Context, transactions []models.Transaction) error {
	return nil
}
func (r *NilTransactionsRepository) DeleteTransactions(ctx context.Context, transactionIds []string) error {
	return nil
}
