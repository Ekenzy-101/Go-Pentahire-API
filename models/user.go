package models

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Ekenzy-101/Pentahire-API/config"
	"github.com/Ekenzy-101/Pentahire-API/services"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"golang.org/x/crypto/argon2"
)

type passwordConfig struct {
	time      uint32
	memory    uint32
	threads   uint8
	keyLength uint32
}

type User struct {
	ID            string    `json:"id"`
	AverageRating float64   `json:"average_rating"`
	CreatedAt     time.Time `json:"created_at"`
	Email         string    `json:"email" binding:"email,max=255"`
	// Favourites      []interface{} `json:"favourites"`
	Firstname       string `json:"firstname" binding:"required,name,max=50"`
	Image           string `json:"image"`
	Is2FAEnabled    bool   `json:"is_2fa_enabled"`
	IsEmailVerified bool   `json:"is_email_verified"`
	IsPhoneVerified bool   `json:"is_phone_verified"`
	Lastname        string `json:"lastname" binding:"required,name,max=50"`
	OTPSecretKey    string `json:"otp_secret_key,omitempty"`
	Password        string `json:"password,omitempty" binding:"required,min=8,max=128,password"`
	TripsCount      int    `json:"trips_count"`
}

func (user *User) NormalizeFields(new bool) {
	user.Email = strings.ToLower(user.Email)
	if new {
		user.CreatedAt = time.Now()
		user.AverageRating = 0.0
	}
}

func (user *User) ComparePassword(password string) (bool, error) {
	parts := strings.Split(user.Password, "$")
	if len(parts) < 4 {
		return false, errors.New("invalid string")
	}

	c := &passwordConfig{}
	_, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &c.memory, &c.time, &c.threads)
	if err != nil {
		return false, err
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}

	decodedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}

	c.keyLength = uint32(len(decodedHash))
	comparisonHash := argon2.IDKey([]byte(password), salt, c.time, c.memory, c.threads, c.keyLength)
	return (subtle.ConstantTimeCompare(decodedHash, comparisonHash) == 1), nil
}

func (user *User) GenerateAccessToken() (string, error) {
	claims := &services.AccessTokenClaim{
		Email: user.Email,
		ID:    user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Second * config.AccessTokenTTLInSeconds)),
			Issuer:    config.ClientOrigin,
		},
	}

	option := services.JWTOption{
		SigningMethod: jwt.SigningMethodHS256,
		Claims:        claims,
		Secret:        config.AccessTokenSecret,
	}
	return services.SignToken(option)
}

func (user *User) HashPassword() error {
	c := &passwordConfig{
		time:      1,
		memory:    64 * 1024,
		threads:   4,
		keyLength: 32,
	}
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return err
	}

	hash := argon2.IDKey([]byte(user.Password), salt, c.time, c.memory, c.threads, c.keyLength)
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)
	format := "$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s"
	user.Password = fmt.Sprintf(format, argon2.Version, c.memory, c.time, c.threads, b64Salt, b64Hash)
	return nil
}

func (user *User) SendEmailVerificationMail(token string) error {
	name := fmt.Sprintf("%v %v", user.Firstname, user.Lastname)
	to := mail.NewEmail(name, user.Email)
	link := fmt.Sprintf("%v/verify-email/%v", config.ClientOrigin, token)
	data := gin.H{"link": link}

	option := services.MailOption{To: to, Data: data, TemplateID: config.VerifyEmailTemplateID}
	response, err := services.SendMail(option)
	log.Printf("SendEmailVerificationMail StatusCode %+v\n", response.StatusCode)
	return err
}

func (user *User) SendPasswordResetMail(token string) error {
	name := fmt.Sprintf("%v %v", user.Firstname, user.Lastname)
	to := mail.NewEmail(name, user.Email)
	link := fmt.Sprintf("%v/reset-password/", config.ClientOrigin)
	data := gin.H{"email": user.Email, "link": link, "token": token}

	option := services.MailOption{To: to, Data: data, TemplateID: config.ResetPasswordTemplateID}
	response, err := services.SendMail(option)
	log.Printf("SendPasswordResetMail StatusCode %+v\n", response.StatusCode)
	return err
}
