package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/Ekenzy-101/Pentahire-API/models"
	"github.com/Ekenzy-101/Pentahire-API/services"
	"github.com/gin-gonic/gin"
)

func CloseAccount(c *gin.Context) {

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

func Disable2FA(c *gin.Context) {

}

func Enable2FA(c *gin.Context) {

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
	response := models.SelectUserRow(ctx, options)
	if response != nil {
		c.JSON(response.StatusCode, response.Body)
		return
	}

	key, err := services.GenerateOTP(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	options.Arguments = []interface{}{key.Secret(), cliams.ID}
	options.AfterTableClauses = "SET otp_secret_key = $1 WHERE id = $2"
	response = models.UpdateAndReturnUserRow(ctx, options)
	if response != nil {
		c.JSON(response.StatusCode, response.Body)
		return
	}

	c.JSON(http.StatusOK, gin.H{"secret": key.Secret(), "url": key.URL()})
}

func UpdatePassword(c *gin.Context) {

}

func UpdateProfile(c *gin.Context) {

}
