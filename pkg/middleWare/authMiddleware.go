package middleWare

// import (
// 	"net/http"

// 	"github.com/gofiber/fiber"
// 	"github.com/maan78612/car_rental_dynamic/pkg/helpers"
// 	"github.com/maan78612/car_rental_dynamic/pkg/models"
// )

// func Authenticate(c *fiber.Ctx) error {

// 	clientToken := c.Fasthttp.Request.Header.Peek("token")
// 	clientTokenString := string(clientToken[:])
// 	if clientTokenString == "" {
// 		c.Context().Done()
// 		// c.Abort()

// 		c.Status(http.StatusInternalServerError).JSON(models.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: "No Authorization Header Provided"})
// 		c.Context().Done()
// 		return
// 	}

// 	claims, err := helpers.ValidateToken(clientTokenString)

// 	if err != "" {

// 		// c.Abort()
// 		c.Status(http.StatusInternalServerError).JSON(models.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: err})
// 		c.Context().Done()
// 		return
// 	}
// 	c.Set("email", claims.Email)
// 	c.Set("name", claims.Name)
// 	c.Set("uid", claims.Uid)
// 	c.Set("user_type", claims.UserType)
// 	c.Next()

// }
