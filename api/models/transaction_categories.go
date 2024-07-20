package models

type TransactionCategory struct {
	Id          string `json:"id"`
	Primary     string `json:"primary"`
	Detailed    string `json:"detailed"`
	Description string `json:"description"`
}
