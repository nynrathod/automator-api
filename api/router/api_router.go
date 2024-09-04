package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nynrathod/automator-api/api/handlers/users"
	"github.com/nynrathod/automator-api/pkg/sms"
	usr "github.com/nynrathod/automator-api/pkg/users"
	"go.mongodb.org/mongo-driver/mongo"
)

type AppServices struct {
	UserService usr.Service
	SmsService  sms.Service
	//AppsService    apps.Service
	UserCollection *mongo.Collection
	SmsCollection  *mongo.Collection
	//AppsCollection *mongo.Collection
}

type ApiRouter struct {
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
	app.Post("/auth/verifyOtp", users.VerifyOtp(services.UserService))
	app.Post("/auth/sendOtp", users.SendOtp(services.UserService))
	app.Post("/auth/register", users.Register(services.UserService))
	app.Post("/user/profile", users.GetUser(services.UserService))

}

//
//func NewApiRouter() *ApiRouter {
//	return &ApiRouter{}
//}
