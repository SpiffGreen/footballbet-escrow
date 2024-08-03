package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github/spiffgreen/footballbet-escrow/dtos"
	"github/spiffgreen/footballbet-escrow/initializers"
	model "github/spiffgreen/footballbet-escrow/models"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitializePayment(c *gin.Context) {
	var body dtos.InitializePayment

	if c.BindJSON(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Incorrect body",
		})
		return
	}

	result, err := json.Marshal(map[string]string{
		"email":        body.Email,
		"amount":       body.Amount,
		"callback_url": "http://localhost:3000/api/payments/transaction-webhook",
	})
	if err != nil {
		log.Printf("Error creating request: %v", err)
	}

	var resp *http.Request
	resp, err = http.NewRequest(http.MethodPost, "https://api.paystack.co/transaction/initialize", bytes.NewBuffer(result))

	//Handle Error
	if err != nil {
		log.Printf("An Error Occured %v", err)
	}

	resp.Header.Add("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("PAYSTACK_SECRET")))
	var response *http.Response
	response, err = http.DefaultClient.Do(resp)

	//Handle Error
	if err != nil {
		log.Printf("An Error Occured %v", err)
	}
	defer response.Body.Close()

	//Read the response body
	body1, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var sb map[string]interface{}
	json.Unmarshal(body1, &sb)
	c.JSON(http.StatusOK, sb)
}

func TransactionWebhook(c *gin.Context) {
	// trxref=e8nw8ok1ky&reference=e8nw8ok1ky
	var query dtos.WebhookQuery
	if c.BindQuery(&query) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Please provide trxref and reference",
		})
		return
	}

	// Do processing
	resp, err := http.NewRequest(http.MethodGet, "https://api.paystack.co/transaction/verify/"+query.Trxref, nil)

	//Handle Error
	if err != nil {
		log.Printf("An Error Occured %v", err)
	}

	resp.Header.Add("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("PAYSTACK_SECRET")))
	var response *http.Response
	response, err = http.DefaultClient.Do(resp)

	//Handle Error
	if err != nil {
		log.Printf("An Error Occured %v", err)
	}
	defer response.Body.Close()

	//Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var sb map[string]interface{}
	json.Unmarshal(body, &sb)

	if sb["status"] == true {
		data, _ := sb["data"].(map[string]interface{})
		customer, _ := data["customer"].(map[string]interface{})
		amount := data["amount"].(float64)

		// Record transaction with unique trxref
		var count int64
		result := initializers.DB.Model(model.Transaction{}).Where("transaction_id = ?", query.Trxref).Count(&count)

		if result.Error != nil {
			log.Printf("Fetch transaction count: %v", result.Error)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Internal server error",
			})
			return // return early
		}

		if count != 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Already processed transaction",
			})
			return // return early
		}
		// Update wallet balance
		var user model.User
		result = initializers.DB.Where("Email = ?", customer["email"]).Find(&user)
		if result.RowsAffected == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid user id",
			})
			return // return early
		}
		trx := model.Transaction{
			UserId:             user.ID,
			TransactionId:      query.Trxref,
			TransactionSummary: "Deposited " + strconv.Itoa(int(amount)),
			Amount:             uint64(amount),
			Action:             "Deposit",
		}
		initializers.DB.Create(&trx)
		initializers.DB.Model(&user).Update("Balance", user.Balance+trx.Amount)

		c.JSON(http.StatusAccepted, gin.H{
			"message": "Processed",
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid transaction reference",
		})
	}
}

// Sets bank for users to withdraw to and
func SetWithdrawBank(c *gin.Context) {
	var body dtos.SetBankDetails

	if c.BindJSON(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Incorrect body",
		})
		return
	}

	result, err := json.Marshal(map[string]string{
		"type":           "nuban",
		"name":           body.Name,
		"account_number": body.AccountNumber,
		"bank_code":      body.BankCode,
		"currency":       "NGN",
	})
	if err != nil {
		log.Printf("Error creating request: %v", err)
	}

	var resp *http.Request
	resp, err = http.NewRequest(http.MethodPost, "https://api.paystack.co/transferrecipient", bytes.NewBuffer(result))

	//Handle Error
	if err != nil {
		log.Printf("An Error Occured %v", err)
	}

	resp.Header.Add("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("PAYSTACK_SECRET")))
	var response *http.Response
	response, err = http.DefaultClient.Do(resp)

	//Handle Error
	if err != nil {
		log.Printf("An Error Occured %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "You may have provided invalid account details",
		})
		return
	}
	defer response.Body.Close()

	//Read the response body
	body1, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "You may have provided invalid account details",
		})
		return
	}
	var sb map[string]interface{}
	json.Unmarshal(body1, &sb)
	data, _ := sb["data"].(map[string]interface{})

	user, _ := c.Get("user")
	var bank model.BankData
	getBankErr := initializers.DB.Where(&model.BankData{UserId: user.(model.User).ID}).First(&bank).Error

	bank.AccountNumber = body.AccountNumber
	bank.BankCode = body.BankCode
	bank.Name = body.Name
	bank.RecipientCode = data["recipient_code"].(string)
	if errors.Is(getBankErr, gorm.ErrRecordNotFound) {
		bank.UserId = user.(model.User).ID
	}
	initializers.DB.Save(&bank)

	c.JSON(http.StatusOK, gin.H{
		"message": "Bank account updated successfully",
	})
}

func WithdrawFunds(c *gin.Context) {
	var body dtos.WithdrawFunds

	if c.BindJSON(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Incorrect body",
		})
		return
	}

	user, _ := c.Get("user")

	var bankData model.BankData
	getBankErr := initializers.DB.Where(&model.BankData{UserId: user.(model.User).ID}).First(&bankData).Error
	if errors.Is(getBankErr, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Please set your bank for withdrawal first",
		})
		return
	}

	if user.(model.User).Balance < body.Amount {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Insufficient funds",
		})
		return
	}

	result, err := json.Marshal(map[string]interface{}{
		"source":    "balance",
		"amount":    body.Amount,
		"recipient": bankData.RecipientCode,
	})
	if err != nil {
		log.Printf("Error creating request: %v", err)
	}

	var resp *http.Request
	resp, err = http.NewRequest(http.MethodPost, "https://api.paystack.co/transfer", bytes.NewBuffer(result))

	//Handle Error
	if err != nil {
		log.Printf("An Error Occured %v", err)
	}

	resp.Header.Add("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("PAYSTACK_SECRET")))
	var response *http.Response
	response, err = http.DefaultClient.Do(resp)

	//Handle Error
	if err != nil {
		log.Printf("An Error Occured %v", err)
	}
	defer response.Body.Close()

	//Read the response body
	body1, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var sb map[string]interface{}
	json.Unmarshal(body1, &sb)
	c.JSON(http.StatusOK, sb)
}

func GetBanks(c *gin.Context) {
	// Do processing
	resp, err := http.NewRequest(http.MethodGet, "https://api.paystack.co/bank", nil)

	//Handle Error
	if err != nil {
		log.Printf("An Error Occured %v", err)
		return
	}

	resp.Header.Add("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("PAYSTACK_SECRET")))
	var response *http.Response
	response, err = http.DefaultClient.Do(resp)

	//Handle Error
	if err != nil {
		log.Printf("An Error Occured %v", err)
	}
	defer response.Body.Close()

	//Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
	}
	var sb map[string]interface{}
	json.Unmarshal(body, &sb)
	c.JSON(http.StatusOK, sb)
}
