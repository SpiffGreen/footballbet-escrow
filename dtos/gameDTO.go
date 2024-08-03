package dtos

type OpenBet struct {
	GameId      uint16 // Id based on index of game
	WhoWins     string // expected value "a" or "h"
	StakeAmount uint64
}

type PlaceBet struct {
	BetId   uint16 // Id of already placed bet
	WhoWins string // expected value "a" or "h"
}
