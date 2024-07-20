package repositories

import (
	"budget/api/dependencies"
	"budget/api/models"
	"context"
)

type PlaidItemsRepository interface {
	GetPlaidItemsByUserId(ctx context.Context, userId string) ([]models.PlaidItem, error)
	CreatePlaidItem(ctx context.Context, i models.PlaidItem) (models.PlaidItem, error)
}

type DatabasePlaidItemsRepository struct {
	db dependencies.Database
}

func NewDatabasePlaidItemsRepository(db dependencies.Database) *DatabasePlaidItemsRepository {
	return &DatabasePlaidItemsRepository{db: db}
}

var nilPlaidItems []models.PlaidItem = []models.PlaidItem{}

func (r *DatabasePlaidItemsRepository) GetPlaidItemsByUserId(ctx context.Context, userId string) ([]models.PlaidItem, error) {
	rows, err := r.db.GetConnection().QueryContext(ctx, "SELECT * FROM plaid_items WHERE user_id = $1", userId)
	if err != nil {
		return nilPlaidItems, err
	}

	items := make([]models.PlaidItem, 0, 10)
	for rows.Next() {
		var item models.PlaidItem
		if err := rows.Scan(&item.Id, &item.UserId, &item.ItemId, &item.AccessToken); err != nil {
			return nilPlaidItems, err
		}
		items = append(items, item)
	}

	return items, nil
}

func (r *DatabasePlaidItemsRepository) CreatePlaidItem(ctx context.Context, i models.PlaidItem) (models.PlaidItem, error) {
	row := r.db.GetConnection().QueryRowContext(ctx, "INSERT INTO plaid_items (id, user_id, item_id, access_token) VALUES (default, $1, $2, $3) RETURNING id", i.UserId, i.ItemId, i.AccessToken)
	if err := row.Scan(&i.Id); err != nil {
		return i, err
	}

	return i, nil
}

type NilPlaidItemsRepository struct{}

func (r *NilPlaidItemsRepository) GetPlaidItemsByUserId(ctx context.Context, userId string) ([]models.PlaidItem, error) {
	return nilPlaidItems, nil
}
func (r *NilPlaidItemsRepository) CreatePlaidItem(ctx context.Context, i models.PlaidItem) (models.PlaidItem, error) {
	return i, nil
}
