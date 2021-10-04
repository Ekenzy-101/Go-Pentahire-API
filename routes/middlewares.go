package routes

import (
	"net/http"

	"github.com/Ekenzy-101/Pentahire-API/config"
	"github.com/Ekenzy-101/Pentahire-API/services"
	"github.com/gin-gonic/gin"
)

func Authorizer(credentialsRequired bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken, err := c.Cookie(config.AccessTokenCookieName)
		if err != nil && credentialsRequired {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "No cookies found"})
			return
		}

		if err != nil {
			c.Next()
			return
		}

		options := services.JWTOptions{
			Secret: config.AccessTokenSecret,
			Token:  accessToken,
			Claims: &services.AccessTokenClaims{},
		}
		user, err := services.VerifyJWTToken(options)
		if err != nil && credentialsRequired {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Token is invalid or has expired"})
			return
		}

		if err != nil {
			c.Next()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}
