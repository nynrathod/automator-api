package bootstrap

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofiber/contrib/socketio"
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

type MessageObject struct {
	Data  string `json:"data"`
	From  string `json:"from"`
	Event string `json:"event"`
	To    string `json:"to"`
}

func NewApplication() *fiber.App {
	//env.SetupEnvFile()
	cfg.InitEnvConfigs()
	InitRedis()
	clients := make(map[string]string)
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

	//app.Use(func(c *fiber.Ctx) error {
	//	// IsWebSocketUpgrade returns true if the client
	//	// requested upgrade to the WebSocket protocol.
	//	if websocket.IsWebSocketUpgrade(c) {
	//		c.Locals("allowed", true)
	//		return c.Next()
	//	}
	//	return fiber.ErrUpgradeRequired
	//})

	// Multiple event handling supported
	socketio.On(socketio.EventConnect, func(ep *socketio.EventPayload) {
		fmt.Printf("Connection event 1 - User: %s", ep.Kws.GetStringAttribute("user_id"))
	})

	// Custom event handling supported
	socketio.On("CUSTOM_EVENT", func(ep *socketio.EventPayload) {
		fmt.Printf("Custom event - User: %s", ep.Kws.GetStringAttribute("user_id"))
		SendNoti()
		// --->

		// DO YOUR BUSINESS HERE

		// --->
	})

	// On message event
	socketio.On(socketio.EventMessage, func(ep *socketio.EventPayload) {

		fmt.Printf("Message event - User: %s - Message: %s", ep.Kws.GetStringAttribute("user_id"), string(ep.Data))

		message := MessageObject{}

		// Unmarshal the json message
		// {
		//  "from": "<user-id>",
		//  "to": "<recipient-user-id>",
		//  "event": "CUSTOM_EVENT",
		//  "data": "hello"
		//}
		err := json.Unmarshal(ep.Data, &message)
		if err != nil {
			fmt.Println(err)
			return
		}

		// Fire custom event based on some
		// business logic
		if message.Event != "" {
			ep.Kws.Fire(message.Event, []byte(message.Data))
		}

		// Emit the message directly to specified user
		err = ep.Kws.EmitTo(clients[message.To], ep.Data, socketio.TextMessage)
		if err != nil {
			fmt.Println(err)
		}
	})

	// On disconnect event
	socketio.On(socketio.EventDisconnect, func(ep *socketio.EventPayload) {
		// Remove the user from the local clients
		delete(clients, ep.Kws.GetStringAttribute("user_id"))
		fmt.Printf("Disconnection event - User: %s", ep.Kws.GetStringAttribute("user_id"))
	})

	// On close event
	// This event is called when the server disconnects the user actively with .Close() method
	socketio.On(socketio.EventClose, func(ep *socketio.EventPayload) {
		// Remove the user from the local clients
		delete(clients, ep.Kws.GetStringAttribute("user_id"))
		fmt.Printf("Close event - User: %s", ep.Kws.GetStringAttribute("user_id"))
	})

	// On error event
	socketio.On(socketio.EventError, func(ep *socketio.EventPayload) {
		fmt.Printf("Error event - User: %s", ep.Kws.GetStringAttribute("user_id"))
	})

	app.Get("/ws/:id", socketio.New(func(kws *socketio.Websocket) {

		// Retrieve the user id from endpoint
		userId := kws.Params("id")

		// Add the connection to the list of the connected clients
		// The UUID is generated randomly and is the key that allow
		// socketio to manage Emit/EmitTo/Broadcast
		clients[userId] = kws.UUID

		// Every websocket connection has an optional session key => value storage
		kws.SetAttribute("user_id", userId)

		//Broadcast to all the connected users the newcomer
		kws.Broadcast([]byte(fmt.Sprintf("New user connected: %s and UUID: %s", userId, kws.UUID)), true, socketio.TextMessage)
		//Write welcome message
		kws.Emit([]byte(fmt.Sprintf("Hello user: %s with UUID: %s", userId, kws.UUID)), socketio.TextMessage)
	}))

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
