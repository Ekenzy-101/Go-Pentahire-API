package handlers

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Ekenzy-101/Pentahire-API/config"
	"github.com/Ekenzy-101/Pentahire-API/helpers"
	"github.com/Ekenzy-101/Pentahire-API/models"
	"github.com/gin-gonic/gin"
)

type LoginRequestBody struct {
	Email    string `json:"email" binding:"email,max=255"`
	Password string `json:"password" binding:"required,min=8,max=128,password"`
}

func Login(c *gin.Context) {
	requestBody := &LoginRequestBody{}
	messages := helpers.ValidateRequestBody(c, requestBody)
	if messages != nil {
		c.JSON(http.StatusBadRequest, messages)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user := &models.User{}
	option := models.SQLOption{
		Arguments:         []interface{}{requestBody.Email},
		AfterTableClauses: `WHERE email = $1`,
		ReturnColumns: []string{
			"id",
			"average_rating",
			"created_at",
			"email",
			"firstname",
			"image",
			`CASE 
  			WHEN otp_secret_key = '' THEN CAST ('false' AS BOOLEAN)
  			ELSE CAST('true' AS BOOLEAN)
			END AS is_2fa_enabled`,
			`CASE 
  			WHEN email_verified_at IS NULL THEN CAST ('false' AS BOOLEAN)
  			ELSE CAST('true' AS BOOLEAN)
			END AS is_email_verified`,
			`CASE 
  			WHEN phone_verified_at IS NULL THEN CAST ('false' AS BOOLEAN)
  			ELSE CAST('true' AS BOOLEAN)
			END AS is_phone_verified`,
			"lastname",
			"password",
			"trips_count",
		},
		Destination: []interface{}{
			&user.ID,
			&user.AverageRating,
			&user.CreatedAt,
			&user.Email,
			&user.Firstname,
			&user.Image,
			&user.Is2FAEnabled,
			&user.IsEmailVerified,
			&user.IsPhoneVerified,
			&user.Lastname,
			&user.Password,
			&user.TripsCount,
		},
	}
	sqlResponse := models.FindUserRow(ctx, option)
	if sqlResponse != nil && sqlResponse.StatusCode == http.StatusNotFound {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid email or password"})
		return
	}

	if sqlResponse != nil {
		c.JSON(sqlResponse.StatusCode, sqlResponse.Body)
		return
	}

	matches, err := user.ComparePassword(requestBody.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	if !matches {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid email or password"})
		return
	}

	if user.Is2FAEnabled {
		log.Println(user.Is2FAEnabled)
		c.JSON(http.StatusNoContent, nil)
		return
	}

	accessToken, err := user.GenerateAccessToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	user.Password = ""
	c.SetCookie(config.AccessTokenCookieName, accessToken, config.AccessTokenTTLInSeconds, "", "", config.IsProduction, true)
	c.JSON(http.StatusOK, gin.H{"user": user})
}

func Logout(c *gin.Context) {

}

func Register(c *gin.Context) {
	user := &models.User{}
	messages := helpers.ValidateRequestBody(c, user)
	if messages != nil {
		c.JSON(http.StatusBadRequest, messages)
		return
	}

	user.NormalizeFields(true)
	err := user.HashPassword()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sqlOption := models.SQLOption{
		Arguments:     []interface{}{strings.ToLower(user.Email), user.Firstname, user.Lastname, user.Password},
		InsertColumns: []string{"email", "firstname", "lastname", "password"},
		ReturnColumns: []string{"id"},
		Destination:   []interface{}{&user.ID},
	}
	response := models.CreateUserRow(ctx, sqlOption)
	if response != nil {
		c.JSON(response.StatusCode, response.Body)
		return
	}

	token, err := helpers.GenerateRandomToken(24)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	response = models.CreateVerifyEmailRow(ctx, []interface{}{user.ID, token})
	if response != nil {
		c.JSON(response.StatusCode, response.Body)
		return
	}

	err = user.SendEmailVerificationMail(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	accessToken, err := user.GenerateAccessToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	user.Password = ""
	c.SetCookie(config.AccessTokenCookieName, accessToken, config.AccessTokenTTLInSeconds, "", "", config.IsProduction, true)
	c.JSON(http.StatusOK, gin.H{"user": user})
}

func ResetPassword(c *gin.Context) {

}
