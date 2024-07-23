package services

import (
	"budget/api/models"
	"budget/api/repositories"
	"context"
	"errors"

	"github.com/plaid/plaid-go/v27/plaid"
)

type BudgetDefinitionBreakdown struct {
	models.BudgetDefinition
	Usage                  int64                        `json:"usage"`
	Categories             []models.TransactionCategory `json:"categories"`
	AssociatedTransactions []models.Transaction         `json:"associated_transactions"`
}

type BudgetBreakdown struct {
	models.Budget              `json:"budget"`
	BudgetDefinitionBreakdowns []BudgetDefinitionBreakdown `json:"budget_definitions"`
}

type BudgetingService interface {
	GetBudgetBreakdown(ctx context.Context, budgetId string, userId string) (BudgetBreakdown, error)
}

type PrimaryBudgetingService struct {
	plaidApi            *plaid.PlaidApiService
	budgetingRepository repositories.BudgetingRepository
	plaidService        PlaidService
}

func NewPrimaryBudgetingService(plaidApi *plaid.PlaidApiService, budgetingRepository repositories.BudgetingRepository, plaidService PlaidService) *PrimaryBudgetingService {
	return &PrimaryBudgetingService{plaidApi: plaidApi, budgetingRepository: budgetingRepository, plaidService: plaidService}
}

func (s *PrimaryBudgetingService) GetBudgetBreakdown(ctx context.Context, budgetId string, userId string) (BudgetBreakdown, error) {
	budget, exists, err := s.budgetingRepository.GetBudget(ctx, budgetId)
	if err != nil {
		return BudgetBreakdown{}, err
	}
	if !exists {
		return BudgetBreakdown{}, errors.New("no budget with id " + budgetId)
	}

	definitions, err := s.budgetingRepository.GetBudgetDefinitionsForBudget(ctx, budgetId)
	if err != nil {
		return BudgetBreakdown{}, err
	}

	err = s.plaidService.SyncUserTransactions(ctx, userId)
	if err != nil {
		return BudgetBreakdown{}, err
	}

	transactions, err := s.plaidService.GetUserTransactions(ctx, userId)
	if err != nil {
		return BudgetBreakdown{}, err
	}

	breakdown := BudgetBreakdown{
		Budget:                     budget,
		BudgetDefinitionBreakdowns: make([]BudgetDefinitionBreakdown, 0, len(definitions)),
	}

	for i := 0; i < len(definitions); i++ {
		definition := definitions[i]
		definitionBreakdown := BudgetDefinitionBreakdown{
			BudgetDefinition: definition,
		}

		categories, err := s.budgetingRepository.GetTransactionCategoriesForBudgetDefinition(ctx, definition.Id)
		if err != nil {
			return BudgetBreakdown{}, err
		}

		definitionBreakdown.Usage = 0
		definitionBreakdown.Categories = categories

		definitionBreakdown.AssociatedTransactions = make([]models.Transaction, 0)

		for j := 0; j < len(definitionBreakdown.Categories); j++ {
			category := definitionBreakdown.Categories[j]
			for k := 0; k < len(transactions); k++ {
				transaction := transactions[k]
				if transaction.TransactionCategoryDetailed.Valid {
					if transaction.TransactionCategoryDetailed.String == category.Detailed {
						definitionBreakdown.Usage += int64(transaction.Amount * 100)
						definitionBreakdown.AssociatedTransactions = append(definitionBreakdown.AssociatedTransactions, transaction)
					}
				}
			}
		}

		breakdown.BudgetDefinitionBreakdowns = append(breakdown.BudgetDefinitionBreakdowns, definitionBreakdown)
	}

	return breakdown, nil
}

type NilBudgetingService struct{}

func (s *NilBudgetingService) GetBudgetBreakdown(ctx context.Context, budgetId string, userId string) (BudgetBreakdown, error) {
	return BudgetBreakdown{}, nil
}
