package config

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var (
	IsDevelopment = gin.Mode() == gin.DebugMode
	IsProduction  = gin.Mode() == gin.ReleaseMode
	IsTesting     = gin.Mode() == gin.TestMode
)

var (
	AccessTokenSecret       string
	AppTokenSecret          string
	AWSBucket               string
	ClientOrigin            string
	DatabaseURL             string
	CaptchaSecretKey        string
	Port                    string
	RedisURL                string
	RefreshTokenSecret      string
	ResetPasswordTemplateID string
	SendgridAPIKey          string
	SendgridSender          string
	VerifyEmailTemplateID   string
)

func init() {
	filename := ""

	if IsTesting {
		filename = "../../.env.test"
	}

	if IsDevelopment {
		filename = ".env"
	}

	if filename != "" {
		if err := godotenv.Load(filename); err != nil {
			log.Fatal(err)
		}
	}

	AccessTokenSecret = os.Getenv("APP_ACCESS_SECRET")
	AppTokenSecret = os.Getenv("APP_TOKEN_SECRET")
	AWSBucket = os.Getenv("AWS_BUCKET")
	ClientOrigin = os.Getenv("CLIENT_ORIGIN")
	DatabaseURL = os.Getenv("DATABASE_URL")
	CaptchaSecretKey = os.Getenv("CAPTCHA_SECRET_KEY")
	Port = os.Getenv("PORT")
	RedisURL = os.Getenv("REDIS_URL")
	RefreshTokenSecret = os.Getenv("REFRESH_TOKEN_SECRET")
	ResetPasswordTemplateID = os.Getenv("RESET_PASSWORD_TEMPLATE_ID")
	SendgridAPIKey = os.Getenv("SENDGRID_API_KEY")
	SendgridSender = os.Getenv("SENDGRID_SENDER")
	VerifyEmailTemplateID = os.Getenv("VERIFY_EMAIL_TEMPLATE_ID")

	if Port == "" {
		Port = "5000"
	}
}
