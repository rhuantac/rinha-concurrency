package model

import "time"

type TransactionType string

const (
	CreditTransaction TransactionType = "c"
	DebitTransaction  TransactionType = "d"
)

type Transaction struct {
	Value           int             `bson:"value"`
	TransactionType TransactionType `bson:"transaction_type"`
	Description     string          `bson:"string"`
	CreatedAt       time.Time       `bson:"created_at"`
	UserId          int             `bson:"user_id"`
}
