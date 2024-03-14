package middleware

import "github.com/gin-gonic/gin"

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		println("before middleware called")
		c.Next()
		println("after middleware")
	}
}
