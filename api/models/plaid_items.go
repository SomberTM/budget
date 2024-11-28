package models

import "time"

type PlaidItem struct {
	Id              string
	UserId          string
	ItemId          string
	InstitutionId   string
	InstitutionName string
	CreatedAt       time.Time
	ModifiedAt      time.Time
	AccessToken     string
	RawPlaidData    string
}

// func (i *PlaidItem) GetAccounts(ctx context.Context) ([]plaid.AccountBase, error) {
// 	accountsGetRequest := plaid.NewAccountsGetRequest(i.AccessToken)

// 	response, _, err := environment.Env.GetApiService().AccountsGet(ctx).AccountsGetRequest(*accountsGetRequest).Execute()
// 	if err != nil {
// 		log.Println("Error getting accounts", err.Error())
// 		return []plaid.AccountBase{}, err
// 	}

// 	return response.GetAccounts(), nil
// }

// func (i *PlaidItem) GetTransactions(ctx context.Context) ([]plaid.Transaction, error) {
// 	var transactions []plaid.Transaction
// 	var cursor string

// 	hasMore := true
// 	for hasMore {
// 		request := plaid.NewTransactionsSyncRequest(i.AccessToken)
// 		if cursor != "" {
// 			request.SetCursor(cursor)
// 		}

// 		resp, _, err := environment.Env.GetApiService().TransactionsSync(
// 			ctx,
// 		).TransactionsSyncRequest(*request).Execute()
// 		if err != nil {
// 			return []plaid.Transaction{}, err
// 		}

// 		// Add this page of results
// 		transactions = append(transactions, resp.GetAdded()...)
// 		// modified = append(modified, resp.GetModified()...)
// 		// removed = append(removed, resp.GetRemoved()...)

// 		hasMore = resp.GetHasMore()
// 		cursor = resp.GetNextCursor()
// 	}

// 	return transactions, nil
// }

// func (u *User) GetAccounts(ctx context.Context) ([]plaid.AccountBase, error) {
// 	rows, err := environment.Env.GetConnection().Query("SELECT * FROM plaid_items WHERE user_id = $1", u.Id)
// 	if err != nil {
// 		return []plaid.AccountBase{}, err
// 	}

// 	items := make([]PlaidItem, 0, 10)
// 	for rows.Next() {
// 		var item PlaidItem
// 		if err := rows.Scan(&item.Id, &item.UserId, &item.ItemId, &item.AccessToken); err != nil {
// 			return []plaid.AccountBase{}, err
// 		}

// 		items = append(items, item)
// 	}

// 	log.Println("Fetching accounts for items", items)

// 	accounts := make([]plaid.AccountBase, 0, 100)
// 	for i := 0; i < len(items); i++ {
// 		itemAccounts, err := items[i].GetAccounts(ctx)
// 		if err != nil {
// 			return []plaid.AccountBase{}, err
// 		}

// 		accounts = append(accounts, itemAccounts...)
// 	}

// 	return accounts, nil
// }

// func (u *User) GetTransactions(ctx context.Context) ([]plaid.Transaction, error) {
// 	rows, err := environment.Env.GetConnection().Query("SELECT * FROM plaid_items WHERE user_id = $1", u.Id)
// 	if err != nil {
// 		return []plaid.Transaction{}, err
// 	}

// 	items := make([]PlaidItem, 0, 10)
// 	for rows.Next() {
// 		var item PlaidItem
// 		if err := rows.Scan(&item.Id, &item.UserId, &item.ItemId, &item.AccessToken); err != nil {
// 			return []plaid.Transaction{}, err
// 		}

// 		items = append(items, item)
// 	}

// 	log.Println("Fetching transactions for items", items)

// 	transactions := make([]plaid.Transaction, 0, 100)
// 	for i := 0; i < len(items); i++ {
// 		itemTransactions, err := items[i].GetTransactions(ctx)
// 		if err != nil {
// 			return []plaid.Transaction{}, err
// 		}

// 		transactions = append(transactions, itemTransactions...)
// 	}

// 	return transactions, nil
// }
