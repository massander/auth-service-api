package server

import (
	"github.com/gofiber/fiber/v2"
)

func okResponse(c *fiber.Ctx, data any) error {
	return c.Status(fiber.StatusOK).JSON(data)
}

func errorResponse(c *fiber.Ctx, status int, code string, message any) error {

	type Error struct {
		Status  int    `json:"status"`
		Code    string `json:"code"`
		Message string `json:"message,omitempty"`
	}

	type Response struct {
		Error Error `json:"error"`
	}

	errorResponse := Response{
		Error: Error{
			Status: status,
			Code:   code,
		},
	}

	if msg, isString := message.(string); isString {
		errorResponse.Error.Message = msg
	}

	return c.Status(status).JSON(errorResponse)

}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
