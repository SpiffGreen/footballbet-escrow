package main

import (
	"github/spiffgreen/footballbet-escrow/initializers"
	model "github/spiffgreen/footballbet-escrow/models"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {
	// Migrate the schema
	initializers.DB.AutoMigrate(&model.User{}, &model.Bet{}, &model.Transaction{}, &model.BankData{})
}
