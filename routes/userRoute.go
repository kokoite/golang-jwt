package routes

import (
	"github.com/gin-gonic/gin"
	controller "github.com/pranjal/jwt/controllers"
	"github.com/pranjal/jwt/middleware"
)

func UserRoutes(router *gin.Engine) {

	router.Use(middleware.Authenticate())
	router.GET("/user", controller.FetchAllUsers())
	router.GET("/user/:id", controller.FetchUserById())
}
