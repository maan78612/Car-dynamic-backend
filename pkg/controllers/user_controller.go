package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/maan78612/car_rental_dynamic/pkg/configs"
	"github.com/maan78612/car_rental_dynamic/pkg/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")

func GetAUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Params("user_id")
	println(userId)
	var user models.User
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(userId)
	// idQuery is used when we call FindOne function of mongoDB
	idQuery := bson.M{"user_id": objId}

	err := userCollection.FindOne(ctx, idQuery).Decode(&user)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(models.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
	}

	return c.Status(http.StatusOK).JSON(models.SuccessfulResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": user}})
}

func EditAUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Params("user_id")
	var user models.User
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(userId)

	//validate the request body
	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(models.ErrorResponse{Status: http.StatusBadRequest, Message: "error", Data: err.Error()})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&user); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(models.ErrorResponse{Status: http.StatusBadRequest, Message: "error", Data: validationErr.Error()})
	}
	//  query is used when we call FindOne function of mongoDB
	idQuery := bson.M{"id": objId}
	//  to update data in mongo we need $set in bson query
	//  we use bson format in query because mongo use snon format
	updateQuery :=
		bson.M{"$set": bson.M{"name": user.Name, "phone": user.Phone, "updated_at": user.Updated_at}}
	result, err := userCollection.UpdateOne(ctx, idQuery, updateQuery)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(models.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
	}

	//get updated user details
	var updatedUser models.User
	if result.MatchedCount == 1 {
		err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedUser)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(models.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
		}
	}

	return c.Status(http.StatusOK).JSON(models.SuccessfulResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": updatedUser}})
}

func DeleteAUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Params("user_id")
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(userId)

	result, err := userCollection.DeleteOne(ctx, bson.M{"id": objId})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(models.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
	}

	if result.DeletedCount < 1 {
		return c.Status(http.StatusNotFound).JSON(
			models.ErrorResponse{Status: http.StatusNotFound, Message: "error", Data: "User with specified ID not found!"},
		)
	}

	return c.Status(http.StatusOK).JSON(
		models.SuccessfulResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": "User successfully deleted!"}},
	)
}

func GetAllUsers(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var users []models.User
	defer cancel()

	results, err := userCollection.Find(ctx, bson.M{})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(models.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleUser models.User
		if err = results.Decode(&singleUser); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(models.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
		}

		users = append(users, singleUser)
	}

	return c.Status(http.StatusOK).JSON(
		models.SuccessfulResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": users}},
	)
}
