package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	route "github.com/pranjal/jwt/routes"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("something went wrong while loading .env files")
	}
}

func main() {
	router := gin.New()
	router.Use(gin.Logger())
	route.AuthRoutes(router)
	route.UserRoutes(router)
	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "all things are working"})
	})
	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
		println("unable to load port number from .env")
	}
	println("Server running on port", port)
	router.Run(":" + port)
}
