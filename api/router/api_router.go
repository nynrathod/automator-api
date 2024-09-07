package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nynrathod/automator-api/api/handlers/access"
	"github.com/nynrathod/automator-api/api/handlers/users"
	"github.com/nynrathod/automator-api/middleware"
	acc "github.com/nynrathod/automator-api/pkg/access"
	"github.com/nynrathod/automator-api/pkg/sms"
	usr "github.com/nynrathod/automator-api/pkg/users"
	"github.com/nynrathod/automator-api/utilities"
	"go.mongodb.org/mongo-driver/mongo"
)

type AppServices struct {
	UserService   usr.Service
	SmsService    sms.Service
	AccessService acc.Service
	//AppsService    apps.Service
	UserCollection *mongo.Collection
	SmsCollection  *mongo.Collection

	AccessCollection *mongo.Collection
	//AppsCollection *mongo.Collection
}

type ApiRouter struct {
}
type RequestBody struct {
	Email string `json:"email"`
}

func InstallRouter(app fiber.Router, services AppServices) {
	//api := app.Group("/api", limiter.New())
	app.Post("/", func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Hello from api",
		})
	})

	app.Post("/auth/verifyEmail", users.VerifyEmail(services.UserService))
	app.Post("/auth/login", users.Login(services.UserService))
	app.Post("/auth/verifyOtp", middleware.Protected(), users.VerifyOtp(services.UserService))
	//app.Post("/auth/sendOtp", middleware.Protected(), users.SendOtp(services.UserService))
	app.Post("/auth/register", middleware.Protected(), users.Register(services.UserService))
	app.Post("/user/profile", middleware.Protected(), users.GetUser(services.UserService))

	app.Post("/user/listusers", middleware.Protected(), users.ListUser(services.UserService))

	app.Post("/user/toggleAccess", middleware.Protected(), access.ToggleAccess(services.AccessService))
	app.Post("/user/sharedUser", middleware.Protected(), access.SharedUser(services.AccessService))

	app.Post("/user/sharedAccess", access.SharedAccess(services.AccessService))

	// other test apis
	app.Post("/user/testname", func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Hello from testname api cicd2",
		})
	})
	app.Post("/user/test", users.TestApi(services.UserService))
	app.Post("/user/list", users.AddUser(services.UserService))

	app.Post("/auth/gentoken", func(ctx *fiber.Ctx) error {
		var body RequestBody
		if err := ctx.BodyParser(&body); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid request body",
			})
		}
		additionalClaims := jwt.MapClaims{
			"email": body.Email,
		}

		res, _ := utilities.GenerateJWT(additionalClaims)
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": res,
		})
	})

}
