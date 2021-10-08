package routes

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Ekenzy-101/Pentahire-API/config"
	"github.com/Ekenzy-101/Pentahire-API/handlers"
	"github.com/Ekenzy-101/Pentahire-API/helpers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func SetupRouter() *gin.Engine {
	binding.Validator = &helpers.DefaultValidator{}
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{config.ClientOrigin},
		AllowMethods:     []string{"GET", "POST", "DELETE", "PUT", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("%v operation is not supported for resource %v", c.Request.Method, c.Request.URL.Path)})
	})

	accountRouter := router.Group("/account").Use(Authorizer(true))
	accountRouter.DELETE("/otp-key", handlers.DeleteOTPKey)
	accountRouter.GET("/otp-key", handlers.GetOTPKey)
	accountRouter.POST("/otp-key/confirm", handlers.ConfirmOTPKey)

	authRouter := router.Group("/auth")
	authRouter.POST("/login", handlers.Login)
	authRouter.POST("/login/verify", handlers.VerifyLogin)
	authRouter.POST("/logout", handlers.Logout)
	authRouter.GET("/me", Authorizer(false), handlers.Me)
	authRouter.POST("/register", handlers.Register)
	authRouter.POST("/reset-password", handlers.ResetPassword)

	notification := router.Group("/notification")
	notification.POST("/verify-email", handlers.VerifyEmail)
	notification.POST("/forgot-password", handlers.ForgotPassword)

	verificationRouter := router.Group("/verification")
	verificationRouter.POST("/email", handlers.EmailVerification)

	return router
}
