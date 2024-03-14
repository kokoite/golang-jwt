package controller

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	helper "github.com/pranjal/jwt/helpers"
	"github.com/pranjal/jwt/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) string {

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Fatal("something went wrong when generating hash for password")
	}
	return string(bytes)
}

func verifyPassword(password string, foundPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(foundPassword), []byte(password))
	if err != nil {
		println("something went wrong when comparing hashed password")
		return false
	}
	return true
}

func HandleLogin() gin.HandlerFunc {

	return func(c *gin.Context) {
		context, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var user models.User
		var foundUser models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := userCollection.FindOne(context, bson.M{"email": user.Email}).Decode(&foundUser)
		if err != nil {
			println("user not found")
			c.JSON(http.StatusAccepted, gin.H{"status": "400", "message": "email not found"})
			return
		}
		isPasswordValid := verifyPassword(user.Password, foundUser.Password)

		if !isPasswordValid {
			println("invalid password entered")
			c.JSON(http.StatusBadRequest, gin.H{"status": "400", "message": "invalid password"})
			log.Fatal("invalid password")
		}

		if foundUser.Email == "" {
			println("email does not exist")
			c.JSON(http.StatusBadRequest, gin.H{"status": "400", "message": "email does not exist"})
			return
		}
		token, refreshToken, _ := helper.GenerateAllToken(foundUser.Email, foundUser.FirstName, foundUser.LastName, foundUser.UserType, foundUser.UserId)
		helper.UpdateAllTokens(token, refreshToken, foundUser.UserId)
		err = userCollection.FindOne(context, bson.M{"user_id": foundUser.UserId}).Decode(&foundUser)
		if err != nil {
			println("something went wrong in login")
			c.JSON(http.StatusBadRequest, gin.H{"status": "400", "message": "something went wrong in login"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "login success"})
	}
}

func HandleRegister() gin.HandlerFunc {
	return func(c *gin.Context) {
		context, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var user models.User
		if err := c.BindJSON(&user); err != nil {
			log.Fatal("error in handle register method", err)
			return
		}
		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": validationErr.Error()})
			return
		}
		hashedPassword := hashPassword(user.Password)
		user.Password = hashedPassword
		println("hashed password is", hashedPassword)
		count, err := userCollection.CountDocuments(context, bson.M{"email": user.Email, "phone": user.Phone})
		defer cancel()
		if err != nil || count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "something went wrong"})
			return
		}

		user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.UserId = user.ID.Hex()
		user.Token, user.RefreshToken, _ = helper.GenerateAllToken(user.Email, user.FirstName, user.LastName, user.UserType, user.UserId)
		insertNumber, insertError := userCollection.InsertOne(context, user)

		if insertError != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": insertError.Error()})
			return
		}
		c.JSON(http.StatusAccepted, gin.H{"message": "register success"})
		println("insert number in collection is", insertNumber)
	}
}
