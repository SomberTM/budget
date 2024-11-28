package repositories

import (
	"budget/api/dependencies"
	"budget/api/models"
	"context"
)

type PlaidItemsRepository struct {
	db dependencies.Database
}

func NewDatabasePlaidItemsRepository(db dependencies.Database) *PlaidItemsRepository {
	return &PlaidItemsRepository{db: db}
}

var nilPlaidItems []models.PlaidItem = []models.PlaidItem{}

func (r *PlaidItemsRepository) GetPlaidItemsByUserId(ctx context.Context, userId string) ([]models.PlaidItem, error) {
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

func (r *PlaidItemsRepository) CreatePlaidItem(ctx context.Context, i models.PlaidItem) (models.PlaidItem, error) {
	row := r.db.GetConnection().QueryRowContext(ctx, "INSERT INTO plaid_items (id, user_id, item_id, institution_id, institution_name, access_token, raw_plaid_data) VALUES (default, $1, $2, $3, $4, $5, $6) RETURNING id", i.UserId, i.ItemId, i.InstitutionId, i.InstitutionName, i.AccessToken, i.RawPlaidData)
	if err := row.Scan(&i.Id); err != nil {
		return i, err
	}

	return i, nil
}
