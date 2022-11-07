package models

import "github.com/gofiber/fiber/v2"

type SuccessfulResponse struct {
	Status  int        `json:"status"`
	Message string     `json:"message"`
	Data    *fiber.Map `json:"data"`
}
type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    string `json:"error"`
}
