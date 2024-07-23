package models

import (
	"database/sql"
	"time"

	"github.com/plaid/plaid-go/v27/plaid"
)

type Transaction struct {
	Id                          string         `json:"id" db:"id"`
	PlaidItemId                 string         `json:"item_id" db:"item_id"`
	UserId                      string         `json:"user_id" db:"user_id"`
	TransactionCategoryDetailed sql.NullString `json:"transaction_category_detailed" db:"transaction_category_detailed"`
	AccountId                   string         `json:"account_id" db:"account_id"`
	TransactionId               string         `json:"transaction_id" db:"transaction_id"`
	Amount                      float64        `json:"amount" db:"amount"`
	Date                        time.Time      `json:"date" db:"date"`
	Data                        []byte         `json:"data" db:"data"`
}

func NewTransaction(plaidTransaction plaid.Transaction) Transaction {
	return NewTransactionForItem(plaidTransaction, PlaidItem{})
}

func NewTransactionForItem(plaidTransaction plaid.Transaction, item PlaidItem) Transaction {
	var transaction Transaction = Transaction{}
	transaction.PlaidItemId = item.Id
	transaction.UserId = item.UserId
	transaction.AccountId = plaidTransaction.AccountId
	transaction.TransactionId = plaidTransaction.TransactionId
	transaction.Amount = plaidTransaction.Amount
	transaction.Date, _ = time.Parse(time.DateOnly, plaidTransaction.Date)
	transaction.Data, _ = plaidTransaction.MarshalJSON()

	var transactionCategoryDetailed sql.NullString
	if tc, ok := plaidTransaction.GetPersonalFinanceCategoryOk(); ok {
		transactionCategoryDetailed.String = tc.Detailed
		transactionCategoryDetailed.Valid = true
	}
	transaction.TransactionCategoryDetailed = transactionCategoryDetailed

	return transaction
}

func NewTransactions(plaidTransactions []plaid.Transaction) []Transaction {
	return NewTransactionsForItem(plaidTransactions, PlaidItem{})
}

func NewTransactionsForItem(plaidTransactions []plaid.Transaction, item PlaidItem) []Transaction {
	transactions := make([]Transaction, 0, len(plaidTransactions))
	for _, pt := range plaidTransactions {
		transactions = append(transactions, NewTransactionForItem(pt, item))
	}
	return transactions
}
