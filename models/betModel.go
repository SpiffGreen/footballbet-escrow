package model

import "gorm.io/gorm"

type Bet struct {
	gorm.Model
	User1      string
	User2      string
	BetSummary string
	BetAmount  uint64
	Done       bool
	GameId     uint16
	ToWin      string
}
