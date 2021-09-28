package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/Ekenzy-101/Pentahire-API/config"
	"github.com/Ekenzy-101/Pentahire-API/helpers"
	"github.com/Ekenzy-101/Pentahire-API/models"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {

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

	response := models.CreateUserRow(ctx, user)
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
