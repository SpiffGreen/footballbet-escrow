package model

import "gorm.io/gorm"

type BankData struct {
	gorm.Model
	UserId        uint
	Name          string
	AccountNumber string
	BankCode      string
	RecipientCode string
}
