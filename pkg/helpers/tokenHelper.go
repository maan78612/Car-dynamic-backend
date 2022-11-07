package helpers

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/maan78612/car_rental_dynamic/pkg/configs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SignedDetails struct {
	Email    string
	Name     string
	Uid      string
	UserType string
	jwt.StandardClaims
}

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")

var SECRETE_KEY = os.Getenv("SECRETE_KEY")

func GenerateAllTokens(email string, name string, userType string, uid string) (signedToken string, signedRefreshToken string, err error) {

	claims := &SignedDetails{
		Email:    email,
		Name:     name,
		Uid:      uid,
		UserType: userType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshClaims := &SignedDetails{

		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRETE_KEY))
	if err != nil {
		log.Panic(err)
		return

	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRETE_KEY))

	if err != nil {
		log.Panic(err)
		return

	}
	return token, refreshToken, err
}

func UpdateAllTokens(signedToken string, signedRefreshToken, userID string) {

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	var updateObj primitive.D
	updateObj = append(updateObj, bson.E{Key: "token", Value: signedToken})
	updateObj = append(updateObj, bson.E{Key: "refresh_token", Value: signedRefreshToken})
	upsert := true
	updatedAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{Key: "updated_at", Value: updatedAt})
	filter := bson.M{"user_id": userID}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}
	_, err := userCollection.UpdateOne(
		ctx,
		filter,
		// bson.D{{"$set", updateObj}},
		bson.D{primitive.E{Key: "$set", Value: updateObj}},
		&opt,
	)
	defer cancel()

	if err != nil {
		log.Panic(err)
		return
	}

}

func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {

	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(SECRETE_KEY), nil
		},
	)
	if err != nil {
		msg = err.Error()
		return
	}

	claims, ok := token.Claims.(*SignedDetails)

	if !ok {
		msg = "invalid token"
		msg = err.Error()
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = "token is expired"
		msg = err.Error()
		return
	}

	return claims, msg

}
