package repositories

import (
	"budget/api/dependencies"
	"budget/api/models"
	"context"
)

type TransactionCategoriesRepository interface {
	GetTransactionCategories(ctx context.Context) ([]models.TransactionCategory, error)
}

type DatabaseTransactionCategoriesRepository struct {
	db dependencies.Database
}

func NewDatabaseTransactionCategoriesRepository(db dependencies.Database) *DatabaseTransactionCategoriesRepository {
	return &DatabaseTransactionCategoriesRepository{db: db}
}

func (r *DatabaseTransactionCategoriesRepository) GetTransactionCategories(ctx context.Context) ([]models.TransactionCategory, error) {
	rows, err := r.db.GetConnection().QueryContext(ctx, "SELECT * FROM transaction_categories")
	if err != nil {
		return []models.TransactionCategory{}, err
	}

	categories := make([]models.TransactionCategory, 0, 104)
	for rows.Next() {
		var category models.TransactionCategory
		if err := rows.Scan(&category.Id, &category.Primary, &category.Detailed, &category.Description); err != nil {
			return []models.TransactionCategory{}, err
		}

		categories = append(categories, category)
	}

	return categories, nil
}

type NilTransactionCategoriesRepository struct{}

func (r *NilTransactionCategoriesRepository) GetTransactionCategories(ctx context.Context) ([]models.TransactionCategory, error) {
	return []models.TransactionCategory{}, nil
}
