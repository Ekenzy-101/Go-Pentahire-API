package handlers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/Ekenzy-101/Pentahire-API/config"
	"github.com/Ekenzy-101/Pentahire-API/helpers"
	"github.com/Ekenzy-101/Pentahire-API/models"
	"github.com/Ekenzy-101/Pentahire-API/services"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func EmailVerification(c *gin.Context) {
	requestBody := &TokenField{}
	messages := helpers.ValidateRequestBody(c, requestBody)

	if messages != nil {
		c.JSON(http.StatusBadRequest, messages)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	redisClient := services.GetRedisClient()
	userId, err := redisClient.GetEx(ctx, config.RedisVerifyEmailPrefix+requestBody.Token, time.Millisecond).Result()
	if errors.Is(err, redis.Nil) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Token has expired or is not valid"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	options := models.SQLOptions{
		Arguments:         []interface{}{time.Now(), userId},
		AfterTableClauses: `SET email_verified_at = $1 WHERE id = $2`,
		ReturnColumns:     []string{"id"},
		Destination:       []interface{}{&userId},
	}
	response := models.UpdateAndReturnUserRow(ctx, options)
	if response != nil {
		c.JSON(response.StatusCode, response.Body)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}

func PhoneVerification(c *gin.Context) {

}

func LoginVerification(c *gin.Context) {

}
