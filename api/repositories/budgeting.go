package repositories

import (
	"budget/api/dependencies"
	"budget/api/models"
	"context"
	"database/sql"
)

type BudgetingRepository interface {
	GetBudgetsForUser(ctx context.Context, userId string) ([]models.Budget, error)
	GetBudget(ctx context.Context, budgetId string) (models.Budget, bool, error)
	GetBudgetDefinitionsForBudget(ctx context.Context, budgetId string) ([]models.BudgetDefinition, error)
	GetTransactionCategoriesForBudgetDefinition(ctx context.Context, budgetDefinitionId string) ([]models.TransactionCategory, error)
	CreateBudget(ctx context.Context, budget models.Budget) (models.Budget, error)
	CreateBudgetDefinition(ctx context.Context, budgetDefinition models.BudgetDefinition) (models.BudgetDefinition, error)
	AssignCategoriesToBudgetDefinition(ctx context.Context, budgetDefinitionId string, transactionCategoryIds []string) error
}

type DatabaseBudgetingRepository struct {
	db dependencies.Database
}

func NewDatabaseBudgetingRepository(db dependencies.Database) *DatabaseBudgetingRepository {
	return &DatabaseBudgetingRepository{db: db}
}

var nilBudget models.Budget = models.Budget{}
var nilBudgetDefinition models.BudgetDefinition = models.BudgetDefinition{}

func (r *DatabaseBudgetingRepository) GetBudgetsForUser(ctx context.Context, userId string) ([]models.Budget, error) {
	rows, err := r.db.GetConnection().QueryContext(ctx, "SELECT * FROM budgets WHERE user_id = $1", userId)
	if err != nil {
		return []models.Budget{}, nil
	}

	budgets := make([]models.Budget, 0)

	for rows.Next() {
		var budget models.Budget
		if err := rows.Scan(&budget.Id, &budget.UserId, &budget.Name, &budget.Color); err != nil {
			return []models.Budget{}, err
		}

		budgets = append(budgets, budget)
	}

	return budgets, nil
}

func (r *DatabaseBudgetingRepository) GetBudget(ctx context.Context, budgetId string) (models.Budget, bool, error) {
	var budget models.Budget
	row := r.db.GetConnection().QueryRowContext(ctx, "SELECT * FROM budgets WHERE id = $1", budgetId)
	if err := row.Scan(&budget.Id, &budget.UserId, &budget.Name, &budget.Color); err != nil {
		if err == sql.ErrNoRows {
			return nilBudget, false, nil
		}

		return nilBudget, false, err
	}

	return budget, true, nil
}

func (r *DatabaseBudgetingRepository) GetBudgetDefinitionsForBudget(ctx context.Context, budgetId string) ([]models.BudgetDefinition, error) {
	rows, err := r.db.GetConnection().QueryContext(ctx, "SELECT * FROM budget_definitions WHERE budget_id = $1", budgetId)
	if err != nil {
		return []models.BudgetDefinition{}, nil
	}

	definitions := make([]models.BudgetDefinition, 0)

	for rows.Next() {
		var definition models.BudgetDefinition
		if err := rows.Scan(&definition.Id, &definition.UserId, &definition.BudgetId, &definition.Name, &definition.Allocation); err != nil {
			return []models.BudgetDefinition{}, err
		}

		definitions = append(definitions, definition)
	}

	return definitions, nil
}

func (r *DatabaseBudgetingRepository) GetTransactionCategoriesForBudgetDefinition(ctx context.Context, budgetDefinitionId string) ([]models.TransactionCategory, error) {
	rows, err := r.db.GetConnection().QueryContext(ctx, "SELECT tc.* FROM transaction_categories as tc LEFT JOIN transaction_categories_to_budget_definitions as jt ON tc.id = jt.category_id WHERE jt.definition_id = $1", budgetDefinitionId)
	if err != nil {
		return []models.TransactionCategory{}, nil
	}

	categories := make([]models.TransactionCategory, 0)

	for rows.Next() {
		var category models.TransactionCategory
		if err := rows.Scan(&category.Id, &category.Primary, &category.Detailed, &category.Description); err != nil {
			return []models.TransactionCategory{}, err
		}

		categories = append(categories, category)
	}

	return categories, nil
}

func (r *DatabaseBudgetingRepository) CreateBudget(ctx context.Context, budget models.Budget) (models.Budget, error) {
	row := r.db.GetConnection().QueryRowContext(ctx, "INSERT INTO budgets (id, user_id, name, color) VALUES (default, $1, $2, $3) RETURNING id", budget.UserId, budget.Name, budget.Color)
	if err := row.Scan(&budget.Id); err != nil {
		return nilBudget, err
	}

	return budget, nil
}

func (r *DatabaseBudgetingRepository) CreateBudgetDefinition(ctx context.Context, budgetDefinition models.BudgetDefinition) (models.BudgetDefinition, error) {
	row := r.db.GetConnection().QueryRowContext(ctx, "INSERT INTO budget_definitions (id, user_id, budget_id, name, allocation) VALUES (default, $1, $2, $3, $4) RETURNING id", budgetDefinition.UserId, budgetDefinition.BudgetId, budgetDefinition.Name, budgetDefinition.Allocation)
	if err := row.Scan(&budgetDefinition.Id); err != nil {
		return nilBudgetDefinition, err
	}

	return budgetDefinition, nil
}

func (r *DatabaseBudgetingRepository) AssignCategoriesToBudgetDefinition(ctx context.Context, budgetDefinitionId string, transactionCategoryIds []string) error {
	_, err := r.db.GetConnection().QueryContext(ctx, "DELETE FROM transaction_categories_to_budget_definitions WHERE definition_id = $1", budgetDefinitionId)
	if err != nil {
		return err
	}

	// bulk insert is too annoying for now we'll just do it 1 at a time rofl
	for _, id := range transactionCategoryIds {
		_, err := r.db.GetConnection().QueryContext(ctx, "INSERT INTO transaction_categories_to_budget_definitions (id, definition_id, category_id) VALUES (default, $1, $2)", budgetDefinitionId, id)
		if err != nil {
			return err
		}
	}

	return nil
}

type NilBudgetingRepository struct{}

func (r *NilBudgetingRepository) GetBudgetsForUser(ctx context.Context, userId string) ([]models.Budget, error) {
	return []models.Budget{}, nil
}
func (r *NilBudgetingRepository) GetBudget(ctx context.Context, budgetId string) (models.Budget, bool, error) {
	return nilBudget, false, nil
}
func (r *NilBudgetingRepository) GetBudgetDefinitionsForBudget(ctx context.Context, budgetId string) ([]models.BudgetDefinition, error) {
	return []models.BudgetDefinition{}, nil
}
func (r *NilBudgetingRepository) GetTransactionCategoriesForBudgetDefinition(ctx context.Context, budgetDefinitionId string) ([]models.TransactionCategory, error) {
	return []models.TransactionCategory{}, nil
}
func (r *NilBudgetingRepository) CreateBudget(ctx context.Context, budget models.Budget) (models.Budget, error) {
	return nilBudget, nil
}
func (r *NilBudgetingRepository) CreateBudgetDefinition(ctx context.Context, budgetDefinition models.BudgetDefinition) (models.BudgetDefinition, error) {
	return nilBudgetDefinition, nil
}
func (r *NilBudgetingRepository) AssignCategoriesToBudgetDefinition(ctx context.Context, budgetDefinitionId string, transactionCategoryIds []string) error {
	return nil
}
