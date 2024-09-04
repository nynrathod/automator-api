package users

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/nynrathod/automator-api/api/presenter"
	"github.com/nynrathod/automator-api/pkg/entities"
	"github.com/nynrathod/automator-api/pkg/users"
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

		//fmt.Println("helloyuser ", requestBody.Email)

		email := requestBody.Email
		// fmt.Println("requestBody", email)

		result, _ := service.GetUser(email)

		response := presenter.UserProfileResponse(result)

		return c.Status(fiber.StatusOK).JSON(response)
	}
}

func TestApi(service users.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).JSON("success")
	}
}

func AddUser(service users.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {

		var requestBody entities.User
		err := c.BodyParser(&requestBody)

		if err != nil {

			return err
		}

		fmt.Println("hanlderbody", requestBody)
		res, _ := service.AddUser(&requestBody)
		return c.JSON(res)
	}
}

func ListUser(service users.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {

		var requestBody entities.User
		err := c.BodyParser(&requestBody)

		if err != nil {

			return err
		}

		fmt.Println("hanlderbody", requestBody)
		res, _ := service.ListUser(requestBody.RequestEmail)
		return c.JSON(res)
	}
}
