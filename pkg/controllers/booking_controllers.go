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

var bookingsCollection *mongo.Collection = configs.GetCollection(configs.DB, "bookings")

func Createbooking(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var booking models.Bookings
	defer cancel()

	//validate the request body
	if err := c.BodyParser(&booking); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(models.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})

	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&booking); validationErr != nil {
		return c.Status(http.StatusInternalServerError).JSON(models.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: validationErr.Error()})

	}

	//  termporary user ID

	userID := "123"

	newBooking := models.Bookings{
		ID:               primitive.NewObjectID(),
		User_id:          &userID,
		Start_time:       booking.Start_time,
		Duration_in_days: booking.Duration_in_days,
		Created_at:       time.Now(),
	}

	newBooking.End_time = time.Date(newBooking.Start_time.Year(), newBooking.Start_time.Month(), newBooking.Start_time.Day()+booking.Duration_in_days, newBooking.Start_time.Hour(), newBooking.Start_time.Minute(), newBooking.Start_time.Second(), 0, time.UTC)

	_, err := bookingsCollection.InsertOne(ctx, newBooking)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(models.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})

	}

	// idQuery is used when we call FindOne function of mongoDB
	idQuery := bson.M{"_id": newBooking.ID}
	println(idQuery)
	err = bookingsCollection.FindOne(ctx, idQuery).Decode(&booking)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(models.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
	}

	return c.Status(http.StatusOK).JSON(models.SuccessfulResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": booking}})

}

func GetBooking(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	bookingID := c.Params("booking_id")
	println(bookingID)
	var booking models.Bookings
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(bookingID)
	// idQuery is used when we call FindOne function of mongoDB
	idQuery := bson.M{"_id": objId}
	println(idQuery)

	err := bookingsCollection.FindOne(ctx, idQuery).Decode(&booking)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(models.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
	}

	return c.Status(http.StatusOK).JSON(models.SuccessfulResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": booking}})
}

// func EditAUser(c *fiber.Ctx) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	userId := c.Params("user_id")
// 	var user models.User
// 	defer cancel()

// 	objId, _ := primitive.ObjectIDFromHex(userId)

// 	//validate the request body
// 	if err := c.BodyParser(&user); err != nil {
// 		return c.Status(http.StatusBadRequest).JSON(models.ErrorResponse{Status: http.StatusBadRequest, Message: "error", Data: err.Error()})
// 	}

// 	//use the validator library to validate required fields
// 	if validationErr := validate.Struct(&user); validationErr != nil {
// 		return c.Status(http.StatusBadRequest).JSON(models.ErrorResponse{Status: http.StatusBadRequest, Message: "error", Data: validationErr.Error()})
// 	}
// 	//  query is used when we call FindOne function of mongoDB
// 	idQuery := bson.M{"id": objId}
// 	//  to update data in mongo we need $set in bson query
// 	//  we use bson format in query because mongo use snon format
// 	updateQuery :=
// 		bson.M{"$set": bson.M{"name": user.Name, "phone": user.Phone, "updated_at": user.Updated_at}}
// 	result, err := userCollection.UpdateOne(ctx, idQuery, updateQuery)
// 	if err != nil {
// 		return c.Status(http.StatusInternalServerError).JSON(models.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
// 	}

// 	//get updated user details
// 	var updatedUser models.User
// 	if result.MatchedCount == 1 {
// 		err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedUser)
// 		if err != nil {
// 			return c.Status(http.StatusInternalServerError).JSON(models.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
// 		}
// 	}

// 	return c.Status(http.StatusOK).JSON(models.SuccessfulResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": updatedUser}})
// }

// func DeleteAUser(c *fiber.Ctx) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	userId := c.Params("user_id")
// 	defer cancel()

// 	objId, _ := primitive.ObjectIDFromHex(userId)

// 	result, err := userCollection.DeleteOne(ctx, bson.M{"id": objId})
// 	if err != nil {
// 		return c.Status(http.StatusInternalServerError).JSON(models.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
// 	}

// 	if result.DeletedCount < 1 {
// 		return c.Status(http.StatusNotFound).JSON(
// 			models.ErrorResponse{Status: http.StatusNotFound, Message: "error", Data: "User with specified ID not found!"},
// 		)
// 	}

// 	return c.Status(http.StatusOK).JSON(
// 		models.SuccessfulResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": "User successfully deleted!"}},
// 	)
// }

// func GetAllUsers(c *fiber.Ctx) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	var users []models.User
// 	defer cancel()

// 	results, err := userCollection.Find(ctx, bson.M{})

// 	if err != nil {
// 		return c.Status(http.StatusInternalServerError).JSON(models.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
// 	}

// 	//reading from the db in an optimal way
// 	defer results.Close(ctx)
// 	for results.Next(ctx) {
// 		var singleUser models.User
// 		if err = results.Decode(&singleUser); err != nil {
// 			return c.Status(http.StatusInternalServerError).JSON(models.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
// 		}

// 		users = append(users, singleUser)
// 	}

// 	return c.Status(http.StatusOK).JSON(
// 		models.SuccessfulResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": users}},
// 	)
// }
