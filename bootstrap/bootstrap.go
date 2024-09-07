package bootstrap

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/nynrathod/automator-api/api/router"
	cfg "github.com/nynrathod/automator-api/config"
	"github.com/nynrathod/automator-api/pkg/services"
	"sync"
)

func NewApplication() *fiber.App {

	cfg.InitEnvConfigs()

	db, cancel, err := cfg.SetupDatabase()
	if err != nil {
		log.Fatal("Database Connection Error: ", err)
	}

	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))
	app.Use(logger.New())

	SetupWebSocket(app, db)
	go ProcessQueue(db)

	api := app.Group("/api")
	v1 := api.Group("/v1")

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		serviceInitializer := services.NewAppServiceInitializer(db)
		appServices := serviceInitializer.InitializeAppServices(db)
		//router.SetupRouter(v1, appServices)

		router.InstallRouter(v1, appServices)
		wg.Done()
	}()

	wg.Wait()
	defer cancel()

	return app
}
