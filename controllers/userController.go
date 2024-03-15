package controller

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/pranjal/jwt/database"
	helper "github.com/pranjal/jwt/helpers"
	"github.com/pranjal/jwt/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var mongoClient *mongo.Client
var userCollection *mongo.Collection
var validate = validator.New()

func init() {
	mongoClient = database.CreateMongoClient()
	userCollection = database.OpenCollection(mongoClient, "user")
}

func FetchAllUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		isAdmin := helper.AuthenticateAdmin(c)
		if !isAdmin {
			println("user is not allowed to see all users")
			c.JSON(http.StatusBadRequest, gin.H{"status": "400", "message": "User is not admin"})
			return
		}

		context, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}

		page, err := strconv.Atoi(c.Query("page"))
		if err != nil {
			page = 1
		}

		startIndex := (page - 1) * recordPerPage
		// startIndex, err = strconv.Atoi(c.Query("startIndex"))

		matching := bson.D{{Key: "$match", Value: bson.D{{}}}}
		grouping := bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "null"},
			{Key: "total_count", Value: bson.D{{Key: "$sum", Value: 1}}},
			{Key: "data", Value: bson.D{{Key: "$push", Value: "$$ROOT"}}},
		}}}
		project := bson.D{
			{Key: "$project", Value: bson.D{
				{Key: "_id", Value: 0},
				{Key: "total_count", Value: 1},
				{Key: "user_items", Value: bson.D{{
					Key: "$slice", Value: []interface{}{"$data", startIndex, recordPerPage},
				}}},
			}},
		}

		result, err := userCollection.Aggregate(context, mongo.Pipeline{matching, grouping, project})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		defer cancel()
		var allUsers []bson.M
		err = result.All(context, &allUsers)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusAccepted, allUsers[0])
	}
}

func FetchUserById() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("id")
		if userId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"status": "400", "message": "user id (id) is empty"})
			return
		}

		if err := helper.AuthenticateAccess(c, userId); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "400", "message": err.Error()})
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var user models.User
		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "400", "message": "something went wrong while decoding"})
			return
		}
		c.JSON(http.StatusAccepted, gin.H{"status": "200", "message": user})
	}
}
