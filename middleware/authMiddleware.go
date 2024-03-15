package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	helper "github.com/pranjal/jwt/helpers"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		println("before middleware called")
		clientToken := c.Request.Header.Get("token")
		if clientToken == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "token not present/invalid"})
			c.Abort()
			return
		}
		claims, err := helper.ValidateToken(clientToken)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			c.Abort()
			return
		}
		c.Set("email", claims.Email)
		c.Set("first_name", claims.FirstName)
		c.Set("last_name", claims.LastName)
		c.Set("uid", claims.Uid)
		c.Set("user_type", claims.UserType)
		c.Next()
		println("after middleware")
	}
}
