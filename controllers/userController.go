package controllers

import (
	"encoding/json"
	"github/spiffgreen/footballbet-escrow/dtos"
	"github/spiffgreen/footballbet-escrow/initializers"
	model "github/spiffgreen/footballbet-escrow/models"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func CreateAccount(c *gin.Context) {
	var body dtos.CreateAccount

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Incorrect body",
		})
		return
	}

	// check for user
	if initializers.DB.First(&model.User{}, "username = ? OR email = ?", body.Username, body.Email).Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Username or email already in use",
		})
		return
	}

	// Create user
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to hash password",
		})
		return
	}

	user := model.User{Username: body.Username, Password: string(hashedPassword), Email: body.Email}
	result := initializers.DB.Create(&user)

	if result.Error != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Account created",
	})
}

func LoginAccount(c *gin.Context) {
	var body dtos.LoginAccount
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Incorrect body",
		})
		return
	}

	var user model.User
	result := initializers.DB.First(&user, "username = ?", body.Username)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid credentials",
		})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid credentials",
		})
		return
	}

	// Generate jwt
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":      user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24 * 30).Unix(),
	})
	var accessTokenString string
	accessTokenString, err = accessToken.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		log.Fatal("Failed to create jwt", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid credentials",
		})
		return
	}

	c.JSON(201, gin.H{
		"message":     "Successful signin",
		"accessToken": accessTokenString,
	})
}

func Profile(c *gin.Context) {
	user, _ := c.Get("user")
	result, _ := json.Marshal(user)

	var tempMap map[string]interface{}
	json.Unmarshal(result, &tempMap)
	delete(tempMap, "Password")

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile",
		"body":    tempMap,
	})
}
