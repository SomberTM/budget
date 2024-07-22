package models

type TransactionCursor struct {
	Id     string `json:"id" db:"id"`
	UserId string `json:"user_id" db:"user_id"`
	Cursor string `json:"cursor" db:"cursor"`
}
