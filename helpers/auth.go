package helper

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func AuthenticateAccess(c *gin.Context, userId string) error {
	userType := c.GetString("user_type")
	uid := c.GetString("uid")
	if userType == "ADMIN" {
		return nil
	} else {
		if uid != userType {
			return errors.New("unauthorized error to access this resource")
		}
		return nil
	}
}

func AuthenticateAdmin(c *gin.Context) bool {
	userType := c.GetString("user_type")
	return userType == "ADMIN"
}
