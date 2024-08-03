package main

import (
	"github/spiffgreen/footballbet-escrow/controllers"
	"github/spiffgreen/footballbet-escrow/initializers"
	"github/spiffgreen/footballbet-escrow/middlewares"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	initializers.LoadGames()
}

func main() {
	app := gin.Default()

	/// Routes
	// 	Auth/signup
	app.GET("/api/auth/profile", middlewares.RequireAuth, controllers.Profile)
	app.POST("/api/auth/signup", controllers.CreateAccount)
	app.POST("/api/auth/signin", controllers.LoginAccount)

	// Games/
	app.GET("/api/games/open-bets", controllers.GetOpenBets)
	app.GET("/api/games/bets", middlewares.RequireAuth, controllers.GetBets)
	app.GET("/api/games", controllers.GetGames)
	app.POST("/api/games/place-bet", middlewares.RequireAuth, controllers.PlaceBet)
	app.POST("/api/games/open-bet", middlewares.RequireAuth, controllers.OpenBet)

	// Payments
	app.GET("/api/payments/transaction-webhook", controllers.TransactionWebhook)
	app.GET("/api/payments/banks", controllers.GetBanks)
	app.POST("/api/payments/fund-wallet", middlewares.RequireAuth, controllers.InitializePayment)
	app.POST("/api/payments/set-bank", middlewares.RequireAuth, controllers.SetWithdrawBank)
	app.POST("/api/payments/withdraw-funds", middlewares.RequireAuth, controllers.WithdrawFunds)

	app.Run("127.0.0.1:3000")
}
