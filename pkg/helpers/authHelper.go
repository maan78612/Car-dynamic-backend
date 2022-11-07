package helpers

import (
	"errors"

	"github.com/gofiber/fiber"
)

func CheckUserType(c *fiber.Ctx, role string) (err error) {

	userType := c.Params("user_type")
	err = nil
	if userType != role {
		err = errors.New("uautherized to access this resource")
		return err
	}
	return err

}

func MatchUserTypeToUid(c *fiber.Ctx, user_id string) (err error) {
	userType := c.Params("user_type")
	uid := c.Params("uid")
	err = nil

	// check that if there is admin or not any other user then fetch data of user

	if userType == "USER" && uid != user_id {
		err := errors.New("unautherized to access this resource")
		return err
	}

	err = CheckUserType(c, userType)
	return err
}
