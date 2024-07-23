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
	GetPlaidItemTransactions(ctx context.Context, item models.PlaidItem) ([]models.Transaction, error)
	GetUserTransactions(ctx context.Context, userId string) ([]models.Transaction, error)
	SyncPlaidItemTransactions(ctx context.Context, item models.PlaidItem) error
	SyncUserTransactions(ctx context.Context, userId string) error
	GetLinkTokenForUser(ctx context.Context, user models.User) (string, error)
	ExchangePublicToken(ctx context.Context, user models.User, publicToken string) (models.PlaidItem, error)
}

type PlaidFreeService struct {
	api                    *plaid.PlaidApiService
	plaidItemsRepository   repositories.PlaidItemsRepository
	transactionsRepository repositories.TransactionsRepository
}

func NewPlaidFreeService(api *plaid.PlaidApiService, plaidItemsRepository repositories.PlaidItemsRepository, transactionsRepository repositories.TransactionsRepository) *PlaidFreeService {
	return &PlaidFreeService{api: api, plaidItemsRepository: plaidItemsRepository, transactionsRepository: transactionsRepository}
}

var nilAccounts []plaid.AccountBase = []plaid.AccountBase{}
var nilTransactions []models.Transaction = []models.Transaction{}

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

func (s *PlaidFreeService) GetPlaidItemTransactions(ctx context.Context, i models.PlaidItem) ([]models.Transaction, error) {
	return s.transactionsRepository.GetTransactionsForItem(ctx, i.Id)
}

func (s *PlaidFreeService) GetUserTransactions(ctx context.Context, userId string) ([]models.Transaction, error) {
	return s.transactionsRepository.GetTransactionsForUser(ctx, userId)
}

func (s *PlaidFreeService) SyncPlaidItemTransactions(ctx context.Context, item models.PlaidItem) error {
	var added []plaid.Transaction
	var modified []plaid.Transaction
	var removed []plaid.RemovedTransaction

	savedCursor, err := s.transactionsRepository.GetTransactionCursorForUser(ctx, item.UserId)
	if err != nil {
		return err
	}

	var cursor string
	if savedCursor != nil {
		cursor = savedCursor.Cursor
	}

	hasMore := true
	for hasMore {
		request := plaid.NewTransactionsSyncRequest(item.AccessToken)
		if cursor != "" {
			request.SetCursor(cursor)
		}

		resp, _, err := s.api.TransactionsSync(
			ctx,
		).TransactionsSyncRequest(*request).Execute()
		if err != nil {
			return err
		}

		// Add this page of results
		added = append(added, resp.GetAdded()...)
		modified = append(modified, resp.GetModified()...)
		removed = append(removed, resp.GetRemoved()...)

		hasMore = resp.GetHasMore()
		cursor = resp.GetNextCursor()
	}

	// log.Printf("%v %v %v", added, modified, removed)

	s.transactionsRepository.UpsertTransactionCursorForUser(ctx, item.UserId, cursor)
	s.transactionsRepository.AddTransactions(ctx, models.NewTransactionsForItem(added, item))
	s.transactionsRepository.ModifyTransactions(ctx, models.NewTransactionsForItem(modified, item))

	transactionIds := make([]string, len(removed))
	for i := 0; i < len(removed); i++ {
		transactionIds = append(transactionIds, removed[i].GetTransactionId())
	}
	s.transactionsRepository.DeleteTransactions(ctx, transactionIds)

	return nil
}

func (s *PlaidFreeService) SyncUserTransactions(ctx context.Context, userId string) error {
	items, err := s.plaidItemsRepository.GetPlaidItemsByUserId(ctx, userId)
	if err != nil {
		return err
	}

	for i := 0; i < len(items); i++ {
		err := s.SyncPlaidItemTransactions(ctx, items[i])
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *PlaidFreeService) GetLinkTokenForUser(ctx context.Context, user models.User) (string, error) {
	request := plaid.NewLinkTokenCreateRequest("Plaid Test App", "en", []plaid.CountryCode{plaid.COUNTRYCODE_US}, *plaid.NewLinkTokenCreateRequestUser(user.Id))
	request.SetProducts([]plaid.Products{plaid.PRODUCTS_AUTH, plaid.PRODUCTS_TRANSACTIONS})
	request.Transactions = plaid.NewLinkTokenTransactions()
	request.Transactions.SetDaysRequested(730)

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
func (s *NilPlaidService) GetPlaidItemTransactions(ctx context.Context, i models.PlaidItem) ([]models.Transaction, error) {
	return nilTransactions, nil
}
func (s *NilPlaidService) GetUserTransactions(ctx context.Context, userId string) ([]models.Transaction, error) {
	return nilTransactions, nil
}
func (s *NilPlaidService) SyncPlaidItemTransactions(ctx context.Context, item models.PlaidItem) error {
	return nil
}
func (s *NilPlaidService) SyncUserTransactions(ctx context.Context, userId string) error {
	return nil
}
func (s *NilPlaidService) GetLinkTokenForUser(ctx context.Context, user models.User) (string, error) {
	return "", nil
}
func (s *NilPlaidService) ExchangePublicToken(ctx context.Context, user models.User, publicToken string) (models.PlaidItem, error) {
	return models.PlaidItem{}, nil
}
