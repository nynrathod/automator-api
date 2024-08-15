package users

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
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

		fmt.Println("helloyuser ", requestBody.Email)

		email := requestBody.Email
		// fmt.Println("requestBody", email)

		result, _ := service.GetUser(email)

		response := presenter.UserProfileResponse(result)
		fmt.Println("res", response)

		additionalClaims := jwt.MapClaims{
			"exp":       14400,
			"id":        result.ID,
			"firstName": result.FirstName,
			"lastName":  result.LastName,
			"email":     result.Email,
		}

		jwtToken, _ := UTL.GenerateJWT(additionalClaims)
		(*response)["token"] = jwtToken

		statusCode, _ := (*response)["statusCode"].(int)
		delete(*response, "statusCode")

		return c.Status(statusCode).JSON(response)
	}
}
