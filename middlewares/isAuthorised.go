package middlewares

import (
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func IsAuthorised() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		if tokenString == "" {
			// Redirect to login page if no token is present
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		// Assuming the token is prefixed with "Bearer "
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		// Validate the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrInvalidKey
			}
			// Return the key used to sign the token
			// Replace "your-secret-key" with your actual secret key
			return []byte("your-secret-key"), nil
		})

		if err != nil || !token.Valid {
			// Redirect to login page if token is invalid
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		// Proceed if everything is fine
		c.Next()
	}
}
