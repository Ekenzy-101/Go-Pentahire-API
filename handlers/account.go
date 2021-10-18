package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Ekenzy-101/Pentahire-API/config"
	"github.com/Ekenzy-101/Pentahire-API/helpers"
	"github.com/Ekenzy-101/Pentahire-API/models"
	"github.com/Ekenzy-101/Pentahire-API/services"
	"github.com/gin-gonic/gin"
)

func CloseAccount(c *gin.Context) {

}

func ConfirmOTPKey(c *gin.Context) {
	authUser := c.MustGet("user")
	cliams := authUser.(*services.AccessTokenClaims)
	requestBody := &CodeField{}
	if messages := helpers.ValidateRequestBody(c, requestBody); messages != nil {
		c.JSON(http.StatusBadRequest, messages)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var emailVerifiedAt interface{}
	secret := ""
	options := models.SQLOptions{
		Arguments:         []interface{}{cliams.ID},
		AfterTableClauses: "WHERE id = $1",
		ReturnColumns:     []string{"otp_secret_key", "email_verified_at"},
		Destination:       []interface{}{&secret, &emailVerifiedAt},
	}
	if response := models.SelectUserRow(ctx, options); response != nil {
		c.JSON(response.StatusCode, response.Body)
		return
	}

	if emailVerifiedAt == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Please verify your email address"})
		return
	}

	if !services.ValidateOTP(requestBody.Code, secret) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid verification code"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}

func DeleteOTPKey(c *gin.Context) {
	authUser := c.MustGet("user")
	cliams := authUser.(*services.AccessTokenClaims)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	options := models.SQLOptions{
		Arguments:         []interface{}{cliams.ID},
		AfterTableClauses: "SET otp_secret_key = '' WHERE id = $1",
		ReturnColumns:     []string{"id"},
		Destination:       []interface{}{&cliams.ID},
	}
	if response := models.UpdateAndReturnUserRow(ctx, options); response != nil {
		c.JSON(response.StatusCode, response.Body)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}

func GetOTPKey(c *gin.Context) {
	authUser := c.MustGet("user")
	cliams := authUser.(*services.AccessTokenClaims)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	email := ""
	options := models.SQLOptions{
		Arguments:         []interface{}{cliams.ID},
		AfterTableClauses: "WHERE id = $1",
		ReturnColumns:     []string{"email"},
		Destination:       []interface{}{&email},
	}
	if response := models.SelectUserRow(ctx, options); response != nil {
		c.JSON(response.StatusCode, response.Body)
		return
	}

	key, err := services.GenerateOTPKey(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	options.Arguments = []interface{}{key.Secret(), cliams.ID}
	options.AfterTableClauses = "SET otp_secret_key = $1 WHERE id = $2"
	if response := models.UpdateAndReturnUserRow(ctx, options); response != nil {
		c.JSON(response.StatusCode, response.Body)
		return
	}

	c.JSON(http.StatusOK, gin.H{"secret": key.Secret(), "url": key.URL()})
}

func UpdatePassword(c *gin.Context) {
	authUser := c.MustGet("user")
	cliams := authUser.(*services.AccessTokenClaims)
	requestBody := &UpdatePasswordRequestBody{}
	if messages := helpers.ValidateRequestBody(c, requestBody); messages != nil {
		c.JSON(http.StatusBadRequest, messages)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user := &models.User{}
	options := models.SQLOptions{
		Arguments:         []interface{}{cliams.ID},
		AfterTableClauses: "WHERE id = $1",
		ReturnColumns:     []string{"password"},
		Destination:       []interface{}{&user.Password},
	}
	if response := models.SelectUserRow(ctx, options); response != nil {
		c.JSON(response.StatusCode, response.Body)
		return
	}

	matches, err := user.ComparePassword(requestBody.OldPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	if !matches {
		c.JSON(http.StatusBadRequest, gin.H{"old_password": "Your current password was entered incorrectly. Please enter it again"})
		return
	}

	if requestBody.NewPassword == requestBody.OldPassword {
		c.JSON(http.StatusBadRequest, gin.H{"new_password": "Create a new password that isn't your current password"})
		return
	}

	user.Password = requestBody.NewPassword
	if err = user.HashPassword(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	options.Arguments = []interface{}{user.Password, cliams.ID}
	options.AfterTableClauses = "SET password = $1 WHERE id = $2"
	if response := models.UpdateAndReturnUserRow(ctx, options); response != nil {
		c.JSON(response.StatusCode, response.Body)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}

func UpdateProfile(c *gin.Context) {
	authUser := c.MustGet("user")
	cliams := authUser.(*services.AccessTokenClaims)
	requestBody := &UpdateProfileRequestBody{}
	if messages := helpers.ValidateRequestBody(c, requestBody); messages != nil {
		c.JSON(http.StatusBadRequest, messages)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user := &models.User{ID: cliams.ID}
	options := models.SQLOptions{
		Arguments:         []interface{}{user.ID},
		AfterTableClauses: "WHERE id = $1",
		ReturnColumns:     []string{"email"},
		Destination:       []interface{}{&user.Email},
	}
	response := models.SelectUserRow(ctx, options)
	if response != nil {
		c.JSON(response.StatusCode, response.Body)
		return
	}

	isSameEmail := strings.ToLower(requestBody.Email) == user.Email
	if isSameEmail {
		options.AfterTableClauses = "SET firstname = $1, lastname = $2 WHERE id = $3"
		options.Arguments = []interface{}{requestBody.Firstname, requestBody.Lastname, user.ID}
	} else {
		options.AfterTableClauses = "SET firstname = $1, lastname = $2, email = $3, email_verified_at = $4 WHERE id = $5"
		options.Arguments = []interface{}{requestBody.Firstname, requestBody.Lastname, requestBody.Email, nil, user.ID}
	}

	response = models.UpdateAndReturnUserRow(ctx, options)
	fmt.Println(response)
	if response != nil {
		c.JSON(response.StatusCode, response.Body)
		return
	}

	if !isSameEmail {
		token, err := helpers.GenerateRandomToken(24)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		redisClient := services.GetRedisClient()
		err = redisClient.Set(ctx, config.RedisVerifyEmailPrefix+token, user.ID, config.RedisVerifyEmailTTL).Err()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		err = user.SendEmailVerificationMail(token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}
