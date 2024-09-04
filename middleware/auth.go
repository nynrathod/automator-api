package middleware

import (
	"fmt"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nynrathod/automator-api/config"
	"github.com/nynrathod/automator-api/pkg/entities"
	"net/http"
)

// Protected protect routes
func Protected() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   jwtware.SigningKey{Key: []byte(config.EnvConfigs.JWTSecrete)},
		ErrorHandler: jwtError,
		SuccessHandler: func(c *fiber.Ctx) error {

			var requestBody entities.LoginEmailOtp
			err := c.BodyParser(&requestBody)
			if err != nil {
				return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "failed to parse request body"})
			}

			userToken := c.Locals("user").(*jwt.Token)
			fmt.Println("userToken", userToken.Raw)
			claims := userToken.Claims.(jwt.MapClaims)

			email, _ := claims["email"].(string)

			if email != requestBody.Email {
				return c.Status(fiber.StatusUnauthorized).
					JSON(fiber.Map{"status": "error", "message": "Email does not match identity", "data": nil})
			}

			return c.Next()
		},
	})
}

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"status": "error", "message": "Missing or malformed JWT", "data": nil})
	}
	return c.Status(fiber.StatusUnauthorized).
		JSON(fiber.Map{"status": "error", "message": "Invalid or expired JWT", "data": nil})
}
