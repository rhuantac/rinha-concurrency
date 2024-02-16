package model

type User struct {
	ID             int `bson:"_id"`
	Limit          int `bson:"limit"`
	InitialBalance int `bson:"initial_balance"`
	CurrentBalance int `bson:"current_balance"`
}
