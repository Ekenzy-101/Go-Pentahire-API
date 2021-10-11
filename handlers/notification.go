package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/Ekenzy-101/Pentahire-API/config"
	"github.com/Ekenzy-101/Pentahire-API/helpers"
	"github.com/Ekenzy-101/Pentahire-API/models"
	"github.com/Ekenzy-101/Pentahire-API/services"
	"github.com/gin-gonic/gin"
)

func ForgotPassword(c *gin.Context) {
	requestBody := &EmailField{}
	messages := helpers.ValidateRequestBody(c, requestBody)
	if messages != nil {
		c.JSON(http.StatusBadRequest, messages)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user := &models.User{Email: requestBody.Email}
	user.NormalizeFields(false)
	options := models.SQLOptions{
		Arguments:         []interface{}{user.Email},
		AfterTableClauses: "WHERE email = $1",
		Destination: []interface{}{
			&user.ID,
			&user.Firstname,
			&user.Lastname,
		},
		ReturnColumns: []string{"id", "firstname", "lastname"},
	}
	sqlResponse := models.SelectUserRow(ctx, options)
	if sqlResponse != nil && sqlResponse.StatusCode == http.StatusNotFound {
		c.JSON(http.StatusNotFound, gin.H{"email": "A user with the given email doesn't exist"})
		return
	}

	token, err := helpers.GenerateRandomToken(24)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	redisClient := services.GetRedisClient()
	err = redisClient.Set(ctx, config.RedisResetPasswordPrefix+token, user.ID, config.RedisResetPasswordTTL).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	err = user.SendPasswordResetMail(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Mail has been sent successfully"})
}

func VerifyEmail(c *gin.Context) {
	requestBody := &EmailField{}
	messages := helpers.ValidateRequestBody(c, requestBody)
	if messages != nil {
		c.JSON(http.StatusBadRequest, messages)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user := &models.User{Email: requestBody.Email}
	user.NormalizeFields(false)
	var emailVerifiedAt interface{}
	options := models.SQLOptions{
		Arguments:         []interface{}{user.Email},
		AfterTableClauses: "WHERE email = $1",
		Destination: []interface{}{
			&emailVerifiedAt,
			&user.Firstname,
			&user.Lastname,
		},
		ReturnColumns: []string{"email_verified_at", "firstname", "lastname"},
	}
	if response := models.SelectUserRow(ctx, options); response != nil {
		if response.StatusCode == http.StatusNotFound {
			c.JSON(http.StatusNotFound, gin.H{"email": "A user with the given email doesn't exist"})
			return
		}

		c.JSON(response.StatusCode, response.Body)
		return
	}

	if emailVerifiedAt != nil {
		c.JSON(http.StatusBadRequest, gin.H{"email": "Your email has been verified"})
		return
	}

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

	c.JSON(http.StatusOK, gin.H{"message": "Mail has been sent successfully"})
}

func VerifyPhone(c *gin.Context) {

}
