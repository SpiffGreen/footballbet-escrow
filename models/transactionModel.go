package model

import "gorm.io/gorm"

type Transaction struct {
	gorm.Model
	UserId             uint
	TransactionId      string
	TransactionSummary string
	Amount             uint64
	Action             string
}
