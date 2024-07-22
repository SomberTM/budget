package models

import "github.com/plaid/plaid-go/v27/plaid"

type Transaction struct {
	plaid.Transaction
	Id                    string `json:"id" db:"id"`
	UserId                string `json:"user_id" db:"user_id"`
	TransactionCategoryId string `json:"transaction_category_id" db:"transaction_category_id"`
}
