package helper

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"github.com/pranjal/jwt/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SignedDetails struct {
	Email     string
	FirstName string
	LastName  string
	Uid       string
	UserType  string
	jwt.StandardClaims
}

var userCollection *mongo.Collection
var secretKey string

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		println("error occured while loading .env files")
	}
	secretKey = os.Getenv("SECRET_KEY")
	userCollection = database.OpenCollection(database.CreateMongoClient(), "users")
}

func GenerateAllToken(email, firstName, lastName, userType, uid string) (token, refreshToken string, err error) {
	claims := &SignedDetails{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Uid:       uid,
		UserType:  userType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}
	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}
	token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secretKey))

	if err != nil {
		log.Fatal("something went wrong")
		return "", "", err
	}

	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(secretKey))
	if err != nil {
		log.Fatal("something went wrong")
		return "", "", err
	}
	return token, "refreshToken", nil
}

func ValidateToken(tokenString string) (*SignedDetails, error) {
	token, err := jwt.ParseWithClaims(tokenString, &SignedDetails{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	},
	)

	if err != nil {
		log.Fatal("something went wrong")
	}

	if !token.Valid {
		log.Fatal("token is invalid")
		println("token is invalid")
	}
	claims, ok := token.Claims.(SignedDetails)

	if !ok {
		log.Fatal("unable to parse claims")
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		return nil, errors.New("token expired")
	}
	return &claims, nil
}

func UpdateAllTokens(signedToken, refreshToken, userId string) {
	context, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var updatedObject primitive.D
	updatedObject = append(updatedObject, bson.E{Key: "token", Value: signedToken})
	updatedObject = append(updatedObject, bson.E{Key: "refresh_token", Value: refreshToken})
	updatedAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updatedObject = append(updatedObject, bson.E{Key: "updated_at", Value: updatedAt})
	upsert := true
	filter := bson.M{"user_id": userId}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}
	_, err := userCollection.UpdateOne(context, filter, bson.D{{Key: "$set", Value: updatedObject}}, &opt)
	if err != nil {
		log.Panic(err)
		return
	}
	println("this is success")
}
