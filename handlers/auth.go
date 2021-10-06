package handlers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Ekenzy-101/Pentahire-API/config"
	"github.com/Ekenzy-101/Pentahire-API/helpers"
	"github.com/Ekenzy-101/Pentahire-API/models"
	"github.com/Ekenzy-101/Pentahire-API/services"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

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
	option := models.SQLOptions{
		Arguments:         []interface{}{requestBody.Email},
		AfterTableClauses: `WHERE email = $1`,
		ReturnColumns:     helpers.GenerateUserReturnColumns([]string{}),
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
			&user.PhoneNo,
			&user.TripsCount,
		},
	}
	sqlResponse := models.SelectUserRow(ctx, option)
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
	c.SetCookie(config.AccessTokenCookieName, "", -1, "/", "", config.IsProduction, true)
	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}

func Me(c *gin.Context) {
	value, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusOK, gin.H{"user": nil})
		return
	}

	payload, ok := value.(*services.AccessTokenClaims)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"user": nil})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user := &models.User{ID: payload.ID}
	options := models.SQLOptions{
		AfterTableClauses: "WHERE id = $1",
		Arguments:         []interface{}{payload.ID},
		ReturnColumns:     helpers.GenerateUserReturnColumns([]string{"id", "password"}),
		Destination: []interface{}{
			&user.AverageRating,
			&user.CreatedAt,
			&user.Email,
			&user.Firstname,
			&user.Image,
			&user.Is2FAEnabled,
			&user.IsEmailVerified,
			&user.IsPhoneVerified,
			&user.Lastname,
			&user.PhoneNo,
			&user.TripsCount,
		},
	}
	sqlResponse := models.SelectUserRow(ctx, options)
	if sqlResponse != nil {
		c.JSON(http.StatusOK, gin.H{"user": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func Register(c *gin.Context) {
	requestBody := &RegisterRequestBody{}
	messages := helpers.ValidateRequestBody(c, requestBody)
	if messages != nil {
		c.JSON(http.StatusBadRequest, messages)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	hCaptchaResponseBody, err := services.VerifyHCaptchaToken(ctx, requestBody.Token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	log.Println("HCaptchaErrorCodes", hCaptchaResponseBody.ErrorCodes)
	if !hCaptchaResponseBody.Success {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Please provide a valid hcaptcha token"})
		return
	}

	user := &models.User{
		Email:     requestBody.Email,
		Firstname: requestBody.Firstname,
		Lastname:  requestBody.Lastname,
		Password:  requestBody.Password,
	}
	user.NormalizeFields(true)
	err = user.HashPassword()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	sqlOption := models.SQLOptions{
		Arguments:     []interface{}{strings.ToLower(user.Email), user.Firstname, user.Lastname, user.Password},
		InsertColumns: []string{"email", "firstname", "lastname", "password"},
		ReturnColumns: []string{"id"},
		Destination:   []interface{}{&user.ID},
	}
	response := models.InsertUserRow(ctx, sqlOption)
	if response != nil {
		c.JSON(response.StatusCode, response.Body)
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
	requestBody := &ResetPasswordRequestBody{}
	messages := helpers.ValidateRequestBody(c, requestBody)
	if messages != nil {
		c.JSON(http.StatusBadRequest, messages)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	redisClient := services.GetRedisClient()
	userId, err := redisClient.GetEx(ctx, config.RedisResetPasswordPrefix+requestBody.Token, time.Millisecond).Result()
	if errors.Is(err, redis.Nil) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Token has expired or is not valid"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	user := &models.User{ID: userId, Password: requestBody.Password}
	err = user.HashPassword()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	options := models.SQLOptions{
		Arguments:         []interface{}{user.Password, user.ID},
		AfterTableClauses: `SET password = $1 WHERE id = $2`,
		ReturnColumns:     helpers.GenerateUserReturnColumns([]string{"id"}),
		Destination: []interface{}{
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
			&user.PhoneNo,
			&user.TripsCount,
		},
	}
	response := models.UpdateAndReturnUserRow(ctx, options)
	if response != nil {
		c.JSON(response.StatusCode, response.Body)
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
