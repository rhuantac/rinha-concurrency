package model

import "time"

type TransactionType string

const (
	CreditTransaction TransactionType = "c"
	DebitTransaction  TransactionType = "d"
)

type Transaction struct {
	Value           int             `bson:"value" json:"valor"`
	TransactionType TransactionType `bson:"transaction_type" json:"tipo"`
	Description     string          `bson:"string" json:"descricao"`
	CreatedAt       time.Time       `bson:"created_at" json:"realizada_em"`
	UserId          int             `bson:"user_id"`
}
