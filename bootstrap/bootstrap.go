package bootstrap

import (
	"context"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/nynrathod/automator-api/api/router"
	cfg "github.com/nynrathod/automator-api/config"
	"github.com/nynrathod/automator-api/pkg/services"
	"github.com/redis/go-redis/v9"
	"sync"
	"time"
)

func NewApplication() *fiber.App {
	//env.SetupEnvFile()
	cfg.InitEnvConfigs()
	InitRedis()

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

	v1.Post("/set", func(c *fiber.Ctx) error {
		key := c.Query("key")
		value := c.Query("value")
		ctx := context.Background()
		err := redisClient.Set(ctx, key, value, 30*time.Second).Err()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to set value")
		}
		return c.SendString("Value set successfully")
	})

	v1.Get("/get", func(c *fiber.Ctx) error {
		key := c.Query("key")
		ctx := context.Background()
		value, err := redisClient.Get(ctx, key).Result()
		if err == redis.Nil {
			return c.Status(fiber.StatusNotFound).SendString("Key does not exist")
		} else if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to get value")
		}
		return c.SendString(value)
	})

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
	//app.Get("/dashboard", monitor.New())

	return app
}
