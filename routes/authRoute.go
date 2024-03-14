package routes

import (
	"github.com/gin-gonic/gin"
	controller "github.com/pranjal/jwt/controllers"
)

func AuthRoutes(router *gin.Engine) {
	router.POST("/login", controller.HandleLogin())
	router.POST("/register", controller.HandleRegister())
}
