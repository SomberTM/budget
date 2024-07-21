package services

import (
	"budget/api/models"
	"budget/api/repositories"
	"context"
	"log"

	"github.com/plaid/plaid-go/v27/plaid"
)

type PlaidService interface {
	GetPlaidItemAccounts(ctx context.Context, item models.PlaidItem) ([]plaid.AccountBase, error)
	GetUserAccounts(ctx context.Context, user models.User) ([]plaid.AccountBase, error)
	GetPlaidItemTransactions(ctx context.Context, item models.PlaidItem) ([]plaid.Transaction, error)
	GetUserTransactions(ctx context.Context, userId string) ([]plaid.Transaction, error)
	GetLinkTokenForUser(ctx context.Context, user models.User) (string, error)
	ExchangePublicToken(ctx context.Context, user models.User, publicToken string) (models.PlaidItem, error)
}

type PlaidFreeService struct {
	api                  *plaid.PlaidApiService
	plaidItemsRepository repositories.PlaidItemsRepository
}

func NewPlaidFreeService(api *plaid.PlaidApiService, plaidItemsRepository repositories.PlaidItemsRepository) *PlaidFreeService {
	return &PlaidFreeService{api: api, plaidItemsRepository: plaidItemsRepository}
}

var nilAccounts []plaid.AccountBase = []plaid.AccountBase{}
var nilTransactions []plaid.Transaction = []plaid.Transaction{}

func (s *PlaidFreeService) GetPlaidItemAccounts(ctx context.Context, i models.PlaidItem) ([]plaid.AccountBase, error) {
	accountsGetRequest := plaid.NewAccountsGetRequest(i.AccessToken)

	response, _, err := s.api.AccountsGet(ctx).AccountsGetRequest(*accountsGetRequest).Execute()
	if err != nil {
		return []plaid.AccountBase{}, err
	}

	return response.GetAccounts(), nil
}

func (s *PlaidFreeService) GetUserAccounts(ctx context.Context, user models.User) ([]plaid.AccountBase, error) {
	items, err := s.plaidItemsRepository.GetPlaidItemsByUserId(ctx, user.Id)
	if err != nil {
		return nilAccounts, err
	}

	accounts := make([]plaid.AccountBase, 0, 20)
	for i := 0; i < len(items); i++ {
		itemAccounts, err := s.GetPlaidItemAccounts(ctx, items[i])
		if err != nil {
			return []plaid.AccountBase{}, err
		}

		accounts = append(accounts, itemAccounts...)
	}

	return accounts, nil
}

func (s *PlaidFreeService) GetPlaidItemTransactions(ctx context.Context, i models.PlaidItem) ([]plaid.Transaction, error) {
	var transactions []plaid.Transaction
	var cursor string

	hasMore := true
	for hasMore {
		request := plaid.NewTransactionsSyncRequest(i.AccessToken)
		if cursor != "" {
			request.SetCursor(cursor)
		}

		resp, _, err := s.api.TransactionsSync(
			ctx,
		).TransactionsSyncRequest(*request).Execute()
		if err != nil {
			return []plaid.Transaction{}, err
		}

		// Add this page of results
		transactions = append(transactions, resp.GetAdded()...)
		// modified = append(modified, resp.GetModified()...)
		// removed = append(removed, resp.GetRemoved()...)

		hasMore = resp.GetHasMore()
		cursor = resp.GetNextCursor()
	}

	return transactions, nil
}

func (s *PlaidFreeService) GetUserTransactions(ctx context.Context, userId string) ([]plaid.Transaction, error) {
	items, err := s.plaidItemsRepository.GetPlaidItemsByUserId(ctx, userId)
	if err != nil {
		return nilTransactions, err
	}

	transactions := make([]plaid.Transaction, 0, 20)
	for i := 0; i < len(items); i++ {
		itemTransactions, err := s.GetPlaidItemTransactions(ctx, items[i])
		if err != nil {
			return []plaid.Transaction{}, err
		}

		transactions = append(transactions, itemTransactions...)
	}

	return transactions, nil
}

func (s *PlaidFreeService) GetLinkTokenForUser(ctx context.Context, user models.User) (string, error) {
	request := plaid.NewLinkTokenCreateRequest("Plaid Test App", "en", []plaid.CountryCode{plaid.COUNTRYCODE_US}, *plaid.NewLinkTokenCreateRequestUser(user.Id))
	request.SetProducts([]plaid.Products{plaid.PRODUCTS_AUTH, plaid.PRODUCTS_TRANSACTIONS})

	response, _, err := s.api.LinkTokenCreate(ctx).LinkTokenCreateRequest(*request).Execute()
	if err != nil {
		return "", err
	}

	return response.GetLinkToken(), nil
}

func (s *PlaidFreeService) ExchangePublicToken(ctx context.Context, user models.User, publicToken string) (models.PlaidItem, error) {
	exchangeRequest := plaid.NewItemPublicTokenExchangeRequest(publicToken)
	response, _, err := s.api.ItemPublicTokenExchange(ctx).ItemPublicTokenExchangeRequest(*exchangeRequest).Execute()
	if err != nil {
		log.Fatalf("Error exchanging public token %v", err)
	}

	itemId := response.GetItemId()
	accessToken := response.GetAccessToken()

	item := models.NewPlaidItem()
	item.SetUserId(user.Id).SetItemId(itemId).SetAccessToken(accessToken)

	_, err = s.plaidItemsRepository.CreatePlaidItem(ctx, item)
	if err != nil {
		return models.PlaidItem{}, err
	}

	return item, nil
}

type NilPlaidService struct{}

func (s *NilPlaidService) GetPlaidItemAccounts(ctx context.Context, i models.PlaidItem) ([]plaid.AccountBase, error) {
	return nilAccounts, nil
}
func (s *NilPlaidService) GetUserAccounts(ctx context.Context, user models.User) ([]plaid.AccountBase, error) {
	return nilAccounts, nil
}
func (s *NilPlaidService) GetPlaidItemTransactions(ctx context.Context, i models.PlaidItem) ([]plaid.Transaction, error) {
	return nilTransactions, nil
}
func (s *NilPlaidService) GetUserTransactions(ctx context.Context, userId string) ([]plaid.Transaction, error) {
	return nilTransactions, nil
}
func (s *NilPlaidService) GetLinkTokenForUser(ctx context.Context, user models.User) (string, error) {
	return "", nil
}
func (s *NilPlaidService) ExchangePublicToken(ctx context.Context, user models.User, publicToken string) (models.PlaidItem, error) {
	return models.PlaidItem{}, nil
}
