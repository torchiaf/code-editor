package api

import (
	"net/http"

	authentication "server/authentication"

	"github.com/gin-gonic/gin"
)

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := authentication.TokenValid(c)
		if err != nil {
			c.String(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}
		user, err := authentication.ExtractUser(c)
		if err != nil {
			c.String(http.StatusNotFound, "User not found")
			c.Abort()
			return
		}
		c.Set("user", user)
		c.Next()
	}
}
