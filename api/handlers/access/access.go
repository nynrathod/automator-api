package access

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nynrathod/automator-api/pkg/access"
	"github.com/nynrathod/automator-api/pkg/entities"
)

func ToggleAccess(service access.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {

		var requestBody *entities.Access

		if err := c.BodyParser(&requestBody); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error",
				"message": "Error on login request", "data": err.Error()})
		}

		_, handlerErr := service.ToggleAccess(requestBody)
		if handlerErr != nil {
			//fmt.Println("handlerErr", handlerErr)
			return nil
		}
		return nil
	}
}

func SharedUser(service access.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {

		var requestBody *entities.Access
		err := c.BodyParser(&requestBody)

		if err != nil {

			return err
		}

		res, _ := service.SharedUser(requestBody)
		return c.JSON(res)
	}
}

func SharedAccess(service access.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {

		var requestBody *entities.Access
		err := c.BodyParser(&requestBody)

		if err != nil {

			return err
		}

		res, _ := service.SharedAccess(requestBody)
		return c.JSON(res)
	}
}
