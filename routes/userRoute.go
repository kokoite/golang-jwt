package routes

import (
	"github.com/gin-gonic/gin"
	controller "github.com/pranjal/jwt/controllers"
)

func UserRoutes(router *gin.Engine) {
	router.GET("/user", controller.FetchAllUsers())
	router.GET("/user/:id", controller.FetchUserById())
}
