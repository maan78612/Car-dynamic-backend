package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/maan78612/car_rental_dynamic/pkg/helpers"
	"github.com/maan78612/car_rental_dynamic/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

var validate = validator.New()

func HashPassword(userPassword string) string {

	pass, err := bcrypt.GenerateFromPassword([]byte(userPassword), 14)

	if err != nil {
		log.Panic(err)
	}
	return string(pass)

}

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""
	if err != nil {
		msg = "Password is incorrect"
		check = false
	}
	return check, msg
}

func SignUp(c *fiber.Ctx) error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var newUser models.User
	defer cancel()

	//validate the request body
	if err := c.BodyParser(&newUser); err != nil {
		return c.Status(http.StatusBadRequest).JSON(models.ErrorResponse{Status: http.StatusBadRequest, Message: "error", Data: err.Error()})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&newUser); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(models.ErrorResponse{Status: http.StatusBadRequest, Message: "error", Data: validationErr.Error()})
	}

	count, err := userCollection.CountDocuments(ctx, bson.M{"email": newUser.Email})
	defer cancel()
	if err != nil {
		log.Panic(err)
		return c.Status(http.StatusBadRequest).JSON(models.ErrorResponse{Status: http.StatusBadRequest, Message: "error", Data: "error occured while checking for the email"})

	}

	if count > 0 {
		return c.Status(http.StatusBadRequest).JSON(models.ErrorResponse{Status: http.StatusBadRequest, Message: "error", Data: "This email  is already taken"})
	}

	password := HashPassword(*newUser.Password)
	newUser.Password = &password
	//Phone number

	countPhone, err := userCollection.CountDocuments(ctx, bson.M{"phone": newUser.Phone})
	defer cancel()
	if err != nil {
		log.Panic(err)
		return c.Status(http.StatusBadRequest).JSON(models.ErrorResponse{Status: http.StatusBadRequest, Message: "error", Data: "error occured while checking for the phone"})

	}

	if countPhone > 0 {
		return c.Status(http.StatusBadRequest).JSON(models.ErrorResponse{Status: http.StatusBadRequest, Message: "error", Data: "This  phone number is already taken"})

	}

	newUser.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	newUser.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	newUser.ID = primitive.NewObjectID()
	newUser.User_id = newUser.ID.Hex()
	token, refreshToken, err := helpers.GenerateAllTokens(*newUser.Email, *newUser.Name, *newUser.User_type, newUser.User_id)

	if err != nil {
		log.Panic(err)
		return c.Status(http.StatusBadRequest).JSON(models.ErrorResponse{Status: http.StatusBadRequest, Message: "error", Data: err.Error()})

	}

	newUser.Token = &token
	newUser.Referesh_token = &refreshToken
	result, err := userCollection.InsertOne(ctx, newUser)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(models.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
	}

	return c.Status(http.StatusCreated).JSON(models.SuccessfulResponse{Status: http.StatusCreated, Message: "success", Data: &fiber.Map{"data": result}})

}

func Login(c *fiber.Ctx) error {

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

	var user models.User
	var foundUser models.User
	defer cancel()
	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(models.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})

	}

	err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
	defer cancel()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(models.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: "email or password is incorrect"})

	}

	// check password is valid or nor
	isPasswordValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
	defer cancel()
	if !isPasswordValid {
		return c.Status(http.StatusInternalServerError).JSON(models.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: msg})

	}

	// check if user found or not

	if foundUser.Email == nil {
		return c.Status(http.StatusInternalServerError).JSON(models.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: "User not found"})

	}

	token, refreshToken, err := helpers.GenerateAllTokens(*foundUser.Email, *foundUser.Name, *foundUser.User_type, foundUser.User_id)

	if err != nil {
		log.Panic(err)
		return c.Status(http.StatusInternalServerError).JSON(models.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})

	}

	helpers.UpdateAllTokens(token, refreshToken, foundUser.User_id)

	err = userCollection.FindOne(ctx, bson.M{"user_id": foundUser.User_id}).Decode(&foundUser)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(models.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})

	}
	return c.Status(http.StatusOK).JSON(models.SuccessfulResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"user": foundUser}})

}
