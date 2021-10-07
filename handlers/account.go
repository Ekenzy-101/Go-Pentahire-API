package handlers

import (
	"context"
	"net/http"
	"time"

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
	response := models.SelectUserRow(ctx, options)
	if response != nil {
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
	response := models.UpdateAndReturnUserRow(ctx, options)
	if response != nil {
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

}

func UpdateProfile(c *gin.Context) {

}
