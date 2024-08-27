package users

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/nynrathod/automator-api/api/presenter"
	"github.com/nynrathod/automator-api/pkg/entities"
	"github.com/nynrathod/automator-api/pkg/users"
	UTL "github.com/nynrathod/automator-api/utilities"
	"net/http"
)

func GetUser(service users.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {

		var requestBody entities.User
		err := c.BodyParser(&requestBody)

		if err != nil {

			c.Status(http.StatusBadRequest)
			return c.JSON(presenter.UserRegisterErrResponse(err))
		}

		authHeader := c.Get("Authorization")
		if authHeader == "" {
			print("anythingempty")
			return c.Status(http.StatusForbidden).JSON(presenter.GetProfileError(http.StatusForbidden))
		}

		_, tokenErr := UTL.VerifyToken(requestBody.Email, authHeader)
		if tokenErr != nil {
			fmt.Println("incorrec")
			return c.Status(http.StatusForbidden).JSON(presenter.GetProfileError(http.StatusForbidden))
		}

		//fmt.Println("helloyuser ", requestBody.Email)

		email := requestBody.Email
		// fmt.Println("requestBody", email)

		result, _ := service.GetUser(email)

		response := presenter.UserProfileResponse(result)

		return c.Status(fiber.StatusOK).JSON(response)
	}
}
