package controllers

import (
	"github/spiffgreen/footballbet-escrow/dtos"
	"github/spiffgreen/footballbet-escrow/initializers"
	model "github/spiffgreen/footballbet-escrow/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetGames(c *gin.Context) {
	// Open the JSON file
	c.JSON(http.StatusOK, gin.H{
		"message": "Games",
		"data":    initializers.GameData,
	})
}

func GetBets(c *gin.Context) {
	user, _ := c.Get("user")
	userName := user.(model.User).Username
	var bets []model.Bet
	initializers.DB.Where("User1 = ? OR User2 = ?", userName, userName).Find(&bets)

	c.JSON(http.StatusOK, gin.H{
		"message": "Bets",
		"data":    bets,
	})
}

func GetOpenBets(c *gin.Context) {
	var bets []model.Bet
	initializers.DB.Where("Done = ?", false).Find(&bets)

	c.JSON(http.StatusOK, gin.H{
		"message": "Open bets",
		"data":    bets,
	})
}

func OpenBet(c *gin.Context) {
	var body dtos.OpenBet

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Incorrect body",
		})
		return
	}

	user, _ := c.Get("user")

	game := initializers.GameData[body.GameId]

	if game == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid game ID",
		})
	}

	if user.(model.User).Balance < body.StakeAmount {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Insufficient funds",
		})
		return
	}

	userData := user.(model.User)
	userData.Balance = userData.Balance - body.StakeAmount
	initializers.DB.Save(&userData)

	betSummary := game[body.WhoWins].(string) + " wins"

	bet := model.Bet{User1: user.(model.User).Username, User2: "", Done: false, BetAmount: body.StakeAmount, BetSummary: betSummary, ToWin: body.WhoWins}
	result := initializers.DB.Create(&bet)

	if result.Error != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Opened bet",
	})
}

func PlaceBet(c *gin.Context) {
	var body dtos.PlaceBet

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Incorrect body",
		})
		return
	}

	// Check balance
	var bet model.Bet
	initializers.DB.First(&bet, body.BetId)
	user, _ := c.Get("user")

	if user.(model.User).Username == bet.User1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "You can't bet with yourself",
		})
		return
	}

	if user.(model.User).Balance < bet.BetAmount {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Insufficient funds",
		})
		return
	}

	// Check win and update
	betData := initializers.GameData[body.BetId]
	if betData == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid bet ID",
		})
		return
	}

	if bet.ToWin == body.WhoWins {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid bet",
		})
		return
	}

	g1, _ := strconv.Atoi(betData[bet.ToWin+"g"].(string))
	g2, _ := strconv.Atoi(betData[body.WhoWins+"g"].(string))
	var user1, user2 model.User
	initializers.DB.Model(&bet).Where("ID = ?", bet.ID).Updates(map[string]interface{}{"Done": true, "User2": user.(model.User).Username})
	if g1 > g2 {
		// Credit user1 - user who opened the bet
		initializers.DB.Model(&user1).Where("UserName = ?", bet.User1).Update("Balance", gorm.Expr("Balance + ?", bet.BetAmount*2))
		// Debit user2 - user who opened the bet
		initializers.DB.Model(&user2).Where("UserName = ?", user.(model.User).Username).Update("Balance", user.(model.User).Balance-bet.BetAmount)
		c.JSON(http.StatusOK, gin.H{
			"message": "Bet lost",
			"data":    betData[bet.ToWin].(string) + " won",
		})
		return
	} else {
		// Credit user2 - user who closed the bet
		initializers.DB.Model(&user2).Where("UserName = ?", user.(model.User).Username).Update("Balance", user.(model.User).Balance+bet.BetAmount)
		c.JSON(http.StatusOK, gin.H{
			"message": "Bet won",
			"data":    betData[body.WhoWins].(string) + " won",
		})
		return
	}
}
